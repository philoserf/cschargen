package chargen_test

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// consequences collects a record's consequence events of one kind.
func consequences(character *chargen.Character, kind chargen.ConsequenceKind) []chargen.ConsequenceEvent {
	found := make([]chargen.ConsequenceEvent, 0, len(character.Events))

	for _, event := range character.Events {
		if event.Kind == chargen.EventConsequence && event.Consequence.Kind == kind {
			found = append(found, *event.Consequence)
		}
	}

	return found
}

// TestAMishapForfeitLeavesNothingToMusterOut is the sixty-seven results
// that read "you are dismissed and lose all Benefits", seen from the far
// end: a mishap ends the career at once, so the muster-out that follows
// finds an empty queue and says so.
//
// The property is stated against the mishap forfeits alone. A handful of
// these results are events rather than mishaps -- Scavenger's law
// enforcement asking where the collection came from is one -- and the
// character serves on afterwards, earning from the terms that follow.
func TestAMishapForfeitLeavesNothingToMusterOut(t *testing.T) {
	t.Parallel()

	seen := 0

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 6)

		forfeitedIn := ""

		for _, event := range consequences(character, chargen.ConsequenceBenefitRolls) {
			if strings.HasPrefix(event.Detail, "lose every benefit roll from this career") {
				forfeitedIn = event.Career
			}

			if forfeitedIn == "" || event.Career != forfeitedIn {
				continue
			}

			if strings.HasPrefix(event.Detail, "no mustering out rolls remain") {
				seen++

				forfeitedIn = ""
			}
		}

		if forfeitedIn != "" {
			t.Errorf("seed %d: %s forfeited its benefit rolls and mustered out anyway",
				seed, forfeitedIn)
		}
	}

	if seen == 0 {
		t.Fatal("no seed in the sample reached a forfeit; the property is vacuous")
	}
}

// TestImmediateCashRollsPayAtOnce is the handful of events that read "take
// two Cash Benefit rolls immediately": they read the career's Cash column
// where they fire rather than queueing for Step 19, so the credits appear
// in the record before the career ends.
func TestImmediateCashRollsPayAtOnce(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample * 4) {
		character := lifepath(t, seed, 6)

		for _, event := range consequences(character, chargen.ConsequenceCredits) {
			if !strings.HasPrefix(event.Detail, "an immediate cash benefit roll") {
				continue
			}

			if event.Career == "" {
				t.Errorf("seed %d: an immediate cash roll naming no career", seed)
			}

			return
		}
	}

	t.Skip("no seed in the sample reached an immediate cash benefit roll")
}
