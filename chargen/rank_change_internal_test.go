package chargen

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/career"
)

// rankedEngine is a generator standing in a career at a rank, which is what
// a demotion needs something to take.
func rankedEngine(t *testing.T, rank int) *Generator {
	t.Helper()

	gen := engine(t, 101)
	colonist := career.Colonist()

	gen.service.career = &colonist
	gen.service.assignment = colonist.Assignments[0]
	gen.service.rank = rank

	return gen
}

// TestADemotionMovesTheNumberAndNothingElse is ERRATA E-36: a demotion
// takes the rank, and the skills the rank granted stay granted, because
// the book has no rule for taking a skill level away.
func TestADemotionMovesTheNumberAndNothingElse(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 3)

	gen.char.State.GainSkill("Drive", "Wheeled", 1)

	err := gen.apply(career.Effect{Kind: career.EffectRank, Levels: -1, Detail: "lose one rank"}, 0)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	if gen.service.rank != 2 {
		t.Errorf("rank is %d, want 2", gen.service.rank)
	}

	if gen.char.State.SkillLevel("Drive") != 1 {
		t.Error("a demotion took a skill level with it")
	}
}

// TestADemotionStopsAtTheBottom. Rank 0 is the bottom of every table, and
// a character already there loses nothing -- but the result fired, so the
// record says what it did.
func TestADemotionStopsAtTheBottom(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 0)

	err := gen.apply(career.Effect{Kind: career.EffectRank, Levels: -2, Detail: "lose two ranks"}, 0)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	if gen.service.rank != 0 {
		t.Errorf("rank is %d, want 0", gen.service.rank)
	}

	if !strings.Contains(lastDetail(t, gen), "already at the lowest rank") {
		t.Errorf("the record does not say why nothing happened: %q", lastDetail(t, gen))
	}
}

// TestAnEventPromotionGrantsTheRanksBenefits. p. 116 attaches a rank
// benefit to holding the rank, not to the advancement throw that usually
// reaches it, so an event that promotes grants it too.
func TestAnEventPromotionGrantsTheRanksBenefits(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 0)

	err := gen.apply(career.Effect{Kind: career.EffectRank, Detail: "gain a rank"}, 0)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	if gen.service.rank != 1 {
		t.Fatalf("rank is %d, want 1", gen.service.rank)
	}

	// Colonist's rank 1 prints a benefit, so something was granted beside
	// the number. Which skill it is belongs to the transcription, not here.
	if len(gen.char.State.Skills) == 0 {
		t.Error("an event promotion granted the rank and none of its benefits")
	}
}

// TestARankChangeOutsideACareer is the guard: there is no rank to move,
// and the record says so rather than moving a number that means nothing.
func TestARankChangeOutsideACareer(t *testing.T) {
	t.Parallel()

	gen := engine(t, 103)

	err := gen.apply(career.Effect{Kind: career.EffectRank, Levels: -1, Detail: "lose one rank"}, 0)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	if !strings.Contains(lastDetail(t, gen), "no rank") {
		t.Errorf("the record does not say there was no rank: %q", lastDetail(t, gen))
	}
}
