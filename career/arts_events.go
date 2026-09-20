package career

// artsEvents is the d66 table of pp. 157-159.
func artsEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "a great deal of time spent learning copyright and intellectual property",
			Effects: []Effect{skill("Advocate", "Legal")},
		},
		13: {
			Summary: "carousing with the art community",
			Effects: []Effect{
				skill("Carouse"),
				rollSub("whether it takes hold",
					onRange(1, 2, "it does", addiction("a drug of the character's choice")),
					onRange(3, 6, "it does not")),
			},
		},
		14: {
			Summary: "a Rival sets out to sully your reputation",
			Effects: []Effect{
				relationship(Rival, 1, ""),
				checkChr("CHA", 8,
					[]Effect{pick("what surviving it earned",
						opt("a benefit roll", benefitRolls(1, 0, ScopeBatch)),
						opt("+1 CHA", chr("CHA", 1)))},
					[]Effect{chr("CHA", -2), throwModifier("next survival roll", -2)}),
			},
		},
		15: {
			Summary: "an attack on someone in the community sends you to a self-defence course",
			Effects: []Effect{pick("what you learned",
				opt("Gun Combat (Slug Pistol)", skill("Gun Combat", "Slug Pistol")),
				opt("Melee (Unarmed Combat)", skill("Melee", "Unarmed Combat")))},
		},
		16: {Summary: "a first aid course", Effects: []Effect{skill("Medic", "First Aid")}},
		21: {
			Summary: "a ranch on a nearby world, and being taught to ride",
			Effects: []Effect{skill("Animals", "Riding")},
		},
		22: {
			Summary: "a criminal wants fake art made to sell offworld",
			Effects: []Effect{pick("take the commission or decline",
				opt("decline", relationship(Enemy, 1, "")),
				opt("accept", checkSkill("Deception", 8,
					[]Effect{relationship(Contact, 1, ""), benefitRolls(2, 0, ScopeBatch)},
					[]Effect{relationship(Enemy, 2, ""), chr("CHA", -2)})))},
		},
		23: {
			Summary: "a friend does something stupid and the coverage lands on you",
			Effects: []Effect{
				chr("CHA", -1),
				loseTie(Contact),
			},
		},
		24: {
			Summary: "a controversial work, and an invitation to exhibit it",
			Effects: []Effect{pick("show it or withhold it",
				opt("exhibit it",
					relationship(Contact, 0, "1d3"),
					benefitRolls(-1, 0, ScopeBatch)),
				opt("withhold it",
					loseTiesRolled("1d3", Contact),
					chr("CHA", 1)))},
		},
		25: {
			Summary: "a gambling group in the art community",
			Effects: []Effect{checkSkill("Gambler", 8,
				[]Effect{
					benefitRolls(2, 0, ScopeBatch),
					pick("what the game taught you",
						opt("Gambler", skill("Gambler")),
						opt("Deception (Lie)", skill("Deception", "Lie"))),
				},
				[]Effect{benefitRolls(-2, 0, ScopeBatch)})},
		},
		26: {
			Summary: "law enforcement asks your help with an art-fraud investigation",
			Effects: []Effect{checkSkill("Investigate", 8,
				[]Effect{chr("CHA", 1), relationship(Contact, 1, "")},
				[]Effect{skill("Deception", "Forgery")})},
		},
		41: {Summary: "a decision to get fit", Effects: []Effect{skill("Athletics", "Any")}},
		42: {
			Summary: "driving yourself to art events",
			Effects: []Effect{pick("what you drove",
				opt("Drive", skill("Drive", "Any")),
				opt("Flyer", skill("Flyer", "Any")))},
		},
		43: {
			Summary: "a local religion, delved into deeply",
			Effects: []Effect{group(
				religion(),
				rollSub("how deeply it takes",
					onRange(1, 5, "it stays an interest"),
					on(6, "it becomes a calling", skill("Science", "Philosophy"))))},
		},
		44: {
			Summary: "legendary dinner parties",
			Effects: []Effect{
				chr("CHA", 2),
				pick("what the parties taught you",
					opt("Carouse", skill("Carouse")),
					opt("Etiquette", skill("Etiquette")),
					opt("Chef", skill("Chef"))),
			},
		},
		45: {
			Summary: "falling in love",
			Effects: []Effect{rollSub("what came of it",
				onRange(1, 2, "unrequited, and it feeds the work",
					benefitRolls(1, 0, ScopeBatch)),
				onRange(3, 4, "it ends amicably", relationship(Contact, 1, "")),
				onRange(5, 6, "it is mutual", relationship(Ally, 1, "")))},
		},
		46: {
			Summary: "space travel as a necessity",
			Effects: []Effect{pick("what the travelling taught you",
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
				opt("Survival (Freefall)", skill("Survival", "Freefall")))},
		},
		51: {
			Summary: "time aboard starships, learning from the crew",
			Effects: []Effect{pick("what you learned",
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Gunner", skill("Gunner", "Any")),
				opt("Pilot (Starship)", skill("Pilot", "Starship")))},
		},
		52: {
			Summary: "in this business it is important to know people",
			Effects: []Effect{relationship(Contact, 0, "1d3")},
		},
		53: {Summary: "how art relates to science", Effects: []Effect{skill("Science", "Any")}},
		54: {Summary: "a great deal of work with electronics", Effects: []Effect{skill("Electronics", "Any")}},
		55: {
			Summary: "time taken to learn a second language",
			Effects: []Effect{skill("Language", "Any")},
		},
		56: {
			Summary: "people rave about you in public and in private",
			Effects: ravedAbout(),
		},
		61: {
			Summary: "your most recent work is very well received",
			Effects: []Effect{advance(), throwModifier("next survival roll", 2)},
		},
		62: {
			Summary: "selling your own work, or the work of others",
			Effects: []Effect{checkSkill("Broker", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch)},
				[]Effect{benefitRolls(-1, 0, ScopeBatch)})},
		},
		63: {
			// ERRATA E-12: the six branches turn on Diplomat checks named
			// by difficulty -- Easy, Routine, Difficult, Very Difficult --
			// which is the Core Rulebook's task system and outside this
			// ruleset.
			Summary: "an invitation onto the subsector's most popular interview show",
			Effects: []Effect{pick("accept or refuse",
				opt("refuse",
					loseRank(1),
					relationship(Enemy, 1, "")),
				opt("accept", rollSub("how the host is disposed towards you (p. 158)",
					on(1, "the host hates your work and means to embarrass you",
						unimplemented("a Diplomat check at Very Difficult: "+
							"failing takes -2 CHA and -2 to the next survival roll, "+
							"succeeding gains +2 CHA and +2 to it (ERRATA E-12)")),
					onRange(2, 3, "the host did not want you on the show",
						unimplemented("a Diplomat check at Difficult: "+
							"failing takes -1 CHA, succeeding gains +1 (ERRATA E-12)")),
					onRange(4, 5, "the host is somewhat impressed",
						unimplemented("a Diplomat check at Routine: "+
							"failing takes -2 to the next survival roll, succeeding gains "+
							"+2 to it and a Contact (ERRATA E-12)")),
					on(6, "the host is already enamoured of your work",
						unimplemented("a Diplomat check at Easy: "+
							"succeeding gains the host as an Ally and +2 to the next survival "+
							"roll, failing takes -2 to it and gains an Enemy (ERRATA E-12)")))))},
		},
		64: {
			Summary: "very talented in your field",
			Effects: []Effect{talentedInYourField()},
		},
		65: {
			Summary: "the planetary government wants to feature you and buy from your portfolio",
			Effects: []Effect{pick("take part or decline",
				opt("take part",
					chr("CHA", 2),
					benefitRolls(1, 0, ScopeBatch),
					relationship(Rival, 1, "")),
				opt("decline", relationship(Ally, 1, "")))},
		},
		66: {
			Summary: "the most prestigious award in your field",
			Effects: append([]Effect{chr("CHA", 2)}, prestigiousAward()...),
		},
	}

	lifeEventRows(table)

	return table
}
