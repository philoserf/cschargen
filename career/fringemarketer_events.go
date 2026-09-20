package career

// fringeMarketerEvents is the d66 table of pp. 200-201.
func fringeMarketerEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "keeping the books for the business",
			Effects: []Effect{skill("Admin")},
		},
		13: {
			Summary: "a valuable item, and time to sell it",
			Effects: []Effect{checkSkill("Broker", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch)},
				[]Effect{benefitRolls(-1, 0, ScopeBatch)})},
		},
		14: {
			Summary: "a group of animals taken on to sell, and cared for in the meantime",
			Effects: []Effect{skill("Animals", "Any")},
		},
		15: {
			Summary: "research, which is often your friend",
			Effects: []Effect{pick("choose Electronics (Computers) or Investigate",
				opt("Electronics (Computers)", skill("Electronics", "Computers")),
				opt("Investigate", skill("Investigate")))},
		},
		16: {
			Summary: "a statue of a bird acquired",
			Effects: []Effect{pick("look it over however you look things over",
				opt("Recon", checkSkill("Recon", 8, fringeDiamonds(), fringeSoldOn())),
				opt("Investigate", checkSkill("Investigate", 8, fringeDiamonds(), fringeSoldOn())))},
		},
		21: {
			Summary: "familiarity with the local laws, which this business requires",
			Effects: []Effect{skill("Advocate", "Legal")},
		},
		22: {
			Summary: "training",
			Effects: []Effect{pick("choose what the training built",
				opt("Athletics", skill("Athletics", "Any")),
				opt("STR", chr("STR", 1)),
				opt("DEX", chr("DEX", 1)),
				opt("END", chr("END", 1)))},
		},
		23: {
			Summary: "it is all about the customers",
			Effects: []Effect{checkSkill("Carouse", 8,
				[]Effect{relationship(Contact, 1, "1d3")},
				[]Effect{unimplemented("lose one Contact")})},
		},
		24: {
			Summary: "a gambling circle in your peer group",
			Effects: []Effect{gamblingCircle("Persuade")},
		},
		25: {
			Summary: "all sorts of people from all sorts of places",
			Effects: []Effect{skill("Language", "Any")},
		},
		26: {Summary: "first aid learned", Effects: []Effect{skill("Medic", "First Aid")}},
		41: {
			Summary: "a hobby taken up",
			Effects: []Effect{pickSkill("Art", "Science")},
		},
		42: {
			Summary: "opportunity that has to be made rather than waited for",
			Effects: []Effect{skill("Deception", "Forgery")},
		},
		43: {
			Summary: "you are the transportation",
			Effects: []Effect{pickSkill("Drive", "Flyer")},
		},
		44: {
			Summary: "a dangerous business",
			Effects: []Effect{skill("Gun Combat", "Any")},
		},
		45: {
			Summary: "who is going to fix this? you are",
			Effects: []Effect{skill("Mechanic")},
		},
		46: {
			Summary: "an expert called in, because nobody knows every item",
			Effects: []Effect{relationship(Contact, 0, "1d3+1")},
		},
		51: {
			Summary: "world to world this term, and a crew that needed a bit of help",
			Effects: []Effect{pick("choose the job they taught you",
				opt("Astrogation", skill("Astrogation")),
				opt("Electronics (Sensors)", skill("Electronics", "Sensors")),
				opt("Engineer", skill("Engineer", "Any")),
				opt("Gunner (Turrets)", skill("Gunner", "Turrets")),
				opt("Pilot", skill("Pilot", "Any")))},
		},
		52: {
			Summary: "the fringe is a seedy place",
			Effects: []Effect{pickSkill("Deception", "Streetwise")},
		},
		53: {
			Summary: "knowing whether someone is lying to you",
			Effects: []Effect{skill("Interrogation", "Any")},
		},
		54: {
			Summary: "fighting your way through some situations",
			Effects: []Effect{pickSkill("Gun Combat", "Melee")},
		},
		55: {
			Summary: "competing for the money of your customers",
			Effects: []Effect{relationship(Rival, 1, "")},
		},
		56: {
			Summary: "blurred lines between the legal and illegal fringe sellers",
			Effects: []Effect{rollOtherAssignment()},
		},
		61: {
			Summary: "selling to the wealthy, which takes a certain flair",
			Effects: []Effect{pickSkill("Diplomat", "Etiquette")},
		},
		62: {
			Summary: "self-sufficiency and flexibility, which the business requires",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		63: {
			Summary: "work done on your education",
			Effects: []Effect{pick("choose EDU or the training",
				opt("EDU", chr("EDU", 1)),
				opt("Advanced Education", rollTable(AdvancedEducation)))},
		},
		64: {
			Summary: "a bit more time in space than usual",
			Effects: []Effect{pick("choose Suit (Vacc Suit) or Survival (Freefall)",
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
				opt("Survival (Freefall)", skill("Survival", "Freefall")))},
		},
		65: {
			Summary: "extremely talented in your field",
			Effects: []Effect{talentedInYourField()},
		},
		66: {
			Summary: "excellent work",
			Effects: []Effect{advance(), benefitRolls(0, 1, ScopeCareer)},
		},
	}

	lifeEventRows(table)

	return table
}

// fringeDiamonds and fringeSoldOn are event 16's two ends: the hole in the
// base of the statue, or a collector who buys it as a statue.
func fringeDiamonds() []Effect {
	return []Effect{benefitRolls(3, 0, ScopeBatch)}
}

func fringeSoldOn() []Effect {
	return []Effect{benefitRolls(1, 0, ScopeBatch)}
}
