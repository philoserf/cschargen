package career

// marineEvents is the d66 table of pp. 227-228.
func marineEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "after-action reports, written often",
			Effects: []Effect{pick("what the writing taught you",
				opt("Admin", skill("Admin")),
				opt("Art (Writing)", skill("Art", "Writing")))},
		},
		13: {
			Summary: "the seedy side of life",
			Effects: []Effect{pick("what it taught you",
				opt("Deception", skill("Deception", "Any")),
				opt("Streetwise", skill("Streetwise")))},
		},
		14: {Summary: "the one the squad relies on to place the explosives", Effects: []Effect{skill("Explosives")}},
		15: {
			Summary: "a heavy weapon called on to defend a ship, port or base",
			Effects: []Effect{checkSkill("Heavy Weapons", 8,
				[]Effect{skill("Heavy Weapons", "Any"), throwModifier("next advancement roll", 2)},
				[]Effect{injury(1)})},
		},
		16: {
			Summary: "hand-to-hand combat during a mission",
			Effects: []Effect{checkSkill("Melee", 8,
				[]Effect{skill("Melee", "Any")}, []Effect{injury(1)})},
		},
		21: {
			Summary: "top physical shape",
			Effects: []Effect{pick("what the training gave you",
				opt("+1 STR", chr("STR", 1)),
				opt("+1 DEX", chr("DEX", 1)),
				opt("+1 END", chr("END", 1)),
				opt("Athletics", skill("Athletics", "Any")))},
		},
		22: {Summary: "practised at drawing your weapon", Effects: []Effect{skill("Draw")}},
		23: {
			Summary: "a group of gamblers",
			Effects: []Effect{checkSkill("Gambler", 8,
				[]Effect{skill("Gambler"), benefitRolls(2, 0, ScopeBatch)},
				[]Effect{benefitRolls(-3, 0, ScopeBatch)})},
		},
		24: {Summary: "chosen for cross-training", Effects: []Effect{rollOtherAssignment()}},
		25: {
			Summary: "tasked with keeping everyone on course",
			Effects: []Effect{pick("what the navigating taught you",
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Navigation", skill("Navigation")))},
		},
		26: {Summary: "quite sneaky", Effects: []Effect{skill("Stealth")}},
		51: {
			Summary: "downtime spent on your hobbies",
			Effects: []Effect{pick("what you studied",
				opt("Art", skill("Art", "Any")),
				opt("Science", skill("Science", "Any")))},
		},
		52: {
			Summary: "tasked with providing transportation",
			Effects: []Effect{pick("what you drove",
				opt("Drive", skill("Drive", "Any")),
				opt("Flyer", skill("Flyer", "Any")),
				opt("Pilot (Small Craft)", skill("Pilot", "Small Craft")),
				opt("Seafarer", skill("Seafarer", "Any")))},
		},
		53: {
			Summary: "a serious firefight during a mission",
			Effects: []Effect{checkSkill("Gun Combat", 8,
				[]Effect{
					throwModifier("next survival roll", 2),
					throwModifier("next advancement roll", 2),
				},
				[]Effect{injury(1), relationship(Ally, 1, "")})},
		},
		54: {
			Summary: "chosen for additional education",
			Effects: []Effect{pick("take a commission or take the schooling",
				opt("become an officer", commission(0)),
				opt("roll twice on Advanced Education",
					rollTable(AdvancedEducation), rollTable(AdvancedEducation)))},
		},
		55: {Summary: "learning to spot threats and dangers", Effects: []Effect{skill("Recon")}},
		56: {
			Summary: "trained to survive in a variety of places",
			Effects: []Effect{skill("Survival", "Any")},
		},
		61: {
			Summary: "quite popular",
			Effects: []Effect{pick("what the popularity brought",
				opt("+1 CHA", chr("CHA", 1)),
				opt("contacts", relationship(Contact, 0, "1d3")),
				opt("Carouse", skill("Carouse")))},
		},
		62: {Summary: "polite, until it is time to not be polite", Effects: []Effect{skill("Etiquette")}},
		63: {
			Summary: "time aboard ship, spent learning to run one",
			Effects: []Effect{pick("what you learned",
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Gunner (Turrets)", skill("Gunner", "Turrets")),
				opt("Mechanic", skill("Mechanic")),
				opt("Pilot", skill("Pilot", "Any")),
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
				opt("Survival (Freefall)", skill("Survival", "Freefall")))},
		},
		64: {Summary: "trained in first aid", Effects: []Effect{skill("Medic", "First Aid")}},
		65: {
			Summary: "extremely talented in your field",
			Effects: []Effect{pick("what the talent earned",
				opt("a level in a skill you already have",
					unimplemented("raise a skill the character already holds")),
				opt("a promotion", Effect{Kind: EffectRank, Detail: "gain a rank"}))},
		},
		66: {
			Summary: "excellent work",
			Effects: []Effect{advance(), benefitRolls(0, 1, ScopeCareer)},
		},
	}

	for result := 31; result <= 36; result++ {
		table[result] = EventRow{Summary: "life event", Effects: []Effect{lifeEvent()}}
	}

	for result := 41; result <= 46; result++ {
		table[result] = EventRow{
			Summary: "military event",
			Effects: []Effect{{Kind: EffectMilitaryEvent, Detail: "roll on the Military Events table"}},
		}
	}

	return table
}
