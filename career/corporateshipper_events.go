package career

// shipperEvents is the d66 table of pp. 179-180.
func shipperEvents() EventTable {
	strandedByPirates := []Effect{pick("what the walk to civilization taught you",
		opt("Animals", skill("Animals", "Any")),
		opt("Navigation", skill("Navigation")),
		opt("Survival", skill("Survival", "Any")))}

	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "pirates attack",
			Effects: []Effect{pick("fight them off with what you have",
				opt("Gun Combat", checkSkill("Gun Combat", 8, nil, strandedByPirates)),
				opt("Melee", checkSkill("Melee", 8, nil, strandedByPirates)))},
		},
		13: {
			Summary: "a barfight",
			Effects: []Effect{checkSkill("Melee", 8,
				[]Effect{relationship(Contact, 1, "")}, []Effect{injury(1)})},
		},
		14: {
			Summary: "the shadier side of the transport business",
			Effects: []Effect{pick("what it taught you",
				opt("Deception", skill("Deception", "Any")),
				opt("Streetwise", skill("Streetwise")))},
		},
		15: {
			Summary: "a distress call from a damaged ship in the outer system",
			Effects: []Effect{pick("answer it, report it, or ignore it",
				opt("ignore it", unimplemented("lose one rank as word gets out")),
				opt("report it and do nothing", chr("CHA", -1)),
				opt("help", unimplemented(
					"roll 1d6: on 1-3 it is a pirate ship and the boarding fight follows; "+
						"on 4-5 a grateful captain becomes a Contact; on 6 an Ally and 5,000 credits")))},
		},
		16: {Summary: "a great time with your coworkers", Effects: []Effect{skill("Carouse")}},
		21: {
			Summary: "knowing the laws and politics of your destination",
			Effects: []Effect{pickSkill("Advocate")},
		},
		22: {
			Summary: "emergency repairs",
			Effects: []Effect{pick("fix it with what you know",
				opt("Engineer", checkSkill("Engineer", 8,
					[]Effect{skill("Engineer", "Any")},
					[]Effect{benefitRolls(-1, 0, ScopeBatch)})),
				opt("Mechanic", checkSkill("Mechanic", 8,
					[]Effect{skill("Mechanic")},
					[]Effect{benefitRolls(-1, 0, ScopeBatch)})))},
		},
		23: {
			Summary: "you spot a crewmate doing something dangerous or illegal",
			Effects: []Effect{spottedSomethingWrong(rollTable(ServiceSkills))},
		},
		24: {
			Summary: "an argument between officers, or in middle management",
			Effects: []Effect{checkSkill("Diplomat", 8,
				[]Effect{unimplemented("three company shares")},
				[]Effect{throwModifier("next advancement roll", -2)})},
		},
		25: {
			Summary: "a gambling circle aboard ship or at a common port",
			Effects: []Effect{checkSkill("Gambler", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch)},
				[]Effect{benefitRolls(-3, 0, ScopeBatch)})},
		},
		26: {
			Summary: "the holographic gun range in the cargo bay",
			Effects: []Effect{skill("Gun Combat", "Any")},
		},
		41: {Summary: "the corporation wants everything on record", Effects: []Effect{skill("Admin")}},
		42: {
			Summary: "a fondness for vehicles",
			Effects: []Effect{pickSkill("Drive", "Flyer", "Seafarer")},
		},
		43: {Summary: "engineering still moonshine", Effects: []Effect{skill("Carouse")}},
		44: {
			Summary: "some days you are the only one aboard who does their job",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		45: {Summary: "a course in first aid", Effects: []Effect{skill("Medic", "First Aid")}},
		46: {
			Summary: "a great many worlds visited",
			Effects: []Effect{pickSkill("Language", "Survival")},
		},
		51: {Summary: "the hobby of riding", Effects: []Effect{skill("Animals", "Riding")}},
		52: {
			Summary: "the crew enjoy making you the shuttle pilot",
			Effects: []Effect{skill("Pilot", "Small Craft")},
		},
		53: {Summary: "fixing your own equipment", Effects: []Effect{skill("Mechanic")}},
		54: {
			Summary: "a great deal of time aboard ship",
			Effects: []Effect{pick("what the time taught you",
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
				opt("Survival (Freefall)", skill("Survival", "Freefall")))},
		},
		55: {
			Summary: "downtime spent on tabletop games",
			Effects: []Effect{pick("what the games taught you",
				opt("Deception (Lie)", skill("Deception", "Lie")),
				opt("Science", skill("Science", "History", "Psychology")),
				opt("Tactics", skill("Tactics", "Any")))},
		},
		56: {
			Summary: "learning what lurks in the shadows",
			Effects: []Effect{pickSkill("Recon", "Stealth")},
		},
		61: {
			Summary: "long stretches in Zimmspace, and time to fill",
			Effects: []Effect{pickSkill("Athletics", "Art", "Science")},
		},
		62: {
			Summary: "short-crewed, and working outside your specialty",
			Effects: []Effect{rollOtherAssignment()},
		},
		63: {Summary: "a profitable voyage you helped make", Effects: []Effect{benefitRolls(2, 0, ScopeBatch)}},
		64: {Summary: "downtime spent working on yourself", Effects: []Effect{rollTable(PersonalDevelopment)}},
		65: {
			Summary: "extremely talented in your field",
			Effects: []Effect{pick("what the talent earned",
				opt("a level in a skill you already have",
					unimplemented("raise a skill the character already holds")),
				opt("a promotion", Effect{Kind: EffectRank, Detail: "gain a rank"}))},
		},
		66: {
			Summary: "excellent work",
			Effects: []Effect{advance(), benefitRolls(0, 1, ScopeCareer)},
		},
	}

	lifeEventRows(table)

	return table
}
