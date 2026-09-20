package career

// Scientist is the career of pp. 269-272: collecting raw data in the field,
// or carrying out research in a lab, a library or a workplace.
//
// Its d66 table has a blank row. Result 26 is printed with a number and
// nothing beside it (p. 271); see ERRATA E-13.
func Scientist() Career {
	return Career{
		Tags:       []Tag{TagAcademic},
		Name:       "Scientist",
		Cite:       "pp. 269-272",
		Enlistment: &Check{Characteristic: "EDU", Number: 12},
		EnlistmentMods: []EnlistmentMod{
			undergraduateDegree(2),
			graduateDegree(4),
		},
		MishapEjects: true,
		Assignments: []Assignment{
			{
				Name:        "Field Researcher",
				Description: "collecting raw data outside a lab, library or workplace",
				Survival:    Check{Characteristic: "END", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Field Researcher", Rows: [6]Effect{
					skill("Electronics", "Any"), skill("Mechanic"), skill("Science", "Any"),
					skill("Survival", "Any"), skill("Recon"),
					skill("Suit", "Hostile Environment", "Vacc Suit"),
				}},
				Ranks: [][]Effect{
					nil,
					{skill("Science", "Any")},
					{skill("Investigate")},
					{skill("Survival", "Any")},
					nil,
					{skill("Science", "Any")},
					{skill("Recon")},
				},
			},
			{
				Name:        "Researcher",
				Description: "academic or scientific research in a lab, library or workplace",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Researcher", Rows: [6]Effect{
					skill("Electronics", "Any"), skill("Diplomat"), skill("Science", "Any"),
					skill("Admin"), skill("Persuade"), skill("Instruction"),
				}},
				Ranks: [][]Effect{
					nil,
					{skill("Science", "Any")},
					{skill("Investigate")},
					{skill("Instruction")},
					nil,
					{skill("Science", "Any")},
					{skill("Diplomat")},
				},
			},
		},
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), skill("Carouse"),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Admin"), skill("Instruction"), skill("Science", "Any"),
				skill("Investigate"), skill("Language", "Any"), skill("Persuade"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Advocate", "Any"), skill("Art", "Any"), skill("Broker"),
				skill("Diplomat"), skill("Etiquette"), skill("Medic", "Any"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 2000, Other: chr("END", 1)},
			{Cash: 3000, Other: chr("INT", 1)},
			{Cash: 4000, Other: chr("EDU", 1)},
			{Cash: 5000, Other: relationship(Contact, 1, "")},
			{Cash: 7500, Other: relationship(Ally, 1, "")},
			{Cash: 10000, Other: stashValued("a tool kit", "1000")},
			{Cash: 25000, Other: credits("2d6x10000")},
		},
		Mishaps: scientistMishaps(),
		Events:  scientistEvents(),
	}
}

func scientistMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{Summary: "budget cuts force you into another profession", Effects: nil},
		{
			Summary: "your research is deemed faulty or useless, and your funder pushes you out",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "your passion for the research is simply gone",
			Effects: []Effect{checkChr("CHA", 8,
				[]Effect{throwModifier("next enlistment attempt", 2)}, nil)},
		},
		{
			Summary: "a romance, and a jealous superior who has you removed",
			Effects: []Effect{
				relationship(Ally, 1, ""),
				relationship(Enemy, 1, ""),
				benefitRolls(-2, 0, ScopeBatch),
			},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "the stress of the work becomes an addiction",
			Effects: []Effect{
				chr("END", -2),
				checkChr("END", 8, nil, []Effect{transfer("Vagabond", "", 0)}),
			},
		},
		{
			Summary: "a Rival accuses you of falsifying your research, and your patron believes it",
			Effects: []Effect{
				relationship(Rival, 1, ""),
				loseAllBenefits("lose every benefit roll from this career"),
				enlistmentPenalty(-2, "-2 to enter any academic career after this",
					enlistmentNarrowing{OnTags: []Tag{TagAcademic}, Standing: true}),
			},
		},
		{
			Summary: "your research is invalidated and you cannot accept it",
			Effects: []Effect{
				group(becomesInCareer(1, Rival, Ally), loseEveryTie(true, Contact)),
				relationship(Rival, 1, ""),
				benefitRolls(-3, 0, ScopeBatch),
			},
		},
		{Summary: "the vehicle you are travelling in crashes", Effects: []Effect{injury(2)}},
		{
			Summary: "your lab or research site is attacked",
			Effects: []Effect{pick("defend it however you can",
				opt("Gun Combat", checkSkill("Gun Combat", 8, scientistHeld(), scientistOverrun())),
				opt("Melee", checkSkill("Melee", 8, scientistHeld(), scientistOverrun())),
				opt("Stealth", checkSkill("Stealth", 8, scientistHeld(), scientistOverrun())),
				opt("Persuade", checkSkill("Persuade", 8, scientistHeld(), scientistOverrun())))},
		},
	}
}

func scientistHeld() []Effect {
	return []Effect{loseTie(Ally, Contact)}
}

func scientistOverrun() []Effect {
	return []Effect{loseTiesRolled("1d3", Ally, Contact), injury(2)}
}
