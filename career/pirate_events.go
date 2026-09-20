package career

// pirateEvents is the d66 table of pp. 253-255.
func pirateEvents() EventTable {
	stranded := []Effect{pick("what the walk to civilization taught you",
		opt("Animals", skill("Animals", "Any")),
		opt("Navigation", skill("Navigation")),
		opt("Survival", skill("Survival", "Any")))}

	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "a great deal of the term spent enjoying yourself",
			Effects: []Effect{pickSkill("Carouse", "Gambler")},
		},
		13: {
			Summary: "a crewmate accused of informing the system defence force",
			Effects: []Effect{pick("ignore it, report it, or look into it",
				opt("ignore it", relationship(Contact, 1, "")),
				opt("tell your superiors", throwModifier("next advancement roll", 2)),
				opt("investigate it", checkSkill("Investigate", 8,
					[]Effect{relationship(Ally, 1, "")},
					[]Effect{relationship(Enemy, 1, "")})))},
		},
		14: {
			Summary: "downtime in Zimmspace, spent studying",
			Effects: []Effect{pickSkill("Art", "Trade", "Science")},
		},
		15: {
			Summary: "a boarding action that goes wrong, and a well-armed merchant crew",
			Effects: []Effect{pick("hold them off with what you have",
				opt("Gun Combat", checkSkill("Gun Combat", 8,
					[]Effect{relationship(Ally, 2, "")}, stranded)),
				opt("Melee", checkSkill("Melee", 8,
					[]Effect{relationship(Ally, 2, "")}, stranded)))},
		},
		16: {
			Summary: "a trophy taken from a rival pirate ship",
			Effects: []Effect{unimplemented("a weapon of the character's choice")},
		},
		21: {
			Summary: "criminal skills, honed",
			Effects: []Effect{pickSkill("Deception", "Streetwise")},
		},
		22: {
			Summary: "arrested by a world's defence forces, with a chance to defend yourself",
			Effects: []Effect{pick("defend yourself however you can",
				opt("Advocate", checkSkill("Advocate", 8,
					[]Effect{skill("Advocate", "Any")},
					[]Effect{benefitRolls(-2, 0, ScopeBatch)})),
				opt("Persuade", checkSkill("Persuade", 8,
					[]Effect{skill("Persuade")},
					[]Effect{benefitRolls(-2, 0, ScopeBatch)})))},
		},
		23: {
			Summary: "crewmates from many worlds and cultures",
			Effects: []Effect{skill("Language", "Any")},
		},
		24: {
			Summary: "your captain's habit of playing possum with the power off",
			Effects: []Effect{pick("what the trick required of you",
				opt("Survival (Freefall)", skill("Survival", "Freefall")),
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")))},
		},
		25: {
			Summary: "a mining company's ships, raided repeatedly and successfully",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch), relationship(Enemy, 1, "")},
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
			Summary: "your leader shot dead during an attack",
			Effects: []Effect{pick("press on or pull back",
				opt("continue the attack", skill("Leadership"), benefitRolls(1, 0, ScopeBatch)),
				opt("retreat", skill("Tactics", "Naval", "Military"), relationship(Ally, 2, "")))},
		},
		43: {
			Summary: "the records and the maintenance, which are most of the life",
			Effects: []Effect{pickSkill("Admin", "Broker")},
		},
		44: {
			Summary: "an ardent fan of a holovid pirate, down to the banter and the coat",
			Effects: []Effect{pick("what the devotion taught you",
				opt("Art (Acting)", skill("Art", "Acting")),
				opt("Carouse", skill("Carouse")),
				opt("Melee (Blade)", skill("Melee", "Blade")))},
		},
		45: {
			Summary: "a captured merchant carrying valuables, and a share of them",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch)},
		},
		46: {
			Summary: "a system government hiring privateers",
			Effects: []Effect{pick("take the commission or decline",
				opt("accept", relationship(Contact, 1, ""), benefitRolls(2, 0, ScopeBatch)),
				opt("decline", relationship(Enemy, 1, "")))},
		},
		51: {
			Summary: "the base attacked, and held",
			Effects: []Effect{pick("hold it with what you have",
				opt("Gun Combat", checkSkill("Gun Combat", 8, pirateHeldBase(), pirateLostBase())),
				opt("Melee", checkSkill("Melee", 8, pirateHeldBase(), pirateLostBase())),
				opt("Gunner", checkSkill("Gunner", 8, pirateHeldBase(), pirateLostBase())),
				opt("Tactics", checkSkill("Tactics", 8, pirateHeldBase(), pirateLostBase())))},
		},
		52: {
			Summary: "an unusual cargo, and a watch to keep on it",
			Effects: []Effect{pickSkill("Animals", "Science")},
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
			Summary: "a reputation, and a nickname to carry it",
			Effects: []Effect{unimplemented(
				"a nickname that follows the character: +2 CHA and a Rival if it impresses the referee, -1 CHA if not")},
		},
		55: {
			Summary: "the raided merchant turns out to be a slave ship",
			Effects: []Effect{relationship(Ally, 1, "")},
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
			Summary: "a gambling circle aboard ship or at the base",
			Effects: []Effect{checkSkill("Gambler", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch)},
				[]Effect{benefitRolls(-3, 0, ScopeBatch)})},
		},
		63: {Summary: "first aid, self-taught", Effects: []Effect{skill("Medic", "First Aid")}},
		64: {
			Summary: "a hidden compartment, and two bottles of an entrancing green liquid",
			Effects: []Effect{unimplemented(
				"roll 1d6: on 1 it is an explosive and goes off; on 2-5 a famous drink to keep " +
					"or sell on Broker 8+; on 6 the home of a strange alien animal")},
		},
		65: {
			Summary: "a famous pirate takes you on as a protege",
			Effects: []Effect{
				chr("CHA", 1),
				pickSkill("Advocate", "Broker", "Carouse", "Deception", "Diplomat",
					"Investigate", "Leadership", "Persuade", "Streetwise", "Tactics"),
			},
		},
		66: {
			Summary: "a very valuable cargo on the latest ship taken",
			Effects: []Effect{
				throwModifier("next advancement roll", 4),
				benefitRolls(3, 0, ScopeBatch),
				benefitRolls(0, 1, ScopeCareer),
			},
		},
	}

	lifeEventRows(table)

	return table
}

func pirateHeldBase() []Effect {
	return []Effect{pick("what holding it earned",
		opt("the next advancement", throwModifier("next advancement roll", 2)),
		opt("a benefit roll", benefitRolls(1, 0, ScopeBatch)))}
}

func pirateLostBase() []Effect {
	return []Effect{
		throwModifier("next advancement roll", -2),
		benefitRolls(-1, 0, ScopeBatch),
	}
}
