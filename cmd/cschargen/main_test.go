package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// capture runs the command with its output redirected to a temporary file
// and returns what it wrote. The command writes to an *os.File so that a
// record can go straight to disk without a buffer in between; tests give it
// a file rather than a pipe so nothing can block.
// cmdNew is the subcommand the tests drive, named because several cases
// spell it.
const cmdNew = "new"

func capture(t *testing.T, args ...string) (string, error) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "out")

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("creating the capture file: %v", err)
	}

	runErr := run(args, file)

	err = file.Close()
	if err != nil {
		t.Fatalf("closing the capture file: %v", err)
	}

	written, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the capture file: %v", err)
	}

	return string(written), runErr
}

// TestUsageErrorsSayUsage is the exit contract: a misused command line and
// an engine failure both exit 1, and the message is what tells them apart.
func TestUsageErrorsSayUsage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{"no command", nil},
		{"an unknown command", []string{"fly"}},
		{"new without --auto", []string{cmdNew}},
		{"new with a positional argument", []string{cmdNew, "--auto", "extra"}},
		{"an unknown flag", []string{cmdNew, "--auto", "--nope"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := capture(t, tc.args...)
			if err == nil {
				t.Fatal("no error")
			}

			if !strings.HasPrefix(err.Error(), "usage:") {
				t.Errorf("error %q does not begin with usage:", err)
			}
		})
	}
}

func TestHelpIsNotAnError(t *testing.T) {
	t.Parallel()

	for _, arg := range []string{"-h", "--help", "help"} {
		out, err := capture(t, arg)
		if err != nil {
			t.Errorf("%s: %v", arg, err)
		}

		if !strings.Contains(out, cmdNew) {
			t.Errorf("%s printed no command list: %q", arg, out)
		}
	}
}

func TestVersionReportsWhatARecordStamps(t *testing.T) {
	t.Parallel()

	out, err := capture(t, "version")
	if err != nil {
		t.Fatalf("version: %v", err)
	}

	for _, want := range []string{"cschargen", "schema", "policy " + policyVersion, "Clement Sector"} {
		if !strings.Contains(out, want) {
			t.Errorf("version output does not mention %q:\n%s", want, out)
		}
	}
}

func TestNewWritesARecord(t *testing.T) {
	t.Parallel()

	out, err := capture(t, cmdNew, "--auto", "--seed", "7", "--name", "Vela", "--tech-level", "11")
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	var record struct {
		Provenance struct {
			RNG struct {
				Seed uint64 `json:"seed"`
			} `json:"rng"`
			Inputs struct {
				Name      string `json:"name"`
				TechLevel int    `json:"techLevel"`
			} `json:"inputs"`
		} `json:"provenance"`
		State struct {
			Characteristics map[string]int `json:"characteristics"`
		} `json:"state"`
		Events []struct {
			Seq  int    `json:"seq"`
			Kind string `json:"kind"`
		} `json:"events"`
	}

	err = json.Unmarshal([]byte(out), &record)
	if err != nil {
		t.Fatalf("the record is not valid JSON: %v\n%s", err, out)
	}

	if record.Provenance.RNG.Seed != 7 {
		t.Errorf("seed = %d", record.Provenance.RNG.Seed)
	}

	if record.Provenance.Inputs.Name != "Vela" {
		t.Errorf("name = %q", record.Provenance.Inputs.Name)
	}

	if record.Provenance.Inputs.TechLevel != 11 {
		t.Errorf("techLevel = %d", record.Provenance.Inputs.TechLevel)
	}

	if len(record.State.Characteristics) != 6 {
		t.Errorf("%d characteristics", len(record.State.Characteristics))
	}

	if len(record.Events) == 0 {
		t.Fatal("the record carries no events")
	}
}

// TestASeedNobodyChoseIsStillRecorded: a character generated without a seed
// has to be reproducible afterwards, so the one the tool picked is in the
// record.
//
// It asks for no career terms. A random seed running a whole lifepath would
// reach different branches of the engine on every run, and under -coverpkg
// that makes the coverage ratchet depend on what the dice did -- which is
// how this test made the ratchet disagree between a developer's machine and
// CI. The seed is what is under test here, not the lifepath.
func TestASeedNobodyChoseIsStillRecorded(t *testing.T) {
	t.Parallel()

	seeds := map[uint64]bool{}

	for range 8 {
		out, err := capture(t, cmdNew, "--auto", "--terms", "-1")
		if err != nil {
			t.Fatalf("new: %v", err)
		}

		var record struct {
			Provenance struct {
				RNG struct {
					Seed uint64 `json:"seed"`
				} `json:"rng"`
			} `json:"provenance"`
		}

		err = json.Unmarshal([]byte(out), &record)
		if err != nil {
			t.Fatalf("unmarshal: %v", err)
		}

		seeds[record.Provenance.RNG.Seed] = true
	}

	// Eight runs landing on one seed would mean the tool is not choosing
	// one at all, which is the failure this test exists for -- a recorded
	// zero that nobody chose replays to a different character than the one
	// the caller saw.
	if len(seeds) == 1 {
		t.Error("eight seedless runs all recorded the same seed")
	}
}

// TestSeedZeroIsASeed: the default is zero, so the tool cannot use the
// default to mean "not given".
func TestSeedZeroIsASeed(t *testing.T) {
	t.Parallel()

	first, err := capture(t, cmdNew, "--auto", "--seed", "0")
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	second, err := capture(t, cmdNew, "--auto", "--seed", "0")
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	if first != second {
		t.Error("--seed 0 did not reproduce the same character")
	}
}

func TestOutputFileIsNotOverwrittenWithoutForce(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "character.json")

	_, err := capture(t, cmdNew, "--auto", "--seed", "1", "-o", path)
	if err != nil {
		t.Fatalf("first write: %v", err)
	}

	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading: %v", err)
	}

	_, err = capture(t, cmdNew, "--auto", "--seed", "2", "-o", path)
	if err == nil {
		t.Fatal("the second write was allowed")
	}

	if !strings.HasPrefix(err.Error(), "usage:") {
		t.Errorf("error %q does not begin with usage:", err)
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading: %v", err)
	}

	if string(before) != string(after) {
		t.Error("the file changed despite the refusal")
	}

	_, err = capture(t, cmdNew, "--auto", "--seed", "2", "-o", path, "--force")
	if err != nil {
		t.Fatalf("--force: %v", err)
	}

	forced, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading: %v", err)
	}

	if string(forced) == string(before) {
		t.Error("--force did not overwrite")
	}
}

func TestRecordEndsWithANewline(t *testing.T) {
	t.Parallel()

	out, err := capture(t, cmdNew, "--auto", "--seed", "1")
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	if !strings.HasSuffix(out, "\n") {
		t.Error("the record does not end with a newline")
	}
}
