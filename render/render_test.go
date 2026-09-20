package render_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/render"
	"github.com/philoserf/cschargen/setting"
)

// update rewrites the golden files instead of comparing against them. It is
// a deliberate act: `go test ./render/ -update`, then read the diff.
var update = flag.Bool("update", false, "rewrite the golden renders")

// testVersion is what the fixtures stamp, so a golden's provenance line does
// not move with the build.
const testVersion = "test"

// character generates the fixture the goldens are taken from. The inputs
// are fixed, so a change to any of them is a change to every golden and
// shows up as one.
func character(t *testing.T, seed uint64, terms int, name string) *chargen.Character {
	t.Helper()

	got, err := chargen.New(chargen.Options{
		Seed:          seed,
		Decider:       chargen.Policy{},
		EngineVersion: testVersion,
		PolicyVersion: testVersion,
		Setting:       sampleSetting(t),
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

	// The four optional pre-career steps are skipped so that "nothing"
	// really is nothing: a character with a family has relatives, and the
	// youth and teenage tables hand out friends.
	built, err := chargen.New(chargen.Options{
		Seed:          3,
		Decider:       chargen.Policy{},
		EngineVersion: testVersion,
		PolicyVersion: testVersion,
		Setting:       sampleSetting(t),
		Inputs: chargen.Inputs{
			Species: "human", TermLimit: -1, TechLevel: 11, MaxTerms: 28,
			SkipFamily: true, SkipYouth: true, SkipTeenage: true, SkipEducation: true,
		},
	}).Run()
	if err != nil {
		t.Fatalf("generating: %v", err)
	}

	got := render.Sheet(built)

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

// TestTheSheetPricesOnlyWhatTheBookPriced is FR11 on the page: a company
// share or a piece of art carries the value the book rolled for it, and a
// weapon of the character's choice carries none, because the book gave
// none. A zero there means the referee sets it, not that it is worthless --
// so the sheet prints nothing rather than "0 credits".
func TestTheSheetPricesOnlyWhatTheBookPriced(t *testing.T) {
	t.Parallel()

	priced, unpriced := false, false

	for seed := range uint64(80) {
		built := character(t, seed, 6, "Test")

		for _, possession := range built.State.Stash {
			line := "**Stash**: "

			got := render.Sheet(built)
			if !strings.Contains(got, line) {
				t.Fatalf("seed %d holds %q and the sheet has no stash", seed, possession.Item)
			}

			if possession.Value > 0 {
				priced = true

				if !strings.Contains(got, possession.Item+" (") {
					t.Errorf("seed %d: %q is priced and the sheet does not say so",
						seed, possession.Item)
				}

				continue
			}

			unpriced = true

			if strings.Contains(got, possession.Item+" (0 credits)") {
				t.Errorf("seed %d: %q the book did not price is shown as worthless",
					seed, possession.Item)
			}
		}

		if priced && unpriced {
			return
		}
	}

	t.Fatalf("no seed in the sample held both a priced and an unpriced possession "+
		"(priced %v, unpriced %v)", priced, unpriced)
}

// TestTheSheetNamesWhatTheCharacterCarries is the conditions line: an
// addiction or a religion is neither a skill, a characteristic nor a
// possession, and the sheet says the character has one without saying what
// it costs them.
func TestTheSheetNamesWhatTheCharacterCarries(t *testing.T) {
	t.Parallel()

	for seed := range uint64(120) {
		built := character(t, seed, 6, "Test")
		if len(built.State.Conditions) == 0 {
			continue
		}

		got := render.Sheet(built)
		for _, carried := range built.State.Conditions {
			if !strings.Contains(got, carried) {
				t.Errorf("seed %d carries %q and the sheet does not say so", seed, carried)
			}
		}

		return
	}

	t.Skip("no seed in the sample carries a condition")
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

		// A cite is "(p. 110)" or "(pp. 68-69)": some tables run over
		// a page and the transcript says so.
		if !strings.Contains(line, "(p. ") && !strings.Contains(line, "(pp. ") {
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

func sampleSetting(t *testing.T) *setting.Data {
	t.Helper()

	data, err := setting.Sample()
	if err != nil {
		t.Fatalf("loading the sample setting: %v", err)
	}

	return data
}

// moved finds a seed whose character was reassigned a homeworld, so that
// the parts of the sheet only such a character reaches are rendered.
func moved(t *testing.T) *chargen.Character {
	t.Helper()

	for seed := range uint64(60) {
		got := character(t, seed, 8, "Wanderer")
		if len(got.State.Homeworlds) > 1 {
			return got
		}
	}

	t.Skip("no character in the sample was ever reassigned a homeworld")

	return nil
}

// TestTheSheetShowsWhereTheyHaveLived: a character who was deported twice
// has a history, and it belongs on the sheet -- where they lived is part of
// what happened to them.
func TestTheSheetShowsWhereTheyHaveLived(t *testing.T) {
	t.Parallel()

	got := render.Sheet(moved(t))

	if !strings.Contains(got, "## Where they have lived") {
		t.Errorf("a character who moved has no history on their sheet:\n%s", got)
	}

	if !strings.Contains(got, "| World | Subsector | TL | From term | Why |") {
		t.Errorf("the history has no table header:\n%s", got)
	}
}

// TestTheSheetOmitsAHistoryOfOne: a character who never moved should not
// get a one-row table repeating the homeworld already in the summary.
func TestTheSheetOmitsAHistoryOfOne(t *testing.T) {
	t.Parallel()

	for seed := range uint64(60) {
		got := character(t, seed, 2, "Settled")
		if len(got.State.Homeworlds) != 1 {
			continue
		}

		sheet := render.Sheet(got)
		if strings.Contains(sheet, "## Where they have lived") {
			t.Errorf("a character who never moved has a history table:\n%s", sheet)
		}

		return
	}

	t.Skip("every character in the sample moved at least once")
}

// TestTheSheetNamesTheHomeworldAndLanguage: both come from Step 4 and both
// are things a referee reads off the sheet.
func TestTheSheetNamesTheHomeworldAndLanguage(t *testing.T) {
	t.Parallel()

	got := render.Sheet(character(t, 7, 3, "Vela Ashgrove"))

	for _, want := range []string{"**Homeworld**:", "**Primary language**:"} {
		if !strings.Contains(got, want) {
			t.Errorf("the sheet has no %q line:\n%s", want, got)
		}
	}
}

// TestTheSheetSaysWhenAgingEndedIt. A character who dies or is
// incapacitated during generation has fewer terms than their homeworld
// allows, and a sheet that did not say why would look like a truncated
// record rather than a finished one (pp. 123-124).
func TestTheSheetSaysWhenAgingEndedIt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		fate chargen.Fate
		want string
	}{
		{chargen.FateDied, "**Died** at age"},
		{chargen.FateIncapacitated, "**Incapacitated** by aging"},
	}

	for _, tc := range tests {
		t.Run(string(tc.fate), func(t *testing.T) {
			t.Parallel()

			got := character(t, 3, 2, "Ended")

			got.State.Fate = tc.fate

			sheet := render.Sheet(got)
			if !strings.Contains(sheet, tc.want) {
				t.Errorf("the sheet of a character who %s does not say so:\n%s", tc.fate, sheet)
			}
		})
	}

	// And a character who finished ordinarily says nothing of the kind.
	alive := render.Sheet(character(t, 3, 2, "Ended"))
	for _, unwanted := range []string{"**Died**", "**Incapacitated**"} {
		if strings.Contains(alive, unwanted) {
			t.Errorf("a living character's sheet says %s", unwanted)
		}
	}
}

// TestApparentAgeOnlyAppearsWhenItSaysSomething. Below tech level 10, and
// below age 30 at any tech level, apparent age is the character's age
// (pp. 124-125) -- and a line repeating the one above it is noise.
func TestApparentAgeOnlyAppearsWhenItSaysSomething(t *testing.T) {
	t.Parallel()

	got := character(t, 3, 2, "Young")
	if strings.Contains(render.Sheet(got), "**Apparent age**") {
		t.Error("a character whose apparent age is their age has a line saying so")
	}

	got.State.Age = 120
	got.State.ApparentAge = chargen.AgeBand{From: 45, To: 50}

	if !strings.Contains(render.Sheet(got), "**Apparent age**: 45-50") {
		t.Error("the sheet does not carry an apparent age the chart supplied")
	}
}

// TestTheSheetPrintsEveryEducationOutcome. Step 8 has four ends -- not
// admitted, admitted and left, a degree, a degree with honours -- and each
// reads differently on the sheet.
func TestTheSheetPrintsEveryEducationOutcome(t *testing.T) {
	t.Parallel()

	got := character(t, 3, 1, "Scholar")

	got.State.Education = []*chargen.Education{
		{
			Institution: "Undergraduate College", Admitted: true, Succeeded: true,
			Field: "Broker", Degree: "bachelor's",
		},
		{
			Institution: "Graduate School", Admitted: true, Succeeded: true, Honors: true,
			Field: "Broker", Degree: "master's",
		},
		{Institution: "Medical School", Admitted: true},
		{Institution: "Military Academy"},
	}

	sheet := render.Sheet(got)

	for _, want := range []string{
		"- Undergraduate College: a bachelor's in Broker",
		"- Graduate School: a master's in Broker, with honours",
		"- Medical School: admitted, left without a degree",
		"- Military Academy: not admitted",
	} {
		if !strings.Contains(sheet, want) {
			t.Errorf("the sheet does not contain %q:\n%s", want, sheet)
		}
	}

	// And a character who never attempted it has no section at all.
	got.State.Education = nil

	if strings.Contains(render.Sheet(got), "## Education") {
		t.Error("a character with no education has an Education section")
	}
}

// TestTheSheetSaysWhereInTheFamilyTheyCame. The parental age die decides
// whether a character was the first child or the last (p. 59), and both
// ends are worth a line on the sheet.
func TestTheSheetSaysWhereInTheFamilyTheyCame(t *testing.T) {
	t.Parallel()

	got := character(t, 3, 1, "Eldest")

	got.State.Family = &chargen.Family{
		Situation: "a heterosexual couple", Detail: "a married couple", Firstborn: true,
	}

	if !strings.Contains(render.Sheet(got), "The eldest child.") {
		t.Error("a firstborn's sheet does not say so")
	}

	got.State.Family.Firstborn = false
	got.State.Family.Lastborn = true

	if !strings.Contains(render.Sheet(got), "The youngest child.") {
		t.Error("a last child's sheet does not say so")
	}

	// A character whose parents rolled neither end gets neither line.
	got.State.Family.Lastborn = false

	sheet := render.Sheet(got)
	for _, unwanted := range []string{"The eldest child.", "The youngest child."} {
		if strings.Contains(sheet, unwanted) {
			t.Errorf("a middle child's sheet says %q", unwanted)
		}
	}
}

// TestTheSheetNamesAnUpliftsClass. It is what their homeworld's technology
// could make of them (p. 66), so it belongs beside the species rather than
// under it -- and a character who is not an uplift has none.
func TestTheSheetNamesAnUpliftsClass(t *testing.T) {
	t.Parallel()

	got := character(t, 3, 1, "Uplifted")

	got.State.UpliftClass = 2

	if !strings.Contains(render.Sheet(got), ", Class 2") {
		t.Error("an uplift's class is not on their sheet")
	}

	got.State.UpliftClass = 0

	if strings.Contains(render.Sheet(got), ", Class") {
		t.Error("a character who is not an uplift has a class")
	}
}

// TestTheSheetPrintsGenetics. An engineered character's parentage is what
// pp. 62-65 made of it; a human and an uplift have none.
func TestTheSheetPrintsGenetics(t *testing.T) {
	t.Parallel()

	got := character(t, 3, 1, "Engineered")

	got.State.Genetics = &chargen.Genetics{
		Kind: "hybrid", Detail: "a hybrid with a baseline human: otherwise human",
	}

	if !strings.Contains(render.Sheet(got), "**Genetics**: a hybrid with a baseline human") {
		t.Error("the sheet does not carry the character's genetics")
	}

	got.State.Genetics = nil

	if strings.Contains(render.Sheet(got), "**Genetics**") {
		t.Error("a character with no genetics has a line saying so")
	}
}

// TestTheSheetCarriesTheFinishingTouches. Step 20's four fields, and only
// the ones somebody supplied -- the engine invents none, and a blank line
// is worse than no line.
func TestTheSheetCarriesTheFinishingTouches(t *testing.T) {
	t.Parallel()

	got := character(t, 3, 1, "Vela Ashgrove")

	got.State.Finishing = chargen.Finishing{
		Name: "Vela Ashgrove", Gender: "she/her",
		Appearance: "tall, with a spacer's stoop",
		Goals:      "to find out what happened to the Kestrel",
	}

	sheet := render.Sheet(got)
	for _, want := range []string{
		"**Gender**: she/her",
		"**Appearance**: tall, with a spacer's stoop",
		"**Goals**: to find out what happened to the Kestrel",
	} {
		if !strings.Contains(sheet, want) {
			t.Errorf("the sheet does not contain %q", want)
		}
	}

	// A character nobody described has none of those lines.
	got.State.Finishing = chargen.Finishing{}

	bare := render.Sheet(got)
	for _, unwanted := range []string{"**Gender**", "**Appearance**", "**Goals**"} {
		if strings.Contains(bare, unwanted) {
			t.Errorf("an undescribed character's sheet has %s", unwanted)
		}
	}
}
