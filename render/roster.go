package render

import (
	"fmt"
	"slices"
	"strings"

	"github.com/philoserf/cschargen/chargen"
)

// unnamed is what to call a character until they have a name. An
// auto-generated one has none: FR12 leaves Step 20's four fields empty
// rather than inventing them.
const unnamed = "(unnamed)"

// topSkills is how many of a character's best skills the roster prints.
// Four is what fits on a line beside the careers, and it is enough to tell
// a pilot from a medic, which is what casting a crew asks of it.
const topSkills = 4

// Roster is a batch seen from above: one entry per character, short enough
// to skim and long enough to cast from.
//
// A referee choosing a crew from a pool of a hundred reads, they do not
// study. The sheet is the right document for one character and the wrong
// one for a hundred, and before this existed the answer was a dozen lines
// of somebody's own Python over the JSONL.
//
// What the line carries is what casting a ship's crew asked for: how old
// they are, what they did, what they are best at, and what they have. A
// referee casting patrons wants CHA and a referee casting antagonists wants
// the Enemies, and a fixed line that is right for the common case beats a
// template language nobody wants to learn.
func Roster(characters []*chargen.Character) string {
	out := &strings.Builder{}

	for i, character := range characters {
		writeRosterEntry(out, i+1, character)
	}

	return out.String()
}

// writeRosterEntry is one character's three lines.
func writeRosterEntry(out *strings.Builder, nth int, character *chargen.Character) {
	state := character.State

	fmt.Fprintf(out, "%3d  %s, age %d, %s  %s\n",
		nth, rosterName(character), state.Age, rosterHomeworld(state),
		rosterCareers(state.Services))
	fmt.Fprintf(out, "     %s\n", rosterSkills(state.Skills))
	fmt.Fprintf(out, "     %s\n", rosterBelongings(state))
}

// rosterName is the character's name, or unnamed until they have one.
func rosterName(character *chargen.Character) string {
	if name := character.State.Finishing.Name; name != "" {
		return name
	}

	return unnamed
}

// rosterHomeworld is where the character was born, which is the first entry
// of the homeworld history: a Colonist mishap can move them later, and where
// they are from is what places them.
//
// It is the most useful word on the line for a setting the referee supplied.
// It carries the primary language, the tech level they grew up at, and
// whether engineered people are free where they were born -- and without it,
// casting a scene from a sector meant opening every record to find out where
// anyone came from.
func rosterHomeworld(state chargen.State) string {
	if len(state.Homeworlds) == 0 {
		return "no homeworld"
	}

	return "of " + state.Homeworlds[0].World
}

// rosterCareers is the career history in one phrase: "Belter 2t → Pirate 1t".
func rosterCareers(services []chargen.Service) string {
	if len(services) == 0 {
		return "no career"
	}

	parts := make([]string, 0, len(services))
	for _, served := range services {
		parts = append(parts, fmt.Sprintf("%s %dt", served.Career, served.Terms))
	}

	return strings.Join(parts, " → ")
}

// rosterSkills is the character's best few, highest first. Ties are broken
// by name so that two runs of the same seed read the same way.
func rosterSkills(held []chargen.Skill) string {
	if len(held) == 0 {
		return "no skills"
	}

	ranked := make([]chargen.Skill, len(held))
	copy(ranked, held)

	sortSkillsByLevel(ranked)

	shown := min(topSkills, len(ranked))
	parts := make([]string, 0, shown)

	for _, skill := range ranked[:shown] {
		parts = append(parts, fmt.Sprintf("%s-%d", skill.Full(), skill.Level))
	}

	return strings.Join(parts, ", ")
}

// rosterBelongings is money, the ties a referee can use, and the stash.
//
// The ties are counted without the family, because 90% of a character's
// relationships are relatives and a referee casting an NPC wants the ones
// a plot can hang on. The sheet itself keeps both; this line is the
// summary.
func rosterBelongings(state chargen.State) string {
	useful := 0

	for _, tie := range state.Ties {
		if tie.Origin != chargen.FamilyOrigin {
			useful++
		}
	}

	return fmt.Sprintf("Cr%d, %s beyond the family, %d carried",
		state.Credits, plural(useful, "tie", "ties"), len(state.Stash))
}

// sortSkillsByLevel puts the highest first, breaking ties by name so that
// two runs of one seed read the same way.
func sortSkillsByLevel(skills []chargen.Skill) {
	slices.SortFunc(skills, func(left, right chargen.Skill) int {
		if left.Level != right.Level {
			return right.Level - left.Level
		}

		return strings.Compare(left.Full(), right.Full())
	})
}

// plural counts a thing and names it, because "1 ties" reads as a bug in
// the tool rather than as a fact about the character.
func plural(n int, one, many string) string {
	if n == 1 {
		return "1 " + one
	}

	return fmt.Sprintf("%d %s", n, many)
}
