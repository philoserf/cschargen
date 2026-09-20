package chargen_test

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// TestARatingMovesInAGeneratedRecord. The Life Events table is reached at
// 31-36 on every career's d66, so a lifepath of any length reaches it -- and
// a rating that only moves in a unit test is a rating the term loop is not
// using.
func TestARatingMovesInAGeneratedRecord(t *testing.T) {
	t.Parallel()

	moved := 0
	lost := 0

	for seed := range uint64(60) {
		character := lifepath(t, seed, 24)

		for _, event := range character.Events {
			if event.Consequence == nil ||
				event.Consequence.Kind != chargen.ConsequenceRelationship {
				continue
			}

			switch {
			case strings.Contains(event.Consequence.Detail, " -> "):
				moved++
			case strings.Contains(event.Consequence.Detail, "is lost"),
				strings.Contains(event.Consequence.Detail, "lost the "):
				lost++
			}
		}
	}

	if moved == 0 {
		t.Error("no seed in 60 moved a Relationship Rating")
	}

	if lost == 0 {
		t.Error("no seed in 60 lost a relationship")
	}
}
