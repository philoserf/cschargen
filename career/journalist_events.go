package career

// prestigiousAward is the field's highest honour, printed the same way in
// Arts (p. 159) and Journalist (p. 223): three benefit rolls and a pool of
// +6 to spend two at a time.
func prestigiousAward() []Effect {
	return []Effect{
		benefitRolls(3, 0, ScopeBatch),
		pool(6, 2, "a +6 modifier in increments of no more than +2, until depleted",
			false, survivalThrow, advancementThrow),
	}
}

// journalistEvents is the d66 table of pp. 222-223.
func journalistEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "extra time in the gym",
			Effects: []Effect{pick("what the training gave you",
				opt("Athletics", skill("Athletics", "Any")),
				opt("+1 STR", chr("STR", 1)),
				opt("+1 DEX", chr("DEX", 1)),
				opt("+1 END", chr("END", 1)))},
		},
		13: {
			Summary: "an assignment to a world whose language you do not speak",
			Effects: []Effect{skill("Language", "Any")},
		},
		14: {Summary: "sneaking around is often necessary", Effects: []Effect{skill("Stealth")}},
		15: {
			Summary: "a great deal of time spent cultivating the contact list",
			Effects: []Effect{relationship(Contact, 0, "1d3")},
		},
		16: {
			Summary: "a gambling group at the news organization",
			Effects: []Effect{
				checkSkill("Gambler", 8,
					[]Effect{benefitRolls(2, 0, ScopeBatch), chr("CHA", 1)},
					[]Effect{benefitRolls(-3, 0, ScopeBatch)}),
				throwModifier("next advancement roll", 2),
			},
		},
		21: {Summary: "downtime spent studying", Effects: []Effect{pickSkill("Art", "Science")}},
		22: {Summary: "repairing your own equipment in the field", Effects: []Effect{skill("Mechanic")}},
		23: {
			Summary: "the seedier side of a culture, which the work requires",
			Effects: []Effect{skill("Streetwise")},
		},
		24: {
			Summary: "damaging information about a politician you get on with",
			Effects: []Effect{pick("expose it or sit on it",
				opt("expose it", pick("how you write it",
					opt("Art (Writing)", checkSkill("Art", 8, journalistScoop(), journalistBackfire())),
					opt("Advocate", checkSkill("Advocate", 8, journalistScoop(), journalistBackfire())))),
				opt("ignore it", relationship(Ally, 1, "")))},
		},
		25: {
			Summary: "a colleague seeking promotion by denigrating your work",
			Effects: denigratedByAColleague(),
		},
		26: {
			Summary: "selected for cross-training with another department",
			Effects: []Effect{rollOtherAssignment()},
		},
		41: {
			Summary: "standing in for a sick crewmember aboard a merchant ship",
			Effects: []Effect{pick("what the month taught you",
				opt("Astrogation", skill("Astrogation")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Gunner", skill("Gunner", "Any")),
				opt("Pilot", skill("Pilot", "Any")),
				opt("Electronics", skill("Electronics", "Any")))},
		},
		42: {Summary: "a self-defence class", Effects: []Effect{skill("Melee", "Any")}},
		43: {Summary: "some worlds are more difficult than others", Effects: []Effect{skill("Survival", "Any")}},
		44: {
			Summary: "a semi-fictional novel about your exploits",
			Effects: []Effect{checkSkill("Art", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch)},
				[]Effect{benefitRolls(-1, 0, ScopeBatch)})},
		},
		45: {
			Summary: "a holovid show built around your views",
			Effects: []Effect{pick("how you carry it",
				opt("Advocate", checkSkill("Advocate", 8, journalistHit(), journalistFlop())),
				opt("Art (Writing)", checkSkill("Art", 8, journalistHit(), journalistFlop())),
				opt("Persuade", checkSkill("Persuade", 8, journalistHit(), journalistFlop())))},
		},
		46: {
			Summary: "a story on a corporation's business practices, widely watched",
			Effects: []Effect{relationship(Contact, 0, "1d3"), relationship(Enemy, 1, "")},
		},
		51: {Summary: "a small bit of fame", Effects: []Effect{chr("CHA", 1)}},
		52: {
			Summary: "a turn as a political pundit on a current-events programme",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch), chr("CHA", 1)},
		},
		53: {
			Summary: "a great deal of time aboard starships",
			Effects: []Effect{pick("what the time taught you",
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics", skill("Electronics", "Any")))},
		},
		54: {
			Summary: "a friendship with a famous gonzo journalist",
			Effects: []Effect{relationship(Ally, 1, ""), pickSkill("Carouse", "Streetwise")},
		},
		55: {
			Summary: "an offer to be your world's official spokesperson",
			Effects: []Effect{pick("accept or decline",
				opt("accept",
					rollTable(ServiceSkills),
					loseTiesRolled("1d3", Contact)),
				opt("decline", relationship(Contact, 0, "1d3")))},
		},
		56: {
			Summary: "driving yourself to meet sources",
			Effects: []Effect{pickSkill("Drive", "Flyer", "Seafarer")},
		},
		61: {
			Summary: "a gun club, joined after a dangerous assignment",
			Effects: []Effect{skill("Gun Combat", "Any")},
		},
		62: {Summary: "getting around on new worlds", Effects: []Effect{skill("Navigation")}},
		63: {
			Summary: "a report that angers someone important",
			Effects: []Effect{relationship(Enemy, 1, "")},
		},
		64: {
			Summary: "accused of plagiarism",
			Effects: []Effect{pick("clear your name however you can",
				opt("Advocate", checkSkill("Advocate", 8, nil, journalistDisgraced())),
				opt("Persuade", checkSkill("Persuade", 8, nil, journalistDisgraced())))},
		},
		65: {
			Summary: "extremely talented in your field",
			Effects: []Effect{talentedInYourField(), benefitRolls(0, 1, ScopeCareer)},
		},
		66: {
			Summary: "the most prestigious journalism award on your world",
			Effects: prestigiousAward(),
		},
	}

	lifeEventRows(table)

	return table
}

func journalistScoop() []Effect {
	return []Effect{throwModifier("next advancement roll", 4), benefitRolls(2, 0, ScopeBatch)}
}

func journalistBackfire() []Effect {
	return []Effect{
		loseTies(2, Contact),
		throwModifier("next advancement roll", -2),
	}
}

func journalistHit() []Effect {
	return []Effect{
		benefitRolls(3, 1, ScopeBatch),
		throwModifier("next advancement roll", 2),
	}
}

func journalistFlop() []Effect {
	return []Effect{
		benefitRolls(-3, 0, ScopeBatch),
		throwModifier("next survival roll", -2),
	}
}

func journalistDisgraced() []Effect {
	return []Effect{
		benefitRolls(-1, 0, ScopeBatch),
		throwModifier("next advancement roll", -2),
	}
}
