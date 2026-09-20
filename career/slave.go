package career

// Slave is the career of pp. 150-154, which the book heads with the term for
// a genetically engineered human that its OGL notice reserves. That word is
// not in this repository; see CLAUDE.md.
//
// It is the third career a mishap does not eject from -- "unless it
// specifically says that you leave the career, you must remain" (p. 152) --
// and the second, after Prisoner, that cannot be chosen. Entry is a result of
// the early life tables, and only an engineered human or an uplift can be
// sent here.
//
// It has no Advanced Education table and no cash worth the name: four of its
// seven benefit rows pay nothing, and the top row pays 2,000 credits. What it
// has instead is a stash, granted at rank 1 in every assignment -- "where a
// slave keeps their personal belongings ... possessions that are forbidden by
// the slave's owner" (p. 151).
func Slave() Career {
	return Career{
		Name: "Engineered/Uplift Slave",
		Cite: "pp. 150-154",
		EnlistmentNote: "no enlistment throw, and it cannot be chosen: entry is a result " +
			"of the early life tables, and the character must be an engineered human " +
			"or an uplift (p. 150)",
		MishapEjects: false,
		Assignments:  slaveAssignments(),
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), chr("CHA", 1),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Athletics", "Any"), skill("Deception", "Any"),
				skill("Animals", "Farming"), skill("Survival", "Any"),
				skill("Melee", "Unarmed Combat"), skill("Persuade"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 0},
			{Cash: 0},
			{Cash: 0, Other: relationship(Contact, 1, "")},
			{Cash: 0, Other: relationship(Ally, 1, "")},
			{Cash: 500, Other: chr("END", 1)},
			{Cash: 1000, Other: chr("STR", 1)},
			{Cash: 2000, Other: chr("CHA", 1)},
		},
		Mishaps: slaveMishaps(),
		Events:  slaveEvents(),
	}
}

// stash is the rank 1 benefit in all three assignments.
func stash() Effect {
	return stashItem("a stash, for possessions the owner forbids")
}

func slaveAssignments() []Assignment {
	return []Assignment{
		{
			Name:        "Miner",
			Description: "mining asteroids or deep caverns for a government or corporation",
			Survival:    Check{Characteristic: "END", Number: 8},
			Advancement: Check{Characteristic: "INT", Number: 8},
			Skills: SkillTable{Kind: AssignmentSkills, Name: "Miner", Rows: [6]Effect{
				skill("Carouse"), skill("Survival", "Any"), skill("Trade", "Prospector"),
				skill("Science", "Geology"), skill("Suit", "Vacc Suit"), skill("Drive", "Mole"),
			}},
			Ranks: [][]Effect{
				{skill("Suit", "Vacc Suit")},
				{stash()},
				{skill("Trade", "Prospector")},
				nil, nil,
				{skill("Drive", "Mole")},
				{skill("Science", "Geology")},
			},
		},
		{
			Name:        "Military/Security",
			Description: "attacking other groups, or defending an owner's assets",
			Survival:    Check{Characteristic: "STR", Number: 8},
			Advancement: Check{Characteristic: "INT", Number: 8},
			Skills: SkillTable{Kind: AssignmentSkills, Name: "Military/Security", Rows: [6]Effect{
				skill("Stealth"), skill("Recon"), skill("Melee", "Unarmed Combat"),
				skill("Gun Combat", "Any"), skill("Navigation"), skill("Explosives"),
			}},
			Ranks: [][]Effect{
				{skill("Melee", "Unarmed Combat")},
				{stash()},
				{skill("Stealth")},
				nil,
				{skill("Recon")},
				nil,
				{skill("Leadership")},
			},
		},
		{
			// The skill table heads this column "Entertainer" and the
			// assignment list calls it Entertainment.
			Name:        "Entertainment",
			Description: "entertaining others for a government or corporation",
			Survival:    Check{Characteristic: "CHA", Number: 8},
			Advancement: Check{Characteristic: "CHA", Number: 8},
			Skills: SkillTable{Kind: AssignmentSkills, Name: "Entertainment", Rows: [6]Effect{
				skill("Persuade"), skill("Athletics", "Any"), skill("Etiquette"),
				skill("Art", "Any"), skill("Chef"), skill("Deception", "Any"),
			}},
			Ranks: [][]Effect{
				{skill("Art", "Any")},
				{stash()},
				{skill("Etiquette")},
				nil,
				{skill("Persuade")},
				nil,
				{skill("Deception", "Any")},
			},
		},
	}
}

// loseStash empties the stash, which five results do on the way out.
func loseStash() Effect {
	return stashItem("")
}

func slaveMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "an enemy who injures you for sport",
			Effects: []Effect{injury(1)},
		},
		{
			Summary: "an owner who has joined a supremacist group and wants their slaves dead",
			Effects: []Effect{
				checkSkill("Melee", 8, nil, []Effect{injury(2)}),
				loseStash(),
				transfer("Vagabond", "", 0),
			},
		},
		{
			Summary: "blamed for an accident that cost your owner",
			Effects: []Effect{
				unimplemented("lose two ranks, retaining any benefit already gained"),
				injury(1),
			},
		},
		{
			Summary: "an addiction",
			Effects: []Effect{
				chr("END", -2),
				unimplemented("choose the alcohol or drug you are addicted to"),
			},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "no longer useful to your owner, and sent away",
			Effects: []Effect{loseStash(), transfer("Vagabond", "", 0)},
		},
		{
			Summary: "being worked like this has taken its toll",
			Effects: []Effect{pick("choose what the work took",
				opt("STR", chr("STR", -2)),
				opt("END", chr("END", -2)))},
		},
		{
			Summary: "a crash during transport of the owner's slaves",
			Effects: []Effect{
				injury(1),
				pick("stay, or run for it on whatever you know",
					opt("stay where you are"),
					opt("Recon", slaveBreakFor("Recon")),
					opt("Stealth", slaveBreakFor("Stealth")),
					opt("Navigation", slaveBreakFor("Navigation"))),
			},
		},
		{
			// ERRATA E-1 and E-2: this result cites the injury table at
			// p. 136 and Vagabond at p. 290, which are the previous
			// edition's pages.
			Summary: "an owner who has decided robots can replace you",
			Effects: []Effect{unimplemented(
				"roll 1d6: on 1-3 you are sold, and your homeworld moves to the nearest " +
					"system with a class B or C port; on 4-6 the owner tries to end your " +
					"suffering, you escape injured -- roll once on the Injury table -- and " +
					"join the Vagabond career. The stash is kept either way")},
		},
		{
			Summary: "accused of a crime that will not get the investigation it should",
			Effects: []Effect{
				transfer("Prisoner", "Prisoner", 1),
				unimplemented("roll 1d6: on a 1 the sentence is two terms rather than one"),
			},
		},
	}
}

// slaveBreakFor is mishap 10's escape, attempted on whichever of the three
// skills the character has.
func slaveBreakFor(name string) Effect {
	return checkSkill(name, 8,
		[]Effect{loseStash(), transfer("Vagabond", "", 0)},
		[]Effect{injury(1), throwModifier("next advancement roll", -2)})
}
