package chargen_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

func TestDeviateRecordsEachReadingOnce(t *testing.T) {
	t.Parallel()

	var prov chargen.Provenance

	prov.Deviate("E-3")
	prov.Deviate("E-4")
	prov.Deviate("E-3")

	want := []string{"E-3", "E-4"}
	if len(prov.Deviations) != len(want) {
		t.Fatalf("Deviations = %v, want %v", prov.Deviations, want)
	}

	for i, id := range want {
		if prov.Deviations[i] != id {
			t.Errorf("Deviations[%d] = %q, want %q", i, prov.Deviations[i], id)
		}
	}
}

// TestRulesetNamesThePrintedArtifact: a page cite is only checkable against
// the edition it was read from, so the record has to say which.
func TestRulesetNamesThePrintedArtifact(t *testing.T) {
	t.Parallel()

	ruleset := chargen.Ruleset
	for _, fragment := range []string{"Clement Sector", "Character Creation", "2026"} {
		if !strings.Contains(ruleset, fragment) {
			t.Errorf("Ruleset %q does not name %q", ruleset, fragment)
		}
	}
}

func TestRecordRoundTrips(t *testing.T) {
	t.Parallel()

	var log chargen.Log

	log.Step("Step 2: Roll Characteristics", "p. 39")

	want := chargen.Character{
		Provenance: chargen.Provenance{
			SchemaVersion: chargen.SchemaVersion,
			Ruleset:       chargen.Ruleset,
			EngineVersion: "0.1.0",
			PolicyVersion: "0.1.0",
			RNG:           chargen.RNG{Algorithm: "math/rand/v2 PCG", Seed: 7},
			SettingData:   chargen.SettingData{Name: "sample", Sample: true},
			Inputs:        chargen.Inputs{Species: "human", TechLevel: 11, MaxTerms: 28},
			Deviations:    []string{"E-3"},
		},
		Events: log.Events(),
	}

	out, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got chargen.Character

	err = json.Unmarshal(out, &got)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Provenance.RNG.Seed != 7 {
		t.Errorf("seed = %d", got.Provenance.RNG.Seed)
	}

	if !got.Provenance.SettingData.Sample {
		t.Error("sample flag did not survive the round trip")
	}

	if len(got.Events) != 1 || got.Events[0].Step == nil {
		t.Fatalf("events = %v", got.Events)
	}

	if got.Events[0].Step.Cite != "p. 39" {
		t.Errorf("cite = %q", got.Events[0].Step.Cite)
	}
}

// TestSampleDataIsAlwaysVisible: a character generated against the invented
// sample is not set on real worlds, and a reader of the JSON must be able
// to tell without knowing what the file was called.
func TestSampleDataIsAlwaysVisible(t *testing.T) {
	t.Parallel()

	out, err := json.Marshal(chargen.SettingData{Name: "sample", Sample: true})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if !strings.Contains(string(out), `"sample":true`) {
		t.Errorf("SettingData serialized as %s, with no visible sample flag", out)
	}
}
