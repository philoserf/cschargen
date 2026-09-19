// Package career holds the career definitions and the two tables every
// career shares: Injury (p. 119) and Life Events (p. 120).
//
// The definitions are hand-typed Go rather than a data file. The format is
// what milestone 3 proves against thirty-four careers whose shapes genuinely
// differ -- National Navy has nine skill tables and a commission, Vagabond
// has four and no enlistment throw -- and designing a schema against three
// and then meeting the variance would mean migrating data twice. Go structs
// make the variance a compile error instead.
//
// This package is consulted, never rolled: it brings no dice. The engine
// reads a table and applies what it finds.
package career

// EffectKind discriminates what a table result does to a character.
//
// The vocabulary is closed on purpose. An event the engine cannot carry out
// is EffectUnimplemented carrying what the book asked for, so that the
// record says what it could not do rather than quietly doing something
// else -- which is what makes a stubbed career transfer honest rather than
// silent.
type EffectKind int

// The effects a table result can have.
const (
	// EffectSkill grants a level in a skill. Specialties lists the
	// specialties to choose between where the book offers a choice, and is
	// empty where it names one or none.
	EffectSkill EffectKind = iota

	// EffectCharacteristic moves a characteristic by Delta.
	EffectCharacteristic

	// EffectChoice is "gain a level in A or B": the Decider picks one of
	// Options, each of which is itself a list of effects.
	EffectChoice

	// EffectCheck is "roll X 8+. If you succeed ... If you fail ...". The
	// two branches are effect lists of their own, so a check inside a
	// branch nests without a special case.
	EffectCheck

	// EffectInjury sends the character to the Injury table (p. 119), Times
	// times.
	EffectInjury

	// EffectLifeEvent sends the character to the Life Events table
	// (p. 120), which every career's d66 reaches at 31-36.
	EffectLifeEvent

	// EffectMishapNoEject rolls on the career's own mishap table without
	// ejecting the character, which is what an event result of 11 does in
	// several careers.
	EffectMishapNoEject

	// EffectBenefitRolls grants or removes mustering-out rolls, optionally
	// carrying a modifier (ERRATA E-3).
	EffectBenefitRolls

	// EffectRelationship gains or loses an Ally, Contact, Rival or Enemy.
	EffectRelationship

	// EffectCredits pays the character.
	EffectCredits

	// EffectStash adds a named item to the character's stash, or empties it.
	EffectStash

	// EffectModifier attaches a modifier to a future throw -- the next
	// advancement roll, the next enlistment attempt.
	EffectModifier

	// EffectAdvance is an automatic advancement, granted without a throw.
	EffectAdvance

	// EffectRank raises the character's rank by one without an advancement
	// throw.
	EffectRank

	// EffectContinue forces the character to serve another term in this
	// career.
	EffectContinue

	// EffectChangeAssignment moves the character to another assignment in
	// the same career.
	EffectChangeAssignment

	// EffectTransfer sends the character to another career by name --
	// Vagabond and Prisoner are the two the rules force, and the two this
	// milestone implements.
	EffectTransfer

	// EffectNewHomeworld reassigns the homeworld, which six of Colonist's
	// eleven mishaps do. Until the setting data lands there is nothing to
	// reassign it to, so the engine records the demand and moves on.
	EffectNewHomeworld

	// EffectUnimplemented is a result this milestone cannot carry out. It
	// carries the book's demand in Detail.
	EffectUnimplemented
)

// Target is what a check is rolled against: a characteristic or a skill,
// and the number to meet.
type Target struct {
	// Characteristic is the abbreviation the book prints -- STR, DEX, END,
	// INT, EDU, CHA. Empty when the check is against a skill.
	Characteristic string

	// Skill is the skill checked, where the book names one ("Roll Melee
	// (Any) 8+"). Empty when the check is against a characteristic.
	Skill string

	// Number is the target to meet or exceed.
	Number int
}

// Option is one branch of an [EffectChoice].
type Option struct {
	Label   string
	Effects []Effect
}

// Relationship is one of the four kinds of NPC tie a character accumulates.
type Relationship string

// The four, in the order Life Event 4 lists them for removal (p. 120).
const (
	Ally    Relationship = "ally"
	Contact Relationship = "contact"
	Rival   Relationship = "rival"
	Enemy   Relationship = "enemy"
)

// BenefitScope says how far an [EffectBenefitRolls] modifier reaches.
// ERRATA E-3: the book grants both "a +1 modifier to all Benefit rolls made
// in this career" and "three Benefit rolls at +1", and the two are not the
// same thing.
type BenefitScope int

// The two scopes.
const (
	// ScopeBatch modifies only the rolls this effect granted.
	ScopeBatch BenefitScope = iota

	// ScopeCareer modifies every benefit roll made in this career,
	// including ones already granted.
	ScopeCareer
)

// Effect is one thing a table result does. Which fields carry meaning
// depends on Kind; the constants above say which.
type Effect struct {
	Kind EffectKind

	// Detail is prose for the record: what the book asked for, in this
	// engine's words. Every effect carries one, because an event's meaning
	// is rarely recoverable from its mechanical parts alone.
	Detail string

	Skill       string
	Specialties []string
	Level       int

	Characteristic string
	Delta          int
	Dice           string // a credit or characteristic amount like "1d6x100"

	Options []Option
	Check   *Target
	Success []Effect
	Failure []Effect

	Times        int
	Count        int
	Modifier     int
	Scope        BenefitScope
	Relationship Relationship
	Item         string

	// Career and Assignment name an [EffectTransfer] destination; Terms is
	// how many the character owes there, zero meaning "until they leave".
	Career     string
	Assignment string
	Terms      int

	// Applies names the throw an [EffectModifier] attaches to.
	Applies string
}
