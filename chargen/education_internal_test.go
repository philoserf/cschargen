package chargen

import (
	"strings"
	"testing"

	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/setting"
)

// educationEngine is a generator with a homeworld and characteristics good
// enough for both institutions, which is what Step 8 needs to have anything
// to offer.
func educationEngine(t *testing.T, seed uint64) *Generator {
	t.Helper()

	gen := New(Options{
		Seed:          seed,
		Decider:       Policy{},
		EngineVersion: testVersion,
		PolicyVersion: testVersion,
		Setting: settingWith(setting.Subsector{
			Name: "New Holdings", OriginRoll: 0,
			Worlds: []setting.World{testWorld("Tinderfall", nil), testWorld("Wake", nil)},
		}),
		Inputs: Inputs{Species: testHuman, TermLimit: -1},
	})

	for _, which := range CharacteristicOrder {
		gen.char.State.Characteristics.Set(which, 12)
	}

	return gen
}

// TestAFailedAcademyClosesTheAcademy is p. 93: "If a character fails their
// success roll while in a military academy, they may not attempt to enter a
// military academy again."
func TestAFailedAcademyClosesTheAcademy(t *testing.T) {
	t.Parallel()

	gen := educationEngine(t, 21)

	gen.academyClosed = true

	institution, attending, err := gen.chooseInstitution(0)
	if err != nil {
		t.Fatalf("chooseInstitution: %v", err)
	}

	if attending && institution.Name == career.MilitaryAcademy().Name {
		t.Error("a character barred from the academy was offered it")
	}
}

// TestNobodyQualifiesForAnything. "The character's EDU must be 6 or higher
// to be accepted into Undergraduate College" (p. 86), and the academy wants
// END 8+ on top.
func TestNobodyQualifiesForAnything(t *testing.T) {
	t.Parallel()

	gen := educationEngine(t, 22)

	for _, which := range CharacteristicOrder {
		gen.char.State.Characteristics.Set(which, 2)
	}

	_, attending, err := gen.chooseInstitution(0)
	if err != nil {
		t.Fatalf("chooseInstitution: %v", err)
	}

	if attending {
		t.Error("a character below every prerequisite was offered higher education")
	}
}

// TestDecliningHigherEducation is the last option the choice offers, which
// the policy never takes: "characters are not required to attend college"
// (p. 85).
func TestDecliningHigherEducation(t *testing.T) {
	t.Parallel()

	gen := educationEngine(t, 23)

	gen.decider = lastOptionDecider{}

	_, attending, err := gen.chooseInstitution(0)
	if err != nil {
		t.Fatalf("chooseInstitution: %v", err)
	}

	if attending {
		t.Error("a character who declined was enrolled anyway")
	}
}

// TestALifeEventWithNoInstitution. The two life-event effects and the
// honours effect all read the institution being rolled on, and a record
// that reached them another way should do nothing rather than panic.
func TestALifeEventWithNoInstitution(t *testing.T) {
	t.Parallel()

	gen := educationEngine(t, 24)

	err := gen.rollInstitutionLifeEvent(0)
	if err != nil {
		t.Fatalf("rollInstitutionLifeEvent with no institution: %v", err)
	}

	gen.honorsByEvent(0)

	if len(gen.char.State.Education) != 0 {
		t.Error("honours were granted to a character who never enrolled")
	}
}

// TestRaisingSomethingThatIsNotACharacteristic. A degree raises EDU to ten
// by name, and a name this ruleset does not have does nothing.
func TestRaisingSomethingThatIsNotACharacteristic(t *testing.T) {
	t.Parallel()

	gen := educationEngine(t, 25)
	before := gen.char.State.Characteristics

	gen.raiseToAtLeast("SOC", medicalEDUFirst, 0)

	if gen.char.State.Characteristics != before {
		t.Error("raising SOC moved something")
	}
}

// TestARefusedEducationChoiceEndsGeneration. Step 8 makes three choices --
// whether to attend, what the degree is in, and what a washout leaves
// behind -- and each of them goes through the Decider.
func TestARefusedEducationChoiceEndsGeneration(t *testing.T) {
	t.Parallel()

	for _, point := range []string{"higher_education", "degree_field", "washout_skill"} {
		t.Run(point, func(t *testing.T) {
			t.Parallel()

			refused := 0

			for seed := range uint64(40) {
				gen := New(Options{
					Seed:          seed,
					Decider:       refuseAt{point: point},
					EngineVersion: testVersion,
					PolicyVersion: testVersion,
					Setting: settingWith(setting.Subsector{
						Name: "New Holdings", OriginRoll: 0,
						Worlds: []setting.World{
							testWorld("Tinderfall", nil),
							testWorld("Wake", nil),
						},
					}),
					Inputs: Inputs{Species: testHuman, TermLimit: -1},
				})

				_, err := gen.Run()
				if err != nil {
					refused++
				}
			}

			if refused == 0 {
				t.Errorf("no seed in 40 reached the %s choice point", point)
			}
		})
	}
}

// academyDecider takes the Military Academy where it is offered and the
// first option everywhere else. The auto policy never reaches the academy,
// because Undergraduate College is printed first and is offered first.
type academyDecider struct{}

func (academyDecider) Choose(c Choice) (int, error) {
	if c.Point == "higher_education" {
		for i, option := range c.Options {
			if option == career.MilitaryAcademy().Name {
				return i, nil
			}
		}
	}

	return 0, nil
}

func (academyDecider) Kind() DeciderKind { return DeciderPolicy }

// TestTheAcademyIsAttemptedAndItsConsequencesRecorded. A graduate owes the
// service two terms as an officer (p. 94) and a washout may never try again
// (p. 93), and neither is reachable under the auto policy.
func TestTheAcademyIsAttemptedAndItsConsequencesRecorded(t *testing.T) {
	t.Parallel()

	var graduated, washedOut int

	for seed := range uint64(80) {
		gen := New(Options{
			Seed:          seed,
			Decider:       academyDecider{},
			EngineVersion: testVersion,
			PolicyVersion: testVersion,
			Setting: settingWith(setting.Subsector{
				Name: "New Holdings", OriginRoll: 0,
				Worlds: []setting.World{testWorld("Tinderfall", nil), testWorld("Wake", nil)},
			}),
			Inputs: Inputs{Species: testHuman, TermLimit: -1},
		})

		character, err := gen.Run()
		if err != nil {
			t.Fatalf("seed %d: %v", seed, err)
		}

		record := academyRecord(character)
		if record == nil {
			continue
		}

		switch {
		case record.Succeeded:
			graduated++

			checkAcademyObligation(t, seed, character)
		case record.Admitted:
			washedOut++

			if !gen.academyClosed {
				t.Errorf("seed %d washed out and the academy is still open to them", seed)
			}
		}
	}

	if graduated == 0 || washedOut == 0 {
		t.Errorf("in 80 seeds, %d graduated the academy and %d washed out; both should happen",
			graduated, washedOut)
	}
}

// checkAcademyObligation: a graduate owes the service two terms as an
// officer (p. 94). The engine records that rather than holding them to it,
// because Step 9 has no way to enter a career at a rank.
func checkAcademyObligation(t *testing.T, seed uint64, character *Character) {
	t.Helper()

	if !mentions(character, "must enter a military career") {
		t.Errorf("seed %d graduated and the record does not carry the obligation", seed)
	}
}

// academyRecord finds the character's attempt at the Military Academy,
// where they made one.
func academyRecord(character *Character) *Education {
	for _, held := range character.State.Education {
		if held.Institution == career.MilitaryAcademy().Name {
			return held
		}
	}

	return nil
}

// mentions reports whether any consequence in a record contains a phrase.
func mentions(character *Character, phrase string) bool {
	for _, event := range character.Events {
		if event.Consequence != nil && strings.Contains(event.Consequence.Detail, phrase) {
			return true
		}
	}

	return false
}
