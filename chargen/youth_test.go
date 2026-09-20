package chargen_test

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/chargen"
)

// TestTheFivePathsAreFiveNineteenRowTables is the shape of pp. 68-74: each
// path is a 2d10 table, so nineteen rows for a result of 2 to 20, and each
// opens on the same orphaning and reaches the Youth Life Events table at
// result 10.
func TestTheFivePathsAreFiveNineteenRowTables(t *testing.T) {
	t.Parallel()

	paths := career.YouthPaths()
	if len(paths) != 5 {
		t.Fatalf("%d youth paths; pp. 68-74 print five", len(paths))
	}

	for _, path := range paths {
		checkYouthPath(t, path)
	}
}

// checkYouthPath holds one path against the shape every one of them has.
func checkYouthPath(t *testing.T, path career.YouthPath) {
	t.Helper()

	if len(path.Rows) != 19 {
		t.Errorf("%s has %d rows; a 2d10 table has 19", path.Name, len(path.Rows))
	}

	if path.Cite == "" {
		t.Errorf("%s has no page cite", path.Name)
	}

	for i, row := range path.Rows {
		if row.Summary == "" {
			t.Errorf("%s result %d has no summary", path.Name, i+2)
		}
	}

	// Result 2 is the orphaning, on every path.
	if !strings.Contains(path.Rows[0].Summary, "orphaned") {
		t.Errorf("%s result 2 is %q, want the orphaning", path.Name, path.Rows[0].Summary)
	}

	// Result 10 is the Youth Life Events table, on every path.
	for _, effect := range path.Rows[8].Effects {
		if effect.Kind == career.EffectYouthLifeEvent {
			return
		}
	}

	t.Errorf("%s result 10 does not reach the Youth Life Events table", path.Name)
}

// TestPathOneIsAlwaysOpenAndTheOthersAreNot is p. 68: "Simply check the
// requirements for each path, determine which path for which your character
// qualifies."
func TestPathOneIsAlwaysOpenAndTheOthersAreNot(t *testing.T) {
	t.Parallel()

	nothing := func(string) int { return 2 }
	everything := func(string) int { return 15 }

	paths := career.YouthPaths()

	if !paths[0].Open(nothing) {
		t.Error("Path 1 has requirements; p. 68 prints none")
	}

	for _, path := range paths[1:] {
		if path.Open(nothing) {
			t.Errorf("%s is open to a character with nothing", path.Name)
		}

		if !path.Open(everything) {
			t.Errorf("%s is closed to a character with everything", path.Name)
		}
	}
}

// TestYouthLifeEventsIsASixRowTable is p. 75, which prints results 3-4 as
// one row and so has six entries for six results.
func TestYouthLifeEventsIsASixRowTable(t *testing.T) {
	t.Parallel()

	table := career.YouthLifeEvents()
	if len(table) != 6 {
		t.Fatalf("%d youth life events; p. 75 prints a d6", len(table))
	}

	// Results 3 and 4 are printed as one row and are the same result.
	if table[2].Summary != table[3].Summary {
		t.Errorf("results 3 and 4 differ: %q and %q", table[2].Summary, table[3].Summary)
	}
}

// TestAYouthIsGenerated. Two rolls happen for every character who did not
// skip the step, and both are on a path the character qualified for.
func TestAYouthIsGenerated(t *testing.T) {
	t.Parallel()

	offered := 0
	lifeEvents := 0

	for seed := range uint64(sample) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = -1

		character := generate(t, opts)

		rolls, chose, events := countYouth(character)

		offered += chose
		lifeEvents += events

		if rolls != 2 {
			t.Errorf("seed %d: %d youth events, want 2 (ages 4-8 and 9-12)", seed, rolls)
		}
	}

	// The auto policy always takes the first option, which is always Path 1,
	// so what varies across seeds is how many paths were open rather than
	// which was taken. A run where nobody ever qualified for a second path
	// would mean the gates were never evaluated.
	if offered == 0 {
		t.Errorf("in %d seeds no character ever qualified for more than one path", sample)
	}

	if lifeEvents == 0 {
		t.Errorf("no seed in %d reached the Youth Life Events table", sample)
	}
}

// countYouth walks a record for the three things Step 6 leaves behind: the
// two rolls, the choices where more than one path was open, and the trips
// to the Youth Life Events table.
func countYouth(character *chargen.Character) (int, int, int) {
	var rolls, chose, events int

	for _, event := range character.Events {
		if event.Kind == chargen.EventChoice && event.Choice.Point == "youth_path" {
			chose++
		}

		if event.Consequence == nil {
			continue
		}

		detail := event.Consequence.Detail

		switch {
		case strings.HasPrefix(detail, "ages 4-8, on "),
			strings.HasPrefix(detail, "ages 9-12, on "):
			rolls++
		case strings.HasPrefix(detail, "youth life event: "):
			events++
		}
	}

	return rolls, chose, events
}

// TestSkippingTheYouth. p. 67: "Players may create this background entirely
// on their own."
func TestSkippingTheYouth(t *testing.T) {
	t.Parallel()

	opts := options(t, 5)

	opts.Inputs.TermLimit = -1
	opts.Inputs.SkipYouth = true

	character := generate(t, opts)

	for _, event := range character.Events {
		if event.Consequence == nil {
			continue
		}

		if strings.HasPrefix(event.Consequence.Detail, "ages ") {
			t.Errorf("a youth event happened although the step was skipped: %s",
				event.Consequence.Detail)
		}
	}

	// And the family still ran, because they are separate steps.
	if character.State.Family == nil {
		t.Error("skipping Step 6 skipped Step 5 as well")
	}
}

// TestAParentsRatingCanBeAimedAt. Two youth results address a parent
// specifically -- "Choose one of your parents and that parent will leave
// your life" -- and a target the engine cannot aim is a result it cannot
// carry out.
func TestAParentsRatingCanBeAimedAt(t *testing.T) {
	t.Parallel()

	moved := 0

	for seed := range uint64(200) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = -1

		character := generate(t, opts)

		for _, event := range character.Events {
			if event.Consequence == nil ||
				event.Consequence.Kind != chargen.ConsequenceRelationship {
				continue
			}

			if event.Consequence.Career == chargen.FamilyOrigin &&
				strings.Contains(event.Consequence.Detail, " -> ") {
				moved++
			}
		}
	}

	if moved == 0 {
		t.Errorf("in 200 seeds no family member's Relationship Rating ever moved")
	}
}
