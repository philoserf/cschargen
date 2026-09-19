package chargen_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/setting"
)

// TestEveryCharacterHasAHomeworld: Steps 3 and 4 are not optional, and
// everything downstream -- background skills, language, tech level, term
// cap -- hangs off the result.
func TestEveryCharacterHasAHomeworld(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 4)

		if len(character.State.Homeworlds) == 0 {
			t.Fatalf("seed %d: no homeworld", seed)
		}

		birth := character.State.Homeworlds[0]
		if birth.FromTerm != 0 {
			t.Errorf("seed %d: the birth world is dated to term %d", seed, birth.FromTerm)
		}

		_, _, found := data.World(birth.World)
		if !found {
			t.Errorf("seed %d: born on %q, which is not in the setting data", seed, birth.World)
		}

		if character.Provenance.Inputs.Homeworld != birth.World {
			t.Errorf("seed %d: the record names %q as the homeworld but the history begins on %q",
				seed, character.Provenance.Inputs.Homeworld, birth.World)
		}
	}
}

// TestEveryCharacterSpeaksSomething is p. 41: every homeworld lists a
// primary language, and a character has the first of them unless a decider
// chose otherwise.
func TestEveryCharacterSpeaksSomething(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 2)

		if character.State.Language == "" {
			t.Fatalf("seed %d: no primary language", seed)
		}

		world, _, found := data.World(character.State.Homeworlds[0].World)
		if !found {
			t.Fatalf("seed %d: unknown homeworld", seed)
		}

		if !slices.Contains(world.PrimaryLanguages, character.State.Language) {
			t.Errorf("seed %d: speaks %q, which %s does not list",
				seed, character.State.Language, world.Name)
		}
	}
}

// TestEveryCharacterHasElectronicsZero is p. 41's one universal line: "All
// characters, regardless of origin, will also receive Electronics 0 as a
// Background Skill."
func TestEveryCharacterHasElectronicsZero(t *testing.T) {
	t.Parallel()

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 2)

		if !character.State.Has("Electronics") {
			t.Errorf("seed %d: no Electronics", seed)
		}
	}
}

// TestBackgroundSkillsComeFromTheHomeworld: a character holds every skill
// their birth world grants, at level 1 or better.
func TestBackgroundSkillsComeFromTheHomeworld(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 1)

		world, _, found := data.World(character.State.Homeworlds[0].World)
		if !found {
			t.Fatalf("seed %d: unknown homeworld", seed)
		}

		for _, requirement := range world.BackgroundSkills {
			satisfied := false

			for _, alternative := range requirement.OneOf {
				if alternative.Item != "" {
					satisfied = satisfied || heldInStash(character, alternative.Item)

					continue
				}

				if character.State.SkillLevel(alternative.Skill) >= 1 {
					satisfied = true
				}
			}

			if !satisfied {
				t.Errorf("seed %d: born on %s and holds none of %v",
					seed, world.Name, describeOneOf(requirement))
			}
		}
	}
}

func heldInStash(character *chargen.Character, item string) bool {
	return slices.Contains(character.State.Stash, item)
}

func describeOneOf(requirement setting.Requirement) []string {
	names := make([]string, 0, len(requirement.OneOf))
	for _, alternative := range requirement.OneOf {
		if alternative.Item != "" {
			names = append(names, alternative.Item)

			continue
		}

		names = append(names, alternative.Skill)
	}

	return names
}

// TestALanguageSpecialtyIsNeverThePrimary is p. 41: "the Language specialty
// chosen for the skill must be different from the character's Primary
// Language. This represents learning an additional language beyond one's
// native tongue."
func TestALanguageSpecialtyIsNeverThePrimary(t *testing.T) {
	t.Parallel()

	checked := 0

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 1)

		for _, held := range character.State.Skills {
			if held.Name != "Language" {
				continue
			}

			checked++

			if held.Specialty == character.State.Language {
				t.Errorf("seed %d: Language (%s) is also the primary language",
					seed, held.Specialty)
			}
		}
	}

	if checked == 0 {
		t.Skip("no character in the sample was granted the Language skill")
	}
}

// TestTheHomeworldCapsTheTerms is p. 125: "If the character has reached the
// maximum number of terms allowed by their homeworld, the character must
// end character generation at this time." The policy limit and the
// homeworld cap are both ceilings and the lower wins.
func TestTheHomeworldCapsTheTerms(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)
	capped := 0

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 40)

		world, _, found := data.World(character.State.Homeworlds[0].World)
		if !found {
			t.Fatalf("seed %d: unknown homeworld", seed)
		}

		served := len(character.State.Terms)
		if served > world.MaximumTerms {
			t.Errorf("seed %d: served %d terms, but %s allows %d",
				seed, served, world.Name, world.MaximumTerms)
		}

		if served == world.MaximumTerms {
			capped++
		}
	}

	if capped == 0 {
		t.Error("no character reached their homeworld's cap, so the cap is untested")
	}
}

// TestReassignmentChangesTheTechLevelAndNothingElse is ERRATA E-9. A
// character deported in adulthood ages by the medicine of where they now
// live, but did not acquire a second childhood.
func TestReassignmentChangesTheTechLevelAndNothingElse(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)
	moved := 0

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 8)

		if len(character.State.Homeworlds) < 2 {
			continue
		}

		moved++

		birth, _, found := data.World(character.State.Homeworlds[0].World)
		if !found {
			t.Fatalf("seed %d: unknown birth world", seed)
		}

		if len(character.State.Terms) > birth.MaximumTerms {
			t.Errorf("seed %d: a reassignment extended the term cap past the birth world's %d",
				seed, birth.MaximumTerms)
		}

		checkMoves(t, seed, data, character)

		if !slices.Contains(character.Provenance.Deviations, "E-9") {
			t.Errorf("seed %d: the homeworld changed without stamping E-9", seed)
		}
	}

	if moved == 0 {
		t.Error("no character was ever reassigned a homeworld, so the reading is untested")
	}
}

// checkMoves holds what a later homeworld records: a world the data knows,
// at the tech level the data gives it, and a reason for the move.
func checkMoves(t *testing.T, seed uint64, data *setting.Data, character *chargen.Character) {
	t.Helper()

	for _, home := range character.State.Homeworlds[1:] {
		later, _, found := data.World(home.World)
		if !found {
			t.Errorf("seed %d: moved to %q, which is not in the data", seed, home.World)

			continue
		}

		if home.TechLevel != later.TechLevel {
			t.Errorf("seed %d: the history records %s at TL %d; the data says %d",
				seed, home.World, home.TechLevel, later.TechLevel)
		}

		if home.Reason == "" {
			t.Errorf("seed %d: a move to %s with no reason recorded", seed, home.World)
		}
	}
}

// TestBeingBornOffworldIsRecordedAsUnimplemented: a character born outside
// the campaign's sector owes it four terms (p. 43), which the engine cannot
// check until it records where a term was served.
func TestBeingBornOffworldIsRecordedAsUnimplemented(t *testing.T) {
	t.Parallel()

	found := false

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 2)

		if !bornOffworld(sampleSetting(t), character) {
			continue
		}

		found = true

		if !recordsTheOffworldRule(character) {
			t.Errorf("seed %d: born offworld with nothing in the record about p. 43", seed)
		}
	}

	if !found {
		t.Skip("no character in the sample was born offworld")
	}
}

func bornOffworld(data *setting.Data, character *chargen.Character) bool {
	for _, sub := range data.Subsectors {
		if sub.Offworld && sub.Name == character.State.Homeworlds[0].Subsector {
			return true
		}
	}

	return false
}

func recordsTheOffworldRule(character *chargen.Character) bool {
	for _, event := range character.Events {
		if event.Kind == chargen.EventConsequence &&
			event.Consequence.Kind == chargen.ConsequenceUnimplemented &&
			strings.Contains(event.Consequence.Detail, "minimum of four terms") {
			return true
		}
	}

	return false
}

// TestGeneratingWithoutSettingDataIsRefused: Steps 3 and 4 have nothing to
// read without it, and a character with no homeworld has no background
// skills, no language, no tech level and no term cap.
func TestGeneratingWithoutSettingDataIsRefused(t *testing.T) {
	t.Parallel()

	opts := options(t, 1)

	opts.Setting = nil

	_, err := chargen.New(opts).Run()
	if err == nil {
		t.Fatal("a character generated with no setting data")
	}

	if !strings.Contains(err.Error(), "setting data") {
		t.Errorf("the error does not name the problem: %v", err)
	}
}
