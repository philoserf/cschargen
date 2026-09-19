package setting

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// sampleJSON is the invented data the repository ships so that the tests
// and the demo run without a transcription of the book. Every name in it
// was made up for this file.
//
//go:embed sample.json
var sampleJSON []byte

// Sample returns the repository's invented data. A character generated
// against it is valid and replayable, and is stamped so that it is never
// mistaken for one set on the published worlds.
func Sample() (*Data, error) {
	data, err := parse(sampleJSON, "sample")
	if err != nil {
		return nil, fmt.Errorf("the embedded sample data is not valid: %w", err)
	}

	data.Sample = true

	return data, nil
}

// Load reads a setting file and validates it. A file that does not validate
// is refused here rather than halfway through a lifepath, which is the
// whole reason the validator is separate from the reader.
func Load(path string) (*Data, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	return parse(content, path)
}

func parse(content []byte, name string) (*Data, error) {
	var data Data

	// DisallowUnknownFields turns a misspelt key into an error the user can
	// fix, rather than a field silently left at its zero value -- which for
	// maximumTerms would be a character who may serve no terms at all.
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("%s is not setting data: %w", name, err)
	}

	problems := Validate(&data)
	if len(problems) > 0 {
		return nil, &InvalidError{Source: name, Problems: problems}
	}

	sum := sha256.Sum256(content)

	data.Hash = hex.EncodeToString(sum[:])

	if data.Name == "" {
		data.Name = name
	}

	return &data, nil
}

// InvalidError carries every problem found in a file, not the first. A
// user transcribing sixty worlds should be told about all of them at once.
type InvalidError struct {
	Source   string
	Problems []string
}

func (e *InvalidError) Error() string {
	var out strings.Builder

	fmt.Fprintf(&out, "%s has %d problem", e.Source, len(e.Problems))

	if len(e.Problems) != 1 {
		out.WriteString("s")
	}

	for _, problem := range e.Problems {
		out.WriteString("\n  - " + problem)
	}

	return out.String()
}
