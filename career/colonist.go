package career

// Colonist is the career of pp. 173-176: "A person who, for one reason or
// another, has decided to leave their home and make a new life on an
// unsettled world."
//
// It is the walking skeleton's career because it is the book's own worked
// example and because its tables reach further than any single career
// suggests: six of its eleven mishaps reassign the homeworld, two send the
// character to Vagabond, and one sends them to Celebrity.
func Colonist() Career {
	return Career{
		Name:         "Colonist",
		Cite:         "pp. 173-176",
		Enlistment:   &Check{Characteristic: "END", Number: 6},
		MishapEjects: true,
		Assignments:  colonistAssignments(),
		Tables:       colonistTables(),
		Benefits:     colonistBenefits(),
		Mishaps:      colonistMishaps(),
		Events:       colonistEvents(),
	}
}

func colonistAssignments() []Assignment {
	return []Assignment{
		{
			Name:        "Settler",
			Description: "left home to make a new life on an unsettled world",
			Survival:    Check{Characteristic: "END", Number: 7},
			Advancement: Check{Characteristic: "INT", Number: 7},
			Skills: SkillTable{
				Kind: AssignmentSkills,
				Name: "Settler",
				Rows: [6]Effect{
					skill("Recon"),
					skill("Gun Combat", "Any"),
					skill("Survival", "Any"),
					skill("Animals", "Any"),
					skill("Navigation"),
					skill("Tactics", "Military"),
				},
			},
			Ranks: [][]Effect{
				{skill("Animals", "Farming")},
				{skill("Animals", "Veterinary")},
				nil,
				{skill("Survival", "Any")},
				nil,
				{skill("Recon")},
				{skill("Jack of All Trades")},
			},
		},
		{
			Name:        "Politician",
			Description: "came to the new world to advance a political career",
			Survival:    Check{Characteristic: "INT", Number: 7},
			Advancement: Check{Characteristic: "CHA", Number: 7},
			Skills: SkillTable{
				Kind: AssignmentSkills,
				Name: "Politician",
				Rows: [6]Effect{
					skill("Admin"),
					skill("Persuade"),
					skill("Advocate", "Any"),
					skill("Deception", "Lie"),
					skill("Diplomat"),
					skill("Carouse"),
				},
			},
			Ranks: [][]Effect{
				{skill("Persuade")},
				{skill("Advocate", "Politics")},
				nil,
				{skill("Deception", "Lie")},
				{skill("Carouse")},
				nil,
				{chr("CHA", 2)},
			},
		},
		{
			Name:        "Commercial",
			Description: "came to the new world to establish a business or trade",
			Survival:    Check{Characteristic: "INT", Number: 7},
			Advancement: Check{Characteristic: "EDU", Number: 7},
			Skills: SkillTable{
				Kind: AssignmentSkills,
				Name: "Commercial",
				Rows: [6]Effect{
					skill("Advocate", "Any"),
					skill("Broker"),
					skill("Persuade"),
					skill("Trade", "Any"),
					skill("Admin"),
					skill("Deception", "Any"),
				},
			},
			Ranks: [][]Effect{
				{skill("Broker")},
				{skill("Admin")},
				nil,
				{skill("Diplomat")},
				nil,
				{skill("Carouse")},
				nil,
			},
		},
	}
}

func colonistTables() []SkillTable {
	return []SkillTable{
		{
			Kind: PersonalDevelopment,
			Name: "Personal Development",
			Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), skill("Athletics", "Any"),
			},
		},
		{
			Kind: ServiceSkills,
			Name: "Service Skills",
			Rows: [6]Effect{
				skill("Animals", "Any"),
				skill("Broker"),
				skill("Mechanic"),
				skill("Gun Combat", "Any"),
				skill("Survival", "Any"),
				skill("Chef"),
			},
		},
		{
			Kind:       AdvancedEducation,
			Name:       "Advanced Education",
			MinimumEDU: 8,
			Rows: [6]Effect{
				skill("Admin"),
				skill("Advocate", "Any"),
				skill("Etiquette"),
				skill("Language", "Any"),
				skill("Science", "Any"),
				skill("Jack of All Trades"),
			},
		},
	}
}

func colonistBenefits() [7]BenefitRow {
	return [7]BenefitRow{
		{Cash: 0, Other: chr("END", 1)},
		{Cash: 0, Other: chr("STR", 1)},
		{Cash: 500, Other: chr("INT", 1)},
		{Cash: 1000, Other: unimplemented("a weapon of the character's choice")},
		{Cash: 2000, Other: relationship(Contact, 1, "")},
		{Cash: 5000, Other: relationship(Ally, 1, "")},
		{Cash: 10000, Other: unimplemented("a land plot: 1d6 square acres, value 2d6 x 10,000 credits")},
	}
}

func colonistMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "the colonists decide you no longer have a place here",
			Effects: []Effect{
				newHomeworld(),
				relationship(Contact, 0, "1d3"),
				chooseNewCareer(),
			},
		},
		{
			Summary: "tired of the constant struggle of colony life",
			Effects: []Effect{newHomeworld(), chooseNewCareer()},
		},
		{
			Summary: "the colony has failed",
			Effects: []Effect{newHomeworld(), chooseNewCareer()},
		},
		{
			Summary: "a rivalry with a fellow colonist ends with you leaving",
			Effects: []Effect{newHomeworld(), chooseNewCareer()},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "accused of negligence resulting in a colonist's death",
			Effects: []Effect{
				benefitRolls(-2, 0, ScopeBatch),
				newHomeworld(),
				chooseNewCareer(),
			},
		},
		{
			// ERRATA E-1: the book cites p. 136 here; the Injury table is
			// on p. 119.
			Summary: "an alien creature nearly kills you",
			Effects: []Effect{injury(2), newHomeworld()},
		},
		{
			Summary: "banished with no money and no destination",
			Effects: []Effect{transfer("Vagabond", "Transient", 0)},
		},
		{Summary: "you contract an alien disease", Effects: []Effect{chr("END", -3)}},
		{
			Summary: "you return to find every colonist vanished without a trace",
			Effects: []Effect{
				newHomeworld(),
				skill("Survival", "Any"),
				unimplemented("lose every Ally, Contact, Rival and Enemy on this world"),
				transfer("Vagabond", "", 0),
			},
		},
	}
}
