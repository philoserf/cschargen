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
// mustGrant applies a benefit-roll effect and fails the test if it errors,
// which it can only do on a dice expression the tables here do not use.
func mustGrant(t *testing.T, gen *Generator, effect career.Effect) {
	t.Helper()

	err := gen.grantBenefits(effect, 1)
	if err != nil {
		t.Fatalf("grantBenefits(%q): %v", effect.Detail, err)
	}
}

func TestBenefitScopesAreNotTheSameThing(t *testing.T) {
	t.Parallel()

	gen := engine(t, 1)
	colonist := career.Colonist()

	gen.service.career = &colonist

	mustGrant(t, gen, career.Effect{Count: 2, Detail: "two rolls"})
	mustGrant(t, gen, career.Effect{Count: 3, Modifier: 1, Detail: "three at +1"})

	if len(gen.char.State.Benefits) != 2 {
		t.Fatalf("%d batches, want 2", len(gen.char.State.Benefits))
	}

	if gen.char.State.Benefits[0].Modifier != 0 {
		t.Errorf("the first batch took a modifier it was not granted")
	}

	if gen.char.State.Benefits[1].Modifier != 1 {
		t.Errorf("the second batch = %+d, want +1", gen.char.State.Benefits[1].Modifier)
	}

	mustGrant(t, gen, career.Effect{Modifier: 1, Scope: career.ScopeCareer, Detail: "career-wide"})

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
		{expr: "1d3+1", low: 2, high: 4},
		{expr: "1d10", low: 1, high: 10},
		{expr: "2d10", low: 2, high: 20},
		{expr: "3d10", low: 3, high: 30},
		{expr: "1d10x15", low: 15, high: 150},
		{expr: "1d6+2", low: 3, high: 8},
		{expr: "10+5", low: 15, high: 15},
		{expr: "1d6+", bad: true},
		{expr: "1d6+one", bad: true},
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

	gen.service.career = &colonist

	settler, ok := colonist.Assignment("Settler")
	if !ok {
		t.Fatal("no Settler assignment")
	}

	gen.service.rank = len(settler.Ranks) - 1

	err := gen.promote(1, settler)
	if err != nil {
		t.Fatalf("promote: %v", err)
	}

	if gen.service.rank != len(settler.Ranks)-1 {
		t.Errorf("rank rose to %d past the last printed row", gen.service.rank)
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

	err := gen.rollMishap(true)
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

func testWorld(name string, roll []int) setting.World {
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

// TestASubsectorRollThatLandsNowhereIsThrownAgain: the book's chart has a
// subsector for every 1d6 result, but a setting may have fewer than six, and
// a throw into the gap must not leave the character nowhere.
//
// The file here is the narrowest one a validator will pass: a single
// subsector, claiming a single result. Five throws in six miss it, so every
// character generated below re-throws, most of them more than once.
func TestASubsectorRollThatLandsNowhereIsThrownAgain(t *testing.T) {
	t.Parallel()

	rolls := []int{1, 100}
	data := settingWith(setting.Subsector{
		Name: testSubsector, OriginRoll: 6, // a 1d6 rarely lands here
		Worlds: []setting.World{testWorld("Somewhere", rolls)},
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

	rolls := []int{1, 100}
	world := testWorld("Monoglot", rolls)

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

// TestAHomeworldItemGoesToTheStash: a world may grant an item alongside its
// background skills (p. 41), and an item is not a skill.
func TestAHomeworldItemGoesToTheStash(t *testing.T) {
	t.Parallel()

	rolls := []int{1, 100}
	world := testWorld("Workshop", rolls)

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

	if !slices.ContainsFunc(character.State.Stash, func(p Possession) bool {
		return p.Item == "a neural companion"
	}) {
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

	gen.service.career = &colonist
	gen.service.assignment = colonist.Assignments[0]

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

	gen.service.career = &colonist
	gen.service.assignment = colonist.Assignments[0]

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

// TestEveryTranscribedExpressionParses walks the whole career corpus and
// throws every Dice string the tables carry.
//
// The expressions are written by hand from the page, and the page writes
// them as prose: "1d6 x ₶10,000", "100,000 credits". [Generator.rollExpression]
// takes none of that -- no spaces, no commas, no currency -- and a mistyped
// one fails at generation time on whichever seed first reaches that row,
// which may be no seed any test runs. Walking the corpus turns that into a
// compile-time-ish check: every expression is thrown once, here, whether or
// not a character ever rolls it.
func TestEveryTranscribedExpressionParses(t *testing.T) {
	t.Parallel()

	gen := engine(t, 2)

	throw := func(where string, effects []career.Effect) {
		for _, effect := range expressionsIn(effects) {
			_, err := gen.rollExpression(effect.Dice, "p. 1")
			if err != nil {
				t.Errorf("%s: %q does not parse: %v", where, effect.Dice, err)
			}
		}
	}

	for _, def := range career.All() {
		for result, row := range def.Events {
			throw(def.Name+" event "+itoa(result), row.Effects)
		}

		for index, row := range def.Mishaps {
			throw(def.Name+" mishap "+itoa(index+2), row.Effects)
		}

		for _, benefit := range def.Benefits {
			throw(def.Name+" benefit", []career.Effect{benefit.Other})
		}

		for _, table := range def.Tables {
			throw(def.Name+" "+table.Name, table.Rows[:])
		}

		throw(def.Name+" ranks", ranksAndSkillsOf(def))
	}

	for _, row := range career.LifeEvents() {
		throw("life event", row.Effects)
	}
}

// ranksAndSkillsOf flattens an assignment's skill table and both rank
// tracks, which nest one level deeper than everything else a career holds.
func ranksAndSkillsOf(def career.Career) []career.Effect {
	var flat []career.Effect

	for _, assignment := range def.Assignments {
		flat = append(flat, assignment.Skills.Rows[:]...)

		for _, rank := range assignment.Ranks {
			flat = append(flat, rank...)
		}

		for _, rank := range assignment.OfficerRanks {
			flat = append(flat, rank...)
		}
	}

	return flat
}

// expressionsIn is the recursion, separated from the check so that neither
// has to be read through the other.
func expressionsIn(effects []career.Effect) []career.Effect {
	var found []career.Effect

	for _, effect := range effects {
		if effect.Dice != "" {
			found = append(found, effect)
		}

		found = append(found, expressionsIn(effect.Success)...)
		found = append(found, expressionsIn(effect.Failure)...)

		for _, option := range effect.Options {
			found = append(found, expressionsIn(option.Effects)...)
		}
	}

	return found
}

// TestTheAgingTablesMatchPages122And123 is a second reading of the four
// human tables, typed independently of chargen/aging.go.
//
// The repetition is the mechanism, as it is in career/transcription_test.go:
// two transcriptions that agree are evidence, and a constant shared between
// them would be one reading wearing two hats. So this table is written as
// the page prints it -- every characteristic on every row -- rather than
// through threePhysical and its two siblings.
func TestTheAgingTablesMatchPages122And123(t *testing.T) {
	t.Parallel()

	type row struct {
		techLevel int
		term      int
		checks    string
	}

	// One string per row of pp. 122-123, read left to right.
	rows := []row{
		{9, 5, ""},
		{9, 6, "STR 8+, DEX 8+, END 8+"},
		{9, 8, "STR 8+, DEX 8+, END 8+"},
		{9, 9, "STR 9+, DEX 9+, END 9+"},
		{9, 10, "STR 9+, DEX 9+, END 9+"},
		{9, 11, "STR 9+, DEX 9+, END 9+, INT 9+, CHA 9+"},
		{9, 12, "STR 10+, DEX 10+, END 10+, INT 10+, EDU 10+, CHA 9+"},
		{9, 40, "STR 10+, DEX 10+, END 10+, INT 10+, EDU 10+, CHA 9+"},
		{10, 17, ""},
		{10, 18, "STR 8+, DEX 8+, END 8+"},
		{10, 28, "STR 8+, DEX 8+, END 8+"},
		{10, 29, "STR 9+, DEX 9+, END 9+"},
		{10, 38, "STR 9+, DEX 9+, END 9+"},
		{10, 39, "STR 9+, DEX 9+, END 9+, INT 9+, CHA 9+"},
		{10, 48, "STR 9+, DEX 9+, END 9+, INT 9+, CHA 9+"},
		{10, 49, "STR 10+, DEX 10+, END 10+, INT 10+, EDU 10+, CHA 9+"},
		{11, 32, ""},
		{11, 33, "STR 8+, DEX 8+, END 8+"},
		{11, 48, "STR 8+, DEX 8+, END 8+"},
		{11, 49, "STR 9+, DEX 9+, END 9+"},
		{11, 58, "STR 9+, DEX 9+, END 9+"},
		// ERRATA E-15: TL 11 and TL 12-13 print no band past 58, where
		// TL 9 and TL 10 each end on an open one.
		{11, 59, ""},
		{12, 43, ""},
		{12, 44, "STR 8+, DEX 8+, END 8+"},
		{12, 58, "STR 8+, DEX 8+, END 8+"},
		{12, 59, ""},
		{13, 44, "STR 8+, DEX 8+, END 8+"},
		// ERRATA E-16: nothing is printed above TL 13, and the validator
		// accepts up to TL 20.
		{14, 44, "STR 8+, DEX 8+, END 8+"},
		{20, 59, ""},
	}

	for _, want := range rows {
		checks, due := agingChecksAt(ProfileTechLevel, want.techLevel, want.term)
		got := ""

		if due {
			parts := make([]string, len(checks))
			for i, check := range checks {
				parts[i] = check.Characteristic + " " + itoa(check.Number) + "+"
			}

			got = strings.Join(parts, ", ")
		}

		if got != want.checks {
			t.Errorf("TL %d term %d: %q, want %q",
				want.techLevel, want.term, got, want.checks)
		}
	}
}

// TestNoAgingBandLeavesAGapOrOverlaps holds the property the rows above
// only sample: within one tech level the bands are contiguous and ordered,
// so no term falls between two of them and no term is covered twice.
func TestNoAgingBandLeavesAGapOrOverlaps(t *testing.T) {
	t.Parallel()

	for _, techLevel := range []int{9, 10, 11, 12, 13, 14} {
		bands := humanAging(techLevel)

		for i, band := range bands {
			if band.Through != 0 && band.Through < band.From {
				t.Errorf("TL %d band %d: %d-%d runs backwards",
					techLevel, i, band.From, band.Through)
			}

			if i == 0 {
				continue
			}

			previous := bands[i-1]
			if previous.Through == 0 {
				t.Errorf("TL %d band %d is open and is not the last", techLevel, i-1)

				continue
			}

			if band.From != previous.Through+1 {
				t.Errorf("TL %d: band %d ends at %d and band %d starts at %d",
					techLevel, i-1, previous.Through, i, band.From)
			}
		}
	}
}

// The Aging Crisis and the four states beneath it (pp. 123-124) are reached
// through a term of aging throws that has to fail, so the tests below set
// the characteristics directly and call the crisis path. A generated
// character can reach them -- TestACharacterActuallyAges proves the throws
// happen -- but not reliably enough for an assertion.

func TestAnAgingCrisisIsPaidForWhenItCanBe(t *testing.T) {
	t.Parallel()

	gen := engine(t, 4)

	gen.char.State.Credits = 100_000
	gen.char.State.Characteristics.Set(STR, 0)

	err := gen.agingCrisis("STR", 0)
	if err != nil {
		t.Fatalf("agingCrisis: %v", err)
	}

	if gen.char.State.Fate != "" {
		t.Errorf("Fate = %q, want empty: the treatment was affordable", gen.char.State.Fate)
	}

	if got := gen.char.State.Characteristics.Get(STR); got != 1 {
		t.Errorf("STR = %d, want 1 restored", got)
	}

	if gen.char.State.Credits >= 100_000 {
		t.Error("the treatment was free")
	}

	if !gen.crisisSurvived {
		t.Error("surviving a crisis did not restrict what comes next (p. 123)")
	}
}

// TestAnAgingCrisisNobodyCanPayForIsFatal is ERRATA E-19: the book sets a
// price and never says what happens to a character who cannot meet it.
func TestAnAgingCrisisNobodyCanPayForIsFatal(t *testing.T) {
	t.Parallel()

	gen := engine(t, 5)

	gen.char.State.Credits = 0
	gen.char.State.Characteristics.Set(END, 0)

	err := gen.agingCrisis("END", 0)
	if err != nil {
		t.Fatalf("agingCrisis: %v", err)
	}

	if gen.char.State.Fate != FateDied {
		t.Errorf("Fate = %q, want %q", gen.char.State.Fate, FateDied)
	}

	if !gen.stopped {
		t.Error("a dead character is still generating")
	}
}

// TestACharacteristicAboveZeroIsNotACrisis: agingCrisis is called after
// every failed check, and all but a handful of them leave the character
// merely older.
func TestACharacteristicAboveZeroIsNotACrisis(t *testing.T) {
	t.Parallel()

	gen := engine(t, 6)

	gen.char.State.Credits = 100_000
	gen.char.State.Characteristics.Set(DEX, 3)

	err := gen.agingCrisis("DEX", 0)
	if err != nil {
		t.Fatalf("agingCrisis: %v", err)
	}

	if gen.char.State.Credits != 100_000 {
		t.Error("treatment was bought for a characteristic that did not need it")
	}

	if gen.crisisSurvived {
		t.Error("a crisis was recorded where none happened")
	}
}

// terminalCase is one row of pp. 123-124's four states, plus the one that
// is survivable.
type terminalCase struct {
	name         string
	zero         []Characteristic
	wantFate     Fate
	wantNoEnlist bool
	wantStopped  bool
	wantVagabond bool
}

func terminalCases() []terminalCase {
	return []terminalCase{
		{
			name:     "all three physical is death",
			zero:     []Characteristic{STR, DEX, END},
			wantFate: FateDied, wantStopped: true,
		},
		{
			name:     "two physical is incapacity",
			zero:     []Characteristic{STR, DEX},
			wantFate: FateIncapacitated, wantStopped: true,
		},
		{
			name:     "two mental is death",
			zero:     []Characteristic{INT, EDU},
			wantFate: FateDied, wantStopped: true,
		},
		{
			name:         "one mental ends enlistment and nothing else",
			zero:         []Characteristic{CHA},
			wantNoEnlist: true, wantVagabond: true,
		},
		{
			name: "one physical is survivable",
			zero: []Characteristic{END},
		},
	}
}

func TestTheTerminalStatesOfPages123And124(t *testing.T) {
	t.Parallel()

	for _, tc := range terminalCases() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gen := engine(t, 7)

			for _, which := range CharacteristicOrder {
				gen.char.State.Characteristics.Set(which, 7)
			}

			for _, which := range tc.zero {
				gen.char.State.Characteristics.Set(which, 0)
			}

			gen.settleTerminalStates(0)

			if gen.char.State.Fate != tc.wantFate {
				t.Errorf("Fate = %q, want %q", gen.char.State.Fate, tc.wantFate)
			}

			if gen.stopped != tc.wantStopped {
				t.Errorf("stopped = %v, want %v", gen.stopped, tc.wantStopped)
			}

			if gen.mentalDecline != tc.wantNoEnlist {
				t.Errorf("mentalDecline = %v, want %v", gen.mentalDecline, tc.wantNoEnlist)
			}

			if !tc.wantVagabond {
				return
			}

			eligible := gen.eligibleCareers()
			if len(eligible) != 1 || eligible[0].Name != "Vagabond" {
				t.Errorf("after a mental characteristic reached 0 the career list is %d long, "+
					"want Vagabond alone", len(eligible))
			}
		})
	}
}

// TestAgingStopsAtTheFirstDeath. A term's checks run one after another, and
// a character killed by the first of them does not make the rest.
func TestAgingStopsAtTheFirstDeath(t *testing.T) {
	t.Parallel()

	gen := engine(t, 8)

	gen.techLevel = 9
	gen.char.State.Credits = 0

	// Term 12 on a tech level 9 world is the last band: six checks, the
	// first of which is STR. At 1 it cannot survive a failure, and the
	// modifier for a characteristic of 1 is -2, so 10+ is unreachable.
	gen.char.State.Terms = make([]Term, 12)

	for _, which := range CharacteristicOrder {
		gen.char.State.Characteristics.Set(which, 1)
	}

	err := gen.agingThrows(0)
	if err != nil {
		t.Fatalf("agingThrows: %v", err)
	}

	if gen.char.State.Fate != FateDied {
		t.Fatalf("Fate = %q, want %q", gen.char.State.Fate, FateDied)
	}

	// Only the first check was made, so only one characteristic moved.
	moved := 0

	for _, which := range CharacteristicOrder {
		if gen.char.State.Characteristics.Get(which) != 1 {
			moved++
		}
	}

	if moved != 1 {
		t.Errorf("%d characteristics moved; the term should have stopped at the first death", moved)
	}
}

// TestAnAbandonedCrisisEndsGeneration: a Decider that declines the payment
// choice -- an interactive session the player walked away from -- is an
// error rather than a decision, and it comes back out.
func TestAnAbandonedCrisisEndsGeneration(t *testing.T) {
	t.Parallel()

	gen := engine(t, 9)

	gen.decider = refusingDecider{}
	gen.char.State.Characteristics.Set(INT, 0)

	err := gen.agingCrisis("INT", 0)
	if err == nil {
		t.Fatal("a refused choice did not come back as an error")
	}
}

// refusingDecider declines every choice, which is what an abandoned
// interactive session looks like to the engine.
type refusingDecider struct{}

func (refusingDecider) Choose(Choice) (int, error) { return 0, ErrNoSetting }

func (refusingDecider) Kind() DeciderKind { return DeciderPolicy }

// lastOptionDecider takes the last option offered, which is the opposite of
// what Policy does and the only way to reach a branch the policy declines.
type lastOptionDecider struct{}

func (lastOptionDecider) Choose(c Choice) (int, error) { return len(c.Options) - 1, nil }

func (lastOptionDecider) Kind() DeciderKind { return DeciderPolicy }

// TestACrisisDeclinedIsFatal. "The player may pay" (p. 123) is a may, and
// the character who does not dies. POLICY.md takes the payment, so this is
// the branch the auto policy never reaches.
func TestACrisisDeclinedIsFatal(t *testing.T) {
	t.Parallel()

	gen := engine(t, 10)

	gen.decider = lastOptionDecider{}
	gen.char.State.Credits = 100_000
	gen.char.State.Characteristics.Set(DEX, 0)

	err := gen.agingCrisis("DEX", 0)
	if err != nil {
		t.Fatalf("agingCrisis: %v", err)
	}

	if gen.char.State.Fate != FateDied {
		t.Errorf("Fate = %q, want %q: the treatment was refused", gen.char.State.Fate, FateDied)
	}

	if gen.char.State.Credits != 100_000 {
		t.Error("a refused treatment was paid for anyway")
	}
}

// TestARefusedCrisisStopsTheTerm: a refused choice is an error, and it has
// to come back out of the term rather than being swallowed by the loop over
// a band's checks.
func TestARefusedCrisisStopsTheTerm(t *testing.T) {
	t.Parallel()

	gen := engine(t, 11)

	gen.decider = refusingDecider{}
	gen.techLevel = 9
	gen.char.State.Terms = make([]Term, 12)

	for _, which := range CharacteristicOrder {
		gen.char.State.Characteristics.Set(which, 1)
	}

	err := gen.agingThrows(0)
	if err == nil {
		t.Fatal("a refused choice inside a term's aging did not come back as an error")
	}
}

// TestTheApparentAgeChartMatchesPage125 is a second reading of the chart,
// typed as the page prints it: the actual-age band on the left, then the
// three tech level columns.
func TestTheApparentAgeChartMatchesPage125(t *testing.T) {
	t.Parallel()

	// "Actual Age | TL 10 | TL 11 | TL 12-13", row by row down p. 125.
	page := []struct {
		age                       int
		tenTL, elevenTL, twelveTL string
	}{
		{30, "20-25", "20-25", "20-25"},
		{41, "20-25", "20-25", "20-25"},
		{51, "30-35", "20-25", "20-25"},
		{61, "30-35", "25-30", "20-25"},
		{71, "35-40", "25-30", "25-30"},
		{81, "35-40", "25-30", "25-30"},
		{91, "40-45", "30-35", "25-30"},
		{101, "40-45", "30-35", "25-30"},
		{111, "45-50", "30-35", "30-35"},
		{121, "45-50", "35-40", "30-35"},
		{131, "50-55", "35-40", "30-35"},
		{141, "50-55", "35-40", "30-35"},
		{151, "55-60", "40-45", "35-40"},
		{161, "55-60", "40-45", "35-40"},
		{171, "60-65", "40-45", "35-40"},
		{181, "60-65", "45-50", "35-40"},
		{191, "65-70", "45-50", "40-45"},
		{201, "65-70", "45-50", "40-45"},
		{211, "70-75", "50-55", "40-45"},
		{221, "70-75", "50-55", "40-45"},
		{231, "75-80", "50-55", "45-50"},
		{241, "75-80", "55-60", "45-50"},
		{251, "80-85", "55-60", "45-50"},
		{261, "80-85", "55-60", "45-50"},
		{271, "85-90", "60-65", "50-55"},
		{281, "85-90", "60-65", "50-55"},
	}

	if len(page) != len(apparentAgeChart) {
		t.Fatalf("the reading above has %d rows and the chart has %d",
			len(page), len(apparentAgeChart))
	}

	for _, want := range page {
		for techLevel, column := range map[int]string{
			10: want.tenTL, 11: want.elevenTL, 12: want.twelveTL,
		} {
			band, fromChart := apparentAge(techLevel, want.age)
			if !fromChart {
				t.Errorf("age %d at TL %d did not read the chart", want.age, techLevel)

				continue
			}

			if band.String() != column {
				t.Errorf("age %d at TL %d: %s, want %s",
					want.age, techLevel, band, column)
			}
		}
	}
}

// TestApparentAgeOutsideTheChart is ERRATA E-20: below tech level 10 and
// below age 30 apparent age is actual age, and above 290 the last printed
// row holds.
func TestApparentAgeOutsideTheChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		techLevel, age int
		want           string
		fromChart      bool
	}{
		{9, 200, "200-200", false},
		{0, 18, "18-18", false},
		{12, 29, "29-29", false},
		{12, 30, "20-25", true},
		{10, 290, "85-90", true},
		{10, 291, "85-90", true},
		{10, 10_000, "85-90", true},
		// ERRATA E-16 again: nothing is printed above TL 13, and a higher
		// homeworld reads the TL 12-13 column.
		{20, 151, "35-40", true},
	}

	for _, tc := range tests {
		band, fromChart := apparentAge(tc.techLevel, tc.age)
		if band.String() != tc.want || fromChart != tc.fromChart {
			t.Errorf("apparentAge(%d, %d) = %s, %v; want %s, %v",
				tc.techLevel, tc.age, band, fromChart, tc.want, tc.fromChart)
		}
	}
}

// TestOverFortyIsABandNotANumber is ERRATA E-21. Twelve careers say "if you
// have an apparent age of over 40", and the chart answers in bands.
func TestOverFortyIsABandNotANumber(t *testing.T) {
	t.Parallel()

	tests := []struct {
		band AgeBand
		want bool
	}{
		{AgeBand{35, 40}, false},
		{AgeBand{40, 45}, true},
		{AgeBand{45, 50}, true},
		{AgeBand{20, 25}, false},
		// A character below the chart carries their own age as a band.
		{AgeBand{39, 39}, false},
		{AgeBand{40, 40}, true},
	}

	for _, tc := range tests {
		if got := apparentAgeOverForty(tc.band); got != tc.want {
			t.Errorf("apparentAgeOverForty(%s) = %v, want %v", tc.band, got, tc.want)
		}
	}
}

// TestRollExpressionSubtracts is what a species' characteristic method
// needs: "roll 2d6-2 for STR and END" (p. 23).
func TestRollExpressionSubtracts(t *testing.T) {
	t.Parallel()

	gen := engine(t, 35)

	tests := []struct {
		expr      string
		low, high int
		bad       bool
	}{
		{expr: "2d6-2", low: 0, high: 10},
		{expr: "2d6+2", low: 4, high: 14},
		{expr: "1d6-1", low: 0, high: 5},
		{expr: "10-4", low: 6, high: 6},
		{expr: "2d6-", bad: true},
		{expr: "2d6-two", bad: true},
	}

	for _, tc := range tests {
		got, err := gen.rollExpression(tc.expr, "p. 23")
		if tc.bad {
			if !errors.Is(err, ErrBadExpression) {
				t.Errorf("%q: err = %v, want ErrBadExpression", tc.expr, err)
			}

			continue
		}

		if err != nil {
			t.Errorf("%q: %v", tc.expr, err)

			continue
		}

		if got < tc.low || got > tc.high {
			t.Errorf("%q = %d, outside [%d,%d]", tc.expr, got, tc.low, tc.high)
		}
	}
}

// TestTheFiveAgingProfiles is a second reading of pp. 122-123's tables,
// once their headings are set aside. The headings are species names the OGL
// notice reserves, so the profiles are named for what they do.
func TestTheFiveAgingProfiles(t *testing.T) {
	t.Parallel()

	type row struct {
		profile AgingProfile
		term    int
		checks  string
	}

	rows := []row{
		// "4-5 / 6-7 / 8+", the only table that opens on five checks and
		// reaches 11+.
		{ProfileRapid, 3, ""},
		{ProfileRapid, 4, "STR 9+, DEX 9+, END 9+, INT 9+, CHA 9+"},
		{ProfileRapid, 5, "STR 9+, DEX 9+, END 9+, INT 9+, CHA 9+"},
		{ProfileRapid, 6, "STR 10+, DEX 10+, END 10+, INT 10+, EDU 10+, CHA 9+"},
		{ProfileRapid, 7, "STR 10+, DEX 10+, END 10+, INT 10+, EDU 10+, CHA 9+"},
		{ProfileRapid, 8, "STR 11+, DEX 11+, END 11+, INT 11+, EDU 11+, CHA 9+"},
		{ProfileRapid, 40, "STR 11+, DEX 11+, END 11+, INT 11+, EDU 11+, CHA 9+"},
		// "6-7 / 8-9 / 10+".
		{ProfileModerate, 5, ""},
		{ProfileModerate, 6, "STR 9+, DEX 9+, END 9+"},
		{ProfileModerate, 7, "STR 9+, DEX 9+, END 9+"},
		{ProfileModerate, 8, "STR 9+, DEX 9+, END 9+, INT 9+, CHA 9+"},
		{ProfileModerate, 10, "STR 10+, DEX 10+, END 10+, INT 10+, EDU 10+, CHA 9+"},
		// "6-8 / 9-11 / 13+" -- and term 12 is in no printed band, which
		// is ERRATA E-30.
		{ProfileSturdy, 5, ""},
		{ProfileSturdy, 6, "STR 9+, DEX 9+, END 9+"},
		{ProfileSturdy, 8, "STR 9+, DEX 9+, END 9+"},
		{ProfileSturdy, 9, "STR 9+, DEX 9+, END 9+, INT 9+, CHA 9+"},
		{ProfileSturdy, 11, "STR 9+, DEX 9+, END 9+, INT 9+, CHA 9+"},
		{ProfileSturdy, 12, "STR 9+, DEX 9+, END 9+, INT 9+, CHA 9+"},
		{ProfileSturdy, 13, "STR 10+, DEX 10+, END 10+, INT 10+, EDU 10+, CHA 9+"},
		// "7 / 8 / 9+", each of its first two bands one term wide.
		{ProfileSudden, 6, ""},
		{ProfileSudden, 7, "STR 9+, DEX 9+, END 9+, INT 9+, CHA 9+"},
		{ProfileSudden, 8, "STR 10+, DEX 10+, END 10+, INT 10+, EDU 10+, CHA 9+"},
		{ProfileSudden, 9, "STR 11+, DEX 11+, END 11+, INT 11+, EDU 11+, CHA 9+"},
		// And the default reads the homeworld's tech level.
		{ProfileTechLevel, 6, "STR 8+, DEX 8+, END 8+"},
	}

	for _, want := range rows {
		// The tech level is only read by the profile that reads it; nine
		// is what the others are given and ignore.
		checks, due := agingChecksAt(want.profile, 9, want.term)

		got := ""
		if due {
			got = renderChecks(checks)
		}

		if got != want.checks {
			t.Errorf("%s at term %d: %q, want %q",
				want.profile, want.term, got, want.checks)
		}
	}
}

// agingReadingCases are hoisted out of the test body: the table is the
// interesting part and funlen counts it against the function otherwise.
var agingReadingCases = []struct {
	name      string
	profile   AgingProfile
	techLevel int
	term      int
	want      []string
}{
	{
		name: "a human on a printed table stamps nothing",
		// TL 10 ends on an open band, so no reading is reached.
		profile: ProfileTechLevel, techLevel: 10, term: 60,
		want: nil,
	},
	{
		name: "past the last printed term on a table that stops",
		// ERRATA E-15: TL 11 and TL 12-13 end at 58.
		profile: ProfileTechLevel, techLevel: 11, term: 59,
		want: []string{"E-15"},
	},
	{
		name: "a homeworld above the highest printed tech level",
		// ERRATA E-16: nothing is printed above TL 13.
		profile: ProfileTechLevel, techLevel: 14, term: 6,
		want: []string{"E-16"},
	},
	{
		name:    "both at once, on a high-tech world and a long life",
		profile: ProfileTechLevel, techLevel: 20, term: 59,
		want: []string{"E-16", "E-15"},
	},
	{
		name: "the term the sturdy table skips",
		// ERRATA E-30: it prints 6-8, 9-11 and 13+.
		profile: ProfileSturdy, techLevel: 12, term: sturdyGapTerm,
		want: []string{"E-30"},
	},
	{
		name:    "a sturdy character either side of the gap",
		profile: ProfileSturdy, techLevel: 12, term: 11,
		want: nil,
	},
}

// TestTheAgingTablesReadingsReachTheRecord is the three readings the printed
// aging tables rest on, each stamped where it fires.
//
// They are checked here rather than through a lifepath because two of them
// need a character the sample setting cannot make: E-30 wants a species on
// the sturdy profile, which setting/sample.json declares none of, and E-15
// wants a character past term 58, which needs a world allowing them. A
// reading nothing exercises is a reading nobody can trust.
func TestTheAgingTablesReadingsReachTheRecord(t *testing.T) {
	t.Parallel()

	for _, test := range agingReadingCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gen := &Generator{
				char:         &Character{},
				agingProfile: test.profile,
				techLevel:    test.techLevel,
			}

			gen.stampAgingReadings(test.term)

			got := gen.char.Provenance.Deviations

			if len(got) != len(test.want) {
				t.Fatalf("stamped %v, want %v", got, test.want)
			}

			for _, id := range test.want {
				if !slices.Contains(got, id) {
					t.Errorf("stamped %v, which does not include %s", got, id)
				}
			}
		})
	}
}

// renderChecks writes a band the way the page prints it.
func renderChecks(checks []AgingCheck) string {
	parts := make([]string, len(checks))
	for i, check := range checks {
		parts[i] = check.Characteristic + " " + itoa(check.Number) + "+"
	}

	return strings.Join(parts, ", ")
}

// TestEveryProfileIsContiguous holds the property the rows above sample,
// for all five: bands in order, no gap and no overlap.
func TestEveryProfileIsContiguous(t *testing.T) {
	t.Parallel()

	for _, profile := range AgingProfiles() {
		bands := agingFor(profile, 9)

		for i, band := range bands {
			if band.Through != 0 && band.Through < band.From {
				t.Errorf("%s band %d: %d-%d runs backwards",
					profile, i, band.From, band.Through)
			}

			if i == 0 {
				continue
			}

			previous := bands[i-1]
			if previous.Through == 0 {
				t.Errorf("%s band %d is open and is not the last", profile, i-1)

				continue
			}

			if band.From != previous.Through+1 {
				t.Errorf("%s: band %d ends at %d and band %d starts at %d",
					profile, i-1, previous.Through, i, band.From)
			}
		}
	}
}

// TestTheProfileNamesAgreeWithTheValidator. The names a setting file may
// use live in setting, because a validator that could not check them would
// pass a file the engine then had to fall back on. The tables live here.
// Nothing holds the two together but this.
func TestTheProfileNamesAgreeWithTheValidator(t *testing.T) {
	t.Parallel()

	if len(setting.AgingProfiles) != len(AgingProfiles()) {
		t.Fatalf("the validator accepts %d profiles and the engine holds %d",
			len(setting.AgingProfiles), len(AgingProfiles()))
	}

	for _, profile := range AgingProfiles() {
		if !setting.AgingProfiles[string(profile)] {
			t.Errorf("the engine holds %q and a file naming it would be rejected", profile)
		}
	}
}

// TestASettingWithNothingToRollForAsksInstead. Every subsector in a file may
// be choose-only: `originRoll` is optional, and the validator does not
// require any subsector to claim a result.
//
// There is then no chart to throw on, so the throw is not thrown again --
// it is never made. The question goes to whoever is deciding, which is the
// one case where a subsector is a choice rather than a fallback from one.
func TestASettingWithNothingToRollForAsksInstead(t *testing.T) {
	t.Parallel()

	rolls := []int{1, 100}
	data := settingWith(setting.Subsector{
		Name:   testSubsector, // no OriginRoll: nothing claims a d6 result
		Worlds: []setting.World{testWorld("Somewhere", rolls)},
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

	if got := character.State.Homeworlds[0].Subsector; got != testSubsector {
		t.Errorf("born in %q, want %q", got, testSubsector)
	}

	// And it was asked, not thrown for: a 1d6 on a chart nothing claims
	// would be a throw whose result means nothing.
	if !askedForASubsector(character) {
		t.Error("no subsector choice was offered")
	}

	if threwForASubsector(character) {
		t.Error("a subsector was thrown for although none claims a result")
	}
}

func askedForASubsector(character *Character) bool {
	for _, event := range character.Events {
		if event.Kind == EventChoice && event.Choice.Point == "subsector" {
			return true
		}
	}

	return false
}

// threwForASubsector finds the throw that is both a 1d6 and cited to Step
// 3's page: the homeworld's d100 shares the cite, and the careers and early
// life share the die.
func threwForASubsector(character *Character) bool {
	for _, event := range character.Events {
		if event.Kind == EventThrow && event.Throw.Expr == "1d6" && event.Throw.Cite == originCite {
			return true
		}
	}

	return false
}
