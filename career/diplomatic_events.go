package career

// denigratedByAColleague is the colleague seeking promotion at your
// expense, printed in Journalist (p. 222) and Diplomatic Service (p. 188).
func denigratedByAColleague() []Effect {
	return []Effect{throwModifier("next advancement roll", -2), relationship(Rival, 1, "")}
}

// theConspiracyRecruits is event 25, which event 53 also routes to: a
// secretive faction inside your own government makes an approach (p. 187).
func theConspiracyRecruits() Effect {
	return pick("accept, decline, or tell your superiors",
		opt("accept",
			benefitRolls(1, 0, ScopeBatch),
			unimplemented(
				"a +6 modifier for the rest of this service, spent in increments of up to +3")),
		opt("decline",
			relationship(Rival, 1, ""),
			throwModifier("next advancement roll", -2)),
		opt("tell your superiors", unimplemented(
			"roll 1d6: on 1 the superior is in on it, costing a rank and gaining a Rival; "+
				"on 2-5 thanks and +2 to the next advancement; on 6 an order to infiltrate, "+
				"which is Deception (Lie) 8+ for an advancement and two benefit rolls, or torture")),
	)
}

// diplomaticEvents is the d66 table of pp. 187-189.
func diplomaticEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "a new assignment, and a new world's language",
			Effects: []Effect{skill("Language", "Any")},
		},
		13: {
			Summary: "recognized by your government for exemplary service",
			Effects: []Effect{benefitRolls(1, 0, ScopeBatch)},
		},
		14: {
			Summary: "a minor cultural faux pas that threatens an important deal",
			Effects: []Effect{pick("smooth it over however you can",
				opt("Diplomat", checkSkill("Diplomat", 8,
					[]Effect{throwModifier("next advancement roll", 2)}, diplomaticFumbled())),
				opt("Persuade", checkSkill("Persuade", 8,
					[]Effect{throwModifier("next advancement roll", 2)}, diplomaticFumbled())))},
		},
		15: {
			Summary: "selected for cross-training with another department",
			Effects: []Effect{rollOtherAssignment()},
		},
		16: {
			Summary: "an amateur sports team among your colleagues",
			Effects: []Effect{pick("what the sport gave you",
				opt("Athletics", skill("Athletics", "Any")),
				opt("+1 STR", chr("STR", 1)),
				opt("+1 DEX", chr("DEX", 1)),
				opt("+1 END", chr("END", 1)))},
		},
		21: {
			Summary: "a mock debate, to prepare for the real negotiation",
			Effects: []Effect{pickSkill("Diplomat", "Advocate", "Persuade")},
		},
		22: {Summary: "selected for special training", Effects: []Effect{rollTable(AdvancedEducation)}},
		23: {
			Summary: "time in the seedier parts of a starport town",
			Effects: []Effect{pickSkill("Carouse", "Deception", "Gambler", "Streetwise")},
		},
		24: {Summary: "a gambling group among your colleagues", Effects: []Effect{gamblingCircle("Deception")}},
		25: {
			Summary: "a secretive group inside your own government makes an approach",
			Effects: []Effect{theConspiracyRecruits()},
		},
		26: {
			Summary: "terrorist awareness drills",
			Effects: []Effect{pickSkill("Gun Combat", "Recon", "Tactics")},
		},
		41: {Summary: "becoming a better cook", Effects: []Effect{skill("Chef")}},
		42: {
			Summary: "time in space between worlds",
			Effects: []Effect{pick("what the time taught you",
				opt("Suit", skill("Suit", "Any")),
				opt("Science (Space)", skill("Science", "Space")))},
		},
		43: {Summary: "a self-defence class", Effects: []Effect{skill("Melee", "Any")}},
		44: {
			Summary: "the negotiation site attacked before you arrive",
			Effects: []Effect{pickSkill("Investigate", "Medic", "Recon", "Leadership")},
		},
		45: {Summary: "a first aid class", Effects: []Effect{skill("Medic", "First Aid")}},
		46: {
			Summary: "deeply involved in a new trade agreement",
			Effects: []Effect{pickSkill("Advocate", "Broker")},
		},
		51: {Summary: "your artistic side, found in downtime", Effects: []Effect{skill("Art", "Any")}},
		52: {
			Summary: "learning to get around in another society",
			Effects: []Effect{pickSkill("Language", "Navigation", "Persuade", "Survival")},
		},
		53: {
			Summary: "political rivals at home try to stop the negotiations",
			Effects: []Effect{pick("uncover the plot with what your posting gives you",
				opt("Diplomat, for an Ambassador", checkSkill("Diplomat", 8,
					[]Effect{rollTable(AssignmentSkills)}, []Effect{theConspiracyRecruits()})),
				opt("Admin, for a Generalist", checkSkill("Admin", 8,
					[]Effect{rollTable(AssignmentSkills)}, []Effect{theConspiracyRecruits()})),
				opt("Investigate, for Security", checkSkill("Investigate", 8,
					[]Effect{rollTable(AssignmentSkills)}, []Effect{theConspiracyRecruits()})))},
		},
		54: {
			Summary: "extensive travel, and a shipboard skill picked up",
			Effects: []Effect{pick("what you learned",
				opt("Pilot", skill("Pilot", "Any")),
				opt("Electronics", skill("Electronics", "Any")),
				opt("Astrogation", skill("Astrogation")),
				opt("Gunner", skill("Gunner", "Any")),
				opt("Engineer", skill("Engineer", "Any")))},
		},
		55: {
			Summary: "a colleague seeking promotion by denigrating your work",
			Effects: denigratedByAColleague(),
		},
		56: {
			Summary: "a leak inside your group, and you are tasked to find it",
			Effects: []Effect{pick("find it however you can",
				opt("Admin", checkSkill("Admin", 8,
					[]Effect{throwModifier("next advancement roll", 2)},
					[]Effect{throwModifier("next advancement roll", -2)})),
				opt("Investigate", checkSkill("Investigate", 8,
					[]Effect{throwModifier("next advancement roll", 2)},
					[]Effect{throwModifier("next advancement roll", -2)})))},
		},
		61: {
			Summary: "improvements made in your personal life",
			Effects: []Effect{rollTable(PersonalDevelopment)},
		},
		62: {
			Summary: "a posting beside a distinguished statesperson",
			Effects: []Effect{benefitRolls(0, 1, ScopeCareer), rollTable(AssignmentSkills)},
		},
		63: {
			Summary: "learning to operate a local vehicle on time off",
			Effects: []Effect{pick("what you drove",
				opt("Drive", skill("Drive", "Any")),
				opt("Flyer", skill("Flyer", "Any")),
				opt("Seafarer", skill("Seafarer", "Any")),
				opt("Pilot (Small Craft)", skill("Pilot", "Small Craft")))},
		},
		64: {
			Summary: "a superior sets out to groom you for higher things",
			Effects: groomedForHigherThings(),
		},
		65: {Summary: "extremely talented in your field", Effects: []Effect{talentedInYourField()}},
		66: {
			Summary: "the right thing said, found or stopped at a tense moment",
			Effects: []Effect{
				throwModifier("next advancement roll", 4),
				relationship(Ally, 1, ""),
				relationship(Contact, 1, ""),
			},
		},
	}

	lifeEventRows(table)

	return table
}

func diplomaticFumbled() []Effect {
	return []Effect{
		throwModifier("next advancement roll", -4),
		benefitRolls(-1, 0, ScopeBatch),
	}
}
