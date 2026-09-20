package chargen_test

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// TestEverySpeciesInTheSampleGenerates is the reach check the milestone 6
// plan makes: four species are declared and, before this, no character had
// ever been generated as one.
func TestEverySpeciesInTheSampleGenerates(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	for _, species := range data.Species {
		t.Run(species.Name, func(t *testing.T) {
			t.Parallel()

			opts := options(t, 5)

			opts.Inputs.Species = species.Name
			opts.Inputs.TermLimit = 4

			character := generate(t, opts)

			if character.Provenance.Inputs.Species != species.Name {
				t.Errorf("the record says %q", character.Provenance.Inputs.Species)
			}

			// Each characteristic the species names a method for is rolled
			// by that method and never assigned freely, so the score is
			// inside the range the expression can produce.
			for name, expr := range species.Characteristics {
				which, ok := chargen.CharacteristicByName(name)
				if !ok {
					t.Fatalf("%s names %q, which is not a characteristic", species.Name, name)
				}

				got := character.State.Characteristics.Get(which)
				if got < 0 || got > species.Maximum && species.Maximum > 0 {
					t.Errorf("%s rolled as %q and came out at %d", name, expr, got)
				}
			}

			// And the species' starting skills are held at level 1 or
			// better -- a career may have raised them.
			for _, start := range species.Skills {
				if !character.State.Has(start.Skill) {
					t.Errorf("a %s does not hold %s", species.Name, start.Skill)
				}
			}
		})
	}
}

// TestASpeciesTheDataDoesNotDeclare. The engine holds no species of its
// own, so there is nothing to fall back to.
func TestASpeciesTheDataDoesNotDeclare(t *testing.T) {
	t.Parallel()

	opts := options(t, 5)

	opts.Inputs.Species = "Basilisk"

	_, err := chargen.New(opts).Run()
	if err == nil {
		t.Fatal("a species the data file does not declare was generated anyway")
	}
}

// TestAHumanNeedsNoSpeciesEntry. "If the character is a human, simply move
// on to Step 2" (p. 21) -- and a human is not one of the file's species.
func TestAHumanNeedsNoSpeciesEntry(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"human", ""} {
		opts := options(t, 5)

		opts.Inputs.Species = name
		opts.Inputs.TermLimit = -1

		character := generate(t, opts)
		if character.Provenance.Inputs.Species != chargen.Human {
			t.Errorf("%q generated as %q", name, character.Provenance.Inputs.Species)
		}
	}
}

// TestASpeciesCeilingHolds. p. 14 caps an unaltered human at 15 and says
// "Some uplifts and engineered humans will have higher maximums".
func TestASpeciesCeilingHolds(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	for _, species := range data.Species {
		if species.Maximum == 0 {
			continue
		}

		for seed := range uint64(30) {
			opts := options(t, seed)

			opts.Inputs.Species = species.Name
			opts.Inputs.TermLimit = 12

			character := generate(t, opts)

			for _, which := range chargen.CharacteristicOrder {
				if got := character.State.Characteristics.Get(which); got > species.Maximum {
					t.Fatalf("a %s reached %s %d, above their ceiling of %d",
						species.Name, which, got, species.Maximum)
				}
			}
		}
	}
}

// TestASpeciesAgesOnItsOwnProfile. Three of the sample's four species name
// a profile that is not the tech-level one, and the record says which.
func TestASpeciesAgesOnItsOwnProfile(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	for _, species := range data.Species {
		if species.Aging == "" {
			continue
		}

		t.Run(species.Name, func(t *testing.T) {
			t.Parallel()

			aged := 0

			for seed := range uint64(40) {
				opts := options(t, seed)

				opts.Inputs.Species = species.Name
				opts.Inputs.TermLimit = 12

				for _, event := range generate(t, opts).Events {
					if event.Consequence == nil {
						continue
					}

					if strings.Contains(event.Consequence.Detail,
						"on the "+species.Aging+" aging profile") {
						aged++
					}
				}
			}

			if aged == 0 {
				t.Errorf("no %s in 40 seeds ever reached the %s profile",
					species.Name, species.Aging)
			}
		})
	}
}
