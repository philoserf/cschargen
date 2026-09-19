package chargen_test

import (
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// TestModifierMatchesThePrintedTable is the table of p. 14, transcribed a
// second time as the value each score confers rather than as bands. Every
// score from 0 to 20 is named, because a band boundary is exactly where a
// transcription goes wrong.
func TestModifierMatchesThePrintedTable(t *testing.T) {
	t.Parallel()

	want := map[int]int{
		0: -3,
		1: -2, 2: -2,
		3: -1, 4: -1, 5: -1,
		6: 0, 7: 0, 8: 0,
		9: +1, 10: +1, 11: +1,
		12: +2, 13: +2, 14: +2,
		15: +3, 16: +3, 17: +3,
		18: +4, 19: +4, 20: +4,
	}

	for score, modifier := range want {
		if got := chargen.Modifier(score); got != modifier {
			t.Errorf("Modifier(%d) = %+d, want %+d", score, got, modifier)
		}
	}
}

// TestModifierClampsOutsideThePrintedRange: the book prints 0 to 20 and no
// further, so a score outside it takes the nearest printed band rather than
// an extrapolated one the book does not offer.
func TestModifierClampsOutsideThePrintedRange(t *testing.T) {
	t.Parallel()

	for _, score := range []int{-1, -5} {
		if got := chargen.Modifier(score); got != -3 {
			t.Errorf("Modifier(%d) = %+d, want -3", score, got)
		}
	}

	for _, score := range []int{21, 99} {
		if got := chargen.Modifier(score); got != +4 {
			t.Errorf("Modifier(%d) = %+d, want +4", score, got)
		}
	}
}

// modEND is the abbreviation the tests assert on in several places.
const modEND = "END"

// TestBonusAndPenaltyThresholds is the prose of p. 13 stated as a property:
// "Characteristics of 9 or higher will give you a bonus ... Characteristics
// of 5 or lower will result in a penalty" (p. 13).
func TestBonusAndPenaltyThresholds(t *testing.T) {
	t.Parallel()

	for score := range 21 {
		got := chargen.Modifier(score)

		switch {
		case score <= 5 && got >= 0:
			t.Errorf("score %d has modifier %+d, but 5 or lower is a penalty", score, got)
		case score >= 9 && got <= 0:
			t.Errorf("score %d has modifier %+d, but 9 or higher is a bonus", score, got)
		case score >= 6 && score <= 8 && got != 0:
			t.Errorf("score %d has modifier %+d, want 0", score, got)
		}
	}
}

func TestCharacteristicNames(t *testing.T) {
	t.Parallel()

	want := []string{"STR", "DEX", modEND, "INT", "EDU", "CHA"}
	for i, which := range chargen.CharacteristicOrder {
		if got := which.String(); got != want[i] {
			t.Errorf("CharacteristicOrder[%d].String() = %q, want %q", i, got, want[i])
		}
	}
}

func TestUnknownCharacteristicStringsItsValue(t *testing.T) {
	t.Parallel()

	if got := chargen.Characteristic(9).String(); got != "Characteristic(9)" {
		t.Errorf("String() = %q", got)
	}
}

// TestPhysicalIsTheThreeOfPage122: aging checks the physical three
// separately from the mental three, and the death conditions of p. 124 turn
// on the same split.
func TestPhysicalIsTheThreeOfPage122(t *testing.T) {
	t.Parallel()

	physical := map[chargen.Characteristic]bool{
		chargen.STR: true, chargen.DEX: true, chargen.END: true,
		chargen.INT: false, chargen.EDU: false, chargen.CHA: false,
	}

	for which, want := range physical {
		if got := which.Physical(); got != want {
			t.Errorf("%s.Physical() = %v, want %v", which, got, want)
		}
	}
}

func TestGetAndSetReachEverySix(t *testing.T) {
	t.Parallel()

	var chars chargen.Characteristics

	for i, which := range chargen.CharacteristicOrder {
		chars.Set(which, i+2)
	}

	for i, which := range chargen.CharacteristicOrder {
		if got := chars.Get(which); got != i+2 {
			t.Errorf("Get(%s) = %d, want %d", which, got, i+2)
		}
	}

	// Every field is distinct, so a Set that wrote the wrong field would
	// show up as two characteristics sharing a score.
	seen := map[int]bool{}

	for _, which := range chargen.CharacteristicOrder {
		score := chars.Get(which)
		if seen[score] {
			t.Errorf("score %d appears twice; Set wrote the wrong field", score)
		}

		seen[score] = true
	}
}

func TestGetUnknownIsZero(t *testing.T) {
	t.Parallel()

	var chars chargen.Characteristics

	if got := chars.Get(chargen.Characteristic(9)); got != 0 {
		t.Errorf("Get(unknown) = %d", got)
	}

	// Set of an unknown characteristic writes nothing rather than panicking
	// or corrupting a real one.
	chars.Set(chargen.Characteristic(9), 12)

	for _, which := range chargen.CharacteristicOrder {
		if chars.Get(which) != 0 {
			t.Errorf("Set(unknown) wrote to %s", which)
		}
	}
}

// TestModifierIsNeverStored is the rule of p. 15 -- "whenever a
// characteristic score changes for any reason, the character's modifier for
// that characteristic must be recalculated immediately" -- held by never
// caching one.
func TestModifierIsNeverStored(t *testing.T) {
	t.Parallel()

	var chars chargen.Characteristics

	chars.Set(chargen.END, 8)

	if got := chars.Modifier(chargen.END); got != 0 {
		t.Fatalf("END 8 modifier = %+d, want 0", got)
	}

	chars.Set(chargen.END, 9)

	if got := chars.Modifier(chargen.END); got != +1 {
		t.Errorf("END 9 modifier = %+d, want +1", got)
	}
}

func TestAllPhysicalZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		chars chargen.Characteristics
		want  bool
	}{
		{"a fresh character", chargen.Characteristics{STR: 7, DEX: 8, END: 9}, false},
		{"two of three at zero", chargen.Characteristics{STR: 0, DEX: 0, END: 4}, false},
		{"all three at zero", chargen.Characteristics{}, true},
		{"mental at zero as well", chargen.Characteristics{INT: 0, EDU: 0, CHA: 0}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.chars.AllPhysicalZero(); got != tc.want {
				t.Errorf("AllPhysicalZero() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHumanMaximum(t *testing.T) {
	t.Parallel()

	if chargen.HumanMaximum != 15 {
		t.Errorf("HumanMaximum = %d, want 15 (p. 14)", chargen.HumanMaximum)
	}
}
