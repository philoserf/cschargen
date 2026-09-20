package career

// OrbitalConstruction is the career of pp. 242-245: building orbital
// habitats from the outside, from the inside, or hauling the pieces.
func OrbitalConstruction() Career {
	return Career{
		Name:           "Orbital Construction",
		Cite:           "pp. 242-245",
		Enlistment:     &Check{Characteristic: "END", Number: 8},
		EnlistmentMods: []EnlistmentMod{apparentAgeOver40()},
		MishapEjects:   true,
		Assignments: []Assignment{
			{
				Name:        "Suitjockey",
				Description: "building the outer shell of orbital habitats",
				Survival:    Check{Characteristic: "END", Number: 7},
				Advancement: Check{Characteristic: "INT", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Suitjockey", Rows: [6]Effect{
					skill("Engineer", "Life Support", "Power"),
					skill("Survival", "Freefall", "Low Gravity", "Low Pressure"),
					skill("Trade", "Space Construction"), skill("Electronics", "Any"),
					skill("Suit", "Vacc Suit"), skill("Jack of All Trades"),
				}},
				Ranks: [][]Effect{
					{skill("Survival", "Freefall")},
					{skill("Suit", "Vacc Suit")},
					{skill("Trade", "Space Construction")},
					nil,
					{skill("Engineer", "Power")},
					{skill("Admin")},
					{skill("Leadership")},
				},
			},
			{
				Name:        "Jobber",
				Description: "building the internal living spaces of orbital habitats",
				Survival:    Check{Characteristic: "END", Number: 6},
				Advancement: Check{Characteristic: "EDU", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Jobber", Rows: [6]Effect{
					skill("Engineer", "Life Support", "Power"),
					skill("Survival", "Freefall", "Low Gravity", "Low Pressure"),
					skill("Electronics", "Any"), skill("Trade", "Construction"),
					skill("Suit", "Vacc Suit"), skill("Jack of All Trades"),
				}},
				Ranks: [][]Effect{
					{skill("Survival", "Freefall")},
					{skill("Suit", "Vacc Suit")},
					{skill("Trade", "Construction")},
					nil,
					{skill("Engineer", "Life Support")},
					{skill("Admin")},
					{skill("Leadership")},
				},
			},
			{
				Name:        "Lifter",
				Description: "small craft carrying cargo and large objects to and from the site",
				Survival:    Check{Characteristic: "DEX", Number: 7},
				Advancement: Check{Characteristic: "EDU", Number: 7},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Lifter", Rows: [6]Effect{
					skill("Mechanic"), skill("Drive", "Wheeled", "Tracked"), skill("Flyer", "Grav"),
					skill("Pilot", "Small Craft"), skill("Suit", "Vacc Suit"), skill("Engineer", "Any"),
				}},
				Ranks: [][]Effect{
					{skill("Pilot", "Small Craft")},
					{skill("Suit", "Vacc Suit")},
					{skill("Mechanic")},
					nil,
					{skill("Admin")},
					nil,
					{skill("Leadership")},
				},
			},
		},
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), chr("CHA", 1),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Admin"), skill("Athletics", "Climbing"), skill("Carouse"),
				skill("Suit", "Vacc Suit"), skill("Trade", "Construction"), skill("Mechanic"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Language", "Any"), skill("Art", "Any"), skill("Broker"),
				skill("Etiquette"), skill("Diplomat"), skill("Instruction"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 0, Other: chr("DEX", 1)},
			{Cash: 250, Other: chr("STR", 1)},
			{Cash: 500, Other: chr("END", 1)},
			{Cash: 1000, Other: chr("CHA", 1)},
			{Cash: 2000, Other: stashValued("a tool kit", "1000")},
			{Cash: 3000, Other: stashItem("a vacc suit")},
			{Cash: 5000, Other: stashValued(companyShare, companyShareValue)},
		},
		Mishaps: orbitalMishaps(),
		Events:  orbitalEvents(),
	}
}

func orbitalMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{Summary: "you are not cut out for this kind of work", Effects: nil},
		{
			Summary: "a downturn cancels every project",
			Effects: []Effect{
				benefitRolls(-1, 0, ScopeBatch),
				unimplemented("-1 to the next non-military enlistment roll"),
			},
		},
		{
			Summary: "constant arguments with your boss end your contract",
			Effects: []Effect{benefitRolls(1, 0, ScopeBatch)},
		},
		{Summary: "accidental exposure to a dangerous atmosphere", Effects: []Effect{chr("END", -2)}},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "an accident sets you adrift, and the psychological damage is permanent",
			Effects: []Effect{chr("INT", -2), skill("Suit", "Vacc Suit")},
		},
		{
			Summary: "accused of negligence resulting in a teammate's death",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "a solar flare catches you working on a new station",
			Effects: []Effect{chr("END", -1)},
		},
		{
			Summary: "a decompression, in your suit or aboard the station",
			Effects: []Effect{chr("END", -1)},
		},
		{
			Summary: "the site is attacked, and the projects end either way",
			Effects: []Effect{pick("get away however you can",
				opt("Gun Combat", checkSkill("Gun Combat", 8, nil, []Effect{injury(2)})),
				opt("Deception", checkSkill("Deception", 8, nil, []Effect{injury(2)})),
				opt("Stealth", checkSkill("Stealth", 8, nil, []Effect{injury(2)})))},
		},
	}
}

// shoddyPractices is the project manager cutting corners, printed
// identically in Craftsperson (p. 184) and Orbital Construction (p. 245).
func shoddyPractices() Effect {
	return pick("report it or ignore it",
		opt("report it",
			skill("Advocate", "Legal"),
			throwModifier("next advancement roll", -4)),
		opt("ignore it",
			relationship(Ally, 1, ""),
			throwModifier("next advancement roll", 4),
			benefitRolls(2, 0, ScopeBatch),
			unimplemented(
				"on a 1d6 of 3 or less, the practice kills several people years later "+
					"and a world bars the character from returning")),
	)
}

// orbitalEvents is the d66 table of pp. 244-245.
func orbitalEvents() EventTable {
	table := EventTable{
		11: {Summary: "disaster occurs", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "a slowdown in the trade, and a move into the head office",
			Effects: []Effect{skill("Admin")},
		},
		13: {
			Summary: "downtime spent getting into trouble",
			Effects: []Effect{pickSkill("Carouse", "Deception", "Gambler", "Streetwise")},
		},
		14: {
			Summary: "a shortage of hands puts you outside your assignment",
			Effects: []Effect{rollOtherAssignment()},
		},
		15: {Summary: "downtime at the gun range", Effects: []Effect{skill("Gun Combat", "Any")}},
		16: {
			Summary: "construction sites are often unruly places",
			Effects: []Effect{pickSkill("Deception", "Melee", "Streetwise")},
		},
		21: {
			Summary: "a sea of regulations and local codes",
			Effects: []Effect{skill("Advocate", "Legal", "Politics")},
		},
		22: {
			Summary: "a planetside construction job",
			Effects: []Effect{pickSkill("Explosives", "Navigation", "Survival")},
		},
		23: {Summary: "a gambling circle", Effects: []Effect{gamblingCircle("Persuade")}},
		24: {
			Summary: "often chosen for jobs outside your comfort zone",
			Effects: []Effect{skill("Jack of All Trades")},
		},
		25: {
			Summary: "a barfight",
			Effects: []Effect{checkSkill("Melee", 8,
				[]Effect{relationship(Contact, 1, "")}, []Effect{injury(1)})},
		},
		26: {
			Summary: "a gossiping coworker spreads something damaging about your past",
			Effects: []Effect{relationship(Rival, 1, ""), throwModifier("next advancement roll", -2)},
		},
		41: {
			Summary: "downtime spent on a hobby",
			Effects: []Effect{pick("what you pursued",
				opt("Animals (Riding)", skill("Animals", "Riding")),
				opt("Art", skill("Art", "Any")),
				opt("Science", skill("Science", "Any")))},
		},
		42: {
			Summary: "you spot a teammate doing something dangerous or illegal",
			Effects: []Effect{spottedSomethingWrong(skill("Leadership"))},
		},
		43: {Summary: "repairing your own equipment", Effects: []Effect{skill("Mechanic")}},
		44: {
			Summary: "a politically active coworker, popular with the workers and not the owners",
			Effects: []Effect{pick("side with them, side with the company, or ignore it",
				opt("side with the coworker",
					relationship(Ally, 1, ""),
					throwModifier("next advancement roll", -2)),
				opt("side with the employers",
					throwModifier("next advancement roll", 2),
					relationship(Rival, 1, "")),
				opt("ignore them", relationship(Rival, 1, ""), skill("Etiquette")))},
		},
		45: {Summary: "working extra hard to get things right", Effects: []Effect{rollTable(ServiceSkills)}},
		46: {
			Summary: "bad planning runs the job over budget, and your pay is cut",
			Effects: []Effect{benefitRolls(-1, 0, ScopeBatch), throwModifier("next advancement roll", 2)},
		},
		51: {
			Summary: "a sports team put together for a local league",
			Effects: []Effect{pick("what the sport gave you",
				opt("+1 STR", chr("STR", 1)),
				opt("+1 DEX", chr("DEX", 1)),
				opt("+1 END", chr("END", 1)))},
		},
		52: {Summary: "chosen for advanced training", Effects: []Effect{rollTable(AdvancedEducation)}},
		53: {Summary: "basic first aid", Effects: []Effect{skill("Medic", "First Aid")}},
		54: {Summary: "assigned to assist the cook staff", Effects: []Effect{skill("Chef")}},
		55: {
			Summary: "chosen for the organization's training staff",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch), skill("Instruction")},
		},
		56: {
			Summary: "shoddy practices by your project manager, in violation of the codes",
			Effects: []Effect{shoddyPractices()},
		},
		61: {Summary: "part of the bidding process for a new project", Effects: []Effect{skill("Broker")}},
		62: {Summary: "a team from many diverse backgrounds", Effects: []Effect{skill("Language", "Any")}},
		63: {
			Summary: "excellent work by the team earns a bonus",
			Effects: []Effect{cashRollsNowRerolling(2)},
		},
		64: {
			Summary: "honing your skillset",
			Effects: []Effect{unimplemented("raise a skill the character already holds")},
		},
		65: {
			Summary: "a superior sets out to groom you for higher things",
			Effects: groomedForHigherThings(),
		},
		66: {
			Summary: "the project becomes one of the sector's most recognized structures",
			Effects: []Effect{benefitRolls(3, 0, ScopeBatch), benefitRolls(0, 1, ScopeCareer)},
		},
	}

	lifeEventRows(table)

	return table
}
