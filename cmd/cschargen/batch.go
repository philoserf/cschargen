package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/setting"
)

// batchCommand is `cschargen batch`: generate several characters and write
// them as JSONL, one record per line.
//
// It requires --auto, because a batch of twenty characters is twenty
// lifepaths of questions and nobody wants to answer four hundred of them.
// The seed of member i is the base seed plus i and is recorded in that
// member's own provenance, so any one of them can be regenerated alone with
// `new --seed`.
func batchCommand(args []string, out *os.File) error {
	flags := bindNewFlags()

	count := flags.set.Int("count", 0, "how many characters to generate")

	seed, err := flags.parse(args)
	if err != nil {
		return err
	}

	if *count < 1 {
		return usagef("batch needs --count N")
	}

	if !*flags.auto {
		return usagef("batch needs --auto: a batch is not an interview")
	}

	world, err := loadSetting(*flags.data)
	if err != nil {
		return err
	}

	encoded, err := generateBatch(*count, seed, world, flags)
	if err != nil {
		return err
	}

	if *flags.output == "" {
		_, err = out.Write(encoded)
		if err != nil {
			return fmt.Errorf("writing the batch: %w", err)
		}

		return nil
	}

	return writeFile(*flags.output, encoded, *flags.force)
}

// generateBatch is the loop. A batch is not a new generator: it is the one
// that already exists, run once per member.
func generateBatch(
	count int, base uint64, world *setting.Data, flags newFlags,
) ([]byte, error) {
	var lines []byte

	for i := range count {
		character, err := chargen.New(chargen.Options{
			Seed:          base + uint64(i),
			Decider:       chargen.Policy{},
			EngineVersion: version(),
			PolicyVersion: policyVersion,
			Setting:       world,
			Inputs:        flags.inputs(),
		}).Run()
		if err != nil {
			return nil, fmt.Errorf("generating character %d of %d: %w", i+1, count, err)
		}

		// JSONL: one record per line, so the file streams and a single
		// member can be pulled out with a line number.
		line, err := json.Marshal(character)
		if err != nil {
			return nil, fmt.Errorf("encoding character %d of %d: %w", i+1, count, err)
		}

		lines = append(lines, line...)
		lines = append(lines, '\n')
	}

	return lines, nil
}
