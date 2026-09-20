package career

// Sports is the career of pp. 273-277: playing a sport for a living, or
// directing from the sidelines.
//
// Three of its mishaps -- cheating, betting on your own games, doping -- say
// the same thing the Gambler's accusation does: it is up to the player
// whether it is true, and the consequence is the same either way.
//
// Its event table is also the one place where the book's usual last two rows
// move: "extremely talented in your field" is at 64 rather than 65, and 65
// and 66 are a standing Celebrity offer and the sport's highest award.
func Sports() Career {
	return Career{
		Name:           "Sports",
		Cite:           "pp. 273-277",
		Enlistment:     &Check{Characteristic: "DEX", Number: 8},
		EnlistmentMods: []EnlistmentMod{apparentAgeOver40(), perPreviousCareer()},
		MishapEjects:   true,
		Assignments: []Assignment{
			{
				Name:        "Athlete",
				Description: "making a living engaging in a sport",
				Survival:    Check{Characteristic: "DEX", Number: 8},
				Advancement: Check{Characteristic: "END", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Athlete", Rows: [6]Effect{
					skill("Carouse"), skill("Discipline"), skill("Medic", "First Aid"),
					skill("Athletics", "Any"), skill("Leadership"), skill("Tactics", "Sport"),
				}},
				Ranks: [][]Effect{
					nil,
					{skill("Discipline")},
					{skill("Tactics", "Sport")},
					nil,
					{skill("Broker")},
					nil,
					{skill("Leadership")},
				},
			},
			{
				Name:        "Coach",
				Description: "directing athletes from the sidelines or the spotter's booth",
				Survival:    Check{Characteristic: "INT", Number: 8},
				Advancement: Check{Characteristic: "EDU", Number: 8},
				Skills: SkillTable{Kind: AssignmentSkills, Name: "Coach", Rows: [6]Effect{
					skill("Admin"), skill("Advocate", "Oratory"), skill("Tactics", "Sport"),
					skill("Instruction"), skill("Broker"), skill("Persuade"),
				}},
				Ranks: [][]Effect{
					{skill("Persuade")},
					{skill("Tactics", "Sport")},
					nil,
					{skill("Broker")},
					{skill("Instruction")},
					nil,
					{skill("Advocate", "Any")},
				},
			},
		},
		Tables: []SkillTable{
			{Kind: PersonalDevelopment, Name: "Personal Development", Rows: [6]Effect{
				chr("STR", 1), chr("DEX", 1), chr("END", 1),
				chr("INT", 1), chr("EDU", 1), chr("CHA", 1),
			}},
			{Kind: ServiceSkills, Name: "Service Skills", Rows: [6]Effect{
				skill("Melee", "Any"), skill("Discipline"), skill("Carouse"),
				skill("Athletics", "Any"), skill("Tactics", "Sport"), skill("Leadership"),
			}},
			{Kind: AdvancedEducation, Name: "Advanced Education", MinimumEDU: 8, Rows: [6]Effect{
				skill("Admin"), skill("Etiquette"), skill("Broker"),
				skill("Tactics", "Sport"), skill("Art", "Any"), skill("Advocate", "Any"),
			}},
		},
		Benefits: [7]BenefitRow{
			{Cash: 1000, Other: chr("DEX", 1)},
			{Cash: 2500, Other: chr("END", 1)},
			{Cash: 5000, Other: chr("CHA", 1)},
			{Cash: 10000, Other: relationship(Contact, 1, "")},
			{Cash: 20000, Other: relationship(Ally, 1, "")},
			{Cash: 50000, Other: credits("2d6x10000")},
			{Cash: 100000, Other: stashValued(companyShare, companyShareValue)},
		},
		Mishaps: sportsMishaps(),
		Events:  sportsEvents(),
	}
}

// sportsSoftCareers is the penalty three of this career's mishaps leave
// behind: a reputation that follows the character into anything they have to
// talk or study their way into.
func sportsSoftCareers(forever bool) Effect {
	scope := "next career"
	if forever {
		scope = "every career after this one"
	}

	return enlistmentPenalty(-2, "-2 to enter "+scope+" where enlistment is based on EDU or CHA",
		enlistmentNarrowing{OnCharacteristics: []string{"EDU", "CHA"}, Standing: forever})
}

func sportsMishaps() MishapTable {
	loseAll := loseAllBenefits("lose every benefit roll from this career")

	return MishapTable{
		{Summary: "severely injured", Effects: []Effect{injury(2)}},
		{
			// Like Gambler's mishap 9, the book leaves the truth of it to
			// the player and applies the damage regardless.
			Summary: "accused of cheating, accurately or not",
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{
			Summary: "a scandal no club will touch you after",
			Effects: []Effect{loseAll},
		},
		{Summary: "tired of the game, and it is showing on the field", Effects: nil},
		{
			Summary: `defeat once too often, and a career "washed up"`,
			Effects: []Effect{benefitRolls(-2, 0, ScopeBatch)},
		},
		{Summary: "injured", Effects: []Effect{injury(1)}},
		{
			Summary: "believed to have bet on games you played in",
			Effects: []Effect{loseAll, skill("Gambler")},
		},
		{
			Summary: "something extremely offensive, said publicly after a watched game",
			Effects: []Effect{sportsSoftCareers(false)},
		},
		{
			Summary: "believed to have used performance enhancers to win",
			Effects: []Effect{loseAll},
		},
		{
			Summary: "an embarrassment that becomes a catch phrase",
			Effects: []Effect{sportsSoftCareers(true)},
		},
		{
			Summary: "a serious disease that does irrevocable damage",
			Effects: []Effect{pick("choose what the disease took",
				opt("STR", chrRolled("STR", "1d3", false)),
				opt("END", chrRolled("END", "1d3", false)))},
		},
	}
}
