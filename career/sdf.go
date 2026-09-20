package career

// The System Defense Forces are a world's own military: a navy, ground
// troopers, and — on a wet world — a surface and undersea navy (pp.
// 278-291).
//
// All three take a commission, and all three route their d66 events to the
// Military Events table (p. 121).

// SystemDefenceNavy is the career of pp. 278-281.
//
// Its tables are National Navy's, cell for cell: the same four assignments
// with the same skill rows, the same Service, Advanced Education and
// Officer tables, and the same two rank tables (compare pp. 233-234 with
// pp. 278-279). What differs is the enlistment throw, the commission
// target, and cash figures at roughly half. They are built from one set of
// tables here because they are one set of tables in the book.
//
// ERRATA E-10: this career prints no Mishap table. Troopers (p. 284) and
// Wet Navy (p. 289) both print one; this one goes straight from its ranks
// to its events. It uses National Navy's, for the reason in ERRATA.md.
func SystemDefenceNavy() Career {
	return Career{
		Name:           "System Defense Forces (Navy)",
		Cite:           "pp. 278-281",
		Enlistment:     &Check{Characteristic: "INT", Number: 6},
		EnlistmentMods: []EnlistmentMod{apparentAgeOver40(), perPreviousCareer()},
		Commission:     &Check{Characteristic: "EDU", Number: 7},
		MishapEjects:   true,
		Assignments:    navyAssignments(),
		Tables:         navyTables(),
		Benefits: [7]BenefitRow{
			{Cash: 125, Other: chr("END", 1)},
			{Cash: 250, Other: chr("DEX", 1)},
			{Cash: 500, Other: chr("INT", 1)},
			{Cash: 1000, Other: chr("EDU", 1)},
			{Cash: 2500, Other: weaponOrItsUse()},
			{Cash: 5000, Other: relationship(Contact, 1, "")},
			{Cash: 10000, Other: relationship(Ally, 1, "")},
		},
		Mishaps: navyMishaps(),
		Events:  sdfNavyEvents(),
	}
}

// SystemDefenceTroopers is the career of pp. 282-286: a world's ground
// forces.
func SystemDefenceTroopers() Career {
	return Career{
		Name:           "System Defense Forces (Troopers)",
		Cite:           "pp. 282-286",
		Enlistment:     &Check{Characteristic: "END", Number: 8},
		EnlistmentMods: []EnlistmentMod{apparentAgeOver40(), perPreviousCareer()},
		Commission:     &Check{Characteristic: "EDU", Number: 7},
		RankTitles: []string{
			"Recruit", "Trooper", "Corporal", "Sergeant",
			"Staff Sergeant", "Sergeant Major", "Command Sergeant Major",
		},
		OfficerTitles: []string{
			"Second Lieutenant", "First Lieutenant", "Captain", "Major",
			"Lieutenant Colonel", "Colonel", "General",
		},
		MishapEjects: true,
		Assignments:  trooperAssignments(),
		Tables:       trooperTables(),
		Benefits: [7]BenefitRow{
			{Cash: 125, Other: chr("END", 1)},
			{Cash: 250, Other: chr("DEX", 1)},
			{Cash: 500, Other: chr("INT", 1)},
			{Cash: 1000, Other: chr("EDU", 1)},
			{Cash: 2500, Other: weaponOrItsUse()},
			{Cash: 5000, Other: relationship(Contact, 1, "")},
			{Cash: 10000, Other: stashItem("armor of the character's choice")},
		},
		Mishaps: trooperMishaps(),
		Events:  trooperEvents(),
	}
}

// SystemDefenceWetNavy is the career of pp. 287-291: a world's surface and
// undersea navy.
//
// It is the one career with a homeworld prerequisite: "Your homeworld must
// have a Hydrographics score of 4+ to join this career" (p. 287). The
// setting data carries no hydrographics figure -- it is part of a world's
// UWP, which is in the Core Setting Book rather than in the origin charts
// -- so the prerequisite is recorded and not enforced.
func SystemDefenceWetNavy() Career {
	return Career{
		Name:           "System Defense Forces (Wet Navy)",
		Cite:           "pp. 287-291",
		Enlistment:     &Check{Characteristic: "END", Number: 8},
		EnlistmentMods: []EnlistmentMod{apparentAgeOver40(), perPreviousCareer()},
		Commission:     &Check{Characteristic: "EDU", Number: 9},
		Prerequisite: "the homeworld must have a Hydrographics score of 4+ (p. 287), " +
			"which the origin charts do not carry",
		RankTitles: []string{
			"Seaman Recruit", "Seaman Apprentice", "Seaman", "Petty Officer 3rd Class",
			"Petty Officer 2nd Class", "Petty Officer 1st Class", "Chief Petty Officer",
		},
		OfficerTitles: []string{
			"Ensign", "Lieutenant, Junior Grade", "Lieutenant", "Lieutenant Commander",
			"Commander", "Captain", "Admiral",
		},
		MishapEjects: true,
		Assignments:  wetNavyAssignments(),
		Tables:       wetNavyTables(),
		Benefits: [7]BenefitRow{
			{Cash: 250, Other: chr("STR", 1)},
			{Cash: 500, Other: chr("DEX", 1)},
			{Cash: 1000, Other: chr("END", 1)},
			{Cash: 2000, Other: chr("EDU", 1)},
			{Cash: 4000, Other: relationship(Contact, 1, "")},
			{Cash: 10000, Other: weaponOrItsUse()},
			{Cash: 20000, Other: credits("2d6x10000")},
		},
		Mishaps: wetNavyMishaps(),
		Events:  wetNavyEvents(),
	}
}

func trooperAssignments() []Assignment {
	officer := [][]Effect{
		nil,
		{skill("Persuade")},
		{skill("Tactics", "Military")},
		{skill("Survival", "Any")},
		{skill("Leadership")},
		nil,
		{skill("Advocate", "Any")},
	}

	rifleman := append([][]Effect{{skill("Gun Combat", "Slug Rifle")}}, officer[1:]...)
	recon := append([][]Effect{{skill("Recon")}}, officer[1:]...)
	heavy := append([][]Effect{{skill("Heavy Weapons", "Any")}}, officer[1:]...)

	return []Assignment{
		{
			Name:        "Rifleman",
			Description: "a combat rifle and the will to fight",
			Survival:    Check{Characteristic: "DEX", Number: 7},
			Advancement: Check{Characteristic: "INT", Number: 8},
			Skills: SkillTable{Kind: AssignmentSkills, Name: "Rifleman", Rows: [6]Effect{
				skill("Athletics", "Any"), skill("Navigation"), skill("Melee", "Any"),
				skill("Gun Combat", "Any"), skill("Survival", "Any"), skill("Electronics", "Any"),
			}},
			Ranks: [][]Effect{
				{skill("Gun Combat", "Slug Rifle")},
				{skill("Discipline")},
				nil,
				{skill("Mechanic")},
				nil,
				{skill("Survival", "Any")},
				{skill("Leadership")},
			},
			OfficerRanks: rifleman,
		},
		{
			Name:        "Recon",
			Description: "lightly armed, scouting the enemy and calling in the forces that deal with them",
			Survival:    Check{Characteristic: "DEX", Number: 8},
			Advancement: Check{Characteristic: "EDU", Number: 8},
			Skills: SkillTable{Kind: AssignmentSkills, Name: "Recon", Rows: [6]Effect{
				skill("Athletics", "Any"), skill("Electronics", "Any"), skill("Recon"),
				skill("Gun Combat", "Any"), skill("Survival", "Any"), skill("Stealth"),
			}},
			Ranks: [][]Effect{
				{skill("Recon")},
				{skill("Discipline")},
				nil,
				{skill("Stealth")},
				nil,
				{skill("Survival", "Any")},
				{skill("Leadership")},
			},
			OfficerRanks: recon,
		},
		{
			Name:        "Hvy Weaps",
			Description: "the big guns, providing fire support",
			Survival:    Check{Characteristic: "END", Number: 7},
			Advancement: Check{Characteristic: "INT", Number: 8},
			Skills: SkillTable{Kind: AssignmentSkills, Name: "Hvy Weaps", Rows: [6]Effect{
				skill("Drive", "Tracked", "Wheeled"), skill("Gun Combat", "Any"),
				skill("Heavy Weapons", "Any"), skill("Explosives"),
				skill("Survival", "Any"), skill("Electronics", "Any"),
			}},
			Ranks: [][]Effect{
				{skill("Heavy Weapons", "Any")},
				{skill("Discipline")},
				{skill("Gun Combat", "Any")},
				nil, nil,
				{skill("Survival", "Any")},
				{skill("Leadership")},
			},
			OfficerRanks: heavy,
		},
	}
}

func trooperTables() []SkillTable {
	return []SkillTable{
		personalDevelopment(),
		{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
			skill("Discipline"), skill("Melee", "Any"), skill("Mechanic"),
			skill("Gun Combat", "Any"), skill("Navigation"), skill("Survival", "Any"),
		}},
		{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
			skill("Medic", "Any"), skill("Flyer", "Any"), skill("Electronics", "Any"),
			skill("Science", "Any"), skill("Advocate", "Any"), skill("Tactics", "Military"),
		}},
		{Kind: OfficerSkills, Name: "Officer Skills", Rows: [6]Effect{
			skill("Admin"), skill("Instruction"), skill("Etiquette"),
			skill("Persuade"), skill("Tactics", "Military"), skill("Leadership"),
		}},
	}
}

func wetNavyAssignments() []Assignment {
	return []Assignment{
		{
			Name:        "Surface Ship",
			Description: "crew or officer aboard a surface ship",
			Survival:    Check{Characteristic: "END", Number: 6},
			Advancement: Check{Characteristic: "EDU", Number: 7},
			Skills: SkillTable{Kind: AssignmentSkills, Name: "Surface Ship", Rows: [6]Effect{
				skill("Electronics", "Any"), skill("Navigation"), skill("Heavy Weapons", "Any"),
				skill("Gunner", "Turrets"), skill("Engineer", "Power"), skill("Survival", "Ocean"),
			}},
			Ranks: [][]Effect{
				{skill("Seafarer", "Ocean Ships")},
				{skill("Survival", "Ocean")},
				nil,
				{skill("Electronics", "Any")},
				nil, nil,
				{skill("Leadership")},
			},
			OfficerRanks: [][]Effect{
				{skill("Seafarer", "Ocean Ships")},
				{skill("Survival", "Ocean")},
				{skill("Admin")},
				nil,
				{skill("Leadership")},
				nil,
				{skill("Advocate", "Any")},
			},
		},
		{
			Name:        "Aviator",
			Description: "the wet navy's airborne strike and reconnaissance force",
			Survival:    Check{Characteristic: "DEX", Number: 7},
			Advancement: Check{Characteristic: "INT", Number: 8},
			Skills: SkillTable{Kind: AssignmentSkills, Name: "Aviator", Rows: [6]Effect{
				skill("Mechanic"), skill("Suit", "Vacc Suit"), skill("Flyer", "Any"),
				skill("Pilot", "Small Craft"), skill("Gunner", "Turrets"), skill("Electronics", "Any"),
			}},
			Ranks: [][]Effect{
				{skill("Mechanic")},
				nil, nil,
				{skill("Electronics", "Any")},
				{skill("Flyer", "Grav")},
				{skill("Suit", "Vacc Suit")},
				{skill("Leadership")},
			},
			OfficerRanks: [][]Effect{
				{skill("Flyer", "Grav")},
				{skill("Suit", "Vacc Suit")},
				{skill("Admin")},
				nil,
				{skill("Leadership")},
				nil,
				{skill("Advocate", "Any")},
			},
		},
		{
			Name:        "Submariner",
			Description: "crew or officer aboard an undersea vessel",
			Survival:    Check{Characteristic: "END", Number: 7},
			Advancement: Check{Characteristic: "EDU", Number: 7},
			Skills: SkillTable{Kind: AssignmentSkills, Name: "Submariner", Rows: [6]Effect{
				skill("Electronics", "Any"), skill("Navigation"), skill("Seafarer", "Submarine"),
				skill("Mechanic"), skill("Engineer", "Life Support", "Power"), skill("Suit", "Vacc Suit"),
			}},
			Ranks: [][]Effect{
				{skill("Seafarer", "Submarine")},
				{skill("Mechanic")},
				nil,
				{skill("Electronics", "Any")},
				nil, nil,
				{skill("Leadership")},
			},
			OfficerRanks: [][]Effect{
				{skill("Seafarer", "Submarine")},
				{skill("Electronics", "Sensors")},
				{skill("Admin")},
				nil,
				{skill("Leadership")},
				nil,
				{skill("Advocate", "Any")},
			},
		},
	}
}

func wetNavyTables() []SkillTable {
	return []SkillTable{
		personalDevelopment(),
		{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
			skill("Discipline"), skill("Mechanic"), skill("Seafarer", "Any"),
			skill("Gun Combat", "Any"), skill("Carouse"), skill("Melee", "Any"),
		}},
		{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
			skill("Admin"), skill("Advocate", "Any"), skill("Science", "Any"),
			skill("Medic", "Any"), skill("Language", "Any"), skill("Etiquette"),
		}},
		{Kind: OfficerSkills, Name: "Officer Skills", Rows: [6]Effect{
			skill("Admin"), skill("Etiquette"), skill("Tactics", "Any"),
			skill("Diplomat"), skill("Persuade"), skill("Leadership"),
		}},
	}
}
