package career

// EnlistmentModKind names a situational modifier on an enlistment throw.
type EnlistmentModKind int

// The kinds. Each is evaluated against something the engine tracks: the
// number of careers entered, the character's apparent age, and the degrees
// Step 8 recorded.
const (
	// PerPreviousCareer is "-1 for each career you have entered before this
	// one", multiplied by the count.
	PerPreviousCareer EnlistmentModKind = iota

	// ApparentAgeOver40 is "-2 if the character's apparent age is 40+".
	ApparentAgeOver40

	// UndergraduateDegree and GraduateDegree are the education bonuses the
	// Instructor career prints (p. 212).
	UndergraduateDegree
	GraduateDegree

	// MedicalSchool is the Medic career's own, and the largest in the book
	// (p. 229).
	MedicalSchool
)

// EnlistmentMod is one such modifier.
type EnlistmentMod struct {
	Kind  EnlistmentModKind
	Value int
}

// perPreviousCareer and the rest build the modifiers a career prints. The
// two military ones take no value because every career that prints them
// prints the same number: -1 for each previous career, -2 for apparent age
// over forty. A career that printed a different figure would take a value,
// and none does.
func perPreviousCareer() EnlistmentMod {
	return EnlistmentMod{Kind: PerPreviousCareer, Value: -1}
}

func apparentAgeOver40() EnlistmentMod {
	return EnlistmentMod{Kind: ApparentAgeOver40, Value: -2}
}

func undergraduateDegree(value int) EnlistmentMod {
	return EnlistmentMod{Kind: UndergraduateDegree, Value: value}
}

func graduateDegree(value int) EnlistmentMod {
	return EnlistmentMod{Kind: GraduateDegree, Value: value}
}

// medicalSchool is the largest modifier in the book: the Medic career
// enlists on EDU 11+, "dropping to EDU 5+" for a character who attended
// medical school (p. 229).
func medicalSchool() EnlistmentMod {
	return EnlistmentMod{Kind: MedicalSchool, Value: 6}
}

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

	// Ranks is indexed from 0. A rank with no printed benefit holds a nil
	// entry, which the engine reads as nothing granted rather than as a
	// missing row.
	//
	// The length varies by career and is not seven. Colonist prints ranks 0
	// to 6 (p. 174); Marine prints nine enlisted ranks, E0 to E8, and eight
	// officer ranks, O0 to O7 (p. 225). A fixed array was the first thing
	// the first transcription broke.
	Ranks [][]Effect

	// OfficerRanks is the second rank table a commissioned career prints
	// per assignment (pp. 225, 234). Nil where the career has no
	// commission.
	OfficerRanks [][]Effect
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

	// EnlistmentMods are the situational modifiers a career prints on its
	// enlistment throw: "Take a -1 modifier for each career you have
	// entered before this one", "If the character's apparent age is 40+,
	// take a -2 modifier" (p. 111).
	//
	// They are carried as a list of kinds rather than as prose because the
	// engine can evaluate some of them and not others, and a modifier it
	// cannot evaluate has to be recorded rather than quietly dropped.
	EnlistmentMods []EnlistmentMod

	// Prerequisite is a condition the book puts on entering the career that
	// the engine cannot check. It is recorded on the character rather than
	// enforced, because refusing a career on a rule this engine cannot
	// evaluate would be worse than admitting it cannot.
	Prerequisite string

	// Commission is the throw to be commissioned, where the career has one
	// (p. 114). Nil otherwise.
	Commission *Check

	Assignments []Assignment

	// Tables are the career-wide skill tables: Personal Development,
	// Service Skills, and Advanced Education where the career has one.
	// Assignment tables live on their assignments.
	Tables []SkillTable

	// Benefits is the Mustering Out Benefits table, rows 1 through 7. This
	// one really is seven everywhere: the roll is 1d6 and the seventh row
	// is reached only through an event's modifier (ERRATA E-3).
	Benefits [7]BenefitRow

	// RankTitles and OfficerTitles name the ranks where the career prints
	// them (p. 115). Descriptive rather than mechanical, and absent from
	// the careers that use no formal titles.
	RankTitles    []string
	OfficerTitles []string

	// Tags are the classes this career belongs to, which twelve results
	// modify or direct an enlistment by. The book defines none of the
	// words, so the assignment is a reading -- see [Tag] and ERRATA E-37.
	Tags []Tag

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

// All returns every career the book names, in the order it lists them
// (pp. 8-9). All thirty-four are transcribed; TestTheBookIsFullyTranscribed
// is what says so.
func All() []Career {
	return []Career{
		Adventurer(),
		Arts(),
		Belter(),
		Celebrity(),
		Clergy(),
		Colonist(),
		CorporateShipper(),
		Craftsperson(),
		DiplomaticService(),
		Exotic(),
		Explorer(),
		FringeMarketer(),
		Gambler(),
		IndependentMerchant(),
		Instructor(),
		Investigator(),
		Journalist(),
		Marine(),
		Medic(),
		NationalNavy(),
		OrbitalConstruction(),
		OrganizedCrime(),
		Pirate(),
		Politician(),
		Prisoner(),
		Scavenger(),
		Scientist(),
		Slave(),
		Sports(),
		SystemDefenceNavy(),
		SystemDefenceTroopers(),
		SystemDefenceWetNavy(),
		Thief(),
		Vagabond(),
	}
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
