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

// TestTheGraduateBonusIsTakenOnce is #132. A doctorate is a master's gone
// further -- degreeFor awards it on a second success at graduate school and
// leaves the master's on the record -- and p. 212 prints one modifier, not
// one per parchment. Reading the record for each in turn and adding both
// gave an Instructor enlistment +8 where the page gives +4, which is the
// difference between EDU 12+ and EDU 8+ on the hardest throw in the book.
//
// The value is the assertion. TestADegreeHelpsEnlistment counts that a
// degree modifier appeared at all, which is what let this through.
func TestTheGraduateBonusIsTakenOnce(t *testing.T) {
	t.Parallel()

	instructor := career.Instructor()

	var graduate career.EnlistmentMod

	for _, mod := range instructor.EnlistmentMods {
		if mod.Kind == career.GraduateDegree {
			graduate = mod
		}
	}

	if graduate.Value == 0 {
		t.Fatal("Instructor prints no graduate modifier; p. 212 gives +4")
	}

	cases := []struct {
		name  string
		held  []career.Degree
		want  int
		label string
	}{
		{"neither", nil, 0, ""},
		{"a master's", []career.Degree{career.Masters}, graduate.Value, "a master's"},
		{"a doctorate", []career.Degree{career.Doctorate}, graduate.Value, "a doctorate"},
		{
			"both",
			[]career.Degree{career.Masters, career.Doctorate},
			graduate.Value,
			"a doctorate",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			total, label := graduateBonus(instructor, tc.held)
			if total != tc.want {
				t.Errorf("held %v: graduate bonus %+d, want %+d", tc.held, total, tc.want)
			}

			if label != tc.label {
				t.Errorf("held %v: named %q, want %q", tc.held, label, tc.label)
			}
		})
	}
}

// graduateBonus sums the graduate-degree modifiers a career's enlistment
// throw would carry for a character holding exactly the degrees given, and
// names the last one. A correct engine names one and sums it once.
func graduateBonus(target career.Career, held []career.Degree) (int, string) {
	gen := &Generator{char: &Character{}}

	for _, degree := range held {
		gen.char.State.Education = append(gen.char.State.Education,
			&Education{Institution: career.GraduateSchool().Name, Degree: degree})
	}

	total, label := 0, ""

	for _, mod := range gen.enlistmentMods(target) {
		if mod.Name == "a master's" || mod.Name == "a doctorate" {
			total += mod.Value

			label = mod.Name
		}
	}

	return total, label
}
