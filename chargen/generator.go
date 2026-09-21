package chargen

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/dice"
	"github.com/philoserf/cschargen/setting"
)

// The page the characteristic rule is printed on, and the identifier of its
// choice point. Named because both appear in several events and a typo in
// one would read as a second, different rule.
const (
	characteristicCite         = "p. 13"
	pointAssignCharacteristics = "assign_characteristics"
)

// Options configure a generation run. Everything here that a seed does not
// capture is recorded in the character's provenance, because replay needs
// it.
type Options struct {
	Seed          uint64
	Decider       Decider
	EngineVersion string
	PolicyVersion string
	Setting       *setting.Data
	Inputs        Inputs
}

// Generator walks the twenty steps of character creation. One run produces
// one character; a Generator is not reusable.
type Generator struct {
	dice    *dice.Dice
	log     *Log
	decider Decider
	char    *Character

	// The current career, and where the character stands in it.
	career        *career.Career
	assignment    career.Assignment
	rank          int
	commissioned  bool
	termsInCareer int
	cite          string

	// stage is where in the twenty steps the character is, for the ties
	// granted before they have a career to name: "youth", "teenage", or
	// the institution being attended. It is what tieOrigin falls back to.
	stage string

	// forfeitedTerms is how many of this career's terms have had their
	// two-per-term mustering-out grant written off by a result that took
	// every benefit roll back. See forfeitBenefits.
	forfeitedTerms int

	// Flags a table result sets for the loop to read. Each is consumed
	// where it is acted on, so a result that fires twice is two effects
	// rather than a latch nobody cleared.
	autoAdvance         bool
	mustContinue        bool
	mayChangeAssignment bool
	ejected             bool

	// A transfer the rules named -- Colonist's mishap 10 sends a character
	// to Vagabond with the Transient assignment -- waiting for the next
	// career decision.
	transfer         *career.Effect
	forcedAssignment string
	forcedTerms      int

	// Enlistment history. Three consecutive failures force Vagabond
	// (p. 110), and a career that turned the character down is closed to
	// them for two terms.
	failedEnlistments int
	lockout           map[string]int
	forced            string

	// stopped is set when a table result names a career career.ByName does
	// not find. All thirty-four the book names are transcribed, so this is
	// a misspelled destination rather than a missing one. Generation ends
	// cleanly; the record carries the unimplemented consequence naming
	// where they went.
	stopped bool

	// crisisSurvived and mentalDecline are the two standing restrictions
	// pp. 123-124 leave behind: a character who has survived an Aging
	// Crisis automatically fails every future enlistment check, and one
	// with a mental characteristic at 0 "may not attempt further Enlistment
	// checks" at all.
	crisisSurvived bool
	mentalDecline  bool

	pending []PendingModifier

	// automatic holds throws a table result has already decided: "Gain an
	// automatic success on your next Survival roll" (p. 176), "you may
	// enlist automatically in a business, military, corporate or colonist
	// career". They are PendingModifiers for the narrowing rather than for
	// the value, which they do not carry: an automatic success needs to
	// know which careers it reaches for the same reason a modifier does.
	// Each is
	// spent by the throw it names.
	automatic []PendingModifier

	// autoFailure is the mirror: throws a result has already decided
	// against the character. pools are the modifiers the character spends
	// themselves, and onFailure the effects waiting on a throw that fails.
	autoFailure      []string
	pools            []Pool
	onFailure        []career.Effect
	careerBenefitMod int
	termLimit        int

	// Where the character is from, and what that world imposes. techLevel
	// follows the current homeworld; homeworldTerms and maximumAge stay
	// with the world they were born on (ERRATA E-9).
	setting   *setting.Data
	techLevel int

	// settledYear is the homeworld's, which Step 7's first two paths
	// divide on (p. 76).
	settledYear int

	// species is what the character is, and agingProfile and maximum are
	// what that makes of them. A human has no species entry in the data
	// file and takes the defaults.
	species      *setting.Species
	agingProfile AgingProfile
	maximum      int

	// class is an uplift's, which the homeworld's tech level decides
	// between (p. 66). Zero for anyone who is not one.
	class int

	// returning is set while a player's Step 18 answer is being carried
	// out, so that Step 8 knows the return was asked for rather than
	// defaulted into.
	returning bool

	// rollAsHuman is set by a hybrid outcome that sends the character back
	// to the human characteristic method: "Roll characteristics as a
	// human" (p. 63).
	rollAsHuman bool

	// enslaved is set where the homeworld owns this character's people:
	// "the character must take the Altrant/Uplift Slave career as their
	// first career term" (p. 42). It is cleared once they have.
	enslaved bool

	// institution is the one Step 8's events table is being rolled on, so
	// that its result 10 knows which life-events table to reach.
	institution *career.Institution

	// academyClosed is p. 93's standing bar: "If a character fails their
	// success roll while in a military academy, they may not attempt to
	// enter a military academy again."
	academyClosed bool

	// educationClosedUntil is the term count before which no institution
	// will admit the character, which two graduate-track failures impose.
	educationClosedUntil int

	homeworldTerms  int
	maximumAge      int
	primaryLanguage string
}

// New returns a Generator ready to run.
func New(opts Options) *Generator {
	return &Generator{
		dice:      dice.New(opts.Seed),
		log:       &Log{},
		decider:   opts.Decider,
		setting:   opts.Setting,
		forced:    opts.Inputs.Career,
		termLimit: termLimit(opts.Inputs.TermLimit),
		char: &Character{
			Provenance: Provenance{
				SchemaVersion: SchemaVersion,
				Ruleset:       Ruleset,
				EngineVersion: opts.EngineVersion,
				PolicyVersion: opts.PolicyVersion,
				RNG:           RNG{Algorithm: "math/rand/v2 PCG", Seed: opts.Seed},
				SettingData:   describeSetting(opts.Setting),
				Inputs:        opts.Inputs,
			},
		},
	}
}

// defaultTermLimit is how many terms auto mode serves when nothing says
// otherwise.
//
// The rules impose no such limit: a character from a long-settled world may
// serve thirty-five terms, and Earth allows fifty-eight (p. 43). Thirty-five
// terms of Colonist is not a usable NPC, so this is a policy decision rather
// than a rule, and POLICY.md records it as one.
const defaultTermLimit = 4

// describeSetting stamps what data a record was generated against. A nil
// Setting is not dereferenced here: Run refuses it, with an error naming
// the flag, rather than panicking inside a constructor.
func describeSetting(data *setting.Data) SettingData {
	if data == nil {
		return SettingData{}
	}

	return SettingData{Name: data.Name, Hash: data.Hash, Sample: data.Sample}
}

// termLimit reads the requested number of terms. Zero means the flag was
// not given and the policy default applies; a negative number means no
// career terms at all, which is how a caller asks for characteristics
// alone.
func termLimit(requested int) int {
	switch {
	case requested > 0:
		return requested
	case requested < 0:
		return 0
	default:
		return defaultTermLimit
	}
}

// Run walks the steps this milestone implements and returns the record.
//
// Milestone 1 implements Step 2 and Steps 9-18; Steps 1 and 3-8 are stubbed,
// and Step 19 arrives with mustering out (docs/MILESTONE-1.md). The order
// here is the book's, so that adding a step is adding a line rather than
// rearranging one.
func (g *Generator) Run() (*Character, error) {
	if g.setting == nil {
		return nil, ErrNoSetting
	}

	g.char.State.Age = startingAge

	err := g.chooseSpecies()
	if err != nil {
		return nil, err
	}

	// The genetics decide which characteristic method Step 2 uses, so they
	// come before it although the book prints them at Step 5 (p. 62).
	err = g.determineGenetics(g.log.Len())
	if err != nil {
		return nil, err
	}

	err = g.rollCharacteristics()
	if err != nil {
		return nil, err
	}

	err = g.determineOrigin()
	if err != nil {
		return nil, err
	}

	err = g.runCareers()
	if err != nil {
		return nil, err
	}

	err = g.finishingTouches()
	if err != nil {
		return nil, err
	}

	g.char.Events = g.log.Events()

	return g.char, nil
}

// runCareers is the loop of Steps 9 through 18: enter a career, serve
// terms, and decide each time whether to go on.
func (g *Generator) runCareers() error {
	for len(g.char.State.Terms) < g.terms() {
		if g.stopped {
			return nil
		}

		if g.career == nil {
			err := g.enterCareer()
			if err != nil {
				return err
			}

			if g.stopped {
				return nil
			}

			if g.career == nil {
				// The enlistment failed. That consumed no term, so the
				// loop tries again -- three failures in a row is what
				// sends a character to Vagabond, not one.
				continue
			}
		}

		err := g.serveTerm()
		if err != nil {
			return err
		}

		err = g.nextTerm()
		if err != nil {
			return err
		}
	}

	// A character who reaches the term limit is still in a career, and
	// leaving it is what mustering out is for.
	return g.leaveCareer(g.log.Len(), "character generation ended")
}

// terms is how many terms this character may serve: the policy's limit and
// the homeworld's cap are both ceilings, and the lower wins. "If the
// character has reached the maximum number of terms allowed by their
// homeworld, the character must end character generation at this time"
// (p. 125).
func (g *Generator) terms() int {
	if g.homeworldTerms > 0 && g.homeworldTerms < g.termLimit {
		return g.homeworldTerms
	}

	return g.termLimit
}

// nextTerm is Step 18 (p. 125): continue, change assignment, change career,
// or stop. A character ejected by a mishap cannot continue in the career
// they were ejected from (p. 126); one who rolled a natural twelve on
// survival must (p. 113).
func (g *Generator) nextTerm() error {
	step := g.log.Step("Step 18: Determining the Next Term", "p. 125")

	switch {
	case g.ejected:
		return g.leaveCareer(step, "ejected by a mishap")
	case g.transfer != nil:
		return g.leaveCareer(step, "sent to another career by a table result")
	case g.forcedTerms > 0 && g.termsInCareer >= g.forcedTerms:
		return g.leaveCareer(step, "the sentence is served")
	case g.mustContinue:
		g.mustContinue = false

		return nil
	}

	return g.mayReturnToEducation(step)
}

// mayReturnToEducation is one of the five things p. 125 lets a character
// decide between terms: "continue in this career, change to a different
// career, change to a different assignment within the same career, return
// to higher education, or exit character generation."
//
// It is offered to a player and not to the policy. The policy takes the
// first option every time, and the first option here would send every
// character back to school after every term -- which is the policy being
// unrefined rather than the rules being odd, and POLICY.md says so.
func (g *Generator) mayReturnToEducation(step int) error {
	// Only a player is asked. The policy takes the first option every
	// time, and this is a choice the policy cannot make well.
	_, interactive := g.decider.(Asker)
	if !interactive {
		return nil
	}

	score := g.characteristicScore()

	open := false

	for _, institution := range career.Institutions() {
		if g.eligibleFor(institution, score) {
			open = true

			break
		}
	}

	if !open {
		return nil
	}

	chosen, err := g.choose(Choice{
		Point:   "return_to_education",
		Prompt:  "Return to higher education before the next term?",
		Options: []string{"stay in the career", "return to higher education"},
		Cite:    "p. 125",
	})
	if err != nil {
		return err
	}

	if chosen == 0 {
		return nil
	}

	err = g.leaveCareer(step, "returning to higher education")
	if err != nil {
		return err
	}

	g.returning = true
	defer func() { g.returning = false }()

	return g.higherEducation()
}

// leaveCareer ends the current service. Mustering out happens here, and
// only here: "When a player decides that a character will leave their
// current career, whether to enter a new career or to conclude character
// generation, the character must first Muster Out" (p. 126). A character
// who changes career twice musters out twice, and each career's cash cap
// is its own.
func (g *Generator) leaveCareer(cause int, why string) error {
	if g.career == nil {
		return nil
	}

	if service, found := g.char.State.Service(g.career.Name); found {
		service.LeftBecause = why
	}

	g.consequence(ConsequenceCareer, cause, "left the "+g.career.Name+" career: "+why, g.career.Name)

	err := g.musterOut(*g.career, g.termsInCareer)
	if err != nil {
		return err
	}

	g.career = nil
	g.ejected = false
	g.mustContinue = false
	g.mayChangeAssignment = false
	g.forcedTerms = 0
	g.termsInCareer = 0
	g.forfeitedTerms = 0
	g.dropCareerModifiers()
	g.dropCareerPools()

	return nil
}

// dropCareerModifiers forgets the modifiers a result granted "for the rest
// of this career". Nothing else ends them: they carry no count, and a
// modifier that outlived the career it was granted in would follow the
// character for good.
func (g *Generator) dropCareerModifiers() {
	kept := make([]PendingModifier, 0, len(g.pending))

	for _, pending := range g.pending {
		if !pending.WhileInThisCareer {
			kept = append(kept, pending)
		}
	}

	g.pending = kept
}

// choose puts a choice point to the decider and returns a checked index.
func (g *Generator) choose(ask Choice) (int, error) {
	chosen, err := g.decider.Choose(ask)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", ask.Point, err)
	}

	if chosen < 0 || chosen >= len(ask.Options) {
		return 0, ErrChoiceOutOfRange
	}

	g.log.Choice(ChoiceEvent{
		Decider: g.decider.Kind(),
		Point:   ask.Point,
		Prompt:  ask.Prompt,
		Options: ask.Options,
		Chosen:  chosen,
		Cite:    ask.Cite,
	})

	return chosen, nil
}

// rollCharacteristics is Step 2 (p. 39), whose rule is on p. 13: "Roll 3d6.
// Drop the score on the lowest of the three dice and add the remaining two
// scores to determine the character's six characteristics. It is best for
// the player to roll six times and then apply the numbers generated to the
// characteristics as they see fit."
//
// The six rolls happen first and the assignment is a separate choice, in
// that order, because that is the order the sentence gives them -- and
// because a player who assigned as they rolled would be making a different
// decision, with less information, at every step.
func (g *Generator) rollCharacteristics() error {
	step := g.log.Step("Step 2: Roll Characteristics", "p. 39")

	// A species' method names which characteristic it is for -- "roll
	// 2d6-2 for STR and END ... DEX should be rolled as 2d6+2" (p. 23) --
	// so those are rolled in place. Only the ones "rolled normally" go
	// into the pool the player assigns freely.
	free := make([]int, 0, len(CharacteristicOrder))

	for _, which := range CharacteristicOrder {
		expr, fixed := g.speciesMethod(which)
		if !fixed {
			roll := g.dice.Characteristic()
			g.log.Roll(roll, characteristicCite)

			free = append(free, roll.Total)

			continue
		}

		score, err := g.rollExpression(expr, speciesCite)
		if err != nil {
			return err
		}

		g.char.State.Characteristics.Set(which, min(score, g.maximum))

		g.log.Consequence(ConsequenceEvent{
			Kind:           ConsequenceCharacteristic,
			Cause:          g.log.Len(),
			Detail:         which.String() + " " + itoa(score) + ", rolled as " + expr,
			Characteristic: which.String(),
			Delta:          score,
			Cite:           speciesCite,
		})
	}

	err := g.assignCharacteristics(free)
	if err != nil {
		return err
	}

	return g.grantSpeciesSkills(step)
}

// speciesCite is the page a species' characteristic method is printed on,
// which differs per species. The engine cites Step 1 instead, because it
// holds no species and cannot name the page of one.
const speciesCite = "p. 21"

// assignCharacteristics puts each characteristic, in printed order, to the
// decider against the scores still unassigned. Six successive choices
// rather than one permutation: a front end can present them one at a time,
// and Nth/Of tells a player answering the fifth that it is the fifth.
func (g *Generator) assignCharacteristics(rolled []int) error {
	remaining := rolled

	for nth, which := range CharacteristicOrder {
		// A characteristic the species rolled in place is already set.
		if _, fixed := g.speciesMethod(which); fixed {
			continue
		}

		options := make([]string, len(remaining))
		for i, score := range remaining {
			options[i] = strconv.Itoa(score)
		}

		prompt := "Assign a rolled score to " + which.String()

		chosen, err := g.choose(Choice{
			Point:   pointAssignCharacteristics,
			Prompt:  prompt,
			Options: options,
			Cite:    characteristicCite,
			Nth:     nth + 1,
			Of:      len(CharacteristicOrder),
		})
		if err != nil {
			return err
		}

		score := remaining[chosen]

		remaining = slices.Delete(remaining, chosen, chosen+1)

		g.char.State.Characteristics.Set(which, score)

		// The choice event g.choose just logged is the cause: the score
		// landed where it did because the decider put it there.
		cause := g.log.Len()

		g.log.Consequence(ConsequenceEvent{
			Kind:           ConsequenceCharacteristic,
			Cause:          cause,
			Detail:         which.String() + " " + strconv.Itoa(score),
			Characteristic: which.String(),
			Delta:          score,
			Cite:           characteristicCite,
		})
	}

	return nil
}
