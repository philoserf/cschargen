package chargen

import (
	"slices"

	"github.com/philoserf/cschargen/career"
)

// skillListCite is the pages the skill list is printed on.
const skillListCite = "pp. 304-314"

// raiseHeldSkill is the thirty-one results that read "gain one level in any
// skill you already possess". The options are whatever the character has
// when the result fires, which is why the table cannot name them.
//
// A character who holds nothing gets nothing, and the record says so. That
// is reachable: the result appears on tables a character can meet before
// any skill has been granted.
func (g *Generator) raiseHeldSkill(effect career.Effect, cause int) error {
	held := g.char.State.Skills

	// "The Science skill you used on this task" is one of the character's
	// Sciences rather than any skill at all.
	if effect.OnSkill != "" {
		held = skillsNamed(held, effect.OnSkill)
	}

	if len(held) == 0 {
		g.consequence(ConsequenceSkill, cause, effect.Detail+": the character holds none", "")

		return nil
	}

	labels := make([]string, len(held))
	for i, skill := range held {
		labels[i] = skill.Full() + "-" + itoa(skill.Level)
	}

	chosen, err := g.choose(Choice{
		Point:   "raise_held_skill",
		Prompt:  asPrompt(effect.Detail),
		Options: labels,
		Cite:    skillListCite,
	})
	if err != nil {
		return err
	}

	// The choice is made against the list as it stands, so the specialty
	// is taken from the entry rather than asked for again: "a skill you
	// already possess" is one entry of the sheet, not a skill name that
	// might have several.
	picked := held[chosen]

	return g.applySkill(career.Effect{
		Kind:        career.EffectSkill,
		Skill:       picked.Name,
		Specialties: []string{picked.Specialty},
		Detail:      effect.Detail,
	}, cause)
}

// skillsNamed is the character's entries for one skill, whatever their
// specialties.
func skillsNamed(held []Skill, name string) []Skill {
	found := make([]Skill, 0, len(held))

	for _, skill := range held {
		if skill.Name == name {
			found = append(found, skill)
		}
	}

	return found
}

// anySkill is "gain a level in any skill of your choice", chosen from the
// list of pp. 304-314. NotHeld narrows it to what the character lacks.
//
// A skill with specialties is offered by name, and applySkill puts the
// specialty choice that follows to the decider: the two questions are the
// book's own, and folding them into one would make a list of hundreds.
func (g *Generator) anySkill(effect career.Effect, cause int) error {
	options := make([]string, 0, len(career.Skills()))
	definitions := make([]career.SkillDefinition, 0, len(career.Skills()))

	for _, definition := range career.Skills() {
		if effect.NotHeld && g.char.State.Has(definition.Name) {
			continue
		}

		options = append(options, definition.Name)
		definitions = append(definitions, definition)
	}

	if len(options) == 0 {
		g.consequence(ConsequenceSkill, cause,
			effect.Detail+": the character already holds every skill", "")

		return nil
	}

	chosen, err := g.choose(Choice{
		Point:   "any_skill",
		Prompt:  asPrompt(effect.Detail),
		Options: options,
		Cite:    skillListCite,
	})
	if err != nil {
		return err
	}

	picked := definitions[chosen]

	specialties := picked.Specialties
	if len(specialties) == 0 {
		specialties = nil
	}

	return g.applySkill(career.Effect{
		Kind:        career.EffectSkill,
		Skill:       picked.Name,
		Specialties: specialties,
		Level:       max(effect.Level, 1),
		Detail:      effect.Detail,
	}, cause)
}

// loseSkill removes every level of a named skill, which one mishap does:
// "You lose all levels of Suit (Vacc Suit)" (Orbital Construction).
func (g *Generator) loseSkill(effect career.Effect, cause int) {
	specialty := ""
	if len(effect.Specialties) > 0 {
		specialty = effect.Specialties[0]
	}

	kept := make([]Skill, 0, len(g.char.State.Skills))
	lost := false

	for _, held := range g.char.State.Skills {
		if held.Name == effect.Skill && held.Specialty == specialty {
			lost = true

			continue
		}

		kept = append(kept, held)
	}

	g.char.State.Skills = kept

	detail := effect.Detail
	if !lost {
		detail += ", which the character did not hold"
	}

	g.consequence(ConsequenceSkill, cause, detail, "")
}

// homeworldSpecialty is the five results that send the character back to
// where they grew up for a specialty: "gain a level in a Survival specialty
// which is used on your homeworld".
//
// The specialties are read from the world's own background skills, which is
// the only place the setting data says what a world is like in those terms.
// A world whose background skills name none -- and several do not -- leaves
// the result with nothing to grant, and the record says which world it was.
func (g *Generator) homeworldSpecialty(effect career.Effect, cause int) error {
	options := g.specialtiesOnTheHomeworld(effect.Skill)
	if len(options) == 0 {
		g.consequence(ConsequenceSkill, cause,
			effect.Detail+": the homeworld's background skills name none", "")

		return nil
	}

	return g.applySkill(career.Effect{
		Kind:        career.EffectSkill,
		Skill:       effect.Skill,
		Specialties: options,
		Detail:      effect.Detail,
	}, cause)
}

// specialtiesOnTheHomeworld collects the specialties one skill is offered
// with in the birth world's background skills, in the order the chart
// prints them and without repeats.
func (g *Generator) specialtiesOnTheHomeworld(name string) []string {
	if len(g.char.State.Homeworlds) == 0 {
		return nil
	}

	world, _, found := g.setting.World(g.char.State.Homeworlds[0].World)
	if !found {
		return nil
	}

	options := make([]string, 0, len(world.BackgroundSkills))

	for _, requirement := range world.BackgroundSkills {
		for _, alternative := range requirement.OneOf {
			if alternative.Skill != name {
				continue
			}

			for _, specialty := range alternative.Specialties {
				if !slices.Contains(options, specialty) {
					options = append(options, specialty)
				}
			}
		}
	}

	return options
}

// addCondition records something a character carries that is neither a
// skill, a characteristic nor a possession: an addiction, a religion.
//
// A condition is recorded once. A character who is already addicted to
// alcohol and takes the result again is no more addicted than before, and
// a second entry on the sheet would read as a second addiction.
func (g *Generator) addCondition(effect career.Effect, cause int) {
	if slices.Contains(g.char.State.Conditions, effect.Condition) {
		g.consequence(ConsequenceCharacteristic, cause,
			effect.Detail+", which the character already carries", "")

		return
	}

	g.char.State.Conditions = append(g.char.State.Conditions, effect.Condition)
	g.consequence(ConsequenceCharacteristic, cause, effect.Detail, "")
}

// ageBy moves the character's age outside the four years a term takes,
// which one result does: "you have lost a year of your life".
func (g *Generator) ageBy(effect career.Effect, cause int) {
	g.char.State.Age += effect.Years
	g.consequence(ConsequenceAge, cause,
		effect.Detail+": age "+itoa(g.char.State.Age), "")
}

// adjustSentence lengthens or shortens a forced spell in a career. A
// sentence cut to one term or less is served: "if you have one term or
// less remaining, you are released".
func (g *Generator) adjustSentence(effect career.Effect, cause int) {
	if g.forcedTerms == 0 {
		g.unimplemented(cause, effect.Detail+" -- the character is serving no sentence")

		return
	}

	g.forcedTerms = max(g.forcedTerms+effect.Terms, g.termsInCareer)

	if g.forcedTerms <= g.termsInCareer {
		g.consequence(ConsequenceCareer, cause, effect.Detail+": released", g.career.Name)

		return
	}

	g.consequence(ConsequenceCareer, cause,
		effect.Detail+": "+itoa(g.forcedTerms-g.termsInCareer)+" terms left", g.career.Name)
}
