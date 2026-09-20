package career

// ravedAbout is the "great person in public and in private" event, printed
// in both Arts (p. 158) and Celebrity (p. 167).
func ravedAbout() []Effect {
	return []Effect{
		relationship(Contact, 0, "1d3"),
		becomesRolled("1d3", Ally, Contact),
	}
}

// celebrityEvents is the d66 table of pp. 166-168.
func celebrityEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "a critic decides your career is a waste, and writes to destroy it",
			Effects: []Effect{
				relationship(Rival, 1, ""),
				checkSkill("Diplomat", 8,
					[]Effect{benefitRolls(1, 0, ScopeBatch)},
					[]Effect{
						benefitRolls(-1, 0, ScopeBatch),
						throwModifier("next survival roll", -2),
					}),
			},
		},
		13: {
			Summary: "a friend does something stupid and the coverage lands on you",
			Effects: []Effect{
				loseTie(Contact),
				benefitRolls(-1, 0, ScopeBatch),
			},
		},
		14: {
			Summary: "a sex scandal",
			Effects: []Effect{pick("manage it however you can",
				opt("Persuade", checkSkill("Persuade", 8, celebritySpun(), celebrityStung())),
				opt("Diplomat", checkSkill("Diplomat", 8, celebritySpun(), celebrityStung())))},
		},
		15: {
			Summary: "a great deal of fun being famous",
			Effects: []Effect{
				skill("Carouse"),
				chr("CHA", 1),
				unimplemented("on a 1d6 of 1-2, an addiction to alcohol or a drug"),
			},
		},
		16: {
			Summary: "your agent absconds with all your money",
			Effects: []Effect{benefitRolls(-3, 0, ScopeBatch)},
		},
		21: {
			Summary: "time with a personal trainer",
			Effects: []Effect{pick("what the training gave you",
				opt("+1 STR", chr("STR", 1)),
				opt("+1 DEX", chr("DEX", 1)),
				opt("+1 END", chr("END", 1)),
				opt("Athletics", skill("Athletics", "Any")))},
		},
		22: {
			Summary: "an invitation onto a celebrity game show",
			Effects: []Effect{pick("take the booking or refuse",
				opt("refuse", throwModifier("next advancement roll", 2)),
				opt("take it",
					credits("5000"),
					unimplemented(
						"roll 1d6 for the kind of show: a farce to embarrass you, a quiz, a test of "+
							"charisma, or a triumph worth 10,000 credits and a benefit roll")))},
		},
		23: {
			Summary: "a religion, delved into deeply",
			Effects: []Effect{unimplemented(
				"choose or invent a religion; on a 1d6 of 6 you become deeply involved and gain Science (Philosophy) 1")},
		},
		24: {
			Summary: "in this business it is important to know people",
			Effects: []Effect{relationship(Contact, 0, "1d3")},
		},
		25: {
			Summary: "a love of operating vehicles",
			Effects: []Effect{pickSkill("Drive", "Flyer", "Seafarer")},
		},
		26: {Summary: "constantly selling your services", Effects: []Effect{skill("Broker")}},
		41: {
			Summary: "the holography side of the work",
			Effects: []Effect{skill("Art", "Holography")},
		},
		42: {Summary: "first aid learned", Effects: []Effect{skill("Medic", "First Aid")}},
		43: {
			Summary: "your agent sends you outside your comfort zone",
			Effects: []Effect{rollOtherAssignment()},
		},
		44: {Summary: "a weapons class", Effects: []Effect{skill("Gun Combat", "Any")}},
		45: {
			Summary: "understanding the common people, for a role or a song",
			Effects: []Effect{unimplemented(
				"roll once on the Service Skills table of any other career")},
		},
		46: {Summary: "frequent travel between star systems", Effects: []Effect{skill("Suit", "Any")}},
		51: {Summary: "riding, taken up to relax", Effects: []Effect{skill("Animals", "Riding")}},
		52: {
			Summary: "a tell-all book about you and your career",
			Effects: []Effect{throwModifier("next survival roll", -2)},
		},
		53: {
			// ERRATA E-1: the failure branch cites p. 136.
			Summary: "a crazed fan gets past security and lunges at you",
			Effects: []Effect{checkSkill("Melee", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch), throwModifier("next survival roll", 2)},
				[]Effect{injury(1)})},
		},
		54: {
			Summary: "time taken to learn a second language",
			Effects: []Effect{{
				Kind: EffectSkill, Detail: "gain Language 2 in a language of the homeworld",
				Skill: "Language", Specialties: []string{"Any"}, Level: 2,
			}},
		},
		55: {
			Summary: "an endorsement deal for your likeness",
			Effects: []Effect{
				credits("25000"),
				unimplemented(
					"roll 1d6 for how the product fares: from an abject failure that drags you " +
						"down to a major success worth 100,000 credits and three benefit rolls"),
			},
		},
		56: {
			Summary: "downtime spent studying a subject as a hobby",
			Effects: []Effect{pickSkill("Art", "Science")},
		},
		61: {
			Summary: "your most recent work is an astounding success",
			Effects: []Effect{advance(), throwModifier("next survival roll", 2)},
		},
		62: {Summary: "your agent is a genius", Effects: []Effect{benefitRolls(2, 0, ScopeBatch)}},
		63: {
			Summary: "people rave about you in public and in private",
			Effects: ravedAbout(),
		},
		64: {Summary: "extremely talented in your field", Effects: []Effect{talentedInYourField()}},
		65: {
			// ERRATA E-12 applies to the 1d6 branches, which name Diplomat
			// and Etiquette checks at a difficulty the Core Rulebook sets.
			Summary: "an invitation onto the subsector's most popular interview show",
			Effects: []Effect{pick("accept or refuse",
				opt("refuse",
					throwModifier("next survival roll", -2),
					throwModifier("next advancement roll", -3)),
				opt("accept",
					credits("5000"),
					unimplemented(
						"roll 1d6 for the host's disposition, then a Diplomat or Etiquette check; "+
							"the outcomes run from severe damage to the career to the host as an Ally")))},
		},
		66: {
			Summary: "recognized as one of the best at what you do",
			Effects: []Effect{
				mustContinue(),
				autoSuccess("next survival roll"),
				autoSuccess("next advancement roll"),
			},
		},
	}

	lifeEventRows(table)

	return table
}

func celebritySpun() []Effect {
	return []Effect{benefitRolls(1, 0, ScopeBatch), skill("Diplomat")}
}

func celebrityStung() []Effect {
	return []Effect{
		benefitRolls(-1, 0, ScopeBatch),
		throwModifier("next survival roll", -2),
	}
}
