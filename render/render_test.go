package render_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/render"
)

// update rewrites the golden files instead of comparing against them. It is
// a deliberate act: `go test ./render/ -update`, then read the diff.
var update = flag.Bool("update", false, "rewrite the golden renders")

// character generates the fixture the goldens are taken from. The inputs
// are fixed, so a change to any of them is a change to every golden and
// shows up as one.
func character(t *testing.T, seed uint64, terms int, name string) *chargen.Character {
	t.Helper()

	got, err := chargen.New(chargen.Options{
		Seed:          seed,
		Decider:       chargen.Policy{},
		EngineVersion: "test",
		PolicyVersion: "test",
		SettingData:   chargen.SettingData{Name: "sample", Sample: true},
		Inputs: chargen.Inputs{
			Name: name, Species: "human", TermLimit: terms, TechLevel: 11, MaxTerms: 28,
		},
	}).Run()
	if err != nil {
		t.Fatalf("generating: %v", err)
	}

	return got
}

func golden(t *testing.T, name, got string) {
	t.Helper()

	path := filepath.Join("testdata", name)

	if *update {
		err := os.WriteFile(path, []byte(got), 0o600)
		if err != nil {
			t.Fatalf("writing %s: %v", path, err)
		}

		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s (run with -update to create it): %v", path, err)
	}

	if string(want) != got {
		t.Errorf("%s differs from the golden.\n--- got ---\n%s", path, got)
	}
}

func TestSheet(t *testing.T) {
	t.Parallel()

	golden(t, "sheet.md", render.Sheet(character(t, 7, 3, "Vela Ashgrove")))
}

func TestTranscript(t *testing.T) {
	t.Parallel()

	golden(t, "transcript.md", render.Transcript(character(t, 7, 3, "Vela Ashgrove")))
}

// TestSheetOfACharacterWithNothing: a record with no career service and no
// belongings still renders a sheet rather than a wall of empty headings.
func TestSheetOfACharacterWithNothing(t *testing.T) {
	t.Parallel()

	got := render.Sheet(character(t, 3, -1, ""))

	for _, want := range []string{"# (unnamed)", "## Characteristics", "_No career service._"} {
		if !strings.Contains(got, want) {
			t.Errorf("the sheet does not contain %q:\n%s", want, got)
		}
	}

	for _, unwanted := range []string{"## Relationships", "## Belongings"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("the sheet contains an empty %q section", unwanted)
		}
	}
}

// TestEveryThrowInTheTranscriptShowsItsDice: the transcript is the audit
// view, and an audit needs the dice as they fell rather than the total.
func TestEveryThrowInTheTranscriptShowsItsDice(t *testing.T) {
	t.Parallel()

	got := render.Transcript(character(t, 11, 4, "Test"))

	for line := range strings.SplitSeq(got, "\n") {
		if !strings.Contains(line, "throw") {
			continue
		}

		if !strings.Contains(line, "[") {
			t.Errorf("a throw line shows no dice: %q", line)
		}

		if !strings.Contains(line, "(p. ") {
			t.Errorf("a throw line carries no page cite: %q", line)
		}
	}
}

// TestTheSheetNamesTheSampleData: a character generated against the
// invented sample is not set on real worlds, and the sheet has to say so
// where a reader will see it.
func TestTheSheetNamesTheSampleData(t *testing.T) {
	t.Parallel()

	got := render.Sheet(character(t, 7, 2, "Test"))

	if !strings.Contains(got, "invented sample") {
		t.Errorf("the sheet does not say the setting data was the sample:\n%s", got)
	}
}

// TestTheSheetNamesTheReadingsApplied: an ERRATA reading changed the
// character, so the sheet says which, where someone checking it against the
// book would want to know.
func TestTheSheetNamesTheReadingsApplied(t *testing.T) {
	t.Parallel()

	got := render.Sheet(character(t, 7, 3, "Test"))

	if !strings.Contains(got, "ERRATA.md") {
		t.Errorf("the sheet does not name the readings applied:\n%s", got)
	}
}

// TestModifiersOnTheSheetAreComputed is p. 15 holding at the point of
// display: the sheet prints a modifier beside every score, and it is
// derived rather than stored.
func TestModifiersOnTheSheetAreComputed(t *testing.T) {
	t.Parallel()

	got := render.Sheet(character(t, 7, 3, "Test"))

	if !strings.Contains(got, "| STR | DEX | END | INT | EDU | CHA |") {
		t.Errorf("no characteristics table:\n%s", got)
	}

	if !strings.Contains(got, "(+") && !strings.Contains(got, "(-") {
		t.Errorf("no modifiers beside the scores:\n%s", got)
	}
}
