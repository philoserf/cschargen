package career

// exoticEvents is the d66 table of pp. 192-193.
func exoticEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "keeping track of your clients",
			Effects: []Effect{skill("Admin")},
		},
		13: {
			Summary: "performing for the client, as part of the entertainment",
			Effects: []Effect{skill("Art", "Any")},
		},
		14: {
			Summary: "a client who skips out without paying",
			Effects: []Effect{pick("let it slide, or go after them",
				opt("let it slide", relationship(Ally, 1, "")),
				opt("Advocate", exoticCollect("Advocate")),
				opt("Gun Combat", exoticCollect("Gun Combat")),
				opt("Melee", exoticCollect("Melee")),
				opt("Streetwise", exoticCollect("Streetwise")))},
		},
		15: {
			Summary: "a lot of influential people, dealt with in secret",
			Effects: []Effect{skill("Diplomat")},
		},
		16: {
			Summary: "a gambling group joined",
			Effects: []Effect{gamblingCircle("Persuade")},
		},
		21: {
			Summary: "accused of engaging in illegal acts",
			Effects: []Effect{checkSkill("Advocate", 8,
				[]Effect{skill("Advocate", "Legal"), throwModifier("next advancement roll", 2)},
				[]Effect{transfer("Prisoner", "Prisoner", 1)})},
		},
		22: {
			Summary: "watching out for your own money, because nobody else will",
			Effects: []Effect{skill("Broker")},
		},
		23: {Summary: "a fun person to have around", Effects: []Effect{skill("Carouse")}},
		24: {
			Summary: "your own transportation",
			Effects: []Effect{pickSkill("Drive", "Flyer")},
		},
		25: {
			Summary: "self-defense, which this career requires",
			Effects: []Effect{pickSkill("Gun Combat", "Melee")},
		},
		26: {
			Summary: "a wide variety of clients",
			Effects: []Effect{skill("Language", "Any")},
		},
		41: {
			Summary: "the laws that govern your profession",
			Effects: []Effect{skill("Advocate", "Legal", "Politics")},
		},
		42: {
			Summary: "a wealthy regular who is careless with their money",
			Effects: []Effect{pick("steal from them, or do not",
				opt("steal", checkSkill("Deception", 8,
					[]Effect{
						relationship(Enemy, 1, ""),
						cashRolls(4),
					},
					[]Effect{pick("the client attacks: fight back with what you have",
						opt("Gun Combat", checkSkill("Gun Combat", 8,
							exoticFoughtOff(), []Effect{injury(2)})),
						opt("Melee", checkSkill("Melee", 8,
							exoticFoughtOff(), []Effect{injury(2)})))})),
				opt("leave their money alone", throwModifier("next advancement roll", 2)))},
		},
		43: {
			Summary: "a client who loves to be at sea, and teaches you to sail",
			Effects: []Effect{skill("Seafarer", "Sail")},
		},
		44: {
			Summary: "knowing what the client really wants",
			Effects: []Effect{skill("Interrogation", "Questioning")},
		},
		45: {
			Summary: "a course in first aid",
			Effects: []Effect{skill("Medic", "First Aid")},
		},
		46: {
			Summary: "a client attacks you",
			Effects: []Effect{checkSkill("Melee", 8,
				[]Effect{relationship(Enemy, 1, "")},
				[]Effect{injury(1)})},
		},
		51: {
			Summary: "a client who farms, and a thing or two about agricultural life",
			Effects: []Effect{skill("Animals", "Any")},
		},
		52: {
			Summary: "experience on the shady side of life",
			Effects: []Effect{skill("Deception", "Any")},
		},
		53: {
			Summary: "a fellow worker who teaches you to broadcast to certain clients",
			Effects: []Effect{skill("Electronics", "Comms", "Computers")},
		},
		54: {
			Summary: "chosen for advanced training",
			Effects: []Effect{rollTable(AdvancedEducation)},
		},
		55: {
			Summary: "frequently in space",
			Effects: []Effect{pick("choose Suit (Vacc Suit) or Survival (Freefall)",
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
				opt("Survival (Freefall)", skill("Survival", "Freefall")))},
		},
		56: {
			Summary: "telling which client is benign and which is a danger",
			Effects: []Effect{skill("Recon")},
		},
		61: {
			Summary: "a body kept in peak condition",
			Effects: []Effect{pick("choose what the gym built",
				opt("STR", chr("STR", 1)),
				opt("DEX", chr("DEX", 1)),
				opt("END", chr("END", 1)),
				opt("Athletics", skill("Athletics", "Any")))},
		},
		62: {
			Summary: "an influential client who wants an argument with a friend settled",
			Effects: []Effect{checkSkill("Diplomat", 8,
				[]Effect{cashRolls(3)},
				[]Effect{chr("CHA", -1)})},
		},
		63: {
			// The table prints "Engineering (Any)"; the skill list (p. 304)
			// has Engineer. See ERRATA E-14.
			Summary: "a lot of time on a starship, and a crew happy to teach their jobs",
			Effects: []Effect{pick("choose the job you learned",
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Gunner (Turrets)", skill("Gunner", "Turrets")),
				opt("Pilot", skill("Pilot", "Any")))},
		},
		64: {
			Summary: "a number of odd jobs",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		65: {
			Summary: "extremely talented in your field",
			Effects: []Effect{talentedInYourField()},
		},
		66: {
			Summary: "excellent work",
			Effects: []Effect{advance(), benefitRolls(0, 1, ScopeCareer)},
		},
	}

	lifeEventRows(table)

	return table
}

// exoticCollect is event 14's second half: whichever way the debt is
// collected, the skill used is the one that improves.
func exoticCollect(name string) Effect {
	return checkSkill(name, 8, []Effect{
		skill(name, "Any"),
		benefitRolls(2, 0, ScopeBatch),
		relationship(Enemy, 1, ""),
	}, nil)
}

// exoticFoughtOff is what surviving event 42's angry client is worth.
func exoticFoughtOff() []Effect {
	return []Effect{relationship(Enemy, 1, ""), cashRolls(2)}
}
