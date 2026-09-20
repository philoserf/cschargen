# Milestone 4: aging

2026-09-19. Status: done. Tracks [#22](https://github.com/philoserf/cschargen/issues/22).

## What it turned out to be

Three PRs in the order planned. Four things the plan did not know:

1. **`Characteristics.AllPhysicalZero` had a test and no production caller.**
   p. 14 states the death rule the aging section restates on p. 124, and the
   method was written for it in milestone 1 and never called. It is called now.

2. **p. 14 also settles half of what E-17 was going to be a reading about.** The
   aging section says "physical" and "mental" and enumerates neither, but p. 14
   names the physical three in passing — "(Strength, Dexterity, and Endurance)" —
   and six minus three is three. The entry is a cross-reference rather than an
   interpretation.

3. **The apparent age chart's top end wanted the opposite reading from its
   bottom.** "Nearest printed band" is right above 290 and wrong below 30: every
   column's first row is 20-25, so a character of 18 would be made _older_ by a
   chart whose whole purpose is to make them younger. Below 30 apparent age is
   the character's age, which is what p. 124 already says for tech level 9 and
   below.

4. **Only one golden moved, and it moved by one line.** The apparent-age stub
   the twelve careers used to record disappeared from the transcript, and the
   event numbering shifted beneath it. No aging throw reached a golden fixture,
   because the sheet's character is from Ironpsalm — tech level 9, but three
   terms in.

Five errata: E-15 and E-16 on the two tables that stop and the tech levels above
them, E-17 to E-19 on the crisis, and E-20 and E-21 on the chart and the bar
twelve careers measure against it.

## The plan as filed

Step 17 (pp. 121–125) is four pages: a per-term age increment the engine already
does, seven aging tables of which four are human, an Aging Crisis with five
terminal states beneath it, and an apparent-age chart that twelve careers'
enlistment throws have been waiting for.

## What the engine already has

`Generator.age` adds four years a term, or 1d3 where a mishap ejected the
character mid-term (p. 121). That is the whole of Step 17 today, and the comment
above it says so.

Twelve careers carry an `ApparentAgeOver40` enlistment modifier that the engine
records as unimplemented, naming this milestone. That is the visible reason to
do aging before the pre-career ninety pages: it closes a stub that already
exists in thirty-four careers rather than opening a new front.

## The tables are indexed by term, and the term is the lifetime's

The column is headed "Term", and for TL 9 and below it reads 6–8, 9–10, 11, 12+
— plainly term numbers. The higher tech levels read 18–28, 29–38, 39–48, 49+
(TL 10), 33–48 and 49–58 (TL 11), and 44–58 (TL 12–13), which are term numbers
too: at four years a term, term 44 is a character around 194, and the apparent
age chart runs to 290.

The index is **terms served in total**, not terms in the current career, and not
derived from age — the 1d3 mishap increment already means age is not 18 + 4n.

## The four human tables

| Homeworld TL | Terms | Checks                                              |
| ------------ | ----- | --------------------------------------------------- |
| 9 or less    | 6–8   | STR 8+, DEX 8+, END 8+                              |
|              | 9–10  | STR 9+, DEX 9+, END 9+                              |
|              | 11    | STR 9+, DEX 9+, END 9+, INT 9+, CHA 9+              |
|              | 12+   | STR 10+, DEX 10+, END 10+, INT 10+, EDU 10+, CHA 9+ |
| 10           | 18–28 | STR 8+, DEX 8+, END 8+                              |
|              | 29–38 | STR 9+, DEX 9+, END 9+                              |
|              | 39–48 | STR 9+, DEX 9+, END 9+, INT 9+, CHA 9+              |
|              | 49+   | STR 10+, DEX 10+, END 10+, INT 10+, EDU 10+, CHA 9+ |
| 11           | 33–48 | STR 8+, DEX 8+, END 8+                              |
|              | 49–58 | STR 9+, DEX 9+, END 9+                              |
| 12–13        | 44–58 | STR 8+, DEX 8+, END 8+                              |

The other three tables on pp. 122–123 belong to species this engine does not
generate yet, and two of their headings are proper names the OGL notice reserves
(p. 335). They arrive with milestone 6, out of the setting data file, which is
where every species name in this repository lives.

**TL 11 and 12–13 stop.** TL 9 and TL 10 each end on an open band — "12+",
"49+" — and the other two do not: nothing is printed for a TL 11 character past
term 58, or a TL 12–13 character past term 58. The engine makes no check there,
which is what the page says, and an errata entry records the reading. No world in
the sample data allows 58 terms, so it is not reachable there either.

**Nothing is printed above TL 13**, and the validator accepts up to TL 20. A TL
14 homeworld uses the TL 12–13 table: each band delays aging and none abolishes
it, so reading the highest printed band forward is the reading the trend
supports. That is an errata entry too.

## Reach, checked against the sample data before writing code

| TL    | First aging term | Sample worlds that reach it                    |
| ----- | ---------------- | ---------------------------------------------- |
| 9     | 6                | Ironpsalm (19 terms), Lowwater (17)            |
| 10    | 18               | Sablecourt (27), Pellucid (25), two more       |
| 11    | 33               | Harrowmere (37), Quillon (33)                  |
| 12–13 | 44               | **none** — the highest cap in the sample is 40 |

So three of the four tables are reachable by a generated character and the fourth
is covered by transcription tests only. That is worth knowing before the feature
ships rather than after: a table no seed reaches is the shape the credit-expression
bug had in milestone 3.

## The Aging Crisis and what sits under it

A characteristic reduced to 0 by aging (p. 123):

- **Crisis.** Pay 1d6 × 1000 credits and the characteristic is restored to 1.
  Without payment the character dies. Two characteristics reaching 0 in the same
  term is two payments, which the page does not say and an errata entry does.
- **After surviving one:** every future enlistment check fails automatically. The
  character may continue in the current career if permitted, end generation, or
  enter Vagabond.
- **Two or more physical at 0, unrestored:** incapable of independent movement,
  generation ends.
- **All three physical at 0:** dead.
- **One mental at 0:** no further enlistment; continue or Vagabond.
- **Two or more mental at 0:** dead.

**Which three are mental is a reading.** The book says "physical" and "mental"
and never enumerates either; STR/DEX/END and INT/EDU/CHA is the Cepheus
convention and what the aging tables themselves group.

**Death is a value, not an error** — the `dusk` rule about polar geometry, applied
here. The record is complete and valid, the sheet says the character died at the
age they died, and the CLI exits 0. Only a misused command line and a malformed
data file exit non-zero.

## Apparent age

The chart (p. 125) maps an actual-age band to an apparent-age band by tech level.
Below TL 10 apparent age is actual age. The chart begins at 30 and ends at 290;
outside that range it says nothing, and the engine reads the nearest printed band.

The generator stamps the band into the record at each Step 17, because apparent
age is derived from setting data the renderer does not have.

`ApparentAgeOver40` then becomes evaluable. Two readings: the chart gives bands,
so "over 40" is a band whose lower bound is 40 or more; and Sports (p. 273) prints
"40+" where the other eleven print "over 40", which is one modifier and one bar.

## The dice stream moves

An aging throw is a throw, so every seed that reaches term 6 on a TL 9 world
rolls differently from here on. Golden files and the coverage ratchet will move.
Regenerate them deliberately and read the diff.

## PR sequence

| PR  | Carries                                                                             |
| --- | ----------------------------------------------------------------------------------- |
| 1   | This plan, the four tables as data, the aging throws, and their transcription tests |
| 2   | The crisis, the five terminal states, and what Step 18 offers after one             |
| 3   | Apparent age, the chart, and the twelve careers' enlistment modifier evaluated      |

## Done when

- A generated character on a TL 9 world loses a characteristic to age, and a test
  says so rather than passing vacuously.
- A character can die of aging, and the record is valid, the sheet says so, and
  the CLI exits 0.
- `ApparentAgeOver40` is evaluated rather than recorded, and no career's
  enlistment carries an unimplemented consequence for it.
- The four human tables have a second-reading transcription test.
- The errata entries above are written, and POLICY.md has a row for the crisis
  payment and for what Step 18 offers after one.
- `task` is green and CI runs exactly `task`.

## Explicitly not in milestone 4

The three species aging tables (pp. 122–123) and the species names in their
headings. They are milestone 6, and the names come from the setting data file.
