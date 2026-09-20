package career

// FringeMarketer is the career of pp. 198-201: selling what is illegal,
// buying what was stolen, or lending against what can be sold.
//
// Two of its mishaps chain a prison term to a second career -- two terms in
// Prisoner and then Vagabond. The engine carries one transfer at a time, so
// the prison term is the transfer and the second leg is recorded.
func FringeMarketer() Career {
	return Career{
		Name:         "Fringe Marketer",
		Cite:         "pp. 198-201",
		Enlistment:   &Check{Characteristic: "INT", Number: 8},
		MishapEjects: true,
		Assignments: []Assignment{
			{
				Name:        "Black Market",
				Description: "selling illegal items",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "CHA", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Black Market", Rows: [6]Effect{
					skill("Admin"), skill("Deception", "Any"), skill("Streetwise"),
					skill("Broker"), skill("Recon"), skill("Gun Combat", "Any"),
				}},
				Ranks: [][]Effect{
					{skill("Streetwise")},
					{relationship(Contact, 1, "")},
					{skill("Recon")},
					nil, nil,
					{skillAt("Streetwise", 2)},
					nil,
				},
			},
			{
				Name:        "Fence",
				Description: "buying what was stolen, and selling it on",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Fence", Rows: [6]Effect{
					skill("Stealth"), skill("Carouse"), skill("Streetwise"),
					skill("Broker"), skill("Recon"), skill("Deception", "Any"),
				}},
				Ranks: [][]Effect{
					{skill("Streetwise")},
					{relationship(Contact, 1, "")},
					{skill("Recon")},
					nil, nil,
					{skillAt("Streetwise", 2)},
					nil,
				},
			},
			{
				Name:        "Pawnbroker",
				Description: "lending against what can be sold if the money is not repaid",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Pawnbroker", Rows: [6]Effect{
					skill("Gun Combat", "Any"), skill("Advocate", "Legal"), skill("Broker"),
					skill("Admin"), skill("Streetwise"), skill("Recon"),
				}},
				Ranks: [][]Effect{
					{skill("Broker")},
					{relationship(Contact, 1, "")},
					{skill("Recon")},
					nil,
					{relationship(Contact, 1, "")},
					nil,
					{chr("CHA", 1)},
				},
			},
		},
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), skill("Athletics", "Any"),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Admin"), skill("Streetwise"), skill("Broker"),
				skill("Carouse"), skill("Deception", "Any"), skill("Persuade"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Interrogation", "Any"), skill("Advocate", "Any"), skill("Investigate"),
				skill("Art", "Any"), skill("Diplomat"), skill("Electronics", "Any"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 0, Other: chr("STR", 1)},
			{Cash: 100, Other: chr("END", 1)},
			{Cash: 250, Other: chr("DEX", 1)},
			{Cash: 500, Other: chr("INT", 1)},
			{Cash: 1000, Other: chr("EDU", 1)},
			{Cash: 5000, Other: stashItem("a weapon")},
			{Cash: 10000, Other: unimplemented("a rare item worth 2d6 x 100,000 credits")},
		},
		Mishaps: fringeMarketerMishaps(),
		Events:  fringeMarketerEvents(),
	}
}

// eightYearsInside is mishaps 9 and 10: two terms in Prisoner, and the
// Vagabond career after them. The second leg is recorded rather than
// transferred to, because one result carries one destination.
func eightYearsInside() []Effect {
	return []Effect{
		relationship(Enemy, 1, ""),
		transfer("Prisoner", "Prisoner", 2),
		unimplemented("on leaving prison, enlist automatically in the Vagabond career"),
	}
}

func fringeMarketerMishaps() MishapTable {
	loseAll := loseAllBenefits("lose every benefit roll from this career")

	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{Summary: "a tough life, and you have had enough", Effects: nil},
		{
			// POLICY: the book branches on whether the character held a
			// career before this one, which no effect can ask. Auto mode
			// takes the first option the book prints.
			Summary: "selling on the fringe is a difficult way to make a living",
			Effects: []Effect{pick("return to your previous career, or take up the road",
				opt("return to the career you held before this one",
					unimplemented("re-enter the previous career without an enlistment roll")),
				opt("there was no previous career", transfer("Vagabond", "", 0)))},
		},
		{Summary: "the local economy will not support you any more", Effects: nil},
		{
			Summary: "the local mafia is moving in on your business",
			Effects: []Effect{pick("sell out, or fight",
				opt("sell out", benefitRolls(2, 0, ScopeBatch), newHomeworld()),
				opt("fight, with a gun", checkSkill("Gun Combat", 8,
					[]Effect{relationship(Enemy, 1, "")},
					[]Effect{loseAll, injury(2)})),
				opt("fight, hand to hand", checkSkill("Melee", 8,
					[]Effect{relationship(Enemy, 1, "")},
					[]Effect{loseAll, injury(2)})))},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "an item you sold turns out to be a poorly made fake",
			Effects: []Effect{
				chr("CHA", -2),
				unimplemented("lose every Contact and Ally gained in this career"),
			},
		},
		{
			// The success is the milder end: a career lost rather than eight
			// years of it.
			Summary: "a rival in the business snitches to the police",
			Effects: []Effect{checkSkill("Advocate", 8,
				[]Effect{loseAll, relationship(Enemy, 1, "")},
				eightYearsInside())},
		},
		{
			Summary: "tied, rightly or wrongly, to the increased arms trade in your area",
			Effects: eightYearsInside(),
		},
		{
			Summary: "bandits who burn the shop down and leave you for dead",
			Effects: []Effect{injury(2), loseAll, relationship(Enemy, 1, "")},
		},
		{
			Summary: "an item sold to an underworld figure, with his enemies' explosives in it",
			Effects: []Effect{
				loseAll,
				relationship(Enemy, 1, ""),
				transfer("Colonist", "", 0),
			},
		},
	}
}
