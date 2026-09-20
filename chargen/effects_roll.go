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

	mods := append(
		[]dice.Mod{{Name: target.Skill, Value: level}},
		g.takeSkillModifiers(target.Skill)...,
	)

	return g.dice.Throw(target.Number, mods...)
}

// takeSkillModifiers is the two results that modify a skill check for a
// spell rather than once: "+1 DM to Melee checks in this career". The
// modifier names the skill it applies to, and is not spent by a check of
// any other skill.
func (g *Generator) takeSkillModifiers(skill string) []dice.Mod {
	var (
		mods []dice.Mod
		kept []PendingModifier
	)

	for _, pending := range g.pending {
		if pending.Applies != skillCheckThrow || pending.OnSkill != skill {
			kept = append(kept, pending)

			continue
		}

		mods = append(mods, dice.Mod{Name: pending.Detail, Value: pending.Value})

		// A skill modifier lasts for the career that granted it, which
		// dropCareerModifiers ends; it is not spent by a check.
		kept = append(kept, pending)
	}

	g.pending = kept

	return mods
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

// rollMilitaryEvent is the shared table of p. 121, which the military
// careers' d66 tables reach at 41-46.
func (g *Generator) rollMilitaryEvent(cause int) error {
	roll := g.dice.TwoD6()
	throw := g.log.Roll(roll, "p. 121")
	row := career.MilitaryEvents()[roll.Total-2]

	g.consequence(ConsequenceCareer, throw, "military event: "+row.Summary, "")

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
func (g *Generator) grantBenefits(effect career.Effect, cause int) error {
	name := ""
	if g.career != nil {
		name = g.career.Name
	}

	switch {
	case effect.ForfeitAll:
		g.forfeitBenefits(effect, name, cause)

		return nil
	case effect.Immediate:
		return g.immediateCashRolls(effect, cause)
	}

	count := effect.Count

	if effect.Dice != "" {
		rolled, err := g.rollExpression(effect.Dice, g.cite)
		if err != nil {
			return err
		}

		count *= rolled
	}

	if count == 0 && effect.Modifier != 0 && effect.Scope == career.ScopeCareer {
		g.widenBenefitModifier(effect, name, cause)

		return nil
	}

	if effect.Modifier != 0 {
		g.char.Provenance.Deviate("E-3")
	}

	g.char.State.Benefits = append(g.char.State.Benefits, BenefitBatch{
		Career:   name,
		Rolls:    count,
		Modifier: effect.Modifier,
		CashOnly: effect.CashOnly,
	})

	g.consequence(ConsequenceBenefitRolls, cause, effect.Detail, name)

	return nil
}

// widenBenefitModifier is the ScopeCareer half of ERRATA E-3: a modifier
// granted with no rolls beside it reaches every batch this career has
// already collected, and every one it collects later.
func (g *Generator) widenBenefitModifier(effect career.Effect, name string, cause int) {
	for i := range g.char.State.Benefits {
		if g.char.State.Benefits[i].Career == name {
			g.char.State.Benefits[i].Modifier += effect.Modifier
		}
	}

	g.careerBenefitMod += effect.Modifier
	g.char.Provenance.Deviate("E-3")
	g.consequence(ConsequenceBenefitRolls, cause, effect.Detail, name)
}

// forfeitBenefits is the sixty-seven results that take back everything
// earned in this career: the queued batches go, and the two-per-term grant
// is written off up to the term the result fired in.
//
// It is a watermark rather than a switch because a handful of these are
// events rather than mishaps -- Scavenger's law-enforcement visit is one --
// and the character stays in the career afterwards. What was earned before
// the result is gone; the terms after it earn as they always did.
func (g *Generator) forfeitBenefits(effect career.Effect, name string, cause int) {
	g.clearBenefits(name)

	if g.termsInCareer > g.forfeitedTerms {
		g.forfeitedTerms = g.termsInCareer
	}

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

	start := effect.Rating
	if start == ratingUnspecified {
		start = defaultRating(effect.Relationship)
	}

	for range count {
		g.char.State.Ties = append(g.char.State.Ties, Tie{
			Kind:   string(effect.Relationship),
			Origin: name,
			Rating: start,
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

// changeStash adds to what a character owns, or takes something away.
//
// Three shapes: an empty Item empties the stash, which is what the results
// that say "lose your stash" mean; Lose removes every entry of one name,
// which is what "lose any Company Shares" means; and otherwise Count
// entries of that name are added, each with whatever value the book rolled
// for it.
func (g *Generator) changeStash(effect career.Effect, cause int) error {
	name := ""
	if g.career != nil {
		name = g.career.Name
	}

	switch {
	case effect.Item == "":
		g.char.State.Stash = nil
		g.consequence(ConsequenceStash, cause, "lose the contents of the stash", name)

		return nil
	case effect.Lose:
		g.loseFromStash(effect, name, cause)

		return nil
	}

	count := effect.Count

	if effect.CountDice != "" {
		rolled, err := g.rollExpression(effect.CountDice, g.cite)
		if err != nil {
			return err
		}

		count = rolled
	}

	// Each is valued on its own throw, because the rows that grant several
	// say "each": "Six Pieces of Art, 2D6 x Cr10000 each" (p. 155).
	for range max(count, 1) {
		value := 0

		if effect.Dice != "" {
			rolled, err := g.rollExpression(effect.Dice, g.cite)
			if err != nil {
				return err
			}

			value = rolled
		}

		g.char.State.Stash = append(g.char.State.Stash,
			Possession{Item: effect.Item, Value: value, Career: name})
	}

	g.consequence(ConsequenceStash, cause, effect.Detail, name)

	return nil
}

// loseFromStash removes every possession of one name, which is what the
// fourteen "lose any Company Shares" results ask for. A character who has
// none loses nothing, and the record says so rather than staying silent:
// the result fired, and what it did is part of what happened.
func (g *Generator) loseFromStash(effect career.Effect, name string, cause int) {
	kept := make([]Possession, 0, len(g.char.State.Stash))
	lost := 0

	for _, held := range g.char.State.Stash {
		if held.Item == effect.Item {
			lost++

			continue
		}

		kept = append(kept, held)
	}

	g.char.State.Stash = kept

	detail := effect.Detail
	if lost == 0 {
		detail += ", and there were none"
	}

	g.consequence(ConsequenceStash, cause, detail, name)
}

// rollExpression evaluates the small language the tables are written in:
// a bare number, "NdM", "NdMxK", "NdM+K" or "NdM-K", where M is 6, 3 or 10
// -- the three dice the book throws for a quantity. Anything else is a
// transcription error rather than a rule, so it errors rather than guessing
// what the page meant.
//
// The multiplier and the offset are exclusive because the book never writes
// both: "1d6 x ₶100,000" is a sum of money and "1d3+1 Contacts" is a count
// of people.
func (g *Generator) rollExpression(expr, cite string) (int, error) {
	expr, offset, ok := suffix(expr, "+", 0)
	if !ok {
		return 0, ErrBadExpression
	}

	// A species' characteristic method subtracts as often as it adds:
	// "roll 2d6-2 for STR and END" (p. 23). Only one of the two suffixes
	// can apply, because no expression in the book carries both.
	if offset == 0 {
		var down int

		expr, down, ok = suffix(expr, "-", 0)
		if !ok {
			return 0, ErrBadExpression
		}

		offset = -down
	}

	expr, multiplier, ok := suffix(expr, "x", 1)
	if !ok {
		return 0, ErrBadExpression
	}

	index := strings.Index(expr, "d")
	if index < 0 {
		value, valid := atoi(expr)
		if !valid {
			return 0, ErrBadExpression
		}

		return value*multiplier + offset, nil
	}

	count, valid := atoi(expr[:index])
	if !valid {
		return 0, ErrBadExpression
	}

	sides, valid := atoi(expr[index+1:])
	if !valid {
		return 0, ErrBadExpression
	}

	total, ok := g.throwDice(count, sides, cite)
	if !ok {
		return 0, ErrBadExpression
	}

	return total*multiplier + offset, nil
}

// throwDice throws count dice of one shape and sums them. Three shapes are
// all the book asks for a quantity: d6, d10 and d3.
func (g *Generator) throwDice(count, sides int, cite string) (int, bool) {
	switch sides {
	case sixSided:
		roll := g.dice.ND6(count)
		g.log.Roll(roll, cite)

		return roll.Total, true
	case tenSided:
		roll := g.dice.ND10(count)
		g.log.Roll(roll, cite)

		return roll.Total, true
	case threeSided:
		total := 0

		for range count {
			roll := g.dice.D3()
			g.log.Roll(roll, cite)

			total += roll.Total
		}

		return total, true
	}

	return 0, false
}

// suffix splits a trailing "<sep><number>" off an expression, returning the
// rest and the number -- or absent is the default, which is what an
// expression without that suffix means. A separator with nothing readable
// after it is a transcription error, not a default.
func suffix(expr, sep string, absent int) (string, int, bool) {
	rest, after, found := strings.Cut(expr, sep)
	if !found {
		return expr, absent, true
	}

	value, ok := atoi(after)
	if !ok {
		return expr, 0, false
	}

	return rest, value, true
}

// The two dice the book throws for a quantity, and the base its numbers
// are written in.
const (
	sixSided   = 6
	threeSided = 3
	tenSided   = 10
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
