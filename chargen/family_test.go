package chargen_test

import (
	"testing"

	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/setting"
)

// TestTheThreeBirthSituationsAgree holds the engine's three result strings
// against the validator's. A world in a data file may force one of them,
// and a fourth word would generate a household with no parents in it.
func TestTheThreeBirthSituationsAgree(t *testing.T) {
	t.Parallel()

	engine := []string{
		chargen.SituationCommunal,
		chargen.SituationSameSex,
		chargen.SituationHeterosexual,
	}

	if len(setting.BirthSituations) != len(engine) {
		t.Fatalf("the validator knows %d birth situations and the engine has %d",
			len(setting.BirthSituations), len(engine))
	}

	for _, situation := range engine {
		if !setting.BirthSituations[situation] {
			t.Errorf("the engine can generate %q and the validator would reject it", situation)
		}
	}
}

// TestEveryCharacterHasAFamily. Step 5 runs by default, so a character who
// did not ask to skip it has parents.
func TestEveryCharacterHasAFamily(t *testing.T) {
	t.Parallel()

	situations := map[string]int{}

	for seed := range uint64(sample) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = -1

		character := generate(t, opts)

		family := character.State.Family
		if family == nil {
			t.Fatalf("seed %d has no family", seed)
		}

		situations[family.Situation]++

		checkParents(t, seed, character)
	}

	// All three results of p. 58 are reachable across the sample, which the
	// data's own modifiers help along: one world adds 30 to the throw and
	// another subtracts 20.
	for _, situation := range []string{
		chargen.SituationCommunal,
		chargen.SituationSameSex,
		chargen.SituationHeterosexual,
	} {
		if situations[situation] == 0 {
			t.Errorf("no seed in %d was born to %s", sample, situation)
		}
	}
}

// checkParents holds one character's parents against p. 58: one tie per
// parental age, all of them from the family, all of them Allies at 125.
func checkParents(t *testing.T, seed uint64, character *chargen.Character) {
	t.Helper()

	family := character.State.Family
	if len(family.ParentAges) == 0 {
		t.Errorf("seed %d was born to %s and has no parents", seed, family.Situation)
	}

	parents := 0

	for _, tie := range character.State.Ties {
		if tie.Role != chargen.RoleParent {
			continue
		}

		parents++

		if tie.Origin != chargen.FamilyOrigin {
			t.Errorf("seed %d: a parent's origin is %q", seed, tie.Origin)
		}

		// ERRATA E-22: Step 5 grants parents as Allies at 125.
		if tie.Rating != 125 {
			t.Errorf("seed %d: a parent is at %d, want 125 (p. 58)", seed, tie.Rating)
		}
	}

	if parents != len(family.ParentAges) {
		t.Errorf("seed %d: %d parent ages and %d parent ties",
			seed, len(family.ParentAges), parents)
	}
}

// TestSkippingTheFamily. p. 57: "Skipping this step can speed up character
// generation, but it may also deprive the character of potential Allies and
// Contacts."
func TestSkippingTheFamily(t *testing.T) {
	t.Parallel()

	opts := options(t, 5)

	opts.Inputs.TermLimit = -1
	opts.Inputs.SkipFamily = true

	character := generate(t, opts)

	if character.State.Family != nil {
		t.Error("the family step ran although it was skipped")
	}

	for _, tie := range character.State.Ties {
		if tie.Origin == chargen.FamilyOrigin {
			t.Errorf("a family tie survived the step being skipped: %s", tie.Detail)
		}
	}
}

// TestBirthOrderConstrainsTheSiblings is p. 60's two re-roll rules, which
// the engine applies as the table's own redirects rather than by re-rolling:
// "If you have already determined that the character is the firstborn child,
// re-roll any result which would make the sibling older."
func TestBirthOrderConstrainsTheSiblings(t *testing.T) {
	t.Parallel()

	firstborn, lastborn := 0, 0

	for seed := range uint64(200) {
		opts := options(t, seed)

		opts.Inputs.TermLimit = -1

		character := generate(t, opts)
		family := character.State.Family

		if !family.Firstborn && !family.Lastborn {
			continue
		}

		checkSiblingOrder(t, seed, character)

		if family.Firstborn {
			firstborn++
		}

		if family.Lastborn {
			lastborn++
		}
	}

	if firstborn == 0 || lastborn == 0 {
		t.Errorf("in 200 seeds, %d firstborn and %d last children; both ends of the "+
			"parental age die should come up", firstborn, lastborn)
	}
}

// checkSiblingOrder is p. 60's two constraints: a firstborn has no older
// siblings and a last child has no younger ones. A twin is neither, and a
// sibling who died at childbirth has no age difference at all.
func checkSiblingOrder(t *testing.T, seed uint64, character *chargen.Character) {
	t.Helper()

	family := character.State.Family

	for _, tie := range character.State.Ties {
		if tie.Role != chargen.RoleSibling {
			continue
		}

		if family.Firstborn && olderThanTheCharacter(tie.Detail) {
			t.Errorf("seed %d: a firstborn has an older sibling (%s)", seed, tie.Detail)
		}

		if family.Lastborn && youngerThanTheCharacter(tie.Detail) {
			t.Errorf("seed %d: a last child has a younger sibling (%s)", seed, tie.Detail)
		}
	}
}

func olderThanTheCharacter(detail string) bool {
	return len(detail) >= 5 && detail[:5] == "older"
}

func youngerThanTheCharacter(detail string) bool {
	return len(detail) >= 7 && detail[:7] == "younger"
}
