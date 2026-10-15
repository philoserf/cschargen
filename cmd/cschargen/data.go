package main

import (
	"fmt"
	"os"

	"github.com/philoserf/cschargen/setting"
)

// dataCommand is `cschargen data validate`: read a setting file and report
// what is wrong with it.
//
// It exists because the file is the user's own transcription of sixty-odd
// rows, and a mistake in it should be reported by what it is rather than
// surfacing as a panic halfway through a lifepath.
func dataCommand(args []string, out *os.File) error {
	if len(args) == 0 {
		return usagef("%s", dataUsage)
	}

	// Asked for, the usage is an answer rather than a misuse: every other
	// command answers --help, and `data` told the user their flag was an
	// unknown subcommand (#173).
	if asksForHelp(args[0]) {
		return dataHelp(out)
	}

	if args[0] != "validate" {
		return usagef("unknown data subcommand %q", args[0])
	}

	flags := flagSet("data validate")

	operands, err := parseOperands(flags, args[1:])
	if err != nil {
		return usagef("data validate: %v", err)
	}

	if len(operands) != 1 {
		return oneOperand("data validate", "setting file", operands)
	}

	data, err := setting.Load(operands[0])
	if err != nil {
		// Not wrapped with a prefix: an InvalidError already names the file
		// and lists every problem in it, and "data validate: x has 3
		// problems" reads worse than the list on its own.
		return fmt.Errorf("%w", err)
	}

	worlds := 0
	for _, sub := range data.Subsectors {
		worlds += len(sub.Worlds)
	}

	_, err = fmt.Fprintf(out, "%s: %d subsectors, %d worlds, %d species\nsha256 %s\n",
		data.Name, len(data.Subsectors), worlds, len(data.Species), data.Hash)
	if err != nil {
		return fmt.Errorf("writing the result: %w", err)
	}

	return nil
}

const dataUsage = "data <subcommand>\n\nsubcommands:\n  validate <file>   check a setting file"

func dataHelp(out *os.File) error {
	_, err := fmt.Fprintf(out, "usage: %s\n", dataUsage)
	if err != nil {
		return fmt.Errorf("writing the data usage: %w", err)
	}

	return nil
}

// asksForHelp is the three spellings the top level answers, so `data`
// answers the same ones.
func asksForHelp(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}
