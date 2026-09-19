package career

// Prisoner is the career of pp. 261-264: "You will only join this career if
// you have been sent here during character generation or within the course
// of a campaign."
//
// It is the one career a player cannot choose. Entry is always a table
// result -- Life Event 3 (p. 120), several of Vagabond's mishaps and events
// -- and the sentence carries a term count the event named. A mishap does
// not eject from it, which p. 263 says in print and in bold.
func Prisoner() Career {
	return Career{
		Name: "Prisoner",
		Cite: "pp. 261-264",
		EnlistmentNote: "no enlistment throw, and it cannot be chosen: entry is always " +
			"a table result, which names the number of terms served (p. 111)",
		MishapEjects: false,
		Assignments:  prisonerAssignments(),
		Tables:       prisonerTables(),
		Benefits:     prisonerBenefits(),
		Mishaps:      prisonerMishaps(),
		Events:       prisonerEvents(),
	}
}

func prisonerAssignments() []Assignment {
	return []Assignment{
		{
			Name:        "Prisoner",
			Description: "imprisoned after being convicted of a crime",
			Survival:    Check{Characteristic: "END", Number: 8},
			Advancement: Check{Characteristic: "STR", Number: 8},
			Skills: SkillTable{
				Kind: AssignmentSkills,
				Name: "Prisoner",
				Rows: [6]Effect{
					skill("Stealth"),
					skill("Interrogation", "Any"),
					skill("Carouse"),
					skill("Melee", "Any"),
					skill("Mechanic"),
					skill("Streetwise"),
				},
			},
			Ranks: [][]Effect{
				{skill("Melee", "Any"), stashItem("a stash")},
				{skill("Streetwise")},
				{chr("END", 1)},
				nil,
				{chr("END", 1)},
				{skill("Streetwise")},
				{skill("Jack of All Trades")},
			},
		},
	}
}

func prisonerTables() []SkillTable {
	return []SkillTable{
		{
			Kind: PersonalDevelopment,
			Name: "Personal Development",
			Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), skill("Athletics", "Any"),
			},
		},
		{
			Kind: ServiceSkills,
			Name: "Service Skills",
			Rows: [6]Effect{
				skill("Carouse"),
				skill("Deception", "Any"),
				skill("Melee", "Any"),
				skill("Persuade"),
				skill("Streetwise"),
				skill("Gambler"),
			},
		},
		{
			Kind:       AdvancedEducation,
			Name:       "Advanced Education",
			MinimumEDU: 8,
			Rows: [6]Effect{
				skill("Art", "Any"),
				skill("Advocate", "Legal"),
				skill("Electronics", "Any"),
				skill("Leadership"),
				chr("EDU", 1),
				skill("Jack of All Trades"),
			},
		},
	}
}

func prisonerBenefits() [7]BenefitRow {
	return [7]BenefitRow{
		{Cash: 0, Other: Effect{Detail: "nothing"}},
		{Cash: 0, Other: Effect{Detail: "nothing"}},
		{Cash: 0, Other: relationship(Contact, 1, "")},
		{Cash: 0, Other: relationship(Ally, 1, "")},
		{Cash: 0, Other: chr("STR", 1)},
		{Cash: 500, Other: chr("END", 1)},
		{Cash: 1000, Other: chr("CHA", 1)},
	}
}

func prisonerMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "toughs in your cell block pick you as their punching bag",
			Effects: []Effect{
				checkSkill("Melee", 8, nil, []Effect{injury(2)}),
				relationship(Enemy, 1, ""),
			},
		},
		{Summary: "you contract a disease", Effects: []Effect{chr("END", -2)}},
		{
			Summary: "aid workers turn out to be running a medical experiment on prisoners",
			Effects: []Effect{unimplemented(
				"roll 1d6: on 1, -2 END permanently; on 2-5 nothing; on 6, +1 STR and +1 END with an addiction")},
		},
		{
			Summary: "you become addicted to what is made in or smuggled into the prison",
			Effects: []Effect{unimplemented("an addiction, with no printed characteristic cost")},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "the guards clean your cell and find your stash",
			Effects: []Effect{
				stashItem(""),
				unimplemented("lose every benefit roll collected so far in this career"),
			},
		},
		{
			// ERRATA E-1: the book cites p. 136 in the failure branch.
			Summary: "a prison lord decides you will be his new servant",
			Effects: []Effect{
				pick("submit, fight, or buy your way out",
					opt("submit", chr("CHA", -2)),
					opt("fight", checkSkill("Melee", 8,
						[]Effect{chr("CHA", 1)},
						[]Effect{injury(1), chr("CHA", -3)})),
					opt("trade the stash", stashItem(""))),
				relationship(Enemy, 1, ""),
			},
		},
		{
			Summary: "the prison's atmosphere or building materials are unsafe",
			Effects: []Effect{chr("END", -2)},
		},
		{
			Summary: "a friend turns on you and the group takes your things",
			Effects: []Effect{relationship(Enemy, 1, ""), injury(1), stashItem("")},
		},
		{
			Summary: "a group of prisoners attempts to escape",
			Effects: []Effect{unimplemented(
				"join them (Stealth 8+ to escape and take a new homeworld, or a failed attempt adds a term), " +
					"stay silent (+2 to the next advancement roll and 1000 credits in the stash), " +
					"or tell the guards (one term off the sentence, and the escapers as Enemies)")},
		},
	}
}
