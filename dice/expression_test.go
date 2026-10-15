package dice_test

import (
	"testing"

	"github.com/philoserf/cschargen/dice"
)

// TestParseExpression holds the reading the engine throws a quantity by and
// the validator checks a species' method against: one parser, so a file
// the validator passes is one the engine can roll (#170).
func TestParseExpression(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expr string
		want dice.Expression
		bad  bool
	}{
		{expr: "100", want: dice.Expression{Multiplier: 1, Offset: 100}},
		{expr: "2d6-2", want: dice.Expression{Count: 2, Sides: 6, Multiplier: 1, Offset: -2}},
		{expr: "1d3+1", want: dice.Expression{Count: 1, Sides: 3, Multiplier: 1, Offset: 1}},
		{expr: "1d10x15", want: dice.Expression{Count: 1, Sides: 10, Multiplier: 15}},
		{expr: "2d6+banana", bad: true},
		{expr: "2d6 - 2", bad: true},
		{expr: "1d4", bad: true},
		{expr: "2d", bad: true},
		// Zero dice is no throw the book makes, and the engine's ND6 panics
		// on it rather than rolling nothing.
		{expr: "0d6", bad: true},
		{expr: "", bad: true},
	}

	for _, tc := range tests {
		got, ok := dice.ParseExpression(tc.expr)
		if ok == tc.bad {
			t.Errorf("%q: ok = %v", tc.expr, ok)

			continue
		}

		if !tc.bad && got != tc.want {
			t.Errorf("%q = %+v, want %+v", tc.expr, got, tc.want)
		}
	}
}
