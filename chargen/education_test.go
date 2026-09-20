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
