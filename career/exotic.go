package career

// Exotic is the career of pp. 190-193: sex work, in a house, as a trained
// long-term companion, or on the street.
//
// Its mishap 4 is the only one in the book whose consequence differs by
// assignment rather than by a throw: the same change in the law costs a
// hetaerae two points of CHA, a brothel worker two benefit rolls and a
// penalty on the way out, and a streetwalker everything.
func Exotic() Career {
	return Career{
		Name:           "Exotic",
		Cite:           "pp. 190-193",
		Enlistment:     &Check{Characteristic: "END", Number: 8},
		EnlistmentMods: []EnlistmentMod{apparentAgeOver40()},
		MishapEjects:   true,
		Assignments: []Assignment{
			{
				Name:        "Brothel",
				Description: "entertainment and companionship, in a house, for a limited time",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Brothel", Rows: [6]Effect{
					skill("Admin"), skill("Persuade"), skill("Broker"),
					skill("Carouse"), skill("Etiquette"), skill("Art", "Any"),
				}},
				Ranks: [][]Effect{
					{skill("Broker")},
					nil,
					{skill("Persuade")},
					nil,
					{skill("Carouse")},
					nil,
					{chr("CHA", 1)},
				},
			},
			{
				Name:        "Hetaerae",
				Description: "a trained and educated companion, for the long term",
				Survival:    Check{Characteristic: "EDU", Number: 8},
				Advancement: Check{Characteristic: "CHA", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Hetaerae", Rows: [6]Effect{
					skill("Art", "Any"), skill("Etiquette"), skill("Broker"),
					skill("Diplomat"), skill("Persuade"), skill("Science", "Any"),
				}},
				Ranks: [][]Effect{
					{skill("Etiquette")},
					{skill("Art", "Any")},
					nil,
					{skill("Carouse")},
					nil,
					{skill("Diplomat")},
					{chr("CHA", 1)},
				},
			},
			{
				Name:        "Streetwalker",
				Description: "working the street, often illegally and often not by choice",
				Survival:    Check{Characteristic: "END", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Streetwalker", Rows: [6]Effect{
					skill("Broker"), skill("Stealth"), skill("Persuade"),
					skill("Gun Combat", "Any"), skill("Melee", "Any"), skill("Streetwise"),
				}},
				Ranks: [][]Effect{
					{skill("Streetwise")},
					{skill("Survival", "Any")},
					nil,
					{skill("Persuade")},
					nil, nil,
					{skill("Broker")},
				},
			},
		},
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), skill("Carouse"),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Broker"), skill("Carouse"), skill("Deception", "Any"),
				skill("Etiquette"), skill("Persuade"), skill("Recon"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Art", "Any"), skill("Science", "Any"), skill("Gambler"),
				skill("Diplomat"), skill("Language", "Any"), skill("Electronics", "Any"),
			}},
		},
		// The smallest benefit table in the book: the top cash result is
		// 5,000 credits where most careers pay tens of thousands.
		Benefits: [7]BenefitRow{
			{Cash: 150, Other: chr("DEX", 1)},
			{Cash: 300, Other: chr("END", 1)},
			{Cash: 600, Other: chr("INT", 1)},
			{Cash: 800, Other: chr("EDU", 1)},
			{Cash: 1000, Other: relationship(Contact, 1, "")},
			{Cash: 2000, Other: relationship(Ally, 1, "")},
			{Cash: 5000, Other: chr("CHA", 2)},
		},
		Mishaps: exoticMishaps(),
		Events:  exoticEvents(),
	}
}

func exoticMishaps() MishapTable {
	loseAll := loseAllBenefits("lose every benefit roll from this career")

	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{Summary: "this is not the career for you, and you want out", Effects: nil},
		{
			Summary: "the laws here have changed, and what that costs depends on the work",
			Effects: []Effect{pick("take what your assignment loses",
				opt("Hetaerae", chr("CHA", -2)),
				opt("Brothel",
					benefitRolls(-2, 0, ScopeBatch),
					throwModifier("next enlistment roll", -2)),
				opt("Streetwalker", loseAll, transfer("Vagabond", "", 0)))},
		},
		{
			Summary: "a dangerous disease, cured, but the damage is done",
			Effects: []Effect{chr("END", -2), benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "love found with a client who wants to take you out of this life",
			Effects: []Effect{relationship(Ally, 1, "")},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "a crackdown, an assault by the officers, and deportation",
			Effects: []Effect{injury(1), newHomeworld()},
		},
		{
			Summary: "attacked by a client",
			Effects: []Effect{injury(2), newHomeworld()},
		},
		{
			Summary: "a client who is a thief, and takes everything",
			Effects: []Effect{
				benefitRolls(-3, 0, ScopeBatch),
				relationship(Enemy, 1, ""),
				injury(1),
			},
		},
		{
			Summary: "a client obsessed past the point of threats",
			Effects: []Effect{
				relationship(Enemy, 1, ""),
				loseAll,
				newHomeworld(),
				transfer("Vagabond", "", 0),
			},
		},
		{
			Summary: "an influential client, an embarrassed one, and a hit ordered on two",
			Effects: []Effect{
				relationship(Enemy, 1, ""),
				relationship(Ally, 1, ""),
				benefitRolls(-3, 0, ScopeBatch),
				newHomeworld(),
			},
		},
	}
}
