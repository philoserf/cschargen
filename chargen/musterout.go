package chargen

import "github.com/philoserf/cschargen/career"

// cashRollsPerCareer is p. 127's cap: "A character may roll on the Cash
// table no more than three times per career, regardless of how many total
// Mustering Out rolls they possess."
const cashRollsPerCareer = 3

// benefitRowsPerCareer is how many rows a benefit table prints. The roll is
// 1d6; the seventh row is reachable only through an event's +1 (ERRATA
// E-3), and a modified roll above the table is clamped to it.
const benefitRowsPerCareer = 7

// rollsPerTerm is p. 126's rule: "The character receives a number of
// Mustering Out Benefit rolls equal to twice the number of terms served in
// that career."
const rollsPerTerm = 2

// musterOut is Step 19 (pp. 126-129). It is called when a character leaves
// a career, whether to enter another or to stop.
func (g *Generator) musterOut(left career.Career, terms int) error {
	step := g.log.Step("Step 19: Muster Out", "p. 126")

	batches := g.benefitsFor(left.Name, terms)

	total := 0
	for _, batch := range batches {
		total += batch.Rolls
	}

	if total <= 0 {
		// The queue is emptied even with nothing to spend. A mishap that
		// removed two rolls leaves a negative batch behind, and leaving it
		// there would charge the character for that mishap a second time
		// the next time they joined this career.
		g.clearBenefits(left.Name)

		g.consequence(ConsequenceBenefitRolls, step,
			"no mustering out rolls remain for "+left.Name, left.Name)

		return nil
	}

	cashTaken := 0

	for _, batch := range batches {
		for range batch.Rolls {
			taken, err := g.benefitRoll(left, batch, cashTaken)
			if err != nil {
				return err
			}

			if taken {
				cashTaken++
			}
		}
	}

	g.clearBenefits(left.Name)

	return nil
}

// benefitsFor assembles the rolls a career is owed: two per term served,
// plus or minus whatever its events and mishaps granted or removed.
//
// The two-per-term grant is its own batch, carrying the career-wide
// modifier events attached to it, and the event batches keep their own --
// which is the distinction ERRATA E-3 rests on.
func (g *Generator) benefitsFor(name string, terms int) []BenefitBatch {
	earned := max(terms-g.forfeitedTerms, 0)

	batches := []BenefitBatch{{
		Career:   name,
		Rolls:    earned * rollsPerTerm,
		Modifier: g.careerBenefitMod,
	}}

	compelled := make([]BenefitBatch, 0, len(g.char.State.Benefits))

	for _, batch := range g.char.State.Benefits {
		if batch.Career != name {
			continue
		}

		if batch.Rolls < 0 {
			// A mishap that removed rolls comes out of the term grant,
			// because that is the only pool a character has when it fires.
			batches[0].Rolls += batch.Rolls

			continue
		}

		if batch.CashOnly {
			// A batch the character may only spend on Cash is settled
			// first, so that the free rolls are chosen knowing what is
			// left of p. 127's three. Taking them the other way round
			// would make the cap unreachable from a compelled roll: every
			// free choice would already have been made.
			compelled = append(compelled, batch)

			continue
		}

		batches = append(batches, batch)
	}

	if batches[0].Rolls < 0 {
		batches[0].Rolls = 0
	}

	return append(compelled, batches...)
}

func (g *Generator) clearBenefits(name string) {
	var kept []BenefitBatch

	for _, batch := range g.char.State.Benefits {
		if batch.Career != name {
			kept = append(kept, batch)
		}
	}

	g.char.State.Benefits = kept
}

// benefitRoll is one mustering-out roll: choose the Cash table or the Other
// Benefits table, then roll 1d6 (p. 127). It reports whether the Cash table
// was the one taken, so the caller can hold the three-per-career cap.
func (g *Generator) benefitRoll(left career.Career, batch BenefitBatch, cashTaken int) (bool, error) {
	takingCash, err := g.chooseBenefitTable(left, batch, cashTaken)
	if err != nil {
		return false, err
	}

	modifier := batch.Modifier

	// A benefit roll has no target: it reads a row rather than passing or
	// failing, so it is logged as a plain roll and the modifier is named in
	// the consequence beside it.
	roll := g.dice.D6()
	cause := g.log.Roll(roll, "p. 127")

	modified := roll.Total + modifier
	if modifier != 0 {
		g.char.Provenance.Deviate("E-3")
		g.consequence(ConsequenceBenefitRolls, cause,
			"a benefit roll of "+itoa(roll.Total)+" at "+itoa(modifier)+" reads row "+
				itoa(min(modified, benefitRowsPerCareer)), left.Name)
	}

	row := min(max(modified, 1), benefitRowsPerCareer)

	benefit := left.Benefits[row-1]

	if takingCash {
		g.char.State.Credits += benefit.Cash
		g.consequence(ConsequenceCredits, cause,
			"mustering out of "+left.Name+": "+itoa(benefit.Cash)+" credits", left.Name)

		return true, nil
	}

	if benefit.Other.Detail == "" || benefit.Other.Detail == "nothing" {
		g.consequence(ConsequenceBenefitRolls, cause,
			"mustering out of "+left.Name+": nothing", left.Name)

		return false, nil
	}

	return false, g.apply(benefit.Other, cause)
}

// chooseBenefitTable is the Cash-or-Other decision of p. 127, and the two
// places it is not a decision.
//
// A batch an event granted as Cash rolls is spent on Cash whatever the
// character would prefer, and it is taken even past the three-per-career
// cap: ERRATA E-32 reads that cap as governing the choice rather than the
// roll, because a result that compels a Cash roll has already made the
// choice. It still counts, so a compelled roll can close the cap against a
// free one later.
func (g *Generator) chooseBenefitTable(
	left career.Career, batch BenefitBatch, cashTaken int,
) (bool, error) {
	if batch.CashOnly {
		if cashTaken >= cashRollsPerCareer {
			g.char.Provenance.Deviate("E-32")
		}

		return true, nil
	}

	cashOpen := cashTaken < cashRollsPerCareer

	options := []string{"Other Benefits"}
	if cashOpen {
		options = []string{"Cash", "Other Benefits"}
	}

	index, err := g.choose(Choice{
		Point:   "benefit_table",
		Prompt:  "Choose a mustering out table for " + left.Name,
		Options: options,
		Cite:    "p. 127",
	})
	if err != nil {
		return false, err
	}

	return cashOpen && index == 0, nil
}

// immediateCashRolls resolves the Cash rolls a result says are taken now
// rather than at Step 19: "excellent work earns a bonus: take two Cash
// Benefit rolls immediately". They read the current career's Cash column
// and pay at once, so they are outside the Step 19 queue and outside the
// three-per-career cap that governs it.
func (g *Generator) immediateCashRolls(effect career.Effect, cause int) error {
	if g.career == nil {
		g.unimplemented(cause, effect.Detail+" -- outside a career, with no cash table to read")

		return nil
	}

	for range effect.Count {
		amount, rollCause := g.oneImmediateCashRoll(effect)

		g.char.State.Credits += amount
		g.consequence(ConsequenceCredits, rollCause,
			"an immediate cash benefit roll in "+g.career.Name+": "+itoa(amount)+" credits",
			g.career.Name)
	}

	return nil
}

// oneImmediateCashRoll reads one row of the current career's Cash column,
// carrying the career-wide modifier of ERRATA E-3 as a Step 19 roll would.
//
// "Re-rolling any result of nothing" means any: Orbital Construction's
// first row pays nothing, so one re-roll can land on it again. The loop is
// bounded by the table's own length so that a career paying nothing on
// every row terminates with nothing rather than spinning.
func (g *Generator) oneImmediateCashRoll(effect career.Effect) (int, int) {
	var (
		amount int
		cause  int
	)

	for range benefitRowsPerCareer {
		roll := g.dice.D6()
		row := min(max(roll.Total+g.careerBenefitMod, 1), benefitRowsPerCareer)

		cause = g.log.Roll(roll, "p. 127")
		amount = g.career.Benefits[row-1].Cash

		if amount != 0 || !effect.RerollNothing {
			break
		}
	}

	if g.careerBenefitMod != 0 {
		g.char.Provenance.Deviate("E-3")
	}

	return amount, cause
}
