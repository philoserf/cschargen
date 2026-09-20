# Milestone 6: species

2026-09-20. Status: plan. Tracks [#33](https://github.com/philoserf/cschargen/issues/33).

Step 1 (pp. 21–38) and everything the previous five milestones deferred because
it belonged to a species rather than to a human: the characteristic methods, the
caps, the aging tables, the genetics, the uplift classes, and the youth and
teenage variants.

## This is what the Product Identity decision was made for

Every species-specific table in the book is **headed by a species name**, and
every one of those names is a proper name the OGL notice reserves (p. 335). The
engine cannot carry them, and the setting data file already carries the species
list — four invented ones in `setting/sample.json`.

So the design decision this milestone turns on is the one to make first:

**A species names a profile; a profile is named by what it does.** The engine
holds the mechanisms and the tables that are not species-specific. A species in
the data file says which characteristic method it uses, which aging profile it
ages on, what skills it starts with, and what its ceiling is. No engine file
names a species and no data file re-implements a rule.

That is the same shape the origin tables already have: the engine knows what a
homeworld is for, and the file says which ones exist.

## The five aging profiles

pp. 122–123 print seven tables, and they collapse to five once the headings are
set aside. Four of the five engineered species share the human ones — the
headings read "Humans, Gaishan, Oskars, Aquans and Sniffers from a World that is
Tech Level 10" — so the tables already built cover most of the book.

| Profile     | Bands (term → checks)                                                   | Who the book gives it to           |
| ----------- | ----------------------------------------------------------------------- | ---------------------------------- |
| `techLevel` | the four tables of milestone 4, keyed by the homeworld's tech level     | humans and four engineered species |
| `rapid`     | 4–5, 6–7, 8+                                                            | one uplift                         |
| `moderate`  | 6–7, 8–9, 10+                                                           | two uplifts                        |
| `sturdy`    | 6–8, 9–11, 13+                                                          | two uplifts                        |
| `sudden`    | 7, 8, 9+ — and the only table that opens at five checks and reaches 11+ | one engineered species             |

The names describe onset and severity, which is what the rows differ in. They
carry no species and a data file picks one by name.

**`sturdy` has a hole.** Its bands are 6–8, 9–11 and 13+, and term 12 is printed
nowhere. That is an errata entry and a reading: the engine holds the last closed
band through the gap, because the alternative is a term that ages nobody on a
table whose whole shape is that ageing accelerates.

## What each species has to say about itself

- **Characteristic method.** "All generations of Gaishan characters should roll
  2d6-2 for STR and END. DEX should be rolled as 2d6+2. CHA should be rolled as
  1d6+3. INT and EDU should be rolled normally" (p. 23). So: a dice expression
  per characteristic, defaulting to the human 3d6-drop-lowest where the file is
  silent.
- **Starting skills.** "All Gaishan should be given Survival (Freefall) and
  Survival (Low Gravity) at level 1" (p. 23).
- **Aging profile**, by name from the table above.
- **A characteristic ceiling.** `HumanMaximum` is 15 and applies to everything
  the engine generates. p. 14's wording is about "an unaltered human", so an
  altered one needs its own.
- **Uplift class**, which is not a species fact but a homeworld one: Class 1
  needs tech level 10, Class 2 needs 11, Class 3 needs 12 (p. 66), and "we
  highly recommend that the highest class available be used."
- **Youth and teenage roll patterns.** Uplifts roll a different number of times
  over different ages: one roll for ages 2–4 for one species, two for most
  (pp. 68, 76).

## Two things that are not species facts

**The homeworld permission is already there.** `setting.World` has carried
`engineered` and `uplifts` permissions since milestone 2, including the "slave"
status that three worlds in the book print. Nothing in the engine reads them
yet, and Step 4 is where it should: a character of a species their homeworld
bars cannot be born there.

**The slave career has no way in.** Its enlistment note says "automatically
enlisted by result of early life tables", and those tables are the enslaved
youth and teenage paths of pp. 74 and 84 — which are milestone 6 because they
are species content. Until they land, the career is transcribed and unreachable.

## The genetics tables

pp. 62–65 are the generation, compound and hybrid rules: a purebred first or
later generation, a compound of two engineered species, a hybrid with a baseline
human. The tables that follow are one per species and their rows are
cross-references to other species by name — the most Product Identity-dense
pages in the book.

They belong in the data file as a per-species table, and the sample's four
invented species get invented ones. The engine holds the three categories and
the rule that a compound crosses two species and a hybrid crosses one with a
baseline human.

## Reach, to check before writing engine code

The sample data has four species — two engineered, two uplifts — and no
character has ever been generated as one, because `Inputs.Species` is a string
the engine compares against "human" and nothing else. Before any of the above:

1. The sample's four species must exercise every profile, method and pattern the
   format can express, the way the aging reach check worked in milestone 4.
2. A test must generate a character of each and walk their record.

## PR sequence

| PR  | Carries                                                                                       |
| --- | --------------------------------------------------------------------------------------------- |
| 1   | This plan, the species data format, characteristic methods, caps, and the five aging profiles |
| 2   | Step 4's homeworld permissions, and uplift classes                                            |
| 3   | The uplift youth and teenage roll patterns                                                    |
| 4   | The enslaved youth and teenage paths, and the way into the slave career                       |
| 5   | The genetics tables: generations, compounds and hybrids                                       |

## Done when

- A character can be generated as any species the data file declares, and the
  sheet says which.
- Characteristics are rolled by that species' method and capped at its ceiling.
- Aging reads the species' profile, and a test reaches every profile.
- A character cannot be born on a world their species is barred from.
- The slave career is reachable.
- No engine file names a species, and `TestTheSampleIsInvented` still passes.
- `task` is green and CI runs exactly `task`.

## Explicitly not in milestone 6

The environmental task modifiers ("+2 to any physical task performed in gravity
of less than 0.50 standard") and the equipment cost multipliers. Both are play
rather than character generation, and neither appears on a character sheet.
