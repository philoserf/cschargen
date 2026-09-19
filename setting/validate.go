package setting

import "fmt"

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

		if kind.Kind != "engineered" && kind.Kind != "uplift" {
			problems = append(problems, fmt.Sprintf(
				"species %s has kind %q; it must be engineered or uplift", kind.Name, kind.Kind))
		}
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
