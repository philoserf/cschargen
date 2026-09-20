package career

// Undergraduate is the college of pp. 86-91.
func Undergraduate() Institution {
	return Institution{
		Name:          "Undergraduate College",
		Cite:          "pp. 86-91",
		Prerequisites: []Check{{Characteristic: "EDU", Number: 6}},
		Admission:     admission(7),
		Success:       successThrow(8),
		Honors:        honorsThrow(10),
		GraduateEDU:   10,
		Skills:        undergraduateSkills(),
		Failure:       universityFailure(),
		Events:        undergraduateEvents(),
		LifeEvents:    collegiateLifeEvents(),
	}
}

// undergraduateSkills is the list of p. 87, which a graduate takes one of
// at level 2 and a washout takes one of at level 1.
func undergraduateSkills() []Effect {
	return []Effect{
		skill("Advocate", "Any"),
		skill("Animals", "Veterinary", "Training", "Farming"),
		skill("Art", "Any"),
		skill("Broker"),
		skill("Chef"),
		skill("Diplomat"),
		skill("Electronics", "Any"),
		skill("Engineer", "Any"),
		skill("Etiquette"),
		skill("Instruction"),
		skill("Investigate"),
		skill("Language", "Any"),
		skill("Mechanic"),
		skill("Pilot", "Any"),
		skill("Science", "Any"),
		skill("Trade", "Any"),
	}
}

// universityFailure is the 2d6 table of p. 88, indexed 0-10.
func universityFailure() MishapTable {
	return MishapTable{
		{
			Summary: "a catastrophic scandal, deserved or not",
			Effects: []Effect{relationshipAt(Enemy, 1, -105)},
		},
		{
			Summary: "a family emergency that forced you home",
			Effects: []Effect{loseRelative(), ratingOfRole("parent", false, 30)},
		},
		{
			Summary: "funding that disappeared",
			Effects: []Effect{unimplemented("a debt of 10,000 credits")},
		},
		{
			Summary: "academic burnout",
			Effects: []Effect{checkChr("END", 8, nil, []Effect{chr("END", -1)})},
		},
		{
			Summary: "the realisation, too late, that this path was not for you",
			Effects: []Effect{chr("CHA", -1)},
		},
		{Summary: "an ordinary failure", Effects: nil},
		{
			Summary: "recruited away before you could graduate",
			Effects: []Effect{
				unimplemented("enlist automatically in a business, military, corporate " +
					"or colonist career"),
				chr("EDU", -1),
			},
		},
		{
			Summary: "the wrong crowd, and the studies that went with them",
			Effects: []Effect{
				checkSkill("Streetwise", 8,
					[]Effect{relationshipAt(Contact, 1, 35)},
					[]Effect{transfer("Prisoner", "Prisoner", 1)}),
				throwModifier("enlistment in any criminal career", 2),
			},
		},
		{
			Summary: "campus politics, and an administration that wanted you gone",
			Effects: []Effect{
				skill("Advocate", "Politics"),
				checkSkill("Advocate", 8,
					[]Effect{relationshipAt(Contact, 1, 35)},
					[]Effect{transfer("Prisoner", "Prisoner", 1)}),
			},
		},
		{
			Summary: "violence or disaster on campus, and a university that closed",
			Effects: []Effect{throwModifier("enlistment in any career", 2)},
		},
		{
			Summary: "brilliant, defiant, and expelled",
			Effects: []Effect{relationshipAt(Enemy, 1, -140)},
		},
	}
}

// undergraduateEvents is the 2d10 table of pp. 89-91, indexed 0-18.
func undergraduateEvents() []EventRow {
	return []EventRow{
		{
			Summary: "a serious disease that does irrevocable damage",
			Effects: []Effect{pick("choose what the disease took",
				opt("STR", chrRolled("STR", "1d3", false)),
				opt("END", chrRolled("END", "1d3", false)))},
		},
		{
			Summary: "an attack by a fellow student severe enough to imprison them",
			Effects: []Effect{
				relationshipAt(Enemy, 1, -120),
				unimplemented("roll 1d6: on 1-3 also roll on the Injury table"),
			},
		},
		{
			Summary: "a group of small-time on-campus criminals",
			Effects: []Effect{
				skill("Deception", "Intrusion"),
				unimplemented("roll 1d6: on 1 a caught associate becomes an Enemy at " +
					"-110; on 2-5 you are forgotten; on 6 one of them becomes a " +
					"detective and a Contact at 55"),
			},
		},
		{
			Summary: "an elective in a subject outside your field, which you excel at",
			Effects: []Effect{pickSkill("Animals", "Athletics", "Art", "Science")},
		},
		{
			Summary: "a safari with fellow students",
			Effects: []Effect{checkSkill("Survival", 8,
				[]Effect{pick("choose what the safari taught you",
					opt("Gun Combat (Slug Rifle)", skill("Gun Combat", "Slug Rifle")),
					opt("Navigation", skill("Navigation")),
					opt("Survival", skill("Survival", "Any")))},
				[]Effect{injury(1)})},
		},
		{
			Summary: "a group that trains future business leaders",
			Effects: []Effect{pickSkill("Advocate", "Broker", "Carouse", "Diplomat", "Persuade")},
		},
		{
			Summary: "a desire to work with your hands",
			Effects: []Effect{pickSkill("Electronics", "Mechanic")},
		},
		{
			Summary: "an interest in weaponry, and a class to go with it",
			Effects: []Effect{pick("choose the class you took",
				opt("Gun Combat", skill("Gun Combat", "Slug Pistol", "Energy Pistol")),
				opt("Melee (Blade)", skill("Melee", "Blade")))},
		},
		{
			Summary: "a collegiate life event",
			Effects: []Effect{{
				Kind:   EffectCollegiateLifeEvent,
				Detail: "roll on the Collegiate Life Events table",
			}},
		},
		{
			Summary: "a local religion delved into deeply",
			Effects: []Effect{unimplemented(
				"roll 1d6: on a 6 you become deeply involved and gain Science " +
					"(Philosophy) 1; where a religion was already taken, this is a change")},
		},
		{
			Summary: "a lot of fun in your college years",
			Effects: []Effect{
				skill("Carouse"),
				unimplemented("roll 1d6: on a 1 an addiction to alcohol or a drug"),
			},
		},
		{
			Summary: "student government, activism, or civic affairs",
			Effects: []Effect{pickSkill("Advocate", "Diplomat", "Persuade")},
		},
		{
			Summary: "a love of operating a vehicle",
			Effects: []Effect{pick("choose the vehicle you took to",
				opt("Drive", skill("Drive", "Wheeled", "Tracked")),
				opt("Flyer (Grav)", skill("Flyer", "Grav")),
				opt("Seafarer", skill("Seafarer", "Sail", "Motorboats")))},
		},
		{
			Summary: "a gambling circle",
			Effects: []Effect{pick("choose what the games taught you",
				opt("Gambler", skill("Gambler")),
				opt("Deception (Lie)", skill("Deception", "Lie")),
				opt("Persuade", skill("Persuade")))},
		},
		{
			Summary: "gaming, at a table, on the worldnet, or in a holonovel",
			Effects: []Effect{pick("choose what the games taught you",
				opt("Art", skill("Art", "Performer", "Holography", "Writing")),
				opt("Carouse", skill("Carouse")),
				opt("Electronics (Computers)", skill("Electronics", "Computers")),
				opt("Science", skill("Science", "Any")),
				opt("Tactics", skill("Tactics", "Military", "Naval")))},
		},
		{
			Summary: "an exchange programme with another institution",
			Effects: []Effect{unimplemented(
				"take the background skills of the exchange world, and a level in its " +
					"primary language where it differs from your own")},
		},
		{
			Summary: "an older student who shows you how the university really works",
			Effects: []Effect{
				relationshipAt(Ally, 1, 130),
				relationshipRolledAt(Contact, "1d3", 40),
				skill("Carouse"),
			},
		},
		{
			Summary: "a colonization programme looking for college-educated settlers",
			Effects: []Effect{
				newHomeworld(),
				pick("choose what the programme taught you",
					opt("Survival", skill("Survival", "Any")),
					opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit"))),
				transfer("Colonist", "", 0),
			},
		},
		{
			Summary: "academic distinction",
			Effects: []Effect{
				{Kind: EffectHonors, Detail: "achieve honours, if the throw did not"},
				chr("EDU", 2),
			},
		},
	}
}
