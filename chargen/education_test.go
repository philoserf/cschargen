package chargen_test

import (
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

		// Admission and honours are helped by EDU, success by INT
		// (pp. 87, 92).
		if got.Admission.Characteristic != "EDU" || got.Honors.Characteristic != "EDU" {
			t.Errorf("%s: admission and honours should read EDU", got.Name)
		}

		if got.Success.Characteristic != "INT" {
			t.Errorf("%s: success should read INT", got.Name)
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
		if len(institution.Failure) != 11 {
			t.Errorf("%s failure table has %d rows; a 2d6 table has 11",
				institution.Name, len(institution.Failure))
		}

		if len(institution.Events) != 19 {
			t.Errorf("%s events table has %d rows; a 2d10 table has 19",
				institution.Name, len(institution.Events))
		}

		if len(institution.LifeEvents) != 6 {
			t.Errorf("%s life events table has %d rows; a d6 table has 6",
				institution.Name, len(institution.LifeEvents))
		}

		if len(institution.Skills) == 0 {
			t.Errorf("%s has no skills to take a degree in", institution.Name)
		}

		for i, row := range institution.Failure {
			if row.Summary == "" {
				t.Errorf("%s failure result %d has no summary", institution.Name, i+2)
			}
		}

		for i, row := range institution.Events {
			if row.Summary == "" {
				t.Errorf("%s event %d has no summary", institution.Name, i+2)
			}
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

		record := generate(t, opts).State.Education
		if record == nil {
			continue
		}

		attempted++

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
	if character.State.Education != nil {
		t.Error("Step 8 ran although it was skipped")
	}

	// And the earlier steps still ran.
	if character.State.Family == nil {
		t.Error("skipping Step 8 skipped Step 5 as well")
	}
}
