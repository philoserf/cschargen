package setting

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// The bounds a world's numbers must fall inside. They are wide on purpose:
// the file is a transcription, and the validator's job is to catch a typo
// or an omission rather than to second-guess the book.
const (
	maxTechLevel = 20
	maxAge       = 400
	maxTerms     = 100
	d100Low      = 1
	d100High     = 100
	d6Low        = 1
	d6High       = 6
)

// Validate reports every problem in a setting file. Every problem, not the
// first: someone transcribing sixty worlds should be told about all of
// them in one pass rather than made to run the validator sixty times.
func Validate(data *Data) []string {
	var problems []string

	if data.SchemaVersion != SchemaVersion {
		problems = append(problems, fmt.Sprintf(
			"schemaVersion is %d; this build reads %d", data.SchemaVersion, SchemaVersion))
	}

	species := map[string]bool{}

	problems = append(problems, validateSpecies(data.Species, species)...)

	if len(data.Subsectors) == 0 {
		problems = append(problems, "no subsectors")
	}

	seenSubsector := map[string]bool{}
	seenOriginRoll := map[int]string{}
	seenWorld := map[string]string{}

	for i, sub := range data.Subsectors {
		where := fmt.Sprintf("subsector %d", i)
		if sub.Name != "" {
			where = "subsector " + sub.Name
		}

		problems = append(problems, validateSubsector(sub, where, seenSubsector, seenOriginRoll)...)
		problems = append(problems, validateWorlds(sub, where, species, seenWorld)...)
		problems = append(problems, validateCoverage(sub, where)...)
	}

	return problems
}

func validateSpecies(all []Species, names map[string]bool) []string {
	var problems []string

	for i, kind := range all {
		if kind.Name == "" {
			problems = append(problems, fmt.Sprintf("species %d has no name", i))

			continue
		}

		if names[kind.Name] {
			problems = append(problems, "species "+kind.Name+" is listed twice")
		}

		names[kind.Name] = true

		problems = append(problems, validateOneSpecies(kind, "species "+kind.Name)...)
	}

	return problems
}

func validateSubsector(sub Subsector, where string, seen map[string]bool, rolls map[int]string) []string {
	var problems []string

	if sub.Name == "" {
		problems = append(problems, where+" has no name")
	} else if seen[sub.Name] {
		problems = append(problems, where+" is listed twice")
	}

	seen[sub.Name] = true

	if sub.OriginRoll != 0 {
		if sub.OriginRoll < d6Low || sub.OriginRoll > d6High {
			problems = append(problems, fmt.Sprintf(
				"%s has originRoll %d; the chart on p. 39 is 1d6", where, sub.OriginRoll))
		}

		if other, taken := rolls[sub.OriginRoll]; taken {
			problems = append(problems, fmt.Sprintf(
				"%s and %s both claim originRoll %d", where, other, sub.OriginRoll))
		}

		rolls[sub.OriginRoll] = where
	}

	if len(sub.Worlds) == 0 {
		problems = append(problems, where+" has no worlds")
	}

	return problems
}

func validateWorlds(sub Subsector, where string, species map[string]bool, worlds map[string]string) []string {
	var problems []string

	covered := map[int]string{}

	for i, world := range sub.Worlds {
		name := world.Name
		if name == "" {
			problems = append(problems, fmt.Sprintf("%s: world %d has no name", where, i))

			continue
		}

		if other, taken := worlds[name]; taken {
			problems = append(problems, fmt.Sprintf(
				"world %s appears in both %s and %s; names identify a homeworld", name, other, where))
		}

		worlds[name] = where

		problems = append(problems, validateWorld(world, where+", "+name, species)...)
		problems = append(problems, validateRange(world, where+", "+name, covered)...)
	}

	return problems
}

func validateRange(world World, where string, covered map[int]string) []string {
	if world.Roll == nil {
		return nil
	}

	low, high := world.Roll[0], world.Roll[1]
	if low < d100Low || high > d100High || low > high {
		return []string{fmt.Sprintf("%s: roll %d-%d is not a d100 range", where, low, high)}
	}

	var problems []string

	for result := low; result <= high; result++ {
		if other, taken := covered[result]; taken {
			problems = append(problems, fmt.Sprintf(
				"%s: a d100 of %d also lands on %s", where, result, other))

			break
		}

		covered[result] = world.Name
	}

	return problems
}

// AgingProfiles is what a species' aging field may name.
//
// The names live here rather than in the engine because they are a fact
// about the file format -- a validator that could not check them would pass
// a file the engine then had to fall back on -- and the engine holds them
// against its own tables in a test. The five are named for what they do,
// because what the book names them by is a species (p. 335).
var AgingProfiles = map[string]bool{
	"techLevel": true,
	"rapid":     true,
	"moderate":  true,
	"sturdy":    true,
	"sudden":    true,
}

// validateOneSpecies holds a species' declarations against what the engine
// can carry out: an aging profile it knows, a characteristic the book has,
// and a ceiling inside the range a score can reach.
func validateOneSpecies(species Species, where string) []string {
	var problems []string

	if species.Kind != KindEngineered && species.Kind != KindUplift {
		problems = append(problems, fmt.Sprintf(
			"%s: kind %q; it must be engineered or uplift (p. 21)", where, species.Kind))
	}

	if species.Aging != "" && !AgingProfiles[species.Aging] {
		problems = append(problems, fmt.Sprintf(
			"%s: aging %q is not a profile the engine holds", where, species.Aging))
	}

	for name := range species.Characteristics {
		if !characteristics[name] {
			problems = append(problems, fmt.Sprintf(
				"%s: %q is not one of the six characteristics (p. 13)", where, name))
		}
	}

	if species.Maximum < 0 || species.Maximum > maxCharacteristic {
		problems = append(problems, fmt.Sprintf(
			"%s: maximum %d", where, species.Maximum))
	}

	problems = append(problems, validateGenetics(species, where)...)

	problems = append(problems, validateRollPattern(
		species.YouthRolls, species.YouthAges, where+", youth")...)
	problems = append(problems, validateRollPattern(
		species.TeenRolls, species.TeenAges, where+", teenage")...)

	return problems
}

// validateGenetics holds a species' three tables of pp. 62-65 against their
// printed shape: each is a 1d6, so each has six rows.
func validateGenetics(species Species, where string) []string {
	if species.Genetics == nil {
		return nil
	}

	var problems []string

	if species.Kind != KindEngineered {
		problems = append(problems, where+
			": genetics, which pp. 62-65 give only to engineered species")
	}

	tables := map[string][]GeneticOutcome{
		"hybrid":   species.Genetics.Hybrid,
		"compound": species.Genetics.Compound,
	}

	for name, table := range tables {
		problems = append(problems, validateOutcomes(table, name, where)...)
	}

	return append(problems, validateGeneticStatus(species, where)...)
}

// validateOutcomes holds one hybrid or compound table against its shape.
func validateOutcomes(table []GeneticOutcome, name, where string) []string {
	var problems []string

	if len(table) != 0 && len(table) != d6Rows {
		problems = append(problems, fmt.Sprintf(
			"%s: the %s table has %d rows; a 1d6 table has %d",
			where, name, len(table), d6Rows))
	}

	for i, outcome := range table {
		if outcome.Detail == "" && !outcome.Choose {
			problems = append(problems, fmt.Sprintf(
				"%s: %s result %d says nothing", where, name, i+1))
		}

		for characteristic := range outcome.Adjust {
			if !characteristics[characteristic] {
				problems = append(problems, fmt.Sprintf(
					"%s: %s result %d adjusts %q, which is not a characteristic",
					where, name, i+1, characteristic))
			}
		}
	}

	return problems
}

// validateGeneticStatus holds the status table, whose rows name which of
// the other two a character is sent to.
func validateGeneticStatus(species Species, where string) []string {
	var problems []string

	status := species.Genetics.Status
	if len(status) != 0 && len(status) != d6Rows {
		problems = append(problems, fmt.Sprintf(
			"%s: the genetic status table has %d rows; a 1d6 table has %d",
			where, len(status), d6Rows))
	}

	for i, row := range status {
		switch row.Kind {
		case Purebred:
		case Hybrid:
			if len(species.Genetics.Hybrid) == 0 {
				problems = append(problems, fmt.Sprintf(
					"%s: status result %d sends the character to a hybrid table it has none of",
					where, i+1))
			}
		case Compound:
			if len(species.Genetics.Compound) == 0 {
				problems = append(problems, fmt.Sprintf(
					"%s: status result %d sends the character to a compound table it has none of",
					where, i+1))
			}
		default:
			problems = append(problems, fmt.Sprintf(
				"%s: status result %d is %q; it must be purebred, hybrid or compound",
				where, i+1, row.Kind))
		}
	}

	return problems
}

// d6Rows is how many rows a 1d6 table has.
const d6Rows = 6

// validateRollPattern: a species that rolls n times on a life-period table
// has to say what each of those n rolls represents.
func validateRollPattern(rolls int, ages []string, where string) []string {
	if rolls == 0 && len(ages) == 0 {
		return nil
	}

	if rolls != len(ages) {
		return []string{fmt.Sprintf(
			"%s: %d rolls and %d age ranges; a roll represents a range", where, rolls, len(ages))}
	}

	return nil
}

// characteristics is the six, as a species' method names them.
var characteristics = map[string]bool{
	"STR": true, "DEX": true, "END": true, "INT": true, "EDU": true, "CHA": true,
}

// maxCharacteristic is a generous ceiling on a species' own ceiling. The
// book's highest is well under it; the validator's job is to catch a typo
// rather than to second-guess a setting.
const maxCharacteristic = 30

func validateWorld(world World, where string, species map[string]bool) []string {
	var problems []string

	if world.TechLevel < 0 || world.TechLevel > maxTechLevel {
		problems = append(problems, fmt.Sprintf("%s: techLevel %d", where, world.TechLevel))
	}

	if world.MaximumAge <= 0 || world.MaximumAge > maxAge {
		problems = append(problems, fmt.Sprintf("%s: maximumAge %d", where, world.MaximumAge))
	}

	if world.MaximumTerms <= 0 || world.MaximumTerms > maxTerms {
		problems = append(problems, fmt.Sprintf("%s: maximumTerms %d", where, world.MaximumTerms))
	}

	if len(world.PrimaryLanguages) == 0 {
		problems = append(problems, where+": no primaryLanguages; every homeworld lists one (p. 41)")
	}

	// A birth situation that is forced has to be one of the three the chart
	// prints (p. 58), because nothing downstream knows what else to do with
	// it: a world that named a fourth would silently generate a household
	// with no parents in it.
	if forced := world.BirthSituationOnly; forced != "" && !BirthSituations[forced] {
		problems = append(problems, fmt.Sprintf(
			"%s: birthSituationOnly %q; the chart of p. 58 prints only %s",
			where, forced, joinSituations()))
	}

	problems = append(problems, validateRequirements(world.BackgroundSkills, where)...)
	problems = append(problems, validatePermission(world.Engineered, where+", engineered", species)...)
	problems = append(problems, validatePermission(world.Uplifts, where+", uplifts", species)...)

	return problems
}

func validateRequirements(requirements []Requirement, where string) []string {
	var problems []string

	for i, requirement := range requirements {
		if len(requirement.OneOf) == 0 {
			problems = append(problems, fmt.Sprintf(
				"%s: background skill %d offers no alternatives", where, i))

			continue
		}

		for j, alternative := range requirement.OneOf {
			at := fmt.Sprintf("%s: background skill %d, alternative %d", where, i, j)

			switch {
			case alternative.Skill == "" && alternative.Item == "":
				problems = append(problems, at+" names neither a skill nor an item")
			case alternative.Skill != "" && alternative.Item != "":
				problems = append(problems, at+" names both a skill and an item")
			case alternative.Item != "" && len(alternative.Specialties) > 0:
				problems = append(problems, at+" is an item with specialties")
			}
		}
	}

	return problems
}

func validatePermission(permission Permission, where string, species map[string]bool) []string {
	var problems []string

	if permission.Allowed && permission.Status != Free && permission.Status != Enslaved {
		problems = append(problems, fmt.Sprintf(
			"%s: status %q; it must be free or enslaved (p. 42)", where, permission.Status))
	}

	if !permission.Allowed && len(permission.Banned) > 0 {
		problems = append(problems, where+": names banned species although none are allowed")
	}

	for _, banned := range permission.Banned {
		if !species[banned] {
			problems = append(problems, fmt.Sprintf(
				"%s: bans %q, which is not in the species list", where, banned))
		}
	}

	return problems
}

// validateCoverage holds the thing a chart is: every d100 result lands
// somewhere. A subsector with any world a throw can reach must have a world
// for all hundred results, or a character rolled into the gap has no
// homeworld and no rule says what to do about it.
//
// A subsector with no selectable worlds at all is choose-only, which is
// what Recently Colonized Worlds is (p. 40), and is not checked.
func validateCoverage(sub Subsector, where string) []string {
	covered := map[int]bool{}
	selectable := false

	for _, world := range sub.Worlds {
		if world.Roll == nil {
			continue
		}

		selectable = true

		for result := world.Roll[0]; result <= world.Roll[1]; result++ {
			covered[result] = true
		}
	}

	if !selectable {
		return nil
	}

	var gaps []int

	for result := d100Low; result <= d100High; result++ {
		if !covered[result] {
			gaps = append(gaps, result)
		}
	}

	if len(gaps) == 0 {
		return nil
	}

	return []string{fmt.Sprintf(
		"%s: %d d100 results land on no world, the first being %d", where, len(gaps), gaps[0])}
}

// BirthSituations is the three results the Human Birth Situation chart
// prints (p. 58). A world may force one of them and no others.
//
// The strings are the engine's, not the book's headings: they are what a
// record says a character was born to, so the data file and the engine have
// to agree on them.
var BirthSituations = map[string]bool{
	"a communal group":      true,
	"a same sex couple":     true,
	"a heterosexual couple": true,
}

// joinSituations names the three in a stable order, for an error message.
func joinSituations() string {
	names := slices.Sorted(maps.Keys(BirthSituations))

	return strings.Join(names, ", ")
}
