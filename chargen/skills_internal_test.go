package chargen

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/setting"
)

// TestRaisingAHeldSkillWithNoneHeld is reachable: the result appears on
// tables a character can meet before any skill has been granted, and a
// choice with no options is not a question to ask.
func TestRaisingAHeldSkillWithNoneHeld(t *testing.T) {
	t.Parallel()

	gen := engine(t, 61)

	err := gen.raiseHeldSkill(career.Effect{
		Kind: career.EffectRaiseHeld, Detail: "raise a skill already held",
	}, 0)
	if err != nil {
		t.Fatalf("raiseHeldSkill: %v", err)
	}

	if !strings.Contains(lastDetail(t, gen), "holds none") {
		t.Errorf("the record does not say why nothing happened: %q", lastDetail(t, gen))
	}
}

// TestRaisingAHeldSkillTakesItsSpecialty. "A skill you already possess" is
// one entry of the sheet, not a name that might have several: raising
// Melee (Blade) must not create a bare Melee beside it.
func TestRaisingAHeldSkillTakesItsSpecialty(t *testing.T) {
	t.Parallel()

	gen := engine(t, 67)

	gen.char.State.GainSkill("Melee", "Blade", 1)

	err := gen.raiseHeldSkill(career.Effect{
		Kind: career.EffectRaiseHeld, Detail: "raise a skill already held",
	}, 0)
	if err != nil {
		t.Fatalf("raiseHeldSkill: %v", err)
	}

	if len(gen.char.State.Skills) != 1 {
		t.Fatalf("%d skills, want 1: %v", len(gen.char.State.Skills), gen.char.State.Skills)
	}

	if gen.char.State.Skills[0].Level != 2 {
		t.Errorf("Melee (Blade) is at %d, want 2", gen.char.State.Skills[0].Level)
	}
}

// TestARefusedSkillChoiceComesBack holds both choice points: a decider that
// answers nothing stops generation rather than picking for the player.
func TestARefusedSkillChoiceComesBack(t *testing.T) {
	t.Parallel()

	for name, effect := range map[string]career.Effect{
		"held": {Kind: career.EffectRaiseHeld, Detail: "raise a skill already held"},
		"any":  {Kind: career.EffectAnySkill, Detail: "any skill"},
	} {
		gen := engine(t, 71)

		gen.char.State.GainSkill("Melee", "Blade", 1)

		gen.decider = refusingDecider{}

		err := gen.apply(effect, 0)
		if err == nil {
			t.Errorf("%s: a refused choice did not come back", name)
		}
	}
}

// TestAnySkillNotHeldWithEverythingHeld is the far end of "any skill you do
// not already possess": a character who holds all forty-four gains nothing,
// and the record says which it was.
func TestAnySkillNotHeldWithEverythingHeld(t *testing.T) {
	t.Parallel()

	gen := engine(t, 73)

	for _, definition := range career.Skills() {
		gen.char.State.GainSkill(definition.Name, "", 1)
	}

	err := gen.anySkill(career.Effect{
		Kind: career.EffectAnySkill, NotHeld: true, Detail: "any skill not held",
	}, 0)
	if err != nil {
		t.Fatalf("anySkill: %v", err)
	}

	if !strings.Contains(lastDetail(t, gen), "already holds every skill") {
		t.Errorf("the record does not say why nothing happened: %q", lastDetail(t, gen))
	}
}

// TestLosingASkillNobodyHas is Orbital Construction's sixth mishap meeting
// a character who never had a vacc suit. The result fired, so the record
// carries it.
func TestLosingASkillNobodyHas(t *testing.T) {
	t.Parallel()

	gen := engine(t, 79)

	gen.loseSkill(career.Effect{
		Kind: career.EffectLoseSkill, Skill: "Suit", Specialties: []string{"Vacc Suit"},
		Detail: "lose every level of Suit (Vacc Suit)",
	}, 0)

	if !strings.Contains(lastDetail(t, gen), "did not hold") {
		t.Errorf("the record does not say the skill was absent: %q", lastDetail(t, gen))
	}
}

// TestAHomeworldSpecialtyComesFromTheChart is the five youth and teenage
// results that ask for "a Survival specialty which is used on your
// homeworld": the options are the ones the world's own background skills
// name, and nowhere else.
func TestAHomeworldSpecialtyComesFromTheChart(t *testing.T) {
	t.Parallel()

	world := testWorld("Tinderfall", nil)

	world.BackgroundSkills = []setting.Requirement{
		{OneOf: []setting.Alternative{
			{Skill: "Survival", Specialties: []string{"Desert", "Heat"}},
			{Skill: "Recon"},
		}},
		{OneOf: []setting.Alternative{{Skill: "Survival", Specialties: []string{"Desert"}}}},
	}

	gen := engine(t, 89)

	gen.setting = settingWith(setting.Subsector{Name: "New Holdings", Worlds: []setting.World{world}})
	gen.char.State.Homeworlds = []Homeworld{{World: "Tinderfall"}}

	options := gen.specialtiesOnTheHomeworld("Survival")
	if !slices.Equal(options, []string{"Desert", "Heat"}) {
		t.Errorf("the homeworld offers %v, want [Desert Heat] once each", options)
	}

	err := gen.homeworldSpecialty(career.Effect{
		Kind: career.EffectHomeworldSkill, Skill: "Survival", Detail: "a Survival specialty",
	}, 0)
	if err != nil {
		t.Fatalf("homeworldSpecialty: %v", err)
	}

	if gen.char.State.SkillLevel("Survival") != 1 {
		t.Errorf("no Survival specialty was granted: %v", gen.char.State.Skills)
	}
}

// TestAHomeworldWithNoSuchSpecialty records rather than invents. Several
// worlds' background skills name no Survival at all; a character may have
// no homeworld yet; and a record may name a world the loaded data does not
// carry. None of the three is an error, and all three say so.
func TestAHomeworldWithNoSuchSpecialty(t *testing.T) {
	t.Parallel()

	for name, homeworlds := range map[string][]Homeworld{
		"none yet":    nil,
		"not in data": {{World: "Somewhere Else"}},
		"no Survival": {{World: "Tinderfall"}},
	} {
		gen := engine(t, 83)

		gen.setting = settingWith(setting.Subsector{
			Name: "New Holdings", Worlds: []setting.World{testWorld("Tinderfall", nil)},
		})
		gen.char.State.Homeworlds = homeworlds

		err := gen.homeworldSpecialty(career.Effect{
			Kind: career.EffectHomeworldSkill, Skill: "Survival",
			Detail: "a Survival specialty from the homeworld",
		}, 0)
		if err != nil {
			t.Fatalf("%s: homeworldSpecialty: %v", name, err)
		}

		if !strings.Contains(lastDetail(t, gen), "name none") {
			t.Errorf("%s: the record does not say why nothing happened: %q",
				name, lastDetail(t, gen))
		}
	}
}

// TestRaisingASkillHeldAtLevelZero. "Any skill you already possess"
// includes one held at level 0: p. 117's first term grants a whole table at
// zero, and those are possessed. The raise must take such an entry to 1
// rather than leaving it where it was.
func TestRaisingASkillHeldAtLevelZero(t *testing.T) {
	t.Parallel()

	gen := engine(t, 97)

	gen.char.State.GainSkill("Melee", "Blade", 0)

	err := gen.raiseHeldSkill(career.Effect{
		Kind: career.EffectRaiseHeld, Detail: "raise a skill already held",
	}, 0)
	if err != nil {
		t.Fatalf("raiseHeldSkill: %v", err)
	}

	if gen.char.State.SkillLevel("Melee") != 1 {
		t.Errorf("Melee (Blade) is at %d, want 1: %v",
			gen.char.State.SkillLevel("Melee"), gen.char.State.Skills)
	}
}

// TestACharacterIsAddictedOnlyOnce. A character already addicted to
// alcohol who takes the result again is no more addicted than before, and
// a second entry on the sheet would read as a second addiction.
func TestACharacterIsAddictedOnlyOnce(t *testing.T) {
	t.Parallel()

	gen := engine(t, 149)

	addicted := career.Effect{
		Kind:      career.EffectCondition,
		Condition: "an addiction to alcohol or a drug of the character's choice",
		Detail:    "an addiction",
	}

	gen.addCondition(addicted, 0)
	gen.addCondition(addicted, 0)

	if len(gen.char.State.Conditions) != 1 {
		t.Errorf("the character carries %v, want one entry", gen.char.State.Conditions)
	}

	if !strings.Contains(lastDetail(t, gen), "already carries") {
		t.Errorf("the record does not say it was already there: %q", lastDetail(t, gen))
	}
}

// TestRaisingTheSkillTheCharacterUsed is the Scientist's breakthrough:
// "gain a level in the Science skill you used on this task" is one of the
// character's Sciences, not any skill they hold.
func TestRaisingTheSkillTheCharacterUsed(t *testing.T) {
	t.Parallel()

	gen := engine(t, 151)

	gen.char.State.GainSkill("Melee", "Blade", 3)
	gen.char.State.GainSkill("Science", "Physics", 2)

	err := gen.raiseHeldSkill(career.Effect{
		Kind: career.EffectRaiseHeld, OnSkill: "Science",
		Detail: "raise the Science skill the character used",
	}, 0)
	if err != nil {
		t.Fatalf("raiseHeldSkill: %v", err)
	}

	if gen.char.State.SkillLevel("Science") != 3 {
		t.Errorf("Science is at %d, want 3: %v",
			gen.char.State.SkillLevel("Science"), gen.char.State.Skills)
	}

	if gen.char.State.SkillLevel("Melee") != 3 {
		t.Error("the raise reached a skill outside the one it named")
	}
}

// anyResults are the cases TestAnAnyResultIsResolved runs.
var anyResults = []struct {
	name    string
	table   bool
	held    []Skill
	grant   career.Effect
	want    []Skill
	stamped bool
}{
	{
		// "If the character does not already have the skill listed, it
		// is added to the character sheet as a zero level skill without
		// a specialty."
		name: "a skill table, skill not held", table: true,
		grant: grantOf("Electronics", 1, "Any"),
		want:  []Skill{{Name: "Electronics", Level: 0}},
	},
	{
		// The policy takes the first the book prints, held or not.
		name: "a skill table, skill held", table: true,
		held:  []Skill{{Name: "Electronics", Specialty: "Sensors", Level: 1}},
		grant: grantOf("Electronics", 1, "Any"),
		want: []Skill{
			{Name: "Electronics", Specialty: "Comms", Level: 1},
			{Name: "Electronics", Specialty: "Sensors", Level: 1},
		},
	},
	{
		name: "a skill table, the first specialty held", table: true,
		held:  []Skill{{Name: "Electronics", Specialty: "Comms", Level: 2}},
		grant: grantOf("Electronics", 1, "Any"),
		want:  []Skill{{Name: "Electronics", Specialty: "Comms", Level: 3}},
	},
	{
		name:    "an event, skill not held",
		grant:   grantOf("Advocate", 1, "Any"),
		want:    []Skill{{Name: "Advocate", Specialty: "Legal", Level: 1}},
		stamped: true,
	},
	{
		name:    "a degree at level 2",
		grant:   grantOf("Advocate", 2, "Any"),
		want:    []Skill{{Name: "Advocate", Specialty: "Legal", Level: 2}},
		stamped: true,
	},
	{
		name: "level 0, skill not held",
		grant: career.Effect{
			Kind: career.EffectSkill, Skill: "Drive", Specialties: []string{"Any"},
			AtLevelZero: true, Detail: "gain Drive at level 0",
		},
		want:    []Skill{{Name: "Drive", Level: 0}},
		stamped: true,
	},
	{
		name: "level 0, skill held",
		held: []Skill{{Name: "Drive", Specialty: "Wheeled", Level: 1}},
		grant: career.Effect{
			Kind: career.EffectSkill, Skill: "Drive", Specialties: []string{"Any"},
			AtLevelZero: true, Detail: "gain Drive at level 0",
		},
		want:    []Skill{{Name: "Drive", Specialty: "Wheeled", Level: 1}},
		stamped: true,
	},
}

// TestAnAnyResultIsResolved is p. 118 and its reading everywhere else
// (ERRATA E-45). "(Any)" is a choice the character makes, so the sheet
// never holds a specialty called "Any".
func TestAnAnyResultIsResolved(t *testing.T) {
	t.Parallel()

	for _, tc := range anyResults {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gen := engine(t, 300)

			for _, held := range tc.held {
				gen.char.State.GainSkill(held.Name, held.Specialty, held.Level)
			}

			if tc.table {
				defer gen.onTable()()
			}

			err := gen.applySkill(tc.grant, 0)
			if err != nil {
				t.Fatalf("applySkill: %v", err)
			}

			if !slices.Equal(gen.char.State.Skills, tc.want) {
				t.Errorf("skills = %v, want %v", gen.char.State.Skills, tc.want)
			}

			if got := slices.Contains(gen.char.Provenance.Deviations, "E-45"); got != tc.stamped {
				t.Errorf("E-45 stamped = %v, want %v", got, tc.stamped)
			}
		})
	}
}

// unheldTableResults are the cases TestASkillTableResultNotHeldGoesOnAtZero
// runs.
var unheldTableResults = []struct {
	name  string
	table bool
	held  []Skill
	grant career.Effect
	want  []Skill
}{
	{
		name: "a skill, not held", table: true, grant: persuadeGrant,
		want: []Skill{{Name: "Persuade", Level: 0}},
	},
	{
		name: "a skill, held at 0", table: true, grant: persuadeGrant,
		held: []Skill{{Name: "Persuade", Level: 0}},
		want: []Skill{{Name: "Persuade", Level: 1}},
	},
	{
		name: "a specific specialty, skill not held", table: true, grant: vaccSuitGrant,
		want: []Skill{{Name: "Suit", Level: 0}},
	},
	{
		// Holding another specialty is holding the skill, so the
		// level goes to the specialty the result names.
		name: "a specific specialty, another held", table: true, grant: vaccSuitGrant,
		held: []Skill{{Name: "Suit", Specialty: "Battle Dress", Level: 0}},
		want: []Skill{
			{Name: "Suit", Specialty: "Battle Dress", Level: 0},
			{Name: "Suit", Specialty: "Vacc Suit", Level: 1},
		},
	},
	{
		name: "an event, not held", grant: vaccSuitGrant,
		want: []Skill{{Name: "Suit", Specialty: "Vacc Suit", Level: 1}},
	},
}

var (
	persuadeGrant = career.Effect{
		Kind: career.EffectSkill, Skill: "Persuade", Level: 1, Detail: "gain Persuade",
	}
	vaccSuitGrant = grantOf("Suit", 1, "Vacc Suit")
)

// TestASkillTableResultNotHeldGoesOnAtZero is p. 118's other two forms:
// "If the character does not already have the skill, they should add the
// skill at level 0", and for a specific specialty, "it is added to the
// character sheet as a zero level skill without a specialty". Away from a
// skill table the level is granted as printed (ERRATA E-45).
func TestASkillTableResultNotHeldGoesOnAtZero(t *testing.T) {
	t.Parallel()

	for _, tc := range unheldTableResults {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gen := engine(t, 301)

			for _, held := range tc.held {
				gen.char.State.GainSkill(held.Name, held.Specialty, held.Level)
			}

			if tc.table {
				defer gen.onTable()()
			}

			err := gen.applySkill(tc.grant, 0)
			if err != nil {
				t.Fatalf("applySkill: %v", err)
			}

			if !slices.Equal(gen.char.State.Skills, tc.want) {
				t.Errorf("skills = %v, want %v", gen.char.State.Skills, tc.want)
			}
		})
	}
}

// TestAnAnyRankBenefitTakesTheSkillsOwnList is p. 116 against "(Any)"
// itself rather than a list the test supplies: the second grant takes a
// specialty the first did not, where recording "Any" made it a no-op.
func TestAnAnyRankBenefitTakesTheSkillsOwnList(t *testing.T) {
	t.Parallel()

	gen := engine(t, 320)

	granted := grantOf("Advocate", 1, "Any")

	for range 4 {
		err := gen.applyRankBenefits([]career.Effect{granted}, 0)
		if err != nil {
			t.Fatalf("applyRankBenefits: %v", err)
		}
	}

	want := []Skill{
		{Name: "Advocate", Specialty: "Legal", Level: 1},
		{Name: "Advocate", Specialty: "Oratory", Level: 1},
		{Name: "Advocate", Specialty: "Politics", Level: 1},
	}

	if !slices.Equal(gen.char.State.Skills, want) {
		t.Errorf("skills = %v, want %v: the fourth grant should have found all three held",
			gen.char.State.Skills, want)
	}
}

// TestAnAnyLanguageIsOneTheSettingSpeaks. pp. 304-314 print no Language
// specialties; they are "any language which can be chosen as a primary
// language", and not the character's own (p. 41).
func TestAnAnyLanguageIsOneTheSettingSpeaks(t *testing.T) {
	t.Parallel()

	gen := engine(t, 321)

	gen.setting = &setting.Data{Subsectors: []setting.Subsector{{Worlds: []setting.World{
		{PrimaryLanguages: []string{"Rivertongue", "Old Signal"}},
		{PrimaryLanguages: []string{"Old Signal", "Marsh Creole"}},
	}}}}
	gen.primaryLanguage = "Rivertongue"

	if got, want := gen.specialtiesOf("Language"), []string{"Old Signal", "Marsh Creole"}; !slices.Equal(got, want) {
		t.Errorf("specialtiesOf(Language) = %v, want %v", got, want)
	}

	gen.setting = &setting.Data{Subsectors: []setting.Subsector{{Worlds: []setting.World{
		{PrimaryLanguages: []string{"Rivertongue"}},
	}}}}

	err := gen.applySkill(grantOf("Language", 1, "Any"), 0)
	if err != nil {
		t.Fatalf("applySkill: %v", err)
	}

	if len(gen.char.State.Skills) != 0 {
		t.Errorf("a setting with no other language granted %v", gen.char.State.Skills)
	}
}

// TestAFirstTermAnyServiceSkillIsBare is p. 118's last sentence at p. 117's
// first term: the level-0 service skills are "added to the character sheet
// as a zero level skill without a specialty".
func TestAFirstTermAnyServiceSkillIsBare(t *testing.T) {
	t.Parallel()

	if got := specialtyOf(grantOf("Survival", 1, "Any")); got != "" {
		t.Errorf("specialtyOf(Survival (Any)) = %q, want none", got)
	}

	if got := specialtyOf(grantOf("Suit", 1, "Vacc Suit")); got != "Vacc Suit" {
		t.Errorf("specialtyOf(Suit (Vacc Suit)) = %q, want Vacc Suit", got)
	}
}

// grantOf is a skill result naming a specialty, at a level.
func grantOf(name string, level int, specialty string) career.Effect {
	return career.Effect{
		Kind: career.EffectSkill, Skill: name, Specialties: []string{specialty},
		Level: level, Detail: "gain " + name,
	}
}

// TestAnAnyLanguageWithOneLanguageOrNone: a setting offering one language
// besides the character's own takes it without a question, and a rank
// benefit with none to offer says so rather than granting a Language that
// names no language. With no setting at all there is nothing to offer.
func TestAnAnyLanguageWithOneLanguageOrNone(t *testing.T) {
	t.Parallel()

	gen := engine(t, 322)

	if got := gen.specialtiesOf("Language"); len(got) != 0 {
		t.Errorf("with no setting, specialtiesOf(Language) = %v", got)
	}

	gen.setting = &setting.Data{Subsectors: []setting.Subsector{{Worlds: []setting.World{
		{PrimaryLanguages: []string{"Rivertongue", "Old Signal"}},
	}}}}
	gen.primaryLanguage = "Rivertongue"

	err := gen.applySkill(grantOf("Language", 2, "Any"), 0)
	if err != nil {
		t.Fatalf("applySkill: %v", err)
	}

	if _, found := gen.char.State.Skill("Language", "Old Signal"); !found {
		t.Errorf("the one language on offer was not taken: %v", gen.char.State.Skills)
	}

	gen.primaryLanguage = "Old Signal"
	gen.setting.Subsectors[0].Worlds[0].PrimaryLanguages = []string{"Old Signal"}
	gen.char.State.Skills = nil

	err = gen.applyRankBenefits([]career.Effect{grantOf("Language", 1, "Any")}, 0)
	if err != nil {
		t.Fatalf("applyRankBenefits: %v", err)
	}

	if len(gen.char.State.Skills) != 0 || !strings.Contains(lastDetail(t, gen), "no language") {
		t.Errorf("granted %v and said %q", gen.char.State.Skills, lastDetail(t, gen))
	}
}

// TestARefusedAnySpecialtyEndsGeneration. "(Any)" is a choice like any
// other, so a player who declines it ends the session.
func TestARefusedAnySpecialtyEndsGeneration(t *testing.T) {
	t.Parallel()

	gen := engine(t, 323)

	gen.decider = refusingDecider{}

	err := gen.applySkill(grantOf("Advocate", 1, "Any"), 0)
	if !errors.Is(err, ErrNoSetting) {
		t.Errorf("applySkill = %v, want the refusal", err)
	}

	// A degree's field is the same question one step later: the field is
	// chosen, and then its specialty.
	gen.decider = refusingAt{point: pointSkillSpecialty}

	_, err = gen.takeFromList(career.Institution{
		Skills: []career.Effect{grantOf("Advocate", 1, "Any")},
	}, 0, 2, "Choose the field of the degree")
	if !errors.Is(err, ErrNoSetting) {
		t.Errorf("takeFromList = %v, want the refusal", err)
	}
}

// refusingAt answers like Policy and refuses one choice point.
type refusingAt struct{ point string }

func (r refusingAt) Choose(c Choice) (int, error) {
	if c.Point == r.point {
		return 0, ErrNoSetting
	}

	return Policy{}.Choose(c)
}

func (refusingAt) Kind() DeciderKind { return DeciderPolicy }

// TestJackOfAllTradesStopsAtThree is p. 309: "a character may never have a
// skill level in Jack of All Trades higher than 3." A grant past it is held
// at 3 and the record says what was lost (#176, ERRATA E-48), and no other
// skill has a ceiling.
func TestJackOfAllTradesStopsAtThree(t *testing.T) {
	t.Parallel()

	gen := engine(t, 62)

	gen.gainSkill("Jack of All Trades", "", 2, 0)
	gen.gainSkill("Jack of All Trades", "", 2, 0)

	if got := gen.char.State.SkillLevel("Jack of All Trades"); got != 3 {
		t.Errorf("Jack of All Trades = %d, want 3", got)
	}

	if got := gen.char.State.GainSkill("Jack of All Trades", "", 0).Level; got != 3 {
		t.Errorf("a level-0 grant moved Jack of All Trades to %d", got)
	}

	held := 0

	for _, event := range gen.log.Events() {
		if event.Consequence != nil && event.Consequence.Cite == ceilingCite {
			held++

			if !strings.Contains(event.Consequence.Detail, "1 level(s) lost") {
				t.Errorf("detail = %q, want the one level lost", event.Consequence.Detail)
			}
		}
	}

	if held != 1 {
		t.Errorf("%d ceiling consequences logged, want 1", held)
	}

	if !slices.Contains(gen.char.Provenance.Deviations, "E-48") {
		t.Error("a held grant did not stamp E-48")
	}
}

// TestOnlyJackOfAllTradesHasACeiling: a first grant past p. 309's 3, as a
// rank row could make, is held too, and no other skill stops anywhere.
func TestOnlyJackOfAllTradesHasACeiling(t *testing.T) {
	t.Parallel()

	var state State

	if got := state.GainSkill("Jack of All Trades", "", 4).Level; got != 3 {
		t.Errorf("a first grant of 4 = %d, want 3", got)
	}

	if got := state.GainSkill("Advocate", "Legal", 5).Level; got != 5 {
		t.Errorf("Advocate (Legal) = %d, want 5", got)
	}
}
