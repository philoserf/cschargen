package career

// colonistEvents is the d66 table of pp. 175-176.
func colonistEvents() EventTable {
	table := EventTable{
		11: {
			Summary: "disaster occurs",
			Effects: []Effect{mishapNoEject()},
		},
		12: {
			Summary: "a disagreement between colonists turns into armed conflict",
			Effects: []Effect{pick("settle it by talking or by shooting",
				opt("Diplomat", checkSkill("Diplomat", 8,
					[]Effect{skill("Leadership")},
					[]Effect{injury(1), relationship(Enemy, 1, "")})),
				opt("Gun Combat", checkSkill("Gun Combat", 8,
					[]Effect{skill("Recon")},
					[]Effect{injury(1), relationship(Enemy, 1, "")})),
			)},
		},
		13: {
			Summary: "a natural disaster strikes the colony",
			Effects: []Effect{pick("save yourself or try to save others",
				opt("save yourself", checkSkill("Survival", 10,
					[]Effect{pick("what the escape taught you",
						opt("Recon", skill("Recon")),
						opt("Survival", skill("Survival", "Any")))},
					[]Effect{injury(1)})),
				opt("save others", checkSkill("Survival", 8,
					[]Effect{
						pick("what the rescue taught you",
							opt("Recon", skill("Recon")),
							opt("Survival", skill("Survival", "Any"))),
						relationship(Ally, 1, ""),
					},
					[]Effect{injury(1), skill("Leadership")})),
			)},
		},
		14: {
			Summary: "caught in a violent weather event on a new world",
			Effects: []Effect{checkSkill("Survival", 8, nil, []Effect{injury(1)})},
		},
		15: {
			Summary: "life as a colonist isn't always honest",
			Effects: []Effect{pick("what the dishonesty taught you",
				opt("Deception", skill("Deception", "Any")),
				opt("Streetwise", skill("Streetwise")),
			)},
		},
		16: {
			Summary: "raiders attack the colony",
			Effects: []Effect{checkSkill("Gun Combat", 8, nil, []Effect{injury(1)})},
		},
		21: {
			Summary: "an odd beast attacks you in the wilderness",
			Effects: []Effect{pick("shoot it or outlast it",
				opt("Gun Combat", checkSkill("Gun Combat", 8, nil, []Effect{injury(1)})),
				opt("Survival", checkSkill("Survival", 8, nil, []Effect{injury(1)})),
			)},
		},
		22: {
			Summary: "you take on work that belongs to another assignment",
			Effects: []Effect{unimplemented("roll once on the skill table of an assignment other than your own")},
		},
		23: {
			Summary: "colonists work hard by day and relieve the stress at night",
			Effects: []Effect{skill("Carouse")},
		},
		24: {
			Summary: "a gambling group forms among your fellow workers",
			Effects: []Effect{checkSkill("Gambler", 8,
				[]Effect{
					benefitRolls(2, 0, ScopeBatch),
					pick("what the game taught you",
						opt("Gambler", skill("Gambler")),
						opt("Persuade", skill("Persuade"))),
				},
				[]Effect{benefitRolls(-2, 0, ScopeBatch)})},
		},
		25: {
			Summary: "you help form a local sports team",
			Effects: []Effect{pick("play or coach",
				opt("Athletics", skill("Athletics", "Any")),
				opt("Tactics", skill("Tactics", "Sport")),
			)},
		},
		26: {
			Summary: "getting anywhere in a new colony means walking",
			Effects: []Effect{pick("what the walking gave you",
				opt("+1 END", chr("END", 1)),
				opt("Navigation", skill("Navigation")),
			)},
		},
		41: {
			Summary: "quiet time on a new world, spent on your own knowledge",
			Effects: []Effect{chr("EDU", 1)},
		},
		42: {
			Summary: "you volunteer with community law enforcement",
			Effects: []Effect{pick("what the work taught you",
				opt("Investigate", skill("Investigate")),
				opt("Streetwise", skill("Streetwise")),
			)},
		},
		43: {
			Summary: "you help with the day-to-day administration of the colony",
			Effects: []Effect{pick("administer or advocate",
				opt("Admin", skill("Admin")),
				opt("Advocate (Politics), for a Politician", skill("Advocate", "Politics")),
			)},
		},
		44: {
			Summary: "you find time to practice your art",
			Effects: []Effect{skill("Art", "Any")},
		},
		45: {
			Summary: "you have something to drive or something to ride",
			Effects: []Effect{pick("what you get about on",
				opt("Drive", skill("Drive", "Wheeled", "Tracked")),
				opt("Flyer", skill("Flyer", "Any")),
				opt("Seafarer", skill("Seafarer", "Any")),
				opt("Animals (Riding)", skill("Animals", "Riding")),
			)},
		},
		46: {
			Summary: "you meet and greet a great many new colonists",
			Effects: []Effect{relationship(Contact, 0, "1d3")},
		},
		51: {
			Summary: "you find a local commodity that might increase trade",
			Effects: []Effect{checkSkill("Broker", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch)},
				[]Effect{throwModifier("next advancement roll", -2)})},
		},
		52: {
			Summary: "you take monitoring duty on the colony's advance warning system",
			Effects: []Effect{skill("Electronics", "Sensors")},
		},
		53: {
			Summary: "the colony wants everyone ready if raiders or animals strike",
			Effects: []Effect{skill("Gun Combat", "Any")},
		},
		54: {
			Summary: "you represent the colony in trade talks with another world",
			Effects: []Effect{unimplemented(
				"roll twice on the Assignment: Ambassador table of the Diplomatic Service career (p. 185)")},
		},
		55: {
			Summary: "your vehicle crashes on an exploratory trip into the wilderness",
			Effects: []Effect{
				injury(1),
				pick("what surviving taught you",
					opt("Medic", skill("Medic", "Any")),
					opt("Recon", skill("Recon")),
					opt("Survival", skill("Survival", "Any")),
					opt("Suit", skill("Suit", "Any"))),
			},
		},
		56: {
			Summary: "a fellow colonist turns out to be a famous actor researching a role",
			Effects: []Effect{
				credits("1d6x100000"),
				transfer("Celebrity", "Star", 0),
			},
		},
		61: {
			Summary: "you are recognized as a valuable member of the community",
			Effects: []Effect{benefitRolls(1, 0, ScopeBatch)},
		},
		62: {
			Summary: "you are determined to improve your scientific skills",
			Effects: []Effect{skill("Science", "Any")},
		},
		63: {
			Summary: "a colonist is often called on to do the work of several people",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		64: {
			Summary: "your contributions to the colony's success are celebrated",
			Effects: []Effect{advance()},
		},
		65: {
			Summary: "an opportunity to better your position in the colony",
			Effects: []Effect{{
				Kind:   EffectChangeAssignment,
				Detail: "change assignment, or -- already a Politician or Commercial -- change assignment or take a rank",
			}},
		},
		66: {
			Summary: "you are recognized as one of the best at what you do",
			Effects: []Effect{
				mustContinue(),
				autoSuccess("next survival roll"),
				autoSuccess("next advancement roll"),
			},
		},
	}

	// 31-36 are all the shared Life Events table (p. 120), which the book
	// prints as one spanning row rather than six.
	for result := 31; result <= 36; result++ {
		table[result] = EventRow{
			Summary: "life event",
			Effects: []Effect{lifeEvent()},
		}
	}

	return table
}
