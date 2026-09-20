package chargen

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/career"
)

// TestSpendingAPoolTakesTheLargestIncrement is the decider's habit rather
// than a rule: the options run from the largest increment down, so a policy
// character spends the pool rather than carrying it to the end of a career.
func TestSpendingAPoolTakesTheLargestIncrement(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 0)

	gen.pools = []Pool{{
		Remaining: 6, PerThrow: 2, Detail: "a +6 modifier",
		Spendable: []string{survivalThrow, advancementThrow},
	}}

	mods, err := gen.spendPools(survivalThrow, 0)
	if err != nil {
		t.Fatalf("spendPools: %v", err)
	}

	if len(mods) != 1 || mods[0].Value != 2 {
		t.Fatalf("spent %v, want one increment of +2", mods)
	}

	if gen.pools[0].Remaining != 4 {
		t.Errorf("%d left in the pool, want 4", gen.pools[0].Remaining)
	}
}

// TestAPoolIsNotSpentOnAThrowItDoesNotName. Both results that grant one
// name survival and advancement; an enlistment throw is not either.
func TestAPoolIsNotSpentOnAThrowItDoesNotName(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 0)

	gen.pools = []Pool{{
		Remaining: 6, PerThrow: 2, Detail: "a +6 modifier",
		Spendable: []string{advancementThrow},
	}}

	mods, err := gen.spendPools(survivalThrow, 0)
	if err != nil {
		t.Fatalf("spendPools: %v", err)
	}

	if len(mods) != 0 {
		t.Errorf("a pool for advancement was spent on a survival roll: %v", mods)
	}

	if gen.pools[0].Remaining != 6 {
		t.Errorf("%d left in the pool, want 6", gen.pools[0].Remaining)
	}
}

// TestAnEmptyPoolIsForgotten, and one granted "over the remainder of your
// career in this service" ends when the service does.
func TestAnEmptyPoolIsForgotten(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 0)

	gen.pools = []Pool{{
		Remaining: 2, PerThrow: 3, Detail: "nearly gone",
		Spendable: []string{survivalThrow}, WhileInThisCareer: true,
	}}

	_, err := gen.spendPools(survivalThrow, 0)
	if err != nil {
		t.Fatalf("spendPools: %v", err)
	}

	if len(gen.pools) != 0 {
		t.Errorf("a spent pool was kept: %v", gen.pools)
	}

	gen.pools = []Pool{{Remaining: 4, PerThrow: 2, WhileInThisCareer: true}}
	gen.dropCareerPools()

	if len(gen.pools) != 0 {
		t.Errorf("a pool tied to the career outlived it: %v", gen.pools)
	}
}

// TestARefusedPoolChoiceComesBack: a decider that answers nothing stops
// generation rather than spending for the player.
func TestARefusedPoolChoiceComesBack(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 0)

	gen.decider = refusingDecider{}
	gen.pools = []Pool{{
		Remaining: 6, PerThrow: 2, Detail: "a +6 modifier",
		Spendable: []string{survivalThrow},
	}}

	_, err := gen.spendPools(survivalThrow, 0)
	if err == nil {
		t.Error("a refused choice did not come back")
	}
}

// TestAFailedAdvancementFiresWhatWaitedOnIt is the nine results that read
// "if you fail that Advancement roll, you lose the Ally and take a -2 DM to
// your next Advancement roll".
func TestAFailedAdvancementFiresWhatWaitedOnIt(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 0)

	gen.char.State.Ties = []Tie{{Kind: string(career.Ally), Origin: "Colonist", Rating: 150}}
	gen.onFailure = []career.Effect{{
		Kind: career.EffectOnFailure, Applies: advancementThrow, Detail: "the Ally is lost",
		Failure: []career.Effect{{
			Kind: career.EffectLoseTie, Order: []career.Relationship{career.Ally},
			Detail: "lose the Ally",
		}},
	}}

	// An automatic failure makes the throw without rolling it, which is the
	// same door the hook comes through.
	gen.autoFailure = []string{advancementThrow}

	err := gen.advance(gen.assignment)
	if err != nil {
		t.Fatalf("advance: %v", err)
	}

	if len(gen.char.State.Ties) != 0 {
		t.Errorf("the Ally survived the failure: %v", gen.char.State.Ties)
	}

	if len(gen.onFailure) != 0 {
		t.Error("the hook outlived the throw it waited on")
	}
}

// TestAnAutomaticFailureIsRecorded: the throw is not made, and the record
// says why rather than showing a throw that never happened.
func TestAnAutomaticFailureIsRecorded(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 0)

	gen.autoFailure = []string{advancementThrow}

	err := gen.advance(gen.assignment)
	if err != nil {
		t.Fatalf("advance: %v", err)
	}

	if gen.rank != 0 {
		t.Errorf("rank is %d after an automatic failure, want 0", gen.rank)
	}

	if !strings.Contains(lastDetail(t, gen), "fails automatically") {
		t.Errorf("the record does not say the throw was decided: %q", lastDetail(t, gen))
	}
}

// TestAPlayerMaySpendNothing. The book says "at any time", which includes
// not now: a decider that declines the increment leaves the pool whole.
func TestAPlayerMaySpendNothing(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 0)

	gen.decider = lastOptionDecider{}
	gen.pools = []Pool{{
		Remaining: 6, PerThrow: 2, Detail: "a +6 modifier",
		Spendable: []string{survivalThrow},
	}}

	mods, err := gen.spendPools(survivalThrow, 0)
	if err != nil {
		t.Fatalf("spendPools: %v", err)
	}

	if len(mods) != 0 {
		t.Errorf("declining the increment still spent %v", mods)
	}

	if gen.pools[0].Remaining != 6 {
		t.Errorf("%d left in the pool, want 6", gen.pools[0].Remaining)
	}
}

// TestAHookWaitsForItsOwnThrow. Nine results attach to an advancement
// roll; a hook is not fired by a throw of another name.
func TestAHookWaitsForItsOwnThrow(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 0)

	gen.onFailure = []career.Effect{
		{Kind: career.EffectOnFailure, Applies: survivalThrow, Detail: "not this one"},
	}

	if fired := gen.takeOnFailure(advancementThrow); len(fired) != 0 {
		t.Errorf("an advancement throw fired a survival hook: %v", fired)
	}

	if len(gen.onFailure) != 1 {
		t.Error("the hook was consumed by a throw it did not name")
	}
}

// TestAnAutomaticFailureIsGrantedThroughApply covers the door the tables
// come through, as against the field the tests above set directly.
func TestAnAutomaticFailureIsGrantedThroughApply(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 0)

	err := gen.apply(career.Effect{
		Kind: career.EffectAutoFailure, Applies: advancementThrow,
		Detail: "suspended: the next advancement roll fails automatically",
	}, 0)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	if !gen.takeAutoFailure(advancementThrow) {
		t.Error("the automatic failure was not recorded")
	}
}

// TestARefusalDuringAThrowComesBack holds the two throws a pool is spent
// on: a player who closes the program mid-term stops generation rather
// than having the increment chosen for them.
func TestARefusalDuringAThrowComesBack(t *testing.T) {
	t.Parallel()

	open := []Pool{{
		Remaining: 6, PerThrow: 2, Detail: "a +6 modifier",
		Spendable: []string{survivalThrow, advancementThrow},
	}}

	gen := rankedEngine(t, 0)

	gen.decider = refusingDecider{}
	gen.pools = open

	_, err := gen.rollSurvival(gen.assignment)
	if err == nil {
		t.Error("a refusal during the survival roll did not come back")
	}

	gen = rankedEngine(t, 0)

	gen.decider = refusingDecider{}
	gen.pools = open
	gen.commissioned = true

	err = gen.advance(gen.assignment)
	if err == nil {
		t.Error("a refusal during the advancement roll did not come back")
	}

	// And the commission offer, which comes before both and is the first
	// question Step 14 asks.
	gen = rankedEngine(t, 0)

	navy := career.NationalNavy()

	gen.career = &navy
	gen.assignment = navy.Assignments[0]
	gen.decider = refusingDecider{}

	err = gen.advance(gen.assignment)
	if err == nil {
		t.Error("a refusal at the commission offer did not come back")
	}
}

// TestASubTableWithAGapSaysSo. TestEverySubTableCoversTheDie keeps the
// corpus from having one, so this is the guard behind that gate: a table
// that covers nothing the die rolled records the gap rather than doing
// nothing quietly.
func TestASubTableWithAGapSaysSo(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 0)

	err := gen.rollSubTable(career.Effect{
		Kind:   career.EffectSubTable,
		Detail: "a table with no rows at all",
	}, 0)
	if err != nil {
		t.Fatalf("rollSubTable: %v", err)
	}

	if !strings.Contains(lastDetail(t, gen), "no row covers") {
		t.Errorf("the record does not name the gap: %q", lastDetail(t, gen))
	}
}

// TestASubTableRollsRatherThanChooses is what separates it from a choice:
// the character has no say in which row comes up, so a decider that
// refuses every question does not stop it.
func TestASubTableRollsRatherThanChooses(t *testing.T) {
	t.Parallel()

	gen := rankedEngine(t, 0)

	gen.decider = refusingDecider{}

	err := gen.rollSubTable(career.Effect{
		Kind:   career.EffectSubTable,
		Detail: "the pay",
		Sub: []career.SubRow{
			{From: 1, To: 3, Summary: "little", Effects: []career.Effect{
				{Kind: career.EffectCredits, Dice: "200", Detail: "200 credits"},
			}},
			{From: 4, To: 6, Summary: "more", Effects: []career.Effect{
				{Kind: career.EffectCredits, Dice: "500", Detail: "500 credits"},
			}},
		},
	}, 0)
	if err != nil {
		t.Fatalf("rollSubTable: %v", err)
	}

	if gen.char.State.Credits != 200 && gen.char.State.Credits != 500 {
		t.Errorf("the table paid %d, which is neither row", gen.char.State.Credits)
	}
}
