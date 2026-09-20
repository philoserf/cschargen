package chargen

import (
	"testing"

	"github.com/philoserf/cschargen/setting"
)

// TestARefusedTeenagePathEndsGeneration. Choosing a path goes through the
// Decider, so an abandoned interactive session comes back as an error.
func TestARefusedTeenagePathEndsGeneration(t *testing.T) {
	t.Parallel()

	gen := engine(t, 19)

	gen.decider = refusingDecider{}

	// Everything at 15 opens Paths 3 and 4 on top of whichever of 1 and 2
	// the homeworld gives, so a choice is offered.
	for _, which := range CharacteristicOrder {
		gen.char.State.Characteristics.Set(which, 15)
	}

	_, err := gen.chooseTeenagePath(0, "ages 13-15")
	if err == nil {
		t.Fatal("a refused choice did not come back as an error")
	}

	err = gen.oneTeenageEvent(0, "ages 13-15")
	if err == nil {
		t.Fatal("a refused path did not stop the event")
	}

	err = gen.teenageEvents()
	if err == nil {
		t.Fatal("a refused path did not stop the step")
	}
}

// TestAWorldWithNoSettledYearIsRecentlyColonized. The field is optional in
// the data format, and a file that leaves it out should not silently put
// every character on the path for long-established worlds.
func TestAWorldWithNoSettledYearIsRecentlyColonized(t *testing.T) {
	t.Parallel()

	gen := engine(t, 20)
	if got := gen.settledYears(); got != 0 {
		t.Errorf("a world with no settled year has been settled %d years", got)
	}
}

// refuseAt declines one named choice point and takes the first option
// everywhere else, which is how a test reaches an error that is raised deep
// in a run rather than at its first choice.
type refuseAt struct{ point string }

func (r refuseAt) Choose(c Choice) (int, error) {
	if c.Point == r.point {
		return 0, ErrNoSetting
	}

	return 0, nil
}

func (refuseAt) Kind() DeciderKind { return DeciderPolicy }

// TestARefusalAnywhereInTheLifepathComesBack. Every step returns its
// errors rather than logging them, so a choice refused at Step 6 or Step 7
// ends the run rather than being swallowed by the step above it.
func TestARefusalAnywhereInTheLifepathComesBack(t *testing.T) {
	t.Parallel()

	for _, point := range []string{"youth_path", "teenage_path", "table_result"} {
		t.Run(point, func(t *testing.T) {
			t.Parallel()

			refused := 0

			// Not every seed reaches every choice point, so this scans
			// until one does rather than asserting on a single lifepath.
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
