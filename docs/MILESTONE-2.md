# Milestone 2: the setting data, and where a character comes from

2026-09-19. Status: plan. Tracks [#9](https://github.com/philoserf/cschargen/issues/9).

Milestone 1 takes the homeworld's two load-bearing numbers — tech level and term
cap — as flags. This milestone replaces that seam with the external file the
Product Identity boundary requires, and implements Steps 3 and 4 (pp. 39–56).

## What the book actually asks for

**Step 3** is one throw: 1d6 for the subsector of origin — Earth, Hub, Cascadia,
Franklin, Sequoyah, Other (p. 39). The referee may choose instead, or let the
player choose.

**Step 4** is a d100 on that subsector's chart, or a choice (p. 39). Every chart
has the same seven columns, which is the whole schema:

| D100 | Planet of Origin | Background Skills | Primary Language | Maximum Age/Terms | engineered-human and uplift Status | Tech Level |

A seventh table, Recently Colonized Worlds (p. 55), has the same columns and one
extra rule: "you cannot randomly be assigned one of these worlds, you may choose
them though maximum ages can be quite low" (p. 40).

## Three columns need more structure than the page shows

**Background skills are an expression.** The page writes them as prose that mixes
two operators and a specialty choice:

- `Suit (Vacc Suit)` — one skill, one named specialty.
- `Suit (Vacc Suit) or Survival (Freefall)` — choose one.
- `Survival (Mountains, Ocean, or Plains)` — one skill, choose the specialty.
- `Suit (Vacc Suit) or Survival (Freefall) and Language (Māori, English, or Hindi)`
  — choose one of the first two, **and** take the Language.
- `Survival (High Pressure, Mountains, Plains, or Ocean) and Electronics (Computers)
and mindcomp` — Kingston (p. 41), where the last item is not a skill at all.

Parsing that prose would be a second place the rules live. The format encodes it
structurally instead: a list of requirements, each a choice between alternatives,
each alternative a skill with the specialties it offers — or a named item.

**Engineered human/uplift status is not a boolean.** Entries read "Yes/Yes, Free/Free" but
also "Achilles banned. Others are allowed and free." So it is a per-kind
permission, a free-or-enslaved status, and a list of banned species.

**Settlement age is not a printed column at all.** Two of the teenage-event paths
gate on whether the homeworld "has been colonized or established for 100+ standard
years" (p. 76), which is milestone 5's problem but this milestone's schema. The
file carries it explicitly.

## The file

One file, two sections, one content hash, one validator, one `--data` flag.
Species live beside worlds because both are Product Identity and splitting them
would mean two hashes to verify and two ways to be out of step.

```json
{
  "schemaVersion": 1,
  "name": "a name for this data set",
  "species": [{ "name": "…", "kind": "engineered human" }],
  "subsectors": [
    {
      "name": "…",
      "originRoll": 2,
      "worlds": [
        {
          "name": "…",
          "roll": [1, 12],
          "techLevel": 11,
          "maximumAge": 160,
          "maximumTerms": 36,
          "settledYear": 2215,
          "primaryLanguages": ["…", "…"],
          "backgroundSkills": [{ "oneOf": [{ "skill": "Suit", "specialties": ["Vacc Suit"] }] }],
          "engineered humans": { "allowed": true, "status": "free", "banned": [] },
          "uplifts": { "allowed": true, "status": "free", "banned": [] }
        }
      ]
    }
  ]
}
```

A subsector with no `originRoll` is choose-only, which is what Recently Colonized
Worlds is. A world with no `roll` is choose-only within its subsector.

## What the engine does with it

- **Step 3** throws 1d6 against the subsectors' `originRoll` values, or offers a
  choice.
- **Step 4** throws d100 against the chosen subsector's worlds, or offers a choice.
- **Background skills** are granted at level 1, each requirement resolved through
  the Decider where it offers a choice (p. 40).
- **Primary language** is the first listed, or a choice; where a homeworld also
  grants the Language skill, its specialty must differ from the primary (p. 41).
- **`Electronics 0`** goes to every character regardless of origin (p. 41).
- **Maximum terms** caps the loop: "If the character has reached the maximum
  number of terms allowed by their homeworld, the character must end character
  generation" (p. 125). The policy term limit and the homeworld cap are both
  ceilings, and the lower wins.
- **`EffectNewHomeworld` stops being a stub.** Six of Colonist's eleven mishaps
  reassign the homeworld, so the record holds a homeworld history — and the new
  world's tech level and term cap take effect from that moment.

## What it does not do

- **Engineered humans and uplifts** are milestone 6. The permission and slavery columns are
  carried in the data and validated, and the engine reads them only far enough to
  say a human is always allowed.
- **The Earth Subsector rule** — "All characters which begin in Earth Subsector …
  must spend a minimum of four terms in Clement Sector" (p. 43) — needs the engine
  to know where a term was served, which nothing yet records. Recorded as an
  unimplemented consequence for a character whose origin subsector is flagged
  `offworld`.
- **Aging** stays milestone 4, though this milestone is what supplies the tech
  level it is gated by.

## PR sequence

| PR  | Carries                                                                                                    |
| --- | ---------------------------------------------------------------------------------------------------------- |
| 1   | The `setting` package: types, loader, validator, content hash, the invented sample, `data validate`        |
| 2   | Steps 3 and 4 in the engine, the `--data` flag, the hash stamped and verified on replay, homeworld history |

## Done when

- A user-supplied file drives Steps 3 and 4, and its hash is in every record and
  checked on replay.
- `cschargen data validate` reports a malformed file by what is wrong with it,
  not by panicking halfway through a lifepath.
- A character's homeworld sets their background skills, primary language, tech
  level and term cap, and a mishap that reassigns it changes all four.
- `setting/sample.json` is invented, and no world, subsector or species name
  from the book is in the repository.
- The gate is green.
