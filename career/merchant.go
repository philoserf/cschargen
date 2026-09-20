package career

// IndependentMerchant is the career of pp. 207-211: crewing a free trader
// rather than a corporation's freighter.
//
// Corporate Shipper's mishap 4 -- "You're going independent" -- sends
// characters here, and it was the last stubbed destination in the book.
func IndependentMerchant() Career {
	return Career{
		Name:         "Independent Merchant",
		Cite:         "pp. 207-211",
		Enlistment:   &Check{Characteristic: "INT", Number: 7},
		MishapEjects: true,
		Assignments: []Assignment{
			{
				Name:        "Crew",
				Description: "the crew of an independent merchant vessel",
				Survival:    Check{Characteristic: "INT", Number: 7},
				Advancement: Check{Characteristic: "EDU", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Crew", Rows: [6]Effect{
					skill("Electronics", "Any"), skill("Pilot", "Any"), skill("Astrogation"),
					skill("Carouse"), skill("Persuade"), skill("Jack of All Trades"),
				}},
				Ranks: [][]Effect{
					nil,
					{skill("Persuade")},
					{skill("Broker")},
					{skill("Etiquette")},
					nil,
					{skill("Leadership")},
					{skill("Jack of All Trades")},
				},
			},
			{
				Name:        "Engineer",
				Description: "the engineering crew of an independent merchant vessel",
				Survival:    Check{Characteristic: "INT", Number: 7},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Engineer", Rows: [6]Effect{
					skill("Admin"), skill("Electronics", "Any"), skill("Engineer", "Any"),
					skill("Mechanic"), skill("Science", "Any"), skill("Jack of All Trades"),
				}},
				Ranks: [][]Effect{
					{skill("Mechanic")},
					{skill("Engineer", "M-Drive")},
					{skill("Engineer", "Z-Drive")},
					nil,
					{skill("Science", "Physics")},
					nil,
					{skill("Jack of All Trades")},
				},
			},
			{
				Name:        "Gunnery",
				Description: "the defence of an independent merchant vessel",
				Survival:    Check{Characteristic: "DEX", Number: 7},
				Advancement: Check{Characteristic: "EDU", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Gunnery", Rows: [6]Effect{
					skill("Mechanic"), skill("Electronics", "Sensors"), skill("Gunner", "Turrets"),
					skill("Gun Combat", "Any"), skill("Tactics", "Naval"), skill("Discipline"),
				}},
				Ranks: [][]Effect{
					nil,
					{skill("Gunner", "Turrets")},
					{skill("Electronics", "Sensors")},
					nil,
					{skill("Discipline")},
					nil,
					{skill("Jack of All Trades")},
				},
			},
			{
				Name:        "Flight",
				Description: "piloting the vessel or the small craft it carries",
				Survival:    Check{Characteristic: "DEX", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Flight", Rows: [6]Effect{
					skill("Astrogation"), skill("Electronics", "Any"), skill("Pilot", "Any"),
					skill("Flyer", "Any"), skill("Tactics", "Naval"), skill("Discipline"),
				}},
				Ranks: [][]Effect{
					nil,
					{skill("Pilot", "Spacecraft")},
					{skill("Astrogation")},
					nil,
					{skill("Discipline")},
					nil,
					{skill("Jack of All Trades")},
				},
			},
		},
		Tables: []SkillTable{
			personalDevelopment(),
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Admin"), skill("Suit", "Vacc Suit"), skill("Broker"),
				skill("Persuade"), skill("Carouse"), skill("Etiquette"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Art", "Any"), skill("Diplomat"), skill("Language", "Any"),
				skill("Medic", "First Aid"), skill("Science", "Any"), skill("Leadership"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 1000, Other: weaponOrItsUse()},
			{Cash: 5000, Other: chr("CHA", 1)},
			{Cash: 7000, Other: chr("DEX", 1)},
			{Cash: 10000, Other: chr("INT", 1)},
			{Cash: 20000, Other: relationship(Contact, 1, "")},
			{Cash: 40000, Other: relationship(Ally, 1, "")},
			{Cash: 75000, Other: stashValued("a Captain's Guild membership for one year", "625000")},
		},
		Mishaps: merchantMishaps(),
		Events:  merchantEvents(),
	}
}

func merchantMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "accused of negligence resulting in a crewmate's death",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "a drive accident strands you in a system with only a minor colony",
			Effects: []Effect{
				pick("what surviving there taught you",
					opt("Survival", skill("Survival", "Any")),
					opt("Animals", skill("Animals", "Any")),
					opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
					opt("Navigation", skill("Navigation"))),
				transfer("Colonist", "", 0),
			},
		},
		{Summary: "accidental exposure to a dangerous atmosphere", Effects: []Effect{chr("END", -1)}},
		{
			Summary: "the ship irreparably damaged in a fight with pirates, and sold for scrap",
			Effects: []Effect{injury(2), benefitRolls(3, 0, ScopeBatch)},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "a rivalry aboard ship, and a fit of anger that ends your service",
			Effects: []Effect{relationship(Rival, 1, "")},
		},
		{Summary: "a shuttle transfer ends in a crash", Effects: []Effect{injury(1)}},
		{Summary: "the nothingness, the cramped quarters, the uncertain pay", Effects: nil},
		{Summary: "a life-threatening illness", Effects: []Effect{chr("END", -2)}},
		{
			Summary: "an act of piracy, and no free trader will have you",
			Effects: []Effect{unimplemented(
				"roll 1d6: on 1 two terms in the Prisoner career; on 2-5 the Pirate career, " +
					"retaining rank; on 6 marooned on a low-port world and into the Vagabond career")},
		},
	}
}
