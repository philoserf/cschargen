package chargen_test

import (
	"encoding/json"
	"testing"

	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/dice"
)

// TestSequenceNumbersAreMonotonic is the invariant every other reader of
// the log rests on: a consequence names its cause by sequence number, so a
// gap or a repeat would point at the wrong event.
func TestSequenceNumbersAreMonotonic(t *testing.T) {
	t.Parallel()

	var log chargen.Log

	d := dice.New(1)

	log.Step("Step 2: Roll Characteristics", "p. 39")
	log.Roll(d.Characteristic(), "p. 13")
	log.Throw(d.Throw(6), "p. 173")
	log.Choice(chargen.ChoiceEvent{Point: "p", Options: []string{"a"}, Chosen: 0})
	log.Consequence(chargen.ConsequenceEvent{Kind: chargen.ConsequenceSkill, Cause: 2})

	events := log.Events()
	if len(events) != 5 {
		t.Fatalf("logged 5 events, got %d", len(events))
	}

	for i, event := range events {
		if event.Seq != i+1 {
			t.Errorf("event %d has Seq %d", i, event.Seq)
		}
	}

	if log.Len() != 5 {
		t.Errorf("Len() = %d, want 5", log.Len())
	}
}

// TestAppendReturnsTheSequenceNumber is what lets a rule reference the
// throw it just made as the cause of what follows, without counting.
func TestAppendReturnsTheSequenceNumber(t *testing.T) {
	t.Parallel()

	var log chargen.Log

	d := dice.New(2)

	if got := log.Step("s", "p. 1"); got != 1 {
		t.Errorf("first Step returned %d", got)
	}

	if got := log.Roll(d.D6(), "p. 119"); got != 2 {
		t.Errorf("Roll returned %d", got)
	}

	if got := log.Throw(d.Throw(8), "p. 110"); got != 3 {
		t.Errorf("Throw returned %d", got)
	}

	if got := log.Choice(chargen.ChoiceEvent{Point: "p", Options: []string{"a"}}); got != 4 {
		t.Errorf("Choice returned %d", got)
	}

	if got := log.Consequence(chargen.ConsequenceEvent{Kind: chargen.ConsequenceAge, Cause: 1}); got != 5 {
		t.Errorf("Consequence returned %d", got)
	}
}

// TestExactlyOnePayload holds the discriminated union: a reader switching
// on Kind must find the matching payload and no other.
func TestExactlyOnePayload(t *testing.T) {
	t.Parallel()

	var log chargen.Log

	d := dice.New(3)

	log.Step("s", "p. 1")
	log.Roll(d.D66(), "p. 175")
	log.Throw(d.Throw(7), "p. 173")
	log.Choice(chargen.ChoiceEvent{Point: "p", Options: []string{"a"}})
	log.Consequence(chargen.ConsequenceEvent{Kind: chargen.ConsequenceCredits, Cause: 1})

	for _, event := range log.Events() {
		payloads := []bool{
			event.Step != nil, event.Throw != nil, event.Choice != nil, event.Consequence != nil,
		}

		present := 0

		for _, nonNil := range payloads {
			if nonNil {
				present++
			}
		}

		if present != 1 {
			t.Errorf("event %d (%s) has %d payloads", event.Seq, event.Kind, present)
		}

		matches := map[chargen.EventKind]bool{
			chargen.EventStep:        event.Step != nil,
			chargen.EventThrow:       event.Throw != nil,
			chargen.EventChoice:      event.Choice != nil,
			chargen.EventConsequence: event.Consequence != nil,
		}

		if !matches[event.Kind] {
			t.Errorf("event %d says %s but carries a different payload", event.Seq, event.Kind)
		}
	}
}

// TestRollRecordsNoTarget is the distinction of docs/PRD.md FR15: a throw
// against no target leaves Target nil rather than zero, because zero is a
// target no throw of dice satisfies and would read as a check that always
// failed.
func TestRollRecordsNoTarget(t *testing.T) {
	t.Parallel()

	var log chargen.Log

	log.Roll(dice.New(4).Characteristic(), "p. 13")

	got := log.Events()[0].Throw
	if got.Target != nil {
		t.Errorf("Target = %v, want nil", *got.Target)
	}

	if got.Success != nil {
		t.Errorf("Success = %v, want nil", *got.Success)
	}

	if got.Expr != "3d6 drop lowest" {
		t.Errorf("Expr = %q", got.Expr)
	}
}

func TestThrowRecordsTargetAndOutcome(t *testing.T) {
	t.Parallel()

	var log chargen.Log

	throw := dice.New(5).Throw(8, dice.Mod{Name: "END", Value: 1})
	log.Throw(throw, "p. 173")

	got := log.Events()[0].Throw
	if got.Target == nil || *got.Target != 8 {
		t.Fatalf("Target = %v, want 8", got.Target)
	}

	if got.Success == nil || *got.Success != throw.Success {
		t.Fatalf("Success = %v, want %v", got.Success, throw.Success)
	}

	if len(got.Mods) != 1 || got.Mods[0].Name != "END" {
		t.Errorf("mods = %v, want the one END mod", got.Mods)
	}

	if got.Total != throw.Total {
		t.Errorf("Total = %d, want %d", got.Total, throw.Total)
	}
}

// TestEveryThrowCarriesACite is the property that makes the log auditable
// against the book, which is the whole reason it exists.
func TestEveryThrowCarriesACite(t *testing.T) {
	t.Parallel()

	var log chargen.Log

	d := dice.New(6)

	log.Roll(d.D6(), "p. 119")
	log.Throw(d.Throw(8), "p. 110")
	log.Step("Step 12: Roll for Survival", "p. 112")

	for _, event := range log.Events() {
		switch event.Kind {
		case chargen.EventThrow:
			if event.Throw.Cite == "" {
				t.Errorf("event %d: throw with no cite", event.Seq)
			}
		case chargen.EventStep:
			if event.Step.Cite == "" {
				t.Errorf("event %d: step with no cite", event.Seq)
			}
		case chargen.EventChoice, chargen.EventConsequence:
		}
	}
}

// TestJSONOmitsAbsentPayloads keeps the record readable: an event that is a
// step should not serialize three null payload keys beside it.
func TestJSONOmitsAbsentPayloads(t *testing.T) {
	t.Parallel()

	var log chargen.Log

	log.Step("Step 1", "p. 21")

	out, err := json.Marshal(log.Events()[0])
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var back map[string]any

	err = json.Unmarshal(out, &back)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	for _, absent := range []string{"throw", "choice", "consequence"} {
		if _, ok := back[absent]; ok {
			t.Errorf("a step event serialized a %q key", absent)
		}
	}

	if _, ok := back["step"]; !ok {
		t.Error("a step event did not serialize its step payload")
	}
}

func TestEmptyLog(t *testing.T) {
	t.Parallel()

	var log chargen.Log

	if log.Len() != 0 {
		t.Errorf("Len() = %d on a zero-value log", log.Len())
	}

	if log.Events() != nil {
		t.Errorf("Events() = %v on a zero-value log", log.Events())
	}
}
