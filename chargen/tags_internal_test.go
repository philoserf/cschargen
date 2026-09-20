package chargen

import (
	"testing"

	"github.com/philoserf/cschargen/career"
)

// TestAnEnlistmentModifierReachesTheCareersItNames is the twelve results
// that modify an enlistment by the class of career being entered rather
// than by name: "-4 DM to enlist in any government related career", "-2 DM
// to enter any non-criminal career", "-2 DM to every enlistment roll that
// involves EDU".
//
// A modifier that names no class reaches every career, which is the
// ordinary case and the one every other modifier in the book is.
func TestAnEnlistmentModifierReachesTheCareersItNames(t *testing.T) {
	t.Parallel()

	marine := career.Marine()
	thief := career.Thief()
	vagabond := career.Vagabond()

	for name, test := range map[string]struct {
		pending PendingModifier
		target  career.Career
		want    bool
	}{
		"no narrowing reaches everything": {
			pending: PendingModifier{}, target: marine, want: true,
		},
		"a class the career holds": {
			pending: PendingModifier{OnTags: []career.Tag{career.TagGovernment}},
			target:  marine, want: true,
		},
		"a class the career does not hold": {
			pending: PendingModifier{OnTags: []career.Tag{career.TagGovernment}},
			target:  thief, want: false,
		},
		"a class excluded, and held": {
			pending: PendingModifier{NotTags: []career.Tag{career.TagCriminal}},
			target:  thief, want: false,
		},
		"a class excluded, and not held": {
			pending: PendingModifier{NotTags: []career.Tag{career.TagCriminal}},
			target:  marine, want: true,
		},
		"a career excluded by name": {
			pending: PendingModifier{NotCareers: []string{vagabond.Name}},
			target:  vagabond, want: false,
		},
		"an enlistment characteristic": {
			pending: PendingModifier{OnCharacteristics: []string{marine.Enlistment.Characteristic}},
			target:  marine, want: true,
		},
		"another enlistment characteristic": {
			pending: PendingModifier{OnCharacteristics: []string{"PSI"}},
			target:  marine, want: false,
		},
		"a career that throws for nothing": {
			pending: PendingModifier{OnCharacteristics: []string{"EDU"}},
			target:  vagabond, want: false,
		},
	} {
		got := test.pending.AppliesTo(test.target)
		if got != test.want {
			t.Errorf("%s: AppliesTo(%s) = %v, want %v", name, test.target.Name, got, test.want)
		}
	}
}

// TestAStandingModifierOutlivesTheThrow is the difference between "-2 DM to
// enter your next career" and "-2 DM to enter every career after this one",
// which several results print each of.
func TestAStandingModifierOutlivesTheThrow(t *testing.T) {
	t.Parallel()

	gen := engine(t, 107)

	gen.pending = []PendingModifier{
		{Applies: enlistmentThrow, Value: -2, Detail: "once", Standing: false},
		{Applies: enlistmentThrow, Value: -2, Detail: "forever", Standing: true},
	}

	first := gen.takeModifiersFor(enlistmentThrow, career.Marine())
	if len(first) != 2 {
		t.Fatalf("%d modifiers on the first throw, want 2", len(first))
	}

	second := gen.takeModifiersFor(enlistmentThrow, career.Marine())
	if len(second) != 1 || second[0].Name != "forever" {
		t.Errorf("the second throw took %v, want the standing one alone", second)
	}
}

// TestAModifierThatDoesNotReachThisCareerIsKept. A penalty aimed at
// government careers is not spent by enlisting in a criminal one.
func TestAModifierThatDoesNotReachThisCareerIsKept(t *testing.T) {
	t.Parallel()

	gen := engine(t, 109)

	gen.pending = []PendingModifier{{
		Applies: enlistmentThrow, Value: -4, Detail: "government only",
		OnTags: []career.Tag{career.TagGovernment},
	}}

	if got := gen.takeModifiersFor(enlistmentThrow, career.Thief()); len(got) != 0 {
		t.Errorf("a government penalty reached the Thief career: %v", got)
	}

	if got := gen.takeModifiersFor(enlistmentThrow, career.Marine()); len(got) != 1 {
		t.Errorf("the penalty was spent on a career it did not reach: %v", got)
	}
}

// TestAnAutomaticEnlistmentIsNarrowedToo is Undergraduate University's
// "you may enlist automatically in a business, military, corporate or
// colonist career" (p. 88): the automatic is not spent by attempting a
// career outside those classes.
func TestAnAutomaticEnlistmentIsNarrowedToo(t *testing.T) {
	t.Parallel()

	gen := engine(t, 113)

	gen.automatic = []PendingModifier{{
		Applies: enlistmentThrow, Detail: "four classes",
		OnTags: []career.Tag{career.TagBusiness, career.TagMilitary},
	}}

	if gen.takeAutomaticFor(enlistmentThrow, career.Thief()) {
		t.Error("an automatic enlistment reached a career outside its classes")
	}

	if !gen.takeAutomaticFor(enlistmentThrow, career.Marine()) {
		t.Error("an automatic enlistment did not reach a career inside its classes")
	}
}

// TestThePreviousCareerIsTheOneBeforeTheLast is Fringe Marketer mishap 8,
// "return to the career you held before this one". With no such career the
// name is empty, and chooseCareer reports a destination it cannot find
// rather than silently choosing one.
func TestThePreviousCareerIsTheOneBeforeTheLast(t *testing.T) {
	t.Parallel()

	gen := engine(t, 127)

	if got := gen.previousCareerName(); got != "" {
		t.Errorf("a character with no service has a previous career of %q", got)
	}

	gen.char.State.Services = []Service{{Career: "Colonist"}}
	if got := gen.previousCareerName(); got != "" {
		t.Errorf("a character in their first career has a previous career of %q", got)
	}

	gen.char.State.Services = []Service{{Career: "Colonist"}, {Career: "Fringe Marketer"}}
	if got := gen.previousCareerName(); got != "Colonist" {
		t.Errorf("the previous career is %q, want Colonist", got)
	}
}

// TestATransferToThePreviousCareerResolvesIt. The table says "the career
// you held before this one" and only the service record knows which that
// is, so the destination is resolved at the moment the transfer is taken.
func TestATransferToThePreviousCareerResolvesIt(t *testing.T) {
	t.Parallel()

	gen := engine(t, 131)

	gen.char.State.Services = []Service{{Career: "Colonist"}, {Career: "Fringe Marketer"}}
	gen.transfer = &career.Effect{
		Kind: career.EffectTransfer, Career: career.PreviousCareer,
		Detail: "re-enter the previous career",
	}

	chosen, forced, err := gen.chooseCareer(0)
	if err != nil {
		t.Fatalf("chooseCareer: %v", err)
	}

	if !forced || chosen.Name != "Colonist" {
		t.Errorf("the transfer resolved to %q (forced %v), want Colonist", chosen.Name, forced)
	}
}
