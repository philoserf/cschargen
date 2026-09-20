package career

// clergyEvents is the d66 table of pp. 171-172.
func clergyEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "someone offended by your religion attacks you",
			Effects: []Effect{pick("defend yourself however you can",
				opt("Gun Combat", checkSkill("Gun Combat", 8,
					[]Effect{skill("Gun Combat", "Any")}, []Effect{injury(1)})),
				opt("Melee", checkSkill("Melee", 8,
					[]Effect{skill("Melee", "Any")}, []Effect{injury(1)})))},
		},
		13: {Summary: "religious work means bureaucracy too", Effects: []Effect{skill("Admin")}},
		14: {
			Summary: "getting along well with parishioners or fellow monks",
			Effects: []Effect{skill("Carouse")},
		},
		15: {Summary: "you were not always a holy person", Effects: []Effect{skill("Deception", "Any")}},
		16: {
			Summary: "communicating and reading in a variety of languages",
			Effects: []Effect{skill("Language", "Any")},
		},
		21: {
			Summary: "seen cosying up to a religion yours calls heresy",
			Effects: []Effect{relationship(Contact, 1, ""), relationship(Enemy, 1, "")},
		},
		22: {
			Summary: "learning more about the laws of this world",
			Effects: []Effect{skill("Advocate", "Legal", "Politics")},
		},
		23: {
			Summary: "your own transportation, planetside",
			Effects: []Effect{pickSkill("Drive", "Flyer", "Navigation")},
		},
		24: {Summary: "a first aid course", Effects: []Effect{skill("Medic", "First Aid")}},
		25: {
			Summary: "a mob corners you and demands you change your religion to theirs",
			Effects: []Effect{checkSkill("Persuade", 8,
				[]Effect{advance(), skill("Persuade")},
				[]Effect{injury(1)})},
		},
		26: {
			Summary: "something said, done or discovered has offended someone powerful",
			Effects: []Effect{relationship(Enemy, 1, "")},
		},
		41: {Summary: "riding, taken up as a hobby", Effects: []Effect{skill("Animals", "Riding")}},
		42: {
			Summary: "frequent travel through space",
			Effects: []Effect{pick("what the travelling taught you",
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics", skill("Electronics", "Comms", "Sensors")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Pilot", skill("Pilot", "Any")))},
		},
		43: {Summary: "a self-defence class", Effects: []Effect{pickSkill("Gun Combat", "Melee")}},
		44: {Summary: "caring for your own machinery", Effects: []Effect{skill("Mechanic")}},
		45: {
			Summary: "convincing many of your beliefs",
			Effects: []Effect{pickSkill("Persuade", "Diplomat")},
		},
		46: {Summary: "sailing, which takes you away", Effects: []Effect{skill("Seafarer", "Any")}},
		51: {
			Summary: "a split in your religion, and a side to choose",
			Effects: []Effect{relationship(Ally, 1, ""), relationship(Enemy, 1, "")},
		},
		52: {
			Summary: "time at the gym",
			Effects: []Effect{pick("what the training gave you",
				opt("+1 STR", chr("STR", 1)),
				opt("+1 DEX", chr("DEX", 1)),
				opt("+1 END", chr("END", 1)),
				opt("Athletics", skill("Athletics", "Any")))},
		},
		53: {
			Summary: "a term managing the church's money",
			Effects: []Effect{skill("Broker")},
		},
		54: {
			Summary: "a murder inside the church, and the clergy look to you to solve it quietly",
			Effects: []Effect{checkSkill("Investigate", 8,
				[]Effect{skill("Investigate"), relationship(Enemy, 1, ""), relationship(Ally, 1, "")},
				[]Effect{relationship(Enemy, 1, ""), benefitRolls(-1, 0, ScopeBatch)})},
		},
		55: {Summary: "a great deal of cooking", Effects: []Effect{skill("Chef")}},
		56: {
			Summary: "understanding the seedier parts of life to reach those in it",
			Effects: []Effect{skill("Streetwise")},
		},
		61: {Summary: "a hobby taken up", Effects: []Effect{pickSkill("Art", "Science")}},
		62: {
			Summary: "studying and noticing the nature of the world",
			Effects: []Effect{pickSkill("Investigate", "Recon")},
		},
		63: {
			Summary: "time spent in space",
			Effects: []Effect{pick("what the time taught you",
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
				opt("Survival (Freefall)", skill("Survival", "Freefall")))},
		},
		64: {
			Summary: "a small house, and someone else's job to do",
			Effects: []Effect{rollOtherAssignment()},
		},
		65: {Summary: "a wide variety of jobs", Effects: []Effect{skill("Jack of All Trades")}},
		66: {
			Summary: "excellent work",
			Effects: []Effect{advance(), benefitRolls(0, 1, ScopeCareer)},
		},
	}

	lifeEventRows(table)

	return table
}
