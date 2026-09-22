package chargen

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/career"
)

// stashEngine is a generator standing inside a career, which is what a
// benefit row is applied from.
func stashEngine(t *testing.T, seed uint64) *Generator {
	t.Helper()

	gen := engine(t, seed)
	colonist := career.Colonist()

	gen.service.career = &colonist

	return gen
}

// mustApply applies one effect and fails the test if it errors.
func mustApply(t *testing.T, gen *Generator, effect career.Effect) {
	t.Helper()

	err := gen.apply(effect, 0)
	if err != nil {
		t.Fatalf("apply(%q): %v", effect.Detail, err)
	}
}

// held counts the possessions of one name.
func held(gen *Generator, item string) int {
	found := 0

	for _, possession := range gen.char.State.Stash {
		if possession.Item == item {
			found++
		}
	}

	return found
}

// TestANamedBenefitReachesTheStash is the finding that opened issue #42:
// a benefit row naming a thing put it in the log and not in the record, so
// there was nothing for a later result to take away.
func TestANamedBenefitReachesTheStash(t *testing.T) {
	t.Parallel()

	gen := stashEngine(t, 29)

	mustApply(t, gen, career.Effect{
		Kind: career.EffectStash, Item: "a company share", Dice: "2d6x100000",
		Count: 1, Detail: "a company share",
	})

	if held(gen, "a company share") != 1 {
		t.Fatalf("the share is not in the stash: %v", gen.char.State.Stash)
	}

	share := gen.char.State.Stash[0]

	if share.Value < 200000 || share.Value > 1200000 {
		t.Errorf("a 2d6 x 100,000 share is worth %d", share.Value)
	}

	if share.Career != "Colonist" {
		t.Errorf("the share does not name where it came from: %q", share.Career)
	}
}

// TestLosingSharesLeavesEverythingElse is the fourteen results that read
// "lose any Company Shares": those and nothing else.
func TestLosingSharesLeavesEverythingElse(t *testing.T) {
	t.Parallel()

	gen := stashEngine(t, 31)

	gen.char.State.Stash = []Possession{
		{Item: "a company share", Value: 400000},
		{Item: "a weapon of the character's choice"},
		{Item: "a company share", Value: 700000},
	}

	mustApply(t, gen, career.Effect{
		Kind: career.EffectStash, Item: "a company share", Lose: true,
		Detail: "lose any company shares held",
	})

	if held(gen, "a company share") != 0 {
		t.Error("a share survived")
	}

	if held(gen, "a weapon of the character's choice") != 1 {
		t.Error("the weapon went with the shares")
	}
}

// TestLosingWhatIsNotHeldIsStillRecorded. The result fired, and what it did
// is part of what happened -- a silent no-op would read as though the
// mishap had never come up.
func TestLosingWhatIsNotHeldIsStillRecorded(t *testing.T) {
	t.Parallel()

	gen := stashEngine(t, 37)

	mustApply(t, gen, career.Effect{
		Kind: career.EffectStash, Item: "a company share", Lose: true,
		Detail: "lose any company shares held",
	})

	if !strings.Contains(lastDetail(t, gen), "there were none") {
		t.Errorf("the record does not say the stash was empty: %q", lastDetail(t, gen))
	}
}

// TestEachOfSeveralIsValuedOnItsOwnThrow is p. 155's "Six Pieces of Art,
// 2D6 x Cr10000 each". Each, not one throw shared six ways.
func TestEachOfSeveralIsValuedOnItsOwnThrow(t *testing.T) {
	t.Parallel()

	gen := stashEngine(t, 41)

	mustApply(t, gen, career.Effect{
		Kind: career.EffectStash, Item: "a piece of art", Dice: "2d6x10000",
		Count: 6, Detail: "six pieces of art",
	})

	if held(gen, "a piece of art") != 6 {
		t.Fatalf("%d pieces of art, want 6", held(gen, "a piece of art"))
	}

	same := true

	for _, piece := range gen.char.State.Stash {
		if piece.Value != gen.char.State.Stash[0].Value {
			same = false
		}
	}

	if same {
		t.Error("six pieces of art all came to the same number")
	}
}

// TestARolledCountOfShares is Scientist's "1D6 Company Shares": how many is
// thrown, and the book prints no value for them.
func TestARolledCountOfShares(t *testing.T) {
	t.Parallel()

	gen := stashEngine(t, 43)

	mustApply(t, gen, career.Effect{
		Kind: career.EffectStash, Item: "a company share", CountDice: "1d6",
		Detail: "1d6 company shares",
	})

	count := held(gen, "a company share")
	if count < 1 || count > 6 {
		t.Errorf("1d6 shares granted %d", count)
	}

	for _, share := range gen.char.State.Stash {
		if share.Value != 0 {
			t.Errorf("a share the book priced at nothing is worth %d", share.Value)
		}
	}
}

// TestAnUnreadableStashExpressionIsAnError holds both dice paths a
// possession has: what it is worth, and how many of it there are. Either
// stops generation rather than quietly granting nothing.
func TestAnUnreadableStashExpressionIsAnError(t *testing.T) {
	t.Parallel()

	for _, effect := range []career.Effect{
		{Kind: career.EffectStash, Item: "a share", Dice: "a fortune", Count: 1, Detail: "worth"},
		{Kind: career.EffectStash, Item: "a share", CountDice: "several", Detail: "how many"},
	} {
		gen := stashEngine(t, 53)

		err := gen.apply(effect, 0)
		if err == nil {
			t.Errorf("%q: an unreadable expression was granted anyway", effect.Detail)
		}
	}
}

// TestAGroupAppliesEachEffect is the vocabulary a benefit row needs: one
// Other effect, and rows that do two things.
func TestAGroupAppliesEachEffect(t *testing.T) {
	t.Parallel()

	gen := stashEngine(t, 47)

	mustApply(t, gen, career.Effect{
		Kind:   career.EffectGroup,
		Detail: "a producer credit, and a yearly payment",
		Group: []career.Effect{
			{Kind: career.EffectStash, Item: "a producer credit", Count: 1, Detail: "a credit"},
			{Kind: career.EffectCharacteristic, Characteristic: "CHA", Delta: 1, Detail: "+1 CHA"},
		},
	})

	if held(gen, "a producer credit") != 1 {
		t.Error("the first effect of the group did not apply")
	}

	if gen.char.State.Characteristics.Get(CHA) == 0 {
		t.Error("the second effect of the group did not apply")
	}
}
