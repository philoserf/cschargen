package render_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/render"
)

// TestARosterEntryCarriesWhatCastingAsks is the line a referee reads when
// choosing a crew from a pool: who, how old, what they did, what they are
// best at, and what they have.
func TestARosterEntryCarriesWhatCastingAsks(t *testing.T) {
	t.Parallel()

	got := render.Roster([]*chargen.Character{character(t, 7, 4, "Vela Ashgrove")})

	for _, want := range []string{"Vela Ashgrove", "age ", "Cr", "beyond the family"} {
		if !strings.Contains(got, want) {
			t.Errorf("the roster entry does not carry %q:\n%s", want, got)
		}
	}

	if lines := strings.Count(strings.TrimRight(got, "\n"), "\n") + 1; lines != 3 {
		t.Errorf("one character is %d lines, want 3:\n%s", lines, got)
	}
}

// TestARosterNumbersEveryCharacter, so that a referee can say "number
// seven" and mean it, and so the numbering matches `batch -o dir`'s files.
func TestARosterNumbersEveryCharacter(t *testing.T) {
	t.Parallel()

	pool := []*chargen.Character{
		character(t, 1, 2, ""),
		character(t, 2, 2, ""),
		character(t, 3, 2, ""),
	}

	got := render.Roster(pool)

	for _, want := range []string{"  1  ", "  2  ", "  3  "} {
		if !strings.Contains(got, want) {
			t.Errorf("the roster does not number %q:\n%s", strings.TrimSpace(want), got)
		}
	}
}

// TestARosterNamesTheUnnamed. An auto-generated character has no name --
// Step 20 is left empty rather than invented -- and a blank where
// the name goes reads as a broken line rather than as a fact.
func TestARosterNamesTheUnnamed(t *testing.T) {
	t.Parallel()

	got := render.Roster([]*chargen.Character{character(t, 11, 2, "")})
	if !strings.Contains(got, "(unnamed)") {
		t.Errorf("an unnamed character renders as:\n%s", got)
	}
}

// TestARosterCountsTiesWithoutTheFamily. Ninety per cent of a character's
// relationships are relatives; the ones a plot can hang on are the rest,
// and a count that buries three colleagues under sixty cousins tells a
// referee nothing.
func TestARosterCountsTiesWithoutTheFamily(t *testing.T) {
	t.Parallel()

	built := character(t, 7, 4, "Test")

	family, other := 0, 0

	for _, tie := range built.State.Ties {
		if tie.Origin == chargen.FamilyOrigin {
			family++

			continue
		}

		other++
	}

	if family == 0 || other == 0 {
		t.Skip("this character has no mix of family and other ties to distinguish")
	}

	got := render.Roster([]*chargen.Character{built})
	if !strings.Contains(got, plural(other)+" beyond the family") {
		t.Errorf("want %d ties beyond the family (of %d), got:\n%s",
			other, family+other, got)
	}
}

// plural matches the roster's own counting, so the assertion reads the way
// the output does.
func plural(n int) string {
	if n == 1 {
		return "1 tie"
	}

	return strconv.Itoa(n) + " ties"
}

// TestOneTieIsATie. "1 ties" reads as a bug in the tool rather than as a
// fact about the character, and a character with exactly one colleague is
// common enough to be worth the branch.
func TestOneTieIsATie(t *testing.T) {
	t.Parallel()

	one := &chargen.Character{}

	one.State.Ties = []chargen.Tie{{Kind: "ally", Origin: "Belter", Rating: 150}}

	got := render.Roster([]*chargen.Character{one})
	if !strings.Contains(got, "1 tie beyond the family") {
		t.Errorf("a character with one tie renders as:\n%s", got)
	}
}

// TestARosterOfNobodyIsEmpty, because a batch that generated nothing is not
// an error the roster should invent a line for.
func TestARosterOfNobodyIsEmpty(t *testing.T) {
	t.Parallel()

	if got := render.Roster(nil); got != "" {
		t.Errorf("an empty pool rendered as %q", got)
	}
}

// TestARosterOfSomebodyWithNothing. A character can reach Step 20 with no
// career and no skills — three failed enlistments and a term limit of one
// will do it — and the line should read as a fact rather than as two gaps.
func TestARosterOfSomebodyWithNothing(t *testing.T) {
	t.Parallel()

	got := render.Roster([]*chargen.Character{{}})

	for _, want := range []string{"no career", "no skills"} {
		if !strings.Contains(got, want) {
			t.Errorf("a character with nothing renders without %q:\n%s", want, got)
		}
	}
}

// TestARosterLineSaysWhereTheyAreFrom. The roster is the casting view, and
// a homeworld is what places an NPC: it carries the primary language, the
// tech level they grew up at, and whether engineered people are free there.
//
// It was on the sheet and not on the line, so casting a scene from a sector
// meant opening every record to find out where anyone came from.
func TestARosterLineSaysWhereTheyAreFrom(t *testing.T) {
	t.Parallel()

	character := &chargen.Character{
		State: chargen.State{
			Age: 34,
			Homeworlds: []chargen.Homeworld{
				{World: "Avicenna", Subsector: "Kestrel Reach"},
				{World: "Elsewhere", Subsector: "The Verge", FromTerm: 2},
			},
		},
	}

	line := render.Roster([]*chargen.Character{character})

	if !strings.Contains(line, "Avicenna") {
		t.Errorf("the roster line does not say where they are from:\n%s", line)
	}

	// The birth world, not wherever a mishap moved them to: the history
	// begins where they were born, and that is what places them.
	if strings.Contains(line, "Elsewhere") {
		t.Errorf("the roster line names a later homeworld:\n%s", line)
	}
}

// TestARosterLineWithNoHomeworld is a record made before Step 4 ran, which
// `--terms -1` produces: the line still has to render.
func TestARosterLineWithNoHomeworld(t *testing.T) {
	t.Parallel()

	character := &chargen.Character{State: chargen.State{Age: 18}}

	line := render.Roster([]*chargen.Character{character})
	if line == "" {
		t.Error("a character with no homeworld rendered no line")
	}
}
