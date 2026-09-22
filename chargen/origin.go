package chargen

import (
	"fmt"

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

	// The heading comes before the work it heads. It used to be written
	// after, which put Step 4's own d100 under Step 3 in the transcript.
	step := g.log.Step("Step 4: Determine Homeworld", "p. 39")

	world, err := g.chooseHomeworld(sub, step)
	if err != nil {
		return err
	}

	world, err = g.worldThatAdmitsThem(world, sub, step)
	if err != nil {
		return err
	}

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

	// The class needs the homeworld's tech level, so it is settled here
	// rather than in Step 1 where the species was chosen (p. 66).
	err = g.upliftClass(step)
	if err != nil {
		return err
	}

	err = g.grantBackground(world, step)
	if err != nil {
		return err
	}

	return g.earlyLife(world, step)
}

// earlyLife is Steps 5 through 8: the family, the youth, the teenage years
// and higher education. They are together because each follows the last and
// only the first of them needs the homeworld.
func (g *Generator) earlyLife(world setting.World, step int) error {
	// Step 5 is optional in the book -- "Skipping this step can speed up
	// character generation, but it may also deprive the character of
	// potential Allies and Contacts" (p. 57) -- so the engine runs it and a
	// flag turns it off.
	if g.char.Provenance.Inputs.SkipFamily {
		g.consequence(ConsequenceFamily, step,
			"Step 5 skipped at the player's request (p. 57)", "")
	} else {
		err := g.determineFamily(world)
		if err != nil {
			return err
		}
	}

	err := g.youthEvents()
	if err != nil {
		return err
	}

	err = g.teenageEvents()
	if err != nil {
		return err
	}

	return g.higherEducation()
}

// chooseSubsector is Step 3 (p. 39): "you start with rolling on the
// subsector of origin chart ... You may roll on the chart, the Referee may
// choose the subsector for the character, or the Referee may choose to
// allow the player to choose."
func (g *Generator) chooseSubsector() (setting.Subsector, error) {
	step := g.log.Step("Step 3: Determine Subsector of Origin", "p. 39")

	// Asked for, by name or by the subsector the asked-for homeworld is
	// in. p. 39 offers this before it offers the throw, and a subsector
	// named on the command line is the Referee having chosen.
	named, asked, err := g.namedSubsector(step)
	if err != nil {
		return setting.Subsector{}, err
	}

	if asked {
		return named, nil
	}

	// A player at the table is offered the same choice the page offers
	// them. The policy is not: it has no character concept to fit, so it
	// rolls, which is what the page describes first and what every record
	// written before this existed did.
	if g.char.Provenance.Inputs.Interactive {
		return g.offerSubsector(step)
	}

	return g.rollSubsector(step)
}

// namedSubsector is the subsector the inputs asked for, and whether they
// asked for one. A homeworld names one implicitly, and wins: it is the more
// specific request.
//
// The inputs say where the character was born, so a reassignment mid-career
// is free of them -- which is what the homeworld history being non-empty
// means here.
func (g *Generator) namedSubsector(step int) (setting.Subsector, bool, error) {
	inputs := g.char.Provenance.Inputs

	if len(g.char.State.Homeworlds) > 0 {
		return setting.Subsector{}, false, nil
	}

	if inputs.Homeworld != "" {
		_, sub, found := g.setting.World(inputs.Homeworld)
		if !found {
			return setting.Subsector{}, false,
				fmt.Errorf("%w: %s", ErrNoSuchHomeworld, inputs.Homeworld)
		}

		g.consequence(ConsequenceHomeworld, step,
			"born in "+sub.Name+", which is where "+inputs.Homeworld+" is ("+originCite+")", "")

		return sub, true, nil
	}

	if inputs.Subsector == "" {
		return setting.Subsector{}, false, nil
	}

	sub, found := g.setting.Subsector(inputs.Subsector)
	if !found {
		return setting.Subsector{}, false,
			fmt.Errorf("%w: %s", ErrNoSuchSubsector, inputs.Subsector)
	}

	g.consequence(ConsequenceHomeworld, step,
		"born in "+sub.Name+", chosen rather than rolled ("+originCite+")", "")

	return sub, true, nil
}

// offerSubsector puts p. 39's own choice to the player: roll on the chart,
// or name the subsector. Rolling is the first option because it is what the
// page describes first.
func (g *Generator) offerSubsector(step int) (setting.Subsector, error) {
	options := make([]string, 0, len(g.setting.Subsectors)+1)

	options = append(options, rollInstead)

	for _, sub := range g.setting.Subsectors {
		options = append(options, sub.Name)
	}

	index, err := g.choose(Choice{
		Point:   "subsector_source",
		Prompt:  "Roll for the subsector of origin, or choose one",
		Options: options,
		Cite:    originCite,
	})
	if err != nil {
		return setting.Subsector{}, err
	}

	if index == 0 {
		return g.rollSubsector(step)
	}

	sub := g.setting.Subsectors[index-1]

	g.consequence(ConsequenceHomeworld, step,
		"born in "+sub.Name+", chosen rather than rolled (p. 39)", "")

	return sub, nil
}

// rollInstead is the first option wherever the page offers a throw and a
// choice, so that the policy takes the throw.
const rollInstead = "roll for it"

// originCite is the page Steps 3 and 4 both begin on.
const originCite = "p. 39"

// rollSubsector throws on the chart of p. 39.
//
// The chart the book prints has an entry for all six results. A file whose
// chart has gaps is not wrong -- a setting may have fewer than six
// subsectors -- and a throw that lands in one is thrown again, which is what
// a table does with a result that is not on the chart. Throwing again keeps
// the chart's own weighting: a subsector claiming two results stays twice as
// likely as one claiming a single result.
//
// It used to become a choice instead, and under the auto policy that choice
// took the first subsector the file lists. On the repository's own sample,
// which claims three of the six results, that put 71 characters in 100 on
// one subsector and none at all on the fourth.
func (g *Generator) rollSubsector(step int) (setting.Subsector, error) {
	if !g.anySubsectorIsRollable() {
		return g.chooseAmongSubsectors(step)
	}

	for range maxOriginThrows {
		roll := g.dice.D6()
		cause := g.log.Roll(roll, originCite)

		for _, sub := range g.setting.Subsectors {
			if sub.OriginRoll == roll.Total {
				g.consequence(ConsequenceHomeworld, cause, "born in "+sub.Name, "")

				return sub, nil
			}
		}

		g.consequence(ConsequenceHomeworld, cause,
			"no subsector claims that result; throw again ("+originCite+")", "")
	}

	// A file claiming one result in six lands in six throws on average, so
	// reaching the cap is the dice rather than the data -- and saying so
	// beats throwing forever.
	return setting.Subsector{}, ErrOriginThrowsExhausted
}

// maxOriginThrows bounds the re-throw so that a pathological file cannot
// spin here. The worst a validated file can be is one subsector claiming a
// single result, where a throw misses five times in six: a hundred misses
// in a row is (5/6)^100, about one in eighty million. Twenty would be one
// in thirty-eight, which is not a bound at all -- a suite throwing this a
// dozen times would trip it about once a run.
const maxOriginThrows = 100

func (g *Generator) anySubsectorIsRollable() bool {
	for _, sub := range g.setting.Subsectors {
		if sub.OriginRoll != 0 {
			return true
		}
	}

	return false
}

// chooseAmongSubsectors is a file whose subsectors are all choose-only.
// Nothing can be thrown for, so the chart is not a chart and the question
// goes to whoever is deciding.
func (g *Generator) chooseAmongSubsectors(step int) (setting.Subsector, error) {
	names := make([]string, 0, len(g.setting.Subsectors))
	for _, sub := range g.setting.Subsectors {
		names = append(names, sub.Name)
	}

	if len(names) == 0 {
		return setting.Subsector{}, ErrNoSubsectors
	}

	index, err := g.choose(Choice{
		Point:   "subsector",
		Prompt:  "No subsector can be rolled for; choose one",
		Options: names,
		Cite:    originCite,
	})
	if err != nil {
		return setting.Subsector{}, err
	}

	chosen := g.setting.Subsectors[index]
	g.consequence(ConsequenceHomeworld, step, "born in "+chosen.Name, "")

	return chosen, nil
}

// worldThatAdmitsThem holds a homeworld against the character's species and
// finds another where it does not admit them (p. 42):
//
//	"If the entry indicates that [engineered humans] or uplifts are not
//	allowed, then the player should select a different homeworld, either by
//	re-rolling on the appropriate table or choosing another world that permits such
//	characters."
//
// The engine chooses rather than re-rolling, because a re-roll can land on
// the same world again and the page offers both. It takes the first world
// in the subsector that admits them, which is what the policy takes
// wherever it is offered a list -- and where the subsector admits them
// nowhere, the whole sector is searched, because a character has to be born
// somewhere.
func (g *Generator) worldThatAdmitsThem(
	world setting.World, sub setting.Subsector, step int,
) (setting.World, error) {
	status, admitted := g.permits(world)
	if admitted {
		g.recordStatus(status, world, step)

		return world, nil
	}

	g.consequence(ConsequenceHomeworld, step,
		world.Name+" does not admit a "+g.species.Name+"; another is chosen (p. 42)", "")

	for _, other := range sub.Worlds {
		status, admitted = g.permits(other)
		if !admitted {
			continue
		}

		g.recordStatus(status, other, step)

		return other, nil
	}

	for _, elsewhere := range g.setting.Subsectors {
		for _, other := range elsewhere.Worlds {
			status, admitted = g.permits(other)
			if !admitted {
				continue
			}

			g.recordStatus(status, other, step)

			return other, nil
		}
	}

	return setting.World{}, ErrNoHomeworldAdmitsThem
}

// recordStatus notes what the character's standing is where they were born,
// and remembers an enslavement: an engineered or uplift character born on a
// world where they are enslaved must take the slave career as their first
// career term (p. 42).
func (g *Generator) recordStatus(status setting.Status, world setting.World, step int) {
	if g.species == nil || status != setting.Enslaved {
		return
	}

	g.enslaved = true

	g.consequence(ConsequenceSpecies, step,
		"a "+g.species.Name+" born on "+world.Name+" is enslaved, and the "+
			"slave career is their first (p. 42)", "")
}

// chooseHomeworld is Step 4 (p. 39): "You may either roll percentile dice
// (d100) to determine your homeworld randomly or simply choose a world that
// fits the character concept you have in mind."
func (g *Generator) chooseHomeworld(sub setting.Subsector, step int) (setting.World, error) {
	// Asked for by name. Checked against this subsector rather than the
	// whole setting, because Step 3 has already put the character in it --
	// if a homeworld was named, that is the subsector it named.
	if named := g.char.Provenance.Inputs.Homeworld; named != "" && len(g.char.State.Homeworlds) == 0 {
		for _, world := range sub.Worlds {
			if world.Name == named {
				g.consequence(ConsequenceHomeworld, step,
					named+", chosen rather than rolled (p. 40)", "")

				return world, nil
			}
		}

		return setting.World{}, fmt.Errorf("%w: %s is not in %s", ErrNoSuchHomeworld, named, sub.Name)
	}

	if g.char.Provenance.Inputs.Interactive {
		return g.offerHomeworld(sub, step)
	}

	return g.rollHomeworld(sub)
}

// offerHomeworld puts p. 40's own choice to the player: "You may either
// roll percentile dice (d100) to determine your homeworld randomly or
// simply choose a world that fits the character concept you have in mind."
//
// Every world in the subsector is offered, including the ones no throw can
// reach -- which is the whole of the Recently Colonized Worlds table, whose
// page says "you cannot randomly be assigned one of these worlds, you may
// choose them" (p. 40). Before this they could not be reached at all.
func (g *Generator) offerHomeworld(sub setting.Subsector, step int) (setting.World, error) {
	options := make([]string, 0, len(sub.Worlds)+1)

	options = append(options, rollInstead)

	for _, world := range sub.Worlds {
		options = append(options, world.Name)
	}

	index, err := g.choose(Choice{
		Point:   "homeworld_source",
		Prompt:  "Roll for a homeworld in " + sub.Name + ", or choose one",
		Options: options,
		Cite:    "p. 40",
	})
	if err != nil {
		return setting.World{}, err
	}

	if index == 0 {
		return g.rollHomeworld(sub)
	}

	world := sub.Worlds[index-1]

	g.consequence(ConsequenceHomeworld, step,
		world.Name+", chosen rather than rolled (p. 40)", "")

	return world, nil
}

func (g *Generator) rollHomeworld(sub setting.Subsector) (setting.World, error) {
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
		Cite:    originCite,
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

	g.techLevel = g.techLevelFor(world, cause)
	g.settledYear = world.SettledYear

	if first {
		// The world a character ended up on is State.Homeworlds, not
		// Inputs. Inputs is what was asked for, the way Inputs.Career
		// keeps the career asked for beside the career served -- and
		// writing the result back here made a record indistinguishable
		// from one that had named its own homeworld, which replay then
		// honoured instead of re-rolling.

		// Engineered and uplift characters age differently and the
		// restrictions do not apply to them (p. 42), so only a human
		// carries their homeworld's caps.
		if g.ageLimitsApply() {
			g.homeworldTerms = g.homeworldTermsFor(world, cause)
			g.maximumAge = world.MaximumAge
		} else {
			g.consequence(ConsequenceHomeworld, cause,
				"the homeworld's maximum age and terms do not bind an "+
					g.species.Kind+" character (p. 42)", "")
		}
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

	world, err := g.chooseHomeworld(sub, cause)
	if err != nil {
		return err
	}

	return g.settleOn(world, sub, cause, detail)
}

// techLevelFor is the tech level the aging tables are read at (pp. 122-123),
// which is the homeworld's unless the command line named one.
//
// `--tech-level` is how a referee generates a character as if from a world
// their setting file does not have. A reassignment mid-career moves the
// character to a world with its own tech level and the flag does not follow
// them there: it describes where they were born, the way `--homeworld`
// does.
func (g *Generator) techLevelFor(world setting.World, cause int) int {
	asked := g.char.Provenance.Inputs.TechLevel
	if asked == 0 || len(g.char.State.Homeworlds) > 1 {
		return world.TechLevel
	}

	g.consequence(ConsequenceHomeworld, cause,
		"the aging tables are read at tech level "+itoa(asked)+
			" rather than "+world.Name+"'s "+itoa(world.TechLevel)+" (p. 122)", "")

	return asked
}

// homeworldTermsFor is the ceiling the homeworld puts on a career (p. 42),
// which is the world's unless the command line named one.
//
// It is one of the term limit's two ceilings and the lower wins (p. 125), so
// naming a smaller one here binds and naming a larger one leaves the
// policy's `--terms` to decide.
func (g *Generator) homeworldTermsFor(world setting.World, cause int) int {
	asked := g.char.Provenance.Inputs.MaxTerms
	if asked == 0 {
		return world.MaximumTerms
	}

	g.consequence(ConsequenceHomeworld, cause,
		"a homeworld maximum of "+itoa(asked)+" terms rather than "+
			world.Name+"'s "+itoa(world.MaximumTerms)+" (p. 42)", "")

	return asked
}
