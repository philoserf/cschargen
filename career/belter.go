package career

// Belter is the career of pp. 160-163: mining the asteroids, finding them,
// or processing what the miners bring back.
func Belter() Career {
	return Career{
		Name:           "Belter",
		Cite:           "pp. 160-163",
		Enlistment:     &Check{Characteristic: "END", Number: 8},
		EnlistmentMods: []EnlistmentMod{apparentAgeOver40()},
		MishapEjects:   true,
		Assignments: []Assignment{
			{
				Name:        "Miner",
				Description: "mining asteroids for their ores and other wealth",
				Survival:    Check{Characteristic: "END", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Miner", Rows: [6]Effect{
					skill("Trade", "Prospector"), skill("Explosives"), skill("Suit", "Vacc Suit"),
					skill("Survival", "Freefall"), skill("Carouse"), skill("Drive", "Mole"),
				}},
				Ranks: [][]Effect{
					nil,
					{skill("Suit", "Vacc Suit")},
					{skill("Survival", "Freefall")},
					nil,
					{skill("Jack of All Trades")},
					nil,
					{skill("Pilot", "Small Craft")},
				},
			},
			{
				Name:        "Prospector",
				Description: "searching for the right asteroid to mine",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Prospector", Rows: [6]Effect{
					skill("Admin"), skill("Suit", "Vacc Suit"), skill("Trade", "Prospector"),
					skill("Science", "Geology"), skill("Survival", "Freefall"),
					skill("Electronics", "Sensors"),
				}},
				Ranks: [][]Effect{
					{skill("Trade", "Prospector")},
					{skill("Suit", "Vacc Suit")},
					nil,
					{skill("Science", "Geology")},
					nil, nil,
					{skill("Admin")},
				},
			},
			{
				Name:        "Worker",
				Description: "the support network that processes what the miners find",
				Survival:    Check{Characteristic: "END", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Worker", Rows: [6]Effect{
					skill("Admin"), skill("Streetwise"), skill("Chef"),
					skill("Electronics", "Any"), skill("Mechanic"), skill("Carouse"),
				}},
				Ranks: [][]Effect{
					nil,
					{skill("Suit", "Vacc Suit")},
					nil,
					{pickSkill("Chef", "Mechanic")},
					{skill("Streetwise")},
					nil,
					{skill("Admin")},
				},
			},
		},
		Tables: []SkillTable{
			// Carouse where every other career prints Athletics, and EDU
			// before INT rather than after.
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("EDU", 1), chr("INT", 1), skill("Carouse"),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Drive", "Mole"), skill("Trade", "Prospector"), skill("Suit", "Vacc Suit"),
				skill("Survival", "Freefall"), skill("Explosives"), skill("Pilot", "Small Craft"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Advocate", "Any"), skill("Admin"), skill("Science", "Any"),
				skill("Electronics", "Any"), pickSkill("Diplomat", "Persuade"), skill("Art", "Any"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 1000, Other: chr("STR", 1)},
			{Cash: 2000, Other: chr("END", 1)},
			{Cash: 3000, Other: chr("DEX", 1)},
			{Cash: 5000, Other: unimplemented("a weapon of the character's choice")},
			{Cash: 10000, Other: chr("CHA", 1)},
			{Cash: 25000, Other: stashItem("a vacc suit")},
			{Cash: 50000, Other: unimplemented("a company share, 2d6 x 100,000 credits")},
		},
		Mishaps: belterMishaps(),
		Events:  belterEvents(),
	}
}

func belterMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "the company closes its operations in this system",
			Effects: []Effect{unimplemented("any Company Share benefit rolled on the way out is re-rolled")},
		},
		{
			Summary: "a psychological exam finds you no longer able to handle the danger",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch)},
		},
		{
			Summary: "accused of negligence resulting in another belter's death",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{Summary: "the nothingness, the cramped quarters, the uncertain pay", Effects: nil},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "a love affair, and a jealous rival with more seniority",
			Effects: []Effect{relationship(Enemy, 1, "")},
		},
		{
			Summary: "a tunnel collapse, a slow rescue, and a claustrophobia that lasts",
			Effects: []Effect{unimplemented("lose every level of Suit (Vacc Suit)")},
		},
		{
			Summary: "an accident sets you adrift, and the psychological damage is permanent",
			Effects: []Effect{chr("INT", -2), skill("Suit", "Vacc Suit")},
		},
		{
			Summary: "a solar flare floods your workplace with radiation",
			Effects: []Effect{chr("END", -1)},
		},
		{
			Summary: "the belt is attacked, and the company leaves the system either way",
			Effects: []Effect{pick("survive it however you can",
				opt("Gun Combat", checkSkill("Gun Combat", 8, nil, []Effect{injury(1)})),
				opt("Deception", checkSkill("Deception", 8, nil, []Effect{injury(1)})),
				opt("Stealth", checkSkill("Stealth", 8, nil, []Effect{injury(1)})))},
		},
	}
}
