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
			Summary: "a death in your peer group",
			Effects: []Effect{unimplemented(
				"roll 1d6: on 1-3 lose an Ally, Contact, Rival or Enemy in that order; on 4-6 the reverse order")},
		},
		{
			Summary: "you have been betrayed",
			Effects: []Effect{unimplemented(
				"an existing Ally or Contact loses 1d6 x 20 Relationship Rating, or, with none, gain an Enemy at -110")},
		},
		{
			Summary: "a relationship changes",
			Effects: []Effect{unimplemented(
				"roll 1d6: on 1-3 an NPC gains 1d6 x 10 Relationship Rating, on 4-6 one loses that much")},
		},
		{
			Summary: "a new contact",
			Effects: []Effect{relationship(Contact, 1, "")},
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
			Effects: []Effect{unimplemented(
				"choose an existing Ally, or raise a Contact by 100 Relationship Rating to make them one; " +
					"with neither, gain an Ally at a Relationship Rating of 180")},
		},
		{
			Summary: "something wonderful",
			Effects: []Effect{pick("what came of it",
				opt("a new Ally at a Relationship Rating of 140", relationship(Ally, 1, "")),
				opt("an inheritance or prize", credits("5000")),
				opt("upgrade a relationship as in result 11",
					unimplemented("raise a Contact by 100 Relationship Rating to make them an Ally")),
			)},
		},
	}
}
