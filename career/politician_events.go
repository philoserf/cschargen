package career

// politicianEvents is the d66 table of pp. 258-260.
func politicianEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {Summary: "surrounded by forms and paperwork", Effects: []Effect{skill("Admin")}},
		13: {
			Summary: "campaign finances are hard to oversee",
			Effects: []Effect{pick("how you oversaw them",
				opt("Advocate (Legal)", skill("Advocate", "Legal")),
				opt("Broker", skill("Broker")))},
		},
		14: {
			Summary: "keeping up with the technology of public life",
			Effects: []Effect{skill("Electronics", "Any")},
		},
		15: {
			Summary: "a high official as an Ally, and a boorish one",
			Effects: []Effect{
				relationship(Ally, 1, ""),
				pick("distance yourself however you can",
					opt("Advocate", checkSkill("Advocate", 8,
						[]Effect{throwModifier("next advancement roll", 2)}, politicianTurnedOn())),
					opt("Diplomat", checkSkill("Diplomat", 8,
						[]Effect{throwModifier("next advancement roll", 2)}, politicianTurnedOn()))),
			},
		},
		16: {
			Summary: "a journalist calls you corrupt",
			Effects: []Effect{
				relationship(Rival, 1, ""),
				checkSkill("Advocate", 8,
					[]Effect{chr("CHA", 1)},
					[]Effect{checkChr("CHA", 8,
						[]Effect{throwModifier("next advancement roll", 2)},
						[]Effect{
							unimplemented("automatically fail the next advancement roll"),
							chr("CHA", -2),
						})}),
			},
		},
		21: {
			Summary: "a debate, soon to be broadcast",
			Effects: []Effect{
				relationship(Rival, 1, ""),
				pick("meet it as your posting requires",
					opt("Advocate (Oratory), as a Candidate", checkSkill("Advocate", 8,
						[]Effect{throwModifier("next advancement roll", 2)},
						[]Effect{throwModifier("next advancement roll", -2), chr("CHA", -2)})),
					opt("Advocate (Politics), as a Manager", checkSkill("Advocate", 8,
						[]Effect{benefitRolls(2, 0, ScopeBatch), throwModifier("next advancement roll", 2)},
						[]Effect{throwModifier("next advancement roll", -2)})),
					opt("a dirty trick, as an Operative", checkSkill("Deception", 8,
						[]Effect{
							relationship(Enemy, 1, ""),
							throwModifier("next advancement roll", 2),
							benefitRolls(1, 0, ScopeBatch),
						},
						[]Effect{
							throwModifier("next advancement roll", -4),
							transfer("Prisoner", "Prisoner", 1),
						})),
					opt("refuse the dirty trick",
						throwModifier("next advancement roll", -2))),
			},
		},
		22: {Summary: "rather popular", Effects: []Effect{skill("Carouse")}},
		23: {
			Summary: "a chance to undermine a member of your own party",
			Effects: []Effect{
				relationship(Rival, 1, ""),
				pick("leak it or sit on it",
					opt("leak it", checkSkill("Persuade", 8,
						[]Effect{throwModifier("next advancement roll", 4), relationship(Ally, 1, "")},
						[]Effect{
							throwModifier("next advancement roll", -4),
							relationship(Contact, 1, ""),
							relationship(Enemy, 1, ""),
						})),
					opt("sit on it")),
			},
		},
		24: {Summary: "a gambling circle in your peer group", Effects: []Effect{gamblingCircle("Persuade")}},
		25: {Summary: "downtime at the firing range", Effects: []Effect{skill("Gun Combat", "Any")}},
		26: {Summary: "a first aid course", Effects: []Effect{skill("Medic", "First Aid")}},
		41: {
			Summary: "a debate over what your party now believes",
			Effects: []Effect{checkSkill("Advocate", 8,
				[]Effect{advance()},
				[]Effect{unimplemented("lose one rank")})},
		},
		42: {
			Summary: "politics is often shady business",
			Effects: []Effect{pickSkill("Deception", "Streetwise")},
		},
		43: {
			Summary: "correspondence mishandled, and a rival party reading it",
			Effects: []Effect{checkSkill("Electronics", 8, nil,
				[]Effect{throwModifier("next advancement roll", -3)})},
		},
		44: {
			Summary: "appealing to a wide variety of people",
			Effects: []Effect{skill("Language", "Any")},
		},
		45: {
			Summary: "recognized by your government for exemplary service",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch)},
		},
		46: {
			Summary: "a secretive group inside your own government makes an approach",
			Effects: []Effect{theConspiracyRecruits()},
		},
		51: {
			Summary: "a hobby, to get away from the political madhouse",
			Effects: []Effect{pickSkill("Art", "Animals", "Drive", "Flyer", "Mechanic", "Science", "Seafarer")},
		},
		52: {
			Summary: "a major disagreement in the party, smoothed over",
			Effects: []Effect{skill("Diplomat")},
		},
		53: {Summary: "a self-defence course", Effects: []Effect{skill("Melee", "Any")}},
		54: {
			Summary: "a mock debate with fellow party members",
			Effects: []Effect{pickSkill("Diplomat", "Advocate", "Persuade")},
		},
		55: {
			Summary: "a colleague seeking promotion by denigrating your work",
			Effects: denigratedByAColleague(),
		},
		56: {
			Summary: "your service is noted, and your government places you in the foreign service",
			Effects: []Effect{pick("which posting",
				opt("Ambassador", transfer("Diplomatic Service", "Ambassador", 0)),
				opt("Generalist", transfer("Diplomatic Service", "Generalist", 0)))},
		},
		61: {
			Summary: "a great deal of starship travel recently",
			Effects: []Effect{pick("what the travelling taught you",
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Gunner (Turrets)", skill("Gunner", "Turrets")),
				opt("Pilot", skill("Pilot", "Any")))},
		},
		62: {
			Summary: "a hands-on approach to everything",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		63: {
			Summary: "appearances as an expert make you famous",
			Effects: []Effect{transfer("Celebrity", "Star", 0)},
		},
		64: {
			Summary: "you may go to hell; I am going to the frontier",
			Effects: []Effect{transfer("Colonist", "", 0)},
		},
		65: {Summary: "extremely talented in your field", Effects: []Effect{talentedInYourField()}},
		66: {
			Summary: "excellent work",
			Effects: []Effect{advance(), benefitRolls(0, 1, ScopeCareer)},
		},
	}

	lifeEventRows(table)

	return table
}

func politicianTurnedOn() []Effect {
	return []Effect{chr("CHA", -2), relationship(Rival, 1, "")}
}
