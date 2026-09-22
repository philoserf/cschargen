# THEORY.md

What you need to hold in your head to change `cschargen` without damaging it.

This is not a tour of the directories — `go doc ./...` does that better, and every
package here carries a doc comment that says what it holds. This is the set of ideas
the code is built on, the invariants it depends on, and the places where changing
something obvious would break something that is not.

---

## 1. What the system is for

The naive reading is that this program makes characters. It does not, quite. It
**executes a printed procedure and writes down what it did**, and a character is what
falls out at the end.

That distinction is the whole design. The Clement Sector lifepath is twenty steps of
dice and table lookups that the book itself calls a minigame requiring a game session
of its own. Done by hand it takes an evening and goes wrong silently — a modifier
missed at Step 14 shows up as a character who is merely slightly too good, and nobody
ever finds it. The thing that is actually hard to buy is not the character. It is the
**certainty that the character is what the book says it is**, and the ability to go
back afterwards and check.

So the deliverable is a record with three parts (`chargen.Character`):

- `Provenance` — under what rules, what code, what data and what decisions this was
  made.
- `State` — the character.
- `Events` — every step entered, every throw made, every choice taken and every
  consequence applied, in order, each carrying the page it came from.

`State` is the part a player cares about and the least interesting part of the file.
Delete it and you could rebuild it from `Events`. Delete `Events` and you have an
unfalsifiable assertion.

Everything downstream follows. `render` is explicitly forbidden from computing
anything: "if a number is not in the record, it does not appear on the sheet"
(`render/sheet.go`). `replay` re-runs the engine from the seed and compares logs.
`ThrowEvent` carries a `Cite` because a throw without a page number cannot be
audited, and being auditable is the reason the log exists.

**If you remember one thing:** a change that makes a character better and the record
less checkable is a regression here, even when it makes the character more correct.

### The vocabulary

The domain words are the book's, used as if their meanings were settled, because in
the book they are:

- A **term** is four years in a career. A **service** is one spell in one career and
  may span several terms; `State.Terms` and `State.Services` are both lists and they
  are not parallel.
- An **assignment** is a track within a career — a career has one or several, each
  with its own survival and advancement targets and its own skill table.
- A **throw** is 2d6 against a target with modifiers; a **roll** is dice with no
  target. `Log.Throw` and `Log.Roll` are different methods for exactly this reason,
  and a throw event with `Target` nil means "this was never a check", not "the check
  had no target". Recording a zero target would read as a check that always failed.
- **d66** is one throw producing 11–16, 21–26 … 61–66. It is not 2d6. Reading it as
  2d6 collapses thirty-six results onto eleven, and every career's event table is a
  d66.
- A **tie** is an Ally, Contact, Rival or Enemy with a Relationship Rating. The kind
  and the rating can disagree and the engine keeps the printed kind when they do.
- A **benefit roll** is a mustering-out roll; benefits are a queue of scoped batches,
  not a count, because the book grants both "a +1 modifier to all Benefit rolls made
  in this career" and "three Benefit rolls at +1", and those are different things.

---

## 2. The Product Identity boundary is an architectural constraint, not a policy note

This is the single least guessable thing about the codebase, and the one most likely
to be damaged by someone doing something otherwise reasonable.

The book's OGL notice declares its mechanics open content and declares as Product
Identity all world names, subsector names, system names, ship names and
organizations — plus one further term, the book's word for a genetically engineered
human. Steps 3 and 4 of character creation are _tables of world names_. Those tables
are load-bearing: a homeworld sets background skills, primary language, maximum age,
maximum terms, tech level, and whether engineered characters may be born there.

So the repository **cannot contain the data that Steps 3–4 read**, and the resolution
is the `setting` package: a schema, a validator, a content hash and a loader, with the
actual worlds supplied by the user from their own copy of the book.

Two consequences that look like ordinary design decisions and are not:

**`setting` is not a data/logic split.** It is a legal one. Career tables are also
data — thousands of lines of hand-typed table — and they are compiled in, because
mechanics are open content and names are not. If you find yourself thinking "the
careers should be a data file too, for symmetry", the symmetry is not there to be had.
The separate argument against it is about shape rather than licensing, and it is in
section 5 below.

**The word is avoided everywhere, including in identifiers and JSON keys.** The
repository says "engineered human" in prose, Go identifiers and data-format keys. The
reasoning in `setting/setting.go` is worth keeping: a data file is something a user
types, and a key is a worse place to put a word we are not entitled to use than a
comment is. `setting/setting_test.go`'s `TestTheSampleIsInvented` is this boundary held
as a test — a list of proper names from the book that must not appear in
`setting/sample.json`, including that term.

Note what that test does _not_ guard: it reads `sample.json` and nothing else. The term
currently appears in about thirty places in Go comments and quoted rules (GitHub issue
#108), because quoting p. 42 verbatim brings it along. The guard is the rule; the test
is a tripwire on the one file most likely to be pasted into.

**Before adding any table, ask whether it names a place, a person, a ship or an
organization.** If it does, it belongs in the external file.

---

## 3. The organizing ideas

### 3.1 The log is written forward, never reconstructed

Events are appended as the rules run. `Log.append` assigns the sequence number and
returns it, so the caller can name the event it just wrote as the `Cause` of what
follows. This is why the engine threads `step`/`cause` integers through nearly every
function signature — it looks like parameter noise and it is the causal spine of the
record.

Nothing anywhere rebuilds a log from a finished character. If you add a rule, the
consequence is logged **at the point the rule fires**, with the identifier of the
throw or choice that caused it.

### 3.2 Every decision that is not dice goes through one interface

`Decider` has one method that matters: `Choose(Choice) (int, error)`. Every fork in
the lifepath — which career, which assignment, which skill table, which specialty,
which characteristic — is a `Choice` with a stable `Point` identifier, an option list,
and a page cite.

Three implementations, and they are **not peers**:

| Decider  | What it is                                                       |
| -------- | ---------------------------------------------------------------- |
| `Player` | Asks a person. The only one that also implements `Asker`.        |
| `Policy` | Answers "the first option", always. Documented by `POLICY.md`.   |
| `Replay` | Decides nothing. Reapplies recorded answers and checks they fit. |

`Replay.Choose` is the interesting one. It does not just return the recorded index — it
checks that the recorded `Point` matches the point being asked and that the recorded
option list is _identical_ to the list now offered. A recorded index means nothing
against a different list. This is what stops a record generated with `--career` forced
(a one-entry list) from replaying against the full career list and picking the wrong
entry.

**The trap this creates**, and it has already been fallen into: `Replay` implements
`Choose` but not `Ask`. So a choice point gated on `_, ok := g.decider.(Asker)` is
offered while generating and skipped while replaying, and the recorded answer is
consumed by whatever choice comes next — which then fails the point check, or worse,
does not. `mayReturnToEducation` in `chargen/generator.go` is gated this way and is
GitHub issue #106. **Gate on `Inputs.Interactive`, which travels in the record.**
`Asker` is for Step 20, whose four free-text answers live in `Inputs` and which replay
never re-asks.

### 3.3 `Inputs` is what was asked for; `State` is what happened

`Provenance.Inputs.Career` holds the career requested. `State.Services` holds the
careers served. `Inputs.Homeworld` holds the world requested; `State.Homeworlds` holds
where the character was actually born and every world they were later reassigned to.

They used to be the same fields, and the bug was subtle enough to be worth carrying
forward: writing a _result_ back into `Inputs` makes a finished record
indistinguishable from one that requested its own outcome. Replay is handed the
record's own `Inputs`, so it then honoured the request instead of re-rolling, and every
auto-generated record diverged.

**Never write an outcome into `Inputs`.** The asymmetry is the point.

### 3.4 The effect vocabulary is a closed alphabet and `apply` is the fold over it

`career.EffectKind` enumerates everything a table result can do to a character —
around thirty kinds, from `EffectSkill` through `EffectCheck` (a nested throw with two
branches that are themselves effect lists) to `EffectUnimplemented`.

`chargen.apply` is one switch over that alphabet, and it is deliberately not split.
The argument in its comment is the load-bearing one: `exhaustive` is configured with
`default-signifies-exhaustive: false`, so one switch with no default guarantees a new
kind is a **build failure**. Two switches need a default between them, and a default
is where a forgotten kind goes to be ignored quietly. The `//nolint:cyclop,funlen,
gocyclo,maintidx` on it is bought with that argument, not with fatigue.

`EffectGroup` and `EffectSubTable` recurse rather than branch, which is how the
alphabet stays small while the tables stay expressive.

`EffectUnimplemented` is the honest kind: a result the engine cannot carry out is
recorded _with what the book asked for_, so the record says what it could not do rather
than quietly doing something else. When you add a rule, the stub is an
`EffectUnimplemented` carrying the page's own words — never a silent no-op.

### 3.5 `career` is data that cannot roll

The package brings no dice and no engine import. Most of the layering is enforced by
the compiler, because the reverse edge would be an import cycle. The one edge that
would otherwise compile — `dice` importing `career` — is held by a `depguard` rule
whose description states the rule directly: "dice knows the shapes of throw, never what
a throw means".

`career/build.go`'s constructors (`skill`, `skillAt`, `skillZero`, `check`, `choice`, …)
exist so a transcribed table reads close to the page. They are the only way effects are
built in the package, because a table written with struct literals buries the one field
that differs between two rows in a wall of field names.

### 3.6 Careers are a graph, not a list

`career.All()` returns thirty-four careers, but the tables reach each other. Colonist
alone reaches Vagabond (mishaps 10 and 12), Celebrity (event 56), Prisoner (Life Event 3) and a Diplomatic Service skill table (event 54). Adding or editing a career means
checking what its tables reach.

Three careers are forced destinations rather than choices, and each is closed off in a
different direction:

- **Vagabond** is where the rules send a character who has nowhere else: three
  consecutive failed enlistments (p. 110), no eligible career, an aging crisis
  survived, a mental characteristic at zero.
- **Prisoner** cannot be chosen voluntarily (p. 111); only a table result puts a
  character there.
- **The slave career** — this repository's name for it, not the book's — has no
  enlistment throw and is reached only from p. 42, as an enslaved character's first
  career. `eligibleCareers` excludes both it and Prisoner, and `ownedFirstCareer`'s
  claim that p. 42 "is the only way into that career" is a claim that list has to hold
  up.

**A mishap does not eject from every career.** Vagabond (p. 298), Prisoner (p. 263) and
the slave career (p. 152) each say so explicitly — and the first two are precisely the
careers a character is most often forced into, so the exception is not rare.

**Being owned is not permanent.** Three results in the early-life tables end an
enslavement outright, so a character born owned may reach Step 9 free. A test asserting
otherwise passes until the dice move.

### 3.7 Three version strings, protecting three different things

They move independently and a record stamps all three:

| String                   | Identifies              | Checked by `replay`?              |
| ------------------------ | ----------------------- | --------------------------------- |
| `chargen.SchemaVersion`  | the shape of the record | yes — refuses to cross            |
| engine version (the tag) | the rules code          | yes, unless `--ignore-provenance` |
| `policyVersion`          | `POLICY.md`             | **no, deliberately**              |

The third is the one worth understanding. Replay reapplies recorded choices and never
consults the policy, so a record made under one policy replays under any other —
`checkProvenance` says so in a comment rather than leaving the omission to be read as
an oversight. But a _different policy is a different character from the same seed_, so
`policyVersion` must be bumped whenever `POLICY.md` changes, or two records claiming
the same provenance describe different procedures.

The setting file's content hash is checked the same way the schema is: a record
generated against one transcription will not replay against another.

### 3.8 ERRATA is a first-class output, not a notes file

`ERRATA.md` carries every place the book is wrong, silent or self-contradicting, split
into typos (the book plainly means otherwise), readings (the book does not say and the
engine had to choose) and limits (the engine knowingly does not model something).
Forty-three entries.

The rule in `CLAUDE.md` is strict and worth obeying: **add the entry in the same change
that implements the reading**, never as a later pass. A reading implemented without an
entry is an undocumented house rule, and the whole auditability argument collapses to
"trust the code".

`Provenance.Deviate(id)` stamps a reading into the record so a character generated
under one reading stays auditable after the reading changes. In practice seven
identifiers of forty-three are ever stamped — see the loose ends below.

### 3.9 Coverage is an integer ratchet, in both directions

`coverage.ratchet` records uncovered statements **per package**, and `task ratchet`
diffs against it. It fails when a count rises (coverage lost), when a count falls
(coverage to lock in), and when a package appears or vanishes.

An integer rather than a percentage because a percentage holds still while a guarded
branch adds one covered statement and one uncovered, and it grows more forgiving as
the repository grows.

**The trap:** under `-coverpkg` every test binary instruments every package. A test
that rolls a whole lifepath from an unpredictable seed reaches different engine
branches on every run, and the ratchet then disagrees between a developer's machine
and CI for no reason anyone can see. A test that needs an unpredictable seed must ask
for **no terms at all** — `--terms 0` — which generates characteristics and stops.

Separately: reflowing blank lines splits coverage blocks, so a pure refactor can move
these counts without changing what the tests reach. Read the diff before assuming a
regression.

---

## 4. The seams

### The setting file

The only substantial input the program does not carry. `setting.Load` validates it
_before_ generation rather than during, so a bad file is refused at the command line
rather than halfway through Step 3. `DisallowUnknownFields` turns a misspelt key into
an error instead of a zero value — which for `maximumTerms` would silently mean a
character who may serve no terms. `InvalidError` reports every problem rather than the
first, because someone transcribing sixty worlds should be told about all of them at
once.

`setting/sample.json` is invented and embedded, and a character generated against it
is stamped `Sample: true` and says so on its own sheet. That stamp exists so an
invented character is never mistaken for one set on the published worlds.

### The command line

`cmd/cschargen` is thin by design and holds one thing that is genuinely load-bearing:
**a name the data does not have is a flag error, not a generation outcome.**
`checkNames` holds `--subsector`, `--homeworld`, `--career`, `--tech-level` and
`--max-terms` against the data and the book before anything is generated. The engine
would refuse them too, but partway through Step 3, where it reads as a generation
failure rather than as the typo it is.

The exit contract: a misused command line exits 1 with a message beginning `usage:`;
any other failure exits 1 with a message that does not. **The two are told apart by
the message, not the status.**

Warnings go to stderr and still exit 0. A referee generating a hundred NPCs should not
have to check a hundred exit codes to find out that nine drifted into Vagabond after a
failed enlistment — the note says so once, or as a count for a batch.

Two flag traps live here:

- **A flag whose zero is a legal value needs `flag.FlagSet.Visit`.** A plain int cannot
  tell an unset flag from one set to zero, which is how `--terms 0` came to mean the
  policy's four. `given()` walks only the flags actually passed. `--max-terms` does not
  need it, because the validator refuses a world with a maximum of zero terms, so zero
  there is free to mean absent.
- **A bound the flag and the data format share is one exported constant.**
  `setting.MaxTechLevel` and `setting.MaxTerms` are exported precisely because
  `--tech-level` and `--max-terms` override a world's own numbers and must be held to
  the same range. Two copies of a bound is how the flag came to accept what the file
  refused.

### Documents held by tests

`POLICY.md` is not documentation about the code — it is checked against it.
`TestPolicyDocumentsEveryChoicePoint` parses the document's table, collects every
choice point named in the package source, and fails **in both directions**: a choice
point the document omits, and a document row the engine no longer offers.

This is the strongest seam in the project and the pattern worth copying. `ERRATA.md`
has no equivalent, and that is where the drift is (§6).

### The sheet and the transcript

`render` reads the record and only the record. The sheet is modelled on the book's
Long Form Character Sheet; the transcript walks the event log. Golden files in
`render/testdata/` are compared byte for byte and are in `.prettierignore` so they are
not reformatted out from under the comparison.

---

## 5. What the system is shaped to accommodate

**Adding a career.** A file in `career/`, a line in `All()`, and a check of what its
tables reach. The compiler catches a shape that does not fit; `TestSkillTablesHaveSixRows`,
`TestEveryCareerReachesLifeEventsAt31To36`, `TestEveryCareerHasItsPrintedShape` and
`TestEveryTranscribedSkillIsOnTheList` catch most transcription faults. What they
cannot catch is a right-shaped table with a wrong number in it.

**Adding an effect kind.** Add the constant, and `apply` fails to compile until it is
handled. This is the mechanism working.

**Adding a choice point.** `TestPolicyDocumentsEveryChoicePoint` fails until `POLICY.md`
has a row, and the row is a decision someone has to write down.

**Adding a rule the book is unclear about.** `ERRATA.md` entry in the same change,
`Deviate` at the site, test citing the entry by identifier.

**A new setting file.** Nothing in the engine changes. A world is data.

Now the other direction — **what would require rethinking something fundamental**:

- **A different ruleset.** Every page cite, every table and `Ruleset` itself are keyed
  to one printed artifact. This is not a generic Cepheus engine and was never meant to
  be. `CLAUDE.md` is explicit that nothing is shared with the sibling `ctchargen` or
  `t5chargen` beyond house conventions.
- **A policy that plays well.** `Policy.Choose` returning the first option is not a
  placeholder; it is a rule a reader can check against the page. A policy that
  optimised would need an argument per choice point, and `POLICY.md` says so twice. A
  strategy belongs beside the policy as a second decider, not inside it.
- **Derived state the log does not produce.** Any `State` field computed at render
  time, or set without a corresponding consequence event, breaks the claim that the
  record is the source of truth and quietly weakens `replay`.
- **Careers as a data file.** Decided against, against all thirty-four. The argument
  is about shape: a table result is a tree — a
  choice between two checks whose success branches differ, a check whose success
  branch holds another check — and that reads close to the page in Go and like nothing
  in JSON.

**Where a maintainer who did not understand the theory would cause damage:**

1. Writing a generation result back into `Inputs` "so the record is complete".
2. Adding a `default:` to `apply` to quiet a linter.
3. Collapsing a repeated cell in a career table into a named constant. `goconst` is
   scoped off `career/` for this reason alone, and `.golangci.yml:150-155` states it:
   "one constant means one typo reaches every cell that shares it, and the second
   reading stops being independent of the first".
4. Gating a new choice point on `Asker` (see §3.2).
5. Adding a test that generates from a random seed with a real term count, breaking the
   ratchet for everyone.
6. Committing a world or species name — or the reserved word — into a table, a test
   fixture or a JSON key.
7. Computing a number in `render` because the record does not carry it.

---

## 6. Where the theory is thin, and where the code disagrees with it

Read this section as the map of what to distrust.

**The ERRATA stamping promise is not kept.** `ERRATA.md` says records stamp the
identifier of every deviation that applied to them. Seven identifiers of forty-three
are ever stamped; thirty-one more appear in Go comments as the reason the code reads
the way it does and reach no record. `E-24` — a relationship ends rather than inverting
— changes outcomes, is implemented in `chargen/ties.go`, and leaves no mark. Nothing
enforces the correspondence the way `TestPolicyDocumentsEveryChoicePoint` enforces the
policy's. This is the largest gap between what the project says about itself and what it
does.

**The second reading covers three careers of thirty-four.** `career/transcription_test.go`
names repetition as the mechanism that makes hand-typed tables trustworthy, and applies
it to Colonist, Vagabond and Prisoner. The structural tests do cover all thirty-four,
but they cannot catch the error transcription actually makes: a right-shaped table with
a wrong number in it. I cannot tell from the code whether this is a deliberate sample or
an unfinished job, and the file does not say.

**`replay` verifies the log completely and the character partially.** `compare` walks
both event logs and then checks exactly one of twenty `State` fields. The single
`Characteristics` comparison concedes that the log alone was not trusted to imply the
character; having conceded it, the check stops there — most likely because
`Characteristics` is the only comparable struct and `==` reached no further.
`State.Age` is set outside any consequence event, so a divergence confined to it
would pass. `Term.Event` and `Term.Mishap` are worse: nothing ever writes them (below).

**Two record fields are declared and never written.** `Term.Event` and `Term.Mishap`
are in the schema, carry `omitempty`, and nothing in the engine assigns either. They
read as the convenience index onto what happened in a term; the information is really
several consequence events away in `Character.Events`. I read this as a leftover from
before the log became the single place a result is recorded, but that is inference.

**`setting`'s package doc points at `data/setting.sample.json`.** The file is
`setting/sample.json`, embedded; `/data/` is gitignored precisely because it is where a
user's transcription goes.

### Things I am inferring rather than reading

- **That `"Any"` is meant to stay unresolved.** `career/academy.go` and many career
  files pass `"Any"` as a single specialty, and both `applySkill` and `rankSpecialty`
  take a one-element list verbatim — so a sheet prints `Gun Combat (Any)` and the
  player picks later. `rankSpecialty` resolves a _multi-element_ list against p. 116's
  "not already at level 1 or higher" rule. So there are two different things called a
  specialty choice, and I believe the split is deliberate — the book writes "(Any)"
  over an unnamed list and writes real lists elsewhere — but nothing states it.
- **That `(limit)` entries in `ERRATA.md` are never meant to be stamped.** None is, and
  the preamble does not acknowledge the category at all.
- **Why `chargen` is one large package rather than several.** The step files are
  independent enough to split, but every one is a method on `Generator`, which holds
  about forty fields of cross-step state (`transfer`, `enslaved`, `academyClosed`,
  `crisisSurvived`, `mustContinue`…). I read this as the honest shape of a procedure
  whose steps genuinely affect each other, rather than as a package that failed to get
  split — but the eighteen `*_internal_test.go` files suggest the seams are felt.

### Things I did not read closely

The thirty-one career files I did not open individually. I read Colonist, Vagabond,
Prisoner, National Navy and the slave career as the shapes `CLAUDE.md` names as
extremes, plus the shared tables and the skill list, and sampled the rest by grep. A
transcription error in the other twenty-nine would not have shown up in this pass.

---

## Index

Loose ends this pass turned up, each filed as a GitHub issue and carried on the
Workbench board. The issue is the durable reference; this table is a map from the
sections above to it.

| #   | Severity | Issue                                                                                     | Primary location                                     |
| --- | -------- | ----------------------------------------------------------------------------------------- | ---------------------------------------------------- |
| 1   | medium   | #118 — ERRATA.md promises every applied reading is stamped; 7 identifiers of 43 ever are  | `ERRATA.md:14`, the 11 `Deviate` sites in `chargen/` |
| 2   | medium   | #111 — `replay` verifies the whole log and one of twenty `State` fields                   | `cmd/cschargen/replay.go:104`                        |
| 3   | —        | #121 — fixed: `Generator.Run`'s doc comment now describes the twenty steps                | `chargen/generator.go`                               |
| 4   | —        | #120 — fixed: the three effect kinds describe the engine, not the plan                    | `career/effect.go`, `career/build.go`                |
| 5   | medium   | #119 — The independent second transcription covers 3 careers of 34                        | `career/transcription_test.go:9`                     |
| 6   | low      | #129 — ERRATA.md says it has two kinds of entry and carries three                         | `ERRATA.md:8`                                        |
| 7   | low      | #127 — `setting`'s package doc names a sample path that is gitignored and does not exist  | `setting/setting.go:19`                              |
| 8   | low      | #128 — The comment justifying one ragged consequence struct counts 13 kinds; there are 17 | `chargen/event.go:105`                               |
| 9   | medium   | #117 — `Term.Event` and `Term.Mishap` are declared in the schema and never written        | `chargen/state.go:57`, `chargen/serve.go:22`         |

**Total: 9 issues (0 critical, 0 high, 6 medium, 3 low); #120 and #121 fixed since.**

Already tracked on GitHub and deliberately not re-filed: **#106** (a `Choice` gated on
`Asker` is skipped on replay — §3.2), **#108** (the reserved word appears in
thirty-three places — §2), **#104** (a communal birth compounds into eighty-one
relatives).
