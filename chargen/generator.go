package chargen

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/philoserf/cschargen/dice"
)

// The page the characteristic rule is printed on, and the identifier of its
// choice point. Named because both appear in several events and a typo in
// one would read as a second, different rule.
const (
	characteristicCite         = "p. 13"
	pointAssignCharacteristics = "assign_characteristics"
)

// Options configure a generation run. Everything here that a seed does not
// capture is recorded in the character's provenance, because replay needs
// it.
type Options struct {
	Seed          uint64
	Decider       Decider
	EngineVersion string
	PolicyVersion string
	SettingData   SettingData
	Inputs        Inputs
}

// Generator walks the twenty steps of character creation. One run produces
// one character; a Generator is not reusable.
type Generator struct {
	dice    *dice.Dice
	log     *Log
	decider Decider
	char    *Character
}

// New returns a Generator ready to run.
func New(opts Options) *Generator {
	return &Generator{
		dice:    dice.New(opts.Seed),
		log:     &Log{},
		decider: opts.Decider,
		char: &Character{
			Provenance: Provenance{
				SchemaVersion: SchemaVersion,
				Ruleset:       Ruleset,
				EngineVersion: opts.EngineVersion,
				PolicyVersion: opts.PolicyVersion,
				RNG:           RNG{Algorithm: "math/rand/v2 PCG", Seed: opts.Seed},
				SettingData:   opts.SettingData,
				Inputs:        opts.Inputs,
			},
		},
	}
}

// Run walks the steps this milestone implements and returns the record.
//
// Milestone 1 implements Step 2 and Steps 9-19; Steps 1 and 3-8 are stubbed
// (docs/MILESTONE-1.md). The order here is the book's, so that adding a
// step is adding a line rather than rearranging one.
func (g *Generator) Run() (*Character, error) {
	err := g.rollCharacteristics()
	if err != nil {
		return nil, err
	}

	g.char.Events = g.log.Events()

	return g.char, nil
}

// rollCharacteristics is Step 2 (p. 39), whose rule is on p. 13: "Roll 3d6.
// Drop the score on the lowest of the three dice and add the remaining two
// scores to determine the character's six characteristics. It is best for
// the player to roll six times and then apply the numbers generated to the
// characteristics as they see fit."
//
// The six rolls happen first and the assignment is a separate choice, in
// that order, because that is the order the sentence gives them -- and
// because a player who assigned as they rolled would be making a different
// decision, with less information, at every step.
func (g *Generator) rollCharacteristics() error {
	g.log.Step("Step 2: Roll Characteristics", "p. 39")

	rolled := make([]int, 0, len(CharacteristicOrder))
	for range CharacteristicOrder {
		roll := g.dice.Characteristic()
		g.log.Roll(roll, characteristicCite)

		rolled = append(rolled, roll.Total)
	}

	return g.assignCharacteristics(rolled)
}

// assignCharacteristics puts each characteristic, in printed order, to the
// decider against the scores still unassigned. Six successive choices
// rather than one permutation: a front end can present them one at a time,
// and Nth/Of tells a player answering the fifth that it is the fifth.
func (g *Generator) assignCharacteristics(rolled []int) error {
	remaining := rolled

	for nth, which := range CharacteristicOrder {
		options := make([]string, len(remaining))
		for i, score := range remaining {
			options[i] = strconv.Itoa(score)
		}

		prompt := "Assign a rolled score to " + which.String()

		chosen, err := g.decider.Choose(Choice{
			Point:   pointAssignCharacteristics,
			Prompt:  prompt,
			Options: options,
			Cite:    characteristicCite,
			Nth:     nth + 1,
			Of:      len(CharacteristicOrder),
		})
		if err != nil {
			return fmt.Errorf("assigning %s: %w", which, err)
		}

		if chosen < 0 || chosen >= len(remaining) {
			return ErrChoiceOutOfRange
		}

		score := remaining[chosen]

		remaining = slices.Delete(remaining, chosen, chosen+1)

		g.char.Characteristics.Set(which, score)

		cause := g.log.Choice(ChoiceEvent{
			Decider: g.decider.Kind(),
			Point:   pointAssignCharacteristics,
			Prompt:  prompt,
			Options: options,
			Chosen:  chosen,
			Cite:    characteristicCite,
		})

		g.log.Consequence(ConsequenceEvent{
			Kind:           ConsequenceCharacteristic,
			Cause:          cause,
			Detail:         which.String() + " " + strconv.Itoa(score),
			Characteristic: which.String(),
			Delta:          score,
			Cite:           characteristicCite,
		})
	}

	return nil
}
