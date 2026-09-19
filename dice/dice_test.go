package dice_test

import (
	"slices"
	"testing"

	"github.com/philoserf/cschargen/dice"
)

// draws is how many rolls each property below takes. Large enough that a
// shape which can produce an out-of-range value almost certainly does, small
// enough that the whole package still runs in well under a second.
const draws = 20000

// twoD6Expr is the Expr an ordinary check carries, named because three
// tests assert on it and a typo in one of them should not read as a
// disagreement between the tests.
const twoD6Expr = "2d6"

func TestNewIsDeterministic(t *testing.T) {
	t.Parallel()

	a, b := dice.New(7), dice.New(7)

	for i := range 100 {
		x, y := a.ND6(3), b.ND6(3)
		if !slices.Equal(x.Dice, y.Dice) {
			t.Fatalf("roll %d: same seed gave %v then %v", i, x.Dice, y.Dice)
		}
	}
}

func TestDifferentSeedsDiverge(t *testing.T) {
	t.Parallel()

	a, b := dice.New(7), dice.New(8)

	same := 0

	for range 100 {
		if a.D6().Total == b.D6().Total {
			same++
		}
	}

	// Two independent d6 streams agree about one time in six. A hundred
	// draws agreeing every time would mean the seed is not reaching the
	// generator.
	if same == 100 {
		t.Fatal("seeds 7 and 8 produced identical rolls for 100 draws")
	}
}

func TestSeedIsReported(t *testing.T) {
	t.Parallel()

	if got := dice.New(42).Seed(); got != 42 {
		t.Errorf("Seed() = %d, want 42", got)
	}
}

func TestShapesStayInRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		roll      func(*dice.Dice) dice.Roll
		expr      string
		dice      int
		low, high int
	}{
		{"D6", (*dice.Dice).D6, "1d6", 1, 1, 6},
		{"D3", (*dice.Dice).D3, "1d3", 1, 1, 3},
		{"TwoD6", (*dice.Dice).TwoD6, twoD6Expr, 2, 2, 12},
		{"Characteristic", (*dice.Dice).Characteristic, "3d6 drop lowest", 3, 2, 12},
		{"D66", (*dice.Dice).D66, "d66", 2, 11, 66},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			d := dice.New(1)

			for i := range draws {
				r := tc.roll(d)

				if r.Expr != tc.expr {
					t.Fatalf("draw %d: Expr = %q, want %q", i, r.Expr, tc.expr)
				}

				if len(r.Dice) != tc.dice {
					t.Fatalf("draw %d: %d dice, want %d", i, len(r.Dice), tc.dice)
				}

				for _, n := range r.Dice {
					if n < 1 || n > 6 {
						t.Fatalf("draw %d: die showed %d", i, n)
					}
				}

				if r.Total < tc.low || r.Total > tc.high {
					t.Fatalf("draw %d: total %d outside [%d,%d]", i, r.Total, tc.low, tc.high)
				}
			}
		})
	}
}

// TestCharacteristicDropsTheLowest is the rule of p. 13 stated as an
// identity: the total is the sum of all three dice less the smallest.
func TestCharacteristicDropsTheLowest(t *testing.T) {
	t.Parallel()

	d := dice.New(99)

	for i := range draws {
		r := d.Characteristic()

		sum := 0
		for _, n := range r.Dice {
			sum += n
		}

		want := sum - slices.Min(r.Dice)
		if r.Total != want {
			t.Fatalf("draw %d: dice %v totalled %d, want %d", i, r.Dice, r.Total, want)
		}
	}
}

// TestCharacteristicReachesItsBounds guards the drop: a plain 2d6 would
// also sit inside [2,12], but only three-dice-drop-lowest reaches 2 as
// rarely as it does and 12 as often. Both ends must be reachable, or the
// shape is wrong in a way the range test cannot see.
func TestCharacteristicReachesItsBounds(t *testing.T) {
	t.Parallel()

	d := dice.New(3)
	seen := map[int]bool{}

	for range draws {
		seen[d.Characteristic().Total] = true
	}

	for total := 2; total <= 12; total++ {
		if !seen[total] {
			t.Errorf("total %d never came up in %d draws", total, draws)
		}
	}
}

// TestD66HasNoZeroOrSevenDigit is what makes d66 a shape of its own rather
// than 2d6: the results are 11-16, 21-26 and so on, so no digit is 0 and
// none is 7 or more (p. 175).
func TestD66HasNoZeroOrSevenDigit(t *testing.T) {
	t.Parallel()

	d := dice.New(5)
	seen := map[int]bool{}

	for i := range draws {
		r := d.D66()

		tens, units := r.Total/10, r.Total%10
		if tens < 1 || tens > 6 || units < 1 || units > 6 {
			t.Fatalf("draw %d: d66 gave %d", i, r.Total)
		}

		if r.Dice[0] != tens || r.Dice[1] != units {
			t.Fatalf("draw %d: dice %v do not read as %d", i, r.Dice, r.Total)
		}

		seen[r.Total] = true
	}

	if len(seen) != 36 {
		t.Errorf("saw %d of the 36 d66 results", len(seen))
	}
}

func TestND6Sums(t *testing.T) {
	t.Parallel()

	d := dice.New(11)

	for n := 1; n <= 12; n++ {
		r := d.ND6(n)

		sum := 0
		for _, die := range r.Dice {
			sum += die
		}

		if r.Total != sum {
			t.Errorf("ND6(%d): total %d, dice sum to %d", n, r.Total, sum)
		}

		if len(r.Dice) != n {
			t.Errorf("ND6(%d): %d dice", n, len(r.Dice))
		}
	}
}

func TestND6ExprNamesTheShape(t *testing.T) {
	t.Parallel()

	d := dice.New(1)

	for _, tc := range []struct {
		n    int
		want string
	}{{1, "1d6"}, {2, twoD6Expr}, {3, "3d6"}, {9, "9d6"}, {10, "10d6"}, {12, "12d6"}} {
		if got := d.ND6(tc.n).Expr; got != tc.want {
			t.Errorf("ND6(%d).Expr = %q, want %q", tc.n, got, tc.want)
		}
	}
}

func TestND6PanicsBelowOne(t *testing.T) {
	t.Parallel()

	for _, n := range []int{0, -1} {
		t.Run(map[bool]string{true: "zero", false: "negative"}[n == 0], func(t *testing.T) {
			t.Parallel()

			defer func() {
				if recover() == nil {
					t.Errorf("ND6(%d) did not panic", n)
				}
			}()

			dice.New(1).ND6(n)
		})
	}
}
