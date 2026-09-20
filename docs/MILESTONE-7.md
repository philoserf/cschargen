# Milestone 7: finishing

2026-09-20. Status: plan. Tracks [#38](https://github.com/philoserf/cschargen/issues/38).

The last milestone before the PRD's v1: Step 20, the interactive mode the PRD
has promised since its second goal, batch generation, and the things earlier
milestones deferred to here.

## What is actually missing

Five milestones have built the twenty steps and left four things:

1. **Step 20** (pp. 129–130). Name, gender, appearance, personal goals. None is
   rolled and none is mechanical: "All player input or policy default … The
   engine records them and leaves them empty rather than inventing them in auto
   mode" (FR12).
2. **Interactive mode.** `Decider` has had three kinds since milestone 1 —
   `player`, `policy`, `replay` — and only two exist. `new` refuses without
   `--auto` and says so.
3. **Batch.** `cschargen batch --count 20 --auto` emitting JSONL, each member's
   seed derived from the base plus index and recorded.
4. **Return to higher education at Step 18**, deferred from milestone 5 with the
   reason that it is a choice only a player can meaningfully make (p. 125 lists
   it among what a character decides between terms).

## The order matters

Interactive mode first, because it is what makes the other three worth having.
Step 20 in auto mode is four empty fields; Step 20 with a player is the point of
the step. Returning to higher education under the auto policy is one more thing
the first option takes; with a player it is a decision.

## Interactive mode

`Player` is a `Decider` that prints the choice and reads an answer. The engine
needs nothing new: every choice already goes through the interface, carries a
prompt, a list, a page cite and — where the choice is one of a series — which of
how many it is.

What it does need is care about **where the prompt goes**. `new` writes the
record to stdout so that it can be piped, so the prompts go to stderr and the
answers come from stdin. A run that is not a terminal on both is a run that
cannot ask, and the engine says so rather than hanging.

The transcript already records who decided each choice, so an interactive record
replays exactly as an auto one does.

## Step 20

Four fields, recorded and never invented. `--name` exists; the other three want
flags too, and in interactive mode they are asked for. A setting data file may
supply name lists, and where it does auto mode may draw from them — which is the
one exception FR12 makes and the only place in the engine where a name comes from
anywhere but the player.

## Batch

`--count N`, `--auto` required, JSONL out. The seed of member _i_ is the base
seed plus _i_, recorded in that member's provenance, so any one of them can be
regenerated alone. A batch is not a new generator: it is a loop over the one that
exists.

## Done when

- `cschargen new` without `--auto` walks the twenty steps asking the player, and
  the record it writes replays.
- Step 20's four fields are on the record and the sheet, empty where nobody
  supplied them.
- `cschargen batch --count 20 --auto` writes twenty records, each regenerable
  from its own recorded seed.
- A character may return to higher education between terms.
- `task` is green and CI runs exactly `task`.

## PR sequence

| PR  | Carries                                                               |
| --- | --------------------------------------------------------------------- |
| 1   | This plan, interactive mode, and the prompts going to the right place |
| 2   | Step 20, and the name lists a data file may supply                    |
| 3   | Batch, and returning to higher education at Step 18                   |

## What v1 is when this lands

The PRD's first run: a tool that generates a character of any species the data
declares, through twenty steps, either by asking or by policy, and writes a
record that replays. What it is not is a tool that knows the published setting —
that file is the reader's to write, and the repository ships an invented sample
so the tool runs without one.
