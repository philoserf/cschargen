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

	indexes, err := g.ratingTargets(effect, cause)
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

	g.dropLostTies(cause)

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
func (g *Generator) dropLostTies(_ int) {
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
func (g *Generator) ratingTargets(effect career.Effect, cause int) ([]int, error) {
	switch effect.Target {
	case career.TargetAll:
		return allIndexes(len(g.char.State.Ties)), nil
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
		Prompt:  effect.Detail,
		Options: labels,
		Cite:    ratingCite,
	})
	if err != nil {
		return nil, err
	}

	_ = cause

	return []int{candidates[chosen]}, nil
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

// FamilyOrigin is the Origin a tie from Step 5 carries, which is how the
// youth tables reach "all family members".
const FamilyOrigin = "family"

func allIndexes(n int) []int {
	found := make([]int, n)
	for i := range found {
		found[i] = i
	}

	return found
}

// loseTie carries out an [career.EffectLoseTie]: it removes one tie, trying
// the kinds in the order the result prints them.
func (g *Generator) loseTie(effect career.Effect, cause int) error {
	order := effect.Order
	if len(order) == 0 {
		order = []career.Relationship{career.Ally, career.Contact, career.Rival, career.Enemy}
	}

	for _, kind := range order {
		candidates := g.tiesOfKind(kind)
		if len(candidates) == 0 {
			continue
		}

		index := candidates[0]
		lost := g.char.State.Ties[index]

		g.char.State.Ties = append(g.char.State.Ties[:index], g.char.State.Ties[index+1:]...)
		g.consequence(ConsequenceRelationship, cause,
			"lost the "+lost.Kind+" at "+itoa(lost.Rating), lost.Origin)

		return nil
	}

	g.consequence(ConsequenceRelationship, cause, effect.Detail+": nobody to lose", "")

	return nil
}

// ratingCite is the page Relationship Ratings are defined on.
const ratingCite = "p. 320"
