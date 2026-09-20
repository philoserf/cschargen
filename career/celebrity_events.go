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
						throwModifier(survivalThrow, -2),
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
				rollSub("whether it takes hold",
					onRange(1, 2, "it does",
						addiction("alcohol or a drug of the character's choice")),
					onRange(3, 6, "it does not")),
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
				opt("refuse", throwModifier(advancementThrow, 2)),
				opt("take it",
					credits("5000"),
					rollSub("what kind of show it turns out to be (p. 166)",
						on(1, "a farce meant to embarrass its guests",
							pick("turn the tables with diplomacy or with charm",
								opt("Diplomat", checkSkill("Diplomat", 8,
									[]Effect{
										benefitRolls(2, 0, ScopeBatch),
										throwModifier(survivalThrow, 2),
									},
									[]Effect{throwModifier(survivalThrow, -2)})),
								opt("Persuade", checkSkill("Persuade", 8,
									[]Effect{
										benefitRolls(2, 0, ScopeBatch),
										throwModifier(survivalThrow, 2),
									},
									[]Effect{throwModifier(survivalThrow, -2)})))),
						onRange(2, 3, "a quiz show of facts and trivia",
							pick("answer on what you know or on how you were taught",
								opt("INT", checkChr("INT", 8,
									[]Effect{benefitRolls(1, 0, ScopeBatch)},
									[]Effect{throwModifier(survivalThrow, -2)})),
								opt("EDU", checkChr("EDU", 8,
									[]Effect{benefitRolls(1, 0, ScopeBatch)},
									[]Effect{throwModifier(survivalThrow, -2)})))),
						onRange(4, 5, "a test of wit and charisma",
							checkChr("CHA", 8,
								[]Effect{benefitRolls(2, 0, ScopeBatch)},
								[]Effect{throwModifier(survivalThrow, -2)})),
						on(6, "a triumph",
							credits("10000"),
							benefitRolls(1, 0, ScopeBatch),
							throwModifier(survivalThrow, 2)))))},
		},
		23: {
			Summary: "a religion, delved into deeply",
			Effects: []Effect{group(
				religion(),
				rollSub("how deeply it takes",
					onRange(1, 5, "it stays an interest"),
					on(6, "it becomes a calling", skill("Science", "Philosophy"))))},
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
			Effects: []Effect{throwModifier(survivalThrow, -2)},
		},
		53: {
			// ERRATA E-1: the failure branch cites p. 136.
			Summary: "a crazed fan gets past security and lunges at you",
			Effects: []Effect{checkSkill("Melee", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch), throwModifier(survivalThrow, 2)},
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
				rollSub("how the product fares (p. 167)",
					on(1, "an abject failure, which drags you down with it",
						throwModifier(advancementThrow, -2)),
					onRange(2, 3, "a success, and an embarrassing one over time",
						throwModifier(survivalThrow, -1)),
					onRange(4, 5, "a success, and a bonus with it",
						credits("10000"),
						throwModifier(survivalThrow, 1)),
					on(6, "a major success, and the contract extended",
						credits("100000"),
						benefitRolls(3, 0, ScopeBatch))),
			},
		},
		56: {
			Summary: "downtime spent studying a subject as a hobby",
			Effects: []Effect{pickSkill("Art", "Science")},
		},
		61: {
			Summary: "your most recent work is an astounding success",
			Effects: []Effect{advance(), throwModifier(survivalThrow, 2)},
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
					throwModifier(survivalThrow, -2),
					throwModifier(advancementThrow, -3)),
				opt("accept",
					credits("5000"),
					rollSub("how the host is disposed towards you (p. 168)",
						on(1, "the host hates your work and means to embarrass you",
							checkSkill("Diplomat", 8,
								[]Effect{skill("Diplomat"), throwModifier(survivalThrow, 2)},
								[]Effect{loseRank(1), throwModifier(advancementThrow, -2)})),
						onRange(2, 3, "the host did not want you on the show",
							checkSkill("Diplomat", 8, nil,
								[]Effect{throwModifier(advancementThrow, -2)})),
						onRange(4, 5, "the host is somewhat impressed",
							checkSkill("Diplomat", 8,
								[]Effect{
									throwModifier(advancementThrow, 2),
									relationship(Contact, 1, ""),
								},
								[]Effect{throwModifier(survivalThrow, -2)})),
						on(6, "the host is already enamoured of your work",
							checkSkill("Etiquette", 8,
								[]Effect{
									relationship(Ally, 1, ""),
									throwModifier(survivalThrow, 2),
									throwModifier(advancementThrow, 2),
								},
								[]Effect{
									throwModifier(survivalThrow, -2),
									relationship(Enemy, 1, ""),
								})))))},
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
		throwModifier(survivalThrow, -2),
	}
}
