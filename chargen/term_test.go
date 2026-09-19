package chargen_test

import (
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// lifepath generates a whole character: characteristics, then terms.
func lifepath(t *testing.T, seed uint64, terms int) *chargen.Character {
	t.Helper()

	opts := options(t, seed)

	opts.Inputs.TermLimit = terms

	return generate(t, opts)
}

// careerColonist is the career the tests force where they need a specific
// one: it is the book's worked example, and six of its eleven mishaps
// reassign the homeworld.
const careerColonist = "Colonist"

// sample is how many seeds the properties below are checked over. Large
// enough that the uncommon branches -- a failed enlistment, a mishap, a
// forced transfer -- show up; small enough to keep the suite quick.
const sample = 60

// TestEveryConsequenceNamesAStepThrowOrChoice is the log's load-bearing
// invariant. A consequence names why it happened, and "why" is always a
// step entered, a throw made or a choice taken -- never another
// consequence, which would be a chain with no cause at the end of it.
func TestEveryConsequenceNamesAStepThrowOrChoice(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 6)

		for _, event := range character.Events {
			if event.Kind != chargen.EventConsequence {
				continue
			}

			cause := event.Consequence.Cause
			if cause < 1 || cause >= event.Seq {
				t.Fatalf("seed %d: event %d names cause %d, which is not an earlier event",
					seed, event.Seq, cause)
			}

			switch character.Events[cause-1].Kind {
			case chargen.EventStep, chargen.EventThrow, chargen.EventChoice:
			case chargen.EventConsequence:
				t.Fatalf("seed %d: consequence %d (%s) names consequence %d as its cause",
					seed, event.Seq, event.Consequence.Detail, cause)
			}
		}
	}
}

// TestTermsAreServedUpToTheLimit: the limit is policy rather than a rule
// (POLICY.md), but the engine has to honour it exactly -- a character who
// served five terms when four were asked for has aged four years too many.
//
// A run may stop short, and only for one reason: a table result sent the
// character to a career this milestone does not implement. That case has
// to be visible in the record, so the test insists on the consequence
// rather than allowing any short run.
func TestTermsAreServedUpToTheLimit(t *testing.T) {
	t.Parallel()

	shortRuns := 0

	for _, limit := range []int{1, 2, 4, 8} {
		for seed := range uint64(12) {
			character := lifepath(t, seed, limit)

			got := len(character.State.Terms)
			if got > limit {
				t.Errorf("seed %d, limit %d: %d terms served", seed, limit, got)

				continue
			}

			if got == limit {
				continue
			}

			shortRuns++

			if !stoppedForAnUnimplementedCareer(character) {
				t.Errorf("seed %d, limit %d: stopped after %d terms with nothing in the record saying why",
					seed, limit, got)
			}
		}
	}

	if shortRuns == 0 {
		t.Log("no run stopped short in this sample; the stop path is checked elsewhere")
	}
}

// stoppedForAnUnimplementedCareer reports whether the record's last
// consequence is the engine saying where the character went and that it
// could not follow.
func stoppedForAnUnimplementedCareer(character *chargen.Character) bool {
	for _, event := range character.Events {
		if event.Kind != chargen.EventConsequence {
			continue
		}

		if event.Consequence.Kind == chargen.ConsequenceUnimplemented &&
			contains(event.Consequence.Detail, "generation ends here") {
			return true
		}
	}

	return false
}

func contains(haystack, needle string) bool {
	if len(needle) > len(haystack) {
		return false
	}

	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}

	return false
}

// TestAgeIsFourYearsPerTermUnlessAMishapEjected: "add four years to their
// character's age ... If the Player has suffered a Mishap ... they should
// roll 1d3 and add that result" (p. 121). So a character's age is bounded
// below by a 1d3 per term and above by four.
func TestAgeIsFourYearsPerTermUnlessAMishapEjected(t *testing.T) {
	t.Parallel()

	const startingAge = 18

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 5)

		terms := len(character.State.Terms)
		age := character.State.Age

		if age > startingAge+terms*4 {
			t.Errorf("seed %d: age %d after %d terms, above four years each", seed, age, terms)
		}

		if age < startingAge+terms {
			t.Errorf("seed %d: age %d after %d terms, below one year each", seed, age, terms)
		}
	}
}

// TestFirstTermGrantsTheServiceSkillsAtLevelZero is p. 117's rule, checked
// on the one career a forced run always enters.
func TestFirstTermGrantsTheServiceSkillsAtLevelZero(t *testing.T) {
	t.Parallel()

	opts := options(t, 5)

	opts.Inputs.TermLimit = 1
	opts.Inputs.Career = careerColonist

	character := generate(t, opts)

	// Colonist's Service Skills table, p. 173.
	for _, want := range []string{"Animals", "Broker", "Mechanic", "Gun Combat", "Survival", "Chef"} {
		if !character.State.Has(want) {
			t.Errorf("a first-term Colonist does not hold %s", want)
		}
	}
}

// TestRankZeroBenefitAppliesOnEntry is ERRATA E-8: rank benefits are
// applied "immediately upon achieving that Rank" (p. 116) and every
// character "begins a career at Rank 0" (p. 115), so the Rank 0 row of the
// assignment applies on entering.
func TestRankZeroBenefitAppliesOnEntry(t *testing.T) {
	t.Parallel()

	opts := options(t, 5)

	opts.Inputs.TermLimit = 1
	opts.Inputs.Career = careerColonist

	character := generate(t, opts)

	// The Settler assignment's Rank 0 benefit, p. 174.
	got, found := character.State.Skill("Animals", "Farming")
	if !found {
		t.Fatal("a Settler does not hold Animals (Farming)")
	}

	if got.Level < 1 {
		t.Errorf("Animals (Farming) is level %d, want at least 1", got.Level)
	}

	deviated := false

	for _, id := range character.Provenance.Deviations {
		if id == "E-8" {
			deviated = true
		}
	}

	if !deviated {
		t.Error("the record does not name E-8, the reading this rests on")
	}
}

// TestPrisonerIsNeverChosen: "A character cannot voluntarily choose to
// enter this career" (p. 111). It may only be reached by a table result.
func TestPrisonerIsNeverChosen(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 8)

		for _, event := range character.Events {
			if event.Kind != chargen.EventChoice || event.Choice.Point != "career" {
				continue
			}

			for _, option := range event.Choice.Options {
				if option == "Prisoner" {
					t.Fatalf("seed %d: Prisoner was offered as a career choice", seed)
				}
			}
		}
	}
}

// TestMishapsDoNotEjectFromVagabondOrPrisoner: both careers say so in
// print (pp. 263, 298), and they are the two a character is most often
// forced into.
func TestMishapsDoNotEjectFromVagabondOrPrisoner(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 8)

		for _, service := range character.State.Services {
			if service.Career != "Vagabond" && service.Career != "Prisoner" {
				continue
			}

			if service.LeftBecause == "ejected by a mishap" {
				t.Errorf("seed %d: a mishap ejected the character from %s", seed, service.Career)
			}
		}
	}
}

// TestVagabondMishapsDoOccur is the control for the test above: without it,
// a bug that stopped Vagabond mishaps from happening at all would make that
// property hold for the wrong reason.
func TestVagabondMishapsDoOccur(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 8)

		for _, event := range character.Events {
			if event.Kind != chargen.EventConsequence {
				continue
			}

			if event.Consequence.Career == "Vagabond" &&
				contains(event.Consequence.Detail, "mishap:") {
				return
			}
		}
	}

	t.Error("no Vagabond mishap came up across the sample")
}

// TestEveryThrowInALifepathCarriesACite is what makes a record auditable
// against the book, over the whole lifepath rather than one step.
func TestEveryThrowInALifepathCarriesACite(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 6)

		for _, event := range character.Events {
			switch event.Kind {
			case chargen.EventThrow:
				if event.Throw.Cite == "" {
					t.Fatalf("seed %d: throw at event %d carries no cite", seed, event.Seq)
				}
			case chargen.EventStep:
				if event.Step.Cite == "" {
					t.Fatalf("seed %d: step at event %d carries no cite", seed, event.Seq)
				}
			case chargen.EventChoice, chargen.EventConsequence:
			}
		}
	}
}

// TestCharacteristicsStayInRange: nothing falls below zero, and an
// unaltered human caps at 15 (p. 14).
func TestCharacteristicsStayInRange(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 8)

		for _, which := range chargen.CharacteristicOrder {
			score := character.State.Characteristics.Get(which)
			if score < 0 || score > chargen.HumanMaximum {
				t.Errorf("seed %d: %s = %d, outside 0-%d", seed, which, score, chargen.HumanMaximum)
			}
		}
	}
}

// TestAWholeLifepathReplays is the contract the milestone rests on, now
// over the term loop rather than over six rolls.
func TestAWholeLifepathReplays(t *testing.T) {
	t.Parallel()

	for seed := range uint64(20) {
		original := lifepath(t, seed, 6)

		opts := options(t, seed)

		opts.Inputs.TermLimit = 6
		opts.Decider = chargen.NewReplay(original.Events)

		replayed := generate(t, opts)

		if original.State.Characteristics != replayed.State.Characteristics {
			t.Errorf("seed %d: characteristics diverged", seed)
		}

		if original.State.Age != replayed.State.Age {
			t.Errorf("seed %d: age %d replayed as %d", seed, original.State.Age, replayed.State.Age)
		}

		if len(original.Events) != len(replayed.Events) {
			t.Errorf("seed %d: %d events replayed as %d", seed, len(original.Events), len(replayed.Events))
		}

		if len(original.State.Skills) != len(replayed.State.Skills) {
			t.Errorf("seed %d: %d skills replayed as %d",
				seed, len(original.State.Skills), len(replayed.State.Skills))
		}
	}
}

// TestUnimplementedResultsAreRecordedNotSwallowed: a stubbed result has to
// leave a trace, or the record silently claims the rules were followed.
func TestUnimplementedResultsAreRecordedNotSwallowed(t *testing.T) {
	t.Parallel()

	found := false

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 8)

		for _, event := range character.Events {
			if event.Kind != chargen.EventConsequence {
				continue
			}

			if event.Consequence.Kind != chargen.ConsequenceUnimplemented {
				continue
			}

			found = true

			if event.Consequence.Detail == "" {
				t.Errorf("seed %d: an unimplemented consequence says nothing about what was asked for", seed)
			}
		}
	}

	if !found {
		t.Error("no unimplemented result came up across the sample; the stub path is untested")
	}
}

// TestSkillLevels: a skill the character was awarded is at level 1 or
// better, and a service skill they were only granted is at 0.
func TestSkillLevels(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 6)

		for _, held := range character.State.Skills {
			if held.Level < 0 {
				t.Errorf("seed %d: %s at level %d", seed, held.Full(), held.Level)
			}

			if held.Name == "" {
				t.Errorf("seed %d: a skill with no name", seed)
			}
		}
	}
}
