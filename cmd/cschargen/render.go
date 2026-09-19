package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/render"
)

// renderCommand is `cschargen render`: turn a record into the Markdown
// sheet, or with --history into the lifepath transcript.
func renderCommand(args []string, out *os.File) error {
	flags := flagSet("render")
	history := flags.Bool("history", false, "render the lifepath transcript instead of the sheet")

	err := flags.Parse(args)
	if err != nil {
		return usagef("render: %v", err)
	}

	if flags.NArg() != 1 {
		return usagef("render takes exactly one character record")
	}

	character, err := readRecord(flags.Arg(0))
	if err != nil {
		return err
	}

	text := render.Sheet(character)
	if *history {
		text = render.Transcript(character)
	}

	_, err = out.WriteString(text)
	if err != nil {
		return fmt.Errorf("writing the render: %w", err)
	}

	return nil
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
