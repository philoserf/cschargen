package chargen

import (
	"slices"
	"strings"
	"testing"

	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/setting"
)

// TestRaisingAHeldSkillWithNoneHeld is reachable: the result appears on
// tables a character can meet before any skill has been granted, and a
// choice with no options is not a question to ask.
func TestRaisingAHeldSkillWithNoneHeld(t *testing.T) {
	t.Parallel()

	gen := engine(t, 61)

	err := gen.raiseHeldSkill(career.Effect{
		Kind: career.EffectRaiseHeld, Detail: "raise a skill already held",
	}, 0)
	if err != nil {
		t.Fatalf("raiseHeldSkill: %v", err)
	}

	if !strings.Contains(lastDetail(t, gen), "holds none") {
		t.Errorf("the record does not say why nothing happened: %q", lastDetail(t, gen))
	}
}

// TestRaisingAHeldSkillTakesItsSpecialty. "A skill you already possess" is
// one entry of the sheet, not a name that might have several: raising
// Melee (Blade) must not create a bare Melee beside it.
func TestRaisingAHeldSkillTakesItsSpecialty(t *testing.T) {
	t.Parallel()

	gen := engine(t, 67)

	gen.char.State.GainSkill("Melee", "Blade", 1)

	err := gen.raiseHeldSkill(career.Effect{
		Kind: career.EffectRaiseHeld, Detail: "raise a skill already held",
	}, 0)
	if err != nil {
		t.Fatalf("raiseHeldSkill: %v", err)
	}

	if len(gen.char.State.Skills) != 1 {
		t.Fatalf("%d skills, want 1: %v", len(gen.char.State.Skills), gen.char.State.Skills)
	}

	if gen.char.State.Skills[0].Level != 2 {
		t.Errorf("Melee (Blade) is at %d, want 2", gen.char.State.Skills[0].Level)
	}
}

// TestARefusedSkillChoiceComesBack holds both choice points: a decider that
// answers nothing stops generation rather than picking for the player.
func TestARefusedSkillChoiceComesBack(t *testing.T) {
	t.Parallel()

	for name, effect := range map[string]career.Effect{
		"held": {Kind: career.EffectRaiseHeld, Detail: "raise a skill already held"},
		"any":  {Kind: career.EffectAnySkill, Detail: "any skill"},
	} {
		gen := engine(t, 71)

		gen.char.State.GainSkill("Melee", "Blade", 1)

		gen.decider = refusingDecider{}

		err := gen.apply(effect, 0)
		if err == nil {
			t.Errorf("%s: a refused choice did not come back", name)
		}
	}
}

// TestAnySkillNotHeldWithEverythingHeld is the far end of "any skill you do
// not already possess": a character who holds all forty-four gains nothing,
// and the record says which it was.
func TestAnySkillNotHeldWithEverythingHeld(t *testing.T) {
	t.Parallel()

	gen := engine(t, 73)

	for _, definition := range career.Skills() {
		gen.char.State.GainSkill(definition.Name, "", 1)
	}

	err := gen.anySkill(career.Effect{
		Kind: career.EffectAnySkill, NotHeld: true, Detail: "any skill not held",
	}, 0)
	if err != nil {
		t.Fatalf("anySkill: %v", err)
	}

	if !strings.Contains(lastDetail(t, gen), "already holds every skill") {
		t.Errorf("the record does not say why nothing happened: %q", lastDetail(t, gen))
	}
}

// TestLosingASkillNobodyHas is Orbital Construction's sixth mishap meeting
// a character who never had a vacc suit. The result fired, so the record
// carries it.
func TestLosingASkillNobodyHas(t *testing.T) {
	t.Parallel()

	gen := engine(t, 79)

	gen.loseSkill(career.Effect{
		Kind: career.EffectLoseSkill, Skill: "Suit", Specialties: []string{"Vacc Suit"},
		Detail: "lose every level of Suit (Vacc Suit)",
	}, 0)

	if !strings.Contains(lastDetail(t, gen), "did not hold") {
		t.Errorf("the record does not say the skill was absent: %q", lastDetail(t, gen))
	}
}

// TestAHomeworldSpecialtyComesFromTheChart is the five youth and teenage
// results that ask for "a Survival specialty which is used on your
// homeworld": the options are the ones the world's own background skills
// name, and nowhere else.
func TestAHomeworldSpecialtyComesFromTheChart(t *testing.T) {
	t.Parallel()

	world := testWorld("Tinderfall", nil)

	world.BackgroundSkills = []setting.Requirement{
		{OneOf: []setting.Alternative{
			{Skill: "Survival", Specialties: []string{"Desert", "Heat"}},
			{Skill: "Recon"},
		}},
		{OneOf: []setting.Alternative{{Skill: "Survival", Specialties: []string{"Desert"}}}},
	}

	gen := engine(t, 89)

	gen.setting = settingWith(setting.Subsector{Name: "New Holdings", Worlds: []setting.World{world}})
	gen.char.State.Homeworlds = []Homeworld{{World: "Tinderfall"}}

	options := gen.specialtiesOnTheHomeworld("Survival")
	if !slices.Equal(options, []string{"Desert", "Heat"}) {
		t.Errorf("the homeworld offers %v, want [Desert Heat] once each", options)
	}

	err := gen.homeworldSpecialty(career.Effect{
		Kind: career.EffectHomeworldSkill, Skill: "Survival", Detail: "a Survival specialty",
	}, 0)
	if err != nil {
		t.Fatalf("homeworldSpecialty: %v", err)
	}

	if gen.char.State.SkillLevel("Survival") != 1 {
		t.Errorf("no Survival specialty was granted: %v", gen.char.State.Skills)
	}
}

// TestAHomeworldWithNoSuchSpecialty records rather than invents. Several
// worlds' background skills name no Survival at all; a character may have
// no homeworld yet; and a record may name a world the loaded data does not
// carry. None of the three is an error, and all three say so.
func TestAHomeworldWithNoSuchSpecialty(t *testing.T) {
	t.Parallel()

	for name, homeworlds := range map[string][]Homeworld{
		"none yet":    nil,
		"not in data": {{World: "Somewhere Else"}},
		"no Survival": {{World: "Tinderfall"}},
	} {
		gen := engine(t, 83)

		gen.setting = settingWith(setting.Subsector{
			Name: "New Holdings", Worlds: []setting.World{testWorld("Tinderfall", nil)},
		})
		gen.char.State.Homeworlds = homeworlds

		err := gen.homeworldSpecialty(career.Effect{
			Kind: career.EffectHomeworldSkill, Skill: "Survival",
			Detail: "a Survival specialty from the homeworld",
		}, 0)
		if err != nil {
			t.Fatalf("%s: homeworldSpecialty: %v", name, err)
		}

		if !strings.Contains(lastDetail(t, gen), "name none") {
			t.Errorf("%s: the record does not say why nothing happened: %q",
				name, lastDetail(t, gen))
		}
	}
}
