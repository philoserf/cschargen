package career

import "strings"

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
// or higher, no additional benefit is gained" -- and the engine currently
// adds instead. The data here is the page; the engine's reading of it is
// issue #20.
func skillAt(name string, level int, specialties ...string) Effect {
	granted := skill(name, specialties...)

	granted.Detail = "gain " + name + " " + itoa(level)
	granted.Level = level

	return granted
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
//
//nolint:unparam // the expression is the page's; 1d3 is merely the only one transcribed so far
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
		return skill(names[0], "Any")
	}

	options := make([]Option, len(names))
	for i, name := range names {
		options[i] = opt(name, skill(name, "Any"))
	}

	return pick("choose "+joinOr(names), options...)
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

	return Effect{Kind: EffectLoseTie, Detail: detail, Order: order}
}

// credits pays the character an amount the book gives as dice.
func credits(rolled string) Effect {
	return Effect{Kind: EffectCredits, Detail: "gain " + rolled + " credits", Dice: rolled}
}

// stashItem adds something to the character's stash.
func stashItem(item string) Effect {
	return Effect{Kind: EffectStash, Detail: "add " + item + " to the stash", Item: item}
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

// unimplemented is a result this milestone cannot carry out, carrying the
// book's demand so the record says what it could not do.
func unimplemented(detail string) Effect {
	return Effect{Kind: EffectUnimplemented, Detail: detail}
}

// itoa avoids importing strconv into a package that is otherwise pure data.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var digits []byte

	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)

		n /= 10
	}

	if negative {
		return "-" + string(digits)
	}

	return string(digits)
}
