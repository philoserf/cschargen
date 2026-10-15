package chargen

import (
	"slices"
	"testing"

	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/dice"
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

// TestAnEnlistmentIsNotCitedToThePreviousTable is #144. g.cite is the page
// the current table came from, and g.consequence stamps it onto everything
// it writes. Nothing cleared it when a career or a school ended, so Steps 9
// and 10 -- which happen with no career in hand -- wrote their consequences
// under whatever had been read last: a character was "accepted into Arts" on
// the pages of the career that had just turned them down, or on a graduate
// school's.
//
// The leak was correct about half the time, whenever the stale table was the
// one being talked about, which is what made it invisible to reading.
func TestAnEnlistmentIsNotCitedToThePreviousTable(t *testing.T) {
	t.Parallel()

	gen := &Generator{
		dice: dice.New(1),
		log:  &Log{},
		char: &Character{},
	}

	// A page from somewhere else entirely, as a career or a school would
	// have left behind.
	gen.cite = "pp. 86-91"

	gen.enlist(career.Vagabond())

	for _, event := range gen.log.Events() {
		if event.Consequence == nil {
			continue
		}

		if event.Consequence.Cite != enlistmentCite {
			t.Errorf("a Step 10 consequence cites %q, want %q: %s",
				event.Consequence.Cite, enlistmentCite, event.Consequence.Detail)
		}
	}
}

// TestLeavingACareerEndsItsPage holds the other half: the cite is cleared
// where the career ends, so a step that forgets to name its own page writes
// an empty one rather than a wrong one. An empty cite is a gap; a wrong one
// is a claim.
func TestLeavingACareerEndsItsPage(t *testing.T) {
	t.Parallel()

	gen := &Generator{
		dice: dice.New(1),
		log:  &Log{},
		char: &Character{},
	}

	entered := career.Vagabond()

	gen.service = serviceState{career: &entered}
	gen.cite = entered.Cite

	err := gen.leaveCareer(0, "testing")
	if err != nil {
		t.Fatalf("leaveCareer: %v", err)
	}

	if gen.cite != "" {
		t.Errorf("the career's page outlived the career: %q", gen.cite)
	}
}

// TestALossToOneClosesEveryCareerButVagabond is p. 119: "If any of the
// character's physical or mental characteristics reach 1, the character
// must either remain within the career they are in, leave character
// generation and join the game, or enter the Vagabond career" (#166).
//
// The cases are the reading ERRATA E-46 records: a floored 0 passed through
// 1, CHA is not one of the page's own "(INT or EDU)", a gain is not a
// reaching, and a later recovery reopens nothing.
func TestALossToOneClosesEveryCareerButVagabond(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		which  string
		from   int
		deltas []int
		closed bool
	}{
		{"a loss to 1", "STR", 4, []int{-3}, true},
		{"a loss floored at 0", "END", 2, []int{-5}, true},
		{"a mental loss to 1", "EDU", 3, []int{-2}, true},
		{"a later recovery", "DEX", 3, []int{-2, +4}, true},
		{"a loss that stops at 2", "INT", 4, []int{-2}, false},
		{"CHA to 1", "CHA", 4, []int{-3}, false},
		{"a gain to 1", "STR", 0, []int{+1}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gen := engine(t, 13)

			which, _ := characteristicByName(tc.which)
			gen.char.State.Characteristics.Set(which, tc.from)

			for _, delta := range tc.deltas {
				gen.adjust(tc.which, delta, "testing", 1)
			}

			eligible := gen.eligibleCareers()
			onlyVagabond := len(eligible) == 1 && eligible[0].Name == career.Vagabond().Name

			if onlyVagabond != tc.closed {
				t.Errorf("only Vagabond = %v, want %v (%d careers offered)",
					onlyVagabond, tc.closed, len(eligible))
			}

			stamped := slices.Contains(gen.char.Provenance.Deviations, "E-46")
			if stamped != tc.closed {
				t.Errorf("E-46 stamped = %v, want %v", stamped, tc.closed)
			}
		})
	}
}

// TestATransferAfterReachingOneIsVagabond is p. 119's second sentence: "If
// they are forced at this point to leave the career through a Mishap, they
// can either exit character generation or continue in the Vagabond career."
// A result that names Prisoner is the commonest case.
func TestATransferAfterReachingOneIsVagabond(t *testing.T) {
	t.Parallel()

	for _, dest := range []string{"Prisoner", career.Vagabond().Name} {
		t.Run(dest, func(t *testing.T) {
			t.Parallel()

			gen := engine(t, 14)

			gen.reachedOne = true
			gen.transfer = &career.Effect{Career: dest, Assignment: "Transient", Terms: 2}

			got, forced, err := gen.chooseCareer(0)
			if err != nil {
				t.Fatalf("chooseCareer: %v", err)
			}

			if got.Name != career.Vagabond().Name || !forced {
				t.Errorf("sent to %s (forced %v), want Vagabond", got.Name, forced)
			}

			// Only a transfer that was already to Vagabond keeps what it
			// carried; a redirected one ordered nothing of Vagabond's.
			kept := gen.forcedAssignment == "Transient" && gen.forcedTerms == 2
			if kept != (dest == career.Vagabond().Name) {
				t.Errorf("assignment %q for %d terms carried over from a transfer to %s",
					gen.forcedAssignment, gen.forcedTerms, dest)
			}
		})
	}
}

// TestReachingOneTakesHigherEducationOffTheTable is the Step 18 half of
// p. 119: of p. 125's five choices it leaves remaining and leaving, and
// returning to higher education is neither. The control proves the offer is
// there to be withheld.
func TestReachingOneTakesHigherEducationOffTheTable(t *testing.T) {
	t.Parallel()

	for _, reached := range []bool{false, true} {
		asked := &countingDecider{}

		gen := New(Options{
			Seed:          15,
			Decider:       asked,
			EngineVersion: testVersion,
			PolicyVersion: testVersion,
			Inputs:        Inputs{Species: testHuman, TermLimit: -1, Interactive: true},
		})

		for _, which := range CharacteristicOrder {
			gen.char.State.Characteristics.Set(which, 12)
		}

		gen.reachedOne = reached

		err := gen.mayReturnToEducation(0)
		if err != nil {
			t.Fatalf("mayReturnToEducation: %v", err)
		}

		if offered := asked.count > 0; offered == reached {
			t.Errorf("reached 1 = %v, offered higher education = %v", reached, offered)
		}
	}
}

// countingDecider declines every choice it is offered, counting them.
type countingDecider struct{ count int }

func (c *countingDecider) Choose(Choice) (int, error) {
	c.count++

	return 0, nil
}

func (*countingDecider) Kind() DeciderKind { return DeciderPlayer }
