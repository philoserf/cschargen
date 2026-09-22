package chargen

import "slices"

// SchemaVersion is the shape of the records this engine writes. It tracks
// the records, not the precision of any document describing them: a change
// that narrows the shape to what the engine already produced is a
// clarification, and one that would invalidate a record the current engine
// writes is a bump.
// Version 2 changed "stash" from a list of names to a list of
// possessions, each of which may carry the value the book rolled for it.
const SchemaVersion = 2

// Ruleset names the printed artifact every page cite in this engine refers
// to. It is stamped into every record, because a cite is only checkable
// against the edition it was read from.
const Ruleset = "Clement Sector Core Character Creation Book, (c) 2026 Independence Games"

// PolicyVersion identifies POLICY.md, the auto-mode decision table, and is
// the third of the three strings a record stamps. It is here with the other
// two because it describes the record rather than the command: a different
// policy is a different character from the same seed.
//
// Bump it when POLICY.md changes. 0.3.2 is the change that showed why the
// rule needs a gate rather than a habit -- the subsector row moved from
// "take the first subsector the file lists" to "a throw that misses is
// thrown again", which is a different character from the same seed against
// a file with choose-only subsectors, and neither this constant nor the
// document's own version line moved with it.
//
// Options.PolicyVersion still carries it into a record, because the tests
// stamp a version of their own rather than claiming a real build's.
const PolicyVersion = "0.3.2"

// RNG records how the dice were driven, so that a record made under one
// generator is not silently replayed under another.
type RNG struct {
	Algorithm string `json:"algorithm"`
	Seed      uint64 `json:"seed"`
}

// SettingData identifies the external file the world and species tables
// were read from. Hash is that file's content hash, so a record refuses to
// replay against a different transcription; Sample reports that the
// character was generated against the invented sample data and is not set
// on real worlds.
type SettingData struct {
	Name   string `json:"name"`
	Hash   string `json:"hash,omitempty"`
	Sample bool   `json:"sample"`
}

// Inputs are the engine inputs a seed does not capture. They are recorded
// because replay needs them: a forced career holds the first career's
// option list to a single entry, so a replay that did not know a record was
// forced would offer the full list and read the recorded index against it.
type Inputs struct {
	Name    string `json:"name,omitempty"`
	Species string `json:"species"`
	Career  string `json:"career,omitempty"`

	// Subsector and Homeworld are Steps 3 and 4 asked for rather than
	// rolled: "the Referee may choose the subsector for the character"
	// (p. 39) and "simply choose a world that fits the character concept
	// you have in mind" (p. 40). A homeworld names its own subsector, so
	// giving both is only useful to pick a subsector and roll inside it.
	//
	// Both are empty on a character who asked for neither, which is most
	// of them. Where a character was actually born is State.Homeworlds,
	// the way Career here is the career asked for and State.Services is
	// the career served: Homeworld used to hold the result as well, and a
	// replay could not then tell a rolled record from a chosen one.
	Subsector     string `json:"subsector,omitempty"`
	Homeworld     string `json:"homeworld,omitempty"`
	TechLevel     int    `json:"techLevel"`
	MaxTerms      int    `json:"maxTerms"`
	TermLimit     int    `json:"termLimit,omitempty"`
	SkipFamily    bool   `json:"skipFamily,omitempty"`
	SkipYouth     bool   `json:"skipYouth,omitempty"`
	SkipTeenage   bool   `json:"skipTeenage,omitempty"`
	SkipEducation bool   `json:"skipEducation,omitempty"`

	// The other three of Step 20's four fields (pp. 129-130). Name is
	// above, because it has been an input from the start.
	Gender      string `json:"gender,omitempty"`
	Appearance  string `json:"appearance,omitempty"`
	Goals       string `json:"goals,omitempty"`
	Interactive bool   `json:"interactive"`
}

// Provenance is everything a record carries about how it was made, as
// distinct from what it says about a character.
type Provenance struct {
	SchemaVersion int         `json:"schemaVersion"`
	Ruleset       string      `json:"ruleset"`
	EngineVersion string      `json:"engineVersion"`
	PolicyVersion string      `json:"policyVersion"`
	RNG           RNG         `json:"rng"`
	SettingData   SettingData `json:"settingData"`
	Inputs        Inputs      `json:"inputs"`

	// Deviations are the ERRATA.md identifiers that applied to this
	// character. A record made under one reading stays auditable after the
	// reading changes, because it names the reading it was made under.
	Deviations []string `json:"deviations,omitempty"`
}

// Character is the generation record. The JSON is the source of truth; the
// Markdown sheet is a render of it.
type Character struct {
	Provenance Provenance `json:"provenance"`
	State      State      `json:"state"`
	Events     []Event    `json:"events"`
}

// Deviate records that an ERRATA.md reading applied to this character,
// once, however many times it fires.
func (p *Provenance) Deviate(id string) {
	if slices.Contains(p.Deviations, id) {
		return
	}

	p.Deviations = append(p.Deviations, id)
}
