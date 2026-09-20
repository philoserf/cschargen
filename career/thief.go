package career

// Thief is the career of pp. 292-295: defrauding by charm, by keyboard, or
// by going in through a window.
//
// Its mishap 8 is the only result in the book that changes the mishap table
// itself: taking the informant deal means a later failed survival throw
// resolves as result 12 rather than being rolled for.
func Thief() Career {
	return Career{
		Tags:         []Tag{TagCriminal},
		Name:         "Thief",
		Cite:         "pp. 292-295",
		Enlistment:   &Check{Characteristic: "DEX", Number: 8},
		MishapEjects: true,
		Assignments: []Assignment{
			{
				Name:        "Con Artist",
				Description: "charisma, used to defraud",
				Survival:    Check{Characteristic: "EDU", Number: 8},
				Advancement: Check{Characteristic: "CHA", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Con Artist", Rows: [6]Effect{
					skill("Gambler"), skill("Etiquette"), skill("Carouse"),
					skill("Persuade"), skill("Deception", "Any"), skill("Broker"),
				}},
				Ranks: [][]Effect{
					{skill("Persuade")},
					{skill("Carouse")},
					nil,
					{skill("Etiquette")},
					nil,
					{relationship(Ally, 1, "")},
					nil,
				},
			},
			{
				Name:        "Hacker",
				Description: "electronics, used to obtain funds illegally",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Hacker", Rows: [6]Effect{
					skill("Advocate", "Any"), skill("Persuade"), skill("Electronics", "Any"),
					skill("Electronics", "Any"), skill("Carouse"), skill("Mechanic"),
				}},
				Ranks: [][]Effect{
					{skill("Electronics", "Computers")},
					nil,
					{skill("Streetwise")},
					nil,
					{skill("Electronics", "Computers")},
					nil,
					{skill("Advocate", "Any")},
				},
			},
			{
				Name:        "Intruder",
				Description: "breaking into homes and businesses",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "DEX", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Intruder", Rows: [6]Effect{
					skill("Streetwise"), skill("Recon"), skill("Deception", "Any"),
					skill("Gun Combat", "Any"), skill("Electronics", "Any"), skill("Stealth"),
				}},
				Ranks: [][]Effect{
					{skill("Deception", "Intrusion")},
					nil,
					{skill("Recon")},
					{skill("Stealth")},
					nil, nil,
					{relationship(Ally, 1, "")},
				},
			},
		},
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), skill("Carouse"),
			}},
			// The second Service Skills table with a characteristic in it.
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Persuade"), skill("Melee", "Any"), skill("Deception", "Any"),
				chr("DEX", 1), skill("Streetwise"), skill("Gun Combat", "Any"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Art", "Any"), skill("Recon"), skill("Advocate", "Legal"),
				skill("Stealth"), skill("Diplomat"), skill("Science", "Any"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 0, Other: chr("STR", 1)},
			{Cash: 0, Other: chr("END", 1)},
			{Cash: 0, Other: chr("DEX", 1)},
			{Cash: 2000, Other: chr("EDU", 1)},
			{Cash: 5000, Other: relationship(Contact, 1, "")},
			{Cash: 7500, Other: relationship(Ally, 1, "")},
			{Cash: 10000, Other: stashValued("a rare item", "2d6x100000")},
		},
		Mishaps: thiefMishaps(),
		Events:  thiefEvents(),
	}
}

func thiefMishaps() MishapTable {
	loseAll := loseAllBenefits("lose every benefit roll from this career")

	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{Summary: "something changed, and you are giving up the criminal life", Effects: nil},
		{
			Summary: "law enforcement is closing in, so you take everything and go",
			Effects: []Effect{newHomeworld()},
		},
		{
			Summary: "organized crime is taking over the area",
			Effects: []Effect{pick("join them or leave the world",
				opt("join them", transfer("Organized Crime", "", 0)),
				opt("leave", throwModifier("next enlistment attempt", -2)))},
		},
		{
			Summary: "the police finally catch up, and it is three terms",
			Effects: []Effect{
				loseAll,
				enlistmentPenalty(-2, "-2 to enter any non-criminal career",
					enlistmentNarrowing{NotTags: []Tag{TagCriminal}, Standing: true}),
				transfer("Prisoner", "Prisoner", 3),
			},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			// The one result that changes the table it is on.
			Summary: "the police offer you a deal: inform, or four terms inside",
			Effects: []Effect{pick("take the deal or take the sentence",
				opt("take the sentence", transfer("Prisoner", "Prisoner", 4)),
				opt("inform for them", unimplemented(
					"stay in the career at -2 to survival and advancement; "+
						"a failed survival roll then resolves as mishap 12 rather than being rolled")))},
		},
		{
			Summary: "a group of toughs takes everything and leaves you for dead",
			Effects: []Effect{injury(2), benefitRolls(-3, 0, ScopeBatch)},
		},
		{
			Summary: "a former associate takes their revenge",
			Effects: []Effect{
				relationship(Rival, 1, ""),
				becomes(1, Rival, Contact, Ally),
				newHomeworld(),
			},
		},
		{
			Summary: "a religious conversion, and a life of crime left behind",
			Effects: []Effect{loseAll, transfer("Clergy", "", 0)},
		},
		{
			Summary: "believed to be an informant, and the criminal community is done with you",
			Effects: []Effect{
				checkSkill("Streetwise", 8, nil, []Effect{injury(2)}),
				newHomeworld(),
				transfer("Vagabond", "", 0),
			},
		},
	}
}
