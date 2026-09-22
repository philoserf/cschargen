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
		return g.takePendingTransfer(step)
	}

	if g.enslaved {
		return g.ownedFirstCareer(step), true, nil
	}

	if g.failedEnlistments >= consecutiveFailuresToVagabond {
		g.consequence(ConsequenceCareer, step,
			"three consecutive failed enlistments: enter the Vagabond career", "Vagabond")

		g.failedEnlistments = 0

		return career.Vagabond(), true, nil
	}

	eligible := g.eligibleCareers()
	if len(eligible) == 0 {
		g.consequence(ConsequenceCareer, step,
			"no career will have them"+g.whyNothingIsOpen()+
				": enter the Vagabond career (p. 111)", career.Vagabond().Name)

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

// previousCareerName is "the career you held before this one", which one
// mishap sends a character back to. The service record is in order, so it
// is the one before the last -- and with no such career the name is empty,
// which chooseCareer reports as a destination it cannot find rather than
// silently choosing one.
func (g *Generator) previousCareerName() string {
	services := g.char.State.Services

	const beforeTheLast = 2

	if len(services) < beforeTheLast {
		return ""
	}

	return services[len(services)-beforeTheLast].Career
}

// whyNothingIsOpen names the reason the career list came back empty, so the
// record says why the character drifted rather than leaving a step with
// nothing under it.
//
// p. 111 is clear that the Vagabond itself is the rules working -- "most
// characters enter the Vagabond career when circumstances leave them with no
// viable alternative ... it ensures that characters always have a path
// forward" -- which is why the silence was the bug and not the destination.
//
// Only two things empty the list. An aging crisis and a mental
// characteristic at 0 do not: eligibleCareers answers both with Vagabond,
// which is a list of one rather than a list of none, and the ordinary choice
// point records it.
func (g *Generator) whyNothingIsOpen() string {
	if g.forced != "" {
		return " but " + g.forced + ", which the flags asked for and which is closed to them"
	}

	return ", every career they might attempt having turned them down"
}

// takePendingTransfer is the career a mishap or an event named, arrived at.
// The result that ordered it said where they were going; this is the step
// where they get there, and a step that records nothing reads as a
// rendering fault rather than as a character's life.
func (g *Generator) takePendingTransfer(step int) (career.Career, bool, error) {
	pending := g.transfer

	g.transfer = nil

	name := pending.Career
	if name == career.PreviousCareer {
		name = g.previousCareerName()
	}

	found, ok := career.ByName(name)
	if !ok {
		// All thirty-four careers the book names are transcribed, so this
		// is now a misspelled destination rather than a missing career.
		// Generation still ends here rather than continuing as though the
		// transfer had not been ordered: the character really is in that
		// career now, and nothing the engine did next would be the rules
		// being followed.
		g.unimplemented(step, pending.Detail+
			" -- generation ends here; no career of that name is transcribed")

		g.stopped = true

		return career.Career{}, false, nil
	}

	g.forcedAssignment = pending.Assignment
	g.forcedTerms = pending.Terms

	g.consequence(ConsequenceCareer, step,
		"entered "+found.Name+" without an enlistment roll: "+pending.Detail,
		found.Name)

	return found, true, nil
}

// ownedFirstCareer is p. 42's rule: an engineered or uplift character born
// on a world where they are enslaved must take the slave career as their
// first career term.
//
// It is the only way into that career. The career prints no enlistment
// throw, cannot be chosen, and says its entry is "a result of the early
// life tables" -- which is the same rule seen from the other side.
func (g *Generator) ownedFirstCareer(step int) career.Career {
	g.enslaved = false

	g.consequence(ConsequenceCareer, step,
		"enslaved on their homeworld: the first career is not a choice (p. 42)",
		career.Slave().Name)

	return career.Slave()
}

// eligibleCareers is every implemented career the character may attempt
// now. Two are never among them, for the same reason worded twice.
//
// Prisoner: "A character cannot voluntarily choose to enter this career"
// (p. 111). The slave career: the book gives it no enlistment throw and
// reaches it only from p. 42, which is why ownedFirstCareer above says it
// "is the only way into that career" -- a claim this list has to hold up.
// Neither career fails an enlistment, so either one left in would be
// entered by whoever drew it.
//
// A career the character failed to enter stays out for two full terms
// (p. 110).
func (g *Generator) eligibleCareers() []career.Career {
	// pp. 123-124: a character who has survived an Aging Crisis
	// "automatically fails all future Enlistment checks", and one with a
	// mental characteristic at 0 "may not attempt further Enlistment
	// checks" at all. Both leave the same three options, and two of them --
	// continue in the current career, end generation -- are not this
	// function's to offer: it is only reached with no career in hand, and
	// the term limit is what ends generation. What remains is Vagabond,
	// which takes no enlistment throw, so the automatic failure never has
	// to be rolled for.
	if g.crisisSurvived || g.mentalDecline {
		return []career.Career{career.Vagabond()}
	}

	var eligible []career.Career

	for _, def := range career.All() {
		if def.Name == "Prisoner" || def.Name == career.Slave().Name {
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

	// "You may enlist automatically", which several results grant and one
	// narrows to a class of career.
	if g.takeAutomaticFor(enlistmentThrow, target) {
		g.failedEnlistments = 0
		g.consequence(ConsequenceCareer, step,
			"accepted into "+target.Name+" without a throw", target.Name)

		return true
	}

	mods := g.takeModifiersFor(enlistmentThrow, target)

	mods = append(mods, g.enlistmentMods(target)...)

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

// enlistmentMods applies the situational modifiers a career prints on its
// enlistment throw (p. 111). The switch is exhaustive over EnlistmentModKind,
// so a kind added there and not handled here is a build failure rather than a
// modifier silently dropped.
func (g *Generator) enlistmentMods(target career.Career) []dice.Mod {
	var mods []dice.Mod

	for _, mod := range target.EnlistmentMods {
		switch mod.Kind {
		case career.PerPreviousCareer:
			entered := len(g.char.State.Services)
			if entered > 0 {
				mods = append(mods, dice.Mod{
					Name:  "previous careers",
					Value: mod.Value * entered,
				})
			}
		case career.ApparentAgeOver40:
			if apparentAgeOverForty(g.char.State.ApparentAge) {
				mods = append(mods, dice.Mod{
					Name:  "apparent age " + g.char.State.ApparentAge.String(),
					Value: mod.Value,
				})
			}
		case career.UndergraduateDegree:
			mods = append(mods, g.degreeModifier(mod, career.Bachelors, "a degree")...)
		case career.GraduateDegree:
			mods = append(mods, g.graduateModifier(mod)...)
		case career.MedicalSchool:
			mods = append(mods, g.degreeModifier(mod, career.MedicalDoctor, "medical school")...)
		}
	}

	return mods
}

// graduateModifier is the graduate education bonus, taken once however far
// the character took the degree.
//
// A doctorate is a master's gone further rather than a second degree:
// degreeFor awards it on a second success at graduate school and leaves the
// master's on the record, so a character can hold both. p. 212 prints one
// modifier, so the highest held is the one that applies -- reading the
// record for each in turn and adding both gave a doctorate-holder +8 where
// the page prints +4 (#132).
func (g *Generator) graduateModifier(mod career.EnlistmentMod) []dice.Mod {
	if held := g.degreeModifier(mod, career.Doctorate, "a doctorate"); held != nil {
		return held
	}

	return g.degreeModifier(mod, career.Masters, "a master's")
}

// degreeModifier is one education bonus, applied where the character holds
// the degree it names.
func (g *Generator) degreeModifier(
	mod career.EnlistmentMod, degree career.Degree, name string,
) []dice.Mod {
	for _, held := range g.char.State.Education {
		if held.Degree == degree {
			return []dice.Mod{{Name: name, Value: mod.Value}}
		}
	}

	return nil
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

	g.service = serviceState{career: &entered, assignment: assignment}
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

	return g.applyRankBenefits(ranksFor(assignment, false)[0], step)
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
