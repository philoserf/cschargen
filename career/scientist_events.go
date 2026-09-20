package career

// scientistEvents is the d66 table of pp. 271-272.
func scientistEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {Summary: "detailed notes kept on your work", Effects: []Effect{skill("Admin")}},
		13: {
			Summary: "a popular person",
			Effects: []Effect{pick("what the popularity brought",
				opt("+1 CHA", chr("CHA", 1)),
				opt("Carouse", skill("Carouse")),
				opt("an Ally", relationship(Ally, 1, "")),
				opt("contacts", relationship(Contact, 0, "1d3")))},
		},
		14: {
			Summary: "always working with electronic devices",
			Effects: []Effect{skill("Electronics", "Any")},
		},
		15: {
			Summary: "asked to speak to a group of students about your research",
			Effects: []Effect{checkSkill("Instruction", 8,
				[]Effect{pick("who you reached",
					opt("an Ally", relationship(Ally, 1, "")),
					opt("contacts", relationship(Contact, 0, "1d3")))},
				[]Effect{throwModifier(survivalThrow, -2)})},
		},
		16: {
			Summary: "a teaching position offered at a local university",
			Effects: []Effect{pick("take the post or stay",
				opt("stay in research"),
				opt("take the post", transfer("Instructor", "Professor", 0)))},
		},
		21: {
			Summary: "a voice for your research and for science in general",
			Effects: []Effect{skill("Advocate", "Oratory", "Politics")},
		},
		22: {Summary: "making your own meals", Effects: []Effect{skill("Chef")}},
		23: {
			Summary: "the art of getting other researchers to tell you what they know",
			Effects: []Effect{skill("Interrogation", "Questioning")},
		},
		24: {Summary: "if it is to be repaired, it is often you who repairs it", Effects: []Effect{skill("Mechanic")}},
		25: {
			Summary: "a discussion panel, broadcast to the public",
			Effects: []Effect{checkSkill("Persuade", 8,
				[]Effect{rollSub("how the panel lands (p. 271)",
					on(1, "you persuade nobody but the host", relationship(Ally, 1, "")),
					onRange(2, 3, "the audience is persuaded",
						relationship(Contact, 0, "1d3")),
					onRange(4, 5, "the panel and the audience both",
						relationship(Ally, 1, ""),
						relationship(Contact, 0, "1d3")),
					on(6, "outstanding, and the host wants you back",
						relationship(Ally, 1, ""),
						chr("CHA", 2),
						relationship(Contact, 0, "1d3"),
						pick("stay a Scientist, or take the offer",
							opt("stay"),
							opt("take it", transfer("Celebrity", "Star", 0)))))},
				[]Effect{throwModifier(advancementThrow, -2)})},
		},
		26: {
			// ERRATA E-13: this row is printed with a number and nothing
			// beside it (p. 271).
			Summary: "the blank row",
			Effects: []Effect{unimplemented(
				"this d66 result is printed with no effect beside it (p. 271); nothing happens")},
		},
		41: {
			Summary: "relaxing with your hobbies",
			Effects: []Effect{pick("what you pursued",
				opt("Animals", skill("Animals", "Riding", "Farming")),
				opt("Athletics", skill("Athletics", "Any")),
				opt("Art", skill("Art", "Any")))},
		},
		42: {
			Summary: "some shady people",
			Effects: []Effect{pick("what came of it",
				opt("contacts in the local underworld", relationship(Contact, 0, "1d3")),
				opt("Deception", skill("Deception", "Any")),
				opt("Streetwise", skill("Streetwise")))},
		},
		43: {
			Summary: "a term spent researching the background of your field",
			Effects: []Effect{checkSkill("Investigate", 8,
				[]Effect{throwModifier(advancementThrow, 2)},
				[]Effect{throwModifier(survivalThrow, -2)})},
		},
		44: {Summary: "first aid learned", Effects: []Effect{skill("Medic", "First Aid")}},
		45: {
			Summary: "the best way to survive your environment",
			Effects: []Effect{pick("which environment that is",
				opt("a lab, for a Researcher", skill("Etiquette")),
				opt("the field, for a Field Researcher", skill("Survival", "Any")))},
		},
		46: {
			Summary: "on the edge of a breakthrough",
			Effects: []Effect{checkSkill("Science", 8,
				[]Effect{rollSub("how large the breakthrough is (p. 272)",
					on(1, "minor, and the work itself is the reward",
						raiseHeldSkillIn("Science")),
					onRange(2, 3, "it leads you into another field",
						skill("Science", "Any")),
					onRange(4, 5, "a success",
						throwModifier(advancementThrow, 2),
						raiseHeldSkillIn("Science")),
					on(6, "a major breakthrough in your field",
						raiseHeldSkillIn("Science"),
						skill("Science", "Any"),
						benefitRolls(2, 0, ScopeBatch)))},
				[]Effect{throwModifier(survivalThrow, -2)})},
		},
		51: {
			Summary: "a term spent aboard a starship, helping with the shipboard duties",
			Effects: []Effect{pick("what you learned",
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Gunner (Turrets)", skill("Gunner", "Turrets")),
				opt("Mechanic", skill("Mechanic")),
				opt("Pilot", skill("Pilot", "Any")))},
		},
		52: {
			Summary: "working with superiors, governments and patrons",
			Effects: []Effect{pickSkill("Diplomat", "Etiquette")},
		},
		53: {Summary: "a gambling group", Effects: []Effect{gamblingCircle("Persuade")}},
		54: {
			Summary: "communicating across cultures and languages",
			Effects: []Effect{skill("Language", "Any")},
		},
		55: {
			Summary: "taken out of your element, into another specialty",
			Effects: []Effect{rollOtherAssignment()},
		},
		56: {
			Summary: "money invested in the corporation that sponsors your research",
			Effects: []Effect{
				skill("Broker"),
				checkSkill("Broker", 8,
					[]Effect{stashCountRolled("1d6", companyShare, companyShareValue)},
					[]Effect{benefitRolls(-1, 0, ScopeBatch)}),
			},
		},
		61: {
			Summary: "getting your research out to the masses",
			Effects: []Effect{pick("how you got it out",
				opt("Art", skill("Art", "Writing", "Holography")),
				opt("Advocate (Oratory)", skill("Advocate", "Oratory")),
				opt("Persuade", skill("Persuade")),
				opt("Instruction", skill("Instruction")))},
		},
		62: {
			Summary: "often your own transportation",
			Effects: []Effect{pick("what you drove",
				opt("Drive", skill("Drive", "Any")),
				opt("Flyer", skill("Flyer", "Any")),
				opt("Pilot (Small Craft)", skill("Pilot", "Small Craft")),
				opt("Seafarer", skill("Seafarer", "Any")))},
		},
		63: {
			Summary: "learning to defend yourself",
			Effects: []Effect{pickSkill("Gun Combat", "Melee")},
		},
		64: {
			Summary: "relying on yourself and your own abilities",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		65: {Summary: "extremely talented in your field", Effects: []Effect{talentedInYourField()}},
		66: {
			Summary: "excellent work",
			Effects: []Effect{advance(), benefitRolls(0, 1, ScopeCareer)},
		},
	}

	lifeEventRows(table)

	return table
}
