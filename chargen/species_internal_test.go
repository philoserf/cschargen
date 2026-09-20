package chargen

import (
	"errors"
	"testing"

	"github.com/philoserf/cschargen/setting"
)

// speciesEngine is a generator whose data file declares one species.
func speciesEngine(t *testing.T, seed uint64, species setting.Species) *Generator {
	t.Helper()

	data := settingWith(setting.Subsector{
		Name: "New Holdings", OriginRoll: 0,
		Worlds: []setting.World{testWorld("Tinderfall", nil), testWorld("Wake", nil)},
	})

	data.Species = append(data.Species, species)

	return New(Options{
		Seed:          seed,
		Decider:       Policy{},
		EngineVersion: testVersion,
		PolicyVersion: testVersion,
		Setting:       data,
		Inputs:        Inputs{Species: species.Name, TermLimit: -1},
	})
}

// TestASpeciesMethodThatCannotBeRolled. A characteristic method is a dice
// expression like any other, so a mistranscribed one fails where every
// other mistranscribed expression fails.
func TestASpeciesMethodThatCannotBeRolled(t *testing.T) {
	t.Parallel()

	gen := speciesEngine(t, 36, setting.Species{
		Name: "Broken", Kind: "uplift",
		Characteristics: map[string]string{"STR": "2d6 - 2"},
	})

	_, err := gen.Run()
	if !errors.Is(err, ErrBadExpression) {
		t.Errorf("err = %v, want ErrBadExpression", err)
	}
}

// TestARefusedSpeciesSkillEndsGeneration. A species' starting skills go
// through the same specialty choice every other skill does.
func TestARefusedSpeciesSkillEndsGeneration(t *testing.T) {
	t.Parallel()

	gen := speciesEngine(t, 37, setting.Species{
		Name: "Choosy", Kind: "uplift",
		Skills: []setting.Alternative{
			{Skill: "Survival", Specialties: []string{"Ocean", "Freefall"}},
		},
	})

	// Step 1's skills are granted after Step 2's assignment, so a decider
	// that refused everything would stop at the assignment and never reach
	// them. This one refuses only the specialty.
	gen.species = &setting.Species{
		Name: "Choosy", Kind: "uplift",
		Skills: []setting.Alternative{
			{Skill: "Survival", Specialties: []string{"Ocean", "Freefall"}},
		},
	}
	gen.decider = refusingDecider{}

	err := gen.grantSpeciesSkills(0)
	if err == nil {
		t.Fatal("a refused specialty did not end the step")
	}
}

// TestASpeciesWithNothingToSay. Every field is optional: a species that
// declares only a name and a kind rolls characteristics the human way, ages
// on the tech-level profile and starts with nothing.
func TestASpeciesWithNothingToSay(t *testing.T) {
	t.Parallel()

	gen := speciesEngine(t, 38, setting.Species{Name: "Plain", Kind: "engineered"})

	character, err := gen.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if gen.agingProfile != ProfileTechLevel {
		t.Errorf("aging profile = %q, want the default", gen.agingProfile)
	}

	if gen.ceiling() != HumanMaximum {
		t.Errorf("ceiling = %d, want the human %d", gen.ceiling(), HumanMaximum)
	}

	if character.Provenance.Inputs.Species != "Plain" {
		t.Errorf("the record says %q", character.Provenance.Inputs.Species)
	}
}
