package chargen_test

import (
	"errors"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// Verify and Reproducible used to live in cmd/cschargen, where the only way
// to exercise them was to build a binary, write a record to a temp directory
// and tamper with the file. These are the tests that could not be written
// then: two records, constructed, differing in one thing.

// lifepathPair generates one character twice from the same seed, which is
// the shape Verify is handed -- a record and a re-run of it.
func lifepathPair(t *testing.T) (*chargen.Character, *chargen.Character) {
	t.Helper()

	opts := options(t, 11)

	opts.Inputs.TermLimit = 2

	original := generate(t, opts)

	again := options(t, 11)

	again.Inputs.TermLimit = 2
	again.Decider = chargen.NewReplay(original.Events)

	return original, generate(t, again)
}

func TestVerifyAcceptsATrueReplay(t *testing.T) {
	t.Parallel()

	original, replayed := lifepathPair(t)

	err := chargen.Verify(original, replayed)
	if err != nil {
		t.Errorf("a record did not verify against its own replay: %v", err)
	}
}

// TestVerifyNamesTheFirstEventThatDiffers walks the kinds of divergence the
// log can carry. Each case changes one event and expects the sequence
// number back, because the number is what a person needs to find it.
func TestVerifyNamesTheFirstEventThatDiffers(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		damage func(*chargen.Character)
	}{
		{"a throw's total", func(c *chargen.Character) {
			for i := range c.Events {
				if c.Events[i].Throw != nil {
					c.Events[i].Throw.Total++

					return
				}
			}
		}},
		{"a choice's index", func(c *chargen.Character) {
			for i := range c.Events {
				if c.Events[i].Choice != nil {
					c.Events[i].Choice.Chosen = 99

					return
				}
			}
		}},
		{"a consequence's detail", func(c *chargen.Character) {
			for i := range c.Events {
				if c.Events[i].Consequence != nil {
					c.Events[i].Consequence.Detail = "something else"

					return
				}
			}
		}},
		{"a shorter log", func(c *chargen.Character) {
			c.Events = c.Events[:len(c.Events)-1]
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			original, replayed := lifepathPair(t)

			tc.damage(replayed)

			err := chargen.Verify(original, replayed)
			if !errors.Is(err, chargen.ErrDiverged) {
				t.Errorf("Verify = %v, want ErrDiverged", err)
			}
		})
	}
}

// TestVerifyCatchesACharacterTheLogDoesNotExplain is why the state check
// exists. The logs are identical and the characters are not, which is the
// case that reported "identical" while three of Step 20's four fields had
// gone missing.
func TestVerifyCatchesACharacterTheLogDoesNotExplain(t *testing.T) {
	t.Parallel()

	original, replayed := lifepathPair(t)

	replayed.State.Finishing.Goals = "a goal the log never mentions"

	err := chargen.Verify(original, replayed)
	if !errors.Is(err, chargen.ErrDiverged) {
		t.Errorf("Verify = %v, want ErrDiverged", err)
	}
}

// TestReproducibleRefusesWhatThisBuildDidNotWrite covers the four things
// provenance is held to, one at a time.
func TestReproducibleRefusesWhatThisBuildDidNotWrite(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	const engine = "test"

	sound := func() *chargen.Character {
		return &chargen.Character{Provenance: chargen.Provenance{
			SchemaVersion: chargen.SchemaVersion,
			Ruleset:       chargen.Ruleset,
			EngineVersion: engine,
			SettingData:   chargen.SettingData{Hash: data.Hash},
		}}
	}

	err := chargen.Reproducible(sound(), engine, data)
	if err != nil {
		t.Fatalf("a sound record was refused: %v", err)
	}

	cases := map[string]func(*chargen.Character){
		"setting data":   func(c *chargen.Character) { c.Provenance.SettingData.Hash = "0000" },
		"schema version": func(c *chargen.Character) { c.Provenance.SchemaVersion = 99 },
		"ruleset":        func(c *chargen.Character) { c.Provenance.Ruleset = "another book" },
		"engine version": func(c *chargen.Character) { c.Provenance.EngineVersion = "someone else's" },
	}

	for name, damage := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			record := sound()

			damage(record)

			err := chargen.Reproducible(record, engine, data)
			if !errors.Is(err, chargen.ErrProvenance) {
				t.Errorf("Reproducible = %v, want ErrProvenance", err)
			}
		})
	}
}

// TestReproducibleIgnoresThePolicyVersion is the deliberate omission.
// Replay reapplies recorded choices and never consults the policy, so a
// record made under one policy replays under any other.
func TestReproducibleIgnoresThePolicyVersion(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	record := &chargen.Character{Provenance: chargen.Provenance{
		SchemaVersion: chargen.SchemaVersion,
		Ruleset:       chargen.Ruleset,
		EngineVersion: "test",
		PolicyVersion: "0.0.1-from-long-ago",
		SettingData:   chargen.SettingData{Hash: data.Hash},
	}}

	err := chargen.Reproducible(record, "test", data)
	if err != nil {
		t.Errorf("an old policy version was refused: %v", err)
	}
}
