# POLICY.md

The auto-mode decision table. Every choice point the engine offers is
resolved here when `--auto` is set, and `policy_version` in each record
identifies which version of this document made the choices.

Version: **0.2.0**

## The rule

**Take the first option the book prints.**

Every table in the _Clement Sector Core Character Creation Book_ prints its
options in a fixed order, so "first listed" is a rule a reader can check against
the page rather than a preference this document is asserting. Where the engine
builds an option list that is not a printed table — the careers a character may
enter, say — the list is built in the book's own order and the same rule applies.

That is the whole policy. It is deliberately unrefined: a policy that chose
otherwise would need an argument per choice point, and two choice points already
need one — they are the two below the table.

## What the policy is not

- **It is not a strategy.** It does not try to build a capable character, and a
  character generated under it is not a recommendation.
- **It is not consulted on replay.** Replay reapplies recorded choices and never
  asks the policy anything, which is why a record made under one policy version
  replays under any other. `policy_version` is therefore recorded but not
  verified.
- **It is not total in the useful sense.** Every choice point has an entry, but
  the entry is almost always "the first option", which is a rule rather than a
  judgement. A policy that chose otherwise would need an argument per choice
  point.

## Every choice point

The engine offers these, and the policy answers every one. A choice point the
policy has no entry for is a build error, not a silent default.

| Point                    | Where                                                          | Policy                                                                                            |
| ------------------------ | -------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| `assign_characteristics` | p. 13, "apply the numbers ... as they see fit"                 | Assign in roll order: the first roll to STR, the second to DEX, and so on                         |
| `generation`             | p. 63, the engineered-human generation tables                  | The first generation the data offers                                                              |
| `genetic_outcome`        | p. 63, a compound outcome the character chooses between        | The first the data prints                                                                         |
| `uplift_class`           | p. 66, the class an uplift character is                        | The highest the homeworld's tech level allows, which is the one offered first                     |
| `subsector`              | p. 39, a chart entry the d100 roll did not reach               | The first subsector the data lists                                                                |
| `homeworld`              | p. 39, a choose-only subsector                                 | The first world the chart prints                                                                  |
| `primary_language`       | p. 41, a world with more than one                              | The first the chart prints                                                                        |
| `background_skill`       | p. 40, "or" between two background skills                      | The first the chart prints                                                                        |
| `background_specialty`   | p. 40, a background skill offering a choice of specialty       | The first the chart prints                                                                        |
| `youth_path`             | p. 68, "the player may choose any path for which they qualify" | Path 1, the first the book prints — which every character qualifies for                           |
| `teenage_path`           | p. 76, the same sentence again                                 | Whichever of Paths 1 and 2 the homeworld gives, which is always printed first                     |
| `higher_education`       | p. 85, "characters are not required to attend college"         | The first institution offered, and never "none" — so an eligible character attends                |
| `degree_field`           | pp. 87, 92, "the character may now take ... at level 2"        | The first skill the list prints                                                                   |
| `washout_skill`          | pp. 87, 92, "gains a level in the player's choice of"          | The first skill the list prints                                                                   |
| `medic_specialty`        | p. 101, "Specialties should all be different"                  | The first specialty still unspoken for, so the four end up different                              |
| `career`                 | p. 107, the career to attempt                                  | The first the book prints that the character may enter                                            |
| `assignment`             | p. 112, the assignment within a career                         | The first the career prints                                                                       |
| `commission`             | p. 114, "this step is optional"                                | Attempt it — the first option, and the one that opens the officer ranks                           |
| `skill_table`            | p. 117, which of a career's skill tables to roll on            | The first the career prints                                                                       |
| `other_assignment`       | p. 117, "an assignment table other than your own"              | The first that is not the character's own                                                         |
| `skill_specialty`        | pp. 116-117, a skill granted with "(Any)"                      | The first specialty the book prints for that skill                                                |
| `table_result`           | a result that offers a choice between effects                  | The first the page prints                                                                         |
| `relationship_target`    | p. 320, "an existing Ally or Contact loses 1d6 × 20"           | The first tie of a matching kind, in the order they were gained                                   |
| `raise_held_skill`       | "a level in any skill you already possess"                     | The first entry on the character's own sheet, which is alphabetical                               |
| `any_skill`              | "a level in any skill of your choice" (pp. 304-314)            | The first on the book's skill list, which is Admin                                                |
| `benefit_table`          | p. 127, Cash or Other Benefits                                 | Cash while p. 127's three remain, then Other — so the cap is spent before the table that has none |
| `spend_pool`             | pp. 188, 218, "+6 ... in increments of up to +3 at any time"   | The largest increment the result allows, every time — see below                                   |
| `aging_crisis`           | p. 123, "the player may pay 1D6 × 1000 credits"                | Pay for the treatment, which is the first option printed — and the only one a character survives  |

Two entries are not "the first option printed", because the page prints no list
for them to be first in.

**`spend_pool`** is the two results that grant a modifier the character spends
themselves. The options run from the largest increment down, so the policy
spends the pool rather than carrying it to the end of a career and wasting it.
A player answering the question may decline, and often should: spending +2 on a
throw that would have passed anyway is +2 gone.

**`benefit_table`** offers Cash first while p. 127's three-per-career cap has
room, so an auto character takes the cash rolls early. That is the same "first
printed" rule applied to a list the engine builds, and its consequence is that
an auto character's Other Benefits come from the rolls left over.

Rows are added as the engine reaches the choice points, in the PR that reaches
them. A row here and no code is as wrong as code and no row — which this table
went twelve PRs without honouring, and `policy_version` 0.2.0 is where it caught
up.

## Choices the book conditions on state the engine cannot read

A handful of results branch on something about the character that no effect can
ask: whether they held a career before this one, which assignment they are in,
how old they look. Where the branch is one the engine could evaluate — the
assignment held, say — the result is written as a choice and the Decider picks.
Where it is not, the branches are still written out, and under the policy the
first one printed wins whether or not it is the one that applies.

| Result                              | Conditions on                             | What the policy takes                                                             |
| ----------------------------------- | ----------------------------------------- | --------------------------------------------------------------------------------- |
| Fringe Marketer mishap 4 (p. 199)   | whether a career was held before this one | "Return to the previous career", which the engine records rather than carries out |
| Exotic mishap 4 (p. 191)            | the assignment held                       | Hetaerae's penalty, the first the page prints                                     |
| Sports events 14 and 56 (pp. 275-6) | the assignment held                       | Athlete's throw, the first the page prints                                        |

None of these is an errata reading: the book is clear about what it wants, and
the engine simply cannot see the state it wants read. They are listed here so a
generated record can be checked against the page by hand.

One consequence of that last row is worth stating plainly: **under the auto
policy a character always takes Youth Path 1**, because it is first and it is
always open. The other four are reached by a player choosing them, and by the
tests that check they are open to the right characters. That is the policy being
unrefined rather than the gates going unevaluated — the record still says which
paths were offered.

The education row has a visible consequence: **under the auto policy an eligible
character attends everything they qualify for.** Step 8 loops, because otherwise
Graduate School and Medical School are unreachable there — the degree they need
is one the character earns in the same step — and the policy never declines. So a
character with the characteristics for it walks out with a bachelor's, a
master's, a doctorate and an MD, sixteen years older.

That is the policy being unrefined rather than the rules being wrong: every one
of those throws was made and could have failed. A player chooses to stop, and
`--skip-education` is how a batch run says so.

## Two choices the policy is not offered

Step 18 lets a character "return to higher education" (p. 125) and Step 20
asks for a name, a gender identity, an appearance and a goal (pp. 129-130).
Neither is put to the auto policy at all.

**Returning to education** would be a question every term, and the policy takes
the first option — which would send every character back to school after every
term, forever. So it is offered to a player and to nobody else.

**Step 20's four fields** have no options to take the first of. FR12 says the
engine "records them and leaves them empty rather than inventing them in auto
mode", so an auto character's four fields are empty and the record says they
were left that way rather than saying nothing.

Both are the policy being unrefined rather than the rules being odd, and both
are the kind of choice this document exists to be honest about.

## After an Aging Crisis

A character who survives one "automatically fails all future Enlistment checks",
and one with a mental characteristic at 0 "may not attempt further Enlistment
checks" at all (pp. 123-124). Both leave three options: continue in the current
career if permitted, end character generation, or enter Vagabond.

The engine offers Vagabond. The other two are not the career list's to offer — it
is reached only with no career in hand, and the term limit is what ends
generation — and Vagabond takes no enlistment throw, so the automatic failure
never has to be rolled for. A player who wants to stop instead sets the term
limit.
