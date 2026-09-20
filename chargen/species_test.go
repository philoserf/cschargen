package chargen_test

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/setting"
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

// TestAWorldThatWillNotHaveThem is p. 42: "If the entry indicates that
// altrants or uplifts are not allowed, then the player should select a
// different homeworld."
func TestAWorldThatWillNotHaveThem(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	for _, species := range data.Species {
		t.Run(species.Name, func(t *testing.T) {
			t.Parallel()

			for seed := range uint64(60) {
				opts := options(t, seed)

				opts.Inputs.Species = species.Name
				opts.Inputs.TermLimit = -1

				character := generate(t, opts)

				born, _, found := data.World(character.State.Homeworlds[0].World)
				if !found {
					t.Fatalf("seed %d: unknown birth world", seed)
				}

				permission := born.Engineered
				if species.Kind == "uplift" {
					permission = born.Uplifts
				}

				if !permission.Admits(species.Name) {
					t.Fatalf("seed %d: a %s was born on %s, which does not admit them",
						seed, species.Name, born.Name)
				}
			}
		})
	}
}

// TestTheHomeworldCapsDoNotBindAnAlteredCharacter is the other half of
// p. 42: "Altrant and Uplift characters age differently and these
// restrictions will not apply to them."
func TestTheHomeworldCapsDoNotBindAnAlteredCharacter(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	// Every world in the sample caps terms well below forty.
	const asked = 40

	longer := 0

	for _, species := range data.Species {
		for seed := range uint64(30) {
			opts := options(t, seed)

			opts.Inputs.Species = species.Name
			opts.Inputs.TermLimit = asked

			character := generate(t, opts)

			born, _, found := data.World(character.State.Homeworlds[0].World)
			if !found {
				t.Fatalf("seed %d: unknown birth world", seed)
			}

			if len(character.State.Terms) > born.MaximumTerms {
				longer++
			}
		}
	}

	if longer == 0 {
		t.Error("no altered character ever served past their homeworld's cap")
	}
}

// TestAnEnslavedCharacterTakesTheSlaveCareerFirst is the way into a career
// the book gives no other entrance: "the character must take the
// Altrant/Uplift Slave career as their first career term" (p. 42).
func TestAnEnslavedCharacterTakesTheSlaveCareerFirst(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)
	enslaved := 0

	for _, species := range data.Species {
		for seed := range uint64(60) {
			opts := options(t, seed)

			opts.Inputs.Species = species.Name
			opts.Inputs.TermLimit = 3

			character := generate(t, opts)

			born, _, found := data.World(character.State.Homeworlds[0].World)
			if !found {
				continue
			}

			permission := born.Engineered
			if species.Kind == "uplift" {
				permission = born.Uplifts
			}

			if permission.Status != "enslaved" {
				continue
			}

			enslaved++

			if len(character.State.Services) == 0 {
				t.Errorf("seed %d: a %s born owned on %s served no career",
					seed, species.Name, born.Name)

				continue
			}

			if first := character.State.Services[0].Career; first != slaveCareerName {
				t.Errorf("seed %d: a %s born owned on %s began in %s",
					seed, species.Name, born.Name, first)
			}
		}
	}

	if enslaved == 0 {
		t.Fatal("no seed produced a character born on a world that owns their people")
	}
}

// slaveCareerName is the career of pp. 150-154, under the name this
// repository gives it.
const slaveCareerName = "Engineered/Uplift Slave"

// TestAnUpliftsClassIsWhatTheirWorldCanMake is p. 66's chart: Class 1 needs
// tech level 10, Class 2 needs 11, Class 3 needs 12. The policy takes the
// highest available, which is what the page recommends.
func TestAnUpliftsClassIsWhatTheirWorldCanMake(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)
	seen := map[int]int{}

	for _, species := range data.Species {
		for seed := range uint64(40) {
			opts := options(t, seed)

			opts.Inputs.Species = species.Name
			opts.Inputs.TermLimit = -1

			character := generate(t, opts)

			if species.Kind != "uplift" {
				if character.State.UpliftClass != 0 {
					t.Errorf("a %s has an uplift class", species.Name)
				}

				continue
			}

			seen[checkUpliftClass(t, data, species.Name, character)]++
		}
	}

	for class := 1; class <= 3; class++ {
		if seen[class] == 0 {
			t.Errorf("no uplift in the sample was ever Class %d", class)
		}
	}
}

// checkUpliftClass holds one uplift's class against the chart on p. 66 --
// Class 1 at tech level 10, Class 2 at 11, Class 3 at 12 -- and returns it.
func checkUpliftClass(
	t *testing.T, data *setting.Data, name string, character *chargen.Character,
) int {
	t.Helper()

	born, _, found := data.World(character.State.Homeworlds[0].World)
	if !found {
		t.Fatal("unknown birth world")
	}

	want := 1

	switch {
	case born.TechLevel >= 12:
		want = 3
	case born.TechLevel >= 11:
		want = 2
	}

	if character.State.UpliftClass != want {
		t.Errorf("a %s from a tech level %d world is Class %d, want Class %d",
			name, born.TechLevel, character.State.UpliftClass, want)
	}

	return character.State.UpliftClass
}

// TestAnUpliftsLifePeriodsAreTheirOwn is p. 68: "Uplifts will often have
// shorter youths than humans. Whereas humans will roll on the following
// charts twice to represent their lives from 4-8 and again from 9-12,
// uplifts will have a different age range and may roll fewer times."
func TestAnUpliftsLifePeriodsAreTheirOwn(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)

	for _, species := range data.Species {
		if species.YouthRolls == 0 {
			continue
		}

		t.Run(species.Name, func(t *testing.T) {
			t.Parallel()

			opts := options(t, 5)

			opts.Inputs.Species = species.Name
			opts.Inputs.TermLimit = -1

			character := generate(t, opts)

			youth := countPeriods(character, species.YouthAges)
			teenage := countPeriods(character, species.TeenAges)

			if youth != species.YouthRolls {
				t.Errorf("%d youth events, want %d", youth, species.YouthRolls)
			}

			if teenage != species.TeenRolls {
				t.Errorf("%d teenage events, want %d", teenage, species.TeenRolls)
			}
		})
	}
}

// countPeriods is how many life-period events a record holds for a set of
// age ranges.
func countPeriods(character *chargen.Character, ages []string) int {
	found := 0

	for _, event := range character.Events {
		if event.Consequence == nil {
			continue
		}

		for _, age := range ages {
			if strings.HasPrefix(event.Consequence.Detail, age+", on ") {
				found++
			}
		}
	}

	return found
}

// TestAnEnslavedCharactersEarlyLife is pp. 74 and 84: a character their
// homeworld owns lives a different childhood, on an eleven-row 2d6 table
// rather than a nineteen-row 2d10 one, and takes no path because there is
// only one.
func TestAnEnslavedCharactersEarlyLife(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)
	owned := 0

	for _, species := range data.Species {
		for seed := range uint64(60) {
			opts := options(t, seed)

			opts.Inputs.Species = species.Name
			opts.Inputs.TermLimit = 1

			character := generate(t, opts)

			born, _, found := data.World(character.State.Homeworlds[0].World)
			if !found {
				continue
			}

			permission := born.Engineered
			if species.Kind == "uplift" {
				permission = born.Uplifts
			}

			if permission.Status != "enslaved" {
				continue
			}

			owned++

			checkNoPathOffered(t, seed, character)
		}
	}

	if owned == 0 {
		t.Fatal("no seed produced a character born owned")
	}
}

// checkNoPathOffered: the enslaved tables are the only life-period tables
// that are not a choice, because there is only one of each.
func checkNoPathOffered(t *testing.T, seed uint64, character *chargen.Character) {
	t.Helper()

	for _, event := range character.Events {
		if event.Kind != chargen.EventChoice {
			continue
		}

		if event.Choice.Point == "youth_path" || event.Choice.Point == "teenage_path" {
			t.Errorf("seed %d: an owned character was offered a path to choose", seed)
		}
	}
}
