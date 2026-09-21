package chargen_test

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// TestEveryTieNamesItsOrigin is FR15: the record carries "relationships
// (allies, contacts, rivals, enemies) with their origin".
//
// A tie granted before Step 9 has no career to name, and for a long while
// said nothing instead — 138 of 2,477 ties across a sample pool, every one
// of them from the youth tables, the teenage tables or a school. A parent
// and a childhood mentor are both Allies at 125, and the origin is the only
// thing that tells them apart.
func TestEveryTieNamesItsOrigin(t *testing.T) {
	t.Parallel()

	checked := 0

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 4)

		for _, tie := range character.State.Ties {
			checked++

			if tie.Origin == "" {
				t.Errorf("seed %d: a %s at %d names no origin",
					seed, tie.Kind, tie.Rating)
			}
		}
	}

	if checked == 0 {
		t.Fatal("no character in the sample has a relationship at all")
	}
}

// TestStep9AlwaysNamesACareer. Step 9 has three ways out — the career the
// decider chose, the one a table result forced, and the Vagabond a character
// drifts into when nothing will have them — and for a while the third said
// nothing at all, leaving a step header with no events under it.
//
// The invariant is per-step rather than general: Step 18 is legitimately
// empty in most records, because deciding to carry on is not an event.
func TestStep9AlwaysNamesACareer(t *testing.T) {
	t.Parallel()

	const step9 = "Step 9: Choose a Career"

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 4)

		inStep9, spoke := false, true

		for _, event := range character.Events {
			if event.Kind == chargen.EventStep {
				if inStep9 && !spoke {
					t.Errorf("seed %d: Step 9 passed without naming a career", seed)
				}

				inStep9 = strings.HasPrefix(event.Step.Name, step9)
				spoke = false

				continue
			}

			if inStep9 && event.Kind != chargen.EventStep {
				spoke = true
			}
		}

		if inStep9 && !spoke {
			t.Errorf("seed %d: the last Step 9 passed without naming a career", seed)
		}
	}
}
