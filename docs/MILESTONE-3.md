# Milestone 3: the other thirty-one careers

2026-09-19. Status: plan. Tracks [#12](https://github.com/philoserf/cschargen/issues/12).

Milestones 1 and 2 transcribed three careers as hand-typed Go, on the argument
that the data format is what this milestone proves. This is that milestone, and
the first thing it has to do is settle the argument.

## The format decision, reversed

I said in the milestone 1 plan that thirty-four careers would probably want a
data file, and when this milestone was filed I said I expected to move them.
Having now written three careers and read a fourth, **I am keeping hand-typed
Go**, and the reason is the shape of an event rather than the number of them.

A table result is a tree. Colonist event 12 is a choice between two skill checks
whose success branches differ, and Prisoner event 13 is a check whose success
branch contains another check. In Go, with the constructors this package already
has, that reads:

```go
12: {
  Summary: "a disagreement between colonists turns into armed conflict",
  Effects: []Effect{pick("settle it by talking or by shooting",
    opt("Diplomat", checkSkill("Diplomat", 8,
      []Effect{skill("Leadership")},
      []Effect{injury(1), relationship(Enemy, 1, "")})),
    opt("Gun Combat", checkSkill("Gun Combat", 8,
      []Effect{skill("Recon")},
      []Effect{injury(1), relationship(Enemy, 1, "")})),
  )},
},
```

The JSON for that is four times the lines and reads nothing like the page. The
argument for a data file was that thirty-four careers is three thousand lines of
Go — which is true, and which is also three thousand lines of table either way.
The transcription is the cost, not the format, and the format that reads closest
to the page is the one where a transcription error is likeliest to be _seen_.

Two things the compiler gives that a schema would have to re-implement: an effect
kind added to the vocabulary and not handled is a build failure, and a mistyped
field is a build failure rather than a value silently absent.

**What would change my mind:** a supplement's careers arriving as user data, the
way world names do. That is a Product Identity question and the careers are not
Product Identity, so it has not arisen.

## The careers

Thirty-one, at pp. 146-300. Three are already built (Colonist, Prisoner,
Vagabond), and one — the career for an enslaved engineered person — is renamed
from the book's term for the same reason everything else is (see `CLAUDE.md`).

## The shapes the first three do not have

- **Commissions.** National Navy (p. 233) and the military careers carry a
  Commission target, and Step 14's commission branch (p. 114) is written and has
  never been reached. It is optional, attemptable once per term, and free to
  fail.
- **Officer Skills.** A fifth career-wide table, open only to the commissioned
  (p. 117).
- **The Military Events table** (p. 121), which several careers' event tables
  route to and which is transcribed nowhere yet.
- **Rank titles** (p. 115, and the chart at p. 238). Descriptive, not mechanical;
  carried as data and printed on the sheet.

## The stubs this closes

- `EffectTransfer` to a career that does not exist currently **ends generation**.
  With thirty-four careers it becomes an ordinary transfer, and the "generation
  ends here" path should become unreachable for any career the book names.
- Colonist event 54 rolls twice on Diplomatic Service's Ambassador table (p. 185);
  Vagabond event 12 enters a System Defense Force career; Prisoner event 14 sends
  the character to the troopers.
- **ERRATA E-6** — "choose another career" as a choice point rather than a named
  transfer — is recorded and deliberately not applied, because with three careers
  there was nothing to choose between. This milestone applies it.

## PR sequence

Careers are grouped so that each PR carries one new mechanism as well as its
share of transcription, rather than saving every hard case for the end.

| PR  | Carries                                                                                                                                                                  |
| --- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 1   | This plan, the Military Events table, the commission and Officer Skills mechanics, and the first military career (National Navy, p. 233)                                 |
| 2   | The three System Defense Force careers and Marine (pp. 224, 278-291) — the rest of the commissioned careers                                                              |
| 3   | The trades and services: Arts, Belter, Craftsperson, Corporate Shipper, Instructor, Medic, Orbital Construction, Scientist                                               |
| 4   | The public lives: Celebrity, Clergy, Diplomatic Service, Journalist, Politician, Sports, Exotic, Adventurer                                                              |
| 5   | The frontier and the underworld: Explorer, Fringe Marketer, Gambler, Independent Merchant, Investigator, Organized Crime, Pirate, Scavenger, Thief, and the slave career |
| 6   | Close the stubs: named transfers resolve, E-6 applied, the completeness test extended to all thirty-four                                                                 |

## Done when

- All thirty-four careers generate, and a character can reach any the rules allow.
- The commission branch and the Military Events table are reached by a generated
  character, not merely compiled.
- No `EffectUnimplemented` remains whose reason is "this career is not
  implemented".
- E-6 is applied, or its entry says why it still is not.
- The completeness test holds every career's printed shape, and the gate is green.
