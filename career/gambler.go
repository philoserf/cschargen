package career

// Gambler is the career of pp. 202-206: making a living at the tables, on
// the sporting book, or at the track.
//
// Its mishap 9 is the only result in the book that hands the decision to the
// player rather than the dice -- "it is the decision of the player if this is
// true or not" -- and then applies the same consequences either way.
func Gambler() Career {
	return Career{
		Name:         "Gambler",
		Cite:         "pp. 202-206",
		Enlistment:   &Check{Characteristic: "INT", Number: 8},
		MishapEjects: true,
		Assignments: []Assignment{
			{
				Name:        "Casino",
				Description: "blackjack, poker and baccarat, played professionally",
				Survival:    Check{Characteristic: "END", Number: 8},
				Advancement: Check{Characteristic: "CHA", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Casino", Rows: [6]Effect{
					skill("Science", "Probability"), skill("Interrogation", "Questioning"),
					skill("Persuade"), skill("Gambler"),
					skill("Deception", "Lie"), skill("Streetwise"),
				}},
				Ranks: [][]Effect{
					{skill("Gambler")},
					{skill("Deception", "Lie")},
					nil,
					{relationship(Contact, 1, "")},
					nil, nil,
					{relationship(Ally, 1, "")},
				},
			},
			{
				Name:        "Sports",
				Description: "betting on sporting events",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "CHA", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Sports", Rows: [6]Effect{
					skill("Investigate"), skill("Tactics", "Sport"),
					skill("Persuade"), skill("Gambler"),
					skill("Streetwise"), skill("Recon"),
				}},
				Ranks: [][]Effect{
					{skill("Gambler")},
					{skill("Tactics", "Sport")},
					nil,
					{relationship(Contact, 1, "")},
					nil, nil,
					{relationship(Ally, 1, "")},
				},
			},
			{
				Name:        "Track",
				Description: "betting on racing events",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Track", Rows: [6]Effect{
					skill("Mechanic"), skill("Animals", "Any"),
					skill("Persuade"), skill("Gambler"),
					skill("Carouse"), skill("Streetwise"),
				}},
				Ranks: [][]Effect{
					{skill("Gambler")},
					{skill("Animals", "Training")},
					nil,
					{relationship(Contact, 1, "")},
					nil, nil,
					{relationship(Ally, 1, "")},
				},
			},
		},
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), skill("Athletics", "Any"),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Investigate"), skill("Broker"), skill("Carouse"),
				skill("Gambler"), skill("Deception", "Any"), skill("Recon"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Admin"), skill("Persuade"), skill("Interrogation", "Questioning"),
				skill("Science", "Psychology"), skill("Science", "Sociology"), skill("Diplomat"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 0, Other: chr("INT", 1)},
			{Cash: 0, Other: chr("EDU", 1)},
			{Cash: 0, Other: chr("CHA", 1)},
			{Cash: 5000, Other: relationship(Contact, 1, "")},
			{Cash: 10000, Other: relationship(Ally, 1, "")},
			{Cash: 20000, Other: stashItem("a weapon")},
			{Cash: 50000, Other: stashValued("a gambler kit", "1d6x1000")},
		},
		Mishaps: gamblerMishaps(),
		Events:  gamblerEvents(),
	}
}

// gamblerLoseAll is the phrase five of this career's eleven mishaps use.
func gamblerLoseAll() Effect {
	return loseAllBenefits("lose every benefit roll gained to this point")
}

// celebrityAsStar is the way out both mishap 12 and event 64 offer: the
// gambling community remembers you, and the holovids want you.
func celebrityAsStar() Effect {
	return transfer("Celebrity", "Star", 0)
}

func gamblerMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "everything on one big run, and everything lost",
			Effects: []Effect{gamblerLoseAll(), transfer("Vagabond", "", 0)},
		},
		{Summary: "the constant grind stopped being fun", Effects: nil},
		{
			Summary: "gambling has stopped being work and become a compulsion",
			Effects: []Effect{unimplemented("add one year to the character's age")},
		},
		{
			Summary: "an addiction that costs you the edge a living requires",
			Effects: []Effect{
				addiction("alcohol or a drug of the character's choice"),
				chr("INT", -2),
				gamblerLoseAll(),
			},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "intimately involved with the spouse of the house, which revokes your access",
			Effects: []Effect{
				relationship(Ally, 1, ""),
				relationship(Enemy, 1, ""),
				benefitRolls(-2, 0, ScopeBatch),
			},
		},
		{
			// The book hands this one to the player: "It is the decision of
			// the player if this is true or not." Either way the room
			// believes it, and the consequences are the same, so the engine
			// applies them without asking.
			Summary: "accused of cheating, and believed",
			Effects: []Effect{
				gamblerLoseAll(),
				becomesTheLastOne(Enemy, Contact, Ally),
			},
		},
		{
			Summary: "the local mafia that controls gambling here wants you dead",
			Effects: []Effect{injury(1), relationship(Enemy, 1, "")},
		},
		{
			Summary: "a dangerous feat, performed for a side bet, and failed",
			Effects: []Effect{injury(2)},
		},
		{
			Summary: "an attack on your establishment that leaves everyone else dead",
			Effects: []Effect{
				injury(2),
				gamblerLoseAll(),
				loseEveryTie(true, Ally, Contact, Rival, Enemy),
				pick("take the legend, or leave it",
					opt("become a Star", celebrityAsStar()),
					opt("leave the tables behind")),
			},
		},
	}
}
