package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/render"
)

// renderCommand is `cschargen render`: turn a record into the Markdown
// sheet, or with --history into the lifepath transcript.
func renderCommand(args []string, out *os.File) error {
	flags := flagSet("render")
	history := flags.Bool("history", false, "render the lifepath transcript instead of the sheet")
	roster := flags.Bool("roster", false, "one line per character, for casting from a batch")

	err := flags.Parse(args)
	if err != nil {
		return usagef("render: %v", err)
	}

	if flags.NArg() != 1 {
		return usagef("render takes exactly one character record")
	}

	characters, err := readRecords(flags.Arg(0))
	if err != nil {
		return err
	}

	_, err = out.WriteString(renderAll(characters, *history, *roster))
	if err != nil {
		return fmt.Errorf("writing the render: %w", err)
	}

	return nil
}

// renderAll turns however many records were read into one document. A batch
// of sheets is separated by a rule, which is what a Markdown reader breaks
// pages on; a roster is not, because it is a list.
func renderAll(characters []*chargen.Character, history, roster bool) string {
	if roster {
		return render.Roster(characters)
	}

	parts := make([]string, 0, len(characters))

	for _, character := range characters {
		if history {
			parts = append(parts, render.Transcript(character))

			continue
		}

		parts = append(parts, render.Sheet(character))
	}

	return strings.Join(parts, "\n---\n\n")
}

// readRecord loads a character record. Flags precede the filename because
// Go's flag package stops parsing at the first non-flag argument, so
// `render char.json --history` would leave --history standing as a second
// positional and exit with a usage error.
func readRecord(path string) (*chargen.Character, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var character chargen.Character

	err = json.Unmarshal(content, &character)
	if err != nil {
		return nil, fmt.Errorf("%s is not a character record: %w", path, err)
	}

	return &character, nil
}

// readRecords reads either one record or the JSONL a batch wrote, because
// the two commands the README prints next to each other should compose.
// Twenty characters used to be 777 KB of JSON and no way to see them.
//
// One record and a JSONL file of one are the same thing to everything
// downstream, so the caller never has to ask which it got.
func readRecords(path string) ([]*chargen.Character, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	// One record first, because `new` writes it indented and a record that
	// spans forty lines is not forty records. Only a file that is not one
	// record is read as a file of many.
	var single chargen.Character

	err = json.Unmarshal(content, &single)
	if err == nil {
		return []*chargen.Character{&single}, nil
	}

	// A file of one line that did not parse above is a broken record
	// rather than a batch, and saying "line 1" about it would be noise.
	const shortestBatch = 2

	lines := splitJSONL(content)
	if len(lines) < shortestBatch {
		return nil, fmt.Errorf("%s is not a character record: %w", path, err)
	}

	characters := make([]*chargen.Character, 0, len(lines))

	for i, line := range lines {
		var character chargen.Character

		err = json.Unmarshal(line, &character)
		if err != nil {
			return nil, fmt.Errorf("%s line %d is not a character record: %w", path, i+1, err)
		}

		characters = append(characters, &character)
	}

	return characters, nil
}
