package career

import (
	"strconv"
	"strings"
)

// The constructors below exist so that a transcribed table reads as close
// to the page as Go allows. They are the only way effects are built in this
// package: a table written with struct literals would bury the one field
// that differs between two rows in a wall of field names.

// skill grants one level in a skill. Specialties are the choices the book
// offers where it writes "(Any)" over a named list, or nothing where it
// names the specialty itself.
func skill(name string, specialties ...string) Effect {
	detail := "gain a level in " + name

	return Effect{
		Kind:        EffectSkill,
		Detail:      detail,
		Skill:       name,
		Specialties: specialties,
		Level:       1,
	}
}

// skillAt is skill at a level the page names above 1, which only rank
// tables do: Fringe Marketer's rank 5 prints "Streetwise 2" over a rank 0
// that already printed "Streetwise 1" (p. 199).
//
// p. 116 makes a rank benefit a floor rather than an increment -- "If a
// Benefit grants a skill the character already possesses at the listed level
// or higher, no additional benefit is gained". The data here is the page;
// grantRankSkill applies the floor.
func skillAt(name string, level int, specialties ...string) Effect {
	granted := skill(name, specialties...)

	granted.Detail = "gain " + name + " " + itoa(level)
	granted.Level = level

	return granted
}

// skillZero is "Gain Streetwise 0": the skill at level 0, which the youth
// and teenage tables grant where a career table grants a level.
//
// Level 0 is not nothing. A character who holds a skill at 0 makes its
// checks without the -3 penalty for having no training at all, which is
// what the pre-career steps are handing out.
func skillZero(name string, specialties ...string) Effect {
	granted := skill(name, specialties...)

	granted.Detail = "gain " + name + " at level 0"
	granted.Level = 0
	granted.AtLevelZero = true

	return granted
}

// raiseHeldSkill is the thirty-one results that read "gain one level in
// any skill you already possess". The engine cannot name the options here:
// they are whatever the character has when the result fires.
func raiseHeldSkill() Effect {
	return Effect{Kind: EffectRaiseHeld, Detail: "raise a skill the character already holds"}
}

// raiseHeldSkillIn is raiseHeldSkill narrowed to one skill's specialties:
// "gain a level in the Science skill you used on this task", where the
// check was made on whichever Science the character has most of.
func raiseHeldSkillIn(name string) Effect {
	return Effect{
		Kind:    EffectRaiseHeld,
		Detail:  "raise the " + name + " skill the character used",
		OnSkill: name,
	}
}

// anySkill is "gain a level in any skill of your choice", chosen from the
// list of pp. 304-314.
func anySkill() Effect {
	return Effect{Kind: EffectAnySkill, Detail: "a level in any skill the character chooses"}
}

// anySkillNotHeld is anySkill narrowed to what the character lacks: "any
// skill at level 1 which you do not already possess".
func anySkillNotHeld() Effect {
	return Effect{
		Kind:    EffectAnySkill,
		Detail:  "any skill at level 1 the character does not already hold",
		NotHeld: true,
	}
}

// gainRank is a promotion granted without an advancement throw, which
// several events and the Military Events table award.
func gainRank() Effect {
	return Effect{Kind: EffectRank, Detail: "gain a rank", Levels: 1}
}

// loseRank is a demotion. Sixteen results across the corpus print one, two
// of them adding "retaining any benefit already gained" -- which is what
// the other fourteen do as well, for the reason ERRATA E-36 gives.
func loseRank(levels int) Effect {
	detail := "lose one rank"
	if levels > 1 {
		detail = "lose " + itoa(levels) + " ranks"
	}

	return Effect{Kind: EffectRank, Detail: detail, Levels: -levels}
}

// homeworldSpecialty is the five results that send the character back to
// where they grew up for a specialty: "a Survival specialty which is used
// on your homeworld". Which specialties those are is in the setting data,
// not in this package.
//
//nolint:unparam // the skill is the page's; Survival is merely the only one any result names
func homeworldSpecialty(name string) Effect {
	return Effect{
		Kind:   EffectHomeworldSkill,
		Detail: "a level in a " + name + " specialty the homeworld's background skills name",
		Skill:  name,
	}
}

// loseSkill removes every level of a skill.
func loseSkill(name, specialty string) Effect {
	full := name
	if specialty != "" {
		full = name + " (" + specialty + ")"
	}

	return Effect{
		Kind:        EffectLoseSkill,
		Detail:      "lose every level of " + full,
		Skill:       name,
		Specialties: []string{specialty},
	}
}

// chr moves a characteristic.
func chr(which string, delta int) Effect {
	sign := "+"
	if delta < 0 {
		sign = ""
	}

	return Effect{
		Kind:           EffectCharacteristic,
		Detail:         sign + itoa(delta) + " " + which,
		Characteristic: which,
		Delta:          delta,
	}
}

// chrRolled moves a characteristic by an amount the book rolls for: "Lose
// 1d3 from your choice of STR or END" (p. 156). The sign comes from up,
// because a table that says "lose 1d3" and one that says "gain 1d6" are the
// same shape.
func chrRolled(which, rolled string, up bool) Effect {
	verb := "lose "
	delta := -1

	if up {
		verb = "gain "
		delta = 1
	}

	return Effect{
		Kind:           EffectCharacteristic,
		Detail:         verb + rolled + " " + which,
		Characteristic: which,
		Delta:          delta,
		Dice:           rolled,
	}
}

// pickSkill is the compact form of a choice between whole skills, which
// rank tables and skill tables use where a cell reads "Chef or Mechanic".
func pickSkill(names ...string) Effect {
	if len(names) == 1 {
		return skill(names[0], anyFor(names[0])...)
	}

	options := make([]Option, len(names))
	for i, name := range names {
		options[i] = opt(name, skill(name, anyFor(name)...))
	}

	return pick("choose "+joinOr(names), options...)
}

// anyFor is the "(Any)" a skill may carry, which is none for a skill the
// book prints no specialties under.
//
// p. 116 is what "(Any)" means: it "allows the character to select a
// specialty within that skill". Seventeen of the skills these tables name
// -- Admin, Broker, Mechanic, Recon and the rest -- have no specialty to
// select, so recording one leaves the player an unanswerable question and
// the sheet a second entry for a skill the book counts once.
//
// The engine already asks this question correctly one layer up, where a
// result names a whole list of skills to choose between: anySkill takes the
// picked definition's own specialties and passes nil where there are none.
func anyFor(name string) []string {
	definition, found := SkillByName(name)
	if !found || (len(definition.Specialties) == 0 && !definition.Open) {
		return nil
	}

	return []string{"Any"}
}

// joinOr renders a list the way the page does.
func joinOr(names []string) string {
	switch len(names) {
	case 0:
		return ""
	case 1:
		return names[0]
	}

	return strings.Join(names[:len(names)-1], ", ") + " or " + names[len(names)-1]
}

// pick is "gain a level in A or B" -- a choice between effect lists.
func pick(detail string, options ...Option) Effect {
	return Effect{Kind: EffectChoice, Detail: detail, Options: options}
}

// opt is one branch of a pick.
func opt(label string, effects ...Effect) Option {
	return Option{Label: label, Effects: effects}
}

// checkSkill is "roll <skill> <n>+", with the two branches the book gives.
func checkSkill(name string, number int, success, failure []Effect) Effect {
	return Effect{
		Kind:    EffectCheck,
		Detail:  "roll " + name + " " + itoa(number) + "+",
		Check:   &Target{Skill: name, Number: number},
		Success: success,
		Failure: failure,
	}
}

// checkChr is "roll <characteristic> <n>+", with the two branches the book
// gives.
//
//nolint:unparam // the target is the page's, not a constant; 8+ is merely what every one so far says
func checkChr(which string, number int, success, failure []Effect) Effect {
	return Effect{
		Kind:    EffectCheck,
		Detail:  "roll " + which + " " + itoa(number) + "+",
		Check:   &Target{Characteristic: which, Number: number},
		Success: success,
		Failure: failure,
	}
}

// injury sends the character to the Injury table (p. 119). ERRATA E-1: three
// references in the book cite p. 136, which carries no table.
func injury(times int) Effect {
	detail := "roll on the Injury table"
	if times > 1 {
		detail = "roll " + itoa(times) + " times on the Injury table"
	}

	return Effect{Kind: EffectInjury, Detail: detail, Times: times}
}

// lifeEvent is the shared table every career's d66 reaches at 31-36.
func lifeEvent() Effect {
	return Effect{Kind: EffectLifeEvent, Detail: "roll on the Life Events table"}
}

// mishapNoEject rolls the career's own mishap table without ejecting the
// character, which is what an event result of 11 does in all three careers.
func mishapNoEject() Effect {
	return Effect{
		Kind:   EffectMishapNoEject,
		Detail: "roll on the Mishap table, without leaving the career",
	}
}

// benefitRolls grants (or, negative, removes) mustering-out rolls. The
// scope of any modifier matters: see ERRATA E-3.
func benefitRolls(count, modifier int, scope BenefitScope) Effect {
	detail := "gain " + itoa(count) + " benefit rolls"
	if count < 0 {
		detail = "lose " + itoa(-count) + " benefit rolls"
	}

	if modifier != 0 {
		detail += " at +" + itoa(modifier)
	}

	return Effect{
		Kind:     EffectBenefitRolls,
		Detail:   detail,
		Count:    count,
		Modifier: modifier,
		Scope:    scope,
	}
}

// loseAllBenefits forfeits every mustering-out roll accrued in this career
// so far. Sixty-seven results across the corpus do this, and they word it
// four ways -- "from this career", "gained to this point", "collected so
// far", "gained before this event" -- so detail carries the book's own
// phrasing. They mean the same thing: what has been earned here is gone,
// and a career the character stays in goes on earning from the next term.
func loseAllBenefits(detail string) Effect {
	return Effect{Kind: EffectBenefitRolls, Detail: detail, ForfeitAll: true}
}

// gainBenefitRollsRolled and loseBenefitRollsRolled move mustering-out
// rolls by a number the book rolls for: "gain 1d3 Benefit rolls", "lose 1d6
// Benefit rolls". Count carries the direction and Dice the expression,
// because the engine cannot know the number until it throws.
func gainBenefitRollsRolled(rolled string) Effect {
	return Effect{
		Kind:   EffectBenefitRolls,
		Detail: "gain " + rolled + " benefit rolls",
		Dice:   rolled,
		Count:  1,
	}
}

func loseBenefitRollsRolled(rolled string) Effect {
	return Effect{
		Kind:   EffectBenefitRolls,
		Detail: "lose " + rolled + " benefit rolls",
		Dice:   rolled,
		Count:  -1,
	}
}

// cashRolls grants mustering-out rolls that must be spent on the Cash
// column: "two Cash Benefit rolls". They queue for Step 19 like any other.
func cashRolls(count int) Effect {
	return Effect{
		Kind:     EffectBenefitRolls,
		Detail:   "gain " + itoa(count) + " benefit rolls, spent on the cash table",
		Count:    count,
		CashOnly: true,
	}
}

// cashRollsNow takes those rolls at once rather than queueing them for Step
// 19, which is what the results that read "immediately" ask for.
func cashRollsNow(count int) Effect {
	return Effect{
		Kind:      EffectBenefitRolls,
		Detail:    "take " + itoa(count) + " cash benefit rolls immediately",
		Count:     count,
		CashOnly:  true,
		Immediate: true,
	}
}

// cashRollsNowRerolling is cashRollsNow with the clause one result adds:
// a row that pays nothing is thrown again rather than wasted.
func cashRollsNowRerolling(count int) Effect {
	effect := cashRollsNow(count)

	effect.Detail += ", re-rolling any result of nothing"

	effect.RerollNothing = true

	return effect
}

// noAdmissionFor closes higher education for a number of terms, which two
// of the graduate tracks' failure results do.
func noAdmissionFor(terms int) Effect {
	return Effect{
		Kind:   EffectEducationLockout,
		Detail: "no admission to any higher learning institution for " + itoa(terms) + " terms",
		Terms:  terms,
	}
}

// relationship gains NPCs of one kind. Count may be a fixed number; where
// the book rolls for it, Dice carries the expression.
func relationship(kind Relationship, count int, rolled string) Effect {
	detail := "gain " + itoa(count) + " " + string(kind)
	if rolled != "" {
		detail = "gain " + rolled + " " + string(kind)
	}

	return Effect{
		Kind:         EffectRelationship,
		Detail:       detail,
		Relationship: kind,
		Count:        count,
		Dice:         rolled,
	}
}

// relationshipAt is relationship where the page names the rating the tie
// starts at: "Gain an Ally with a Relationship Rating of 125" (Youth Path 1
// result 12, p. 69).
func relationshipAt(kind Relationship, count, rating int) Effect {
	tie := relationship(kind, count, "")

	tie.Detail += " at a Relationship Rating of " + itoa(rating)

	tie.Rating = rating

	return tie
}

// rating moves Relationship Ratings (p. 320). Amount is the fixed change;
// where the book rolls for it, rolled carries the expression and amount's
// sign says which way it goes.
func rating(target TieTarget, only Relationship, amount int, rolled string) Effect {
	detail := "Relationship Rating"

	switch target {
	case TargetAll:
		detail = "every relationship's " + detail
	case TargetFamily:
		detail = "every family member's " + detail
	case TargetRole, TargetAllOfRole:
		detail = "a relative's " + detail
	case TargetOne:
		detail = "one " + string(only) + "'s " + detail
		if only == "" {
			detail = "one relationship's Relationship Rating"
		}
	}

	move := " by " + itoa(amount)
	if rolled != "" {
		move = " by " + rolled
	}

	if amount < 0 {
		detail = "lower " + detail + move
	} else {
		detail = "raise " + detail + move
	}

	return Effect{
		Kind:         EffectRating,
		Detail:       detail,
		Target:       target,
		Relationship: only,
		Modifier:     amount,
		Dice:         rolled,
	}
}

// ratingOfRole moves the Relationship Rating of one relative of a named
// role, or of all of them.
//
//nolint:unparam // the role is the page's; every result that names one happens to name a parent
func ratingOfRole(role string, all bool, amount int) Effect {
	which := "one " + role + "'s"
	if all {
		which = "every " + role + "'s"
	}

	verb := "raise "
	if amount < 0 {
		verb = "lower "
	}

	target := TargetRole
	if all {
		target = TargetAllOfRole
	}

	return Effect{
		Kind:     EffectRating,
		Detail:   verb + which + " Relationship Rating by " + itoa(amount),
		Target:   target,
		Role:     role,
		Modifier: amount,
	}
}

// relationshipRolledAt is relationshipAt where the count is a throw: "1d3
// Contacts with a Relationship Rating of 50" (Youth Path 5 result 11).
//
//nolint:unparam // the kind is the page's; every rolled count in the book happens to be Contacts
func relationshipRolledAt(kind Relationship, rolled string, rating int) Effect {
	tie := relationship(kind, 0, rolled)

	tie.Detail += " at a Relationship Rating of " + itoa(rating)

	tie.Rating = rating

	return tie
}

// loseRelative removes one family tie: "Choose a family member from your
// existing list and lose them" (Teenage Path 2 result 3, p. 78).
func loseRelative() Effect {
	return Effect{
		Kind:   EffectLoseTie,
		Detail: "lose a family member",
		Target: TargetFamily,
	}
}

// loseTie removes a relationship, trying the kinds in the order given --
// which is the order Life Event 4 prints them in, and reverses on a 4-6.
func loseTie(order ...Relationship) Effect {
	detail := "lose a relationship"

	if len(order) > 0 {
		names := make([]string, len(order))
		for i, kind := range order {
			names[i] = string(kind)
		}

		detail = "lose " + joinOr(names) + ", in that order"
	}

	return Effect{Kind: EffectLoseTie, Detail: detail, Order: order, Count: 1}
}

// loseThatTie is the nine hooks that read "if you fail that Advancement
// roll, you lose the Ally": the Ally the event granted a line earlier,
// which the engine reaches as the most recently gained one that is not
// family. ERRATA E-39.
func loseThatTie(kind Relationship) Effect {
	effect := loseTie(kind)

	effect.Newest = true
	effect.ExcludeFamily = true
	effect.Detail = "lose that " + string(kind)

	return effect
}

// loseTies is loseTie several times over: "lose 1D3 Allies and Contacts",
// "lose two Contacts". Each removal tries the kinds in the order given, so
// a character with one Ally and three Contacts losing three of "Ally and
// Contact" loses the Ally and two Contacts.
func loseTies(count int, order ...Relationship) Effect {
	effect := loseTie(order...)

	effect.Count = count
	effect.Detail = "lose " + itoa(count) + " x (" + effect.Detail + ")"

	return effect
}

// loseTiesRolled is loseTies where the book rolls for how many.
//
//nolint:unparam // the expression is the page's; 1d3 is merely the only one transcribed so far
func loseTiesRolled(rolled string, order ...Relationship) Effect {
	effect := loseTie(order...)

	effect.Count = 0
	effect.CountDice = rolled
	effect.Detail = "lose " + rolled + " x (" + effect.Detail + ")"

	return effect
}

// loseEveryTie removes every tie of the kinds given. ThisCareer narrows it
// to the ties this career gave, which is what "every Contact made in this
// career" means; without it the result reaches the character's whole life.
func loseEveryTie(thisCareer bool, kinds ...Relationship) Effect {
	effect := loseTie(kinds...)

	effect.Count = everyTie
	effect.ThisCareer = thisCareer

	where := ""
	if thisCareer {
		where = " gained in this career"
	}

	names := make([]string, len(kinds))
	for i, kind := range kinds {
		names[i] = string(kind) + "s"
	}

	effect.Detail = "lose every " + joinOr(names) + where

	return effect
}

// everyTie is the Count that means "all of them" rather than a number.
const everyTie = -1

// loseTieOrElse is the results that say what happens when there is nobody
// to lose: "Lose one Ally or Contact. If you have none, gain an Enemy with
// a Relationship Rating of -110."
func loseTieOrElse(order []Relationship, fallback ...Effect) Effect {
	effect := loseTie(order...)

	details := make([]string, len(fallback))
	for i, e := range fallback {
		details[i] = e.Detail
	}

	effect.Fallback = fallback

	effect.Detail += ", or with nobody to lose, " + strings.Join(details, " and ")

	return effect
}

// becomes changes what ties are: "one Contact or Ally from this career
// becomes an Enemy", "1D3 existing Contacts become Allies". The new rating
// is the middle of the band the new kind sits in, because a change of kind
// is not a change of rating and the book gives no number for it -- the same
// argument as ERRATA E-23.
func becomes(count int, to Relationship, from ...Relationship) Effect {
	names := make([]string, len(from))
	for i, kind := range from {
		names[i] = string(kind)
	}

	return Effect{
		Kind:         EffectBecome,
		Detail:       itoa(count) + " x " + joinOr(names) + " becomes " + string(to),
		From:         from,
		Relationship: to,
		Count:        count,
	}
}

// becomesRolled is becomes where the book rolls for how many.
func becomesRolled(rolled string, to Relationship, from ...Relationship) Effect {
	effect := becomes(1, to, from...)

	effect.Count = 0
	effect.CountDice = rolled
	effect.Detail = rolled + " x " + effect.Detail[len("1 x "):]

	return effect
}

// becomesInCareer is becomes, narrowed to the ties this career gave.
func becomesInCareer(count int, to Relationship, from ...Relationship) Effect {
	effect := becomes(count, to, from...)

	effect.ThisCareer = true

	effect.Detail += " gained in this career"

	return effect
}

// becomesEveryOneInCareer turns every tie this career gave: "every Contact,
// Rival and Ally made in this career becomes an Enemy".
func becomesEveryOneInCareer(to Relationship, from ...Relationship) Effect {
	effect := becomes(everyTie, to, from...)

	effect.ThisCareer = true
	effect.Detail = "every " + effect.Detail[len("-1 x "):] + " gained in this career"

	return effect
}

// becomesTheLastOne is Gambler mishap 11: "one Contact or Ally from this
// career becomes an Enemy, and every other relationship gained here is
// lost". The one that changed is the one that survives.
func becomesTheLastOne(to Relationship, from ...Relationship) Effect {
	effect := becomes(1, to, from...)

	effect.ThisCareer = true
	effect.LoseTheRest = true

	effect.Detail += " and every other relationship gained in this career is lost"

	return effect
}

// ratingOfSome is Teenage Path 3 result 4: a rolled number of ties, none of
// them family, move by a fixed amount.
func ratingOfSome(rolled string, minimum, amount int, first Relationship, rest ...Relationship) Effect {
	only := append([]Relationship{first}, rest...)

	names := make([]string, len(only))
	for i, kind := range only {
		names[i] = string(kind)
	}

	direction, size := "gain ", amount
	if amount < 0 {
		direction, size = "lose ", -amount
	}

	return Effect{
		Kind: EffectRating,
		Detail: rolled + " (minimum " + itoa(minimum) + ") non-family " + joinOr(names) +
			" " + direction + itoa(size) + " relationship rating",
		Target:        TargetAll,
		Relationship:  first,
		From:          only,
		CountDice:     rolled,
		Minimum:       minimum,
		Modifier:      amount,
		ExcludeFamily: true,
	}
}

// becomesOrElse is "a Rival becomes an Enemy, or a Rival is gained where
// there was none".
func becomesOrElse(to Relationship, from []Relationship, fallback ...Effect) Effect {
	effect := becomes(1, to, from...)

	details := make([]string, len(fallback))
	for i, e := range fallback {
		details[i] = e.Detail
	}

	effect.Fallback = fallback

	effect.Detail += ", or with none, " + strings.Join(details, " and ")

	return effect
}

// credits pays the character an amount the book gives as dice.
func credits(rolled string) Effect {
	return Effect{Kind: EffectCredits, Detail: "gain " + rolled + " credits", Dice: rolled}
}

// stashItem adds something to the character's stash. An empty item is the
// results that say to lose the whole of it.
func stashItem(item string) Effect {
	return Effect{Kind: EffectStash, Detail: "add " + item + " to the stash", Item: item}
}

// stashValued adds something the book prints a credit value for: "a Company
// Share worth 2D6 x Cr100000". The value is rolled and recorded, which is
// not an invented value -- the book gave the dice.
func stashValued(item, rolled string) Effect {
	return Effect{
		Kind:      EffectStash,
		Detail:    "add " + item + ", worth " + rolled + " credits, to the stash",
		Item:      item,
		Dice:      rolled,
		Count:     1,
		Deviation: stashDeviation(item),
	}
}

// stashDeviation is the ERRATA identifier a stashed item's valuation rests
// on, where it rests on one.
//
// Only the company share does. The rows granting several print no value at
// all -- "Two Company Shares", "Ten Shares in the Company" -- so every one
// is worth what p. 128 defines rather than what its row omits, and a
// character's wealth rests on that reading.
func stashDeviation(item string) string {
	if item == companyShare {
		return "E-33"
	}

	return ""
}

// stashCountValued is stashValued for the results that grant several at
// once: "Six Pieces of Art, 2D6 x Cr10000 each". Each is valued on its own
// throw, which is what "each" means.
func stashCountValued(count int, item, rolled string) Effect {
	effect := stashValued(item, rolled)

	effect.Count = count
	effect.Detail = "add " + itoa(count) + " x " + item +
		", each worth " + rolled + " credits, to the stash"

	return effect
}

// stashCountRolled is stashCountValued where the book rolls for how many:
// Scientist's "1D6 Company Shares".
func stashCountRolled(rolled, item, value string) Effect {
	return Effect{
		Kind:      EffectStash,
		Detail:    "add " + rolled + " x " + item + ", each worth " + value + " credits, to the stash",
		Item:      item,
		Dice:      value,
		CountDice: rolled,
		Deviation: stashDeviation(item),
	}
}

// companyShare is p. 128's definition: "a 1% share of the corporation with
// which the character has been associated. The value of this share is 2d6
// x 100,000 credits."
//
// The value lives here rather than at each site because the definition is
// the book's: a career row reading "Three Company Shares" means three of
// these, whether or not that row reprints the number. See ERRATA E-33.
const (
	companyShare      = "a company share"
	companyShareValue = "2d6x100000"
)

// weaponOrItsUse is p. 129's Weapon benefit in full: "The player may choose
// a weapon ... If the player wishes, they may choose to take a level in
// Melee (Any) or Gun Combat (Any) in lieu of a weapon."
//
// Fourteen benefit rows print it. The engine cannot name a weapon -- the
// charts are in another book -- but the choice is a rule, and two of its
// three branches are things this engine can carry out.
func weaponOrItsUse() Effect {
	return pick("a weapon, or a level in using one (p. 129)",
		opt("a weapon", stashItem("a weapon of the character's choice")),
		opt("Melee (Any)", skill("Melee", "Any")),
		opt("Gun Combat (Any)", skill("Gun Combat", "Any")))
}

// loseStashItem removes every possession of one name.
func loseStashItem(item string) Effect {
	return Effect{
		Kind:   EffectStash,
		Detail: "lose any " + item + " held",
		Item:   item,
		Lose:   true,
	}
}

// group applies several effects as one, for the benefit rows that carry a
// single effect and say two things.
func group(effects ...Effect) Effect {
	details := make([]string, len(effects))
	for i, e := range effects {
		details[i] = e.Detail
	}

	return Effect{Kind: EffectGroup, Detail: strings.Join(details, ", and "), Group: effects}
}

// throwModifier attaches a modifier to a named future throw.
func throwModifier(applies string, value int) Effect {
	sign := "+"
	if value < 0 {
		sign = ""
	}

	return Effect{
		Kind:     EffectModifier,
		Detail:   sign + itoa(value) + " to the " + applies,
		Modifier: value,
		Applies:  applies,
	}
}

// modifierFor attaches a modifier to a named throw for a number of uses,
// narrowed the way an enlistment penalty is. Every other constructor in
// this file is a shorthand for a shape of it.
func modifierFor(applies string, value, uses int, detail string, narrow enlistmentNarrowing) Effect {
	effect := enlistmentPenalty(value, detail, narrow)

	effect.Applies = applies

	// One use and none are the same thing -- the throw it was granted for
	// -- so the record keeps one spelling of it rather than two.
	if uses > 1 {
		effect.Uses = uses
	}

	return effect
}

// enlistmentPenalty is the twelve results that modify an enlistment by the
// class of career being entered rather than by name. The narrowing sets
// may each be empty, which means the modifier reaches every career.
//
// standing says whether it outlives the throw it modifies: "-2 DM to enter
// every career after this one where enlistment is based on EDU or CHA" is
// a standing modifier, and "-2 DM to enter your next career" is not.
func enlistmentPenalty(value int, detail string, narrow enlistmentNarrowing) Effect {
	return Effect{
		Kind:              EffectModifier,
		Detail:            detail,
		Modifier:          value,
		Applies:           enlistmentThrow,
		OnTags:            narrow.OnTags,
		NotTags:           narrow.NotTags,
		OnCareers:         narrow.OnCareers,
		NotCareers:        narrow.NotCareers,
		OnCharacteristics: narrow.OnCharacteristics,
		OnSkill:           narrow.OnSkill,
		WhileInThisCareer: narrow.WhileInThisCareer,
		Standing:          narrow.Standing,
	}
}

// enlistmentNarrowing is what an enlistmentPenalty applies to, gathered
// into one argument so a call site reads as the page does.
type enlistmentNarrowing struct {
	OnTags            []Tag
	NotTags           []Tag
	OnCareers         []string
	NotCareers        []string
	OnCharacteristics []string
	OnSkill           string
	WhileInThisCareer bool
	Standing          bool
}

// The names a modifier or an automatic can be filed under. A result grants
// one of these and the engine consumes it at the throw of that name; a name
// no part of the engine reads is a modifier that silently never applies,
// which TestEveryModifierNamesAThrowTheEngineTakes is what prevents.
const (
	enlistmentThrow  = "next enlistment attempt"
	survivalThrow    = "next survival roll"
	advancementThrow = "next advancement roll"
	commissionThrow  = "next commission roll"
	admissionThrow   = "admission to any higher education"
	skillCheckThrow  = "a skill check"
)

// throws is every name a modifier or an automatic may be filed under. The
// engine consumes each of them at the throw it names, and
// TestEveryModifierNamesAThrowTheEngineTakes holds the corpus to the list:
// a result filed under a name nothing reads is a modifier that is recorded,
// shown in the transcript, and never applied.
//
// It is a small list because the book has few throws. A result that wants
// to narrow one -- "every advancement roll in a military career" -- says so
// in the narrowing fields rather than in the name.
func throws() []string {
	return []string{
		enlistmentThrow,
		survivalThrow,
		advancementThrow,
		commissionThrow,
		admissionThrow,
		skillCheckThrow,
	}
}

// autoEnlist is "you may enlist automatically", which several results
// grant: the next enlistment succeeds without a throw. The narrowing is an
// enlistment penalty's, because one result restricts it to a class of
// career -- "a business, military, corporate or colonist career".
func autoEnlist(detail string, narrow enlistmentNarrowing) Effect {
	return Effect{
		Kind:              EffectAutoSuccess,
		Detail:            detail,
		Applies:           enlistmentThrow,
		OnTags:            narrow.OnTags,
		NotTags:           narrow.NotTags,
		OnCareers:         narrow.OnCareers,
		NotCareers:        narrow.NotCareers,
		OnCharacteristics: narrow.OnCharacteristics,
	}
}

// rejoinPreviousCareer is Fringe Marketer mishap 8: "return to the career
// you held before this one", without an enlistment roll. Which career that
// was is in the record, not in this package, so the engine names it.
func rejoinPreviousCareer() Effect {
	return Effect{
		Kind:       EffectTransfer,
		Detail:     "re-enter the previous career without an enlistment roll",
		Career:     PreviousCareer,
		Assignment: "",
	}
}

// PreviousCareer is the Career an [EffectTransfer] names when the
// destination is "the career you held before this one" rather than a
// career the page names. The engine resolves it from the service record.
const PreviousCareer = "\x00previous"

// ageBy moves the character's age outside the four years a term takes.
func ageBy(years int) Effect {
	return Effect{Kind: EffectAge, Detail: "age " + itoa(years) + " years", Years: years}
}

// sentenceBy lengthens or shortens a forced spell in a career. Shortening
// it to one term or less is a release, which the engine reads as the
// sentence being served.
func sentenceBy(terms int) Effect {
	detail := "add " + itoa(terms) + " terms to the sentence"
	if terms < 0 {
		detail = "take " + itoa(-terms) + " terms off the sentence"
	}

	return Effect{Kind: EffectSentence, Detail: detail, Terms: terms}
}

// debt is money owed rather than money held, which two results impose and
// nothing in the engine can spend. It is recorded as a condition for the
// same reason an addiction is: the book names it and prices nothing else.
func debt(amount string) Effect {
	return Effect{
		Kind:      EffectCondition,
		Detail:    "a debt of " + amount + " credits",
		Condition: "a debt of " + amount + " credits",
	}
}

// addiction is the thirteen results that leave a character dependent on
// something. The book names the kind and leaves the substance to the
// player -- "an addiction to alcohol or a drug of your choice" -- and
// prints no characteristic cost beside it (ERRATA E-38).
func addiction(to string) Effect {
	return Effect{
		Kind:      EffectCondition,
		Detail:    "an addiction to " + to,
		Condition: "an addiction to " + to,
	}
}

// religion is the nine results that give a character one. The book says
// "choose or invent a religion", which is the player's to answer on the
// sheet; the engine records that they have one.
func religion() Effect {
	return Effect{
		Kind:      EffectCondition,
		Detail:    "a religion of the character's choosing",
		Condition: "a religion of the character's choosing",
	}
}

// rollSub is a 1d6 table printed inside a result. The rows are given in
// the order the page prints them and must cover 1 to 6 exactly once, which
// TestEverySubTableCoversTheDie holds them to.
func rollSub(detail string, rows ...SubRow) Effect {
	return Effect{Kind: EffectSubTable, Detail: detail, Sub: rows}
}

// rollSubOn2d6 is a sub-table thrown on 2d6 with a characteristic's
// modifier added, which the enslaved tables' escape result asks for: "Roll
// 2d6 and add your DEX bonus to the roll."
func rollSubOn2d6(which, detail string, rows ...SubRow) Effect {
	effect := rollSub(detail, rows...)

	effect.Dice = "2d6"
	effect.Characteristic = which

	return effect
}

// on is one row of a rollSub covering a single result of the die.
func on(result int, summary string, effects ...Effect) SubRow {
	return SubRow{From: result, To: result, Summary: summary, Effects: effects}
}

// onRange is one row covering a span: "on a 2-5".
func onRange(from, to int, summary string, effects ...Effect) SubRow {
	return SubRow{From: from, To: to, Summary: summary, Effects: effects}
}

// pool grants a modifier the character spends themselves, a little at a
// time. total is the whole of it and perThrow the most that may go on one
// throw; applies names the throws it may be spent on.
func pool(total, perThrow int, detail string, whileInCareer bool, applies ...string) Effect {
	return Effect{
		Kind:              EffectPool,
		Detail:            detail,
		Modifier:          total,
		Uses:              perThrow,
		Spendable:         applies,
		WhileInThisCareer: whileInCareer,
	}
}

// onFailure attaches effects to the failure of a named throw, which nine
// results do: "if you fail that Advancement roll, you lose the Ally and
// take a -2 DM to your next Advancement roll".
func onFailure(applies, detail string, effects ...Effect) Effect {
	return Effect{
		Kind:    EffectOnFailure,
		Detail:  detail,
		Applies: applies,
		Failure: effects,
	}
}

// autoFailure is the mirror of autoSuccess: a named throw fails without
// being rolled.
func autoFailure(applies, detail string) Effect {
	return Effect{Kind: EffectAutoFailure, Detail: detail, Applies: applies}
}

// advance is an automatic advancement, granted without a throw.
func advance() Effect {
	return Effect{Kind: EffectAdvance, Detail: "gain an automatic advancement"}
}

// mustContinue forces another term in this career.
func mustContinue() Effect {
	return Effect{Kind: EffectContinue, Detail: "remain in this career for another term"}
}

// chooseNewCareer is "choose another career", which five Colonist mishaps
// end on (ERRATA E-6). It leaves the career without
// naming the next one, so the character enlists in whatever they choose on
// the ordinary terms -- which is not what transfer does.
func chooseNewCareer() Effect {
	return Effect{Kind: EffectChooseCareer, Detail: "leave this career and choose another"}
}

// transfer sends the character to another career. Terms is how many they
// owe there; zero means until they choose to leave.
func transfer(name, assignment string, terms int) Effect {
	detail := "enter the " + name + " career"
	if assignment != "" {
		detail += " with the " + assignment + " assignment"
	}

	if terms > 0 {
		detail += " for " + itoa(terms) + " terms"
	}

	return Effect{
		Kind:       EffectTransfer,
		Detail:     detail,
		Career:     name,
		Assignment: assignment,
		Terms:      terms,
	}
}

// newHomeworld reassigns the homeworld, which six of Colonist's eleven
// mishaps do. Until the setting data lands there is nothing to reassign it
// to, so the engine records the demand.
func newHomeworld() Effect {
	return Effect{Kind: EffectNewHomeworld, Detail: "take a new homeworld"}
}

// autoSuccess makes a named future throw succeed without rolling.
func autoSuccess(applies string) Effect {
	return Effect{
		Kind:    EffectAutoSuccess,
		Detail:  "an automatic success on the " + applies,
		Applies: applies,
	}
}

// rollTable rolls on one of the career's skill tables by name.
func rollTable(kind SkillTableKind) Effect {
	return Effect{
		Kind:   EffectRollTable,
		Detail: "roll once on the " + kind.String() + " table",
		Table:  kind,
	}
}

// rollAssignmentOf rolls on a named assignment's skill table in another
// career, which Colonist event 54 asks for: "Roll twice on the Assignment:
// Ambassador table of the Diplomatic Service career" (p. 176).
func rollAssignmentOf(careerName, assignment string) Effect {
	return Effect{
		Kind:       EffectRollTable,
		Detail:     "roll on the " + assignment + " table of the " + careerName + " career",
		Table:      AssignmentSkills,
		Career:     careerName,
		Assignment: assignment,
	}
}

// rollOtherAssignment rolls on the skill table of an assignment the
// character is not in, which two careers' events ask for.
func rollOtherAssignment() Effect {
	return Effect{
		Kind:            EffectRollTable,
		Detail:          "roll once on the skill table of an assignment other than your own",
		Table:           AssignmentSkills,
		OtherAssignment: true,
	}
}

// commission attempts a commission outside the ordinary Step 14 offer.
func commission(modifier int) Effect {
	return Effect{
		Kind:     EffectCommission,
		Detail:   "attempt a commission at +" + itoa(modifier),
		Modifier: modifier,
	}
}

// personalDevelopment is the Personal Development table. Every career
// transcribed so far prints the same six rows -- five characteristics and
// Athletics -- and a career that prints different ones writes its own.
func personalDevelopment() SkillTable {
	return SkillTable{
		Kind: PersonalDevelopment,
		Name: "Personal Development",
		Rows: [6]Effect{
			chr("STR", 1), chr("DEX", 1), chr("END", 1),
			chr("INT", 1), chr("EDU", 1), skill("Athletics", "Any"),
		},
	}
}

// unimplemented is a result the engine cannot carry out, carrying the
// book's demand so the record says what it could not do.
func unimplemented(detail string) Effect {
	return Effect{Kind: EffectUnimplemented, Detail: detail}
}

// itoa is strconv.Itoa under a shorter name, because these detail strings
// are mostly concatenation and `+ itoa(n) +` reads better than the whole
// package path. It was a hand-rolled implementation until #116: four times
// slower, five allocations against none, and "-" for math.MinInt.
func itoa(n int) string { return strconv.Itoa(n) }
