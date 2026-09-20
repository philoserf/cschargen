package career

// CorporateShipper is the career of pp. 177-180: crewing a corporation's
// freighters.
//
// Its rank table is the first in the book that is not per-assignment: one
// column of benefits and titles serves all three assignments (p. 178).
func CorporateShipper() Career {
	ranks := [][]Effect{
		nil,
		{rollTable(AssignmentSkills)},
		{skill("Persuade")},
		{skill("Etiquette")},
		{skill("Broker")},
		{skill("Admin")},
		{chr("CHA", 1)},
		{skill("Diplomat")},
	}

	return Career{
		Tags:           []Tag{TagCorporate},
		Name:           "Corporate Shipper",
		Cite:           "pp. 177-180",
		Enlistment:     &Check{Characteristic: "INT", Number: 8},
		EnlistmentMods: []EnlistmentMod{perPreviousCareer()},
		RankTitles: []string{
			"Crewman", "Senior Crewman", "4th Officer", "3rd Officer",
			"2nd Officer", "1st Officer", "Captain", "Senior Captain",
		},
		MishapEjects: true,
		Assignments: []Assignment{
			{
				Name:        "Crew",
				Description: "the general crew or officer ranks of a corporate shipper",
				Survival:    Check{Characteristic: "INT", Number: 7},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Crew", Rows: [6]Effect{
					skill("Chef"), pickSkill("Drive", "Flyer"), skill("Pilot", "Any"),
					skill("Astrogation"), skill("Gun Combat", "Any"), skill("Electronics", "Any"),
				}},
				Ranks: ranks,
			},
			{
				Name:        "Engineer",
				Description: "the engineering crew of a corporate shipper",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Engineer", Rows: [6]Effect{
					skill("Suit", "Vacc Suit"), skill("Survival", "Freefall"), skill("Engineer", "Any"),
					skill("Engineer", "Any"), skill("Mechanic"), skill("Jack of All Trades"),
				}},
				Ranks: ranks,
			},
			{
				Name:        "Gunnery",
				Description: "the gunnery crew that protects the ship",
				Survival:    Check{Characteristic: "DEX", Number: 8},
				Advancement: Check{Characteristic: "INT", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Gunnery", Rows: [6]Effect{
					skill("Gun Combat", "Any"), skill("Electronics", "Sensors"), skill("Gunner", "Turrets"),
					skill("Gunner", "Turrets"), skill("Tactics", "Naval"), skill("Electronics", "Any"),
				}},
				Ranks: ranks,
			},
		},
		Tables: []SkillTable{
			personalDevelopment(),
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Admin"), skill("Suit", "Vacc Suit"), skill("Broker"),
				skill("Persuade"), skill("Survival", "Freefall"), skill("Carouse"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Art", "Any"), skill("Diplomat"), skill("Language", "Any"),
				skill("Broker"), skill("Science", "Any"), skill("Advocate", "Any"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 500, Other: chr("DEX", 1)},
			{Cash: 1000, Other: chr("INT", 1)},
			{Cash: 2000, Other: weaponOrItsUse()},
			{Cash: 4000, Other: chr("EDU", 1)},
			{Cash: 10000, Other: stashValued(companyShare, companyShareValue)},
			{Cash: 15000, Other: stashCountValued(2, companyShare, companyShareValue)},
			{Cash: 20000, Other: stashCountValued(3, companyShare, companyShareValue)},
		},
		Mishaps: shipperMishaps(),
		Events:  shipperEvents(),
	}
}

func shipperMishaps() MishapTable {
	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			Summary: "cutbacks eliminate your position",
			Effects: []Effect{
				benefitRolls(-2, 0, ScopeBatch),
				loseStashItem(companyShare),
			},
		},
		{
			Summary: "enough of this; you are going independent",
			Effects: []Effect{
				loseAllBenefits("lose every benefit roll from this career"),
				loseStashItem(companyShare),
				transfer("Independent Merchant", "", 0),
			},
		},
		{Summary: "accidental exposure to a dangerous atmosphere", Effects: []Effect{chr("END", -1)}},
		{
			Summary: "a senior officer decides you should be fired, and management agrees",
			Effects: []Effect{
				loseAllBenefits("lose every benefit roll from this career"),
				relationship(Enemy, 1, ""),
			},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "a rivalry aboard ship, and a fit of anger that ends your service",
			Effects: []Effect{relationship(Rival, 1, "")},
		},
		{
			Summary: "accused of negligence resulting in a crewmate's death",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "a drive accident strands you in a system with only a minor colony",
			Effects: []Effect{pick("what surviving there taught you",
				opt("Survival", skill("Survival", "Any")),
				opt("Animals", skill("Animals", "Any")),
				opt("Suit (Vacc Suit)", skill("Suit", "Vacc Suit")),
				opt("Navigation", skill("Navigation")))},
		},
		{
			Summary: "the company accuses you of theft",
			Effects: []Effect{checkSkill("Advocate", 8,
				[]Effect{benefitRolls(2, 0, ScopeBatch), relationship(Enemy, 1, "")},
				[]Effect{
					relationship(Enemy, 1, ""),
					loseStashItem(companyShare),
					transfer("Prisoner", "Prisoner", 1),
				})},
		},
		{
			Summary: "your ship is destroyed and, with nobody else to blame, the company blames you",
			Effects: []Effect{
				loseAllBenefits("lose every benefit roll from this career"),
				loseStashItem(companyShare),
				debt("2,300,000"),
				throwModifier("next enlistment attempt", -4),
			},
		},
	}
}
