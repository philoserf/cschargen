package chargen

import (
	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/dice"
)

// serveTerm is Steps 12 to 16 (pp. 112-121): survival, mishap, advancement,
// skills, events -- one four-year term.
func (g *Generator) serveTerm() error {
	if g.career == nil {
		return ErrNoCareer
	}

	assignment, ok := g.career.Assignment(g.assignment.Name)
	if !ok {
		return ErrNoCareer
	}

	g.termsInCareer++

	term := Term{
		Career:     g.career.Name,
		Assignment: assignment.Name,
		Number:     len(g.char.State.Terms) + 1,
	}

	survived, err := g.rollSurvival(assignment)
	if err != nil {
		return err
	}

	term.Survived = survived

	err = g.resolveTerm(assignment, term.Survived)
	if err != nil {
		return err
	}

	term.Rank = g.rank
	g.char.State.Terms = append(g.char.State.Terms, term)

	if service, found := g.char.State.Service(g.career.Name); found {
		service.Terms = g.termsInCareer
		service.Rank = g.rank
	}

	return g.age()
}

// resolveTerm is the half of a term that depends on the survival throw: a
// mishap, or advancement, a skill and an event.
func (g *Generator) resolveTerm(assignment career.Assignment, survived bool) error {
	if !survived {
		g.log.Step("Step 13: Roll A Mishap", "p. 113")

		return g.rollMishap(true)
	}

	err := g.advance(assignment)
	if err != nil {
		return err
	}

	err = g.rollSkill(assignment)
	if err != nil {
		return err
	}

	return g.rollEvent()
}

// rollSurvival is Step 12 (p. 112). A natural twelve before modifiers
// requires another term in the same career (p. 113).
func (g *Generator) rollSurvival(assignment career.Assignment) (bool, error) {
	step := g.log.Step("Step 12: Roll for Survival", "p. 112")

	if g.takeAutomaticFor(survivalThrow, *g.career) {
		g.consequence(ConsequenceCareer, step,
			"an automatic success on the survival roll, granted earlier", g.career.Name)

		return true, nil
	}

	mods := g.takeCareerModifiers(survivalThrow)

	spent, err := g.spendPools(survivalThrow, step)
	if err != nil {
		return false, err
	}

	throw := g.characteristicThrow(assignment.Survival, append(mods, spent...)...)
	cause := g.log.Throw(throw, "p. 112")

	if throw.Natural() == naturalTwelve {
		g.mustContinue = true
		g.consequence(ConsequenceCareer, cause,
			"a natural twelve on the survival roll: another term in this career", g.career.Name)
	}

	if !throw.Success {
		g.consequence(ConsequenceCareer, cause, "the survival roll failed", g.career.Name)

		return false, nil
	}

	return true, nil
}

// naturalTwelve is the exceptional success the book reads as a result in
// its own right (p. 113).
const naturalTwelve = 12

// advance is Step 14 (pp. 114-116): the optional commission, then the
// advancement throw.
func (g *Generator) advance(assignment career.Assignment) error {
	step := g.log.Step("Step 14: Roll for Advancement", "p. 114")

	err := g.offerCommission()
	if err != nil {
		return err
	}

	if g.autoAdvance {
		g.autoAdvance = false

		g.takeOnFailure(advancementThrow)

		return g.promote(step, assignment)
	}

	// "You are suspended, and your next Advancement roll fails
	// automatically": the throw is not made, and what waits on its failure
	// still fires.
	if g.takeAutoFailure(advancementThrow) {
		g.consequence(ConsequenceRank, step,
			"the advancement roll fails automatically, decided earlier", g.career.Name)

		return g.applyAll(failureEffects(g.takeOnFailure(advancementThrow)), step)
	}

	mods := g.takeCareerModifiers(advancementThrow)

	spent, err := g.spendPools(advancementThrow, step)
	if err != nil {
		return err
	}

	mods = append(mods, spent...)

	which, ok := characteristicByName(assignment.Advancement.Characteristic)
	if ok {
		mods = append(mods, dice.Mod{
			Name:  assignment.Advancement.Characteristic,
			Value: g.char.State.Characteristics.Modifier(which),
		})
	}

	throw := g.dice.Throw(assignment.Advancement.Number, mods...)
	cause := g.log.Throw(throw, "p. 115")

	if !throw.Success {
		return g.applyAll(failureEffects(g.takeOnFailure(advancementThrow)), cause)
	}

	g.takeOnFailure(advancementThrow)

	return g.promote(cause, assignment)
}

// failureEffects is the effects nine results attach to a failed advancement
// roll, flattened out of the hooks that carried them.
func failureEffects(hooks []career.Effect) []career.Effect {
	found := make([]career.Effect, 0, len(hooks))
	for _, hook := range hooks {
		found = append(found, hook.Failure...)
	}

	return found
}

// offerCommission is the commission half of Step 14 (p. 114). "This step is
// optional. If the character has already received a Commission ... or does
// not wish to pursue officer rank, this step is skipped", and "If the
// Commission roll fails, there is no penalty."
func (g *Generator) offerCommission() error {
	if g.career.Commission == nil || g.commissioned {
		return nil
	}

	index, err := g.choose(Choice{
		Point:   "commission",
		Prompt:  "Attempt a commission in " + g.career.Name + "?",
		Options: []string{"attempt it", "stay enlisted"},
		Cite:    "p. 114",
	})
	if err != nil {
		return err
	}

	if index != 0 {
		return nil
	}

	return g.attemptCommission(0, g.log.Len())
}

// attemptCommission makes the throw, whether Step 14 offered it or an event
// granted one with a modifier.
func (g *Generator) attemptCommission(modifier, cause int) error {
	if g.career == nil || g.career.Commission == nil {
		g.unimplemented(cause, "a commission in a career that offers none")

		return nil
	}

	if g.commissioned {
		return nil
	}

	if g.takeAutomaticFor(commissionThrow, *g.career) {
		return g.grantCommission(cause)
	}

	mods := []dice.Mod{}
	if modifier != 0 {
		mods = append(mods, dice.Mod{Name: "event", Value: modifier})
	}

	which, ok := characteristicByName(g.career.Commission.Characteristic)
	if ok {
		mods = append(mods, dice.Mod{
			Name:  g.career.Commission.Characteristic,
			Value: g.char.State.Characteristics.Modifier(which),
		})
	}

	throw := g.dice.Throw(g.career.Commission.Number, mods...)
	thrown := g.log.Throw(throw, "p. 114")

	if !throw.Success {
		return nil
	}

	return g.grantCommission(thrown)
}

// grantCommission moves the character onto the officer track, whether the
// throw made it or a table result decided it.
func (g *Generator) grantCommission(thrown int) error {
	g.commissioned = true

	if service, found := g.char.State.Service(g.career.Name); found {
		service.Commissioned = true
	}

	g.consequence(ConsequenceRank, thrown, "commissioned in "+g.career.Name, g.career.Name)

	// A commission resets rank to the officer track's Rank 0, whose benefit
	// is granted like any other (p. 116).
	g.rank = 0

	assignment, found := g.career.Assignment(g.assignment.Name)
	if !found {
		return nil
	}

	return g.applyAll(ranksFor(assignment, true)[0], thrown)
}

// ranksFor is the rank table that applies: the officer track once the
// character is commissioned, the enlisted one otherwise. A career with a
// commission prints two per assignment (pp. 225, 234).
func ranksFor(assignment career.Assignment, commissioned bool) [][]career.Effect {
	if commissioned && assignment.OfficerRanks != nil {
		return assignment.OfficerRanks
	}

	return assignment.Ranks
}

// tableInAnotherCareer finds a named assignment's skill table in a career
// the character is not in. Colonist event 54 is the only result in the book
// that asks for one (p. 176).
func tableInAnotherCareer(effect career.Effect) (career.SkillTable, bool) {
	other, found := career.ByName(effect.Career)
	if !found {
		return career.SkillTable{}, false
	}

	assignment, found := other.Assignment(effect.Assignment)
	if !found {
		return career.SkillTable{}, false
	}

	return assignment.Skills, true
}

// rollNamedTable carries out "make a roll on the <named> table".
func (g *Generator) rollNamedTable(effect career.Effect, cause int) error {
	if g.career == nil {
		return ErrNoCareer
	}

	assignment, found := g.career.Assignment(g.assignment.Name)
	if !found {
		return ErrNoCareer
	}

	table, found, err := g.namedTable(effect, assignment)
	if err != nil {
		return err
	}

	if !found {
		g.unimplemented(cause, effect.Detail+" -- this career prints no such table")

		return nil
	}

	roll := g.dice.D6()
	thrown := g.log.Roll(roll, skillStepCite)

	return g.apply(table.Rows[roll.Total-1], thrown)
}

// namedTable finds the table an effect asks for, including "an assignment
// other than your own", which is a choice where the career has more than
// two assignments.
func (g *Generator) namedTable(
	effect career.Effect, own career.Assignment,
) (career.SkillTable, bool, error) {
	// A table in another career, which one event asks for by name.
	if effect.Career != "" {
		table, found := tableInAnotherCareer(effect)

		return table, found, nil
	}

	if !effect.OtherAssignment {
		table, found := g.career.Table(effect.Table)

		return table, found, nil
	}

	var others []career.Assignment

	for _, other := range g.career.Assignments {
		if other.Name != own.Name {
			others = append(others, other)
		}
	}

	if len(others) == 0 {
		return career.SkillTable{}, false, nil
	}

	index := 0

	if len(others) > 1 {
		names := make([]string, len(others))
		for i, other := range others {
			names[i] = other.Name
		}

		chosen, err := g.choose(Choice{
			Point:   "other_assignment",
			Prompt:  "Choose another assignment's skill table",
			Options: names,
			Cite:    skillStepCite,
		})
		if err != nil {
			return career.SkillTable{}, false, err
		}

		index = chosen
	}

	return others[index].Skills, true, nil
}

// promote raises the rank by one and applies that rank's benefit, which is
// granted "immediately upon achieving that Rank" (p. 116).
func (g *Generator) promote(cause int, assignment career.Assignment) error {
	ranks := ranksFor(assignment, g.commissioned)

	if g.rank >= len(ranks)-1 {
		g.consequence(ConsequenceRank, cause,
			"already at the highest printed rank in this assignment", g.career.Name)

		return nil
	}

	g.rank++

	g.consequence(ConsequenceRank, cause,
		"advanced to rank "+itoa(g.rank)+" in "+assignment.Name, g.career.Name)

	return g.applyRankBenefits(ranks[g.rank], cause)
}

// changeRank carries out an [career.EffectRank]: a promotion granted
// without an advancement throw, or a demotion.
//
// A promotion goes through promote, so the rank's printed benefits are
// granted: p. 116 attaches them to holding the rank, not to the throw that
// reached it, and an event that says "you are promoted" reaches it.
//
// A demotion grants nothing back and takes nothing away. Two results add
// "retaining any benefit already gained" and fourteen do not; ERRATA E-36
// reads them as the same thing, because the book has no rule anywhere for
// un-granting a skill.
func (g *Generator) changeRank(effect career.Effect, cause int) error {
	if g.career == nil {
		g.unimplemented(cause, effect.Detail+" -- outside a career, where there is no rank")

		return nil
	}

	levels := effect.Levels
	if levels == 0 {
		levels = 1
	}

	if levels > 0 {
		for range levels {
			err := g.promote(cause, g.assignment)
			if err != nil {
				return err
			}
		}

		return nil
	}

	g.demote(-levels, cause)

	return nil
}

// demote lowers the rank, stopping at the bottom of the table rather than
// going below it. A character already at rank 0 loses nothing, and the
// record says so: the result fired, and what it did is what happened.
func (g *Generator) demote(levels, cause int) {
	before := g.rank

	g.rank = max(g.rank-levels, 0)

	if service, found := g.char.State.Service(g.career.Name); found {
		service.Rank = g.rank
	}

	if g.rank == before {
		g.consequence(ConsequenceRank, cause,
			"already at the lowest rank in "+g.assignment.Name, g.career.Name)

		return
	}

	g.consequence(ConsequenceRank, cause,
		"reduced to rank "+itoa(g.rank)+" in "+g.assignment.Name, g.career.Name)
}

// applyRankBenefits is p. 116's rule, which an ordinary applyAll gets
// wrong:
//
//	"If a Benefit grants a skill the character already possesses at the
//	listed level or higher, no additional benefit is gained. Rank Benefits
//	represent required competence, not bonus stacking."
//
// So a rank benefit is a floor rather than an increment. Everything else a
// rank row can carry -- a characteristic, a Contact, a stash -- is applied
// as printed; it is only the skills the rule speaks about.
func (g *Generator) applyRankBenefits(effects []career.Effect, cause int) error {
	for _, effect := range effects {
		if effect.Kind != career.EffectSkill {
			err := g.apply(effect, cause)
			if err != nil {
				return err
			}

			continue
		}

		err := g.grantRankSkill(effect, cause)
		if err != nil {
			return err
		}
	}

	return nil
}

// grantRankSkill raises a skill to the level the rank row prints, or does
// nothing where the character is already there.
//
// "(Any)" is its own rule: "This allows the character to select a specialty
// within that skill, provided they do not already possess that specialty at
// level 1 or higher ... If the character already possesses all available
// specialties at level 1 or higher, they gain no additional benefit"
// (p. 116). The engine records "Any" as the specialty rather than resolving
// it, so a second grant of the same skill at "Any" is one the character
// already has.
func (g *Generator) grantRankSkill(effect career.Effect, cause int) error {
	level := effect.Level
	if level == 0 {
		level = 1
	}

	specialty, err := g.rankSpecialty(effect)
	if err != nil {
		return err
	}

	held, found := g.char.State.Skill(effect.Skill, specialty)
	if found && held.Level >= level {
		name := ""
		if g.career != nil {
			name = g.career.Name
		}

		g.consequence(ConsequenceSkill, cause,
			"rank benefit: "+held.Full()+" is already "+itoa(held.Level)+
				", which meets the required "+itoa(level), name)

		return nil
	}

	got := g.char.State.GainSkill(effect.Skill, specialty, level-held.Level)

	g.log.Consequence(ConsequenceEvent{
		Kind:   ConsequenceSkill,
		Cause:  cause,
		Detail: "rank benefit: " + got.Full() + " " + itoa(got.Level),
		Skill:  got.Full(),
		Level:  got.Level,
		Cite:   g.cite,
	})

	return nil
}

// rankSpecialty resolves a rank benefit's specialty the way applySkill
// does, except that a choice excludes what the character already holds at
// level 1 or higher (p. 116).
func (g *Generator) rankSpecialty(effect career.Effect) (string, error) {
	switch len(effect.Specialties) {
	case 0:
		return "", nil
	case 1:
		return effect.Specialties[0], nil
	}

	var open []string

	for _, specialty := range effect.Specialties {
		held, found := g.char.State.Skill(effect.Skill, specialty)
		if found && held.Level >= 1 {
			continue
		}

		open = append(open, specialty)
	}

	// "If the character already possesses all available specialties at level
	// 1 or higher, they gain no additional benefit from that Rank Benefit."
	// Falling back to the printed list is what makes that a no-op: the
	// grant below finds the skill already at or above its level.
	if len(open) == 0 {
		open = effect.Specialties
	}

	chosen, err := g.choose(Choice{
		Point:   "skill_specialty",
		Prompt:  "Choose a specialty for " + effect.Skill,
		Options: open,
		Cite:    "p. 116",
	})
	if err != nil {
		return "", err
	}

	return open[chosen], nil
}

// rollSkill is Step 15 (p. 117): "The Player should choose one of the
// tables and roll a d6."
func (g *Generator) rollSkill(assignment career.Assignment) error {
	g.log.Step("Step 15: Roll for Skills", skillStepCite)

	tables := g.availableSkillTables(assignment)

	names := make([]string, len(tables))
	for i, table := range tables {
		names[i] = table.Name
	}

	index, err := g.choose(Choice{
		Point:   "skill_table",
		Prompt:  "Choose a skill table to roll on",
		Options: names,
		Cite:    skillStepCite,
	})
	if err != nil {
		return err
	}

	roll := g.dice.D6()
	cause := g.log.Roll(roll, skillStepCite)

	return g.apply(tables[index].Rows[roll.Total-1], cause)
}

// availableSkillTables is the career's own tables, minus any the character
// is not eligible for, plus the assignment's. The Advanced Education table
// is "restricted to characters with an EDU of a certain level or higher"
// (p. 117), so an ineligible character is not offered it rather than being
// allowed to roll and then told no.
func (g *Generator) availableSkillTables(assignment career.Assignment) []career.SkillTable {
	var tables []career.SkillTable

	for _, table := range g.career.Tables {
		if table.MinimumEDU > 0 && g.char.State.Characteristics.EDU < table.MinimumEDU {
			continue
		}

		tables = append(tables, table)
	}

	return append(tables, assignment.Skills)
}

// rollEvent is Step 16 (p. 118).
func (g *Generator) rollEvent() error {
	g.log.Step("Step 16: Roll for Events", "p. 118")

	roll := g.dice.D66()
	cause := g.log.Roll(roll, "p. 118")

	row, ok := g.career.Events[roll.Total]
	if !ok {
		return ErrMissingEventRow
	}

	g.consequence(ConsequenceCareer, cause, "event: "+row.Summary, g.career.Name)

	return g.applyAll(row.Effects, cause)
}

// age is Step 17 (pp. 121-123): four years per term, or 1d3 where a mishap
// ejected the character mid-term, and then whatever aging throws the
// homeworld's tech level has made due.
func (g *Generator) age() error {
	step := g.log.Step("Step 17: Aging", "p. 121")
	cause := step

	years := termYears

	if g.ejected {
		roll := g.dice.D3()

		cause = g.log.Roll(roll, "p. 121")

		years = roll.Total
	}

	g.char.State.Age += years
	g.consequence(ConsequenceAge, cause, "age "+itoa(g.char.State.Age), "")

	g.stampApparentAge(cause)

	return g.agingThrows(step)
}

// agingThrows is the aging table itself (pp. 122-123). Each check is an
// ordinary characteristic check -- 2d6 plus that characteristic's own
// modifier -- and failing one costs a point of it.
//
// The index is the character's lifetime term count, which is why this reads
// len(State.Terms) rather than g.termsInCareer.
func (g *Generator) agingThrows(step int) error {
	term := len(g.char.State.Terms)

	checks, due := agingChecksAt(g.agingProfile, g.techLevel, term)
	if !due {
		return nil
	}

	g.consequence(ConsequenceAge, step,
		"term "+itoa(term)+", on the "+string(g.agingProfile)+
			" aging profile: checks come due", "")

	for _, check := range checks {
		target := career.Check{Characteristic: check.Characteristic, Number: check.Number}
		throw := g.characteristicThrow(target)

		rolled := g.log.Roll(throw.Roll, agingCite)

		if throw.Success {
			continue
		}

		g.adjust(check.Characteristic, -1,
			"aged: "+check.Characteristic+" "+itoa(check.Number)+"+ failed", rolled)

		err := g.agingCrisis(check.Characteristic, rolled)
		if err != nil {
			return err
		}

		if g.char.State.Fate != "" {
			return nil
		}
	}

	g.settleTerminalStates(step)

	return nil
}

// agingCrisis is p. 123: a characteristic aging has reduced to 0 leaves the
// character "on the verge of death or permanent incapacity", and 1d6 x 1000
// credits of emergency treatment restores it to 1.
//
// Two characteristics reaching 0 in the same term is two payments, which the
// page does not say and ERRATA E-18 records.
func (g *Generator) agingCrisis(which string, cause int) error {
	target, ok := characteristicByName(which)
	if !ok || g.char.State.Characteristics.Get(target) > 0 {
		return nil
	}

	g.consequence(ConsequenceAge, cause,
		"aging crisis: "+which+" is at 0, and without treatment the character dies",
		"")

	chosen, err := g.choose(Choice{
		Point:   "aging_crisis",
		Prompt:  "Pay for emergency treatment to restore " + which + " to 1?",
		Options: []string{"pay for treatment", "go without"},
		Cite:    crisisCite,
	})
	if err != nil {
		return err
	}

	if chosen != 0 {
		g.die(cause, which+" was left at 0")

		return nil
	}

	price, err := g.rollExpression(crisisPayment, crisisCite)
	if err != nil {
		return err
	}

	// ERRATA E-19: the book sets a price and never says what happens to a
	// character who cannot meet it. Treatment not paid for is treatment not
	// given, so this is the same outcome as declining it.
	if g.char.State.Credits < price {
		g.die(cause, "emergency treatment costs "+itoa(price)+
			" credits and the character has "+itoa(g.char.State.Credits))

		return nil
	}

	g.char.State.Credits -= price
	g.char.State.Characteristics.Set(target, 1)

	g.crisisSurvived = true

	g.consequence(ConsequenceAge, cause,
		"paid "+itoa(price)+" credits for emergency treatment: "+which+" restored to 1",
		"")

	return nil
}

// settleTerminalStates is pp. 123-124, read after a term's aging checks:
// three physical characteristics at 0 is death, two is incapacity, two
// mental is death, and one mental ends enlistment.
func (g *Generator) settleTerminalStates(cause int) {
	physical := g.char.State.Characteristics.zeroed(physicalCharacteristics)
	mental := g.char.State.Characteristics.zeroed(mentalCharacteristics)

	switch {
	case g.char.State.Characteristics.AllPhysicalZero():
		g.die(cause, "all three physical characteristics are at 0")

		return
	case mental >= atZeroEndsIt:
		g.die(cause, "two or more mental characteristics are at 0")

		return
	case physical >= atZeroEndsIt:
		g.char.State.Fate = FateIncapacitated
		g.stopped = true

		g.consequence(ConsequenceAge, cause,
			"two physical characteristics are at 0: incapable of independent movement, "+
				"and character generation ends here", "")

		return
	}

	if mental == 1 && !g.mentalDecline {
		g.mentalDecline = true

		g.consequence(ConsequenceAge, cause,
			"a mental characteristic is at 0: no further enlistment may be attempted", "")
	}
}

// die ends generation. There is no mustering out: the character did not
// leave a career, they stopped.
func (g *Generator) die(cause int, why string) {
	g.char.State.Fate = FateDied
	g.stopped = true

	g.consequence(ConsequenceAge, cause, "died at age "+itoa(g.char.State.Age)+": "+why, "")
}

// crisisCite is the page the Aging Crisis is printed on.
const crisisCite = "p. 123"

// stampApparentAge records what p. 125's chart makes of the character's
// years. It is recorded rather than computed on demand because it needs the
// homeworld's tech level, which the record does not otherwise carry into
// the renderer.
func (g *Generator) stampApparentAge(cause int) {
	band, fromChart := apparentAge(g.techLevel, g.char.State.Age)
	changed := band != g.char.State.ApparentAge

	g.char.State.ApparentAge = band

	// Below the chart, apparent age is the character's age, which the line
	// above this one already said. Only the chart's own answer is worth a
	// consequence of its own.
	if fromChart && changed {
		g.consequence(ConsequenceAge, cause, "apparent age "+band.String(), "")
	}
}

// agingCite is the page the aging tables are printed on. The step above is
// cited to p. 121, where the step begins.
const agingCite = "pp. 122-123"

// characteristicThrow rolls 2d6 plus a characteristic's modifier against a
// target, which is how every check in the book resolves (p. 110).
func (g *Generator) characteristicThrow(check career.Check, extra ...dice.Mod) dice.Throw {
	which, ok := characteristicByName(check.Characteristic)
	if ok {
		extra = append(extra, dice.Mod{
			Name:  check.Characteristic,
			Value: g.char.State.Characteristics.Modifier(which),
		})
	}

	return g.dice.Throw(check.Number, extra...)
}

// takeAutomatic reports whether a named throw has already been decided by a
// table result, and spends it if so.
// enlistmentThrow is the name enlistment modifiers and automatics are
// filed under, matching the career package's own constant.
const (
	enlistmentThrow  = "next enlistment attempt"
	survivalThrow    = "next survival roll"
	advancementThrow = "next advancement roll"
	commissionThrow  = "next commission roll"
	admissionThrow   = "admission to any higher education"
	skillCheckThrow  = "a skill check"
)

// takeAutomaticFor is takeAutomatic for an enlistment, where the result
// that granted it may have named the class of career it reaches: "you may
// enlist automatically in a business, military, corporate or colonist
// career" (Undergraduate University, p. 88).
func (g *Generator) takeAutomaticFor(applies string, target career.Career) bool {
	for i, pending := range g.automatic {
		if pending.Applies != applies || !pending.AppliesTo(target) {
			continue
		}

		g.automatic = append(g.automatic[:i], g.automatic[i+1:]...)

		return true
	}

	return false
}

// takeCareerModifiers and rollSurvival's takeAutomaticFor are the pair: a
// throw made inside a career reads the class of that career, whether what
// it reads is a modifier or an already-decided success.
//
// takeCareerModifiers is takeModifiers for the throws made inside a career
// -- survival and advancement -- where a modifier may name the class of
// career it applies to: "-2 DM to the first two Advancement rolls in a
// military career".
func (g *Generator) takeCareerModifiers(applies string) []dice.Mod {
	return g.takeModifiersFor(applies, *g.career)
}

// takeModifiers consumes every pending modifier that applies to a named
// throw. Consuming rather than reading: "take a -2 on your next
// Advancement roll" is spent once.
func (g *Generator) takeModifiers(applies string) []dice.Mod {
	return g.takeModifiersFor(applies, career.Career{})
}

// takeModifiersFor is takeModifiers for an enlistment throw, where a
// pending modifier may name the class of career it applies to: "-4 DM to
// enlist in any government related career".
//
// A modifier that does not reach this career is kept for the next one, and
// a standing modifier is kept whether it reached this one or not -- which
// is the difference between "your next career" and "every career after
// this one".
func (g *Generator) takeModifiersFor(applies string, target career.Career) []dice.Mod {
	var (
		mods []dice.Mod
		kept []PendingModifier
	)

	for _, pending := range g.pending {
		switch {
		case pending.Applies != applies:
			kept = append(kept, pending)
		case !pending.AppliesTo(target):
			kept = append(kept, pending)
		default:
			mods = append(mods, dice.Mod{Name: pending.Detail, Value: pending.Value})

			// Standing is read on the enlistment throw alone. No result in
			// the book grants a standing modifier to any other throw, and
			// one that reached every future advancement roll would be
			// permanent with nothing to end it.
			if pending.Standing && applies == enlistmentThrow {
				kept = append(kept, pending)

				continue
			}

			// "Your next two Advancement rolls" is one modifier spent
			// twice rather than two modifiers.
			if pending.Uses > 1 {
				pending.Uses--

				kept = append(kept, pending)
			}
		}
	}

	g.pending = kept

	return mods
}
