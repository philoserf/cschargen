package chargen_test

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/setting"
)

// TestTheFourTeenagePathsAreFourNineteenRowTables is the shape of
// pp. 76-83, and that each reaches the Teenage Life Events table at result
// 10.
func TestTheFourTeenagePathsAreFourNineteenRowTables(t *testing.T) {
	t.Parallel()

	paths := career.TeenagePaths()
	if len(paths) != 4 {
		t.Fatalf("%d teenage paths; pp. 76-83 print four", len(paths))
	}

	for _, path := range paths {
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

		reaches := false

		for _, effect := range path.Rows[8].Effects {
			if effect.Kind == career.EffectTeenageLifeEvent {
				reaches = true
			}
		}

		if !reaches {
			t.Errorf("%s result 10 does not reach the Teenage Life Events table", path.Name)
		}
	}
}

// TestPathsOneAndTwoPartitionTheHomeworld. p. 76 gates them on whether the
// homeworld "has been colonized or established for 100+ standard years",
// which is a condition with exactly two sides: every character has one of
// them and never both.
func TestPathsOneAndTwoPartitionTheHomeworld(t *testing.T) {
	t.Parallel()

	nothing := func(string) int { return 2 }
	paths := career.TeenagePaths()

	for _, years := range []int{0, 1, 99, 100, 101, 500} {
		one := paths[0].Open(nothing, years)
		two := paths[1].Open(nothing, years)

		if one == two {
			t.Errorf("a homeworld settled %d years opens Path 1 = %v and Path 2 = %v; "+
				"exactly one should be open", years, one, two)
		}
	}

	// And the hundredth year is on Path 1's side: "100+".
	if !paths[0].Open(nothing, career.SettledYearsForPathOne) {
		t.Error("a homeworld settled exactly 100 years does not open Path 1")
	}
}

// TestPathsThreeAndFourAreCharacteristicGates.
func TestPathsThreeAndFourAreCharacteristicGates(t *testing.T) {
	t.Parallel()

	nothing := func(string) int { return 8 }
	everything := func(string) int { return 9 }

	for _, path := range career.TeenagePaths()[2:] {
		if path.Open(nothing, 0) {
			t.Errorf("%s is open to a character at 8 across the board; it wants 9+", path.Name)
		}

		if !path.Open(everything, 0) {
			t.Errorf("%s is closed to a character at 9 across the board", path.Name)
		}
	}
}

// TestTeenageLifeEventsIsASixRowTable is p. 85.
func TestTeenageLifeEventsIsASixRowTable(t *testing.T) {
	t.Parallel()

	if got := len(career.TeenageLifeEvents()); got != 6 {
		t.Fatalf("%d teenage life events; p. 85 prints a d6", got)
	}
}

// TestTheCampaignsPresentYear is ERRATA E-27: the book gives 2350 on p. 43
// and 2345 on p. 124, and a data file may give its own.
func TestTheCampaignsPresentYear(t *testing.T) {
	t.Parallel()

	var empty setting.Data

	if got := empty.Now(); got != setting.DefaultPresentYear {
		t.Errorf("a file with no present year says %d, want %d", got, setting.DefaultPresentYear)
	}

	stated := setting.Data{PresentYear: 2400}
	if got := stated.Now(); got != 2400 {
		t.Errorf("a file that says 2400 says %d", got)
	}
}

// TestATeenageIsGenerated. Two rolls happen for every character who did not
// skip the step, and both paths of the homeworld condition are reached
// across the sample data's thirteen worlds.
func TestATeenageIsGenerated(t *testing.T) {
	t.Parallel()

	onPath := map[string]int{}
	lifeEvents := 0

	for seed := range uint64(sample) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = -1

		character := generate(t, opts)

		rolls := 0

		for _, event := range character.Events {
			if event.Consequence == nil {
				continue
			}

			detail := event.Consequence.Detail

			switch {
			case strings.HasPrefix(detail, "ages 13-15, on "),
				strings.HasPrefix(detail, "ages 16-18, on "):
				rolls++

				onPath[detail[strings.Index(detail, "on ")+3:][:6]]++
			case strings.HasPrefix(detail, "teenage life event: "):
				lifeEvents++
			}
		}

		if rolls != 2 {
			t.Errorf("seed %d: %d teenage events, want 2", seed, rolls)
		}
	}

	// Paths 1 and 2 are chosen by the homeworld rather than by the policy,
	// so both should appear across the sample's worlds.
	if onPath["Path 1"] == 0 || onPath["Path 2"] == 0 {
		t.Errorf("Path 1 taken %d times and Path 2 %d; the sample's worlds should reach both",
			onPath["Path 1"], onPath["Path 2"])
	}

	if lifeEvents == 0 {
		t.Errorf("no seed in %d reached the Teenage Life Events table", sample)
	}
}

// TestSkippingTheTeenage.
func TestSkippingTheTeenage(t *testing.T) {
	t.Parallel()

	opts := options(t, 5)

	opts.Inputs.TermLimit = -1
	opts.Inputs.SkipTeenage = true

	character := generate(t, opts)

	for _, event := range character.Events {
		if event.Consequence == nil {
			continue
		}

		detail := event.Consequence.Detail
		if strings.HasPrefix(detail, "ages 13-15") || strings.HasPrefix(detail, "ages 16-18") {
			t.Errorf("a teenage event happened although the step was skipped: %s", detail)
		}
	}

	// Steps 5 and 6 still ran, because they are separate steps.
	if character.State.Family == nil {
		t.Error("skipping Step 7 skipped Step 5 as well")
	}
}
