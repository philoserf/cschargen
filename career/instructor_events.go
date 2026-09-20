package career

// gamblingCircle is the gambling group, printed with the same two branches
// in six careers.
func gamblingCircle(second string) Effect {
	return checkSkill("Gambler", 8,
		[]Effect{
			benefitRolls(2, 0, ScopeBatch),
			pick("what the game taught you",
				opt("Gambler", skill("Gambler")),
				opt(second, skill(second))),
		},
		[]Effect{benefitRolls(-2, 0, ScopeBatch)},
	)
}

// instructorEvents is the d66 table of pp. 214-215.
func instructorEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "a Rival challenges your work",
			Effects: []Effect{
				relationship(Rival, 1, ""),
				checkSkill("Admin", 8,
					[]Effect{advance()},
					[]Effect{unimplemented("lose one rank")}),
			},
		},
		13: {
			Summary: "some of your money, invested",
			Effects: []Effect{checkSkill("Broker", 8,
				[]Effect{benefitRolls(3, 0, ScopeBatch)},
				[]Effect{benefitRolls(-3, 0, ScopeBatch)})},
		},
		14: {
			Summary: "money is tight, and the methods get shady",
			Effects: []Effect{pickSkill("Deception", "Streetwise")},
		},
		15: {Summary: "a day out on the water", Effects: []Effect{skill("Seafarer", "Any")}},
		16: {
			Summary: "a student who wants to argue about everything",
			Effects: []Effect{checkSkill("Persuade", 8,
				[]Effect{throwModifier("next advancement roll", 2)},
				[]Effect{throwModifier("next advancement roll", -2)})},
		},
		21: {
			Summary: "your position used to advance your political causes",
			Effects: []Effect{checkSkill("Advocate", 8,
				[]Effect{throwModifier("next advancement roll", 2)},
				[]Effect{throwModifier("next advancement roll", -2)})},
		},
		22: {Summary: "good times with faculty and students", Effects: []Effect{skill("Carouse")}},
		23: {
			Summary: "a term spent smoothing over disagreements",
			Effects: []Effect{pickSkill("Diplomat", "Etiquette")},
		},
		24: {Summary: "skills honed with electronics", Effects: []Effect{skill("Electronics", "Any")}},
		25: {Summary: "a gambling circle in your peer group", Effects: []Effect{gamblingCircle("Persuade")}},
		26: {Summary: "a cooking class, taken", Effects: []Effect{skill("Chef")}},
		41: {Summary: "riding, taken up as a hobby", Effects: []Effect{skill("Animals", "Riding")}},
		42: {
			Summary: "a self-defence class",
			Effects: []Effect{pick("what you learned",
				opt("Gun Combat", skill("Gun Combat", "Slug Pistol", "Shotgun", "Energy Pistol")),
				opt("Melee", skill("Melee", "Any")))},
		},
		43: {
			Summary: "people from a wide variety of backgrounds",
			Effects: []Effect{skill("Language", "Any")},
		},
		44: {
			Summary: "time spent organizing your fellow educators into a union",
			Effects: []Effect{pick("how you organized them",
				opt("Advocate", skill("Advocate", "Oratory", "Politics")),
				opt("Leadership", skill("Leadership")),
				opt("Persuade", skill("Persuade")))},
		},
		45: {Summary: "a survival course", Effects: []Effect{pickSkill("Survival", "Navigation")}},
		46: {
			Summary: "education is a social affair",
			Effects: []Effect{pick("what came of it",
				opt("Carouse", skill("Carouse")),
				opt("contacts", relationship(Contact, 0, "1d3")))},
		},
		51: {
			Summary: "hitting the gym",
			Effects: []Effect{pick("what the training gave you",
				opt("+1 STR", chr("STR", 1)),
				opt("+1 DEX", chr("DEX", 1)),
				opt("+1 END", chr("END", 1)),
				opt("Athletics", skill("Athletics", "Any")))},
		},
		52: {
			Summary: "providing your own transportation",
			Effects: []Effect{pickSkill("Drive", "Mechanic", "Flyer")},
		},
		53: {Summary: "a first aid class", Effects: []Effect{skill("Medic", "First Aid")}},
		54: {
			Summary: "a local amateur sports team",
			Effects: []Effect{pick("what the team taught you",
				opt("Athletics", skill("Athletics", "Any")),
				opt("Tactics (Sport)", skill("Tactics", "Sport")))},
		},
		55: {
			Summary: "the academic bowl team needs a coach",
			Effects: []Effect{pick("accept or decline",
				opt("decline", throwModifier("next advancement roll", -2)),
				opt("coach them", pick("how you coached",
					opt("Instruction", checkSkill("Instruction", 8,
						[]Effect{benefitRolls(1, 0, ScopeBatch)},
						[]Effect{relationship(Contact, 0, "1d3")})),
					opt("Persuade", checkSkill("Persuade", 8,
						[]Effect{benefitRolls(1, 0, ScopeBatch)},
						[]Effect{relationship(Contact, 0, "1d3")})),
					opt("EDU", checkChr("EDU", 8,
						[]Effect{benefitRolls(1, 0, ScopeBatch)},
						[]Effect{relationship(Contact, 0, "1d3")})))))},
		},
		56: {
			Summary: "another educator thinks you less of a professional than they are",
			Effects: []Effect{relationship(Rival, 1, "")},
		},
		61: {
			Summary: "downtime spent aboard a research vessel",
			Effects: []Effect{pick("what you learned",
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Pilot", skill("Pilot", "Any")),
				opt("Science", skill("Science", "Any")),
				opt("Survival (Freefall)", skill("Survival", "Freefall")),
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")))},
		},
		62: {
			Summary: "a lot of odd jobs around your place of employment",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		63: {
			Summary: "constantly improving your skill set and maintaining standards",
			Effects: []Effect{chr("EDU", 1)},
		},
		64: {
			Summary: "honoured as teacher of the year",
			Effects: []Effect{benefitRolls(1, 0, ScopeBatch), advance()},
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

// talentedInYourField is event 65 in a dozen careers, printed the same way
// each time: a level in a skill already held, or a promotion.
func talentedInYourField() Effect {
	return pick("what the talent earned",
		opt("a level in a skill you already have",
			raiseHeldSkill()),
		opt("a promotion", Effect{Kind: EffectRank, Detail: "gain a rank"}),
	)
}
