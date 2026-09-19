package chargen

import (
	"errors"
	"testing"

	"github.com/philoserf/cschargen/career"
)

// The tests here reach the engine's own helpers. They are the paths a
// generated character rarely takes -- a malformed table, an unknown
// characteristic, a rank already at its ceiling -- and each of them is a
// branch that would otherwise be reasoned about rather than run.

func engine(t *testing.T, seed uint64) *Generator {
	t.Helper()

	return New(Options{
		Seed:          seed,
		Decider:       Policy{},
		EngineVersion: "test",
		PolicyVersion: "test",
		Inputs:        Inputs{Species: "human", TermLimit: -1},
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
