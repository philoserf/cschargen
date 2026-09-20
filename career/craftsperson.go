package career

// Craftsperson is the career of pp. 181-184: designing construction
// projects, installing what they need, or doing the work.
func Craftsperson() Career {
	tradesRanks := [][]Effect{
		nil,
		{skill("Trade", "Construction")},
		{skill("Mechanic")},
		nil,
		{skill("Electronics", "Any")},
		nil,
		{skill("Jack of All Trades")},
	}
	architectRanks := [][]Effect{
		nil,
		{skill("Trade", "Architect")},
		{skill("Art", "Any")},
		nil,
		{skill("Electronics", "Any")},
		nil,
		{skill("Broker")},
	}

	return Career{
		Name:           "Craftsperson",
		Cite:           "pp. 181-184",
		Enlistment:     &Check{Characteristic: "END", Number: 8},
		EnlistmentMods: []EnlistmentMod{apparentAgeOver40()},
		MishapEjects:   true,
		Assignments: []Assignment{
			{
				Name:        "Architect",
				Description: "designing the newest construction projects",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Architect", Rows: [6]Effect{
					skill("Art", "Any"), skill("Engineer", "Power"), skill("Electronics", "Any"),
					skill("Trade", "Architect"), skill("Broker"), skill("Persuade"),
				}},
				Ranks: architectRanks,
			},
			{
				Name:        "Tradesperson",
				Description: "installing plumbing, electrical or holography on a site",
				Survival:    Check{Characteristic: "DEX", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Tradesperson", Rows: [6]Effect{
					skill("Language", "Any"), skill("Mechanic"), skill("Trade", "Construction"),
					skill("Electronics", "Any"), skill("Survival", "Any"), skill("Jack of All Trades"),
				}},
				Ranks: tradesRanks,
			},
			{
				Name:        "Labor",
				Description: "the backbone of the site",
				Survival:    Check{Characteristic: "STR", Number: 8},
				Advancement: Check{Characteristic: "END", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Labor", Rows: [6]Effect{
					skill("Language", "Any"), skill("Electronics", "Any"), skill("Mechanic"),
					skill("Trade", "Construction"), skill("Survival", "Any"), skill("Jack of All Trades"),
				}},
				Ranks: tradesRanks,
			},
		},
		Tables: []SkillTable{
			personalDevelopment(),
			// The one Service Skills table in the book with a
			// characteristic in it rather than six skills.
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Carouse"), chr("END", 1), skill("Trade", "Construction"),
				skill("Engineer", "Power"), pickSkill("Drive", "Flyer"), skill("Mechanic"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Admin"), skill("Broker"), skill("Diplomat"),
				skill("Explosives"), skill("Advocate", "Any"), skill("Science", "Any"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 100, Other: chr("EDU", 1)},
			{Cash: 250, Other: chr("INT", 1)},
			{Cash: 500, Other: chr("STR", 1)},
			{Cash: 1000, Other: chr("DEX", 1)},
			{Cash: 2000, Other: chr("END", 1)},
			{Cash: 5000, Other: stashItem("a tool kit")},
			{Cash: 10000, Other: unimplemented("a company share, 2d6 x 100,000 credits")},
		},
		Mishaps: craftMishaps(),
		Events:  craftEvents(),
	}
}

func craftMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{Summary: "you are not cut out for this kind of work", Effects: nil},
		{
			Summary: "a downturn cancels every project",
			Effects: []Effect{benefitRolls(-1, 0, ScopeBatch)},
		},
		{
			Summary: "constant arguments with your boss end your contract",
			Effects: []Effect{
				unimplemented("lose any company shares"),
				benefitRolls(1, 0, ScopeBatch),
			},
		},
		{Summary: "accidental exposure to a dangerous chemical", Effects: []Effect{chr("END", -2)}},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "a romance with a coworker's spouse",
			Effects: []Effect{relationship(Ally, 1, ""), relationship(Enemy, 1, "")},
		},
		{
			Summary: "accused of negligence resulting in a teammate's death",
			Effects: []Effect{
				benefitRolls(-2, 0, ScopeBatch),
				unimplemented("lose any company shares"),
			},
		},
		{
			Summary: "your boss believes you have been stealing materials",
			Effects: []Effect{pick("defend yourself however you can",
				opt("Advocate (Legal)", checkSkill("Advocate", 8,
					[]Effect{unimplemented("lose any company shares")},
					craftPrison())),
				opt("Deception (Lie)", checkSkill("Deception", 8,
					[]Effect{unimplemented("lose any company shares")},
					craftPrison())))},
		},
		{
			Summary: "a construction mistake kills a family of four, and you are blamed",
			Effects: []Effect{pick("defend yourself however you can",
				opt("Advocate (Legal)", checkSkill("Advocate", 8,
					craftBlamed(), append(craftBlamed(), transfer("Prisoner", "Prisoner", 1)))),
				opt("Advocate (Oratory)", checkSkill("Advocate", 8,
					craftBlamed(), append(craftBlamed(), transfer("Prisoner", "Prisoner", 1)))))},
		},
		{
			Summary: "the site is attacked, and the projects end either way",
			Effects: []Effect{pick("get away however you can",
				opt("Gun Combat", checkSkill("Gun Combat", 8, nil, []Effect{injury(1)})),
				opt("Deception", checkSkill("Deception", 8, nil, []Effect{injury(1)})),
				opt("Stealth", checkSkill("Stealth", 8, nil, []Effect{injury(1)})))},
		},
	}
}

func craftPrison() []Effect {
	return []Effect{
		loseAllBenefits("lose every benefit roll from this career"),
		unimplemented("lose any company shares"),
		relationship(Enemy, 1, ""),
		transfer("Prisoner", "Prisoner", 1),
	}
}

func craftBlamed() []Effect {
	return []Effect{
		loseAllBenefits("lose every benefit roll from this career"),
		unimplemented("lose any company shares"),
		relationship(Enemy, 1, ""),
	}
}
