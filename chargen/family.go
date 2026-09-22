package chargen

import (
	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/setting"
)

// Step 5: Determine Family (pp. 57-61).
//
// The book opens the step by saying it is optional -- "The player and
// Referee should determine as much or as little about the character's family
// as they wish ... Skipping this step can speed up character generation" --
// and then notes what skipping costs: "it may also deprive the character of
// potential Allies and Contacts that could become meaningful parts of the
// campaign". So the engine runs it and a flag turns it off, rather than the
// other way around.
//
// It stops at cousins. The book does not: "this process can then be
// repeated to determine great-grandparents and further ancestors ... as far
// as one wishes". Cousins are where it stops giving specific instructions,
// and an unbounded recursion is not a rule.

// Family is the household Step 5 built.
type Family struct {
	// Situation is the Human Birth Situation result (p. 58): a communal
	// group, a same sex couple, or a heterosexual couple.
	Situation string `json:"situation"`

	// Detail is what the result's own sub-rolls made of it -- whether the
	// couple married, whether an artificial womb was used.
	Detail string `json:"detail,omitempty"`

	// ParentAges are how old each parent was at the character's birth
	// (p. 59).
	ParentAges []int `json:"parentAges,omitempty"`

	// Firstborn and Lastborn come from the face of the parental age die: a
	// 1 makes the character the first child and a 10 the last (p. 59).
	// Both can be false, and both can be true where two parents rolled the
	// two ends -- the book does not say what happens then, and the engine
	// lets the constraint on siblings' ages be satisfied by neither.
	Firstborn bool `json:"firstborn,omitempty"`
	Lastborn  bool `json:"lastborn,omitempty"`
}

// The three Human Birth Situation results, in the order the chart prints
// them (p. 58).
//
// setting.BirthSituations holds the same three, because a world in the data
// file may force one and the validator has to know which words are valid.
// A test holds the two lists against each other.
const (
	SituationCommunal     = "a communal group"
	SituationSameSex      = "a same sex couple"
	SituationHeterosexual = "a heterosexual couple"
)

// The five roles a family tie can have.
const (
	RoleParent      = "parent"
	RoleSibling     = "sibling"
	RoleGrandparent = "grandparent"
	RoleAuntOrUncle = "aunt or uncle"
	RoleCousin      = "cousin"
)

// The Relationship Ratings Step 5 grants. ERRATA E-22: p. 320 says family
// members "will always begin as Allies with a Relationship Rating of 150"
// and the step itself says 125 and 100. The step wins, and 100 is an Ally
// by the step's word although it is a Contact by the band.
const (
	closeFamilyRating    = 125
	extendedFamilyRating = 100
)

// familyCite is the page the birth situation chart is printed on.
const familyCite = "p. 58"

// determineFamily is Step 5. It runs after the homeworld is known, because
// the birth situation throw is modified by where the character was born and
// the parental age table is read by its tech level.
func (g *Generator) determineFamily(world setting.World) error {
	step := g.log.Step("Step 5: Determine Family", "p. 57")

	family := &Family{}

	g.char.State.Family = family

	parents := g.birthSituation(world, family, step)
	g.parentalAges(family, parents, step)
	g.addFamily(RoleParent, parents, closeFamilyRating, step, "a parent")

	g.siblings(family, step)
	g.extendedFamily(step)

	return nil
}

// birthSituation is the chart of p. 58, and returns how many parents the
// household holds.
//
// The modifiers are a list of world names -- "Add 30 to the roll if born
// on ..." -- so they live in the setting data beside the world, and so does
// the handful of worlds where the result is not modified but replaced.
func (g *Generator) birthSituation(world setting.World, family *Family, step int) int {
	if forced := world.BirthSituationOnly; forced != "" {
		family.Situation = forced
		g.consequence(ConsequenceFamily, step,
			world.Name+" allows only one birth situation: "+forced, "")

		return g.householdSize(family)
	}

	roll := g.dice.D100()
	cause := g.log.Roll(roll, familyCite)

	total := roll.Total + world.BirthSituation

	switch {
	case total <= communalCeiling:
		family.Situation = SituationCommunal
	case total <= sameSexCeiling:
		family.Situation = SituationSameSex
	default:
		family.Situation = SituationHeterosexual
	}

	g.consequence(ConsequenceFamily, cause,
		"born to "+family.Situation, "")

	return g.householdSize(family)
}

// The two boundaries of the Human Birth Situation chart: "30 or less",
// "31-40", "41+".
const (
	communalCeiling = 30
	sameSexCeiling  = 40
)

// householdSize resolves each situation's own sub-rolls and returns how
// many parents the character gained.
func (g *Generator) householdSize(family *Family) int {
	switch family.Situation {
	case SituationCommunal:
		return g.communalHousehold(family)
	case SituationSameSex:
		return g.sameSexHousehold(family)
	}

	return g.heterosexualHousehold(family)
}

// communalHousehold is the first result: a polyamorous group on 1-3 of a
// d6, a commune or kibbutz on 4-6.
func (g *Generator) communalHousehold(family *Family) int {
	roll := g.dice.D6()
	cause := g.log.Roll(roll, familyCite)

	if roll.Total <= halfADie {
		// "1d6-1 (minimum 1) fathers in the home and 1d6+1 (minimum 1)
		// mothers in the home."
		fathers := max(g.dice.D6().Total-1, 1)
		mothers := max(g.dice.D6().Total+1, 1)

		family.Detail = "a polyamorous group of " + itoa(fathers) + " fathers and " +
			itoa(mothers) + " mothers"
		g.consequence(ConsequenceFamily, cause, family.Detail, "")

		return fathers + mothers
	}

	// "Gain 2d6-4 (minimum 4) parents as Allies."
	parents := max(g.dice.TwoD6().Total-4, 4)

	family.Detail = "a commune raising " + itoa(parents) + " parents' children together"
	g.consequence(ConsequenceFamily, cause, family.Detail, "")

	return parents
}

// sameSexHousehold is the second result: two parents, and a sub-roll for
// how the character came to be.
func (g *Generator) sameSexHousehold(family *Family) int {
	roll := g.dice.D6()
	cause := g.log.Roll(roll, familyCite)

	male := roll.Total%2 == 1
	how := g.dice.D6().Total

	switch {
	case male && how <= halfADie:
		family.Detail = "two fathers, and an artificial womb"
	case male && how <= 5:
		family.Detail = "two fathers, and a friend who carried the pregnancy"
	case male:
		family.Detail = "two fathers, and a surrogate they hired"
	case how == 1:
		family.Detail = "two mothers, and an artificial womb"
	case how <= 5:
		family.Detail = "two mothers, and an anonymous donor at a clinic"
	default:
		family.Detail = "two mothers, and a male friend who donated"
	}

	g.consequence(ConsequenceFamily, cause, family.Detail, "")

	return 2
}

// heterosexualHousehold is the third result, and the only one where the
// character may end up in the care of one parent rather than two -- though
// both are still gained as Allies, which is what the chart says.
func (g *Generator) heterosexualHousehold(family *Family) int {
	roll := g.dice.TwoD6()
	cause := g.log.Roll(roll, familyCite)

	switch {
	case roll.Total >= 7:
		family.Detail = "a married couple"
	case roll.Total >= 4:
		family.Detail = "a couple who never married"
	default:
		with := "their mother"
		if g.dice.D6().Total == 1 {
			with = "their father"
		}

		family.Detail = "a couple who separated, leaving the character with " + with
	}

	g.consequence(ConsequenceFamily, cause, family.Detail, "")

	return 2
}

// halfADie is the 1-3 half of a d6, which four of this step's sub-rolls
// split on.
const halfADie = 3

// parentalAges is the table of p. 59: "15 x 1d10" at tech level 10,
// "(23 x 1d10) - 6" at 11, "(28 x 1d10) - 12" at 12-13.
//
// The die's face is read as well as its value. "A roll of 1 on the 1d10
// indicates that your character was the first child born to their parents
// ... A roll of 10 indicates that your character was the final child", and
// both matter when the sibling ages are rolled.
func (g *Generator) parentalAges(family *Family, parents, step int) {
	multiplier, offset := parentalAgeTable(g.techLevel)

	for range parents {
		roll := g.dice.D10()
		cause := g.log.Roll(roll, "p. 59")

		age := roll.Total*multiplier + offset

		family.ParentAges = append(family.ParentAges, age)

		switch roll.Total {
		case 1:
			family.Firstborn = true
		case 10:
			family.Lastborn = true
		}

		g.consequence(ConsequenceFamily, cause,
			"a parent was "+itoa(age)+" at the character's birth", "")
	}

	switch {
	case family.Firstborn && family.Lastborn:
		// The book does not say what happens when two parents roll the two
		// ends. Neither constraint is applied, which is the only reading
		// that does not make the sibling table unsatisfiable.
		g.consequence(ConsequenceFamily, step,
			"one parent's die made the character the first child and another's the last; "+
				"neither constrains the siblings", "")

		family.Firstborn = false
		family.Lastborn = false
	case family.Firstborn:
		g.consequence(ConsequenceFamily, step, "the character is the firstborn", "")
	case family.Lastborn:
		g.consequence(ConsequenceFamily, step, "the character is the last child", "")
	}
}

// parentalAgeTable is p. 59's three columns. Nothing is printed below tech
// level 10 or above 13, and the engine reads the nearest printed column for
// the reason ERRATA E-16 gives for the aging tables.
func parentalAgeTable(techLevel int) (int, int) {
	switch {
	case techLevel <= 10:
		return 15, 0
	case techLevel == 11:
		return 23, -6
	}

	return 28, -12
}

// addFamily gains n relatives of one role as Allies at the rating the step
// grants them.
func (g *Generator) addFamily(role string, count, rating, cause int, detail string) {
	for range count {
		g.char.State.Ties = append(g.char.State.Ties, Tie{
			Kind:   string(career.Ally),
			Origin: FamilyOrigin,
			Rating: rating,
			Role:   role,
			Detail: detail,
		})
	}

	if count > 0 {
		g.consequence(ConsequenceRelationship, cause,
			"gain "+itoa(count)+" "+role+" as Allies at a Relationship Rating of "+
				itoa(rating), FamilyOrigin)
	}
}

// siblings is p. 60: 1d6-2 for a character born to a couple, 2d6 for one
// born into a communal situation, and then a 2d6 age table for each.
func (g *Generator) siblings(family *Family, step int) {
	count := g.siblingCount(family)
	if count == 0 {
		g.consequence(ConsequenceFamily, step, "an only child", "")

		return
	}

	twins := 0

	for range count {
		detail, isTwin := g.siblingAge(family)
		if isTwin {
			twins++

			detail = multipleBirth(twins)
		}

		g.char.State.Ties = append(g.char.State.Ties, Tie{
			Kind:   string(career.Ally),
			Origin: FamilyOrigin,
			Rating: closeFamilyRating,
			Role:   RoleSibling,
			Detail: detail,
		})
	}

	g.consequence(ConsequenceRelationship, step,
		"gain "+itoa(count)+" siblings as Allies at a Relationship Rating of "+
			itoa(closeFamilyRating), FamilyOrigin)
}

// siblingCount is the first sentence of p. 60. Note the asymmetry: a couple
// gives 1d6-2, which is nothing more than half the time, and a communal
// upbringing gives 2d6, which is never nothing.
func (g *Generator) siblingCount(family *Family) int {
	if family.Situation == SituationCommunal {
		roll := g.dice.TwoD6()
		g.log.Roll(roll, "p. 60")

		return roll.Total
	}

	roll := g.dice.D6()
	g.log.Roll(roll, "p. 60")

	return max(roll.Total-2, 0)
}

// siblingAge is the 2d6 table of p. 60, with the two re-roll rules folded
// in: a firstborn character has no older siblings and a last child has no
// younger ones, and the table says which result each redirects to rather
// than asking for a re-roll.
//
// The age difference is not checked against the parents' ages. p. 60 asks
// for that -- "if one of the parents will have been less than 16 years of
// age when the child will have been born, then reduce the age difference" --
// and it needs a parent's age at the sibling's birth rather than at the
// character's, which the step does not record. ERRATA E-25.
func (g *Generator) siblingAge(family *Family) (string, bool) {
	roll := g.dice.TwoD6()
	g.log.Roll(roll, "p. 60")

	result := redirectSiblingResult(roll.Total, family)

	switch result {
	case 2:
		return "died at childbirth", false
	case 3:
		return "younger by " + years(g.dice.ND10(3).Total), false
	case 4:
		return "younger by " + years(g.dice.TwoD6().Total), false
	case 10:
		return "older by " + years(g.dice.TwoD6().Total), false
	case 11:
		return "older by " + years(g.dice.ND10(3).Total), false
	case 12:
		return "", true
	}

	return g.nearSibling(family), false
}

// redirectSiblingResult applies the table's own redirects rather than
// re-rolling: result 3 becomes 11 for a youngest child and 11 becomes 3 for
// an eldest, 4 becomes 10 and 10 becomes 4.
func redirectSiblingResult(result int, family *Family) int {
	if family.Lastborn {
		switch result {
		case 3:
			return 11
		case 4:
			return 10
		}
	}

	if family.Firstborn {
		switch result {
		case 10:
			return 4
		case 11:
			return 3
		}
	}

	return result
}

// nearSibling is results 5-9: a d6 decides the direction and another the
// distance, unless the character's birth order has already decided which
// way it goes.
func (g *Generator) nearSibling(family *Family) string {
	older := g.dice.D6().Total%2 == 1

	switch {
	case family.Lastborn:
		older = true
	case family.Firstborn:
		older = false
	}

	apart := years(g.dice.D6().Total)
	if older {
		return "older by " + apart
	}

	return "younger by " + apart
}

// years renders an age difference, which is one year often enough to be
// worth not printing "1 years".
func years(n int) string {
	if n == 1 {
		return "1 year"
	}

	return itoa(n) + " years"
}

// multipleBirth names what result 12 makes of the character and their
// siblings: "If this is the second sibling to get this result, the character
// and the two siblings are triplets. A third result will mean quadruplets."
func multipleBirth(n int) string {
	switch n {
	case 1:
		return "the character's twin"
	case 2:
		return "one of three, with the character"
	case 3:
		return "one of four, with the character"
	}

	return "one of " + itoa(n+1) + ", with the character"
}

// extendedFamily is pp. 61's grandparents, aunts, uncles and cousins. The
// book builds them by repeating the birth situation for each parent; the
// engine repeats the counts rather than the narrative, because the second
// generation's household is not something the character's sheet records.
func (g *Generator) extendedFamily(step int) {
	// Two grandparents per parent, and their siblings are the character's
	// aunts and uncles.
	parents := len(g.char.State.Family.ParentAges)

	grandparents := 0
	auntsAndUncles := 0

	for range parents {
		grandparents += 2

		roll := g.dice.D6()
		g.log.Roll(roll, "p. 61")

		auntsAndUncles += max(roll.Total-2, 0)
	}

	g.addFamily(RoleGrandparent, grandparents, extendedFamilyRating, step, "a grandparent")
	g.addFamily(RoleAuntOrUncle, auntsAndUncles, extendedFamilyRating, step, "an aunt or uncle")

	cousins := 0

	for range auntsAndUncles {
		roll := g.dice.D6()
		g.log.Roll(roll, "p. 61")

		cousins += max(roll.Total-2, 0)
	}

	g.addFamily(RoleCousin, cousins, extendedFamilyRating, step, "a cousin")
}
