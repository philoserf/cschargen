// Command cschargen generates rules-accurate Clement Sector characters.
//
// The exit contract: a misused command line exits 1 with a message
// beginning "usage:", and any other failure exits 1 with a message that
// does not. The two are told apart by the message rather than the status,
// so a script wrapping the tool can distinguish its own bug from the
// engine's.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	err := run(os.Args[1:], os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// errUsage is returned for a misused command line. Its message carries the
// "usage:" prefix the exit contract is read by.
var errUsage = errors.New("usage")

func usagef(format string, args ...any) error {
	return fmt.Errorf("%w: %s", errUsage, fmt.Sprintf(format, args...))
}

func run(args []string, out *os.File) error {
	if len(args) == 0 {
		return usagef("cschargen <command> [flags]\n\n%s", commands)
	}

	command, known := dispatch()[args[0]]
	if known {
		return command(args[1:], out)
	}

	switch args[0] {
	case "version":
		return versionCommand(out)
	case "-h", "--help", "help":
		_, err := fmt.Fprintf(out, "usage: cschargen <command> [flags]\n\n%s\n", commands)
		if err != nil {
			return fmt.Errorf("writing the command list: %w", err)
		}

		return nil
	}

	return usagef("unknown command %q\n\n%s", args[0], commands)
}

// dispatch is the commands that take arguments, which is all of them but
// version and the help. It is a map rather than a switch because the switch
// grew past the point where reading it told you anything the map does not.
func dispatch() map[string]func([]string, *os.File) error {
	return map[string]func([]string, *os.File) error{
		"new":    newCommand,
		"batch":  batchCommand,
		"data":   dataCommand,
		"render": renderCommand,
		"replay": replayCommand,
	}
}

const commands = `commands:
  new       generate a character
  batch     generate several, as one JSON record per line
  data      validate a setting data file
  render    turn a record into a character sheet, or its lifepath
  replay    re-run a record from its seed and recorded choices
  version   report the build and the versions a record stamps`

// flagSet builds a flag set that reports its own errors through the usage
// contract rather than printing and exiting on its own.
func flagSet(name string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	return flags
}

// parseOperands parses a command line whose flags may sit before, between or
// after its operands, and returns the operands in the order given.
//
// Go's flag package stops at the first argument that is not a flag, so
// `replay character.json --data setting.json` used to read as three records
// (#171). The record is the thing already in hand and the flag is the
// afterthought, so that is the order people type. Parsing again after each
// operand accepts it with the package as it is; a `--` still ends the flags,
// and everything after it is an operand however it is spelled.
func parseOperands(flags *flag.FlagSet, args []string) ([]string, error) {
	var operands []string

	for {
		err := flags.Parse(args)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		rest := flags.Args()
		if len(rest) == 0 {
			return operands, nil
		}

		if consumed := len(args) - len(rest); consumed > 0 && args[consumed-1] == "--" {
			return append(operands, rest...), nil
		}

		operands = append(operands, rest[0])
		args = rest[1:]
	}
}

// oneOperand is the usage error for a command that takes exactly one, which
// names what it was given: a count alone told a reader looking at one record
// that there was not one record on the line.
func oneOperand(command, what string, operands []string) error {
	if len(operands) == 0 {
		return usagef("%s takes exactly one %s", command, what)
	}

	return usagef("%s takes exactly one %s, got %d: %s",
		command, what, len(operands), strings.Join(quoteAll(operands), " "))
}

func quoteAll(words []string) []string {
	quoted := make([]string, len(words))
	for i, word := range words {
		quoted[i] = strconv.Quote(word)
	}

	return quoted
}
