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

`docs/PRD.md` is the v1 contract. `docs/MILESTONE-1.md` is the current plan.
`ERRATA.md` carries every place the book is wrong or silent and what the engine
does about it — **add to it in the same change that implements the reading**, never
as a later pass.

## Commands

Run `task --list` for the current set.

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

| Package         | Holds                                                                                   |
| --------------- | --------------------------------------------------------------------------------------- |
| `dice`          | Seeded stream and the shapes of throw: 1d6, 1d3, Nd6, 2d6, 3d6-drop-lowest, d66, target |
| `chargen`       | The engine: character record, event log, `Decider`, and the rules, a file per step      |
| `career`        | Career definitions and the shared Injury and Life Events tables                         |
| `setting`       | The external world, subsector and species data: types, validator, content hash          |
| `render`        | Record → Markdown sheet and lifepath transcript                                         |
| `cmd/cschargen` | Flags, subcommands, exit statuses                                                       |

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
- **Errors are checked on their own line**, never inline: `err := f()` then
  `if err != nil`.
- **Every test calls `t.Parallel()`**, top level and subtest.
- Table-driven tests throughout; expected values from the book, cited.

## Releases

No release automation. `CHANGELOG.md` is written by hand and **the tag goes last**,
after the gate is green and the README's Status section says what the release
actually is — that section went stale for five milestones and was caught only at
the first release.

Versions are prereleases until something has been played: `v0.1.0-alpha.N`. Two
other version strings ship in every record and move independently of the tag:

- `chargen.SchemaVersion`, which a record stamps and a replay refuses to cross.
- `policyVersion` in `cmd/cschargen/version.go`, which identifies `POLICY.md`.
  **Bump it when that document changes**, because a different policy is a
  different character from the same seed.

`cschargen version` prints all three plus the ruleset, so a bug report can be
matched against any record the binary wrote.

## Lint posture

`.golangci.yml` runs `default: all`. Three of its disables are linters deprecated
upstream and replaced by one that is enabled, so they disable no check that is not
still being made. Everything else carries a measured count:

- **Disabled outright:** `exhaustruct_v5` (26 findings, all on types whose unset
  fields are the design — `Event` is a discriminated union where three of four
  payload pointers are nil by construction) and `gochecknoglobals` (2, immutable
  tables Go has no const form for).
- **Scoped to `career/`:** `mnd` (51) and `goconst` (27) and `funlen` (4). A
  transcribed table is a literal of printed numbers and repeated cells, and
  collapsing a repeated cell into a constant is exactly the coupling the
  transcription exists to prevent — one constant means one typo reaches every cell
  that shares it, and the second reading stops being independent of the first.
- **Scoped to `career/*_test.go`:** `cyclop` and `gocognit` (4). A test that checks
  a table cell by cell has the complexity of the table.
- **Configured, not disabled:** `gosec` waives G404 and G304, `godot` accepts a
  comment ending in a quotation mark.

That is the bar for adding another: count the findings, read them, and write down
why they are wrong here. Every other finding the gate has produced — a hundred and
thirty so far — was answered in the code.

`nolintlint` requires a specific linter and an explanation and fails on an unused
`//nolint`, so a blanket directive will not pass.

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
  unpredictable seed should ask for `--terms -1`, which generates characteristics
  and stops.
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
  own maximum are both caps, and the lower wins (p. 125).
