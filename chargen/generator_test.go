package chargen_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// The strings the tests repeat, named so that a typo in one case reads as a
// typo rather than as a disagreement between cases.
const (
	sampleData         = "sample"
	characteristicRoll = "3d6 drop lowest"
)

func options(seed uint64) chargen.Options {
	return chargen.Options{
		Seed:          seed,
		Decider:       chargen.Policy{},
		EngineVersion: "test",
		PolicyVersion: "test",
		SettingData:   chargen.SettingData{Name: sampleData, Sample: true},
		Inputs:        chargen.Inputs{Species: "human"},
	}
}

func generate(t *testing.T, opts chargen.Options) *chargen.Character {
	t.Helper()

	character, err := chargen.New(opts).Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	return character
}

// TestSameSeedSameCharacter is the determinism contract, checked on the
// characteristics rather than on the dice: the engine could consume the
// stream identically and still assign differently.
func TestSameSeedSameCharacter(t *testing.T) {
	t.Parallel()

	first := generate(t, options(7))
	second := generate(t, options(7))

	if first.Characteristics != second.Characteristics {
		t.Errorf("seed 7 gave %+v then %+v", first.Characteristics, second.Characteristics)
	}

	if len(first.Events) != len(second.Events) {
		t.Errorf("event counts %d and %d", len(first.Events), len(second.Events))
	}
}

func TestDifferentSeedsDivergeSomewhere(t *testing.T) {
	t.Parallel()

	same := 0

	for seed := range uint64(40) {
		if generate(t, options(seed)).Characteristics == generate(t, options(0)).Characteristics {
			same++
		}
	}

	if same == 40 {
		t.Error("forty seeds produced the same six characteristics")
	}
}

// TestCharacteristicsArePermutationOfTheRolls: the assignment moves scores
// around, it does not invent them. Six rolls in, the same six out.
func TestCharacteristicsArePermutationOfTheRolls(t *testing.T) {
	t.Parallel()

	for seed := range uint64(50) {
		character := generate(t, options(seed))

		rolled := make([]int, 0, len(chargen.CharacteristicOrder))

		for _, event := range character.Events {
			if event.Kind == chargen.EventThrow && event.Throw.Expr == characteristicRoll {
				rolled = append(rolled, event.Throw.Total)
			}
		}

		if len(rolled) != 6 {
			t.Fatalf("seed %d: %d characteristic rolls logged, want 6", seed, len(rolled))
		}

		assigned := make([]int, 0, len(chargen.CharacteristicOrder))
		for _, which := range chargen.CharacteristicOrder {
			assigned = append(assigned, character.Characteristics.Get(which))
		}

		slices.Sort(rolled)
		slices.Sort(assigned)

		if !slices.Equal(rolled, assigned) {
			t.Errorf("seed %d: rolled %v, assigned %v", seed, rolled, assigned)
		}
	}
}

// TestEveryScoreIsReachable: 3d6-drop-lowest runs 2 to 12, and a
// characteristic that never left a narrower band would mean the assignment
// or the roll was wrong.
func TestEveryScoreIsReachable(t *testing.T) {
	t.Parallel()

	seen := map[int]bool{}

	for seed := range uint64(400) {
		character := generate(t, options(seed))
		for _, which := range chargen.CharacteristicOrder {
			seen[character.Characteristics.Get(which)] = true
		}
	}

	for score := 2; score <= 12; score++ {
		if !seen[score] {
			t.Errorf("score %d never came up across 400 characters", score)
		}
	}
}

// TestPolicyAssignsInRollOrder is the row POLICY.md carries for this choice
// point: taking the first remaining option every time is assignment in the
// order the dice fell.
func TestPolicyAssignsInRollOrder(t *testing.T) {
	t.Parallel()

	character := generate(t, options(11))

	rolled := make([]int, 0, len(chargen.CharacteristicOrder))

	for _, event := range character.Events {
		if event.Kind == chargen.EventThrow && event.Throw.Expr == characteristicRoll {
			rolled = append(rolled, event.Throw.Total)
		}
	}

	if len(rolled) != len(chargen.CharacteristicOrder) {
		t.Fatalf("%d characteristic rolls logged, want %d", len(rolled), len(chargen.CharacteristicOrder))
	}

	for i, which := range chargen.CharacteristicOrder {
		if got := character.Characteristics.Get(which); got != rolled[i] {
			t.Errorf("%s = %d, want roll %d which was %d", which, got, i+1, rolled[i])
		}
	}
}

// TestTheLogNarratesTheStep: six rolls, six choices and six consequences,
// under one step header, each consequence naming the choice that caused it.
func TestTheLogNarratesTheStep(t *testing.T) {
	t.Parallel()

	character := generate(t, options(3))

	counts := map[chargen.EventKind]int{}
	for _, event := range character.Events {
		counts[event.Kind]++
	}

	want := map[chargen.EventKind]int{
		chargen.EventStep:        1,
		chargen.EventThrow:       6,
		chargen.EventChoice:      6,
		chargen.EventConsequence: 6,
	}

	for kind, n := range want {
		if counts[kind] != n {
			t.Errorf("%s events: %d, want %d", kind, counts[kind], n)
		}
	}

	for _, event := range character.Events {
		if event.Kind != chargen.EventConsequence {
			continue
		}

		cause := event.Consequence.Cause
		if cause < 1 || cause >= event.Seq {
			t.Errorf("event %d names cause %d, which is not an earlier event", event.Seq, cause)

			continue
		}

		if character.Events[cause-1].Kind != chargen.EventChoice {
			t.Errorf("event %d names event %d as its cause, which is a %s",
				event.Seq, cause, character.Events[cause-1].Kind)
		}
	}
}

// TestChoicesShrinkTheOptionList: each characteristic takes one of the
// scores still unassigned, so the sixth choice has exactly one option.
func TestChoicesShrinkTheOptionList(t *testing.T) {
	t.Parallel()

	character := generate(t, options(5))

	nth := 0

	for _, event := range character.Events {
		if event.Kind != chargen.EventChoice {
			continue
		}

		nth++

		want := 7 - nth
		if len(event.Choice.Options) != want {
			t.Errorf("choice %d offered %d options, want %d", nth, len(event.Choice.Options), want)
		}
	}

	if nth != 6 {
		t.Errorf("%d choices logged, want 6", nth)
	}
}

// watching records the Choice values the engine put to it, which the record
// deliberately does not keep: Nth and Of are engine-provided context for a
// front end, not part of the printed rule, so replay never compares them.
type watching struct {
	seen []chargen.Choice
}

func (*watching) Kind() chargen.DeciderKind { return chargen.DeciderPolicy }

func (w *watching) Choose(asked chargen.Choice) (int, error) {
	w.seen = append(w.seen, asked)

	return 0, nil
}

// TestChoicesArePlacedInTheirRun: a player answering the fifth of six
// identical questions cannot otherwise tell it from the first.
func TestChoicesArePlacedInTheirRun(t *testing.T) {
	t.Parallel()

	watcher := &watching{}

	opts := options(5)

	opts.Decider = watcher

	generate(t, opts)

	if len(watcher.seen) != 6 {
		t.Fatalf("%d choices put to the decider, want 6", len(watcher.seen))
	}

	for i, asked := range watcher.seen {
		if asked.Nth != i+1 {
			t.Errorf("choice %d has Nth %d", i+1, asked.Nth)
		}

		if asked.Of != 6 {
			t.Errorf("choice %d has Of %d, want 6", i+1, asked.Of)
		}

		if asked.Cite == "" {
			t.Errorf("choice %d carries no cite", i+1)
		}

		if asked.Prompt == "" {
			t.Errorf("choice %d carries no prompt", i+1)
		}
	}
}

func TestProvenanceIsStamped(t *testing.T) {
	t.Parallel()

	character := generate(t, options(13))

	got := character.Provenance
	if got.SchemaVersion != chargen.SchemaVersion {
		t.Errorf("schemaVersion = %d", got.SchemaVersion)
	}

	if got.Ruleset != chargen.Ruleset {
		t.Errorf("ruleset = %q", got.Ruleset)
	}

	if got.RNG.Seed != 13 {
		t.Errorf("seed = %d", got.RNG.Seed)
	}

	if got.RNG.Algorithm == "" {
		t.Error("the record does not name the generator")
	}

	if !got.SettingData.Sample {
		t.Error("a character built on the sample data is not flagged as such")
	}
}

// refusing is a decider that declines, which is an interactive session the
// player abandoned.
type refusing struct{}

func (refusing) Kind() chargen.DeciderKind { return chargen.DeciderPlayer }

var errRefused = errors.New("refused")

func (refusing) Choose(chargen.Choice) (int, error) { return 0, errRefused }

func TestARefusedChoiceEndsGeneration(t *testing.T) {
	t.Parallel()

	opts := options(1)

	opts.Decider = refusing{}

	_, err := chargen.New(opts).Run()
	if !errors.Is(err, errRefused) {
		t.Errorf("err = %v, want the decider's own error", err)
	}
}

// wrong is a decider that answers with an index the list cannot hold, which
// is a decider answering wrongly rather than declining.
type wrong struct{}

func (wrong) Kind() chargen.DeciderKind { return chargen.DeciderPolicy }

func (wrong) Choose(chargen.Choice) (int, error) { return 99, nil }

func TestAnOutOfRangeChoiceIsDistinctFromARefusal(t *testing.T) {
	t.Parallel()

	opts := options(1)

	opts.Decider = wrong{}

	_, err := chargen.New(opts).Run()
	if !errors.Is(err, chargen.ErrChoiceOutOfRange) {
		t.Errorf("err = %v, want ErrChoiceOutOfRange", err)
	}
}

// TestReplayReproducesTheCharacter closes the loop this milestone exists to
// establish: the record replays to itself.
func TestReplayReproducesTheCharacter(t *testing.T) {
	t.Parallel()

	original := generate(t, options(23))

	opts := options(23)

	opts.Decider = chargen.NewReplay(original.Events)

	replayed := generate(t, opts)

	if original.Characteristics != replayed.Characteristics {
		t.Errorf("replay gave %+v, record holds %+v", replayed.Characteristics, original.Characteristics)
	}

	if len(original.Events) != len(replayed.Events) {
		t.Errorf("replay logged %d events, record holds %d", len(replayed.Events), len(original.Events))
	}
}
