package career

// Check is a throw of 2d6 plus a characteristic modifier against a target
// (p. 110).
type Check struct {
	Characteristic string
	Number         int
}

// SkillTableKind names the tables of p. 117. Every career has Personal
// Development and Service Skills; most have Advanced Education, gated on
// EDU; careers with commissions add Officer Skills; and each assignment has
// one of its own.
type SkillTableKind int

// The kinds.
const (
	PersonalDevelopment SkillTableKind = iota
	ServiceSkills
	AdvancedEducation
	OfficerSkills
	AssignmentSkills
)

// String names the table as the book heads it.
func (k SkillTableKind) String() string {
	switch k {
	case PersonalDevelopment:
		return "Personal Development"
	case ServiceSkills:
		return "Service Skills"
	case AdvancedEducation:
		return "Advanced Education"
	case OfficerSkills:
		return "Officer Skills"
	case AssignmentSkills:
		return "Assignment"
	}

	return "unknown"
}

// SkillTable is one of the six-row tables a character rolls 1d6 on (p. 117).
// Rows is indexed 0-5 for results 1-6.
type SkillTable struct {
	Kind SkillTableKind
	Name string

	// MinimumEDU gates the Advanced Education table: "restricted to
	// characters with an EDU of a certain level or higher (most often 8 or
	// higher)" (p. 117). Zero means ungated.
	MinimumEDU int

	Rows [6]Effect
}

// Assignment is a track within a career (p. 112): its own survival and
// advancement targets, its own skill table, and its own rank benefits.
type Assignment struct {
	Name        string
	Description string
	Survival    Check
	Advancement Check
	Skills      SkillTable

	// Ranks is indexed 0-6 for ranks 0 through 6. A rank with no printed
	// benefit holds a zero Effect, which the engine reads as nothing
	// granted rather than as a missing row.
	Ranks [7][]Effect
}

// BenefitRow is one row of a career's Mustering Out Benefits table
// (p. 127). Rows run 1 through 7 although the roll is 1d6; see ERRATA E-3.
type BenefitRow struct {
	Cash  int
	Other Effect
}

// MishapTable is the 2d6 table of p. 113, indexed 0-10 for results 2-12.
type MishapTable [11]MishapRow

// MishapRow is one mishap.
type MishapRow struct {
	Summary string
	Effects []Effect
}

// EventTable is the d66 table of p. 118. Keys are the printed results --
// 11-16, 21-26, ... 61-66 -- so a reader can find a row by the number they
// rolled.
type EventTable map[int]EventRow

// EventRow is one event.
type EventRow struct {
	Summary string
	Effects []Effect
}

// D66Results is every result a d66 can produce, in printed order. A table
// is complete when it has a row for each.
func D66Results() []int {
	results := make([]int, 0, 36)

	for tens := 1; tens <= 6; tens++ {
		for units := 1; units <= 6; units++ {
			results = append(results, tens*10+units)
		}
	}

	return results
}

// Career is one career definition.
type Career struct {
	Name string
	Cite string

	// Enlistment is the throw to enter the career (p. 110). Nil for the two
	// careers entered without one: Prisoner, which cannot be chosen at all,
	// and Vagabond, which is entered by choice or by circumstance (p. 111).
	Enlistment *Check

	// EnlistmentNote explains a nil Enlistment, so that a reader of the
	// data is not left to infer why the throw is missing.
	EnlistmentNote string

	// Commission is the throw to be commissioned, where the career has one
	// (p. 114). Nil otherwise.
	Commission *Check

	Assignments []Assignment

	// Tables are the career-wide skill tables: Personal Development,
	// Service Skills, and Advanced Education where the career has one.
	// Assignment tables live on their assignments.
	Tables []SkillTable

	// Benefits is the Mustering Out Benefits table, rows 1 through 7.
	Benefits [7]BenefitRow

	Mishaps MishapTable
	Events  EventTable

	// MishapEjects is false for the two careers that say so in print:
	// Vagabond (p. 298) and Prisoner (p. 263) both note that, unlike other
	// careers, a mishap does not force the character out. They are also the
	// two careers a character is most often forced into, so the exception
	// is not a corner.
	MishapEjects bool
}

// Assignment finds an assignment by name.
func (c Career) Assignment(name string) (Assignment, bool) {
	for _, a := range c.Assignments {
		if a.Name == name {
			return a, true
		}
	}

	return Assignment{}, false
}

// Table finds a career-wide skill table by kind.
func (c Career) Table(kind SkillTableKind) (SkillTable, bool) {
	for _, t := range c.Tables {
		if t.Kind == kind {
			return t, true
		}
	}

	return SkillTable{}, false
}

// All returns every career this milestone implements, in the order the book
// lists them (p. 107).
func All() []Career {
	return []Career{Colonist(), Prisoner(), Vagabond()}
}

// ByName finds an implemented career.
func ByName(name string) (Career, bool) {
	for _, c := range All() {
		if c.Name == name {
			return c, true
		}
	}

	return Career{}, false
}
