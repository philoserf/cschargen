package career

// The events tables of the two graduate tracks, pp. 99-100 and 102-103.
//
// Both reach the career Life Events table of p. 120 at result 10, where the
// two undergraduate tracks reach one of their own. So the graduate schools
// have no LifeEvents of their own: the row is an ordinary lifeEvent().

// graduateEvents is the 2d10 table of pp. 99-100.
func graduateEvents() []EventRow {
	return []EventRow{
		{
			Summary: "an academic scandal you remain in the programme despite",
			Effects: []Effect{relationshipAt(Enemy, 1, -140), chr("EDU", -2), chr("CHA", -1)},
		},
		{
			Summary: "a primary project that fails catastrophically",
			Effects: []Effect{chr("EDU", -1), chr("CHA", -1)},
		},
		{
			Summary: "a supervising academic who proves obstructive",
			Effects: []Effect{relationshipAt(Rival, 1, -60)},
		},
		{
			Summary: "uncertain funding, juggled against the research",
			Effects: []Effect{chr("EDU", -1), skill("Trade", "Any")},
		},
		{
			Summary: "an academic load that pushes you past comfort",
			Effects: []Effect{chr("EDU", -1), chr("END", 1)},
		},
		{
			Summary: "a strong rivalry with another student",
			Effects: []Effect{relationshipAt(Rival, 1, -40), skill("Persuade")},
		},
		{
			Summary: "instruction and tutoring assisted with",
			Effects: []Effect{pick("choose what the teaching built",
				opt("CHA", chr("CHA", 1)),
				opt("Instruction", skill("Instruction")))},
		},
		{
			Summary: "coursework that advances effectively",
			Effects: []Effect{pickSkill("Science", "Investigate", "Admin")},
		},
		{Summary: "a life event", Effects: []Effect{lifeEvent()}},
		{
			Summary: "a lasting intellectual partnership",
			Effects: []Effect{relationshipAt(Contact, 1, 40)},
		},
		{
			Summary: "professional gatherings presented at or attended",
			Effects: []Effect{pickSkill("Broker", "Carouse")},
		},
		{
			Summary: "a respected academic who takes an interest in your work",
			Effects: []Effect{relationshipAt(Ally, 1, 110)},
		},
		{
			Summary: "research formally published",
			Effects: []Effect{
				pick("choose what the publication built",
					opt("Art (Writing)", skill("Art", "Writing")),
					opt("Advocate", skill("Advocate", "Any")),
					opt("Admin", skill("Admin")),
					opt("Broker", skill("Broker"))),
				chr("CHA", 1),
			},
		},
		{
			Summary: "work that leads to real world use",
			Effects: []Effect{pickSkill("Admin", "Broker", "Trade")},
		},
		{
			Summary: "a major grant secured",
			Effects: []Effect{skill("Broker"), chr("EDU", 1)},
		},
		{
			Summary: "a genuine research breakthrough",
			Effects: []Effect{pick("choose what the insight built",
				opt("INT", chr("INT", 1)),
				opt("Science", skill("Science", "Any")))},
		},
		{
			Summary: "strong academic connections to leave with",
			Effects: []Effect{pick("choose what the network is worth",
				opt("contacts", relationshipAt(Contact, 2, 50)),
				opt("CHA", chr("CHA", 1)),
				opt("Carouse", skill("Carouse")))},
		},
		{
			Summary: "an institution that marks you as exceptional",
			Effects: []Effect{
				pickSkill("Leadership", "Science"),
				pick("choose what the recognition built",
					opt("EDU", chr("EDU", 1)),
					opt("CHA", chr("CHA", 1))),
			},
		},
		distinguishedGraduate(),
	}
}

// medicalEvents is the 2d10 table of pp. 102-103.
func medicalEvents() []EventRow {
	return []EventRow{
		{
			Summary: "a supervised case that ends badly despite your best efforts",
			Effects: []Effect{chr("CHA", -2)},
		},
		{
			Summary: "long rotations that degrade your own health",
			Effects: []Effect{
				pick("choose what the exhaustion took",
					opt("STR", chr("STR", -1)),
					opt("DEX", chr("DEX", -1))),
				chr("END", -1),
			},
		},
		{
			Summary: "a senior instructor who publicly disputes your judgement",
			Effects: []Effect{chr("EDU", -1), relationshipAt(Rival, 1, -80)},
		},
		{
			Summary: "extensive documentation and support duties",
			Effects: []Effect{skill("Admin"), chr("EDU", -1)},
		},
		{
			Summary: "assignment to a large emergency or outbreak",
			Effects: []Effect{chr("EDU", -1), chr("END", 1), skill("Survival", "Any")},
		},
		{
			Summary: "strong competition among your peers",
			Effects: []Effect{relationshipAt(Rival, 1, -20)},
		},
		{
			Summary: "junior students tutored",
			Effects: []Effect{skill("Instruction")},
		},
		{
			Summary: "extended exposure to equipment and analysis",
			Effects: []Effect{skill("Electronics", "Any")},
		},
		{Summary: "a life event", Effects: []Effect{lifeEvent()}},
		{
			Summary: "a dependable partnership with a fellow student",
			Effects: []Effect{relationshipAt(Contact, 1, 70)},
		},
		{
			Summary: "praise for your performance during clinical rotation",
			Effects: []Effect{chr("CHA", 1), relationshipAt(Contact, 1, 60)},
		},
		{
			Summary: "familiarity with an additional field of practice",
			Effects: []Effect{skill("Medic", "Any")},
		},
		{
			Summary: "formal study or trials assisted in",
			Effects: []Effect{pick("choose what the research built",
				opt("EDU", chr("EDU", 1)),
				opt("Science", skill("Science", "Any")))},
		},
		{
			Summary: "care provided in underserved or remote settings",
			Effects: []Effect{pickSkill("Carouse", "Diplomat", "Persuade", "Streetwise")},
		},
		{
			Summary: "a respected instructor who advocates for your advancement",
			Effects: []Effect{relationshipAt(Ally, 1, 130)},
		},
		{
			Summary: "small teams coordinated under supervision",
			Effects: []Effect{pickSkill("Admin", "Leadership")},
		},
		{
			Summary: "a network of comrades formed",
			Effects: []Effect{relationshipAt(Contact, 2, 60)},
		},
		{
			Summary: "publication credit from a study or case review",
			Effects: []Effect{
				pick("choose what the credit built",
					opt("Art (Writing)", skill("Art", "Writing")),
					opt("Science", skill("Science", "Any"))),
				chr("INT", 1),
			},
		},
		distinguishedGraduate(),
	}
}

// distinguishedGraduate is result 20 on three of the four events tables,
// printed identically each time.
func distinguishedGraduate() EventRow {
	return EventRow{
		Summary: "among the finest of your class",
		Effects: []Effect{
			relationshipAt(Ally, 1, 110),
			relationshipAt(Contact, 1, 60),
			{Kind: EffectHonors, Detail: "achieve honours, if the throw did not"},
			chr("CHA", 1),
			pickSkill("Discipline", "Diplomat", "Admin", "Leadership", "Recon"),
		},
	}
}
