package career

import "maps"

// sharedDefenceMishaps is the mishap table Troopers (p. 284) and Wet Navy
// (p. 289) print. The two differ in three rows and are otherwise the same
// table with "soldier" for "sailor"; the differences are passed in.
func sharedDefenceMishaps(pushedRecruit, tribunal []Effect, budgetCut Effect) MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "a senior officer's negative reports end your service",
			Effects: []Effect{relationship(Enemy, 1, "")},
		},
		{Summary: "budget cuts cost you your place in the force", Effects: []Effect{budgetCut}},
		{
			Summary: "a love affair, and a jealous rival with more seniority",
			Effects: []Effect{relationship(Enemy, 1, "")},
		},
		{
			Summary: "accused of negligence resulting in another's death",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "a virulent local virus leaves you permanently weakened",
			Effects: []Effect{pick("what the illness took",
				opt("STR", chr("STR", -1)),
				opt("END", chr("END", -1)))},
		},
		{Summary: "pushing a recruit past their limits kills them", Effects: pushedRecruit},
		{
			Summary: "a psychological profile deems you unfit, and stains your record",
			Effects: []Effect{throwModifier("next enlistment attempt", -2)},
		},
		{
			Summary: "a training accident kills several, and one survivor blames you",
			Effects: []Effect{injury(1), relationship(Enemy, 1, "")},
		},
		{Summary: "a tribunal finds you guilty over a friendly-fire incident", Effects: tribunal},
	}
}

func trooperMishaps() MishapTable {
	return sharedDefenceMishaps(
		[]Effect{relationship(Enemy, 1, "")},
		[]Effect{
			loseAllBenefits("lose every benefit roll from this career"),
			checkChr("CHA", 8, nil, []Effect{transfer("Prisoner", "Prisoner", 2)}),
		},
		pick("argue your case however you can",
			opt("CHA", checkChr("CHA", 8,
				[]Effect{throwModifier("next enlistment attempt", 2)}, nil)),
			opt("Persuade", checkSkill("Persuade", 8,
				[]Effect{throwModifier("next enlistment attempt", 2)}, nil))),
	)
}

func wetNavyMishaps() MishapTable {
	return sharedDefenceMishaps(
		nil,
		[]Effect{
			loseAllBenefits("lose every benefit roll from this career"),
			checkChr("CHA", 8, nil, []Effect{
				skill("Streetwise"),
				transfer("Prisoner", "Prisoner", 2),
			}),
		},
		checkChr("CHA", 8, []Effect{throwModifier("next enlistment attempt", 2)}, nil),
	)
}

// sharedDefenceEvents is the d66 table Troopers (p. 285) and Wet Navy
// (p. 290) print. Fifteen of the thirty-six results are identical; the rest
// are passed in by result.
func sharedDefenceEvents(own map[int]EventRow, prisonTerms int) EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "a group of criminals wants you to get them weapons",
			Effects: []Effect{
				skill("Deception", "Any"),
				relationship(Contact, 1, ""),
				pick("help them or refuse",
					opt("refuse", relationship(Enemy, 1, "")),
					opt("help", checkSkill("Deception", 8,
						[]Effect{credits("50000"), relationship(Ally, 1, "")},
						[]Effect{transfer("Prisoner", "Prisoner", prisonTerms)}))),
			},
		},
		13: {
			Summary: "an operetta appreciation group, drunkenly joined, and an aptitude for the piano",
			Effects: []Effect{skill("Art", "Instrument")},
		},
		14: {
			Summary: "a great deal of fun on leave",
			Effects: []Effect{
				skill("Carouse"),
				unimplemented("on a 1d6 of 1-2, an addiction to alcohol or a drug"),
			},
		},
		15: {
			Summary: "downtime spent on your hobbies",
			Effects: []Effect{pick("what you studied",
				opt("Art", skill("Art", "Any")),
				opt("Science", skill("Science", "Any")))},
		},
		21: {Summary: "time spent working as part of the base staff", Effects: []Effect{skill("Admin")}},
		22: {
			Summary: "a local religion, delved into deeply",
			Effects: []Effect{unimplemented(
				"choose or invent a religion; on a 1d6 of 6 you become deeply involved and gain Science (Philosophy) 1")},
		},
		24: {Summary: "chosen for advanced training", Effects: []Effect{rollTable(AdvancedEducation)}},
		25: {
			Summary: "a wounded squadmate in an exposed position, before reinforcements arrive",
			Effects: []Effect{pick("attempt the rescue or not",
				opt("attempt it", checkChr("END", 8,
					[]Effect{relationship(Contact, 1, "")}, []Effect{injury(1)})),
				opt("wait for reinforcements"))},
		},
		26: {
			Summary: "a love of operating a vehicle",
			Effects: []Effect{pick("what you drove",
				opt("Drive", skill("Drive", "Wheeled", "Tracked")),
				opt("Flyer (Grav)", skill("Flyer", "Grav")),
				opt("Seafarer", skill("Seafarer", "Sail", "Motorboat")))},
		},
		52: {
			Summary: "time taken to learn a second language",
			Effects: []Effect{skill("Language", "Any")},
		},
		53: {Summary: "a quick learner", Effects: []Effect{rollTable(ServiceSkills)}},
		55: {
			Summary: "a term working logistics in the purchasing department",
			Effects: []Effect{skill("Broker")},
		},
		61: {
			Summary: "a starring turn in a recruiting holovid",
			Effects: []Effect{
				relationship(Contact, 0, "1d3"),
				throwModifier("next enlistment attempt", 2),
			},
		},
		62: {
			Summary: "personal training on the side",
			Effects: []Effect{pick("what the training gave you",
				opt("Athletics", skill("Athletics", "Any")),
				opt("Melee", skill("Melee", "Any")),
				opt("+1 STR", chr("STR", 1)),
				opt("+1 DEX", chr("DEX", 1)),
				opt("+1 END", chr("END", 1)))},
		},
		64: {
			Summary: "a superior officer sets out to groom you for higher things",
			Effects: groomedForHigherThings(),
		},
		66: {
			Summary: "great heroism in a recent action",
			Effects: []Effect{pick("what the heroism earned",
				opt("a promotion", gainRank()),
				opt("a commission", commission(0)))},
		},
	}

	maps.Copy(table, own)

	lifeEventRows(table)
	militaryEventRows(table, "military event")

	return table
}

func trooperEvents() EventTable {
	return sharedDefenceEvents(map[int]EventRow{
		16: {
			Summary: "pirates or scavengers break into the installation you are guarding",
			Effects: []Effect{pick("stop them with what you have",
				opt("Gun Combat", checkSkill("Gun Combat", 8,
					[]Effect{benefitRolls(1, 0, ScopeBatch), autoSuccess("next advancement roll")},
					[]Effect{injury(1), throwModifier("next survival roll", -2)})),
				opt("Melee", checkSkill("Melee", 8,
					[]Effect{benefitRolls(1, 0, ScopeBatch), autoSuccess("next advancement roll")},
					[]Effect{injury(1), throwModifier("next survival roll", -2)})))},
		},
		23: {Summary: "sometimes the cook for your unit", Effects: []Effect{skill("Chef")}},
		51: {
			Summary: "rules and regulations, military and civilian",
			Effects: []Effect{skill("Advocate", "Legal")},
		},
		54: {
			Summary: "a crash course as a combat medic for a field exercise",
			Effects: []Effect{skill("Medic", "First Aid")},
		},
		56: {
			Summary: "a merchant captain offers a bribe at the downport",
			Effects: []Effect{bribeOffered(pick("talk your way through it",
				opt("Deception", checkSkill("Deception", 8,
					[]Effect{cashRolls(1)},
					[]Effect{benefitRolls(-1, 0, ScopeBatch), throwModifier("next advancement roll", -2)})),
				opt("Persuade", checkSkill("Persuade", 8,
					[]Effect{cashRolls(1)},
					[]Effect{benefitRolls(-1, 0, ScopeBatch), throwModifier("next advancement roll", -2)}))))},
		},
		63: {
			Summary: "a visiting senior officer, impressed, offers a transfer",
			Effects: []Effect{pick("transfer or stay",
				opt("stay where you are"),
				opt("the system navy", transfer("System Defense Forces (Navy)", "", 0)),
				opt("the wet navy", transfer("System Defense Forces (Wet Navy)", "", 0)))},
		},
		65: {
			Summary: "extremely talented in your field",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch), benefitRolls(0, 1, ScopeCareer)},
		},
	}, 3)
}

func wetNavyEvents() EventTable {
	return sharedDefenceEvents(map[int]EventRow{
		16: {
			Summary: "dangerous conditions for your ship, submarine or aircraft",
			Effects: []Effect{checkChr("END", 8,
				[]Effect{pick("what the weather taught you",
					opt("Seafarer", skill("Seafarer", "Any")),
					opt("Flyer", skill("Flyer", "Any")))},
				[]Effect{throwModifier("next advancement roll", -2)})},
		},
		23: {
			Summary: "placed in charge of communications",
			Effects: []Effect{skill("Electronics", "Comms")},
		},
		51: {
			Summary: "rules and regulations, military and civilian",
			Effects: []Effect{skill("Advocate", "Any")},
		},
		54: {
			Summary: "trained to perform first aid on yourself and others",
			Effects: []Effect{skill("Medic", "First Aid")},
		},
		56: {
			Summary: "a posting to the government's naval command office",
			Effects: []Effect{pick("what the posting taught you",
				opt("Carouse", skill("Carouse")),
				opt("Diplomat", skill("Diplomat")),
				opt("Persuade", skill("Persuade")),
				opt("Chef", skill("Chef")))},
		},
		63: {
			Summary: "a visiting senior officer, impressed, offers a transfer with your rank retained",
			Effects: []Effect{pick("transfer or stay",
				opt("stay where you are"),
				opt("the system navy", transfer("System Defense Forces (Navy)", "", 0)),
				opt("the troopers", transfer("System Defense Forces (Troopers)", "", 0)),
				opt("a national navy", transfer("National Navy", "", 0)))},
		},
		65: {
			Summary: "extremely talented in your field",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch), benefitRolls(0, 1, ScopeCareer)},
		},
	}, 2)
}
