package career

// youthPathFour is pp. 72-73, open to a character with INT 8+.
func youthPathFour() []EventRow {
	return []EventRow{
		orphaned(),
		{
			Summary: "curiosity treated as defiance, and punished for it",
			Effects: []Effect{
				pick("choose what the thinking built anyway",
					opt("INT", chr("INT", 1)),
					opt("EDU", chr("EDU", 1))),
				chr("CHA", -1),
			},
		},
		{
			Summary: "working out how grading and reward actually worked",
			Effects: []Effect{pick("choose what the cynicism taught you",
				opt("Admin", skill("Admin")),
				opt("Deception (Lie)", skill("Deception", "Lie")))},
		},
		{
			Summary: "nothing around you that engaged your mind",
			Effects: []Effect{chr("INT", 1), chr("EDU", -1)},
		},
		{
			Summary: `noticed as "different", and the wrong conclusions drawn`,
			Effects: []Effect{chr("INT", 1), chr("CHA", -1)},
		},
		{
			Summary: "facts without frameworks, and the confusion that lingered",
			Effects: []Effect{pick("choose what you salvaged",
				opt("EDU", chr("EDU", 1)),
				opt("Admin", skillZero("Admin")))},
		},
		{
			Summary: "listening more than speaking, and staying informed",
			Effects: []Effect{pick("choose what the listening built",
				opt("INT", chr("INT", 1)),
				opt("Recon", skillZero("Recon")))},
		},
		{
			Summary: `told to "work it out yourself", and doing so`,
			Effects: []Effect{pick("choose what independence built",
				opt("INT", chr("INT", 1)),
				opt("Jack of All Trades", skill("Jack of All Trades")))},
		},
		youthLifeEvent(),
		{
			Summary: "questions asked, answers learned, life in balance",
			Effects: []Effect{pickSkill("Admin", "Science"), chr("EDU", 1)},
		},
		{
			Summary: "teaching the other children, and the clarity it brought",
			Effects: []Effect{
				pick("choose what the teaching built",
					opt("Instruction", skillZero("Instruction")),
					opt("Leadership", skillZero("Leadership"))),
				chr("CHA", 1),
			},
		},
		{
			Summary: `described as "bright", with expectations that rose gently`,
			Effects: []Effect{chr("INT", 1), chr("EDU", 1)},
		},
		{
			Summary: "access to records and explanations earlier than others",
			Effects: []Effect{pick("choose what the context opened",
				opt("Admin", skill("Admin")),
				opt("Science", skill("Science", "Any")),
				opt("Electronics (Computers)", skill("Electronics", "Computers")))},
		},
		{
			Summary: "someone who challenged your thinking without dismissing it",
			Effects: []Effect{
				chr("EDU", 1),
				mentorAt125(),
				pick("choose what they taught",
					opt("Admin", skill("Admin")),
					opt("Advocate", skill("Advocate", "Legal", "Politics")),
					opt("Art", skill("Art", "Any")),
					opt("Broker", skill("Broker")),
					opt("Diplomat", skill("Diplomat")),
					opt("Electronics", skill("Electronics", "Any")),
					opt("Engineer", skill("Engineer", "Any")),
					opt("Etiquette", skill("Etiquette")),
					opt("Investigate", skill("Investigate")),
					opt("Science", skill("Science", "Any"))),
			},
		},
		{
			Summary: "a problem worked through that mattered to someone",
			Effects: []Effect{
				pick("choose what the solution built",
					opt("INT", chr("INT", 1)),
					opt("EDU", chr("EDU", 1))),
				chr("CHA", 1),
			},
		},
		{
			Summary: "an adult who trusted your reasoning in a real decision",
			Effects: []Effect{
				chr("CHA", 1),
				relationshipAt(Contact, 1, 50),
				pick("choose what accountability taught you",
					opt("Leadership", skillZero("Leadership")),
					opt("Science", skill("Science", "Any"))),
			},
		},
		{
			Summary: "freedom to explore interests without pressure",
			Effects: []Effect{
				pick("choose what curiosity built",
					opt("INT", chr("INT", 1)),
					opt("EDU", chr("EDU", 1)),
					opt("CHA", chr("CHA", 1))),
				pickSkill("Art", "Science"),
			},
		},
		{
			Summary: "an intelligence others quietly deferred to",
			Effects: []Effect{chr("INT", 1), chr("CHA", 1)},
		},
		{
			// The only result in the book that grants a skill the
			// character chooses without a list.
			Summary: "intelligence that created opportunity, and a step ahead",
			Effects: []Effect{
				chr("INT", 1), chr("EDU", 1),
				anySkillNotHeld(),
			},
		},
	}
}

// youthPathFive is pp. 73-74, open to a character with CHA 8+.
func youthPathFive() []EventRow {
	return []EventRow{
		orphaned(),
		{
			Summary: "shamed in front of others, and the fear of being judged",
			Effects: []Effect{skill("Deception", "Lie"), chr("CHA", -1)},
		},
		{
			Summary: "speaking and being ignored, until you stopped",
			Effects: []Effect{pickSkill("Recon", "Investigate")},
		},
		{
			Summary: "a name that came up first whenever something went wrong",
			Effects: []Effect{chr("CHA", -1), skillZero("Streetwise")},
		},
		{
			Summary: "conflict you learned to calm, a word here and a joke there",
			Effects: []Effect{pick("choose what the buffering taught you",
				opt("Persuade", skillZero("Persuade")),
				opt("Carouse", skillZero("Carouse")))},
		},
		{
			Summary: "changing tone and posture without thinking, because survival required it",
			Effects: []Effect{pick("choose what the flexibility taught you",
				opt("Art (Acting)", skill("Art", "Acting")),
				opt("Deception (Lie)", skill("Deception", "Lie")))},
		},
		{
			Summary: "not popular, but welcome",
			Effects: []Effect{pick("choose what the comfort built",
				opt("CHA", chr("CHA", 1)),
				opt("Carouse", skill("Carouse")))},
		},
		{
			Summary: "people who talked to you because you listened",
			Effects: []Effect{pick("choose what the trust built",
				opt("CHA", chr("CHA", 1)),
				opt("Persuade", skill("Persuade")))},
		},
		youthLifeEvent(),
		{
			Summary: "a dependable group that endured its arguments",
			Effects: []Effect{pick("choose what the group taught you",
				opt("CHA", chr("CHA", 1)),
				opt("Leadership", skillZero("Leadership")),
				opt("contacts", relationshipRolledAt(Contact, "1d3", 50)))},
		},
		{
			Summary: "chosen to speak on behalf of a group",
			Effects: []Effect{skillZero("Leadership"), chr("CHA", 1)},
		},
		{
			Summary: "reasonable and even-handed, and trusted to settle disputes",
			Effects: []Effect{skill("Persuade"), chr("CHA", 1)},
		},
		{
			Summary: "entering conversations without hesitation",
			Effects: []Effect{skill("Carouse"), chr("CHA", 1)},
		},
		{
			Summary: "teachers and officials who remembered your name and face",
			Effects: []Effect{relationshipAt(Contact, 2, 50), chr("CHA", 1)},
		},
		{
			Summary: "informal authority when the adults stepped away",
			Effects: []Effect{skillZero("Leadership"), chr("CHA", 1)},
		},
		{
			Summary: "explaining adults to your fellow children",
			Effects: []Effect{skill("Persuade"), chr("CHA", 1)},
		},
		{
			Summary: "social strengths noticed and nurtured by peers and adults alike",
			Effects: []Effect{chr("EDU", 1), chr("CHA", 1)},
		},
		{
			Summary: "sought out privately for advice that was not always taken",
			Effects: []Effect{pickSkill("Leadership", "Carouse"), chr("CHA", 1)},
		},
		{
			Summary: "growing up knowing you mattered to the people around you",
			Effects: []Effect{skill("Carouse"), skill("Persuade"), skillZero("Leadership")},
		},
	}
}
