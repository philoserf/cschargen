package chargen_test

import (
	"encoding/json"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// TestARecordReloadsWithoutLoss: replay regenerates from the seed rather
// than reading state back, so nothing else would catch a field that does
// not survive being written and read again.
func TestARecordReloadsWithoutLoss(t *testing.T) {
	t.Parallel()

	for seed := range uint64(30) {
		original := lifepath(t, seed, 6)

		first, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("seed %d: marshal: %v", seed, err)
		}

		var reloaded chargen.Character

		err = json.Unmarshal(first, &reloaded)
		if err != nil {
			t.Fatalf("seed %d: unmarshal: %v", seed, err)
		}

		second, err := json.Marshal(&reloaded)
		if err != nil {
			t.Fatalf("seed %d: re-marshal: %v", seed, err)
		}

		if string(first) != string(second) {
			t.Errorf("seed %d: the record changed on reload", seed)
		}

		for _, part := range differences(original, &reloaded) {
			t.Errorf("seed %d: %s did not survive the reload", seed, part)
		}
	}
}

// differences names the parts of a record that did not survive a reload.
// Counting each separately is the point: "the record changed" says nothing
// about which field lost its tag.
func differences(original, reloaded *chargen.Character) []string {
	var changed []string

	counts := []struct {
		part        string
		want, found int
	}{
		{"age", original.State.Age, reloaded.State.Age},
		{"skills", len(original.State.Skills), len(reloaded.State.Skills)},
		{"services", len(original.State.Services), len(reloaded.State.Services)},
		{"terms", len(original.State.Terms), len(reloaded.State.Terms)},
		{"ties", len(original.State.Ties), len(reloaded.State.Ties)},
		{"injuries", len(original.State.Injuries), len(reloaded.State.Injuries)},
		{"stash", len(original.State.Stash), len(reloaded.State.Stash)},
		{"credits", original.State.Credits, reloaded.State.Credits},
		{"events", len(original.Events), len(reloaded.Events)},
	}

	for _, count := range counts {
		if count.want != count.found {
			changed = append(changed, count.part)
		}
	}

	if original.State.Characteristics != reloaded.State.Characteristics {
		changed = append(changed, "characteristics")
	}

	return changed
}
