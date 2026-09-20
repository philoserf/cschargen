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
	output    *string
	force     *bool
	data      *string

	set *flag.FlagSet
}

func bindNewFlags() newFlags {
	flags := flagSet("new")

	return newFlags{
		seed:      flags.Uint64("seed", 0, "seed for the dice; omitted means a random one, recorded in the output"),
		auto:      flags.Bool("auto", false, "resolve every choice with the auto policy (POLICY.md)"),
		name:      flags.String("name", "", "the character's name"),
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
		return 0, usagef("new: %v", err)
	}

	if f.set.NArg() > 0 {
		return 0, usagef("new takes no positional arguments, got %q", f.set.Arg(0))
	}

	if !*f.auto {
		return 0, usagef("new needs --auto until the interactive mode lands")
	}

	if isSet(f.set, "seed") {
		return *f.seed, nil
	}

	return rand.Uint64(), nil
}

// newCommand is `cschargen new`: generate one character and write its
// record.
func newCommand(args []string, out *os.File) error {
	flags := bindNewFlags()

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
		Decider:       chargen.Policy{},
		EngineVersion: version(),
		PolicyVersion: policyVersion,
		Setting:       world,
		Inputs: chargen.Inputs{
			Name:          *flags.name,
			Species:       *flags.species,
			TechLevel:     *flags.techLevel,
			MaxTerms:      *flags.maxTerms,
			TermLimit:     *flags.terms,
			Career:        *flags.forceCar,
			SkipFamily:    *flags.noFamily,
			SkipYouth:     *flags.noYouth,
			SkipTeenage:   *flags.noTeenage,
			SkipEducation: *flags.noSchool,
		},
	}).Run()
	if err != nil {
		return fmt.Errorf("generating: %w", err)
	}

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
