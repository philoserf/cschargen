package career

// Investigator is the career of pp. 216-219: forensics, a detective's
// caseload, or private work for whoever is paying.
//
// Its rank table is the only one in the book that prints one column of rank
// titles across two assignments and separate benefit columns beneath -- an
// Inspector and a Forensics officer climb the same ladder and are paid
// differently for it.
func Investigator() Career {
	inspectorTitles := []string{
		"Trainee Investigator",
		"Detective Constable",
		"Detective Sergeant",
		"Detective Inspector",
		"Detective Chief Inspector",
		"Detective Superintendent",
		"Detective Chief Superintendent",
	}

	return Career{
		Name:           "Investigator",
		Cite:           "pp. 216-219",
		Enlistment:     &Check{Characteristic: "INT", Number: 8},
		EnlistmentMods: []EnlistmentMod{apparentAgeOver40()},
		MishapEjects:   true,
		RankTitles:     inspectorTitles,
		Assignments: []Assignment{
			{
				Name:        "Forensics",
				Description: "law enforcement, solving crimes with science",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Forensics", Rows: [6]Effect{
					skill("Art", "Any"), skill("Electronics", "Any"),
					skill("Science", "Forensics"), skill("Investigate"),
					skill("Medic", "Any"), skill("Science", "Any"),
				}},
				Ranks: [][]Effect{
					nil,
					{skill("Science", "Forensics")},
					{skill("Investigate")},
					{skill("Recon")},
					nil,
					{skill("Medic", "Diagnosis")},
					nil,
				},
			},
			{
				Name:        "Inspector",
				Description: "law enforcement, tasked with solving crimes",
				Survival:    Check{Characteristic: "EDU", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Inspector", Rows: [6]Effect{
					skill("Advocate", "Legal"), skill("Carouse"), skill("Investigate"),
					skill("Gun Combat", "Any"), skill("Interrogation", "Questioning"),
					skill("Recon"),
				}},
				Ranks: [][]Effect{
					{skill("Advocate", "Legal")},
					{skill("Investigate")},
					{skill("Etiquette")},
					nil,
					{skill("Admin")},
					nil,
					{skill("Advocate", "Politics")},
				},
			},
			{
				Name:        "PI",
				Description: "a private investigator, in a firm or independent",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "CHA", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "PI", Rows: [6]Effect{
					skill("Melee", "Any"), skill("Carouse"), skill("Gun Combat", "Any"),
					skill("Interrogation", "Questioning"), skill("Stealth"), skill("Recon"),
				}},
				Ranks: [][]Effect{
					{skill("Recon")},
					{skill("Investigate")},
					{skill("Etiquette")},
					{skill("Stealth")},
					nil, nil,
					{relationship(Contact, 2, "")},
				},
			},
		},
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), skill("Carouse"),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Carouse"), skill("Advocate", "Legal"), skill("Investigate"),
				skill("Interrogation", "Questioning"), skill("Deception", "Any"), skill("Recon"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Admin"), skill("Diplomat"), skill("Advocate", "Any"),
				skill("Investigate"), skill("Science", "Forensics"), skill("Electronics", "Any"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 0, Other: chr("CHA", 1)},
			{Cash: 250, Other: chr("EDU", 1)},
			{Cash: 500, Other: chr("INT", 1)},
			{Cash: 1000, Other: relationship(Contact, 1, "")},
			{Cash: 2000, Other: stashItem("a weapon")},
			{Cash: 5000, Other: relationship(Ally, 1, "")},
			{Cash: 10000, Other: relationship(Contact, 2, "")},
		},
		Mishaps: investigatorMishaps(),
		Events:  investigatorEvents(),
	}
}

func investigatorMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			// The only mishap in the book that pays: budget cuts come with
			// severance.
			Summary: "budget cuts, and a release with severance",
			Effects: []Effect{benefitRolls(1, 0, ScopeBatch)},
		},
		{Summary: "tired of the stress and the strain", Effects: nil},
		{
			Summary: "an obsession with a case the authorities took away from you",
			Effects: []Effect{relationship(Enemy, 1, "")},
		},
		{
			Summary: "someone you jailed escapes and goes after the people around you",
			Effects: []Effect{
				loseTiesRolled("1d3", Ally, Contact),
				relationship(Enemy, 1, ""),
			},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "one of the people you put away belongs to a prominent family",
			Effects: []Effect{relationship(Enemy, 1, ""), benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "an addiction that impairs you past the point of staying",
			Effects: []Effect{
				chr("END", -2),
				unimplemented("choose the alcohol or drug you are addicted to"),
			},
		},
		{
			Summary: "your partner is killed, and you vow revenge",
			Effects: []Effect{relationship(Enemy, 1, ""), loseTie(Ally)},
		},
		{
			Summary: "a dangerous disease contracted",
			Effects: []Effect{chr("STR", -1), chr("DEX", -1), chr("END", -2)},
		},
		{
			Summary: "an investigation that finds you little more than a criminal",
			Effects: []Effect{
				loseAllBenefits("lose every benefit roll from this career"),
				newHomeworld(),
				transfer("Prisoner", "Prisoner", 2),
				transfer("Vagabond", "", 0),
			},
		},
	}
}
