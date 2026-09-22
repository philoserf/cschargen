package chargen

import "github.com/philoserf/cschargen/setting"

// The genetics of pp. 62-65.
//
// An engineered person is a purebred of some generation, a hybrid with a
// baseline human, or a compound of two engineered species. The status is a
// 1d6 four of the book's five species share; the hybrid and compound tables
// it sends a character to are one per species, and every row of a compound
// table names another species by name. So the tables are data and only the
// shape is here.
//
// The one thing that is mechanical is which method the characteristics are
// rolled with: a hybrid may come out "otherwise human", and one row of one
// table says the opposite -- appears human, rolls as the species.

// Genetics is what Step 5 established about an engineered character's
// parentage.
type Genetics struct {
	// Kind is "purebred", "hybrid" or "compound".
	Kind string `json:"kind"`

	// Generation is which one, for a purebred.
	Generation int `json:"generation,omitempty"`

	// Detail is the row that was rolled, in the data file's words.
	Detail string `json:"detail"`

	// RolledAsHuman records that the outcome sent the character back to
	// the human characteristic method.
	RolledAsHuman bool `json:"rolledAsHuman,omitempty"`
}

// laterGenerations is what "third, fourth or later generation ... The
// player chooses which generation" offers (p. 63). The engine offers three
// and stops, because the page gives no ceiling and a list has to end.
var laterGenerations = []string{"third", "fourth", "fifth"}

// firstChosenGeneration is the first the list above offers.
const firstChosenGeneration = 3

// determineGenetics is the first half of Step 5 for an engineered
// character, and runs before the characteristics because one of its
// outcomes decides how they are rolled.
func (g *Generator) determineGenetics() error {
	if g.species == nil || g.species.Kind != setting.KindEngineered {
		return nil
	}

	table := g.species.GeneticStatusTable()

	roll := g.dice.D6()
	cause := g.log.Roll(roll, geneticsCite)

	status := table[roll.Total-1]
	record := &Genetics{
		Kind: status.Kind, Generation: status.Generation, Detail: status.Detail,
	}

	g.char.State.Genetics = record

	if status.Kind == setting.Purebred && status.Generation == 0 {
		chosen, err := g.choose(Choice{
			Point:   "generation",
			Prompt:  "Choose which generation",
			Options: laterGenerations,
			Cite:    geneticsCite,
		})
		if err != nil {
			return err
		}

		record.Generation = chosen + firstChosenGeneration

		record.Detail += " -- the " + laterGenerations[chosen]
	}

	g.consequence(ConsequenceSpecies, cause, "genetics: "+record.Detail, "")

	return g.geneticOutcome(record)
}

// geneticsCite is where the genetic status table is printed. The hybrid and
// compound tables are one per species and the engine cannot name their
// pages, for the reason it cannot name the species.
const geneticsCite = "pp. 62-65"

// geneticOutcome rolls on the hybrid or compound table the status sent the
// character to.
func (g *Generator) geneticOutcome(record *Genetics) error {
	table := g.outcomeTable(record.Kind)
	if len(table) == 0 {
		return nil
	}

	roll := g.dice.D6()
	throw := g.log.Roll(roll, geneticsCite)

	outcome := table[roll.Total-1]

	// "The player chooses from the above results" is the last row of every
	// compound table (pp. 64-65).
	if outcome.Choose {
		chosen, err := g.chooseOutcome(table)
		if err != nil {
			return err
		}

		outcome = table[chosen]
	}

	record.Detail += ": " + outcome.Detail

	if outcome.RollAs == humanMethod {
		record.RolledAsHuman = true

		g.rollAsHuman = true
	}

	g.consequence(ConsequenceSpecies, throw, outcome.Detail, "")

	for name, delta := range outcome.Adjust {
		g.adjust(name, delta, outcome.Detail, throw)
	}

	return nil
}

// humanMethod is what a hybrid row means by "Roll characteristics as a
// human".
const humanMethod = "human"

// outcomeTable is the hybrid or compound table, or nothing for a purebred.
func (g *Generator) outcomeTable(kind string) []setting.GeneticOutcome {
	if g.species.Genetics == nil {
		return nil
	}

	switch kind {
	case setting.Hybrid:
		return g.species.Genetics.Hybrid
	case setting.Compound:
		return g.species.Genetics.Compound
	}

	return nil
}

// chooseOutcome is the last row of a compound table, which hands the
// decision back rather than being one.
func (g *Generator) chooseOutcome(table []setting.GeneticOutcome) (int, error) {
	labels := make([]string, 0, len(table))
	indexes := make([]int, 0, len(table))

	for i, outcome := range table {
		if outcome.Choose {
			continue
		}

		labels = append(labels, outcome.Detail)
		indexes = append(indexes, i)
	}

	if len(labels) == 0 {
		return 0, ErrNoOptions
	}

	chosen, err := g.choose(Choice{
		Point:   "genetic_outcome",
		Prompt:  "Choose a compound outcome",
		Options: labels,
		Cite:    geneticsCite,
	})
	if err != nil {
		return 0, err
	}

	return indexes[chosen], nil
}
