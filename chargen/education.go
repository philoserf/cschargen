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
}

// The years each outcome costs (pp. 87, 92).
const (
	graduateYears = 4
	washoutDice   = "1d3"
)

// The EDU floors and ceilings success and honours impose (pp. 87, 89).
const (
	graduateEDU   = 10
	honorsEDUCap  = 14
	honorsEDUGain = 2
)

// higherEducation is Step 8.
func (g *Generator) higherEducation() error {
	step := g.log.Step("Step 8: Higher Education", "p. 85")

	// "Unless a previous life period event states otherwise, characters
	// are not required to attend college" (p. 85), so the engine asks.
	institution, attending, err := g.chooseInstitution(step)
	if err != nil || !attending {
		return err
	}

	record := &Education{Institution: institution.Name}

	g.char.State.Education = record

	g.cite = institution.Cite

	if !g.admitted(institution, step) {
		return nil
	}

	record.Admitted = true

	succeeded, err := g.attemptDegree(institution, step, record)
	if err != nil {
		return err
	}

	if !succeeded {
		return g.washOut(institution, step)
	}

	return g.graduateEvents(institution, step)
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
		if !g.meetsPrerequisites(institution, score) {
			continue
		}

		// p. 93: a failed academy closes the academy for good.
		if g.academyClosed && institution.Name == career.MilitaryAcademy().Name {
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

	if g.char.Provenance.Inputs.SkipEducation {
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

	g.raiseToAtLeast("EDU", graduateEDU, step)

	g.char.State.Age += graduateYears
	g.consequence(ConsequenceAge, step,
		"four years at the "+institution.Name+": age "+itoa(g.char.State.Age), "")

	field, err := g.chooseField(institution, step)
	if err != nil {
		return false, err
	}

	record.Field = field

	g.consequence(ConsequenceEducation, step,
		"graduated from the "+institution.Name+" in "+field, "")

	if institution.Note != "" {
		g.unimplemented(step, institution.Note)
	}

	g.attemptHonors(institution, step, record)

	return true, nil
}

// chooseField picks the skill taken at level 2, which is what the degree is
// in.
func (g *Generator) chooseField(institution career.Institution, step int) (string, error) {
	names := make([]string, len(institution.Skills))
	for i, option := range institution.Skills {
		names[i] = option.Skill
	}

	chosen, err := g.choose(Choice{
		Point:   "degree_field",
		Prompt:  "Choose the field of the degree",
		Options: names,
		Cite:    institution.Cite,
	})
	if err != nil {
		return "", err
	}

	taken := institution.Skills[chosen]

	taken.Level = degreeLevel

	err = g.apply(taken, step)
	if err != nil {
		return "", err
	}

	g.unimplemented(step, "any other skill at level 1 (p. 87)")

	return names[chosen], nil
}

// degreeLevel is the level a graduate takes their field at.
const degreeLevel = 2

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
	// character's EDU by 1."
	gain := honorsEDUGain
	if g.char.State.Characteristics.Get(EDU) >= honorsEDUCap {
		gain = 1
	} else if g.char.State.Characteristics.Get(EDU)+gain > honorsEDUCap {
		gain = honorsEDUCap - g.char.State.Characteristics.Get(EDU)
	}

	g.adjust("EDU", gain, "honours at the "+institution.Name, step)
	g.unimplemented(step, "the skills gained at success rise by one level (p. 89)")
}

// educationThrow is a 2d6 against a target, plus one characteristic's
// modifier -- and plus whatever a teenage result attached to it.
func (g *Generator) educationThrow(target career.Throw, cite string) bool {
	mods := g.takeModifiers("admission to any higher education")

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
	if g.institution == nil {
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
	if g.institution == nil || g.char.State.Education == nil {
		return
	}

	g.grantHonors(*g.institution, cause, g.char.State.Education)
}
