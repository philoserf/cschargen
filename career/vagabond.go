package career

// Vagabond is the career of pp. 296-300: "There is no Enlistment roll to
// join the Vagabond career. Chances are high that you are being forced into
// this career by an event in character generation or by failing to enlist
// in other careers."
//
// It is the floor of the whole system. Three consecutive failed enlistments
// send a character here (p. 110), several careers' mishaps name it, and the
// aging crisis leaves it as one of three options (p. 123). A mishap does
// not eject from it, which p. 298 says in print.
func Vagabond() Career {
	return Career{
		Name: "Vagabond",
		Cite: "pp. 296-300",
		EnlistmentNote: "no enlistment throw: entered by choice, by a table result, " +
			"or by failing three enlistments in a row (pp. 110-111)",
		MishapEjects: false,
		Assignments:  vagabondAssignments(),
		Tables:       vagabondTables(),
		Benefits:     vagabondBenefits(),
		Mishaps:      vagabondMishaps(),
		Events:       vagabondEvents(),
	}
}

func vagabondAssignments() []Assignment {
	return []Assignment{
		{
			Name:        "Destitute",
			Description: "living off the surroundings, rough on a planet, in an orbital city or on a starport",
			Survival:    Check{Characteristic: "END", Number: 8},
			Advancement: Check{Characteristic: "INT", Number: 8},
			Skills: SkillTable{
				Kind: AssignmentSkills,
				Name: "Destitute",
				Rows: [6]Effect{
					skill("Deception", "Any"),
					skill("Navigation"),
					skill("Survival", "Any"),
					skill("Streetwise"),
					skill("Gun Combat", "Any"),
					skill("Jack of All Trades"),
				},
			},
			Ranks: [7][]Effect{
				{skill("Survival", "Any")},
				{stashItem("a stash")},
				nil,
				{skill("Recon")},
				nil,
				{skill("Persuade")},
				{skill("Carouse")},
			},
		},
		{
			Name:        "Transient",
			Description: "moving through space to see the greatness of it, on charity and odd jobs",
			Survival:    Check{Characteristic: "END", Number: 8},
			Advancement: Check{Characteristic: "INT", Number: 8},
			Skills: SkillTable{
				Kind: AssignmentSkills,
				Name: "Transient",
				Rows: [6]Effect{
					skill("Gun Combat", "Any"),
					skill("Deception", "Any"),
					skill("Survival", "Any"),
					skill("Suit", "Vacc Suit"),
					skill("Mechanic"),
					skill("Jack of All Trades"),
				},
			},
			Ranks: [7][]Effect{
				{skill("Suit", "Vacc Suit")},
				{stashItem("a stash")},
				nil,
				{skill("Persuade")},
				nil,
				{skill("Recon")},
				{skill("Carouse")},
			},
		},
	}
}

// vagabondTables has no Advanced Education row, and that absence is the
// data format's first real test: p. 117 says a career has "at least four
// skill tables" and names Advanced Education among the usual ones, but
// p. 296 prints only two career-wide tables.
func vagabondTables() []SkillTable {
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
				skill("Stealth"),
				skill("Recon"),
				skill("Survival", "Any"),
				skill("Streetwise"),
				skill("Melee", "Any"),
				skill("Persuade"),
			},
		},
	}
}

func vagabondBenefits() [7]BenefitRow {
	return [7]BenefitRow{
		{Cash: 0, Other: Effect{Detail: "nothing"}},
		{Cash: 0, Other: Effect{Detail: "nothing"}},
		{Cash: 0, Other: chr("END", 1)},
		{Cash: 0, Other: chr("STR", 1)},
		{Cash: 500, Other: relationship(Contact, 1, "")},
		{Cash: 1000, Other: relationship(Ally, 1, "")},
		{Cash: 2000, Other: relationship(Ally, 2, "")},
	}
}

func vagabondMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "local thugs pick you as their punching bag for the evening",
			Effects: []Effect{
				pick("fight back with what you have",
					opt("Melee", checkSkill("Melee", 8,
						[]Effect{injury(1)},
						[]Effect{injury(2), stashItem(""), relationship(Enemy, 1, "")})),
					opt("Gun Combat", checkSkill("Gun Combat", 8,
						[]Effect{injury(1)},
						[]Effect{injury(2), stashItem(""), relationship(Enemy, 1, "")}))),
			},
		},
		{Summary: "you contract a deadly disease and recover badly", Effects: []Effect{chr("END", -2)}},
		{
			Summary: "an aid station turns out to be a secret drug trial on the expendable",
			Effects: []Effect{unimplemented(
				"roll 1d6: on 1, lose 2 from two physical characteristics; on 2-5 nothing; on 6, +1 END and an addiction")},
		},
		{
			Summary: "you become addicted to drugs or alcohol",
			Effects: []Effect{chr("CHA", -1), chr("END", -2)},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "a politician decides to clean up the area",
			Effects: []Effect{unimplemented(
				"roll 1d6: on 1, a term in the Prisoner career; on 2-5, deportation to a new " +
					"homeworld; on 6, wounded escaping a culling")},
		},
		{
			Summary: "you fail to prepare for the weather and suffer exposure",
			Effects: []Effect{chr("END", -1)},
		},
		{
			Summary: "the shelter you chose turns out to be unsafe",
			Effects: []Effect{injury(1)},
		},
		{
			Summary: "swept up into a mental facility for the term, and marked insane after",
			Effects: []Effect{throwModifier("next enlistment attempt", -2)},
		},
		{
			Summary: "accused of a crime that gets no investigation",
			Effects: []Effect{unimplemented(
				"roll 1d6: on 1, two terms in the Prisoner career; on 2-5, one term; on 6, you escape it")},
		},
	}
}
