package career

// Step 8: Higher Education (pp. 85-104).
//
// Four institutions, each with the same shape: a prerequisite, an admission
// throw, a success throw, an honours throw, a failure table for those who do
// not succeed and an events table for those who do.
//
// It is not strictly a pre-career step. "College or university may be
// entered at any point after the character reaches the age of 18 ... Some
// characters may attend immediately after their teenage years, while others
// may return later in life after military service" (p. 85), and Step 18
// lists "return to higher education" among what a character decides between
// terms (p. 125).

// Institution is one of the four tracks.
type Institution struct {
	Name string
	Cite string

	// Prerequisites are what a character must already be to apply at all,
	// as against what they must throw. Undergraduate wants EDU 6+; the
	// Military Academy wants EDU 6+ and END 8+.
	Prerequisites []Check

	// RequiresDegree is the prerequisite that is not a characteristic:
	// "the character must have achieved Success in an Undergraduate
	// University or Military Academy" (pp. 97, 101).
	RequiresDegree bool

	// GraduateEDU is the floor success raises EDU to (pp. 87, 97): ten for
	// the two undergraduate tracks, twelve for the two graduate ones.
	GraduateEDU int

	// Admission, Success and Honors are 2d6 throws, each modified by one
	// characteristic's modifier.
	Admission Throw
	Success   Throw
	Honors    Throw

	// Skills is the list a successful graduate takes one of at level 2,
	// and that a failed one takes one of at level 1.
	Skills []Effect

	// Failure is the 2d6 table for a character who does not succeed,
	// indexed 0-10. Events is the 2d10 table for one who does, indexed
	// 0-18. LifeEvents is the d6 table result 10 of Events reaches.
	Failure    MishapTable
	Events     []EventRow
	LifeEvents []EventRow

	// Note is what the institution does beyond the tables: the Military
	// Academy commits a graduate to a military career as an officer.
	Note string
}

// Throw is a 2d6 against a target, modified by one characteristic's
// modifier. The institutions' three throws differ only in the number and
// which characteristic helps.
type Throw struct {
	Number         int
	Characteristic string
}

// Institutions is the four, in the order the book prints them. Graduate
// School and Medical School require a degree, which is why they come last
// in more than presentation.
func Institutions() []Institution {
	return []Institution{
		Undergraduate(), MilitaryAcademy(), GraduateSchool(), MedSchool(),
	}
}

// Degree names what a character holds after succeeding at an institution,
// which is what three careers' enlistment modifiers read.
type Degree string

const (
	// Bachelors is what Undergraduate College and the Military Academy
	// grant (pp. 89, 94).
	Bachelors Degree = "bachelor's"

	// Masters is a first success at graduate school (p. 99).
	Masters Degree = "master's"

	// Doctorate is a second: "If the character chooses to return to
	// graduate school for a second term and they achieve success, then the
	// character has a doctorate" (p. 99).
	Doctorate Degree = "doctorate"

	// MedicalDoctor is what medical school grants (p. 103).
	MedicalDoctor Degree = "medical doctor"
)

// The three throws every institution makes, named because the numbers
// differ between them and the shape does not.
func admission(number int) Throw {
	return Throw{Number: number, Characteristic: "EDU"}
}

func successThrow(number int) Throw {
	return Throw{Number: number, Characteristic: "INT"}
}

// The two graduate tracks read the other characteristic on each throw:
// admission and honours on INT, success on EDU (pp. 97, 101).
func graduateAdmission(number int) Throw {
	return Throw{Number: number, Characteristic: "INT"}
}

func graduateSuccess(number int) Throw {
	return Throw{Number: number, Characteristic: "EDU"}
}

func graduateHonors(number int) Throw {
	return Throw{Number: number, Characteristic: "INT"}
}

func honorsThrow(number int) Throw {
	return Throw{Number: number, Characteristic: "EDU"}
}

// collegiateLifeEvents is the d6 table of p. 91, which the undergraduate
// events table reaches at result 10.
func collegiateLifeEvents() []EventRow {
	return []EventRow{
		{Summary: "an injury", Effects: []Effect{injury(1)}},
		{
			Summary: "a death among family or friends",
			Effects: []Effect{loseTie(Contact, Ally)},
		},
		{
			Summary: "a relationship that ends badly",
			Effects: []Effect{loseTie(Contact, Ally)},
		},
		{Summary: "a new friend", Effects: []Effect{relationship(Contact, 1, "")}},
		{Summary: "a relationship that improves", Effects: []Effect{improveARelationship()}},
		{
			Summary: "something wonderful",
			Effects: []Effect{pick("what came of it",
				opt("CHA", chr("CHA", 2)),
				opt("a new Ally", relationship(Ally, 1, "")),
				opt("a boon", credits("100000")),
				opt("a relationship improves", improveARelationship()))},
		},
	}
}

// improveARelationship is the upgrade three life-event tables print in the
// same words: every band moves up one, and a character with nothing gains a
// Contact.
func improveARelationship() Effect {
	return unimplemented(
		"one relationship improves by a band -- an Enemy becomes a Rival, a Rival a " +
			"Contact, a Contact an Ally -- and with none, gain a Contact")
}

// RaiseSkill is a skill raised by a named number of levels, which Step 8
// grants where a career table grants one: "increase the skill they
// increased in Undergraduate University by two levels" (p. 97).
func RaiseSkill(name string, levels int) Effect {
	granted := skill(name)

	granted.Detail = "gain " + itoa(levels) + " levels in " + name
	granted.Level = levels

	return granted
}

// MedicSpecialties is the list a medical degree's four levels are spread
// across, from the skill list on p. 310.
//
// The book prints eight. One of them is the specialty for treating
// genetically engineered humans, and its name is the term the OGL notice
// reserves, so it is here under the name this repository uses for that
// people -- see CLAUDE.md.
func MedicSpecialties() []string {
	return []string{
		"Alteration",
		"Cryogenics",
		"Cybernetics",
		"Diagnosis",
		"Engineered Humans",
		"First Aid",
		"Surgery",
		"Uplifts",
	}
}
