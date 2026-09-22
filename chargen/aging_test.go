package chargen_test

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// TestACharacterActuallyAges is the reach check against the sample data:
// aging begins at term 6 on a tech level 9 world,
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

// TestTheEnlistmentModifierIsApplied. Twelve careers take -2 to enlistment
// from an apparent age over 40, and until aging landed the engine recorded
// that as unimplemented instead of applying it. The record has to show it
// being applied, and nowhere show it being recorded.
func TestTheEnlistmentModifierIsApplied(t *testing.T) {
	t.Parallel()

	recorded := 0

	for seed := range uint64(sample) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = 40

		character := generate(t, opts)

		for _, event := range character.Events {
			if event.Consequence == nil {
				continue
			}

			if strings.Contains(event.Consequence.Detail, "apparent age, which arrives") {
				recorded++
			}
		}
	}

	if recorded > 0 {
		t.Errorf("%d records still say the apparent age modifier is unimplemented", recorded)
	}
}

// TestApparentAgeReachesTheSheet. It is derived from the homeworld's tech
// level, which the record does not otherwise carry, so the generator stamps
// it -- and a character old enough for the chart to say something the age
// does not is the case worth checking.
//
// The homeworld it reads is the last one, not the first. ERRATA E-9: a
// homeworld reassigned during a career changes the tech level that gates
// the aging throws, so a character born on a tech level 11 world and
// deported to a tech level 8 one looks their age again.
func TestApparentAgeReachesTheSheet(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)
	checked := 0

	for seed := range uint64(sample) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = 40

		character := generate(t, opts)
		homeworlds := character.State.Homeworlds

		born, _, found := data.World(homeworlds[len(homeworlds)-1].World)
		if !found || born.TechLevel < 10 || character.State.Age < 30 {
			continue
		}

		checked++

		band := character.State.ApparentAge
		if band.From >= character.State.Age {
			t.Errorf("seed %d: age %d on a tech level %d world, apparent age %s -- "+
				"the chart should have made them look younger",
				seed, character.State.Age, born.TechLevel, band)
		}
	}

	if checked == 0 {
		t.Fatalf("no seed in %d reached age 30 from a tech level 10 or higher world", sample)
	}
}
