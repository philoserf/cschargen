package career

// Celebrity is the career of pp. 164-168: acting, performing, or simply
// being famous.
//
// It is the career the first three built kept sending characters to and
// could not reach: Colonist event 56 and Vagabond event 66 both end with
// "immediately enter the Celebrity career as a Star", and until now that
// ended generation.
func Celebrity() Career {
	return Career{
		Name:           "Celebrity",
		Cite:           "pp. 164-168",
		Enlistment:     &Check{Characteristic: "CHA", Number: 8},
		EnlistmentMods: []EnlistmentMod{apparentAgeOver40()},
		MishapEjects:   true,
		Assignments: []Assignment{
			{
				Name:        "Actor",
				Description: "portraying others in holovids, holonovels and on the stage",
				Survival:    Check{Characteristic: "CHA", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Actor", Rows: [6]Effect{
					skill("Art", "Holography"), skill("Art", "Writing"), skill("Art", "Acting"),
					skill("Persuade"), skill("Deception", "Disguise"), skill("Advocate", "Oratory"),
				}},
				Ranks: [][]Effect{
					{skill("Art", "Acting")},
					{skill("Art", "Writing")},
					{skill("Carouse")},
					nil,
					{skill("Advocate", "Oratory")},
					nil,
					{skill("Art", "Acting")},
				},
			},
			{
				Name:        "Musician",
				Description: "performing the music that entertains the masses",
				Survival:    Check{Characteristic: "INT", Number: 7},
				Advancement: Check{Characteristic: "CHA", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Musician", Rows: [6]Effect{
					skill("Language", "Any"), skill("Carouse"), skill("Art", "Instrument"),
					skill("Art", "Writing"), skill("Persuade"), skill("Advocate", "Oratory"),
				}},
				Ranks: [][]Effect{
					{skill("Art", "Instrument")},
					{skill("Art", "Writing")},
					{skill("Carouse")},
					{skill("Persuade")},
					nil,
					{skill("Carouse")},
					nil,
				},
			},
			{
				Name:        "Star",
				Description: "famous for being famous, and beloved by millions",
				Survival:    Check{Characteristic: "CHA", Number: 8},
				Advancement: Check{Characteristic: "CHA", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Star", Rows: [6]Effect{
					skill("Advocate", "Any"), skill("Deception", "Disguise"), skill("Diplomat"),
					skill("Carouse"), skill("Persuade"), skill("Etiquette"),
				}},
				Ranks: [][]Effect{
					{skill("Carouse")},
					{skill("Etiquette")},
					{skill("Persuade")},
					nil,
					{skill("Diplomat")},
					nil,
					{skill("Carouse")},
				},
			},
		},
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), chr("CHA", 1),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Advocate", "Oratory"), skill("Persuade"), skill("Art", "Any"),
				skill("Carouse"), skill("Etiquette"), skill("Broker"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Electronics", "Any"), skill("Language", "Any"), skill("Advocate", "Any"),
				skill("Broker"), skill("Diplomat"), skill("Science", "Any"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 0, Other: relationship(Contact, 1, "")},
			{Cash: 0, Other: relationship(Contact, 1, "")},
			{Cash: 10000, Other: relationship(Ally, 1, "")},
			{Cash: 10000, Other: chr("CHA", 1)},
			{Cash: 20000, Other: chr("CHA", 1)},
			{Cash: 50000, Other: chr("EDU", 1)},
			{Cash: 100000, Other: group(
				stashItem("a producer credit on a well-known entertainment"),
				credits("2d6x100000"),
				unimplemented("2d6 x 100 credits a year for 1d6 years"))},
		},
		Mishaps: celebrityMishaps(),
		Events:  celebrityEvents(),
	}
}

func celebrityMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{Summary: "disenchanted with fame, and wanting a simple life", Effects: nil},
		{
			Summary: "so associated with one part, one song or one moment that you cannot grow past it",
			Effects: nil,
		},
		{Summary: "a disease that affects your ability to perform", Effects: []Effect{chr("END", -2)}},
		{
			Summary: "something so offensive or embarrassing that your fame does not recover",
			Effects: []Effect{chr("CHA", -2), benefitRolls(-3, 0, ScopeBatch)},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "a fracas with a member of the production crew, and the publicity that follows",
			Effects: []Effect{benefitRolls(-3, 0, ScopeBatch)},
		},
		{
			Summary: "a customs inspection finds a great deal of an illegal substance in your luggage",
			Effects: []Effect{transfer("Prisoner", "Prisoner", 1)},
		},
		{
			Summary: "a political or religious movement your fanbase does not care for",
			Effects: []Effect{benefitRolls(-3, 0, ScopeBatch)},
		},
		{
			Summary: "a Contact or Ally dies violently, and you cannot perform in public again",
			Effects: []Effect{loseTie(Contact, Ally)},
		},
		{
			Summary: "an assassin cites your career as their influence, and the backlash is palpable",
			Effects: []Effect{
				benefitRolls(-4, 0, ScopeBatch),
				enlistmentPenalty(-2, "-2 to enter any career other than Vagabond",
					enlistmentNarrowing{NotCareers: []string{"Vagabond"}, Standing: true}),
			},
		},
	}
}
