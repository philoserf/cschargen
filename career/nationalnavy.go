package career

// NationalNavy is the career of pp. 233-241: service in the navy of one of
// the sector's nations.
//
// It is the first career with a commission, so it carries two rank tables
// per assignment and an Officer Skills table (p. 117), and its events route
// to the Military Events table (p. 121).
func NationalNavy() Career {
	return Career{
		Name:         "National Navy",
		Cite:         "pp. 233-241",
		Enlistment:   &Check{Characteristic: "INT", Number: 7},
		Commission:   &Check{Characteristic: "EDU", Number: 8},
		MishapEjects: true,
		Assignments:  navyAssignments(),
		Tables:       navyTables(),
		Benefits:     navyBenefits(),
		Mishaps:      navyMishaps(),
		Events:       navyEvents(),
	}
}

func navyAssignments() []Assignment {
	return []Assignment{
		navyAssignment("Crew/Line", "a crew member or officer in a national navy",
			Check{Characteristic: "INT", Number: 7}, Check{Characteristic: "EDU", Number: 7},
			[6]Effect{
				skill("Admin"), skill("Astrogation"), skill("Pilot", "Any"),
				skill("Electronics", "Any"), skill("Mechanic"), skill("Gunner", "Any"),
			},
			[7][]Effect{
				{skill("Electronics", "Any")},
				{skill("Discipline")},
				nil,
				{skill("Astrogation")},
				nil,
				{skill("Pilot", "Spacecraft")},
				{skill("Jack of All Trades")},
			},
			[7][]Effect{
				{skill("Admin")},
				{skill("Electronics", "Any")},
				nil,
				{skill("Tactics", "Naval")},
				{skill("Astrogation")},
				{skill("Leadership")},
				{skill("Diplomat")},
			}),
		navyAssignment("Engineer", "an engineer maintaining a naval vessel's engines",
			Check{Characteristic: "INT", Number: 7}, Check{Characteristic: "EDU", Number: 8},
			[6]Effect{
				skill("Electronics", "Any"), skill("Mechanic"), skill("Engineer", "Any"),
				skill("Engineer", "Any"), skill("Science", "Any"), skill("Jack of All Trades"),
			},
			[7][]Effect{
				{skill("Mechanic")},
				{skill("Discipline")},
				nil,
				{skill("Engineer", "M-Drive")},
				nil,
				{skill("Engineer", "Power")},
				{skill("Engineer", "Life Support")},
			},
			[7][]Effect{
				{skill("Admin")},
				{skill("Mechanic")},
				{skill("Engineer", "Z-Drive")},
				{skill("Engineer", "M-Drive")},
				{skill("Engineer", "Life Support")},
				nil,
				{skill("Leadership")},
			}),
		navyAssignment("Gunnery", "a weapons specialist on a warship",
			Check{Characteristic: "DEX", Number: 7}, Check{Characteristic: "EDU", Number: 7},
			[6]Effect{
				skill("Mechanic"), skill("Electronics", "Any"), skill("Gunner", "Any"),
				skill("Gunner", "Any"), skill("Tactics", "Naval"), skill("Gun Combat", "Any"),
			},
			[7][]Effect{
				{skill("Gunner", "Turrets")},
				{skill("Discipline")},
				nil,
				{skill("Electronics", "Sensors")},
				nil,
				{skill("Mechanic")},
				nil,
			},
			[7][]Effect{
				{skill("Admin")},
				{skill("Gunner", "Any")},
				nil,
				{skill("Tactics", "Naval")},
				nil,
				{skill("Electronics", "Sensors")},
				{skill("Leadership")},
			}),
		navyAssignment("Flight", "a pilot of a fighter or shuttle",
			Check{Characteristic: "DEX", Number: 8}, Check{Characteristic: "EDU", Number: 7},
			[6]Effect{
				skill("Tactics", "Naval"), skill("Astrogation"), skill("Pilot", "Any"),
				skill("Pilot", "Any"), skill("Electronics", "Any"), skill("Flyer", "Any"),
			},
			[7][]Effect{
				{skill("Mechanic")},
				{skill("Discipline")},
				nil,
				{skill("Electronics", "Sensors")},
				nil,
				{skill("Pilot", "Small Craft")},
				nil,
			},
			[7][]Effect{
				{skill("Admin")},
				{skill("Pilot", "Small Craft")},
				nil,
				{skill("Electronics", "Sensors")},
				nil,
				{skill("Pilot", "Spacecraft")},
				{skill("Leadership")},
			}),
	}
}

// navyAssignment builds one of the four, which differ only in their tables.
func navyAssignment(
	name, description string,
	survival, advancement Check,
	skills [6]Effect,
	enlisted, officer [7][]Effect,
) Assignment {
	ranks := officer

	return Assignment{
		Name:         name,
		Description:  description,
		Survival:     survival,
		Advancement:  advancement,
		Skills:       SkillTable{Kind: AssignmentSkills, Name: name, Rows: skills},
		Ranks:        enlisted,
		OfficerRanks: &ranks,
	}
}

func navyTables() []SkillTable {
	return []SkillTable{
		{
			Kind: PersonalDevelopment,
			Name: "Personal Development",
			Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), skill("Athletics", "Any"),
			},
		},
		{
			Kind: ServiceSkills,
			Name: "Service Skills",
			Rows: [6]Effect{
				skill("Carouse"), skill("Discipline"), skill("Survival", "Any"),
				skill("Gun Combat", "Any"), skill("Melee", "Any"), skill("Suit", "Vacc Suit"),
			},
		},
		{
			Kind:       AdvancedEducation,
			Name:       "Advanced Education",
			MinimumEDU: 8,
			Rows: [6]Effect{
				skill("Admin"), skill("Advocate", "Any"), skill("Science", "Any"),
				skill("Medic", "Any"), skill("Language", "Any"), skill("Etiquette"),
			},
		},
		{
			Kind: OfficerSkills,
			Name: "Officer Skills",
			Rows: [6]Effect{
				skill("Admin"), skill("Etiquette"), skill("Tactics", "Naval"),
				skill("Diplomat"), skill("Persuade"), skill("Leadership"),
			},
		},
	}
}

func navyBenefits() [7]BenefitRow {
	return [7]BenefitRow{
		{Cash: 250, Other: chr("END", 1)},
		{Cash: 500, Other: chr("DEX", 1)},
		{Cash: 1000, Other: chr("INT", 1)},
		{Cash: 2000, Other: chr("EDU", 1)},
		{Cash: 5000, Other: unimplemented("a weapon of the character's choice")},
		{Cash: 10000, Other: relationship(Contact, 1, "")},
		{Cash: 20000, Other: relationship(Ally, 1, "")},
	}
}

func navyMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "accused of negligence resulting in a crewmate's death",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "budget cuts cost you your place in the navy",
			Effects: []Effect{checkChr("CHA", 8,
				[]Effect{throwModifier("next enlistment attempt", 2)},
				nil)},
		},
		{
			Summary: "accidental exposure to a dangerous atmosphere",
			Effects: []Effect{chr("END", -1)},
		},
		{
			Summary: "a rivalry with a superior officer finally ends your career",
			Effects: []Effect{relationship(Rival, 1, "")},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "a love affair puts you in the sights of a rival with more seniority",
			Effects: []Effect{relationship(Enemy, 1, "")},
		},
		{
			Summary: "a shuttle transfer ends in a crash",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch), injury(1)},
		},
		{
			Summary: "your ship is destroyed in battle and you are blamed for the loss",
			Effects: []Effect{
				injury(1),
				checkChr("CHA", 8, nil, []Effect{benefitRolls(-2, 0, ScopeBatch)}),
			},
		},
		{
			Summary: "a psychological profile deems you unfit, and stains your record",
			Effects: []Effect{
				throwModifier("next enlistment attempt", -2),
			},
		},
		{
			Summary: "a tribunal finds you guilty over a friendly-fire incident",
			Effects: []Effect{
				unimplemented("lose every benefit roll from this career"),
				checkChr("CHA", 8, nil, []Effect{transfer("Prisoner", "Prisoner", 2)}),
			},
		},
	}
}
