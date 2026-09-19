package career_test

import (
	"testing"

	"github.com/philoserf/cschargen/career"
)

// TestEveryCareerHasItsPrintedShape is the completeness check of
// docs/MILESTONE-1.md. What is invariant is narrower than "every career has
// every table": a d66 event table has 36 entries, a 2d6 mishap table has
// 11, any skill table that exists has 6 rows, and a benefit table has 7.
// Which tables exist varies by career, and the test must not assume
// otherwise -- Vagabond has no Advanced Education table at all.
func TestEveryCareerHasItsPrintedShape(t *testing.T) {
	t.Parallel()

	for _, def := range career.All() {
		t.Run(def.Name, func(t *testing.T) {
			t.Parallel()

			if len(def.Events) != 36 {
				t.Errorf("%d event rows, want 36", len(def.Events))
			}

			for _, result := range career.D66Results() {
				row, ok := def.Events[result]
				if !ok {
					t.Errorf("no event row for a d66 of %d", result)

					continue
				}

				if row.Summary == "" {
					t.Errorf("event %d has no summary", result)
				}
			}

			for i, row := range def.Mishaps {
				if row.Summary == "" {
					t.Errorf("mishap %d (2d6 result %d) has no summary", i, i+2)
				}
			}

			for i, row := range def.Benefits {
				if row.Cash == 0 && row.Other.Detail == "" {
					t.Errorf("benefit row %d is empty", i+1)
				}
			}

			if len(def.Assignments) == 0 {
				t.Error("no assignments")
			}

			if def.Cite == "" {
				t.Error("no page cite")
			}
		})
	}
}

// TestSkillTablesHaveSixRows: a row is a d6 result, and a table missing one
// would silently award whatever the zero Effect means.
func TestSkillTablesHaveSixRows(t *testing.T) {
	t.Parallel()

	for _, def := range career.All() {
		t.Run(def.Name, func(t *testing.T) {
			t.Parallel()

			tables := def.Tables
			for _, assignment := range def.Assignments {
				tables = append(tables, assignment.Skills)
			}

			for _, table := range tables {
				for i, row := range table.Rows {
					if row.Detail == "" {
						t.Errorf("%s table, row %d (a d6 of %d) is empty", table.Name, i, i+1)
					}
				}
			}
		})
	}
}

// TestTheTwoCareersWithoutEnlistmentSayWhy: a nil Enlistment is a rule, not
// an omission, and a reader of the data should not have to infer which.
func TestTheTwoCareersWithoutEnlistmentSayWhy(t *testing.T) {
	t.Parallel()

	for _, def := range career.All() {
		t.Run(def.Name, func(t *testing.T) {
			t.Parallel()

			if def.Enlistment == nil && def.EnlistmentNote == "" {
				t.Error("no enlistment throw and no note saying why")
			}

			if def.Enlistment != nil && def.EnlistmentNote != "" {
				t.Error("an enlistment throw and a note explaining its absence")
			}
		})
	}
}

// TestMishapEjectsMatchesThePage: Vagabond (p. 298) and Prisoner (p. 263)
// each say in print that a mishap does not force the character out, and
// they are the two careers a character is most often forced into -- so the
// exception is not a corner and a regression here would be quiet.
func TestMishapEjectsMatchesThePage(t *testing.T) {
	t.Parallel()

	want := map[string]bool{"Colonist": true, "Vagabond": false, "Prisoner": false}

	for _, def := range career.All() {
		expected, known := want[def.Name]
		if !known {
			t.Errorf("no expectation recorded for %s", def.Name)

			continue
		}

		if def.MishapEjects != expected {
			t.Errorf("%s: MishapEjects = %v, want %v", def.Name, def.MishapEjects, expected)
		}
	}
}

// TestVagabondHasNoAdvancedEducationTable is the data format's first real
// test. p. 117 says a career has "at least four skill tables" and names
// Advanced Education among the usual ones; p. 296 prints two career-wide
// tables for Vagabond and no more.
func TestVagabondHasNoAdvancedEducationTable(t *testing.T) {
	t.Parallel()

	vagabond := career.Vagabond()

	_, found := vagabond.Table(career.AdvancedEducation)
	if found {
		t.Error("Vagabond has an Advanced Education table; p. 296 prints none")
	}

	for _, kind := range []career.SkillTableKind{career.PersonalDevelopment, career.ServiceSkills} {
		_, found = vagabond.Table(kind)
		if !found {
			t.Errorf("Vagabond has no %s table", kind)
		}
	}
}

func TestAdvancedEducationIsGatedOnEDU(t *testing.T) {
	t.Parallel()

	for _, def := range career.All() {
		table, found := def.Table(career.AdvancedEducation)
		if !found {
			continue
		}

		if table.MinimumEDU == 0 {
			t.Errorf("%s: the Advanced Education table is ungated (p. 117)", def.Name)
		}
	}
}

// TestEveryCareerReachesLifeEventsAt31To36: the book prints the span as one
// row, and a career missing it would silently have six dead results.
func TestEveryCareerReachesLifeEventsAt31To36(t *testing.T) {
	t.Parallel()

	for _, def := range career.All() {
		for result := 31; result <= 36; result++ {
			row := def.Events[result]
			found := false

			for _, effect := range row.Effects {
				if effect.Kind == career.EffectLifeEvent {
					found = true
				}
			}

			if !found {
				t.Errorf("%s: event %d does not reach the Life Events table", def.Name, result)
			}
		}
	}
}

func TestLookups(t *testing.T) {
	t.Parallel()

	colonist, ok := career.ByName("Colonist")
	if !ok {
		t.Fatal("Colonist not found")
	}

	settler, ok := colonist.Assignment("Settler")
	if !ok {
		t.Fatal("Settler not found")
	}

	if settler.Survival.Characteristic != "END" || settler.Survival.Number != 7 {
		t.Errorf("Settler survival = %s %d+, want END 7+ (p. 173)",
			settler.Survival.Characteristic, settler.Survival.Number)
	}

	_, ok = career.ByName("Celebrity")
	if ok {
		t.Error("Celebrity is not implemented in this milestone but was found")
	}

	_, ok = colonist.Assignment("Ambassador")
	if ok {
		t.Error("Colonist has no Ambassador assignment")
	}
}

func TestD66ResultsIsThirtySix(t *testing.T) {
	t.Parallel()

	results := career.D66Results()
	if len(results) != 36 {
		t.Fatalf("%d results, want 36", len(results))
	}

	if results[0] != 11 || results[35] != 66 {
		t.Errorf("results run %d to %d, want 11 to 66", results[0], results[35])
	}

	for _, result := range results {
		tens, units := result/10, result%10
		if tens < 1 || tens > 6 || units < 1 || units > 6 {
			t.Errorf("%d is not a d66 result", result)
		}
	}
}

func TestSharedTableShapes(t *testing.T) {
	t.Parallel()

	injury := career.Injury()
	for i, row := range injury {
		if row.Summary == "" {
			t.Errorf("Injury row %d (a d6 of %d) has no summary", i, i+1)
		}
	}

	// The sixth result is "lightly injured. No permanent effect" (p. 119),
	// so it is the one row that legitimately carries no effects.
	if len(injury[5].Effects) != 0 {
		t.Errorf("Injury 6 has %d effects; p. 119 gives it none", len(injury[5].Effects))
	}

	for i, row := range career.LifeEvents() {
		if row.Summary == "" {
			t.Errorf("Life Events row %d (a 2d6 of %d) has no summary", i, i+2)
		}

		if len(row.Effects) == 0 {
			t.Errorf("Life Events row %d (a 2d6 of %d) has no effects", i, i+2)
		}
	}
}

// TestEveryEffectCarriesDetail: an event's meaning is rarely recoverable
// from its mechanical parts, and the record prints the detail.
func TestEveryEffectCarriesDetail(t *testing.T) {
	t.Parallel()

	var walk func(t *testing.T, where string, effects []career.Effect)

	walk = func(t *testing.T, where string, effects []career.Effect) {
		t.Helper()

		for i, effect := range effects {
			if effect.Detail == "" {
				t.Errorf("%s effect %d has no detail", where, i)
			}

			walk(t, where, effect.Success)
			walk(t, where, effect.Failure)

			for _, option := range effect.Options {
				walk(t, where, option.Effects)
			}
		}
	}

	for _, def := range career.All() {
		for result, row := range def.Events {
			walk(t, def.Name+" event "+itoa(result), row.Effects)
		}

		for i, row := range def.Mishaps {
			walk(t, def.Name+" mishap "+itoa(i+2), row.Effects)
		}
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	var out []byte

	for n > 0 {
		out = append([]byte{byte('0' + n%10)}, out...)

		n /= 10
	}

	return string(out)
}
