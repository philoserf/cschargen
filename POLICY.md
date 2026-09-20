# POLICY.md

The auto-mode decision table. Every choice point the engine offers is
resolved here when `--auto` is set, and `policy_version` in each record
identifies which version of this document made the choices.

Version: **0.1.0**

## The rule

**Take the first option the book prints.**

Every table in the _Clement Sector Core Character Creation Book_ prints its
options in a fixed order, so "first listed" is a rule a reader can check against
the page rather than a preference this document is asserting. Where the engine
builds an option list that is not a printed table — the careers a character may
enter, say — the list is built in the book's own order and the same rule applies.

That is the whole policy at milestone 1. It is deliberately unrefined: a policy
that chose otherwise would need an argument per choice point, and those arguments
belong with the choice points, which later milestones add.

## What the policy is not

- **It is not a strategy.** It does not try to build a capable character, and a
  character generated under it is not a recommendation.
- **It is not consulted on replay.** Replay reapplies recorded choices and never
  asks the policy anything, which is why a record made under one policy version
  replays under any other. `policy_version` is therefore recorded but not
  verified.
- **It is not total yet.** Milestone 1 covers the choice points the three
  implemented careers reach. A choice point the policy has no entry for is a
  build error, not a silent default.

## Choice points at milestone 1

| Point                    | Where                                                                     | Policy                                                                                           |
| ------------------------ | ------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| `assign_characteristics` | p. 13, "apply the numbers ... as they see fit"                            | Assign in roll order: the first roll to STR, the second to DEX, and so on                        |
| `aging_crisis`           | p. 123, "the player may pay 1D6 × 1000 credits"                           | Pay for the treatment, which is the first option printed — and the only one a character survives |
| `relationship_target`    | p. 320, "an existing Ally or Contact loses 1d6 × 20"                      | The first tie of a matching kind, in the order they were gained                                  |
| `youth_path`             | p. 68, "the player may choose any path for which the character qualifies" | Path 1, the first the book prints — which every character qualifies for                          |
| `teenage_path`           | p. 76, the same sentence again                                            | Whichever of Paths 1 and 2 the homeworld gives, which is always printed first                    |
| `higher_education`       | p. 85, "characters are not required to attend college"                    | The first institution offered, and never "none" — so an eligible character attends               |
| `degree_field`           | pp. 87, 92, "the character may now take ... at level 2"                   | The first skill the list prints                                                                  |
| `washout_skill`          | pp. 87, 92, "gains a level in the player's choice of"                     | The first skill the list prints                                                                  |
| `medic_specialty`        | p. 101, "Specialties should all be different"                             | The first specialty still unspoken for, so the four end up different                             |

Rows are added as the engine reaches the choice points, in the PR that reaches
them. A row here and no code is as wrong as code and no row.

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
