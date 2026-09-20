package career

// Adventurer is the career of pp. 146-149: hunting, exploring for its own
// sake, or digging up what earlier people left behind.
func Adventurer() Career {
	return Career{
		Name:           "Adventurer",
		Cite:           "pp. 146-149",
		Enlistment:     &Check{Characteristic: "DEX", Number: 8},
		EnlistmentMods: []EnlistmentMod{apparentAgeOver40(), perPreviousCareer()},
		MishapEjects:   true,
		Assignments: []Assignment{
			{
				Name:        "Hunter",
				Description: "hunting animals to kill them, capture them or holograph them",
				Survival:    Check{Characteristic: "END", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Hunter", Rows: [6]Effect{
					skill("Tactics", "Military"), skill("Gun Combat", "Any"), skill("Animals", "Any"),
					skill("Stealth"), skill("Recon"), skill("Survival", "Any"),
				}},
				Ranks: [][]Effect{
					{skill("Survival", "Any")},
					nil,
					{skill("Navigation")},
					{skill("Tactics", "Military")},
					nil,
					{skill("Recon")},
					nil,
				},
			},
			{
				Name:        "Pioneer",
				Description: "exploring worlds for the pleasure of the unknown, for nobody",
				Survival:    Check{Characteristic: "END", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Pioneer", Rows: [6]Effect{
					skill("Seafarer", "Any"), skill("Gun Combat", "Any"), skill("Navigation"),
					skill("Recon"), skill("Survival", "Any"), skill("Survival", "Any"),
				}},
				Ranks: [][]Effect{
					{skill("Survival", "Any")},
					nil, nil,
					{skill("Navigation")},
					nil,
					{skill("Recon")},
					{skill("Drive", "Wheeled", "Tracked")},
				},
			},
			{
				Name:        "Relic Raider",
				Description: "hunting alien artifacts and relics of the early colonies",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Relic Raider", Rows: [6]Effect{
					skill("Carouse"), skill("Science", "Any"), skill("Recon"),
					skill("Investigate"), skill("Navigation"), skill("Survival", "Any"),
				}},
				Ranks: [][]Effect{
					{skill("Survival", "Any")},
					nil, nil,
					{skill("Science", "Archeology")},
					nil,
					{skill("Investigate")},
					{skill("Jack of All Trades")},
				},
			},
		},
		Tables: []SkillTable{
			personalDevelopment(),
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Drive", "Any"), skill("Investigate"), skill("Survival", "Any"),
				skill("Recon"), skill("Navigation"), skill("Flyer", "Any"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Electronics", "Any"), skill("Art", "Any"), skill("Science", "Any"),
				skill("Language", "Any"), skill("Jack of All Trades"), skill("Medic", "First Aid"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 500, Other: chr("END", 1)},
			{Cash: 1000, Other: chr("INT", 1)},
			{Cash: 2000, Other: chr("EDU", 1)},
			{Cash: 5000, Other: relationship(Contact, 1, "")},
			{Cash: 10000, Other: weaponOrItsUse()},
			{Cash: 25000, Other: relationship(Ally, 1, "")},
			{Cash: 50000, Other: chr("CHA", 1)},
		},
		Mishaps: adventurerMishaps(),
		Events:  adventurerEvents(),
	}
}

func adventurerMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "a Rival takes credit for some of your work",
			Effects: []Effect{benefitRolls(-1, 0, ScopeBatch)},
		},
		{Summary: "the long stretches alone have got to you", Effects: nil},
		{Summary: "accidental exposure to a dangerous atmosphere", Effects: []Effect{chr("END", -1)}},
		{
			Summary: "a Rival accuses you of being a fraud, and your career is ruined",
			Effects: []Effect{
				loseAllBenefits("lose every benefit roll from this career"),
				transfer("Vagabond", "", 0),
			},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "the last world had more background radiation than you thought",
			Effects: []Effect{chr("END", -2)},
		},
		{
			Summary: "an alien creature takes you by surprise and nearly kills you",
			Effects: []Effect{injury(2)},
		},
		{Summary: "an alien disease contracted", Effects: []Effect{chr("END", -3)}},
		{
			Summary: "your ship crashes on a wild planet and everyone else aboard is killed",
			Effects: []Effect{
				injury(1),
				checkSkill("Survival", 8,
					[]Effect{rollOtherAssignment(), rollOtherAssignment()},
					[]Effect{skill("Survival", "Any"), skill("Navigation"), injury(1)}),
			},
		},
		{
			Summary: "an alien artifact puts you into suspended animation for years",
			Effects: nil,
		},
	}
}
