package render

import (
	"fmt"
	"strings"

	"github.com/philoserf/cschargen/chargen"
)

// Transcript walks the generation record and renders it as the character's
// lifepath in game terms: every step entered, every throw made, every
// choice taken and what came of it.
//
// It is the third use the log is kept for. The other two are audit -- walk
// it against the book, which is what the page cites are for -- and replay.
func Transcript(character *chargen.Character) string {
	var out strings.Builder

	name := character.Provenance.Inputs.Name
	if name == "" {
		name = "(unnamed)"
	}

	fmt.Fprintf(&out, "# Lifepath: %s\n\n", name)
	fmt.Fprintf(&out, "Seed %d, %d events.\n\n", character.Provenance.RNG.Seed, len(character.Events))

	for _, event := range character.Events {
		writeEvent(&out, event)
	}

	return out.String()
}

func writeEvent(out *strings.Builder, event chargen.Event) {
	switch event.Kind {
	case chargen.EventStep:
		fmt.Fprintf(out, "\n## %s (%s)\n\n", event.Step.Name, event.Step.Cite)
	case chargen.EventThrow:
		fmt.Fprintf(out, "- %4d  throw  %s\n", event.Seq, throwLine(event.Throw))
	case chargen.EventChoice:
		fmt.Fprintf(out, "- %4d  choice %s: %s (%s)\n",
			event.Seq, event.Choice.Prompt, event.Choice.Options[event.Choice.Chosen],
			event.Choice.Decider)
	case chargen.EventConsequence:
		fmt.Fprintf(out, "- %4d  -> %s [%s, from %d]\n",
			event.Seq, event.Consequence.Detail, event.Consequence.Kind, event.Consequence.Cause)
	}
}

// throwLine renders a throw the way the page asks for it: the dice as they
// fell, the modifiers itemized, and -- where there was a target -- whether
// it was met.
func throwLine(throw *chargen.ThrowEvent) string {
	var out strings.Builder

	fmt.Fprintf(&out, "%s %v = %d", throw.Expr, throw.Dice, throw.Total)

	for _, mod := range throw.Mods {
		fmt.Fprintf(&out, " %+d(%s)", mod.Value, mod.Name)
	}

	if throw.Target != nil {
		outcome := "failed"
		if throw.Success != nil && *throw.Success {
			outcome = "made it"
		}

		fmt.Fprintf(&out, " vs %d+ -- %s", *throw.Target, outcome)
	}

	fmt.Fprintf(&out, "  (%s)", throw.Cite)

	return out.String()
}
