# Milestone 5: the pre-career steps

2026-09-19. Status: done. Tracks [#24](https://github.com/philoserf/cschargen/issues/24).

## What it turned out to be

Six PRs in the order planned. Five things the plan did not know:

1. **A tie never crosses zero.** p. 320 says so from three sides, and the obvious
   implementation -- move the number, look up the band -- gets exactly those
   cases wrong. An Ally at 110 who takes -120 is gone; they are not a Rival at
   -10. That is E-24, and it was found by a table-driven test rather than by
   reading.

2. **The engine granted level-0 skills at level 1.** "Gain Streetwise 0" is not
   "gain a level in Streetwise", and an unset `Level` meant one. Belter event 24
   has said "gain Medic at level 0" since milestone 3 and has been granting Medic
   1 the whole time.

3. **Step 8 has to loop.** Graduate School and Medical School require a degree
   the character earns in the same step, so a single pass makes them unreachable.
   The loop then raised a question the book does not (E-29) and made the auto
   policy collect four degrees, which POLICY.md now states plainly.

4. **The apparent-age chart's two ends want opposite readings**, which is E-20 in
   milestone 4 but was found here, while writing the teenage gates.

5. **Eleven tests moved rather than broke**, and each move says something. The
   Step 2 tests skip the three pre-career steps because those legitimately change
   characteristics. Two Colonist tests scanned for a magic seed that stopped
   enlisting once a step was added ahead of Step 10. The muster-out test learned
   that a character who died in service did not leave a career. The age test
   learned that a degree costs four years on a different page.

Eight errata: E-22 to E-24 on the ratings, E-25 on a sibling's age, E-26 on which
characteristic opened a path, E-27 on the campaign's present year, and E-28 and
E-29 on the two schools.

## The plan as filed

Steps 5 to 8 (pp. 57–104): family, youth events, teenage events, and the four
higher-education tracks. The PRD calls it "the largest milestone by volume and
the one that can be deferred longest without the tool being useless", and both
halves of that are still true — milestones 1–4 generate a human adult with a
career history, an aging record and a cause of death where there is one.

About half of those forty-seven pages are species content: the engineered-human
genetics tables (pp. 62–65), uplift classes and origins (pp. 65–67), and the
uplift and enslaved variants of each event path. Those are milestone 6, and two
of their headings are proper names the OGL notice reserves. What is left is the
human path through each step.

## The thing that has to land first

**Relationship Ratings are a mechanic this engine does not have.** p. 320
defines them:

> Contacts are NPCs with a Relationship Rating of 1-100 ... NPCs that have a
> Relationship Rating of 101-200 are Allies ... Rivals are NPCs with a
> Relationship Rating of -1 through -100 ... Enemies are NPCs with a Relationship
> Rating of -101 or less.

and every band has an exit: a Contact that falls to 0 is lost, one that rises past
100 becomes an Ally, an Enemy that rises above -100 becomes a Rival.

The engine creates **every tie at a rating of 0**, which by that rule means every
tie a character has ever gained is one the page would say they had already lost.
ERRATA E-5 recorded, in milestone 1, that a rating had to exist because four Life
Events write to one; it exists as a field and nothing sets it.

Steps 5, 6 and 7 cannot be transcribed on top of that. The family step grants
parents at 125 and grandparents at 100; Youth Path 1 alone moves a named parent
by -150, everyone in the character's life by -25, every family member by +20, and
both parents by +25. A rating that is always zero makes all of those no-ops.

So the first PR is the ratings, not the family.

### What the ratings need

- A starting rating per tie, from the result that granted it.
- **The kind is stored, not derived.** p. 320's band sentences are all motion —
  "Contacts which fall to a Relationship Rating of 0", "Allies whose Relationship
  Rating drops to 1-100 become Contacts" — and the tables that grant a tie name
  its kind outright. Deriving the kind from the rating would make every
  grandparent a Contact the moment they were created, because Step 5 grants them
  "as an Ally with a Relationship Rating of 100" and 100 is the top of the
  Contact band. The bands govern movement; the granting result governs birth.
- Band transitions applied whenever a rating moves, including the two ways a tie
  is lost outright.
- A tie that can be **named**, because the tables address them: "choose one of
  your parents", "everyone in your life", "all family members". A rating that
  cannot be aimed is a rating the youth tables cannot use.
- The four Life Events that write to ratings (p. 120), which are recorded as
  unimplemented today.

**Three numbers that disagree, in one errata entry.** p. 320 says "Family members
will always begin as Allies with a Relationship Rating of 150". Step 5 grants
parents and siblings at 125 and grandparents, aunts, uncles and cousins at
100 — and 100 is a Contact by p. 320's own bands, not an Ally. The reading is to
follow the result that grants the tie: the specific instruction over the general
rule, and the named kind over the band.

**And one the book does not give at all.** Most career results say "gain an Ally"
with no rating, and the engine has to start them somewhere. That is a reading too,
and distinct from the ratings already written into `career/shared.go` where the
page does name one (140, 180).

## Step 5: family (pp. 57–61)

The book opens it by saying it is optional — "determine as much or as little
about the character's family as they wish ... Skipping this step can speed up
character generation". So the engine's default is to run it and a flag turns it
off, rather than the other way around.

Four tables: the Human Birth Situation chart (d100 with world modifiers), the
parental age expression by tech level, the sibling count and the 2d6 sibling age
table. Then grandparents, aunts and uncles by repeating the birth situation for
each parent, and cousins by repeating it again.

**Two things bound it.** The recursion is unbounded in the book — "this process
can then be repeated to determine great-grandparents and further ancestors ... as
far as one wishes" — so the engine stops at grandparents, aunts, uncles and
cousins, which is where the book stops giving specific instructions. And the
birth-situation modifiers are a list of world names, every one of them Product
Identity: they belong in the setting data file beside the origin tables, as a
per-world field.

## Step 6: youth events (pp. 67–75)

Five 2d10 paths of nineteen rows, gated on characteristics, rolled twice — ages
4–8 and 9–12 — with a re-qualification between the two rolls if the first changed
a characteristic. Plus the Youth Life Events table.

New vocabulary: a skill granted at **level 0** (the tables say "Gain Streetwise 0"
where a career table says "gain a level in Streetwise"), and a rating adjustment
aimed at a named tie or a group of them.

## Step 7: teenage events (pp. 75–85)

Four paths, two of them gated on how long the homeworld has been settled — which
the setting data already carries as `settledYear`. Plus the Teenage Life Events
table, whose result 5 is a relationship upgrade across every band at once and
whose result 6 offers it as one of four choices.

## Step 8: higher education (pp. 85–104)

Four tracks — undergraduate, military academy, graduate school, medical school —
each with a prerequisite, an admission throw, a success throw, an honours throw, a
failure table and an events table.

**It is not strictly a pre-career step.** p. 85 says university "may be entered at
any point after the character reaches the age of 18 ... Some characters may attend
immediately after their teenage years, while others may return later in life after
military service", and Step 18 (p. 125) lists "return to higher education" among
what a character decides between terms. So the education step hooks into the term
loop as well as sitting before Step 9. This is the half of the milestone that feeds
back into work already done: the `UndergraduateDegree`, `GraduateDegree` and
`MedicalSchool` enlistment modifiers are recorded as unimplemented in every career
that carries one, exactly as `ApparentAgeOver40` was before milestone 4.

**One thing the setting data does not carry.** "Any world with a population code
of 4 or higher will have an undergraduate college present on the planet" (p. 86).
The world schema has no population code. Either it gains one, or the reading is
that a character may always attend, on their homeworld or "a nearby world" as the
same page allows.

## Two things to settle before the tables, not during them

**There is no d10 in this repository.** `dice` has D6, D3, D66 and D100, and
`rollExpression` accepts six and three sides. The youth paths are 2d10, sibling
age is 2d6 and 3d10, and parental age is 1d10 whose **face** carries meaning — a 1
makes the character the firstborn and a 10 the last. `D10` and the parser's
support for it land before any table that needs them.

**The present day is 2350 on p. 43 and 2345 on p. 124.** Teenage paths 1 and 2 are
gated on whether the homeworld was colonized more than a hundred standard years
ago, which needs a reference year to measure `settledYear` against. One of the two
is picked, cited, and the other recorded.

## PR sequence

| PR  | Carries                                                                 |
| --- | ----------------------------------------------------------------------- |
| 1   | This plan, Relationship Ratings as a mechanic, and the four Life Events |
| 2   | Step 5, the family, and the setting data field its modifiers need       |
| 3   | Step 6, the five youth paths and the Youth Life Events table            |
| 4   | Step 7, the four teenage paths and the Teenage Life Events table        |
| 5   | Step 8, undergraduate and the military academy                          |
| 6   | Step 8, graduate and medical school, and the three enlistment modifiers |

## Done when

- A tie's kind is a function of its rating, ties are lost at the two boundaries,
  and no tie is created at a rating the page would call lost.
- The four Life Events that move ratings do so.
- A generated character has a family, a childhood and a youth, and the sheet
  shows them.
- `UndergraduateDegree`, `GraduateDegree` and `MedicalSchool` are evaluated
  rather than recorded, and no career's enlistment carries an unimplemented
  consequence for them.
- Every path table has a second-reading transcription test, and every table a
  generated character can reach is reached by one.
- `task` is green and CI runs exactly `task`.

## Explicitly not in milestone 5

The engineered-human genetics tables (pp. 62–65), uplift classes and origins
(pp. 65–67), and the uplift and enslaved variants of the youth and teenage
paths. They are milestone 6, and their names come from the setting data file.
