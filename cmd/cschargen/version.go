package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/philoserf/cschargen/chargen"
)

// policyVersion identifies POLICY.md, the auto-mode decision table. It is
// stamped into every record and bumped when that document changes, because
// a different policy is a different character from the same seed.
//
// 0.3.1 is the exception that proves it: the document gained rows for
// `subsector_source` and `homeworld_source`, which only an interactive run
// reaches, so the policy decides nothing it did not decide before and the
// same seed yields the same character. The document changed, so the version
// moves; it moves in the last place because the decisions did not.
const policyVersion = "0.3.1"

// version reports what this binary is, read from the build info the
// toolchain embeds rather than from a build flag -- so it names what was
// built rather than what a flag was told to say.
func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" {
		return "(devel)"
	}

	return info.Main.Version
}

// versionCommand reports the build and the versions a record stamps, so
// that a bug report can be matched against any record this binary wrote.
func versionCommand(out *os.File) error {
	_, err := fmt.Fprintf(out,
		"cschargen %s\nschema %d\npolicy %s\nruleset %s\n",
		version(), chargen.SchemaVersion, policyVersion, chargen.Ruleset)
	if err != nil {
		return fmt.Errorf("writing the version: %w", err)
	}

	return nil
}

// isSet reports whether a flag was given on the command line, as against
// left at its default. A seed of zero is a legitimate seed, so the default
// cannot stand in for "not given".
func isSet(flags *flag.FlagSet, name string) bool {
	found := false

	flags.Visit(func(given *flag.Flag) {
		if given.Name == name {
			found = true
		}
	})

	return found
}
