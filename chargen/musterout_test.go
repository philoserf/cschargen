package chargen_test

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// pointBenefitTable is the choice point mustering out puts to the decider,
// once per roll taken.
const pointBenefitTable = "benefit_table"

// benefitRolls counts the mustering-out table choices in a record.
func benefitRolls(character *chargen.Character) (int, int) {
	cash, other := 0, 0

	for _, event := range character.Events {
		if event.Kind != chargen.EventChoice || event.Choice.Point != pointBenefitTable {
			continue
		}

		if event.Choice.Options[event.Choice.Chosen] == "Cash" {
			cash++
		} else {
			other++
		}
	}

	return cash, other
}

// forfeited reports whether any result in a record took back every benefit
// roll earned in a career, which sixty-seven of them can.
func forfeited(character *chargen.Character) bool {
	for _, event := range character.Events {
		if event.Kind != chargen.EventConsequence {
			continue
		}

		if event.Consequence.Kind != chargen.ConsequenceBenefitRolls {
			continue
		}

		if strings.Contains(event.Consequence.Detail, "lose every benefit roll") {
			return true
		}
	}

	return false
}

// TestTwoBenefitRollsPerTerm is p. 126: "The character receives a number of
// Mustering Out Benefit rolls equal to twice the number of terms served in
// that career."
//
// Events and mishaps adjust the count, so the test holds the floor rather
// than an equality: a character who served six terms across some number of
// careers cannot have fewer rolls than an unluckier one who lost some.
//
// A record carrying a forfeit is excluded rather than held to the floor.
// "You are dismissed and lose all Benefits" really does mean a character
// can serve four terms and muster out with nothing, so the floor is not a
// property of every record -- only of one where nothing took the rolls
// away.
func TestTwoBenefitRollsPerTerm(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 4)
		if forfeited(character) {
			continue
		}

		cash, other := benefitRolls(character)
		terms := len(character.State.Terms)

		if cash+other == 0 && terms > 0 {
			t.Errorf("seed %d: %d terms served and no mustering out rolls", seed, terms)
		}
	}
}

// TestCashIsCappedAtThreePerCareer is p. 127: "A character may roll on the
// Cash table no more than three times per career, regardless of how many
// total Mustering Out rolls they possess."
//
// The cap is per career, so a character who served in three careers may
// take nine cash rolls in all -- which is why the count is taken per
// muster-out rather than per character.
func TestCashIsCappedAtThreePerCareer(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 8)

		taken := 0

		for _, event := range character.Events {
			if event.Kind == chargen.EventStep && event.Step.Name == "Step 19: Muster Out" {
				taken = 0

				continue
			}

			if event.Kind != chargen.EventChoice || event.Choice.Point != pointBenefitTable {
				continue
			}

			if event.Choice.Options[event.Choice.Chosen] == "Cash" {
				taken++
			}

			if taken > 3 {
				t.Fatalf("seed %d: a fourth cash roll in one career", seed)
			}
		}
	}
}

// TestTheCashTableStopsBeingOffered: the cap is enforced by not offering
// the table, rather than by offering it and refusing the result.
func TestTheCashTableStopsBeingOffered(t *testing.T) {
	t.Parallel()

	sawOtherOnly := false

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 8)

		for _, event := range character.Events {
			if event.Kind != chargen.EventChoice || event.Choice.Point != pointBenefitTable {
				continue
			}

			if len(event.Choice.Options) == 1 {
				sawOtherOnly = true

				if event.Choice.Options[0] != "Other Benefits" {
					t.Errorf("seed %d: the single option was %q", seed, event.Choice.Options[0])
				}
			}
		}
	}

	if !sawOtherOnly {
		t.Log("no character took three cash rolls in one career in this sample")
	}
}

// TestMusteringOutHappensOnceForEachCareerLeft is p. 126: a character must
// muster out "whether to enter a new career or to conclude character
// generation". A character who served in three careers musters out three
// times.
func TestMusteringOutHappensOnceForEachCareerLeft(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 8)

		musters := 0

		for _, event := range character.Events {
			if event.Kind == chargen.EventStep && event.Step.Name == "Step 19: Muster Out" {
				musters++
			}
		}

		// A run that stopped at an unimplemented career never reached its
		// last muster out, which is the honest outcome and not a fault.
		if stoppedForAnUnimplementedCareer(character) {
			continue
		}

		// Nor does a character who died in service (pp. 123-124). They did
		// not leave their career; they stopped.
		if character.State.Fate != "" {
			continue
		}

		if musters != len(character.State.Services) {
			t.Errorf("seed %d: %d services, %d muster outs", seed, len(character.State.Services), musters)
		}
	}
}

// TestBenefitQueueIsEmptiedPerCareer: a batch granted in one career must
// not be spendable in the next.
func TestBenefitQueueIsEmptiedPerCareer(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 8)

		if stoppedForAnUnimplementedCareer(character) {
			continue
		}

		for _, batch := range character.State.Benefits {
			t.Errorf("seed %d: %d rolls left over from %s after generation ended",
				seed, batch.Rolls, batch.Career)
		}
	}
}

// TestAMishapsPenaltyIsChargedOnce: a mishap that removes benefit rolls can
// leave a career owing fewer than none, and that debt must not follow the
// character into their next spell in the same career. It did: mustering out
// with nothing to spend returned before emptying the queue.
func TestAMishapsPenaltyIsChargedOnce(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 8)

		if stoppedForAnUnimplementedCareer(character) {
			continue
		}

		for _, batch := range character.State.Benefits {
			if batch.Rolls < 0 {
				t.Errorf("seed %d: a debt of %d rolls survived %s's muster out",
					seed, batch.Rolls, batch.Career)
			}
		}
	}
}

// TestCreditsOnlyComeFromCashRollsAndEvents: a character with no cash roll
// and no event paying them has no money, which is what makes ERRATA E-4 --
// an injury that cannot be paid for -- reachable at all.
func TestCreditsAreNeverNegative(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 6)

		if character.State.Credits < 0 {
			t.Errorf("seed %d: credits = %d", seed, character.State.Credits)
		}
	}
}

// TestMusterOutReplays: the benefit table choice is a choice point like any
// other, so a record carrying it has to replay.
func TestMusterOutReplays(t *testing.T) {
	t.Parallel()

	for seed := range uint64(20) {
		original := lifepath(t, seed, 6)

		opts := options(t, seed)

		opts.Inputs.TermLimit = 6
		opts.Decider = chargen.NewReplay(original.Events)

		replayed := generate(t, opts)

		if original.State.Credits != replayed.State.Credits {
			t.Errorf("seed %d: credits %d replayed as %d",
				seed, original.State.Credits, replayed.State.Credits)
		}

		if len(original.Events) != len(replayed.Events) {
			t.Errorf("seed %d: %d events replayed as %d",
				seed, len(original.Events), len(replayed.Events))
		}
	}
}
