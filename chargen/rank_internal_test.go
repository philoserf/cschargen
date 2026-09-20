package chargen

import (
	"errors"
	"strings"
	"testing"

	"github.com/philoserf/cschargen/career"
	"github.com/philoserf/cschargen/setting"
)

// Rank benefits are a floor, not an increment (p. 116):
//
//	If a Benefit grants a skill the character already possesses at the
//	listed level or higher, no additional benefit is gained. Rank Benefits
//	represent required competence, not bonus stacking.
//
// The engine added instead, in every career, from milestone 1 until this
// test existed.

func TestARankBenefitIsAFloor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		held  int
		grant career.Effect
		want  int
	}{
		{"nothing held", 0, career.RaiseSkill("Recon", 1), 1},
		{"held below the line", 0, career.RaiseSkill("Recon", 2), 2},
		{"held at the line", 1, career.RaiseSkill("Recon", 1), 1},
		{"held above the line", 3, career.RaiseSkill("Recon", 1), 3},
		{"held between two ranks", 1, career.RaiseSkill("Recon", 2), 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gen := engine(t, 28)

			if tc.held > 0 {
				gen.char.State.GainSkill("Recon", "", tc.held)
			}

			err := gen.applyRankBenefits([]career.Effect{tc.grant}, 0)
			if err != nil {
				t.Fatalf("applyRankBenefits: %v", err)
			}

			got, found := gen.char.State.Skill("Recon", "")
			if !found {
				t.Fatal("Recon was not granted at all")
			}

			if got.Level != tc.want {
				t.Errorf("Recon = %d, want %d", got.Level, tc.want)
			}
		})
	}
}

// TestARankBenefitThatIsNotASkill. p. 116's rule is about skills. A rank
// row can also grant a characteristic, a Contact or a stash, and those are
// applied as printed -- Scavenger's Junker rank 2 is "+1 END" (p. 266), and
// a character who already has END does not stop gaining it.
func TestARankBenefitThatIsNotASkill(t *testing.T) {
	t.Parallel()

	gen := engine(t, 29)

	gen.char.State.Characteristics.Set(END, 7)

	err := gen.applyRankBenefits([]career.Effect{
		{Kind: career.EffectCharacteristic, Detail: "+1 END", Characteristic: "END", Delta: 1},
	}, 0)
	if err != nil {
		t.Fatalf("applyRankBenefits: %v", err)
	}

	if got := gen.char.State.Characteristics.Get(END); got != 8 {
		t.Errorf("END = %d, want 8: p. 116's rule is about skills", got)
	}
}

// TestAnAnySpecialtyRankBenefitPicksOneNotHeld is p. 116's other half:
//
//	This allows the character to select a specialty within that skill,
//	provided they do not already possess that specialty at level 1 or
//	higher.
func TestAnAnySpecialtyRankBenefitPicksOneNotHeld(t *testing.T) {
	t.Parallel()

	gen := engine(t, 30)

	gen.char.State.GainSkill("Advocate", "Legal", 1)

	granted := career.RaiseSkill("Advocate", 1)

	granted.Specialties = []string{"Legal", "Politics", "Oratory"}

	err := gen.applyRankBenefits([]career.Effect{granted}, 0)
	if err != nil {
		t.Fatalf("applyRankBenefits: %v", err)
	}

	// The policy takes the first option, which is now Politics rather than
	// the Legal the character already holds.
	got, found := gen.char.State.Skill("Advocate", "Politics")
	if !found || got.Level != 1 {
		t.Error("the rank benefit did not pick a specialty the character lacked")
	}

	if held, _ := gen.char.State.Skill("Advocate", "Legal"); held.Level != 1 {
		t.Errorf("Advocate (Legal) is at %d; the benefit should not have touched it", held.Level)
	}
}

// TestAnAnySpecialtyRankBenefitWithNothingLeft is the last sentence of
// p. 116: "If the character already possesses all available specialties at
// level 1 or higher, they gain no additional benefit from that Rank
// Benefit."
func TestAnAnySpecialtyRankBenefitWithNothingLeft(t *testing.T) {
	t.Parallel()

	gen := engine(t, 31)

	for _, specialty := range []string{"Legal", "Politics"} {
		gen.char.State.GainSkill("Advocate", specialty, 1)
	}

	granted := career.RaiseSkill("Advocate", 1)

	granted.Specialties = []string{"Legal", "Politics"}

	err := gen.applyRankBenefits([]career.Effect{granted}, 0)
	if err != nil {
		t.Fatalf("applyRankBenefits: %v", err)
	}

	for _, specialty := range []string{"Legal", "Politics"} {
		held, _ := gen.char.State.Skill("Advocate", specialty)
		if held.Level != 1 {
			t.Errorf("Advocate (%s) is at %d, want 1", specialty, held.Level)
		}
	}
}

// TestARefusedRankSpecialtyEndsGeneration.
func TestARefusedRankSpecialtyEndsGeneration(t *testing.T) {
	t.Parallel()

	gen := engine(t, 32)

	gen.decider = refusingDecider{}

	granted := career.RaiseSkill("Advocate", 1)

	granted.Specialties = []string{"Legal", "Politics"}

	err := gen.applyRankBenefits([]career.Effect{granted}, 0)
	if err == nil {
		t.Fatal("a refused specialty did not come back as an error")
	}
}

// TestTheFloorIsReachedInAGeneratedRecord. A unit test proves the rule; a
// generated character proves the rule is reached. A career whose rank
// benefits repeat a skill the character picked up from its own skill
// tables is common, and before this the second grant stacked.
func TestTheFloorIsReachedInAGeneratedRecord(t *testing.T) {
	t.Parallel()

	met := 0

	for seed := range uint64(80) {
		gen := New(Options{
			Seed:          seed,
			Decider:       Policy{},
			EngineVersion: testVersion,
			PolicyVersion: testVersion,
			Setting: settingWith(setting.Subsector{
				Name: "New Holdings", OriginRoll: 0,
				Worlds: []setting.World{testWorld("Tinderfall", nil), testWorld("Wake", nil)},
			}),
			Inputs: Inputs{Species: testHuman, TermLimit: 12},
		})

		character, err := gen.Run()
		if err != nil {
			t.Fatalf("seed %d: %v", seed, err)
		}

		for _, event := range character.Events {
			if event.Consequence == nil {
				continue
			}

			if strings.Contains(event.Consequence.Detail, "which meets the required") {
				met++
			}
		}
	}

	if met == 0 {
		t.Error("in 80 seeds no rank benefit was ever met by a skill already held")
	}
}

// TestARankBenefitWithNoLevelIsOneLevel. Every rank row in the book is
// built by career's skill constructor, which sets a level. A row that did
// not is read as the ordinary "gain a level in X".
func TestARankBenefitWithNoLevelIsOneLevel(t *testing.T) {
	t.Parallel()

	gen := engine(t, 33)

	err := gen.applyRankBenefits([]career.Effect{
		{Kind: career.EffectSkill, Detail: "Recon", Skill: "Recon"},
	}, 0)
	if err != nil {
		t.Fatalf("applyRankBenefits: %v", err)
	}

	got, found := gen.char.State.Skill("Recon", "")
	if !found || got.Level != 1 {
		t.Errorf("a rank benefit with no level granted %d, want 1", got.Level)
	}
}

// TestARankBenefitThatCannotBeApplied. A rank row's non-skill effects go
// through the ordinary fold, and an effect kind that does not exist comes
// back as an error rather than being skipped.
func TestARankBenefitThatCannotBeApplied(t *testing.T) {
	t.Parallel()

	gen := engine(t, 34)

	err := gen.applyRankBenefits([]career.Effect{{Kind: career.EffectKind(-1)}}, 0)
	if !errors.Is(err, ErrUnknownEffect) {
		t.Errorf("err = %v, want ErrUnknownEffect", err)
	}
}
