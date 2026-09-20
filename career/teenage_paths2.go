package career

// teenagePathThree is pp. 80-81, for a teenager with STR, DEX or END at 9
// or higher.
func teenagePathThree() []EventRow {
	return []EventRow{
		{
			Summary: "a serious violation of local law that stays on the record",
			Effects: []Effect{
				throwModifier("enlistment in any non-criminal, non-military career", -2),
				pickSkill("Streetwise", "Deception"),
			},
		},
		familyObligation(),
		{
			Summary: "time spent in hazardous conditions, out of necessity",
			Effects: []Effect{
				chr("END", 1),
				unimplemented("a level in a Survival specialty from the homeworld's " +
					"background skills"),
			},
		},
		authorityConflict(),
		{
			Summary: "taken on informally by a worker, coach or supervisor",
			Effects: []Effect{
				relationshipAt(Ally, 1, 120),
				unimplemented("choose a career, roll once on its Service Skills table, " +
					"and take +4 to enlist in it at Step 9"),
			},
		},
		{
			Summary: "a demanding expedition, and something proved to yourself",
			Effects: []Effect{
				pick("choose what the effort built",
					opt("STR", chr("STR", 1)),
					opt("DEX", chr("DEX", 1)),
					opt("END", chr("END", 1))),
				skill("Athletics", "Any"),
			},
		},
		{
			Summary: "a careless comment that damages a relationship",
			Effects: []Effect{rating(TargetOne, "", -50, ""), skill("Persuade")},
		},
		educationalDirection(),
		teenageLifeEvent(),
		broadeningHorizons(),
		{
			Summary: "chosen for a team or youth program on ability",
			Effects: []Effect{
				skill("Athletics", "Any"),
				throwModifier("enlistment in the Sports career", 2),
			},
		},
		{
			Summary: "someone others rely on, through presence rather than authority",
			Effects: []Effect{skill("Leadership"), chr("CHA", 1)},
		},
		{
			Summary: "competence that is not purely physical",
			Effects: []Effect{pick("choose what it brought",
				opt("contacts", relationshipRolledAt(Contact, "1d3", 50)),
				opt("CHA", chr("CHA", 2)))},
		},
		earnedRecommendation(),
		{
			Summary: "time alongside adults in real labor and emergency response",
			Effects: []Effect{throwModifier(
				"enlistment in the National Navy, a System Defense Force, or any "+
					"military or law enforcement career", 4)},
		},
		socialNetwork(),
		opportunityKnocks(),
		{
			Summary: "local sports, with your school or an organized league",
			Effects: []Effect{
				skill("Athletics", "Any"),
				unimplemented("roll 1d6 for how it went: on 1 nothing, on 2-3 gain 1d3 " +
					"Contacts at 50, on 4-5 gain 1d6 Contacts at 50, on 6 gain another " +
					"level in Athletics, +2 CHA and +2 to enter the Sports career"),
			},
		},
		{
			Summary: "a local university that offers a sports scholarship",
			Effects: []Effect{
				unimplemented("skip to Step 8 with automatic admission anywhere you " +
					"qualify, and to Undergraduate University whether you qualify or not"),
				skill("Athletics", "Any"),
				skill("Tactics", "Sport"),
				unimplemented("after Step 8, you may enlist in the Sports career automatically"),
			},
		},
	}
}

// teenagePathFour is pp. 82-83, for a teenager with INT or EDU at 9 or
// higher.
func teenagePathFour() []EventRow {
	return []EventRow{
		{
			Summary: "a serious disease that does irrevocable damage",
			Effects: []Effect{pick("choose what the disease took",
				opt("STR", chrRolled("STR", "1d3", false)),
				opt("END", chrRolled("END", "1d3", false)))},
		},
		{
			Summary: "few peers who share your interests or pace of thought",
			Effects: []Effect{chr("CHA", -1), pickSkill("Art", "Science", "Recon")},
		},
		{
			Summary: "failure at something important, despite the intelligence",
			Effects: []Effect{chr("EDU", -1), chr("CHA", -1)},
		},
		familyObligation(),
		authorityConflict(),
		{
			Summary: "a retreat into reading and quiet learning",
			Effects: []Effect{chr("EDU", 1), chr("CHA", -1)},
		},
		educationalDirection(),
		{
			Summary: "a group of elite computer hackers",
			Effects: []Effect{
				skill("Electronics", "Computers"),
				unimplemented("roll 1d6: on 1 a term in the Prisoner career; on 2-3 the " +
					"group as a Contact at 50; on 4-5 another level in Electronics " +
					"(Computers) and the group as an Ally at 110; on 6 that and " +
					"1d6 x 100,000 credits, but a homeworld you can never return to"),
			},
		},
		teenageLifeEvent(),
		{
			Summary: "people who begin to seek your help with practical questions",
			Effects: []Effect{pickSkill("Investigate", "Recon")},
		},
		{
			Summary: "subjects delved into more deeply than most teenagers attempt",
			Effects: []Effect{skill("Science", "Any")},
		},
		{
			Summary: "gaming, at a table, on the worldnet, or in a holonovel",
			Effects: []Effect{pick("choose what the games taught you",
				opt("Art", skill("Art", "Acting", "Holography", "Writing")),
				opt("Carouse", skill("Carouse")),
				opt("Electronics (Computers)", skill("Electronics", "Computers")),
				opt("Tactics", skill("Tactics", "Military", "Naval")))},
		},
		{
			Summary: "the realisation that intellect alone is not enough",
			Effects: []Effect{chr("CHA", 1), pickSkill("Diplomat", "Etiquette", "Persuade")},
		},
		earnedRecommendation(),
		{
			Summary: "an interest in flying noticed, and lessons paid for",
			Effects: []Effect{pick("choose what they taught you to fly",
				opt("Flyer", skill("Flyer", "Grav", "Rotor", "Wing")),
				opt("Pilot (Small Craft)", skill("Pilot", "Small Craft")))},
		},
		socialNetwork(),
		opportunityKnocks(),
		{
			Summary: "a band joined",
			Effects: []Effect{
				skill("Art", "Instrument"),
				unimplemented("roll 1d6 for how it went: on 1 a Rival at -60; on 2-3 the " +
					"other three as Contacts at 50; on 4-5 +2 CHA and the other three as " +
					"Allies at 125; on 6 gain 1d3+1 CHA and +2 to enter any " +
					"entertainment-based career"),
			},
		},
		{
			Summary: "a local university impressed by your mental characteristics",
			Effects: []Effect{
				unimplemented("skip to Step 8 with automatic admission anywhere you " +
					"qualify, and to Undergraduate University whether you qualify or not"),
				skill("Admin"),
				skill("Instruction"),
				unimplemented("after Step 8, you may enlist in the Instructor career " +
					"automatically"),
			},
		},
	}
}
