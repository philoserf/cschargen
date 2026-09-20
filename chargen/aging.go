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

// Apparent age (pp. 124-125).
//
// "Apparent age is a game term used to indicate to the Players how old a
// character might appear to be to the early 21st century onlooker" (p. 124).
// Below tech level 10 it is the character's actual age; from TL 10 up the
// chart on p. 125 maps a ten-year band of actual age onto a five-year band
// of apparent age, and the higher the tech level the slower it climbs.
//
// The chart runs from 30 to 290. Below and above it the book says nothing,
// and the engine reads the nearest printed band. ERRATA E-20.

// AgeBand is an inclusive range of years, which is how both columns of the
// chart are written.
type AgeBand struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// String renders a band the way the chart prints it.
func (b AgeBand) String() string {
	return itoa(b.From) + "-" + itoa(b.To)
}

// apparentAgeChart is p. 125, one row per ten years of actual age, in the
// order the page prints them. The three columns are TL 10, TL 11 and
// TL 12-13.
//
// Each row's actual-age band starts at 30 and runs in tens, so the row is
// found by arithmetic rather than stored: row i covers 30+10i to 39+10i,
// which the page writes as 30-40, 41-50, 51-60 and so on. The overlap at
// each boundary is the page's; 40 and 41 land in different rows and both
// give the same answer in every column, so nothing turns on it.
var apparentAgeChart = [26][3]AgeBand{
	{{20, 25}, {20, 25}, {20, 25}}, // 30-40
	{{20, 25}, {20, 25}, {20, 25}}, // 41-50
	{{30, 35}, {20, 25}, {20, 25}}, // 51-60
	{{30, 35}, {25, 30}, {20, 25}}, // 61-70
	{{35, 40}, {25, 30}, {25, 30}}, // 71-80
	{{35, 40}, {25, 30}, {25, 30}}, // 81-90
	{{40, 45}, {30, 35}, {25, 30}}, // 91-100
	{{40, 45}, {30, 35}, {25, 30}}, // 101-110
	{{45, 50}, {30, 35}, {30, 35}}, // 111-120
	{{45, 50}, {35, 40}, {30, 35}}, // 121-130
	{{50, 55}, {35, 40}, {30, 35}}, // 131-140
	{{50, 55}, {35, 40}, {30, 35}}, // 141-150
	{{55, 60}, {40, 45}, {35, 40}}, // 151-160
	{{55, 60}, {40, 45}, {35, 40}}, // 161-170
	{{60, 65}, {40, 45}, {35, 40}}, // 171-180
	{{60, 65}, {45, 50}, {35, 40}}, // 181-190
	{{65, 70}, {45, 50}, {40, 45}}, // 191-200
	{{65, 70}, {45, 50}, {40, 45}}, // 201-210
	{{70, 75}, {50, 55}, {40, 45}}, // 211-220
	{{70, 75}, {50, 55}, {40, 45}}, // 221-230
	{{75, 80}, {50, 55}, {45, 50}}, // 231-240
	{{75, 80}, {55, 60}, {45, 50}}, // 241-250
	{{80, 85}, {55, 60}, {45, 50}}, // 251-260
	{{80, 85}, {55, 60}, {45, 50}}, // 261-270
	{{85, 90}, {60, 65}, {50, 55}}, // 271-280
	{{85, 90}, {60, 65}, {50, 55}}, // 281-290
}

// The chart's first row and the width of each.
const (
	chartFirstAge = 30
	chartRowYears = 10
)

// apparentAge returns the band a character of this age appears to be in, on
// a homeworld of this tech level, and whether the chart is what said so.
//
// Below TL 10, "your apparent age and your real age are the same" (p. 124).
// The chart itself begins at 30, and below that the same thing is true for
// a different reason -- nobody's apparent age has diverged from their real
// one yet. Above 290 the last printed row holds. ERRATA E-20.
func apparentAge(techLevel, age int) (AgeBand, bool) {
	if techLevel < 10 || age < chartFirstAge {
		return AgeBand{From: age, To: age}, false
	}

	column := 2

	switch techLevel {
	case 10:
		column = 0
	case 11:
		column = 1
	}

	row := min((age-chartFirstAge)/chartRowYears, len(apparentAgeChart)-1)

	return apparentAgeChart[row][column], true
}

// apparentAgeOverForty is the bar twelve careers' enlistment throws are
// measured against: "If you have an apparent age of over 40, take a -2
// modifier to this roll." ERRATA E-21 reads that against a chart that gives
// bands rather than numbers.
func apparentAgeOverForty(band AgeBand) bool {
	return band.From >= 40
}
