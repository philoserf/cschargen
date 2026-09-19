package chargen

import (
	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/dice"
)

// enterCareer is Steps 9 to 11 (pp. 104-112): choose a career, enlist in
// it, and choose an assignment.
func (g *Generator) enterCareer() error {
	step := g.log.Step("Step 9: Choose a Career", "p. 104")

	chosen, forced, err := g.chooseCareer(step)
	if err != nil {
		return err
	}

	if g.stopped {
		return nil
	}

	if !forced && !g.enlist(chosen) {
		return nil
	}

	return g.beginService(chosen)
}

// chooseCareer resolves which career the character attempts. A pending
// transfer overrides the choice: the rules named the career, and Prisoner
// cannot be chosen at all (p. 111).
func (g *Generator) chooseCareer(step int) (career.Career, bool, error) {
	if g.transfer != nil {
		pending := g.transfer

		g.transfer = nil

		found, ok := career.ByName(pending.Career)
		if !ok {
			// The rules named a career this milestone does not implement.
			// Generation ends here rather than continuing as though the
			// transfer had not been ordered: the character really is in
			// that career now, and nothing the engine did next would be
			// the rules being followed (docs/MILESTONE-1.md).
			g.unimplemented(step, pending.Detail+
				" -- generation ends here; this career is not implemented in milestone 1")

			g.stopped = true

			return career.Career{}, false, nil
		}

		g.forcedAssignment = pending.Assignment
		g.forcedTerms = pending.Terms

		return found, true, nil
	}

	if g.failedEnlistments >= consecutiveFailuresToVagabond {
		g.consequence(ConsequenceCareer, step,
			"three consecutive failed enlistments: enter the Vagabond career", "Vagabond")

		g.failedEnlistments = 0

		return career.Vagabond(), true, nil
	}

	eligible := g.eligibleCareers()
	if len(eligible) == 0 {
		return career.Vagabond(), true, nil
	}

	names := make([]string, len(eligible))
	for i, def := range eligible {
		names[i] = def.Name
	}

	index, err := g.choose(Choice{
		Point:   "career",
		Prompt:  "Choose a career to attempt",
		Options: names,
		Cite:    "p. 107",
	})
	if err != nil {
		return career.Career{}, false, err
	}

	return eligible[index], false, nil
}

// eligibleCareers is every implemented career the character may attempt
// now. Prisoner is never among them: "A character cannot voluntarily choose
// to enter this career" (p. 111). A career the character failed to enter
// stays out for two full terms (p. 110).
func (g *Generator) eligibleCareers() []career.Career {
	var eligible []career.Career

	for _, def := range career.All() {
		if def.Name == "Prisoner" {
			continue
		}

		if g.forced != "" && def.Name != g.forced {
			continue
		}

		if until, locked := g.lockout[def.Name]; locked && len(g.char.State.Terms) < until {
			continue
		}

		eligible = append(eligible, def)
	}

	return eligible
}

// consecutiveFailuresToVagabond is p. 110's rule: "If a character fails
// three consecutive enlistment attempts, regardless of which careers were
// chosen, they automatically enter the Vagabond career."
const consecutiveFailuresToVagabond = 3

// enlistmentLockoutTerms is p. 110's other rule: "A character may not
// attempt to enter the same career again for two full terms."
const enlistmentLockoutTerms = 2

// enlist is Step 10 (p. 110).
func (g *Generator) enlist(target career.Career) bool {
	step := g.log.Step("Step 10: Enlist in a Career", "p. 110")

	if target.Enlistment == nil {
		g.consequence(ConsequenceCareer, step, target.EnlistmentNote, target.Name)

		return true
	}

	mods := g.takeModifiers("next enlistment attempt")

	which, ok := characteristicByName(target.Enlistment.Characteristic)
	if ok {
		mods = append(mods, dice.Mod{
			Name:  target.Enlistment.Characteristic,
			Value: g.char.State.Characteristics.Modifier(which),
		})
	}

	throw := g.dice.Throw(target.Enlistment.Number, mods...)
	cause := g.log.Throw(throw, "p. 110")

	if throw.Success {
		g.failedEnlistments = 0
		g.consequence(ConsequenceCareer, cause, "accepted into "+target.Name, target.Name)

		return true
	}

	g.failedEnlistments++

	if g.lockout == nil {
		g.lockout = map[string]int{}
	}

	g.lockout[target.Name] = len(g.char.State.Terms) + enlistmentLockoutTerms

	g.consequence(ConsequenceCareer, cause,
		"rejected by "+target.Name+"; it cannot be attempted again for two terms", target.Name)

	return false
}

// beginService is Step 11 (p. 112) plus the two things that happen on
// entering a career: the Rank 0 benefit, and the first term's level-0
// service skills.
func (g *Generator) beginService(entered career.Career) error {
	step := g.log.Step("Step 11: Choose an Assignment", "p. 112")

	assignment, err := g.chooseAssignment(entered)
	if err != nil {
		return err
	}

	g.career = &entered
	g.assignment = assignment
	g.rank = 0
	g.commissioned = false
	g.termsInCareer = 0
	g.careerBenefitMod = 0
	g.cite = entered.Cite

	g.char.State.Services = append(g.char.State.Services, Service{
		Career:     entered.Name,
		Assignment: assignment.Name,
	})

	// "If this is the character's first term in this career, the character
	// gains all the skills in the Service Skills table at level 0 that they
	// do not already possess at level 0 or higher" (p. 117).
	table, ok := entered.Table(career.ServiceSkills)
	if ok {
		for _, row := range table.Rows {
			if row.Kind != career.EffectSkill || g.char.State.Has(row.Skill) {
				continue
			}

			granted := g.char.State.GainSkill(row.Skill, specialtyOf(row), 0)
			g.log.Consequence(ConsequenceEvent{
				Kind:   ConsequenceSkill,
				Cause:  step,
				Detail: "first term in the career: " + granted.Full() + " at level 0",
				Skill:  granted.Full(),
				Cite:   skillStepCite,
			})
		}
	}

	// ERRATA E-8: rank benefits "are applied immediately upon achieving
	// that Rank" (p. 116), and every character "begins a career at Rank 0"
	// (p. 115), so the Rank 0 row applies on entry.
	g.char.Provenance.Deviate("E-8")

	return g.applyAll(assignment.Ranks[0], step)
}

func (g *Generator) chooseAssignment(entered career.Career) (career.Assignment, error) {
	if g.forcedAssignment != "" {
		found, ok := entered.Assignment(g.forcedAssignment)

		g.forcedAssignment = ""

		if ok {
			return found, nil
		}
	}

	if len(entered.Assignments) == 1 {
		return entered.Assignments[0], nil
	}

	names := make([]string, len(entered.Assignments))
	for i, a := range entered.Assignments {
		names[i] = a.Name
	}

	index, err := g.choose(Choice{
		Point:   "assignment",
		Prompt:  "Choose an assignment within " + entered.Name,
		Options: names,
		Cite:    "p. 112",
	})
	if err != nil {
		return career.Assignment{}, err
	}

	return entered.Assignments[index], nil
}

func specialtyOf(effect career.Effect) string {
	if len(effect.Specialties) == 1 {
		return effect.Specialties[0]
	}

	return ""
}
