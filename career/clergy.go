package career

// Clergy is the career of pp. 169-172: carrying a religion out to the
// stars, tending a congregation, or studying the thing itself.
func Clergy() Career {
	return Career{
		Name:         "Clergy",
		Cite:         "pp. 169-172",
		Enlistment:   &Check{Characteristic: "INT", Number: 8},
		MishapEjects: true,
		Assignments: []Assignment{
			{
				Name:        "Evangelist",
				Description: "going out into the stars to spread your religion",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "CHA", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Evangelist", Rows: [6]Effect{
					skill("Diplomat"), skill("Electronics", "Any"), skill("Advocate", "Oratory"),
					skill("Persuade"), pickSkill("Drive", "Flyer"), skill("Carouse"),
				}},
				Ranks: [][]Effect{
					nil,
					{skill("Advocate", "Oratory")},
					{skill("Persuade")},
					nil,
					{skill("Carouse")},
					nil,
					{chr("CHA", 1)},
				},
			},
			{
				Name:        "Minister",
				Description: "running a church and tending to its parishioners",
				Survival:    Check{Characteristic: "EDU", Number: 8},
				Advancement: Check{Characteristic: "CHA", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Minister", Rows: [6]Effect{
					skill("Admin"), skill("Persuade"), skill("Advocate", "Oratory"),
					skill("Diplomat"), skill("Electronics", "Any"), skill("Etiquette"),
				}},
				Ranks: [][]Effect{
					{skill("Etiquette")},
					{skill("Advocate", "Oratory")},
					{skill("Diplomat")},
					nil,
					{skill("Persuade")},
					nil,
					{chr("CHA", 1)},
				},
			},
			{
				Name:        "Monk",
				Description: "a monastic order studying the intricacies of the faith",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Monk", Rows: [6]Effect{
					skill("Admin"), skill("Art", "Any"), skill("Persuade"),
					skill("Chef"), skill("Language", "Any"), skill("Diplomat"),
				}},
				Ranks: [][]Effect{
					nil,
					{skill("Language", "Any")},
					nil,
					{skill("Investigate")},
					nil,
					{skill("Art", "Any")},
					{skill("Science", "Any")},
				},
			},
		},
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), skill("Persuade"),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Admin"), skill("Language", "Any"), skill("Persuade"),
				skill("Advocate", "Oratory"), skill("Streetwise"), skill("Diplomat"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Admin"), skill("Science", "Psychology"), skill("Art", "Any"),
				skill("Science", "Any"), skill("Diplomat"), skill("Investigate"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 200, Other: chr("CHA", 1)},
			{Cash: 500, Other: chr("EDU", 1)},
			{Cash: 1000, Other: chr("INT", 1)},
			{Cash: 2000, Other: relationship(Contact, 1, "")},
			{Cash: 5000, Other: relationship(Ally, 1, "")},
			{Cash: 10000, Other: stashItem("a holy book")},
			{Cash: 20000, Other: churchPension()},
		},
		Mishaps: clergyMishaps(),
		Events:  clergyEvents(),
	}
}

// churchPension is p. 128: 2d6 x 10,000 credits, taken either as a lump
// sum "of 10% less than the full value" -- which is 2d6 x 9,000 -- or as
// 1,000 a year, for which "the character must return to their homeworld
// annually to pick up the payment". The engine has no model for a yearly
// income tied to a place, so the second branch is recorded rather than run.
func churchPension() Effect {
	return pick("take the pension as a lump sum or as a yearly payment (p. 128)",
		opt("a lump sum, 10% less", credits("2d6x9000")),
		opt("1,000 a year",
			unimplemented("1,000 credits a year, collected on the homeworld in person")))
}

func clergyMishaps() MishapTable {
	loseAll := loseAllBenefits("lose every benefit roll from this career")

	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{Summary: "a severe crisis of faith, and the discovery that you no longer believe", Effects: nil},
		{
			Summary: "a romance the church cannot accept, and you choose the lover",
			Effects: []Effect{relationship(Ally, 1, ""), loseAll},
		},
		{
			Summary: "caught skimming the church's funds; no charges, but no place either",
			Effects: []Effect{loseAll},
		},
		{
			Summary: "a rivalry with a superior that finally ends your career",
			Effects: []Effect{relationship(Rival, 1, "")},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "the human condition, seen too closely, becomes a breakdown",
			Effects: []Effect{chr("INT", -2)},
		},
		{
			Summary: "the vehicle or ship you are travelling in crashes",
			Effects: []Effect{injury(2)},
		},
		{
			Summary: "something so offensive that you are excommunicated",
			Effects: []Effect{
				loseAll,
				unimplemented("every Contact, Rival and Ally made in this career becomes an Enemy"),
			},
		},
		{
			Summary: "the local government decides your religion threatens social stability",
			Effects: []Effect{
				newHomeworld(),
				skill("Language", "Any"),
				transfer("Vagabond", "", 0),
			},
		},
		{
			Summary: "a revelation: the religion is wrong, and a new sect needs a new world",
			Effects: []Effect{
				loseAll,
				pick("which life the new colony takes",
					opt("Settler", transfer("Colonist", "Settler", 0)),
					opt("Politician", transfer("Colonist", "Politician", 0))),
			},
		},
	}
}
