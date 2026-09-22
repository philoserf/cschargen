// Package setting holds the world, subsector and species data character
// generation is set in, read from a file the user supplies.
//
// None of it is in this repository. The Clement Sector Core Character
// Creation Book declares as Product Identity "all subsector names, world
// maps, world names, system names ... and organizations" (OGL section 16,
// p. 335), and those names are exactly what the origin charts of pp. 43-56
// are made of. The mechanics that read this data are open content and live
// in the engine; the names are the user's own transcription of a book they
// own.
//
// The same clause names one further term, the book's word for a
// genetically engineered human, and says it is not open content. This
// package therefore says "engineered human" throughout -- in its prose, its
// field names and its JSON keys -- rather than the book's term. A data file
// is something a user types, and a key is a worse place to put a word we
// are not entitled to use than a comment is.
//
// data/setting.sample.json is invented, so the tests and the demo run
// without it. A character generated against the sample is stamped as such
// and says so on its sheet.
package setting

import "slices"

// SchemaVersion is the shape of file this package reads. A file declaring a
// version this build does not know is refused rather than read hopefully:
// a field that moved would otherwise be read as absent.
const SchemaVersion = 1

// Status is whether engineered people on a world are free or owned. Most
// worlds do not enslave them; a small number continue the practice of
// slavery or forced servitude, and on those worlds they may be treated as
// property rather than people (p. 42).
type Status string

// The two statuses.
const (
	Free     Status = "free"
	Enslaved Status = "enslaved"
)

// The two kinds of species a data file may declare.
const (
	KindEngineered = "engineered"
	KindUplift     = "uplift"
)

// Permission is one world's treatment of one kind of engineered person.
//
// It is not a boolean, because the page is not: an entry may bar one named
// species while admitting the rest, and say whether those admitted are free
// or owned. So a world carries a permission, a status, and a list of
// species barred by name.
type Permission struct {
	Allowed bool     `json:"allowed"`
	Status  Status   `json:"status,omitempty"`
	Banned  []string `json:"banned,omitempty"`
}

// Admits reports whether a named species may be born here.
func (p Permission) Admits(species string) bool {
	if !p.Allowed {
		return false
	}

	return !slices.Contains(p.Banned, species)
}

// Alternative is one way of satisfying a background-skill requirement: a
// skill with the specialties it offers, or a named item that is not a skill
// at all.
//
// One world grants "Survival (High Pressure, Mountains, Plains, or Ocean)
// and Electronics (Computers) and a mindcomp" (p. 41), and the mindcomp is
// the reason Item exists.
type Alternative struct {
	Skill       string   `json:"skill,omitempty"`
	Specialties []string `json:"specialties,omitempty"`
	Item        string   `json:"item,omitempty"`
}

// Requirement is one background skill the homeworld grants. OneOf holds the
// alternatives the page joins with "or"; a requirement with a single
// alternative is granted outright.
//
// A world's background skills are a list of these, which is the page's
// "and". Encoding the expression rather than the prose keeps the operators
// out of a parser: "Suit (Vacc Suit) or Survival (Freefall) and Language
// (Maori, English, or Hindi)" is two requirements, the first with two
// alternatives and the second with one that offers three specialties.
type Requirement struct {
	OneOf []Alternative `json:"oneOf"`
}

// World is one row of an origin chart (pp. 43-56).
type World struct {
	Name string `json:"name"`

	// Roll is the inclusive d100 range this world occupies on its chart. A
	// world with no range is choose-only, which is what every world on the
	// Recently Colonized Worlds table is: "you cannot randomly be assigned
	// one of these worlds, you may choose them" (p. 40).
	//
	// A slice rather than a fixed-size array so that its length is a thing the
	// validator can check. encoding/json fills a fixed-size array
	// positionally and discards what will not fit, so "roll": [1, 4, 9]
	// decoded silently as 1-4 -- a transcription typo read as a narrower
	// range, in the one input this program treats as untrusted.
	Roll []int `json:"roll,omitempty"`

	TechLevel    int `json:"techLevel"`
	MaximumAge   int `json:"maximumAge"`
	MaximumTerms int `json:"maximumTerms"`

	// SettledYear is not a printed column. Two of the teenage-event paths
	// gate on whether the homeworld "has been colonized or established for
	// 100+ standard years" (p. 76), so the file carries it explicitly.
	SettledYear int `json:"settledYear,omitempty"`

	PrimaryLanguages []string      `json:"primaryLanguages"`
	BackgroundSkills []Requirement `json:"backgroundSkills,omitempty"`

	// BirthSituation modifies the Human Birth Situation throw (p. 58). The
	// book gives that modifier as a list of world names -- "Add 30 to the
	// roll if born on ..." -- and every one of those names is Product
	// Identity, so the number lives with the world rather than in the
	// engine. Zero is no modifier.
	BirthSituation int `json:"birthSituation,omitempty"`

	// BirthSituationOnly forces the throw's result where a world does not
	// modify it but replaces it: "If born on Kyiv, Erlik, Galawdewos,
	// Sarawak, or Superior, then heterosexual couple will be the only
	// result possible" (p. 58), and three worlds where it is always a
	// communal group. Empty means the throw decides.
	BirthSituationOnly string `json:"birthSituationOnly,omitempty"`

	Engineered Permission `json:"engineered"`
	Uplifts    Permission `json:"uplifts"`
}

// Selectable reports whether this world can be reached by a d100 throw, as
// against only by choosing it.
func (w World) Selectable() bool {
	return len(w.Roll) == rollBounds
}

// Covers reports whether a d100 result lands on this world.
func (w World) Covers(roll int) bool {
	if !w.Selectable() {
		return false
	}

	return roll >= w.Roll[0] && roll <= w.Roll[1]
}

// Subsector is one origin chart.
type Subsector struct {
	Name string `json:"name"`

	// OriginRoll is this subsector's result on the 1d6 chart of p. 39. A
	// subsector with none is choose-only.
	OriginRoll int `json:"originRoll,omitempty"`

	// Offworld marks a subsector outside the sector the campaign is set in.
	// The book's rule for those characters -- "must spend a minimum of four
	// terms" in the campaign's own sector (p. 43) -- needs the engine to
	// know where a term was served, which it does not yet.
	Offworld bool `json:"offworld,omitempty"`

	Worlds []World `json:"worlds"`
}

// Species is one engineered kind, named so that a world's Banned list can
// refer to it. This is the vocabulary the permission columns are written
// against; the rules for playing one are in chargen.
type Species struct {
	Name string `json:"name"`

	// Kind is "engineered" for an engineered human or "uplift" for an
	// uplifted animal.
	Kind string `json:"kind"`

	// Characteristics is the method Step 2 uses, one dice expression per
	// characteristic: "All generations should roll 2d6-2 for STR and END.
	// DEX should be rolled as 2d6+2 ... INT and EDU should be rolled
	// normally" (p. 23, of a species this repository cannot name).
	//
	// A characteristic the map does not mention is rolled the human way,
	// and a character's freely assigned scores are only the ones rolled
	// that way -- a method names which characteristic it is for.
	Characteristics map[string]string `json:"characteristics,omitempty"`

	// Skills is what the species starts with at level 1, before any
	// background or career: one species is owed "Survival (Freefall) and
	// Survival (Low Gravity) at level 1" (p. 23), and others their own.
	Skills []Alternative `json:"skills,omitempty"`

	// Aging names one of the engine's aging profiles. Empty is the
	// tech-level profile, which humans and four of the book's five
	// engineered species use.
	Aging string `json:"aging,omitempty"`

	// Maximum is the ceiling on this species' characteristics. p. 14 sets
	// fifteen "for an unaltered human" and says "Some uplifts and
	// engineered humans will have higher maximums". Zero is the human's.
	Maximum int `json:"maximum,omitempty"`

	// YouthRolls and TeenRolls are how many times this species rolls on
	// the youth and teenage tables, and YouthAges and TeenAges are what
	// each roll represents. "Uplifts will often have shorter youths than
	// humans" (p. 68): one species rolls once for ages 2-4 where a human
	// rolls twice for 4-8 and 9-12.
	//
	// Zero rolls means the human pattern.
	YouthRolls int      `json:"youthRolls,omitempty"`
	YouthAges  []string `json:"youthAges,omitempty"`
	TeenRolls  int      `json:"teenRolls,omitempty"`
	TeenAges   []string `json:"teenAges,omitempty"`

	// Genetics is how this species is born (pp. 62-65): a purebred of some
	// generation, a hybrid with a baseline human, or a compound of two
	// engineered species. Absent for a species the book gives no such
	// table -- uplifts have classes instead.
	Genetics *Genetics `json:"genetics,omitempty"`
}

// Genetics is a species' three tables of pp. 62-65.
//
// Every row of the compound table names another species, and every one of
// those names is Product Identity -- which is why the whole thing is data
// rather than engine.
type Genetics struct {
	// Status is the 1d6 genetic status table, six rows. Empty takes the
	// shared table of p. 63, which four of the book's five engineered
	// species use; the fifth prints its own, with no hybrids or compounds
	// in it (p. 62).
	Status []GeneticStatus `json:"status,omitempty"`

	// Hybrid and Compound are the 1d6 tables the status table's last three
	// rows send a character to, six rows each.
	Hybrid   []GeneticOutcome `json:"hybrid,omitempty"`
	Compound []GeneticOutcome `json:"compound,omitempty"`
}

// GeneticStatus is one row of the status table: what the character is, and
// which of the two other tables it sends them to.
type GeneticStatus struct {
	// Kind is "purebred", "hybrid" or "compound".
	Kind string `json:"kind"`

	// Generation is which one, for a purebred. Zero on the row that lets
	// the player choose "third, fourth or later" (p. 63).
	Generation int `json:"generation,omitempty"`

	Detail string `json:"detail,omitempty"`
}

// GeneticOutcome is one row of a hybrid or compound table.
type GeneticOutcome struct {
	Detail string `json:"detail"`

	// RollAs is "human" where the row says "Roll characteristics as a
	// human". Empty rolls them as the species, which is what every other
	// row means.
	RollAs string `json:"rollAs,omitempty"`

	// Adjust is a one-off change a row prints: "take a -1 to their DEX".
	Adjust map[string]int `json:"adjust,omitempty"`

	// Choose is the last row of every compound table: "The player chooses
	// from the above results."
	Choose bool `json:"choose,omitempty"`
}

// The three kinds of genetic status (pp. 62-65).
const (
	Purebred = "purebred"
	Hybrid   = "hybrid"
	Compound = "compound"
)

// DefaultGeneticStatus is the 1d6 table of p. 63, which four of the book's
// five engineered species share: "1 first generation purebred, 2 second
// generation, 3 third or later, 4-5 hybrid, 6 compound".
func DefaultGeneticStatus() []GeneticStatus {
	const secondGeneration = 2

	return []GeneticStatus{
		{Kind: Purebred, Generation: 1, Detail: "a first generation purebred"},
		{Kind: Purebred, Generation: secondGeneration, Detail: "a second generation purebred"},
		{Kind: Purebred, Detail: "a third, fourth or later generation purebred"},
		{Kind: Hybrid, Detail: "a hybrid with a baseline human"},
		{Kind: Hybrid, Detail: "a hybrid with a baseline human"},
		{Kind: Compound, Detail: "a compound of two engineered species"},
	}
}

// GeneticStatusTable is this species' own status table, or the shared one.
func (s Species) GeneticStatusTable() []GeneticStatus {
	if s.Genetics == nil || len(s.Genetics.Status) == 0 {
		return DefaultGeneticStatus()
	}

	return s.Genetics.Status
}

// Data is a whole setting file.
type Data struct {
	SchemaVersion int         `json:"schemaVersion"`
	Name          string      `json:"name"`
	Species       []Species   `json:"species,omitempty"`
	Subsectors    []Subsector `json:"subsectors"`

	// PresentYear is the campaign's now, which Step 7 measures a world's
	// settledYear against: two of its paths are gated on whether the
	// homeworld "has been colonized or established for 100+ standard years"
	// (p. 76).
	//
	// The book gives two answers. p. 43 says "the present day of 2350" and
	// p. 124 says "the people living in 2345". Zero means p. 43's, which is
	// the page that talks about dates. ERRATA E-27.
	PresentYear int `json:"presentYear,omitempty"`

	// Hash is the content hash of the file this was read from, stamped into
	// every record so that a replay against different data is refused
	// rather than silently diverging at Step 3.
	Hash string `json:"-"`

	// Sample reports that this is the repository's invented data rather than
	// a transcription of the book.
	Sample bool `json:"-"`
}

// DefaultPresentYear is p. 43's: "the character could not have immigrated
// to Clement Sector after 2331 when the Conduit collapsed and they must
// have present in Clement Sector for the 19 years between the 2331 collapse
// and the present day of 2350".
const DefaultPresentYear = 2350

// Now is the campaign's present year, which is the file's where it gives
// one and p. 43's where it does not.
func (d *Data) Now() int {
	if d.PresentYear == 0 {
		return DefaultPresentYear
	}

	return d.PresentYear
}

// SpeciesNamed finds a species by name.
func (d *Data) SpeciesNamed(name string) (Species, bool) {
	for _, species := range d.Species {
		if species.Name == name {
			return species, true
		}
	}

	return Species{}, false
}

// Subsector finds a subsector by name.
func (d *Data) Subsector(name string) (Subsector, bool) {
	for _, sub := range d.Subsectors {
		if sub.Name == name {
			return sub, true
		}
	}

	return Subsector{}, false
}

// World finds a world by name, anywhere in the data.
func (d *Data) World(name string) (World, Subsector, bool) {
	for _, sub := range d.Subsectors {
		for _, world := range sub.Worlds {
			if world.Name == name {
				return world, sub, true
			}
		}
	}

	return World{}, Subsector{}, false
}
