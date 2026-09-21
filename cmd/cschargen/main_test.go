package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// capture runs the command with its output redirected to a temporary file
// and returns what it wrote. The command writes to an *os.File so that a
// record can go straight to disk without a buffer in between; tests give it
// a file rather than a pipe so nothing can block.
// cmdNew is the subcommand the tests drive, named because several cases
// spell it.
const (
	cmdNew   = "new"
	cmdBatch = "batch"
	cmdData  = "data"
	fileA    = "a.json"
	fileB    = "b.json"
)

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

// A generated record, written to a temporary file, for the commands that
// take one.
func record(t *testing.T, seed string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "character.json")

	_, err := capture(t, cmdNew, "--auto", "--seed", seed, "--terms", "3", "-o", path)
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	return path
}

func TestRenderWritesASheet(t *testing.T) {
	t.Parallel()

	out, err := capture(t, "render", record(t, "7"))
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	for _, want := range []string{"## Characteristics", "## Skills", "## Career history"} {
		if !strings.Contains(out, want) {
			t.Errorf("the sheet has no %q section:\n%s", want, out)
		}
	}
}

func TestRenderHistoryWritesATranscript(t *testing.T) {
	t.Parallel()

	out, err := capture(t, "render", "--history", record(t, "7"))
	if err != nil {
		t.Fatalf("render --history: %v", err)
	}

	if !strings.Contains(out, "# Lifepath") {
		t.Errorf("the transcript has no heading:\n%s", out)
	}
}

// TestFlagsPrecedeTheFilename: Go's flag package stops parsing at the first
// non-flag argument, so `render char.json --history` leaves --history
// standing as a second positional. That is a usage error, not a silent
// sheet where a transcript was asked for.
func TestFlagsPrecedeTheFilename(t *testing.T) {
	t.Parallel()

	_, err := capture(t, "render", record(t, "7"), "--history")
	if err == nil {
		t.Fatal("no error")
	}

	if !strings.HasPrefix(err.Error(), "usage:") {
		t.Errorf("error %q does not begin with usage:", err)
	}
}

func TestReplayReproducesARecord(t *testing.T) {
	t.Parallel()

	out, err := capture(t, "replay", record(t, "23"))
	if err != nil {
		t.Fatalf("replay: %v", err)
	}

	if !strings.Contains(out, "identical") {
		t.Errorf("replay did not report a match: %q", out)
	}
}

// TestReplayCatchesATamperedRecord: the log is verification data, so
// altering it has to be caught rather than reapplied.
func TestReplayCatchesATamperedRecord(t *testing.T) {
	t.Parallel()

	path := record(t, "23")

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading: %v", err)
	}

	// Change a recorded choice to one the engine did not make.
	tampered := strings.Replace(string(content), `"chosen": 0`, `"chosen": 1`, 1)
	if tampered == string(content) {
		t.Skip("the record carries no choice at index 0 to tamper with")
	}

	err = os.WriteFile(path, []byte(tampered), 0o600)
	if err != nil {
		t.Fatalf("writing: %v", err)
	}

	_, err = capture(t, "replay", path)
	if err == nil {
		t.Fatal("a tampered record replayed clean")
	}
}

func TestReplayRefusesARecordThisBuildDidNotWrite(t *testing.T) {
	t.Parallel()

	path := record(t, "23")

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading: %v", err)
	}

	stale := strings.Replace(string(content),
		`"engineVersion": "`+version()+`"`, `"engineVersion": "0.0.1-ancient"`, 1)
	if stale == string(content) {
		t.Skipf("the record does not carry engineVersion %q", version())
	}

	err = os.WriteFile(path, []byte(stale), 0o600)
	if err != nil {
		t.Fatalf("writing: %v", err)
	}

	_, err = capture(t, "replay", path)
	if err == nil {
		t.Fatal("a record from another engine replayed without complaint")
	}

	// ...and the flag waives that one check and no others.
	out, err := capture(t, "replay", "--ignore-provenance", path)
	if err != nil {
		t.Fatalf("--ignore-provenance: %v", err)
	}

	if !strings.Contains(out, "identical") {
		t.Errorf("--ignore-provenance did not then replay: %q", out)
	}
}

func TestRenderAndReplayUsageErrors(t *testing.T) {
	t.Parallel()

	tests := [][]string{
		{"render"},
		{"render", fileA, fileB},
		{"replay"},
		{"replay", fileA, fileB},
	}

	for _, args := range tests {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			t.Parallel()

			_, err := capture(t, args...)
			if err == nil {
				t.Fatal("no error")
			}

			if !strings.HasPrefix(err.Error(), "usage:") {
				t.Errorf("error %q does not begin with usage:", err)
			}
		})
	}
}

func TestReadingSomethingThatIsNotARecord(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "not.json")

	err := os.WriteFile(path, []byte("this is not JSON"), 0o600)
	if err != nil {
		t.Fatalf("writing: %v", err)
	}

	_, err = capture(t, "render", path)
	if err == nil {
		t.Fatal("a file that is not a record rendered")
	}

	if strings.HasPrefix(err.Error(), "usage:") {
		t.Errorf("a malformed file is not a usage error: %q", err)
	}

	_, err = capture(t, "render", filepath.Join(t.TempDir(), "absent.json"))
	if err == nil {
		t.Fatal("a file that does not exist rendered")
	}
}

func TestDataValidateReportsTheSample(t *testing.T) {
	t.Parallel()

	out, err := capture(t, cmdData, "validate", filepath.Join("..", "..", "setting", "sample.json"))
	if err != nil {
		t.Fatalf("data validate: %v", err)
	}

	for _, want := range []string{"subsectors", "worlds", "species", "sha256"} {
		if !strings.Contains(out, want) {
			t.Errorf("the summary does not mention %q:\n%s", want, out)
		}
	}
}

// TestDataValidateReportsAProblemRatherThanPanicking is why the command
// exists: the file is the user's own transcription, and a mistake in it
// should be named rather than surfacing halfway through a lifepath.
func TestDataValidateReportsAProblemRatherThanPanicking(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "broken.json")

	err := os.WriteFile(path, []byte(`{"schemaVersion":1,"name":"x","subsectors":[]}`), 0o600)
	if err != nil {
		t.Fatalf("writing: %v", err)
	}

	_, err = capture(t, cmdData, "validate", path)
	if err == nil {
		t.Fatal("a file with no subsectors validated")
	}

	if !strings.Contains(err.Error(), "no subsectors") {
		t.Errorf("the error does not name the problem: %v", err)
	}
}

func TestDataUsageErrors(t *testing.T) {
	t.Parallel()

	for _, args := range [][]string{
		{cmdData},
		{cmdData, "inspect"},
		{cmdData, "validate"},
		{cmdData, "validate", fileA, fileB},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			t.Parallel()

			_, err := capture(t, args...)
			if err == nil {
				t.Fatal("no error")
			}

			if !strings.HasPrefix(err.Error(), "usage:") {
				t.Errorf("error %q does not begin with usage:", err)
			}
		})
	}
}

// TestReplayRefusesDifferentSettingData is the check the content hash
// exists for. With the tables outside the binary, a seed alone no longer
// determines a character: a record replayed against a different file would
// diverge at Step 3 and never recover.
func TestReplayRefusesDifferentSettingData(t *testing.T) {
	t.Parallel()

	path := record(t, "31")

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading: %v", err)
	}

	// A record that claims to have come from data this build does not have.
	tampered := strings.Replace(string(content),
		`"sample": true`, `"sample": true, "hash": "0000000000000000"`, 1)
	if tampered == string(content) {
		t.Skip("the record does not carry the field this test alters")
	}

	err = os.WriteFile(path, []byte(tampered), 0o600)
	if err != nil {
		t.Fatalf("writing: %v", err)
	}

	_, err = capture(t, "replay", path)
	if err == nil {
		t.Fatal("a record generated against other setting data replayed clean")
	}

	if !strings.Contains(err.Error(), "setting data") {
		t.Errorf("the error does not name the mismatch: %v", err)
	}
}

func TestNewWithASettingFile(t *testing.T) {
	t.Parallel()

	sample := filepath.Join("..", "..", "setting", "sample.json")

	out, err := capture(t, cmdNew, "--auto", "--seed", "5", "--terms", "2", "--data", sample)
	if err != nil {
		t.Fatalf("new --data: %v", err)
	}

	if !strings.Contains(out, `"homeworlds"`) {
		t.Errorf("the record carries no homeworld:\n%s", out)
	}
}

func TestNewWithASettingFileThatIsNotThere(t *testing.T) {
	t.Parallel()

	_, err := capture(t, cmdNew, "--auto", "--data", filepath.Join(t.TempDir(), "absent.json"))
	if err == nil {
		t.Fatal("a missing setting file was accepted")
	}
}

// TestSkipFamilyReachesTheRecord. Step 5 is optional in the book (p. 57)
// and the flag is how a player says so, so it has to reach the Inputs the
// record stamps rather than stopping at the flag set.
func TestSkipFamilyReachesTheRecord(t *testing.T) {
	t.Parallel()

	out, err := capture(t, cmdNew, "--auto", "--seed", "3", "--terms", "1", "--skip-family")
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	var record struct {
		Provenance struct {
			Inputs struct {
				SkipFamily bool `json:"skipFamily"`
			} `json:"inputs"`
		} `json:"provenance"`
		State struct {
			Family *struct{} `json:"family"`
		} `json:"state"`
	}

	err = json.Unmarshal([]byte(out), &record)
	if err != nil {
		t.Fatalf("decoding the record: %v", err)
	}

	if !record.Provenance.Inputs.SkipFamily {
		t.Error("the record does not say the family step was skipped")
	}

	if record.State.Family != nil {
		t.Error("the family step ran although --skip-family was given")
	}
}

// TestNewWithoutAutoAsksTheirPlayer. `new` refused without --auto until
// interactive mode landed; it now asks, and a run with nothing to answer
// with ends rather than falling back on the policy.
func TestNewWithoutAutoAsksTheirPlayer(t *testing.T) {
	t.Parallel()

	_, err := capture(t, cmdNew, "--seed", "3", "--terms", "1")
	if err == nil {
		t.Fatal("a run with no answers generated a character anyway")
	}

	if !strings.Contains(err.Error(), "the player stopped answering") {
		t.Errorf("err = %v, want the abandoned session", err)
	}

	// And it is not a usage error: the command line was fine.
	if strings.HasPrefix(err.Error(), "usage:") {
		t.Errorf("err = %v, which blames the command line", err)
	}
}

// TestBatchWritesOneRecordPerLine. "batch emits JSONL, requires --auto, and
// derives each member's seed from the base seed plus index, recorded per
// record" (the PRD's CLI sketch).
func TestBatchWritesOneRecordPerLine(t *testing.T) {
	t.Parallel()

	out, err := capture(t, "batch", "--auto", "--count", "5", "--seed", "100", "--terms", "1")
	if err != nil {
		t.Fatalf("batch: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 5 {
		t.Fatalf("%d lines, want 5", len(lines))
	}

	seen := map[uint64]bool{}

	for i, line := range lines {
		var record struct {
			Provenance struct {
				RNG struct {
					Seed uint64 `json:"seed"`
				} `json:"rng"`
			} `json:"provenance"`
		}

		err = json.Unmarshal([]byte(line), &record)
		if err != nil {
			t.Fatalf("line %d does not decode: %v", i+1, err)
		}

		// The seed of member i is the base plus i, and is on that member's
		// own record -- so any one of them regenerates alone.
		want := uint64(100 + i)
		if record.Provenance.RNG.Seed != want {
			t.Errorf("line %d carries seed %d, want %d", i+1, record.Provenance.RNG.Seed, want)
		}

		seen[record.Provenance.RNG.Seed] = true
	}

	if len(seen) != 5 {
		t.Errorf("%d distinct seeds across five records", len(seen))
	}
}

// TestBatchNeedsWhatItNeeds.
func TestBatchNeedsWhatItNeeds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want string
	}{
		{"no count", []string{"batch", "--auto"}, "batch needs --count N"},
		{"a count of zero", []string{"batch", "--auto", "--count", "0"}, "batch needs --count N"},
		{"no auto", []string{"batch", "--count", "2"}, "batch needs --auto"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := capture(t, tc.args...)
			if err == nil {
				t.Fatal("it ran anyway")
			}

			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("err = %v, want %q", err, tc.want)
			}
		})
	}
}

// TestABatchMemberRegeneratesAlone is the point of recording each seed.
func TestABatchMemberRegeneratesAlone(t *testing.T) {
	t.Parallel()

	batch, err := capture(t, "batch", "--auto", "--count", "3", "--seed", "200", "--terms", "1")
	if err != nil {
		t.Fatalf("batch: %v", err)
	}

	third := strings.Split(strings.TrimSpace(batch), "\n")[2]

	alone, err := capture(t, cmdNew, "--auto", "--seed", "202", "--terms", "1")
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	// The batch writes compact JSON and `new` writes it indented, so the
	// two are compared as values rather than as text.
	var fromBatch, fromNew any

	err = json.Unmarshal([]byte(third), &fromBatch)
	if err != nil {
		t.Fatalf("decoding the batch member: %v", err)
	}

	err = json.Unmarshal([]byte(alone), &fromNew)
	if err != nil {
		t.Fatalf("decoding the single record: %v", err)
	}

	if !reflect.DeepEqual(fromBatch, fromNew) {
		t.Error("the third member of the batch is not what its own seed generates")
	}
}

// TestBatchWritesAFile, and refuses to overwrite one without --force,
// which is `new`'s rule and should be the same rule.
func TestBatchWritesAFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "crew.jsonl")

	_, err := capture(t, "batch", "--auto", "--count", "2", "--seed", "1", "--terms", "1",
		"-o", path)
	if err != nil {
		t.Fatalf("batch: %v", err)
	}

	written, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the batch: %v", err)
	}

	if lines := strings.Count(string(written), "\n"); lines != 2 {
		t.Errorf("%d lines in the file, want 2", lines)
	}

	// A second run over the same path is refused.
	_, err = capture(t, "batch", "--auto", "--count", "2", "--seed", "1", "--terms", "1",
		"-o", path)
	if err == nil {
		t.Error("a batch overwrote a file without --force")
	}
}

// TestBatchRefusesABadDataFile. The setting is loaded once for the whole
// batch, so a bad one fails before any character is generated.
func TestBatchRefusesABadDataFile(t *testing.T) {
	t.Parallel()

	_, err := capture(t, "batch", "--auto", "--count", "2", "--data", "nowhere.json")
	if err == nil {
		t.Fatal("a batch ran against a data file that does not exist")
	}
}

// TestBatchRefusesAPositionalArgument, for the reason `new` does.
func TestBatchRefusesAPositionalArgument(t *testing.T) {
	t.Parallel()

	_, err := capture(t, "batch", "--auto", "--count", "2", "extra")
	if err == nil {
		t.Fatal("a batch took a positional argument")
	}

	if !strings.HasPrefix(err.Error(), "usage:") {
		t.Errorf("err = %v, want a usage error", err)
	}
}

// TestBatchRefusesACharacterItCannotGenerate. A species the data does not
// declare fails the first member, and the error says which.
func TestBatchRefusesACharacterItCannotGenerate(t *testing.T) {
	t.Parallel()

	_, err := capture(t, "batch", "--auto", "--count", "3", "--species", "Basilisk")
	if err == nil {
		t.Fatal("a batch generated a species the data does not declare")
	}

	if !strings.Contains(err.Error(), "character 1 of 3") {
		t.Errorf("err = %v, which does not say which member failed", err)
	}
}

// TestBatchToAStreamThatIsClosed. `batch` writes to stdout unless -o, and a
// stream that will not take it is an error rather than a silent success.
func TestBatchToAStreamThatIsClosed(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "closed")

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("creating the file: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("closing the file: %v", err)
	}

	err = run([]string{"batch", "--auto", "--count", "1", "--terms", "1"}, file)
	if err == nil {
		t.Error("a batch written to a closed file succeeded")
	}
}

// TestBatchRefusesTheStep20Flags. `--name`, `--gender`, `--appearance` and
// `--goals` set one character's Step 20 fields. `batch` shares the flag set
// because the two commands take the same inputs otherwise, and applying
// them there would give twenty people the same name and the same face.
//
// Refusing costs the caller one line. Applying them quietly costs twenty
// wrong sheets and the time to work out why.
func TestBatchRefusesTheStep20Flags(t *testing.T) {
	t.Parallel()

	for _, flag := range []string{"name", "gender", "appearance", "goals"} {
		t.Run(flag, func(t *testing.T) {
			t.Parallel()

			_, err := capture(t, cmdBatch, "--auto", "--count", "2", "--"+flag, "x")
			if err == nil {
				t.Fatal("batch accepted a Step 20 flag")
			}

			if !strings.HasPrefix(err.Error(), "usage:") {
				t.Errorf("error %q does not begin with usage:", err)
			}

			if !strings.Contains(err.Error(), "--"+flag) {
				t.Errorf("error %q does not name the flag", err)
			}
		})
	}
}

// TestEachCommandsHelpNamesItself. A flag set carries the name it was built
// with, and `new` and `batch` share theirs — so `batch --help` announced
// `Usage of new:` for as long as nothing passed the command's own name in.
func TestEachCommandsHelpNamesItself(t *testing.T) {
	t.Parallel()

	for _, command := range []string{cmdNew, cmdBatch} {
		t.Run(command, func(t *testing.T) {
			t.Parallel()

			_, err := capture(t, command, "--nope")
			if err == nil {
				t.Fatal("an unknown flag was accepted")
			}

			if !strings.Contains(err.Error(), command+":") {
				t.Errorf("%s reported the error as %q, which does not name it", command, err)
			}
		})
	}
}
