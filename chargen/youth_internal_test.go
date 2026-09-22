package chargen

import (
	"testing"

	"github.com/philoserf/cschargen/career"
)

// TestARefusedYouthPathEndsGeneration. Choosing a path goes through the
// Decider like every other choice, so an abandoned interactive session
// comes back as an error rather than a default.
func TestARefusedYouthPathEndsGeneration(t *testing.T) {
	t.Parallel()

	gen := engine(t, 17)

	gen.decider = refusingDecider{}

	// Everything at 15 opens all five paths, so a choice is offered.
	for _, which := range CharacteristicOrder {
		gen.char.State.Characteristics.Set(which, 15)
	}

	_, err := gen.chooseYouthPath("ages 4-8")
	if err == nil {
		t.Fatal("a refused choice did not come back as an error")
	}

	err = gen.oneYouthEvent("ages 4-8")
	if err == nil {
		t.Fatal("a refused path did not stop the event")
	}

	err = gen.youthEvents()
	if err == nil {
		t.Fatal("a refused path did not stop the step")
	}
}

// TestAPathRequirementNamesSomethingReal. Open() asks the caller for a
// characteristic by name, and a name nothing answers to scores zero rather
// than opening the path by accident.
func TestAPathRequirementNamesSomethingReal(t *testing.T) {
	t.Parallel()

	gen := engine(t, 18)

	for _, which := range CharacteristicOrder {
		gen.char.State.Characteristics.Set(which, 15)
	}

	nonsense := career.YouthPath{
		Name:     "Path 0",
		Requires: []career.Check{{Characteristic: "SOC", Number: 2}},
	}

	// SOC is Traveller's, not this ruleset's, so it scores zero.
	if nonsense.Open(gen.characteristicScore()) {
		t.Error("a path requiring a characteristic this ruleset does not have is open")
	}
}
