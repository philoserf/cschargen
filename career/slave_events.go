package career

// slaveEvents is the d66 table of pp. 152-154.
//
// Its result 11 is the one in the book that does not spare the character:
// every other career's 11 says "but you are not ejected from the career",
// and this one just says "roll on the Mishap table" -- which comes to the
// same thing here, because a mishap in this career does not eject either.
func slaveEvents() EventTable {
	table := EventTable{
		11: {Summary: "something terrible happens", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "war on this world, and the military gathering everyone it can find",
			Effects: []Effect{
				unimplemented("in Military/Security this is an immediate rank increase; " +
					"otherwise a transfer to Military/Security without its service skills"),
				checkChr("STR", 8,
					[]Effect{rollTable(ServiceSkills), rollTable(ServiceSkills)},
					[]Effect{injury(2)}),
			},
		},
		13: {
			Summary: "judged worthy of helping with the logistics",
			Effects: []Effect{skill("Admin")},
		},
		14: {
			Summary: "to survive, you have to be willing to do it for yourself",
			Effects: []Effect{pickSkill("Deception", "Streetwise")},
		},
		// The longest choice in the book: four ways to answer another
		// slave's grave error, and the first costs nothing.
		15: {
			Summary: "a grave error by another of your group, and someone must be punished",
			Effects: []Effect{pick("let it proceed, look away, take it, or help",
				opt("let the punishment proceed"),
				opt("turn a blind eye",
					relationship(Enemy, 2, ""),
					throwModifier("next advancement roll", -2)),
				opt("take responsibility",
					skill("Leadership"),
					relationship(Ally, 3, ""),
					relationship(Enemy, 1, "")),
				opt("assist the taskmaster",
					throwModifier("next advancement roll", 4),
					unimplemented("gain your entire group as Enemies")))},
		},
		16: {
			Summary: "the chance to become taskmaster for your group",
			Effects: []Effect{pick("accept the job, or decline it",
				opt("accept",
					autoSuccess("next survival roll"),
					autoSuccess("next advancement roll"),
					relationship(Ally, 1, ""),
					chr("CHA", -2),
					relationship(Enemy, 0, "1d6+2")),
				opt("decline",
					unimplemented("lose one rank, retaining any benefit already gained"),
					relationship(Enemy, 1, ""),
					chr("CHA", 1),
					relationship(Ally, 0, "1d3")))},
		},
		21: {
			Summary: "a discarded law text, read during the day",
			Effects: []Effect{skill("Advocate", "Legal")},
		},
		22: {
			Summary: "a great deal of time among people who speak another language",
			Effects: []Effect{skill("Language", "Any")},
		},
		23: {
			Summary: "a spirit that still wants to speak out against the oppression",
			Effects: []Effect{pick("choose Advocate (Oratory) or Advocate (Politics)",
				opt("Advocate (Oratory)", skill("Advocate", "Oratory")),
				opt("Advocate (Politics)", skill("Advocate", "Politics")))},
		},
		24: {
			Summary: "training on a vehicle",
			Effects: []Effect{pickSkill("Drive", "Flyer", "Seafarer")},
		},
		25: {
			Summary: "keeping your own equipment running",
			Effects: []Effect{skill("Mechanic")},
		},
		26: {
			Summary: "a discarded handcomp found",
			Effects: []Effect{stashItem("a handcomp")},
		},
		41: {
			Summary: "allowed to help with the owner's farm or ranch",
			Effects: []Effect{pick("choose Animals (Veterinary) or Animals (Farming)",
				opt("Animals (Veterinary)", skill("Animals", "Veterinary")),
				opt("Animals (Farming)", skill("Animals", "Farming")))},
		},
		42: {
			Summary: "trained to assist in the operation of a starship",
			Effects: []Effect{pick("choose the job they trained you for",
				opt("Astrogation", skill("Astrogation")),
				opt("Gunner", skill("Gunner", "Any")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Pilot", skill("Pilot", "Any")))},
		},
		43: {
			Summary: "trained to operate several pieces of equipment",
			Effects: []Effect{skill("Electronics", "Any")},
		},
		44: {
			Summary: "maintaining the health of your group",
			Effects: []Effect{skill("Medic", "First Aid")},
		},
		// The table says "Survive (Any)"; the skill list (p. 304) has
		// Survival. See ERRATA E-14.
		45: {Summary: "learning to survive", Effects: []Effect{skill("Survival", "Any")}},
		46: {
			Summary: "a lot of time in space",
			Effects: []Effect{pick("choose Suit (Vacc Suit) or Survival (Freefall)",
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
				opt("Survival (Freefall)", skill("Survival", "Freefall")))},
		},
		51: {
			Summary: "money found",
			Effects: []Effect{credits("1d6x100")},
		},
		52: {
			Summary: "keeping in good physical condition where you can",
			Effects: []Effect{pick("choose what the effort built",
				opt("Athletics", skill("Athletics", "Any")),
				opt("STR", chr("STR", 1)),
				opt("DEX", chr("DEX", 1)),
				opt("END", chr("END", 1)))},
		},
		53: {
			Summary: "a game of chance, wagering food and what the stash holds",
			Effects: []Effect{
				skill("Gambler"),
				checkSkill("Gambler", 8,
					[]Effect{credits("500")},
					[]Effect{unimplemented("lose any money in the stash")}),
			},
		},
		// The table says "Melee (Unarmed Combat)"; the specialty is
		// Unarmed.
		54: {
			Summary: "fights breaking out among your group, often enough to learn from",
			Effects: []Effect{skill("Melee", "Unarmed Combat")},
		},
		55: {Summary: "getting very sneaky", Effects: []Effect{skill("Stealth")}},
		56: {
			Summary: "an escape attempted",
			Effects: []Effect{unimplemented(
				"roll 1d6: on 1 you are caught within moments, at -2 to the next survival " +
					"roll; on 2-3 you are lost, and Survival 8+ makes the escape; on 4-5 the " +
					"guards find you, and Gun Combat or Melee 8+ makes it; on 6 you are away. " +
					"Any escape enlists you in the Vagabond career; any failure is -2 to the " +
					"next survival roll")},
		},
		61: {
			Summary: "a slug pistol found",
			Effects: []Effect{pick("be rid of it, or keep it",
				opt("be rid of it"),
				opt("keep it",
					stashItem("a slug pistol"),
					unimplemented("roll 1d6 for what it is: on 1 it was used in a recent "+
						"murder and puts you in prison for 1d3 terms until the real killer "+
						"is found; on 2-3 it is old and unreliable, at -2 to use, with six "+
						"rounds; on 4-5 it is sound but empty; on 6 it is sound with ten "+
						"rounds")))},
		},
		62: {
			Summary: "the small downtime you have earned, spent on the arts",
			Effects: []Effect{skill("Art", "Any")},
		},
		63: {
			Summary: "trained by your owner in the use of a weapon",
			Effects: []Effect{skill("Gun Combat", "Any")},
		},
		64: {
			Summary: "what a hard life earns you",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		// Neither of this table's last two rows is the usual pair: both are
		// ways out of the career.
		65: {
			Summary: "slavery outlawed on this world, and your freedom with it",
			Effects: []Effect{unimplemented(
				"roll 1d6: on 1-2 enlist automatically in the Vagabond career, keeping the " +
					"stash; on 3-6 attempt any career, losing the stash")},
		},
		66: {
			Summary: "an owner who sets you free, and helps you find another job",
			Effects: []Effect{
				credits("500"),
				unimplemented("enter the career of your choice without an enlistment roll"),
			},
		},
	}

	lifeEventRows(table)

	return table
}
