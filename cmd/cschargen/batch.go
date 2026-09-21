package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/setting"
)

// splitJSONL is the lines of a JSONL document, without the empty tail the
// trailing newline leaves.
func splitJSONL(document []byte) [][]byte {
	var records [][]byte

	for line := range bytes.SplitSeq(document, []byte("\n")) {
		if len(bytes.TrimSpace(line)) > 0 {
			records = append(records, line)
		}
	}

	return records
}

// itoa keeps the width arithmetic in writeRecordsInto readable.
func itoa(n int) string { return strconv.Itoa(n) }

// batchCommand is `cschargen batch`: generate several characters and write
// them as JSONL, one record per line.
//
// It requires --auto, because a batch of twenty characters is twenty
// lifepaths of questions and nobody wants to answer four hundred of them.
// The seed of member i is the base seed plus i and is recorded in that
// member's own provenance, which is what `replay` reproduces it from.
//
// It is not what `new --seed` reproduces it from. Since policy_version
// 0.3.0 a batch draws its career and assignment where `new --auto` takes
// the first of each, so one seed makes two different characters under the
// two commands, deliberately. `-o dir` keeps each member as its own file
// for exactly this reason.
// checkBatchFlags holds the flags a batch needs that a single character
// does not, so that batchCommand reads as the sequence it is.
func checkBatchFlags(flags newFlags, count int) error {
	err := rejectFinishingFlags(flags)
	if err != nil {
		return err
	}

	if count < 1 {
		return usagef("batch needs --count N")
	}

	if !*flags.auto {
		return usagef("batch needs --auto: a batch is not an interview")
	}

	return nil
}

func batchCommand(args []string, out *os.File) error {
	flags := bindNewFlags("batch")

	count := flags.set.Int("count", 0, "how many characters to generate")

	seed, err := flags.parse(args)
	if err != nil {
		return err
	}

	err = checkBatchFlags(flags, *count)
	if err != nil {
		return err
	}

	world, err := loadSetting(*flags.data)
	if err != nil {
		return err
	}

	err = checkOrigin(world, *flags.subsector, *flags.homeworld)
	if err != nil {
		return err
	}

	encoded, err := generateFrom(*count, seed, world, flags)
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

	if isDirectory(*flags.output) {
		return writeRecordsInto(*flags.output, encoded, *flags.force)
	}

	return writeFile(*flags.output, encoded, *flags.force)
}

// isDirectory reports whether the output path names a directory that
// already exists. The PRD's CLI sketch has always read
// `-o dir|file.jsonl`, and an existing directory is the least surprising
// way to ask for the first: nothing is created behind the caller's back,
// and a path ending in `.jsonl` can never be mistaken for one.
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// writeRecordsInto splits a batch into one file per record, numbered in the
// order they were generated, so that every other command works on them
// unchanged: `render dir/npc-03.json` needs no new code at all.
//
// The width of the number is the width of the count, so a batch of a
// hundred sorts correctly in a listing and in a glob.
func writeRecordsInto(dir string, batch []byte, force bool) error {
	records := splitJSONL(batch)
	width := len(itoa(len(records)))

	for i, record := range records {
		name := fmt.Sprintf("npc-%0*d.json", width, i+1)

		err := writeFile(filepath.Join(dir, name), append(record, '\n'), force)
		if err != nil {
			return err
		}
	}

	return nil
}

// generateFrom reads the name list, if there is one, and runs the batch.
// It exists so that batchCommand reads as the sequence of things it does
// rather than as the sum of their error handling.
func generateFrom(
	count int, base uint64, world *setting.Data, flags newFlags,
) ([]byte, error) {
	names, err := readNames(*flags.names)
	if err != nil {
		return nil, err
	}

	return generateBatch(count, base, world, flags, names)
}

// rejectFinishingFlags refuses the four that set one character's Step 20
// fields. They are shared with `new` because the two commands take the same
// inputs otherwise, and on a batch they would give twenty people the same
// name and the same face.
//
// Refusing costs the caller one line; applying them quietly costs them
// twenty wrong sheets and the time to work out why.
func rejectFinishingFlags(flags newFlags) error {
	for _, name := range finishingFlags {
		if !isSet(flags.set, name) {
			continue
		}

		return usagef("--%s sets one character's Step 20 fields (pp. 129-130); "+
			"batch generates many. Drop it, or use `new`", name)
	}

	return nil
}

// generateBatch is the loop. A batch is not a new generator: it is the one
// that already exists, run once per member.
func generateBatch(
	count int, base uint64, world *setting.Data, flags newFlags, names []string,
) ([]byte, error) {
	var (
		lines []byte
		made  []*chargen.Character
	)

	for i := range count {
		seed := base + uint64(i)

		inputs := flags.inputs()

		inputs.Name = nameFor(names, seed)

		character, err := chargen.New(chargen.Options{
			Seed:          seed,
			Decider:       chargen.NewBatchPolicy(seed),
			EngineVersion: version(),
			PolicyVersion: policyVersion,
			Setting:       world,
			Inputs:        inputs,
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
		made = append(made, character)
	}

	warnHowManyCareersChanged(os.Stderr, made, *flags.forceCar)

	return lines, nil
}
