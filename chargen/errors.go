package chargen

// Sentinel errors are const, not var: they are declared as the unexported
// stringError type so that they cannot be reassigned.
type stringError string

func (e stringError) Error() string { return string(e) }

// The engine's sentinels.
const (
	// ErrNoOptions is a choice point offered with an empty option list. It
	// is a bug in the engine rather than a rule: every choice the book
	// prints has at least one answer.
	ErrNoOptions stringError = "choice point offered no options"

	// ErrReplayExhausted is a replay that ran past the end of the record's
	// choices -- the engine asked a question the record does not answer.
	ErrReplayExhausted stringError = "replay ran out of recorded choices"

	// ErrReplayDiverged is a replay whose recorded choice does not answer
	// the question the engine asked.
	ErrReplayDiverged stringError = "replay diverged from the record"

	// ErrChoiceOutOfRange is a decider that answered with an index the
	// option list cannot hold. That is a decider answering wrongly, which
	// is distinct from one declining to answer.
	ErrChoiceOutOfRange stringError = "decider chose an option that was not offered"
)
