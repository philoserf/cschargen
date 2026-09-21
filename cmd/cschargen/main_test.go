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
func TestABatchMemberCarriesItsOwnSeedAndReplays(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	pool := filepath.Join(dir, "crew.jsonl")

	_, err := capture(t, cmdBatch, "--auto", "--count", "3", "--seed", "200",
		"--terms", "1", "-o", pool)
	if err != nil {
		t.Fatalf("batch: %v", err)
	}

	records, err := readRecords(pool)
	if err != nil {
		t.Fatalf("reading the batch: %v", err)
	}

	// The seed of member i is the base plus i, and it is recorded.
	for i, record := range records {
		if got := record.Provenance.RNG.Seed; got != 200+uint64(i) {
			t.Errorf("member %d carries seed %d, want %d", i+1, got, 200+uint64(i))
		}
	}

	// And that is enough to reproduce it exactly, which is what the seed is
	// for. `new --seed` is not: since policy_version 0.3.0 a batch draws its
	// career and assignment where `new --auto` takes the first of each, so
	// the two commands make different characters from one seed on purpose.
	// The record is the input to reproducing a member, and `replay` is how.
	third := filepath.Join(dir, "third.json")

	err = writeFile(third, mustMarshal(t, records[2]), false)
	if err != nil {
		t.Fatalf("writing the member: %v", err)
	}

	out, err := capture(t, "replay", third)
	if err != nil {
		t.Fatalf("replay: %v", err)
	}

	if !strings.Contains(out, "identical") {
		t.Errorf("replaying a batch member gave %q", out)
	}
}

// mustMarshal encodes a record or fails the test.
func mustMarshal(t *testing.T, value any) []byte {
	t.Helper()

	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("encoding: %v", err)
	}

	return encoded
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

// TestRenderReadsWhatBatchWrote is the two commands the README prints next
// to each other, made to compose. Twenty NPCs used to be 777 KB of JSON and
// no way to see them: `batch` wrote JSONL and `render` read one record.
func TestRenderReadsWhatBatchWrote(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	pool := filepath.Join(dir, "crew.jsonl")

	_, err := capture(t, cmdBatch, "--auto", "--count", "4", "--seed", "7", "-o", pool)
	if err != nil {
		t.Fatalf("batch: %v", err)
	}

	sheets, err := capture(t, "render", pool)
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if got := strings.Count(sheets, "## Characteristics"); got != 4 {
		t.Errorf("rendered %d sheets from a batch of four", got)
	}

	roster, err := capture(t, "render", "--roster", pool)
	if err != nil {
		t.Fatalf("render --roster: %v", err)
	}

	if got := strings.Count(strings.TrimRight(roster, "\n"), "\n") + 1; got != 4*3 {
		t.Errorf("the roster is %d lines for four characters, want 12", got)
	}
}

// TestBatchWritesADirectoryWhenGivenOne is the PRD's own CLI sketch --
// `-o dir|file.jsonl` -- which had never been built. A directory of records
// makes every other command work on them unchanged.
func TestBatchWritesADirectoryWhenGivenOne(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	_, err := capture(t, cmdBatch, "--auto", "--count", "12", "--seed", "7", "-o", dir)
	if err != nil {
		t.Fatalf("batch: %v", err)
	}

	// Two digits for twelve, so a listing and a glob sort correctly.
	for _, name := range []string{"npc-01.json", "npc-12.json"} {
		sheet, err := capture(t, "render", filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("render %s: %v", name, err)
		}

		if !strings.Contains(sheet, "## Characteristics") {
			t.Errorf("%s did not render as a sheet", name)
		}
	}
}

// TestAFileOfNoRecordsIsAnError, and says so as one record rather than as a
// batch whose first line is broken.
func TestAFileOfNoRecordsIsAnError(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "junk.json")

	err := os.WriteFile(path, []byte("not json at all\n"), 0o600)
	if err != nil {
		t.Fatalf("writing the fixture: %v", err)
	}

	_, err = capture(t, "render", path)
	if err == nil {
		t.Fatal("render accepted a file that is not a record")
	}

	if !strings.Contains(err.Error(), "is not a character record") {
		t.Errorf("error %q does not say what is wrong", err)
	}
}

// TestRenderSaysWhichLineOfABatchIsBroken, because "line 9" is what a
// person needs in order to find it in a file of a hundred.
func TestRenderSaysWhichLineOfABatchIsBroken(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	pool := filepath.Join(dir, "crew.jsonl")

	_, err := capture(t, cmdBatch, "--auto", "--count", "3", "--seed", "7", "-o", pool)
	if err != nil {
		t.Fatalf("batch: %v", err)
	}

	content, err := os.ReadFile(pool)
	if err != nil {
		t.Fatalf("reading the batch: %v", err)
	}

	const wanted = 3

	lines := strings.SplitN(string(content), "\n", wanted)
	if len(lines) < wanted {
		t.Fatalf("the batch is %d lines, want at least %d", len(lines), wanted)
	}

	broken := lines[0] + "\n{not a record}\n" + lines[2]

	err = os.WriteFile(pool, []byte(broken), 0o600)
	if err != nil {
		t.Fatalf("writing the broken batch: %v", err)
	}

	_, err = capture(t, "render", pool)
	if err == nil {
		t.Fatal("render accepted a batch with a broken line")
	}

	if !strings.Contains(err.Error(), "line 2") {
		t.Errorf("error %q does not name the line", err)
	}
}

// TestBatchIntoADirectoryRespectsForce. A directory of records is written
// one file at a time, and each goes through the same guard a single record
// does: writing over somebody's cast is not something to do quietly.
func TestBatchIntoADirectoryRespectsForce(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	_, err := capture(t, cmdBatch, "--auto", "--count", "2", "--seed", "7", "-o", dir)
	if err != nil {
		t.Fatalf("batch: %v", err)
	}

	_, err = capture(t, cmdBatch, "--auto", "--count", "2", "--seed", "9", "-o", dir)
	if err == nil {
		t.Fatal("a second batch overwrote the first without --force")
	}

	if !strings.Contains(err.Error(), "--force") {
		t.Errorf("error %q does not say how to proceed", err)
	}

	_, err = capture(t, cmdBatch, "--auto", "--count", "2", "--seed", "9", "--force", "-o", dir)
	if err != nil {
		t.Errorf("--force did not allow the overwrite: %v", err)
	}
}

// TestRenderRejectsAnUnknownFlag, reported under its own name.
func TestRenderRejectsAnUnknownFlag(t *testing.T) {
	t.Parallel()

	_, err := capture(t, "render", "--nope", "x.json")
	if err == nil {
		t.Fatal("render accepted an unknown flag")
	}

	if !strings.HasPrefix(err.Error(), "usage: render:") {
		t.Errorf("error %q is not reported under render", err)
	}
}

// TestReplayRejectsAFileThatIsNotARecord. `render` learned to read a batch
// and stopped sharing `readRecord` with `replay`, which took this path's
// only test with it — replay is the one that needs it most, because it is
// the command people reach for when a record looks wrong.
func TestReplayRejectsAFileThatIsNotARecord(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "junk.json")

	err := os.WriteFile(path, []byte("{\"nope\":\n"), 0o600)
	if err != nil {
		t.Fatalf("writing the fixture: %v", err)
	}

	_, err = capture(t, "replay", path)
	if err == nil {
		t.Fatal("replay accepted a file that is not a record")
	}

	if !strings.Contains(err.Error(), "is not a character record") {
		t.Errorf("error %q does not say what is wrong", err)
	}
}

// TestABadNameFileStopsBothCommands, rather than generating a hundred
// unnamed characters and leaving the referee to notice.
func TestABadNameFileStopsBothCommands(t *testing.T) {
	t.Parallel()

	absent := filepath.Join(t.TempDir(), "absent.txt")

	for name, args := range map[string][]string{
		"new":   {cmdNew, "--auto", "--seed", "1", "--terms", "1", "--names", absent},
		"batch": {cmdBatch, "--auto", "--count", "2", "--seed", "1", "--names", absent},
	} {
		_, err := capture(t, args...)
		if err == nil {
			t.Errorf("%s accepted a name file that is not there", name)
		}
	}
}

// TestAnOriginFlagIsCheckedBeforeGenerating. A name the setting data does
// not have is a flag error, not a generation failure: the engine refuses it
// too, but partway through Step 3, which reads as the run having gone wrong
// rather than the command line having a typo in it.
func TestAnOriginFlagIsCheckedBeforeGenerating(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		mention string
	}{
		{"an unknown world", []string{"--homeworld", "Nowhere At All"}, "Nowhere At All"},
		{"an unknown subsector", []string{"--subsector", "No Such Reach"}, "No Such Reach"},
		{
			"a world that is not in the subsector named beside it",
			[]string{"--homeworld", "Quillon", "--subsector", "Tallow Drift"},
			"Quillon",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			args := append([]string{cmdNew, "--auto"}, test.args...)

			_, err := capture(t, args...)
			if err == nil {
				t.Fatal("the flag was accepted")
			}

			if !strings.HasPrefix(err.Error(), "usage:") {
				t.Errorf("error %q does not begin with usage:", err)
			}

			if !strings.Contains(err.Error(), test.mention) {
				t.Errorf("error %q does not name %q", err, test.mention)
			}
		})
	}
}

// TestAnOriginFlagReachesABatch: a cast is the reason the flags exist, so
// the check has to be in front of `batch` as well as `new`.
func TestAnOriginFlagReachesABatch(t *testing.T) {
	t.Parallel()

	_, err := capture(t, cmdBatch, "--auto", "--count", "2", "--homeworld", "Nowhere At All")
	if err == nil {
		t.Fatal("the flag was accepted")
	}

	if !strings.HasPrefix(err.Error(), "usage:") {
		t.Errorf("error %q does not begin with usage:", err)
	}
}

// TestAHomeworldFlagPutsTheCharacterThere, all the way through the command
// rather than in the engine alone.
func TestAHomeworldFlagPutsTheCharacterThere(t *testing.T) {
	t.Parallel()

	out, err := capture(t, cmdNew, "--auto", "--seed", "3", "--terms", "1",
		"--homeworld", "Quillon")
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	var record struct {
		Provenance struct {
			Inputs struct {
				Homeworld string `json:"homeworld"`
			} `json:"inputs"`
		} `json:"provenance"`
		State struct {
			Homeworlds []struct {
				World string `json:"world"`
			} `json:"homeworlds"`
		} `json:"state"`
	}

	err = json.Unmarshal([]byte(out), &record)
	if err != nil {
		t.Fatalf("the record is not valid JSON: %v", err)
	}

	if got := record.State.Homeworlds[0].World; got != "Quillon" {
		t.Errorf("asked for Quillon, born on %s", got)
	}

	if record.Provenance.Inputs.Homeworld != "Quillon" {
		t.Error("the record does not say the homeworld was asked for")
	}
}

// TestTermsZeroMeansZero. The engine reads TermLimit as zero for "not given"
// and negative for "no terms at all", because a plain int cannot tell an
// unset flag from one set to its zero value. That left `--terms 0` asking
// for none and getting the policy's four, silently.
func TestTermsZeroMeansZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want int
	}{
		{"zero is none", []string{"--terms", "0"}, 0},
		{"the old sentinel still is", []string{"--terms", "-1"}, 0},
		{"a number is itself", []string{"--terms", "2"}, 2},
		{"omitted is the policy default", nil, 4},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			args := append([]string{cmdNew, "--auto", "--seed", "21"}, test.args...)

			out, err := capture(t, args...)
			if err != nil {
				t.Fatalf("new: %v", err)
			}

			var record struct {
				State struct {
					Terms []struct{} `json:"terms"`
				} `json:"state"`
			}

			err = json.Unmarshal([]byte(out), &record)
			if err != nil {
				t.Fatalf("the record is not valid JSON: %v", err)
			}

			if got := len(record.State.Terms); got != test.want {
				t.Errorf("served %d terms, want %d", got, test.want)
			}
		})
	}
}
