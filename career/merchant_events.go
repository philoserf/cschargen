package career

// piratesAttack is the boarding action printed in Corporate Shipper
// (p. 179) and Independent Merchant (p. 209) with the same two branches.
func piratesAttack() Effect {
	stranded := []Effect{pick("what the walk to civilization taught you",
		opt("Animals", skill("Animals", "Any")),
		opt("Navigation", skill("Navigation")),
		opt("Survival", skill("Survival", "Any")))}

	return pick("fight them off with what you have",
		opt("Gun Combat", checkSkill("Gun Combat", 8, nil, stranded)),
		opt("Melee", checkSkill("Melee", 8, nil, stranded)),
	)
}

// distressCall is the damaged ship in the outer system, printed in both
// merchant careers with the same three choices and the same 1d6 behind the
// third.
func distressCall(ignored, reported Effect) Effect {
	return pick("answer it, report it, or ignore it",
		opt("ignore it", ignored),
		opt("report it and do nothing", reported),
		opt("help", rollSub("what the distress call was",
			onRange(1, 3, "a trap, and the boarding fight follows", injury(1)),
			onRange(4, 5, "genuine, and the captain is grateful",
				relationship(Contact, 1, "")),
			on(6, "genuine, and the captain is generous",
				relationship(Ally, 1, ""),
				credits("5000")))),
	)
}

// merchantEvents is the d66 table of pp. 209-211.
func merchantEvents() EventTable {
	table := EventTable{
		11: {Summary: "something terrible happens", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "a great deal of the term spent enjoying yourself",
			Effects: []Effect{pickSkill("Carouse", "Gambler")},
		},
		13: {
			Summary: "another trader accuses you of cheating them out of cargo",
			Effects: []Effect{relationship(Rival, 1, "")},
		},
		14: {
			Summary: "downtime in Zimmspace, spent studying",
			Effects: []Effect{pickSkill("Art", "Trade", "Science")},
		},
		15: {Summary: "pirates attack", Effects: []Effect{piratesAttack()}},
		16: {
			Summary: "a smuggling job offered",
			Effects: []Effect{pick("take the job or turn it down",
				opt("turn it down", throwModifier("next advancement roll", 2)),
				opt("take it", unimplemented(
					"a Deception or Persuade check against the target world's law level: "+
						"8+ for law 8-10 and 10+ above that, for two benefit rolls; "+
						"failing it is a term in the Prisoner career and a world closed to you")))},
		},
		21: {
			Summary: "a great deal of time sorting cargo records and flight plans",
			Effects: []Effect{skill("Admin")},
		},
		22: {
			Summary: "a barfight",
			Effects: []Effect{checkSkill("Melee", 8,
				[]Effect{relationship(Contact, 1, "")}, []Effect{injury(1)})},
		},
		23: {
			Summary: "crewmates from many worlds and cultures",
			Effects: []Effect{skill("Language", "Any")},
		},
		24: {
			Summary: "a distress call from a damaged ship in the outer system",
			Effects: []Effect{distressCall(
				throwModifier("next advancement roll", -2),
				throwModifier("next advancement roll", -1))},
		},
		25: {
			Summary: "local contacts, cultivated",
			Effects: []Effect{relationship(Contact, 0, "1d3")},
		},
		26: {
			Summary: "crew shortages, and several jobs aboard ship",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		41: {
			Summary: "extra time in the holographic gun range",
			Effects: []Effect{skill("Gun Combat", "Any")},
		},
		42: {
			Summary: "taught to ride a local animal at a colony world",
			Effects: []Effect{skill("Animals", "Riding")},
		},
		43: {
			Summary: "a holographic self-defence class in Zimmspace",
			Effects: []Effect{skill("Melee", "Unarmed Combat")},
		},
		44: {
			Summary: "a rescue in the outer reaches of a thinly defended system",
			Effects: []Effect{pick("assist or refuse",
				opt("refuse", throwModifier("next advancement roll", -2)),
				opt("assist", checkSkill("Suit", 8,
					[]Effect{relationship(Contact, 1, "")}, []Effect{injury(1)})))},
		},
		45: {
			Summary: "a favourable story in a Real Lives of Traders expose",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch), throwModifier("next advancement roll", 2)},
		},
		46: {
			Summary: "you spot a crewmate doing something dangerous or illegal",
			Effects: []Effect{spottedSomethingWrong(rollTable(ServiceSkills))},
		},
		51: {
			Summary: "the shadier side of the transport business",
			Effects: []Effect{pickSkill("Deception", "Streetwise")},
		},
		52: {
			Summary: "a love of operating a vehicle",
			Effects: []Effect{pick("what you drove",
				opt("Drive", skill("Drive", "Wheeled", "Tracked")),
				opt("Flyer (Grav)", skill("Flyer", "Grav")),
				opt("Seafarer", skill("Seafarer", "Sail", "Motorboat")))},
		},
		53: {
			Summary: "emergency repairs",
			Effects: []Effect{pick("fix it with what you know",
				opt("Engineer", checkSkill("Engineer", 8,
					[]Effect{skill("Engineer", "Any")},
					[]Effect{benefitRolls(-1, 0, ScopeBatch)})),
				opt("Mechanic", checkSkill("Mechanic", 8,
					[]Effect{skill("Mechanic")},
					[]Effect{benefitRolls(-1, 0, ScopeBatch)})))},
		},
		54: {
			Summary: "learning what lurks in the shadows",
			Effects: []Effect{skill("Recon")},
		},
		55: {
			Summary: "a production company wants a reality holovid of a trader's life",
			Effects: []Effect{pick("agree or opt out",
				opt("opt out"),
				opt("agree",
					credits("5000"),
					rollSub("how the holovid lands (p. 210)",
						on(1, "an absolute disaster, watched only for the ridicule",
							modifierFor(advancementThrow, -2, 2,
								"-2 to the next two advancement rolls",
								enlistmentNarrowing{})),
						onRange(2, 3, "a failure nobody admits to watching",
							throwModifier(advancementThrow, -2)),
						onRange(4, 5, "a success", credits("5000")),
						on(6, "a runaway success",
							benefitRolls(3, 0, ScopeBatch),
							transfer("Celebrity", "Star", 1)))))},
		},
		56: {
			Summary: "training during a long stretch in Zimmspace",
			Effects: []Effect{pick("what the training gave you",
				opt("Athletics", skill("Athletics", "Any")),
				opt("+1 STR", chr("STR", 1)),
				opt("+1 DEX", chr("DEX", 1)),
				opt("+1 END", chr("END", 1)))},
		},
		61: {
			Summary: "short-crewed, and working outside your specialty",
			Effects: []Effect{rollOtherAssignment()},
		},
		62: {
			Summary: "a profitable voyage you helped make",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch)},
		},
		63: {
			Summary: "downtime spent working on yourself",
			Effects: []Effect{rollTable(PersonalDevelopment)},
		},
		64: {
			Summary: "a gambling circle aboard ship or at a common port",
			Effects: []Effect{checkSkill("Gambler", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch)},
				[]Effect{benefitRolls(-3, 0, ScopeBatch)})},
		},
		65: {Summary: "extremely talented in your field", Effects: []Effect{talentedInYourField()}},
		66: {
			Summary: "an exclusive shipping contract, and a base to work it from",
			Effects: []Effect{
				benefitRolls(2, 0, ScopeBatch),
				benefitRolls(0, 1, ScopeCareer),
				unimplemented("at rank 4 or higher, 20,000 credits instead of the two benefit rolls"),
			},
		},
	}

	lifeEventRows(table)

	return table
}
