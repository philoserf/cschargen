package career

// Step 7: Teenage Events (pp. 75-85).
//
// Four paths of nineteen rows each, rolled twice -- ages 13-15 and 16-18 --
// and, like the youth paths, "this does not have to be the same table as the
// previous roll" (p. 76).
//
// The gates differ from Step 6's. Paths 1 and 2 are the two halves of one
// condition -- whether the homeworld has been settled a hundred standard
// years -- so exactly one of them is open to any character, and paths 3 and
// 4 are characteristic gates on top of that. So a character always has one
// path and may have three.
//
// The uplift and enslaved variants (pp. 76, 84) belong to the species work.

// TeenagePathGate is what opens a teenage path, which is not always a
// characteristic.
type TeenagePathGate int

const (
	// GateSettledLong is Path 1: "Homeworld is a location which has been
	// colonized or established for 100+ standard years."
	GateSettledLong TeenagePathGate = iota

	// GateSettledRecently is Path 2, the other half of the same condition.
	GateSettledRecently

	// GateCharacteristic is Paths 3 and 4, which read Requires.
	GateCharacteristic
)

// SettledYearsForPathOne is the hundred standard years of p. 76.
const SettledYearsForPathOne = 100

// TeenagePath is one of the four columns of pp. 76-83.
type TeenagePath struct {
	Name     string
	Cite     string
	Gate     TeenagePathGate
	Requires []Check
	Rows     []EventRow
}

// Open reports whether this path is available. settledYears is how long the
// homeworld has been colonized, which decides Paths 1 and 2 between them.
func (p TeenagePath) Open(score func(string) int, settledYears int) bool {
	switch p.Gate {
	case GateSettledLong:
		return settledYears >= SettledYearsForPathOne
	case GateSettledRecently:
		return settledYears < SettledYearsForPathOne
	case GateCharacteristic:
	}

	for _, check := range p.Requires {
		if score(check.Characteristic) >= check.Number {
			return true
		}
	}

	return false
}

// TeenagePaths is the four, in the order the book prints them.
func TeenagePaths() []TeenagePath {
	return []TeenagePath{
		{Name: "Path 1", Cite: "pp. 76-78", Gate: GateSettledLong, Rows: teenagePathOne()},
		{Name: "Path 2", Cite: "pp. 78-79", Gate: GateSettledRecently, Rows: teenagePathTwo()},
		{
			Name: "Path 3", Cite: "pp. 80-81", Gate: GateCharacteristic,
			Requires: []Check{
				{Characteristic: "STR", Number: 9},
				{Characteristic: "DEX", Number: 9},
				{Characteristic: "END", Number: 9},
			},
			Rows: teenagePathThree(),
		},
		{
			Name: "Path 4", Cite: "pp. 82-83", Gate: GateCharacteristic,
			Requires: []Check{
				{Characteristic: "INT", Number: 9},
				{Characteristic: "EDU", Number: 9},
			},
			Rows: teenagePathFour(),
		},
	}
}

// teenageLifeEvent is result 10 on every one of the four paths.
func teenageLifeEvent() EventRow {
	return EventRow{
		Summary: "a teenage life event",
		Effects: []Effect{{
			Kind:   EffectTeenageLifeEvent,
			Detail: "roll on the Teenage Life Events table",
		}},
	}
}

// Four rows are printed identically on more than one path, so they are
// written once here.

// localRivalry is result 8 on Paths 1 and 2.
func localRivalry() EventRow {
	return EventRow{
		Summary: "a rivalry with a peer that sharpens you",
		Effects: []Effect{relationshipAt(Rival, 1, -50), anyCharacteristic()},
	}
}

// aGoodFriend is result 9 on Paths 1 and 2.
func aGoodFriend() EventRow {
	return EventRow{
		Summary: "a friendship that makes the teenage years bearable",
		Effects: []Effect{relationshipAt(Contact, 1, 90)},
	}
}

// competence is result 13 on Paths 1 and 2.
func competence() EventRow {
	return EventRow{
		Summary: "doing something well, and being noticed for it",
		Effects: []Effect{chr("CHA", 1), aSkillYouAlreadyHave()},
	}
}

// success is result 14 on Paths 1 and 2.
func success() EventRow {
	return EventRow{
		Summary: "something accomplished and genuinely proud of",
		Effects: []Effect{chr("EDU", 1), aSkillYouAlreadyHave()},
	}
}

// opportunity is result 16 on Paths 1 and 2.
func opportunity() EventRow {
	return EventRow{
		Summary: "selected for a program, apprenticeship, or special responsibility",
		Effects: []Effect{
			relationshipAt(Contact, 1, 50),
			unimplemented("choose any career, roll once on its Service Skills table, " +
				"and take +2 to enlist in it later"),
		},
	}
}

// trustedRole is result 17 on Paths 1 and 2.
func trustedRole() EventRow {
	return EventRow{
		Summary: "adults who begin seeing you as dependable",
		Effects: []Effect{skill("Leadership"), chr("CHA", 1)},
	}
}

// strongSocialFoundation is result 18 on Paths 1 and 2.
func strongSocialFoundation() EventRow {
	return EventRow{
		Summary: "lasting connections built during these years",
		Effects: []Effect{chr("CHA", 1), relationshipRolledAt(Contact, "1d3", 50)},
	}
}

// fortunateYouth is result 20 on Paths 1 and 2.
func fortunateYouth() EventRow {
	return EventRow{
		Summary: "teenage years shaped by support, success and belonging",
		Effects: []Effect{
			chr("CHA", 1), chr("EDU", 1),
			anySkill(),
		},
	}
}

// broadeningHorizons is result 15 on Path 1, 19 on Path 2 and 11 on Path 3:
// the first journey offworld.
func broadeningHorizons() EventRow {
	return EventRow{
		Summary: "travel beyond your world for the first time",
		Effects: []Effect{unimplemented(
			"choose a world within three parsecs, take one of its background skills, " +
				"and a level in its primary language where it differs from your own")},
	}
}

// familyObligation is result 3 on Path 3 and result 5 on Path 4.
func familyObligation() EventRow {
	return EventRow{
		Summary: "a relative who needs your support, and adulthood learned early",
		Effects: []Effect{
			chr("EDU", -1),
			pick("choose what the duty taught you",
				opt("Admin", skill("Admin")),
				opt("Medic (First Aid)", skill("Medic", "First Aid"))),
		},
	}
}

// authorityConflict is result 5 on Path 3 and result 6 on Path 4.
func authorityConflict() EventRow {
	return EventRow{
		Summary: "an official who takes a disliking to you",
		Effects: []Effect{relationshipAt(Enemy, 1, -115)},
	}
}

// educationalDirection is result 9 on Path 3 and result 8 on Path 4.
func educationalDirection() EventRow {
	return EventRow{
		Summary: "a mentor who helps you clarify your next steps",
		Effects: []Effect{
			relationshipAt(Contact, 1, 50),
			throwModifier("admission to any higher education", 2),
		},
	}
}

// earnedRecommendation is result 15 on Paths 3 and 4.
func earnedRecommendation() EventRow {
	return EventRow{
		Summary: "a recommendation from somebody with standing",
		Effects: []Effect{
			relationshipAt(Contact, 1, 90),
			throwModifier("enlistment in the first career you choose", 2),
		},
	}
}

// socialNetwork is result 17 on Paths 3 and 4.
func socialNetwork() EventRow {
	return EventRow{
		Summary: "lasting connections by adolescence",
		Effects: []Effect{chr("CHA", 2), relationshipRolledAt(Contact, "1d3", 70)},
	}
}

// opportunityKnocks is result 18 on Paths 3 and 4.
func opportunityKnocks() EventRow {
	return EventRow{
		Summary: "a scholarship, recruitment offer, or transport slot",
		Effects: []Effect{unimplemented(
			"move immediately to Step 8 or Step 9, and enlist or be admitted " +
				"automatically wherever you qualify")},
	}
}

// aSkillYouAlreadyHave is "a level in any skill which you already possess",
// which several results grant and the engine records rather than resolves:
// the choice is over whatever the character holds, which is not a list any
// table prints.
func aSkillYouAlreadyHave() Effect {
	return raiseHeldSkill()
}

// anyCharacteristic is "+1 to any of your Characteristics".
func anyCharacteristic() Effect {
	return pick("raise a characteristic",
		opt("STR", chr("STR", 1)),
		opt("DEX", chr("DEX", 1)),
		opt("END", chr("END", 1)),
		opt("INT", chr("INT", 1)),
		opt("EDU", chr("EDU", 1)),
		opt("CHA", chr("CHA", 1)))
}

// TeenageLifeEvents is the d6 table of p. 85, indexed 0-5. Every path
// reaches it at result 10.
func TeenageLifeEvents() []EventRow {
	return []EventRow{
		{Summary: "an injury", Effects: []Effect{injury(1)}},
		{
			Summary: "a death among family or friends",
			Effects: []Effect{loseTie(Contact, Ally)},
		},
		{
			Summary: "a relationship that ends badly",
			Effects: []Effect{loseTie(Contact, Ally)},
		},
		{Summary: "a new friend", Effects: []Effect{relationship(Contact, 1, "")}},
		{
			// Every band moves up one: an Enemy becomes a Rival, a Rival a
			// Contact, a Contact an Ally, and with none of them a new
			// Contact arrives.
			Summary: "a relationship that improves",
			Effects: []Effect{improveARelationship()},
		},
		{
			Summary: "something wonderful",
			Effects: []Effect{pick("what came of it",
				opt("CHA", chr("CHA", 2)),
				opt("a new Ally", relationship(Ally, 1, "")),
				opt("a boon", credits("100000")),
				opt("a relationship improves", improveARelationship()))},
		},
	}
}
