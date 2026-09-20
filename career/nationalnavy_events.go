package career

// navyEvents is the d66 table of pp. 236-237.
func navyEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "recognized by your nation's navy for exemplary service",
			Effects: []Effect{benefitRolls(1, 0, ScopeBatch)},
		},
		13: {
			Summary: "a term spent as an instructor at the naval academy",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch), skill("Instruction")},
		},
		14: {
			Summary: "you spot a crewmate doing something dangerous or illegal",
			Effects: []Effect{pick("turn them in, correct them, or look away",
				opt("turn them in", throwModifier("next advancement roll", 2)),
				opt("instruct them otherwise", skill("Leadership")),
				opt("ignore it", skill("Streetwise"), relationship(Ally, 1, "")),
			)},
		},
		15: {
			Summary: "selected for cross-training in an alternate assignment",
			Effects: []Effect{rollOtherAssignment()},
		},
		16: {
			Summary: "a gambling group forms on board",
			Effects: []Effect{checkSkill("Gambler", 8,
				[]Effect{
					benefitRolls(2, 0, ScopeBatch),
					pick("what the game taught you",
						opt("Gambler", skill("Gambler")),
						opt("Deception (Lie)", skill("Deception", "Lie"))),
				},
				[]Effect{benefitRolls(-2, 0, ScopeBatch)})},
		},
		21: {
			Summary: "chosen for advanced training",
			Effects: []Effect{rollTable(AdvancedEducation)},
		},
		22: {
			Summary: "a boarding raid against a pirate vessel",
			Effects: []Effect{pick("board with a gun or a blade",
				opt("Gun Combat", checkSkill("Gun Combat", 8,
					[]Effect{skill("Gun Combat", "Any")}, []Effect{injury(1)})),
				opt("Melee", checkSkill("Melee", 8,
					[]Effect{skill("Melee", "Any")}, []Effect{injury(1)})),
			)},
		},
		23: {
			Summary: "a strict new commander, and you pull through",
			Effects: []Effect{unimplemented("raise a skill the character already holds")},
		},
		24: {
			Summary: "your ship joins an exploration expedition",
			Effects: []Effect{pick("what the expedition taught you",
				opt("Navigation", skill("Navigation")),
				opt("Recon", skill("Recon")),
				opt("Electronics", skill("Electronics", "Any")),
				opt("Survival", skill("Survival", "Any")),
			)},
		},
		25: {
			Summary: "accused of a crime, investigated and exonerated; the accuser is unconvinced",
			Effects: []Effect{relationship(Enemy, 1, "")},
		},
		26: {
			Summary: "a barfight",
			Effects: []Effect{checkSkill("Melee", 8,
				[]Effect{relationship(Contact, 1, "")},
				[]Effect{injury(1)})},
		},
		51: {
			Summary: "a cache of ancient coins found during a boarding action, and a suggestion to steal them",
			Effects: []Effect{pick("take the coins or refuse",
				opt("steal them", checkSkill("Deception", 8,
					[]Effect{
						pick("what the theft taught you",
							opt("Deception", skill("Deception", "Any")),
							opt("Stealth", skill("Stealth")),
							opt("Streetwise", skill("Streetwise"))),
						benefitRolls(2, 0, ScopeBatch),
					},
					[]Effect{
						relationship(Ally, 1, ""),
						unimplemented("lose two levels of rank as a result of the reprimand"),
					})),
				opt("refuse",
					throwModifier("next advancement roll", 2),
					throwModifier("next survival roll", -1),
					relationship(Enemy, 1, "")),
			)},
		},
		52: {
			Summary: "downtime on board, spent studying",
			Effects: []Effect{pick("what you studied",
				opt("Art", skill("Art", "Any")),
				opt("Science", skill("Science", "Any")),
			)},
		},
		53: {
			Summary: "a short posting to another nation's navy",
			Effects: []Effect{pick("what the posting taught you",
				opt("Carouse", skill("Carouse")),
				opt("Diplomat", skill("Diplomat")),
				opt("Language", skill("Language", "Any")),
				opt("Persuade", skill("Persuade")),
				opt("Science", skill("Science", "Any")),
			)},
		},
		54: {
			Summary: "a crewmate trapped during a battle or an accident",
			Effects: []Effect{pick("attempt the rescue or not",
				opt("attempt it", checkChr("END", 8,
					[]Effect{relationship(Contact, 1, "")},
					[]Effect{injury(1)})),
				opt("leave them", checkChr("CHA", 8,
					nil,
					[]Effect{unimplemented("a demotion in rank")})),
			)},
		},
		55: {
			Summary: "a rescue in the outer reaches of a system",
			Effects: []Effect{checkSkill("Suit", 8,
				[]Effect{relationship(Contact, 1, "")},
				[]Effect{injury(1)})},
		},
		56: {
			Summary: "a shuttle crash leaves you surviving in the wilderness",
			Effects: []Effect{
				injury(1),
				pick("what surviving taught you",
					opt("Medic (First Aid)", skill("Medic", "First Aid")),
					opt("Recon", skill("Recon")),
					opt("Survival", skill("Survival", "Any")),
					opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit"))),
			},
		},
		61: {
			Summary: "a merchant captain offers a bribe during an inspection",
			Effects: []Effect{pick("take the bribe or refuse",
				opt("take it", checkSkill("Deception", 8,
					[]Effect{cashRolls(1)},
					[]Effect{
						benefitRolls(-1, 0, ScopeBatch),
						throwModifier("next advancement roll", -2),
					})),
				opt("refuse",
					Effect{Kind: EffectRank, Detail: "an instant promotion"},
					relationship(Enemy, 1, "")),
			)},
		},
		62: {
			Summary: "a posting to the navy's most famous flagship",
			Effects: []Effect{
				throwModifier("next advancement roll", 2),
				rollTable(AssignmentSkills),
			},
		},
		63: {
			Summary: "downtime spent working on yourself",
			Effects: []Effect{rollTable(PersonalDevelopment)},
		},
		64: {
			Summary: "a superior officer sets out to groom you for higher things",
			Effects: []Effect{
				relationship(Ally, 1, ""),
				throwModifier("next advancement roll", 4),
				unimplemented(
					"failing that advancement roll loses the Ally and carries -2 to the one after"),
			},
		},
		65: {
			Summary: "you prove extremely talented in your field",
			Effects: []Effect{
				benefitRolls(0, 1, ScopeCareer),
				pick("what the talent earned",
					opt("a level in a skill you already have",
						unimplemented("raise a skill the character already holds")),
					opt("a promotion", Effect{Kind: EffectRank, Detail: "gain a rank"})),
			},
		},
		66: {
			Summary: "great heroism in battle",
			Effects: []Effect{
				pick("what the heroism earned",
					opt("a promotion", Effect{Kind: EffectRank, Detail: "gain a rank"}),
					opt("a commission", commission(0))),
			},
		},
	}

	for result := 31; result <= 36; result++ {
		table[result] = EventRow{Summary: "life event", Effects: []Effect{lifeEvent()}}
	}

	// 41-46 is the Military Events table, which p. 241 reprints as the
	// Naval Events Table with naval wording and identical mechanics.
	for result := 41; result <= 46; result++ {
		table[result] = EventRow{
			Summary: "naval event",
			Effects: []Effect{{Kind: EffectMilitaryEvent, Detail: "roll on the Military Events table"}},
		}
	}

	return table
}
