package chargen

import (
	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/dice"
)

// Step 8: Higher Education (pp. 85-104).
//
// Four throws deep at most: a prerequisite that is not a throw at all, then
// admission, then success, then honours. A character who fails admission
// goes straight to Step 9; one who fails success rolls on the failure table
// and then goes to Step 9; one who succeeds rolls for honours and then on
// the events table.

// Education is what Step 8 established, and is what the three enlistment
// modifiers that twelve careers carry will read.
type Education struct {
	// Institution is the track attempted, by name.
	Institution string `json:"institution"`

	// Admitted, Succeeded and Honors are the three throws, in order. A
	// character who was not admitted has the other two false.
	Admitted  bool `json:"admitted"`
	Succeeded bool `json:"succeeded"`
	Honors    bool `json:"honors,omitempty"`

	// Field is the skill taken at level 2, which is what the degree is in:
	// "the character has a bachelor's degree or the equivalent in their
	// chosen field which should match the skill they chose at level 2"
	// (p. 89).
	Field string `json:"field,omitempty"`

	// Degree is what a success is worth, and is what three careers'
	// enlistment modifiers read. Empty where the character did not
	// graduate.
	Degree career.Degree `json:"degree,omitempty"`
}

// The years each outcome costs (pp. 87, 92).
const (
	graduateYears = 4
	washoutDice   = "1d3"
)

// The EDU ceiling and gain honours impose (p. 89).
const (
	honorsEDUCap  = 14
	honorsEDUGain = 2
)

// educationAttempts caps the loop below. Graduate School is the only
// institution worth entering twice -- "If the character chooses to return
// to graduate school for a second term and they achieve success, then the
// character has a doctorate" (p. 99) -- so four passes is more than the
// step can use.
const educationAttempts = 4

// higherEducation is Step 8, which runs until the character stops
// qualifying or stops wanting to.
//
// It loops because the two graduate tracks require a degree the character
// may only just have earned: "the character must have achieved Success in
// an Undergraduate University or Military Academy" (pp. 97, 101). A single
// pass would make them unreachable at this step.
func (g *Generator) higherEducation() error {
	step := g.log.Step("Step 8: Higher Education", "p. 85")

	for range educationAttempts {
		again, err := g.oneEducation(step)
		if err != nil || !again {
			return err
		}
	}

	return nil
}

// oneEducation is one attempt, and reports whether another is worth
// offering: only a graduate can go further.
func (g *Generator) oneEducation(step int) (bool, error) {
	// "Unless a previous life period event states otherwise, characters
	// are not required to attend college" (p. 85), so the engine asks.
	institution, attending, err := g.chooseInstitution(step)
	if err != nil || !attending {
		return false, err
	}

	record := &Education{Institution: institution.Name}

	g.char.State.Education = append(g.char.State.Education, record)

	g.cite = institution.Cite

	// A tie made at school names the school, the way one made in a career
	// names the career (FR15).
	g.stage = institution.Name
	defer func() { g.stage = "" }()

	if !g.admitted(institution, step) {
		return false, nil
	}

	record.Admitted = true

	succeeded, err := g.attemptDegree(institution, step, record)
	if err != nil {
		return false, err
	}

	if !succeeded {
		return false, g.washOut(institution, step)
	}

	return true, g.graduateEvents(institution, step)
}

// chooseInstitution asks whether to attend, and where. A character who
// meets no prerequisite is not offered the question.
func (g *Generator) chooseInstitution(step int) (career.Institution, bool, error) {
	score := g.characteristicScore()

	var (
		open  []career.Institution
		names []string
	)

	for _, institution := range career.Institutions() {
		if !g.eligibleFor(institution, score) {
			continue
		}

		open = append(open, institution)
		names = append(names, institution.Name)
	}

	if len(open) == 0 {
		g.consequence(ConsequenceEducation, step,
			"no higher education the character is eligible for (pp. 86, 92)", "")

		return career.Institution{}, false, nil
	}

	// --skip-education says not to attempt it before Step 9. A player who
	// asked to return at Step 18 has said otherwise, and their answer is
	// later and more specific than the flag.
	if g.char.Provenance.Inputs.SkipEducation && !g.returning {
		g.consequence(ConsequenceEducation, step,
			"higher education declined; it is never required (p. 85)", "")

		return career.Institution{}, false, nil
	}

	names = append(names, "none")

	chosen, err := g.choose(Choice{
		Point:   "higher_education",
		Prompt:  "Attempt higher education?",
		Options: names,
		Cite:    "p. 85",
	})
	if err != nil {
		return career.Institution{}, false, err
	}

	if chosen == len(open) {
		g.consequence(ConsequenceEducation, step,
			"higher education declined; it is never required (p. 85)", "")

		return career.Institution{}, false, nil
	}

	return open[chosen], true, nil
}

// eligibleFor collects every reason a character may not attempt an
// institution: the prerequisites it prints, the academy it has already
// failed (p. 93), the degree the graduate tracks require (pp. 97, 101), an
// institution already attended, and the second bachelor's of ERRATA E-29.
func (g *Generator) eligibleFor(institution career.Institution, score func(string) int) bool {
	switch {
	case !g.meetsPrerequisites(institution, score):
		return false
	case g.academyClosed && institution.Name == career.MilitaryAcademy().Name:
		return false
	case len(g.char.State.Terms) < g.educationClosedUntil:
		// A graduate track's worst two failures close every institution
		// rather than the one that expelled the character (pp. 98, 102).
		return false
	case institution.RequiresDegree != g.holdsADegree():
		// A graduate track needs a degree, and an undergraduate one is not
		// gone back to once a degree is held.
		return false
	case g.attendedAlready(institution):
		return false
	}

	return true
}

// closeEducation bars every institution for a number of terms, measured in
// career terms served as the enlistment lockout is. A second such result
// extends the bar rather than restarting it, so the later date wins.
func (g *Generator) closeEducation(effect career.Effect, cause int) {
	until := len(g.char.State.Terms) + effect.Terms
	if until > g.educationClosedUntil {
		g.educationClosedUntil = until
	}

	g.consequence(ConsequenceEducation, cause, effect.Detail, "")
}

// holdsADegree is the graduate tracks' other prerequisite: "the character
// must have achieved Success in an Undergraduate University or Military
// Academy" (pp. 97, 101).
func (g *Generator) holdsADegree() bool {
	for _, held := range g.char.State.Education {
		if held.Degree == career.Bachelors {
			return true
		}
	}

	return false
}

// attendedAlready keeps a character from re-entering an institution in the
// same step. Graduate School is the exception the book makes -- a second
// success there is a doctorate (p. 99) -- so it is not excluded until the
// second time.
func (g *Generator) attendedAlready(institution career.Institution) bool {
	attempts := 0

	for _, held := range g.char.State.Education {
		if held.Institution == institution.Name {
			attempts++
		}
	}

	if institution.Name == career.GraduateSchool().Name {
		return attempts >= graduateSchoolTwice
	}

	return attempts > 0
}

// graduateSchoolTwice is a master's and then a doctorate, which is as far
// as p. 99 goes.
const graduateSchoolTwice = 2

// meetsPrerequisites is the gate before the admission throw: "The
// character's EDU must be 6 or higher to be accepted into Undergraduate
// College" (p. 86).
func (g *Generator) meetsPrerequisites(institution career.Institution, score func(string) int) bool {
	for _, check := range institution.Prerequisites {
		if score(check.Characteristic) < check.Number {
			return false
		}
	}

	return true
}

// admitted makes the admission throw.
func (g *Generator) admitted(institution career.Institution, step int) bool {
	made := g.educationThrow(institution.Admission, institution.Cite)

	if !made {
		g.consequence(ConsequenceEducation, step,
			"not admitted to the "+institution.Name, "")
	}

	_ = step

	return made
}

// attemptDegree makes the success throw and, where it is made, applies what
// p. 87 grants: EDU to ten, four years, a field at level 2 and any other
// skill at level 1, then the honours throw.
func (g *Generator) attemptDegree(
	institution career.Institution, step int, record *Education,
) (bool, error) {
	if !g.educationThrow(institution.Success, institution.Cite) {
		return false, nil
	}

	record.Succeeded = true

	g.raiseEDUForDegree(institution, step)

	g.char.State.Age += graduateYears
	g.consequence(ConsequenceAge, step,
		"four years at the "+institution.Name+": age "+itoa(g.char.State.Age), "")

	field, err := g.chooseField(institution, step)
	if err != nil {
		return false, err
	}

	record.Field = field
	record.Degree = g.degreeFor(institution)

	g.consequence(ConsequenceEducation, step,
		"graduated from the "+institution.Name+" with a "+string(record.Degree)+
			" in "+field, "")

	if institution.Note != "" {
		g.unimplemented(step, institution.Note)
	}

	g.attemptHonors(institution, step, record)

	return true, nil
}

// chooseField picks the skill taken at level 2, which is what the degree is
// in.
// chooseField picks the skill taken at level 2, which is what the degree is
// in, and applies whatever else success grants at that institution.
//
// The three tracks grant different things. An undergraduate one gives "one
// of the following at level 2 ... any other skill at level 1" (pp. 87, 92).
// Graduate school gives "the skill they increased in Undergraduate
// University by two levels ... a level in one of the following" (p. 97).
// Medical school gives no choice at all: "Medic (Any) at level 3, Medic
// (Any) at level 2, Medic (Any) at level 2, and Medic (Any) at level 1.
// Specialties should all be different" (p. 101).
func (g *Generator) chooseField(institution career.Institution, step int) (string, error) {
	switch institution.Name {
	case career.MedSchool().Name:
		return g.medicalDegree(institution, step)
	case career.GraduateSchool().Name:
		return g.graduateDegree(institution, step)
	}

	return g.bachelorsDegree(institution, step)
}

// bachelorsDegree is pp. 87 and 92: one skill from the list at level 2, and
// any other skill at level 1.
func (g *Generator) bachelorsDegree(
	institution career.Institution, step int,
) (string, error) {
	field, err := g.takeFromList(institution, step, degreeLevel, "Choose the field of the degree")
	if err != nil {
		return "", err
	}

	// "The character may also choose any other skill at level 1", which
	// the skill list of pp. 304-314 is what makes answerable.
	err = g.apply(career.Effect{
		Kind:    career.EffectAnySkill,
		NotHeld: true,
		Level:   1,
		Detail:  "any other skill at level 1 (p. 87)",
	}, step)
	if err != nil {
		return "", err
	}

	return field, nil
}

// graduateDegree is p. 97: the undergraduate field rises by two levels, and
// one more skill from the list is taken at level 1.
//
// The raise has to name the same specialty the bachelor's took, so it is
// found in the institution's own list rather than rebuilt from the field's
// name -- "Advocate" and "Advocate (Any)" are two skills, and raising the
// wrong one would leave a character with both.
func (g *Generator) graduateDegree(
	institution career.Institution, step int,
) (string, error) {
	field := g.previousField()

	raise, found := skillNamed(institution.Skills, field)
	if !found {
		// The prerequisite makes this unreachable -- a graduate track
		// cannot be entered without a bachelor's -- but a record that got
		// here anyway says so rather than silently raising nothing.
		g.unimplemented(step, "raise the undergraduate field by two levels (p. 97)")
	} else {
		raise.Level = graduateGain
		raise.Detail = "gain " + itoa(graduateGain) + " levels in " + field

		err := g.apply(raise, step)
		if err != nil {
			return "", err
		}
	}

	taken, err := g.takeFromList(institution, step, 1, "Choose a second field to take a level in")
	if err != nil {
		return "", err
	}

	if !found {
		field = taken
	}

	return field, nil
}

// skillNamed finds a skill in an institution's list by name.
func skillNamed(skills []career.Effect, name string) (career.Effect, bool) {
	for _, option := range skills {
		if option.Skill == name {
			return option, true
		}
	}

	return career.Effect{}, false
}

// graduateGain is the two levels a master's adds to the field a bachelor's
// established (p. 97).
const graduateGain = 2

// previousField is the skill the character's bachelor's was in, which
// graduate school raises rather than replacing.
func (g *Generator) previousField() string {
	for _, held := range g.char.State.Education {
		if held.Degree == career.Bachelors {
			return held.Field
		}
	}

	return ""
}

// medicalDegree is p. 101: "Gain Medic (Any) at level 3, Medic (Any) at
// level 2, Medic (Any) at level 2, and Medic (Any) at level 1. Specialties
// should all be different."
//
// The four are chosen one at a time from what is left, which is what makes
// them different. Offering the whole list four times and letting the
// Decider repeat itself would give a character Medic 8 in one specialty,
// which is the opposite of what the page asks for.
func (g *Generator) medicalDegree(
	institution career.Institution, step int,
) (string, error) {
	left := career.MedicSpecialties()

	for _, level := range medicalLevels {
		if len(left) == 0 {
			break
		}

		chosen, err := g.choose(Choice{
			Point:   "medic_specialty",
			Prompt:  "Choose a Medic specialty to take at level " + itoa(level),
			Options: left,
			Cite:    institution.Cite,
		})
		if err != nil {
			return "", err
		}

		granted := career.RaiseSkill("Medic", level)

		granted.Specialties = []string{left[chosen]}

		err = g.apply(granted, step)
		if err != nil {
			return "", err
		}

		left = append(left[:chosen], left[chosen+1:]...)
	}

	return "Medic", nil
}

// medicalLevels is what p. 101 grants, in the order it prints them.
var medicalLevels = [...]int{3, 2, 2, 1}

// takeFromList offers an institution's skill list and grants the chosen one
// at a level.
func (g *Generator) takeFromList(
	institution career.Institution, step, level int, prompt string,
) (string, error) {
	names := make([]string, len(institution.Skills))
	for i, option := range institution.Skills {
		names[i] = option.Skill
	}

	chosen, err := g.choose(Choice{
		Point:   "degree_field",
		Prompt:  prompt,
		Options: names,
		Cite:    institution.Cite,
	})
	if err != nil {
		return "", err
	}

	taken := institution.Skills[chosen]

	taken.Level = level

	return names[chosen], g.apply(taken, step)
}

// degreeLevel is the level a graduate takes their field at.
const degreeLevel = 2

// degreeFor is what a success at an institution is worth. A second success
// at graduate school is a doctorate rather than a second master's (p. 99).
func (g *Generator) degreeFor(institution career.Institution) career.Degree {
	switch institution.Name {
	case career.GraduateSchool().Name:
		for _, held := range g.char.State.Education {
			if held.Degree == career.Masters {
				return career.Doctorate
			}
		}

		return career.Masters
	case career.MedSchool().Name:
		return career.MedicalDoctor
	}

	return career.Bachelors
}

// raiseEDUForDegree is the EDU floor a success imposes, which is ten for
// the undergraduate tracks and twelve for the graduate ones (pp. 87, 97).
//
// Medical school's is a ladder rather than a floor: "Increase the
// character's EDU to 12. If the character's EDU is 12-13, make it 14. If it
// is already 14 or higher, add 1 to a maximum of 16" (p. 101).
func (g *Generator) raiseEDUForDegree(institution career.Institution, step int) {
	if institution.Name != career.MedSchool().Name {
		g.raiseToAtLeast("EDU", institution.GraduateEDU, step)

		return
	}

	held := g.char.State.Characteristics.Get(EDU)

	switch {
	case held < medicalEDUFirst:
		g.adjust("EDU", medicalEDUFirst-held, "a medical degree raises EDU", step)
	case held < medicalEDUSecond:
		g.adjust("EDU", medicalEDUSecond-held, "a medical degree raises EDU", step)
	case held < medicalEDUCap:
		g.adjust("EDU", 1, "a medical degree raises EDU", step)
	}
}

// The three rungs of medical school's EDU ladder (p. 101).
const (
	medicalEDUFirst  = 12
	medicalEDUSecond = 14
	medicalEDUCap    = 16
)

// attemptHonors is the third throw.
func (g *Generator) attemptHonors(institution career.Institution, step int, record *Education) {
	if !g.educationThrow(institution.Honors, institution.Cite) {
		return
	}

	g.grantHonors(institution, step, record)
}

// grantHonors applies what honours are worth, which one result on each
// events table can also do after the throw has failed (pp. 89, 91).
func (g *Generator) grantHonors(institution career.Institution, step int, record *Education) {
	if record.Honors {
		return
	}

	record.Honors = true

	// "the character's EDU should be increased by 2 to a maximum of 14. If
	// the character's EDU is already 14 or higher, increase the
	// character's EDU by 1" (p. 89) -- and medical school's version adds
	// "to a maximum of 16" to that second sentence (p. 102).
	gain := honorsGain(g.char.State.Characteristics.Get(EDU), institution)
	if gain <= 0 {
		return
	}

	g.adjust("EDU", gain, "honours at the "+institution.Name, step)
	g.unimplemented(step, "the skills gained at success rise by one level (p. 89)")
}

// honorsGain is how much EDU honours are worth at this institution: "+2 to
// a maximum of 14. If the character's EDU is already 14 or higher, increase
// the character's EDU by 1" (p. 89), with medical school's version adding
// "to a maximum of 16" to that second sentence (p. 102).
func honorsGain(held int, institution career.Institution) int {
	if held < honorsEDUCap {
		return min(honorsEDUGain, honorsEDUCap-held)
	}

	if institution.Name == career.MedSchool().Name {
		return max(min(1, medicalEDUCap-held), 0)
	}

	return 1
}

// educationThrow is a 2d6 against a target, plus one characteristic's
// modifier -- and plus whatever a teenage result attached to it.
func (g *Generator) educationThrow(target career.Throw, cite string) bool {
	mods := g.takeModifiers(admissionThrow)

	which, ok := characteristicByName(target.Characteristic)
	if ok {
		mods = append(mods, dice.Mod{
			Name:  target.Characteristic,
			Value: g.char.State.Characteristics.Modifier(which),
		})
	}

	throw := g.dice.Throw(target.Number, mods...)
	g.log.Throw(throw, cite)

	return throw.Success
}

// raiseToAtLeast is p. 87's "increase their EDU to 10 ... If their EDU was
// already higher than 10, the character should increase their EDU by 1".
func (g *Generator) raiseToAtLeast(which string, floor, cause int) {
	target, ok := characteristicByName(which)
	if !ok {
		return
	}

	held := g.char.State.Characteristics.Get(target)

	gain := max(floor-held, 1)

	g.adjust(which, gain, "a degree raises "+which, cause)
}

// washOut is the failure table, and the years it costs.
func (g *Generator) washOut(institution career.Institution, step int) error {
	roll := g.dice.TwoD6()
	cause := g.log.Roll(roll, institution.Cite)
	row := institution.Failure[roll.Total-2]

	g.consequence(ConsequenceEducation, cause,
		"left the "+institution.Name+": "+row.Summary, "")

	years, err := g.rollExpression(washoutDice, institution.Cite)
	if err != nil {
		return err
	}

	g.char.State.Age += years
	g.consequence(ConsequenceAge, cause, "age "+itoa(g.char.State.Age), "")

	// "If a character fails their success roll while in a military
	// academy, they may not attempt to enter a military academy again"
	// (p. 93).
	if institution.Name == career.MilitaryAcademy().Name {
		g.academyClosed = true
	}

	_, err = g.washoutSkill(institution, step)

	return err
}

// washoutSkill is the consolation p. 87 gives: a level in one of the same
// skills, rather than two of them and a degree.
func (g *Generator) washoutSkill(institution career.Institution, step int) (string, error) {
	names := make([]string, len(institution.Skills))
	for i, option := range institution.Skills {
		names[i] = option.Skill
	}

	chosen, err := g.choose(Choice{
		Point:   "washout_skill",
		Prompt:  "Choose the skill the unfinished degree left behind",
		Options: names,
		Cite:    institution.Cite,
	})
	if err != nil {
		return "", err
	}

	return names[chosen], g.apply(institution.Skills[chosen], step)
}

// graduateEvents is the 2d10 table a successful graduate rolls on.
func (g *Generator) graduateEvents(institution career.Institution, step int) error {
	roll := g.dice.ND10(youthDice)
	cause := g.log.Roll(roll, institution.Cite)
	row := institution.Events[roll.Total-youthLow]

	g.consequence(ConsequenceEducation, cause,
		"at the "+institution.Name+": "+row.Summary, "")

	g.institution = &institution

	_ = step

	return g.applyAll(row.Effects, cause)
}

// rollInstitutionLifeEvent is the d6 table an institution's events table
// reaches at result 10.
func (g *Generator) rollInstitutionLifeEvent(cause int) error {
	if g.institution == nil || len(g.institution.LifeEvents) == 0 {
		return nil
	}

	roll := g.dice.D6()
	throw := g.log.Roll(roll, g.institution.Cite)
	row := g.institution.LifeEvents[roll.Total-1]

	g.consequence(ConsequenceEducation, throw,
		"life at the "+g.institution.Name+": "+row.Summary, "")

	_ = cause

	return g.applyAll(row.Effects, throw)
}

// honorsByEvent is the events-table result that grants honours a failed
// throw did not: "If you failed on your Honors roll, you have now been
// successful" (p. 91).
func (g *Generator) honorsByEvent(cause int) {
	if g.institution == nil || len(g.char.State.Education) == 0 {
		return
	}

	g.grantHonors(*g.institution, cause, g.char.State.Education[len(g.char.State.Education)-1])
}
