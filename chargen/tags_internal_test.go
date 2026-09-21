package chargen

import (
	"strings"
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

// TestAModifierSpentTwiceIsOneModifier is "take a -2 DM on your next two
// Advancement rolls": one modifier with two uses, not two modifiers.
func TestAModifierSpentTwiceIsOneModifier(t *testing.T) {
	t.Parallel()

	gen := engine(t, 137)

	gen.pending = []PendingModifier{
		{Applies: advancementThrow, Value: -2, Detail: "twice", Uses: 2},
	}

	for i := range 2 {
		got := gen.takeModifiers(advancementThrow)
		if len(got) != 1 {
			t.Fatalf("throw %d took %d modifiers, want 1", i+1, len(got))
		}
	}

	if got := gen.takeModifiers(advancementThrow); len(got) != 0 {
		t.Errorf("a third throw took %v, want nothing", got)
	}
}

// TestAModifierNamingACareerReachesThatCareerOnly is "+2 DM to enlistment
// in the Sports career", which is narrower than any class.
func TestAModifierNamingACareerReachesThatCareerOnly(t *testing.T) {
	t.Parallel()

	pending := PendingModifier{OnCareers: []string{"Sports"}}

	if pending.AppliesTo(career.Marine()) {
		t.Error("a modifier naming Sports reached the Marine career")
	}

	if !pending.AppliesTo(career.Sports()) {
		t.Error("a modifier naming Sports did not reach it")
	}
}

// TestASkillModifierLastsTheCareerAndIsNotSpent is "+1 DM to Melee checks
// in this career": every Melee check takes it, no other skill's does, and
// leaving the career ends it.
func TestASkillModifierLastsTheCareerAndIsNotSpent(t *testing.T) {
	t.Parallel()

	gen := engine(t, 139)

	gen.pending = []PendingModifier{{
		Applies: skillCheckThrow, Value: 1, Detail: "Melee in this career",
		OnSkill: "Melee", WhileInThisCareer: true,
	}}

	if got := gen.takeSkillModifiers("Stealth"); len(got) != 0 {
		t.Errorf("a Melee modifier reached a Stealth check: %v", got)
	}

	for i := range 2 {
		if got := gen.takeSkillModifiers("Melee"); len(got) != 1 {
			t.Errorf("Melee check %d took %d modifiers, want 1", i+1, len(got))
		}
	}

	gen.dropCareerModifiers()

	if got := gen.takeSkillModifiers("Melee"); len(got) != 0 {
		t.Errorf("the modifier outlived the career: %v", got)
	}
}

// TestASurvivalModifierReachesTheSurvivalRoll is the bug this pass was
// written for: twenty-five results grant one and nothing consumed it.
func TestASurvivalModifierReachesTheSurvivalRoll(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 0)

	gen.pending = []PendingModifier{
		{Applies: survivalThrow, Value: -2, Detail: "a hard term"},
	}

	_, err := gen.rollSurvival(gen.assignment)
	if err != nil {
		t.Fatalf("rollSurvival: %v", err)
	}

	for _, event := range gen.log.Events() {
		if event.Kind != EventThrow {
			continue
		}

		for _, mod := range event.Throw.Mods {
			if mod.Name == "a hard term" && mod.Value == -2 {
				return
			}
		}
	}

	t.Error("the survival throw did not carry the modifier granted to it")
}

// TestNothingWillHaveThemSaysWhy is the record naming the reason a career
// list came back empty, rather than leaving a step that happened and said
// nothing.
//
// Two things empty it, and writing this test is what showed that an aging
// crisis and a mental characteristic at 0 are not among them: eligibleCareers
// answers both with Vagabond, a list of one, which the ordinary choice point
// records.
func TestNothingWillHaveThemSaysWhy(t *testing.T) {
	t.Parallel()

	for name, arrange := range map[string]func(*Generator){
		"a career the flags asked for, closed": func(g *Generator) {
			g.forced = "Colonist"
			g.lockout = map[string]int{"Colonist": 9}
		},
		"every career has turned them down": func(g *Generator) {
			g.lockout = map[string]int{}
			for _, def := range career.All() {
				g.lockout[def.Name] = 9
			}
		},
	} {
		gen := engine(t, 163)
		arrange(gen)

		chosen, forced, err := gen.chooseCareer(0)
		if err != nil {
			t.Fatalf("%s: chooseCareer: %v", name, err)
		}

		if chosen.Name != career.Vagabond().Name || !forced {
			t.Errorf("%s: went to %q (forced %v), want Vagabond", name, chosen.Name, forced)
		}

		detail := lastDetail(t, gen)
		if !strings.Contains(detail, "no career will have them") {
			t.Errorf("%s: the record says %q", name, detail)
		}

		if !strings.Contains(detail, "(p. 111)") {
			t.Errorf("%s: the record does not cite the page: %q", name, detail)
		}
	}
}
