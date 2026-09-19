package chargen

import "slices"

// SchemaVersion is the shape of the records this engine writes. It tracks
// the records, not the precision of any document describing them: a change
// that narrows the shape to what the engine already produced is a
// clarification, and one that would invalidate a record the current engine
// writes is a bump.
const SchemaVersion = 1

// Ruleset names the printed artifact every page cite in this engine refers
// to. It is stamped into every record, because a cite is only checkable
// against the edition it was read from.
const Ruleset = "Clement Sector Core Character Creation Book, (c) 2026 Independence Games"

// RNG records how the dice were driven, so that a record made under one
// generator is not silently replayed under another.
type RNG struct {
	Algorithm string `json:"algorithm"`
	Seed      uint64 `json:"seed"`
}

// SettingData identifies the external file the world and species tables
// were read from (docs/PRD.md, Product Identity constraint). Milestone 2
// fills Hash; until then Sample reports that the character was generated
// against the invented sample data and is not set on real worlds.
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
	Name        string `json:"name,omitempty"`
	Species     string `json:"species"`
	Career      string `json:"career,omitempty"`
	Homeworld   string `json:"homeworld,omitempty"`
	TechLevel   int    `json:"techLevel"`
	MaxTerms    int    `json:"maxTerms"`
	TermLimit   int    `json:"termLimit,omitempty"`
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
