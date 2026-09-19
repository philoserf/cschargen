package chargen

import "fmt"

// Characteristic is one of the six abilities: "Strength (STR), Dexterity
// (DEX), Endurance (END), Intelligence (INT), Education (EDU), and Charisma
// (CHA)" (p. 13). Charisma stands where Traveller and Cepheus put Social
// Standing; nothing in this engine has a SOC.
type Characteristic int

// The six, in the order the book prints them, which is also the order they
// are assigned in and the order they print on the sheet.
const (
	STR Characteristic = iota
	DEX
	END
	INT
	EDU
	CHA
)

// CharacteristicOrder is the six in the order the book prints them, which
// is also the order they are assigned in and the order they print on the
// sheet. Iterating it is how the engine asks the six questions in the order
// the page asks them.
var CharacteristicOrder = [...]Characteristic{STR, DEX, END, INT, EDU, CHA}

// Physical reports whether this is one of the three physical
// characteristics. The distinction is load-bearing: aging checks them
// separately from the mental three (p. 122), and all three at zero is death
// during generation where two is immobility (p. 124).
func (c Characteristic) Physical() bool {
	return c == STR || c == DEX || c == END
}

// String names the characteristic as the book abbreviates it.
func (c Characteristic) String() string {
	switch c {
	case STR:
		return "STR"
	case DEX:
		return "DEX"
	case END:
		return "END"
	case INT:
		return "INT"
	case EDU:
		return "EDU"
	case CHA:
		return "CHA"
	}

	return fmt.Sprintf("Characteristic(%d)", int(c))
}

// HumanMaximum is the ceiling on an unaltered human's characteristics: "For
// an unaltered human, these characteristics may never rise higher than 15.
// Some uplifts and altrants will have higher maximums" (p. 14).
const HumanMaximum = 15

// modifierBands is the table of p. 14, as upper bounds. A score at or below
// a band's bound and above the previous one takes that band's modifier.
var modifierBands = [...]struct {
	upTo     int
	modifier int
}{
	{0, -3},
	{2, -2},
	{5, -1},
	{8, 0},
	{11, +1},
	{14, +2},
	{17, +3},
	{20, +4},
}

// Modifier is the die modifier a score confers (p. 14). "Characteristics of
// 9 or higher will give you a bonus on tasks which use that characteristic.
// Characteristics of 5 or lower will result in a penalty" (p. 13).
//
// The printed table runs 0 to 20 and no further, because 20 is above any
// score the rules produce -- humans cap at 15 and the highest engineered
// maximum is lower than 21. A score outside that range is clamped to the
// nearest end rather than extrapolated: the book prints no band to
// extrapolate from, and inventing one would be a rule this engine made up.
func Modifier(score int) int {
	if score < 0 {
		return modifierBands[0].modifier
	}

	for _, band := range modifierBands {
		if score <= band.upTo {
			return band.modifier
		}
	}

	return modifierBands[len(modifierBands)-1].modifier
}

// Characteristics is a character's six scores.
type Characteristics struct {
	STR int `json:"str"`
	DEX int `json:"dex"`
	END int `json:"end"`
	INT int `json:"int"`
	EDU int `json:"edu"`
	CHA int `json:"cha"`
}

// Get reads one score.
func (c *Characteristics) Get(which Characteristic) int {
	switch which {
	case STR:
		return c.STR
	case DEX:
		return c.DEX
	case END:
		return c.END
	case INT:
		return c.INT
	case EDU:
		return c.EDU
	case CHA:
		return c.CHA
	}

	return 0
}

// Set writes one score.
func (c *Characteristics) Set(which Characteristic, score int) {
	switch which {
	case STR:
		c.STR = score
	case DEX:
		c.DEX = score
	case END:
		c.END = score
	case INT:
		c.INT = score
	case EDU:
		c.EDU = score
	case CHA:
		c.CHA = score
	}
}

// Modifier is the die modifier for one characteristic. "Whenever a
// characteristic score changes for any reason, the character's modifier for
// that characteristic must be recalculated immediately" (p. 15), which this
// gets for free by never storing one.
func (c *Characteristics) Modifier(which Characteristic) int {
	return Modifier(c.Get(which))
}

// AllPhysicalZero reports the condition of p. 14: "In the highly unlikely
// event that all three physical characteristics (Strength, Dexterity, and
// Endurance) are reduced to 0 during character creation, the character is
// considered to have died before play begins." (The closing quotation mark
// is the sentence's end; godot is configured to accept it.)
func (c *Characteristics) AllPhysicalZero() bool {
	return c.STR == 0 && c.DEX == 0 && c.END == 0
}
