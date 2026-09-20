package career

// DiplomaticService is the career of pp. 185-189: the foreign service of a
// government, its embassy staff, or the people who keep both safe.
//
// Colonist event 54 sends a character to roll twice on this career's
// Ambassador table (p. 176), which is why it is the one career another
// career reaches into rather than transfers to.
func DiplomaticService() Career {
	return Career{
		Name:           "Diplomatic Service",
		Cite:           "pp. 185-189",
		Enlistment:     &Check{Characteristic: "EDU", Number: 7},
		EnlistmentMods: []EnlistmentMod{perPreviousCareer()},
		MishapEjects:   true,
		Assignments: []Assignment{
			{
				Name:        "Ambassador",
				Description: "the foreign service, dealing with other foreign services",
				Survival:    Check{Characteristic: "EDU", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Ambassador", Rows: [6]Effect{
					skill("Etiquette"), skill("Carouse"), skill("Diplomat"),
					skill("Language", "Any"), skill("Advocate", "Any"), skill("Persuade"),
				}},
				Ranks: [][]Effect{
					{skill("Etiquette")},
					{skill("Language", "Any")},
					nil,
					{skill("Diplomat")},
					nil,
					{skill("Advocate", "Any")},
					{skill("Diplomat")},
				},
			},
			{
				Name:        "Generalist",
				Description: "the support staff of an embassy or consulate",
				Survival:    Check{Characteristic: "INT", Number: 7},
				Advancement: Check{Characteristic: "EDU", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Generalist", Rows: [6]Effect{
					skill("Electronics", "Any"), skill("Admin"), skill("Advocate", "Any"),
					skill("Etiquette"), skill("Mechanic"), skill("Chef"),
				}},
				Ranks: [][]Effect{
					{skill("Admin")},
					nil,
					{skill("Language", "Any")},
					nil,
					{skill("Etiquette")},
					nil,
					{skill("Leadership")},
				},
			},
			{
				Name:        "Diplomatic Security",
				Description: "criminal investigation and security abroad, in your government's name",
				Survival:    Check{Characteristic: "END", Number: 7},
				Advancement: Check{Characteristic: "EDU", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Diplomatic Security", Rows: [6]Effect{
					skill("Draw"), skill("Gun Combat", "Any"), skill("Investigate"),
					skill("Recon"), skill("Interrogation", "Any"), skill("Streetwise"),
				}},
				Ranks: [][]Effect{
					{skill("Gun Combat", "Any")},
					nil,
					{skill("Recon")},
					nil,
					{skill("Investigate")},
					{skill("Interrogation", "Any")},
					{skill("Leadership")},
				},
			},
		},
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), skill("Carouse"),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Advocate", "Any"), skill("Admin"), skill("Language", "Any"),
				skill("Persuade"), skill("Survival", "Any"), skill("Etiquette"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Science", "Any"), skill("Advocate", "Any"), skill("Diplomat"),
				skill("Leadership"), skill("Medic", "Any"), skill("Art", "Any"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 500, Other: chr("EDU", 1)},
			{Cash: 1000, Other: chr("INT", 1)},
			{Cash: 5000, Other: weaponOrItsUse()},
			{Cash: 10000, Other: relationship(Contact, 1, "")},
			{Cash: 25000, Other: chr("CHA", 1)},
			{Cash: 50000, Other: relationship(Ally, 1, "")},
			{Cash: 100000, Other: relationship(Ally, 2, "")},
		},
		Mishaps: diplomaticMishaps(),
		Events:  diplomaticEvents(),
	}
}

func diplomaticMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "accused of negligence in your duties",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "a media story implicates you in a crime; no charges, but no reputation either",
			Effects: []Effect{benefitRolls(-1, 0, ScopeBatch)},
		},
		{
			Summary: "a change of government, and every current employee removed",
			Effects: []Effect{checkSkill("Diplomat", 8, nil,
				[]Effect{loseAllBenefits("lose every benefit roll from this career")})},
		},
		{
			Summary: "a rivalry with a superior that finally ends your career",
			Effects: []Effect{relationship(Rival, 1, "")},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "caught up in a corruption sting",
			Effects: []Effect{checkSkill("Advocate", 8, nil,
				[]Effect{loseBenefitRollsRolled("1d6")})},
		},
		{
			Summary: "a clandestine meeting raided by a rival rebel group",
			Effects: []Effect{
				pick("escape however you can",
					opt("Gun Combat", checkSkill("Gun Combat", 8,
						[]Effect{injury(1)}, []Effect{injury(3)})),
					opt("Deception", checkSkill("Deception", 8,
						[]Effect{injury(1)}, []Effect{injury(3)}))),
				unimplemented("-2 to enlist in any government-related career"),
			},
		},
		{Summary: "your ship crashes on landing", Effects: []Effect{injury(2)}},
		{Summary: "a life-threatening illness", Effects: []Effect{chr("END", -2)}},
		{
			Summary: "your actions start an interstellar war, and you are its first prisoner",
			Effects: []Effect{
				injury(2),
				transfer("Prisoner", "Prisoner", 2),
				unimplemented("-4 to enlist in any government-related career"),
			},
		},
	}
}
