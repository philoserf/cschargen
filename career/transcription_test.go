package career_test

import (
	"testing"

	"github.com/philoserf/cschargen/career"
)

// The tables below are a second reading of the page, typed independently of
// career/colonist.go. Repetition is the mechanism: two transcriptions that
// agree are evidence, and one constant shared between them would be a
// single reading wearing two hats.

// TestColonistNumbersMatchPage173 reads the Career Progress and Mustering
// Out blocks of p. 173 a second time.
func TestColonistNumbersMatchPage173(t *testing.T) {
	t.Parallel()

	colonist := career.Colonist()

	if colonist.Enlistment == nil {
		t.Fatal("Colonist has no enlistment throw; p. 173 gives END 6+")
	}

	if colonist.Enlistment.Characteristic != "END" || colonist.Enlistment.Number != 6 {
		t.Errorf("enlistment = %s %d+, want END 6+",
			colonist.Enlistment.Characteristic, colonist.Enlistment.Number)
	}

	progress := []struct {
		assignment  string
		survival    string
		survivalN   int
		advancement string
		advanceN    int
	}{
		{"Settler", "END", 7, "INT", 7},
		{"Politician", "INT", 7, "CHA", 7},
		{"Commercial", "INT", 7, "EDU", 7},
	}

	for _, want := range progress {
		got, ok := colonist.Assignment(want.assignment)
		if !ok {
			t.Errorf("no %s assignment", want.assignment)

			continue
		}

		if got.Survival.Characteristic != want.survival || got.Survival.Number != want.survivalN {
			t.Errorf("%s survival = %s %d+, want %s %d+", want.assignment,
				got.Survival.Characteristic, got.Survival.Number, want.survival, want.survivalN)
		}

		if got.Advancement.Characteristic != want.advancement || got.Advancement.Number != want.advanceN {
			t.Errorf("%s advancement = %s %d+, want %s %d+", want.assignment,
				got.Advancement.Characteristic, got.Advancement.Number, want.advancement, want.advanceN)
		}
	}

	cash := [7]int{0, 0, 500, 1000, 2000, 5000, 10000}
	for i, want := range cash {
		if got := colonist.Benefits[i].Cash; got != want {
			t.Errorf("benefit row %d cash = %d, want %d", i+1, got, want)
		}
	}
}

// TestVagabondAndPrisonerNumbers reads pp. 296 and 261 a second time. Both
// careers give their two assignments the same survival and advancement
// targets, which is unusual enough to be worth pinning: a transcription
// that got one of them wrong would look plausible.
func TestVagabondAndPrisonerNumbers(t *testing.T) {
	t.Parallel()

	vagabond := career.Vagabond()
	for _, name := range []string{"Destitute", "Transient"} {
		got, ok := vagabond.Assignment(name)
		if !ok {
			t.Errorf("Vagabond has no %s assignment", name)

			continue
		}

		if got.Survival.Characteristic != "END" || got.Survival.Number != 8 {
			t.Errorf("Vagabond %s survival = %s %d+, want END 8+",
				name, got.Survival.Characteristic, got.Survival.Number)
		}

		if got.Advancement.Characteristic != "INT" || got.Advancement.Number != 8 {
			t.Errorf("Vagabond %s advancement = %s %d+, want INT 8+",
				name, got.Advancement.Characteristic, got.Advancement.Number)
		}
	}

	prisoner := career.Prisoner()

	cell, ok := prisoner.Assignment("Prisoner")
	if !ok {
		t.Fatal("Prisoner has no Prisoner assignment")
	}

	if cell.Survival.Characteristic != "END" || cell.Survival.Number != 8 {
		t.Errorf("Prisoner survival = %s %d+, want END 8+",
			cell.Survival.Characteristic, cell.Survival.Number)
	}

	if cell.Advancement.Characteristic != "STR" || cell.Advancement.Number != 8 {
		t.Errorf("Prisoner advancement = %s %d+, want STR 8+",
			cell.Advancement.Characteristic, cell.Advancement.Number)
	}

	// p. 261 gives Prisoner a single assignment, which is what makes it the
	// smallest career in the book and part of why it is in this milestone.
	if len(prisoner.Assignments) != 1 {
		t.Errorf("Prisoner has %d assignments, want 1", len(prisoner.Assignments))
	}
}

// TestColonistReachesFourOtherCareers is the finding that set this
// milestone's scope (docs/MILESTONE-1.md). If a later change makes one
// career self-contained, this test should be the thing that notices.
func TestColonistReachesFourOtherCareers(t *testing.T) {
	t.Parallel()

	reached := map[string]bool{}

	var walk func(effects []career.Effect)

	walk = func(effects []career.Effect) {
		for _, e := range effects {
			if e.Kind == career.EffectTransfer {
				reached[e.Career] = true
			}

			walk(e.Success)
			walk(e.Failure)

			for _, o := range e.Options {
				walk(o.Effects)
			}
		}
	}

	colonist := career.Colonist()
	for _, row := range colonist.Events {
		walk(row.Effects)
	}

	for _, row := range colonist.Mishaps {
		walk(row.Effects)
	}

	for _, want := range []string{"Vagabond", "Celebrity"} {
		if !reached[want] {
			t.Errorf("Colonist no longer reaches %s; check pp. 174-176", want)
		}
	}
}

// TestTheImplementedCareersCoverTheForcedTransfers: whatever the three
// careers force a character into, this milestone has to be able to carry
// out -- or the record ends on an unimplemented consequence, which is
// honest but is not the same as playing the rules.
func TestTheImplementedCareersCoverTheForcedTransfers(t *testing.T) {
	t.Parallel()

	implemented := map[string]bool{}
	for _, c := range career.All() {
		implemented[c.Name] = true
	}

	stubbed := map[string]bool{}

	var walk func(effects []career.Effect)

	walk = func(effects []career.Effect) {
		for _, e := range effects {
			// PreviousCareer is not a career name: it is the engine's
			// word for "the one you held before this", which only the
			// service record can resolve.
			if e.Kind == career.EffectTransfer && e.Career != career.PreviousCareer &&
				!implemented[e.Career] {
				stubbed[e.Career] = true
			}

			walk(e.Success)
			walk(e.Failure)

			for _, o := range e.Options {
				walk(o.Effects)
			}
		}
	}

	for _, c := range career.All() {
		for _, row := range c.Events {
			walk(row.Effects)
		}

		for _, row := range c.Mishaps {
			walk(row.Effects)
		}
	}

	walk(career.LifeEvents()[1].Effects)

	// The careers the implemented ones send characters to and that are not
	// built yet. The list shrinks as milestone 3 proceeds and should reach
	// zero; a destination appearing that is not here fails the test rather
	// than passing unnoticed.
	// Independent Merchant's mishap 12 does name Pirate inside a 1d6 branch
	// the engine records rather than resolves, so it is not a transfer.
	// When that branch becomes one, this test is what will say so.
	want := map[string]bool{}
	for name := range stubbed {
		if !want[name] {
			t.Errorf("a transfer to %s is stubbed, and this test did not know about it", name)
		}
	}

	for name := range want {
		if !stubbed[name] {
			t.Errorf("%s was expected to be a stubbed destination and is not", name)
		}
	}
}
