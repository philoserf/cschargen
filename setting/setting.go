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
// Kingston grants "Survival (High Pressure, Mountains, Plains, or Ocean)
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
	Roll *[2]int `json:"roll,omitempty"`

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
	return w.Roll != nil
}

// Covers reports whether a d100 result lands on this world.
func (w World) Covers(roll int) bool {
	if w.Roll == nil {
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
// refer to it. The rules for playing one are milestone 6; this is the
// vocabulary the permission columns are written against.
type Species struct {
	Name string `json:"name"`

	// Kind is "engineered" for an engineered human or "uplift" for an
	// uplifted animal.
	Kind string `json:"kind"`
}

// Data is a whole setting file.
type Data struct {
	SchemaVersion int         `json:"schemaVersion"`
	Name          string      `json:"name"`
	Species       []Species   `json:"species,omitempty"`
	Subsectors    []Subsector `json:"subsectors"`

	// Hash is the content hash of the file this was read from, stamped into
	// every record so that a replay against different data is refused
	// rather than silently diverging at Step 3.
	Hash string `json:"-"`

	// Sample reports that this is the repository's invented data rather than
	// a transcription of the book.
	Sample bool `json:"-"`
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
