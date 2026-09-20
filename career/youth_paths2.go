package career

// youthPathTwo is pp. 70-71, open to a character with STR 8+ or END 8+.
func youthPathTwo() []EventRow {
	return []EventRow{
		orphaned(),
		{
			Summary: "adults who noted your strength and took advantage of it",
			Effects: []Effect{
				pick("choose what the work built",
					opt("STR", chr("STR", 1)),
					opt("END", chr("END", 1))),
				pick("choose what it cost",
					opt("EDU", chr("EDU", -1)),
					opt("CHA", chr("CHA", -1))),
				rating(TargetOne, Ally, -20, ""),
				skill("Trade", "Laborer"),
			},
		},
		{
			Summary: "violence, which tends to find the strong",
			Effects: []Effect{pick("choose what trouble taught you",
				opt("Melee (Unarmed)", skill("Melee", "Unarmed Combat")),
				opt("Streetwise", skill("Streetwise")))},
		},
		{
			Summary: "a serious accident many thought had killed you",
			Effects: []Effect{chr("DEX", -1), chr("END", 1)},
		},
		{
			// "Gain +1 in the characteristic that you used to get into this
			// path but lose 1 from another physical characteristic." The
			// engine does not record which characteristic opened the path,
			// so both halves are a choice. ERRATA E-26.
			Summary: "the genetics without the training, and growing fast but not cleanly",
			Effects: []Effect{
				pick("raise the characteristic that opened this path",
					opt("STR", chr("STR", 1), pick("and lose another physical one",
						opt("DEX", chr("DEX", -1)), opt("END", chr("END", -1)))),
					opt("END", chr("END", 1), pick("and lose another physical one",
						opt("STR", chr("STR", -1)), opt("DEX", chr("DEX", -1))))),
			},
		},
		{
			Summary: "too big for your age, and watched carefully for it",
			Effects: []Effect{chr("CHA", -1), chr("END", 1)},
		},
		{
			Summary: "an older child's attempt at bullying that did not go as expected",
			Effects: []Effect{pick("choose what you took from it",
				opt("Melee (Unarmed)", skill("Melee", "Unarmed Combat")),
				opt("Recon", skill("Recon")))},
		},
		{
			Summary: "tough times your physical capabilities carried you through",
			Effects: []Effect{
				chr("EDU", -1),
				homeworldSpecialty("Survival"),
				chr("END", 1),
			},
		},
		youthLifeEvent(),
		{
			Summary: "athletic talent spotted early and encouraged",
			Effects: []Effect{skill("Athletics", "Any"), skill("Carouse")},
		},
		{
			Summary: "trusted by adults to watch over the younger children",
			Effects: []Effect{skillZero("Leadership"), chr("CHA", 1)},
		},
		{
			Summary: "steadiness rather than muscle, and the confidence it brought",
			Effects: []Effect{skill("Persuade"), chr("CHA", 1)},
		},
		{
			Summary: "an adult who saw the potential as well as the strength",
			Effects: []Effect{
				mentorAt125(),
				pick("choose what they trained",
					opt("Melee (Unarmed)", skill("Melee", "Unarmed Combat")),
					opt("Athletics", skill("Athletics", "Any")),
					opt("Discipline", skill("Discipline")),
					opt("Tactics (Sport)", skill("Tactics", "Sport"))),
			},
		},
		{
			Summary: "mental toughness as well as physical",
			Effects: []Effect{chr("END", 1), chr("INT", 1)},
		},
		{
			Summary: "physical strength that got the adults through some tough times",
			Effects: []Effect{
				chr("STR", 1),
				pick("choose what the help taught you",
					opt("Animals", skill("Animals", "Any")),
					opt("Athletics", skill("Athletics", "Any")),
					opt("Mechanic", skill("Mechanic")),
					opt("Trade", skill("Trade", "Carpenter", "Laborer", "Lumberjack",
						"Prospector", "Space Construction"))),
			},
		},
		{
			Summary: "help given during a crisis in your community",
			Effects: []Effect{relationshipAt(Ally, 1, 110), chr("END", 1), chr("CHA", 1)},
		},
		{
			Summary: "seen by the adults as dependable",
			Effects: []Effect{relationshipAt(Contact, 2, 40), chr("CHA", 1)},
		},
		{
			Summary: "confidence that inspired and helped others",
			Effects: []Effect{
				pick("choose what the confidence built",
					opt("STR", chr("STR", 1)),
					opt("END", chr("END", 1)),
					opt("CHA", chr("CHA", 1))),
				pickSkill("Athletics", "Carouse", "Leadership", "Persuade"),
			},
		},
		{
			Summary: "physical gifts nurtured rather than exploited",
			Effects: []Effect{
				chr("STR", 1), chr("END", 1),
				skill("Athletics", "Any"),
				modifierFor(enlistmentThrow, 4, 0, "+4 to enlistment in the Sports career",
					enlistmentNarrowing{OnCareers: []string{"Sports"}, Standing: true}),
			},
		},
	}
}

// youthPathThree is pp. 71-72, open to a character with DEX 8+.
func youthPathThree() []EventRow {
	return []EventRow{
		orphaned(),
		{
			Summary: "too much confidence, and one slip too many",
			Effects: []Effect{chr("DEX", -1), chr("END", -1), skillZero("Medic")},
		},
		{
			Summary: "the discovery that most doors are less secure than adults think",
			Effects: []Effect{
				pick("choose what you learned before someone noticed",
					opt("Deception (Intrusion)", skill("Deception", "Intrusion")),
					opt("Stealth", skill("Stealth")),
					opt("Streetwise", skill("Streetwise"))),
				chr("EDU", -1),
			},
		},
		{
			Summary: "a careless moment with a delicate item, and the effort of mending it",
			Effects: []Effect{skill("Mechanic"), chr("CHA", -1)},
		},
		{
			Summary: "dexterity that was creative rather than mechanical",
			Effects: []Effect{
				skill("Art", "Any"),
				pick("choose what the making built",
					opt("EDU", chr("EDU", 1)),
					opt("CHA", chr("CHA", 1))),
			},
		},
		{
			Summary: "always underfoot, always listening, always forgotten",
			Effects: []Effect{
				pick("choose what the listening taught you",
					opt("Recon", skill("Recon")),
					opt("Stealth", skill("Stealth"))),
				chr("INT", 1),
			},
		},
		{
			Summary: "a toy that was good, and that your hands made better",
			Effects: []Effect{pick("choose what came of it",
				opt("Mechanic", skill("Mechanic")),
				opt("EDU", chr("EDU", 1)))},
		},
		{
			Summary: "the one who fixed things for the other children",
			Effects: []Effect{pick("choose what the fixing taught you",
				opt("Mechanic", skill("Mechanic")),
				opt("Medic (First Aid)", skill("Medic", "First Aid")),
				opt("Persuade", skill("Persuade")))},
		},
		youthLifeEvent(),
		{
			Summary: "quick, and so the runner",
			Effects: []Effect{pickSkill("Athletics", "Carouse", "Streetwise")},
		},
		{
			Summary: "not a bad child, just quick and mischievous",
			Effects: []Effect{
				pickSkill("Carouse", "Deception", "Streetwise"),
				chr("CHA", 1),
			},
		},
		{
			Summary: "a lot of time around the stalls and the loading ramps",
			Effects: []Effect{pick("choose what the trade taught you",
				opt("Broker", skillZero("Broker")),
				opt("Streetwise", skillZero("Streetwise")),
				opt("INT", chr("INT", 1)))},
		},
		{
			Summary: "not the smartest or the fastest, but consistent",
			Effects: []Effect{pick("choose what reliability built",
				opt("EDU", chr("EDU", 1)),
				opt("Jack of All Trades", skill("Jack of All Trades")))},
		},
		{
			Summary: "a local gambler who took you under his wing",
			Effects: []Effect{pick("choose what he taught you",
				opt("DEX", chr("DEX", 1)),
				opt("Gambler", skill("Gambler")),
				opt("Stealth", skill("Stealth")),
				opt("Recon", skill("Recon")))},
		},
		{
			Summary: "good with your hands, and finally a responsible adult noticed",
			Effects: []Effect{
				mentorAt125(),
				pick("choose what they taught",
					opt("Athletics", skill("Athletics", "Any")),
					opt("Art", skill("Art", "Any")),
					opt("Chef", skill("Chef")),
					opt("Mechanic", skill("Mechanic")),
					opt("Melee (Unarmed)", skill("Melee", "Unarmed Combat")),
					opt("Stealth", skill("Stealth"))),
			},
		},
		{
			Summary: "games that rewarded timing and coordination, including vehicle sims",
			Effects: []Effect{pick("choose the sim you played",
				opt("Drive", skillZero("Drive", "Any")),
				opt("Flyer", skillZero("Flyer", "Any")),
				opt("Pilot", skillZero("Pilot", "Any")))},
		},
		{
			Summary: "steadiness when something went wrong, noticed by others",
			Effects: []Effect{pick("choose what the steadiness built",
				opt("DEX", chr("DEX", 1)),
				opt("Leadership", skillZero("Leadership")))},
		},
		{
			Summary: "comfortable in your own body, and movement that felt natural",
			Effects: []Effect{chr("DEX", 1), chr("END", 1)},
		},
		{
			Summary: "talents noticed early and supported appropriately",
			Effects: []Effect{
				chr("DEX", 1),
				pick("choose what the support built",
					opt("EDU", chr("EDU", 1)),
					opt("CHA", chr("CHA", 1))),
			},
		},
	}
}
