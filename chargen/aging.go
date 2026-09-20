package chargen

// Step 17's aging tables (pp. 122-123).
//
// The column is headed "Term" and for TL 9 and below it reads 6-8, 9-10, 11,
// 12+ -- plainly term numbers. The higher tech levels read 18-28, 33-48,
// 44-58, which are term numbers too: at four years a term, term 44 is a
// character around 194, and the apparent age chart (p. 125) runs to 290.
//
// The index is terms served in total, not terms in the current career, and
// not derived from age -- the 1d3 mishap increment (p. 121) already means
// age is not 18 + 4n.
//
// The book prints seven of these tables. The four here are the human ones;
// the other three belong to species this engine does not generate yet, and
// two of their headings are proper names the OGL notice reserves (p. 335).
// They arrive with the species work, out of the setting data file.

// AgingCheck is one characteristic and the number it has to make.
type AgingCheck struct {
	Characteristic string
	Number         int
}

// AgingBand is one row of an aging table: the terms it covers and the
// checks that come due in each of them. Through is zero on a band the page
// leaves open -- "12+", "49+".
type AgingBand struct {
	From    int
	Through int
	Checks  []AgingCheck
}

// covers reports whether a band applies at this term count.
func (b AgingBand) covers(term int) bool {
	if term < b.From {
		return false
	}

	return b.Through == 0 || term <= b.Through
}

// The three shapes the bands are written in. Every table is built from
// these four rows in some order, which is the pattern worth naming: the
// tables differ in when the bands start, not in what they ask.
func threePhysical(number int) []AgingCheck {
	return []AgingCheck{
		{Characteristic: "STR", Number: number},
		{Characteristic: "DEX", Number: number},
		{Characteristic: "END", Number: number},
	}
}

func physicalPlusMind(number int) []AgingCheck {
	return append(threePhysical(number),
		AgingCheck{Characteristic: "INT", Number: number},
		AgingCheck{Characteristic: "CHA", Number: number},
	)
}

// everything is the only row where the six checks do not share a number:
// CHA stays at 9+ while the rest move to 10+ or 11+.
func everything(number int) []AgingCheck {
	return append(threePhysical(number),
		AgingCheck{Characteristic: "INT", Number: number},
		AgingCheck{Characteristic: "EDU", Number: number},
		AgingCheck{Characteristic: "CHA", Number: 9},
	)
}

// humanAging returns the table for a homeworld tech level (pp. 122-123).
//
// TL 11 and TL 12-13 stop where they stop: TL 9 and TL 10 each end on an
// open band and the other two do not, so nothing is printed for a TL 11
// character past term 58. The engine makes no check there. ERRATA E-15.
//
// Nothing is printed above TL 13 and the setting validator accepts up to
// TL 20, so a higher homeworld reads the TL 12-13 table: each band delays
// aging and none abolishes it. ERRATA E-16.
func humanAging(techLevel int) []AgingBand {
	switch {
	case techLevel <= 9:
		return []AgingBand{
			{From: 6, Through: 8, Checks: threePhysical(8)},
			{From: 9, Through: 10, Checks: threePhysical(9)},
			{From: 11, Through: 11, Checks: physicalPlusMind(9)},
			{From: 12, Checks: everything(10)},
		}
	case techLevel == 10:
		return []AgingBand{
			{From: 18, Through: 28, Checks: threePhysical(8)},
			{From: 29, Through: 38, Checks: threePhysical(9)},
			{From: 39, Through: 48, Checks: physicalPlusMind(9)},
			{From: 49, Checks: everything(10)},
		}
	case techLevel == 11:
		return []AgingBand{
			{From: 33, Through: 48, Checks: threePhysical(8)},
			{From: 49, Through: 58, Checks: threePhysical(9)},
		}
	default:
		return []AgingBand{
			{From: 44, Through: 58, Checks: threePhysical(8)},
		}
	}
}

// agingChecksAt returns the checks due at a term count, and whether any
// band covers it at all.
func agingChecksAt(techLevel, term int) ([]AgingCheck, bool) {
	for _, band := range humanAging(techLevel) {
		if band.covers(term) {
			return band.Checks, true
		}
	}

	return nil, false
}

// The Aging Crisis (p. 123) and the four states beneath it (pp. 123-124).
//
// A characteristic reduced to 0 by aging puts the character "on the verge of
// death or permanent incapacity": 1d6 x 1000 credits of emergency treatment
// restores it to 1, and without the payment they die.
//
// The crisis rules say "physical" and "mental" and enumerate neither, but
// p. 14 does name one of them -- "all three physical characteristics
// (Strength, Dexterity, and Endurance)" -- and there are six, so the other
// three follow. ERRATA E-17 records the cross-reference.
var (
	physicalCharacteristics = []Characteristic{STR, DEX, END}
	mentalCharacteristics   = []Characteristic{INT, EDU, CHA}
)

// atZeroEndsIt is "two or more", which three of the four terminal states
// turn on (pp. 123-124).
const atZeroEndsIt = 2

// Fate is a terminal state aging reached. Empty is the ordinary case.
//
// Death is a value rather than an error: the record is complete and valid,
// the sheet says the character died at the age they died, and the CLI exits
// 0. Only a misused command line and a malformed data file exit non-zero.
type Fate string

const (
	// FateDied is p. 124's two deaths: all three physical characteristics
	// at 0, or two or more mental ones.
	FateDied Fate = "died"

	// FateIncapacitated is two or more physical characteristics at 0 and
	// unrestored: "incapable of independent movement ... may not continue
	// in career generation".
	FateIncapacitated Fate = "incapacitated"
)

// crisisPayment is the emergency treatment's price, as an expression
// rollExpression reads.
const crisisPayment = "1d6x1000"

// zeroed counts how many of a group have been reduced to 0.
func (c *Characteristics) zeroed(group []Characteristic) int {
	count := 0

	for _, which := range group {
		if c.Get(which) == 0 {
			count++
		}
	}

	return count
}
