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

// BatchPolicy is the auto policy for a batch: Policy at every choice point
// but one.
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
// career is drawn from the eligible list instead, and every other choice
// point is Policy unchanged -- because varying the rest would be a
// different policy rather than the same one applied to a crowd.
//
// The draw is seeded from the record's own seed, so `batch --seed N` member
// i and `new --seed N+i --policy-batch` are the same character: the PRD's
// promise that any member can be regenerated alone survives.
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

// Choose takes the first option, except at the career.
func (b BatchPolicy) Choose(c Choice) (int, error) {
	if len(c.Options) == 0 {
		return 0, fmt.Errorf("%w: %s", ErrNoOptions, c.Point)
	}

	if c.Point != pointCareer && c.Point != pointAssignment {
		return Policy{}.Choose(c)
	}

	return b.draw.IntN(len(c.Options)), nil
}
