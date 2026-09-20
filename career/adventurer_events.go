package career

// adventurerEvents is the d66 table of pp. 148-149.
func adventurerEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "caught in a violent weather event",
			Effects: []Effect{checkSkill("Survival", 8, nil, []Effect{injury(1)})},
		},
		13: {Summary: "records kept of your adventures", Effects: []Effect{skill("Admin")}},
		14: {
			Summary: "an animal befriended on the latest expedition",
			Effects: []Effect{stashItem("a pet"), skill("Animals", "Any")},
		},
		15: {
			Summary: "back in colonized space, and time to party with old friends",
			Effects: []Effect{skill("Carouse")},
		},
		16: {
			Summary: "a term spent as another type of adventurer",
			Effects: []Effect{rollOtherAssignment()},
		},
		21: {
			Summary: "the political landscape at home has changed, and you take a stand",
			Effects: []Effect{skill("Advocate", "Politics", "Oratory")},
		},
		22: {Summary: "some advanced training", Effects: []Effect{rollTable(AdvancedEducation)}},
		23: {
			Summary: "downtime spent studying",
			Effects: []Effect{pick("what you studied",
				opt("+1 EDU", chr("EDU", 1)),
				opt("Art", skill("Art", "Any")),
				opt("Science", skill("Science", "Any")))},
		},
		24: {
			Summary: "lying, cheating or stealing to make ends meet",
			Effects: []Effect{pick("how you made them meet",
				opt("Deception (Forgery)", skill("Deception", "Forgery")),
				opt("Deception (Lie)", skill("Deception", "Lie")),
				opt("Deception (Pickpocket)", skill("Deception", "Pickpocket")),
				opt("Streetwise", skill("Streetwise")))},
		},
		25: {Summary: "your own cook", Effects: []Effect{skill("Chef")}},
		26: {
			Summary: "growing proficient with a weapon",
			Effects: []Effect{pick("which weapon",
				opt("Draw", skill("Draw")),
				opt("Gun Combat", skill("Gun Combat", "Any")),
				opt("Melee", skill("Melee", "Blade", "Bludgeon", "Bow")))},
		},
		41: {
			// The Persuade branch is a 10+ check rather than the usual 8+.
			Summary: "a pirate base, stumbled upon",
			Effects: []Effect{
				aPirateBase(
					[]Option{opt("Persuade at 10+", checkSkill("Persuade", 10, nil,
						[]Effect{checkSkill("Gun Combat", 8,
							[]Effect{skill("Tactics", "Military")}, []Effect{injury(1)})}))},
					checkSkill("Gun Combat", 8,
						[]Effect{skill("Tactics", "Military")}, []Effect{injury(1)})),
				relationship(Enemy, 1, ""),
			},
		},
		42: {Summary: "a memoir of your exploits", Effects: []Effect{skill("Art", "Writing")}},
		43: {
			Summary: "smoothing things over after trespassing on owned ground",
			Effects: []Effect{pickSkill("Diplomat", "Persuade")},
		},
		44: {
			Summary: "time spent on worlds with dangerous atmospheres",
			Effects: []Effect{skill("Suit", "Vacc Suit")},
		},
		45: {Summary: "someone else intruding on your discoveries", Effects: []Effect{relationship(Rival, 1, "")}},
		46: {
			Summary: "lost one time too many",
			Effects: []Effect{injury(1), skill("Navigation")},
		},
		51: {
			Summary: "learning to operate your own ship",
			Effects: []Effect{pick("what you learned",
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Pilot (Spacecraft)", skill("Pilot", "Spacecraft")))},
		},
		52: {
			Summary: "maintaining and repairing your own equipment",
			Effects: []Effect{pickSkill("Electronics", "Engineer", "Mechanic")},
		},
		53: {
			Summary: "a crash out in the wilds",
			Effects: []Effect{injury(1), pickSkill("Drive", "Flyer")},
		},
		54: {
			Summary: "a conference on exploring the frontier",
			Effects: []Effect{
				unimplemented("gain 1d6-3 Contacts"),
				relationship(Rival, 1, ""),
				chr("CHA", 1),
			},
		},
		55: {
			Summary: "hired to take someone important on an adventure holiday",
			Effects: []Effect{benefitRolls(3, 1, ScopeBatch), relationship(Contact, 1, "")},
		},
		56: {
			Summary: "you keep getting better at this",
			Effects: []Effect{autoSuccess("next advancement roll")},
		},
		61: {Summary: "learning to sell what the expeditions bring back", Effects: []Effect{skill("Broker")}},
		62: {
			Summary: "a new colony, and a colonist who wants you to stay",
			Effects: []Effect{pick("stay or move on",
				opt("stay", relationship(Ally, 1, ""), transfer("Colonist", "", 0)),
				opt("move on", relationship(Contact, 1, "")))},
		},
		63: {
			Summary: "on the frontier you do most things for yourself",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		64: {
			Summary: "it belongs in a museum",
			Effects: []Effect{
				benefitRolls(-2, 0, ScopeBatch),
				relationship(Contact, 1, ""),
				advance(),
			},
		},
		65: {
			Summary: "a world discovered that would be perfect for a colony",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch), advance()},
		},
		66: {
			Summary: "the tales of your exploits have made you famous",
			Effects: []Effect{transfer("Celebrity", "", 0)},
		},
	}

	lifeEventRows(table)

	return table
}
