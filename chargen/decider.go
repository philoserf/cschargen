package chargen

import (
	"fmt"
	"slices"
)

// DeciderKind identifies who resolved a choice.
type DeciderKind string

// The deciders. Player and Policy are the two modes of docs/PRD.md goal 2;
// Replay is the third, and it never decides anything -- it reapplies what a
// record already holds.
const (
	DeciderPlayer DeciderKind = "player"
	DeciderPolicy DeciderKind = "policy"
	DeciderReplay DeciderKind = "replay"
)

// Choice is a choice point put to a [Decider].
//
// Nth and Of place a choice in a run of identical ones: a term's skill
// selections are the same question asked several times, and a player
// answering the fifth cannot otherwise tell it from the first. They are
// engine-provided context, not part of the printed rule.
type Choice struct {
	Point   string // stable identifier, e.g. "assign_characteristics"
	Prompt  string
	Options []string
	Cite    string
	Nth, Of int
}

// Decider resolves choice points. Every choice in the engine goes through
// this interface so that replay can reapply recorded ones.
type Decider interface {
	// Choose returns the index of the selected option. An error refuses the
	// choice outright and ends generation: an interactive session the player
	// abandoned, or a replay whose recorded choice no longer matches the
	// choice the engine presents. Both are distinct from an out-of-range
	// index, which is a decider that answered wrongly rather than one that
	// declined to answer.
	Choose(c Choice) (int, error)

	// Kind identifies the decider in choice events.
	Kind() DeciderKind
}

// Policy is the auto-mode decider: deterministic, total, and tie-breaking
// by the order the book prints its options in. The decision table is
// POLICY.md, identified in every record by policy_version.
type Policy struct{}

// Kind implements [Decider].
func (Policy) Kind() DeciderKind { return DeciderPolicy }

// Choose takes the first option. Every table in the book prints its options
// in a fixed order, so "first listed" is a rule a reader can check rather
// than a preference -- and a policy that chose otherwise would need a
// reason per choice point, which is what later milestones add.
func (Policy) Choose(c Choice) (int, error) {
	if len(c.Options) == 0 {
		return 0, fmt.Errorf("%w: %s", ErrNoOptions, c.Point)
	}

	return 0, nil
}

// Replay reapplies the choices a record already holds, in order. It decides
// nothing: a record that does not match the choices the engine now presents
// is a record the current engine cannot reproduce, and saying so is the
// whole point of replaying it.
type Replay struct {
	choices []ChoiceEvent
	next    int
}

// NewReplay returns a decider that reapplies the choice events of a record,
// in the order they were logged.
func NewReplay(events []Event) *Replay {
	var choices []ChoiceEvent

	for _, e := range events {
		if e.Kind == EventChoice && e.Choice != nil {
			choices = append(choices, *e.Choice)
		}
	}

	return &Replay{choices: choices}
}

// Kind implements [Decider].
func (*Replay) Kind() DeciderKind { return DeciderReplay }

// Choose returns the recorded answer to this choice, after checking that it
// is an answer to the same question. The check is the point and the options
// are part of it: a recorded index means nothing against a different list,
// which is how a record generated with a forced career would otherwise
// replay against the full list of eligible ones and pick the wrong entry.
func (r *Replay) Choose(ask Choice) (int, error) {
	if r.next >= len(r.choices) {
		return 0, fmt.Errorf("%w: %s, after %d choices", ErrReplayExhausted, ask.Point, len(r.choices))
	}

	rec := r.choices[r.next]
	r.next++

	if rec.Point != ask.Point {
		return 0, fmt.Errorf("%w: choice %d is %q, engine asked %q",
			ErrReplayDiverged, r.next, rec.Point, ask.Point)
	}

	if !slices.Equal(rec.Options, ask.Options) {
		return 0, fmt.Errorf("%w: choice %d (%s) recorded options %v, engine offered %v",
			ErrReplayDiverged, r.next, ask.Point, rec.Options, ask.Options)
	}

	if rec.Chosen < 0 || rec.Chosen >= len(ask.Options) {
		return 0, fmt.Errorf("%w: choice %d (%s) recorded index %d of %d options",
			ErrReplayDiverged, r.next, ask.Point, rec.Chosen, len(ask.Options))
	}

	return rec.Chosen, nil
}

// Remaining reports how many recorded choices have not been reapplied. A
// replay that ends with any left over reproduced a shorter lifepath than
// the record holds, which is a divergence the caller has to see.
func (r *Replay) Remaining() int {
	return len(r.choices) - r.next
}
