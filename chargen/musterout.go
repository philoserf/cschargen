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
			taken, err := g.benefitRoll(left, batch.Modifier, cashTaken)
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
	batches := []BenefitBatch{{
		Career:   name,
		Rolls:    terms * rollsPerTerm,
		Modifier: g.careerBenefitMod,
	}}

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

		batches = append(batches, batch)
	}

	if batches[0].Rolls < 0 {
		batches[0].Rolls = 0
	}

	return batches
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
func (g *Generator) benefitRoll(left career.Career, modifier, cashTaken int) (bool, error) {
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

	takingCash := cashOpen && index == 0

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
