package career

// belterEvents is the d66 table of pp. 162-163.
func belterEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "the company sheds staff, and moves you into the head office",
			Effects: []Effect{skill("Admin")},
		},
		13: {
			Summary: "downtime spent getting into trouble",
			Effects: []Effect{pick("what the trouble taught you",
				opt("Carouse", skill("Carouse")),
				opt("Deception", skill("Deception", "Any")),
				opt("Gambler", skill("Gambler")),
				opt("Streetwise", skill("Streetwise")))},
		},
		14: {
			Summary: "a shortage of hands puts you outside your specialty",
			Effects: []Effect{rollOtherAssignment()},
		},
		15: {
			Summary: "you spot a teammate doing something dangerous or illegal",
			Effects: []Effect{spottedSomethingWrong(skill("Leadership"))},
		},
		16: {
			Summary: "a barfight",
			Effects: []Effect{checkSkill("Melee", 8,
				[]Effect{relationship(Contact, 1, "")}, []Effect{injury(1)})},
		},
		21: {
			Summary: "a gambling circle aboard ship",
			Effects: []Effect{checkSkill("Gambler", 8,
				[]Effect{
					benefitRolls(2, 0, ScopeBatch),
					pick("what the game taught you",
						opt("Gambler", skill("Gambler")),
						opt("Persuade", skill("Persuade"))),
				},
				[]Effect{benefitRolls(-2, 0, ScopeBatch)})},
		},
		22: {Summary: "repairing your own equipment as you work", Effects: []Effect{skill("Mechanic")}},
		23: {
			Summary: "belters are often an unruly sort",
			Effects: []Effect{pick("what the company taught you",
				opt("Deception", skill("Deception", "Any")),
				opt("Melee", skill("Melee", "Any")),
				opt("Streetwise", skill("Streetwise")))},
		},
		24: {
			Summary: "basic first aid",
			// The engine granted this at level 1 until AtLevelZero
			// existed: an unset Level means "the result did not say",
			// which is one.
			Effects: []Effect{skillZero("Medic")},
		},
		25: {
			Summary: "chosen for the organization's training staff",
			Effects: []Effect{benefitRolls(0, 1, ScopeCareer), chr("CHA", 1)},
		},
		26: {
			Summary: "downtime spent on a hobby",
			Effects: []Effect{pick("what you pursued",
				opt("Art", skill("Art", "Any")),
				opt("Science", skill("Science", "Any")))},
		},
		41: {Summary: "downtime at the gun range", Effects: []Effect{skill("Gun Combat", "Any")}},
		42: {Summary: "belters come from many backgrounds", Effects: []Effect{skill("Language", "Any")}},
		43: {
			Summary: "a gossiping coworker spreads something damaging about your past",
			Effects: []Effect{relationship(Rival, 1, ""), throwModifier("next advancement roll", -2)},
		},
		44: {
			Summary: "a politically active coworker, popular with the workers and not the owners",
			Effects: []Effect{politicalCoworker()},
		},
		45: {
			Summary: "a coworker's book about belter life, with a thinly veiled you at its centre",
			Effects: []Effect{benefitRolls(1, 0, ScopeBatch)},
		},
		46: {
			Summary: "often chosen for jobs outside your comfort zone",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		51: {
			Summary: "belters come and go, and you have met a great many",
			Effects: []Effect{relationship(Contact, 0, "1d3")},
		},
		52: {
			Summary: "two corporations go to war and you are caught in the middle",
			Effects: []Effect{pick("profit from it, fight through it, or hide from it",
				opt("Deception", checkSkill("Deception", 8,
					[]Effect{cashRollsNow(2)},
					[]Effect{
						unimplemented("lose one rank"),
						benefitRolls(-2, 0, ScopeBatch),
					})),
				opt("Gun Combat", checkSkill("Gun Combat", 8,
					[]Effect{skill("Gun Combat", "Any"), throwModifier("next advancement roll", 2)},
					[]Effect{injury(1)})),
				opt("Stealth", checkSkill("Stealth", 8,
					[]Effect{skill("Survival", "Freefall")},
					[]Effect{injury(1)})))},
		},
		53: {Summary: "chosen for advanced training", Effects: []Effect{rollTable(AdvancedEducation)}},
		54: {Summary: "working extra hard to get things right", Effects: []Effect{rollTable(ServiceSkills)}},
		55: {
			Summary: "religion, delved into deeply",
			Effects: []Effect{unimplemented(
				"choose or invent a religion; on a 1d6 of 6 you become deeply involved and gain Science (Philosophy) 1")},
		},
		56: {
			Summary: "you strike it rich, and the company is thankful",
			Effects: []Effect{benefitRolls(0, 1, ScopeCareer), autoSuccess("next advancement roll")},
		},
		61: {
			Summary: "excellent work earns a bonus",
			Effects: []Effect{cashRollsNow(2)},
		},
		62: {
			Summary: "honing your skillset",
			Effects: []Effect{raiseHeldSkill()},
		},
		63: {
			Summary: "a tunnel collapse with your friends trapped inside",
			Effects: []Effect{pick("reach them however you can",
				opt("Suit (Vacc Suit)", checkSkill("Suit", 8,
					[]Effect{relationship(Contact, 0, "1d3"), chr("CHA", 1)},
					[]Effect{chr("CHA", -2)})),
				opt("Survival (Freefall)", checkSkill("Survival", 8,
					[]Effect{relationship(Contact, 0, "1d3"), chr("CHA", 1)},
					[]Effect{chr("CHA", -2)})))},
		},
		64: {
			Summary: "a small sphere, a flash of blue light, and dreams that never stop",
			Effects: []Effect{chr("INT", 2)},
		},
		65: {
			Summary: "a superior sets out to groom you for higher things",
			Effects: groomedForHigherThings(),
		},
		66: {
			Summary: "a massive deposit uncovered, and the company is grateful",
			Effects: []Effect{
				stashCountValued(10, companyShare, companyShareValue),
				{Kind: EffectRank, Detail: "gain a rank"},
				{Kind: EffectRank, Detail: "gain a second rank"},
				relationship(Ally, 1, ""),
			},
		},
	}

	lifeEventRows(table)

	return table
}
