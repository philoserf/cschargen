package chargen

import (
	"github.com/philoserf/cschargen/setting"
)

// Homeworld is one entry in a character's homeworld history. A character
// has more than one because six of Colonist's eleven mishaps reassign it
// (p. 174), and the record has to say which world they were on when.
type Homeworld struct {
	World     string `json:"world"`
	Subsector string `json:"subsector"`
	TechLevel int    `json:"techLevel"`
	FromTerm  int    `json:"fromTerm"`
	Reason    string `json:"reason,omitempty"`
}

// determineOrigin is Steps 3 and 4 (pp. 39-56): the subsector a character
// was born in, and the world within it.
func (g *Generator) determineOrigin() error {
	sub, err := g.chooseSubsector()
	if err != nil {
		return err
	}

	world, err := g.chooseHomeworld(sub)
	if err != nil {
		return err
	}

	step := g.log.Step("Step 4: Determine Homeworld", "p. 39")

	err = g.settleOn(world, sub, step, "born there")
	if err != nil {
		return err
	}

	// A character born outside the sector the campaign is set in owes it a
	// minimum number of terms (p. 43). The engine does not yet record where
	// a term was served, so it says so rather than pretending.
	if sub.Offworld {
		g.unimplemented(step,
			"a character born outside the sector must spend a minimum of four terms inside it (p. 43), "+
				"which the engine cannot check until it records where a term was served")
	}

	err = g.grantBackground(world, step)
	if err != nil {
		return err
	}

	// Step 5 is optional in the book -- "Skipping this step can speed up
	// character generation, but it may also deprive the character of
	// potential Allies and Contacts" (p. 57) -- so the engine runs it and a
	// flag turns it off.
	if g.char.Provenance.Inputs.SkipFamily {
		g.consequence(ConsequenceFamily, step,
			"Step 5 skipped at the player's request (p. 57)", "")

		return nil
	}

	return g.determineFamily(world)
}

// chooseSubsector is Step 3 (p. 39): "you start with rolling on the
// subsector of origin chart ... You may roll on the chart, the Referee may
// choose the subsector for the character, or the Referee may choose to
// allow the player to choose."
func (g *Generator) chooseSubsector() (setting.Subsector, error) {
	step := g.log.Step("Step 3: Determine Subsector of Origin", "p. 39")

	roll := g.dice.D6()
	cause := g.log.Roll(roll, "p. 39")

	for _, sub := range g.setting.Subsectors {
		if sub.OriginRoll == roll.Total {
			g.consequence(ConsequenceHomeworld, cause, "born in "+sub.Name, "")

			return sub, nil
		}
	}

	// The chart the book prints has a subsector for every result. A file
	// whose chart has a gap is not wrong -- a setting may have fewer than
	// six subsectors -- so the throw becomes a choice rather than an error.
	names := make([]string, 0, len(g.setting.Subsectors))
	for _, sub := range g.setting.Subsectors {
		names = append(names, sub.Name)
	}

	if len(names) == 0 {
		return setting.Subsector{}, ErrNoSubsectors
	}

	index, err := g.choose(Choice{
		Point:   "subsector",
		Prompt:  "The subsector roll landed on no chart entry; choose one",
		Options: names,
		Cite:    "p. 39",
	})
	if err != nil {
		return setting.Subsector{}, err
	}

	chosen := g.setting.Subsectors[index]
	g.consequence(ConsequenceHomeworld, step, "born in "+chosen.Name, "")

	return chosen, nil
}

// chooseHomeworld is Step 4 (p. 39): "You may either roll percentile dice
// (d100) to determine your homeworld randomly or simply choose a world that
// fits the character concept you have in mind."
func (g *Generator) chooseHomeworld(sub setting.Subsector) (setting.World, error) {
	percentile := g.dice.D100()
	cause := g.log.Roll(percentile, "p. 39")

	for _, world := range sub.Worlds {
		if world.Covers(percentile.Total) {
			g.consequence(ConsequenceHomeworld, cause,
				"a d100 of "+itoa(percentile.Total)+" lands on "+world.Name, "")

			return world, nil
		}
	}

	// Every subsector with a selectable world covers all hundred results --
	// the validator holds that -- so this is a subsector whose worlds are
	// all choose-only, which is what Recently Colonized Worlds is (p. 40).
	names := make([]string, 0, len(sub.Worlds))
	for _, world := range sub.Worlds {
		names = append(names, world.Name)
	}

	if len(names) == 0 {
		return setting.World{}, ErrNoWorlds
	}

	index, err := g.choose(Choice{
		Point:   "homeworld",
		Prompt:  "Choose a homeworld in " + sub.Name,
		Options: names,
		Cite:    "p. 39",
	})
	if err != nil {
		return setting.World{}, err
	}

	return sub.Worlds[index], nil
}

// settleOn records a homeworld, whether at birth or after a mishap sent the
// character elsewhere.
//
// ERRATA E-9: a homeworld reassigned during a career changes the tech level
// that gates the aging throws, and nothing else. Background skills are
// "learned through living their normal life on this world" (p. 40), which
// is a childhood, and maximum terms "represents the longest possible career
// history that character could have accumulated" from where they were born
// (p. 42). Neither is a property of where an adult happens to be living.
func (g *Generator) settleOn(world setting.World, sub setting.Subsector, cause int, reason string) error {
	first := len(g.char.State.Homeworlds) == 0

	g.char.State.Homeworlds = append(g.char.State.Homeworlds, Homeworld{
		World:     world.Name,
		Subsector: sub.Name,
		TechLevel: world.TechLevel,
		FromTerm:  len(g.char.State.Terms),
		Reason:    reason,
	})

	g.techLevel = world.TechLevel

	if first {
		g.char.Provenance.Inputs.Homeworld = world.Name
		g.homeworldTerms = world.MaximumTerms
		g.maximumAge = world.MaximumAge
	} else {
		g.char.Provenance.Deviate("E-9")
	}

	g.consequence(ConsequenceHomeworld, cause,
		world.Name+" in "+sub.Name+": "+reason, "")

	return nil
}

// reassignHomeworld carries out the six Colonist mishaps and the several
// other results that send a character to a new world (p. 174).
//
// The new world is thrown for as a birth world would be, because the book
// says "roll on a homeworld chart" and gives no other procedure. What
// changes is narrower than what a birth world sets: see ERRATA E-9 on
// [Generator.settleOn].
func (g *Generator) reassignHomeworld(cause int, detail string) error {
	sub, err := g.chooseSubsector()
	if err != nil {
		return err
	}

	world, err := g.chooseHomeworld(sub)
	if err != nil {
		return err
	}

	return g.settleOn(world, sub, cause, detail)
}
