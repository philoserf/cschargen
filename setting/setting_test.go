package setting_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/philoserf/cschargen/setting"
)

func TestSampleLoadsAndValidates(t *testing.T) {
	t.Parallel()

	data, err := setting.Sample()
	if err != nil {
		t.Fatalf("the embedded sample does not validate: %v", err)
	}

	if !data.Sample {
		t.Error("the sample is not flagged as the sample")
	}

	if data.Hash == "" {
		t.Error("the sample carries no content hash")
	}

	if len(data.Subsectors) == 0 {
		t.Fatal("the sample has no subsectors")
	}
}

// TestTheSampleIsInvented is the Product Identity boundary, held as a test.
// No world, subsector or species name from the book may be in this
// repository (CLAUDE.md; OGL section 16, p. 335).
//
// The list is every proper name the origin charts and species sections of
// the book use that a careless copy would most likely bring with it,
// including the one term the OGL notice singles out by name. It is not
// exhaustive and cannot be -- the real guard is the rule -- but a
// transcription pasted in by accident would almost certainly trip it.
//
// The term appears here, in a test that exists to keep it out. That is the
// one place naming it does the opposite of reproducing it.
func TestTheSampleIsInvented(t *testing.T) {
	t.Parallel()

	fromTheBook := []string{
		"Hub", "Cascadia", "Franklin", "Sequoyah", "Earth Subsector",
		"Toku", "Kingston", "Reuschle", "Viteges", "Sheba", "Totaro",
		"Nyahururu", "Bingxue Shijie", "Ti Nsan Oke", "Sone Ke Amsu",
		"Aishan Ko Arama", "Gaishan", "Oskar", "Aquan", "Sniffer",
		"Achilles", "Altrant", "altrant",
	}

	content, err := os.ReadFile("sample.json")
	if err != nil {
		t.Fatalf("reading the sample: %v", err)
	}

	text := string(content)
	for _, name := range fromTheBook {
		if strings.Contains(text, name) {
			t.Errorf("the sample contains %q, which is Product Identity", name)
		}
	}
}

// TestSampleCoversEveryD100Result: a chart is a thing where every result
// lands somewhere, and a gap would leave a character with no homeworld.
func TestSampleCoversEveryD100Result(t *testing.T) {
	t.Parallel()

	data, err := setting.Sample()
	if err != nil {
		t.Fatalf("sample: %v", err)
	}

	for _, sub := range data.Subsectors {
		covered := map[int]bool{}
		selectable := false

		for _, world := range sub.Worlds {
			if !world.Selectable() {
				continue
			}

			selectable = true

			for result := 1; result <= 100; result++ {
				if world.Covers(result) {
					covered[result] = true
				}
			}
		}

		if !selectable {
			continue
		}

		for result := 1; result <= 100; result++ {
			if !covered[result] {
				t.Errorf("%s: a d100 of %d lands on no world", sub.Name, result)
			}
		}
	}
}

// The schema's keys, named once. They are a vocabulary the tests write
// against, and a typo in one of them would read as a validator bug rather
// than as a broken fixture.
const (
	keyName       = "name"
	keyWorlds     = "worlds"
	keySubsector  = "subsectors"
	keyTerms      = "maximumTerms"
	keyLanguages  = "primaryLanguages"
	keyAllowed    = "allowed"
	keyOneOf      = "oneOf"
	keyEngineered = "engineered"
	statusFree    = "free"
	keyStatus     = "status"
	keyAge        = "maximumAge"
	keyKind       = "kind"
	keyUplifts    = "uplifts"
	langTrade     = "Trade"
	wantNoName    = "has no name"
	keyTech       = "techLevel"
	kindUplift    = "uplift"
	speciesCorvid = "Corvid"
)

// write puts a setting file in a temporary directory and returns its path.
func write(t *testing.T, data map[string]any) string {
	t.Helper()

	encoded, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	path := filepath.Join(t.TempDir(), "setting.json")

	err = os.WriteFile(path, encoded, 0o600)
	if err != nil {
		t.Fatalf("write: %v", err)
	}

	return path
}

// fixtureWorld is a world that validates, written once so that every test
// breaking one breaks a copy of the same thing. A nil roll makes it
// choose-only.
func fixtureWorld(name string, roll []int) map[string]any {
	world := map[string]any{
		keyName:       name,
		keyTech:       10,
		keyAge:        100,
		keyTerms:      20,
		keyLanguages:  []string{langTrade},
		keyEngineered: map[string]any{keyAllowed: true, keyStatus: statusFree},
		keyUplifts:    map[string]any{keyAllowed: true, keyStatus: statusFree},
	}

	if roll != nil {
		world["roll"] = roll
	}

	return world
}

// minimal is the smallest file that validates, for a test to break in one
// named way.
func minimal() map[string]any {
	return map[string]any{
		"schemaVersion": setting.SchemaVersion,
		keyName:         "test",
		keySubsector: []any{map[string]any{
			keyName:      "Only",
			"originRoll": 1,
			keyWorlds:    []any{fixtureWorld("Somewhere", []int{1, 100})},
		}},
	}
}

func TestMinimalFileValidates(t *testing.T) {
	t.Parallel()

	_, err := setting.Load(write(t, minimal()))
	if err != nil {
		t.Fatalf("the minimal file does not validate: %v", err)
	}
}

// TestTheValidatorNamesWhatIsWrong walks one broken file per rule. Each case
// says which problem it expects to see, because a validator that reports
// *something* for every broken file is not the same as one that reports the
// right thing.
func TestTheValidatorNamesWhatIsWrong(t *testing.T) {
	t.Parallel()

	for _, tc := range brokenFiles() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			broken := minimal()
			tc.breakIt(broken)

			_, err := setting.Load(write(t, broken))
			if err == nil {
				t.Fatal("the file validated")
			}

			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error does not mention %q:\n%v", tc.want, err)
			}
		})
	}
}

// brokenFile is one file broken in one named way, and the problem the
// validator must name for it. A validator that reports *something* for
// every broken file is not the same as one that reports the right thing.
type brokenFile struct {
	name    string
	breakIt func(m map[string]any)
	want    string
}

func brokenFiles() []brokenFile {
	return []brokenFile{
		{name: "a schema version this build does not read", breakIt: func(m map[string]any) {
			m["schemaVersion"] = 99
		}, want: "schemaVersion"},
		{name: "no subsectors", breakIt: func(m map[string]any) {
			m[keySubsector] = []any{}
		}, want: "no subsectors"},
		{name: "a subsector with no worlds", breakIt: func(m map[string]any) {
			sub(m)[keyWorlds] = []any{}
		}, want: "no worlds"},
		{name: "an originRoll off the 1d6 chart", breakIt: func(m map[string]any) {
			sub(m)["originRoll"] = 9
		}, want: "1d6"},
		{name: "a d100 range that is not one", breakIt: func(m map[string]any) {
			world(m)["roll"] = []int{40, 20}
		}, want: "not a d100 range"},
		{name: "a gap in the chart", breakIt: func(m map[string]any) {
			world(m)["roll"] = []int{1, 50}
		}, want: "land on no world"},
		{name: "no primary language", breakIt: func(m map[string]any) {
			world(m)[keyLanguages] = []string{}
		}, want: "primaryLanguages"},
		// A world may force one of the three results the Human Birth
		// Situation chart prints (p. 58). A fourth would generate a
		// household with no parents in it.
		{name: "a birth situation the chart does not print", breakIt: func(m map[string]any) {
			world(m)["birthSituationOnly"] = "a wolf pack"
		}, want: "birthSituationOnly"},
		{name: "no maximum terms", breakIt: func(m map[string]any) {
			world(m)[keyTerms] = 0
		}, want: "maximumTerms"},
		{name: "a status that is neither free nor enslaved", breakIt: func(m map[string]any) {
			world(m)[keyEngineered] = map[string]any{keyAllowed: true, keyStatus: "tolerated"}
		}, want: "free or enslaved"},
		{name: "a ban on a species that is not listed", breakIt: func(m map[string]any) {
			world(m)[keyEngineered] = map[string]any{
				keyAllowed: true, keyStatus: statusFree, "banned": []string{"Nobody"},
			}
		}, want: "not in the species list"},
		{name: "a background skill offering nothing", breakIt: func(m map[string]any) {
			world(m)["backgroundSkills"] = []any{map[string]any{keyOneOf: []any{}}}
		}, want: "offers no alternatives"},
		{name: "a background alternative that is neither", breakIt: func(m map[string]any) {
			world(m)["backgroundSkills"] = []any{
				map[string]any{keyOneOf: []any{map[string]any{}}},
			}
		}, want: "neither a skill nor an item"},
		{name: "an item with specialties", breakIt: func(m map[string]any) {
			world(m)["backgroundSkills"] = []any{map[string]any{keyOneOf: []any{
				map[string]any{"item": "a thing", "specialties": []string{"Any"}},
			}}}
		}, want: "item with specialties"},
	}
}

// sub and world reach into the fixture. The assertions are checked so that
// a fixture whose shape drifts fails as "the fixture is wrong" rather than
// as a panic inside a table-driven test, where the case that panicked is
// the hardest thing to identify.
func sub(m map[string]any) map[string]any {
	subsectors, ok := m[keySubsector].([]any)
	if !ok || len(subsectors) == 0 {
		panic("fixture has no subsectors")
	}

	first, ok := subsectors[0].(map[string]any)
	if !ok {
		panic("fixture subsector is not an object")
	}

	return first
}

func world(m map[string]any) map[string]any {
	worlds, ok := sub(m)[keyWorlds].([]any)
	if !ok || len(worlds) == 0 {
		panic("fixture has no worlds")
	}

	first, ok := worlds[0].(map[string]any)
	if !ok {
		panic("fixture world is not an object")
	}

	return first
}

// TestTheValidatorReportsEveryProblem: someone transcribing sixty worlds
// should be told about all of them in one pass rather than made to run the
// validator sixty times.
func TestTheValidatorReportsEveryProblem(t *testing.T) {
	t.Parallel()

	broken := minimal()

	world(broken)[keyTerms] = 0
	world(broken)[keyAge] = 0
	world(broken)[keyLanguages] = []string{}

	_, err := setting.Load(write(t, broken))
	if err == nil {
		t.Fatal("the file validated")
	}

	for _, want := range []string{"maximumTerms", keyAge, "primaryLanguages"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("only some problems were reported; %q is missing:\n%v", want, err)
		}
	}
}

// TestAMisspeltKeyIsAnError: a field silently left at its zero value would
// be a world whose character may serve no terms.
func TestAMisspeltKeyIsAnError(t *testing.T) {
	t.Parallel()

	broken := minimal()

	world(broken)["maximumTurns"] = 20

	_, err := setting.Load(write(t, broken))
	if err == nil {
		t.Fatal("a misspelt key was accepted")
	}
}

func TestLoadingAFileThatIsNotThere(t *testing.T) {
	t.Parallel()

	_, err := setting.Load(filepath.Join(t.TempDir(), "absent.json"))
	if err == nil {
		t.Fatal("a file that does not exist loaded")
	}
}

func TestTheHashFollowsTheContent(t *testing.T) {
	t.Parallel()

	first, err := setting.Load(write(t, minimal()))
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	same, err := setting.Load(write(t, minimal()))
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if first.Hash != same.Hash {
		t.Error("the same content hashed differently")
	}

	changed := minimal()

	world(changed)[keyTech] = 11

	other, err := setting.Load(write(t, changed))
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if first.Hash == other.Hash {
		t.Error("a changed tech level did not change the hash")
	}
}

func TestPermissionAdmits(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		permission setting.Permission
		species    string
		want       bool
	}{
		{"allowed, nothing banned", setting.Permission{Allowed: true, Status: setting.Free}, speciesCorvid, true},
		{"not allowed at all", setting.Permission{}, speciesCorvid, false},
		{"allowed but this one banned", setting.Permission{
			Allowed: true, Status: setting.Free, Banned: []string{speciesCorvid},
		}, speciesCorvid, false},
		{"allowed, a different one banned", setting.Permission{
			Allowed: true, Status: setting.Free, Banned: []string{"Otter"},
		}, speciesCorvid, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.permission.Admits(tc.species); got != tc.want {
				t.Errorf("Admits(%q) = %v, want %v", tc.species, got, tc.want)
			}
		})
	}
}

func TestLookups(t *testing.T) {
	t.Parallel()

	data, err := setting.Sample()
	if err != nil {
		t.Fatalf("sample: %v", err)
	}

	first := data.Subsectors[0]

	found, ok := data.Subsector(first.Name)
	if !ok || found.Name != first.Name {
		t.Errorf("Subsector(%q) did not find it", first.Name)
	}

	_, ok = data.Subsector("Nowhere")
	if ok {
		t.Error("a subsector that does not exist was found")
	}

	world, inSub, ok := data.World(first.Worlds[0].Name)
	if !ok || inSub.Name != first.Name || world.Name != first.Worlds[0].Name {
		t.Errorf("World(%q) found %q in %q", first.Worlds[0].Name, world.Name, inSub.Name)
	}

	_, _, ok = data.World("Nowhere")
	if ok {
		t.Error("a world that does not exist was found")
	}
}

// The cases below are the validator's remaining branches: the ones a file
// reaches only by being wrong in a way the table above does not cover.

func TestMoreWaysAFileCanBeWrong(t *testing.T) {
	t.Parallel()

	for _, tc := range []brokenFile{
		{name: "a species with no name", breakIt: func(m map[string]any) {
			m["species"] = []any{map[string]any{keyKind: kindUplift}}
		}, want: wantNoName},
		{name: "a species listed twice", breakIt: func(m map[string]any) {
			m["species"] = []any{
				map[string]any{keyName: speciesCorvid, keyKind: kindUplift},
				map[string]any{keyName: speciesCorvid, keyKind: kindUplift},
			}
		}, want: "listed twice"},
		{name: "a species that is neither kind", breakIt: func(m map[string]any) {
			m["species"] = []any{map[string]any{keyName: speciesCorvid, keyKind: "robot"}}
		}, want: "engineered or uplift"},
		{name: "a subsector with no name", breakIt: func(m map[string]any) {
			delete(sub(m), keyName)
		}, want: wantNoName},
		{name: "a world with no name", breakIt: func(m map[string]any) {
			delete(world(m), keyName)
		}, want: wantNoName},
		{name: "a tech level off any scale", breakIt: func(m map[string]any) {
			world(m)[keyTech] = 99
		}, want: keyTech},
		{name: "a maximum age of nothing", breakIt: func(m map[string]any) {
			world(m)[keyAge] = 0
		}, want: keyAge},
		{name: "banned species although none are allowed", breakIt: func(m map[string]any) {
			world(m)[keyEngineered] = map[string]any{
				keyAllowed: false, "banned": []string{speciesCorvid},
			}
		}, want: "although none are allowed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			broken := minimal()

			tc.breakIt(broken)

			_, err := setting.Load(write(t, broken))
			if err == nil {
				t.Fatal("the file validated")
			}

			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error does not mention %q:\n%v", tc.want, err)
			}
		})
	}
}

// TestTwoSubsectorsCannotShareAnOriginRoll: the 1d6 of p. 39 has to land on
// one subsector, and two claiming the same result would make one of them
// unreachable.
func TestTwoSubsectorsCannotShareAnOriginRoll(t *testing.T) {
	t.Parallel()

	broken := minimal()

	subsectors, ok := broken[keySubsector].([]any)
	if !ok {
		t.Fatal("fixture")
	}

	second := map[string]any{
		keyName:      "Second",
		"originRoll": 1,
		keyWorlds:    []any{fixtureWorld("Elsewhere", []int{1, 100})},
	}

	broken[keySubsector] = append(subsectors, second)

	_, err := setting.Load(write(t, broken))
	if err == nil {
		t.Fatal("two subsectors shared an originRoll")
	}

	if !strings.Contains(err.Error(), "both claim originRoll") {
		t.Errorf("error does not name the clash:\n%v", err)
	}
}

// TestAWorldNameIdentifiesOneWorld: a mishap reassigns a homeworld by name,
// so two worlds sharing one would make the reassignment ambiguous.
func TestAWorldNameIdentifiesOneWorld(t *testing.T) {
	t.Parallel()

	broken := minimal()

	subsectors, ok := broken[keySubsector].([]any)
	if !ok {
		t.Fatal("fixture has no subsectors")
	}

	// the same name the minimal file's world already has
	twin := fixtureWorld("Somewhere", nil)

	broken[keySubsector] = append(subsectors, map[string]any{
		keyName:   "Second",
		keyWorlds: []any{twin},
	})

	_, err := setting.Load(write(t, broken))
	if err == nil {
		t.Fatal("two worlds shared a name")
	}

	if !strings.Contains(err.Error(), "names identify a homeworld") {
		t.Errorf("error does not name the clash:\n%v", err)
	}
}

func TestCoversIsFalseForAChooseOnlyWorld(t *testing.T) {
	t.Parallel()

	var world setting.World

	if world.Selectable() {
		t.Error("a world with no roll reported itself selectable")
	}

	for result := 1; result <= 100; result++ {
		if world.Covers(result) {
			t.Fatalf("a world with no roll covers %d", result)
		}
	}
}

func TestCoversTheEndsOfItsRange(t *testing.T) {
	t.Parallel()

	data, err := setting.Sample()
	if err != nil {
		t.Fatalf("sample: %v", err)
	}

	for _, sub := range data.Subsectors {
		for _, world := range sub.Worlds {
			if !world.Selectable() {
				continue
			}

			low, high := world.Roll[0], world.Roll[1]
			if !world.Covers(low) || !world.Covers(high) {
				t.Errorf("%s does not cover its own range %d-%d", world.Name, low, high)
			}

			if world.Covers(low-1) || world.Covers(high+1) {
				t.Errorf("%s covers beyond its range %d-%d", world.Name, low, high)
			}
		}
	}
}

// aSpecies is the name the broken species below all carry; what is broken
// about them is never the name.
const aSpecies = "Thing"

// TestASpeciesMustDeclareThingsTheEngineCanCarryOut. A species names an
// aging profile the engine holds, characteristics the book has, and a
// ceiling inside the range a score can reach.
func TestASpeciesMustDeclareThingsTheEngineCanCarryOut(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		species setting.Species
		want    string
	}{
		{
			name:    "a kind that is neither",
			species: setting.Species{Name: aSpecies, Kind: "robot"},
			want:    "engineered or uplift",
		},
		{
			name: "an aging profile the engine does not hold",
			species: setting.Species{
				Name: aSpecies, Kind: kindUplift, Aging: "glacial",
			},
			want: "not a profile the engine holds",
		},
		{
			name: "a characteristic this ruleset does not have",
			species: setting.Species{
				Name: aSpecies, Kind: kindUplift,
				Characteristics: map[string]string{"SOC": "2d6"},
			},
			want: "not one of the six characteristics",
		},
		{
			name:    "a ceiling outside the range",
			species: setting.Species{Name: aSpecies, Kind: kindUplift, Maximum: 400},
			want:    "maximum 400",
		},
		{
			name: "rolls with no age ranges to match",
			species: setting.Species{
				Name: aSpecies, Kind: kindUplift, YouthRolls: 2,
				YouthAges: []string{"ages 2-4"},
			},
			want: "a roll represents a range",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			broken := minimal()

			broken["species"] = []any{tc.species}

			_, err := setting.Load(write(t, broken))
			if err == nil {
				t.Fatal("the file validated")
			}

			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error does not mention %q:\n%v", tc.want, err)
			}
		})
	}
}
