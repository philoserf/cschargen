package dice_test

import (
	"testing"

	"github.com/philoserf/cschargen/dice"
)

func TestThrowResolvesAgainstTheTarget(t *testing.T) {
	t.Parallel()

	d := dice.New(17)

	for i := range draws {
		got := d.Throw(8)

		if got.Target != 8 {
			t.Fatalf("draw %d: Target = %d", i, got.Target)
		}

		if got.Total != got.Roll.Total {
			t.Fatalf("draw %d: no mods, but Total %d != roll %d", i, got.Total, got.Roll.Total)
		}

		if got.Success != (got.Total >= 8) {
			t.Fatalf("draw %d: total %d, Success = %v", i, got.Total, got.Success)
		}
	}
}

func TestThrowAppliesEveryMod(t *testing.T) {
	t.Parallel()

	// The two modifiers the Celebrity career applies to its enlistment
	// throw (p. 111), itemized as the record keeps them.
	mods := []dice.Mod{
		{Name: "apparent age 40+", Value: -2},
		{Name: "one previous career", Value: -1},
	}

	d := dice.New(23)

	for i := range draws {
		got := d.Throw(8, mods...)

		if got.Total != got.Roll.Total-3 {
			t.Fatalf("draw %d: roll %d with -3 of mods gave %d", i, got.Roll.Total, got.Total)
		}

		if len(got.Mods) != 2 {
			t.Fatalf("draw %d: %d mods recorded, want 2", i, len(got.Mods))
		}

		if got.Success != (got.Total >= 8) {
			t.Fatalf("draw %d: total %d, Success = %v", i, got.Total, got.Success)
		}
	}
}

// TestNaturalIgnoresMods is the distinction p. 125 rests on: the natural 12
// that forces a character to continue in a career is the dice, not the
// modified total, so a +2 can produce a 12 that is not a natural one.
func TestNaturalIgnoresMods(t *testing.T) {
	t.Parallel()

	d := dice.New(29)

	naturals, modified := 0, 0

	for range draws {
		got := d.Throw(8, dice.Mod{Name: "commendation", Value: 2})

		if got.Natural() == 12 {
			naturals++
		}

		if got.Total == 12 && got.Natural() != 12 {
			modified++
		}
	}

	if naturals == 0 {
		t.Error("no natural 12 in the sample")
	}

	if modified == 0 {
		t.Error("no modified 12 that was not a natural 12 -- Natural is reading the total")
	}
}

func TestThrowRecordsTwoDice(t *testing.T) {
	t.Parallel()

	got := dice.New(1).Throw(8)

	if got.Roll.Expr != twoD6Expr {
		t.Errorf("Expr = %q, want 2d6", got.Roll.Expr)
	}

	if len(got.Roll.Dice) != 2 {
		t.Errorf("%d dice, want 2", len(got.Roll.Dice))
	}
}

// A target no throw can reach and one every throw clears: the boundaries
// have to fail and succeed closed, because an unset target would otherwise
// read as "anything passes".
func TestThrowBoundaries(t *testing.T) {
	t.Parallel()

	d := dice.New(31)

	for range 500 {
		if d.Throw(13).Success {
			t.Fatal("2d6 reached 13")
		}

		if !d.Throw(2).Success {
			t.Fatal("2d6 failed to reach 2")
		}
	}
}
