package chargen

import (
	"fmt"
	"math/rand/v2"
)

// pointCareer is the choice point BatchPolicy answers differently from
// Policy, named because two files spell it.
const (
	pointCareer     = "career"
	pointAssignment = "assignment"
)

// batchDrawn are the choice points a batch draws rather than taking the
// first option. Each chooses among options the book offers on equal terms
// -- which table to roll on, which path to take, which specialty to raise
// -- so a draw cannot make a worse character, only a different one.
//
// The rest stay Policy because Policy's answer is argued and a draw's
// would not be: a loss comes out of the highest score, an aging crisis is
// paid for, benefits and pools follow their own rows in POLICY.md, and
// characteristics are assigned in roll order because the rolls are already
// random.
var batchDrawn = map[string]bool{
	pointCareer:            true,
	pointAssignment:        true,
	"other_assignment":     true,
	"skill_table":          true,
	pointSkillSpecialty:    true,
	"any_skill":            true,
	"raise_held_skill":     true,
	"youth_path":           true,
	"teenage_path":         true,
	"degree_field":         true,
	"washout_skill":        true,
	"background_skill":     true,
	"background_specialty": true,
}

// BatchPolicy is the auto policy for a batch: Policy, except at the choice
// points in batchDrawn.
//
// Policy takes the first option the book prints, which is a rule a reader
// can check against the page rather than a preference. Applied to the
// career list a hundred times it is also the reason a hundred characters
// drew from eight careers of thirty-four, fifty-five of them Adventurer,
// with nobody holding Engineer or Gunner at any level: the list is built in
// the book's order, and the book's order is alphabetical.
//
// One character generated that way is a character. A hundred are not a
// population, and `batch` is the command a population is asked of. So the
// career was drawn from the eligible list from 0.3.0. That left a hundred
// characters in thirty-four careers who all rolled on Personal Development
// every term and all took Youth Path 1, so from 0.5.0 the draw reaches
// every choice point where the options stand level.
//
// The draw is seeded from the record's own seed, so a member is a function
// of the seed it records rather than of its position in the batch, and
// `replay` on its record regenerates it alone. It is not the character
// `new --seed N+i` makes, which takes the first option at every point.
type BatchPolicy struct {
	draw *rand.Rand
}

// careerDrawStream keeps the policy's own randomness away from the dice.
// The engine's determinism contract is a seed plus recorded choices, and a
// decider that pulled from the dice stream would shift every throw after
// the first career choice.
const careerDrawStream = 0x63617265_65727370

// NewBatchPolicy is the policy for one member of a batch, drawing from the
// seed that member was generated with.
func NewBatchPolicy(seed uint64) BatchPolicy {
	return BatchPolicy{draw: rand.New(rand.NewPCG(seed, careerDrawStream))}
}

// Kind implements [Decider]. A batch character's choices were still made by
// a policy; which policy is what policy_version says.
func (BatchPolicy) Kind() DeciderKind { return DeciderPolicy }

// Choose draws at a point in batchDrawn and is Policy everywhere else.
func (b BatchPolicy) Choose(c Choice) (int, error) {
	if len(c.Options) == 0 {
		return 0, fmt.Errorf("%w: %s", ErrNoOptions, c.Point)
	}

	if !batchDrawn[c.Point] {
		return Policy{}.Choose(c)
	}

	return b.draw.IntN(len(c.Options)), nil
}
