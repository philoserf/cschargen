// Package career holds the career definitions and the two tables every
// career shares: Injury (p. 119) and Life Events (p. 120).
//
// The definitions are hand-typed Go rather than a data file. The format is
// proven against all thirty-four careers, whose shapes genuinely
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
	// carrying a modifier (ERRATA E-3). ForfeitAll removes every roll
	// accrued in the career so far rather than a counted number, and
	// CashOnly restricts the rolls granted to the Cash column.
	EffectBenefitRolls

	// EffectRelationship gains or loses an Ally, Contact, Rival or Enemy.
	EffectRelationship

	// EffectCredits pays the character.
	EffectCredits

	// EffectStash adds named things to what a character owns, removes every
	// one of a name, or empties the lot. Item names it, Count is how many,
	// Dice carries the value where the book prints one, and Lose says the
	// result takes them rather than gives them.
	EffectStash

	// EffectModifier attaches a modifier to a future throw -- the next
	// advancement roll, the next enlistment attempt.
	EffectModifier

	// EffectAdvance is an automatic advancement, granted without a throw.
	EffectAdvance

	// EffectRank moves the character's rank without an advancement throw.
	// Levels says how far and which way: a promotion is positive and a
	// demotion negative, and zero means one rank up, which is what every
	// bare promotion result asks for.
	EffectRank

	// EffectContinue forces the character to serve another term in this
	// career.
	EffectContinue

	// EffectChangeAssignment moves the character to another assignment in
	// the same career.
	EffectChangeAssignment

	// EffectChooseCareer ends the current career without naming the next
	// one: "choose another career and a new homeworld" (Colonist mishaps
	// 3-6 and 8, p. 174). ERRATA E-6 -- distinct from EffectTransfer,
	// because the character enlists in whatever they choose rather than
	// being placed there.
	EffectChooseCareer

	// EffectTransfer sends the character to another career by name. The
	// destination is resolved when the transfer is taken; a name no career
	// answers to ends generation rather than being skipped.
	EffectTransfer

	// EffectNewHomeworld reassigns the homeworld, which six of Colonist's
	// eleven mishaps do. The new world is thrown for as a birth world is,
	// and changes the tech level and nothing else (ERRATA E-9).
	EffectNewHomeworld

	// EffectRollTable rolls on one of the career's skill tables by name,
	// which several events ask for: "Make a roll on the Advanced Education
	// table", "Choose an Assignment skill table for an assignment other
	// than your own and roll once".
	EffectRollTable

	// EffectAutoSuccess makes a named future throw succeed without rolling:
	// "Gain an automatic success on your next Survival roll" (p. 176). It is
	// not a large modifier -- a modifier can still fail, and the book says
	// the throw succeeds.
	EffectAutoSuccess

	// EffectMilitaryEvent sends the character to the Military Events table
	// (p. 121), which the military careers' d66 tables reach at 41-46.
	EffectMilitaryEvent

	// EffectCommission attempts a commission outside the ordinary Step 14
	// offer, which the Military Events table's tenth result does: "If you
	// are enlisted, roll to be Commissioned and take +2 DM. If you are
	// currently an Officer, roll once on the Officer Skills Table and take
	// a +2 DM to your next Advancement roll" (p. 121).
	EffectCommission

	// EffectRating moves a tie's Relationship Rating (p. 320): "an existing
	// Ally or Contact loses 1d6 x 20 Relationship Rating" (Life Event 5),
	// "decrease your Relationship Rating with that parent by 150" (Youth
	// Path 1 result 5). Target says which tie, Relationship narrows it to
	// one kind, and Modifier or Dice carries the amount.
	EffectRating

	// EffectLoseTie removes ties outright: "lose one Contact or Ally"
	// (Teenage Life Event 3), "lose an Ally, Contact, Rival or Enemy in
	// that order" (Life Event 4), "lose 1D3 Allies and Contacts". Order
	// lists the kinds to try, in the order to try them; Count or CountDice
	// is how many to take; ThisCareer narrows them to this career's.
	EffectLoseTie

	// EffectAge moves the character's age by a number of years outside the
	// four a term takes: "you have lost a year of your life", which one
	// result prints.
	EffectAge

	// EffectSentence lengthens or shortens a forced spell in a career:
	// "add one term to your sentence", "your sentence is reduced by one
	// term; if you have one term or less remaining, you are released".
	EffectSentence

	// EffectCondition records something a character carries that is
	// neither a skill, a characteristic nor a possession: an addiction,
	// a religion. Thirteen results give the first and nine the second.
	//
	// The engine names it and says no more. What an addiction costs a
	// character is the referee's, and p. 152's own "an addiction to
	// alcohol or a drug of your choice" prints no characteristic loss
	// beside it -- ERRATA E-38.
	EffectCondition

	// EffectSubTable is a 1d6 table printed inside a result rather than as
	// a table of its own: "Roll 1d6. On a 1 ... on a 2-5 ... on a 6 ...".
	// Forty-two results across the corpus carry one.
	//
	// It is a roll rather than a choice, which is what distinguishes it
	// from [EffectChoice]: the character has no say in which row comes up.
	// Sub holds the rows, and every result of a d6 falls in exactly one --
	// TestEverySubTableCoversTheDie is what says so.
	EffectSubTable

	// EffectPool grants a modifier the character spends themselves, a
	// little at a time: "gain a modifier of +6 which must be split up into
	// increments of not more than +2 to use on any Survival or Advancement
	// rolls until it is depleted" (Journalist event 66), "+6 to be used
	// over the remainder of your career in this service which you can
	// break into increments of up to +3 at any time" (Diplomatic Service).
	//
	// Modifier is the whole of it and Uses the most that may go on one
	// throw; Applies names the throws it may be spent on, and
	// WhileInThisCareer ends it when the career does.
	EffectPool

	// EffectOnFailure attaches effects to the failure of a named throw:
	// "if you fail that Advancement roll, you lose the Ally and take a -2
	// DM to your next Advancement roll", which nine results print. It is
	// not an EffectCheck, because the throw is the career's own rather
	// than one the result rolls.
	EffectOnFailure

	// EffectAutoFailure is the mirror of [EffectAutoSuccess]: a named
	// throw fails without being rolled. Two results print it -- "you are
	// suspended, and your next Advancement roll fails automatically".
	EffectAutoFailure

	// EffectRaiseHeld raises a skill the character already has, without
	// naming which: "Gain one level in any skill you already possess",
	// which thirty-one results across the corpus print. The choice is
	// among what the character holds, so the engine cannot know the
	// options until it is applied.
	EffectRaiseHeld

	// EffectAnySkill grants a level in any skill at all, chosen from the
	// list of pp. 304-314: "gain a level in any skill of your choice".
	// NotHeld narrows it to the skills the character does not have, which
	// one result and both undergraduate degrees ask for.
	EffectAnySkill

	// EffectHomeworldSkill grants a level in a specialty of one skill that
	// the homeworld's own background skills name: "gain a level in a
	// Survival specialty which is used on your homeworld" (Youth Path 1
	// result 2, and four more like it). Skill says which skill's
	// specialties to look through.
	EffectHomeworldSkill

	// EffectLoseSkill removes a skill outright: "You lose all levels of
	// Suit (Vacc Suit)" (Orbital Construction mishap 6).
	EffectLoseSkill

	// EffectBecome changes what a tie is rather than what it is worth: "If
	// you currently have Enemies, one of those is now a Rival" (p. 91),
	// "one Contact or Ally from this career becomes an Enemy". From lists
	// the kinds eligible, Relationship is what they become, Count or
	// CountDice how many, and ThisCareer narrows them to the ties this
	// career gave. Fallback is what happens when there are none.
	EffectBecome

	// EffectImproveTies is the upgrade three life-event tables print in
	// the same words (pp. 91, 96, and the teenage table of p. 85): "If you
	// have no Contacts, then you will gain one Contact. If you currently
	// have Enemies, one of those is now a Rival. If you have Rivals, one of
	// those is now a Contact. If you have Contacts, one of those is now an
	// Ally."
	//
	// It is its own kind rather than four EffectBecomes in a row because
	// the four clauses are read against the state as it was: running them
	// in sequence would let one NPC climb three bands on a result the book
	// calls "an improvement to a relationship". ERRATA E-34.
	EffectImproveTies

	// EffectYouthLifeEvent sends the character to the Youth Life Events
	// table (p. 75), which result 10 of every youth path reaches. It is
	// distinct from EffectLifeEvent because the two tables are different:
	// the youth one is a d6 of six rows, the career one a 2d6 of eleven.
	EffectYouthLifeEvent

	// EffectTeenageLifeEvent sends the character to the Teenage Life Events
	// table (p. 85), which result 10 of every teenage path reaches.
	EffectTeenageLifeEvent

	// EffectCollegiateLifeEvent and EffectAcademyLifeEvent send the
	// character to an institution's own life-events table (pp. 91, 96),
	// which result 10 of its events table reaches. They are separate kinds
	// because the two tables differ: the academy's result 2 is a battle.
	EffectCollegiateLifeEvent
	EffectAcademyLifeEvent

	// EffectHonors grants honours to a character whose throw for them
	// failed, which one result on each institution's events table does:
	// "If you failed on your Honors roll, you have now been successful"
	// (p. 91).
	EffectHonors

	// EffectFreed ends a character's enslavement: "Your owner has decided
	// to set you free. Continue your character as a free altrant or uplift"
	// (pp. 75, 84, and the slave career's own event 66).
	EffectFreed

	// EffectGroup applies several effects as one, which the benefit tables
	// need because a row carries a single Other effect and several rows do
	// two things: a Producer credit pays at once and then pays yearly, and
	// a mishap that takes the benefit rolls also takes the shares.
	EffectGroup

	// EffectEducationLockout closes every higher-learning institution for a
	// number of terms: "you may not apply to any institute of higher
	// learning for the next two terms" (pp. 98, 102). It is the education
	// side of the enlistment lockout a refused career leaves behind.
	EffectEducationLockout

	// EffectUnimplemented is a result the engine cannot carry out. It
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

// SubRow is one row of a 1d6 table printed inside a result. From and To
// are the inclusive range of the die it covers, which is one number on most
// rows and a span on the rest: "on a 2-5".
type SubRow struct {
	From    int
	To      int
	Summary string
	Effects []Effect
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

// TieTarget says which ties an [EffectRating] moves.
type TieTarget string

const (
	// TargetOne moves a single tie, chosen at the choice point. The
	// Relationship field narrows the candidates where the page does:
	// "an existing Ally or Contact".
	TargetOne TieTarget = "one"

	// TargetAll moves every tie the character has: "Lower the Relationship
	// Rating with everyone in your life by 25" (Youth Path 1 result 7).
	TargetAll TieTarget = "all"

	// TargetFamily moves every tie that came from Step 5: "Gain +20 to the
	// Relationship Rating of all family members" (Youth Path 1 result 16).
	TargetFamily TieTarget = "family"

	// TargetRole moves one tie of a named family role: "Choose one of your
	// parents and that parent will leave your life" (Youth Path 1 result
	// 5). The Role field says which.
	TargetRole TieTarget = "role"

	// TargetAllOfRole moves every tie of a named role: "increase your
	// Relationship Rating with both parents by 25" (Youth Path 1 result
	// 20).
	TargetAllOfRole TieTarget = "all of role"
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
	Dice           string // a rolled amount or count like "1d6x100", "1d3"

	Options []Option
	Check   *Target
	Success []Effect
	Failure []Effect

	Times int

	// Count is how many: benefit rolls, possessions, ties. Where a result
	// names no number the default differs by kind, because the pages do:
	// an [EffectLoseTie] or an [EffectBecome] with no count takes one
	// ("lose one Contact or Ally"), and an [EffectRating] with no count
	// moves every tie its Target names ("lower the Relationship Rating
	// with everyone in your life by 25").
	Count        int
	Modifier     int
	Scope        BenefitScope
	Relationship Relationship
	Item         string

	// ForfeitAll, CashOnly and Immediate belong to [EffectBenefitRolls].
	//
	// ForfeitAll is the sixty-seven results that take back everything the
	// character has earned in this career -- "you are dismissed and lose
	// all Benefits". It is a flag rather than a large negative Count
	// because the number is not knowable when the table is written: it
	// depends on how many terms have been served when the result fires.
	//
	// CashOnly restricts the granted rolls to the Cash column, which the
	// "two Cash Benefit rolls" results ask for. Immediate takes them at
	// once instead of queueing them for Step 19, and RerollNothing re-rolls
	// a row that pays nothing, which one such result prints.
	ForfeitAll    bool
	CashOnly      bool
	Immediate     bool
	RerollNothing bool

	// Lose turns an [EffectStash] from a grant into a removal: every
	// possession named Item goes, which is what "lose any Company Shares"
	// asks for.
	Lose bool

	// CountDice is how many, where the book rolls for it rather than
	// printing a number: "1D6 Company Shares". Dice on an [EffectStash] is
	// the value of one of them, which is a different question.
	CountDice string

	// Group is the list an [EffectGroup] applies, in order.
	Group []Effect

	// Career and Assignment name an [EffectTransfer] destination; Terms is
	// how many the character owes there, zero meaning "until they leave".
	Career     string
	Assignment string
	Terms      int

	// Target and Order belong to [EffectRating] and [EffectLoseTie]: which
	// ties are moved, and the order in which kinds are tried for removal.
	Target TieTarget
	Order  []Relationship

	// From is the kinds an [EffectBecome] may change, in the order to try
	// them. ThisCareer narrows an [EffectLoseTie] or an [EffectBecome] to
	// the ties this career gave, which is what "every Contact made in this
	// career" means. Fallback is what a result does when it finds nobody:
	// "with no Ally or Contact to lose, gain an Enemy at -110 instead".
	From       []Relationship
	ThisCareer bool
	Fallback   []Effect

	// Levels is how many ranks an [EffectRank] moves, and which way. Zero
	// is one rank up.
	Levels int

	// NotHeld narrows an [EffectAnySkill] to the skills the character does
	// not already have: "any skill at level 1 which you do not already
	// possess".
	NotHeld bool

	// LoseTheRest belongs to [EffectBecome]: "one Contact or Ally from this
	// career becomes an Enemy, and every other relationship gained here is
	// lost" (Gambler mishap 11). The one that changed survives, which is
	// why the two halves cannot be written as two effects in a row.
	LoseTheRest bool

	// Newest takes the most recently gained tie rather than the first,
	// which is what a result meaning "that Ally" needs: the nine hooks
	// that read "if you fail that Advancement roll, you lose the Ally"
	// mean the one the event granted a line earlier. ERRATA E-39.
	Newest bool

	// ExcludeFamily belongs to [EffectRating] and [EffectLoseTie]: "1D6-2
	// (minimum 1) of your Contacts or Allies who are not family members
	// lose 50 Relationship Rating" (Teenage Path 3 result 4), and the
	// hooks of E-39, which mean a colleague rather than a parent.
	//
	// Minimum belongs to [EffectRating] alone.
	ExcludeFamily bool
	Minimum       int

	// Role narrows a [TargetRole] to one kind of relative: "parent",
	// "sibling", and the rest of Step 5's roles.
	Role string

	// AtLevelZero makes an [EffectSkill] grant the skill at level 0 rather
	// than raising it by Level: "Gain Streetwise 0 or Admin 0" (Youth
	// Path 1 result 5, p. 68), "gain Medic at level 0" (Belter event 24,
	// p. 162).
	//
	// It is a flag rather than a Level of zero because Level zero means
	// "the result did not say", which is one level -- every ordinary result
	// reads "gain a level in X".
	AtLevelZero bool

	// Rating is the Relationship Rating an [EffectRelationship] starts a
	// tie at, where the result that grants it names one. Zero means the
	// result did not, and the engine supplies a default (ERRATA E-23).
	Rating int

	// OnTags, NotTags, OnCharacteristics and Standing narrow an
	// [EffectModifier] that attaches to an enlistment throw.
	//
	// Twelve results modify an enlistment by the class of career being
	// entered rather than by name: "-4 DM to enlist in any government
	// related career", "-2 DM to enter any non-criminal career", "-2 DM to
	// every enlistment roll that involves EDU". OnTags is the classes it
	// applies to and NotTags the classes it does not;
	// OnCharacteristics is the enlistment characteristics it applies to.
	// An empty set means the modifier applies to every career, which is
	// the ordinary case.
	//
	// Standing distinguishes "every career after this one" from "the next
	// career": a standing modifier is not consumed by the throw it
	// modifies. Several results print each, and the difference is the
	// rest of a character's life.
	// OnCareers and NotCareers name careers rather than classes, which two
	// results do: "+2 DM to enlistment in the Sports career", "-2 DM to
	// enter any career other than Vagabond".
	//
	// OnSkill narrows a check modifier to one skill: "+1 DM to Melee
	// checks in this career". WhileInThisCareer ends a modifier when the
	// career does, which is what "for the rest of this career" means.
	OnTags            []Tag
	NotTags           []Tag
	OnCareers         []string
	NotCareers        []string
	OnCharacteristics []string
	OnSkill           string
	WhileInThisCareer bool
	Standing          bool

	// Years is how far an [EffectAge] moves the character's age, and Terms
	// how far an [EffectSentence] moves a sentence. Both carry their sign.
	Years int

	// Condition names what an [EffectCondition] records.
	Condition string

	// Sub is the rows of an [EffectSubTable], in the order the page prints
	// them. Dice names the throw where it is not 1d6 -- two results roll
	// 2d6 -- and Characteristic adds that characteristic's modifier to it,
	// which one of them does: "Roll 2d6 and add your DEX bonus".
	Sub []SubRow

	// Spendable is the throws an [EffectPool] may be spent on. It is a
	// list because both results that grant one name two: "any Survival or
	// Advancement rolls".
	Spendable []string

	// Uses is how many throws of that name an [EffectModifier] survives:
	// "take a -2 DM on your next two Advancement rolls", "-2 DM to all
	// Advancement rolls for the remainder of this career". Zero is one
	// throw, which is what most results grant. Standing is the same idea
	// without a count, and the two do not combine.
	Uses int

	// Applies names the throw an [EffectModifier] attaches to.
	Applies string

	// Table names the skill table an [EffectRollTable] rolls on, and
	// OtherAssignment asks for an assignment other than the character's
	// own, which two events do.
	Table           SkillTableKind
	OtherAssignment bool
}
