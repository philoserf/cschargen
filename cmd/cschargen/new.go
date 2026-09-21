package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/setting"
)

// recordMode is the permission a written record takes: readable and
// writable by its owner and nobody else. A character record carries
// nothing secret, but it is the user's file and there is no reason for it
// to be world-readable by default.
const recordMode = 0o600

// newFlags are the flags `cschargen new` accepts, kept together so that the
// command body reads as what it does rather than as what it parses.
type newFlags struct {
	seed      *uint64
	auto      *bool
	name      *string
	species   *string
	techLevel *int
	maxTerms  *int
	terms     *int
	forceCar  *string
	noFamily  *bool
	noYouth   *bool
	noTeenage *bool
	noSchool  *bool
	gender    *string
	look      *string
	goals     *string
	output    *string
	force     *bool
	data      *string

	set *flag.FlagSet
}

// bindNewFlags binds the flags `new` and `batch` share. The command's own
// name is passed in so that its help announces itself: a flag set carries
// the name it was built with, and `batch --help` announced `new` for as
// long as the two shared one.
func bindNewFlags(command string) newFlags {
	flags := flagSet(command)

	return newFlags{
		seed:      flags.Uint64("seed", 0, "seed for the dice; omitted means a random one, recorded in the output"),
		auto:      flags.Bool("auto", false, "resolve every choice with the auto policy (POLICY.md) instead of asking"),
		name:      flags.String("name", "", "the character's name (Step 20, p. 129)"),
		gender:    flags.String("gender", "", "the character's gender identity (Step 20, p. 130)"),
		look:      flags.String("appearance", "", "what the character looks like (Step 20, p. 130)"),
		goals:     flags.String("goals", "", "what the character wants (Step 20, p. 130)"),
		species:   flags.String("species", "human", "human, or an engineered species from the setting data"),
		techLevel: flags.Int("tech-level", 0, "the homeworld's tech level, which gates the aging tables (p. 122)"),
		maxTerms:  flags.Int("max-terms", 0, "the homeworld's maximum terms (p. 42)"),
		terms:     flags.Int("terms", 0, "how many terms to serve; the rules impose no limit, so this is policy (POLICY.md)"),
		forceCar:  flags.String("career", "", "attempt only this career"),
		noFamily:  flags.Bool("skip-family", false, "skip Step 5, which the book allows (p. 57)"),
		noYouth:   flags.Bool("skip-youth", false, "skip Step 6, which the book allows (p. 67)"),
		noTeenage: flags.Bool("skip-teenage", false, "skip Step 7, which the book allows (p. 76)"),
		noSchool:  flags.Bool("skip-education", false, "skip Step 8; higher education is never required (p. 85)"),
		data:      flags.String("data", "", "setting data file; omitted means the repository's invented sample"),
		output:    flags.String("o", "", "write the record here instead of stdout"),
		force:     flags.Bool("force", false, "overwrite the output file if it exists"),
		set:       flags,
	}
}

// parse reads the command line and returns the seed to use, which is the
// one given or, when none was, one chosen here and recorded so that the
// character can be reproduced afterwards.
func (f newFlags) parse(args []string) (uint64, error) {
	err := f.set.Parse(args)
	if err != nil {
		return 0, usagef("%s: %v", f.set.Name(), err)
	}

	if f.set.NArg() > 0 {
		return 0, usagef("%s takes no positional arguments, got %q", f.set.Name(), f.set.Arg(0))
	}

	if isSet(f.set, "seed") {
		return *f.seed, nil
	}

	return rand.Uint64(), nil
}

// finishingFlags are the four that set one character's Step 20 fields. They
// are meaningless on a batch: the same name, gender, appearance and goal on
// twenty people is not a thing anyone means, and applying them quietly is
// worse than refusing. See the rejection in batchCommand.
var finishingFlags = []string{"name", "gender", "appearance", "goals"}

// inputs is everything the command line says about the character, which is
// the same set whether one is generated or twenty.
func (f newFlags) inputs() chargen.Inputs {
	return chargen.Inputs{
		Name:          *f.name,
		Species:       *f.species,
		TechLevel:     *f.techLevel,
		MaxTerms:      *f.maxTerms,
		TermLimit:     *f.terms,
		Career:        *f.forceCar,
		SkipFamily:    *f.noFamily,
		SkipYouth:     *f.noYouth,
		SkipTeenage:   *f.noTeenage,
		SkipEducation: *f.noSchool,
		Gender:        *f.gender,
		Appearance:    *f.look,
		Goals:         *f.goals,
	}
}

// decider is the mode the run is in. Auto applies the fixed policy of
// POLICY.md; without it the player is asked.
//
// The prompts go to stderr and the answers come from stdin, because the
// record goes to stdout so that it can be piped. A record generated either
// way is the same file, and replays the same.
//
//nolint:ireturn // the whole point is to choose between two implementations of one interface
func decider(auto bool) chargen.Decider {
	if auto {
		return chargen.Policy{}
	}

	return chargen.NewPlayer(os.Stdin, os.Stderr)
}

// newCommand is `cschargen new`: generate one character and write its
// record.
func newCommand(args []string, out *os.File) error {
	flags := bindNewFlags("new")

	seed, err := flags.parse(args)
	if err != nil {
		return err
	}

	world, err := loadSetting(*flags.data)
	if err != nil {
		return err
	}

	character, err := chargen.New(chargen.Options{
		Seed:          seed,
		Decider:       decider(*flags.auto),
		EngineVersion: version(),
		PolicyVersion: policyVersion,
		Setting:       world,
		Inputs:        flags.inputs(),
	}).Run()
	if err != nil {
		return fmt.Errorf("generating: %w", err)
	}

	warnIfTheCareerChanged(os.Stderr, character, *flags.forceCar)

	encoded, err := json.MarshalIndent(character, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding the record: %w", err)
	}

	encoded = append(encoded, '\n')

	if *flags.output == "" {
		_, err = out.Write(encoded)
		if err != nil {
			return fmt.Errorf("writing the record: %w", err)
		}

		return nil
	}

	return writeFile(*flags.output, encoded, *flags.force)
}

// loadSetting reads the user's setting file, or falls back to the
// repository's invented sample. A character built on the sample is stamped
// as such and says so on its sheet, so the fallback can never be mistaken
// for the published worlds.
func loadSetting(path string) (*setting.Data, error) {
	if path == "" {
		data, err := setting.Sample()
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		return data, nil
	}

	data, err := setting.Load(path)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}

// writeFile refuses to overwrite without --force, because a generated
// character is not reproducible from anything the caller still holds unless
// they kept the seed.
func writeFile(path string, content []byte, force bool) error {
	if !force {
		_, err := os.Stat(path)
		if err == nil {
			return usagef("%s exists; pass --force to overwrite it", path)
		}

		if !os.IsNotExist(err) {
			return fmt.Errorf("checking %s: %w", path, err)
		}
	}

	err := os.WriteFile(path, content, recordMode)
	if err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	return nil
}
