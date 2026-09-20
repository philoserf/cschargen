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
	"strings"

	"github.com/philoserf/cschargen/chargen"
)

// Sheet renders the character sheet, modelled on the Long Form Character
// Sheet of pp. 324-333.
func Sheet(character *chargen.Character) string {
	var out strings.Builder

	name := character.Provenance.Inputs.Name
	if name == "" {
		name = "(unnamed)"
	}

	fmt.Fprintf(&out, "# %s\n\n", name)

	writeSummary(&out, character)
	writeOrigin(&out, character)
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

	fmt.Fprintf(out, "**Species**: %s  \n", character.Provenance.Inputs.Species)
	fmt.Fprintf(out, "**Age**: %d  \n", state.Age)

	// "Apparent age is a game term used to indicate to the Players how old
	// a character might appear to be to the early 21st century onlooker"
	// (p. 124). It only says something the age does not once the chart on
	// p. 125 has an answer of its own.
	if band := state.ApparentAge; band.From != state.Age {
		fmt.Fprintf(out, "**Apparent age**: %s  \n", band)
	}

	if world := character.Provenance.Inputs.Homeworld; world != "" {
		fmt.Fprintf(out, "**Homeworld**: %s  \n", world)
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

	counts := map[string]int{}
	for _, tie := range character.State.Ties {
		counts[tie.Kind]++
	}

	out.WriteString("## Relationships\n\n")

	// The four in the order p. 120 lists them, so the sheet reads the same
	// way twice running whatever the map does. Plurals are spelt out rather
	// than made by adding an s, which would give "Allys".
	plurals := []struct{ kind, plural string }{
		{"ally", "Allies"},
		{"contact", "Contacts"},
		{"rival", "Rivals"},
		{"enemy", "Enemies"},
	}

	for _, entry := range plurals {
		if counts[entry.kind] == 0 {
			continue
		}

		fmt.Fprintf(out, "- %s: %d\n", entry.plural, counts[entry.kind])
	}

	out.WriteString("\n")
}

func writeBelongings(out *strings.Builder, character *chargen.Character) {
	state := character.State

	if state.Credits == 0 && len(state.Stash) == 0 && len(state.Injuries) == 0 {
		return
	}

	out.WriteString("## Belongings and injuries\n\n")

	if state.Credits != 0 {
		fmt.Fprintf(out, "**Credits**: %d\n\n", state.Credits)
	}

	if len(state.Stash) > 0 {
		fmt.Fprintf(out, "**Stash**: %s\n\n", strings.Join(state.Stash, ", "))
	}

	for _, injury := range state.Injuries {
		permanence := ""
		if injury.Permanent {
			permanence = " (permanent)"
		}

		fmt.Fprintf(out, "- Term %d: %s%s\n", injury.Term, injury.Detail, permanence)
	}

	if len(state.Injuries) > 0 {
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
		fmt.Fprintf(out, "Readings applied (see ERRATA.md): %s\n", strings.Join(prov.Deviations, ", "))
	}
}
