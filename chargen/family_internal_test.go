package chargen

import "testing"

// TestMultipleBirthCountsUp is result 12 of the Sibling Age table (p. 60):
// "If this is the second sibling to get this result, the character and the
// two siblings are triplets. A third result will mean quadruplets, etc."
//
// The etc. is the part worth testing. A character can roll 12 more than
// three times, and a table that stopped at quadruplets would print nothing
// for the fifth. Past four the count is a numeral: the book stops giving
// words there and so does this.
func TestMultipleBirthCountsUp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		n    int
		want string
	}{
		{1, "the character's twin"},
		{2, "one of three, with the character"},
		{3, "one of four, with the character"},
		{4, "one of 5, with the character"},
		{9, "one of 10, with the character"},
	}

	for _, tc := range tests {
		if got := multipleBirth(tc.n); got != tc.want {
			t.Errorf("multipleBirth(%d) = %q, want %q", tc.n, got, tc.want)
		}
	}
}
