package career

// Instructor is the career of pp. 212-215: teaching at a university, in a
// school, or in people's homes.
//
// It has the hardest enlistment throw in the book — EDU 12+ — and the only
// one the education tracks modify: +2 for an undergraduate degree, +4 for a
// graduate one (p. 212). Both arrive with milestone 5, and until then the
// engine records that it could not apply them.
func Instructor() Career {
	ranks := func(third, fifth, sixth []Effect) [][]Effect {
		return [][]Effect{
			{pickSkill("Art", "Science")},
			{skill("Instruction")},
			{skill("Etiquette")},
			third,
			nil,
			fifth,
			sixth,
		}
	}

	return Career{
		Name:       "Instructor",
		Cite:       "pp. 212-215",
		Enlistment: &Check{Characteristic: "EDU", Number: 12},
		EnlistmentMods: []EnlistmentMod{
			undergraduateDegree(2),
			graduateDegree(4),
		},
		MishapEjects: true,
		Assignments: []Assignment{
			{
				Name:        "Professor",
				Description: "teaching at a college or university",
				Survival:    Check{Characteristic: "EDU", Number: 8},
				Advancement: Check{Characteristic: "CHA", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Professor", Rows: [6]Effect{
					skill("Advocate", "Any"), skill("Art", "Writing"), skill("Instruction"),
					skill("Science", "Any"), skill("Persuade"), skill("Art", "Any"),
				}},
				Ranks: ranks([]Effect{skill("Art", "Writing")}, nil, []Effect{skill("Admin")}),
			},
			{
				Name:        "Teacher",
				Description: "teaching children between four and eighteen",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Teacher", Rows: [6]Effect{
					skill("Electronics", "Computers"), skill("Art", "Any"), skill("Instruction"),
					skill("Science", "Any"), skill("Persuade"), skill("Jack of All Trades"),
				}},
				Ranks: ranks([]Effect{skill("Persuade")}, []Effect{skill("Admin")}, nil),
			},
			{
				Name:        "Tutor",
				Description: "a private instructor teaching people in their homes",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "CHA", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Tutor", Rows: [6]Effect{
					skill("Broker"), skill("Art", "Any"), skill("Instruction"),
					skill("Science", "Any"), skill("Persuade"),
					pick("what you drove", opt("Drive (Wheeled)", skill("Drive", "Wheeled")),
						opt("Flyer (Grav)", skill("Flyer", "Grav"))),
				}},
				Ranks: ranks([]Effect{skill("Broker")}, []Effect{skill("Persuade")}, nil),
			},
		},
		Tables: []SkillTable{
			personalDevelopment(),
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Admin"), skill("Art", "Any"), skill("Instruction"),
				skill("Science", "Any"), skill("Language", "Any"), skill("Persuade"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Diplomat"), skill("Advocate", "Any"), skill("Art", "Any"),
				skill("Science", "Any"), skill("Instruction"), skill("Investigate"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 50, Other: chr("INT", 1)},
			{Cash: 100, Other: chr("EDU", 1)},
			{Cash: 200, Other: chr("CHA", 1)},
			{Cash: 400, Other: relationship(Contact, 1, "")},
			{Cash: 800, Other: relationship(Contact, 1, "")},
			{Cash: 2000, Other: relationship(Ally, 1, "")},
			{Cash: 5000, Other: chr("EDU", 2)},
		},
		Mishaps: instructorMishaps(),
		Events:  instructorEvents(),
	}
}

func instructorMishaps() MishapTable {
	loseAll := unimplemented("lose every benefit roll from this career")

	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{Summary: "a second career, worked between lessons, is going better", Effects: nil},
		{
			Summary: "your students have not met core curriculum standards",
			Effects: []Effect{unimplemented(
				"-2 to the next enlistment roll for any career that enlists on EDU")},
		},
		{
			Summary: "budget cuts, and you are a casualty of them",
			Effects: []Effect{benefitRolls(-3, 0, ScopeBatch)},
		},
		{
			Summary: "a romance the education authority deems inappropriate",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch), relationship(Ally, 1, "")},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "the stress of the job becomes an addiction",
			Effects: []Effect{
				chr("END", -2),
				checkChr("END", 8, nil, []Effect{transfer("Vagabond", "", 0)}),
			},
		},
		{
			Summary: "a Rival accuses you of faking your students' scores, and the authority believes them",
			Effects: []Effect{relationship(Rival, 1, ""), loseAll},
		},
		{
			Summary: "a Rival accuses you of plagiarism",
			Effects: []Effect{
				relationship(Rival, 1, ""),
				pick("prove your innocence however you can",
					opt("Advocate (Legal)", checkSkill("Advocate", 8, nil,
						[]Effect{loseAll, throwModifier("next enlistment attempt", -4)})),
					opt("Art (Writing)", checkSkill("Art", 8, nil,
						[]Effect{loseAll, throwModifier("next enlistment attempt", -4)}))),
			},
		},
		{
			Summary: "the authority decides you are a crackpot teaching hokum",
			Effects: []Effect{loseAll, transfer("Vagabond", "", 0)},
		},
		{
			Summary: "your school is the scene of a mass killing",
			Effects: []Effect{pick("stop it however you can",
				opt("Stealth", checkSkill("Stealth", 8,
					[]Effect{relationship(Contact, 0, "1d6")}, []Effect{injury(2)})),
				opt("Gun Combat", checkSkill("Gun Combat", 8,
					[]Effect{relationship(Contact, 0, "1d6")}, []Effect{injury(2)})),
				opt("Persuade", checkSkill("Persuade", 8,
					[]Effect{relationship(Contact, 0, "1d6")}, []Effect{injury(2)})),
				opt("Melee", checkSkill("Melee", 8,
					[]Effect{relationship(Contact, 0, "1d6")}, []Effect{injury(2)})))},
		},
	}
}
