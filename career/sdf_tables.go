package career

// militaryEventRows fills a d66 table's 41-46 span, which every commissioned
// career routes to the Military Events table (p. 121).
func militaryEventRows(table EventTable, label string) {
	for result := 41; result <= 46; result++ {
		table[result] = EventRow{
			Summary: label,
			Effects: []Effect{{Kind: EffectMilitaryEvent, Detail: "roll on the Military Events table"}},
		}
	}
}

// lifeEventRows fills a d66 table's 31-36 span, which every career's does.
func lifeEventRows(table EventTable) {
	for result := 31; result <= 36; result++ {
		table[result] = EventRow{Summary: "life event", Effects: []Effect{lifeEvent()}}
	}
}

// bribeOffered is the merchant-captain bribe, which three naval careers
// print with the same three outcomes (pp. 237, 281, 286).
func bribeOffered(check Effect) Effect {
	return pick("take the bribe or refuse",
		opt("take it", check),
		opt("refuse",
			gainRank(),
			relationship(Enemy, 1, "")),
	)
}

// groomedForHigherThings is the superior officer who takes an interest,
// printed identically in four careers (pp. 237, 281, 286, 291).
func groomedForHigherThings() []Effect {
	return []Effect{
		relationship(Ally, 1, ""),
		throwModifier("next advancement roll", 4),
		onFailure(advancementThrow,
			"failing that advancement roll loses the Ally and carries -2 to the one after",
			loseThatTie(Ally), throwModifier(advancementThrow, -2)),
	}
}

// sdfNavyEvents is the d66 table of pp. 280-281.
func sdfNavyEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "a great deal of time spent at the system's main starport",
			Effects: []Effect{pick("what the port taught you",
				opt("Carouse", skill("Carouse")),
				opt("Gambler", skill("Gambler")),
				opt("Streetwise", skill("Streetwise")))},
		},
		13: {
			Summary: "recognized by the defense force for exemplary service",
			Effects: []Effect{benefitRolls(1, 0, ScopeBatch)},
		},
		14: {
			Summary: "you spot a crewmate doing something dangerous or illegal",
			Effects: []Effect{pick("turn them in, correct them, or look away",
				opt("turn them in", throwModifier("next advancement roll", 2)),
				opt("instruct them otherwise", skill("Leadership")),
				opt("ignore it", skill("Streetwise"), relationship(Ally, 1, "")))},
		},
		15: {
			Summary: "a gambling group at the base or the main starport",
			Effects: []Effect{checkSkill("Gambler", 8,
				[]Effect{
					benefitRolls(2, 0, ScopeBatch),
					pick("what the game taught you",
						opt("Gambler", skill("Gambler")),
						opt("Persuade", skill("Persuade"))),
				},
				[]Effect{benefitRolls(-2, 0, ScopeBatch)})},
		},
		16: {Summary: "chosen for advanced training", Effects: []Effect{rollTable(AdvancedEducation)}},
		21: {Summary: "selected for cross-training with another unit", Effects: []Effect{rollOtherAssignment()}},
		22: {
			Summary: "downtime at the base, spent studying",
			Effects: []Effect{pick("what you studied",
				opt("Art", skill("Art", "Any")),
				opt("Language", skill("Language", "Any")),
				opt("Science", skill("Science", "Any")))},
		},
		23: {
			Summary: "a term spent as an instructor at the system's naval academy",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch), skill("Instruction")},
		},
		24: {
			Summary: "accused of a crime, investigated and exonerated; the accuser is unconvinced",
			Effects: []Effect{relationship(Enemy, 1, "")},
		},
		25: {
			Summary: "a barfight",
			Effects: []Effect{checkSkill("Melee", 8,
				[]Effect{relationship(Contact, 1, "")}, []Effect{injury(1)})},
		},
		26: {
			Summary: "a strict new commander, and you pull through",
			Effects: []Effect{raiseHeldSkill()},
		},
		51: {Summary: "quiet years, and a great deal of deskwork", Effects: []Effect{skill("Admin")}},
		52: {
			Summary: "combat against a major pirate group",
			Effects: []Effect{pick("fight it with what you know",
				opt("Gunner", checkSkill("Gunner", 8,
					[]Effect{throwModifier("next advancement roll", 2)}, []Effect{injury(1)})),
				opt("Pilot", checkSkill("Pilot", 8,
					[]Effect{throwModifier("next advancement roll", 2)}, []Effect{injury(1)})),
				opt("Engineer", checkSkill("Engineer", 8,
					[]Effect{throwModifier("next advancement roll", 2)}, []Effect{injury(1)})))},
		},
		53: {Summary: "extra time at the shooting range", Effects: []Effect{skill("Gun Combat", "Any")}},
		54: {
			Summary: "a short posting to another system's defense forces",
			Effects: []Effect{pick("what the posting taught you",
				opt("Carouse", skill("Carouse")),
				opt("Diplomat", skill("Diplomat")),
				opt("Language", skill("Language", "Any")),
				opt("Persuade", skill("Persuade")),
				opt("Science", skill("Science", "Any")))},
		},
		55: {
			Summary: "a rescue in the outer reaches of the system",
			Effects: []Effect{pick("reach them however you can",
				opt("Survival (Freefall)", checkSkill("Survival", 8,
					[]Effect{relationship(Contact, 1, ""), throwModifier("next advancement roll", 2)},
					[]Effect{injury(1)})),
				opt("Suit (Vacc Suit)", checkSkill("Suit", 8,
					[]Effect{relationship(Contact, 1, ""), throwModifier("next advancement roll", 2)},
					[]Effect{injury(1)})))},
		},
		56: {
			Summary: "a crewmate trapped during a battle or an accident",
			Effects: []Effect{pick("attempt the rescue or not",
				opt("attempt it", checkChr("END", 8,
					[]Effect{relationship(Contact, 1, "")}, []Effect{injury(1)})),
				opt("leave them", checkSkill("Persuade", 8,
					nil, []Effect{loseRank(1)})))},
		},
		61: {
			Summary: "personal training on the side",
			Effects: []Effect{pick("what the training gave you",
				opt("Athletics", skill("Athletics", "Any")),
				opt("Melee", skill("Melee", "Any")),
				opt("+1 STR", chr("STR", 1)),
				opt("+1 DEX", chr("DEX", 1)),
				opt("+1 END", chr("END", 1)))},
		},
		62: {
			Summary: "a merchant captain offers a bribe during an inspection",
			Effects: []Effect{bribeOffered(checkSkill("Deception", 8,
				[]Effect{cashRolls(1)},
				[]Effect{
					benefitRolls(-1, 0, ScopeBatch),
					throwModifier("next advancement roll", -2),
				}))},
		},
		63: {
			Summary: "a posting to the defense force's flagship",
			Effects: []Effect{
				benefitRolls(1, 0, ScopeBatch),
				rollTable(AssignmentSkills),
				benefitRolls(0, 1, ScopeCareer),
			},
		},
		64: {
			Summary: "a superior officer sets out to groom you for higher things",
			Effects: groomedForHigherThings(),
		},
		65: {
			Summary: "extremely talented in your field",
			Effects: []Effect{pick("what the talent earned",
				opt("a level in a skill you already have",
					raiseHeldSkill()),
				opt("an automatic advancement", advance()))},
		},
		66: {
			Summary: "great heroism in a recent action",
			Effects: []Effect{
				benefitRolls(0, 1, ScopeCareer),
				pick("where the certainty is spent",
					opt("the next survival roll", autoSuccess("next survival roll")),
					opt("the next advancement roll", advance()),
					opt("the next commission roll", autoSuccess("next commission roll"))),
			},
		},
	}

	lifeEventRows(table)
	militaryEventRows(table, "naval event")

	return table
}
