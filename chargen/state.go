package chargen

import (
	"slices"
	"strings"
)

// Skill is one skill a character holds. The book writes a skill as a name
// and, often, a specialty in parentheses -- Survival (Desert), Gun Combat
// (Any) -- and treats the two as one thing for the purpose of levelling.
type Skill struct {
	Name      string `json:"name"`
	Specialty string `json:"specialty,omitempty"`
	Level     int    `json:"level"`
}

// Full is the skill as the book writes it.
func (s Skill) Full() string {
	if s.Specialty == "" {
		return s.Name
	}

	return s.Name + " (" + s.Specialty + ")"
}

// Tie is an Ally, Contact, Rival or Enemy, with the Relationship Rating the
// Life Events table writes to (ERRATA E-5).
type Tie struct {
	Kind   string `json:"kind"`
	Origin string `json:"origin"` // the career the tie came from
	Rating int    `json:"rating"`
}

// Injury is one result of the Injury table (p. 119). Permanent is set at
// the moment of the injury rather than recomputed at muster out: recovery
// is charged against credits held "at the time of the injury", and a
// character who has not mustered out has none (ERRATA E-4).
type Injury struct {
	Detail    string `json:"detail"`
	Permanent bool   `json:"permanent"`
	Term      int    `json:"term"`
}

// Term is one four-year period in a career (p. 12).
type Term struct {
	Career     string `json:"career"`
	Assignment string `json:"assignment"`
	Number     int    `json:"number"` // the character's term count, not the career's
	Survived   bool   `json:"survived"`
	Rank       int    `json:"rank"`
	Event      string `json:"event,omitempty"`
	Mishap     string `json:"mishap,omitempty"`
}

// Service is one spell in one career.
type Service struct {
	Career       string `json:"career"`
	Assignment   string `json:"assignment"`
	Terms        int    `json:"terms"`
	Rank         int    `json:"rank"`
	Commissioned bool   `json:"commissioned,omitempty"`
	LeftBecause  string `json:"leftBecause,omitempty"`
}

// BenefitBatch is a run of mustering-out rolls carrying one modifier.
// Benefit rolls are a queue of scoped batches rather than a count plus a
// career modifier, because the book grants both "a +1 modifier to all
// Benefit rolls made in this career" and "three Benefit rolls at +1", and
// the two are not the same thing (ERRATA E-3).
type BenefitBatch struct {
	Career   string `json:"career"`
	Rolls    int    `json:"rolls"`
	Modifier int    `json:"modifier"`
}

// PendingModifier is an adjustment waiting for the throw it applies to.
// Named for the waiting rather than the adjustment, because [Modifier] is
// already the die modifier a characteristic confers.
type PendingModifier struct {
	Applies string `json:"applies"`
	Value   int    `json:"value"`
	Detail  string `json:"detail"`
}

// State is everything generation has established about a character, as
// distinct from how it was established -- which is the event log's job.
type State struct {
	Characteristics Characteristics `json:"characteristics"`
	Language        string          `json:"language,omitempty"`
	Skills          []Skill         `json:"skills"`
	Ties            []Tie           `json:"ties,omitempty"`
	Credits         int             `json:"credits"`
	Stash           []string        `json:"stash,omitempty"`
	Injuries        []Injury        `json:"injuries,omitempty"`
	Homeworlds      []Homeworld     `json:"homeworlds,omitempty"`
	Services        []Service       `json:"services,omitempty"`
	Terms           []Term          `json:"terms,omitempty"`
	Benefits        []BenefitBatch  `json:"benefits,omitempty"`
	Age             int             `json:"age"`

	// Fate is set only where aging ended generation (pp. 123-124). Empty
	// is the ordinary case, and death is a value rather than an error.
	Fate Fate `json:"fate,omitempty"`
}

// startingAge is where a lifepath begins: "career terms are normally
// assumed to begin at age 18, the typical age of majority" (p. 42).
const startingAge = 18

// termYears is the length of a term: "Each career is divided into four-year
// periods known as terms" (p. 12).
const termYears = 4

// Skill finds a skill by name and specialty.
func (s *State) Skill(name, specialty string) (Skill, bool) {
	for _, held := range s.Skills {
		if held.Name == name && held.Specialty == specialty {
			return held, true
		}
	}

	return Skill{}, false
}

// SkillLevel reports the best level the character holds in a skill under
// any specialty, which is what a check like "roll Melee (Any) 8+" is
// rolled against.
func (s *State) SkillLevel(name string) int {
	best := -1

	for _, held := range s.Skills {
		if held.Name == name && held.Level > best {
			best = held.Level
		}
	}

	return best
}

// Has reports whether the character holds a skill at level 0 or better
// under any specialty.
func (s *State) Has(name string) bool {
	return s.SkillLevel(name) >= 0
}

// GainSkill grants a level, or records the skill at level 0 where the
// character does not hold it at all. The two are different: p. 117's first
// term grants "all the skills in the Service Skills table at level 0 that
// they do not already possess", where every other award is a level.
func (s *State) GainSkill(name, specialty string, level int) Skill {
	for i, held := range s.Skills {
		if held.Name == name && held.Specialty == specialty {
			s.Skills[i].Level += level

			return s.Skills[i]
		}
	}

	added := Skill{Name: name, Specialty: specialty, Level: level}

	s.Skills = append(s.Skills, added)

	slices.SortFunc(s.Skills, func(left, right Skill) int {
		if left.Name != right.Name {
			return strings.Compare(left.Name, right.Name)
		}

		return strings.Compare(left.Specialty, right.Specialty)
	})

	return added
}

// Service finds the character's current spell in a career, which is the
// last one recorded. A character may serve in the same career twice --
// ejected by a mishap and later accepted back -- and those are two spells,
// not one: they have their own term counts, their own rank, and their own
// reason for ending. Returning the first would write every update to the
// spell that is already over.
func (s *State) Service(name string) (*Service, bool) {
	for i := len(s.Services) - 1; i >= 0; i-- {
		if s.Services[i].Career == name {
			return &s.Services[i], true
		}
	}

	return nil, false
}
