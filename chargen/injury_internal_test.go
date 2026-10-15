package chargen

import (
	"testing"

	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/dice"
)

// injured is a generator holding a character with the given physical
// characteristics, answering every choice with the auto policy.
func injured(str, dex, end int) *Generator {
	gen := &Generator{
		dice:    dice.New(1),
		log:     &Log{},
		decider: Policy{},
		char:    &Character{},
	}

	gen.char.State.Characteristics.Set(STR, str)
	gen.char.State.Characteristics.Set(DEX, dex)
	gen.char.State.Characteristics.Set(END, end)

	return gen
}

// TestSevereInjuryComesOutOfTheHighest is p. 119's second row under the
// auto policy: "Reduce one physical characteristic by 1d6", and the policy
// places the loss on the highest score rather than on STR, which is printed
// first (POLICY.md).
func TestSevereInjuryComesOutOfTheHighest(t *testing.T) {
	t.Parallel()

	gen := injured(3, 6, 11)

	err := gen.applyAll(career.Injury()[1].Effects, 0)
	if err != nil {
		t.Fatalf("applyAll: %v", err)
	}

	got := gen.char.State.Characteristics

	if got.Get(STR) != 3 || got.Get(DEX) != 6 {
		t.Errorf("STR %d, DEX %d: the loss left END, the highest", got.Get(STR), got.Get(DEX))
	}

	if lost := 11 - got.Get(END); lost < 1 || lost > 6 {
		t.Errorf("END lost %d, want 1d6", lost)
	}
}

// TestNearlyKilledReducesAllThree is p. 119's first row: 1d6 from one
// physical characteristic, and "other physical characteristics by 2 (or one
// of them by 4)". The 1d6 goes to the highest, and the second choice --
// which names a characteristic twice -- takes the first option printed,
// both by 2.
func TestNearlyKilledReducesAllThree(t *testing.T) {
	t.Parallel()

	gen := injured(9, 12, 7)

	err := gen.applyAll(career.Injury()[0].Effects, 0)
	if err != nil {
		t.Fatalf("applyAll: %v", err)
	}

	got := gen.char.State.Characteristics

	if got.Get(STR) != 7 || got.Get(END) != 5 {
		t.Errorf("STR %d, END %d, want 7 and 5: the other two lose 2 each", got.Get(STR), got.Get(END))
	}

	if lost := 12 - got.Get(DEX); lost < 1 || lost > 6 {
		t.Errorf("DEX lost %d, want 1d6", lost)
	}
}

// TestLossScoresOnlyWhereEveryOptionIsALoss holds the boundary of the
// policy's one exception: a choice is scored only when every option opens
// with a loss from a different characteristic.
func TestLossScoresOnlyWhereEveryOptionIsALoss(t *testing.T) {
	t.Parallel()

	gen := injured(4, 8, 6)
	nearlyKilled := career.Injury()[0].Effects[0]

	cases := []struct {
		name    string
		options []career.Option
		want    []int
	}{
		{"which characteristic took the 1d6", nearlyKilled.Options, []int{4, 8, 6}},
		{"heavily injured", career.Injury()[3].Effects[0].Options, []int{4, 8, 6}},
		{"an eye or a limb", career.Injury()[2].Effects[0].Options, []int{4, 8}},
		{"both by 2 or one by 4", nearlyKilled.Options[0].Effects[1].Options, nil},
		{"a gain", []career.Option{
			{Label: "STR", Effects: []career.Effect{{Kind: career.EffectCharacteristic, Characteristic: "STR", Delta: 1}}},
			{Label: "DEX", Effects: []career.Effect{{Kind: career.EffectCharacteristic, Characteristic: "DEX", Delta: 1}}},
		}, nil},
		{"an empty option", []career.Option{{Label: "nothing"}}, nil},
	}

	for _, c := range cases {
		got := gen.lossScores(c.options)

		if len(got) != len(c.want) {
			t.Errorf("%s: scores %v, want %v", c.name, got, c.want)

			continue
		}

		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("%s: scores %v, want %v", c.name, got, c.want)

				break
			}
		}
	}
}
