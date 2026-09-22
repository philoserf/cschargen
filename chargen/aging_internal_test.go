package chargen

import (
	"slices"
	"testing"
)

// agingReadingCases are hoisted out of the test body: the table is the
// interesting part and funlen counts it against the function otherwise.
var agingReadingCases = []struct {
	name      string
	profile   AgingProfile
	techLevel int
	term      int
	want      []string
}{
	{
		name: "a human on a printed table stamps nothing",
		// TL 10 ends on an open band, so no reading is reached.
		profile: ProfileTechLevel, techLevel: 10, term: 60,
		want: nil,
	},
	{
		name: "past the last printed term on a table that stops",
		// ERRATA E-15: TL 11 and TL 12-13 end at 58.
		profile: ProfileTechLevel, techLevel: 11, term: 59,
		want: []string{"E-15"},
	},
	{
		name: "a homeworld above the highest printed tech level",
		// ERRATA E-16: nothing is printed above TL 13.
		profile: ProfileTechLevel, techLevel: 14, term: 6,
		want: []string{"E-16"},
	},
	{
		name:    "both at once, on a high-tech world and a long life",
		profile: ProfileTechLevel, techLevel: 20, term: 59,
		want: []string{"E-16", "E-15"},
	},
	{
		name: "the term the sturdy table skips",
		// ERRATA E-30: it prints 6-8, 9-11 and 13+.
		profile: ProfileSturdy, techLevel: 12, term: sturdyGapTerm,
		want: []string{"E-30"},
	},
	{
		name:    "a sturdy character either side of the gap",
		profile: ProfileSturdy, techLevel: 12, term: 11,
		want: nil,
	},
}

// TestTheAgingTablesReadingsReachTheRecord is the three readings the printed
// aging tables rest on, each stamped where it fires.
//
// They are checked here rather than through a lifepath because two of them
// need a character the sample setting cannot make: E-30 wants a species on
// the sturdy profile, which setting/sample.json declares none of, and E-15
// wants a character past term 58, which needs a world allowing them. A
// reading nothing exercises is a reading nobody can trust.
func TestTheAgingTablesReadingsReachTheRecord(t *testing.T) {
	t.Parallel()

	for _, test := range agingReadingCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gen := &Generator{
				char:         &Character{},
				agingProfile: test.profile,
				techLevel:    test.techLevel,
			}

			gen.stampAgingReadings(test.term)

			got := gen.char.Provenance.Deviations

			if len(got) != len(test.want) {
				t.Fatalf("stamped %v, want %v", got, test.want)
			}

			for _, id := range test.want {
				if !slices.Contains(got, id) {
					t.Errorf("stamped %v, which does not include %s", got, id)
				}
			}
		})
	}
}
