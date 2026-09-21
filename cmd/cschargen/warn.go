package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/setting"
)

// warnIfTheCareerChanged says so when the career a character entered is not
// the career `--career` asked for.
//
// That happens for a reason the rules give: a failed enlistment closes a
// career for two terms (p. 110), and with nothing else the flag allows, the
// character drifts into Vagabond (p. 111). Across twenty seeds asking for
// Corporate Shipper, nine got one.
//
// The record has always been honest about it -- provenance keeps the career
// asked for beside the career served -- but a referee scripting a cast reads
// the sheets, not the provenance, and finds out twenty NPCs later. So the
// command says it once, on stderr, and still exits 0: the character
// generated successfully, and somebody running a hundred should not have to
// check a hundred exit codes.
func warnIfTheCareerChanged(notes io.Writer, character *chargen.Character, asked string) {
	if asked == "" || servedIn(character, asked) {
		return
	}

	entered := "no career at all"
	if len(character.State.Services) > 0 {
		entered = character.State.Services[0].Career
	}

	// A note nobody can read is not worth failing a generation over: the
	// record is already written and correct.
	_, _ = fmt.Fprintf(notes,
		"note: asked for %s; the character failed to enlist and entered %s"+
			" (p. 110). Try another seed.\n", asked, entered)
}

// warnHowManyCareersChanged is the same note for a batch, counted rather
// than repeated: a hundred records asking for one career would otherwise
// print up to fifty-five identical lines, and a count is what a referee
// needs to decide whether to draw again.
func warnHowManyCareersChanged(
	notes io.Writer, characters []*chargen.Character, asked string,
) {
	if asked == "" {
		return
	}

	missed := 0

	for _, character := range characters {
		if !servedIn(character, asked) {
			missed++
		}
	}

	if missed == 0 {
		return
	}

	_, _ = fmt.Fprintf(notes,
		"note: %d of %d failed to enlist in %s and entered another career"+
			" (p. 110).\n", missed, len(characters), asked)
}

// servedIn reports whether a character ever entered a named career.
func servedIn(character *chargen.Character, name string) bool {
	for _, served := range character.State.Services {
		if served.Career == name {
			return true
		}
	}

	return false
}

// checkOrigin holds `--subsector` and `--homeworld` against the setting
// data before a character is generated.
//
// The engine refuses an unknown name too, but it does so partway through
// Step 3, which reads as a generation failure rather than as the typo it
// is. #81 is the same argument about `--career`: a name the data does not
// have is a flag error, and the command is where a flag error belongs.
//
// The message lists what the data does have. A referee who misremembers a
// world's spelling is better served by the list than by being told their
// spelling is wrong.
func checkOrigin(world *setting.Data, subsector, homeworld string) error {
	if subsector != "" {
		if _, found := world.Subsector(subsector); !found {
			return usagef("no subsector named %q in %s; it has %s",
				subsector, world.Name, joinNames(subsectorNames(world)))
		}
	}

	if homeworld == "" {
		return nil
	}

	_, sub, found := world.World(homeworld)
	if !found {
		return usagef("no world named %q in %s", homeworld, world.Name)
	}

	// Both given and disagreeing is a mistake worth naming: silently
	// preferring one would generate a character from a place the command
	// line did not ask for.
	if subsector != "" && sub.Name != subsector {
		return usagef("%s is in %s, not %s", homeworld, sub.Name, subsector)
	}

	return nil
}

func subsectorNames(world *setting.Data) []string {
	names := make([]string, 0, len(world.Subsectors))
	for _, sub := range world.Subsectors {
		names = append(names, sub.Name)
	}

	return names
}

// joinNames lists what the data does have. A validated setting file always
// has at least one subsector, so there is no empty case to write, and one
// name needs no "and" in front of it.
func joinNames(names []string) string {
	const needsAnAnd = 2

	if len(names) < needsAnAnd {
		return strings.Join(names, "")
	}

	return strings.Join(names[:len(names)-1], ", ") + " and " + names[len(names)-1]
}

// checkCareer holds `--career` against the careers the book prints, before
// a character is generated.
//
// #63 gave the flag a note for the case where the career asked for is not
// the career entered, because a failed enlistment closes a career for two
// terms (p. 110) and the character drifts into Vagabond. That note is right
// for a career that exists. For one that does not, it said the character
// "failed to enlist" and to "try another seed" -- when no enlistment was
// attempted and no seed would help.
//
// A name the book does not print is a flag error, so it fails here rather
// than becoming a generation outcome. Matching is case-insensitive: every
// career name the program prints is capitalised, and a referee types what
// they read.
func checkCareer(asked string) (string, error) {
	if asked == "" {
		return "", nil
	}

	for _, def := range career.All() {
		if strings.EqualFold(def.Name, asked) {
			return def.Name, nil
		}
	}

	if nearest := nearestCareer(asked); nearest != "" {
		return "", usagef("no career named %q; did you mean %s?", asked, nearest)
	}

	return "", usagef("no career named %q", asked)
}

// nearestCareer is the one career whose name contains what was asked for,
// or is contained by it, and only where exactly one does: "navy" reaches
// nothing useful, because three careers have it in their name.
func nearestCareer(asked string) string {
	folded := strings.ToLower(asked)

	var found string

	for _, def := range career.All() {
		name := strings.ToLower(def.Name)
		if !strings.Contains(name, folded) && !strings.Contains(folded, name) {
			continue
		}

		if found != "" {
			return ""
		}

		found = def.Name
	}

	return found
}

// checkNames holds every flag naming something the setting data or the book
// has to have, before a character is generated. A name that is not there is
// a flag error, and a flag error belongs to the command rather than to a
// lifepath that is already halfway through Step 3.
//
// The career is rewritten to the book's own spelling, so that a note about
// it later names it the way every other line does.
func checkNames(world *setting.Data, flags newFlags) error {
	err := checkOrigin(world, *flags.subsector, *flags.homeworld)
	if err != nil {
		return err
	}

	*flags.forceCar, err = checkCareer(*flags.forceCar)

	return err
}
