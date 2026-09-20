package chargen

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/setting"
)

// The tests here reach the engine's own helpers. They are the paths a
// generated character rarely takes -- a malformed table, an unknown
// characteristic, a rank already at its ceiling -- and each of them is a
// branch that would otherwise be reasoned about rather than run.

// The strings these fixtures repeat, named once.
const (
	testSubsector = "Only"
	testLanguage  = "Trade"
	testHuman     = "human"
	testVersion   = "test"
)

func engine(t *testing.T, seed uint64) *Generator {
	t.Helper()

	return New(Options{
		Seed:          seed,
		Decider:       Policy{},
		EngineVersion: testVersion,
		PolicyVersion: testVersion,
		Inputs:        Inputs{Species: testHuman, TermLimit: -1},
	})
}

func TestCharacteristicByName(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"STR", "dex", "End", "INT", "edu", "CHA"} {
		_, ok := characteristicByName(name)
		if !ok {
			t.Errorf("characteristicByName(%q) did not resolve", name)
		}
	}

	_, ok := characteristicByName("SOC")
	if ok {
		t.Error("SOC resolved; Clement Sector has CHA where Traveller has SOC")
	}
}

// TestAdjustAnUnknownCharacteristicIsRecordedNotApplied: a transcription
// naming a characteristic that does not exist must not silently move a
// real one.
func TestAdjustAnUnknownCharacteristicIsRecordedNotApplied(t *testing.T) {
	t.Parallel()

	gen := engine(t, 1)
	before := gen.char.State.Characteristics

	gen.adjust("SOC", 2, "a characteristic that does not exist", 1)

	if gen.char.State.Characteristics != before {
		t.Error("an unknown characteristic changed the record")
	}

	last := gen.log.Events()[gen.log.Len()-1]
	if last.Consequence == nil || last.Consequence.Kind != ConsequenceUnimplemented {
		t.Error("the unknown characteristic was not recorded as unimplemented")
	}
}

// TestAdjustHoldsTheRange: nothing falls below zero, and an unaltered human
// caps at 15 (p. 14).
func TestAdjustHoldsTheRange(t *testing.T) {
	t.Parallel()

	gen := engine(t, 1)

	gen.char.State.Characteristics.Set(END, 3)
	gen.adjust("END", -9, "a heavy loss", 1)

	if got := gen.char.State.Characteristics.END; got != 0 {
		t.Errorf("END = %d after falling past zero, want 0", got)
	}

	gen.char.State.Characteristics.Set(STR, 14)
	gen.adjust("STR", +9, "a large gain", 1)

	if got := gen.char.State.Characteristics.STR; got != HumanMaximum {
		t.Errorf("STR = %d after rising past the cap, want %d", got, HumanMaximum)
	}
}

// TestRollCheckFallsBackWhenTheTargetNamesNothingReal: a malformed target
// still produces a throw rather than a panic, because a table result the
// engine cannot read must not take the whole lifepath with it.
func TestRollCheckFallsBackWhenTheTargetNamesNothingReal(t *testing.T) {
	t.Parallel()

	gen := engine(t, 1)

	throw := gen.rollCheck(career.Target{Characteristic: "SOC", Number: 8})
	if throw.Target != 8 {
		t.Errorf("target = %d", throw.Target)
	}

	if len(throw.Mods) != 0 {
		t.Errorf("an unresolvable characteristic contributed %v", throw.Mods)
	}
}

// TestRollCheckOnASkillUsesItsLevel is ERRATA E-7, and it stamps the
// reading on the record.
func TestRollCheckOnASkillUsesItsLevel(t *testing.T) {
	t.Parallel()

	gen := engine(t, 1)
	gen.char.State.GainSkill("Gambler", "", 3)

	throw := gen.rollCheck(career.Target{Skill: "Gambler", Number: 8})
	if len(throw.Mods) != 1 || throw.Mods[0].Value != 3 {
		t.Errorf("mods = %v, want one of +3", throw.Mods)
	}

	deviated := false

	for _, id := range gen.char.Provenance.Deviations {
		if id == "E-7" {
			deviated = true
		}
	}

	if !deviated {
		t.Error("a skill check did not stamp E-7, the reading it rests on")
	}
}

func TestRollCheckOnASkillTheCharacterLacks(t *testing.T) {
	t.Parallel()

	gen := engine(t, 1)

	throw := gen.rollCheck(career.Target{Skill: "Gambler", Number: 8})
	if len(throw.Mods) != 1 || throw.Mods[0].Value != 0 {
		t.Errorf("mods = %v, want one of +0", throw.Mods)
	}
}

// TestBenefitScopesAreNotTheSameThing is ERRATA E-3: "a +1 modifier to all
// Benefit rolls made in this career" reaches batches already granted;
// "three Benefit rolls at +1" reaches only its own three.
func TestBenefitScopesAreNotTheSameThing(t *testing.T) {
	t.Parallel()

	gen := engine(t, 1)
	colonist := career.Colonist()

	gen.career = &colonist

	gen.grantBenefits(career.Effect{Count: 2, Detail: "two rolls"}, 1)
	gen.grantBenefits(career.Effect{Count: 3, Modifier: 1, Detail: "three at +1"}, 1)

	if len(gen.char.State.Benefits) != 2 {
		t.Fatalf("%d batches, want 2", len(gen.char.State.Benefits))
	}

	if gen.char.State.Benefits[0].Modifier != 0 {
		t.Errorf("the first batch took a modifier it was not granted")
	}

	if gen.char.State.Benefits[1].Modifier != 1 {
		t.Errorf("the second batch = %+d, want +1", gen.char.State.Benefits[1].Modifier)
	}

	gen.grantBenefits(career.Effect{Modifier: 1, Scope: career.ScopeCareer, Detail: "career-wide"}, 1)

	if gen.char.State.Benefits[0].Modifier != 1 {
		t.Errorf("a career-wide modifier did not reach a batch already granted")
	}

	if gen.char.State.Benefits[1].Modifier != 2 {
		t.Errorf("a career-wide modifier did not stack on a batch that had one")
	}
}

func TestRollExpression(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expr      string
		low, high int
		bad       bool
	}{
		{expr: "100", low: 100, high: 100},
		{expr: "5000", low: 5000, high: 5000},
		{expr: "1d6", low: 1, high: 6},
		{expr: "2d6", low: 2, high: 12},
		{expr: "1d3", low: 1, high: 3},
		{expr: "1d6x100", low: 100, high: 600},
		{expr: "2d6x10000", low: 20000, high: 120000},
		{expr: "1d4", bad: true},
		{expr: "d6", bad: true},
		{expr: "1d6x", bad: true},
		{expr: "", bad: true},
		{expr: "one", bad: true},
		{expr: "1d6xtwo", bad: true},
	}

	for _, tc := range tests {
		t.Run(tc.expr, func(t *testing.T) {
			t.Parallel()

			gen := engine(t, 2)

			got, err := gen.rollExpression(tc.expr, "p. 1")
			if tc.bad {
				if !errors.Is(err, ErrBadExpression) {
					t.Errorf("err = %v, want ErrBadExpression", err)
				}

				return
			}

			if err != nil {
				t.Fatalf("err = %v", err)
			}

			if got < tc.low || got > tc.high {
				t.Errorf("= %d, outside [%d,%d]", got, tc.low, tc.high)
			}
		})
	}
}

// TestPromoteStopsAtTheHighestPrintedRank: the rank tables run 0 to 6, and
// a character who kept advancing past the last printed row would be reading
// off the end of the table.
func TestPromoteStopsAtTheHighestPrintedRank(t *testing.T) {
	t.Parallel()

	gen := engine(t, 3)
	colonist := career.Colonist()

	gen.career = &colonist

	settler, ok := colonist.Assignment("Settler")
	if !ok {
		t.Fatal("no Settler assignment")
	}

	gen.rank = len(settler.Ranks) - 1

	err := gen.promote(1, settler)
	if err != nil {
		t.Fatalf("promote: %v", err)
	}

	if gen.rank != len(settler.Ranks)-1 {
		t.Errorf("rank rose to %d past the last printed row", gen.rank)
	}
}

func TestApplyRejectsAnUnknownEffectKind(t *testing.T) {
	t.Parallel()

	gen := engine(t, 4)

	err := gen.apply(career.Effect{Kind: career.EffectKind(999)}, 1)
	if !errors.Is(err, ErrUnknownEffect) {
		t.Errorf("err = %v, want ErrUnknownEffect", err)
	}
}

func TestApplyRejectsACheckWithNoTarget(t *testing.T) {
	t.Parallel()

	gen := engine(t, 4)

	err := gen.apply(career.Effect{Kind: career.EffectCheck, Detail: "malformed"}, 1)
	if !errors.Is(err, ErrMalformedCheck) {
		t.Errorf("err = %v, want ErrMalformedCheck", err)
	}
}

func TestServeTermNeedsACareer(t *testing.T) {
	t.Parallel()

	gen := engine(t, 5)

	err := gen.serveTerm()
	if !errors.Is(err, ErrNoCareer) {
		t.Errorf("err = %v, want ErrNoCareer", err)
	}
}

func TestRollMishapNeedsACareer(t *testing.T) {
	t.Parallel()

	gen := engine(t, 5)

	err := gen.rollMishap(1, true)
	if !errors.Is(err, ErrNoCareer) {
		t.Errorf("err = %v, want ErrNoCareer", err)
	}
}

func TestItoa(t *testing.T) {
	t.Parallel()

	for value, want := range map[int]string{0: "0", 7: "7", 12: "12", 100: "100", -3: "-3", -47: "-47"} {
		if got := itoa(value); got != want {
			t.Errorf("itoa(%d) = %q, want %q", value, got, want)
		}
	}
}

func TestTermLimit(t *testing.T) {
	t.Parallel()

	for requested, want := range map[int]int{0: defaultTermLimit, 6: 6, 1: 1, -1: 0, -9: 0} {
		if got := termLimit(requested); got != want {
			t.Errorf("termLimit(%d) = %d, want %d", requested, got, want)
		}
	}
}

func TestSkillTableKindNames(t *testing.T) {
	t.Parallel()

	want := map[career.SkillTableKind]string{
		career.PersonalDevelopment: "Personal Development",
		career.ServiceSkills:       "Service Skills",
		career.AdvancedEducation:   "Advanced Education",
		career.OfficerSkills:       "Officer Skills",
		career.AssignmentSkills:    "Assignment",
	}

	for kind, name := range want {
		if got := kind.String(); got != name {
			t.Errorf("String() = %q, want %q", got, name)
		}
	}

	if got := career.SkillTableKind(99).String(); got != "unknown" {
		t.Errorf("an unnamed kind stringed as %q", got)
	}
}

// settingWith builds a one-subsector file for a test to reach a branch a
// generated character rarely takes.
func settingWith(sub setting.Subsector) *setting.Data {
	return &setting.Data{
		SchemaVersion: setting.SchemaVersion,
		Name:          testVersion,
		Hash:          testVersion,
		Subsectors:    []setting.Subsector{sub},
	}
}

func testWorld(name string, roll *[2]int) setting.World {
	return setting.World{
		Name: name, Roll: roll, TechLevel: 10, MaximumAge: 100, MaximumTerms: 20,
		PrimaryLanguages: []string{testLanguage},
		Engineered:       setting.Permission{Allowed: true, Status: setting.Free},
		Uplifts:          setting.Permission{Allowed: true, Status: setting.Free},
	}
}

// TestAChooseOnlySubsectorIsChosenFrom is the Recently Colonized Worlds
// case: "you cannot randomly be assigned one of these worlds, you may
// choose them" (p. 40). A subsector whose worlds all lack a d100 range
// becomes a choice rather than an error.
func TestAChooseOnlySubsectorIsChosenFrom(t *testing.T) {
	t.Parallel()

	data := settingWith(setting.Subsector{
		Name: "New Holdings", OriginRoll: 0,
		Worlds: []setting.World{testWorld("Tinderfall", nil), testWorld("Wake", nil)},
	})

	character, err := New(Options{
		Seed: 3, Decider: Policy{}, Setting: data,
		Inputs: Inputs{Species: testHuman, TermLimit: 1},
	}).Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(character.State.Homeworlds) == 0 {
		t.Fatal("no homeworld")
	}

	// The policy takes the first option at every choice point.
	if got := character.State.Homeworlds[0].World; got != "Tinderfall" {
		t.Errorf("homeworld = %q, want the first world offered", got)
	}
}

// TestASubsectorRollThatLandsNowhereBecomesAChoice: the book's chart has a
// subsector for every 1d6 result, but a setting may have fewer than six, and
// a throw into the gap must not leave the character nowhere.
func TestASubsectorRollThatLandsNowhereBecomesAChoice(t *testing.T) {
	t.Parallel()

	rolls := [2]int{1, 100}
	data := settingWith(setting.Subsector{
		Name: testSubsector, OriginRoll: 6, // a 1d6 rarely lands here
		Worlds: []setting.World{testWorld("Somewhere", &rolls)},
	})

	for seed := range uint64(12) {
		character, err := New(Options{
			Seed: seed, Decider: Policy{}, Setting: data,
			Inputs: Inputs{Species: testHuman, TermLimit: 1},
		}).Run()
		if err != nil {
			t.Fatalf("seed %d: %v", seed, err)
		}

		if len(character.State.Homeworlds) == 0 {
			t.Fatalf("seed %d: no homeworld", seed)
		}
	}
}

// TestALanguageWithNoAlternativeIsRecorded: p. 41 requires the Language
// specialty to differ from the primary language, and a homeworld that
// offers only the primary leaves nothing to grant.
func TestALanguageWithNoAlternativeIsRecorded(t *testing.T) {
	t.Parallel()

	rolls := [2]int{1, 100}
	world := testWorld("Monoglot", &rolls)

	world.PrimaryLanguages = []string{testLanguage}
	world.BackgroundSkills = []setting.Requirement{{
		OneOf: []setting.Alternative{{Skill: "Language", Specialties: []string{testLanguage}}},
	}}

	data := settingWith(setting.Subsector{
		Name: testSubsector, OriginRoll: 1, Worlds: []setting.World{world},
	})

	character, err := New(Options{
		Seed: 1, Decider: Policy{}, Setting: data,
		Inputs: Inputs{Species: testHuman, TermLimit: -1},
	}).Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	said := false

	for _, event := range character.Events {
		if event.Kind == EventConsequence &&
			event.Consequence.Kind == ConsequenceUnimplemented &&
			strings.Contains(event.Consequence.Detail, "no specialty other than the primary") {
			said = true
		}
	}

	if !said {
		t.Error("the record says nothing about a Language with no alternative specialty")
	}
}

// TestAHomeworldItemGoesToTheStash: Kingston grants a mindcomp alongside
// its skills (p. 41), and an item is not a skill.
func TestAHomeworldItemGoesToTheStash(t *testing.T) {
	t.Parallel()

	rolls := [2]int{1, 100}
	world := testWorld("Workshop", &rolls)

	world.BackgroundSkills = []setting.Requirement{{
		OneOf: []setting.Alternative{{Item: "a neural companion"}},
	}}

	data := settingWith(setting.Subsector{
		Name: testSubsector, OriginRoll: 1, Worlds: []setting.World{world},
	})

	character, err := New(Options{
		Seed: 1, Decider: Policy{}, Setting: data,
		Inputs: Inputs{Species: testHuman, TermLimit: -1},
	}).Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !slices.Contains(character.State.Stash, "a neural companion") {
		t.Errorf("the item is not in the stash: %v", character.State.Stash)
	}
}

// TestARollOnAnotherCareersTable is Colonist event 54: "Roll twice on the
// Assignment: Ambassador table of the Diplomatic Service career" (p. 176).
// It is the one result in the book that reaches into a career the character
// is not in, and no seed reaches it often enough to rely on.
func TestARollOnAnotherCareersTable(t *testing.T) {
	t.Parallel()

	gen := engine(t, 5)
	colonist := career.Colonist()

	gen.career = &colonist
	gen.assignment = colonist.Assignments[0]

	err := gen.apply(career.Colonist().Events[54].Effects[0], 1)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	// The Ambassador table is Etiquette, Carouse, Diplomat, Language,
	// Advocate, Persuade (p. 185). Whichever row came up, the character now
	// holds one of them, and none is on any Colonist table.
	ambassador := []string{"Etiquette", "Carouse", "Diplomat", "Language", "Advocate", "Persuade"}

	held := false

	for _, name := range ambassador {
		if gen.char.State.Has(name) {
			held = true
		}
	}

	if !held {
		t.Errorf("no Ambassador skill was granted; the character holds %v", gen.char.State.Skills)
	}
}

// TestARollOnACareerThatIsNotBuiltIsRecorded: the same effect naming a
// career this repository has not transcribed must say so rather than
// silently granting nothing.
func TestARollOnACareerThatIsNotBuiltIsRecorded(t *testing.T) {
	t.Parallel()

	gen := engine(t, 5)
	colonist := career.Colonist()

	gen.career = &colonist
	gen.assignment = colonist.Assignments[0]

	err := gen.apply(career.Effect{
		Kind:       career.EffectRollTable,
		Detail:     "roll on the Ambassador table of the Nowhere career",
		Table:      career.AssignmentSkills,
		Career:     "Nowhere",
		Assignment: "Ambassador",
	}, 1)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	last := gen.log.Events()[gen.log.Len()-1]
	if last.Consequence == nil || last.Consequence.Kind != ConsequenceUnimplemented {
		t.Error("a table in a career that does not exist was not recorded as unimplemented")
	}
}
