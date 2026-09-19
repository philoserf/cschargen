package dice

// Mod is one named modifier on a throw's total. Modifiers are itemized
// rather than summed at the call site so that the record can show where
// each came from -- "-2 if the character's Apparent Age is forty or
// greater, -1 for each previous career" is two mods, not a -3 (p. 111).
type Mod struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// Throw is a 2d6 roll resolved against a target number: "the player rolls
// 2d6 and adds the Characteristic Modifier associated with the relevant
// characteristic. If the modified result meets or exceeds the target number
// listed for the career, the character is accepted" (p. 110).
type Throw struct {
	Roll    Roll  `json:"roll"`
	Mods    []Mod `json:"mods,omitempty"`
	Target  int   `json:"target"`
	Total   int   `json:"total"` // Roll.Total plus every mod
	Success bool  `json:"success"`
}

// Natural reports the dice before modifiers. The book reads two naturals as
// results in their own right: a natural 12 on a survival roll forces the
// character to continue in the career (p. 125), and Clement Sector reads a
// natural 12 as an exceptional success and a natural 2 as an exceptional
// failure.
func (t Throw) Natural() int {
	return t.Roll.Total
}

// Throw rolls 2d6 against target, applying mods. A throw with no target
// is not one of these -- use [Dice.TwoD6] and read the total.
func (d *Dice) Throw(target int, mods ...Mod) Throw {
	roll := d.TwoD6()

	total := roll.Total
	for _, m := range mods {
		total += m.Value
	}

	return Throw{
		Roll:    roll,
		Mods:    mods,
		Target:  target,
		Total:   total,
		Success: total >= target,
	}
}
