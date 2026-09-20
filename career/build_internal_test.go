package career

import "testing"

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
