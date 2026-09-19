package chargen

import (
	"strings"

	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/dice"
)

// rollCheck resolves a target from a table result.
//
// A characteristic check is 2d6 plus that characteristic's modifier, as
// p. 110 states for enlistment and the book uses throughout.
//
// A skill check -- "Roll Gambler 8+", "Roll Melee (Any) 8+" -- is stated
// nowhere in this book, which points at the Core Rulebook's task system
// instead (p. 13). ERRATA E-7 records the reading: 2d6 plus the
// character's level in the named skill, with an unheld skill contributing
// nothing.
func (g *Generator) rollCheck(target career.Target) dice.Throw {
	if target.Characteristic != "" {
		which, ok := characteristicByName(target.Characteristic)
		if !ok {
			return g.dice.Throw(target.Number)
		}

		modifier := g.char.State.Characteristics.Modifier(which)

		return g.dice.Throw(target.Number, dice.Mod{Name: target.Characteristic, Value: modifier})
	}

	g.char.Provenance.Deviate("E-7")

	level := max(g.char.State.SkillLevel(target.Skill), 0)

	return g.dice.Throw(target.Number, dice.Mod{Name: target.Skill, Value: level})
}

// applyInjury rolls the Injury table (p. 119) as many times as the result
// asked for.
func (g *Generator) applyInjury(effect career.Effect, cause int) error {
	times := max(effect.Times, 1)

	table := career.Injury()

	for range times {
		roll := g.dice.D6()
		throw := g.log.Roll(roll, "p. 119")
		row := table[roll.Total-1]

		// ERRATA E-4: recovery is charged against credits held at the time
		// of the injury, so a character who has not mustered out cannot
		// afford it and the loss is permanent from this moment.
		permanent := g.char.State.Credits <= 0
		if permanent {
			g.char.Provenance.Deviate("E-4")
		}

		g.char.State.Injuries = append(g.char.State.Injuries, Injury{
			Detail:    row.Summary,
			Permanent: permanent,
			Term:      len(g.char.State.Terms),
		})

		g.log.Consequence(ConsequenceEvent{
			Kind:      ConsequenceInjury,
			Cause:     throw,
			Detail:    row.Summary,
			Permanent: permanent,
			Cite:      "p. 119",
		})

		err := g.applyAll(row.Effects, throw)
		if err != nil {
			return err
		}
	}

	_ = cause

	return nil
}

// rollLifeEvent is the shared table of p. 120, which every career's d66
// reaches at 31-36.
func (g *Generator) rollLifeEvent(cause int) error {
	roll := g.dice.TwoD6()
	throw := g.log.Roll(roll, "p. 120")
	row := career.LifeEvents()[roll.Total-2]

	g.consequence(ConsequenceCareer, throw, "life event: "+row.Summary, "")

	_ = cause

	return g.applyAll(row.Effects, throw)
}

// rollMishap rolls the current career's mishap table. eject is false when
// an event sent the character here rather than a failed survival throw, and
// false again for the two careers that say in print that a mishap does not
// force a character out (pp. 263, 298).
func (g *Generator) rollMishap(cause int, eject bool) error {
	if g.career == nil {
		return ErrNoCareer
	}

	roll := g.dice.TwoD6()
	throw := g.log.Roll(roll, "p. 113")
	row := g.career.Mishaps[roll.Total-2]

	g.consequence(ConsequenceCareer, throw, "mishap: "+row.Summary, g.career.Name)

	if eject && g.career.MishapEjects {
		g.ejected = true
	}

	_ = cause

	return g.applyAll(row.Effects, throw)
}

// grantBenefits pushes a batch of mustering-out rolls, or, where the effect
// carries only a modifier, applies it to every batch of this career
// (ERRATA E-3).
func (g *Generator) grantBenefits(effect career.Effect, cause int) {
	name := ""
	if g.career != nil {
		name = g.career.Name
	}

	if effect.Count == 0 && effect.Modifier != 0 && effect.Scope == career.ScopeCareer {
		for i := range g.char.State.Benefits {
			if g.char.State.Benefits[i].Career == name {
				g.char.State.Benefits[i].Modifier += effect.Modifier
			}
		}

		g.careerBenefitMod += effect.Modifier
		g.char.Provenance.Deviate("E-3")
		g.consequence(ConsequenceBenefitRolls, cause, effect.Detail, name)

		return
	}

	if effect.Modifier != 0 {
		g.char.Provenance.Deviate("E-3")
	}

	g.char.State.Benefits = append(g.char.State.Benefits, BenefitBatch{
		Career:   name,
		Rolls:    effect.Count,
		Modifier: effect.Modifier,
	})

	g.consequence(ConsequenceBenefitRolls, cause, effect.Detail, name)
}

// gainTies adds Allies, Contacts, Rivals or Enemies. Where the book rolls
// for how many, Dice carries the expression.
func (g *Generator) gainTies(effect career.Effect, cause int) error {
	count := effect.Count

	if effect.Dice != "" {
		rolled, err := g.rollExpression(effect.Dice, "p. 120")
		if err != nil {
			return err
		}

		count = rolled
	}

	name := ""
	if g.career != nil {
		name = g.career.Name
	}

	for range count {
		g.char.State.Ties = append(g.char.State.Ties, Tie{
			Kind:   string(effect.Relationship),
			Origin: name,
			Rating: 0,
		})
	}

	g.consequence(ConsequenceRelationship, cause, effect.Detail, name)

	return nil
}

func (g *Generator) payCredits(effect career.Effect, cause int) error {
	amount, err := g.rollExpression(effect.Dice, g.cite)
	if err != nil {
		return err
	}

	g.char.State.Credits += amount
	g.consequence(ConsequenceCredits, cause, effect.Detail, "")

	return nil
}

func (g *Generator) changeStash(effect career.Effect, cause int) {
	if effect.Item == "" {
		g.char.State.Stash = nil
		g.consequence(ConsequenceStash, cause, "lose the contents of the stash", "")

		return
	}

	g.char.State.Stash = append(g.char.State.Stash, effect.Item)
	g.consequence(ConsequenceStash, cause, effect.Detail, "")
}

// rollExpression evaluates the small language the tables are written in:
// a bare number, "NdM", or "NdMxK", where M is 6 or 3 -- the only two dice
// the book throws for a quantity. Anything else is a transcription error
// rather than a rule, so it errors rather than guessing what the page
// meant.
func (g *Generator) rollExpression(expr, cite string) (int, error) {
	multiplier := 1

	if index := strings.Index(expr, "x"); index >= 0 {
		parsed, ok := atoi(expr[index+1:])
		if !ok {
			return 0, ErrBadExpression
		}

		multiplier = parsed
		expr = expr[:index]
	}

	index := strings.Index(expr, "d")
	if index < 0 {
		value, ok := atoi(expr)
		if !ok {
			return 0, ErrBadExpression
		}

		return value * multiplier, nil
	}

	count, ok := atoi(expr[:index])
	if !ok {
		return 0, ErrBadExpression
	}

	sides, ok := atoi(expr[index+1:])
	if !ok {
		return 0, ErrBadExpression
	}

	switch sides {
	case sixSided:
		roll := g.dice.ND6(count)
		g.log.Roll(roll, cite)

		return roll.Total * multiplier, nil
	case threeSided:
		total := 0

		for range count {
			roll := g.dice.D3()
			g.log.Roll(roll, cite)

			total += roll.Total
		}

		return total * multiplier, nil
	}

	return 0, ErrBadExpression
}

// The two dice the book throws for a quantity, and the base its numbers
// are written in.
const (
	sixSided   = 6
	threeSided = 3
	decimal    = 10
)

func atoi(s string) (int, bool) {
	if s == "" {
		return 0, false
	}

	value := 0

	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, false
		}

		value = value*decimal + int(r-'0')
	}

	return value, true
}
