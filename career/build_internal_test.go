package career

import (
	"slices"
	"testing"
)

// The constructors in build.go are how every table in this package is
// written, so most of them are covered a thousand times over by the tables
// themselves. These are the branches no printed result happens to take.

func TestJoinOrRendersAListTheWayThePageDoes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		names []string
		want  string
	}{
		{nil, ""},
		{[]string{"Carouse"}, "Carouse"},
		{[]string{"Carouse", "Gambler"}, "Carouse or Gambler"},
		{[]string{"Broker", "Carouse", "Recon"}, "Broker, Carouse or Recon"},
	}

	for _, tc := range tests {
		if got := joinOr(tc.names); got != tc.want {
			t.Errorf("joinOr(%v) = %q, want %q", tc.names, got, tc.want)
		}
	}
}

// TestRatingNamesWhatItMoves. The three targets of p. 320 read differently
// on the page -- "an existing Ally or Contact", "everyone in your life",
// "all family members" -- and the detail is what a transcript prints, so
// each of them has to say which it was.
func TestRatingNamesWhatItMoves(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		effect Effect
		want   string
	}{
		{"one of a kind", rating(TargetOne, Contact, -20, ""), "lower one contact's Relationship Rating by -20"},
		{"any one", rating(TargetOne, "", 50, ""), "raise one relationship's Relationship Rating by 50"},
		{"everyone", rating(TargetAll, "", -25, ""), "lower every relationship's Relationship Rating by -25"},
		{"the family", rating(TargetFamily, "", 20, ""), "raise every family member's Relationship Rating by 20"},
		{"a rolled amount", rating(TargetOne, Ally, -1, "1d6x20"), "lower one ally's Relationship Rating by 1d6x20"},
	}

	for _, tc := range tests {
		if tc.effect.Detail != tc.want {
			t.Errorf("%s: %q, want %q", tc.name, tc.effect.Detail, tc.want)
		}
	}
}

// TestLoseTieNamesItsOrder. Life Event 4 tries four kinds in order and the
// same result reversed, so a transcript that did not print the order would
// not say which of the two happened.
func TestLoseTieNamesItsOrder(t *testing.T) {
	t.Parallel()

	forward := loseTie(Ally, Contact, Rival, Enemy)
	if forward.Detail != "lose ally, contact, rival or enemy, in that order" {
		t.Errorf("forward = %q", forward.Detail)
	}

	unordered := loseTie()
	if unordered.Detail != "lose a relationship" {
		t.Errorf("unordered = %q", unordered.Detail)
	}
}

// TestTheMilitaryEventsTableIsAnElevenRowTable, which nothing else checks:
// it is reached through EffectMilitaryEvent rather than by name, so the
// shape test over All() never sees it.
// TestEveryWeaponBenefitOffersItsAlternatives is p. 129: "If the player
// wishes, they may choose to take a level in Melee (Any) or Gun Combat
// (Any) in lieu of a weapon." Fifteen rows across the corpus print the
// Weapon benefit, and the alternative is part of the benefit rather than a
// house rule -- so no row may offer the weapon alone.
func TestEveryWeaponBenefitOffersItsAlternatives(t *testing.T) {
	t.Parallel()

	found := 0

	for _, def := range All() {
		for row, benefit := range def.Benefits {
			if !offersAWeapon(benefit.Other) {
				continue
			}

			found++

			labels := make([]string, 0, len(benefit.Other.Options))
			for _, option := range benefit.Other.Options {
				labels = append(labels, option.Label)
			}

			for _, want := range []string{"a weapon", "Melee (Any)", "Gun Combat (Any)"} {
				if !slices.Contains(labels, want) {
					t.Errorf("%s benefit row %d does not offer %q: %v",
						def.Name, row+1, want, labels)
				}
			}
		}
	}

	if found == 0 {
		t.Fatal("no career prints the Weapon benefit; the test checks nothing")
	}
}

// offersAWeapon reports whether an effect is the Weapon benefit's choice.
func offersAWeapon(effect Effect) bool {
	if effect.Kind != EffectChoice {
		return false
	}

	for _, option := range effect.Options {
		if option.Label == "a weapon" {
			return true
		}
	}

	return false
}

func TestTheMilitaryEventsTableIsAnElevenRowTable(t *testing.T) {
	t.Parallel()

	table := MilitaryEvents()
	if len(table) != 11 {
		t.Fatalf("MilitaryEvents has %d rows; a 2d6 table has 11 (p. 121)", len(table))
	}

	for i, row := range table {
		if row.Summary == "" {
			t.Errorf("military event %d (2d6 result %d) has no summary", i, i+2)
		}
	}
}

// TestEveryTranscribedSkillIsOnTheList walks every effect in the corpus
// and holds the skill it names against the list of pp. 304-314. A name
// that is not on it is either a typo in the transcription or a typo in the
// book, and the two are told apart by reading the page -- which is how
// ERRATA E-35 came to exist.
func TestEveryTranscribedSkillIsOnTheList(t *testing.T) {
	t.Parallel()

	named := map[string]bool{}
	collect(everyEffect(), named)

	if len(named) == 0 {
		t.Fatal("no effect in the corpus names a skill")
	}

	for name := range named {
		_, found := SkillByName(name)
		if !found {
			t.Errorf("%q is named by a table and is not on the skill list", name)
		}
	}
}

// TestEveryTranscribedSpecialtyIsOnItsSkill holds the specialties too. A
// specialty of "Any" is the book's own word for "choose one", and is not a
// specialty.
func TestEveryTranscribedSpecialtyIsOnItsSkill(t *testing.T) {
	t.Parallel()

	for _, effect := range everyEffect() {
		if effect.Kind != EffectSkill || effect.Skill == "" {
			continue
		}

		definition, found := SkillByName(effect.Skill)
		if !found {
			continue
		}

		if definition.Open {
			continue
		}

		for _, specialty := range effect.Specialties {
			if specialty == "Any" || slices.Contains(definition.Specialties, specialty) {
				continue
			}

			t.Errorf("%s (%s) is not a specialty the book prints for %s",
				effect.Skill, specialty, effect.Skill)
		}
	}
}

// collect gathers the skill names an effect tree mentions.
func collect(effects []Effect, into map[string]bool) {
	for _, effect := range effects {
		if effect.Kind == EffectSkill && effect.Skill != "" {
			into[effect.Skill] = true
		}

		collect(effect.Success, into)
		collect(effect.Failure, into)
		collect(effect.Group, into)
		collect(effect.Fallback, into)

		for _, row := range effect.Sub {
			collect(row.Effects, into)
		}

		for _, option := range effect.Options {
			collect(option.Effects, into)
		}
	}
}

// everyEffect is every effect in every transcribed table.
func everyEffect() []Effect {
	found := make([]Effect, 0, 4096)

	for _, def := range All() {
		for _, row := range def.Mishaps {
			found = append(found, row.Effects...)
		}

		for _, row := range def.Events {
			found = append(found, row.Effects...)
		}

		for _, benefit := range def.Benefits {
			found = append(found, benefit.Other)
		}

		for _, table := range def.Tables {
			found = append(found, table.Rows[:]...)
		}

		for _, assignment := range def.Assignments {
			found = append(found, assignment.Skills.Rows[:]...)

			for _, rank := range assignment.Ranks {
				found = append(found, rank...)
			}

			for _, rank := range assignment.OfficerRanks {
				found = append(found, rank...)
			}
		}
	}

	for _, row := range LifeEvents() {
		found = append(found, row.Effects...)
	}

	for _, row := range MilitaryEvents() {
		found = append(found, row.Effects...)
	}

	for _, rows := range [][]EventRow{
		YouthLifeEvents(), TeenageLifeEvents(),
		EnslavedYouth().Rows, EnslavedTeenage().Rows,
	} {
		for _, row := range rows {
			found = append(found, row.Effects...)
		}
	}

	for _, path := range YouthPaths() {
		for _, row := range path.Rows {
			found = append(found, row.Effects...)
		}
	}

	for _, path := range TeenagePaths() {
		for _, row := range path.Rows {
			found = append(found, row.Effects...)
		}
	}

	for _, institution := range Institutions() {
		found = append(found, institution.Skills...)

		for _, row := range institution.Failure {
			found = append(found, row.Effects...)
		}

		for _, rows := range [][]EventRow{institution.Events, institution.LifeEvents} {
			for _, row := range rows {
				found = append(found, row.Effects...)
			}
		}
	}

	return found
}

// TestSkillByNameFindsNothingForAName holds the miss: a name the list does
// not carry is reported as absent rather than returned as an empty entry
// the caller might use.
func TestSkillByNameFindsNothingForAName(t *testing.T) {
	t.Parallel()

	_, found := SkillByName("Sabotage")
	if found {
		t.Error("a skill the book does not print was found on the list")
	}
}

// TestEveryModifierNamesAThrowTheEngineTakes is the gate on a bug this
// test was written to find: twenty-five results granted a modifier to
// "next survival roll" and nothing in the engine consumed it, and about
// twenty more were filed under names -- "next two advancement rolls",
// "enlistment in any criminal career", "Melee checks in this career" --
// that no part of the engine reads.
//
// A modifier under such a name is granted, recorded, shown in the
// transcript, and never applied. Nothing else in the gate can tell: the
// record looks right, and the throw is simply easier than the page says.
func TestEveryModifierNamesAThrowTheEngineTakes(t *testing.T) {
	t.Parallel()

	known := Throws()
	seen := 0

	for _, effect := range everyEffect() {
		collectThrowNames(t, effect, known, &seen)
	}

	if seen == 0 {
		t.Fatal("no effect in the corpus names a throw")
	}
}

// collectThrowNames walks one effect tree, holding every throw name it
// finds against the list.
func collectThrowNames(t *testing.T, effect Effect, known []string, seen *int) {
	t.Helper()

	if effect.Applies != "" {
		*seen++

		if !slices.Contains(known, effect.Applies) {
			t.Errorf("%q is a throw no part of the engine takes: %s", effect.Applies, effect.Detail)
		}
	}

	for _, nested := range [][]Effect{effect.Success, effect.Failure, effect.Group, effect.Fallback} {
		for _, inner := range nested {
			collectThrowNames(t, inner, known, seen)
		}
	}

	for _, row := range effect.Sub {
		for _, inner := range row.Effects {
			collectThrowNames(t, inner, known, seen)
		}
	}

	for _, option := range effect.Options {
		for _, inner := range option.Effects {
			collectThrowNames(t, inner, known, seen)
		}
	}
}

// TestEverySubTableCoversTheDie holds the forty-two 1d6 tables printed
// inside results: every result of the die falls in exactly one row. A gap
// would leave the engine with nothing to apply, and an overlap would make
// the first row printed win silently.
func TestEverySubTableCoversTheDie(t *testing.T) {
	t.Parallel()

	seen := 0

	for _, effect := range everyEffect() {
		checkSubTables(t, effect, &seen)
	}

	if seen == 0 {
		t.Fatal("no result in the corpus prints a 1d6 table")
	}
}

// checkSubTables walks one effect tree, holding each sub-table it finds to
// covering 1 through 6 exactly once.
func checkSubTables(t *testing.T, effect Effect, seen *int) {
	t.Helper()

	if effect.Kind == EffectSubTable {
		*seen++

		covered := map[int]int{}

		for _, row := range effect.Sub {
			for result := row.From; result <= row.To; result++ {
				covered[result]++
			}
		}

		// A 2d6 table runs 2 to 12 before a characteristic's modifier,
		// which can carry it either way, so its rows run past both ends.
		low, high := 1, 6
		if effect.Dice == "2d6" {
			low, high = 2, 12
		}

		for result := low; result <= high; result++ {
			if covered[result] != 1 {
				t.Errorf("a table covers %d %d times: %s", result, covered[result], effect.Detail)
			}
		}
	}

	for _, nested := range [][]Effect{effect.Success, effect.Failure, effect.Group, effect.Fallback} {
		for _, inner := range nested {
			checkSubTables(t, inner, seen)
		}
	}

	for _, row := range effect.Sub {
		for _, inner := range row.Effects {
			checkSubTables(t, inner, seen)
		}
	}

	for _, option := range effect.Options {
		for _, inner := range option.Effects {
			checkSubTables(t, inner, seen)
		}
	}
}
