package main

import (
	"fmt"
	"io"

	"github.com/philoserf/cschargen/chargen"
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
