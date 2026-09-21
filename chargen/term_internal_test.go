package chargen

import (
	"testing"

	"github.com/philoserf/cschargen/career"
)

// TestTheUnchoosableCareersAreNeverOffered holds the claim
// ownedFirstCareer makes in prose -- that p. 42 "is the only way into that
// career" -- against the list a character actually chooses from.
//
// Both careers state it themselves. Prisoner: "A character cannot
// voluntarily choose to enter this career" (p. 111). The slave career:
// "no enlistment throw, and it cannot be chosen: entry is a result of the
// early life tables, and the character must be an engineered human or an
// uplift" (p. 150). Neither takes an enlistment throw, so neither can fail
// one -- a character offered either would simply enter it.
func TestTheUnchoosableCareersAreNeverOffered(t *testing.T) {
	t.Parallel()

	unchoosable := []string{"Prisoner", career.Slave().Name}

	gen := engine(t, 11)

	for _, def := range gen.eligibleCareers() {
		for _, name := range unchoosable {
			if def.Name == name {
				t.Errorf("%s is on the list of careers a character may attempt", name)
			}
		}
	}
}

// TestEveryOtherCareerIsOffered is the other half: the exclusion above is
// two names, not a filter that quietly drops more. The book prints
// thirty-four careers (pp. 8-9) and a character with no history may attempt
// every one that is not one of the two.
func TestEveryOtherCareerIsOffered(t *testing.T) {
	t.Parallel()

	const unchoosable = 2

	gen := engine(t, 12)

	offered := map[string]bool{}
	for _, def := range gen.eligibleCareers() {
		offered[def.Name] = true
	}

	for _, def := range career.All() {
		if def.Name == "Prisoner" || def.Name == career.Slave().Name {
			continue
		}

		if !offered[def.Name] {
			t.Errorf("%s is not offered to a character with no history", def.Name)
		}
	}

	if want := len(career.All()) - unchoosable; len(offered) != want {
		t.Errorf("offered %d careers, want %d: the book prints %d and two cannot be chosen",
			len(offered), want, len(career.All()))
	}
}
