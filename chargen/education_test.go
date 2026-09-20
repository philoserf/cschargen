package chargen_test

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/chargen"
)

// TestTheInstitutionsMatchTheirBoxes is a second reading of the three-row
// box each institution prints, and of the prerequisite beneath it.
// institutionBox is the three-row box an institution prints, plus the
// prerequisite beneath it.
type institutionBox struct {
	admission, success, honors int
	prerequisites              []career.Check

	// graduate marks the two tracks that read the other characteristic on
	// each throw and require a degree.
	graduate bool
}

func institutionBoxes() map[string]institutionBox {
	return map[string]institutionBox{
		// "Admission EDU 7+ / Success INT 8+ / Honors EDU 10+" (p. 86),
		// and "The character's EDU must be 6 or higher" (p. 86).
		"Undergraduate College": {
			admission: 7, success: 8, honors: 10,
			prerequisites: []career.Check{{Characteristic: "EDU", Number: 6}},
		},
		// "Admission EDU 9+ / Success INT 8+ / Honors EDU 10+" (p. 92),
		// and "EDU must be 6 or higher and their END must be 8 or higher"
		// (p. 92).
		"Military Academy": {
			admission: 9, success: 8, honors: 10,
			prerequisites: []career.Check{
				{Characteristic: "EDU", Number: 6},
				{Characteristic: "END", Number: 8},
			},
		},
		// "Admission INT 8+ / Success EDU 8+ / Honors INT 10+" (p. 97),
		// and "EDU must be 10 or higher and the character must have
		// achieved Success in an Undergraduate University or Military
		// Academy" (p. 97).
		"Graduate School": {
			admission: 8, success: 8, honors: 10, graduate: true,
			prerequisites: []career.Check{{Characteristic: "EDU", Number: 10}},
		},
		// "Admission INT 8+ / Success EDU 9+ / Honors INT 10+" (p. 100) --
		// and then p. 101 says "The player must roll 8 or higher on 2d6
		// for the character to be admitted ... The player must roll 8 or
		// higher on 2d6 for the character to succeed", which is ERRATA
		// E-28. The prose wins.
		"Medical School": {
			admission: 8, success: 8, honors: 10, graduate: true,
			prerequisites: []career.Check{{Characteristic: "EDU", Number: 8}},
		},
	}
}

func TestTheInstitutionsMatchTheirBoxes(t *testing.T) {
	t.Parallel()

	want := institutionBoxes()

	institutions := career.Institutions()
	if len(institutions) != len(want) {
		t.Fatalf("%d institutions built, %d read", len(institutions), len(want))
	}

	for _, got := range institutions {
		expected, known := want[got.Name]
		if !known {
			t.Errorf("no second reading for %s", got.Name)

			continue
		}

		checkInstitutionBox(t, got, expected)
	}
}

func checkInstitutionBox(t *testing.T, got career.Institution, expected institutionBox) {
	t.Helper()

	{
		if got.Admission.Number != expected.admission {
			t.Errorf("%s admission = %d+, want %d+", got.Name, got.Admission.Number, expected.admission)
		}

		if got.Success.Number != expected.success {
			t.Errorf("%s success = %d+, want %d+", got.Name, got.Success.Number, expected.success)
		}

		if got.Honors.Number != expected.honors {
			t.Errorf("%s honors = %d+, want %d+", got.Name, got.Honors.Number, expected.honors)
		}

		// The undergraduate tracks are helped by EDU on admission and
		// honours and by INT on success (pp. 87, 92). The graduate tracks
		// read the other one on each (pp. 97, 101).
		helps, succeeds := "EDU", "INT"
		if expected.graduate {
			helps, succeeds = "INT", "EDU"
		}

		if got.Admission.Characteristic != helps || got.Honors.Characteristic != helps {
			t.Errorf("%s: admission and honours should read %s", got.Name, helps)
		}

		if got.Success.Characteristic != succeeds {
			t.Errorf("%s: success should read %s", got.Name, succeeds)
		}

		if got.RequiresDegree != expected.graduate {
			t.Errorf("%s: RequiresDegree = %v, want %v",
				got.Name, got.RequiresDegree, expected.graduate)
		}

		if len(got.Prerequisites) != len(expected.prerequisites) {
			t.Errorf("%s has %d prerequisites, want %d",
				got.Name, len(got.Prerequisites), len(expected.prerequisites))
		}
	}
}

// TestEveryInstitutionHasItsPrintedTables: a 2d6 failure table has eleven
// rows, a 2d10 events table nineteen, and a d6 life-events table six.
func TestEveryInstitutionHasItsPrintedTables(t *testing.T) {
	t.Parallel()

	for _, institution := range career.Institutions() {
		checkInstitutionTables(t, institution)
	}
}

func checkInstitutionTables(t *testing.T, institution career.Institution) {
	t.Helper()

	{
		if len(institution.Failure) != 11 {
			t.Errorf("%s failure table has %d rows; a 2d6 table has 11",
				institution.Name, len(institution.Failure))
		}

		if len(institution.Events) != 19 {
			t.Errorf("%s events table has %d rows; a 2d10 table has 19",
				institution.Name, len(institution.Events))
		}

		// The graduate tracks reach the career Life Events table of
		// p. 120 rather than one of their own, so they print none.
		if institution.RequiresDegree {
			if len(institution.LifeEvents) != 0 {
				t.Errorf("%s prints a life events table; pp. 99, 103 send it to p. 120",
					institution.Name)
			}
		} else if len(institution.LifeEvents) != 6 {
			t.Errorf("%s life events table has %d rows; a d6 table has 6",
				institution.Name, len(institution.LifeEvents))
		}

		if len(institution.Skills) == 0 {
			t.Errorf("%s has no skills to take a degree in", institution.Name)
		}

		summaries := make([]string, 0, len(institution.Failure)+len(institution.Events))
		for _, row := range institution.Failure {
			summaries = append(summaries, row.Summary)
		}

		for _, row := range institution.Events {
			summaries = append(summaries, row.Summary)
		}

		checkSummaries(t, institution.Name, summaries)
	}
}

// checkSummaries: a row with no summary is a row the transcript cannot
// print, which is how a table with a gap in it shows up.
func checkSummaries(t *testing.T, where string, summaries []string) {
	t.Helper()

	for i, summary := range summaries {
		if summary == "" {
			t.Errorf("%s: a table row at index %d has no summary", where, i)
		}
	}
}

// TestTheAcademyLifeEventsDifferFromTheCollegiateOnes. p. 96's table is
// p. 91's with one result changed: where the college prints a death, the
// academy prints a battle the cadets were called into.
func TestTheAcademyLifeEventsDifferFromTheCollegiateOnes(t *testing.T) {
	t.Parallel()

	academy := career.MilitaryAcademy().LifeEvents
	college := career.Undergraduate().LifeEvents

	if academy[1].Summary == college[1].Summary {
		t.Error("the academy's result 2 is the college's; p. 96 prints a battle")
	}

	for i := range academy {
		if i == 1 {
			continue
		}

		if academy[i].Summary != college[i].Summary {
			t.Errorf("result %d differs between the two tables: %q and %q",
				i+1, academy[i].Summary, college[i].Summary)
		}
	}
}

// TestHigherEducationHappens. Every outcome of Step 8 -- declined, not
// admitted, washed out, graduated, graduated with honours -- should be
// reachable across the sample.
func TestHigherEducationHappens(t *testing.T) {
	t.Parallel()

	var attempted, admitted, graduated, honors int

	for seed := range uint64(200) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = -1

		history := generate(t, opts).State.Education
		if len(history) == 0 {
			continue
		}

		attempted++

		record := history[0]

		if record.Admitted {
			admitted++
		}

		if record.Succeeded {
			graduated++
		}

		if record.Honors {
			honors++
		}

		checkEducationRecord(t, seed, record)
	}

	for _, count := range []struct {
		name string
		got  int
	}{
		{"attempted", attempted},
		{"admitted", admitted},
		{"graduated", graduated},
		{"took honours", honors},
	} {
		if count.got == 0 {
			t.Errorf("in 200 seeds nobody %s higher education", count.name)
		}
	}
}

// checkEducationRecord holds the three throws in the order the page makes
// them: nobody takes honours without a degree, and nobody graduates from an
// institution that did not admit them.
func checkEducationRecord(t *testing.T, seed uint64, record *chargen.Education) {
	t.Helper()

	if record.Succeeded && record.Field == "" {
		t.Errorf("seed %d graduated in nothing", seed)
	}

	if record.Honors && !record.Succeeded {
		t.Errorf("seed %d took honours without a degree", seed)
	}

	if record.Succeeded && !record.Admitted {
		t.Errorf("seed %d graduated without being admitted", seed)
	}
}

// TestSkippingHigherEducation. p. 85: "characters are not required to
// attend college."
func TestSkippingHigherEducation(t *testing.T) {
	t.Parallel()

	opts := options(t, 5)

	opts.Inputs.TermLimit = -1
	opts.Inputs.SkipEducation = true

	character := generate(t, opts)
	if len(character.State.Education) != 0 {
		t.Error("Step 8 ran although it was skipped")
	}

	// And the earlier steps still ran.
	if character.State.Family == nil {
		t.Error("skipping Step 8 skipped Step 5 as well")
	}
}

// TestADegreeHelpsEnlistment. Four careers modify enlistment on a degree --
// Instructor, Journalist, Medic and Scientist -- and the engine recorded
// all of them as unimplemented until Step 8 landed. The record has to show
// the modifier applied, and nowhere show it recorded.
func TestADegreeHelpsEnlistment(t *testing.T) {
	t.Parallel()

	applied := 0

	for seed := range uint64(120) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = 3
		opts.Inputs.Career = "Instructor"

		character := generate(t, opts)

		applied += countDegreeModifiers(t, seed, character)
	}

	if applied == 0 {
		t.Error("in 120 seeds no degree ever helped an enlistment throw")
	}
}

// countDegreeModifiers walks a record for the degree bonuses, and fails on
// any that are still recorded as unimplemented.
func countDegreeModifiers(t *testing.T, seed uint64, character *chargen.Character) int {
	t.Helper()

	degrees := map[string]bool{
		"a degree": true, "a master's": true, "a doctorate": true, "medical school": true,
	}

	applied := 0

	for _, event := range character.Events {
		if event.Consequence != nil &&
			strings.Contains(event.Consequence.Detail, "modifies enlistment on a degree") {
			t.Errorf("seed %d still records the degree modifier as unimplemented", seed)
		}

		if event.Kind != chargen.EventThrow || event.Throw == nil {
			continue
		}

		for _, mod := range event.Throw.Mods {
			if degrees[mod.Name] {
				applied++
			}
		}
	}

	return applied
}

// TestAMedicalDegreeIsFourDifferentSpecialties is p. 101: "Gain Medic (Any)
// at level 3, Medic (Any) at level 2, Medic (Any) at level 2, and Medic
// (Any) at level 1. Specialties should all be different."
func TestAMedicalDegreeIsFourDifferentSpecialties(t *testing.T) {
	t.Parallel()

	checked := 0

	for seed := range uint64(120) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = -1

		character := generate(t, opts)
		if !holdsA(character, "medical doctor") {
			continue
		}

		checked++

		levels := map[int]int{}

		for _, held := range character.State.Skills {
			if held.Name != "Medic" {
				continue
			}

			levels[held.Level]++
		}

		// Four specialties at 3, 2, 2 and 1. A character who held Medic
		// before medical school may have more.
		for level, want := range map[int]int{3: 1, 2: 2, 1: 1} {
			if levels[level] < want {
				t.Errorf("seed %d holds %d Medic specialties at level %d, want %d",
					seed, levels[level], level, want)
			}
		}
	}

	if checked == 0 {
		t.Fatal("in 120 seeds nobody finished medical school")
	}
}

// TestAGraduateDegreeRaisesTheUndergraduateField is p. 97: "The character
// may now increase the skill they increased in Undergraduate University by
// two levels."
//
// The raise has to name the same specialty the bachelor's took. "Advocate"
// and "Advocate (Any)" are two skills, and raising the wrong one leaves a
// character holding both.
func TestAGraduateDegreeRaisesTheUndergraduateField(t *testing.T) {
	t.Parallel()

	checked := 0

	for seed := range uint64(120) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = -1

		character := generate(t, opts)
		if !holdsA(character, "master's") {
			continue
		}

		checked++

		field := ""

		for _, held := range character.State.Education {
			if held.Degree == "bachelor's" {
				field = held.Field
			}
		}

		if field == "" {
			t.Errorf("seed %d holds a master's and no bachelor's", seed)

			continue
		}

		checkTheFieldRose(t, seed, character, field)
	}

	if checked == 0 {
		t.Fatal("in 120 seeds nobody finished graduate school")
	}
}

// checkTheFieldRose: the bachelor's took the field to 2 and the master's
// adds two, so it is at 4 or higher -- and it is one skill, not two under
// different specialties.
func checkTheFieldRose(t *testing.T, seed uint64, character *chargen.Character, field string) {
	t.Helper()

	names := 0

	for _, held := range character.State.Skills {
		if held.Name != field {
			continue
		}

		names++

		if held.Level < 4 {
			t.Errorf("seed %d: %s is at %d after a master's, want 4 or more",
				seed, held.Full(), held.Level)
		}
	}

	if names != 1 {
		t.Errorf("seed %d holds %s under %d different specialties", seed, field, names)
	}
}

// holdsA reports whether a character finished an institution with a degree
// of this name.
func holdsA(character *chargen.Character, degree career.Degree) bool {
	for _, held := range character.State.Education {
		if held.Degree == degree {
			return true
		}
	}

	return false
}

// TestAPlayerMayReturnToHigherEducation is one of the five things p. 125
// lets a character decide between terms: "continue in this career, change
// to a different career, change to a different assignment within the same
// career, return to higher education, or exit character generation."
//
// It was deferred from milestone 5 with the reason that it is a choice only
// a player can make well, and that is why it is offered to a player and not
// to the policy.
func TestAPlayerMayReturnToHigherEducation(t *testing.T) {
	t.Parallel()

	// A decider that declines education before Step 9 and takes it after,
	// so the return is the only way this character goes to school.
	returner := &returnToSchool{}

	opts := options(t, 5)

	opts.Decider = returner
	opts.Inputs.TermLimit = 4
	opts.Inputs.SkipEducation = true

	character := generate(t, opts)

	if len(character.State.Education) == 0 {
		t.Fatal("a character who returned to education has no education")
	}

	// And they had a career before it, which is what makes it a return.
	if len(character.State.Services) == 0 {
		t.Error("they went back to school without ever leaving")
	}
}

// TestThePolicyIsNotAskedToReturn. The policy takes the first option every
// time, and being asked this one every term would send every character back
// to school forever.
func TestThePolicyIsNotAskedToReturn(t *testing.T) {
	t.Parallel()

	opts := options(t, 5)

	opts.Inputs.TermLimit = 6

	for _, event := range generate(t, opts).Events {
		if event.Kind == chargen.EventChoice &&
			event.Choice.Point == "return_to_education" {
			t.Error("the auto policy was asked whether to return to education")
		}
	}
}

// returnToSchool answers every choice with the first option and says yes to
// the one question this test is about.
type returnToSchool struct{}

func (returnToSchool) Choose(ask chargen.Choice) (int, error) {
	if ask.Point == "return_to_education" {
		return 1, nil
	}

	return 0, nil
}

func (returnToSchool) Kind() chargen.DeciderKind { return chargen.DeciderPlayer }

func (returnToSchool) Ask(chargen.Question) (string, error) { return "", nil }

// TestARefusedReturnEndsGeneration. The Step 18 question goes through the
// Decider like every choice, and a refusal ends the run rather than
// defaulting to staying in the career.
func TestARefusedReturnEndsGeneration(t *testing.T) {
	t.Parallel()

	opts := options(t, 5)

	opts.Decider = refuseReturn{}
	opts.Inputs.TermLimit = 4

	// Skipping Step 8 leaves the institutions open, so Step 18 has
	// something to offer.
	opts.Inputs.SkipEducation = true

	_, err := chargen.New(opts).Run()
	if err == nil {
		t.Fatal("a refused question did not end generation")
	}
}

// refuseReturn answers everything but the Step 18 question, which it
// declines.
type refuseReturn struct{}

func (refuseReturn) Choose(ask chargen.Choice) (int, error) {
	if ask.Point == "return_to_education" {
		return 0, chargen.ErrPlayerGone
	}

	return 0, nil
}

func (refuseReturn) Kind() chargen.DeciderKind { return chargen.DeciderPlayer }

func (refuseReturn) Ask(chargen.Question) (string, error) { return "", nil }
