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

| Point                    | Where                                          | Policy                                                                    |
| ------------------------ | ---------------------------------------------- | ------------------------------------------------------------------------- |
| `assign_characteristics` | p. 13, "apply the numbers ... as they see fit" | Assign in roll order: the first roll to STR, the second to DEX, and so on |

Rows are added as the engine reaches the choice points, in the PR that reaches
them. A row here and no code is as wrong as code and no row.
