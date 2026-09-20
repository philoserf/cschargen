package career

// Explorer is the career of pp. 194-197: surveying new systems, walking
// the planets, or keeping both alive.
func Explorer() Career {
	return Career{
		Name:         "Explorer",
		Cite:         "pp. 194-197",
		Enlistment:   &Check{Characteristic: "END", Number: 8},
		MishapEjects: true,
		Assignments: []Assignment{
			{
				Name:        "Survey",
				Description: "the science team that surveys new star systems",
				Survival:    Check{Characteristic: "END", Number: 7},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Survey", Rows: [6]Effect{
					skill("Astrogation"), skill("Pilot", "Any"), skill("Electronics", "Any"),
					skill("Science", "Any"), skill("Engineer", "Any"), skill("Investigate"),
				}},
				Ranks: [][]Effect{
					{skill("Science", "Orbital Mechanics")},
					{skill("Electronics", "Sensors")},
					nil,
					{skill("Astrogation")},
					{skill("Investigate")},
					nil,
					{skill("Leadership")},
				},
			},
			{
				Name:        "Explorer",
				Description: "the planetary exploration team",
				Survival:    Check{Characteristic: "END", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Explorer", Rows: [6]Effect{
					skill("Animals", "Any"), skill("Navigation"), skill("Science", "Any"),
					skill("Survival", "Any"), pickSkill("Flyer", "Drive"), skill("Jack of All Trades"),
				}},
				Ranks: [][]Effect{
					{skill("Science", "Planetology")},
					{skill("Electronics", "Sensors")},
					nil,
					{skill("Investigate")},
					{skill("Recon")},
					nil,
					{skill("Leadership")},
				},
			},
			{
				Name:        "Escort",
				Description: "the armed escort that assists the survey and explorer teams",
				Survival:    Check{Characteristic: "DEX", Number: 7},
				Advancement: Check{Characteristic: "EDU", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Escort", Rows: [6]Effect{
					skill("Gunner", "Any"), skill("Pilot", "Small Craft"), skill("Gun Combat", "Any"),
					skill("Recon"), skill("Navigation"), skill("Draw"),
				}},
				Ranks: [][]Effect{
					{skill("Gun Combat", "Any")},
					{skill("Recon")},
					nil,
					{skill("Gunner", "Turrets")},
					nil,
					{skill("Leadership")},
					nil,
				},
			},
		},
		Tables: []SkillTable{
			personalDevelopment(),
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Pilot", "Any"), skill("Investigate"), skill("Survival", "Any"),
				skill("Astrogation"), skill("Navigation"), skill("Suit", "Any"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Medic", "Any"), skill("Language", "Any"), skill("Science", "Any"),
				skill("Survival", "Any"), skill("Engineer", "Any"), skill("Diplomat"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 1000, Other: chr("END", 1)},
			{Cash: 2000, Other: chr("INT", 1)},
			{Cash: 3000, Other: chr("EDU", 1)},
			{Cash: 5000, Other: weaponOrItsUse()},
			{Cash: 7000, Other: stashItem("a suit")},
			{Cash: 10000, Other: relationship(Contact, 1, "")},
			{Cash: 25000, Other: relationship(Ally, 1, "")},
		},
		Mishaps: explorerMishaps(),
		Events:  explorerEvents(),
	}
}

func explorerMishaps() MishapTable {
	return MishapTable{
		// ERRATA E-1: the book cites p. 136 here.
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{Summary: "accidental exposure to a dangerous atmosphere", Effects: []Effect{chr("END", -2)}},
		{
			Summary: "a spacewalk accident sets you adrift, and the damage is permanent",
			Effects: []Effect{chr("INT", -2), skill("Suit", "Vacc Suit")},
		},
		{
			Summary: "accused of negligence resulting in a crewmate's death",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{Summary: "a decompression aboard ship", Effects: []Effect{chr("END", -1)}},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "your ship found dead in space years later, with you in a cold berth and no memory",
			Effects: []Effect{relationship(Contact, 0, "1d6")},
		},
		{
			Summary: "an alien creature takes you by surprise and nearly kills you",
			Effects: []Effect{injury(2)},
		},
		{
			Summary: "more radiation than the ship was designed for",
			Effects: []Effect{chr("END", -1)},
		},
		{Summary: "an alien disease or virus, caught on a new world", Effects: []Effect{chr("END", -3)}},
		{
			Summary: "your ship destroyed, and years living on an alien planet",
			Effects: []Effect{
				skill("Survival", "Any"),
				injury(1),
				unimplemented("for every four of the 2d6 years, rounding up, a roll on the Explorer table"),
			},
		},
	}
}

// aPirateBase is the base stumbled upon, printed in Adventurer (p. 148) and
// Explorer (p. 196) with the same shape: a 10+ attempt to talk, a fight if
// it fails, and an Enemy either way.
func aPirateBase(talk []Option, fight Effect) Effect {
	return pick("talk them down or fight", append(talk, opt("fight", fight))...)
}

// explorerEvents is the d66 table of pp. 196-197.
func explorerEvents() EventTable {
	fight := checkSkill("Gun Combat", 8, []Effect{skill("Tactics", "Any")}, []Effect{injury(1)})

	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {Summary: "a great deal of the term spent filling out forms", Effects: []Effect{skill("Admin")}},
		13: {
			Summary: "selected for cross-training in an alternate assignment",
			Effects: []Effect{rollOtherAssignment()},
		},
		14: {
			Summary: "selected as the team driver",
			Effects: []Effect{pickSkill("Drive", "Flyer", "Seafarer")},
		},
		15: {Summary: "a gambling circle aboard ship", Effects: []Effect{gamblingCircle("Persuade")}},
		16: {
			Summary: "chosen for the organization's training staff",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch), skill("Instruction")},
		},
		21: {
			Summary: "an animal befriended on a world you were exploring",
			Effects: []Effect{stashItem("a pet"), skill("Animals", "Any")},
		},
		22: {
			Summary: "you spot a crewmate doing something dangerous or illegal",
			Effects: []Effect{spottedSomethingWrong(skill("Leadership"))},
		},
		23: {
			Summary: "an odd beast attacks you on a new world",
			Effects: []Effect{pick("meet it however you can",
				opt("Gun Combat", checkSkill("Gun Combat", 8, explorerBeat(), []Effect{injury(1)})),
				opt("Survival", checkSkill("Survival", 8, explorerBeat(), []Effect{injury(1)})))},
		},
		24: {
			Summary: "caught in a violent weather event on a new world",
			Effects: []Effect{checkSkill("Survival", 8, nil, []Effect{injury(1)})},
		},
		25: {
			Summary: "a rescue in the outer reaches of a frontier system",
			Effects: []Effect{checkSkill("Suit", 8,
				[]Effect{relationship(Contact, 1, "")}, []Effect{injury(1)})},
		},
		26: {Summary: "chosen for advanced training", Effects: []Effect{rollTable(AdvancedEducation)}},
		41: {
			Summary: "downtime aboard ship spent on physical training",
			Effects: []Effect{pick("what the training gave you",
				opt("+1 STR", chr("STR", 1)),
				opt("+1 DEX", chr("DEX", 1)),
				opt("+1 END", chr("END", 1)),
				opt("Athletics", skill("Athletics", "Any")),
				opt("Melee", skill("Melee", "Any")))},
		},
		42: {
			Summary: "extra time on the holographic gun range",
			Effects: []Effect{skill("Gun Combat", "Any")},
		},
		43: {
			Summary: "time spent honing your skillset",
			Effects: []Effect{raiseHeldSkill()},
		},
		44: {
			Summary: "a crew from many diverse backgrounds",
			Effects: []Effect{skill("Language", "Any")},
		},
		45: {Summary: "often called on for a variety of jobs", Effects: []Effect{skill("Jack of All Trades")}},
		46: {
			Summary: "a pirate base, stumbled upon by the exploration team",
			Effects: []Effect{
				aPirateBase([]Option{
					opt("Diplomat at 10+", checkSkill("Diplomat", 10, nil, []Effect{fight})),
					opt("Persuade at 10+", checkSkill("Persuade", 10, nil, []Effect{fight})),
				}, fight),
				relationship(Enemy, 1, ""),
			},
		},
		51: {Summary: "downtime aboard ship, spent on hobbies", Effects: []Effect{skill("Art", "Any")}},
		52: {
			Summary: "teams in the field repair their own equipment",
			Effects: []Effect{skill("Mechanic")},
		},
		53: {
			Summary: "a shuttle crash leaves you surviving in the wilderness",
			Effects: []Effect{
				injury(1),
				pick("what surviving taught you",
					opt("Medic (First Aid)", skill("Medic", "First Aid")),
					opt("Recon", skill("Recon")),
					opt("Survival", skill("Survival", "Any")),
					opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit"))),
			},
		},
		54: {
			Summary: "hiding from large alien animals planetside",
			Effects: []Effect{skill("Stealth")},
		},
		55: {
			Summary: "a squabble over the last mission becomes a full-time rivalry",
			Effects: []Effect{relationship(Rival, 1, "")},
		},
		56: {
			Summary: "a conference on exploring the frontier, and a jealous crew",
			Effects: []Effect{
				relationship(Contact, 0, "1d3"),
				relationship(Rival, 0, "1d3"),
			},
		},
		61: {
			Summary: "downtime spent in the company of your crewmates",
			Effects: []Effect{skill("Carouse")},
		},
		62: {
			Summary: "extra time on a new planet",
			Effects: []Effect{rollAssignmentOf("Explorer", "Explorer")},
		},
		63: {
			Summary: "determined to improve your scientific skills",
			Effects: []Effect{skill("Science", "Any")},
		},
		64: {
			Summary: "a superior sets out to groom you for higher things",
			Effects: groomedForHigherThings(),
		},
		65: {
			Summary: "recognized for exemplary service",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch), benefitRolls(0, 1, ScopeCareer)},
		},
		66: {
			Summary: "a new world found, suitable for a colony and astrogationally advantageous",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch), autoSuccess("next advancement roll")},
		},
	}

	lifeEventRows(table)

	return table
}

func explorerBeat() []Effect {
	return []Effect{pick("what surviving it taught you",
		opt("Gun Combat", skill("Gun Combat", "Any")),
		opt("Survival", skill("Survival", "Any")),
		opt("Animals", skill("Animals", "Any")))}
}
