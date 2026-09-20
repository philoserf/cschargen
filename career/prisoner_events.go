package career

// prisonerEvents is the d66 table of pp. 263-264.
//
// Its Life Events row carries a rider no other career's does: "Treat a
// result of 3 as adding a term to your sentence" (p. 264), because Life
// Event 3's own offer -- lose two benefit rolls or spend a term in the
// Prisoner career -- means nothing to someone already there.
func prisonerEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {Summary: "a class in bookkeeping", Effects: []Effect{skill("Admin")}},
		13: {
			Summary: "you work towards an appeal",
			Effects: []Effect{checkSkill("Advocate", 8,
				[]Effect{checkSkill("Advocate", 8,
					[]Effect{unimplemented("lessen the sentence by one term; one term or less left means release")},
					[]Effect{unimplemented("add one term to the sentence")})},
				[]Effect{unimplemented("add one term to the sentence")})},
		},
		14: {
			Summary: "war breaks out and amnesty is offered to those who serve",
			Effects: []Effect{unimplemented(
				"accepting commutes the sentence and spends two terms in the System Defense Forces " +
					"troopers (p. 282); refusing continues the career")},
		},
		15: {
			Summary: "a model prisoner, freed early and penniless",
			Effects: []Effect{
				loseAllBenefits("lose every benefit roll from this career"),
				transfer("Vagabond", "", 0),
			},
		},
		16: {
			Summary: "trouble inside costs you an additional term",
			Effects: []Effect{unimplemented("add one term to the sentence, and an Enemy if you were framed")},
		},
		21: {
			Summary: "you are allowed study time",
			Effects: []Effect{pick("what you studied",
				opt("Advocate", skill("Advocate", "Any")),
				opt("Language", skill("Language", "Any")),
				opt("Science", skill("Science", "Any")),
			)},
		},
		22: {
			Summary: "you make friends with several of your cell block",
			Effects: []Effect{unimplemented(
				"gain 1d3+1 Contacts, each becoming an Ally on a 1d6 roll of 5-6")},
		},
		23: {
			Summary: "you find a knife",
			Effects: []Effect{
				stashItem("a knife"),
				throwModifier("Melee checks in this career", 1),
			},
		},
		24: {
			Summary: "a disagreement with a fellow inmate",
			Effects: []Effect{checkSkill("Melee", 8,
				[]Effect{relationship(Enemy, 1, "")},
				[]Effect{relationship(Enemy, 1, ""), stashItem("")})},
		},
		25: {Summary: "you learn to patch yourself up", Effects: []Effect{skill("Medic", "First Aid")}},
		26: {
			// ERRATA E-1: the failure branch cites p. 136.
			Summary: "you anger a group of prisoners and they take revenge",
			Effects: []Effect{checkSkill("Melee", 8,
				[]Effect{skill("Melee", "Any")},
				[]Effect{injury(2)})},
		},
		41: {Summary: "the prison teaches you a trade", Effects: []Effect{skill("Trade", "Any")}},
		42: {
			Summary: "time in the prison yard",
			Effects: []Effect{pick("what the exercise gave you",
				opt("+1 STR", chr("STR", 1)),
				opt("+1 DEX", chr("DEX", 1)),
				opt("+1 END", chr("END", 1)),
				opt("Athletics", skill("Athletics", "Any")),
			)},
		},
		43: {
			Summary: "you are placed on a work detail",
			Effects: []Effect{pick("what the detail was",
				opt("Animals (Farming)", skill("Animals", "Farming")),
				opt("Electronics (Electrical Repair)", skill("Electronics", "Electrical Repair")),
				opt("Engineer", skill("Engineer", "Life Support", "Power")),
				opt("Leadership", skill("Leadership")),
				opt("Mechanic", skill("Mechanic")),
				opt("Chef", skill("Chef")),
				opt("Trade (Construction)", skill("Trade", "Construction")),
			)},
		},
		44: {
			// ERRATA E-1: the empty-stash branch cites p. 136.
			Summary: "a gambling group forms, and you may bet what is in your stash",
			Effects: []Effect{unimplemented(
				"bet the stash on Gambler 8+: succeeding doubles what went in, failing loses it and " +
					"one more item, and an empty stash costs an Injury roll instead")},
		},
		45: {
			Summary: "you steal a science textbook",
			Effects: []Effect{
				stashItem("a science textbook"),
				pick("what the book gave you",
					opt("Science", skill("Science", "Any")),
					opt("+1 EDU per term the book is in the stash", chr("EDU", 1))),
			},
		},
		46: {
			Summary: "you pick up new methods for criminal activity",
			Effects: []Effect{pick("what you learned",
				opt("Deception", skill("Deception", "Any")),
				opt("Streetwise", skill("Streetwise")),
			)},
		},
		51: {Summary: "you must be sneaky", Effects: []Effect{skill("Stealth")}},
		52: {
			Summary: "you become friendly with some fellow inmates",
			Effects: []Effect{pick("drink with them or know them",
				opt("Carouse", skill("Carouse")),
				opt("Contacts", relationship(Contact, 0, "1d3")),
			)},
		},
		53: {Summary: "you find 100 credits", Effects: []Effect{credits("100")}},
		54: {
			Summary: "a guard lets you call a Contact, and keeps half of what they send",
			Effects: []Effect{credits("1d6x50")},
		},
		55: {
			Summary: "a prison riot breaks out",
			Effects: []Effect{unimplemented(
				"participating is Melee 8+ -- succeeding grants a skill, +1 CHA or a knife, failing " +
					"costs two Injury rolls and an additional term; standing aside gains a Contact among " +
					"the guards")},
		},
		56: {
			Summary: "you get into a lot of fights",
			Effects: []Effect{pick("what the fighting gave you",
				opt("Melee", skill("Melee", "Any")),
				opt("+1 STR", chr("STR", 1)),
				opt("+1 END", chr("END", 1)),
			)},
		},
		61: {Summary: "you find contraband worth selling", Effects: []Effect{credits("1d6x100")}},
		62: {
			Summary: "you do what you must to survive",
			Effects: []Effect{unimplemented("roll on any skill table from this career")},
		},
		63: {Summary: "bad situations demand mental flexibility", Effects: []Effect{skill("Jack of All Trades")}},
		64: {
			Summary: "your lawyer shows you were wrongly accused",
			Effects: []Effect{unimplemented("immediate release, and a new career of your choice")},
		},
		65: {
			Summary: "you have worked hard",
			Effects: []Effect{pick("what the work earned",
				opt("a level in a skill you already have", raiseHeldSkill()),
				opt("a promotion", gainRank()),
			)},
		},
		66: {
			Summary: "you stay on your best behaviour for the term",
			Effects: []Effect{advance(), benefitRolls(0, 1, ScopeCareer)},
		},
	}

	for result := 31; result <= 36; result++ {
		table[result] = EventRow{
			Summary: "life event, where a result of 3 adds a term to the sentence (p. 264)",
			Effects: []Effect{lifeEvent()},
		}
	}

	return table
}
