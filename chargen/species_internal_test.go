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

// TestAWorldWithNoStatedStatus. A permission that admits a species and says
// nothing about their standing is a world where they are free: most worlds
// "do not enslave altrants or uplifts, and such characters live as free
// citizens or residents" (p. 42).
func TestAWorldWithNoStatedStatus(t *testing.T) {
	t.Parallel()

	gen := speciesEngine(t, 39, setting.Species{Name: "Quiet", Kind: setting.KindUplift})

	gen.species = &setting.Species{Name: "Quiet", Kind: setting.KindUplift}

	status, admitted := gen.permits(setting.World{
		Name:    "Nowhere",
		Uplifts: setting.Permission{Allowed: true},
	})
	if !admitted || status != setting.Free {
		t.Errorf("a world that admits them without a status gave %q, %v", status, admitted)
	}

	// And one that bars them admits nobody.
	_, admitted = gen.permits(setting.World{Name: "Closed"})
	if admitted {
		t.Error("a world that allows no uplifts admitted one")
	}
}

// TestARefusedUpliftClassEndsGeneration. The class is a choice wherever the
// tech level allows one (p. 66).
func TestARefusedUpliftClassEndsGeneration(t *testing.T) {
	t.Parallel()

	gen := speciesEngine(t, 40, setting.Species{Name: "Choosy", Kind: setting.KindUplift})

	gen.species = &setting.Species{Name: "Choosy", Kind: setting.KindUplift}
	gen.techLevel = classThreeTech
	gen.decider = refusingDecider{}

	err := gen.upliftClass(0)
	if err == nil {
		t.Fatal("a refused class did not come back as an error")
	}
}

// TestAWorldBelowTenThatAdmitsUplifts is ERRATA E-31: p. 66's chart starts
// at tech level 10, and a data file may put an uplift on a world below it.
func TestAWorldBelowTenThatAdmitsUplifts(t *testing.T) {
	t.Parallel()

	gen := speciesEngine(t, 41, setting.Species{Name: "Early", Kind: setting.KindUplift})

	gen.species = &setting.Species{Name: "Early", Kind: setting.KindUplift}
	gen.techLevel = 8

	err := gen.upliftClass(0)
	if err != nil {
		t.Fatalf("upliftClass: %v", err)
	}

	if gen.class != 1 {
		t.Errorf("class = %d, want the least a class can be", gen.class)
	}
}

// TestASubsectorThatWillNotHaveThem. p. 42 says to choose another world.
// Where the character's own subsector has none that admits them, the
// search widens to the whole sector -- and where the sector has none
// either, there is no character to generate.
func TestASubsectorThatWillNotHaveThem(t *testing.T) {
	t.Parallel()

	closed := testWorld("Closed", &[2]int{1, 100})

	closed.Uplifts = setting.Permission{Allowed: false}
	closed.Engineered = setting.Permission{Allowed: false}

	open := testWorld("Open", &[2]int{1, 100})

	data := settingWith(setting.Subsector{
		Name: "Shut", OriginRoll: 1, Worlds: []setting.World{closed},
	})

	data.Subsectors = append(data.Subsectors, setting.Subsector{
		Name: "Elsewhere", OriginRoll: 2, Worlds: []setting.World{open},
	})

	data.Species = append(data.Species, setting.Species{
		Name: "Wanderer", Kind: setting.KindUplift,
	})

	character, err := New(Options{
		Seed:          42,
		Decider:       Policy{},
		EngineVersion: testVersion,
		PolicyVersion: testVersion,
		Setting:       data,
		Inputs:        Inputs{Species: "Wanderer", TermLimit: -1},
	}).Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if born := character.State.Homeworlds[0].World; born != "Open" {
		t.Errorf("born on %s; only Open admits them", born)
	}
}

// TestNoWorldAnywhereWillHaveThem is the end of that search.
func TestNoWorldAnywhereWillHaveThem(t *testing.T) {
	t.Parallel()

	closed := testWorld("Closed", &[2]int{1, 100})

	closed.Uplifts = setting.Permission{Allowed: false}

	data := settingWith(setting.Subsector{
		Name: "Shut", OriginRoll: 1, Worlds: []setting.World{closed},
	})

	data.Species = append(data.Species, setting.Species{
		Name: "Homeless", Kind: setting.KindUplift,
	})

	_, err := New(Options{
		Seed:          43,
		Decider:       Policy{},
		EngineVersion: testVersion,
		PolicyVersion: testVersion,
		Setting:       data,
		Inputs:        Inputs{Species: "Homeless", TermLimit: -1},
	}).Run()
	if !errors.Is(err, ErrNoHomeworldAdmitsThem) {
		t.Errorf("err = %v, want ErrNoHomeworldAdmitsThem", err)
	}
}

// TestBeingFreed is result 12 on both enslaved tables and result 66 of the
// slave career: "Your owner has decided to set you free. Continue your
// character as a free altrant or uplift."
func TestBeingFreed(t *testing.T) {
	t.Parallel()

	gen := speciesEngine(t, 44, setting.Species{Name: "Owned", Kind: setting.KindUplift})

	gen.enslaved = true

	gen.freed(0)

	if gen.enslaved {
		t.Error("a freed character is still owned")
	}

	// And freeing a character who was never owned changes nothing and says
	// nothing, because there was nothing to end.
	before := gen.log.Len()

	gen.freed(0)

	if gen.log.Len() != before {
		t.Error("freeing a free character wrote to the record")
	}
}
