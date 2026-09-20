package career

// Scavenger is the career of pp. 265-268: collecting what others discard,
// hunting for what is old enough to be worth something, or putting it back
// the way it was.
func Scavenger() Career {
	return Career{
		Name:         "Scavenger",
		Cite:         "pp. 265-268",
		Enlistment:   &Check{Characteristic: "INT", Number: 8},
		MishapEjects: true,
		Assignments: []Assignment{
			{
				Name:        "Junker",
				Description: "collecting and storing discarded material",
				Survival:    Check{Characteristic: "END", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Junker", Rows: [6]Effect{
					skill("Admin"), skill("Mechanic"), skill("Broker"),
					skill("Streetwise"), skill("Deception", "Any"), skill("Jack of All Trades"),
				}},
				Ranks: [][]Effect{
					{skill("Mechanic")},
					nil,
					{chr("END", 1)},
					nil,
					{skill("Broker")},
					nil,
					{skill("Jack of All Trades")},
				},
			},
			{
				Name:        "Picker",
				Description: "searching out older material to sell or have restored",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Picker", Rows: [6]Effect{
					skill("Admin"), skill("Carouse"), skill("Broker"),
					skill("Persuade"), skill("Deception", "Any"), skill("Advocate", "Legal"),
				}},
				Ranks: [][]Effect{
					{skill("Broker")},
					nil,
					{skill("Persuade")},
					nil,
					{skill("Science", "History")},
					nil,
					{skill("Investigate")},
				},
			},
			{
				Name:        "Restorer",
				Description: "putting dated material back into its original condition",
				Survival:    Check{Characteristic: "DEX", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Restorer", Rows: [6]Effect{
					skill("Admin"), skill("Electronics", "Any"), skill("Broker"),
					skill("Mechanic"), skill("Science", "History"), skill("Art", "Any"),
				}},
				Ranks: [][]Effect{
					{skill("Mechanic")},
					nil,
					{skill("Electronics", "Any")},
					nil,
					{skill("Broker")},
					nil,
					{skill("Investigate")},
				},
			},
		},
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), skill("Carouse"),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Science", "History"), skill("Carouse"), skill("Broker"),
				skill("Persuade"), skill("Mechanic"), skill("Language", "Any"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Art", "Any"), skill("Science", "Any"), skill("Diplomat"),
				skill("Investigate"), skill("Recon"), skill("Advocate", "Any"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 0, Other: chr("STR", 1)},
			{Cash: 0, Other: chr("END", 1)},
			{Cash: 0, Other: chr("DEX", 1)},
			{Cash: 2000, Other: chr("EDU", 1)},
			{Cash: 5000, Other: relationship(Contact, 1, "")},
			{Cash: 7500, Other: relationship(Ally, 1, "")},
			{Cash: 10000, Other: unimplemented("a rare item worth 2d6 x 100,000 credits")},
		},
		Mishaps: scavengerMishaps(),
		Events:  scavengerEvents(),
	}
}

func scavengerMishaps() MishapTable {
	loseAll := unimplemented("lose every benefit roll from this career")

	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{Summary: "the local economy will no longer support this business", Effects: nil},
		{
			Summary: "something you sold turns out to have been stolen",
			Effects: []Effect{benefitRolls(-3, 0, ScopeBatch)},
		},
		{Summary: "tired of the underbelly of civilization this career sees", Effects: nil},
		{
			Summary: "a rival who has finally pushed you out of this community",
			Effects: []Effect{relationship(Rival, 1, "")},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{Summary: "a thief breaks into your home and attacks you", Effects: []Effect{injury(1)}},
		{
			Summary: "an addiction that costs you the career",
			Effects: []Effect{
				unimplemented("choose the alcohol or drug you are addicted to"),
				chr("INT", -2),
				loseAll,
			},
		},
		{
			Summary: "a celebrity accuses you of overcharging, and the charge sticks",
			Effects: []Effect{relationship(Enemy, 1, "")},
		},
		{
			// The d6 decides only whether the DEX throw happens; the fire
			// takes everything either way.
			Summary: "a fire in your storehouse",
			Effects: []Effect{
				unimplemented("roll 1d6: on 1-2 you were present, and roll DEX 8+ " +
					"or take two rolls on the Injury table"),
				loseAll,
			},
		},
		{
			Summary: "an item in your collection said to be cursed",
			Effects: []Effect{pick("discard the item, or keep it",
				opt("discard it", benefitRolls(-3, 0, ScopeBatch)),
				opt("keep it",
					unimplemented("lose every Ally and 1d3 Contacts"),
					transfer("Vagabond", "", 0)))},
		},
	}
}
