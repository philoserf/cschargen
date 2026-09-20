package chargen

import (
	"fmt"
	"strings"

	"github.com/philoserf/cschargen/career"
)

// skillStepCite is p. 117, the page that governs skill tables and the
// specialty choices they offer. Three call sites reach it.
const skillStepCite = "p. 117"

// applyAll carries out a table result's effects in order.
func (g *Generator) applyAll(effects []career.Effect, cause int) error {
	for _, effect := range effects {
		err := g.apply(effect, cause)
		if err != nil {
			return err
		}
	}

	return nil
}

// apply carries out one effect. The switch is exhaustive over the effect
// vocabulary; a kind added to career/ and not handled here is a build
// failure rather than a silent no-op.
//
// The effect vocabulary is a closed alphabet of nineteen kinds and this is
// the fold over it; splitting it would scatter one table across files, and
// a fold has the size of the alphabet it folds over.
//
// and EffectGroup adds a letter that recurses rather than a branch
//
//nolint:cyclop,funlen,gocyclo // a fold over a closed alphabet has the size of the alphabet,
func (g *Generator) apply(effect career.Effect, cause int) error {
	switch effect.Kind {
	case career.EffectSkill:
		return g.applySkill(effect, cause)
	case career.EffectCharacteristic:
		return g.adjustBy(effect, cause)
	case career.EffectChoice:
		return g.applyChoice(effect)
	case career.EffectCheck:
		return g.applyCheck(effect)
	case career.EffectInjury:
		return g.applyInjury(effect, cause)
	case career.EffectLifeEvent:
		return g.rollLifeEvent(cause)
	case career.EffectMishapNoEject:
		return g.rollMishap(cause, false)
	case career.EffectBenefitRolls:
		return g.grantBenefits(effect, cause)
	case career.EffectEducationLockout:
		g.closeEducation(effect, cause)

		return nil
	case career.EffectRelationship:
		return g.gainTies(effect, cause)
	case career.EffectCredits:
		return g.payCredits(effect, cause)
	case career.EffectStash:
		return g.changeStash(effect, cause)
	case career.EffectGroup:
		return g.applyAll(effect.Group, cause)
	case career.EffectRaiseHeld:
		return g.raiseHeldSkill(effect, cause)
	case career.EffectAnySkill:
		return g.anySkill(effect, cause)
	case career.EffectHomeworldSkill:
		return g.homeworldSpecialty(effect, cause)
	case career.EffectLoseSkill:
		g.loseSkill(effect, cause)

		return nil
	case career.EffectBecome:
		return g.becomeTies(effect, cause)
	case career.EffectImproveTies:
		return g.improveTies(cause)
	case career.EffectModifier:
		g.pending = append(g.pending, PendingModifier{
			Applies: effect.Applies, Value: effect.Modifier, Detail: effect.Detail,
		})
		g.consequence(ConsequenceModifier, cause, effect.Detail, "")

		return nil
	case career.EffectAdvance:
		g.autoAdvance = true
		g.consequence(ConsequenceRank, cause, effect.Detail, "")

		return nil
	case career.EffectRank:
		return g.changeRank(effect, cause)
	case career.EffectContinue:
		g.mustContinue = true
		g.consequence(ConsequenceCareer, cause, effect.Detail, "")

		return nil
	case career.EffectChangeAssignment:
		g.mayChangeAssignment = true
		g.consequence(ConsequenceCareer, cause, effect.Detail, "")

		return nil
	case career.EffectCollegiateLifeEvent, career.EffectAcademyLifeEvent:
		return g.rollInstitutionLifeEvent(cause)
	case career.EffectHonors:
		g.honorsByEvent(cause)

		return nil
	case career.EffectTeenageLifeEvent:
		return g.rollTeenageLifeEvent(cause)
	case career.EffectFreed:
		g.freed(cause)

		return nil
	case career.EffectYouthLifeEvent:
		return g.rollYouthLifeEvent(cause)
	case career.EffectRating:
		return g.moveRatings(effect, cause)
	case career.EffectLoseTie:
		return g.loseTie(effect, cause)
	case career.EffectChooseCareer:
		// ERRATA E-6: an unnamed career change is the ordinary path out of
		// a career, not a transfer. Ejecting is the whole of it -- Step 18
		// (p. 125) then offers the career list, and the character enlists
		// in what they choose on the ordinary terms.
		g.ejected = true

		g.char.Provenance.Deviate("E-6")
		g.consequence(ConsequenceCareer, cause, effect.Detail, "")

		return nil
	case career.EffectTransfer:
		g.transfer = &career.Effect{
			Career: effect.Career, Assignment: effect.Assignment,
			Terms: effect.Terms, Detail: effect.Detail,
		}
		g.consequence(ConsequenceCareer, cause, effect.Detail, "")

		return nil
	case career.EffectNewHomeworld:
		return g.reassignHomeworld(cause, effect.Detail)
	case career.EffectAutoSuccess:
		g.automatic = append(g.automatic, effect.Applies)
		g.consequence(ConsequenceModifier, cause, effect.Detail, "")

		return nil
	case career.EffectMilitaryEvent:
		return g.rollMilitaryEvent(cause)
	case career.EffectRollTable:
		return g.rollNamedTable(effect, cause)
	case career.EffectCommission:
		return g.attemptCommission(effect.Modifier, cause)
	case career.EffectUnimplemented:
		g.unimplemented(cause, effect.Detail)

		return nil
	}

	return fmt.Errorf("%w: %d", ErrUnknownEffect, int(effect.Kind))
}

// applySkill grants a level. Where the book offers a choice of specialty it
// is put to the decider; where it writes "(Any)" it is recorded as "Any",
// which is what the sheet prints and what the player resolves at the table.
func (g *Generator) applySkill(effect career.Effect, cause int) error {
	specialty := ""

	switch len(effect.Specialties) {
	case 0:
	case 1:
		specialty = effect.Specialties[0]
	default:
		chosen, err := g.choose(Choice{
			Point:   "skill_specialty",
			Prompt:  "Choose a specialty for " + effect.Skill,
			Options: effect.Specialties,
			Cite:    skillStepCite,
		})
		if err != nil {
			return err
		}

		specialty = effect.Specialties[chosen]
	}

	// Level zero on an ordinary result means the result did not say, and
	// every ordinary result reads "gain a level in X". AtLevelZero is what
	// a result that really means level 0 sets.
	level := effect.Level
	if level == 0 && !effect.AtLevelZero {
		level = 1
	}

	got := g.char.State.GainSkill(effect.Skill, specialty, level)

	g.log.Consequence(ConsequenceEvent{
		Kind:   ConsequenceSkill,
		Cause:  cause,
		Detail: got.Full() + " " + itoa(got.Level),
		Skill:  got.Full(),
		Level:  got.Level,
		Cite:   g.cite,
	})

	return nil
}

// adjustBy moves a characteristic by a fixed amount, or by one the book
// rolls for: "Lose 1d3 from your choice of STR or END" (p. 156). Delta
// carries the sign and Dice the magnitude.
func (g *Generator) adjustBy(effect career.Effect, cause int) error {
	delta := effect.Delta

	if effect.Dice != "" {
		rolled, err := g.rollExpression(effect.Dice, g.cite)
		if err != nil {
			return err
		}

		if delta < 0 {
			rolled = -rolled
		}

		delta = rolled
	}

	g.adjust(effect.Characteristic, delta, effect.Detail, cause)

	return nil
}

// adjust moves a characteristic, holding it inside the range the rules
// allow: nothing falls below zero, and an unaltered human caps at 15
// (p. 14).
func (g *Generator) adjust(which string, delta int, detail string, cause int) {
	target, ok := characteristicByName(which)
	if !ok {
		g.unimplemented(cause, detail)

		return
	}

	before := g.char.State.Characteristics.Get(target)

	// The ceiling is the species': "For an unaltered human, these
	// characteristics may never rise higher than 15. Some uplifts and
	// engineered humans will have higher maximums" (p. 14).
	after := min(max(before+delta, 0), g.ceiling())

	g.char.State.Characteristics.Set(target, after)

	g.log.Consequence(ConsequenceEvent{
		Kind:           ConsequenceCharacteristic,
		Cause:          cause,
		Detail:         detail,
		Characteristic: which,
		Delta:          after - before,
		Cite:           g.cite,
	})
}

func (g *Generator) applyChoice(effect career.Effect) error {
	labels := make([]string, len(effect.Options))
	for i, option := range effect.Options {
		labels[i] = option.Label
	}

	chosen, err := g.choose(Choice{
		Point:   "table_result",
		Prompt:  effect.Detail,
		Options: labels,
		Cite:    g.cite,
	})
	if err != nil {
		return err
	}

	// g.choose logged the choice; it is the cause of whatever the branch
	// does.
	return g.applyAll(effect.Options[chosen].Effects, g.log.Len())
}

// applyCheck resolves "roll X 8+" and follows the branch it lands in.
//
// A characteristic check is 2d6 plus that characteristic's modifier, which
// p. 110 states for enlistment and the book uses throughout. A *skill*
// check is not stated anywhere in this book: it points at the Core
// Rulebook's task system (p. 13). See ERRATA E-7 for the reading.
func (g *Generator) applyCheck(effect career.Effect) error {
	if effect.Check == nil {
		return ErrMalformedCheck
	}

	throw := g.rollCheck(*effect.Check)

	g.log.Throw(throw, g.cite)

	branch := effect.Failure
	if throw.Success {
		branch = effect.Success
	}

	return g.applyAll(branch, g.log.Len())
}

// unimplemented records a demand the engine could not carry out. It stamps
// no deviation: ERRATA identifiers name readings the engine applied, and
// declining to act on a result is not a reading of it. The record already
// says what was not done.
func (g *Generator) unimplemented(cause int, detail string) {
	g.log.Consequence(ConsequenceEvent{
		Kind:   ConsequenceUnimplemented,
		Cause:  cause,
		Detail: detail,
		Cite:   g.cite,
	})
}

func (g *Generator) consequence(kind ConsequenceKind, cause int, detail, career string) {
	g.log.Consequence(ConsequenceEvent{
		Kind:   kind,
		Cause:  cause,
		Detail: detail,
		Career: career,
		Cite:   g.cite,
	})
}

// characteristicByName maps the abbreviation the tables are written in.
func characteristicByName(name string) (Characteristic, bool) {
	for _, which := range CharacteristicOrder {
		if which.String() == strings.ToUpper(name) {
			return which, true
		}
	}

	return 0, false
}

// itoa avoids strconv in the hot path of building detail strings.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var digits []byte

	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)

		n /= 10
	}

	if negative {
		return "-" + string(digits)
	}

	return string(digits)
}
