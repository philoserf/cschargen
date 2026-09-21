# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# cschargen

Go CLI that generates rules-accurate Clement Sector characters.

**Ruleset baseline:** _Clement Sector Core Character Creation Book_, © 2026
Independence Games, author John Watts, 337 pp., held at
`~/Documents/Traveller/Clement/`. All page cites are to that artifact's **printed**
page numbers — printed page N is PDF page N+1, which matters when checking one.

Clement Sector is Cepheus Engine-compatible, **not Traveller**: Mongoose Traveller
SRD (2008) → Cepheus Engine SRD (2016) → _Clement Sector: The Rules_ (2016) → this
book. Nothing here is shared with `ctchargen` or `t5chargen` beyond house
conventions; do not reason from those rulesets.

`docs/PRD.md` is the v1 contract.
`ERRATA.md` carries every place the book is wrong or silent and what the engine
does about it — **add to it in the same change that implements the reading**, never
as a later pass.

`THEORY.md` and `WALKTHROUGH.md` are standing documents on **different cadences**,
and treating them as one thing gets both wrong. `THEORY.md` is about intent — the
record is the deliverable, `Inputs` is what was asked for and `State` is what
happened — and those claims outlive the code expressing them, so change it when a
decision invalidates one, never because a function moved. `WALKTHROUGH.md` quotes
forty snippets by file and symbol and rots the moment one is renamed, so do not
patch it per change: **regenerate it at release time**. Neither is checked by
anything — `POLICY.md` has a test holding it to the engine in both directions and
these have no equivalent — which is why the cadence is written down rather than
left to judgement.

## Commands

Run `task --list` for the current set.

Iterate with `go test ./chargen/ -run TestName` — plain, no flags. A subset run
writes no coverage profile, so it cannot move the ratchet; `task test` is the only
thing that should write `coverage.out`.

**CI runs exactly `task`.** Never add a check to CI that the local gate does not
run, and never add a tool to the gate without also installing it in the workflow.
The Go toolchain, golangci-lint, nilaway and prettier are deliberately unpinned: a
red gate on untouched code is the signal working — answer the finding rather than
pinning the tool.

Coverage is held by a **ratchet**, not a percentage: `coverage.ratchet` records the
count of uncovered statements per package, and `task ratchet` diffs current counts
against it, so the gate fails in both directions and on a package appearing or
vanishing. An integer rather than a percentage because a percentage holds still
while a guarded branch adds one covered statement and one uncovered, and it grows
more forgiving as the repository grows. Reflowing blank lines splits coverage
blocks, so a refactor can move these counts without changing what the tests reach
— read the diff before assuming a regression.

## The Product Identity boundary

The book declares its mechanics Open Game Content and declares as Product Identity
"all subsector names, world maps, world names, system names, system maps, vehicle
names, starship names, starship classes, artistic depictions of ships and vehicles,
and organizations" (OGL §16, p. 335). The same notice singles out one further term
— the book's word for a genetically engineered human — and says it too is not open
content.

**No Product Identity is committed to this repository.** World names, subsector
names and engineered-species names live in an external data file the user supplies;
`setting/sample.json` is invented. Career tables are here, because they are
mechanics rather than names.

**That includes the word.** This repository says "engineered human" throughout —
in its prose, its Go identifiers and its JSON keys — rather than the book's term,
which the notice names. A data file is something a user types, and a key is a
worse place to put a word we are not entitled to use than a comment is. The one
place the term appears is `setting/setting_test.go`, in the list of strings that
must never reach `sample.json`.

This is a hard rule, not a preference. Before committing a table, ask whether it
names a place, a person, a ship or an organization from the setting. If it does, it
belongs in the external file.

## Architecture

Each package's doc comment says what it holds.

`dice` knows the shapes of throw and nothing about what one means. `career` is data
consulted by the engine; it does not roll. Most of that layering is enforced by the
compiler because the reverse edge is an import cycle; `depguard` holds the edges
that would otherwise compile.

## Key conventions

- **Career definitions are hand-typed Go, not JSON.** Milestone 3 was filed to
  decide whether they belonged in a data file and decided they do not, against
  all thirty-four: a table result is a tree — a choice between two checks whose
  success branches differ, a check whose success branch holds another check — and
  that reads close to the page in Go and like nothing in JSON. The transcription
  is the cost either way. The careers' shapes genuinely differ, too: National Navy
  has nine skill tables and a commission, Vagabond has four and no enlistment
  throw, and Go structs make that variance a compile error. See
  `docs/MILESTONE-3.md`.
- **Every throw carries the page it came from.** A `ThrowEvent` without a cite is
  not auditable, which is the whole reason the log exists.
- **The event log is written as rules run, never reconstructed afterward.**
- Table-driven tests throughout; expected values from the book, cited.

## Releases

No release automation. `CHANGELOG.md` is written by hand and **the tag goes last**,
after the gate is green and the README's Status section says what the release
actually is — that section went stale for five milestones and was caught only at
the first release. `WALKTHROUGH.md` is regenerated in that same pass, for the
same reason.

Versions are prereleases until something has been played: `v0.1.0-alpha.N`. Two
other version strings ship in every record and move independently of the tag:

- `chargen.SchemaVersion`, which a record stamps and a replay refuses to cross.
- `policyVersion` in `cmd/cschargen/version.go`, which identifies `POLICY.md`.
  **Bump it when that document changes**, because a different policy is a
  different character from the same seed.

`cschargen version` prints all three plus the ruleset, so a bug report can be
matched against any record the binary wrote.

## Formatting

`.golangci.yml` runs `default: all`, and every disable and scoped exclusion in it
carries its reason in a comment beside itself. That is the bar for adding another.

gofumpt and goimports run **inside** golangci-lint, which is the single definition
of formatted for Go here. prettier is the same for every file that is not Go —
Markdown, JSON, YAML — run by `task docs`. `embeddedLanguageFormatting: "off"` is
load-bearing: prettier's default rewrites source inside fenced blocks, and this
repository's documents quote their own examples.

## Gotchas

- **A test that generates from a random seed makes the coverage ratchet
  unreproducible.** Under `-coverpkg` every test binary instruments every package,
  so a command test that rolls a whole lifepath from `rand.Uint64()` reaches
  different engine branches on every run — and the ratchet then disagrees between
  a developer's machine and CI for no reason anyone can see. A test that needs an
  unpredictable seed should ask for no terms at all, which generates
  characteristics and stops. `--terms 0` says that since alpha.4; `--terms -1`
  is the older spelling and still works, which is why the tests still use it.
- **`d66` is one throw, not two d6 rolls.** The results are 11-16, 21-26 … 61-66,
  so no digit is 0 or above 6. Reading it as 2d6 would collapse thirty-six results
  onto eleven.
- **Careers are a graph.** Colonist alone reaches Vagabond (mishaps 10, 12),
  Celebrity (event 56), Prisoner (Life Event 3) and a Diplomatic Service skill
  table (event 54). Adding a career means checking what its tables reach.
- **A mishap does not eject from every career.** Vagabond (p. 298), Prisoner
  (p. 263) and the slave career (p. 152) each say so explicitly, and the first two
  are the careers a character is most often forced into.
- **The slave career's name is not the book's.** The book heads pp. 150-154 with
  a term its OGL notice reserves, so this repository calls it
  `Engineered/Uplift Slave`. That word appears here only in the test that keeps it
  out.
- **Rank benefits are a floor, not an increment** (p. 116). A rank row's skills
  go through `applyRankBenefits`, which raises a skill to the level the row
  prints and does nothing where the character is already there; everything else
  a rank row carries is applied as printed. "(Any)" picks a specialty the
  character does not hold at level 1 or higher. Data written against a rank
  table carries what the page prints.
- **Aging is indexed by term number and gated by homeworld tech level** (pp.
  122-123), not by age. Apparent age is a derived lookup (p. 125), never rolled.
- **Homeworld is not fixed after character creation.** Six of Colonist's eleven
  mishaps reassign it, so the record holds a homeworld history — and a
  reassignment changes the tech level and nothing else (ERRATA E-9).
- **The term limit has two ceilings.** The policy's `--terms` and the homeworld's
  own maximum are both caps, and the lower wins (p. 125). `--max-terms` sets the
  second of them, standing in for a world the setting file does not have.
- **`Inputs` is what was asked for; `State` is what happened.** `Inputs.Career`
  keeps the career asked for beside the career served, and `Inputs.Homeworld`
  keeps the world asked for beside `State.Homeworlds`. Writing a result back
  into `Inputs` made a finished record indistinguishable from one that had
  requested its own outcome, and replay — which is handed the record's own
  inputs — then honoured the request instead of re-rolling. Every auto record
  diverged.
- **A `Choice` cannot be gated on `Asker`.** `Replay` implements `Choose` and
  not `Ask`, so a choice point gated that way is offered while generating and
  skipped while replaying, and the recorded answer is never consumed. Gate on
  `Inputs.Interactive`, which travels in the record. `Asker` is for Step 20,
  whose answers live in `Inputs` and which replay does not re-ask.
- **A flag whose zero is a legal value needs `flag.FlagSet.Visit`.** A plain int
  cannot tell an unset flag from one set to its zero value, which is how
  `--terms 0` came to mean the policy's four. `Visit` walks only the flags
  given. `--max-terms` does not need it: the validator refuses a world with a
  maximum of zero terms, so zero there is free to mean absent.
- **A bound a flag and the data format share is one exported constant.**
  `setting.MaxTechLevel` and `setting.MaxTerms` are exported because
  `--tech-level` and `--max-terms` override a world's own numbers and have to be
  held to the range a world's numbers are held to. Two copies of a bound is how
  the flag came to take what the file refused.
- **Being owned is not permanent.** p. 42 makes the slave career an enslaved
  character's first, and three results in the early life tables end an
  enslavement outright — "Continue your character as a free altrant or uplift".
  A character born owned may reach Step 9 free, and a test asserting otherwise
  passes only until the dice move.
- **A subsector throw that lands on no chart entry is thrown again.** p. 39's
  chart is a 1d6 and a setting file may name fewer than six subsectors. Throwing
  again keeps the chart's weighting; turning the miss into a choice hands it to
  the policy, which answers with the first subsector the file lists.
- **`lifepath` in the tests keeps what it makes.** Thirty-two call sites sweep
  the same seeds at seven term counts, so the suite walked 1,860 lifepaths to
  see 420 characters. A character is a pure function of its seed, term limit,
  setting and decider, and nothing in the tests writes to one — so a test that
  needs to modify a character builds its own with `options()` and `generate()`,
  which is what the tests taking a different decider already do.
