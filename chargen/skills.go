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

	g.forcedTerms = max(g.forcedTerms+effect.Terms, g.service.termsInCareer)

	if g.forcedTerms <= g.service.termsInCareer {
		g.consequence(ConsequenceCareer, cause, effect.Detail+": released", g.service.career.Name)

		return
	}

	g.consequence(ConsequenceCareer, cause,
		effect.Detail+": "+itoa(g.forcedTerms-g.service.termsInCareer)+" terms left", g.service.career.Name)
}

// pointSkillSpecialty is the choice of a specialty, from a printed list or
// from the skill's own where a result writes "(Any)". Three places ask it.
const pointSkillSpecialty = "skill_specialty"

// anySpecialty is what a table writes where it names a skill and leaves the
// specialty to the character: "Electronics (Any)".
const anySpecialty = "Any"

// isAny reports whether a skill result leaves its specialty to the
// character.
func isAny(effect career.Effect) bool {
	return len(effect.Specialties) == 1 && effect.Specialties[0] == anySpecialty
}

// grantAny resolves "(Any)" outside a rank row, which p. 118 says how to do
// for a skill table's result:
//
//	"the character may choose to add a level to any specialty under that
//	skill. They may choose to add a level to a skill specialty that they
//	already have or add a new specialty ... If the character does not
//	already have the skill listed, it is added to the character sheet as a
//	zero level skill without a specialty."
//
// So the choice is every specialty the skill has, held or not, and a
// character without the skill gets it bare at level 0 instead. A grant at
// level 0 -- "Gain Drive (Any) at level 0" -- is that same bare entry,
// because a level-0 skill has no specialty to name.
//
// The book gives "(Any)" no other rule outside a rank row, so an event's,
// a degree's or a youth path's is resolved the same way, with the one
// difference ERRATA E-45 records: away from a skill table the level is
// granted as printed, which is what the engine does for every other skill
// a result names, and the level-0 sentence is read as p. 118's own.
func (g *Generator) grantAny(effect career.Effect, level, cause int) (Skill, error) {
	if !g.onSkillTable {
		g.char.Provenance.Deviate("E-45")
	}

	held := g.char.State.Has(effect.Skill)

	if level == 0 || (g.onSkillTable && !held) {
		if held && level == 0 {
			g.consequence(ConsequenceSkill, cause,
				effect.Detail+": the character already holds "+effect.Skill, "")

			return Skill{Name: effect.Skill}, nil
		}

		return g.gainSkill(effect.Skill, "", 0, cause), nil
	}

	options := g.specialtiesOf(effect.Skill)
	if len(options) == 0 {
		g.unimplemented(cause, effect.Detail+" -- the setting names no language to learn (p. 41)")

		return Skill{Name: effect.Skill}, nil
	}

	specialty, err := g.chooseAny(effect.Skill, options, skillStepCite)
	if err != nil {
		return Skill{}, err
	}

	return g.gainSkill(effect.Skill, specialty, level, cause), nil
}

// chooseAny puts an "(Any)" to the decider. It is the same question as a
// printed list of specialties, so it is the same choice point.
func (g *Generator) chooseAny(skill string, options []string, cite string) (string, error) {
	if len(options) == 1 {
		return options[0], nil
	}

	chosen, err := g.choose(Choice{
		Point:   pointSkillSpecialty,
		Prompt:  "Choose a specialty for " + skill + " (Any)",
		Options: options,
		Cite:    cite,
	})
	if err != nil {
		return "", err
	}

	return options[chosen], nil
}

// specialtiesOf is what "(Any)" chooses from: the specialties pp. 304-314
// print beneath the skill, in that order, then any the character holds that
// the list does not print -- Science and Trade print examples, "by no means
// a full list".
//
// Language prints none. Its specialties are "any language which can be
// chosen as a primary language", which is the setting's to say, and the
// character's own is not one: "the Language specialty chosen for the skill
// must be different from the character's Primary Language" (p. 41). ERRATA
// E-45 reads the setting's primary languages as the list.
func (g *Generator) specialtiesOf(name string) []string {
	var options []string

	if definition, found := career.SkillByName(name); found {
		options = slices.Clone(definition.Specialties)
	}

	if name == "Language" {
		g.char.Provenance.Deviate("E-45")

		options = g.languagesOfTheSetting()
	}

	for _, held := range g.char.State.Skills {
		if held.Name == name && held.Specialty != "" && !slices.Contains(options, held.Specialty) {
			options = append(options, held.Specialty)
		}
	}

	return options
}

// languagesOfTheSetting is every primary language the setting's worlds
// list, in the order the file lists them and without repeats, less the
// character's own.
func (g *Generator) languagesOfTheSetting() []string {
	var languages []string

	if g.setting == nil {
		return languages
	}

	for _, subsector := range g.setting.Subsectors {
		for _, world := range subsector.Worlds {
			for _, language := range world.PrimaryLanguages {
				if language != g.primaryLanguage && !slices.Contains(languages, language) {
					languages = append(languages, language)
				}
			}
		}
	}

	return languages
}
