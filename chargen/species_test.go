package chargen_test

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
	"github.com/philoserf/cschargen/setting"
)

// TestEverySpeciesInTheSampleGenerates is the reach check: four species are
// declared and, before this, no character had ever been generated as one.
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

// TestAWorldThatWillNotHaveThem is p. 42: a world whose entry bars this kind
// of character sends the player to a different homeworld.
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
// p. 42: engineered and uplift characters age differently, and the
// restrictions do not apply to them.
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
// the book gives no other entrance: an enslaved character must take the
// slave career as their first career term (p. 42).
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

			if wasFreed(character) {
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

			if wasFreed(character) {
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

// wasFreed reports whether an enslavement ended during early life. Three
// results in the book end one -- the character continues free -- and one
// they reach is owned at birth and free by the time the first career is
// chosen, so p. 42 no longer binds them.
//
// The tests below are about what being owned does, so a character it stopped
// doing anything to is not one of their cases.
func wasFreed(character *chargen.Character) bool {
	for _, event := range character.Events {
		if event.Kind != chargen.EventConsequence || event.Consequence == nil {
			continue
		}

		if strings.HasPrefix(event.Consequence.Detail, "freed") {
			return true
		}
	}

	return false
}

// TestAnEngineeredCharacterHasGenetics is pp. 62-65: they are a purebred of
// some generation, a hybrid with a baseline human, or a compound of two
// engineered species.
func TestAnEngineeredCharacterHasGenetics(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)
	kinds := map[string]int{}

	for _, species := range data.Species {
		for seed := range uint64(60) {
			opts := options(t, seed)

			opts.Inputs.Species = species.Name
			opts.Inputs.TermLimit = -1

			genetics := generate(t, opts).State.Genetics

			if species.Kind != "engineered" {
				if genetics != nil {
					t.Errorf("a %s has genetics; pp. 62-65 give them to engineered "+
						"species only", species.Name)
				}

				continue
			}

			if genetics == nil {
				t.Fatalf("seed %d: a %s has no genetics", seed, species.Name)
			}

			kinds[genetics.Kind]++

			checkGenetics(t, seed, species.Name, genetics)
		}
	}

	for _, kind := range []string{"purebred", "hybrid", "compound"} {
		if kinds[kind] == 0 {
			t.Errorf("no character in the sample was ever a %s", kind)
		}
	}
}

// TestAHybridMayBeRolledAsAHuman is the one mechanical thing the genetics
// tables do: "Roll characteristics as a human" (p. 63). A character whose
// hybrid outcome says so does not use their species' method at all.
func TestAHybridMayBeRolledAsAHuman(t *testing.T) {
	t.Parallel()

	data := sampleSetting(t)
	asHuman := 0

	for _, species := range data.Species {
		if species.Genetics == nil || len(species.Genetics.Hybrid) == 0 {
			continue
		}

		for seed := range uint64(80) {
			opts := options(t, seed)

			opts.Inputs.Species = species.Name
			opts.Inputs.TermLimit = -1

			character := generate(t, opts)
			if character.State.Genetics == nil || !character.State.Genetics.RolledAsHuman {
				continue
			}

			asHuman++

			// A human's six scores come from 3d6-drop-lowest, which runs 2
			// to 12 -- so none of them can exceed what that can produce,
			// whatever the species' own method would have given.
			//
			// Step 2's own record, not the finished sheet: Steps 6 to 8
			// raise characteristics legitimately, and a degree can put EDU
			// above 12 (pp. 87, 92) on a character whose EDU was rolled
			// the human way and rolled low.
			for _, rolled := range rolledCharacteristics(character) {
				if rolled.Delta > maxHumanScore {
					t.Errorf("seed %d: a hybrid rolled as a human rolled %s %d",
						seed, rolled.Characteristic, rolled.Delta)
				}
			}
		}
	}

	if asHuman == 0 {
		t.Error("no hybrid in the sample was ever rolled as a human")
	}
}

// checkGenetics: every outcome says something, and a purebred is a
// purebred of some generation.
func checkGenetics(t *testing.T, seed uint64, name string, genetics *chargen.Genetics) {
	t.Helper()

	if genetics.Detail == "" {
		t.Errorf("seed %d: a %s's genetics say nothing", seed, name)
	}

	if genetics.Kind == "purebred" && genetics.Generation == 0 {
		t.Errorf("seed %d: a purebred of no generation", seed)
	}
}

// maxHumanScore is what 3d6-drop-lowest can produce (p. 13).
const maxHumanScore = 12

// rolledCharacteristics is what Step 2 put on the sheet, before anything
// later raised it. Each assignment is recorded as it is made, cited to the
// page the method is on.
func rolledCharacteristics(character *chargen.Character) []chargen.ConsequenceEvent {
	var rolled []chargen.ConsequenceEvent

	for _, event := range character.Events {
		if event.Kind != chargen.EventConsequence || event.Consequence == nil {
			continue
		}

		if event.Consequence.Kind == chargen.ConsequenceCharacteristic &&
			event.Consequence.Cite == "p. 13" {
			rolled = append(rolled, *event.Consequence)
		}
	}

	return rolled
}
