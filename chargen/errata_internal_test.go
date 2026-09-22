package chargen

import (
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// errataDoc records every place the book is silent, ambiguous or wrong. Like
// POLICY.md it lives beside the module rather than inside this package.
const errataDoc = "../ERRATA.md"

// errataEntryPattern finds an entry heading and its kind. The three kinds are
// the file's own vocabulary: a typo is the book meaning what it plainly meant,
// a reading is a choice the engine had to make, and a limit is a rule the
// engine knowingly does not implement.
var errataEntryPattern = regexp.MustCompile(`(?m)^## (E-\d+) \((typo|reading|limit)\)`)

// deviatePattern finds an identifier this package stamps directly. Every
// argument to Deviate here is a string literal, so this sees all of them.
var deviatePattern = regexp.MustCompile(`Deviate\("(E-\d+)"\)`)

// careerPackage is the table data, which reaches a record through the other
// channel: it cannot call Deviate, so it names the reading on an effect or a
// career and the engine stamps the value when it applies it (#154).
const careerPackage = "../career"

// declaredPattern finds an identifier the table data declares. Any quoted
// identifier in career/ is one: that package cites ERRATA in prose as
// "ERRATA E-1", unquoted, so a quoted one is always a declaration for the
// engine to record.
var declaredPattern = regexp.MustCompile(`"(E-\d+)"`)

// unstampedReadings are the readings a record does not carry, and why.
//
// A reading reaches a record to say that it changed *this* character. One
// that changed every character alike distinguishes nothing, and that is the
// only reason left on this list -- three entries were here because they were
// implemented in career/, which sees no character, until career.Effect and
// career.Career gained a field naming the reading they rest on (#154).
//
// Each entry is a claim that has to stay true, which is what
// TestUnstampedReadingsIsStillTrue is for: an allowlist that outlives its
// reason is how a gate quietly stops checking.
var unstampedReadings = map[string]string{
	// Applies to every character alike. Stamping it on every record would
	// add a line to each and distinguish none of them.
	"E-5":  "every tie carries a rating from the first term; there is no character without one",
	"E-17": "which characteristics are physical and which mental is a definition, read the same way for everyone",
	"E-22": "the result that grants a tie decides its kind, at every grant, for every character",
	"E-23": "a tie with no printed rating starts mid-band, and most ties have no printed rating",
	"E-27": "the campaign's present year is one number, the same for every character generated",
}

// TestEveryReadingIsStampedOrSaysWhyNot is ERRATA.md's own promise, made a
// gate: "A record stamps the identifiers of the entries that changed that
// character, so a character generated under one reading stays auditable
// after the reading changes."
//
// That is the whole argument for Provenance.Deviations existing, and it went
// unheld until #118: eight identifiers of forty-four ever reached a record,
// and thirty-one more appeared in Go comments as the reason the code reads
// the way it does and left no mark on any character.
//
// POLICY.md has had this shape of test since #110 and ERRATA.md had none,
// which is why the two documents drifted differently.
func TestEveryReadingIsStampedOrSaysWhyNot(t *testing.T) {
	t.Parallel()

	entries := errataEntries(t)
	stamped := stampedIdentifiers(t)

	for _, id := range slices.Sorted(maps.Keys(entries)) {
		if entries[id] != "reading" || stamped[id] || unstampedReadings[id] != "" {
			continue
		}

		t.Errorf("%s is a reading that no record stamps, and unstampedReadings does not say why", id)
	}
}

// TestEveryStampedIdentifierHasAnEntry is the other half: a record naming a
// reading a reader cannot look up is worse than one not naming it at all.
func TestEveryStampedIdentifierHasAnEntry(t *testing.T) {
	t.Parallel()

	entries := errataEntries(t)
	stamped := stampedIdentifiers(t)

	if len(entries) == 0 || len(stamped) == 0 {
		t.Fatal("no entry or no stamp found; a pattern has stopped matching")
	}

	for _, id := range slices.Sorted(maps.Keys(stamped)) {
		if entries[id] == "" {
			t.Errorf("%s is stamped into records and %s has no entry for it", id, errataDoc)
		}
	}
}

// TestUnstampedReadingsIsStillTrue holds the allowlist to its own claims. An
// entry that outlives its reason -- renumbered, relabelled, deleted, or
// stamped after all -- silently excuses nothing, which is how a gate stops
// checking without anyone noticing.
func TestUnstampedReadingsIsStillTrue(t *testing.T) {
	t.Parallel()

	entries := errataEntries(t)
	stamped := stampedIdentifiers(t)

	for _, id := range slices.Sorted(maps.Keys(unstampedReadings)) {
		switch {
		case entries[id] == "":
			t.Errorf("unstampedReadings names %s, which %s no longer carries", id, errataDoc)
		case entries[id] != "reading":
			t.Errorf("unstampedReadings names %s, which is now a (%s) rather than a reading",
				id, entries[id])
		case stamped[id]:
			t.Errorf("unstampedReadings says %s is never stamped, and it is; remove the entry", id)
		}
	}
}

// errataEntries is every entry in ERRATA.md, by identifier, to its kind.
func errataEntries(t *testing.T) map[string]string {
	t.Helper()

	body, err := os.ReadFile(errataDoc)
	if err != nil {
		t.Fatalf("reading %s: %v", errataDoc, err)
	}

	found := map[string]string{}
	for _, match := range errataEntryPattern.FindAllStringSubmatch(string(body), -1) {
		found[match[1]] = match[2]
	}

	return found
}

// stampedIdentifiers is every identifier that can reach a record, by either
// channel: stamped here, or declared by the table data for the engine to
// stamp when it applies it.
func stampedIdentifiers(t *testing.T) map[string]bool {
	t.Helper()

	found := map[string]bool{}

	scan(t, ".", deviatePattern, found)
	scan(t, careerPackage, declaredPattern, found)

	return found
}

// scan collects every identifier a pattern finds in a package's own source.
func scan(t *testing.T, dir string, pattern *regexp.Regexp, found map[string]bool) {
	t.Helper()

	sources, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		t.Fatalf("listing %s: %v", dir, err)
	}

	if len(sources) == 0 {
		t.Fatalf("no source in %s; the package has moved", dir)
	}

	for _, source := range sources {
		// The tests call Deviate directly to check that it deduplicates.
		// Counting those would let a stamp exist in the suite and nowhere
		// a character can reach.
		if strings.HasSuffix(source, "_test.go") {
			continue
		}

		body, err := os.ReadFile(source)
		if err != nil {
			t.Fatalf("reading %s: %v", source, err)
		}

		for _, match := range pattern.FindAllStringSubmatch(string(body), -1) {
			found[match[1]] = true
		}
	}
}
