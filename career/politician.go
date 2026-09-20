package career

// Politician is the career of pp. 256-260: standing for office, running
// the campaign, or doing whatever the campaign needs done.
//
// It reaches further than any other career in the book. Its mishaps send
// characters to Colonist and Journalist; its events send them to
// Diplomatic Service, Celebrity, Colonist and Prisoner.
func Politician() Career {
	return Career{
		Tags:           []Tag{TagGovernment},
		Name:           "Politician",
		Cite:           "pp. 256-260",
		Enlistment:     &Check{Characteristic: "INT", Number: 8},
		EnlistmentMods: []EnlistmentMod{perPreviousCareer()},
		MishapEjects:   true,
		Assignments: []Assignment{
			{
				Name:        "Candidate",
				Description: "standing for election on a world that elects people",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "CHA", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Candidate", Rows: [6]Effect{
					skill("Advocate", "Legal"), skill("Persuade"), skill("Advocate", "Politics"),
					skill("Advocate", "Oratory"), skill("Deception", "Lie"), skill("Leadership"),
				}},
				Ranks: [][]Effect{
					{skill("Advocate", "Politics")},
					{skill("Deception", "Lie")},
					nil,
					{skill("Persuade")},
					nil,
					{relationship(Ally, 1, "")},
					{skill("Leadership")},
				},
			},
			{
				Name:        "Manager",
				Description: "organizing and advancing a favoured politician or cause",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Manager", Rows: [6]Effect{
					skill("Art", "Writing"), skill("Admin"), skill("Broker"),
					skill("Advocate", "Politics"), skill("Electronics", "Comms"),
					skill("Leadership"),
				}},
				Ranks: [][]Effect{
					{skill("Admin")},
					{skill("Broker")},
					nil,
					{skill("Advocate", "Politics")},
					nil,
					{relationship(Contact, 1, "")},
					{relationship(Ally, 1, "")},
				},
			},
			{
				Name:        "Operative",
				Description: "everything from getting the coffee to doing the dirty work",
				Survival:    Check{Characteristic: "END", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Operative", Rows: [6]Effect{
					skill("Etiquette"), skill("Diplomat"), skill("Carouse"),
					skill("Broker"), skill("Deception", "Any"), skill("Streetwise"),
				}},
				Ranks: [][]Effect{
					{skill("Etiquette")},
					nil,
					{skill("Broker")},
					nil,
					{skill("Deception", "Any")},
					{relationship(Contact, 1, "")},
					{relationship(Ally, 1, "")},
				},
			},
		},
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), skill("Carouse"),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Etiquette"), skill("Persuade"), skill("Advocate", "Any"),
				skill("Carouse"), skill("Electronics", "Comms"), skill("Deception", "Lie"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Art", "Any"), skill("Advocate", "Politics"), skill("Diplomat"),
				skill("Advocate", "Oratory"), skill("Leadership"), skill("Science", "Any"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 0, Other: chr("INT", 1)},
			{Cash: 1000, Other: chr("EDU", 1)},
			{Cash: 2000, Other: relationship(Contact, 1, "")},
			{Cash: 6000, Other: relationship(Ally, 1, "")},
			{Cash: 10000, Other: chr("CHA", 1)},
			{Cash: 20000, Other: relationship(Contact, 2, "")},
			{Cash: 100000, Other: chr("CHA", 2)},
		},
		Mishaps: politicianMishaps(),
		Events:  politicianEvents(),
	}
}

func politicianMishaps() MishapTable {
	loseAll := loseAllBenefits("lose every benefit roll from this career")

	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{Summary: "politics is a dirty business and you have had enough", Effects: nil},
		{
			Summary: "the civilized worlds are no longer for you",
			Effects: []Effect{transfer("Colonist", "Politician", 0)},
		},
		{
			Summary: "something so offensive that your reputation cannot recover",
			Effects: []Effect{loseAll},
		},
		{
			Summary: "a change of government or party, and your dismissal with it",
			Effects: []Effect{checkSkill("Diplomat", 8, nil, []Effect{loseAll})},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "your usefulness here is at an end, but a media company wants a consultant",
			Effects: []Effect{transfer("Journalist", "Producer", 0)},
		},
		{
			Summary: "a movement your supporters do not care for",
			Effects: []Effect{chr("CHA", -2)},
		},
		{Summary: "your ship or vehicle crashes", Effects: []Effect{injury(2)}},
		{
			Summary: "a disease that ends your work in this field",
			Effects: []Effect{chr("END", -2)},
		},
		{
			Summary: "a social faux pas your rivals have been waiting for",
			Effects: []Effect{checkSkill("Advocate", 8,
				[]Effect{benefitRolls(-2, 0, ScopeBatch)},
				[]Effect{
					loseAll,
					enlistmentPenalty(-4, "-4 to enter any career that enlists on CHA or EDU",
						enlistmentNarrowing{OnCharacteristics: []string{"CHA", "EDU"}, Standing: true}),
				})},
		},
	}
}
