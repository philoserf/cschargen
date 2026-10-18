# WALKTHROUGH.md

How `cschargen` runs, from the command line to the written record, in the order the
code executes.

This is the linear tour. For _why_ the system is shaped this way — the invariants, the
Product Identity constraint, what would break if you changed something — read
[`THEORY.md`](THEORY.md). This one assumes no knowledge of it, and starts at `main`.

Regenerated for `v0.1.0-alpha.7`. It quotes the code by file and symbol, never by
line, so it rots the moment one is renamed; `CLAUDE.md` records that it is brought
current at each release.

---

## 1. Overview

`cschargen` is a single Go binary, no dependencies outside the standard library
(`go.mod` declares the module and a Go version and nothing else). It generates
characters for the _Clement Sector Core Character Creation Book_ by walking the book's
twenty-step lifepath and writing down everything it did.

Six commands:

`cmd/cschargen/main.go` — `commands`

```go
const commands = `commands:
  new       generate a character
  batch     generate several, as one JSON record per line
  data      validate a setting data file
  render    turn a record into a character sheet, or its lifepath
  replay    re-run a record from its seed and recorded choices
  version   report the build and the versions a record stamps`
```

The unit of work is a **record**: one JSON document holding the provenance, the
character, and an ordered log of every step, throw, choice and consequence. `render`
turns a record into Markdown; `replay` re-runs it from its seed and checks that the
current build still produces it. Every other command produces records.

---

## 2. Architecture

Six packages, and the dependency arrows all point the same way:

```
cmd/cschargen  ──►  chargen  ──►  career
                       │            │
                       ├──►  dice ◄─┘  (denied by depguard)
                       └──►  setting
render         ──►  chargen
```

| Package         | What it holds                                                                            |
| --------------- | ---------------------------------------------------------------------------------------- |
| `dice`          | A seeded stream and the shapes of throw: 2d6, d66, d100, 3d6-drop-lowest                 |
| `career`        | Thirty-four career definitions plus the shared Injury and Life Events tables, hand-typed |
| `setting`       | The world/subsector/species schema, validator and loader                                 |
| `chargen`       | The engine — the record, the event log, the `Decider`, and the rules                     |
| `render`        | Record → Markdown sheet, lifepath transcript, or roster line                             |
| `cmd/cschargen` | Flags, subcommands, exit statuses                                                        |

Most of the layering is enforced by the compiler, because the reverse edge would be an
import cycle. The one edge that would otherwise compile is held by a linter:

`.golangci.yml` — `depguard`

```yaml
    depguard:
      # depguard's own default is a deny-all list that flags this module's
      # own imports, so the choice was never on or off but configured or
      # meaningless. Most of the layering is enforced by the compiler,
      # because the reverse edge is an import cycle. These are the edges
      # that would otherwise compile.
      rules:
        dice-knows-no-rules:
          files: ["**/dice/*.go"]
          deny:
            - pkg: github.com/philoserf/cschargen/career
              desc: dice knows the shapes of throw, never what a throw means
```

One thing to know before reading further: **the worlds are not in this repository.**
The book declares world, subsector and organization names as Product Identity, and
Steps 3–4 of character creation are tables of world names. So `setting` is a schema
with no data — the user supplies a file transcribed from their own copy. A small
invented sample (`setting/sample.json`) is embedded so the tests and the demo run.

---

## 3. Entry: `main` and the exit contract

`cmd/cschargen/main.go` — `main`

```go
func main() {
	err := run(os.Args[1:], os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

Everything real is in `run`, which takes its arguments and its output stream as
parameters — that is what makes the whole CLI testable in-process, and
`cmd/cschargen/main_test.go` drives it through one helper, `capture`, rather than
through a built binary.

Note what the exit status does **not** encode. Both kinds of failure exit 1:

`cmd/cschargen/main.go` — package doc

```go
// The exit contract: a misused command line exits 1 with a message
// beginning "usage:", and any other failure exits 1 with a message that
// does not. The two are told apart by the message rather than the status,
// so a script wrapping the tool can distinguish its own bug from the
// engine's.
```

`errUsage` carries that prefix and `usagef` wraps it. Dispatch is a map:

`cmd/cschargen/main.go` — `dispatch`

```go
func dispatch() map[string]func([]string, *os.File) error {
	return map[string]func([]string, *os.File) error{
		"new":    newCommand,
		"batch":  batchCommand,
		"data":   dataCommand,
		"render": renderCommand,
		"replay": replayCommand,
	}
}
```

`version` and the help forms are handled by a `switch` after the map misses, because
they take no arguments.

The three commands that take a file — `render`, `replay`, `data validate` — read
their command lines through one function, because Go's `flag` package stops at the
first argument that is not a flag:

`cmd/cschargen/main.go` — `parseOperands`

```go
	for {
		err := flags.Parse(args)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		rest := flags.Args()
		if len(rest) == 0 {
			return operands, nil
		}

		if consumed := len(args) - len(rest); consumed > 0 && args[consumed-1] == "--" {
			return append(operands, rest...), nil
		}

		operands = append(operands, rest[0])
		args = rest[1:]
	}
```

It parses again after each operand, so `replay character.json --data setting.json`
is one record and a flag rather than three records; a `--` still ends the flags.

---

## 4. `new`: from flags to a `Generator`

`newCommand` is the spine of the program. Follow it in order.

### 4.1 The seed

`cmd/cschargen/new.go` — `newFlags.parse`

```go
	if given(f.set, "seed") {
		return *f.seed, nil
	}

	return rand.Uint64(), nil
```

A seed given is used; a seed omitted is _chosen here and recorded_, so a character
generated without a seed is still reproducible afterwards. This is the first place the
program's central bargain shows up: nothing is allowed to happen that the record
cannot explain.

`given` walks `flag.FlagSet.Visit`, which visits only the flags actually given. That
matters for one flag in particular:

`cmd/cschargen/new.go` — `newFlags.termLimit`

```go
func (f newFlags) termLimit() int {
	if given(f.set, "terms") && *f.terms == 0 {
		return noTerms
	}

	return *f.terms
}
```

The engine reads `TermLimit` as zero for "not given, use the policy default" and
negative for "no career terms at all". A plain `int` cannot tell an unset flag from one
set to zero, so `--terms 0` used to silently mean the policy's four. `Visit` is what
lets the command tell them apart and hand the engine the sentinel it already
understands — `noTerms`, which is -1.

The engine reads _any_ negative as "none", so the other edge is held by refusing it:

`cmd/cschargen/new.go` — `checkTerms`

```go
func checkTerms(terms int) error {
	if terms < noTerms {
		return usagef("--terms %d; 0 serves none, a positive number serves that many,"+
			" and omitted is the policy default", terms)
	}

	return nil
}
```

`--terms -5` used to serve no terms and write `termLimit: -5` into the record — a
value that is neither of the command's two spellings of "none". Normalising it would
write -1 where -5 was asked for, and `Inputs` is what was asked for, so it is refused.

### 4.2 The setting data

`cmd/cschargen/new.go` — `loadSetting`

```go
func loadSetting(path string) (*setting.Data, error) {
	if path == "" {
		data, err := setting.Sample()
		// ...
		return data, nil
	}

	data, err := setting.Load(path)
	// ...
}
```

Both paths end in `setting.parse`, which does three things in a fixed order:

`setting/load.go` — `parse`

```go
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&data)
	// ...
	problems := Validate(&data)
	if len(problems) > 0 {
		return nil, &InvalidError{Source: name, Problems: problems}
	}

	sum := sha256.Sum256(content)

	data.Hash = hex.EncodeToString(sum[:])
```

`DisallowUnknownFields` turns a misspelt key into an error rather than a zero value —
which for `maximumTerms` would mean a character allowed no terms at all. Validation is
complete rather than first-failure, because someone transcribing sixty worlds should be
told about all the problems at once. And the **content hash** is computed here: it goes
into the record, and `replay` refuses a record whose hash does not match the file it is
handed.

### 4.3 Names are checked before anything is generated

`cmd/cschargen/warn.go` — `checkNames`

```go
func checkNames(world *setting.Data, flags newFlags) error {
	err := checkOrigin(world, *flags.subsector, *flags.homeworld)
	if err != nil {
		return err
	}

	*flags.forceCar, err = checkCareer(*flags.forceCar)
	if err != nil {
		return err
	}

	err = checkTerms(*flags.terms)
	if err != nil {
		return err
	}

	return checkWorldOverrides(flags)
}
```

The engine would refuse a bad subsector or career too — but partway through Step 3,
where it reads as a generation failure rather than as the typo it is. So every flag
that names something is held against the data here, and the error lists what the data
_does_ have. `checkCareer` also rewrites the name to the book's own capitalisation, so
later messages spell it the way every other line does.

### 4.4 Inputs, and the one field the decider does not decide

`cmd/cschargen/new.go` — `newFlags.inputs`

```go
		// Recorded because the engine reads it back: Steps 3 and 4 offer
		// a choice the policy is not asked to make, and a replay has to
		// reach the same choice points the run did. Deciding it from the
		// decider would not survive replay, whose decider is neither the
		// player nor the policy.
		Interactive: !*f.auto,
```

This is worth pausing on, because it is the shape of a bug the project has already hit.
A step that behaves differently for an interactive run cannot ask the _decider_ whether
it is interactive — during replay the decider is `Replay`, which is neither. So the
answer travels in the record.

### 4.5 Choosing a decider

`cmd/cschargen/new.go` — `decider`

```go
func decider(auto bool) chargen.Decider {
	if auto {
		return chargen.Policy{}
	}

	return chargen.NewPlayer(os.Stdin, os.Stderr)
}
```

The player reads stdin and prompts to **stderr**, because the record goes to stdout so
it can be piped. A record generated interactively is byte-identical in form to an auto
one and replays the same way.

`batch` does not go through `decider`. Each member gets a `BatchPolicy` seeded from
its own seed, which is `Policy` everywhere except the choice points named in one set:

`chargen/batchpolicy.go` — `BatchPolicy.Choose`

```go
	if !batchDrawn[c.Point] {
		return Policy{}.Choose(c)
	}

	return b.draw.IntN(len(c.Options)), nil
```

`batchDrawn` holds thirteen points — the career, the assignment, the skill table, the
youth and teenage paths, the specialty and degree-field choices among them — where
the book offers the options on equal terms, so a draw makes a different character
rather than a worse one. The points with an argued answer stay `Policy`'s: a loss
goes on the highest score, an aging crisis is paid for. The draw has its own PCG
stream rather than the dice's, because a decider that pulled from the dice would
shift every throw after its first answer. `chargen/policy_internal_test.go` —
`TestPolicyDocumentsWhatABatchDraws` holds the set to `POLICY.md`'s rows in both
directions.

### 4.6 Handing off

`cmd/cschargen/new.go` — `newCommand`

```go
	character, err := chargen.New(chargen.Options{
		Seed:          seed,
		Decider:       decider(*flags.auto),
		EngineVersion: version(),
		PolicyVersion: chargen.PolicyVersion,
		Setting:       world,
		Inputs:        inputs,
	}).Run()
```

Two notes can follow on stderr before the record is written:

`cmd/cschargen/new.go` — `newCommand`

```go
	warnIfTheCareerChanged(os.Stderr, character, *flags.forceCar)
	warnIfTheHomeworldChanged(os.Stderr, character)
```

Both are cases where the rules overrode a flag and the record already says so: a
`--career` the character failed to enlist in (p. 110), and a `--homeworld` that does
not admit the character's species, which sends them to another world (p. 42). The
command says it once and still exits 0, because nothing failed — but a referee who
asked for a world should not have to render the history to learn the character was
born somewhere else. `batch` counts each in one line rather than repeating it.

---

## 5. The engine

### 5.1 Provenance is stamped at construction, not at the end

`chargen/generator.go` — `New`

```go
		char: &Character{
			Provenance: Provenance{
				SchemaVersion: SchemaVersion,
				Ruleset:       Ruleset,
				EngineVersion: opts.EngineVersion,
				PolicyVersion: opts.PolicyVersion,
				RNG:           RNG{Algorithm: "math/rand/v2 PCG", Seed: opts.Seed},
				SettingData:   describeSetting(opts.Setting),
				Inputs:        opts.Inputs,
			},
		},
```

A record that fails halfway still knows what it was trying to be.

### 5.2 `Run` is the twenty steps, in the book's order

`chargen/generator.go` — `Run`

```go
	err := g.chooseSpecies()
	// ...
	// The genetics decide which characteristic method Step 2 uses, so they
	// come before it although the book prints them at Step 5 (p. 62).
	err = g.determineGenetics()
	// ...
	err = g.rollCharacteristics()
	// ...
	err = g.determineOrigin()
	// ...
	err = g.runCareers()
	// ...
	err = g.finishingTouches()
	// ...
	g.char.Events = g.log.Events()
```

Six calls. The one departure from the printed order is commented where it happens: the
genetics of Step 5 decide which characteristic method Step 2 uses, so they run first.

Running first is only half of it. A hybrid or compound outcome can also adjust a
characteristic — "take a -1 to their DEX" — and when the genetics run, every score is
still 0. So `geneticOutcome` decides the outcome and holds it rather than applying it:

`chargen/genetics.go` — `geneticOutcome`

```go
	g.genetic = outcome
	g.geneticThrow = throw
```

and Step 2 applies it once `assignCharacteristics` has put the scores in place, with
the genetics throw still the cause:

`chargen/genetics.go` — `applyGeneticAdjustment`

```go
func (g *Generator) applyGeneticAdjustment() {
	defer g.citing(geneticsCite)()

	for _, name := range slices.Sorted(maps.Keys(g.genetic.Adjust)) {
		g.adjust(name, g.genetic.Adjust[name], g.genetic.Detail, g.geneticThrow)
	}
}
```

It used to apply where it was decided, moving a 0 that the roll then overwrote, so a
record logged a DEX +1 the sheet did not have. The keys are sorted because a Go map's
order is not, and an outcome adjusting two characteristics would otherwise log them in
a different order from one run to the next.

### 5.3 Every roll goes through `dice`, and every roll is logged

`dice/throw.go` — `Dice.Throw`

```go
func (d *Dice) Throw(target int, mods ...Mod) Throw {
	roll := d.TwoD6()

	total := roll.Total
	for _, m := range mods {
		total += m.Value
	}

	return Throw{
		Roll:    roll,
		Mods:    mods,
		Target:  target,
		Total:   total,
		Success: total >= target,
	}
}
```

`Throw.Natural()` returns `Roll.Total` — the dice _before_ modifiers — because several
rules turn on a natural 12 rather than a modified one.

`dice` knows nothing about what any of this means. The shapes are named at the call
site so the rule is legible:

`dice/dice.go` — `Dice.D66`

```go
func (d *Dice) D66() Roll {
	tens, units := d.die(), d.die()

	return Roll{Expr: "d66", Dice: []int{tens, units}, Total: tens*10 + units}
}
```

**d66 is one throw, not two d6 rolls.** Results run 11–16, 21–26 … 61–66 — thirty-six
of them, which is what every career's event table is indexed by. Reading it as 2d6
would collapse those thirty-six onto eleven.

### 5.4 Every decision goes through the `Decider`

`chargen/generator.go` — `Generator.choose`

```go
func (g *Generator) choose(ask Choice) (int, error) {
	chosen, err := g.decider.Choose(ask)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", ask.Point, err)
	}

	if chosen < 0 || chosen >= len(ask.Options) {
		return 0, ErrChoiceOutOfRange
	}

	g.log.Choice(ChoiceEvent{
		Decider: g.decider.Kind(),
		Point:   ask.Point,
		Prompt:  ask.Prompt,
		Options: slices.Clone(ask.Options),
		Chosen:  chosen,
		Cite:    ask.Cite,
	})

	return chosen, nil
}
```

One funnel, and the answer is logged with the full option list. The list is _copied_
into the log, not shared with the caller: `medicalDegree` narrows its own list of Medic
specialties in place after each of its four picks, and a shared slice let it rewrite
what all four recorded choices said had been offered — so a medical graduate's record
named specialties they never held and would not replay. Copying here closes the class
rather than the one caller. That list is what makes replay possible:

`chargen/decider.go` — `Replay.Choose`

```go
	if rec.Point != ask.Point {
		return 0, fmt.Errorf("%w: choice %d is %q, engine asked %q",
			ErrReplayDiverged, r.next, rec.Point, ask.Point)
	}

	if !slices.Equal(rec.Options, ask.Options) {
		return 0, fmt.Errorf("%w: choice %d (%s) recorded options %v, engine offered %v",
			ErrReplayDiverged, r.next, ask.Point, rec.Options, ask.Options)
	}
```

A recorded index means nothing against a different list. A record made with `--career`
forced saw a one-entry career list; replayed against the full list, index 0 would name
a different career. The options check catches that.

The auto decider is nearly one line of logic:

`chargen/decider.go` — `Policy.Choose`

```go
func (Policy) Choose(c Choice) (int, error) {
	if len(c.Options) == 0 {
		return 0, fmt.Errorf("%w: %s", ErrNoOptions, c.Point)
	}

	if len(c.Scores) != len(c.Options) {
		return 0, nil
	}

	highest := 0

	for i, score := range c.Scores {
		if score > c.Scores[highest] {
			highest = i
		}
	}

	return highest, nil
}
```

The first option, unless the engine has scored the options. `POLICY.md` argues "first
listed" is a rule rather than a preference, because every table in the book prints its
options in a fixed order — so it is checkable against the page. The one exception is a
loss the character places: the Injury table prints STR, DEX, END, so "first listed"
took every such loss from STR, and a character injured three times in twelve terms
reached STR 0 while DEX and END sat untouched — the outcome p. 14 asks a Referee to
head off. The engine fills `Scores` only where every option opens with a loss from a
different characteristic, with the score each would take from; the policy takes the
highest, and a tie goes to the first printed. Every other choice leaves `Scores` empty
and gets the first option, as before. That rule has stood since policy 0.4.0; 0.5.0
changed only what `batch` draws (§4.5).

### 5.5 The event log

`chargen/event.go` — `Log.append`

```go
func (l *Log) append(e Event) int {
	e.Seq = len(l.events) + 1
	l.events = append(l.events, e)

	return e.Seq
}
```

Sequence numbers start at 1 and increase by 1. **`append` returns the number**, which is
how the caller names the event it just wrote as the `Cause` of what follows. That is
why `step` and `cause` integers are threaded through nearly every function in the
engine: they are the causal spine of the record, not parameter noise.

Four event kinds — step, throw, choice, consequence — as a discriminated union:

`chargen/event.go` — `Event`

```go
type Event struct {
	Seq  int       `json:"seq"`
	Kind EventKind `json:"kind"`

	Step        *StepEvent        `json:"step,omitempty"`
	Throw       *ThrowEvent       `json:"throw,omitempty"`
	Choice      *ChoiceEvent      `json:"choice,omitempty"`
	Consequence *ConsequenceEvent `json:"consequence,omitempty"`
}
```

Exactly one payload is non-nil. Note `ThrowEvent.Target` and `.Success` are pointers:

`chargen/event.go` — `ThrowEvent`

```go
// ThrowEvent records one throw. A throw resolved against no target -- the
// six characteristic rolls of p. 13, say -- leaves Target and Success nil,
// because zero is a target no throw of dice satisfies and recording one
// would read as a check that always failed.
```

`Log.Roll` writes the first kind, `Log.Throw` the second.

### 5.6 Steps 3 and 4: where the setting data is read

`chargen/origin.go` — `determineOrigin`

```go
	sub, err := g.chooseSubsector()
	// ...
	// The heading comes before the work it heads. It used to be written
	// after, which put Step 4's own d100 under Step 3 in the transcript.
	step := g.log.Step("Step 4: Determine Homeworld", "p. 39")

	world, err := g.chooseHomeworld(sub, step)
	// ...
	world, err = g.worldThatAdmitsThem(world, sub, step)
	// ...
	err = g.settleOn(world, sub, step, "born there")
```

`worldThatAdmitsThem` is the step you would not guess at: a world may bar an engineered
species by name, or admit them as property rather than as people. A character whose
rolled homeworld will not have them is moved rather than re-rolled, because a re-roll
can land on the same world forever.

One subtlety visible in the transcript below: the subsector chart is a **1d6**, and a
setting file may name fewer than six subsectors. A throw landing on no entry is thrown
again, which keeps the chart's weighting — turning the miss into a choice would hand it
to the policy, which answers with the first subsector the file lists.

### 5.7 Steps 5–8, in one call

`chargen/origin.go` — `earlyLife`

```go
	if g.char.Provenance.Inputs.SkipFamily {
		restore := g.citing("p. 57")

		g.consequence(ConsequenceFamily, step,
			"Step 5 skipped at the player's request (p. 57)", "")

		restore()
	} else {
		err := g.determineFamily(world)
		// ...
	}

	err := g.youthEvents()
	// ...
	err = g.teenageEvents()
	// ...
	return g.higherEducation()
```

Each of the four is optional in the book, and each `--skip-*` flag writes a consequence
saying so. A skipped step is recorded as skipped, not omitted.

`g.citing` is how a consequence gets its page. A step is not one page — Step 4 writes
results from five — so the page is set by the function whose work it is and `restore`
puts the previous one back on return. The pattern recurs throughout the engine as
`defer g.citing(...)()`.

### 5.8 The career loop

`chargen/generator.go` — `runCareers`

```go
	for len(g.char.State.Terms) < g.terms() {
		if g.stopped {
			return nil
		}

		if g.service.career == nil {
			err := g.enterCareer()
			// ...
			if g.stopped {
				return nil
			}

			if g.service.career == nil {
				// The enlistment failed. That consumed no term, so the
				// loop tries again -- three failures in a row is what
				// sends a character to Vagabond, not one.
				continue
			}
		}

		err := g.serveTerm()
		// ...
		err = g.nextTerm()
		// ...
	}

	// A character who reaches the term limit is still in a career, and
	// leaving it is what mustering out is for.
	return g.leaveCareer(g.log.Len(), "character generation ended")
```

The loop bound has two ceilings:

`chargen/generator.go` — `Generator.terms`

```go
func (g *Generator) terms() int {
	if g.homeworldTerms > 0 && g.homeworldTerms < g.termLimit {
		return g.homeworldTerms
	}

	return g.termLimit
}
```

The policy's `--terms` and the homeworld's own maximum are both caps and the lower
wins (p. 125).

**Entering a career** — `chargen/term.go`, `enterCareer` → `chooseCareer` → `enlist` →
`beginService`. `chooseCareer` has four ways to answer before it ever offers a choice:
a pending transfer a table result ordered, an enslaved character's forced first career,
three consecutive failed enlistments, or an empty eligible list. Only then does it ask.

Before the list is built, a standing restriction can close it to one career:

`chargen/term.go` — `onlyVagabond`

```go
func (g *Generator) onlyVagabond() bool {
	return g.crisisSurvived || g.mentalDecline || g.reachedOne
}
```

Three flags set long before Step 9 reads them. An Aging Crisis survived and a mental
characteristic at 0 are pp. 123–124's; `reachedOne` is p. 119's, set by
`noteReachedOne` in `chargen/effects.go` the first time a loss brings a physical
characteristic, INT or EDU to 1 or below. It never clears — a score that recovers
reopens nothing (ERRATA E-46) — and `takePendingTransfer` honours it too, so a mishap
that would send the character to Prisoner sends them to Vagabond instead.

`chargen/term.go` — `eligibleCareers`

```go
	for _, def := range career.All() {
		if def.Name == "Prisoner" || def.Name == career.Slave().Name {
			continue
		}

		if g.forced != "" && def.Name != g.forced {
			continue
		}

		if until, locked := g.lockout[def.Name]; locked && len(g.char.State.Terms) < until {
			continue
		}

		eligible = append(eligible, def)
	}
```

Two careers are never in the list — neither has an enlistment throw, so either one left
in would be entered by whoever drew it. A career that turned the character down is
locked out for two full terms (p. 110).

**Serving a term** — `chargen/serve.go`, `serveTerm`:

`chargen/serve.go` — `resolveTerm`

```go
func (g *Generator) resolveTerm(assignment career.Assignment, survived bool) error {
	if !survived {
		g.log.Step("Step 13: Roll A Mishap", "p. 113")

		return g.rollMishap(true)
	}

	err := g.advance(assignment)
	// ...
	err = g.rollSkill(assignment)
	// ...
	return g.rollEvent()
}
```

Survival decides everything. Survive and you get Steps 14, 15 and 16 — advancement, a
skill, an event. Fail and you get Step 13 alone.

`chargen/serve.go` — `rollEvent`

```go
func (g *Generator) rollEvent() error {
	g.log.Step("Step 16: Roll for Events", "p. 118")

	roll := g.dice.D66()
	cause := g.log.Roll(roll, "p. 118")

	row, ok := g.service.career.Events[roll.Total]
	if !ok {
		return ErrMissingEventRow
	}

	g.consequence(ConsequenceCareer, cause, "event: "+row.Summary, g.service.career.Name)

	return g.applyAll(row.Effects, cause)
}
```

`Events` is a `map[int]` keyed by the d66 result, which is why the results are
11–66 with gaps rather than 1–36.

After the term, Step 17 ages the character, and it is where §5.7's `g.citing` earns
its keep. The step runs with the career's page still set from the term, and nothing
it writes is the career's:

`chargen/serve.go` — `age`

```go
	// The career's page is still set from the term, and nothing this step
	// writes is the career's. Each function below names its own page.
	defer g.citing(ageCite)()
```

Each function beneath it names its own page in turn — `stampApparentAge` p. 125,
`agingThrows` pp. 122-123, and inside that `agingCrisis` p. 123 and
`settleTerminalStates` pp. 123-124.
Before they did, an aging loss cited the career just served, one line below a throw
citing the aging table.

### 5.9 Effects: the closed alphabet and the fold over it

A table row does not describe what happens — it _is_ a list of `career.Effect` values,
and the engine folds over them:

`chargen/effects.go` — `applyAll`

```go
func (g *Generator) applyAll(effects []career.Effect, cause int) error {
	for _, effect := range effects {
		err := g.apply(effect, cause)
		if err != nil {
			return err
		}
	}

	return nil
}
```

`chargen/effects.go` — `apply`

```go
// The effect vocabulary is a closed alphabet and this is the fold over it.
// A fold has the size of the alphabet it folds over, and EffectGroup and
// EffectSubTable add letters that recurse rather than branch.
//
// Splitting it would cost more than it saved. exhaustive checks a switch,
// and one switch over EffectKind is what guarantees a new kind is handled
// somewhere; two switches need a default between them, and a default is
// where a forgotten kind goes to be ignored quietly.
//
//nolint:cyclop,funlen,gocyclo,maintidx // see the paragraph above
func (g *Generator) apply(effect career.Effect, cause int) error {
	// career/ builds table data and never sees a character, so a reading
	// implemented there declares itself on the effect and is recorded
	// here (#154).
	if effect.Deviation != "" {
		g.char.Provenance.Deviate(effect.Deviation)
	}

	switch effect.Kind {
	case career.EffectSkill:
		return g.applySkill(effect, cause)
	case career.EffectCharacteristic:
		return g.adjustBy(effect, cause)
	case career.EffectChoice:
		return g.applyChoice(effect)
	case career.EffectCheck:
		return g.applyCheck(effect)
	...
```

This is a big switch on purpose. `.golangci.yml` sets
`exhaustive.default-signifies-exhaustive: false`, and the switch has no `default` — so
adding a kind to `career` and forgetting it here is a **build failure**. Two smaller
switches would need a default between them, and a default is where a forgotten kind
goes to be ignored quietly. The `nolint` is bought with that argument.

The lines before the switch are how a reading implemented in `career/` reaches a
record: that package never sees a character, so an effect names the `ERRATA.md` entry
it rests on and the engine stamps it when it applies the value.

`EffectCheck` and `EffectGroup` recurse: a check's success and failure branches are
themselves effect lists, so "roll 8+; on a success, roll again" nests without a special
case.

The commonest letter is a skill, and two cases in it are worth following: "(Any)", and
a skill the character does not hold yet.

`chargen/effects.go` — `grantSkill`

```go
	if isAny(effect) {
		return g.grantAny(effect, level, cause)
	}

	// p. 118, for a skill table's result: "If the character does not
	// already have the skill, they should add the skill at level 0", and
	// where it names a specialty, "it is added to the character sheet as a
	// zero level skill without a specialty". Any specialty counts as having
	// the skill, which is the reading grantAny makes of the same sentence.
	if g.onSkillTable && !g.char.State.Has(effect.Skill) {
		return g.gainSkill(effect.Skill, "", 0, cause), nil
	}
```

`g.onSkillTable` is set while Step 15's roll, or a result that rolls on a named
table, is being applied; an event or a degree that grants a skill leaves it clear and
grants the level it prints.

`chargen/skills.go` — `grantAny`

```go
	held := g.char.State.Has(effect.Skill)

	if level == 0 || (g.onSkillTable && !held) {
		if held && level == 0 {
			g.consequence(ConsequenceSkill, cause,
				effect.Detail+": the character already holds "+effect.Skill, "")

			return Skill{Name: effect.Skill}, nil
		}

		return g.gainSkill(effect.Skill, "", 0, cause), nil
	}

	options := g.specialtiesOf(effect.Skill)
	// ...
	specialty, err := g.chooseAny(effect.Skill, options, skillStepCite)
```

p. 118 says what a skill table's result does for a character without the skill: it
goes on bare at level 0, whatever form the result takes — a plain skill, a named
specialty, or "(Any)". The branch above applies that to the first two; until alpha.7
only "(Any)" followed it, and a first Persuade from a table went on at 1. For a
character who holds the skill, p. 118's "(Any)" raises any specialty of it, held or
not. `grantAny` does both halves, and `chooseAny` puts the specialty to the decider as
a `skill_specialty` choice, the same point a printed list of specialties uses. The engine used to record
the word "Any" as if it were a specialty, which put `Advocate (Any)-2` on a sheet and
made p. 116's rank-benefit proviso — a specialty not already held at level 1 — find
"Any" already held and grant nothing. Away from a skill table and a rank row the book
gives "(Any)" no rule of its own; reading it as p. 118's is ERRATA E-45, stamped where
it applies.

The setting file writes "(Any)" too — a world's background skill, a species' starting
skill — and those go through `chooseSpecialty` in `chargen/background.go`, which
reaches the same list:

`chargen/background.go` — `chooseSpecialty`

```go
	if len(options) == 1 && options[0] == anySpecialty {
		g.char.Provenance.Deviate("E-45")

		options = g.specialtiesOf(alternative.Skill)
	}
```

It used to take a one-item list as the answer, so `["Any"]` put `Electronics (Any)` on
a sheet. The sample's one "(Any)" background skill is the second of a pair, which the
first-option policy never took, and it surfaced when `batch` began drawing that choice.

Every grant ends in `State.GainSkill`, which holds Jack of All Trades at p. 309's
printed maximum of 3; `noteCeiling` in `chargen/effects.go` logs the levels a grant
past it lost, cited to p. 309, and stamps ERRATA E-48.

The careers themselves are hand-typed Go, built through constructors that keep a
transcribed table close to the page:

`career/build.go` — `skill`

```go
func skill(name string, specialties ...string) Effect {
	detail := "gain a level in " + name

	return Effect{
		Kind:        EffectSkill,
		Detail:      detail,
		Skill:       name,
		Specialties: specialties,
		Level:       1,
	}
}
```

### 5.10 Leaving a career, and Step 19

`chargen/generator.go` — `leaveCareer`

```go
	g.consequence(ConsequenceCareer, cause, "left the "+g.service.career.Name+" career: "+why, g.service.career.Name)

	err := g.musterOut(*g.service.career, g.service.termsInCareer)
	// ...
	g.service = serviceState{}

	// The career's page ends with the career. Steps 9 and 10 set their own
	// now, so nothing downstream reads what this used to leave behind.
	g.cite = ""

	// Not in serviceState, for the reason its doc comment gives: a pending
	// transfer sets this before the career begins.
	g.forcedTerms = 0

	g.dropCareerModifiers()
	g.dropCareerPools()
```

Leaving a career is where the per-career state is torn down. This used to be seven
assignments to seven `Generator` fields, initialised by one function and cleared by
a different one — so a field added to the first and forgotten in the second leaked
into the next career, silently. They are one `serviceState` value now, and clearing
it is one line that cannot be partially written. `forcedTerms` is deliberately not
in it: a pending transfer sets it _before_ the career begins, so it does not live
exactly as long as one.

`g.cite = ""` is the same argument for the page a consequence carries. It used not
to be cleared, and Steps 9 and 10 wrote their consequences under whatever table had
been read last — a character "accepted into Scavenger" on a graduate school's
pages. The right failure mode is an empty cite rather than a confident wrong one.

`dropCareerModifiers` exists because a modifier granted "for the rest of this
career" carries no count — nothing else would ever end it, and it would follow the
character for good.

`chargen/musterout.go` — `musterOut`

```go
	if total <= 0 {
		// The queue is emptied even with nothing to spend. A mishap that
		// removed two rolls leaves a negative batch behind, and leaving it
		// there would charge the character for that mishap a second time
		// the next time they joined this career.
		g.clearBenefits(left.Name)
		// ...
	}
```

Benefits are a **queue of batches**, not a count, because the book grants both "a +1
modifier to all Benefit rolls made in this career" and "three Benefit rolls at +1", and
those are different things.

---

## 6. A real run

Transcript of commands run while writing this document, against the embedded sample
data, with a binary built from the release branch. Nothing re-runs these.

```console
$ cschargen version
cschargen v0.1.0-alpha.6.0.20261018204727-77fb95fde823+dirty
schema 2
policy 0.5.0
ruleset Clement Sector Core Character Creation Book, (c) 2026 Independence Games

$ cschargen new --auto --seed 42 --terms 2 -o c42.json
$ head -29 c42.json
{
  "provenance": {
    "schemaVersion": 2,
    "ruleset": "Clement Sector Core Character Creation Book, (c) 2026 Independence Games",
    "engineVersion": "v0.1.0-alpha.6.0.20261018204727-77fb95fde823+dirty",
    "policyVersion": "0.5.0",
    "rng": {
      "algorithm": "math/rand/v2 PCG",
      "seed": 42
    },
    "settingData": {
      "name": "sample",
      "hash": "37f77096e3b28c00b026e2b675657fb7b924077f0136e7cdbd37239f0ebbcc67",
      "sample": true
    },
    "inputs": {
      "species": "human",
      "techLevel": 0,
      "maxTerms": 0,
      "termLimit": 2,
      "interactive": false
    },
    "deviations": [
      "E-44",
      "E-21",
      "E-8",
      "E-20"
    ]
  },
```

Every field there is a promise the record makes about itself. `inputs` is what was
asked for, zeros included, because a replay reconstructs the run from it. `deviations`
is the `ERRATA.md` readings that changed this character,: E-44 truncates p. 61's family recursion at one generation, E-21 and E-20 are
readings of the apparent-age chart, and E-8 is "rank benefits are applied immediately
upon achieving that Rank", which fires on entering any career because every character
begins at Rank 0. A reading that applies to every character alike is not stamped.

Now the log, rendered:

```console
$ cschargen render --history c42.json
# Lifepath: (unnamed)

Seed 42, 152 events.


## Step 1: Choose Human, Engineered Human, or Uplift (p. 21)

-    2  -> a baseline human [species, from 1]

## Step 2: Roll Characteristics (p. 39)

-    4  throw  3d6 drop lowest [3 6 5] = 11  (p. 13)
-    5  throw  3d6 drop lowest [3 6 6] = 12  (p. 13)
-    6  throw  3d6 drop lowest [2 3 2] = 5  (p. 13)
-    7  throw  3d6 drop lowest [6 4 5] = 11  (p. 13)
-    8  throw  3d6 drop lowest [5 4 6] = 11  (p. 13)
-    9  throw  3d6 drop lowest [4 3 5] = 9  (p. 13)
-   10  choice Assign a rolled score to STR: 11 (policy)
-   11  -> STR 11 [characteristic, from 10]
-   12  choice Assign a rolled score to DEX: 12 (policy)
-   13  -> DEX 12 [characteristic, from 12]
-   14  choice Assign a rolled score to END: 5 (policy)
-   15  -> END 5 [characteristic, from 14]
-   16  choice Assign a rolled score to INT: 11 (policy)
-   17  -> INT 11 [characteristic, from 16]
-   18  choice Assign a rolled score to EDU: 11 (policy)
...
## Step 3: Determine Subsector of Origin (p. 39)

-   23  throw  1d6 [4] = 4  (p. 39)
-   24  -> no subsector claims that result; throw again (p. 39) [homeworld, from 23]
-   25  throw  1d6 [5] = 5  (p. 39)
-   26  -> no subsector claims that result; throw again (p. 39) [homeworld, from 25]
-   27  throw  1d6 [2] = 2  (p. 39)
-   28  -> born in Kestrel Reach [homeworld, from 27]
```

Read that against the code. Events 4–9 are `Log.Roll` — six throws with no target, so
no target is recorded. Events 10–21 are twelve events for six assignments: a
`ChoiceEvent` and then the `ConsequenceEvent` it caused, each consequence naming the
sequence number of its choice (`from 10`, `from 12`). Events 23–27 are the 1d6 subsector
chart missing twice on a four-subsector file and being thrown again, exactly as §5.6
described. The blank line doubled under the header is the renderer's, verbatim: the
header ends in a blank line and every step heading begins with one.

Step 1's heading is the engine's own string, and it reads "Engineered Human" because
the book's word for a genetically engineered human is one the OGL notice reserves as
Product Identity. It used to be the book's word, which put a reserved term in every
record and every rendered sheet; `setting/setting_test.go`'s
`TestNoProductIdentityIsCommitted` now reads every file git tracks for it.

Finally, the claim the whole design exists to support:

```console
$ cschargen replay c42.json
replayed 152 events from seed 42: identical
```

---

## 7. `render`: the record is the only input

`render/sheet.go` — package doc

```go
// The JSON record is the source of truth. Everything here is a view of it,
// so a render never computes a rule -- if a number is not in the record,
// it does not appear on the sheet.
```

Three views. `Sheet` is modelled on the book's Long Form Character Sheet;
`Roster` is one line per character for casting from a batch; `Transcript` walks the log:

`render/transcript.go` — `writeEvent`

```go
func writeEvent(out *strings.Builder, event chargen.Event) {
	switch event.Kind {
	case chargen.EventStep:
		fmt.Fprintf(out, "\n## %s (%s)\n\n", event.Step.Name, event.Step.Cite)
	case chargen.EventThrow:
		fmt.Fprintf(out, "- %4d  throw  %s\n", event.Seq, throwLine(event.Throw))
	case chargen.EventChoice:
		fmt.Fprintf(out, "- %4d  choice %s: %s (%s)\n",
			event.Seq, event.Choice.Prompt, event.Choice.Options[event.Choice.Chosen],
			event.Choice.Decider)
	case chargen.EventConsequence:
		fmt.Fprintf(out, "- %4d  -> %s [%s, from %d]\n",
			event.Seq, event.Consequence.Detail, event.Consequence.Kind, event.Consequence.Cause)
	}
}
```

One `switch` of `Fprintf` calls produces the transcript above. That is the payoff for
writing the log forward: the renderer has nothing to work out.

The sheet's one piece of emphasis is a characteristic at 0, which it puts above
everything else:

`render/sheet.go` — `writeZeroed`

```go
	for _, which := range chargen.CharacteristicOrder {
		if scores.Get(which) == 0 {
			zeroed = append(zeroed, which.String())
		}
	}
```

An **At 0** line names it and quotes p. 14 — "an extreme limitation", and a strong
recommendation that the Referee intervene — and the cell in the characteristics table
is bold. A 0 in a row of six numbers is what an eye slides over, and the referee holding
the sheet is the reader who must not miss it. It is still a view and not a rule: the
score is in the record, and the sheet only declines to bury it.

Golden files in `render/testdata/` are compared byte for byte, and are listed in
`.prettierignore` so the formatter does not rewrite them out from under the comparison.

---

## 8. `replay`: the verification path

`cmd/cschargen/replay.go` — `replayCommand`

```go
	replayed, err := chargen.New(chargen.Options{
		Seed:          original.Provenance.RNG.Seed,
		Decider:       chargen.NewReplay(original.Events),
		EngineVersion: version(),
		PolicyVersion: chargen.PolicyVersion,
		Setting:       world,
		Inputs:        original.Provenance.Inputs,
	}).Run()
```

Note what is reused from the record: the seed, the recorded choices and the _inputs_.
The stored throws are **not** input — every throw is recomputed from the seed. The log
is verification data.

The command is a page of flag handling around two calls into `chargen`.
Both used to live here; they moved to `chargen/verify.go` because verification is
the record's own business, and eight statements no test in `cmd` could reach came
with them.

`chargen/verify.go` — `Reproducible`

```go
// policy_version is deliberately not among the checks. Replay reapplies
// recorded choices and never consults the policy, so a record made under one
// policy replays under any other -- which is why PolicyVersion is recorded
// and not verified.
```

The setting hash, the schema version and the engine version are all checked;
`--ignore-provenance` skips the lot. It is the provenance half of a replay: a
refusal to try, as against a result.

`chargen/verify.go` — `Verify`

```go
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
```

Comparing logs rather than characters is deliberate: a divergence shows up in the log
several events before it reaches the character, and the sequence number is what a person
needs in order to find it.

The character-level check that follows used to be thinner than it looked. It ran on
`Characteristics` alone — the only comparable struct in `State`, so `==` reached
exactly as far as Go let it and the line was never revisited. A replay that lost
three of Step 20's four fields reported "identical". `verifyState` now compares the
whole character, through JSON rather than `reflect.DeepEqual`, because the two sides
are not built the same way: one was decoded from a file and the other constructed.

---

## 9. The gate

`Taskfile.yml` — `default`

```yaml
  default:
    desc: The whole gate, in the order a failure is cheapest to read
    cmds:
      - task: tidy
      - task: vet
      - task: lint
      - task: nilaway
      - task: test
      - task: ratchet
      - task: docs
```

**There is no CI.** The gate runs locally, before a pull request is opened, and nothing
else checks the tree; `task` is the whole of it.

Two parts of the gate are unusual enough to explain.

**The coverage ratchet.** Not a percentage — a count of uncovered statements per
package, checked in:

`coverage.ratchet`

```
github.com/philoserf/cschargen/career 1
github.com/philoserf/cschargen/chargen 79
github.com/philoserf/cschargen/cmd/cschargen 21
github.com/philoserf/cschargen/dice 0
github.com/philoserf/cschargen/render 3
github.com/philoserf/cschargen/setting 5
```

It fails in **both** directions — a count rising is coverage lost, a count falling is
coverage to lock in — and on a package appearing or vanishing. An integer rather than a
percentage because a percentage holds still while a guarded branch adds one covered
statement and one uncovered, and it grows more forgiving as the repository grows.

The trap: under `-coverpkg` every test binary instruments every package, so a test that
rolls a whole lifepath from an unpredictable seed reaches different branches on every
run and the ratchet disagrees between two runs on one machine for no visible reason. A test that
needs an unpredictable seed asks for `--terms 0`, which generates characteristics and
stops.

**A document held by a test.** `POLICY.md` is not commentary — it is checked:

`chargen/policy_internal_test.go` — `TestPolicyDocumentsEveryChoicePoint`

```go
	points := choicePointsInSource(t)
	// ...
			t.Errorf("%q is a choice point the engine offers and POLICY.md does not name", point)
	// ...
			t.Errorf("POLICY.md has a row for %q, which the engine does not offer", row[1])
```

Bidirectional: a choice point the document omits fails, and a document row the engine no
longer offers fails. `ERRATA.md` is held the same way, by
`chargen/errata_internal_test.go`: a reading no record stamps fails unless an allowlist
says why, a stamped identifier with no entry fails, and the allowlist is itself checked
against the document so an excuse cannot outlive its reason.

---

## 10. Where the linear order broke down

Three places where following the call chain meant holding two things at once, which is
worth knowing before you go reading:

- **Step 5 runs before Step 2, and finishes inside it.** `Run` calls
  `determineGenetics` before `rollCharacteristics`, because the genetics decide which
  characteristic method applies — and an outcome's adjustment is held until
  `rollCharacteristics` has assigned the scores (§5.2). The step numbers in the
  transcript are the book's and are not the execution order.
- **`Generator` carries about forty fields of cross-step state.** `transfer`,
  `enslaved`, `academyClosed`, `crisisSurvived`, `reachedOne`, `mustContinue`,
  `autoAdvance`, `rollAsHuman` … A result in Step 16 sets a flag that Step 9 reads three
  events later, and a loss that brings a score to 1 closes every career but one for the
  rest of the run (§5.8).
  Reading any one step in isolation will not tell you what can reach it; `grep` for the
  field.
- **Careers are a graph.** Colonist alone reaches Vagabond, Celebrity, Prisoner and a
  Diplomatic Service skill table through its mishaps and events. A linear read of
  `career/` gives no sense of this; `career/transcription_test.go` —
  `TestColonistReachesFourOtherCareers` is the closest thing to a map.
