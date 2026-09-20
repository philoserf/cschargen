package career

// politicalCoworker is the coworker with causes, printed identically in
// Belter (p. 162) and Craftsperson (p. 183).
func politicalCoworker() Effect {
	return pick("side with them, side with the company, or ignore it",
		opt("side with the coworker",
			relationship(Ally, 1, ""),
			throwModifier("next advancement roll", -2)),
		opt("side with the employers",
			throwModifier("next advancement roll", 2),
			relationship(Rival, 1, "")),
		opt("ignore them", relationship(Rival, 1, ""), skill("Persuade")),
	)
}

// spottedSomethingWrong is the teammate doing something dangerous or
// illegal, printed in five careers with the same three choices.
func spottedSomethingWrong(middle Effect) Effect {
	return pick("turn them in, correct them, or look away",
		opt("turn them in", throwModifier("next advancement roll", 2)),
		opt("instruct them otherwise", middle),
		opt("ignore it", relationship(Ally, 1, "")),
	)
}

// craftEvents is the d66 table of pp. 183-184.
func craftEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {Summary: "keeping the books for the business", Effects: []Effect{skill("Admin")}},
		13: {
			Summary: "a farm or ranch under construction, and an interest in what happens there",
			Effects: []Effect{skill("Animals", "Any")},
		},
		14: {
			Summary: "a chosen drinking establishment after a long day",
			Effects: []Effect{pickSkill("Carouse", "Streetwise")},
		},
		15: {Summary: "workers from a variety of backgrounds", Effects: []Effect{skill("Language", "Any")}},
		16: {
			Summary: "a shortage of hands puts you outside your specialty",
			Effects: []Effect{rollOtherAssignment()},
		},
		21: {
			Summary: "the laws and codes of construction work",
			Effects: []Effect{skill("Advocate", "Legal")},
		},
		22: {
			Summary: "training",
			Effects: []Effect{pick("what the training gave you",
				opt("Athletics", skill("Athletics", "Any")),
				opt("+1 STR", chr("STR", 1)),
				opt("+1 DEX", chr("DEX", 1)),
				opt("+1 END", chr("END", 1)))},
		},
		23: {
			Summary: "driving yourself to the jobsite",
			Effects: []Effect{pick("what you drove",
				opt("Drive", skill("Drive", "Wheeled", "Tracked")),
				opt("Flyer (Grav)", skill("Flyer", "Grav")))},
		},
		24: {Summary: "downtime at the range", Effects: []Effect{skill("Gun Combat", "Any")}},
		25: {
			Summary: "bad planning runs the job over budget, and your pay is cut",
			Effects: []Effect{benefitRolls(-1, 0, ScopeBatch), throwModifier("next advancement roll", 2)},
		},
		26: {
			Summary: "a gossiping coworker spreads something damaging about your past",
			Effects: []Effect{relationship(Rival, 1, ""), throwModifier("next advancement roll", -2)},
		},
		41: {
			Summary: "a politically active coworker, popular with the workers and not the owners",
			Effects: []Effect{politicalCoworker()},
		},
		42: {Summary: "a hobby taken up", Effects: []Effect{pickSkill("Art", "Science")}},
		43: {Summary: "particularly stubborn rock", Effects: []Effect{skill("Explosives")}},
		44: {
			Summary: "you spot a teammate doing something dangerous or illegal",
			Effects: []Effect{spottedSomethingWrong(skill("Leadership"))},
		},
		45: {Summary: "repairing your own equipment", Effects: []Effect{skill("Mechanic")}},
		46: {
			Summary: "local codes want someone trained in first aid, and your boss picks you",
			Effects: []Effect{skill("Medic", "First Aid")},
		},
		51: {Summary: "the hobby of riding", Effects: []Effect{skill("Animals", "Riding")}},
		52: {
			Summary: "selling your own skills",
			Effects: []Effect{pick("how you sold them",
				opt("Broker", skill("Broker")),
				opt("Advocate (Oratory)", skill("Advocate", "Oratory")))},
		},
		53: {
			Summary: "a gambling circle in your peer group",
			Effects: []Effect{checkSkill("Gambler", 8,
				[]Effect{
					benefitRolls(2, 0, ScopeBatch),
					pick("what the game taught you",
						opt("Gambler", skill("Gambler")),
						opt("Persuade", skill("Persuade"))),
				},
				[]Effect{benefitRolls(-2, 0, ScopeBatch)})},
		},
		54: {Summary: "studying self-defence", Effects: []Effect{skill("Melee", "Any")}},
		55: {
			Summary: "a jobsite full of people, uplifts, robots and machinery",
			Effects: []Effect{skill("Recon")},
		},
		56: {
			Summary: "shoddy practices by your project manager, in violation of the codes",
			Effects: []Effect{shoddyPractices()},
		},
		61: {Summary: "friends in this business", Effects: []Effect{relationship(Contact, 0, "1d3")}},
		62: {
			Summary: "another normal day doing your job and someone else's",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		63: {
			Summary: "a superior sets out to groom you for higher things",
			Effects: groomedForHigherThings(),
		},
		64: {
			Summary: "excellent work by the team earns a bonus",
			Effects: []Effect{cashRollsNowRerolling(2)},
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
