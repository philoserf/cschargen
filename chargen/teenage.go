package chargen

import "github.com/philoserf/cschargen/career"

// Step 7: Teenage Events (pp. 75-85).
//
// Two rolls, ages 13-15 and 16-18, on one of four paths. Its gates differ
// from Step 6's: Paths 1 and 2 are the two halves of one condition -- has
// the homeworld been settled a hundred standard years -- so exactly one of
// them is always open, and Paths 3 and 4 are characteristic gates on top of
// that. A character always has one path and may have three.

// teenagePeriods is what a human's two rolls represent, in the book's own
// words. An uplift's are in the setting data: "Uplifts will often have a
// shorter adolescence than humans" (p. 76).
var teenagePeriods = [...]string{"ages 13-15", "ages 16-18"}

// teenageEvents is Step 7.
func (g *Generator) teenageEvents() error {
	step := g.log.Step("Step 7: Teenage Events", "p. 75")

	// The book offers this one as a choice too: "Some Referees may prefer
	// to allow players to select events rather than rolling randomly"
	// (p. 76).
	if g.char.Provenance.Inputs.SkipTeenage {
		g.consequence(ConsequenceFamily, step,
			"Step 7 skipped at the player's request (p. 76)", "")

		return nil
	}

	periods := g.lifePeriods(
		speciesTeenRolls(g.species), speciesTeenAges(g.species), teenagePeriods[:])

	for _, period := range periods {
		err := g.oneTeenageEvent(step, period)
		if err != nil {
			return err
		}
	}

	return nil
}

// oneTeenageEvent chooses a path, rolls 2d10 on it, and applies the result.
func (g *Generator) oneTeenageEvent(step int, period string) error {
	if g.enslaved {
		return g.enslavedEvent(career.EnslavedTeenage(), step, period)
	}

	path, err := g.chooseTeenagePath(step, period)
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

// chooseTeenagePath offers the paths this character qualifies for. One of
// Paths 1 and 2 is always open, so there is always at least one.
func (g *Generator) chooseTeenagePath(step int, period string) (career.TeenagePath, error) {
	score := g.characteristicScore()
	settled := g.settledYears()

	paths := career.TeenagePaths()

	// Paths 1 and 2 are the two sides of one condition, so exactly one of
	// them is open to every character. Choosing between those two first
	// says that, rather than leaving the loop below to discover it.
	homeworldPath := paths[1]
	if paths[0].Open(score, settled) {
		homeworldPath = paths[0]
	}

	open := []career.TeenagePath{homeworldPath}
	names := []string{homeworldPath.Name}

	for _, path := range paths[2:] {
		if !path.Open(score, settled) {
			continue
		}

		open = append(open, path)
		names = append(names, path.Name)
	}

	if len(open) == 1 {
		return open[0], nil
	}

	chosen, err := g.choose(Choice{
		Point:   "teenage_path",
		Prompt:  "Choose a teenage events path for " + period,
		Options: names,
		Cite:    "p. 76",
	})
	if err != nil {
		return career.TeenagePath{}, err
	}

	_ = step

	return open[chosen], nil
}

// settledYears is how long the character's homeworld has been colonized,
// which is what Paths 1 and 2 divide on. A world whose data gives no
// settled year reads as recently colonized, which is the harsher of the two
// paths and the safer default for a file that forgot the field.
func (g *Generator) settledYears() int {
	if g.settledYear == 0 || g.setting == nil {
		return 0
	}

	return g.setting.Now() - g.settledYear
}

// rollTeenageLifeEvent is the d6 table of p. 85.
func (g *Generator) rollTeenageLifeEvent(cause int) error {
	roll := g.dice.D6()
	throw := g.log.Roll(roll, "p. 85")
	row := career.TeenageLifeEvents()[roll.Total-1]

	g.consequence(ConsequenceFamily, throw, "teenage life event: "+row.Summary, "")

	_ = cause

	return g.applyAll(row.Effects, throw)
}
