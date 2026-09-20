package career

// thiefEvents is the d66 table of pp. 294-295.
//
// Five of its rows are a jail term with a way back: the character enters
// Prisoner for a fixed number of terms and may return here afterwards. That
// return is the character's choice at the time, so the effect records the
// sentence and the term loop offers this career again like any other.
func thiefEvents() EventTable {
	prison := func(terms int) []Effect {
		return []Effect{transfer("Prisoner", "Prisoner", terms)}
	}

	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: `a "normal" job, for a while`,
			Effects: []Effect{pickSkill("Admin", "Animals", "Chef", "Drive", "Flyer", "Mechanic")},
		},
		13: {
			Summary: "what you steal still has to be sold, and the money laundered",
			Effects: []Effect{relationship(Ally, 1, ""), skill("Broker")},
		},
		14: {
			Summary: "law enforcement is onto you",
			Effects: []Effect{checkSkill("Advocate", 8,
				[]Effect{throwModifier("next advancement roll", 2)},
				prison(1))},
		},
		15: {
			Summary: "the next job calls for a disguise",
			Effects: []Effect{checkSkill("Deception", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch)},
				prison(1))},
		},
		16: {
			Summary: "sometimes you just have to blow the safe",
			Effects: []Effect{skill("Explosives")},
		},
		21: {
			Summary: "well-versed in the local laws",
			Effects: []Effect{skill("Advocate", "Legal")},
		},
		22: {
			Summary: "a popular person",
			Effects: []Effect{skill("Carouse"), relationship(Contact, 1, "1d3")},
		},
		23: {
			Summary: "the getaway driver",
			Effects: []Effect{pick("choose Drive (Wheeled) or Flyer (Grav)",
				opt("Drive (Wheeled)", skill("Drive", "Wheeled")),
				opt("Flyer (Grav)", skill("Flyer", "Grav")))},
		},
		24: {
			Summary: "a job goes wrong, and ends in a shootout where you live",
			Effects: []Effect{
				pick("get away however you can",
					opt("Gun Combat", checkSkill("Gun Combat", 8, nil, []Effect{injury(1)})),
					opt("Stealth", checkSkill("Stealth", 8, nil, []Effect{injury(1)}))),
				thiefLoseAll(),
				relationship(Enemy, 1, ""),
			},
		},
		25: {
			Summary: "a local casino you frequent",
			Effects: []Effect{checkSkill("Gambler", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch), pickSkill("Gambler", "Persuade")},
				[]Effect{benefitRolls(-2, 0, ScopeBatch)})},
		},
		// The only row whose check depends on the assignment held.
		26: {
			Summary: "a wealthy mark walks into your local starport",
			Effects: []Effect{pick("work the mark the way your assignment does",
				opt("Con Artist", checkSkill("Carouse", 8, thiefBigTake(), thiefEscape(prison(1)))),
				opt("Hacker", checkSkill("Electronics", 8, thiefBigTake(), thiefEscape(prison(1)))),
				opt("Intruder", checkSkill("Deception", 8, thiefBigTake(), thiefEscape(prison(1)))))},
		},
		41: {
			Summary: "working with a trainer",
			Effects: []Effect{pick("choose what the training built",
				opt("STR", chr("STR", 1)),
				opt("DEX", chr("DEX", 1)),
				opt("END", chr("END", 1)),
				opt("Athletics", skill("Athletics", "Any")))},
		},
		42: {
			Summary: "merchandise instead of money, and a fence to move it",
			Effects: []Effect{
				relationship(Contact, 1, ""),
				checkSkill("Broker", 8,
					[]Effect{skill("Broker"), benefitRolls(2, 0, ScopeBatch)},
					[]Effect{
						benefitRolls(-1, 0, ScopeBatch),
						unimplemented("the fence Contact becomes an Enemy"),
					}),
			},
		},
		43: {
			Summary: "a big score, and another thief who was working the same angle",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch), relationship(Rival, 1, "")},
		},
		44: {
			Summary: "time spent at the local gun range",
			Effects: []Effect{skill("Gun Combat", "Any")},
		},
		45: {
			Summary: "a wealthy individual wants their own security tested",
			Effects: []Effect{pick("decline the proposition or take it",
				opt("decline", relationship(Contact, 1, "")),
				opt("accept, by deception", checkSkill("Deception", 8,
					[]Effect{benefitRolls(2, 0, ScopeBatch), relationship(Ally, 1, "")},
					[]Effect{benefitRolls(1, 0, ScopeBatch)})),
				opt("accept, by computer", checkSkill("Electronics", 8,
					[]Effect{benefitRolls(2, 0, ScopeBatch), relationship(Ally, 1, "")},
					[]Effect{benefitRolls(1, 0, ScopeBatch)})))},
		},
		46: {
			Summary: "people from a variety of backgrounds",
			Effects: []Effect{skill("Language", "Any")},
		},
		51: {
			Summary: "an expert at being somebody else",
			Effects: []Effect{pick("choose what the act taught you",
				opt("Art (Acting)", skill("Art", "Acting")),
				opt("Deception", skill("Deception", "Disguise", "Lie")),
				opt("Persuade", skill("Persuade")))},
		},
		52: {
			Summary: "a romance begun for the passcode it would give you",
			Effects: []Effect{checkSkill("Carouse", 8,
				[]Effect{relationship(Ally, 1, "")},
				[]Effect{relationship(Enemy, 1, "")})},
		},
		53: {
			Summary: "your way around the seedy parts of the world",
			Effects: []Effect{skill("Streetwise"), relationship(Contact, 1, "1d3")},
		},
		54: {
			Summary: "a trap set for you by law enforcement",
			Effects: []Effect{checkSkill("Recon", 8, nil, prison(2))},
		},
		55: {
			Summary: "a lot of time spent in space",
			Effects: []Effect{pick("choose Suit (Vacc Suit) or Survival (Freefall)",
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
				opt("Survival (Freefall)", skill("Survival", "Freefall")))},
		},
		56: {
			Summary: "an expert at getting people to do what you want",
			Effects: []Effect{skill("Persuade")},
		},
		61: {
			Summary: "between jobs, and working a starship berth",
			Effects: []Effect{pick("choose the berth you took",
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Gunner (Turrets)", skill("Gunner", "Turrets")),
				opt("Pilot", skill("Pilot", "Any")))},
		},
		62: {
			Summary: "a big score ahead",
			Effects: []Effect{pick("set the score up however you can",
				opt("by deception", checkSkill("Deception", 8,
					[]Effect{benefitRolls(3, 0, ScopeBatch), advance()}, prison(3))),
				opt("by computer", checkSkill("Electronics", 8,
					[]Effect{benefitRolls(3, 0, ScopeBatch), advance()}, prison(3))))},
		},
		63: {
			Summary: "a dirty cop befriended",
			Effects: []Effect{
				relationship(Ally, 1, ""),
				unimplemented("while that Ally lives, ignore any Event -- but not Mishap -- " +
					"that sends you to prison, and halve every cash benefit from this career"),
			},
		},
		64: {
			Summary: "a number of odd jobs",
			Effects: []Effect{skill("Jack of All Trades")},
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

// thiefBigTake is what the mark is worth on a success (event 26).
func thiefBigTake() []Effect {
	return []Effect{benefitRolls(3, 0, ScopeBatch)}
}

// thiefEscape wraps event 26's second chance: port security is on to you,
// and Stealth 8+ is the way out of whatever failure brought them.
func thiefEscape(caught []Effect) []Effect {
	return []Effect{checkSkill("Stealth", 8, nil, caught)}
}

// thiefLoseAll is the phrase two of this career's results use, where the
// engine has no way to reach benefit rolls already banked.
func thiefLoseAll() Effect {
	return loseAllBenefits("lose every benefit roll gained to this point")
}
