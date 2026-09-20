package chargen_test

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// TestACharacterActuallyAges is the reach check the milestone 4 plan makes
// against the sample data: aging begins at term 6 on a tech level 9 world,
// and two of the sample's thirteen worlds are tech level 9 with term caps
// of 17 and 19.
//
// The closing assertion is what keeps this honest. A test that only looked
// for a lost point would pass silently the day no seed reaches a tech level
// 9 homeworld, which is how the Colonist reassignment test learned to fail
// instead.
func TestACharacterActuallyAges(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)
	lowTech := 0
	aged := 0

	for seed := range uint64(sample) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = 40

		character := generate(t, opts)

		born, _, found := data.World(character.State.Homeworlds[0].World)
		if !found {
			t.Fatalf("seed %d: unknown birth world", seed)
		}

		if born.TechLevel > 9 || len(character.State.Terms) < 6 {
			continue
		}

		lowTech++

		if agedInLog(character) {
			aged++
		}
	}

	if lowTech == 0 {
		t.Fatalf("no seed in %d reached six terms from a tech level 9 world; "+
			"the aging tables are untested by generation", sample)
	}

	if aged == 0 {
		t.Errorf("%d characters served six terms from a tech level 9 world and "+
			"none ever failed an aging check", lowTech)
	}
}

// agedInLog reports whether the record holds a characteristic lost to an
// aging check, which is the consequence agingThrows writes.
func agedInLog(character *chargen.Character) bool {
	for _, event := range character.Events {
		if event.Consequence == nil {
			continue
		}

		if event.Consequence.Kind != chargen.ConsequenceCharacteristic {
			continue
		}

		if strings.HasPrefix(event.Consequence.Detail, "aged: ") {
			return true
		}
	}

	return false
}
