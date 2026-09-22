package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/setting"
)

// replayCommand is `cschargen replay`: re-run the engine from a record's
// seed, inputs and recorded choices, and report the first place the
// re-run disagrees with what the record says happened.
//
// The stored event log is verification data, not input: every throw is
// recomputed. A record that replays clean is a record the current engine
// can still produce.
func replayCommand(args []string, out *os.File) error {
	flags := flagSet("replay")
	ignore := flags.Bool("ignore-provenance", false,
		"replay a record whose engine or schema version does not match this build")
	data := flags.String("data", "", "setting data file; omitted means the repository's invented sample")

	err := flags.Parse(args)
	if err != nil {
		return usagef("replay: %v", err)
	}

	if flags.NArg() != 1 {
		return usagef("replay takes exactly one character record")
	}

	original, err := readRecord(flags.Arg(0))
	if err != nil {
		return err
	}

	world, err := loadSetting(*data)
	if err != nil {
		return err
	}

	if !*ignore {
		err = checkProvenance(original, world)
		if err != nil {
			return err
		}
	}

	replayed, err := chargen.New(chargen.Options{
		Seed:          original.Provenance.RNG.Seed,
		Decider:       chargen.NewReplay(original.Events),
		EngineVersion: version(),
		PolicyVersion: chargen.PolicyVersion,
		Setting:       world,
		Inputs:        original.Provenance.Inputs,
	}).Run()
	if err != nil {
		return fmt.Errorf("replaying: %w", err)
	}

	err = compare(original, replayed)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "replayed %d events from seed %d: identical\n",
		len(original.Events), original.Provenance.RNG.Seed)
	if err != nil {
		return fmt.Errorf("writing the result: %w", err)
	}

	return nil
}

// checkProvenance refuses a record this build cannot claim to reproduce.
//
// policy_version is deliberately not among the checks: replay reapplies
// recorded choices and never consults the policy, so a record made under
// one policy replays under any other.
func checkProvenance(original *chargen.Character, world *setting.Data) error {
	if world.Hash != original.Provenance.SettingData.Hash {
		return fmt.Errorf(
			"%w: the record was generated against setting data hashed %.12s, this is %.12s",
			errProvenance, original.Provenance.SettingData.Hash, world.Hash)
	}

	if original.Provenance.SchemaVersion != chargen.SchemaVersion {
		return fmt.Errorf("%w: record is schema %d, this build writes %d",
			errProvenance, original.Provenance.SchemaVersion, chargen.SchemaVersion)
	}

	if original.Provenance.Ruleset != chargen.Ruleset {
		return fmt.Errorf("%w: record was read against %q",
			errProvenance, original.Provenance.Ruleset)
	}

	if original.Provenance.EngineVersion != version() {
		return fmt.Errorf("%w: record was written by engine %s, this is %s",
			errProvenance, original.Provenance.EngineVersion, version())
	}

	return nil
}

// compare walks the two logs together and names the first event that
// differs. Comparing the logs rather than the characters is deliberate: a
// divergence shows up in the log several events before it reaches the
// character, and the sequence number is what a person needs to find it.
func compare(original, replayed *chargen.Character) error {
	for i, want := range original.Events {
		if i >= len(replayed.Events) {
			return fmt.Errorf("%w: the replay ended after %d events; the record holds %d",
				errDiverged, len(replayed.Events), len(original.Events))
		}

		got := replayed.Events[i]

		difference := eventDifference(want, got)
		if difference != "" {
			return fmt.Errorf("%w at event %d: %s", errDiverged, want.Seq, difference)
		}
	}

	if len(replayed.Events) > len(original.Events) {
		return fmt.Errorf("%w: the replay ran to %d events; the record holds %d",
			errDiverged, len(replayed.Events), len(original.Events))
	}

	return compareState(original, replayed)
}

// compareState holds the whole character, after the log. It ran on
// Characteristics alone for a long while -- the only comparable struct in
// State, so `==` reached exactly as far as Go let it and the line was never
// revisited. A replay that lost three of Step 20's four fields reported
// "identical" (#109).
//
// The comparison goes through JSON rather than reflect.DeepEqual because the
// two sides are not built the same way: original was decoded from a file and
// replayed came fresh from the engine, so an empty Skills is a non-nil empty
// slice on one side and nil on the other. Marshalling both is also comparing
// them in the representation the record is stored in, which is the one that
// has to match.
func compareState(original, replayed *chargen.Character) error {
	before, err := json.Marshal(original.State)
	if err != nil {
		return fmt.Errorf("encoding the record's character: %w", err)
	}

	after, err := json.Marshal(replayed.State)
	if err != nil {
		return fmt.Errorf("encoding the replayed character: %w", err)
	}

	if !bytes.Equal(before, after) {
		return fmt.Errorf("%w: the logs agree but the characters do not", errDiverged)
	}

	return nil
}

// eventDifference describes how two events differ, or returns empty when
// they do not. Each kind is compared on the fields that identify it: a
// throw by its expression and total, a choice by its point and the index
// taken, a consequence by what it says happened.
func eventDifference(want, got chargen.Event) string {
	if want.Kind != got.Kind {
		return fmt.Sprintf("record has a %s, the replay a %s", want.Kind, got.Kind)
	}

	switch want.Kind {
	case chargen.EventStep:
		return stepDifference(want, got)
	case chargen.EventThrow:
		return throwDifference(want, got)
	case chargen.EventChoice:
		return choiceDifference(want, got)
	case chargen.EventConsequence:
		if want.Consequence.Detail != got.Consequence.Detail {
			return fmt.Sprintf("consequence %q against %q",
				want.Consequence.Detail, got.Consequence.Detail)
		}
	}

	return ""
}

func stepDifference(want, got chargen.Event) string {
	if want.Step.Name != got.Step.Name {
		return fmt.Sprintf("step %q against %q", want.Step.Name, got.Step.Name)
	}

	return ""
}

func throwDifference(want, got chargen.Event) string {
	if want.Throw.Total != got.Throw.Total || want.Throw.Expr != got.Throw.Expr {
		return fmt.Sprintf("throw %s=%d against %s=%d",
			want.Throw.Expr, want.Throw.Total, got.Throw.Expr, got.Throw.Total)
	}

	return ""
}

func choiceDifference(want, got chargen.Event) string {
	if want.Choice.Chosen != got.Choice.Chosen || want.Choice.Point != got.Choice.Point {
		return fmt.Sprintf("choice %s#%d against %s#%d",
			want.Choice.Point, want.Choice.Chosen, got.Choice.Point, got.Choice.Chosen)
	}

	return ""
}
