package career

// vagabondEvents is the d66 table of pp. 298-300.
// oddJobPay is the manager's parting payment, printed identically in
// Vagabond events 23, 53 and 56 (pp. 298-300): "Roll 1d6. On a roll of 1,
// you are paid 200 for your time there. On a roll of 2-5, you are paid 500.
// On a roll of 6, you are paid 500 and the manager recommends you for your
// choice of career."
func oddJobPay() Effect {
	return rollSub("the manager's parting payment",
		on(1, "paid 200 credits", credits("200")),
		onRange(2, 5, "paid 500 credits", credits("500")),
		on(6, "paid 500 credits and recommended",
			credits("500"),
			throwModifier(enlistmentThrow, 2)))
}

// farmhandPay is Vagabond event 42, which pays differently: "On a roll of
// 1, the farmer reneges on his promise of payment. On a 2-5, you are paid
// 500 credits. On a 6, the farmer pays you 750 and offers to help you find
// a job in another field."
func farmhandPay() Effect {
	return rollSub("the farmer's payment",
		on(1, "the farmer reneges on the promise of payment"),
		onRange(2, 5, "paid 500 credits", credits("500")),
		on(6, "paid 750 credits and helped to another field",
			credits("750"),
			throwModifier(enlistmentThrow, 2)))
}

func vagabondEvents() EventTable {
	table := EventTable{
		11: {Summary: "something terrible happens", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "war breaks out and the military gathers everyone it can find",
			Effects: []Effect{pick("which of the three system defence forces",
				opt("Naval", transfer("System Defense Forces (Navy)", "", 0)),
				opt("Troopers", transfer("System Defense Forces (Troopers)", "", 0)),
				opt("Wet Navy", transfer("System Defense Forces (Wet Navy)", "", 0)))},
		},
		13: {
			Summary: "you become a low-level operative in local illegal activity",
			Effects: []Effect{pick("talk your way through it or lie",
				opt("Streetwise", checkSkill("Streetwise", 8,
					[]Effect{skill("Streetwise")},
					[]Effect{transfer("Prisoner", "Prisoner", 1)})),
				opt("Deception", checkSkill("Deception", 8,
					[]Effect{skill("Streetwise")},
					[]Effect{transfer("Prisoner", "Prisoner", 1)})),
			)},
		},
		14: {
			Summary: "you find a knife",
			Effects: []Effect{
				stashItem("a knife"),
				modifierFor(skillCheckThrow, 1, 0, "+1 to Melee checks in this career",
					enlistmentNarrowing{OnSkill: "Melee", WhileInThisCareer: true}),
			},
		},
		15: {Summary: "you learn to ask others for aid", Effects: []Effect{skill("Persuade")}},
		16: {Summary: "you have defended your stash more than once", Effects: []Effect{skill("Melee", "Any")}},
		21: {Summary: "you find a discarded legal or political text", Effects: []Effect{skill("Advocate", "Any")}},
		22: {
			Summary: "time among people who speak another language",
			Effects: []Effect{skill("Language", "Any")},
		},
		23: {
			Summary: "a local shipping company needs a driver",
			Effects: []Effect{
				pick("what you drove",
					opt("Drive", skill("Drive", "Wheeled")),
					opt("Flyer", skill("Flyer", "Grav"))),
				oddJobPay(),
			},
		},
		24: {
			Summary: "a disagreement with another Destitute or Transient turns violent",
			Effects: []Effect{checkSkill("Melee", 8,
				[]Effect{relationship(Enemy, 1, "")},
				[]Effect{relationship(Enemy, 1, ""), stashItem("")})},
		},
		25: {Summary: "you learn first aid the hard way", Effects: []Effect{skill("Medic", "First Aid")}},
		26: {Summary: "you find a discarded handcomp", Effects: []Effect{stashItem("a handcomp")}},
		41: {
			Summary: "you befriend the owner of a local eating establishment",
			Effects: []Effect{relationship(Contact, 1, "")},
		},
		42: {
			Summary: "a local farmer takes you on as an extra hand",
			Effects: []Effect{
				skill("Animals", "Any"),
				farmhandPay(),
			},
		},
		43: {
			Summary: "a game of dice among the locals",
			Effects: []Effect{skill("Gambler")},
		},
		44: {Summary: "you learn to notice what is out of place", Effects: []Effect{skill("Recon")}},
		45: {
			Summary: "you find a textbook and study it when you can",
			Effects: []Effect{stashItem("a science textbook"), skill("Science", "Any")},
		},
		46: {Summary: "it is often best not to be seen or heard", Effects: []Effect{skill("Stealth")}},
		51: {Summary: "you locate some money", Effects: []Effect{credits("1d6x100")}},
		52: {
			Summary: "you keep yourself in condition when you can",
			Effects: []Effect{pick("what the effort gave you",
				opt("Athletics", skill("Athletics", "Any")),
				opt("+1 STR", chr("STR", 1)),
				opt("+1 DEX", chr("DEX", 1)),
				opt("+1 END", chr("END", 1)),
			)},
		},
		53: {
			Summary: "a vehicle repair shop wants an extra worker",
			Effects: []Effect{
				skill("Mechanic"),
				oddJobPay(),
			},
		},
		54: {
			// The specialty rule differs by assignment -- a Destitute
			// deepens one they have, a Transient takes one they do not --
			// which the engine cannot choose without knowing the
			// assignment's own survival history.
			Summary: "surviving is sometimes all one can do",
			Effects: []Effect{skill("Survival", "Any")},
		},
		55: {Summary: "you read up on how to behave in society", Effects: []Effect{skill("Etiquette")}},
		56: {
			Summary: "a local diner wants an extra worker",
			Effects: []Effect{
				skill("Chef"),
				oddJobPay(),
			},
		},
		61: {
			Summary: "you find a slug pistol and may keep it or be rid of it",
			Effects: []Effect{unimplemented(
				"keeping it means rolling 1d6: on 1, arrest for a murder you did not commit and 1d3 " +
					"terms in the Prisoner career; otherwise a pistol in varying condition")},
		},
		62: {
			Summary: "you befriend a group of other Destitutes or Transients",
			Effects: []Effect{relationship(Contact, 0, "1d3"), skill("Carouse")},
		},
		63: {
			Summary: "you find a discarded middle passage to a nearby world",
			Effects: []Effect{unimplemented(
				"take the nearest A-class port as a new homeworld, and either change assignment or " +
					"gain one of that world's background skills or languages")},
		},
		64: {Summary: "you have earned this through a hard life", Effects: []Effect{skill("Jack of All Trades")}},
		65: {
			Summary: "a continuing education programme puts you into a local university",
			Effects: []Effect{unimplemented(
				"spend the next term in higher education (p. 86); succeeding means a degree and a new career of your choice")},
		},
		66: {
			Summary: "the person living near you turns out to be a famous actor researching a role",
			Effects: []Effect{credits("1d6x100000"), transfer("Celebrity", "Star", 0)},
		},
	}

	for result := 31; result <= 36; result++ {
		table[result] = EventRow{Summary: "life event", Effects: []Effect{lifeEvent()}}
	}

	return table
}
