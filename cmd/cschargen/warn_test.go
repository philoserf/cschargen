package main

import (
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
