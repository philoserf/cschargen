package career

// sportsEvents is the d66 table of pp. 275-277.
//
// Two of its results -- 14 and 56 -- ask a different skill of an Athlete
// than of a Coach, which is the only place in the book where a result's
// check depends on the assignment held. They are written as a choice so that
// the Decider picks the one the character can actually make.
func sportsEvents() EventTable {
	// The pairing those three results share.
	byAssignment := func(success, failureAthlete, failureCoach []Effect) Effect {
		return pick("answer with what your assignment trained",
			opt("Athlete", checkSkill("Athletics", 8, success, failureAthlete)),
			opt("Coach", checkSkill("Tactics", 8, success, failureCoach)))
	}

	table := EventTable{
		11: {Summary: "something terrible happens", Effects: []Effect{mishapNoEject()}},
		12: {
			Summary: "an attack by a teammate severe enough that the attacker is imprisoned",
			Effects: []Effect{relationship(Enemy, 1, "")},
		},
		// Unlike 14 and 56, this row's throw is a free choice between two
		// skills rather than one the assignment dictates.
		13: {
			Summary: "local criminals with an offer to throw an upcoming game",
			Effects: []Effect{pick("report the offer, or take it",
				opt("report it",
					modifierFor(advancementThrow, 2, 2, "+2 to the next two advancement rolls",
						enlistmentNarrowing{}),
					relationship(Enemy, 1, "")),
				opt("take it, on Athletics",
					checkSkill("Athletics", 8, []Effect{credits("100000")}, sportsThrewIt())),
				opt("take it, on Tactics",
					checkSkill("Tactics", 8, []Effect{credits("100000")}, sportsThrewIt())))},
		},
		14: {
			Summary: "targeted by a player on another team, or an opposing coach",
			Effects: []Effect{
				relationship(Rival, 1, ""),
				byAssignment(
					[]Effect{benefitRolls(2, 0, ScopeBatch), chr("CHA", 1)},
					[]Effect{injury(1)},
					[]Effect{benefitRolls(-1, 0, ScopeBatch)}),
			},
		},
		15: {
			Summary: "a lot of fun being involved in professional sports",
			Effects: []Effect{
				skill("Carouse"),
				rollSub("whether it takes hold",
					onRange(1, 2, "it does",
						addiction("alcohol or a drug of the character's choice")),
					onRange(3, 6, "it does not")),
			},
		},
		16: {
			Summary: "downtime taken to study",
			Effects: []Effect{pickSkill("Art", "Science")},
		},
		21: {
			Summary: "a lot of time spent keeping your statistics correct",
			Effects: []Effect{skill("Admin")},
		},
		22: {
			Summary: "questionable people in your retinue",
			Effects: []Effect{pickSkill("Deception", "Streetwise")},
		},
		23: {
			Summary: "a local religion, delved into deeply",
			Effects: []Effect{unimplemented(
				"roll 1d6: on a 6, gain a level in Science (Philosophy); " +
					"where a religion was already taken, this is a change of religion")},
		},
		24: {
			Summary: "a weapons class",
			Effects: []Effect{pick("choose Gun Combat or Melee (Blade)",
				opt("Gun Combat", skill("Gun Combat", "Any")),
				opt("Melee (Blade)", skill("Melee", "Blade")))},
		},
		25: {
			Summary: "a love of operating a vehicle",
			Effects: []Effect{pick("choose the vehicle you took to",
				opt("Drive", skill("Drive", "Wheeled", "Tracked")),
				opt("Flyer (Grav)", skill("Flyer", "Grav")),
				opt("Seafarer", skill("Seafarer", "Sail", "Motorboat")))},
		},
		26: {
			Summary: "staying on top of your agent, your contract and your finances",
			Effects: []Effect{skill("Broker")},
		},
		41: {
			Summary: "the laws and regulations of your league, and of the world you play on",
			Effects: []Effect{pick("choose Admin or Advocate (Legal)",
				opt("Admin", skill("Admin")),
				opt("Advocate (Legal)", skill("Advocate", "Legal")))},
		},
		42: {Summary: "first aid learned", Effects: []Effect{skill("Medic", "First Aid")}},
		43: {
			Summary: "another player making a name by denigrating your accomplishments",
			Effects: []Effect{relationship(Rival, 1, "")},
		},
		44: {
			Summary: "a self-defense class",
			Effects: []Effect{pickSkill("Gun Combat", "Melee")},
		},
		45: {
			Summary: "a nickname earned",
			Effects: []Effect{rollSub("how the public took it",
				on(1, "badly", chr("CHA", -2)),
				onRange(2, 3, "not at all"),
				onRange(4, 5, "well", chr("CHA", 1)),
				on(6, "very well", chr("CHA", 2)))},
		},
		46: {Summary: "a lot of traveling", Effects: []Effect{skill("Suit", "Vacc Suit")}},
		51: {
			Summary: "riding, taken up to relax",
			Effects: []Effect{skill("Animals", "Riding")},
		},
		52: {
			Summary: "a seat on the league's competition committee",
			Effects: []Effect{
				throwModifier("next advancement roll", 2),
				benefitRolls(2, 0, ScopeBatch),
			},
		},
		53: {
			Summary: "your sport outlawed on your world as too dangerous",
			Effects: []Effect{pick("leave the career, or relocate",
				opt("leave the game", benefitRolls(-1, 0, ScopeBatch)),
				opt("relocate", newHomeworld()))},
		},
		54: {
			Summary: "a second language learned",
			Effects: []Effect{skill("Language", "Any")},
		},
		55: {
			Summary: "an endorsement deal, and your likeness in a widely seen advertisement",
			Effects: []Effect{
				credits("25000"),
				rollSub("how the product does",
					on(1, "it embarrasses you", throwModifier(advancementThrow, -2)),
					onRange(2, 3, "it does not sell", benefitRolls(-1, 0, ScopeBatch)),
					onRange(4, 5, "it sells", credits("25000"), benefitRolls(1, 0, ScopeBatch)),
					on(6, "it is everywhere",
						credits("100000"),
						benefitRolls(2, 0, ScopeBatch),
						benefitRolls(0, 1, ScopeCareer))),
			},
		},
		56: {
			Summary: "a victory guaranteed in your next match",
			Effects: []Effect{byAssignment(
				[]Effect{autoSuccess("next survival roll"), throwModifier("next advancement roll", 4)},
				sportsGuaranteeFailed(),
				sportsGuaranteeFailed())},
		},
		61: {
			Summary: "one of the leaders of your team or coaching staff",
			Effects: []Effect{skill("Leadership")},
		},
		62: {
			Summary: "an agent who is a genius",
			Effects: []Effect{benefitRolls(2, 0, ScopeBatch)},
		},
		63: {
			Summary: "something to fall back on: entertainment, or the media",
			Effects: []Effect{pickSkill("Advocate", "Art")},
		},
		// 64 rather than 65: this table's last two rows are its own.
		64: {
			Summary: "extremely talented in your field",
			Effects: []Effect{talentedInYourField()},
		},
		65: {
			Summary: "so popular with the fans that the holovid studios want you",
			Effects: []Effect{pick("take the offer, or let it stand",
				opt("enlist in Celebrity", transfer("Celebrity", "", 0)),
				opt("let it stand", unimplemented(
					"the offer must be accepted within 1d3 terms, and expires "+
						"if a survival roll is failed in that time")))},
		},
		66: {
			Summary: "the most prestigious award in your sport",
			Effects: append(prestigiousAward(), benefitRolls(0, 1, ScopeCareer)),
		},
	}

	lifeEventRows(table)

	return table
}

// sportsThrewIt is what event 13 costs when the thrown game is found out.
func sportsThrewIt() []Effect {
	return []Effect{loseAllBenefits("lose every benefit roll from this career")}
}

// sportsGuaranteeFailed is what event 56 costs when the guarantee does not
// hold. The book gives the same penalty to an Athlete and a Coach.
func sportsGuaranteeFailed() []Effect {
	return []Effect{
		throwModifier("next survival roll", -2),
		throwModifier("next advancement roll", -2),
	}
}
