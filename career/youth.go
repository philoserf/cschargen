package career

// Step 6: Youth Events (pp. 67-75).
//
// Five paths of nineteen rows each, gated on characteristics, rolled twice:
// once for ages 4-8 and once for 9-12. "If something happens which improves
// or decreases the characteristic scores of the character which allows them
// to change their path, the player may choose to do so for the second roll"
// (p. 68), so the gate is checked again between the rolls.
//
// The tables are 2d10 rather than 2d6 or d66, which is why the d10 landed
// with the family step. Result 10 sends the character to the Youth Life
// Events table on every path.
//
// The uplift and enslaved variants of this step (pp. 68, 74) belong to the
// species work, and two of their headings are proper names the OGL notice
// reserves.

// YouthPath is one of the five columns of pp. 68-74: what a character has
// to be to take it, and what nineteen things can happen to them on it.
type YouthPath struct {
	Name string

	// Cite is the page the path is printed on.
	Cite string

	// Requires is the characteristics that open the path, any one of which
	// is enough: Path 2 is "STR 8+ or END 8+". Nil is Path 1, which has no
	// requirements and is always open.
	Requires []Check

	// Rows is the 2d10 table, indexed 0-18 for a result of 2-20.
	Rows []EventRow
}

// Open reports whether a character with these characteristics may take this
// path. "Simply check the requirements for each path, determine which path
// for which your character qualifies, and then roll on that chart" (p. 68).
func (p YouthPath) Open(score func(string) int) bool {
	if len(p.Requires) == 0 {
		return true
	}

	for _, check := range p.Requires {
		if score(check.Characteristic) >= check.Number {
			return true
		}
	}

	return false
}

// YouthPaths is the five, in the order the book prints them -- which is
// also the order the auto policy walks, so it takes Path 1 unless something
// narrows the choice.
func YouthPaths() []YouthPath {
	return []YouthPath{
		{Name: "Path 1", Cite: "pp. 68-69", Rows: youthPathOne()},
		{
			Name: "Path 2", Cite: "pp. 70-71",
			Requires: []Check{
				{Characteristic: "STR", Number: 8},
				{Characteristic: "END", Number: 8},
			},
			Rows: youthPathTwo(),
		},
		{
			Name: "Path 3", Cite: "pp. 71-72",
			Requires: []Check{{Characteristic: "DEX", Number: 8}},
			Rows:     youthPathThree(),
		},
		{
			Name: "Path 4", Cite: "pp. 72-73",
			Requires: []Check{{Characteristic: "INT", Number: 8}},
			Rows:     youthPathFour(),
		},
		{
			Name: "Path 5", Cite: "pp. 73-74",
			Requires: []Check{{Characteristic: "CHA", Number: 8}},
			Rows:     youthPathFive(),
		},
	}
}

// orphaned is result 2 on every one of the five paths, printed identically:
// "Your parents have been killed in a serious accident ... Lose both
// parents and 1d3 from EDU as the orphanage does not provide the best
// education."
func orphaned() EventRow {
	return EventRow{
		Summary: "orphaned, and sent to a care home",
		Effects: []Effect{
			loseTie(), loseTie(),
			chrRolled("EDU", "1d3", false),
		},
	}
}

// youthLifeEvent is result 10 on every one of the five paths.
func youthLifeEvent() EventRow {
	return EventRow{
		Summary: "a youth life event",
		Effects: []Effect{{Kind: EffectYouthLifeEvent, Detail: "roll on the Youth Life Events table"}},
	}
}

// mentorAt125 is the Ally four of the five paths grant by the same words:
// "You have gained a mentor. Gain an older person in your life as an Ally
// with a Relationship Rating of 125."
func mentorAt125() Effect {
	return relationshipAt(Ally, 1, 125)
}

// YouthLifeEvents is the d6 table of p. 75, indexed 0-5 for a result of
// 1-6. Every path reaches it at result 10.
func YouthLifeEvents() []EventRow {
	return []EventRow{
		{
			Summary: "a childhood illness or physical trauma",
			Effects: []Effect{injury(1)},
		},
		{
			Summary: "someone close to you has died",
			Effects: []Effect{loseTie(Contact, Ally)},
		},
		{Summary: "a new friend", Effects: []Effect{relationship(Contact, 1, "")}},
		{Summary: "a new friend", Effects: []Effect{relationship(Contact, 1, "")}},
		{
			Summary: "an older person takes an interest in your life",
			Effects: []Effect{relationship(Ally, 1, "")},
		},
		{
			Summary: "something wonderful, and the confidence it brings",
			Effects: []Effect{
				chr("CHA", 2),
				unimplemented("if the character was orphaned by another event, they are adopted"),
			},
		},
	}
}
