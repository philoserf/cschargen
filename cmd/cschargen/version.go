package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/philoserf/cschargen/chargen"
)

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
		version(), chargen.SchemaVersion, chargen.PolicyVersion, chargen.Ruleset)
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
