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

		if reloaded.State.Age != original.State.Age ||
			len(reloaded.State.Skills) != len(original.State.Skills) ||
			len(reloaded.State.Services) != len(original.State.Services) ||
			len(reloaded.State.Ties) != len(original.State.Ties) ||
			len(reloaded.State.Injuries) != len(original.State.Injuries) ||
			len(reloaded.Events) != len(original.Events) {
			t.Errorf("seed %d: state did not survive the reload", seed)
		}
	}
}
