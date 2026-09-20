package career

// The enslaved youth and teenage paths (pp. 74, 84).
//
// A character born on a world that owns their people takes the slave career
// as their first (p. 42), and these are the tables they live their
// childhood and adolescence on instead of the five and four the free take.
//
// They are 2d6 tables of eleven rows rather than 2d10 tables of nineteen,
// and both reach the ordinary Youth or Teenage Life Events table at 5-6
// rather than at 10.
//
// The two tables share nine of their eleven rows. The teenage one prints a
// different result 4 and is otherwise the youth table word for word.

// EnslavedPath is one of the two, as a 2d6 table indexed 0-10.
type EnslavedPath struct {
	Name string
	Cite string
	Rows []EventRow
}

// EnslavedYouth is the table of p. 74.
func EnslavedYouth() EnslavedPath {
	return EnslavedPath{
		Name: "Enslaved Youth",
		Cite: "p. 74",
		Rows: enslavedRows(EventRow{
			Summary: "your parents, who are also owned, beaten for a mistake",
			Effects: []Effect{ownerAsEnemy(-175)},
		}, EffectYouthLifeEvent),
	}
}

// EnslavedTeenage is the table of p. 84.
func EnslavedTeenage() EnslavedPath {
	return EnslavedPath{
		Name: "Enslaved Adolescence",
		Cite: "p. 84",
		Rows: enslavedRows(EventRow{
			Summary: "forced to mate with another of your people, for your owners' entertainment",
			Effects: []Effect{relationshipAt(Contact, 1, 60), ownerAsEnemy(-200)},
		}, EffectTeenageLifeEvent),
	}
}

// ownerAsEnemy is the tie three of these results grant, at the rating each
// of them prints.
func ownerAsEnemy(rating int) Effect {
	return relationshipAt(Enemy, 1, rating)
}

// enslavedRows is the nine rows the two tables share, with the one that
// differs passed in and the life-events table each reaches.
func enslavedRows(fourth EventRow, lifeEvent EffectKind) []EventRow {
	return []EventRow{
		{
			// ERRATA E-12 again: "Roll Melee (Unarmed) at Difficult" names
			// a difficulty rather than a target number.
			Summary: "forced to fight another of your people, for your owners' entertainment",
			Effects: []Effect{
				unimplemented("roll Melee (Unarmed) at Difficult: on a success gain the " +
					"defeated as an Enemy, on a failure roll twice on the Injury table"),
				ownerAsEnemy(-175),
			},
		},
		{
			Summary: "beaten by the taskmaster for a minor mistake",
			Effects: []Effect{
				unimplemented("roll 1d6: on 1-3 roll twice on the Injury table, on 4-6 once"),
				ownerAsEnemy(-175),
			},
		},
		fourth,
		{Summary: "a life event", Effects: []Effect{{Kind: lifeEvent, Detail: "roll on the life events table"}}},
		{Summary: "a life event", Effects: []Effect{{Kind: lifeEvent, Detail: "roll on the life events table"}}},
		{
			Summary: "the work is making you strong",
			Effects: []Effect{pick("choose what the work built",
				opt("STR", chr("STR", 1)),
				opt("DEX", chr("DEX", 1)),
				opt("END", chr("END", 1)))},
		},
		{
			Summary: "the arts, as a way of not thinking about the horrors",
			Effects: []Effect{skill("Art", "Any")},
		},
		{
			Summary: "a religion, as a way to search for hope",
			Effects: []Effect{unimplemented("choose a religion")},
		},
		{
			Summary: "sold to another corporation or government",
			Effects: []Effect{
				newHomeworld(),
				unimplemented("take the new world's primary language and background skills"),
			},
		},
		{
			Summary: "an escape attempted",
			Effects: []Effect{unimplemented(
				"roll 2d6 and add the DEX modifier: under 5 you are caught and punished, " +
					"and roll twice on the Injury table; 6-10 you are caught and returned; " +
					"11+ you are free")},
		},
		{
			Summary: "your owner sets you free",
			Effects: []Effect{{Kind: EffectFreed, Detail: "continue as a free character"}},
		},
	}
}
