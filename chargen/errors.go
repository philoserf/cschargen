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

	// ErrNoCareer is a term served with no career entered. It is a bug in
	// the engine, not a rule.
	ErrNoCareer stringError = "no career to serve a term in"

	// ErrUnknownEffect is an effect kind the engine does not handle, which
	// means career/ grew a kind and chargen/ did not.
	ErrUnknownEffect stringError = "unknown effect kind"

	// ErrMalformedCheck is a check effect carrying no target, which is a
	// transcription error rather than a rule.
	ErrMalformedCheck stringError = "check effect carries no target"

	// ErrBadExpression is a dice expression in a table the engine cannot
	// read. Like ErrMalformedCheck it is a transcription error: the engine
	// errors rather than guessing what the page meant.
	ErrBadExpression stringError = "table carries a dice expression the engine cannot read"

	// ErrMissingEventRow is a d66 result with no row, which the career
	// package's completeness test exists to prevent.
	ErrMissingEventRow stringError = "no event row for that d66 result"
)
