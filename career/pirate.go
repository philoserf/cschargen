package career

// Pirate is the career of pp. 250-255: crewing a raider, boarding what it
// catches, shooting at what it cannot, or keeping it running.
//
// It takes a commission, which makes it the only non-military career in the
// book with an officer track.
func Pirate() Career {
	return Career{
		Name:         "Pirate",
		Cite:         "pp. 250-255",
		Enlistment:   &Check{Characteristic: "INT", Number: 6},
		Commission:   &Check{Characteristic: "INT", Number: 8},
		MishapEjects: true,
		Assignments: []Assignment{
			{
				Name:        "Crew",
				Description: "the command crew of a pirate vessel",
				Survival:    Check{Characteristic: "END", Number: 7},
				Advancement: Check{Characteristic: "EDU", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Crew", Rows: [6]Effect{
					skill("Admin"), skill("Astrogation"), skill("Electronics", "Any"),
					skill("Pilot", "Any"), skill("Gunner", "Any"), skill("Mechanic"),
				}},
				Ranks: [][]Effect{
					{skill("Electronics", "Any")},
					nil,
					{skill("Astrogation")},
					nil,
					{skill("Pilot", "Any")},
					nil,
					{skill("Leadership")},
				},
				OfficerRanks: [][]Effect{
					{skill("Admin")},
					nil,
					{skill("Astrogation")},
					{skill("Broker")},
					{skill("Pilot", "Any")},
					nil,
					{skill("Leadership")},
				},
			},
			{
				Name:        "Boarder",
				Description: "going over to take the cargo or subdue the crew",
				Survival:    Check{Characteristic: "DEX", Number: 7},
				Advancement: Check{Characteristic: "INT", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Boarder", Rows: [6]Effect{
					skill("Pilot", "Small Craft"), skill("Tactics", "Military"),
					skill("Gun Combat", "Any"), skill("Melee", "Any"),
					skill("Explosives"), skill("Draw"),
				}},
				Ranks: [][]Effect{
					{skill("Suit", "Vacc Suit")},
					{skill("Gun Combat", "Any")},
					nil,
					{skill("Melee", "Any")},
					nil,
					{skill("Tactics", "Military")},
					{skill("Leadership")},
				},
				OfficerRanks: [][]Effect{
					{skill("Suit", "Vacc Suit")},
					{skill("Gun Combat", "Any")},
					{skill("Tactics", "Military")},
					nil,
					{skill("Explosives")},
					nil,
					{skill("Leadership")},
				},
			},
			{
				Name:        "Gunner",
				Description: "the ship's weaponry, defending or damaging",
				Survival:    Check{Characteristic: "DEX", Number: 7},
				Advancement: Check{Characteristic: "EDU", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Gunner", Rows: [6]Effect{
					skill("Admin"), skill("Gun Combat", "Any"), skill("Gunner", "Any"),
					skill("Electronics", "Any"), skill("Mechanic"), skill("Tactics", "Naval"),
				}},
				Ranks: [][]Effect{
					{skill("Gunner", "Any")},
					nil,
					{skill("Electronics", "Sensors")},
					nil,
					{skill("Tactics", "Naval")},
					nil,
					{skill("Leadership")},
				},
				OfficerRanks: [][]Effect{
					{skill("Gunner", "Any")},
					{skill("Admin")},
					{skill("Electronics", "Sensors")},
					{skill("Tactics", "Naval")},
					nil, nil,
					{skill("Leadership")},
				},
			},
			{
				Name:        "Engineer",
				Description: "keeping the pirate vessel running",
				Survival:    Check{Characteristic: "INT", Number: 7},
				Advancement: Check{Characteristic: "EDU", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Engineer", Rows: [6]Effect{
					skill("Admin"), skill("Mechanic"), skill("Electronics", "Any"),
					skill("Engineer", "Any"), skill("Science", "Any"), skill("Suit", "Vacc Suit"),
				}},
				Ranks: [][]Effect{
					{skill("Mechanic")},
					nil,
					{skill("Engineer", "M-Drive")},
					nil,
					{skill("Engineer", "Z-Drive")},
					nil,
					{skill("Leadership")},
				},
				OfficerRanks: [][]Effect{
					{skill("Mechanic")},
					{skill("Admin")},
					{skill("Engineer", "M-Drive")},
					{skill("Engineer", "Z-Drive")},
					nil, nil,
					{skill("Leadership")},
				},
			},
		},
		Tables: []SkillTable{
			personalDevelopment(),
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Suit", "Vacc Suit"), skill("Gun Combat", "Any"), skill("Deception", "Any"),
				skill("Survival", "Any"), skill("Streetwise"), skill("Melee", "Any"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Science", "Any"), skill("Medic", "Any"), skill("Broker"),
				skill("Investigate"), skill("Language", "Any"), skill("Interrogation", "Any"),
			}},
			{Kind: OfficerSkills, Name: "Officer Skills", Rows: [6]Effect{
				skill("Language", "Any"), skill("Persuade"), skill("Tactics", "Naval"),
				skill("Broker"), skill("Instruction"), skill("Leadership"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 0, Other: chr("END", 1)},
			{Cash: 0, Other: chr("INT", 1)},
			{Cash: 5000, Other: chr("DEX", 1)},
			{Cash: 10000, Other: stashItem("a vacc suit")},
			{Cash: 10000, Other: unimplemented("a weapon of the character's choice")},
			{Cash: 25000, Other: unimplemented("a prize share of 2d6 x 10,000 credits")},
			{Cash: 50000, Other: unimplemented("two prize shares of 2d6 x 10,000 credits each")},
		},
		Mishaps: pirateMishaps(),
		Events:  pirateEvents(),
	}
}

func pirateMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "accused of negligence resulting in a crewmate's death",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "a drive accident strands you in a system with no colony at all",
			Effects: []Effect{pick("what surviving there taught you",
				opt("Survival", skill("Survival", "Any")),
				opt("Animals", skill("Animals", "Any")),
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
				opt("Navigation", skill("Navigation")))},
		},
		{Summary: "accidental exposure to a dangerous atmosphere", Effects: []Effect{chr("END", -1)}},
		{
			Summary: "the ship wrecked by a system defence ship, and sold for scrap",
			Effects: []Effect{injury(2), benefitRolls(3, 0, ScopeBatch)},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "your crewmates sell you out, and you are beaten for it",
			Effects: []Effect{injury(1), relationship(Enemy, 1, "")},
		},
		{
			Summary: "marines raid the secret base, and you leave piracy behind",
			Effects: []Effect{injury(1), skill("Stealth"), benefitRolls(-3, 0, ScopeBatch)},
		},
		{
			Summary: "the ship destroyed in battle and you blamed for it",
			Effects: []Effect{
				injury(1),
				benefitRolls(-3, 0, ScopeBatch),
				transfer("Vagabond", "", 0),
			},
		},
		{
			Summary: "made to walk the plank in a faulty suit, and rescued by a passing merchant",
			Effects: []Effect{
				injury(2),
				skill("Suit", "Vacc Suit"),
				loseAllBenefits("lose every benefit roll from this career"),
				relationship(Enemy, 1, ""),
			},
		},
		{
			Summary: "arrested, and two terms for it",
			Effects: []Effect{transfer("Prisoner", "Prisoner", 2)},
		},
	}
}
