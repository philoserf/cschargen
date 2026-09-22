# Changelog

## Unreleased

Bullets as the work lands. The release pass rewrites this into whatever the
release turns out to be.

### Breaking

- **Step 1's label changed, so every record written before it is refused on
  replay.** It read the book's word for a genetically engineered human, which
  the OGL notice reserves; it now reads "Step 1: Choose Human, Engineered
  Human, or Uplift". The refusal is the engine-version check doing its job and
  comes before any divergence — `--ignore-provenance` replays anyway and names
  the step whose label moved ([#108](https://github.com/philoserf/cschargen/issues/108), [#134](https://github.com/philoserf/cschargen/issues/134)).
- **A character holding both a master's and a doctorate no longer takes the
  graduate enlistment bonus twice.** p. 212 prints one modifier and the engine
  added one per degree, so a doctorate-holder enlisted at +8 where the page
  gives +4 — on Instructor, the difference between EDU 12+ and EDU 4+. Ten in
  three hundred auto characters hold both, and every one of them that tried
  Instructor, Journalist, Medic or Scientist got the wrong throw
  ([#132](https://github.com/philoserf/cschargen/issues/132)).

### Fixed

- **Three of Step 20's four answers did not survive a replay.** Only the name
  was written back to the record's inputs, and `Replay` cannot be asked, so an
  interactive character replayed with no gender, appearance or goals — and
  `replay` said "identical" while it happened ([#109](https://github.com/philoserf/cschargen/issues/109)).
- **`replay` compared the whole event log and one of the character's twenty
  fields.** It now compares the whole character ([#111](https://github.com/philoserf/cschargen/issues/111)).
- **A character who returned to higher education could not be replayed.** The
  choice was gated on the decider being askable rather than on the run being
  interactive, so it was offered while generating and skipped while replaying,
  and every choice after it read the wrong index ([#106](https://github.com/philoserf/cschargen/issues/106)).
- **Every record stamped `policy 0.3.1` and `POLICY.md` declared itself
  0.3.0.** The document is now 0.3.2 — the subsector rule changed behaviour
  between those two versions — and a test holds the two together
  ([#110](https://github.com/philoserf/cschargen/issues/110), [#124](https://github.com/philoserf/cschargen/issues/124)).
- **The two lists of Medic specialties disagreed**, so medical school granted a
  specialty the setting-data validator rejected ([#131](https://github.com/philoserf/cschargen/issues/131)).

### Changed

- **The Product Identity boundary is a gate rather than a rule people
  remember.** It reads every tracked file instead of only `setting/sample.json`,
  which is how thirty-eight occurrences across eighteen files went unnoticed —
  two of them test fixtures using real subsector names ([#134](https://github.com/philoserf/cschargen/issues/134)).
- **Record verification moved into the package that owns the record.**
  `chargen.Reproducible` and `chargen.Verify`; `cmd/cschargen/replay.go` went
  from 215 lines to 73, and eight statements no test could reach are now
  covered ([#113](https://github.com/philoserf/cschargen/issues/113)).
- **`docs/` is retired.** Eight files and 1,614 lines of process history that
  git and this changelog already held. `THEORY.md` and `WALKTHROUGH.md` take
  its place, on cadences `CLAUDE.md` records ([#120](https://github.com/philoserf/cschargen/issues/120),
  [#121](https://github.com/philoserf/cschargen/issues/121)).

## v0.1.0-alpha.4 — 2026-09-21

Thirteen issues out of one afternoon of playtesting the tool the way a referee
would: a transcribed setting file loaded with `--data`, a cast of NPCs, a
character built to a concept, and a list of the mistakes anybody makes on their
first day.

Three of them put a wrong character on the sheet and said nothing about it.

### The character the book would not recognise

**A baseline human could choose the slave career.** `eligibleCareers` excluded
Prisoner and not the career of pp. 150-154, which the book reaches only from
p. 42 and which `career/slave.go` says in its own words "cannot be chosen".
Neither takes an enlistment throw, so neither can fail one, and a career on
that list that cannot fail is a career whoever draws it enters. 14 characters
in 100 against the repository's own sample data
([#75](https://github.com/philoserf/cschargen/issues/75)).

**`(Any)` landed on skills that have none.** p. 116 gives "(Any)" its meaning —
it "allows the character to select a specialty within that skill" — and
`pickSkill` attached it to all thirty-six skill names it is called with.
Seventeen of those have no specialty to select. On a sheet that read
`Admin-1, Admin (Any)-3`: two skills where the book means one, on 19 characters
in 100 ([#77](https://github.com/philoserf/cschargen/issues/77)).

**A subsector throw into a gap became a choice, and the policy answers a choice
with its first option.** p. 39's chart is a 1d6 and a setting file may name
fewer than six subsectors. On the sample, which claims three of the six
results, that put 71 characters in 100 on one subsector and none at all on the
fourth. A throw that lands nowhere is thrown again now, which is what a table
does and what keeps the chart's own weighting
([#85](https://github.com/philoserf/cschargen/issues/85)).

### What the book permits and the engine did not

**A homeworld can be chosen.** p. 39 says "the Referee may choose the subsector
for the character" and p. 40 says a player may "simply choose a world that fits
the character concept you have in mind". The engine offered neither, and
`chooseHomeworld` had quoted that second sentence in its doc comment since
milestone 1, above a function that only rolled. There are three ways in now:
`--homeworld` and `--subsector` pin it, an interactive run is offered the
page's own choice with the throw first, and `--auto` rolls as before
([#76](https://github.com/philoserf/cschargen/issues/76)).

That reaches a whole page that nothing could. The Recently Colonized Worlds
table (p. 55) has no `originRoll`, so the 1d6 never selected it: across 200
generated characters, not one was born on any of its fourteen worlds. p. 40
introduces that table with "you cannot randomly be assigned one of these
worlds, **you may choose them**", and the choosing is what was missing.

**A path is offered by what opens it.** Steps 6 and 7 offered `Path 1, Path 3,
Path 4` — ordinals for the character's childhood and adolescence, with the gaps
reading as a bug rather than a requirement the character missed. The book gives
the paths no titles; it heads each with what opens it, so the heading is both
the description and the answer
([#80](https://github.com/philoserf/cschargen/issues/80)).

**`--tech-level` and `--max-terms` do what they say.** Both were read from the
command line, stamped into the record, and used by nothing — a record naming a
tech level that never gated an aging throw, against a `version` command that
exists so a bug report can be matched to the record the binary wrote. The
override is logged where it happens and cites its page
([#101](https://github.com/philoserf/cschargen/issues/101)).

### The referee's table

- **`--names` deals rather than draws.** Twelve names over twelve NPCs produced
  nine people, three of them twice, and left four names unused. The duplicates
  were holding up an identity alpha.3 gave away on purpose
  ([#78](https://github.com/philoserf/cschargen/issues/78)).
- **`render --roster` says where they are from.** The casting view left out the
  field that places an NPC and that the setting file exists to supply
  ([#79](https://github.com/philoserf/cschargen/issues/79)).
- **An unknown `--career` is a flag error.** It used to report that the
  character "failed to enlist" and to "try another seed", when nothing enlisted
  and no seed would help. `--career belter` entered the career on none of eight
  seeds where `Belter` entered it on three. Matching is case-insensitive now,
  and a near miss is suggested where exactly one career matches
  ([#81](https://github.com/philoserf/cschargen/issues/81)).
- **`--terms 0` means zero terms**, where it used to mean the policy's four
  ([#82](https://github.com/philoserf/cschargen/issues/82)).
- **`--tech-level` and `--max-terms` are held to the range a world's own
  numbers are held to**, against the same exported constants rather than a
  second copy of them ([#84](https://github.com/philoserf/cschargen/issues/84)).
- **A prompt reads as a question.** Four choice points take their prompt from
  the table result that raised them, and a result is written as the page writes
  it ([#83](https://github.com/philoserf/cschargen/issues/83)).

### The gate

**Four and a half minutes to ninety-three seconds**, and `chargen` was 99.6% of
it. Thirty-two `lifepath` call sites sweep the same seeds at seven distinct
term counts, so the suite walked 1,860 lifepaths to look at 420 characters —
ten separate sweeps each generating sixty eight-term characters. A character is
a pure function of its seed, term limit, setting and decider, and nothing in
the tests writes to one, so `lifepath` keeps what it makes
([#92](https://github.com/philoserf/cschargen/issues/92)).

CI also stopped linting twice. `golangci-lint-action@v8` dropped the
`install-only` input and ignored it with a warning, so the action linted the
tree and `task lint` linted it again — which meant CI was running a lint the
Taskfile did not control, with the action's own arguments, against a workflow
whose first line says CI runs exactly `task`. Nothing failed, so nothing said
anything. The actions are current and Dependabot watches them now, because that
half of the build cannot report its own drift the way the unpinned toolchain
can.

### Replay

**A record written before this release will not replay.** `(Any)` no longer
lands on a skill that has no specialties, so a seed produces a different
character than it did under alpha.3. The divergence names the value that moved:

```
replay diverged at event 80: consequence "Admin (Any) 1" against "Admin 1"
```

`schema_version` stays at 2. Its own definition is the shape of the record, and
the shape has not changed — `specialty` is still an optional string. A bump
would replace a message that names the skill with a blanket refusal that names
nothing, and `--ignore-provenance` does not help either: it skips the version
checks, not the comparison.

`policy_version` moves 0.3.0 → 0.3.1. POLICY.md gained rows for the two choice
points Step 3 and Step 4 now offer, and the rule is to bump when that document
changes. It moves in the last place because the decisions did not: only an
interactive run reaches either point, so the policy decides nothing new and the
same seed yields the same character.

## v0.1.0-alpha.3 — 2026-09-21

Ten issues out of one afternoon of using the tool the way a referee would.
[#57](https://github.com/philoserf/cschargen/issues/57) cast a campaign's worth
of NPCs from the alpha.2 binary and wrote down what happened; #58 to #67 are that
report split into work, and this is the work.

None of it was reachable from the gate. Every one of the ten is something the
engine did correctly and said badly, two commands the README prints next to each
other that did not compose, or a policy that was right for one character and
wrong for a hundred of them.

### Breaking

- **Policy version 0.2.0 → 0.3.0.** `batch` draws the career and the assignment
  from each character's own seed, where it used to take the first option the book
  prints at every choice point. `new --auto` is unchanged. The table below is what
  that did to a population.
- **A batch member is no longer `new --seed base+i`.** The two used the same
  decider and produced the same character; they now produce different characters
  from one seed, deliberately. Reproducing a member is `replay` on its record,
  which is exact — or `batch -o dir/`, which keeps each one as its own file.
- A record made under 0.2.0 still replays either way: `policy_version` is recorded
  and never verified, and replay reapplies recorded choices without asking the
  policy.

### A batch is a population — `policy_version` 0.3.0

Two policy changes, one version bump. `new --auto` is unchanged: it still takes
the first option the book prints at every choice point.

**`batch` draws the career and the assignment.** Measured on the same hundred
characters the section below counted under 0.2.0:

|                               | 0.2.0   | 0.3.0        |
| ----------------------------- | ------- | ------------ |
| careers represented           | 8 of 34 | **34 of 34** |
| started as Adventurer         | 55      | 1            |
| holding Gunner (in 300)       | 0       | **18**       |
| distinct assignments (in 300) | 10      | **88**       |

The assignment is in for the reason in that third row. With the career alone, a
hundred characters covered every career and **nobody in three hundred held
Gunner** — every National Navy character took the first assignment the career
prints, which is not Gunnery. With the career drawn and the assignment not — the
shape first planned — three hundred characters reached thirty-two assignments
between them, about one per career, and still nobody held Gunner. A crew can be
cast from the second column and could not be cast from the first.

Everything else still takes the first option, deliberately: varying the rest
would be a different policy rather than the same one applied to a crowd.

**`--names file`** draws a character's name from a list the caller supplies, one
name per line, by that character's own seed. The engine still invents nothing —
FR12 leaves Step 20's fields empty in auto mode, and this fills one of them only
when somebody hands it a list. `--name` beats `--names`, and `batch` refuses
`--name`.

### POLICY.md says what the rule does to a population

`POLICY.md` was candid about two consequences of "take the first option the book
prints" — every character takes Youth Path 1, and an eligible character attends
every school. It said nothing about the third, which is the one that matters for
the command the auto policy exists to serve.

Measured and written down: 100 auto characters draw from **8 of 34 careers**, 55
of them Adventurer, and nobody in the hundred has Engineer or Gunner at any
level. The `career` row is the only one of the twenty-eight a batch repeats a
hundred times, and every later choice point inherits its narrowness.

The behaviour is unchanged. This documents `policy_version` 0.2.0, which is the
version every record generated so far points at, and it goes in before the
behaviour changes rather than after — a version string that identifies a document
is only worth having if the document was true when the records were made.

### The sheet a referee reads

The relationships section is split. Across a sample pool **90% of every
character's relationships were family** — parents, siblings, grandparents,
aunts, uncles and cousins, each an Ally at 100 to 150 from Step 5 — so a line
reading `Allies: 76 (100 x62, 125 x13, 190)` buried the one colleague a plot
could hang on among sixty-two cousins.

```
## Relationships

**Made along the way**

- Ally (125) — youth
- Contact (50) — Adventurer
- Enemy (-150) — Belter

**Family**

- Allies: 15 (100 x11, 125 x4)
```

The people a character met are a list, each naming where it came from. The
family is a count, because a referee wants to know a character has a large
family and what it thinks of them, not to read sixty-two lines each saying
"cousin". Both headings are omitted when there is nothing under them.

The record has always known which was which — every tie carries an origin
(FR15). The sheet was throwing it away.

### A batch you can read

`batch` wrote JSONL and `render` read one record, so the two commands the README
prints next to each other did not compose: twenty NPCs were 777 KB of JSON and no
sheets.

- **`batch -o dir/` writes one file per record**, numbered to the width of the
  count so a listing and a glob sort correctly. This was in the PRD's own CLI
  sketch — `-o dir|file.jsonl` — and had never been built. Every other command
  works on the result unchanged: `render crew/npc-03.json` needed no new code.
- **`render` reads a batch.** Given JSONL it writes each sheet with a rule
  between them. A record `new` wrote is indented across many lines, so one record
  is tried first and only a file that is not one is read as a file of many.
- **`render --roster`** is one entry per character — age, careers, best four
  skills, money, and the ties that are not family. Casting a crew from a pool
  means skimming, and the sheet is the right document for one character and the
  wrong one for a hundred.

The roster counts ties without the family because 90% of a character's
relationships are relatives, and the one ally at 190 is the one a plot can hang
on.

### The command says what it did

- **`--career` says when it did not get the career.** A failed enlistment closes
  a career for two terms (p. 110) and the character drifts, which is the rules
  working — but the command was silent about it, so a referee scripting a cast
  got silent substitutions and found out twenty sheets later. Nine of twenty seeds
  asking for Corporate Shipper produced one. Now a note on stderr, and exit 0
  unchanged: the character generated successfully, and somebody running a hundred
  should not have to check a hundred exit codes. `batch` counts the substitutions
  rather than repeating the line.
- **`batch` refuses the four Step 20 flags.** `--name`, `--gender`, `--appearance`
  and `--goals` set one character's fields; applied to twenty they set the same
  four values on everybody. Refusing costs one line, applying them quietly costs
  twenty wrong sheets.
- **Each command's help names itself.** A flag set carries the name it was built
  with, and `new` and `batch` share theirs, so `batch --help` announced
  `Usage of new:`.

### The record says what happened

Two places where the record knew something and did not say it, both found by
using the tool as a referee ([#57](https://github.com/philoserf/cschargen/issues/57)).

- **Step 9 always names the career entered.** Three ways lead out of it — the
  career the decider chose, the one a table result ordered, and the Vagabond a
  character drifts into when nothing will have them — and two of the three
  recorded nothing, leaving a step header with no events beneath it. A reader
  could not tell why the character changed course, and the empty step read like a
  rendering fault. p. 111 is clear that the Vagabond is the rules working: "it
  ensures that characters always have a path forward, even when opportunity is
  scarce." The silence was the bug, not the destination.
- **Every relationship names its origin**, which FR15 asks for and 138 of 2,477
  ties in a sample pool did not do. A tie made before Step 9 has no career to
  name, so it names the stage of life instead — `youth`, `teenage`, or the school
  attended. A parent and a childhood mentor are both Allies at 125, and the
  origin is the only thing that tells them apart.

Two gates hold both: every tie in a generated record carries an origin, and Step
9 never passes without naming a career. The second is deliberately narrow rather
than "no step is ever empty" — Step 18 is legitimately empty in most records,
because deciding to carry on is not an event.

Writing the first gate turned up a third silent path nobody had noticed: a career
a mishap ordered arrived with no line at Step 9 either, because the result that
sent the character there had spoken one step earlier. Writing the second showed
that two of the four reasons a career list can empty cannot actually happen — an
aging crisis and a mental characteristic at 0 both answer with Vagabond, which is
a list of one rather than a list of none.

## v0.1.0-alpha.2 — 2026-09-20

One pass, in twelve parts: [#42](https://github.com/philoserf/cschargen/issues/42),
the results the engine recorded rather than resolved. **387 occurrences of "the
engine could not carry this out" became 60**, and 190 distinct results became 51.

Five of them turned out to be bugs rather than gaps, and three new gate tests
exist because of what they showed — a corpus this size hides a silent modifier
well.

### Breaking

- **Record schema 1 → 2.** `stash` is a list of possessions rather than a list of
  names, because a possession may carry the value the book rolled for it. `replay`
  refuses a version-1 record, which is what the version is for: **a character
  generated under alpha.1 will not replay.** Regenerate from the seed instead —
  the record carries it, and the sheet prints it.
- **Policy version 0.1.0 → 0.2.0.** `POLICY.md` now carries a row for every one of
  the twenty-eight choice points the engine offers, where it had nine. Two of them
  are not "the first option printed" and say why. A test holds the document to the
  source both ways, because the document's own rule — "a row here and no code is as
  wrong as code and no row" — went twelve pull requests without being honoured,
  which is what a rule with nothing checking it does.

  A record made under 0.1.0 still replays: replay reapplies recorded choices and
  never asks the policy.

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
