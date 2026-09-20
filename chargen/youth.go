package chargen

import (
	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/setting"
)

// Step 6: Youth Events (pp. 67-69).
//
// Two rolls, one for ages 4-8 and one for 9-12, on one of five paths the
// character's characteristics open. "If something happens which improves or
// decreases the characteristic scores of the character which allows them to
// change their path, the player may choose to do so for the second roll"
// (p. 68) -- so the path is chosen again between them, against whatever the
// first roll left behind.

// youthPeriods is what a human's two rolls represent, in the book's own
// words. An uplift's are in the setting data: "Uplifts will often have
// shorter youths than humans. Whereas humans will roll on the following
// charts twice to represent their lives from 4-8 and again from 9-12,
// uplifts will have a different age range and may roll fewer times"
// (p. 68).
var youthPeriods = [...]string{"ages 4-8", "ages 9-12"}

// lifePeriods is what this character's rolls on a life-period table
// represent. The species' own pattern wins where it declares one.
func (g *Generator) lifePeriods(rolls int, ages, human []string) []string {
	if g.species == nil || rolls == 0 || len(ages) != rolls {
		return human
	}

	return ages
}

// youthEvents is Step 6.
func (g *Generator) youthEvents() error {
	step := g.log.Step("Step 6: Youth Events", "p. 67")

	// Like Step 5, the book offers this one as a choice: "Players may
	// create this background entirely on their own, developing a childhood
	// history that fits the character concept they have in mind"
	// (p. 67). A player who has is not served by rolling on top of it.
	if g.char.Provenance.Inputs.SkipYouth {
		g.consequence(ConsequenceFamily, step,
			"Step 6 skipped at the player's request (p. 67)", "")

		return nil
	}

	periods := g.lifePeriods(
		speciesYouthRolls(g.species), speciesYouthAges(g.species), youthPeriods[:])

	for _, period := range periods {
		err := g.oneYouthEvent(step, period)
		if err != nil {
			return err
		}
	}

	return nil
}

// oneYouthEvent chooses a path, rolls 2d10 on it, and applies the result.
//
// A character their homeworld owns takes a different table: an eleven-row
// 2d6 rather than a nineteen-row 2d10, and the only one of the six that is
// not a choice (p. 74).
func (g *Generator) oneYouthEvent(step int, period string) error {
	if g.enslaved {
		return g.enslavedEvent(career.EnslavedYouth(), step, period)
	}

	path, err := g.chooseYouthPath(step, period)
	if err != nil {
		return err
	}

	g.cite = path.Cite

	roll := g.dice.ND10(youthDice)
	cause := g.log.Roll(roll, path.Cite)

	row := path.Rows[roll.Total-youthLow]

	g.consequence(ConsequenceFamily, cause,
		period+", on "+path.Name+": "+row.Summary, "")

	return g.applyAll(row.Effects, cause)
}

// The shape of a youth path's table: two ten-sided dice, so the lowest
// result is 2 and that is where its first row sits.
const (
	youthDice = 2
	youthLow  = 2
)

// chooseYouthPath offers the paths this character qualifies for. Path 1 has
// no requirements, so there is always at least one.
func (g *Generator) chooseYouthPath(step int, period string) (career.YouthPath, error) {
	score := g.characteristicScore()

	paths := career.YouthPaths()

	// Path 1 has no requirements (p. 68), so it is always open and is
	// always the first the book prints. Starting the list with it says so
	// rather than leaving it to the loop below to discover.
	open := []career.YouthPath{paths[0]}
	names := []string{paths[0].Name}

	for _, path := range paths[1:] {
		if !path.Open(score) {
			continue
		}

		open = append(open, path)
		names = append(names, path.Name)
	}

	if len(open) == 1 {
		return open[0], nil
	}

	chosen, err := g.choose(Choice{
		Point:   "youth_path",
		Prompt:  "Choose a youth events path for " + period,
		Options: names,
		Cite:    "p. 68",
	})
	if err != nil {
		return career.YouthPath{}, err
	}

	_ = step

	return open[chosen], nil
}

// characteristicScore looks a characteristic up by the name a table writes
// it in. A name this ruleset does not have -- Traveller's SOC, say --
// scores zero rather than opening a path by accident.
func (g *Generator) characteristicScore() func(string) int {
	return func(which string) int {
		target, ok := characteristicByName(which)
		if !ok {
			return 0
		}

		return g.char.State.Characteristics.Get(target)
	}
}

// rollYouthLifeEvent is the d6 table of p. 75, which result 10 of every
// path reaches.
func (g *Generator) rollYouthLifeEvent(cause int) error {
	roll := g.dice.D6()
	throw := g.log.Roll(roll, "p. 75")
	row := career.YouthLifeEvents()[roll.Total-1]

	g.consequence(ConsequenceFamily, throw, "youth life event: "+row.Summary, "")

	_ = cause

	return g.applyAll(row.Effects, throw)
}

// speciesYouthRolls and speciesYouthAges read a species' youth pattern
// without the caller having to know whether there is a species at all.
func speciesYouthRolls(species *setting.Species) int {
	if species == nil {
		return 0
	}

	return species.YouthRolls
}

func speciesYouthAges(species *setting.Species) []string {
	if species == nil {
		return nil
	}

	return species.YouthAges
}

// speciesTeenRolls and speciesTeenAges are the same for Step 7.
func speciesTeenRolls(species *setting.Species) int {
	if species == nil {
		return 0
	}

	return species.TeenRolls
}

func speciesTeenAges(species *setting.Species) []string {
	if species == nil {
		return nil
	}

	return species.TeenAges
}

// enslavedEvent rolls on one of the two tables of pp. 74 and 84.
func (g *Generator) enslavedEvent(
	path career.EnslavedPath, step int, period string,
) error {
	g.cite = path.Cite

	roll := g.dice.TwoD6()
	cause := g.log.Roll(roll, path.Cite)

	row := path.Rows[roll.Total-2]

	g.consequence(ConsequenceFamily, cause,
		period+", on the "+path.Name+" table: "+row.Summary, "")

	_ = step

	return g.applyAll(row.Effects, cause)
}

// freed ends an enslavement, which three results in the book do: "Continue
// your character as a free altrant or uplift."
//
// It clears the obligation to take the slave career first, where the
// character has not taken it yet.
func (g *Generator) freed(cause int) {
	if !g.enslaved {
		return
	}

	g.enslaved = false

	g.consequence(ConsequenceSpecies, cause, "freed, and no longer owned", "")
}
