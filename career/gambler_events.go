package career

// gamblerEvents is the d66 table of pp. 205-206.
func gamblerEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "a bankroll and expenses, kept track of",
			Effects: []Effect{skill("Admin")},
		},
		13: {Summary: "a friendly gambler", Effects: []Effect{skill("Carouse")}},
		14: {
			Summary: "knowing when to hold them and when to walk away",
			Effects: []Effect{pickSkill("Broker", "Carouse", "Recon", "Streetwise")},
		},
		15: {
			Summary: "a low bankroll, and work on a starship",
			Effects: []Effect{pick("choose the berth you took",
				opt("Astrogation", skill("Astrogation")),
				opt("Chef", skill("Chef")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Gunner (Turrets)", skill("Gunner", "Turrets")),
				opt("Mechanic", skill("Mechanic")),
				opt("Pilot", skill("Pilot", "Any")),
				opt("Survival (Freefall)", skill("Survival", "Freefall")),
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")))},
		},
		16: {
			Summary: "a rival, taken on at a game of chance",
			Effects: []Effect{
				relationship(Rival, 1, ""),
				checkSkill("Gambler", 8,
					[]Effect{skill("Gambler"), benefitRolls(2, 0, ScopeBatch)},
					[]Effect{loseRank(1)}),
			},
		},
		21: {
			Summary: `a book written about your "system"`,
			Effects: []Effect{
				skill("Art", "Writing"),
				checkSkill("Art", 8,
					[]Effect{benefitRolls(3, 0, ScopeBatch)},
					[]Effect{chr("CHA", -2)}),
			},
		},
		22: {
			Summary: "deceiving the opposition",
			Effects: []Effect{skill("Deception", "Lie")},
		},
		23: {
			Summary: "your own transportation",
			Effects: []Effect{pickSkill("Drive", "Flyer")},
		},
		24: {
			Summary: "downtime at the range",
			Effects: []Effect{skill("Gun Combat", "Any")},
		},
		25: {
			Summary: "establishments that draw people from all cultures",
			Effects: []Effect{skill("Language", "Any")},
		},
		26: {
			Summary: "a berth as the house gambler on a cruise ship",
			Effects: []Effect{
				relationship(Contact, 1, "1d3"),
				pickSkill("Carouse", "Diplomat", "Etiquette", "Gambler", "Persuade"),
			},
		},
		41: {
			Summary: "an activist for the right to gamble",
			Effects: []Effect{skill("Advocate", "Politics")},
		},
		42: {
			Summary: "an empty bankroll, and something illegal to fill it",
			Effects: []Effect{checkSkill("Deception", 8,
				[]Effect{gainBenefitRollsRolled("1d3")},
				[]Effect{
					benefitRolls(-3, 0, ScopeBatch),
					transfer("Prisoner", "Prisoner", 1),
					unimplemented("after that term, rejoin Gambler or Vagabond automatically, " +
						"or attempt any other career at -2 to enlist"),
				})},
		},
		43: {
			Summary: "information, which is how you know where to put the money",
			Effects: []Effect{pickSkill("Interrogation", "Investigate", "Recon")},
		},
		44: {
			Summary: "a self-defense course",
			Effects: []Effect{skill("Melee", "Any")},
		},
		45: {
			Summary: "a week in an uninhabited region, for a side bet, survived",
			Effects: []Effect{pickSkill("Navigation", "Survival")},
		},
		46: {
			Summary: "a lot of time on ships lately",
			Effects: []Effect{pick("choose Suit (Vacc Suit) or Survival (Freefall)",
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
				opt("Survival (Freefall)", skill("Survival", "Freefall")))},
		},
		51: {
			Summary: "riding, taken up as a hobby",
			Effects: []Effect{skill("Animals", "Riding")},
		},
		52: {
			Summary: "the underworld, which gambling is never far from",
			Effects: []Effect{relationship(Contact, 1, "1d3"), skill("Streetwise")},
		},
		53: {
			Summary: "a hobby decided on",
			Effects: []Effect{pickSkill("Art", "Science")},
		},
		54: {
			Summary: "a wealthy gambler with more enthusiasm than talent",
			Effects: []Effect{pick("fleece them, or do not",
				opt("fleece them", checkSkill("Gambler", 8,
					[]Effect{benefitRolls(4, 0, ScopeBatch)},
					[]Effect{benefitRolls(-2, 0, ScopeBatch), relationship(Enemy, 1, "")})),
				opt("leave them their money", relationship(Ally, 1, "")))},
		},
		55: {
			Summary: "a first aid class",
			Effects: []Effect{skill("Medic", "First Aid")},
		},
		56: {
			Summary: "three days on a small platform atop a pole, for a side bet",
			Effects: []Effect{checkChr("END", 8,
				[]Effect{skill("Survival", "Any"), benefitRolls(3, 0, ScopeBatch)},
				[]Effect{injury(1)})},
		},
		61: {
			Summary: "training",
			Effects: []Effect{pick("choose what the training built",
				opt("Athletics", skill("Athletics", "Any")),
				opt("STR", chr("STR", 1)),
				opt("DEX", chr("DEX", 1)),
				opt("END", chr("END", 1)))},
		},
		62: {
			Summary: "a government agent, and a rival government's man at the local casino",
			Effects: []Effect{pick("take the job, or decline it",
				opt("take it", checkSkill("Gambler", 8,
					[]Effect{
						benefitRolls(2, 0, ScopeBatch),
						relationship(Contact, 1, ""),
						relationship(Enemy, 1, ""),
					},
					[]Effect{injury(2), relationship(Enemy, 1, "")})),
				opt("decline it",
					modifierFor(advancementThrow, -2, 0,
						"-2 to every advancement roll for the rest of this career",
						enlistmentNarrowing{WhileInThisCareer: true}),
					relationship(Enemy, 1, "")))},
		},
		63: {
			Summary: "odd jobs, to keep the bankroll going",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		64: {
			Summary: "gambling so popular on your world that gamblers are holovid personalities",
			Effects: []Effect{
				chrRolled("CHA", "1d3", true),
				benefitRolls(2, 0, ScopeBatch),
				celebrityAsStar(),
			},
		},
		65: {
			Summary: "extremely talented in your field",
			Effects: []Effect{talentedInYourField()},
		},
		66: {
			Summary: "excellent work",
			Effects: []Effect{advance(), benefitRolls(0, 1, ScopeCareer)},
		},
	}

	lifeEventRows(table)

	return table
}
