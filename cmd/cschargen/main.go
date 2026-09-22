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
