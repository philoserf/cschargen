# Changelog

## Unreleased

### Resolved rather than recorded

The first of the passes issue #42 lays out. Seventy-three of the 387 results the
engine recorded rather than resolved now resolve, all of them about mustering-out
rolls:

- **Every benefit roll forfeited** (67 occurrences). "You are dismissed and lose
  all Benefits" now takes the queued batches _and_ writes off the two-per-term
  grant up to the term the result fired in. It is a watermark rather than a
  switch, because a few of these are events rather than mishaps and the character
  serves on afterwards: what was earned before the result is gone, and the terms
  after it earn as they always did.
- **Cash rolls** (12). A result that grants rolls on the Cash column grants them
  there; one that says "immediately" reads the career's Cash table where it fires
  instead of queueing for Step 19. ERRATA E-32 records what a compelled Cash roll
  does when p. 127's three-per-career cap is already spent.
- **Rolled counts** (2). "Lose 1d6 Benefit rolls" and "gain 1d3" are thrown.
- **Expulsion from a graduate track** (2) now closes every institution for two
  terms (pp. 98, 102), not just the one that expelled the character.

### Named benefits reach the record

A second pass on #42, and the finding the plan leads with. Fifty-four more
occurrences resolve — 176 distinct results become 159, 314 occurrences become 260.

`benefitRoll` passed a weapon, a company share, a pension or a rare item to
`apply`, which recorded it as something the engine could not do. They are
possessions, and the record holds them now: a Corporate Shipper who mustered out
three times shows three company shares rather than three apologies.

- **Possessions carry the value p. 128 gives them** (39 occurrences). A company
  share is worth 2d6 × 100,000 wherever it was granted, including on the rows
  that say "Three Company Shares" and print no number — the definition is of the
  benefit, not of one row, which is ERRATA E-33. Pieces of art are valued one
  throw each, because p. 128 says the value "should be rolled at the time the
  benefit is received". A weapon and a suit of armour carry none, because the
  book gives none: FR11's "names those and stores no invented value".
- **The Weapon benefit is a choice, not an item** (15 rows). p. 129 ends it with
  "If the player wishes, they may choose to take a level in Melee (Any) or Gun
  Combat (Any) in lieu of a weapon", and two of those three branches are things
  this engine can carry out.
- **"Lose any company shares"** (14) has something to remove. It takes the
  shares and leaves everything else, and says so even when there were none.
- **Money is money.** A Pension, a Prize Share and a Church Pension are amounts,
  not objects, so they pay credits. The church lump sum is 10% less than the full
  value (p. 128), which is 2d6 × 9,000.
- **A benefit row can do two things.** `EffectGroup` exists because a row
  carries one effect and Celebrity's Producer credit is an object, a payment and
  an annuity.

### Relationships lost, changed and improved

The third pass on #42. Forty more occurrences resolve — 159 distinct results
become 135, 260 occurrences become 220.

- **Losses take a number.** "Lose 1d3 Allies and Contacts", "lose two Contacts",
  "lose every Contact and Ally gained in this career": each removal works down
  the printed order afresh, and "in this career" is measured against the origin a
  tie already carried.
- **A tie can change what it is.** "One Contact or Ally from this career becomes
  an Enemy", "1d3 existing Contacts become Allies", "a Rival becomes an Enemy, or
  a Rival is gained where there was none". A changed tie takes the middle of its
  new band, for the reason E-23 gives.
- **`improveARelationship` is implemented** (6 occurrences), which is the result
  the issue named first. **ERRATA E-34** records why the four clauses
  pp. 85, 91 and 96 print are read against the state as it was: in sequence, one
  Enemy would climb to Ally on a single result.
- **Results that find nobody do what they say instead.** "Lose one Ally or
  Contact. If you have none, gain an Enemy with a Relationship Rating of -110."

### The skill list, and the results that need it

The fourth pass on #42. Forty more occurrences resolve — 135 distinct results
become 127, 220 occurrences become 180.

- **The skill list of pp. 304–314 is transcribed** (44 skills with their
  specialties), which is what makes "gain a level in any skill of your choice"
  and "any skill at level 1 which you do not already have" answerable. It also
  fixes the free skill an undergraduate degree grants, which had been recorded
  rather than offered.
- **"Raise a skill the character already holds"** (31 occurrences) offers what
  the character has, and raises that entry — specialty included, so a character
  who holds Melee (Blade) does not end up with a bare Melee beside it.
- **The whole corpus is held against the list.** Two new tests walk every
  transcribed effect and check the skill and specialty names. They found six
  spellings the book itself writes two ways, which **ERRATA E-35** records: a
  Vagabond event grants "Gambling" where the list calls the skill Gambler, and
  "Melee (Unarmed)", "Gunner (Turret)" and "Electronics (Communications)" were
  each a second skill on the sheet.
- **"A Survival specialty used on your homeworld"** (5) reads the world's own
  background skills, which is the only place the setting data says what a world
  is like in those terms.
- **"Lose all levels of Suit (Vacc Suit)"** (1) does.
- **`data validate` holds a setting file's background skills to the same list**,
  because that is the other way a second skill reaches a sheet. It caught
  `Gun Combat (Slug)` in this repository's own sample data the first time it ran.

### Rank lost, and rank gained

The fifth pass on #42. Fourteen more occurrences resolve — 127 distinct results
become 121, 180 occurrences become 166.

- **A demotion moves the number and nothing else.** Two of the sixteen results
  add "retaining any benefit already gained" and fourteen do not; **ERRATA E-36**
  reads them as the same thing, because the book has no rule anywhere for taking
  a skill level back. A character at rank 0 loses nothing, and the record says so.
- **An event promotion grants that rank's benefits.** p. 116 attaches them to
  holding the rank, not to the advancement throw that usually reaches it, and a
  bare `EffectRank` had been skipping them.

### Careers have classes, and enlistment reads them

The sixth pass on #42. Twenty more occurrences resolve — 121 distinct results
become 102, 166 occurrences become 146.

Twelve results modify or direct an enlistment by the class of career being
entered rather than by name — "-4 DM to enlist in any government related career",
"-6 DM to any enlistment roll for a career involving violence" — and the book
defines none of the words it uses. **ERRATA E-37** carries the assignment, one
row per career, argued from the career descriptions of pp. 132–145 rather than
from the names. Fifteen careers carry no class at all, and two of those are
deliberate: an Investigator works in "legal systems, corporate structures, and
private practice" without the record saying which, and an Adventurer's career is
danger rather than violence.

- **An enlistment modifier knows which careers it reaches.** A penalty aimed at
  government careers is not spent by enlisting in a criminal one.
- **A standing modifier outlives the throw.** "Every career after this one" and
  "your next career" are different sentences, and several results print each.
- **"You may enlist automatically"** succeeds without a throw, narrowed the same
  way — Undergraduate University's offer reaches four classes and no others.
- **"Return to the career you held before this one"** resolves against the
  service record at the moment the transfer is taken.
- Three results that send a released prisoner to Vagabond are transfers, which
  the engine already did.

### Modifiers that reach the throw they name

Found while scoping the next pass, and a bug rather than a gap: the count of
unresolved results does not move.

**Twenty-five results grant a modifier to "next survival roll" and nothing in the
engine consumed it.** A survival throw was made without them. The record showed
the modifier granted, the transcript showed it in the log, and the throw was
simply easier than the page said. About twenty more were filed under names no
part of the engine reads — "next two advancement rolls", "enlistment in any
criminal career", "Melee checks in this career", "advancement rolls for the rest
of this career".

`TestEveryModifierNamesAThrowTheEngineTakes` holds every modifier in the corpus
against the five throws the engine takes. It is what found these, and what stops
another one being written.

Closing them needed three things beyond PR 6's narrowing:

- **`Uses`** — "your next two Advancement rolls" is one modifier spent twice.
- **`WhileInThisCareer`** — "for the rest of this career" ends when the career
  does, because nothing else would ever end it.
- **`OnSkill`** — "+1 DM to Melee checks in this career" is not spent by a check,
  and no other skill's check takes it.

Survival and advancement throws now read the class of the career they are made
in, so "-2 DM to the first two Advancement rolls in a military career" reaches
only those.

### Pools, hooks and decided throws

Nineteen more occurrences resolve — 102 distinct results become 97, 146
occurrences become 127.

- **A pool the character spends themselves** (8). "Gain a modifier of +6 which
  must be split up into increments of not more than +2 to use on any Survival or
  Advancement rolls until it is depleted" is a question asked at each of those
  throws while the pool has credit. The options run from the largest increment
  down, so a policy character spends the pool rather than carrying it to the end
  of a career — a decider's habit rather than a rule, because the book says "at
  any time" and a player may say nothing.
- **Effects that wait on a failed throw** (9). "If you fail that Advancement
  roll, you lose the Ally and take a -2 DM to your next Advancement roll" is the
  largest single line the corpus had left. The hook is taken whether the throw
  passed or failed, so it cannot fire against the next one.
- **A throw decided against the character** (2). "You are suspended, and your
  next Advancement roll fails automatically": the throw is not made, the record
  says why, and what waited on its failure still fires.

### A 1d6 table printed inside a result

Twenty-nine more occurrences resolve — 97 distinct results become 76, 127
occurrences become 98.

(The count itself is more honest than it was: the walk that produces it did not
descend into an effect's grouped, fallback or sub-table branches, which the last
few passes gave it. Earlier figures in this changelog were therefore a little
optimistic.)

Forty-two results print a small table of their own: "Roll 1d6. On a 1 … on a
2-5 … on a 6 …". `EffectSubTable` is a **roll**, not a choice — the character
has no say in which row comes up, so the record carries the die as it fell and
a replay of the same seed finds the same row.

`TestEverySubTableCoversTheDie` holds every one of them to covering 1 through 6
exactly once. A gap would leave the engine with nothing to apply and an overlap
would let the first row printed win silently.

Twenty-five of the forty-two are transcribed against mechanisms the earlier
passes built — the Vagabond's parting pay, the Diplomatic Service's conspiracy,
the Medic's lawsuit, the Independent Merchant's distress call, the slave's
escape attempt. The rest wait on the two things they need: a list of the
conditions a character carries, for the thirteen addiction results, and a
player's own words for the religion ones.

### What a character carries

Twenty-three more occurrences resolve — 76 distinct results become 64, 98
occurrences become 75.

`State.Conditions` is the third list beside skills and possessions: what a
character carries that is neither a number nor a thing. An addiction goes there,
and a religion.

- **Thirteen addiction results.** **ERRATA E-38** records why nothing else
  happens: most of them read like "you have also gained an addiction to alcohol
  or a drug which will haunt you for years to come", and the book prints no
  characteristic loss, no modifier and no rule beside it. A cost invented here
  would be a house rule wearing a page number.
- **Nine religion results.** "Choose or invent a religion" is the player's answer
  on the sheet; the engine records that there is one, and carries out the
  mechanical half the results do print — a level in Science (Philosophy) on a
  1d6 of 6.
- **Four more 1d6 tables** the two unblocked, including the drug trials the
  Vagabond and Prisoner careers each run on the people they can spare.
- A character is addicted once. Taking the result again records that they
  already were.

### The last of the 1d6 tables that could be read

Five more occurrences resolve — 64 distinct results become 59, 75 occurrences
become 70 — and six of the forty-two sub-tables are now transcribed from the page
rather than summarised: the Celebrity's game show and interview show and
endorsement deal, the Scientist's discussion panel and breakthrough, the
Independent Merchant's reality holovid.

"The Science skill you used on this task" raises one of the character's
Sciences rather than offering a new specialty, which could have given Science
(Chemistry) 1 to a character whose check was made on Science (Physics) 3.

The Arts career's own interview show is transcribed as far as it can be. Its four
rows each call for a Diplomat check at a named **difficulty**, which is the Core
Rulebook's task system and outside this ruleset (**ERRATA E-12**) — so the die is
rolled, the row is recorded, and only the check is left open. That is four small
unresolved results where there was one large one, and a reader can see which part
is missing.

### The close of #42

Ten more occurrences resolve — 59 distinct results become 51, 70 occurrences
become 60, from the 190 distinct results across 387 occurrences the pass began
with.

- **A sub-table on 2d6.** The enslaved tables' escape result rolls 2d6 and adds
  the character's DEX modifier. **ERRATA E-40** records the gap it leaves: the
  bands are "less than 5", "6-10" and "11+", and a total of exactly 5 falls in
  none of them.
- **Age, sentences and debts.** A year lost outside a term, a term added to or
  taken off a sentence, and money owed rather than held — a debt goes beside an
  addiction, for the same reason: the book names it and prices nothing else.
- **Three careers named as one choice.** "Enter the System Defense Force career —
  Naval, Troopers or Wet Navy" is a pick of three transfers.

### What is left, and why

Three ERRATA entries close the pass by naming the causes rather than the results:

- **E-41** — twelve results need a map of the sector. The setting data is a list
  of origin charts; it carries no coordinates and no starport class, so "the
  nearest A-class port" would be a world picked at random wearing a rule's name.
- **E-42** — five results send a character back through the steps. The engine
  walks the twenty once, and a step loop any result could re-enter is a different
  generator from the one the PRD describes.
- **E-43** — eight outcomes are the referee's. A bounty hunter is a campaign, not
  a modifier; "your entire group" is however many the referee says; and two are
  conditions that last while something else remains true, which is a shape the
  engine has no room for.

The rest are the small bespoke results that were always going to be last, and
E-11, E-12 and E-13 cover their own.

### Fixed

- **"You lose the Ally" could take a parent.** The nine results that read "gain an
  Ally and a +4 DM to your next Advancement roll; if you fail that roll, you lose
  the Ally" mean the superior officer the result granted a sentence earlier.
  **ERRATA E-39** reads it as the most recently gained Ally that is not family: a
  promotion that goes wrong at work does not cost a character their mother.

- A character dismissed from a career could still muster out of it.
- An instant promotion from an event granted the rank and none of its benefits.
- A character could hold both Gambler and Gambling, or both Melee (Unarmed) and
  Melee (Unarmed Combat), as separate skills.

### Changed

- **Record schema 1 → 2.** `stash` is a list of possessions rather than a list
  of names. `replay` refuses a version-1 record, which is what it is for.

## v0.1.0-alpha.1 — 2026-09-20

The first release. Every step of the book's character generation is built, and
nothing in it has been played.

### What it does

Twenty steps, in the order the book prints them:

| Steps | What                                                                  |
| ----- | --------------------------------------------------------------------- |
| 1–2   | Species, and characteristics by that species' own method              |
| 3–4   | Subsector, homeworld, background skills, primary language             |
| 5–8   | Family, youth events, teenage events, four tracks of higher education |
| 9–19  | All thirty-four careers, aging, mustering out                         |
| 20    | Name, gender identity, appearance, personal goals                     |

```sh
cschargen new --auto --seed 7 --terms 4        # or without --auto, and it asks
cschargen batch --auto --count 20 --seed 100   # twenty, as JSONL
cschargen render character.json                # the sheet
cschargen render --history character.json      # every throw, with its page
cschargen replay character.json                # re-run it from the seed
cschargen data validate setting.json           # before it fails mid-lifepath
```

Every throw carries the page it came from. Every choice records who made it —
player, policy or replay — so a record made either way replays the same.

### What it does not do

**It ships no Product Identity.** World, subsector and species names live in a
data file you write from your own copy of the book. An invented sample ships so
the tool runs without one, and characters built on it are stamped as such and say
so on their sheet. `README.md` explains what the OGL notice reserves and why the
project says "engineered human".

**It resolves 151 fewer results than the book prints.** Those are recorded rather
than skipped, so a record always says what the engine could not carry out —
"a weapon of the character's choice", "roll Melee at Difficult", "a minimum of
four terms inside the sector". None of them is silently approximated.
[#42](https://github.com/philoserf/cschargen/issues/42) sorts them.

**It does not know the published setting**, compute anything about play, or
value the equipment a character musters out with. The record names a Rare Item
and stores no invented price.

### Thirty errata

`ERRATA.md` records thirty places the book is mistaken, silent, or asks for
something the engine cannot evaluate. Each says which kind it is, what the engine
does, and — for a reading — what changes if the reading is wrong.

Three are worth naming here because they change every character:

- **Rank benefits are a floor, not an increment** (p. 116). A character who
  already holds a skill at the rank's level gains nothing from it.
- **A relationship never crosses zero** (p. 320). An Ally who falls far enough is
  lost rather than becoming a Rival.
- **Aging is indexed by term number and gated by the homeworld's tech level**
  (pp. 122–123), not by age. A character from a high-technology world may serve
  forty terms before their first aging throw.

### Why alpha

Nothing has been played. The engine follows the book and cites it, but no
character generated by it has sat at a table, and the readings are one person's
reading of a book published this year. A reading that turns out wrong changes
what the engine generates, and records carry the readings they were made under —
so an old record stays auditable after a reading changes.

The schema is versioned, the policy is versioned, and a replay against a
different setting file is refused rather than allowed to diverge quietly.

### Provenance

Ruleset baseline: **_Clement Sector Core Character Creation Book_**, © 2026
Independence Games, author John Watts. This project is not affiliated with
Independence Games, Samardan Press, Mongoose Publishing or Far Future
Enterprises.
