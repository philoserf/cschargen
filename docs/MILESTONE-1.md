# Milestone 1: the walking skeleton

2026-09-19. Status: **shipped**, in seven pull requests (#2 through #8). Tracks
[#1](https://github.com/philoserf/cschargen/issues/1).

## What it turned out to be

The plan below is the one that was followed, including the scope finding that
opens it. Seven PRs, in the order listed, each leaving the gate green.

Five things the plan did not know:

1. **A test that generates from a random seed makes the coverage ratchet
   unreproducible.** CI and a developer's machine disagreed by six statements with
   nothing in the diff to explain it, because a command test rolled eight whole
   lifepaths from `rand.Uint64()` and under `-coverpkg` those reached different
   engine branches every run. The test now asks for characteristics alone.
2. **A career a character re-enters is a second spell of service, not a
   continuation.** `State.Service` returned the first match, so every update after
   an ejection-and-return landed on the spell that was already over.
3. **A career with no benefit rolls left to spend was not clearing its queue**, so
   a mishap's "lose two benefit rolls" stayed behind as a debt and would have been
   charged again the next time the character joined that career.
4. **Two more readings were forced** and are in `ERRATA.md`: an event's skill check
   is stated nowhere in this book (E-7), and the Rank 0 row of every rank table is
   unreachable unless entering a career counts as achieving Rank 0 (E-8).
5. **"Every consequence names a step, a throw or a choice as its cause"** is worth
   writing as a test. It caught three places where the cause was whatever happened
   to be logged last, which in a nested effect tree is often another consequence.

Measured after the fact: across 120 generated characters, all three careers are
reached in play -- 203 spells of Colonist, 92 of Vagabond, 6 of Prisoner -- and 14
runs end at an unimplemented career, each with the consequence that names it.

---

## The scope finding that came first

`docs/PRD.md` proposed milestone 1 as "the term loop for one career (Colonist),
human only, other steps stubbed — the whole engine shape in one career."

**Colonist is not closed.** Reading its tables through, a character in it can be
sent out of it four ways:

| Route                          | Where it goes                                                           | Forced?                                          |
| ------------------------------ | ----------------------------------------------------------------------- | ------------------------------------------------ |
| Mishap 10, Mishap 12 (p. 174)  | Vagabond, with the Transient assignment named in 10                     | **Yes**                                          |
| Event 56 (p. 176)              | Celebrity, as a Star                                                    | **Yes**                                          |
| Life Event 3 (p. 120)          | Prisoner, one term                                                      | No — "lose two benefit rolls" is the alternative |
| Event 54 (p. 176)              | Two rolls on Diplomatic Service's Assignment: Ambassador table (p. 185) | Yes, but it needs one table, not the career      |
| Mishaps 3, 4, 5, 6, 8 (p. 174) | "Choose another career" — any career                                    | Yes, but the career is the player's              |

So a one-career skeleton cannot reach its own "done". This is a property of the
system rather than of Colonist: the careers are a graph, and the terminal careers
are where every failure path drains.

### The closure taken

**Milestone 1 implements three careers: Colonist, Vagabond and Prisoner.**

Vagabond and Prisoner are the two the rules force a character into, they are the
smallest in the book (Vagabond has no Advanced Education table, p. 296; Prisoner
has a single assignment, p. 261), and neither has an enlistment throw — so they
exercise the career-entry path's exceptional branch without adding an enlistment
variant. They are also two of the four the PRD's own milestone 3 named as the ones
that set the data format, so the work is not duplicated later.

Celebrity (event 56) and "choose another career" (mishaps 3–8) are **stubbed**: the
engine records a `career_transfer_unimplemented` consequence naming the target and
ends generation cleanly. The record says what it could not do rather than quietly
doing something else. Event 54 is stubbed the same way — it needs one skill table
from a career milestone 3 will bring in, and inventing it now would mean writing
data this milestone cannot test against the book.

Two closures were considered and not taken. **Stub every exit** keeps milestone 1
at one career but leaves the mishap path — the half of the term loop that is
hardest to get right — unexercised end to end. **Reorder the milestones** so that
career transfer becomes milestone 2 and origins slip to 3 defers the same work by
renaming it. Either is a one-line change to this document's scope if the judgment
turns out wrong; the PR sequence below is unaffected except at PR6.

## Two corrections to the PRD that came out of the same reading

1. **Relationship Ratings are in scope for the term loop.** The PRD lists NPC
   generation (pp. 315–323) as a non-goal, which is right for _generating_ NPCs
   and wrong for the numeric rating an Ally or Contact carries: Life Events 5, 6
   and 7 (p. 120) all write to it, and Life Event 4 orders the relationship types
   for removal. An Ally/Contact/Rival/Enemy is therefore a record with a rating
   from milestone 1 onward. Generating an NPC around one stays out.

2. **Injury recovery is charged against credits that do not exist yet.** The
   Injury table charges 1d6 × 10,000 HFCredits to recover a lost characteristic,
   "if the character can afford this cost", and makes the loss permanent otherwise
   (p. 119). But mustering out is what produces credits (p. 126), and it happens
   after the term the injury occurred in. A character injured in their first term
   can never afford recovery. **Recorded reading:** the engine charges against
   credits held at the moment of the injury, so an early injury is permanent, and
   the injury record carries a `permanent` flag set at that moment rather than
   recomputed later. The alternative — deferring the decision to muster out —
   would let a benefit roll retroactively heal an injury that the mishap table has
   already acted on.

Both go into `ERRATA.md` and `docs/PRD.md` when PR1 lands, not as a later pass.

## Package layout

```
dice/       Throws: 2d6, 1d6, 3d6-drop-lowest, 1d3, d66, target-number.
            Seeded stream, freestanding, no knowledge of the rules.
chargen/    The engine. Character record, the event log, the Decider
            interface, and the rules in files named for the step they
            implement. Consumes career definitions; owns no tables.
career/     The three career definitions and the two shared tables
            (Injury p.119, Life Events p.120). Hand-typed Go, not JSON
            (see below).
render/     JSON record -> Markdown sheet and lifepath transcript.
cmd/cschargen/  Flags, subcommands, exit statuses.
```

`setting`, `origin` and `species` do not appear until milestone 2. Milestone 1
takes the homeworld's two load-bearing numbers — tech level and term cap — as
flags, which is exactly the seam the data file replaces.

**Career definitions are hand-typed Go structs in this milestone, not a JSON data
format.** The format is what milestone 3 proves against thirty-four careers, and
the shape genuinely varies — National Navy has nine skill tables and a commission,
Vagabond has four tables and no enlistment. Designing the schema against three
careers and then discovering the variance would mean migrating data twice. Go
structs make the variance a compile error instead.

## PR sequence

Each row is one PR, and each leaves the gate green.

| PR  | Carries                                                                                                                                                      | Notes                                                                                                                                                                |
| --- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | `Taskfile.yml`, `.golangci.yml`, `.github/workflows/ci.yml`, `coverage.ratchet`, `CLAUDE.md`, `ERRATA.md` — and the `dice` package with property tests       | The gate has to arrive with code it can actually run. `go.mod` currently declares a module with no Go files; `dice` closes that in the same change                   |
| 2   | `chargen`: the `Character` record, the `Event` log, the `Decider` interface, and the auto and interactive implementations — with no rules in them yet        | The log comes before the first rule so nothing is retrofitted into it. Tests are the log's own invariants: monotonic sequence, consequences referencing a real cause |
| 3   | Characteristics (FR1) and a minimal `cmd/cschargen new` that emits the JSON record                                                                           | End-to-end from here on. The six-rolls-then-assign choice is the first `Decider` call, so the choice-event path is exercised by the first rule                       |
| 4   | `career`: Colonist (pp. 173–176), Vagabond (p. 296), Prisoner (p. 261), plus the Injury and Life Events tables                                               | Data only, with a completeness test per career against the printed shape                                                                                             |
| 5   | The term loop (FR8): enlist, assignment, survival, mishap, advancement, skills, events                                                                       | The largest PR. Splittable at the skills/events boundary if review gets unwieldy                                                                                     |
| 6   | Muster out (FR11) and career transfer: the forced routes into Vagabond and Prisoner, and the `career_transfer_unimplemented` consequence for everything else | Where the closure above lives. Benefit rolls are a queue of scoped batches, not a count                                                                              |
| 7   | `render`: the Markdown sheet and the lifepath transcript; `cschargen replay`; golden fixtures                                                                | Replay lands last because it verifies everything before it                                                                                                           |

## Design decisions carried from `t5chargen`

Three things are worth copying rather than re-deriving, because they are proven
and the reasoning is already written down in that repo:

- **The `Decider` interface**, with `Choose(Choice) (int, error)` and a `Kind()`.
  An error refuses the choice and ends generation — an abandoned interactive
  session, or a replay whose recorded choice no longer matches what the engine
  offers — which is distinct from an out-of-range index, a decider answering
  wrongly. `Choice` carries `Nth`/`Of` so that a run of identical questions (the
  skill selections of one term) is answerable.
- **Four event kinds** — step, throw, choice, consequence — with exactly one
  payload non-nil, monotonic sequence numbers from 1, and consequences carrying
  the sequence number of their cause.
- **Throw events record the expression, the individual dice, the target, the
  itemized modifiers and the page cite.** The cite is what makes the log auditable
  against the book, which is the whole point of keeping it.

What is not copied: `t5chargen`'s characteristic roll is fixed-order, and this
book's is six rolls then a free assignment (p. 13). That difference is the first
place the engines diverge and it is load-bearing for the choice-event path.

## Done when

- `cschargen new --seed N --career colonist --tech-level 11 --max-terms 28`
  produces a JSON record and, with `render`, a Markdown sheet.
- The record's event log walks against the book: every throw carries the page it
  came from.
- `cschargen replay` reproduces the record exactly from its seed and recorded
  choices, and exits non-zero at the first divergence naming the sequence number.
- A character forced into Vagabond or Prisoner completes generation; a character
  sent toward Celebrity, Diplomatic Service or a free career choice ends with a
  `career_transfer_unimplemented` consequence naming the target, and the record
  is valid.
- `dice` has property tests; the three careers have completeness tests against
  their printed shapes (Colonist: 36 event entries, 11 mishaps, 6 rows per skill
  table, 7 benefit rows).
- `task` is green, and CI runs exactly `task`.

## Explicitly not in milestone 1

Steps 1 and 3–8 — species, subsector and homeworld, family, youth events, teenage
events, higher education. Aging (step 17) is milestone 4: it needs the homeworld
tech level that milestone 2 supplies from data, and a character who cannot serve
enough terms to reach an aging band gains nothing from it. The other
thirty-one careers are milestone 3.
