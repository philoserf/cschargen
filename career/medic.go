package career

// Medic is the career of pp. 229-232: a doctor in a healthcare system,
// aboard a civilian ship, or on a first-response team.
//
// Its enlistment is EDU 11+, "dropping to EDU 5+ if the character has
// attended Medical School" (p. 229) -- the largest single modifier in the
// book, and one that arrives with higher education at milestone 5.
//
// Its mishap table is the one that does not begin with an injury: result 2
// is being caught using patients' drugs, and the injury is at 6 where every
// other career puts it at 7.
func Medic() Career {
	return Career{
		Name:           "Medic",
		Cite:           "pp. 229-232",
		Enlistment:     &Check{Characteristic: "EDU", Number: 11},
		EnlistmentMods: []EnlistmentMod{medicalSchool()},
		MishapEjects:   true,
		Assignments: []Assignment{
			{
				Name:        "Physician",
				Description: "a doctor seeing patients privately or in a healthcare system",
				Survival:    Check{Characteristic: "EDU", Number: 7},
				Advancement: Check{Characteristic: "CHA", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Physician", Rows: [6]Effect{
					skill("Broker"), skill("Carouse"), skill("Electronics", "Any"),
					skill("Medic", "Any"), skill("Investigate"), skill("Persuade"),
				}},
				Ranks: [][]Effect{
					{skill("Medic", "Diagnosis")},
					{skill("Medic", "Surgery")},
					nil,
					{skill("Broker")},
					nil,
					{skill("Medic", "Surgery")},
					{skill("Persuade")},
				},
			},
			{
				Name:        "Ship's Doctor",
				Description: "the health of a civilian vessel's crew and passengers",
				Survival:    Check{Characteristic: "EDU", Number: 7},
				Advancement: Check{Characteristic: "CHA", Number: 6},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Ship's Doctor", Rows: [6]Effect{
					skill("Suit", "Vacc Suit"), skill("Electronics", "Any"), skill("Persuade"),
					skill("Medic", "Any"), skill("Investigate"), skill("Carouse"),
				}},
				Ranks: [][]Effect{
					{skill("Medic", "First Aid")},
					{skill("Medic", "Diagnosis")},
					nil,
					{skill("Medic", "Cryogenics")},
					nil,
					{skill("Medic", "Surgery")},
					{skill("Persuade")},
				},
			},
			{
				Name:        "EMT",
				Description: "a first-response team",
				Survival:    Check{Characteristic: "END", Number: 7},
				Advancement: Check{Characteristic: "EDU", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "EMT", Rows: [6]Effect{
					skill("Survival", "Any"), skill("Discipline"), skill("Electronics", "Any"),
					skill("Medic", "Any"), skill("Investigate"), skill("Persuade"),
				}},
				Ranks: [][]Effect{
					{skill("Medic", "First Aid")},
					{skill("Medic", "Diagnosis")},
					nil,
					{skill("Discipline")},
					nil,
					{skill("Medic", "Surgery")},
					{skill("Jack of All Trades")},
				},
			},
		},
		Tables: []SkillTable{
			personalDevelopment(),
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Persuade"), skill("Investigate"), skill("Science", "Any"),
				skill("Medic", "Any"), skill("Broker"), skill("Etiquette"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Instruction"), skill("Admin"), skill("Medic", "Any"),
				skill("Diplomat"), skill("Advocate", "Any"), skill("Language", "Any"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 2500, Other: chr("DEX", 1)},
			{Cash: 5000, Other: chr("INT", 1)},
			{Cash: 7500, Other: chr("EDU", 1)},
			{Cash: 10000, Other: chr("CHA", 1)},
			{Cash: 20000, Other: stashItem("a medical kit")},
			{Cash: 50000, Other: relationship(Contact, 1, "")},
			{Cash: 100000, Other: stashValued(companyShare, companyShareValue)},
		},
		Mishaps: medicMishaps(),
		Events:  medicEvents(),
	}
}

func medicMishaps() MishapTable {
	loseAll := loseAllBenefits("lose every benefit roll from this career")

	return MishapTable{
		{
			// The one career whose mishap table does not open with an
			// injury, and whose injury result sits at 6 rather than 7.
			Summary: "caught using drugs meant for patients",
			Effects: []Effect{
				checkSkill("Advocate", 8, nil, []Effect{transfer("Prisoner", "Prisoner", 1)}),
				unimplemented("an addiction to a drug of the character's choice"),
			},
		},
		{
			Summary: "a patient dead of a simple mistake, despite your best efforts",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "you cannot see death and suffering any longer",
			Effects: []Effect{unimplemented(
				"-6 to any enlistment roll for a career involving violence")},
		},
		{
			Summary: "malpractice suits take your credentials and much of your pay",
			Effects: []Effect{benefitRolls(-3, 0, ScopeBatch)},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "the stress of the job becomes a nervous breakdown",
			Effects: []Effect{benefitRolls(-1, 0, ScopeBatch)},
		},
		{
			Summary: "perceived to have conducted illegal or immoral procedures on patients",
			Effects: []Effect{
				loseAll,
				unimplemented("-2 to every enlistment roll that involves EDU"),
			},
		},
		{
			Summary: "a deadly disease caught from a patient, survived but not shrugged off",
			Effects: []Effect{injury(1)},
		},
		{
			Summary: "a personal scandal tarnishes your reputation",
			Effects: []Effect{
				loseTie(Ally),
				throwModifier("next enlistment attempt", -2),
				benefitRolls(-2, 0, ScopeBatch),
			},
		},
		{
			Summary: "a family crisis takes you out of medicine",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "recklessness, and a patient dead of it",
			Effects: []Effect{relationship(Enemy, 1, "")},
		},
	}
}
