package chargen

import (
	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/dice"
)

// serveTerm is Steps 12 to 16 (pp. 112-121): survival, mishap, advancement,
// skills, events -- one four-year term.
func (g *Generator) serveTerm() error {
	if g.career == nil {
		return ErrNoCareer
	}

	assignment, ok := g.career.Assignment(g.assignment.Name)
	if !ok {
		return ErrNoCareer
	}

	g.termsInCareer++

	term := Term{
		Career:     g.career.Name,
		Assignment: assignment.Name,
		Number:     len(g.char.State.Terms) + 1,
	}

	term.Survived = g.rollSurvival(assignment)

	err := g.resolveTerm(assignment, term.Survived)
	if err != nil {
		return err
	}

	term.Rank = g.rank
	g.char.State.Terms = append(g.char.State.Terms, term)

	if service, found := g.char.State.Service(g.career.Name); found {
		service.Terms = g.termsInCareer
		service.Rank = g.rank
	}

	return g.age()
}

// resolveTerm is the half of a term that depends on the survival throw: a
// mishap, or advancement, a skill and an event.
func (g *Generator) resolveTerm(assignment career.Assignment, survived bool) error {
	if !survived {
		step := g.log.Step("Step 13: Roll A Mishap", "p. 113")

		return g.rollMishap(step, true)
	}

	err := g.advance(assignment)
	if err != nil {
		return err
	}

	err = g.rollSkill(assignment)
	if err != nil {
		return err
	}

	return g.rollEvent()
}

// rollSurvival is Step 12 (p. 112). A natural twelve before modifiers
// requires another term in the same career (p. 113).
func (g *Generator) rollSurvival(assignment career.Assignment) bool {
	g.log.Step("Step 12: Roll for Survival", "p. 112")

	throw := g.characteristicThrow(assignment.Survival)
	cause := g.log.Throw(throw, "p. 112")

	if throw.Natural() == naturalTwelve {
		g.mustContinue = true
		g.consequence(ConsequenceCareer, cause,
			"a natural twelve on the survival roll: another term in this career", g.career.Name)
	}

	if !throw.Success {
		g.consequence(ConsequenceCareer, cause, "the survival roll failed", g.career.Name)

		return false
	}

	return true
}

// naturalTwelve is the exceptional success the book reads as a result in
// its own right (p. 113).
const naturalTwelve = 12

// advance is Step 14 (pp. 114-116). None of this milestone's three careers
// offers a commission, so the commission branch is a recorded absence
// rather than an implementation.
func (g *Generator) advance(assignment career.Assignment) error {
	step := g.log.Step("Step 14: Roll for Advancement", "p. 114")

	if g.autoAdvance {
		g.autoAdvance = false

		return g.promote(step, assignment)
	}

	mods := g.takeModifiers("next advancement roll")

	which, ok := characteristicByName(assignment.Advancement.Characteristic)
	if ok {
		mods = append(mods, dice.Mod{
			Name:  assignment.Advancement.Characteristic,
			Value: g.char.State.Characteristics.Modifier(which),
		})
	}

	throw := g.dice.Throw(assignment.Advancement.Number, mods...)
	cause := g.log.Throw(throw, "p. 115")

	if !throw.Success {
		return nil
	}

	return g.promote(cause, assignment)
}

// promote raises the rank by one and applies that rank's benefit, which is
// granted "immediately upon achieving that Rank" (p. 116).
func (g *Generator) promote(cause int, assignment career.Assignment) error {
	if g.rank >= len(assignment.Ranks)-1 {
		g.consequence(ConsequenceRank, cause,
			"already at the highest printed rank in this assignment", g.career.Name)

		return nil
	}

	g.rank++

	g.consequence(ConsequenceRank, cause,
		"advanced to rank "+itoa(g.rank)+" in "+assignment.Name, g.career.Name)

	return g.applyAll(assignment.Ranks[g.rank], cause)
}

// rollSkill is Step 15 (p. 117): "The Player should choose one of the
// tables and roll a d6."
func (g *Generator) rollSkill(assignment career.Assignment) error {
	g.log.Step("Step 15: Roll for Skills", skillStepCite)

	tables := g.availableSkillTables(assignment)

	names := make([]string, len(tables))
	for i, table := range tables {
		names[i] = table.Name
	}

	index, err := g.choose(Choice{
		Point:   "skill_table",
		Prompt:  "Choose a skill table to roll on",
		Options: names,
		Cite:    skillStepCite,
	})
	if err != nil {
		return err
	}

	roll := g.dice.D6()
	cause := g.log.Roll(roll, skillStepCite)

	return g.apply(tables[index].Rows[roll.Total-1], cause)
}

// availableSkillTables is the career's own tables, minus any the character
// is not eligible for, plus the assignment's. The Advanced Education table
// is "restricted to characters with an EDU of a certain level or higher"
// (p. 117), so an ineligible character is not offered it rather than being
// allowed to roll and then told no.
func (g *Generator) availableSkillTables(assignment career.Assignment) []career.SkillTable {
	var tables []career.SkillTable

	for _, table := range g.career.Tables {
		if table.MinimumEDU > 0 && g.char.State.Characteristics.EDU < table.MinimumEDU {
			continue
		}

		tables = append(tables, table)
	}

	return append(tables, assignment.Skills)
}

// rollEvent is Step 16 (p. 118).
func (g *Generator) rollEvent() error {
	g.log.Step("Step 16: Roll for Events", "p. 118")

	roll := g.dice.D66()
	cause := g.log.Roll(roll, "p. 118")

	row, ok := g.career.Events[roll.Total]
	if !ok {
		return ErrMissingEventRow
	}

	g.consequence(ConsequenceCareer, cause, "event: "+row.Summary, g.career.Name)

	return g.applyAll(row.Effects, cause)
}

// age is Step 17 (p. 121) as far as this milestone goes: four years per
// term, or 1d3 where a mishap ejected the character mid-term. The aging
// throws themselves are milestone 4 -- they are indexed by term number and
// gated by the homeworld's tech level (pp. 122-123), which is setting data
// milestone 2 brings in.
func (g *Generator) age() error {
	cause := g.log.Step("Step 17: Aging", "p. 121")

	years := termYears

	if g.ejected {
		roll := g.dice.D3()

		cause = g.log.Roll(roll, "p. 121")

		years = roll.Total
	}

	g.char.State.Age += years
	g.consequence(ConsequenceAge, cause, "age "+itoa(g.char.State.Age), "")

	return nil
}

// characteristicThrow rolls 2d6 plus a characteristic's modifier against a
// target, which is how every check in the book resolves (p. 110).
func (g *Generator) characteristicThrow(check career.Check) dice.Throw {
	which, ok := characteristicByName(check.Characteristic)
	if !ok {
		return g.dice.Throw(check.Number)
	}

	return g.dice.Throw(check.Number, dice.Mod{
		Name:  check.Characteristic,
		Value: g.char.State.Characteristics.Modifier(which),
	})
}

// takeModifiers consumes every pending modifier that applies to a named
// throw. Consuming rather than reading: "take a -2 on your next
// Advancement roll" is spent once.
func (g *Generator) takeModifiers(applies string) []dice.Mod {
	var (
		mods []dice.Mod
		kept []PendingModifier
	)

	for _, pending := range g.pending {
		if pending.Applies != applies {
			kept = append(kept, pending)

			continue
		}

		mods = append(mods, dice.Mod{Name: pending.Detail, Value: pending.Value})
	}

	g.pending = kept

	return mods
}
