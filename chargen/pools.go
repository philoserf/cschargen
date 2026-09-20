package chargen

import (
	"slices"

	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/dice"
)

// Pool is a modifier the character spends themselves, a little at a time:
// "gain a modifier of +6 which must be split up into increments of not more
// than +2 to use on any Survival or Advancement rolls until it is depleted"
// (Journalist event 66).
//
// It is not a PendingModifier because nothing about it is pending: the
// character decides how much of it to spend at each throw, and the decision
// is a choice point like any other.
type Pool struct {
	Remaining int      `json:"remaining"`
	PerThrow  int      `json:"perThrow"`
	Detail    string   `json:"detail"`
	Spendable []string `json:"spendable"`

	// WhileInThisCareer ends the pool when the career does, which is what
	// "over the remainder of your career in this service" means. The other
	// of the two runs "until it is depleted", and outlives the career.
	WhileInThisCareer bool `json:"whileInThisCareer,omitempty"`
}

// grantPool records a pool the character may spend.
func (g *Generator) grantPool(effect career.Effect, cause int) {
	g.pools = append(g.pools, Pool{
		Remaining:         effect.Modifier,
		PerThrow:          effect.Uses,
		Detail:            effect.Detail,
		Spendable:         effect.Spendable,
		WhileInThisCareer: effect.WhileInThisCareer,
	})

	g.consequence(ConsequenceModifier, cause, effect.Detail, "")
}

// spendPools offers each open pool to the decider before a throw and
// returns what was spent.
//
// The options run from the largest increment down, so a policy character --
// which takes the first option offered -- spends the pool rather than
// carrying it to the end of a career and wasting it. That is a decider's
// habit rather than a rule: the book says "at any time", and a player
// answering the question themselves may say nothing.
func (g *Generator) spendPools(applies string, cause int) ([]dice.Mod, error) {
	mods := make([]dice.Mod, 0, len(g.pools))

	for i := range g.pools {
		if g.pools[i].Remaining <= 0 || !slices.Contains(g.pools[i].Spendable, applies) {
			continue
		}

		most := min(g.pools[i].PerThrow, g.pools[i].Remaining)

		options := make([]string, 0, most+1)
		for spend := most; spend >= 0; spend-- {
			options = append(options, "+"+itoa(spend))
		}

		chosen, err := g.choose(Choice{
			Point:   "spend_pool",
			Prompt:  "How much of " + g.pools[i].Detail + " to spend on this throw?",
			Options: options,
			Cite:    poolCite,
		})
		if err != nil {
			return nil, err
		}

		spend := most - chosen
		if spend == 0 {
			continue
		}

		g.pools[i].Remaining -= spend

		mods = append(mods, dice.Mod{Name: g.pools[i].Detail, Value: spend})
		g.consequence(ConsequenceModifier, cause,
			"spent +"+itoa(spend)+" of "+g.pools[i].Detail+", leaving "+
				itoa(g.pools[i].Remaining), "")
	}

	g.dropSpentPools()

	return mods, nil
}

// poolCite is where the two results that grant a pool are printed.
const poolCite = "pp. 188, 218"

// dropSpentPools forgets the pools with nothing left in them.
func (g *Generator) dropSpentPools() {
	kept := make([]Pool, 0, len(g.pools))

	for _, held := range g.pools {
		if held.Remaining > 0 {
			kept = append(kept, held)
		}
	}

	g.pools = kept
}

// dropCareerPools ends the pools granted "over the remainder of your career
// in this service" when that service ends.
func (g *Generator) dropCareerPools() {
	kept := make([]Pool, 0, len(g.pools))

	for _, held := range g.pools {
		if !held.WhileInThisCareer {
			kept = append(kept, held)
		}
	}

	g.pools = kept
}

// onFailure records effects that fire if a named throw fails, which nine
// results attach to an advancement roll they have just improved.
func (g *Generator) grantOnFailure(effect career.Effect, cause int) {
	g.onFailure = append(g.onFailure, effect)
	g.consequence(ConsequenceModifier, cause, effect.Detail, "")
}

// takeOnFailure collects and clears the effects waiting on a named throw.
// They are taken whether the throw passed or failed: a hook that outlived
// the throw it was attached to would fire against the next one.
//
//nolint:unparam // the throw is the page's; nine results attach to an advancement roll and none to another
func (g *Generator) takeOnFailure(applies string) []career.Effect {
	var (
		fired []career.Effect
		kept  []career.Effect
	)

	for _, pending := range g.onFailure {
		if pending.Applies == applies {
			fired = append(fired, pending)

			continue
		}

		kept = append(kept, pending)
	}

	g.onFailure = kept

	return fired
}

// takeAutoFailure reports whether a named throw has already been decided
// against the character, and spends it if so.
func (g *Generator) takeAutoFailure(applies string) bool {
	for i, pending := range g.autoFailure {
		if pending == applies {
			g.autoFailure = append(g.autoFailure[:i], g.autoFailure[i+1:]...)

			return true
		}
	}

	return false
}
