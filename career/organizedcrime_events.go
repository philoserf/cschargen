package career

// organizedCrimeEvents is the d66 table of pp. 248-249.
//
// Four of its rows are a prison term the character may come back from, and
// two of them turn on a single skill throw with nothing in between.
func organizedCrimeEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "records to maintain for the boss",
			Effects: []Effect{skill("Admin")},
		},
		13: {
			Summary: "the boss wants to know whether you have been skimming",
			Effects: []Effect{checkSkill("Broker", 8,
				[]Effect{throwModifier("next advancement roll", 2)},
				[]Effect{throwModifier("next survival roll", -2)})},
		},
		14: {
			Summary: "brought in on a negotiation with another organization",
			Effects: []Effect{checkSkill("Persuade", 8,
				[]Effect{
					pickSkill("Etiquette", "Diplomat"),
					throwModifier("next two advancement rolls", 2),
				},
				[]Effect{throwModifier("next two advancement rolls", -2)})},
		},
		15: {
			Summary: "a bomb is sometimes the best way to annoy the opposition",
			Effects: []Effect{skill("Explosives")},
		},
		16: {
			Summary: "a member of a rival organization, captured",
			Effects: []Effect{
				relationship(Enemy, 1, ""),
				checkSkill("Interrogation", 8,
					[]Effect{benefitRolls(2, 0, ScopeBatch)},
					[]Effect{chr("CHA", -2)}),
			},
		},
		21: {
			Summary: "law enforcement is onto you",
			Effects: []Effect{checkSkill("Advocate", 8,
				[]Effect{throwModifier("next advancement roll", 2)},
				[]Effect{transfer("Prisoner", "Prisoner", 1)})},
		},
		22: {
			Summary: "the family loves you",
			Effects: []Effect{pick("choose what their affection is worth",
				opt("Carouse", skill("Carouse")),
				opt("contacts", relationship(Contact, 1, "1d3")))},
		},
		23: {
			Summary: "sometimes the getaway driver",
			Effects: []Effect{pick("choose Drive (Wheeled) or Flyer (Grav)",
				opt("Drive (Wheeled)", skill("Drive", "Wheeled")),
				opt("Flyer (Grav)", skill("Flyer", "Grav")))},
		},
		24: {
			Summary: "the boss has decided you are a degenerate gambler",
			Effects: []Effect{
				skill("Gambler"),
				throwModifier("next two advancement rolls", -2),
			},
		},
		25: {
			Summary: "toughs from a rival organization corner you",
			Effects: []Effect{checkSkill("Melee", 8,
				[]Effect{chr("CHA", 2), benefitRolls(1, 0, ScopeBatch)},
				[]Effect{chr("CHA", -2), throwModifier("next advancement roll", -2)})},
		},
		26: {
			Summary: "an ear kept to the streets",
			Effects: []Effect{skill("Streetwise")},
		},
		41: {
			Summary: "time away from the grind of the city",
			Effects: []Effect{skill("Animals", "Any")},
		},
		42: {
			Summary: "time on a spaceship for the boss",
			Effects: []Effect{pick("choose the berth you worked",
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Gunner (Turrets)", skill("Gunner", "Turrets")),
				opt("Pilot", skill("Pilot", "Any")))},
		},
		43: {
			Summary: "keeping the boss's office clear of surveillance",
			Effects: []Effect{skill("Electronics", "Any")},
		},
		44: {
			Summary: "a shootout with another gang",
			Effects: []Effect{checkSkill("Gun Combat", 8,
				[]Effect{throwModifier("next advancement roll", 2)},
				[]Effect{injury(1)})},
		},
		45: {
			Summary: "wits, used to talk your way out of things",
			Effects: []Effect{skill("Persuade")},
		},
		46: {
			Summary: "the opposition attempts to get the drop on you",
			Effects: []Effect{checkSkill("Tactics", 8,
				[]Effect{pickSkill("Tactics", "Recon")},
				[]Effect{injury(1)})},
		},
		51: {
			Summary: "training",
			Effects: []Effect{pick("choose what the training built",
				opt("STR", chr("STR", 1)),
				opt("DEX", chr("DEX", 1)),
				opt("END", chr("END", 1)),
				opt("Athletics", skill("Athletics", "Any")))},
		},
		52: {
			Summary: "a job to steal a precious jewel",
			Effects: []Effect{checkSkill("Deception", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch)},
				[]Effect{transfer("Prisoner", "Prisoner", 2)})},
		},
		53: {
			Summary: "a bomb, for the underside of a rival boss's grav vehicle",
			Effects: []Effect{checkSkill("Explosives", 8,
				[]Effect{skill("Explosives"), benefitRolls(2, 0, ScopeBatch)},
				[]Effect{injury(2), skill("Explosives")})},
		},
		54: {
			Summary: "awareness of your surroundings",
			Effects: []Effect{skill("Recon")},
		},
		55: {
			Summary: "a shootout with police",
			Effects: []Effect{checkSkill("Gun Combat", 8,
				[]Effect{throwModifier("next advancement roll", 2)},
				[]Effect{transfer("Prisoner", "Prisoner", 1)})},
		},
		56: {
			Summary: "time on a starship recently",
			Effects: []Effect{pick("choose Suit (Vacc Suit) or Survival (Freefall)",
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
				opt("Survival (Freefall)", skill("Survival", "Freefall")))},
		},
		61: {
			Summary: "a hobby taken up",
			Effects: []Effect{pickSkill("Art", "Science")},
		},
		62: {
			Summary: "side work, outside what the family sanctions",
			Effects: []Effect{checkSkill("Deception", 8,
				[]Effect{benefitRolls(3, 0, ScopeBatch)},
				[]Effect{
					loseRank(1),
					throwModifier("advancement rolls for the rest of this career", -2),
				})},
		},
		63: {
			Summary: "a variety of jobs",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		64: {
			Summary: "the need to be somewhat sneaky",
			Effects: []Effect{skill("Stealth")},
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
