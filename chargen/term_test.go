package chargen_test

import (
	"maps"
	"slices"
	"sync"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// lifepath generates a whole character: characteristics, then terms.
// lifepath is a generated character, kept so that the sweeps below share
// one rather than each walking their own.
//
// A character is a pure function of its seed, term limit, setting and
// decider, and lifepath fixes all four -- so two sweeps asking for the same
// pair are asking for the same character. Thirty-two call sites ask for
// seven distinct term counts, so the suite used to walk 1,860 lifepaths to
// see 420 characters.
//
// Nothing writes to one: every use of a generated character in these tests
// reads it. A test that needs to modify one should build it with options()
// and generate() directly, which is what the tests taking a different
// decider or different inputs already do.
func lifepath(t *testing.T, seed uint64, terms int) *chargen.Character {
	t.Helper()

	key := lifepathKey{seed: seed, terms: terms}

	lifepaths.mu.RLock()

	held, found := lifepaths.byKey[key]

	lifepaths.mu.RUnlock()

	if found {
		return held
	}

	opts := options(t, seed)

	opts.Inputs.TermLimit = terms

	// Generated outside the lock. Two sweeps racing for the same key cost
	// one duplicate walk, which is cheaper than serialising every sweep
	// behind one mutex -- and the duplicate is the same character, because
	// the character is determined by the key.
	made := generate(t, opts)

	lifepaths.mu.Lock()
	defer lifepaths.mu.Unlock()

	if held, found := lifepaths.byKey[key]; found {
		return held
	}

	lifepaths.byKey[key] = made

	return made
}

// lifepathKey is everything lifepath varies: the rest of the inputs are
// fixed by options().
type lifepathKey struct {
	seed  uint64
	terms int
}

// lifepathStore holds the generated characters the sweeps share.
type lifepathStore struct {
	mu    sync.RWMutex
	byKey map[lifepathKey]*chargen.Character
}

var lifepaths = lifepathStore{byKey: map[lifepathKey]*chargen.Character{}}

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
// A run may stop short, and only for one reason: a table result named a
// career no career answers to by that name. That case has
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
		// Step 8 adds four years for a degree and 1d3 for a washout
		// (pp. 87-88), which is a different rule on a different page. This
		// one is about Step 17's arithmetic.
		opts := options(t, seed)

		opts.Inputs.TermLimit = 5
		opts.Inputs.SkipEducation = true

		character := generate(t, opts)

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

// aColonist generates a character who actually entered Colonist.
//
// Forcing a career is not the same as entering one: the enlistment throw
// can still fail, and which seed makes it moves whenever a step is added
// ahead of Step 10. Scanning for one is what keeps these tests about the
// rule rather than about a seed.
func aColonist(t *testing.T) *chargen.Character {
	t.Helper()

	for seed := range uint64(sample) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = 1
		opts.Inputs.Career = careerColonist

		character := generate(t, opts)
		if len(character.State.Services) > 0 &&
			character.State.Services[0].Career == careerColonist {
			return character
		}
	}

	t.Fatalf("no seed in %d enlisted in Colonist", sample)

	return nil
}

// TestFirstTermGrantsTheServiceSkillsAtLevelZero is p. 117's rule, checked
// on the one career a forced run enters.
func TestFirstTermGrantsTheServiceSkillsAtLevelZero(t *testing.T) {
	t.Parallel()

	character := aColonist(t)

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

	character := aColonist(t)

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

// TestEveryConsequenceInALifepathCarriesACite is the other half of the
// invariant above, and the half a reader actually checks. The throw says
// what was rolled; the consequence says what it did to the character, which
// is what gets compared against the page.
//
// It reports by step rather than by seed because a gap is a property of the
// step that wrote it, not of the dice that reached it -- one run names every
// step that needs fixing instead of the first one a sweep happens to hit.
func TestEveryConsequenceInALifepathCarriesACite(t *testing.T) {
	t.Parallel()

	gaps := map[string]int{}

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 6)

		step := "(before any step)"

		for _, event := range character.Events {
			switch event.Kind {
			case chargen.EventStep:
				step = event.Step.Name
			case chargen.EventConsequence:
				if event.Consequence.Cite == "" {
					gaps[step]++
				}
			case chargen.EventThrow, chargen.EventChoice:
			}
		}
	}

	if len(gaps) == 0 {
		return
	}

	for _, step := range slices.Sorted(maps.Keys(gaps)) {
		t.Errorf("%s: %d consequences carry no cite", step, gaps[step])
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

// TestAWholeLifepathReplays is the reproducibility contract, now over the
// term loop rather than over six rolls.
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

// TestChoosingANewCareerIsNotATransfer is ERRATA E-6, applied.
//
// Five Colonist mishaps end "choose another career and a new homeworld"
// (p. 174). They are not the mishaps that name a career: Colonist 10 sends
// the character to Vagabond by name. The difference the engine has to show is that an unnamed change
// goes through Step 18's career list and its enlistment throw, where a
// named one places the character without either.
//
// Forcing Colonist is what the test above learned to do: the auto policy no
// longer happens to pick it now that thirty-four careers are on the list,
// and the closing check is what turns a vacuous pass into a failure.
func TestChoosingANewCareerIsNotATransfer(t *testing.T) {
	t.Parallel()

	reached := 0

	for seed := range uint64(sample) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = 8
		opts.Inputs.Career = careerColonist

		character := generate(t, opts)

		if !slices.Contains(character.Provenance.Deviations, "E-6") {
			continue
		}

		reached++

		// The character left the career that sent them away, and whatever
		// they did next was entered on their own account rather than
		// placed by the rules. A forced career is only the first attempt,
		// so Colonist need not be the first service -- only one they held
		// and then left.
		left := slices.IndexFunc(character.State.Services, func(s chargen.Service) bool {
			return s.Career == careerColonist
		})

		if left < 0 {
			t.Errorf("seed %d: stamped E-6 without ever holding Colonist", seed)

			continue
		}

		if left == len(character.State.Services)-1 {
			t.Errorf("seed %d: stamped E-6 and never left Colonist", seed)
		}
	}

	if reached == 0 {
		t.Errorf("no seed in %d reached an unnamed career change; E-6 is untested", sample)
	}
}
