package career

// MilitaryAcademy is the academy of pp. 91-96.
//
// It is the only institution that commits a graduate to what comes next:
// "After attending a military academy, the character must enter a military
// career in Step 9 ... the character must spend a minimum of two terms in
// that career" (p. 94). Success enters as an officer at rank 0, success
// with honours at rank 1.
func MilitaryAcademy() Institution {
	return Institution{
		Name: "Military Academy",
		Cite: "pp. 91-96",
		Prerequisites: []Check{
			{Characteristic: "EDU", Number: 6},
			{Characteristic: "END", Number: 8},
		},
		Admission:  admission(9),
		Success:    successThrow(8),
		Honors:     honorsThrow(10),
		Skills:     academySkills(),
		Failure:    academyFailure(),
		Events:     academyEvents(),
		LifeEvents: academyLifeEvents(),
		Note: "a graduate must enter a military career at Step 9 as an officer -- rank 0, " +
			"or rank 1 with honours -- and serve at least two terms in it (p. 94)",
	}
}

// academySkills is the list of p. 92.
func academySkills() []Effect {
	return []Effect{
		skill("Admin"),
		skill("Advocate", "Any"),
		skill("Athletics", "Any"),
		skill("Art", "Any"),
		skill("Astrogation"),
		skill("Diplomat"),
		skill("Discipline"),
		skill("Drive", "Any"),
		skill("Electronics", "Any"),
		skill("Engineer", "Any"),
		skill("Explosives"),
		skill("Flyer", "Any"),
		skill("Gun Combat", "Any"),
		skill("Gunner", "Any"),
		skill("Heavy Weapons", "Any"),
		skill("Interrogation", "Any"),
		skill("Leadership"),
		skill("Mechanic"),
		skill("Melee", "Any"),
		skill("Navigation"),
		skill("Persuade"),
		skill("Pilot", "Any"),
		skill("Recon"),
		skill("Science", "Any"),
		skill("Seafarer", "Any"),
		skill("Stealth"),
		skill("Suit", "Any"),
		skill("Survival", "Any"),
		skill("Tactics", "Any"),
	}
}

// academyFailure is the 2d6 table of p. 93. A character who reaches it may
// never attempt a military academy again, which the engine records on the
// character rather than on the table.
func academyFailure() MishapTable {
	return MishapTable{
		{
			Summary: "a court martial, and dismissal in disgrace",
			Effects: []Effect{relationshipAt(Enemy, 1, -170)},
		},
		{Summary: "a severe injury during training", Effects: []Effect{injury(1)}},
		{
			Summary: "dismissed for insubordination",
			Effects: []Effect{throwModifier("enlistment in any military career", -4)},
		},
		{
			Summary: "caught in academy politics, and made the scapegoat",
			Effects: []Effect{relationshipAt(Enemy, 1, -120), relationshipAt(Contact, 1, 60)},
		},
		{
			Summary: "the tactical and technical demands proved too much",
			Effects: []Effect{chr("EDU", -2)},
		},
		{Summary: "a quiet washout, without distinction", Effects: nil},
		{
			Summary: "a crisis at home that outweighed duty to the service",
			Effects: []Effect{loseRelative(), ratingOfRole("parent", false, 50)},
		},
		{
			Summary: "a breakdown under constant pressure",
			Effects: []Effect{checkChr("END", 8,
				[]Effect{chr("CHA", -1)},
				[]Effect{chr("CHA", -2)})},
		},
		{
			Summary: "not an officer, but the military kept you",
			Effects: []Effect{unimplemented(
				"enlist automatically in a military career at the lowest enlisted rank")},
		},
		{
			Summary: "war, and cadets rushed into service before finishing",
			Effects: []Effect{unimplemented(
				"enlist automatically in a military career at enlisted rank 1")},
		},
		{
			Summary: "expelled for breaking the rules for the right reason",
			Effects: []Effect{relationshipAt(Ally, 1, 120), relationshipAt(Enemy, 1, -120)},
		},
	}
}

// academyEvents is the 2d10 table of pp. 94-95.
func academyEvents() []EventRow {
	return []EventRow{
		{
			Summary: "a training accident you stayed in the academy despite",
			Effects: []Effect{injury(1)},
		},
		{
			Summary: "a formal punishment that stains any future military career",
			Effects: []Effect{
				relationshipAt(Rival, 1, -35),
				throwModifier("first two advancement rolls in a military career", -2),
			},
		},
		{
			Summary: "an important course failed, repeated, and passed",
			Effects: []Effect{chr("EDU", -2)},
		},
		{
			Summary: "an instructor who sees you as a problem and works it out of you",
			Effects: []Effect{chr("END", 1), relationshipAt(Rival, 1, -35)},
		},
		{
			Summary: "a rivalry with another cadet that becomes personal",
			Effects: []Effect{relationshipAt(Rival, 1, -70)},
		},
		{
			Summary: "unpleasant duties, and the humility they teach",
			Effects: []Effect{skill("Discipline"), chr("CHA", -1)},
		},
		{
			Summary: "a friendship with a fellow cadet",
			Effects: []Effect{relationshipAt(Contact, 1, 75)},
		},
		{
			Summary: "extended field training",
			Effects: []Effect{pickSkill("Survival", "Recon")},
		},
		{
			Summary: "an academy life event",
			Effects: []Effect{{
				Kind:   EffectAcademyLifeEvent,
				Detail: "roll on the Military Academy Life Events table",
			}},
		},
		{
			Summary: "an instructor who notices competence",
			Effects: []Effect{relationshipAt(Ally, 1, 110)},
		},
		{
			Summary: "specialized training for the career you will enter",
			Effects: []Effect{unimplemented(
				"choose the military career you will enter and roll once on its " +
					"Service Skills table")},
		},
		{
			Summary: "placed in charge of other cadets",
			Effects: []Effect{skill("Leadership"), chr("CHA", 1)},
		},
		{
			Summary: "distinguished performance in an evaluation",
			Effects: []Effect{
				skill("Discipline"),
				pick("choose what the reputation built",
					opt("EDU", chr("EDU", 1)),
					opt("CHA", chr("CHA", 1))),
			},
		},
		{
			Summary: "sharp thinking under pressure during exercises",
			Effects: []Effect{skill("Tactics", "Military", "Naval"), chr("EDU", 1)},
		},
		{
			Summary: "identified as suited to a specialized service role",
			Effects: []Effect{pickSkill("Electronics", "Engineer")},
		},
		{
			Summary: "a formal award for service, courage or excellence",
			Effects: []Effect{
				relationshipAt(Ally, 1, 110),
				pick("choose what the citation built",
					opt("CHA", chr("CHA", 1)),
					opt("EDU", chr("EDU", 1))),
			},
		},
		{
			Summary: "a network of comrades to leave the academy with",
			Effects: []Effect{relationshipAt(Contact, 2, 60)},
		},
		{
			Summary: "a senior officer who becomes a patron",
			Effects: []Effect{
				relationshipAt(Ally, 1, 130),
				throwModifier("first two survival and advancement rolls in a military career", 2),
			},
		},
		{
			Summary: "among the finest of your class",
			Effects: []Effect{
				relationshipAt(Ally, 1, 110),
				relationshipAt(Contact, 1, 60),
				{Kind: EffectHonors, Detail: "achieve honours, if the throw did not"},
				chr("CHA", 1),
				pickSkill("Discipline", "Diplomat", "Admin", "Leadership", "Recon"),
			},
		},
	}
}

// academyLifeEvents is the d6 table of p. 96. It is the collegiate table
// with one result changed: where p. 91's result 2 is a death, this one is a
// battle the cadets were called into.
func academyLifeEvents() []EventRow {
	table := collegiateLifeEvents()

	table[1] = EventRow{
		Summary: "a military action the academy was called into",
		Effects: []Effect{pick("answer the call with what you have",
			opt("Gun Combat", checkSkill("Gun Combat", 8,
				[]Effect{throwModifier("first advancement roll in a military career", 2)},
				[]Effect{injury(1)})),
			opt("Tactics", checkSkill("Tactics", 8,
				[]Effect{throwModifier("first advancement roll in a military career", 2)},
				[]Effect{injury(1)})))},
	}

	return table
}
