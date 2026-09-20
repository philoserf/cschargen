// Package chargen is the character-generation engine. This file defines the
// generation record: the ordered event log of every step, throw, choice and
// consequence in a lifepath (docs/PRD.md FR15).
//
// Events carry monotonic sequence numbers starting at 1 and increasing by
// 1. Consequence events reference the sequence number of the throw, choice
// or step that caused them. The stored log is verification data for replay:
// the engine re-runs from the seed and the choice events and recomputes
// every throw (docs/PRD.md, Replay and provenance contract).
package chargen

import "github.com/philoserf/cschargen/dice"

// EventKind discriminates the four event kinds.
type EventKind string

// The four event kinds: a step of the twenty-step process entered, a throw,
// a choice, and a consequence of one of them.
const (
	EventStep        EventKind = "step"
	EventThrow       EventKind = "throw"
	EventChoice      EventKind = "choice"
	EventConsequence EventKind = "consequence"
)

// Event is one entry in the generation record. Exactly one payload field is
// non-nil, matching Kind.
type Event struct {
	Seq  int       `json:"seq"`
	Kind EventKind `json:"kind"`

	Step        *StepEvent        `json:"step,omitempty"`
	Throw       *ThrowEvent       `json:"throw,omitempty"`
	Choice      *ChoiceEvent      `json:"choice,omitempty"`
	Consequence *ConsequenceEvent `json:"consequence,omitempty"`
}

// StepEvent records entering one of the twenty steps of character creation
// (pp. 16-130).
type StepEvent struct {
	Name string `json:"name"`
	Cite string `json:"cite"`
}

// ThrowEvent records one throw. A throw resolved against no target -- the
// six characteristic rolls of p. 13, say -- leaves Target and Success nil,
// because zero is a target no throw of dice satisfies and recording one
// would read as a check that always failed.
type ThrowEvent struct {
	Expr    string     `json:"expr"`
	Dice    []int      `json:"dice"`
	Total   int        `json:"total"`
	Mods    []dice.Mod `json:"mods,omitempty"`
	Target  *int       `json:"target,omitempty"`
	Success *bool      `json:"success,omitempty"`
	Cite    string     `json:"cite"`
}

// ChoiceEvent records one resolved choice point: who decided, what the
// alternatives were, and which was taken. Choice events are replay input:
// re-running the engine reapplies them.
type ChoiceEvent struct {
	Decider DeciderKind `json:"decider"`
	Point   string      `json:"point"`
	Prompt  string      `json:"prompt"`
	Options []string    `json:"options"`
	Chosen  int         `json:"chosen"` // index into Options
	Cite    string      `json:"cite"`
}

// ConsequenceKind discriminates consequence payloads. Kinds are added as
// the engine grows; each names what changed about the character.
type ConsequenceKind string

// The consequence kinds. ConsequenceUnimplemented is the honest one: a
// table result the engine cannot carry out, recorded with what the book
// asked for, so the record says what it could not do rather than quietly
// doing something else (docs/MILESTONE-1.md).
const (
	ConsequenceCharacteristic ConsequenceKind = "characteristic"
	ConsequenceSkill          ConsequenceKind = "skill"
	ConsequenceRelationship   ConsequenceKind = "relationship"
	ConsequenceCredits        ConsequenceKind = "credits"
	ConsequenceStash          ConsequenceKind = "stash"
	ConsequenceBenefitRolls   ConsequenceKind = "benefit_rolls"
	ConsequenceInjury         ConsequenceKind = "injury"
	ConsequenceRank           ConsequenceKind = "rank"
	ConsequenceCareer         ConsequenceKind = "career"
	ConsequenceHomeworld      ConsequenceKind = "homeworld"
	ConsequenceFamily         ConsequenceKind = "family"
	ConsequenceEducation      ConsequenceKind = "education"
	ConsequenceSpecies        ConsequenceKind = "species"
	ConsequenceFinishing      ConsequenceKind = "finishing"
	ConsequenceAge            ConsequenceKind = "age"
	ConsequenceModifier       ConsequenceKind = "modifier"
	ConsequenceUnimplemented  ConsequenceKind = "unimplemented"
)

// ConsequenceEvent records one change to the character. Cause is the
// sequence number of the throw or choice that produced it, or of the step
// that established the state where no throw or choice did.
//
// The payload fields are deliberately ragged: a kind uses the ones it
// needs. Which fields each kind uses is documented on the kinds above
// rather than in the type, because a struct per kind would be thirteen
// types that differ by two fields each.
type ConsequenceEvent struct {
	Kind   ConsequenceKind `json:"kind"`
	Cause  int             `json:"cause"`
	Detail string          `json:"detail"`

	Characteristic string `json:"characteristic,omitempty"`
	Skill          string `json:"skill,omitempty"`
	Career         string `json:"career,omitempty"`
	Delta          int    `json:"delta,omitempty"`
	Level          int    `json:"level,omitempty"`
	Permanent      bool   `json:"permanent,omitempty"`
	Cite           string `json:"cite"`
}

// Log is the ordered event log. The zero value is an empty log ready to
// append to.
type Log struct {
	events []Event
}

// Events returns the log in order. The slice is the log's own: callers
// read it, and every append goes through the methods below so that the
// sequence numbering has one owner.
func (l *Log) Events() []Event {
	return l.events
}

// Len reports how many events the log holds, which is also the sequence
// number of the last one.
func (l *Log) Len() int {
	return len(l.events)
}

// Step records entering a step of character creation.
func (l *Log) Step(name, cite string) int {
	return l.append(Event{Kind: EventStep, Step: &StepEvent{Name: name, Cite: cite}})
}

// Roll records a throw resolved against no target.
func (l *Log) Roll(r dice.Roll, cite string) int {
	return l.append(Event{Kind: EventThrow, Throw: &ThrowEvent{
		Expr:  r.Expr,
		Dice:  r.Dice,
		Total: r.Total,
		Cite:  cite,
	}})
}

// Throw records a throw resolved against a target number.
func (l *Log) Throw(t dice.Throw, cite string) int {
	target, success := t.Target, t.Success

	return l.append(Event{Kind: EventThrow, Throw: &ThrowEvent{
		Expr:    t.Roll.Expr,
		Dice:    t.Roll.Dice,
		Total:   t.Total,
		Mods:    t.Mods,
		Target:  &target,
		Success: &success,
		Cite:    cite,
	}})
}

// Choice records a resolved choice point.
func (l *Log) Choice(c ChoiceEvent) int {
	return l.append(Event{Kind: EventChoice, Choice: &c})
}

// Consequence records one change to the character.
func (l *Log) Consequence(c ConsequenceEvent) int {
	return l.append(Event{Kind: EventConsequence, Consequence: &c})
}

// append assigns the next sequence number and returns it, so that a caller
// can reference the event it just wrote as the cause of what follows.
func (l *Log) append(e Event) int {
	e.Seq = len(l.events) + 1
	l.events = append(l.events, e)

	return e.Seq
}
