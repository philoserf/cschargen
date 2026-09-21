package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// served builds a character who entered the careers named, which is all
// these notes read.
func served(careers ...string) *chargen.Character {
	character := &chargen.Character{}
	for _, name := range careers {
		character.State.Services = append(character.State.Services,
			chargen.Service{Career: name})
	}

	return character
}

// TestTheNoteFiresOnlyWhenTheCareerChanged. `--career` restricts which
// career may be attempted; a failed enlistment still closes it for two
// terms (p. 110) and the character drifts. Across twenty seeds asking for
// Corporate Shipper, nine got one — and the command said nothing.
func TestTheNoteFiresOnlyWhenTheCareerChanged(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		character *chargen.Character
		asked     string
		want      string
	}{
		"no career was asked for": {
			character: served("Vagabond"), asked: "", want: "",
		},
		"the career asked for was entered": {
			character: served("Belter"), asked: "Belter", want: "",
		},
		"entered later, after another": {
			character: served("Vagabond", "Belter"), asked: "Belter", want: "",
		},
		"substituted": {
			character: served("Vagabond"), asked: "Belter", want: "entered Vagabond",
		},
		"no career at all": {
			character: served(), asked: "Belter", want: "no career at all",
		},
	} {
		got := &strings.Builder{}
		warnIfTheCareerChanged(got, test.character, test.asked)

		switch {
		case test.want == "" && got.Len() != 0:
			t.Errorf("%s: said %q, want silence", name, got)
		case test.want != "" && !strings.Contains(got.String(), test.want):
			t.Errorf("%s: said %q, want it to mention %q", name, got, test.want)
		}
	}
}

// TestABatchCountsTheSubstitutionsRatherThanRepeatingThem. A hundred
// records asking for one career could print fifty-five identical lines; a
// count is what a referee needs in order to decide whether to draw again.
func TestABatchCountsTheSubstitutionsRatherThanRepeatingThem(t *testing.T) {
	t.Parallel()

	batch := []*chargen.Character{
		served("Belter"), served("Vagabond"), served("Vagabond"), served("Belter"),
	}

	got := &strings.Builder{}
	warnHowManyCareersChanged(got, batch, "Belter")

	if !strings.Contains(got.String(), "2 of 4") {
		t.Errorf("said %q, want it to count 2 of 4", got)
	}

	silent := &strings.Builder{}
	warnHowManyCareersChanged(silent, []*chargen.Character{served("Belter")}, "Belter")

	if silent.Len() != 0 {
		t.Errorf("said %q when every character got the career asked for", silent)
	}

	unasked := &strings.Builder{}
	warnHowManyCareersChanged(unasked, batch, "")

	if unasked.Len() != 0 {
		t.Errorf("said %q when no career was asked for", unasked)
	}
}

// TestNamesComeFromTheFileAndNowhereElse. FR12 leaves Step 20's fields
// empty in auto mode rather than inventing them, which is right for a
// player character and leaves a cast of NPCs a hundred sheets headed
// "(unnamed)". The engine will use names it is handed; it still invents
// none.
func TestNamesComeFromTheFileAndNowhereElse(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "names.txt")

	err := os.WriteFile(path, []byte("# a comment\nVela Ashgrove\n\n  Toma Iri  \n"), 0o600)
	if err != nil {
		t.Fatalf("writing the fixture: %v", err)
	}

	names, err := readNames(path)
	if err != nil {
		t.Fatalf("readNames: %v", err)
	}

	if len(names) != 2 {
		t.Fatalf("read %v, want two names with the comment and blank skipped", names)
	}

	if names[1] != "Toma Iri" {
		t.Errorf("read %q, want the surrounding space trimmed", names[1])
	}

	// No file, no name — not a made-up one.
	if got := nameFor(nil, 7); got != "" {
		t.Errorf("with no list the engine named a character %q", got)
	}
}

// TestANameIsAFunctionOfTheSeed, so that a member of a batch and the same
// seed alone draw the same name, and two runs agree.
func TestANameIsAFunctionOfTheSeed(t *testing.T) {
	t.Parallel()

	names := []string{"Vela", "Toma", "Kess", "Nine"}

	for seed := range uint64(20) {
		first := nameFor(names, seed)
		if first != nameFor(names, seed) {
			t.Fatalf("seed %d drew two different names", seed)
		}

		if !slices.Contains(names, first) {
			t.Fatalf("seed %d drew %q, which is not on the list", seed, first)
		}
	}
}

// TestAnEmptyNameFileIsAUsageError, because a referee who passes one meant
// to name their cast and would otherwise get a silent hundred unnamed.
func TestAnEmptyNameFileIsAUsageError(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "empty.txt")

	err := os.WriteFile(path, []byte("# nothing but a comment\n"), 0o600)
	if err != nil {
		t.Fatalf("writing the fixture: %v", err)
	}

	_, err = readNames(path)
	if err == nil {
		t.Fatal("an empty name file was accepted")
	}

	if !strings.HasPrefix(err.Error(), "usage:") {
		t.Errorf("error %q does not begin with usage:", err)
	}

	// And a file that is not there at all.
	_, err = readNames(filepath.Join(t.TempDir(), "absent.txt"))
	if err == nil {
		t.Error("a missing name file was accepted")
	}
}

// TestADirectoryIsNotANameFile. `os.Open` opens a directory happily and
// fails on the first read, which is a realistic slip — a referee with a
// names file and a records directory beside each other.
func TestADirectoryIsNotANameFile(t *testing.T) {
	t.Parallel()

	_, err := readNames(t.TempDir())
	if err == nil {
		t.Fatal("a directory was accepted as a list of names")
	}

	if !strings.Contains(err.Error(), "reading") {
		t.Errorf("error %q does not say what failed", err)
	}
}

// TestJoinNames. The list goes into the message naming what a setting file
// does have, so it has to read as a sentence at every length the validator
// permits -- which includes one, since a file is valid with a single
// subsector in it and " and Earth" is not a list.
func TestJoinNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		names []string
		want  string
	}{
		{"one subsector, which a file may have", []string{"Earth"}, "Earth"},
		{"two", []string{"Earth", "Hub"}, "Earth and Hub"},
		{"more", []string{"Earth", "Hub", "Cascadia"}, "Earth, Hub and Cascadia"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := joinNames(test.names); got != test.want {
				t.Errorf("joinNames(%v) = %q, want %q", test.names, got, test.want)
			}
		})
	}
}
