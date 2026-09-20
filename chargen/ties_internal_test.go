package chargen

import (
	"errors"
	"strings"
	"testing"

	"github.com/philoserf/cschargen/career"
)

// TestTheBandsOfPage320 is a second reading of the four bands and their
// boundaries:
//
//	Contacts are NPCs with a Relationship Rating of 1-100 ... NPCs that have
//	a Relationship Rating of 101-200 are Allies ... Rivals are NPCs with a
//	Relationship Rating of -1 through -100 ... Enemies are NPCs with a
//	Relationship Rating of -101 or less.
func TestTheBandsOfPage320(t *testing.T) {
	t.Parallel()

	tests := []struct {
		rating int
		want   career.Relationship
		held   bool
	}{
		{200, career.Ally, true},
		{101, career.Ally, true},
		{100, career.Contact, true},
		{1, career.Contact, true},
		// "Contacts which fall to a Relationship Rating of 0 have decided
		// to no longer be associated with the character and are lost."
		{0, "", false},
		{-1, career.Rival, true},
		{-100, career.Rival, true},
		{-101, career.Enemy, true},
		{-200, career.Enemy, true},
	}

	for _, tc := range tests {
		got, held := kindForRating(tc.rating)
		if got != tc.want || held != tc.held {
			t.Errorf("kindForRating(%d) = %q, %v; want %q, %v",
				tc.rating, got, held, tc.want, tc.held)
		}
	}
}

// TestABigFallIsALossNotADemotion. p. 320: "Allies which drop immediately to
// 0 or less are no longer associated with the character and are lost." The
// classification happens once, on the final rating, so an Ally who takes a
// loss big enough to pass through the Contact band and out the far side is
// gone rather than a Rival.
// ratingMove is one row of p. 320's movement rules.
type ratingMove struct {
	name       string
	start      Tie
	delta      int
	wantKind   career.Relationship
	wantRating int
	wantHeld   bool
}

func ratingMoves() []ratingMove {
	return []ratingMove{
		{
			name:  "an ally who becomes a contact",
			start: Tie{Kind: string(career.Ally), Rating: 110}, delta: -20,
			wantKind: career.Contact, wantRating: 90, wantHeld: true,
		},
		{
			name:  "an ally who is lost outright",
			start: Tie{Kind: string(career.Ally), Rating: 110}, delta: -120,
			wantHeld: false,
		},
		{
			name:  "a contact who becomes an ally",
			start: Tie{Kind: string(career.Contact), Rating: 90}, delta: 20,
			wantKind: career.Ally, wantRating: 110, wantHeld: true,
		},
		{
			name:  "an enemy who becomes a rival",
			start: Tie{Kind: string(career.Enemy), Rating: -110}, delta: 20,
			wantKind: career.Rival, wantRating: -90, wantHeld: true,
		},
		{
			name:  "an enemy who stops caring",
			start: Tie{Kind: string(career.Enemy), Rating: -110}, delta: 120,
			wantHeld: false,
		},
		{
			// "Allies cannot advance to higher than 200."
			name:  "an ally at the ceiling",
			start: Tie{Kind: string(career.Ally), Rating: 190}, delta: 50,
			wantKind: career.Ally, wantRating: 200, wantHeld: true,
		},
		{
			// "Enemies cannot drop to a rating lower than -200."
			name:  "an enemy at the floor",
			start: Tie{Kind: string(career.Enemy), Rating: -190}, delta: -50,
			wantKind: career.Enemy, wantRating: -200, wantHeld: true,
		},
	}
}

func TestABigFallIsALossNotADemotion(t *testing.T) {
	t.Parallel()

	for _, tc := range ratingMoves() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, held := adjustRating(tc.start, tc.delta)
			if held != tc.wantHeld {
				t.Fatalf("held = %v, want %v", held, tc.wantHeld)
			}

			if !held {
				return
			}

			if got.Kind != string(tc.wantKind) || got.Rating != tc.wantRating {
				t.Errorf("= %s at %d, want %s at %d",
					got.Kind, got.Rating, tc.wantKind, tc.wantRating)
			}
		})
	}
}

// TestNoTieIsBornLost is the bug this milestone opened with: every tie was
// created at a rating of 0, which p. 320 calls lost.
func TestNoTieIsBornLost(t *testing.T) {
	t.Parallel()

	for _, kind := range []career.Relationship{
		career.Ally, career.Contact, career.Rival, career.Enemy,
	} {
		rating := defaultRating(kind)

		got, held := kindForRating(rating)
		if !held {
			t.Errorf("a %s starts at %d, which p. 320 calls lost", kind, rating)

			continue
		}

		if got != kind {
			t.Errorf("a %s starts at %d, which p. 320 calls a %s", kind, rating, got)
		}
	}
}

// TestATieNeverCrossesZero holds the property the table above samples. Both
// halves of it are printed on p. 320, and together they mean a relationship
// ends rather than inverting: nobody's Enemy becomes their Contact by one
// good turn.
func TestATieNeverCrossesZero(t *testing.T) {
	t.Parallel()

	for start := ratingFloor; start <= ratingCeiling; start++ {
		if start == 0 {
			continue
		}

		for _, delta := range []int{-400, -200, -101, -1, 1, 101, 200, 400} {
			got, held := adjustRating(Tie{Rating: start}, delta)
			if !held {
				continue
			}

			if (got.Rating > 0) != (start > 0) {
				t.Fatalf("a tie at %d moved by %d and is still held at %d",
					start, delta, got.Rating)
			}
		}
	}
}

// TestTheLifeEventsThatMoveRatingsDoSo. Four results on p. 120 write to a
// Relationship Rating, and until this milestone every one of them recorded
// itself as unimplemented.
func TestTheLifeEventsThatMoveRatingsDoSo(t *testing.T) {
	t.Parallel()

	moves := 0

	var walk func(effects []career.Effect)

	walk = func(effects []career.Effect) {
		for _, effect := range effects {
			if effect.Kind == career.EffectRating || effect.Kind == career.EffectLoseTie {
				moves++
			}

			if effect.Kind == career.EffectUnimplemented &&
				strings.Contains(effect.Detail, "Relationship Rating") &&
				!strings.Contains(effect.Detail, "with no Ally or Contact") {
				t.Errorf("a Life Event still records a rating as unimplemented: %s",
					effect.Detail)
			}

			walk(effect.Success)
			walk(effect.Failure)

			for _, option := range effect.Options {
				walk(option.Effects)
			}
		}
	}

	for _, row := range career.LifeEvents() {
		walk(row.Effects)
	}

	// Results 4, 5, 6 and 11 -- and 12, which offers 11's upgrade as one of
	// three choices.
	if moves < 5 {
		t.Errorf("%d rating effects in the Life Events table; p. 120 has five", moves)
	}
}

// fromACareer is an Origin that is not FamilyOrigin, which is the only
// thing the family target cares about.
const fromACareer = "Colonist"

// tiedEngine is a generator with a set of ties already in hand, which is
// what the rating effects operate on.
func tiedEngine(t *testing.T, ties ...Tie) *Generator {
	t.Helper()

	gen := engine(t, 12)

	gen.char.State.Ties = ties

	return gen
}

func TestMovingEveryRelationshipAtOnce(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t,
		Tie{Kind: string(career.Ally), Origin: FamilyOrigin, Rating: 125},
		Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 40},
		Tie{Kind: string(career.Enemy), Origin: fromACareer, Rating: -150},
	)

	// Youth Path 1 result 7: "Lower the Relationship Rating with everyone
	// in your life by 25."
	err := gen.moveRatings(career.Effect{
		Kind: career.EffectRating, Target: career.TargetAll, Modifier: -25,
	}, 0)
	if err != nil {
		t.Fatalf("moveRatings: %v", err)
	}

	// The Contact at 40 falls to 15 and survives; the Enemy at -150 is
	// moved further from zero, not closer, because the delta is applied as
	// written rather than as a distance.
	want := []int{100, 15, -175}
	if len(gen.char.State.Ties) != len(want) {
		t.Fatalf("%d ties remain, want %d", len(gen.char.State.Ties), len(want))
	}

	for i, rating := range want {
		if gen.char.State.Ties[i].Rating != rating {
			t.Errorf("tie %d is at %d, want %d", i, gen.char.State.Ties[i].Rating, rating)
		}
	}

	// And the Ally at 125 is now a Contact at 100.
	if gen.char.State.Ties[0].Kind != string(career.Contact) {
		t.Errorf("an Ally who fell to 100 is still %s", gen.char.State.Ties[0].Kind)
	}
}

func TestMovingOnlyTheFamily(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t,
		Tie{Kind: string(career.Ally), Origin: FamilyOrigin, Rating: 125},
		Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 40},
	)

	// Youth Path 1 result 16: "Gain +20 to the Relationship Rating of all
	// family members."
	err := gen.moveRatings(career.Effect{
		Kind: career.EffectRating, Target: career.TargetFamily, Modifier: 20,
	}, 0)
	if err != nil {
		t.Fatalf("moveRatings: %v", err)
	}

	if gen.char.State.Ties[0].Rating != 145 {
		t.Errorf("the family tie is at %d, want 145", gen.char.State.Ties[0].Rating)
	}

	if gen.char.State.Ties[1].Rating != 40 {
		t.Errorf("a tie from a career moved on a family result")
	}
}

func TestMovingARelationshipNobodyHas(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t)

	err := gen.moveRatings(career.Effect{
		Kind: career.EffectRating, Target: career.TargetOne,
		Relationship: career.Ally, Modifier: -50,
	}, 0)
	if err != nil {
		t.Fatalf("moveRatings: %v", err)
	}

	if len(gen.char.State.Ties) != 0 {
		t.Error("a tie appeared out of nowhere")
	}
}

// TestChoosingWhichRelationshipMoves is the choice point POLICY.md records:
// with more than one candidate the Decider picks, and the policy takes the
// first.
func TestChoosingWhichRelationshipMoves(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t,
		Tie{Kind: string(career.Contact), Origin: "first", Rating: 60},
		Tie{Kind: string(career.Contact), Origin: "second", Rating: 60},
	)

	err := gen.moveRatings(career.Effect{
		Kind: career.EffectRating, Target: career.TargetOne,
		Relationship: career.Contact, Modifier: 50,
	}, 0)
	if err != nil {
		t.Fatalf("moveRatings: %v", err)
	}

	if gen.char.State.Ties[0].Rating != 110 || gen.char.State.Ties[1].Rating != 60 {
		t.Errorf("the policy moved %d and %d; it should have taken the first",
			gen.char.State.Ties[0].Rating, gen.char.State.Ties[1].Rating)
	}

	// 110 is over the line, so the Contact is an Ally now.
	if gen.char.State.Ties[0].Kind != string(career.Ally) {
		t.Errorf("a Contact raised to 110 is still %s", gen.char.State.Ties[0].Kind)
	}
}

// TestARolledRatingChangeTakesItsSignFromTheEffect. "an existing Ally or
// Contact loses 1d6 x 20 Relationship Rating" (Life Event 5) rolls for the
// amount and the page says which way it goes.
func TestARolledRatingChangeTakesItsSignFromTheEffect(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t, Tie{Kind: string(career.Ally), Origin: fromACareer, Rating: 200})

	err := gen.moveRatings(career.Effect{
		Kind: career.EffectRating, Target: career.TargetOne,
		Relationship: career.Ally, Modifier: -1, Dice: "1d6x20",
	}, 0)
	if err != nil {
		t.Fatalf("moveRatings: %v", err)
	}

	if len(gen.char.State.Ties) == 0 {
		t.Fatal("an Ally at 200 was lost to a roll of at most 120")
	}

	if got := gen.char.State.Ties[0].Rating; got >= 200 {
		t.Errorf("the Ally is at %d; a loss should have moved them down", got)
	}
}

func TestLosingATieTriesTheKindsInOrder(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t,
		Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 40},
		Tie{Kind: string(career.Enemy), Origin: fromACareer, Rating: -150},
	)

	// Life Event 4, the forward order: Ally first, then Contact.
	err := gen.loseTie(career.Effect{
		Kind:  career.EffectLoseTie,
		Order: []career.Relationship{career.Ally, career.Contact, career.Rival, career.Enemy},
	}, 0)
	if err != nil {
		t.Fatalf("loseTie: %v", err)
	}

	if len(gen.char.State.Ties) != 1 || gen.char.State.Ties[0].Kind != string(career.Enemy) {
		t.Errorf("the Contact should have gone before the Enemy; %d ties remain",
			len(gen.char.State.Ties))
	}

	// And the reverse order takes what is left.
	err = gen.loseTie(career.Effect{Kind: career.EffectLoseTie}, 0)
	if err != nil {
		t.Fatalf("loseTie: %v", err)
	}

	if len(gen.char.State.Ties) != 0 {
		t.Error("a default order left a tie behind")
	}

	// A third attempt with nothing left is not an error.
	err = gen.loseTie(career.Effect{Kind: career.EffectLoseTie}, 0)
	if err != nil {
		t.Fatalf("loseTie with nothing to lose: %v", err)
	}
}

// TestARefusedRelationshipChoiceEndsGeneration: choosing which tie a result
// moves goes through the Decider like every other choice, so an abandoned
// interactive session comes back as an error rather than a default.
func TestARefusedRelationshipChoiceEndsGeneration(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t,
		Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 60},
		Tie{Kind: string(career.Contact), Origin: "another", Rating: 60},
	)

	gen.decider = refusingDecider{}

	err := gen.moveRatings(career.Effect{
		Kind: career.EffectRating, Target: career.TargetOne,
		Relationship: career.Contact, Modifier: 10,
	}, 0)
	if err == nil {
		t.Fatal("a refused choice did not come back as an error")
	}
}

// TestDefaultRatingOfSomethingThatIsNotARelationship. Tie.Kind is a string
// in the record, so a hand-edited or future one can carry a word the four
// bands do not cover. It starts at the rating that means "lost", which is
// the safe answer for a relationship the engine cannot place.
func TestDefaultRatingOfSomethingThatIsNotARelationship(t *testing.T) {
	t.Parallel()

	if got := defaultRating("colleague"); got != ratingUnspecified {
		t.Errorf("defaultRating of an unknown kind = %d, want %d", got, ratingUnspecified)
	}
}

// TestARatingChangeCanLoseATie. The three ends of a rating change are all
// worth a consequence of their own: a move within a band, a move across
// one, and a move that ends the relationship (p. 320).
func TestARatingChangeCanLoseATie(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t,
		Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 20},
		Tie{Kind: string(career.Ally), Origin: fromACareer, Rating: 150},
	)

	// Youth Path 1 result 7 at its harshest: everyone falls, and the
	// Contact at 20 falls out of the scale.
	err := gen.moveRatings(career.Effect{
		Kind: career.EffectRating, Target: career.TargetAll, Modifier: -60,
	}, 0)
	if err != nil {
		t.Fatalf("moveRatings: %v", err)
	}

	if len(gen.char.State.Ties) != 1 {
		t.Fatalf("%d ties remain, want 1", len(gen.char.State.Ties))
	}

	// The Ally at 150 is a Contact at 90.
	remaining := gen.char.State.Ties[0]
	if remaining.Kind != string(career.Contact) || remaining.Rating != 90 {
		t.Errorf("the survivor is a %s at %d, want a contact at 90",
			remaining.Kind, remaining.Rating)
	}
}

// TestARatingChangeWithAnUnreadableAmount. The amount comes from the same
// expression language the rest of the tables are written in, so a
// mistranscribed one fails here the way it fails anywhere else -- and
// TestEveryTranscribedExpressionParses is what keeps it from reaching a
// record.
func TestARatingChangeWithAnUnreadableAmount(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t, Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 40})

	err := gen.moveRatings(career.Effect{
		Kind: career.EffectRating, Target: career.TargetOne,
		Relationship: career.Contact, Modifier: -1, Dice: "1d6 x 20",
	}, 0)
	if !errors.Is(err, ErrBadExpression) {
		t.Errorf("err = %v, want ErrBadExpression", err)
	}
}

// TestAimingAtARelativeNobodyHas. Two youth results address a parent
// specifically, and a character who was orphaned by result 2 has none.
func TestAimingAtARelativeNobodyHas(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t, Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 40})

	err := gen.moveRatings(career.Effect{
		Kind: career.EffectRating, Target: career.TargetRole,
		Role: "parent", Modifier: -150,
	}, 0)
	if err != nil {
		t.Fatalf("moveRatings: %v", err)
	}

	if len(gen.char.State.Ties) != 1 || gen.char.State.Ties[0].Rating != 40 {
		t.Error("a result aimed at a parent moved somebody who is not one")
	}
}

// TestLosingAFamilyMember is Teenage Path 2 result 3: "Choose a family
// member from your existing list and lose them." It is the only result in
// the book that removes a tie by what it is to the character rather than by
// its kind.
func TestLosingAFamilyMember(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t,
		Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 40},
		Tie{Kind: string(career.Ally), Origin: FamilyOrigin, Role: "sibling", Rating: 125},
	)

	lose := career.Effect{
		Kind: career.EffectLoseTie, Detail: "lose a family member",
		Target: career.TargetFamily,
	}

	err := gen.loseTie(lose, 0)
	if err != nil {
		t.Fatalf("loseTie: %v", err)
	}

	if len(gen.char.State.Ties) != 1 || gen.char.State.Ties[0].Origin != fromACareer {
		t.Error("the family member was not the one lost")
	}

	// And a character with no family left keeps what they have.
	err = gen.loseTie(lose, 0)
	if err != nil {
		t.Fatalf("loseTie with no family: %v", err)
	}

	if len(gen.char.State.Ties) != 1 {
		t.Error("a family result took somebody who is not family")
	}
}

// inCareerEngine is tiedEngine standing inside the career the ties came
// from, which is what "gained in this career" is measured against.
func inCareerEngine(t *testing.T, ties ...Tie) *Generator {
	t.Helper()

	gen := tiedEngine(t, ties...)
	colonist := career.Colonist()

	gen.career = &colonist

	return gen
}

// TestLosingSeveralTiesTriesTheOrderEachTime is "lose 1D3 Allies and
// Contacts": each removal works down the order afresh, so a character with
// one Ally and two Contacts losing three loses the Ally first.
func TestLosingSeveralTiesTriesTheOrderEachTime(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t,
		Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 40},
		Tie{Kind: string(career.Ally), Origin: fromACareer, Rating: 150},
		Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 60},
	)

	err := gen.loseTie(career.Effect{
		Kind: career.EffectLoseTie, Count: 2,
		Order: []career.Relationship{career.Ally, career.Contact},
	}, 0)
	if err != nil {
		t.Fatalf("loseTie: %v", err)
	}

	if len(gen.char.State.Ties) != 1 {
		t.Fatalf("%d ties remain, want 1", len(gen.char.State.Ties))
	}

	if gen.char.State.Ties[0].Kind != string(career.Contact) {
		t.Errorf("the surviving tie is a %s", gen.char.State.Ties[0].Kind)
	}
}

// TestLosingEveryTieOfThisCareerLeavesTheRest is "lose every Contact and
// Ally gained in this career": a Contact made at school is not in it.
func TestLosingEveryTieOfThisCareerLeavesTheRest(t *testing.T) {
	t.Parallel()

	gen := inCareerEngine(t,
		Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 40},
		Tie{Kind: string(career.Contact), Origin: FamilyOrigin, Rating: 60},
		Tie{Kind: string(career.Ally), Origin: fromACareer, Rating: 150},
	)

	err := gen.loseTie(career.Effect{
		Kind: career.EffectLoseTie, Count: -1, ThisCareer: true,
		Order: []career.Relationship{career.Contact, career.Ally},
	}, 0)
	if err != nil {
		t.Fatalf("loseTie: %v", err)
	}

	if len(gen.char.State.Ties) != 1 || gen.char.State.Ties[0].Origin != FamilyOrigin {
		t.Errorf("the family Contact did not survive: %v", gen.char.State.Ties)
	}
}

// TestLosingNobodyTakesTheFallback is "Lose one Ally or Contact. If you
// have none, gain an Enemy with a Relationship Rating of -110."
func TestLosingNobodyTakesTheFallback(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t)

	err := gen.loseTie(career.Effect{
		Kind:  career.EffectLoseTie,
		Order: []career.Relationship{career.Ally, career.Contact},
		Fallback: []career.Effect{{
			Kind: career.EffectRelationship, Relationship: career.Enemy,
			Count: 1, Rating: -110, Detail: "an Enemy at -110",
		}},
	}, 0)
	if err != nil {
		t.Fatalf("loseTie: %v", err)
	}

	if len(gen.char.State.Ties) != 1 || gen.char.State.Ties[0].Rating != -110 {
		t.Errorf("the fallback did not fire: %v", gen.char.State.Ties)
	}
}

// TestTheTieThatChangedSurvivesTheRest is Gambler mishap 11: "one Contact
// or Ally from this career becomes an Enemy, and every other relationship
// gained here is lost". The one that changed is the one that is left.
func TestTheTieThatChangedSurvivesTheRest(t *testing.T) {
	t.Parallel()

	gen := inCareerEngine(t,
		Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 40},
		Tie{Kind: string(career.Ally), Origin: fromACareer, Rating: 150},
		Tie{Kind: string(career.Rival), Origin: FamilyOrigin, Rating: -50},
	)

	err := gen.becomeTies(career.Effect{
		Kind: career.EffectBecome, Count: 1, ThisCareer: true, LoseTheRest: true,
		From: []career.Relationship{career.Contact, career.Ally}, Relationship: career.Enemy,
		Detail: "one becomes an Enemy",
	}, 0)
	if err != nil {
		t.Fatalf("becomeTies: %v", err)
	}

	if len(gen.char.State.Ties) != 2 {
		t.Fatalf("%d ties remain, want 2: %v", len(gen.char.State.Ties), gen.char.State.Ties)
	}

	if gen.char.State.Ties[0].Kind != string(career.Enemy) {
		t.Errorf("the changed tie is a %s", gen.char.State.Ties[0].Kind)
	}

	if gen.char.State.Ties[1].Origin != FamilyOrigin {
		t.Error("a tie from outside the career was taken with the rest")
	}
}

// TestAnImprovementReadsTheStateAsItWas is ERRATA E-34. An Enemy, a Rival
// and a Contact each move up one band; the Enemy does not climb to Ally by
// being found again as a Rival.
func TestAnImprovementReadsTheStateAsItWas(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t,
		Tie{Kind: string(career.Enemy), Origin: fromACareer, Rating: -150},
		Tie{Kind: string(career.Rival), Origin: fromACareer, Rating: -50},
		Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 40},
	)

	err := gen.improveTies(0)
	if err != nil {
		t.Fatalf("improveTies: %v", err)
	}

	want := []career.Relationship{career.Rival, career.Contact, career.Ally}
	for i, kind := range want {
		if gen.char.State.Ties[i].Kind != string(kind) {
			t.Errorf("tie %d is a %s, want %s", i, gen.char.State.Ties[i].Kind, kind)
		}
	}
}

// TestAnImprovementWithNoContactsGainsOne is the clause the other three sit
// beside: "If you have no Contacts, then you will gain one Contact."
func TestAnImprovementWithNoContactsGainsOne(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t, Tie{Kind: string(career.Ally), Origin: fromACareer, Rating: 150})

	err := gen.improveTies(0)
	if err != nil {
		t.Fatalf("improveTies: %v", err)
	}

	if len(gen.tiesOfKind(career.Contact)) != 1 {
		t.Errorf("no Contact was gained: %v", gen.char.State.Ties)
	}
}

// TestChangingNobodyTakesTheFallback is "a Rival becomes an Enemy, or a
// Rival is gained where there was none".
func TestChangingNobodyTakesTheFallback(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t)

	err := gen.becomeTies(career.Effect{
		Kind: career.EffectBecome, Count: 1, Relationship: career.Enemy,
		From:   []career.Relationship{career.Rival},
		Detail: "a Rival becomes an Enemy",
		Fallback: []career.Effect{{
			Kind: career.EffectRelationship, Relationship: career.Rival,
			Count: 1, Detail: "a Rival is gained",
		}},
	}, 0)
	if err != nil {
		t.Fatalf("becomeTies: %v", err)
	}

	if len(gen.tiesOfKind(career.Rival)) != 1 {
		t.Errorf("the fallback did not fire: %v", gen.char.State.Ties)
	}
}

// TestARolledNumberOfNonFamilyTies is Teenage Path 3 result 4: "1D6-2
// (minimum 1) of your Contacts or Allies who are not family members lose 50
// Relationship Rating".
func TestARolledNumberOfNonFamilyTies(t *testing.T) {
	t.Parallel()

	gen := tiedEngine(t,
		Tie{Kind: string(career.Contact), Origin: FamilyOrigin, Rating: 60, Role: "parent"},
		Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 60},
	)

	err := gen.moveRatings(career.Effect{
		Kind: career.EffectRating, Target: career.TargetAll, ExcludeFamily: true,
		From:      []career.Relationship{career.Contact, career.Ally},
		CountDice: "1d6-2", Minimum: 1, Modifier: -50,
		Detail: "non-family ties lose 50",
	}, 0)
	if err != nil {
		t.Fatalf("moveRatings: %v", err)
	}

	for _, tie := range gen.char.State.Ties {
		if tie.Origin == FamilyOrigin && tie.Rating != 60 {
			t.Errorf("a family tie moved to %d", tie.Rating)
		}
	}
}

// TestAnUnreadableTieCountIsAnError holds the three places a relationship
// result throws for how many: how many to lose, how many to change, and how
// many to move. An expression the parser cannot read stops generation
// rather than quietly doing nothing.
func TestAnUnreadableTieCountIsAnError(t *testing.T) {
	t.Parallel()

	for name, apply := range map[string]func(*Generator) error{
		"lose": func(gen *Generator) error {
			return gen.loseTie(career.Effect{
				Kind: career.EffectLoseTie, CountDice: "a few",
				Order: []career.Relationship{career.Contact},
			}, 0)
		},
		"become": func(gen *Generator) error {
			return gen.becomeTies(career.Effect{
				Kind: career.EffectBecome, CountDice: "a few",
				From: []career.Relationship{career.Contact}, Relationship: career.Ally,
			}, 0)
		},
		"move": func(gen *Generator) error {
			return gen.moveRatings(career.Effect{
				Kind: career.EffectRating, Target: career.TargetAll,
				CountDice: "a few", Modifier: -50,
			}, 0)
		},
	} {
		gen := tiedEngine(t, Tie{Kind: string(career.Contact), Origin: fromACareer, Rating: 40})

		err := apply(gen)
		if err == nil {
			t.Errorf("%s: an unreadable count was carried out anyway", name)
		}
	}
}
