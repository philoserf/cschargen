package career

// investigatorEvents is the d66 table of pp. 218-219.
func investigatorEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {Summary: "paperwork, and so much of it", Effects: []Effect{skill("Admin")}},
		13: {
			Summary: "a fund for current and former members of the profession, to invest wisely",
			Effects: []Effect{
				skill("Broker"),
				checkSkill("Broker", 8,
					[]Effect{
						relationship(Contact, 1, "1d3"),
						relationship(Ally, 1, ""),
						benefitRolls(2, 0, ScopeBatch),
					},
					[]Effect{relationship(Enemy, 2, ""), benefitRolls(-2, 0, ScopeBatch)}),
			},
		},
		14: {
			Summary: "a lot of surveillance",
			Effects: []Effect{pick("choose what the watching taught you",
				opt("Electronics (Computers)", skill("Electronics", "Computers")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Stealth", skill("Stealth")))},
		},
		15: {
			Summary: "getting information out of suspects",
			Effects: []Effect{pick("choose how you got it out of them",
				opt("Carouse", skill("Carouse")),
				opt("Interrogation (Questioning)", skill("Interrogation", "Questioning")),
				opt("Persuade", skill("Persuade")))},
		},
		16: {
			Summary: "learning to work many angles",
			Effects: []Effect{rollOtherAssignment()},
		},
		21: {
			// Two throws deep: the journalist, then the investigation the
			// journalist opens.
			Summary: "an intrepid journalist hounding you as corrupt",
			Effects: []Effect{
				relationship(Rival, 1, ""),
				checkSkill("Advocate", 8,
					[]Effect{chr("CHA", 1)},
					[]Effect{checkChr("CHA", 8,
						[]Effect{throwModifier("next advancement roll", 2)},
						[]Effect{
							autoFailure(advancementThrow, "suspended: the next advancement roll fails automatically"),
							benefitRolls(-2, 0, ScopeBatch),
						})}),
			},
		},
		22: {
			Summary: "a drink after a hard day",
			Effects: []Effect{relationship(Contact, 1, ""), skill("Carouse")},
		},
		23: {
			Summary: "a case that puts you aboard a starship as crew",
			Effects: []Effect{pick("choose the berth the case put you in",
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Gunner (Turrets)", skill("Gunner", "Turrets")),
				opt("Pilot", skill("Pilot", "Any")),
				opt("Survival (Freefall)", skill("Survival", "Freefall")),
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")))},
		},
		24: {
			Summary: "a lot of time at the range",
			Effects: []Effect{skill("Gun Combat", "Any")},
		},
		25: {
			Summary: "a cold case looked into",
			Effects: []Effect{checkSkill("Investigate", 8,
				[]Effect{throwModifier("next advancement roll", 2), relationship(Contact, 1, "")},
				[]Effect{relationship(Contact, 1, "")})},
		},
		26: {
			Summary: "a trap set for you by a group of criminals",
			Effects: []Effect{checkSkill("Gun Combat", 8,
				[]Effect{throwModifier("next advancement roll", 2), chr("CHA", 1)},
				[]Effect{injury(2)})},
		},
		41: {
			Summary: "riding an animal, taken up as a hobby",
			Effects: []Effect{skill("Animals", "Riding")},
		},
		42: {
			Summary: "knowing how criminals work in order to investigate them",
			Effects: []Effect{pickSkill("Deception", "Streetwise")},
		},
		43: {
			Summary: "a variety of people from different cultures",
			Effects: []Effect{skill("Language", "Any")},
		},
		44: {
			Summary: "working on your own equipment",
			Effects: []Effect{skill("Mechanic")},
		},
		45: {
			Summary: "a first aid course",
			Effects: []Effect{skill("Medic", "First Aid")},
		},
		46: {
			Summary: "always being aware of your surroundings",
			Effects: []Effect{skill("Recon")},
		},
		51: {
			Summary: "an amateur sports league joined",
			Effects: []Effect{
				relationship(Contact, 1, ""),
				pick("choose what the league built",
					opt("STR", chr("STR", 1)),
					opt("DEX", chr("DEX", 1)),
					opt("END", chr("END", 1)),
					opt("Athletics", skill("Athletics", "Any"))),
			},
		},
		52: {
			Summary: "arguments with the establishment, smoothed over",
			Effects: []Effect{pickSkill("Diplomat", "Etiquette")},
		},
		53: {
			Summary: "a gambling circle in your peer group",
			Effects: []Effect{gamblingCircle("Persuade")},
		},
		54: {Summary: `you are the "guv"`, Effects: []Effect{skill("Leadership")}},
		55: {
			Summary: "a martial arts course",
			Effects: []Effect{pick("choose Athletics or Melee (Unarmed)",
				opt("Athletics", skill("Athletics", "Any")),
				opt("Melee (Unarmed)", skill("Melee", "Unarmed Combat")))},
		},
		56: {
			Summary: "a case in a wild and remote part of your world",
			Effects: []Effect{pickSkill("Navigation", "Survival")},
		},
		61: {
			Summary: "a hobby taken up",
			Effects: []Effect{pickSkill("Art", "Science")},
		},
		62: {
			Summary: "a vehicle, which is an important tool in this work",
			Effects: []Effect{pick("choose Drive (Wheeled) or Flyer (Grav)",
				opt("Drive (Wheeled)", skill("Drive", "Wheeled")),
				opt("Flyer (Grav)", skill("Flyer", "Grav")))},
		},
		63: {
			Summary: "a variety of jobs, done because you must",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		64: {
			Summary: "a high-profile case between two of the system's prominent families",
			Effects: []Effect{
				checkSkill("Investigate", 8,
					[]Effect{chr("CHA", 2), benefitRolls(3, 0, ScopeBatch)},
					[]Effect{chr("CHA", -2), throwModifier("next advancement roll", -2)}),
				relationship(Enemy, 1, ""),
			},
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
