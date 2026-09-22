package chargen

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/philoserf/cschargen/setting"
)

// Reproducible refuses a record this build cannot claim to reproduce. It is
// the provenance half of a replay: a refusal to try, as against a result.
//
// policy_version is deliberately not among the checks. Replay reapplies
// recorded choices and never consults the policy, so a record made under one
// policy replays under any other -- which is why PolicyVersion is recorded
// and not verified.
func Reproducible(record *Character, engineVersion string, data *setting.Data) error {
	if data.Hash != record.Provenance.SettingData.Hash {
		return fmt.Errorf(
			"%w: the record was generated against setting data hashed %.12s, this is %.12s",
			ErrProvenance, record.Provenance.SettingData.Hash, data.Hash)
	}

	if record.Provenance.SchemaVersion != SchemaVersion {
		return fmt.Errorf("%w: record is schema %d, this build writes %d",
			ErrProvenance, record.Provenance.SchemaVersion, SchemaVersion)
	}

	if record.Provenance.Ruleset != Ruleset {
		return fmt.Errorf("%w: record was read against %q",
			ErrProvenance, record.Provenance.Ruleset)
	}

	if record.Provenance.EngineVersion != engineVersion {
		return fmt.Errorf("%w: record was written by engine %s, this is %s",
			ErrProvenance, record.Provenance.EngineVersion, engineVersion)
	}

	return nil
}

// Verify names the first place a re-run disagrees with the record it came
// from: the event log first, then the character.
//
// The log comes first on purpose. A divergence shows up there several events
// before it reaches the character, and the sequence number is what a person
// needs in order to find it -- "choice 27 is return_to_education, engine
// asked skill_table" says where to look, and a differing character does not.
func Verify(original, replayed *Character) error {
	for i, want := range original.Events {
		if i >= len(replayed.Events) {
			return fmt.Errorf("%w: the replay ended after %d events; the record holds %d",
				ErrDiverged, len(replayed.Events), len(original.Events))
		}

		difference := eventDifference(want, replayed.Events[i])
		if difference != "" {
			return fmt.Errorf("%w at event %d: %s", ErrDiverged, want.Seq, difference)
		}
	}

	if len(replayed.Events) > len(original.Events) {
		return fmt.Errorf("%w: the replay ran to %d events; the record holds %d",
			ErrDiverged, len(replayed.Events), len(original.Events))
	}

	return verifyState(original, replayed)
}

// verifyState holds the whole character, after the log. It ran on
// Characteristics alone for a long while -- the only comparable struct in
// State, so `==` reached exactly as far as Go let it and the line was never
// revisited. A replay that lost three of Step 20's four fields reported
// "identical".
//
// The comparison goes through JSON rather than reflect.DeepEqual because the
// two sides are not built the same way: one was decoded from a file and the
// other came fresh from the engine, so an empty Skills is a non-nil empty
// slice on one side and nil on the other. Marshalling both is also comparing
// them in the representation the record is stored in, which is the one that
// has to match.
func verifyState(original, replayed *Character) error {
	before, err := json.Marshal(original.State)
	if err != nil {
		return fmt.Errorf("encoding the record's character: %w", err)
	}

	after, err := json.Marshal(replayed.State)
	if err != nil {
		return fmt.Errorf("encoding the replayed character: %w", err)
	}

	if !bytes.Equal(before, after) {
		return fmt.Errorf("%w: the logs agree but the characters do not", ErrDiverged)
	}

	return nil
}

// eventDifference describes how two events differ, or returns empty when
// they do not. Each kind is compared on the fields that identify it: a throw
// by its expression and total, a choice by its point and the index taken, a
// consequence by what it says happened.
func eventDifference(want, got Event) string {
	if want.Kind != got.Kind {
		return fmt.Sprintf("record has a %s, the replay a %s", want.Kind, got.Kind)
	}

	switch want.Kind {
	case EventStep:
		return stepDifference(want, got)
	case EventThrow:
		return throwDifference(want, got)
	case EventChoice:
		return choiceDifference(want, got)
	case EventConsequence:
		if want.Consequence.Detail != got.Consequence.Detail {
			return fmt.Sprintf("consequence %q against %q",
				want.Consequence.Detail, got.Consequence.Detail)
		}
	}

	return ""
}

func stepDifference(want, got Event) string {
	if want.Step.Name != got.Step.Name {
		return fmt.Sprintf("step %q against %q", want.Step.Name, got.Step.Name)
	}

	return ""
}

func throwDifference(want, got Event) string {
	if want.Throw.Total != got.Throw.Total || want.Throw.Expr != got.Throw.Expr {
		return fmt.Sprintf("throw %s=%d against %s=%d",
			want.Throw.Expr, want.Throw.Total, got.Throw.Expr, got.Throw.Total)
	}

	return ""
}

func choiceDifference(want, got Event) string {
	if want.Choice.Chosen != got.Choice.Chosen || want.Choice.Point != got.Choice.Point {
		return fmt.Sprintf("choice %s#%d against %s#%d",
			want.Choice.Point, want.Choice.Chosen, got.Choice.Point, got.Choice.Chosen)
	}

	return ""
}
