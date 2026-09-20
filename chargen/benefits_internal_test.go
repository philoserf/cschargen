package chargen

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/philoserf/cschargen/career"
)

// pennilessCareer is a career whose every Cash row pays nothing, which is
// what a re-rolling result needs to be held against: the re-roll is forced
// whichever row comes up, and the answer is still nothing.
func pennilessCareer() career.Career {
	return career.Career{Name: "Penniless"}
}

// TestARolledBenefitCountIsThrown is the two results that move benefit
// rolls by a number the book rolls for -- "lose 1d6 Benefit rolls", "gain
// 1d3" -- rather than by a printed one.
func TestARolledBenefitCountIsThrown(t *testing.T) {
	t.Parallel()

	gen := engine(t, 3)
	colonist := career.Colonist()

	gen.career = &colonist

	mustGrant(t, gen, career.Effect{Count: 1, Dice: "1d3", Detail: "gain 1d3 benefit rolls"})

	if len(gen.char.State.Benefits) != 1 {
		t.Fatalf("%d batches, want 1", len(gen.char.State.Benefits))
	}

	rolls := gen.char.State.Benefits[0].Rolls
	if rolls < 1 || rolls > 3 {
		t.Errorf("1d3 benefit rolls granted %d", rolls)
	}

	mustGrant(t, gen, career.Effect{Count: -1, Dice: "1d6", Detail: "lose 1d6 benefit rolls"})

	lost := gen.char.State.Benefits[1].Rolls
	if lost > -1 || lost < -6 {
		t.Errorf("lose 1d6 benefit rolls removed %d", -lost)
	}
}

// TestACompelledCashRollIgnoresTheCap is ERRATA E-32: a batch an event
// granted as Cash rolls is spent on Cash even past p. 127's three, because
// the cap governs a choice the result has already made.
func TestACompelledCashRollIgnoresTheCap(t *testing.T) {
	t.Parallel()

	gen := engine(t, 5)

	taking, err := gen.chooseBenefitTable(career.Colonist(),
		BenefitBatch{Career: "Colonist", Rolls: 1, CashOnly: true}, cashRollsPerCareer)
	if err != nil {
		t.Fatalf("chooseBenefitTable: %v", err)
	}

	if !taking {
		t.Error("a compelled cash roll was sent to the Other table")
	}

	if !slices.Contains(gen.char.Provenance.Deviations, "E-32") {
		t.Error("a compelled roll past the cap did not stamp E-32")
	}
}

// TestAnImmediateCashRollOutsideACareer is the guard on a result that can
// only fire inside one: there is no Cash column to read, so the record says
// so rather than reading somebody else's.
func TestAnImmediateCashRollOutsideACareer(t *testing.T) {
	t.Parallel()

	gen := engine(t, 7)

	err := gen.immediateCashRolls(career.Effect{Count: 1, Detail: "two cash rolls"}, 0)
	if err != nil {
		t.Fatalf("immediateCashRolls: %v", err)
	}

	if !strings.Contains(lastDetail(t, gen), "no cash table to read") {
		t.Errorf("no record of the missing career: %q", lastDetail(t, gen))
	}
}

// TestARerolledNothingIsStillNothing covers the re-roll clause one result
// adds. A career that pays nothing on every row is the only way to force
// it without depending on which row the dice find.
func TestARerolledNothingIsStillNothing(t *testing.T) {
	t.Parallel()

	gen := engine(t, 11)
	penniless := pennilessCareer()

	gen.career = &penniless

	err := gen.immediateCashRolls(
		career.Effect{Count: 1, RerollNothing: true, Detail: "one cash roll"}, 0)
	if err != nil {
		t.Fatalf("immediateCashRolls: %v", err)
	}

	if gen.char.State.Credits != 0 {
		t.Errorf("a career that pays nothing paid %d", gen.char.State.Credits)
	}

	if !strings.Contains(lastDetail(t, gen), "0 credits") {
		t.Errorf("the record does not say nothing was paid: %q", lastDetail(t, gen))
	}
}

// TestABenefitRollCarriesADeciderRefusalOut is the path a player quitting
// mid-muster-out takes: the refusal is returned rather than swallowed, so
// the record stops where the character stopped.
func TestABenefitRollCarriesADeciderRefusalOut(t *testing.T) {
	t.Parallel()

	gen := engine(t, 13)

	gen.decider = refusingDecider{}

	_, err := gen.benefitRoll(career.Colonist(), BenefitBatch{Career: "Colonist", Rolls: 1}, 0)
	if !errors.Is(err, ErrNoSetting) {
		t.Errorf("benefitRoll swallowed the refusal: %v", err)
	}
}

// TestABadBenefitCountIsAnError holds the dice-expression path: a rolled
// count the parser cannot read stops generation rather than granting zero.
func TestABadBenefitCountIsAnError(t *testing.T) {
	t.Parallel()

	gen := engine(t, 17)
	colonist := career.Colonist()

	gen.career = &colonist

	err := gen.grantBenefits(career.Effect{Count: 1, Dice: "several", Detail: "several"}, 0)
	if err == nil {
		t.Error("an unreadable benefit count was granted anyway")
	}
}

// lastDetail is the detail of the most recent consequence in a log.
func lastDetail(t *testing.T, gen *Generator) string {
	t.Helper()

	for _, event := range slices.Backward(gen.log.Events()) {
		if event.Kind == EventConsequence {
			return event.Consequence.Detail
		}
	}

	t.Fatal("no consequence in the log")

	return ""
}
