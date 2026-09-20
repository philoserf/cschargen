package career

// Injury is the table of p. 119, indexed 0-5 for a 1d6 result of 1-6.
//
// Three references in the book send a character to "the Injury table
// (p.136)" -- Colonist mishap 9, and two of Prisoner's events. p. 136
// carries no table and every other reference says p. 119, so all of them
// resolve here (ERRATA E-1).
//
// The recovery clause is not in this table because it is not a result: a
// lost characteristic can be bought back for 1d6 x 10,000 credits "if the
// character can afford this cost ... at the time of the injury", which for
// a character who has not yet mustered out means it cannot (ERRATA E-4).
func Injury() [6]MishapRow {
	return [6]MishapRow{
		{
			Summary: "nearly killed",
			Effects: []Effect{unimplemented(
				"reduce one physical characteristic by 1d6, and the other two by 2 each (or one of them by 4)")},
		},
		{
			Summary: "severely injured",
			Effects: []Effect{unimplemented("reduce one physical characteristic by 1d6")},
		},
		{
			Summary: "heavily injured: an eye or a limb is gone",
			Effects: []Effect{pick("what was lost",
				opt("an arm or hand", chr("STR", -2)),
				opt("an eye", chr("DEX", -2)),
			)},
		},
		{
			Summary: "heavily injured",
			Effects: []Effect{pick("which physical characteristic took it",
				opt("STR", chr("STR", -2)),
				opt("DEX", chr("DEX", -2)),
				opt("END", chr("END", -2)),
			)},
		},
		{
			Summary: "injured",
			Effects: []Effect{pick("which physical characteristic took it",
				opt("STR", chr("STR", -1)),
				opt("DEX", chr("DEX", -1)),
				opt("END", chr("END", -1)),
			)},
		},
		{Summary: "lightly injured, with no permanent effect", Effects: nil},
	}
}

// LifeEvents is the table of p. 120, indexed 0-10 for a 2d6 result of 2-12.
// Every career's d66 reaches it at 31-36.
//
// Four of its eleven results write to a Relationship Rating, which is why
// an Ally or Contact carries one from the first term although generating an
// NPC around it is out of scope (ERRATA E-5).
func LifeEvents() MishapTable {
	return MishapTable{
		{Summary: "severe sickness or injury", Effects: []Effect{injury(1)}},
		{
			Summary: "involved in or accused of a crime",
			Effects: []Effect{pick("take the loss or take the sentence",
				opt("lose two of this career's benefit rolls", benefitRolls(-2, 0, ScopeBatch)),
				opt("a term in prison", transfer("Prisoner", "Prisoner", 1)),
			)},
		},
		{
			// The 1d6 decides only which end of the list to start from,
			// so it is a choice between two orders rather than a roll the
			// engine has to make.
			Summary: "a death in your peer group",
			Effects: []Effect{pick("which of them it was",
				opt("the closest first", loseTie(Ally, Contact, Rival, Enemy)),
				opt("the most distant first", loseTie(Enemy, Rival, Contact, Ally)))},
		},
		{
			Summary: "you have been betrayed",
			Effects: []Effect{
				rating(TargetOne, "", -1, "1d6x20"),
				unimplemented("with no Ally or Contact to lose, gain an Enemy at -110 instead"),
			},
		},
		{
			// Again the 1d6 chooses a direction rather than a quantity.
			Summary: "a relationship changes",
			Effects: []Effect{pick("which way it went",
				opt("closer", rating(TargetOne, "", 1, "1d6x10")),
				opt("further away", rating(TargetOne, "", -1, "1d6x10")))},
		},
		{
			Summary: "a new contact",
			Effects: []Effect{relationshipAt(Contact, 1, 40)},
		},
		{
			Summary: "dissatisfaction with your career",
			Effects: []Effect{
				throwModifier("next advancement roll", -2),
				throwModifier("next enlistment attempt", 2),
			},
		},
		{
			Summary: "a religious conversion",
			Effects: []Effect{unimplemented(
				"choose or invent a religion; on a 1d6 roll of 6 you become deeply involved and gain Science (Philosophy) 1")},
		},
		{
			Summary: "a minor life goal achieved",
			Effects: []Effect{benefitRolls(0, 2, ScopeBatch)},
		},
		{
			Summary: "a new romantic relationship",
			Effects: []Effect{pick("who it turned out to be",
				opt("a Contact who becomes more", rating(TargetOne, Contact, 100, "")),
				opt("somebody new", relationshipAt(Ally, 1, 180)))},
		},
		{
			Summary: "something wonderful",
			Effects: []Effect{pick("what came of it",
				opt("a new Ally", relationshipAt(Ally, 1, 140)),
				opt("an inheritance or prize", credits("5000")),
				opt("a Contact who becomes more", rating(TargetOne, Contact, 100, "")),
			)},
		},
	}
}

// MilitaryEvents is the table of p. 121, indexed 0-10 for a 2d6 result of
// 2-12. Several careers' d66 tables route to it.
//
// The book prints this table twice. p. 241's "Naval Events Table" is the
// same twelve results with naval wording -- "one of your crewmates" for
// "another force member", "shore leave" for "leave" -- and identical
// mechanics throughout. One table, transcribed once.
func MilitaryEvents() MishapTable {
	return MishapTable{
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "selected for special training",
			Effects: []Effect{rollTable(AdvancedEducation)},
		},
		{
			Summary: "a special bond with another member of the service",
			Effects: []Effect{relationship(Ally, 1, "")},
		},
		{
			Summary: "stuck at the base, with time to study",
			Effects: []Effect{checkChr("EDU", 8,
				[]Effect{pick("what the studying gave you",
					opt("Art", skill("Art", "Any")),
					opt("Science", skill("Science", "Any")))},
				nil)},
		},
		{
			Summary: "the government pays a military bonus",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch)},
		},
		{
			Summary: "good friends made in the service",
			Effects: []Effect{relationship(Contact, 1, "")},
		},
		{
			Summary: "a great time on leave",
			Effects: []Effect{pick("what the leave taught you",
				opt("Carouse", skill("Carouse")),
				opt("Gambler", skill("Gambler")),
				opt("Streetwise", skill("Streetwise")))},
		},
		{
			Summary: "you have greatly angered a superior",
			Effects: []Effect{
				relationship(Enemy, 1, ""),
				throwModifier("next advancement roll", -2),
			},
		},
		{
			// The branch turns on whether the character is commissioned,
			// which the engine knows and the table cannot express: the
			// commission effect is a no-op for an officer, and the officer
			// branch is the alternative.
			Summary: "time spent studying for a promotion",
			Effects: []Effect{pick("what the study was for",
				opt("a commission, if still enlisted", commission(2)),
				opt("officer skills, if already commissioned",
					rollTable(OfficerSkills),
					throwModifier("next advancement roll", 2)))},
		},
		{
			Summary: "a rivalry with someone in your unit",
			Effects: []Effect{relationship(Rival, 1, "")},
		},
		{
			Summary: "a special commendation for actions above and beyond the call of duty",
			Effects: []Effect{throwModifier("next advancement roll", 2)},
		},
	}
}
