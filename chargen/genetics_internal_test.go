package chargen

import (
	"errors"
	"testing"

	"github.com/philoserf/cschargen/setting"
)

// engineered is a species with the three genetics tables of pp. 62-65.
func engineered(name string) setting.Species {
	six := func(detail string) []setting.GeneticOutcome {
		table := make([]setting.GeneticOutcome, 6)
		for i := range table {
			table[i] = setting.GeneticOutcome{Detail: detail}
		}

		return table
	}

	compound := six("a compound")

	compound[len(compound)-1] = setting.GeneticOutcome{
		Detail: "the player chooses", Choose: true,
	}

	return setting.Species{
		Name: name, Kind: setting.KindEngineered,
		Genetics: &setting.Genetics{Hybrid: six("a hybrid"), Compound: compound},
	}
}

// TestARefusedGenerationEndsGeneration. "The character is a third, fourth
// or later generation purebred altrant. The player chooses which
// generation" (p. 63) is a choice like any other.
func TestARefusedGenerationEndsGeneration(t *testing.T) {
	t.Parallel()

	species := engineered("Refused")

	species.Genetics.Status = []setting.GeneticStatus{
		{Kind: setting.Purebred},
		{Kind: setting.Purebred},
		{Kind: setting.Purebred},
		{Kind: setting.Purebred},
		{Kind: setting.Purebred},
		{Kind: setting.Purebred},
	}

	gen := speciesEngine(t, 45, species)

	gen.decider = refusingDecider{}

	_, err := gen.Run()
	if err == nil {
		t.Fatal("a refused generation did not end generation")
	}
}

// TestARefusedCompoundOutcomeEndsGeneration is the last row of every
// compound table: "The player chooses from the above results" (pp. 64-65).
func TestARefusedCompoundOutcomeEndsGeneration(t *testing.T) {
	t.Parallel()

	species := engineered("Compounded")

	species.Genetics.Status = []setting.GeneticStatus{
		{Kind: setting.Compound},
		{Kind: setting.Compound},
		{Kind: setting.Compound},
		{Kind: setting.Compound},
		{Kind: setting.Compound},
		{Kind: setting.Compound},
	}

	gen := speciesEngine(t, 46, species)

	gen.decider = refusingDecider{}

	_, err := gen.Run()
	if err == nil {
		t.Fatal("a refused compound outcome did not end generation")
	}
}

// TestACompoundTableThatOnlyDefers. "The player chooses from the above
// results" with nothing above it is a table with no answer in it.
func TestACompoundTableThatOnlyDefers(t *testing.T) {
	t.Parallel()

	species := engineered("Empty")

	for i := range species.Genetics.Compound {
		species.Genetics.Compound[i] = setting.GeneticOutcome{
			Detail: "the player chooses", Choose: true,
		}
	}

	gen := speciesEngine(t, 47, species)

	gen.species = &species

	_, err := gen.chooseOutcome(species.Genetics.Compound, 0)
	if !errors.Is(err, ErrNoOptions) {
		t.Errorf("err = %v, want ErrNoOptions", err)
	}
}

// TestACompoundOutcomeIsChosen. The last row hands the decision back, and
// the policy takes the first of the rows above it.
func TestACompoundOutcomeIsChosen(t *testing.T) {
	t.Parallel()

	species := engineered("Chosen")

	species.Genetics.Status = []setting.GeneticStatus{
		{Kind: setting.Compound},
		{Kind: setting.Compound},
		{Kind: setting.Compound},
		{Kind: setting.Compound},
		{Kind: setting.Compound},
		{Kind: setting.Compound},
	}
	species.Genetics.Compound[0] = setting.GeneticOutcome{
		Detail: "the first result", Adjust: map[string]int{"DEX": 1},
	}

	// A decider that takes the last option reaches the row that defers.
	gen := speciesEngine(t, 48, species)

	gen.decider = lastOptionDecider{}

	character, err := gen.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if character.State.Genetics == nil {
		t.Fatal("no genetics")
	}

	if character.State.Genetics.Kind != setting.Compound {
		t.Errorf("kind = %q, want a compound", character.State.Genetics.Kind)
	}
}
