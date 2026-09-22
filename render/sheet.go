// Package render turns a character record into something a person reads:
// the Markdown character sheet, and the lifepath transcript that walks the
// generation record event by event.
//
// The JSON record is the source of truth. Everything here is a view of it,
// so a render never computes a rule -- if a number is not in the record,
// it does not appear on the sheet.
package render

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/philoserf/cschargen/chargen"
)

// Sheet renders the character sheet, modelled on the Long Form Character
// Sheet of pp. 324-333.
func Sheet(character *chargen.Character) string {
	var out strings.Builder

	name := character.Provenance.Inputs.Name
	if name == "" {
		name = unnamed
	}

	fmt.Fprintf(&out, "# %s\n\n", name)

	writeSummary(&out, character)
	writeOrigin(&out, character)
	writeFinishing(&out, character)
	writeGenetics(&out, character)
	writeFamily(&out, character)
	writeEducation(&out, character)
	writeCharacteristics(&out, character)
	writeSkills(&out, character)
	writeCareers(&out, character)
	writeTies(&out, character)
	writeBelongings(&out, character)
	writeProvenance(&out, character)

	return out.String()
}

func writeSummary(out *strings.Builder, character *chargen.Character) {
	state := character.State

	fmt.Fprintf(out, "**Species**: %s", character.Provenance.Inputs.Species)

	// An uplift's class is what their homeworld's technology could make of
	// them (p. 66), so it belongs beside the species rather than under it.
	if state.UpliftClass > 0 {
		fmt.Fprintf(out, ", Class %d", state.UpliftClass)
	}

	out.WriteString("  \n")
	fmt.Fprintf(out, "**Age**: %d  \n", state.Age)

	// "Apparent age is a game term used to indicate to the Players how old
	// a character might appear to be to the early 21st century onlooker"
	// (p. 124). It only says something the age does not once the chart on
	// p. 125 has an answer of its own.
	if band := state.ApparentAge; band.From != state.Age {
		fmt.Fprintf(out, "**Apparent age**: %s  \n", band)
	}

	// Where the character was born, which is the first entry of the
	// homeworld history rather than the inputs: the inputs say what was
	// asked for, and most characters ask for nothing.
	if len(state.Homeworlds) > 0 {
		fmt.Fprintf(out, "**Homeworld**: %s  \n", state.Homeworlds[0].World)
	}

	if state.Language != "" {
		fmt.Fprintf(out, "**Primary language**: %s  \n", state.Language)
	}

	fmt.Fprintf(out, "**Terms**: %d\n", len(state.Terms))

	// Aging is the one thing that can end a lifepath short of the term
	// limit (pp. 123-124), and a sheet that did not say so would show a
	// character with fewer terms than their homeworld allows and no reason
	// for it.
	switch state.Fate {
	case chargen.FateDied:
		fmt.Fprintf(out, "**Died** at age %d, during generation  \n", state.Age)
	case chargen.FateIncapacitated:
		fmt.Fprint(out, "**Incapacitated** by aging: generation ended here  \n")
	}

	fmt.Fprint(out, "\n")
}

// writeFinishing prints Step 20's four fields (pp. 129-130), and only the
// ones somebody supplied: the engine invents none, and a blank line is
// worse than no line.
func writeFinishing(out *strings.Builder, character *chargen.Character) {
	finishing := character.State.Finishing

	for _, field := range []struct{ label, value string }{
		{"Gender", finishing.Gender},
		{"Appearance", finishing.Appearance},
		{"Goals", finishing.Goals},
	} {
		if field.value != "" {
			fmt.Fprintf(out, "**%s**: %s  \n", field.label, field.value)
		}
	}

	if !finishing.Empty() {
		out.WriteString("\n")
	}
}

// writeGenetics prints how an engineered character came to be (pp. 62-65).
// A human and an uplift have none.
func writeGenetics(out *strings.Builder, character *chargen.Character) {
	genetics := character.State.Genetics
	if genetics == nil {
		return
	}

	fmt.Fprintf(out, "**Genetics**: %s  \n\n", genetics.Detail)
}

// writeFamily prints what Step 5 established (pp. 57-61). The relatives
// themselves are counted under Relationships with everyone else, because
// that is where their Relationship Ratings are; this is the household they
// came from.
func writeFamily(out *strings.Builder, character *chargen.Character) {
	family := character.State.Family
	if family == nil {
		return
	}

	out.WriteString("## Family\n\n")

	fmt.Fprintf(out, "Born to %s", family.Situation)

	if family.Detail != "" {
		fmt.Fprintf(out, ": %s", family.Detail)
	}

	out.WriteString(".\n")

	switch {
	case family.Firstborn:
		out.WriteString("The eldest child.\n")
	case family.Lastborn:
		out.WriteString("The youngest child.\n")
	}

	siblings := rolesOf(character, chargen.RoleSibling)
	if len(siblings) == 0 {
		out.WriteString("\nAn only child.\n\n")

		return
	}

	out.WriteString("\n")

	for _, sibling := range siblings {
		fmt.Fprintf(out, "- A sibling, %s\n", sibling.Detail)
	}

	out.WriteString("\n")
}

// rolesOf is every family tie of one role, in the order they were gained.
func rolesOf(character *chargen.Character, role string) []chargen.Tie {
	var found []chargen.Tie

	for _, tie := range character.State.Ties {
		if tie.Role == role {
			found = append(found, tie)
		}
	}

	return found
}

// writeEducation prints what Step 8 came to (pp. 85-104): every
// institution attempted, whether they were admitted, and what they left
// with.
func writeEducation(out *strings.Builder, character *chargen.Character) {
	history := character.State.Education
	if len(history) == 0 {
		return
	}

	out.WriteString("## Education\n\n")

	for _, held := range history {
		switch {
		case held.Degree != "" && held.Honors:
			fmt.Fprintf(out, "- %s: a %s in %s, with honours\n",
				held.Institution, held.Degree, held.Field)
		case held.Degree != "":
			fmt.Fprintf(out, "- %s: a %s in %s\n", held.Institution, held.Degree, held.Field)
		case held.Admitted:
			fmt.Fprintf(out, "- %s: admitted, left without a degree\n", held.Institution)
		default:
			fmt.Fprintf(out, "- %s: not admitted\n", held.Institution)
		}
	}

	out.WriteString("\n")
}

// writeCharacteristics prints the six with their modifiers. The modifier is
// computed here rather than stored, which is p. 15's "recalculated
// immediately" holding at the point of display too.
// writeOrigin prints the homeworld history, and only when there is a
// history to print: six of Colonist's eleven mishaps reassign the
// homeworld, so where a character has lived is part of what happened to
// them.
func writeOrigin(out *strings.Builder, character *chargen.Character) {
	// One homeworld is the birth world, already printed in the summary. A
	// table is worth its space only once there is a move to show.
	const aMove = 2

	if len(character.State.Homeworlds) < aMove {
		return
	}

	out.WriteString("## Where they have lived\n\n")
	out.WriteString("| World | Subsector | TL | From term | Why |\n")
	out.WriteString("| ----- | --------- | -- | --------- | --- |\n")

	for _, home := range character.State.Homeworlds {
		fmt.Fprintf(out, "| %s | %s | %d | %d | %s |\n",
			home.World, home.Subsector, home.TechLevel, home.FromTerm, home.Reason)
	}

	out.WriteString("\n")
}

func writeCharacteristics(out *strings.Builder, character *chargen.Character) {
	out.WriteString("## Characteristics\n\n")
	out.WriteString("| STR | DEX | END | INT | EDU | CHA |\n")
	out.WriteString("| --- | --- | --- | --- | --- | --- |\n|")

	for _, which := range chargen.CharacteristicOrder {
		score := character.State.Characteristics.Get(which)
		fmt.Fprintf(out, " %d (%+d) |", score, chargen.Modifier(score))
	}

	out.WriteString("\n\n")
}

func writeSkills(out *strings.Builder, character *chargen.Character) {
	out.WriteString("## Skills\n\n")

	if len(character.State.Skills) == 0 {
		out.WriteString("_None._\n\n")

		return
	}

	parts := make([]string, 0, len(character.State.Skills))
	for _, skill := range character.State.Skills {
		parts = append(parts, fmt.Sprintf("%s-%d", skill.Full(), skill.Level))
	}

	fmt.Fprintf(out, "%s\n\n", strings.Join(parts, ", "))
}

func writeCareers(out *strings.Builder, character *chargen.Character) {
	out.WriteString("## Career history\n\n")

	if len(character.State.Services) == 0 {
		out.WriteString("_No career service._\n\n")

		return
	}

	out.WriteString("| Career | Assignment | Terms | Rank | Left because |\n")
	out.WriteString("| ------ | ---------- | ----- | ---- | ------------ |\n")

	for _, service := range character.State.Services {
		reason := service.LeftBecause
		if reason == "" {
			reason = "still serving"
		}

		fmt.Fprintf(out, "| %s | %s | %d | %d | %s |\n",
			service.Career, service.Assignment, service.Terms, service.Rank, reason)
	}

	out.WriteString("\n")
}

func writeTies(out *strings.Builder, character *chargen.Character) {
	if len(character.State.Ties) == 0 {
		return
	}

	out.WriteString("## Relationships\n\n")

	made, family := splitTies(character.State.Ties)

	writeTiesMade(out, made)
	writeFamilyTies(out, family)

	out.WriteString("\n")
}

// splitTies separates the people a character met from the people they were
// born to.
//
// Across a sample pool, 90% of every character's relationships were family:
// parents, siblings, grandparents, aunts, uncles and cousins, each an Ally
// at 100 to 150 from Step 5. One line reading "Allies: 76" buried the one
// ally at 190 -- the colleague a plot can hang on -- among sixty-two
// cousins, and gave a referee no way to find them.
//
// The record has always known which is which: every tie carries an origin
// on it. The sheet was throwing it away.
func splitTies(ties []chargen.Tie) ([]chargen.Tie, []chargen.Tie) {
	made := make([]chargen.Tie, 0, len(ties))
	family := make([]chargen.Tie, 0, len(ties))

	for _, tie := range ties {
		if tie.Origin == chargen.FamilyOrigin {
			family = append(family, tie)

			continue
		}

		made = append(made, tie)
	}

	return made, family
}

// writeTiesMade lists the relationships a character made, one per line and
// each naming where it came from -- a career, a school, or the stage of life
// it was made in. These are the ones a referee reads.
func writeTiesMade(out *strings.Builder, made []chargen.Tie) {
	if len(made) == 0 {
		return
	}

	out.WriteString("**Made along the way**\n\n")

	for _, kind := range tieKinds() {
		for _, tie := range made {
			if tie.Kind != kind.kind {
				continue
			}

			// The Relationship Rating is what says whether a Contact is
			// nearly an Ally or nearly lost (p. 320), so a name on its own
			// leaves out the half of a relationship that moves.
			fmt.Fprintf(out, "- %s (%d) — %s\n", kind.singular, tie.Rating, tie.Origin)
		}
	}

	out.WriteString("\n")
}

// writeFamilyTies counts the family rather than listing them. A referee
// wants to know a character has a large family and what it thinks of them,
// not to read sixty-two lines each saying "cousin".
func writeFamilyTies(out *strings.Builder, family []chargen.Tie) {
	if len(family) == 0 {
		return
	}

	counts := map[string]int{}
	ratings := map[string][]int{}

	for _, tie := range family {
		counts[tie.Kind]++

		ratings[tie.Kind] = append(ratings[tie.Kind], tie.Rating)
	}

	out.WriteString("**Family**\n\n")

	for _, kind := range tieKinds() {
		if counts[kind.kind] == 0 {
			continue
		}

		fmt.Fprintf(out, "- %s: %d (%s)\n",
			kind.plural, counts[kind.kind], joinRatings(ratings[kind.kind]))
	}
}

// tieKinds is the four in the order p. 120 lists them, so the sheet reads
// the same way twice running whatever a map does. The plurals are spelt out
// rather than made by adding an s, which would give "Allys".
func tieKinds() []struct{ kind, singular, plural string } {
	return []struct{ kind, singular, plural string }{
		{"ally", "Ally", "Allies"},
		{"contact", "Contact", "Contacts"},
		{"rival", "Rival", "Rivals"},
		{"enemy", "Enemy", "Enemies"},
	}
}

// joinRatings renders a kind's Relationship Ratings, sorted so the sheet
// reads the same way twice running and grouped so that a family of twelve
// does not print the same number seven times.
func joinRatings(ratings []int) string {
	sorted := slices.Sorted(slices.Values(ratings))

	var (
		parts []string
		run   int
	)

	for i, rating := range sorted {
		run++

		if i+1 < len(sorted) && sorted[i+1] == rating {
			continue
		}

		part := strconv.Itoa(rating)
		if run > 1 {
			part += " x" + strconv.Itoa(run)
		}

		parts = append(parts, part)
		run = 0
	}

	return strings.Join(parts, ", ")
}

// stash lists what a character owns, naming the credit value of anything
// the book priced and nothing else: a weapon and a suit of armour are
// referee-valued, and a zero there means the book gave no number rather
// than that the thing is worthless.
func stash(held []chargen.Possession) string {
	parts := make([]string, 0, len(held))

	for _, item := range held {
		if item.Value == 0 {
			parts = append(parts, item.Item)

			continue
		}

		parts = append(parts, fmt.Sprintf("%s (%d credits)", item.Item, item.Value))
	}

	return strings.Join(parts, ", ")
}

func writeBelongings(out *strings.Builder, character *chargen.Character) {
	state := character.State

	if state.Credits == 0 && len(state.Stash) == 0 &&
		len(state.Conditions) == 0 && len(state.Injuries) == 0 {
		return
	}

	out.WriteString("## Belongings and injuries\n\n")

	if state.Credits != 0 {
		fmt.Fprintf(out, "**Credits**: %d\n\n", state.Credits)
	}

	if len(state.Stash) > 0 {
		fmt.Fprintf(out, "**Stash**: %s\n\n", stash(state.Stash))
	}

	if len(state.Conditions) > 0 {
		fmt.Fprintf(out, "**Carrying**: %s\n\n", strings.Join(state.Conditions, ", "))
	}

	writeInjuries(out, state.Injuries)
}

// writeInjuries lists what a character carries from the Injury table,
// marking the ones that never healed.
func writeInjuries(out *strings.Builder, injuries []chargen.Injury) {
	for _, injury := range injuries {
		permanence := ""
		if injury.Permanent {
			permanence = " (permanent)"
		}

		fmt.Fprintf(out, "- Term %d: %s%s\n", injury.Term, injury.Detail, permanence)
	}

	if len(injuries) > 0 {
		out.WriteString("\n")
	}
}

// writeProvenance is what makes a sheet checkable: the seed it came from,
// the ruleset it was read against, and any ERRATA readings that applied.
func writeProvenance(out *strings.Builder, character *chargen.Character) {
	prov := character.Provenance

	out.WriteString("---\n\n")
	fmt.Fprintf(out, "Generated by cschargen %s from seed %d.  \n",
		prov.EngineVersion, prov.RNG.Seed)
	fmt.Fprintf(out, "Ruleset: %s  \n", prov.Ruleset)

	if prov.SettingData.Sample {
		out.WriteString("Setting data: the repository's invented sample, not the published worlds.  \n")
	}

	if len(prov.Deviations) > 0 {
		fmt.Fprintf(out, "Deviations applied (see ERRATA.md): %s\n", strings.Join(prov.Deviations, ", "))
	}
}
