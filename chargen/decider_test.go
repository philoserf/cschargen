package chargen_test

import (
	"errors"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// The two option strings the replay tests repeat. Named so that a typo in
// one case reads as a typo rather than as a disagreement between cases.
const (
	pointCareer = "career"
	optColonist = "colonist"
)

func choice(point string, options ...string) chargen.Choice {
	return chargen.Choice{Point: point, Prompt: point, Options: options, Cite: "p. 13"}
}

func TestPolicyTakesTheFirstOption(t *testing.T) {
	t.Parallel()

	got, err := chargen.Policy{}.Choose(choice("assign_characteristics", "STR", "DEX", "END"))
	if err != nil {
		t.Fatalf("Choose: %v", err)
	}

	if got != 0 {
		t.Errorf("Choose = %d, want 0", got)
	}
}

func TestPolicyRefusesAnEmptyOptionList(t *testing.T) {
	t.Parallel()

	_, err := chargen.Policy{}.Choose(choice("nothing"))
	if !errors.Is(err, chargen.ErrNoOptions) {
		t.Errorf("err = %v, want ErrNoOptions", err)
	}
}

func TestDeciderKinds(t *testing.T) {
	t.Parallel()

	if got := (chargen.Policy{}).Kind(); got != chargen.DeciderPolicy {
		t.Errorf("Policy.Kind() = %q", got)
	}

	if got := chargen.NewReplay(nil).Kind(); got != chargen.DeciderReplay {
		t.Errorf("Replay.Kind() = %q", got)
	}
}

// recorded builds the log a record would carry for the given choices.
func recorded(choices ...chargen.ChoiceEvent) []chargen.Event {
	var log chargen.Log

	log.Step("Step 2", "p. 39")

	for _, c := range choices {
		log.Choice(c)
	}

	return log.Events()
}

func TestReplayReappliesRecordedChoices(t *testing.T) {
	t.Parallel()

	events := recorded(
		chargen.ChoiceEvent{Point: "first", Options: []string{"a", "b"}, Chosen: 1},
		chargen.ChoiceEvent{Point: "second", Options: []string{"x", "y", "z"}, Chosen: 2},
	)

	replay := chargen.NewReplay(events)

	got, err := replay.Choose(choice("first", "a", "b"))
	if err != nil {
		t.Fatalf("first: %v", err)
	}

	if got != 1 {
		t.Errorf("first = %d, want 1", got)
	}

	got, err = replay.Choose(choice("second", "x", "y", "z"))
	if err != nil {
		t.Fatalf("second: %v", err)
	}

	if got != 2 {
		t.Errorf("second = %d, want 2", got)
	}

	if replay.Remaining() != 0 {
		t.Errorf("Remaining() = %d after reapplying both", replay.Remaining())
	}
}

// TestReplayChecksTheQuestion is why a recorded index is not enough on its
// own: the index means nothing against a different list, which is how a
// record generated with a forced career would otherwise replay against the
// full list of eligible ones and pick the wrong entry.
func TestReplayChecksTheQuestion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		recordedC chargen.ChoiceEvent
		asked     chargen.Choice
	}{
		{
			"a different choice point",
			chargen.ChoiceEvent{Point: pointCareer, Options: []string{"a"}, Chosen: 0},
			choice("assignment", "a"),
		},
		{
			"the same point, more options than were recorded",
			chargen.ChoiceEvent{Point: pointCareer, Options: []string{optColonist}, Chosen: 0},
			choice(pointCareer, optColonist, "vagabond"),
		},
		{
			"the same point, the options reordered",
			chargen.ChoiceEvent{Point: pointCareer, Options: []string{optColonist, "vagabond"}, Chosen: 0},
			choice(pointCareer, "vagabond", optColonist),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			replay := chargen.NewReplay(recorded(tc.recordedC))

			_, err := replay.Choose(tc.asked)
			if !errors.Is(err, chargen.ErrReplayDiverged) {
				t.Errorf("err = %v, want ErrReplayDiverged", err)
			}
		})
	}
}

// TestReplayAcceptsTheSameQuestion is the control for the table above: the
// same point with the same options in the same order replays cleanly, so
// those three failures are about the mismatch and not about the mechanism.
func TestReplayAcceptsTheSameQuestion(t *testing.T) {
	t.Parallel()

	replay := chargen.NewReplay(recorded(
		chargen.ChoiceEvent{Point: pointCareer, Options: []string{"a", "b", "c"}, Chosen: 2},
	))

	got, err := replay.Choose(choice(pointCareer, "a", "b", "c"))
	if err != nil {
		t.Fatalf("Choose: %v", err)
	}

	if got != 2 {
		t.Errorf("Choose = %d, want 2", got)
	}
}

func TestReplayRejectsAnOutOfRangeIndex(t *testing.T) {
	t.Parallel()

	replay := chargen.NewReplay(recorded(
		chargen.ChoiceEvent{Point: pointCareer, Options: []string{"a", "b"}, Chosen: 5},
	))

	_, err := replay.Choose(choice(pointCareer, "a", "b"))
	if !errors.Is(err, chargen.ErrReplayDiverged) {
		t.Errorf("err = %v, want ErrReplayDiverged", err)
	}
}

func TestReplayRunsOut(t *testing.T) {
	t.Parallel()

	replay := chargen.NewReplay(recorded())

	_, err := replay.Choose(choice(pointCareer, "a"))
	if !errors.Is(err, chargen.ErrReplayExhausted) {
		t.Errorf("err = %v, want ErrReplayExhausted", err)
	}
}

// TestReplayRemainingReportsAShortLifepath: a replay that ends with
// recorded choices left over reproduced less than the record holds, and the
// caller has to be able to see that.
func TestReplayRemainingReportsAShortLifepath(t *testing.T) {
	t.Parallel()

	replay := chargen.NewReplay(recorded(
		chargen.ChoiceEvent{Point: "one", Options: []string{"a"}, Chosen: 0},
		chargen.ChoiceEvent{Point: "two", Options: []string{"a"}, Chosen: 0},
	))

	_, err := replay.Choose(choice("one", "a"))
	if err != nil {
		t.Fatalf("Choose: %v", err)
	}

	if got := replay.Remaining(); got != 1 {
		t.Errorf("Remaining() = %d, want 1", got)
	}
}

// TestReplayIgnoresNonChoiceEvents: a record is mostly throws and
// consequences, and the replay decider must walk past them.
func TestReplayIgnoresNonChoiceEvents(t *testing.T) {
	t.Parallel()

	var log chargen.Log

	log.Step("Step 9", "p. 104")
	log.Consequence(chargen.ConsequenceEvent{Kind: chargen.ConsequenceAge, Cause: 1})
	log.Choice(chargen.ChoiceEvent{Point: pointCareer, Options: []string{optColonist}, Chosen: 0})
	log.Consequence(chargen.ConsequenceEvent{Kind: chargen.ConsequenceCareer, Cause: 3})

	replay := chargen.NewReplay(log.Events())
	if replay.Remaining() != 1 {
		t.Fatalf("Remaining() = %d, want 1", replay.Remaining())
	}

	got, err := replay.Choose(choice(pointCareer, optColonist))
	if err != nil || got != 0 {
		t.Errorf("Choose = %d, %v", got, err)
	}
}

func TestSentinelsAreDistinct(t *testing.T) {
	t.Parallel()

	all := []error{chargen.ErrNoOptions, chargen.ErrReplayExhausted, chargen.ErrReplayDiverged}
	for i, a := range all {
		for j, b := range all {
			if i != j && errors.Is(a, b) {
				t.Errorf("sentinel %d matches sentinel %d", i, j)
			}
		}
	}
}
