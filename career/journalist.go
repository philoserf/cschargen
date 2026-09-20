package career

// Journalist is the career of pp. 220-223: reporting from other systems,
// digging for what is hidden, or driving the news from behind it.
func Journalist() Career {
	return Career{
		Name:           "Journalist",
		Cite:           "pp. 220-223",
		Enlistment:     &Check{Characteristic: "EDU", Number: 10},
		EnlistmentMods: []EnlistmentMod{undergraduateDegree(2)},
		MishapEjects:   true,
		Assignments: []Assignment{
			{
				Name:        "Reporter",
				Description: "going out to other systems and reporting back",
				Survival:    Check{Characteristic: "EDU", Number: 7},
				Advancement: Check{Characteristic: "CHA", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Reporter", Rows: [6]Effect{
					skill("Electronics", "Any"), skill("Art", "Writing"), skill("Streetwise"),
					skill("Carouse"), skill("Persuade"), skill("Investigate"),
				}},
				Ranks: [][]Effect{
					{skill("Art", "Writing")},
					{skill("Carouse")},
					nil,
					{skill("Investigate")},
					{skill("Advocate", "Any")},
					nil,
					{skill("Streetwise")},
				},
			},
			{
				Name:        "Investigator",
				Description: "digging deep for the truth and exposing it",
				Survival:    Check{Characteristic: "INT", Number: 7},
				Advancement: Check{Characteristic: "CHA", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Investigator", Rows: [6]Effect{
					skill("Recon"), skill("Advocate", "Any"), skill("Deception", "Any"),
					skill("Investigate"), skill("Carouse"), skill("Interrogation", "Questioning"),
				}},
				Ranks: [][]Effect{
					{skill("Art", "Writing")},
					{skill("Persuade")},
					{skill("Investigate")},
					nil,
					{skill("Interrogation", "Questioning")},
					nil,
					{skill("Carouse")},
				},
			},
			{
				Name:        "Producer",
				Description: "pushing the journalists, directing the content, managing the business",
				Survival:    Check{Characteristic: "INT", Number: 7},
				Advancement: Check{Characteristic: "EDU", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Producer", Rows: [6]Effect{
					skill("Electronics", "Any"), skill("Diplomat"), skill("Admin"),
					skill("Broker"), skill("Advocate", "Any"), skill("Leadership"),
				}},
				Ranks: [][]Effect{
					{skill("Admin")},
					nil,
					{skill("Broker")},
					nil,
					{skill("Persuade")},
					{skill("Advocate", "Any")},
					{skill("Leadership")},
				},
			},
		},
		Tables: []SkillTable{
			personalDevelopment(),
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Admin"), skill("Investigate"), skill("Art", "Writing"),
				skill("Persuade"), skill("Language", "Any"), skill("Etiquette"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Interrogation", "Questioning"), skill("Electronics", "Any"),
				skill("Advocate", "Any"), skill("Diplomat"),
				skill("Art", "Any"), skill("Science", "Any"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 100, Other: chr("INT", 1)},
			{Cash: 250, Other: chr("EDU", 1)},
			{Cash: 500, Other: chr("CHA", 1)},
			{Cash: 1000, Other: relationship(Contact, 1, "")},
			{Cash: 2500, Other: relationship(Contact, 2, "")},
			{Cash: 5000, Other: relationship(Ally, 1, "")},
			{Cash: 10000, Other: relationship(Ally, 2, "")},
		},
		Mishaps: journalistMishaps(),
		Events:  journalistEvents(),
	}
}

func journalistMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "a severe injury in a combat zone, and local medicine to talk your way into",
			Effects: []Effect{pick("talk your way in however you can",
				opt("Persuade", checkSkill("Persuade", 8, []Effect{injury(1)}, []Effect{injury(2)})),
				opt("Diplomat", checkSkill("Diplomat", 8, []Effect{injury(1)}, []Effect{injury(2)})))},
		},
		{
			Summary: "budget cuts at the agency",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "an award-winning report that months later turns out to have been wrong",
			Effects: []Effect{relationship(Enemy, 1, ""), benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "something offensive, said accidentally and broadcast by a rival",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "an uncomfortable truth uncovered, and a news service that makes you the scapegoat",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "caught up in the natural disaster you were reporting on",
			Effects: []Effect{injury(1)},
		},
		{
			Summary: "a lawsuit over false information in a report",
			Effects: []Effect{checkSkill("Advocate", 8, nil,
				[]Effect{benefitRolls(-3, 0, ScopeBatch)})},
		},
		{
			Summary: "captured by local rebels and held until the agency secures your release",
			Effects: []Effect{unimplemented(
				"held 2d6 weeks, with an Injury roll for every four weeks rounded down")},
		},
		{
			Summary: "a bomb in your vehicle, escaped in time but not unhurt",
			Effects: []Effect{injury(2)},
		},
	}
}
