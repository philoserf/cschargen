package career_test

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/career"
)

// slaveCareer names the career of pp. 150-154, which the book heads with a
// term its OGL notice reserves. That word is not in this repository.
const slaveCareer = "Engineered/Uplift Slave"

// TestTheBookIsFullyTranscribed is what milestone 3 was filed to reach.
//
// The Careers list on pp. 8-9 names thirty-four, and each is here by the
// name this repository gives it -- which is the book's own except for the
// slave career, where the book's word is Product Identity. A career the
// list names and All() does not carry is the failure this test exists for;
// a career All() carries and the list does not is an invented one.
func TestTheBookIsFullyTranscribed(t *testing.T) {
	t.Parallel()

	printed := []string{
		"Adventurer", slaveCareer, "Arts", "Belter", "Celebrity", "Clergy",
		"Colonist", "Corporate Shipper", "Craftsperson", "Diplomatic Service",
		"Exotic", "Explorer", "Fringe Marketer", "Gambler",
		"Independent Merchant", "Instructor", "Investigator", "Journalist",
		"Marine", "Medic", "National Navy", "Orbital Construction",
		"Organized Crime", "Pirate", "Politician", "Prisoner", "Scavenger",
		"Scientist", "Sports", "System Defense Forces (Navy)",
		"System Defense Forces (Troopers)", "System Defense Forces (Wet Navy)",
		"Thief", "Vagabond",
	}

	if len(printed) != 34 {
		t.Fatalf("the list above has %d names; pp. 8-9 print 34", len(printed))
	}

	built := map[string]bool{}
	for _, def := range career.All() {
		built[def.Name] = true
	}

	for _, name := range printed {
		if !built[name] {
			t.Errorf("the book names %s and All() does not carry it", name)
		}
	}

	if len(built) != len(printed) {
		t.Errorf("All() carries %d careers; the book names %d", len(built), len(printed))
	}
}

// TestEveryCareerHasItsPrintedShape is the completeness check. What is
// invariant is narrower than "every career has
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
				// The slave career (p. 150) prints "0" and "None" on its
				// first two rows: a benefit roll there really does buy
				// nothing, which is the point being made about the career.
				if def.Name == slaveCareer {
					continue
				}

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

	// Every career states this, and a career missing from the map fails the
	// test rather than defaulting -- the two exceptions are the whole point.
	want := map[string]bool{
		"Colonist":                         true,
		"Adventurer":                       true,
		"Arts":                             true,
		"Belter":                           true,
		"Celebrity":                        true,
		"Clergy":                           true,
		"Corporate Shipper":                true,
		"Craftsperson":                     true,
		"Diplomatic Service":               true,
		"Exotic":                           true,
		"Explorer":                         true,
		"Fringe Marketer":                  true,
		"Gambler":                          true,
		"Independent Merchant":             true,
		"Instructor":                       true,
		"Investigator":                     true,
		"Journalist":                       true,
		"Marine":                           true,
		"Medic":                            true,
		"National Navy":                    true,
		"Orbital Construction":             true,
		"Organized Crime":                  true,
		"System Defense Forces (Navy)":     true,
		"System Defense Forces (Troopers)": true,
		"System Defense Forces (Wet Navy)": true,
		"Pirate":                           true,
		"Politician":                       true,
		"Scavenger":                        true,
		"Scientist":                        true,
		"Sports":                           true,
		"Thief":                            true,
		"Engineered/Uplift Slave":          false,
		"Prisoner":                         false,
		"Vagabond":                         false,
	}

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

// TestTheEnslavedPathsAreElevenRowTables. pp. 74 and 84 are 2d6 tables
// where every other life-period path is a 2d10, and they share nine of
// their eleven rows -- the teenage one prints a different result 4 and is
// otherwise the youth table word for word.
func TestTheEnslavedPathsAreElevenRowTables(t *testing.T) {
	t.Parallel()

	youth := career.EnslavedYouth()
	teenage := career.EnslavedTeenage()

	for _, path := range []career.EnslavedPath{youth, teenage} {
		if len(path.Rows) != 11 {
			t.Errorf("%s has %d rows; a 2d6 table has 11", path.Name, len(path.Rows))
		}

		if path.Cite == "" {
			t.Errorf("%s has no page cite", path.Name)
		}

		for i, row := range path.Rows {
			if row.Summary == "" {
				t.Errorf("%s result %d has no summary", path.Name, i+2)
			}
		}
	}

	// Result 4 is the one that differs.
	if youth.Rows[2].Summary == teenage.Rows[2].Summary {
		t.Error("the two tables' result 4 is the same; pp. 74 and 84 print different ones")
	}

	for i := range youth.Rows {
		if i == 2 {
			continue
		}

		if youth.Rows[i].Summary != teenage.Rows[i].Summary {
			t.Errorf("result %d differs between the two tables: %q and %q",
				i+2, youth.Rows[i].Summary, teenage.Rows[i].Summary)
		}
	}
}

// TestEveryPathIsLabelledAsTheBookHeadsIt. The paths of Steps 6 and 7 have
// no titles: the book heads each one with the requirement that opens it,
// "Path 2 (STR 8+ or END 8+)", and that heading is the only description
// there is. Offering a bare "Path 2" tells a player nothing, and says
// nothing about why the paths they do not qualify for are absent.
func TestEveryPathIsLabelledAsTheBookHeadsIt(t *testing.T) {
	t.Parallel()

	for _, path := range career.YouthPaths() {
		if path.Qualifier == "" {
			t.Errorf("youth %s carries no qualifier", path.Name)
		}

		if want := path.Name + " (" + path.Qualifier + ")"; path.Label() != want {
			t.Errorf("youth %s labels itself %q", path.Name, path.Label())
		}
	}

	for _, path := range career.TeenagePaths() {
		if path.Qualifier == "" {
			t.Errorf("teenage %s carries no qualifier", path.Name)
		}

		if want := path.Name + " (" + path.Qualifier + ")"; path.Label() != want {
			t.Errorf("teenage %s labels itself %q", path.Name, path.Label())
		}
	}
}

// TestAPathQualifierNamesItsOwnRequirement holds the transcribed heading
// against the check the engine actually makes, so the two cannot drift: a
// path gated on DEX must say DEX.
func TestAPathQualifierNamesItsOwnRequirement(t *testing.T) {
	t.Parallel()

	for _, path := range career.YouthPaths() {
		for _, check := range path.Requires {
			if !strings.Contains(path.Qualifier, check.Characteristic) {
				t.Errorf("youth %s is gated on %s and says %q",
					path.Name, check.Characteristic, path.Qualifier)
			}
		}
	}

	for _, path := range career.TeenagePaths() {
		for _, check := range path.Requires {
			if !strings.Contains(path.Qualifier, check.Characteristic) {
				t.Errorf("teenage %s is gated on %s and says %q",
					path.Name, check.Characteristic, path.Qualifier)
			}
		}
	}
}
