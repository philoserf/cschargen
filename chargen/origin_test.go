package chargen_test

import (
	"slices"
	"strconv"
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

		// Nobody asked for one, so the inputs name none. Inputs carry
		// what was requested, not what happened -- the history is where
		// the world a character was born on is recorded, and the only
		// place, because it is also where a reassignment goes.
		if asked := character.Provenance.Inputs.Homeworld; asked != "" {
			t.Errorf("seed %d: the inputs name %q as a homeworld that was never asked for",
				seed, asked)
		}
	}
}

// TestAHomeworldAskedForIsTheHomeworld is p. 40's other half: "You may
// either roll percentile dice (d100) to determine your homeworld randomly
// or simply choose a world that fits the character concept you have in
// mind." A world named in the inputs is chosen, and no throw decides it.
func TestAHomeworldAskedForIsTheHomeworld(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	// Every world in the sample, including the ones on the choose-only
	// subsector that no d100 can reach.
	for _, sub := range data.Subsectors {
		for _, world := range sub.Worlds {
			t.Run(world.Name, func(t *testing.T) {
				t.Parallel()

				opts := options(t, 3)

				opts.Inputs.TermLimit = 1
				opts.Inputs.Homeworld = world.Name

				character := generate(t, opts)

				birth := character.State.Homeworlds[0]
				if birth.World != world.Name {
					t.Errorf("asked to be born on %s, born on %s", world.Name, birth.World)
				}

				if birth.Subsector != sub.Name {
					t.Errorf("%s is in %s, recorded as %s", world.Name, sub.Name, birth.Subsector)
				}
			})
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
//
// A record where something took a skill away is excluded rather than held
// to it. Orbital Construction's sixth mishap really does read "you lose all
// levels of Suit (Vacc Suit)", and a Vacc Suit from the homeworld is still
// a Vacc Suit.
func TestBackgroundSkillsComeFromTheHomeworld(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	for seed := range uint64(sample) {
		character := lifepath(t, seed, 1)
		if lostASkill(character) {
			continue
		}

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

// lostASkill reports whether any result in a record took a skill away.
func lostASkill(character *chargen.Character) bool {
	for _, event := range character.Events {
		if event.Kind != chargen.EventConsequence {
			continue
		}

		if strings.HasPrefix(event.Consequence.Detail, "lose every level of") {
			return true
		}
	}

	return false
}

func heldInStash(character *chargen.Character, item string) bool {
	return slices.ContainsFunc(character.State.Stash, func(held chargen.Possession) bool {
		return held.Item == item
	})
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
//
// It forces the Colonist career, six of whose eleven mishaps reassign the
// homeworld (p. 174). It did not, until the career list grew past the point
// where the auto policy happened to pick Colonist -- and the test then
// failed rather than passing vacuously, which is what the closing check is
// for.
func TestReassignmentChangesTheTechLevelAndNothingElse(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)
	moved := 0

	for seed := range uint64(sample) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = 8
		opts.Inputs.Career = careerColonist

		character := generate(t, opts)

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

// picking is a decider that answers the two origin choice points by the
// option's own text and takes the first option everywhere else, so a test
// can name the subsector or world it wants without counting how many
// choices the lifepath made before it.
type picking struct {
	want map[string]string
	seen map[string][]string
}

func (p *picking) Kind() chargen.DeciderKind { return chargen.DeciderPlayer }

func (p *picking) Choose(ask chargen.Choice) (int, error) {
	if p.seen == nil {
		p.seen = map[string][]string{}
	}

	p.seen[ask.Point] = ask.Options

	label, wanted := p.want[ask.Point]
	if !wanted {
		return 0, nil
	}

	for i, option := range ask.Options {
		if option == label {
			return i, nil
		}
	}

	return 0, nil
}

// TestAnInteractiveRunIsOfferedTheOriginItIsPromised is p. 39 and p. 40
// both: the Referee may choose the subsector, and a player may "simply
// choose a world that fits the character concept". Neither was offered
// before, so a setting file's worlds could only be reached by dice.
func TestAnInteractiveRunIsOfferedTheOriginItIsPromised(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	// The choose-only subsector: no world in it carries a roll, so nothing
	// but a choice can reach it, and before this nothing could.
	var target setting.Subsector

	for _, sub := range data.Subsectors {
		if sub.OriginRoll == 0 {
			target = sub

			break
		}
	}

	if target.Name == "" {
		t.Skip("the sample has no choose-only subsector")
	}

	world := target.Worlds[len(target.Worlds)-1]

	decider := &picking{want: map[string]string{
		"subsector_source": target.Name,
		"homeworld_source": world.Name,
	}}

	opts := options(t, 4)

	opts.Decider = decider
	opts.Inputs.Interactive = true

	character := generate(t, opts)

	birth := character.State.Homeworlds[0]
	if birth.World != world.Name || birth.Subsector != target.Name {
		t.Errorf("chose %s in %s, born on %s in %s",
			world.Name, target.Name, birth.World, birth.Subsector)
	}

	// Both prompts offer the throw first, so that the policy -- which takes
	// the first option everywhere -- keeps rolling.
	for _, point := range []string{"subsector_source", "homeworld_source"} {
		offered := decider.seen[point]
		if len(offered) == 0 {
			t.Fatalf("%s was never offered", point)
		}

		if offered[0] != "roll for it" {
			t.Errorf("%s offers %q first, not the throw", point, offered[0])
		}
	}
}

// TestAnInteractiveRunMayStillRoll: the choice the pages offer includes
// rolling, and taking it has to land where the dice land rather than
// falling through to some first entry.
func TestAnInteractiveRunMayStillRoll(t *testing.T) {
	t.Parallel()

	rolled := options(t, 9)

	rolled.Inputs.TermLimit = 1

	fromPolicy := generate(t, rolled)

	chose := options(t, 9)

	chose.Inputs.TermLimit = 1
	chose.Inputs.Interactive = true
	chose.Decider = &picking{want: map[string]string{
		"subsector_source": "roll for it",
		"homeworld_source": "roll for it",
	}}

	fromPlayer := generate(t, chose)

	if fromPolicy.State.Homeworlds[0].World != fromPlayer.State.Homeworlds[0].World {
		t.Errorf("the policy rolled %s and a player choosing the throw got %s",
			fromPolicy.State.Homeworlds[0].World, fromPlayer.State.Homeworlds[0].World)
	}
}

// TestAnOriginTheDataDoesNotHaveIsRefused. The command checks both names
// before generating, so this is the engine's own guard against being driven
// directly with a name the setting never had.
func TestAnOriginTheDataDoesNotHaveIsRefused(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name   string
		inputs func(*chargen.Inputs)
	}{
		{"homeworld", func(in *chargen.Inputs) { in.Homeworld = "Nowhere At All" }},
		{"subsector", func(in *chargen.Inputs) { in.Subsector = "No Such Reach" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			opts := options(t, 2)
			test.inputs(&opts.Inputs)

			_, err := chargen.New(opts).Run()
			if err == nil {
				t.Fatal("a name the setting data does not have was accepted")
			}
		})
	}
}

// TestASubsectorAskedForRollsInsideIt is the half of p. 39 that is not a
// whole homeworld: the Referee names the subsector and the d100 still
// decides the world.
func TestASubsectorAskedForRollsInsideIt(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	for _, sub := range data.Subsectors {
		if sub.OriginRoll == 0 {
			continue
		}

		t.Run(sub.Name, func(t *testing.T) {
			t.Parallel()

			opts := options(t, 6)

			opts.Inputs.TermLimit = 1
			opts.Inputs.Subsector = sub.Name

			character := generate(t, opts)

			if got := character.State.Homeworlds[0].Subsector; got != sub.Name {
				t.Errorf("asked for %s, born in %s", sub.Name, got)
			}
		})
	}
}

// refusingAt is a player who stops answering at one named choice point,
// which is the abandoned session the Decider interface names.
type refusingAt struct{ point string }

func (refusingAt) Kind() chargen.DeciderKind { return chargen.DeciderPlayer }

func (r refusingAt) Choose(ask chargen.Choice) (int, error) {
	if ask.Point == r.point {
		return 0, chargen.ErrPlayerGone
	}

	return 0, nil
}

// TestAbandoningTheOriginEndsGeneration. Both new prompts come before any
// career, so a session abandoned at either one has to stop rather than fall
// back on a default homeworld -- which would put the character somewhere
// nobody chose.
func TestAbandoningTheOriginEndsGeneration(t *testing.T) {
	t.Parallel()

	for _, point := range []string{"subsector_source", "homeworld_source"} {
		t.Run(point, func(t *testing.T) {
			t.Parallel()

			opts := options(t, 8)

			opts.Inputs.Interactive = true
			opts.Decider = refusingAt{point: point}

			_, err := chargen.New(opts).Run()
			if err == nil {
				t.Fatalf("generation continued after %s was abandoned", point)
			}
		})
	}
}

// TestAThrowThatMissesIsThrownAgain is the chart of p. 39 against a file
// that does not fill it.
//
// The book's chart has an entry for all six results; a setting file may
// have fewer subsectors than that. A throw landing in the gap used to
// become a choice, and the auto policy answers a choice with its first
// option -- which put most of a batch on whichever subsector the file
// happened to list first, and the skew was invisible one record at a time.
func TestAThrowThatMissesIsThrownAgain(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	rollable := map[string]bool{}

	for _, sub := range data.Subsectors {
		if sub.OriginRoll != 0 {
			rollable[sub.Name] = true
		}
	}

	if len(rollable) < 2 || len(rollable) == len(data.Subsectors) {
		t.Skip("the sample fills the chart, so there is no gap to land in")
	}

	seen := map[string]int{}

	for seed := range uint64(sample) {
		born := lifepath(t, seed, 1).State.Homeworlds[0].Subsector
		seen[born]++

		if !rollable[born] {
			t.Errorf("seed %d: a throw put the character in %s, which claims no result",
				seed, born)
		}
	}

	// Every rollable subsector claims one result here, so none of them
	// should be starved and none should take most of the sample. The
	// bounds are wide: this is checking that the throw is a throw, not
	// that sixty seeds are uniform.
	for name := range rollable {
		if seen[name] == 0 {
			t.Errorf("%s claims a result and no character in %d was born there", name, sample)
		}

		if share := float64(seen[name]) / float64(sample); share > 0.6 {
			t.Errorf("%s took %.0f%% of the sample, which is the old first-option skew",
				name, share*100)
		}
	}
}

// TestTheTechLevelAskedForGatesTheAging is pp. 122-123: the aging bands are
// keyed by the homeworld's tech level, and `--tech-level` is how a referee
// generates a character as if from a world their setting file does not have.
//
// It was recorded in provenance and read by nothing, so a record named a
// tech level that never gated a throw.
func TestTheTechLevelAskedForGatesTheAging(t *testing.T) {
	t.Parallel()

	// A low tech level starts the throws at 18 rather than 33 or later, so
	// a character of the same age shows it in the apparent-age lookup
	// (p. 125) without needing the aging throws themselves read out.
	low := options(t, 21)

	low.Inputs.TermLimit = 4
	low.Inputs.TechLevel = 9

	high := options(t, 21)

	high.Inputs.TermLimit = 4

	aged := generate(t, low)
	ordinary := generate(t, high)

	if aged.State.ApparentAge == ordinary.State.ApparentAge {
		t.Errorf("tech level 9 and the homeworld's own gave the same apparent age %v",
			aged.State.ApparentAge)
	}
}

// TestTheMaximumTermsAskedForBinds is p. 42's homeworld ceiling, which
// CLAUDE.md describes as one of the term limit's two and the lower winning
// (p. 125).
func TestTheMaximumTermsAskedForBinds(t *testing.T) {
	t.Parallel()

	for _, want := range []int{1, 2, 3} {
		t.Run(strconv.Itoa(want), func(t *testing.T) {
			t.Parallel()

			opts := options(t, 21)

			// Higher than the cap, so the cap is what binds.
			opts.Inputs.TermLimit = 8
			opts.Inputs.MaxTerms = want

			if got := len(generate(t, opts).State.Terms); got != want {
				t.Errorf("a homeworld maximum of %d terms let the character serve %d", want, got)
			}
		})
	}
}

// TestAReassignmentTakesTheNewWorldsTechLevel. `--tech-level` says where a
// character was born, the way `--homeworld` does. A Colonist mishap moves
// them to a world with its own tech level, and ERRATA E-9 is that the
// reassignment changes exactly that -- so the flag must not follow them.
func TestAReassignmentTakesTheNewWorldsTechLevel(t *testing.T) {
	t.Parallel()

	opts := options(t, 21)

	opts.Inputs.TermLimit = 8
	opts.Inputs.TechLevel = 9
	opts.Inputs.Career = careerColonist

	character := generate(t, opts)
	if len(character.State.Homeworlds) < 2 {
		t.Skip("this seed was never reassigned a homeworld")
	}

	moved := character.State.Homeworlds[1]
	if moved.TechLevel == 9 && moved.World != character.State.Homeworlds[0].World {
		t.Error("the tech level asked for followed the character to a world that has its own")
	}
}
