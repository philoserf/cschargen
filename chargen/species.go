package chargen

import "github.com/philoserf/cschargen/setting"

// Step 1: Choose Human, Engineered Human, or Uplift (p. 21).
//
// The page offers three: a human moves straight on to Step 2, and an
// engineered human or an uplift picks a type, whose characteristic method
// Step 2 then uses.
//
// The engine holds no species. The data file declares them, and each one
// says which characteristic method it uses, which aging profile it ages on,
// what it starts with and what its ceiling is. No file here names a
// species, which is what the Product Identity boundary requires: every
// species heading in the book is a proper name the OGL notice reserves
// (p. 335).

// Human is what Inputs.Species holds for a character who is not one of the
// data file's species. It is not a species in the file and never needs to
// be.
const Human = "human"

// chooseSpecies is Step 1. It resolves the requested species against the
// data file and settles what the rest of generation reads from it.
func (g *Generator) chooseSpecies() error {
	step := g.log.Step("Step 1: Choose Human, Engineered Human, or Uplift", "p. 21")

	g.agingProfile = ProfileTechLevel
	g.maximum = HumanMaximum

	name := g.char.Provenance.Inputs.Species
	if name == "" || name == Human {
		g.char.Provenance.Inputs.Species = Human
		g.consequence(ConsequenceSpecies, step, "a baseline human", "")

		return nil
	}

	found, ok := g.setting.SpeciesNamed(name)
	if !ok {
		return ErrUnknownSpecies
	}

	g.species = &found

	if found.Aging != "" {
		g.agingProfile = AgingProfile(found.Aging)
	}

	if found.Maximum > 0 {
		g.maximum = found.Maximum
	}

	g.consequence(ConsequenceSpecies, step,
		"a "+found.Kind+" species: "+found.Name+", ageing on the "+
			string(g.agingProfile)+" profile, with a characteristic ceiling of "+
			itoa(g.maximum), "")

	return nil
}

// speciesMethod is the dice expression this species rolls a characteristic
// with, and whether it has one at all. A characteristic the species does
// not mention is "rolled normally" (p. 23), which is the human method and
// the free assignment that goes with it.
func (g *Generator) speciesMethod(which Characteristic) (string, bool) {
	// A hybrid may come out "otherwise human", and then the species' own
	// method does not apply at all (p. 63).
	if g.species == nil || g.rollAsHuman {
		return "", false
	}

	expr, ok := g.species.Characteristics[which.String()]

	return expr, ok
}

// grantSpeciesSkills is the last part of Step 1 (p. 23): a species may be
// owed skills at level 1 by what it is, named in the data file rather than
// here, because the page names them by species.
//
// They are granted after the characteristics rather than before, because a
// species' skills are a fact about the species and the characteristics are
// what Step 2 rolls -- and because a skill granted before Step 2 would be
// the only thing on a sheet above the six scores.
func (g *Generator) grantSpeciesSkills(step int) error {
	if g.species == nil || len(g.species.Skills) == 0 {
		return nil
	}

	for _, alternative := range g.species.Skills {
		specialty, err := g.chooseSpecialty(alternative, step)
		if err != nil {
			return err
		}

		granted := g.char.State.GainSkill(alternative.Skill, specialty, 1)

		g.log.Consequence(ConsequenceEvent{
			Kind:   ConsequenceSkill,
			Cause:  step,
			Detail: "a " + g.species.Name + " begins with " + granted.Full() + " 1",
			Skill:  granted.Full(),
			Level:  granted.Level,
			Cite:   "p. 21",
		})
	}

	return nil
}

// ceiling is this character's characteristic maximum, which is the
// species' where they have one and p. 14's fifteen where they do not.
func (g *Generator) ceiling() int {
	if g.maximum > 0 {
		return g.maximum
	}

	return HumanMaximum
}

// permits reports whether this world admits the character's species, and
// what their status there would be (p. 42).
//
// A world whose entry bars this kind of character sends the player to a
// different homeworld, by re-rolling on the chart or choosing another world
// that permits them.
func (g *Generator) permits(world setting.World) (setting.Status, bool) {
	if g.species == nil {
		return setting.Free, true
	}

	permission := world.Engineered
	if g.species.Kind == setting.KindUplift {
		permission = world.Uplifts
	}

	if !permission.Admits(g.species.Name) {
		return "", false
	}

	if permission.Status == "" {
		return setting.Free, true
	}

	return permission.Status, true
}

// ageLimitsApply reports whether a world's maximum age and terms bind this
// character. They do not bind an engineered person or an uplift, who age
// differently and to whom the restrictions do not apply (p. 42).
func (g *Generator) ageLimitsApply() bool {
	return g.species == nil
}

// Classes of uplifts (p. 66).
//
// "If the character's homeworld has a tech level of 10, they must be a Class
// 1 uplift. If the character's homeworld has a tech level of 11, the Player
// may choose from Class 1 or Class 2. If the character's homeworld has a
// tech level of 12, the Player may choose from any of the three classes."
//
// The classes are not species facts and are not in the data file: they are
// what a world's technology can make of any uplift.

// The tech level each class needs, from the chart on p. 66.
const (
	classOneTech   = 10
	classTwoTech   = 11
	classThreeTech = 12
)

// upliftClass is the class a character of this homeworld may be. It offers
// the choice where the tech level allows one, and takes the highest under
// the policy: "We highly recommend that the highest class available be
// used" (p. 66).
func (g *Generator) upliftClass(step int) error {
	if g.species == nil || g.species.Kind != setting.KindUplift {
		return nil
	}

	available := 1

	switch {
	case g.techLevel >= classThreeTech:
		available = 3
	case g.techLevel >= classTwoTech:
		available = 2
	}

	// A world below tech level 10 can make no uplift at all, and the
	// permission rules are what should have kept the character off it. A
	// data file that admits uplifts to such a world gets Class 1, which is
	// the least it can mean. ERRATA E-31.
	if g.techLevel < classOneTech {
		g.unimplemented(step,
			"the homeworld's tech level is below 10, which p. 66 says can make no uplift; "+
				"Class 1 is taken")
	}

	options := make([]string, available)
	for i := range options {
		options[i] = "Class " + itoa(available-i)
	}

	chosen := 0

	if available > 1 {
		picked, err := g.choose(Choice{
			Point:   "uplift_class",
			Prompt:  "Choose an uplift class",
			Options: options,
			Cite:    "p. 66",
		})
		if err != nil {
			return err
		}

		chosen = picked
	}

	g.class = available - chosen
	g.char.State.UpliftClass = g.class

	g.consequence(ConsequenceSpecies, step,
		"a Class "+itoa(g.class)+" uplift, which a tech level "+itoa(g.techLevel)+
			" homeworld can make", "")

	return nil
}
