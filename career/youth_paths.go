package career

// The five youth paths of pp. 68-74, one function each. Every one is
// nineteen rows for a 2d10 result of 2 to 20, and every one opens on result
// 2 with the same orphaning and reaches the Youth Life Events table at 10.

// youthPathOne is p. 68-69: no requirements, and the path every character
// can take.
func youthPathOne() []EventRow {
	return []EventRow{
		orphaned(),
		{
			Summary: "a year without adequate food or medication",
			Effects: []Effect{
				chr("STR", -1), chr("END", 1),
				unimplemented("a level in a Survival specialty available on the homeworld"),
			},
		},
		{
			Summary: "violence or chaos nearby, and the wariness it leaves",
			Effects: []Effect{skill("Recon"), chr("DEX", 1), chr("CHA", -1)},
		},
		{
			Summary: "the family collapses: debt, addiction, or imprisonment",
			Effects: []Effect{
				ratingOfRole("parent", false, -150),
				pick("choose what you grew accustomed to",
					opt("Streetwise", skillZero("Streetwise")),
					opt("Admin", skillZero("Admin"))),
				chr("EDU", -1),
			},
		},
		{
			Summary: "a parent who believes in harsh discipline",
			Effects: []Effect{
				ratingOfRole("parent", false, -50),
				pick("choose how you answered it",
					opt("Discipline", skill("Discipline")),
					opt("Melee (Unarmed)", skill("Melee", "Unarmed"))),
				chr("CHA", -1),
			},
		},
		{
			Summary: "isolation, and a feeling of not belonging",
			Effects: []Effect{
				rating(TargetAll, "", -25, ""),
				chr("CHA", -1),
				pick("choose what you compensated with",
					opt("Art", skill("Art", "Any")),
					opt("Electronics (Computers)", skill("Electronics", "Computers"))),
			},
		},
		{
			Summary: "a childhood injury or disease, and a long time in hospital",
			Effects: []Effect{
				pick("choose what it cost",
					opt("STR", chr("STR", -1)),
					opt("DEX", chr("DEX", -1)),
					opt("END", chr("END", -1))),
				pick("choose what the exposure taught you",
					opt("Admin", skillZero("Admin")),
					opt("Medic", skillZero("Medic"))),
			},
		},
		{
			Summary: "undue blame from authority figures outside the home",
			Effects: []Effect{pick("choose what chafing taught you",
				opt("Deception", skillZero("Deception", "Any")),
				opt("Streetwise", skillZero("Streetwise")))},
		},
		youthLifeEvent(),
		{
			Summary: "forced into labor at a young age",
			Effects: []Effect{
				pick("choose what the work built",
					opt("STR", chr("STR", 1)),
					opt("END", chr("END", 1))),
				pickSkill("Admin", "Broker", "Chef", "Electronics", "Mechanic", "Trade"),
				chr("EDU", -1),
			},
		},
		{
			Summary: "a close childhood friend",
			Effects: []Effect{
				relationshipAt(Ally, 1, 125),
				pick("choose what the friendship taught you",
					opt("Carouse", skillZero("Carouse")),
					opt("Persuade", skillZero("Persuade"))),
				chr("CHA", 1),
			},
		},
		{
			Summary: "a mentor",
			Effects: []Effect{
				mentorAt125(),
				pickSkill("Admin", "Art", "Broker", "Chef", "Etiquette", "Mechanic",
					"Science", "Trade"),
			},
		},
		{
			Summary: "the focus of attention at a local celebration",
			Effects: []Effect{chr("CHA", 1), skillZero("Carouse")},
		},
		{
			Summary: "a natural knack for something practical or mental",
			Effects: []Effect{
				pick("choose which",
					opt("DEX", chr("DEX", 1)),
					opt("INT", chr("INT", 1))),
				pickSkill("Admin", "Art", "Broker", "Chef", "Electronics", "Mechanic",
					"Science", "Trade"),
			},
		},
		{
			Summary: "a stable and supportive family",
			Effects: []Effect{rating(TargetFamily, "", 20, ""), chr("EDU", 1)},
		},
		{
			Summary: "work in school praised and encouraged",
			Effects: []Effect{
				pick("choose what the encouragement built",
					opt("EDU", chr("EDU", 1)),
					opt("CHA", chr("CHA", 1))),
				skillZero("Leadership"),
			},
		},
		{
			Summary: "a person or place that gives you a sense of belonging",
			Effects: []Effect{
				relationshipAt(Contact, 1, 90),
				anyTwoCharacteristics(),
				pickSkill("Admin", "Animals", "Athletics", "Art", "Broker", "Carouse",
					"Chef", "Etiquette", "Language", "Mechanic", "Science", "Trade"),
			},
		},
		{
			Summary: "a time of plenty, and living without stress or worry",
			Effects: []Effect{
				pick("choose what the comfort brought",
					opt("EDU", chr("EDU", 1)),
					opt("CHA", chr("CHA", 1))),
				chr("END", -1),
			},
		},
		{
			Summary: "deeply loved by your parents",
			Effects: []Effect{chr("EDU", 1), ratingOfRole("parent", true, 25)},
		},
	}
}

// anyTwoCharacteristics is "+1 to any two characteristics", which one result
// grants and nothing else in the book does. The two choices are made
// separately, so a character may raise the same one twice -- the page does
// not forbid it, and forbidding it would be a rule this reading invented.
func anyTwoCharacteristics() Effect {
	return pick("raise two characteristics, the first",
		opt("STR", chr("STR", 1), oneMoreCharacteristic()),
		opt("DEX", chr("DEX", 1), oneMoreCharacteristic()),
		opt("END", chr("END", 1), oneMoreCharacteristic()),
		opt("INT", chr("INT", 1), oneMoreCharacteristic()),
		opt("EDU", chr("EDU", 1), oneMoreCharacteristic()),
		opt("CHA", chr("CHA", 1), oneMoreCharacteristic()))
}

func oneMoreCharacteristic() Effect {
	return pick("and the second",
		opt("STR", chr("STR", 1)),
		opt("DEX", chr("DEX", 1)),
		opt("END", chr("END", 1)),
		opt("INT", chr("INT", 1)),
		opt("EDU", chr("EDU", 1)),
		opt("CHA", chr("CHA", 1)))
}
