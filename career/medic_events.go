package career

// medicEvents is the d66 table of pp. 231-232.
func medicEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {Summary: "free time spent studying", Effects: []Effect{chr("EDU", 1)}},
		13: {
			Summary: "an insurance fraud scheme, and an invitation to join it",
			Effects: []Effect{pick("join it, ignore it, or report it",
				opt("join it", pick("how you cover it",
					opt("Deception (Lie)", checkSkill("Deception", 8,
						[]Effect{benefitRolls(4, 0, ScopeBatch)}, medicCaught())),
					opt("Admin", checkSkill("Admin", 8,
						[]Effect{benefitRolls(4, 0, ScopeBatch)}, medicCaught())))),
				opt("ignore it", relationship(Ally, 1, "")),
				opt("report them",
					relationship(Ally, 1, ""),
					throwModifier("next advancement roll", 2)))},
		},
		14: {
			Summary: "a group of gambling doctors",
			Effects: []Effect{
				checkSkill("Gambler", 8,
					[]Effect{benefitRolls(2, 0, ScopeBatch), chr("CHA", 1)},
					[]Effect{benefitRolls(-3, 0, ScopeBatch)}),
				throwModifier("next advancement roll", 2),
			},
		},
		15: {
			Summary: "documentation on a patient, left improperly filed",
			Effects: []Effect{checkSkill("Admin", 8, nil,
				[]Effect{pick("what the lapse cost",
					opt("a benefit roll", benefitRolls(-1, 0, ScopeBatch)),
					opt("the next advancement", throwModifier("next advancement roll", -2)))})},
		},
		16: {Summary: "another doctor wants your position", Effects: []Effect{relationship(Rival, 1, "")}},
		21: {Summary: "called on to help during a disaster", Effects: []Effect{skill("Survival", "Any")}},
		22: {
			Summary: "a love of vehicles discovered",
			Effects: []Effect{pickSkill("Drive", "Flyer", "Seafarer")},
		},
		23: {Summary: "volunteering at a free clinic", Effects: []Effect{skill("Streetwise")}},
		24: {
			Summary: "an assistantship to a renowned and curmudgeonly doctor",
			Effects: []Effect{pick("take it or decline",
				opt("decline"),
				opt("take it", checkSkill("Medic", 8,
					[]Effect{pick("what you learned",
						opt("Science", skill("Science", "Any")),
						opt("Interrogation (Questioning)", skill("Interrogation", "Questioning")),
						opt("Investigate", skill("Investigate")))},
					[]Effect{
						relationship(Enemy, 1, ""),
						throwModifier("next advancement roll", -2),
					})))},
		},
		25: {
			Summary: "law enforcement asks your help with drug theft at your hospital",
			Effects: []Effect{
				relationship(Contact, 1, ""),
				pick("investigate it however you can",
					opt("Investigate", checkSkill("Investigate", 8, medicCaughtThief(), medicMissedThief())),
					opt("Admin", checkSkill("Admin", 8, medicCaughtThief(), medicMissedThief()))),
			},
		},
		26: {
			Summary: "a social clique",
			Effects: []Effect{pick("what the clique gave you",
				opt("+1 CHA", chr("CHA", 1)),
				opt("Carouse", skill("Carouse")),
				opt("contacts", relationship(Contact, 0, "1d3")))},
		},
		41: {
			Summary: "patient confidentiality, violated",
			Effects: []Effect{checkSkill("Admin", 8, nil,
				[]Effect{throwModifier("next survival roll", -2)})},
		},
		42: {
			Summary: "sued for malpractice, and a lawyer to hire",
			Effects: []Effect{
				benefitRolls(-1, 0, ScopeBatch),
				relationship(Contact, 1, ""),
				rollSub("how the suit goes",
					on(1, "it is settled against you",
						benefitRolls(-3, 0, ScopeBatch),
						throwModifier(survivalThrow, -2)),
					onRange(2, 5, "it is dismissed"),
					on(6, "the plaintiff is shown to be lying",
						benefitRolls(2, 0, ScopeBatch))),
			},
		},
		43: {
			Summary: "treating a local celebrity",
			Effects: []Effect{checkSkill("Medic", 8,
				[]Effect{relationship(Contact, 1, ""), benefitRolls(2, 0, ScopeBatch)},
				[]Effect{throwModifier("next advancement roll", -4)})},
		},
		44: {
			Summary: "the life of a wealthy and thankful person, saved",
			Effects: []Effect{benefitRolls(3, 0, ScopeBatch), relationship(Contact, 1, "")},
		},
		45: {
			Summary: "hobbies, to get away from the stress",
			Effects: []Effect{pickSkill("Art", "Animals", "Drive", "Flyer", "Science", "Seafarer")},
		},
		46: {
			Summary: "a local amateur sports group",
			Effects: []Effect{checkSkill("Athletics", 8,
				[]Effect{pick("what the sport gave you",
					opt("Athletics", skill("Athletics", "Any")),
					opt("+1 STR", chr("STR", 1)),
					opt("+1 DEX", chr("DEX", 1)),
					opt("+1 END", chr("END", 1)))},
				[]Effect{injury(1)})},
		},
		51: {
			Summary: "medical journals studied, and skills improved",
			Effects: []Effect{skill("Medic", "Any")},
		},
		52: {Summary: "money invested poorly", Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)}},
		53: {
			Summary: "a lucrative offer from another company, ship or hospital",
			Effects: []Effect{pick("take it or turn it down",
				opt("take it", benefitRolls(1, 0, ScopeBatch)),
				opt("turn it down",
					relationship(Ally, 1, ""),
					throwModifier("next advancement roll", 4)))},
		},
		54: {
			Summary: "a misdiagnosis",
			Effects: []Effect{checkSkill("Medic", 8, nil, []Effect{relationship(Rival, 1, "")})},
		},
		55: {
			Summary: "an embarrassing problem treated with discretion",
			Effects: []Effect{relationship(Contact, 1, "")},
		},
		56: {
			Summary: "a self-defence course",
			Effects: []Effect{pickSkill("Gun Combat", "Melee")},
		},
		61: {Summary: "several surgeries performed", Effects: []Effect{skill("Medic", "Surgery")}},
		62: {
			// ERRATA E-1: the failure branch cites p. 136.
			Summary: "an offer to treat an infamous mafia boss with a fatal illness",
			Effects: []Effect{pick("take the case or turn it down",
				opt("turn it down", relationship(Enemy, 1, "")),
				opt("take it", checkSkill("Medic", 8,
					[]Effect{benefitRolls(1, 0, ScopeBatch), relationship(Ally, 1, "")},
					[]Effect{relationship(Enemy, 1, ""), injury(1)})))},
		},
		63: {
			Summary: "despite all protestations, you are an engineer",
			Effects: []Effect{skill("Engineer", "Any")},
		},
		64: {
			Summary: "your superiors note your dedication",
			Effects: []Effect{advance(), benefitRolls(0, 1, ScopeCareer)},
		},
		65: {Summary: "flying, taken up as a hobby", Effects: []Effect{skill("Flyer", "Any")}},
		66: {
			Summary: "findings published in a medical journal",
			Effects: []Effect{benefitRolls(1, 0, ScopeBatch)},
		},
	}

	lifeEventRows(table)

	return table
}

func medicCaught() []Effect {
	return []Effect{
		loseRank(1),
		throwModifier("next advancement roll", -2),
	}
}

func medicCaughtThief() []Effect {
	return []Effect{
		skill("Streetwise"),
		throwModifier("next advancement roll", 2),
		relationship(Enemy, 1, ""),
	}
}

func medicMissedThief() []Effect {
	return []Effect{
		skill("Admin"),
		throwModifier("next survival roll", -2),
		relationship(Rival, 1, ""),
	}
}
