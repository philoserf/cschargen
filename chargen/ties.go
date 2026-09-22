package chargen

import "github.com/philoserf/cschargen/career"

// Relationship Ratings (p. 320).
//
// "Each NPC should be given a Relationship Rating. The Relationship Rating
// reflects the strength of the relationship between the character and the
// NPC on a scale of -200 through 200."
//
// Four bands divide that scale, and every one of them has an exit: a
// Contact that falls to 0 is lost, one that rises past 100 becomes an Ally,
// an Enemy that rises above -100 becomes a Rival.
//
// The band governs movement, not birth. A result that grants a tie names
// its kind, and the engine keeps that name even where the rating disagrees:
// Step 5 grants a grandparent "as an Ally with a Relationship Rating of 100"
// (p. 61) and 100 is the top of the Contact band. ERRATA E-22.

// The bounds of p. 320's four bands, and the scale they sit on.
const (
	ratingCeiling     = 200
	ratingFloor       = -200
	allyFloor         = 101
	contactCeiling    = 100
	contactFloor      = 1
	rivalCeiling      = -1
	rivalFloor        = -100
	enemyCeiling      = -101
	ratingUnspecified = 0
)

// kindForRating is the band a rating falls in, and whether it falls in one
// at all. A rating of 0 is in no band: "Contacts which fall to a
// Relationship Rating of 0 have decided to no longer be associated with the
// character and are lost."
func kindForRating(rating int) (career.Relationship, bool) {
	switch {
	case rating >= allyFloor:
		return career.Ally, true
	case rating >= contactFloor:
		return career.Contact, true
	case rating <= enemyCeiling:
		return career.Enemy, true
	case rating <= rivalCeiling:
		return career.Rival, true
	}

	return "", false
}

// defaultRating is where a tie starts when the result that granted it does
// not say. Most career results do not: "Gain an Ally" is the whole of it.
//
// The book gives no number for that case, so the engine takes the middle of
// the band, which is the only choice that is not an argument about which end
// of it a nameless relationship belongs at. ERRATA E-23.
func defaultRating(kind career.Relationship) int {
	switch kind {
	case career.Ally:
		return midpoint(allyFloor, ratingCeiling)
	case career.Contact:
		return midpoint(contactFloor, contactCeiling)
	case career.Rival:
		return midpoint(rivalFloor, rivalCeiling)
	case career.Enemy:
		return midpoint(ratingFloor, enemyCeiling)
	}

	return ratingUnspecified
}

// midpoint is the middle of a band, rounded toward zero.
func midpoint(low, high int) int {
	const halves = 2

	return (low + high) / halves
}

// adjustRating moves a tie's rating and returns what became of it: the kind
// it now is, and whether the character still has it.
//
// A tie never crosses zero. p. 320 says so from both sides -- "Allies which
// drop immediately to 0 or less are no longer associated with the character
// and are lost", "Enemies whose rating rises immediately to 0 or higher have
// decided they no longer care about the character and are lost" -- so an
// Ally at 110 who takes -120 is gone rather than a Rival at -10, and an
// Enemy at -110 who is forgiven by 120 is gone rather than a Contact at 10.
//
// Within a sign the bands apply as printed, and the classification happens
// once on the final rating rather than band by band.
func adjustRating(tie Tie, delta int) (Tie, bool) {
	before := tie.Rating

	tie.Rating = min(max(before+delta, ratingFloor), ratingCeiling)

	if tie.Rating == 0 || (before > 0) != (tie.Rating > 0) {
		return tie, false
	}

	// A rating that is neither zero nor across zero is inside a band, and
	// kindForRating's "no band" answer is unreachable from here: it exists
	// for the rating 0 the line above has already taken.
	kind, _ := kindForRating(tie.Rating)

	tie.Kind = string(kind)

	return tie, true
}

// moveRatings carries out an [career.EffectRating]: it moves one tie, every
// tie, or every family tie, and drops any that fall out of their band.
func (g *Generator) moveRatings(effect career.Effect, cause int) error {
	delta := effect.Modifier

	if effect.Dice != "" {
		rolled, err := g.rollExpression(effect.Dice, g.cite)
		if err != nil {
			return err
		}

		if delta < 0 {
			rolled = -rolled
		}

		delta = rolled
	}

	indexes, err := g.ratingTargets(effect)
	if err != nil {
		return err
	}

	if len(indexes) == 0 {
		g.consequence(ConsequenceRelationship, cause,
			effect.Detail+": nobody to move", "")

		return nil
	}

	for _, index := range indexes {
		g.moveOneRating(index, delta, cause)
	}

	g.dropLostTies()

	return nil
}

// moveOneRating applies the change to a single tie and logs what it did.
// A tie that falls out of every band is marked rather than removed, so the
// indexes the caller is walking stay valid.
func (g *Generator) moveOneRating(index, delta, cause int) {
	before := g.char.State.Ties[index]

	after, held := adjustRating(before, delta)

	g.char.State.Ties[index] = after

	switch {
	case !held:
		g.char.State.Ties[index].Kind = lostTie
		g.consequence(ConsequenceRelationship, cause,
			"the "+before.Kind+" at "+itoa(before.Rating)+" is lost", before.Origin)
	case after.Kind != before.Kind:
		g.consequence(ConsequenceRelationship, cause,
			"the "+before.Kind+" at "+itoa(before.Rating)+" is now a "+
				after.Kind+" at "+itoa(after.Rating), before.Origin)
	default:
		g.consequence(ConsequenceRelationship, cause,
			before.Kind+" "+itoa(before.Rating)+" -> "+itoa(after.Rating), before.Origin)
	}
}

// lostTie marks a tie between falling out of its band and being removed.
// It is never a value a record carries: dropLostTies runs before the
// effect returns.
const lostTie = ""

// dropLostTies removes what moveOneRating marked.
func (g *Generator) dropLostTies() {
	kept := g.char.State.Ties[:0]

	for _, tie := range g.char.State.Ties {
		if tie.Kind != lostTie {
			kept = append(kept, tie)
		}
	}

	g.char.State.Ties = kept
}

// ratingTargets is which ties an effect moves. TargetOne is a choice point,
// narrowed to one kind where the page narrows it.
func (g *Generator) ratingTargets(effect career.Effect) ([]int, error) {
	switch effect.Target {
	case career.TargetAll:
		return g.narrowedIndexes(effect)
	case career.TargetFamily:
		return g.familyIndexes(), nil
	case career.TargetAllOfRole:
		return g.roleIndexes(effect.Role), nil
	case career.TargetRole:
		found := g.roleIndexes(effect.Role)
		if len(found) == 0 {
			return nil, nil
		}

		// The page says "choose one of your parents", and the policy takes
		// the first. A relative is not distinguishable from another of the
		// same role in the record, so there is nothing to choose between.
		return found[:1], nil
	case career.TargetOne:
	}

	candidates := g.tiesOfKind(effect.Relationship)
	if len(candidates) <= 1 {
		return candidates, nil
	}

	labels := make([]string, len(candidates))
	for i, index := range candidates {
		tie := g.char.State.Ties[index]

		labels[i] = tie.Kind + " from " + tie.Origin + " at " + itoa(tie.Rating)
	}

	chosen, err := g.choose(Choice{
		Point:   "relationship_target",
		Prompt:  asPrompt(effect.Detail),
		Options: labels,
		Cite:    ratingCite,
	})
	if err != nil {
		return nil, err
	}

	return []int{candidates[chosen]}, nil
}

// narrowedIndexes is TargetAll, with the two narrowings Teenage Path 3
// result 4 puts on it: "1D6-2 (minimum 1) of your Contacts or Allies who
// are not family members". With neither it is every tie, which is what
// "everyone in your life" means.
func (g *Generator) narrowedIndexes(effect career.Effect) ([]int, error) {
	found := make([]int, 0, len(g.char.State.Ties))

	for i, tie := range g.char.State.Ties {
		if effect.ExcludeFamily && tie.Origin == FamilyOrigin {
			continue
		}

		if len(effect.From) > 0 && !matchesAnyKind(tie, effect.From) {
			continue
		}

		found = append(found, i)
	}

	if effect.CountDice == "" {
		return found, nil
	}

	rolled, err := g.rollExpression(effect.CountDice, ratingCite)
	if err != nil {
		return nil, err
	}

	return found[:min(max(rolled, effect.Minimum), len(found))], nil
}

// matchesAnyKind reports whether a tie is one of the kinds given.
func matchesAnyKind(tie Tie, kinds []career.Relationship) bool {
	for _, kind := range kinds {
		if tie.Kind == string(kind) {
			return true
		}
	}

	return false
}

// tiesOfKind is the indexes of every tie of one kind, or of every tie where
// the effect names none.
func (g *Generator) tiesOfKind(only career.Relationship) []int {
	var found []int

	for i, tie := range g.char.State.Ties {
		if only == "" || tie.Kind == string(only) {
			found = append(found, i)
		}
	}

	return found
}

// familyIndexes is every tie Step 5 granted, which the family step marks
// by origin.
func (g *Generator) familyIndexes() []int {
	var found []int

	for i, tie := range g.char.State.Ties {
		if tie.Origin == FamilyOrigin {
			found = append(found, i)
		}
	}

	return found
}

// roleIndexes is every family tie of one role.
func (g *Generator) roleIndexes(role string) []int {
	var found []int

	for i, tie := range g.char.State.Ties {
		if tie.Role == role {
			found = append(found, i)
		}
	}

	return found
}

// The Origins a tie carries when it did not come from a career. Every
// relationship in the record names where it came from, and a tie
// granted before Step 9 has no career to name -- so it names the stage of
// life it was made in, which is what a reader wants of a childhood mentor
// or a school rival.
//
// FamilyOrigin is also how the youth tables reach "all family members".
const (
	FamilyOrigin    = "family"
	YouthOrigin     = "youth"
	TeenageOrigin   = "teenage"
	ChildhoodOrigin = "childhood"
)

// everyTie is the Count a result carries when it means all of them.
const everyTie = -1

// loseTie carries out an [career.EffectLoseTie]: it removes ties, trying
// the kinds in the order the result prints them, as many times as the
// result says.
func (g *Generator) loseTie(effect career.Effect, cause int) error {
	// A result may name the family rather than a kind: "Choose a family
	// member from your existing list and lose them" (p. 78).
	if effect.Target == career.TargetFamily {
		family := g.familyIndexes()
		if len(family) == 0 {
			g.consequence(ConsequenceRelationship, cause, effect.Detail+": nobody to lose", "")

			return nil
		}

		// The page says "choose a family member", and the policy takes the
		// first: a relative is not distinguishable from another of the same
		// role in the record, so there is nothing to choose between.
		g.dropTie(family[0], cause)

		return nil
	}

	order := effect.Order
	if len(order) == 0 {
		order = []career.Relationship{career.Ally, career.Contact, career.Rival, career.Enemy}
	}

	count, err := g.tieCount(effect)
	if err != nil {
		return err
	}

	taken := 0

	for count == everyTie || taken < count {
		index, found := g.firstOfAnyKind(order, effect)
		if !found {
			break
		}

		g.dropTie(index, cause)

		taken++
	}

	if taken > 0 {
		return nil
	}

	if len(effect.Fallback) > 0 {
		g.consequence(ConsequenceRelationship, cause, effect.Detail+": nobody to lose", "")

		return g.applyAll(effect.Fallback, cause)
	}

	g.consequence(ConsequenceRelationship, cause, effect.Detail+": nobody to lose", "")

	return nil
}

// tieCount is how many ties a result takes or changes: a printed number, a
// thrown one, or everyTie.
func (g *Generator) tieCount(effect career.Effect) (int, error) {
	if effect.CountDice != "" {
		rolled, err := g.rollExpression(effect.CountDice, ratingCite)
		if err != nil {
			return 0, err
		}

		// "1D6-3 Contacts" can throw at or below zero, and a result that
		// takes a negative number of people does nothing rather than
		// reversing itself.
		return max(rolled, 0), nil
	}

	if effect.Count == everyTie {
		return everyTie, nil
	}

	// A result that names no number names one: "Lose one Contact or Ally".
	return max(effect.Count, 1), nil
}

// firstOfAnyKind is the index of a tie matching the kinds given, working
// through them in the order given. The narrowing fields decide which one:
// ThisCareer and ExcludeFamily rule ties out, and Newest takes the last
// match rather than the first, which is what "that Ally" means (E-39).
func (g *Generator) firstOfAnyKind(order []career.Relationship, effect career.Effect) (int, bool) {
	for _, kind := range order {
		found, ok := -1, false

		for i, tie := range g.char.State.Ties {
			if tie.Kind != string(kind) {
				continue
			}

			if effect.ThisCareer && !g.fromThisCareer(tie) {
				continue
			}

			if effect.ExcludeFamily && tie.Origin == FamilyOrigin {
				continue
			}

			found, ok = i, true

			// Newest keeps looking; the ordinary case stops at the first.
			if !effect.Newest {
				break
			}
		}

		if ok {
			return found, true
		}
	}

	return 0, false
}

// fromThisCareer reports whether a tie came from the career being served.
func (g *Generator) fromThisCareer(tie Tie) bool {
	return g.service.career != nil && tie.Origin == g.service.career.Name
}

// dropTie removes one tie by index and records what went.
func (g *Generator) dropTie(index, cause int) {
	lost := g.char.State.Ties[index]

	g.char.State.Ties = append(g.char.State.Ties[:index], g.char.State.Ties[index+1:]...)

	detail := "lost the " + lost.Kind + " at " + itoa(lost.Rating)
	if lost.Role != "" {
		detail = "lost a " + lost.Role + ", the " + lost.Kind + " at " + itoa(lost.Rating)
	}

	g.consequence(ConsequenceRelationship, cause, detail, lost.Origin)
}

// becomeTies carries out an [career.EffectBecome]: ties of one kind become
// another. The new rating is the middle of the new band, because a change
// of kind is not a change of rating and the book gives no number -- the
// same argument as ERRATA E-23.
func (g *Generator) becomeTies(effect career.Effect, cause int) error {
	count, err := g.tieCount(effect)
	if err != nil {
		return err
	}

	changed := 0
	skip := map[int]bool{}

	for count == everyTie || changed < count {
		index, found := g.firstOfAnyKindExcept(effect.From, effect.ThisCareer, skip)
		if !found {
			break
		}

		g.turnTie(index, effect.Relationship, cause)

		skip[index] = true
		changed++
	}

	if changed > 0 {
		if effect.LoseTheRest {
			g.loseEveryOtherTieInCareer(skip, cause)
		}

		return nil
	}

	g.consequence(ConsequenceRelationship, cause, effect.Detail+": nobody to change", "")

	if len(effect.Fallback) > 0 {
		return g.applyAll(effect.Fallback, cause)
	}

	return nil
}

// loseEveryOtherTieInCareer is the second half of Gambler mishap 11: every
// relationship this career gave goes except the ones the result has just
// changed, which are the ones that survive it.
func (g *Generator) loseEveryOtherTieInCareer(spared map[int]bool, cause int) {
	kept := make([]Tie, 0, len(g.char.State.Ties))

	for i, tie := range g.char.State.Ties {
		if spared[i] || !g.fromThisCareer(tie) {
			kept = append(kept, tie)

			continue
		}

		g.consequence(ConsequenceRelationship, cause,
			"lost the "+tie.Kind+" at "+itoa(tie.Rating), tie.Origin)
	}

	g.char.State.Ties = kept
}

// firstOfAnyKindExcept is firstOfAnyKind with a set of indexes already
// used. A tie that has just become an Enemy must not be found again as an
// Enemy to change: the result names the kinds as they were.
func (g *Generator) firstOfAnyKindExcept(
	order []career.Relationship, thisCareer bool, skip map[int]bool,
) (int, bool) {
	for _, kind := range order {
		for i, tie := range g.char.State.Ties {
			if skip[i] || tie.Kind != string(kind) {
				continue
			}

			if thisCareer && !g.fromThisCareer(tie) {
				continue
			}

			return i, true
		}
	}

	return 0, false
}

// turnTie changes one tie's kind and reseats its rating in the new band.
func (g *Generator) turnTie(index int, to career.Relationship, cause int) {
	before := g.char.State.Ties[index]

	g.char.State.Ties[index].Kind = string(to)
	g.char.State.Ties[index].Rating = defaultRating(to)

	g.consequence(ConsequenceRelationship, cause,
		"the "+before.Kind+" at "+itoa(before.Rating)+" is now a "+string(to)+
			" at "+itoa(defaultRating(to)), before.Origin)
}

// improveTies is the four clauses pp. 85, 91 and 96 print together:
//
//	"If you have no Contacts, then you will gain one Contact. If you
//	currently have Enemies, one of those is now a Rival. If you have
//	Rivals, one of those is now a Contact. If you have Contacts, one of
//	those is now an Ally."
//
// They are read against the state as it was, not one after another: a
// sequential reading would let a single Enemy climb to Ally on a result the
// book calls "an improvement to a relationship". ERRATA E-34.
func (g *Generator) improveTies(cause int) error {
	steps := []struct {
		from career.Relationship
		to   career.Relationship
	}{
		{career.Enemy, career.Rival},
		{career.Rival, career.Contact},
		{career.Contact, career.Ally},
	}

	chosen := make([]int, 0, len(steps))
	taken := map[int]bool{}
	hadAContact := len(g.tiesOfKind(career.Contact)) > 0

	for _, step := range steps {
		index, found := g.firstOfAnyKindExcept([]career.Relationship{step.from}, false, taken)
		if !found {
			chosen = append(chosen, -1)

			continue
		}

		chosen = append(chosen, index)
		taken[index] = true
	}

	for i, index := range chosen {
		if index >= 0 {
			g.turnTie(index, steps[i].to, cause)
		}
	}

	if hadAContact {
		return nil
	}

	return g.gainTies(career.Effect{
		Kind: career.EffectRelationship, Relationship: career.Contact, Count: 1,
		Detail: "with no Contacts, one is gained",
	}, cause)
}

// ratingCite is the page Relationship Ratings are defined on.
const ratingCite = "p. 320"

// tieOrigin is where a relationship granted right now came from: the career
// being served, or -- before Step 9, where there is no career -- the stage
// of life the character is in.
//
// Every relationship the record carries names its origin, and a blank one
// is the record failing to. A childhood mentor and a parent are
// both Allies at 125, and only the origin tells them apart.
func (g *Generator) tieOrigin() string {
	if g.service.career != nil {
		return g.service.career.Name
	}

	return g.stage
}
