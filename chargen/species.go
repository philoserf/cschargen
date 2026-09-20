package chargen

// Step 1: Choose Human, Altrant, or Uplift (p. 21).
//
// "If the character is a human, simply move on to Step 2. If the character
// is to be an altrant or an uplift, the Player should choose the type ...
// and note the method used for that altrant or uplift for characteristic
// generation to be used in Step 2."
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
	step := g.log.Step("Step 1: Choose Human, Altrant, or Uplift", "p. 21")

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
	if g.species == nil {
		return "", false
	}

	expr, ok := g.species.Characteristics[which.String()]

	return expr, ok
}

// grantSpeciesSkills is the last part of Step 1: "All Gaishan should be
// given Survival (Freefall) and Survival (Low Gravity) at level 1" (p. 23).
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
