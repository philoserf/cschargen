package career

// Marine is the career of pp. 224-228: the troops a nation puts first into
// a battle, aboard its ships, and into its covert operations.
//
// Its rank tables are what showed that a rank table is not seven rows:
// Marine prints nine enlisted ranks, E0 to E8, and eight officer ranks, O0
// to O7 (p. 225).
func Marine() Career {
	return Career{
		Name:       "Marine",
		Cite:       "pp. 224-228",
		Enlistment: &Check{Characteristic: "END", Number: 8},
		Commission: &Check{Characteristic: "EDU", Number: 8},
		RankTitles: []string{
			"Private", "Private First Class", "Lance Corporal", "Corporal", "Sergeant",
			"Staff Sergeant", "Gunnery Sergeant", "Master Sergeant", "Sergeant Major",
		},
		OfficerTitles: []string{
			"Second Lieutenant", "First Lieutenant", "Captain", "Major",
			"Lieutenant Colonel", "Colonel", "Brigadier General", "Major General",
		},
		MishapEjects: true,
		Assignments:  marineAssignments(),
		Tables:       marineTables(),
		Benefits:     marineBenefits(),
		Mishaps:      marineMishaps(),
		Events:       marineEvents(),
	}
}

func marineAssignments() []Assignment {
	return []Assignment{
		{
			Name:        "Assault",
			Description: "first into the scene of a battle",
			Survival:    Check{Characteristic: "DEX", Number: 8},
			Advancement: Check{Characteristic: "INT", Number: 8},
			Skills: SkillTable{Kind: AssignmentSkills, Name: "Assault", Rows: [6]Effect{
				skill("Explosives"), skill("Recon"), skill("Gun Combat", "Any"),
				skill("Heavy Weapons", "Any"), skill("Suit", "Any"), skill("Survival", "Any"),
			}},
			Ranks: [][]Effect{
				{skill("Discipline")},
				{skill("Gun Combat", "Any")},
				nil,
				{skill("Recon")},
				nil,
				{skill("Admin")},
				{skill("Tactics", "Military")},
				{skill("Leadership")},
				nil,
			},
			OfficerRanks: [][]Effect{
				{skill("Admin")},
				nil,
				{skill("Leadership")},
				{skill("Tactics", "Military")},
				nil, nil,
				{skill("Tactics", "Military")},
				{skill("Advocate", "Any")},
			},
		},
		{
			Name:        "Ship's Troops",
			Description: "security aboard a ship, repelling boarders and boarding others",
			Survival:    Check{Characteristic: "DEX", Number: 8},
			Advancement: Check{Characteristic: "EDU", Number: 8},
			Skills: SkillTable{Kind: AssignmentSkills, Name: "Ship's Troops", Rows: [6]Effect{
				skill("Recon"), skill("Draw"), skill("Gun Combat", "Any"),
				skill("Melee", "Any"), skill("Survival", "Freefall"),
				skill("Suit", "Battle Armor", "Vacc Suit"),
			}},
			Ranks: [][]Effect{
				{skill("Discipline")},
				{skill("Suit", "Vacc Suit")},
				{skill("Gun Combat", "Any")},
				{skill("Survival", "Freefall")},
				nil,
				{skill("Admin")},
				{skill("Tactics", "Military")},
				{skill("Leadership")},
				nil,
			},
			OfficerRanks: [][]Effect{
				{skill("Admin")},
				{skill("Survival", "Freefall")},
				{skill("Leadership")},
				{skill("Tactics", "Military")},
				nil, nil,
				{skill("Tactics", "Military")},
				{skill("Advocate", "Any")},
			},
		},
		{
			Name:        "Spec Ops",
			Description: "small teams specializing in covert action",
			Survival:    Check{Characteristic: "DEX", Number: 8},
			Advancement: Check{Characteristic: "INT", Number: 8},
			Skills: SkillTable{Kind: AssignmentSkills, Name: "Spec Ops", Rows: [6]Effect{
				skill("Deception", "Any"), skill("Explosives"), skill("Recon"),
				skill("Stealth"), skill("Gun Combat", "Any"), skill("Melee", "Any"),
			}},
			Ranks: [][]Effect{
				{skill("Discipline")},
				{skill("Gun Combat", "Any")},
				{skill("Draw")},
				{skill("Melee", "Any")},
				nil,
				{skill("Admin")},
				{skill("Tactics", "Military")},
				{skill("Leadership")},
				nil,
			},
			OfficerRanks: [][]Effect{
				{skill("Admin")},
				nil,
				{skill("Leadership")},
				{skill("Tactics", "Military")},
				nil, nil,
				{skill("Tactics", "Military")},
				{skill("Advocate", "Any")},
			},
		},
	}
}

func marineTables() []SkillTable {
	return []SkillTable{
		personalDevelopment(),
		{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
			skill("Discipline"), skill("Suit", "Any"), skill("Gun Combat", "Any"),
			skill("Heavy Weapons", "Any"), skill("Melee", "Any"), skill("Survival", "Any"),
		}},
		{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
			skill("Admin"), skill("Electronics", "Any"), skill("Art", "Any"),
			skill("Science", "Any"), skill("Medic", "Any"), skill("Tactics", "Military"),
		}},
		{Kind: OfficerSkills, Name: "Officer Skills", Rows: [6]Effect{
			skill("Advocate", "Any"), skill("Etiquette"), skill("Instruction"),
			skill("Interrogation", "Any"), skill("Persuade"), skill("Leadership"),
		}},
	}
}

func marineBenefits() [7]BenefitRow {
	return [7]BenefitRow{
		{Cash: 250, Other: chr("END", 1)},
		{Cash: 500, Other: chr("DEX", 1)},
		{Cash: 1000, Other: chr("INT", 1)},
		{Cash: 2000, Other: chr("EDU", 1)},
		{Cash: 5000, Other: unimplemented("a weapon of the character's choice")},
		{Cash: 10000, Other: relationship(Contact, 1, "")},
		{Cash: 25000, Other: unimplemented("armor of the character's choice")},
	}
}

func marineMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "a senior officer files a series of negative performance reports",
			Effects: []Effect{relationship(Enemy, 1, "")},
		},
		{
			Summary: "budget cuts cost you your place in the marines",
			Effects: []Effect{pick("argue your case however you can",
				opt("CHA", checkChr("CHA", 8,
					[]Effect{throwModifier("next enlistment attempt", 2)}, nil)),
				opt("Persuade", checkSkill("Persuade", 8,
					[]Effect{throwModifier("next enlistment attempt", 2)}, nil)),
			)},
		},
		{
			Summary: "a romance with another marine, and a jealous superior who ends your service",
			Effects: []Effect{relationship(Ally, 1, ""), relationship(Enemy, 1, "")},
		},
		{
			Summary: "accused of negligence resulting in another marine's death",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "a deadly virus survived, and carried for life",
			Effects: []Effect{chr("STR", -2), chr("END", -2)},
		},
		{
			Summary: "pushing a recruit to their limits kills them, and you are found guilty",
			Effects: []Effect{
				unimplemented("lose every benefit roll from this career"),
				pick("defend yourself however you can",
					opt("CHA", checkChr("CHA", 8, nil,
						[]Effect{transfer("Prisoner", "Prisoner", 3)})),
					opt("Advocate (Legal)", checkSkill("Advocate", 8, nil,
						[]Effect{transfer("Prisoner", "Prisoner", 3)}))),
			},
		},
		{
			Summary: "an operation goes horribly wrong and the blame is placed on you",
			Effects: []Effect{
				injury(1),
				pick("accept the blame or shift it",
					opt("accept it",
						unimplemented("lose every benefit roll from this career"),
						throwModifier("next enlistment attempt", 2)),
					opt("shift it", checkSkill("Persuade", 8,
						[]Effect{benefitRolls(-2, 0, ScopeBatch), relationship(Enemy, 1, "")},
						[]Effect{
							unimplemented("lose every benefit roll from this career"),
							transfer("Prisoner", "Prisoner", 4),
						}))),
			},
		},
		{
			Summary: "a training accident kills several marines",
			Effects: []Effect{
				unimplemented("lose every benefit roll from this career"),
				relationship(Enemy, 1, ""),
			},
		},
		{
			Summary: "a tribunal finds you guilty over a friendly-fire incident",
			Effects: []Effect{
				unimplemented("lose every benefit roll from this career"),
				transfer("Prisoner", "Prisoner", 3),
				throwModifier("next enlistment attempt", -4),
			},
		},
	}
}
