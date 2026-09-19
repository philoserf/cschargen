package chargen

import (
	"strings"

	"github.com/philoserf/cschargen/setting"
)

// universalSkill is p. 41's one line that applies to everyone: "All
// characters, regardless of origin, will also receive Electronics 0 as a
// Background Skill."
const universalSkill = "Electronics"

// backgroundCite is p. 40, the page that governs background skills and the
// choices they offer. Three call sites reach it.
const backgroundCite = "p. 40"

// grantBackground gives the character what their homeworld taught them:
// background skills at level 1, a primary language, and the Electronics 0
// every character has regardless of origin (pp. 40-41).
func (g *Generator) grantBackground(world setting.World, step int) error {
	err := g.choosePrimaryLanguage(world, step)
	if err != nil {
		return err
	}

	for _, requirement := range world.BackgroundSkills {
		err = g.grantRequirement(requirement, step)
		if err != nil {
			return err
		}
	}

	granted := g.char.State.GainSkill(universalSkill, "", 0)

	g.log.Consequence(ConsequenceEvent{
		Kind:   ConsequenceSkill,
		Cause:  step,
		Detail: "every character, regardless of origin: " + granted.Full() + " at level 0",
		Skill:  granted.Full(),
		Cite:   "p. 41",
	})

	return nil
}

// choosePrimaryLanguage takes the homeworld's language. "Some worlds list
// more than one language ... the most common language spoken on that world
// is listed first, followed by other widely used languages in descending
// order of commonality" (p. 41), so the order is meaningful and the policy
// takes the first.
func (g *Generator) choosePrimaryLanguage(world setting.World, step int) error {
	if len(world.PrimaryLanguages) == 0 {
		return ErrNoLanguage
	}

	chosen := 0

	if len(world.PrimaryLanguages) > 1 {
		index, err := g.choose(Choice{
			Point:   "primary_language",
			Prompt:  "Choose a primary language for " + world.Name,
			Options: world.PrimaryLanguages,
			Cite:    "p. 41",
		})
		if err != nil {
			return err
		}

		chosen = index
	}

	g.primaryLanguage = world.PrimaryLanguages[chosen]
	g.char.State.Language = g.primaryLanguage

	g.consequence(ConsequenceSkill, step,
		"primary language: "+g.primaryLanguage, "")

	return nil
}

// grantRequirement gives one background skill, resolving the choice where
// the page offers one.
func (g *Generator) grantRequirement(requirement setting.Requirement, step int) error {
	chosen := 0

	if len(requirement.OneOf) > 1 {
		labels := make([]string, len(requirement.OneOf))
		for i, alternative := range requirement.OneOf {
			labels[i] = describe(alternative)
		}

		index, err := g.choose(Choice{
			Point:   "background_skill",
			Prompt:  "Choose a background skill",
			Options: labels,
			Cite:    backgroundCite,
		})
		if err != nil {
			return err
		}

		chosen = index
	}

	return g.grantAlternative(requirement.OneOf[chosen], step)
}

func (g *Generator) grantAlternative(alternative setting.Alternative, step int) error {
	if alternative.Item != "" {
		g.char.State.Stash = append(g.char.State.Stash, alternative.Item)
		g.consequence(ConsequenceStash, step, "from the homeworld: "+alternative.Item, "")

		return nil
	}

	specialty, err := g.chooseSpecialty(alternative, step)
	if err != nil {
		return err
	}

	granted := g.char.State.GainSkill(alternative.Skill, specialty, 1)

	g.log.Consequence(ConsequenceEvent{
		Kind:   ConsequenceSkill,
		Cause:  step,
		Detail: "background skill: " + granted.Full() + " " + itoa(granted.Level),
		Skill:  granted.Full(),
		Level:  granted.Level,
		Cite:   backgroundCite,
	})

	return nil
}

// chooseSpecialty resolves a specialty list. Where the skill granted is
// Language, the primary language is not among the options: "the Language
// specialty chosen for the skill must be different from the character's
// Primary Language. This represents learning an additional language beyond
// one's native tongue" (p. 41).
func (g *Generator) chooseSpecialty(alternative setting.Alternative, step int) (string, error) {
	options := alternative.Specialties

	if alternative.Skill == "Language" {
		options = withoutPrimary(options, g.primaryLanguage)
		if len(options) == 0 {
			g.unimplemented(step,
				"the homeworld grants Language but offers no specialty other than the primary language (p. 41)")

			return "", nil
		}
	}

	switch len(options) {
	case 0:
		return "", nil
	case 1:
		return options[0], nil
	}

	index, err := g.choose(Choice{
		Point:   "background_specialty",
		Prompt:  "Choose a specialty for " + alternative.Skill,
		Options: options,
		Cite:    backgroundCite,
	})
	if err != nil {
		return "", err
	}

	return options[index], nil
}

func withoutPrimary(options []string, primary string) []string {
	kept := make([]string, 0, len(options))

	for _, option := range options {
		if option != primary {
			kept = append(kept, option)
		}
	}

	return kept
}

// describe renders an alternative for a choice prompt, the way the page
// writes it.
func describe(alternative setting.Alternative) string {
	if alternative.Item != "" {
		return alternative.Item
	}

	if len(alternative.Specialties) == 0 {
		return alternative.Skill
	}

	return alternative.Skill + " (" + strings.Join(alternative.Specialties, " or ") + ")"
}
