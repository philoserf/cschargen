package career

// teenagePathOne is pp. 76-78, for a homeworld settled a hundred standard
// years or more.
func teenagePathOne() []EventRow {
	return []EventRow{
		{
			Summary: "a personal scandal, deserved or not",
			Effects: []Effect{chr("CHA", -1), skill("Streetwise")},
		},
		{
			Summary: "family hardship, and maturing quickly for it",
			Effects: []Effect{chr("EDU", -1)},
		},
		{
			Summary: "pushed out of a circle of friends",
			Effects: []Effect{
				skill("Recon"), chr("CHA", -1),
				unimplemented("1d6-2 (minimum 1) non-family Contacts or Allies lose 50 " +
					"Relationship Rating, which may change what they are"),
			},
		},
		{
			Summary: "people who are not good for you, and the mark they leave",
			Effects: []Effect{
				pickSkill("Deception", "Streetwise"),
				throwModifier("enlistment in any criminal career", 2),
			},
		},
		{
			Summary: "school that does not come easy",
			Effects: []Effect{chr("EDU", -2)},
		},
		{
			Summary: "steady work after school, taken on out of expectation",
			Effects: []Effect{
				pickSkill("Admin", "Animals", "Broker", "Chef", "Electronics",
					"Instruction", "Mechanic", "Medic", "Science", "Trade"),
				chr("EDU", -1),
			},
		},
		localRivalry(),
		aGoodFriend(),
		teenageLifeEvent(),
		{
			Summary: "something civic or social, and a sense of belonging",
			Effects: []Effect{pickSkill("Admin", "Carouse", "Etiquette")},
		},
		{
			Summary: "an adult who takes the time to teach you something useful",
			Effects: []Effect{
				relationshipAt(Ally, 1, 100),
				pickSkill("Admin", "Advocate", "Animals", "Athletics", "Art", "Broker",
					"Carouse", "Chef", "Drive", "Electronics", "Etiquette", "Instruction",
					"Mechanic", "Medic", "Melee", "Persuade", "Pilot", "Science",
					"Seafarer", "Tactics", "Trade"),
			},
		},
		competence(),
		success(),
		broadeningHorizons(),
		opportunity(),
		trustedRole(),
		strongSocialFoundation(),
		{
			Summary: "a sponsorship or scholarship to a local university",
			Effects: []Effect{
				relationshipAt(Contact, 1, 45),
				throwModifier("admission to Undergraduate College or the Military Academy", 2),
				chr("EDU", 2),
			},
		},
		fortunateYouth(),
	}
}

// teenagePathTwo is pp. 78-79, for a homeworld settled less than a hundred
// standard years.
func teenagePathTwo() []EventRow {
	return []EventRow{
		{
			Summary: "a serious crisis in the settlement, and a world that feels fragile",
			Effects: []Effect{
				chr("END", 1),
				pick("choose what it cost",
					opt("EDU", chr("EDU", -1)),
					opt("CHA", chr("CHA", -1))),
				unimplemented("a level in a Survival specialty used on the homeworld"),
			},
		},
		{
			Summary: "a death in the family",
			Effects: []Effect{loseRelative()},
		},
		{
			Summary: "a colony divided by faction, labor, politics or war",
			Effects: []Effect{
				pickSkill("Advocate", "Broker", "Deception", "Diplomat", "Discipline",
					"Draw", "Gun Combat", "Leadership", "Medic", "Melee", "Recon",
					"Stealth", "Streetwise", "Tactics"),
				chr("EDU", -2),
			},
		},
		{
			Summary: "hazardous labor, because the colony needs all hands",
			Effects: []Effect{chr("END", 1), chr("EDU", -1)},
		},
		{
			Summary: "schooling interrupted by shortages and labor demands",
			Effects: []Effect{chr("INT", 1), chr("EDU", -1)},
		},
		{
			Summary: "taken on as a helper by an older worker",
			Effects: []Effect{
				relationshipAt(Contact, 1, 90),
				pick("choose what they taught you",
					opt("Admin", skill("Admin")),
					opt("Animals (Farming)", skill("Animals", "Farming")),
					opt("Broker", skill("Broker")),
					opt("Chef", skill("Chef")),
					opt("Electronics", skill("Electronics", "Any")),
					opt("Engineer", skill("Engineer", "Life Support", "Power")),
					opt("Mechanic", skill("Mechanic")),
					opt("Trade", skill("Trade", "Any"))),
				chr("EDU", -1),
			},
		},
		localRivalry(),
		aGoodFriend(),
		teenageLifeEvent(),
		{
			Summary: "close bonds with your fellow settlers",
			Effects: []Effect{pick("choose what the bonds became",
				opt("two contacts", relationshipAt(Contact, 2, 70)),
				opt("one ally", relationshipAt(Ally, 1, 125)))},
		},
		{
			Summary: "offworlders, and a first look at how things work elsewhere",
			Effects: []Effect{pickSkill("Broker", "Carouse", "Streetwise")},
		},
		competence(),
		success(),
		{
			Summary: "experience with the unsettled regions of your world",
			Effects: []Effect{
				pick("choose what the frontier built",
					opt("END", chr("END", 1)),
					opt("CHA", chr("CHA", 1))),
				unimplemented("a level in a Survival specialty from the homeworld's " +
					"background skills"),
			},
		},
		opportunity(),
		trustedRole(),
		strongSocialFoundation(),
		broadeningHorizons(),
		fortunateYouth(),
	}
}
