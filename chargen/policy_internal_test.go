package chargen

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// policyDoc is the auto-mode decision table, which lives beside the module
// rather than inside this package.
const policyDoc = "../POLICY.md"

// choicePointPattern finds the identifier of a choice point in the source:
// every Choice literal names one, and it is the stable key a record carries.
// One is written as a constant rather than inline, so the constant's own
// declaration is matched too.
var choicePointPattern = regexp.MustCompile(`(?:Point:\s+|point[A-Za-z]+ += )"([a-z_]+)"`)

// askedPoints are the two put to a player and never to the policy. POLICY.md
// has a section on each rather than a row, because a row would have to say
// what the policy takes and the policy is not offered them.
var askedPoints = map[string]bool{
	"finishing":           true,
	"return_to_education": true,
}

// TestPolicyDocumentsEveryChoicePoint is POLICY.md's own rule, made a gate:
// "Rows are added as the engine reaches the choice points, in the PR that
// reaches them. A row here and no code is as wrong as code and no row."
//
// It went twelve pull requests without being honoured, which is what a rule
// with nothing checking it does. The document is what a record's
// policy_version identifies, so a choice point missing from it is a
// decision no reader of a record can look up.
func TestPolicyDocumentsEveryChoicePoint(t *testing.T) {
	t.Parallel()

	doc, err := os.ReadFile(policyDoc)
	if err != nil {
		t.Fatalf("reading %s: %v", policyDoc, err)
	}

	documented := string(doc)
	points := choicePointsInSource(t)

	if len(points) == 0 {
		t.Fatal("no choice point found in the source; the pattern has stopped matching")
	}

	for point := range points {
		if askedPoints[point] {
			continue
		}

		if !strings.Contains(documented, "`"+point+"`") {
			t.Errorf("%q is a choice point the engine offers and POLICY.md does not name", point)
		}
	}

	// And the other way: a row for a point the engine no longer offers is a
	// decision a reader would go looking for and never find in a record.
	for _, row := range regexp.MustCompile("\n\\| `([a-z_]+)`").FindAllStringSubmatch(documented, -1) {
		if !points[row[1]] {
			t.Errorf("POLICY.md has a row for %q, which the engine does not offer", row[1])
		}
	}
}

// choicePointsInSource is every choice point named in this package.
func choicePointsInSource(t *testing.T) map[string]bool {
	t.Helper()

	found := map[string]bool{}

	sources, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("listing the package: %v", err)
	}

	for _, source := range sources {
		if strings.HasSuffix(source, "_test.go") {
			continue
		}

		body, err := os.ReadFile(source)
		if err != nil {
			t.Fatalf("reading %s: %v", source, err)
		}

		for _, match := range choicePointPattern.FindAllStringSubmatch(string(body), -1) {
			found[match[1]] = true
		}
	}

	return found
}
