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

	// ErrNoSetting is a run with no setting data. Steps 3 and 4 have
	// nothing to read without it, and a character with no homeworld has no
	// background skills, no language, no tech level and no term cap.
	ErrNoSetting stringError = "no setting data; pass --data or use the sample"

	// ErrNoSubsectors and ErrNoWorlds are setting data with nothing to be
	// from. The validator refuses both, so reaching one means the engine
	// was handed data that never went through it.
	ErrNoSubsectors stringError = "the setting data has no subsectors"
	ErrNoWorlds     stringError = "the subsector has no worlds"

	// ErrNoLanguage is a world with no primary language, which the
	// validator also refuses.
	ErrNoLanguage stringError = "the homeworld lists no primary language"

	// ErrOriginThrowsExhausted is the chart of p. 39 thrown the maximum
	// number of times without landing on an entry, which is the dice
	// rather than the data: any file with a rollable subsector lands
	// eventually.
	ErrOriginThrowsExhausted stringError = "the subsector chart was thrown its limit without landing"

	// ErrNoSuchSubsector and ErrNoSuchHomeworld are an origin asked for by
	// name that the setting data does not have. The command checks both
	// before generating, so reaching one here means the engine was driven
	// directly.
	ErrNoSuchSubsector stringError = "the setting data has no such subsector"
	ErrNoSuchHomeworld stringError = "the setting data has no such homeworld"

	// ErrMissingEventRow is a d66 result with no row, which the career
	// package's completeness test exists to prevent.
	ErrMissingEventRow stringError = "no event row for that d66 result"
)

// ErrUnknownSpecies is a character asked for a species the setting data
// file does not declare. The engine holds no species of its own, so there
// is nothing to fall back to.
const ErrUnknownSpecies = stringError(
	"species: not declared in the setting data")

// ErrNoHomeworldAdmitsThem is a species no world in the setting data will
// have. p. 42 says to choose another world; where there is none, there is
// no character to generate.
const ErrNoHomeworldAdmitsThem = stringError(
	"no world in the setting data admits this species")

// ErrPlayerGone is an interactive session whose input ended before the
// character did. It is the abandoned session the Decider interface names,
// and it ends generation rather than falling back on a default.
const ErrPlayerGone = stringError("the player stopped answering")
