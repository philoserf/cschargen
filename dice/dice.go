// Package dice is the seeded random source the engine draws every throw
// from. It knows the shapes of throw the Clement Sector Core Character
// Creation Book asks for and nothing about what any of them mean: a
// characteristic roll here is three dice with the lowest dropped, not a
// characteristic.
//
// Every method returns a [Roll] carrying the individual dice, because the
// generation record logs them (docs/PRD.md FR15) and an audit against the
// book needs to see what was thrown, not only what it totalled.
package dice

import (
	"math/rand/v2"
	"strconv"
)

// The shapes of throw the book uses, named so that the rule is legible at
// the call site rather than appearing as a bare number.
const (
	sides              = 6 // a die
	d3Sides            = 3
	d10Sides           = 10
	percentileHigh     = 100
	ordinaryCheck      = 2 // "the player rolls 2d6" (p. 110)
	characteristicDice = 3 // "Roll 3d6. Drop the score on the lowest" (p. 13)
)

// scramble derives PCG's second state word from the single seed a record
// carries. The constant is the golden-ratio odd number SplitMix64 uses; any
// fixed derivation would do, but it has to be fixed, or the same seed would
// not reproduce the same character.
const scramble = 0x9E37_79B9_7F4A_7C15

// Roll is one resolved roll: the expression that produced it, the dice as
// they fell, and their total.
type Roll struct {
	Expr  string `json:"expr"`
	Dice  []int  `json:"dice"`
	Total int    `json:"total"`
}

// Dice is a seeded stream. The zero value is not usable; construct with
// [New].
type Dice struct {
	rng  *rand.Rand
	seed uint64
}

// New returns a stream seeded with seed. Two streams built from the same
// seed produce the same rolls in the same order, which is what makes a
// record replayable (docs/PRD.md, Replay and provenance contract).
func New(seed uint64) *Dice {
	return &Dice{
		rng:  rand.New(rand.NewPCG(seed, seed^scramble)),
		seed: seed,
	}
}

// Seed reports the seed the stream was built from, for stamping into a
// record.
func (d *Dice) Seed() uint64 {
	return d.seed
}

// D6 throws one six-sided die. Skill tables, benefit tables and the Injury
// table are all 1d6 (pp. 117, 127, 119).
func (d *Dice) D6() Roll {
	n := d.die()

	return Roll{Expr: "1d6", Dice: []int{n}, Total: n}
}

// D3 throws 1d3, which the book asks for only when a mishap ejects a
// character mid-term and the years still have to be counted (p. 121).
func (d *Dice) D3() Roll {
	n := d.rng.IntN(d3Sides) + 1

	return Roll{Expr: "1d3", Dice: []int{n}, Total: n}
}

// D10 throws 1d10, which Step 5 reads as more than a number: parental age
// is "15 x 1d10" and a 1 on that die makes the character the firstborn
// while a 10 makes them the last (p. 59). The face is in Dice, as it is for
// every other shape here.
func (d *Dice) D10() Roll {
	n := d.rng.IntN(d10Sides) + 1

	return Roll{Expr: "1d10", Dice: []int{n}, Total: n}
}

// ND10 throws n ten-sided dice and sums them: the youth and teenage paths
// are 2d10 tables (pp. 68-84) and a sibling can be 3d10 years older than
// the character (p. 60). Panics for n below one, for ND6's reason.
func (d *Dice) ND10(n int) Roll {
	if n < 1 {
		panic("dice: ND10 needs at least one die")
	}

	rolled := make([]int, n)
	total := 0

	for i := range rolled {
		rolled[i] = d.rng.IntN(d10Sides) + 1

		total += rolled[i]
	}

	return Roll{Expr: strconv.Itoa(n) + "d10", Dice: rolled, Total: total}
}

// ND6 throws n six-sided dice and sums them. Panics for n below one: there
// is no such throw in the book, so a caller asking for one is a bug rather
// than a rule.
func (d *Dice) ND6(n int) Roll {
	if n < 1 {
		panic("dice: ND6 needs at least one die")
	}

	rolled := make([]int, n)
	total := 0

	for i := range rolled {
		rolled[i] = d.die()

		total += rolled[i]
	}

	return Roll{Expr: strconv.Itoa(n) + "d6", Dice: rolled, Total: total}
}

// TwoD6 throws 2d6, the book's ordinary check (p. 110).
func (d *Dice) TwoD6() Roll {
	return d.ND6(ordinaryCheck)
}

// Characteristic throws a characteristic: "Roll 3d6. Drop the score on the
// lowest of the three dice and add the remaining two scores" (p. 13). The
// returned Dice hold all three as they fell, lowest included, because the
// record shows the throw rather than its conclusion.
func (d *Dice) Characteristic() Roll {
	r := d.ND6(characteristicDice)

	lowest := r.Dice[0]
	for _, n := range r.Dice[1:] {
		if n < lowest {
			lowest = n
		}
	}

	return Roll{Expr: "3d6 drop lowest", Dice: r.Dice, Total: r.Total - lowest}
}

// D66 throws the book's event tables: two dice read as tens and units, so
// the results run 11-16, 21-26, ... 61-66 rather than 2-12 (p. 175). It is
// one throw rather than two d6 rolls so that the log reads as the page
// does.
func (d *Dice) D66() Roll {
	tens, units := d.die(), d.die()

	return Roll{Expr: "d66", Dice: []int{tens, units}, Total: tens*10 + units}
}

// D100 throws percentile dice, which the origin charts of pp. 43-56 are
// read with. Two ten-sided dice, tens and units, where a double zero is a
// hundred rather than nothing -- so the results run 1 to 100 and every row
// of a chart printed "01-04" through "97-00" is reachable.
func (d *Dice) D100() Roll {
	tens, units := d.rng.IntN(d10Sides), d.rng.IntN(d10Sides)

	total := tens*d10Sides + units
	if total == 0 {
		total = percentileHigh
	}

	return Roll{Expr: "d100", Dice: []int{tens, units}, Total: total}
}

// die draws one die. Every shape above is built from it, so the stream is
// consumed one die at a time and a change in how a shape is assembled is a
// change in the whole sequence after it.
func (d *Dice) die() int {
	return d.rng.IntN(sides) + 1
}
