package main

import (
	"fmt"
	"os"

	"github.com/philoserf/cschargen/chargen"
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
		err = chargen.Reproducible(original, version(), world)
		if err != nil {
			// The message is the whole value here -- it names which
			// version or hash disagreed -- so this wraps without adding
			// to it, the way loadSetting does.
			return fmt.Errorf("%w", err)
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

	err = chargen.Verify(original, replayed)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_, err = fmt.Fprintf(out, "replayed %d events from seed %d: identical\n",
		len(original.Events), original.Provenance.RNG.Seed)
	if err != nil {
		return fmt.Errorf("writing the result: %w", err)
	}

	return nil
}
