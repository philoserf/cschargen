package career

// scavengerEvents is the d66 table of pp. 267-268.
//
// Result 41 prints "Gain a level in Advocate (Legal) 8+" (ERRATA E-14): there
// is no throw in it, and the "8+" is stray.
func scavengerEvents() EventTable {
	// Two rows -- a painting at 13, a book at 42 -- are the same result over
	// a different specialty: recognise it and it is worth three benefit
	// rolls, miss it and it is one more piece in the collection.
	find := func(specialty string) []Effect {
		return []Effect{checkSkill("Art", 8,
			[]Effect{skill("Art", specialty), benefitRolls(3, 0, ScopeBatch)},
			[]Effect{benefitRolls(1, 0, ScopeBatch)})}
	}

	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "law enforcement asks where certain items came from",
			Effects: []Effect{checkSkill("Admin", 8, nil,
				[]Effect{unimplemented("the police confiscate the collection: " +
					"lose every benefit roll gained before this event")})},
		},
		13: {Summary: "a painting obtained", Effects: find("Painting")},
		14: {
			Summary: "a buyer for something in the collection",
			Effects: []Effect{checkSkill("Broker", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch)},
				[]Effect{benefitRolls(-2, 0, ScopeBatch)})},
		},
		15: {
			Summary: "a commission to forge the signatures of Earth celebrities",
			Effects: []Effect{pick("forge them, or turn the buyer in",
				opt("forge them", checkSkill("Deception", 8,
					[]Effect{relationship(Ally, 1, ""), benefitRolls(3, 0, ScopeBatch)},
					[]Effect{transfer("Prisoner", "Prisoner", 1)})),
				opt("turn them in",
					throwModifier("next advancement roll", 2),
					relationship(Enemy, 1, "")))},
		},
		16: {
			Summary: "a collection of vintage computers, and a collector who wants it",
			Effects: []Effect{checkSkill("Broker", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch), skill("Electronics", "Computers")},
				[]Effect{benefitRolls(-1, 0, ScopeBatch)})},
		},
		21: {
			Summary: "a buyer for something in the depths of the warehouse",
			Effects: []Effect{checkSkill("Admin", 8, []Effect{credits("1d6x10000")}, nil)},
		},
		22: {Summary: "it pays to be aware of the arts", Effects: []Effect{skill("Art", "Any")}},
		23: {
			Summary: "a popular person",
			Effects: []Effect{pick("choose what the popularity is worth",
				opt("Carouse", skill("Carouse")),
				opt("contacts", relationship(Contact, 1, "1d3")))},
		},
		24: {
			Summary: "something shady, done in order to succeed",
			Effects: []Effect{skill("Deception", "Any")},
		},
		25: {
			Summary: "electronics, which you have to understand to deal in them",
			Effects: []Effect{skill("Electronics", "Any")},
		},
		26: {
			// The table prints "Engineering (Any)"; the skill list (p. 304)
			// has Engineer. See ERRATA E-14.
			Summary: "an old starship, torn apart or put back together",
			Effects: []Effect{skill("Engineer", "Any")},
		},
		41: {
			Summary: "the local laws and regulations your business runs under",
			Effects: []Effect{skill("Advocate", "Legal")},
		},
		42: {Summary: "an old book obtained", Effects: find("Writing")},
		43: {
			Summary: "a lot of arguments in the community, smoothed over",
			Effects: []Effect{skill("Diplomat")},
		},
		44: {Summary: "a first aid class", Effects: []Effect{skill("Medic", "First Aid")}},
		45: {
			Summary: "older grav vehicles collected, sold or restored",
			Effects: []Effect{pick("choose Flyer (Grav) or Mechanic",
				opt("Flyer (Grav)", skill("Flyer", "Grav")),
				opt("Mechanic", skill("Mechanic")))},
		},
		46: {
			Summary: "people of varied backgrounds, met in this business",
			Effects: []Effect{skill("Language", "Any")},
		},
		51: {
			Summary: "animals, obtained as part of the collection",
			Effects: []Effect{skill("Animals", "Any")},
		},
		52: {
			Summary: "a starship berth worked for passage to a piece worth having",
			Effects: []Effect{pick("choose the berth you worked",
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Gunner (Turrets)", skill("Gunner", "Turrets")),
				opt("Mechanic", skill("Mechanic")),
				opt("Pilot", skill("Pilot", "Any")))},
		},
		53: {
			Summary: "rare automobiles collected, sold or restored",
			Effects: []Effect{pick("choose Drive (Wheeled) or Mechanic",
				opt("Drive (Wheeled)", skill("Drive", "Wheeled")),
				opt("Mechanic", skill("Mechanic")))},
		},
		54: {
			Summary: "a self-defense course",
			Effects: []Effect{pickSkill("Gun Combat", "Melee")},
		},
		55: {
			Summary: "a lot of time spent in space recently",
			Effects: []Effect{pick("choose Suit (Vacc Suit) or Survival (Freefall)",
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
				opt("Survival (Freefall)", skill("Survival", "Freefall")))},
		},
		56: {
			Summary: "the provenance of an item, which it is important to learn",
			Effects: []Effect{pick("choose Investigate or Science (History)",
				opt("Investigate", skill("Investigate")),
				opt("Science (History)", skill("Science", "History")))},
		},
		61: {
			Summary: "training",
			Effects: []Effect{pick("choose what the training built",
				opt("STR", chr("STR", 1)),
				opt("DEX", chr("DEX", 1)),
				opt("END", chr("END", 1)),
				opt("Athletics", skill("Athletics", "Any")))},
		},
		62: {
			Summary: "if there is a deal to be made, you are the one to make it",
			Effects: []Effect{skill("Broker")},
		},
		63: {
			Summary: "a holovid about your exploits, and a stunning success",
			Effects: []Effect{chr("CHA", 2), celebrityAsStar()},
		},
		64: {
			Summary: "a number of odd jobs",
			Effects: []Effect{skill("Jack of All Trades")},
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
