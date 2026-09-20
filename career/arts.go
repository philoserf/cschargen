package career

// Arts is the career of pp. 155-159: making fine art, collecting it, or
// judging it.
//
// Its event 63 is the first result in the book to name a difficulty rather
// than a target number -- a Diplomat check "at Very Difficult" -- which is
// the Core Rulebook's task system and outside this ruleset (ERRATA E-12).
func Arts() Career {
	return Career{
		Name:         "Arts",
		Cite:         "pp. 155-159",
		Enlistment:   &Check{Characteristic: "INT", Number: 8},
		MishapEjects: true,
		Assignments: []Assignment{
			{
				Name:        "Artist",
				Description: "a painter, a sculptor, or a holographer",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Artist", Rows: [6]Effect{
					skill("Advocate", "Any"), skill("Carouse"), skill("Art", "Any"),
					skill("Art", "Any"), skill("Persuade"), skill("Broker"),
				}},
				Ranks: [][]Effect{
					{skill("Art", "Any")},
					nil,
					{skill("Carouse")},
					nil,
					{skill("Admin")},
					nil,
					{skill("Broker")},
				},
			},
			{
				Name:        "Collector",
				Description: "obtaining and collecting fine art, alone or for a museum",
				Survival:    Check{Characteristic: "EDU", Number: 8},
				Advancement: Check{Characteristic: "CHA", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Collector", Rows: [6]Effect{
					skill("Carouse"), skill("Admin"), skill("Broker"),
					skill("Persuade"), skill("Art", "Any"), skill("Advocate", "Legal"),
				}},
				Ranks: [][]Effect{
					{skill("Broker")},
					{skill("Admin")},
					nil,
					{skill("Persuade")},
					nil,
					{skill("Investigate")},
					{chr("CHA", 1)},
				},
			},
			{
				Name:        "Critic",
				Description: "evaluating and analysing the work of others",
				Survival:    Check{Characteristic: "CHA", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Critic", Rows: [6]Effect{
					skill("Advocate", "Oratory"), skill("Investigate"), skill("Art", "Writing"),
					skill("Persuade"), skill("Carouse"), skill("Art", "Any"),
				}},
				Ranks: [][]Effect{
					{skill("Art", "Writing")},
					nil,
					{chr("CHA", 1)},
					nil,
					{skill("Investigate")},
					nil,
					{skill("Persuade")},
				},
			},
		},
		Tables: []SkillTable{
			// The one career so far whose Personal Development table ends
			// in +1 CHA rather than Athletics.
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), chr("CHA", 1),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Persuade"), skill("Art", "Any"), skill("Carouse"),
				skill("Language", "Any"), skill("Advocate", "Any"), skill("Investigate"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Deception", "Any"), skill("Electronics", "Any"), skill("Art", "Any"),
				skill("Investigate"), skill("Diplomat"), skill("Jack of All Trades"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 1000, Other: chr("END", 1)},
			{Cash: 2000, Other: chr("DEX", 1)},
			{Cash: 5000, Other: chr("EDU", 1)},
			{Cash: 10000, Other: chr("INT", 1)},
			{Cash: 15000, Other: chr("CHA", 1)},
			{Cash: 20000, Other: stashCountValued(3, "a piece of art", "2d6x10000")},
			{Cash: 50000, Other: stashCountValued(6, "a piece of art", "2d6x10000")},
		},
		Mishaps: artsMishaps(),
		Events:  artsEvents(),
	}
}

func artsMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{Summary: "your style has gone out of style", Effects: nil},
		{
			Summary: "a work or exhibition nobody cared about, and a rejection too hard to take",
			Effects: []Effect{chr("CHA", -1)},
		},
		{
			Summary: "your muse is dead or will have nothing more to do with you",
			Effects: []Effect{loseTie(Ally, Contact)},
		},
		{
			Summary: "popular still, but too deep in a depression to work",
			Effects: nil,
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "a scathing and untrue article destroys your reputation",
			Effects: []Effect{chr("CHA", -2), benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "a Rival frames you for a crime you did not commit",
			Effects: []Effect{
				transfer("Prisoner", "Prisoner", 1),
				becomesOrElse(Enemy, []Relationship{Rival}, relationship(Rival, 1, "")),
			},
		},
		{
			Summary: "an art thief breaks into your home and attacks you",
			Effects: []Effect{injury(1)},
		},
		{Summary: "the vehicle you are travelling in crashes", Effects: []Effect{injury(2)}},
		{
			Summary: "a serious disease ends your time in the art world",
			Effects: []Effect{pick("what the illness took",
				opt("STR", chrRolled("STR", "1d3", false)),
				opt("END", chrRolled("END", "1d3", false)))},
		},
	}
}
