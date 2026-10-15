package dice

import "strings"

// Expression is one quantity as the tables write it: a bare number, "NdM",
// "NdMxK", "NdM+K" or "NdM-K", where M is 6, 3 or 10 -- the three dice the
// book throws for a quantity. Count is zero for a bare number, which is
// then the whole of Offset.
//
// The multiplier and the offset are exclusive because the book never writes
// both: "1d6 x ₶100,000" is a sum of money and "1d3+1 Contacts" is a count
// of people.
type Expression struct {
	Count      int
	Sides      int
	Multiplier int
	Offset     int
}

// ParseExpression reads one expression, reporting false for anything else.
// Anything else is a transcription error rather than a rule, so the caller
// refuses it rather than guessing what the page meant.
//
// It is here rather than beside the engine that throws it because the
// setting validator holds a species' characteristic methods to the same
// reading generation gives them, and the validator cannot import the
// engine. One parser means a file the validator passes is a file the engine
// can roll.
func ParseExpression(expr string) (Expression, bool) {
	expr, multiplier, offset, ok := modifiers(expr)
	if !ok {
		return Expression{}, false
	}

	index := strings.Index(expr, "d")
	if index < 0 {
		value, valid := atoi(expr)
		if !valid {
			return Expression{}, false
		}

		return Expression{Multiplier: 1, Offset: value*multiplier + offset}, true
	}

	count, valid := atoi(expr[:index])
	if !valid || count < 1 {
		return Expression{}, false
	}

	faces, valid := atoi(expr[index+1:])
	if !valid {
		return Expression{}, false
	}

	switch faces {
	case sides, d3Sides, d10Sides:
	default:
		return Expression{}, false
	}

	return Expression{Count: count, Sides: faces, Multiplier: multiplier, Offset: offset}, true
}

// modifiers splits the multiplier and the offset off an expression, which
// default to 1 and 0 where it carries none.
func modifiers(expr string) (string, int, int, bool) {
	expr, offset, ok := suffix(expr, "+", 0)
	if !ok {
		return "", 0, 0, false
	}

	// A species' characteristic method subtracts as often as it adds:
	// "roll 2d6-2 for STR and END" (p. 23). Only one of the two suffixes
	// can apply, because no expression in the book carries both.
	if offset == 0 {
		var down int

		expr, down, ok = suffix(expr, "-", 0)
		if !ok {
			return "", 0, 0, false
		}

		offset = -down
	}

	expr, multiplier, ok := suffix(expr, "x", 1)
	if !ok {
		return "", 0, 0, false
	}

	return expr, multiplier, offset, true
}

// suffix splits a trailing "<sep><number>" off an expression, returning the
// rest and the number -- or absent is the default, which is what an
// expression without that suffix means. A separator with nothing readable
// after it is a transcription error, not a default.
func suffix(expr, sep string, absent int) (string, int, bool) {
	rest, after, found := strings.Cut(expr, sep)
	if !found {
		return expr, absent, true
	}

	value, ok := atoi(after)
	if !ok {
		return expr, 0, false
	}

	return rest, value, true
}

// decimal is the base the book's numbers are written in.
const decimal = 10

func atoi(s string) (int, bool) {
	if s == "" {
		return 0, false
	}

	value := 0

	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, false
		}

		value = value*decimal + int(r-'0')
	}

	return value, true
}
