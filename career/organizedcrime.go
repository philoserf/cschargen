package career

// OrganizedCrime is the career of pp. 246-249: killing for the organization,
// collecting for it, or running one of its businesses.
//
// Six of its eleven mishaps take away every benefit roll banked so far, and
// three of them put the character in prison for four terms. It is the
// harshest mishap table in the book.
func OrganizedCrime() Career {
	return Career{
		Tags:           []Tag{TagCriminal, TagViolent},
		Name:           "Organized Crime",
		Cite:           "pp. 246-249",
		Enlistment:     &Check{Characteristic: "END", Number: 8},
		EnlistmentMods: []EnlistmentMod{apparentAgeOver40(), perPreviousCareer()},
		MishapEjects:   true,
		Assignments: []Assignment{
			{
				Name:        "Hitman",
				Description: "killing members of the opposition",
				Survival:    Check{Characteristic: "DEX", Number: 8},
				Advancement: Check{Characteristic: "END", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Hitman", Rows: [6]Effect{
					skill("Draw"), skill("Deception", "Any"), skill("Melee", "Any"),
					skill("Gun Combat", "Any"), skill("Stealth"), skill("Tactics", "Military"),
				}},
				Ranks: [][]Effect{
					{skill("Gun Combat", "Any")},
					{skill("Melee", "Any")},
					nil,
					{skill("Stealth")},
					nil,
					{relationship(Ally, 1, "")},
					nil,
				},
			},
			{
				Name:        "Enforcer",
				Description: "a foot soldier, taking care of the organization's earnings",
				Survival:    Check{Characteristic: "END", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Enforcer", Rows: [6]Effect{
					skill("Broker"), skill("Persuade"), skill("Streetwise"),
					skill("Gun Combat", "Slug Pistol", "Energy Pistol"), skill("Melee", "Any"), skill("Recon"),
				}},
				Ranks: [][]Effect{
					{skill("Streetwise")},
					nil,
					{skill("Deception", "Any")},
					nil,
					{relationship(Contact, 1, "")},
					nil,
					{relationship(Ally, 1, "")},
				},
			},
			{
				Name:        "Operator",
				Description: "running a business for a crime organization",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Operator", Rows: [6]Effect{
					skill("Trade", "Any"), skill("Broker"), skill("Streetwise"),
					skill("Gun Combat", "Shotgun"), skill("Melee", "Any"), skill("Persuade"),
				}},
				Ranks: [][]Effect{
					{skill("Trade", "Any")},
					{skill("Streetwise")},
					nil,
					{skill("Broker")},
					nil,
					{relationship(Contact, 1, "")},
					nil,
				},
			},
		},
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), skill("Carouse"),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				pick("choose Drive (Wheeled) or Flyer (Grav)",
					opt("Drive (Wheeled)", skill("Drive", "Wheeled")),
					opt("Flyer (Grav)", skill("Flyer", "Grav"))),
				skill("Persuade"), skill("Deception", "Any"),
				skill("Gun Combat", "Any"), skill("Melee", "Any"), skill("Streetwise"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Explosives"), skill("Interrogation", "Any"), skill("Advocate", "Legal"),
				skill("Broker"), skill("Recon"), skill("Investigate"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 0, Other: chr("STR", 1)},
			{Cash: 0, Other: chr("DEX", 1)},
			{Cash: 1000, Other: chr("END", 1)},
			{Cash: 2500, Other: stashItem("a weapon")},
			{Cash: 5000, Other: relationship(Contact, 1, "")},
			{Cash: 10000, Other: relationship(Ally, 1, "")},
			{Cash: 25000, Other: relationship(Contact, 2, "")},
		},
		Mishaps: organizedCrimeMishaps(),
		Events:  organizedCrimeEvents(),
	}
}

// crimeLoseAll is the phrase six of this career's eleven mishaps use.
func crimeLoseAll() Effect {
	return loseAllBenefits("lose every benefit roll gained to this point")
}

// bountyOnYourHead is the organization's answer to leaving badly: a hunter
// who is an Enemy from here on.
func bountyOnYourHead() Effect {
	return unimplemented("a bounty is placed on your head, and a bounty hunter begins pursuit")
}

func organizedCrimeMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "you want out of the life",
			Effects: []Effect{checkSkill("Persuade", 8,
				[]Effect{benefitRolls(-3, 0, ScopeBatch), newHomeworld()},
				[]Effect{bountyOnYourHead(), relationship(Enemy, 1, "")})},
		},
		{
			Summary: "arrested, and offered the chance to turn on your bosses",
			Effects: []Effect{pick("testify, or take the four terms",
				opt("testify and be relocated",
					benefitRolls(2, 0, ScopeBatch), newHomeworld()),
				opt(`take the sentence rather than be a "rat"`,
					transfer("Prisoner", "Prisoner", 4), benefitRolls(3, 0, ScopeBatch)))},
		},
		{
			Summary: "an addiction that impairs you past the point of staying",
			Effects: []Effect{
				chr("END", -2),
				unimplemented("choose the alcohol or drug you are addicted to"),
			},
		},
		{
			Summary: "intimately involved with another member's spouse, which is a cardinal rule",
			Effects: []Effect{relationship(Ally, 1, ""), relationship(Enemy, 1, ""), crimeLoseAll()},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "your boss turns on you",
			Effects: []Effect{
				pick("fight your way out",
					opt("Gun Combat", checkSkill("Gun Combat", 8, nil, []Effect{injury(2)})),
					opt("Melee", checkSkill("Melee", 8, nil, []Effect{injury(2)}))),
				newHomeworld(),
			},
		},
		{
			Summary: `another member is a "rat", and law enforcement is investigating you`,
			Effects: []Effect{
				pick("answer the investigation however you can",
					opt("Persuade", checkSkill("Persuade", 8, crimeRatEscaped(), crimeRatCaught())),
					opt("Gun Combat", checkSkill("Gun Combat", 8, crimeRatEscaped(), crimeRatCaught())),
					opt("Melee", checkSkill("Melee", 8, crimeRatEscaped(), crimeRatCaught())),
					opt("Stealth", checkSkill("Stealth", 8, crimeRatEscaped(), crimeRatCaught()))),
			},
		},
		{
			Summary: "a rival organization hits your boss, who is killed",
			Effects: []Effect{injury(2), relationship(Rival, 1, "")},
		},
		{
			Summary: "the boss suspects a betrayal, and will not be convinced otherwise",
			Effects: []Effect{
				crimeLoseAll(),
				bountyOnYourHead(),
				relationship(Enemy, 2, ""),
			},
		},
		{
			Summary: "a turf war, and a surprise attack you are on the wrong end of",
			Effects: []Effect{injury(2), crimeLoseAll(), relationship(Enemy, 1, "")},
		},
	}
}

// crimeRatEscaped and crimeRatCaught are mishap 9's two ends: a different
// career on another world, or four terms inside with nothing to show.
func crimeRatEscaped() []Effect {
	return []Effect{newHomeworld()}
}

func crimeRatCaught() []Effect {
	return []Effect{transfer("Prisoner", "Prisoner", 4), crimeLoseAll()}
}
