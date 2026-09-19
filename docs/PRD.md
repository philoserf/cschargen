# PRD: Clement Sector Character Generator (Go CLI)

2026-09-19. Status: draft for review; nothing built yet.

## Problem

Clement Sector character generation is a twenty-step lifepath the book itself calls
a minigame that "often requires a game session of its own" (p. 16). It is
substantially deeper than Classic Traveller or baseline Cepheus: eight pre-career
steps before a character enters their first career, thirty-four careers with
assignment tracks, a d66 event table and a 2d6 mishap table per career, and an
aging system keyed to homeworld tech level that lets a character serve thirty-five
terms. Doing it by hand is slow and error-prone. Build a Go CLI that generates
rules-accurate Clement Sector characters, as a sibling to `philoserf/ctchargen`
and `philoserf/t5chargen`.

**Ruleset baseline:** _Clement Sector Core Character Creation Book_, © 2026
Independence Games, author John Watts, 337 pp., as held at
`~/Documents/Traveller/Clement/`. All page cites are to that artifact's printed
page numbers (printed page N is PDF page N+1). It supersedes both the 2021
single-volume _Clement Sector_ and the 2015–2016 two-book arrangement of
_Clement Sector: The Rules_ plus _Core Setting Book 2.0_.

The book is Cepheus Engine-compatible, not Traveller: the OGL chain runs Mongoose
Traveller SRD (2008) → Cepheus Engine SRD (Samardan Press, 2016) → _Clement
Sector: The Rules_ (2016) → this book. The engine shares no code and no rules
model with `ctchargen` or `t5chargen`; only the house conventions carry over.

## User

Mark: solo referee and developer. Secondary: any Clement Sector or Cepheus referee
needing NPCs in bulk, or a player who wants the lifepath walked rather than
looked up.

## Product Identity constraint

The book declares game mechanics Open Game Content, and declares as Product
Identity — quoting the book's own OGL section 16, p. 335 —

> all subsector names, world maps, world names, system names, system maps,
> vehicle names, starship names, starship classes, artistic depictions of ships
> and vehicles, and organizations. Please note that the term "altrant" is not
> open content.

That lands squarely on Steps 3–4, whose tables are tables of world names, and on
the engineered-species names in Step 1. Those tables are also load-bearing: a
homeworld sets background skills, primary language, maximum age, maximum terms,
and whether engineered characters may be born there at all.

**Resolution: the repo ships no Product Identity.** World, subsector and species
names live in an external data file the user supplies from their own copy of the
book; the repo carries a small synthetic sample so that tests, fixtures and the
README demo run without it. The engine, the mechanics and the career tables —
which are mechanics, not names — are in the repo. This is a data/PI boundary, not
a data/logic boundary; the career tables are embedded precisely because they are
open content.

Consequence for the CLI: a real character needs `--data <file>`. A character
generated against the sample data is valid and replayable but set on worlds that
do not exist.

## Goals

1. Generate a complete character per the twenty-step process (pp. 16–130), from
   species through finishing touches.
2. Two modes: **interactive** (player makes each choice) and **auto** (tool
   decides; supports batch NPC generation).
3. Deterministic replay: re-running from the recorded seed, inputs and choices
   reproduces the identical character.
4. Output a character record as JSON (canonical) and a Markdown character sheet
   modeled on the book's Long Form Character Sheet (pp. 324–333).
5. Emit a generation record: every throw, choice and outcome, in order, embedded
   in the JSON and renderable as a Markdown transcript.

## Non-goals (v1)

- Careers drawn from supplements. The book prints its own thirty-four careers and
  then, at pp. 107–110, catalogs dozens more that live in _Manhunters_, _Covert_,
  _Badge_, _Port of Entry_, _Outlaw_, _Cybersneaks_, _Interface_, _Unmerciful
  Frontier_, _Hub Federation Navy_, _Hub Federation Ground Forces_ and _The Dark
  Lonesome_ — most of which are in the collection. v1 implements only this book's
  thirty-four. The career data format should make supplement careers additive, and
  that catalog is the argument for getting the format right the first time.
- The _Diverse Roles_ career catalog.
- The quick-generation method in _Clement Sector Core Rulebook_ p. 153.
- Earth Sector origins beyond this book's single Earth Subsector Origin Chart
  (p. 43); the book refers those to the _Earth Sector_ expansion.
- In-play advancement: Adventure Points, Success Points, Instruction (pp. 301–303).
  The skill list (pp. 304–315) is in scope as a vocabulary; the advancement rules
  that spend points are not.
- NPC generation (pp. 315–323), robots, psionics, equipment purchase beyond
  mustering-out benefits, combat.

## Cross-book references

The book points into the _Clement Sector Core Setting Book_ at several places:
the task system (pp. 20–27), Hub Federation Credits (p. 184), mindcomp (p. 94),
armor and medical kit listings (pp. 90, 100), Allies (p. 188), the Seidel
Treatment (p. 26). **The engine needs none of them.** Each is either a play
reference or a description of a benefit the record names rather than resolves.
The record stores `mindcomp: true`, `armor: "player choice"` and a credit total;
what those mean at the table is the referee's business. This is stated so that no
one later reads the cross-references as a missing dependency.

## Functional requirements

Page cites are to the Core Character Creation Book.

**FR1 — Characteristics.** Six characteristics: STR, DEX, END, INT, EDU, CHA
(p. 13). Note CHA where Traveller has SOC. Generate by rolling 3d6, dropping the
lowest die, summing the remaining two; roll six times and **assign the results to
characteristics as the player sees fit** (p. 13). That assignment is a choice
point the Decider owns — unlike the fixed-order roll in `t5chargen`. The modifier
table runs from 0 to 20 in bands (p. 14): 0 → −3, 1–2 → −2, 3–5 → −1, 6–8 → +0,
9–11 → +1, 12–14 → +2, 15–17 → +3, 18–20 → +4. Unaltered humans cap at 15; some
engineered species cap higher (p. 14). A modifier is recomputed the instant its
characteristic changes (p. 15).

Characteristic zero is a real state, not an error: all three physical at 0 is
death during generation (p. 14), and the aging rules in FR9 attach further
consequences.

**FR2 — Species (Step 1, pp. 21–38).** Human, altrant, or uplift. The book prints
five engineered-human types and five uplift types, each with its own
characteristic-generation method, characteristic caps, and age restrictions.
Species definitions are **external data** (FR14), because the names are Product
Identity; the engine handles a species as a record of dice methods, caps, age
rules and aging-table key. Human is the built-in default and needs no data file.

**FR3 — Origin (Steps 3–4, pp. 39–56).** Subsector of origin, then homeworld. A
homeworld supplies: background skills at level 1, with `and`/`or` composition and
specialty choices (p. 40); primary language, possibly a ranked list, where a
homeworld granting the Language skill must take a specialty other than the primary
(p. 41); maximum age and maximum terms (pp. 41–42); and whether altrants and
uplifts may live openly there, plus their legal status (p. 42). Every character
additionally gains Electronics 0 (p. 41).

If a homeworld enslaves engineered people and the character is one, the character
**must** take the Altrant/Uplift Slave career as their first career (p. 42).
Maximum age and maximum terms do not apply to altrants and uplifts, who carry
their own restrictions from FR2 (p. 42).

Homeworld is **not fixed after Step 4.** Six of the Colonist career's eleven
mishaps reassign it (p. 174), and other careers do the same. The record therefore
holds a homeworld history, not a homeworld.

All origin tables are external data.

**FR4 — Family (Step 5, pp. 57–66).** Parents with a birth-situation chart, parent
ages, siblings with ages, grandparents, uncles, aunts and cousins, each carrying
quirk and relationship rolls. Altrants instead roll genetic status — purebred
first generation, purebred later generation, compound, hybrid — on the random
genetics tables (p. 62). Uplifts roll class and origin (pp. 65–67).

**FR5 — Youth and teenage events (Steps 6–7, pp. 67–85).** Youth events on five
paths, gated by characteristic (no requirement; STR 8+ or END 8+; DEX 8+; INT 8+;
CHA 8+), plus a separate uplift youth table and an enslaved path. Teenage events
on four paths, two gated by how long the homeworld has been settled and two by
characteristic, again with uplift and enslaved variants. Each stage has its own
Life Events table.

The settlement-age gate on the teenage paths is a property of the homeworld data,
so this requirement depends on FR14 carrying a founding date or settlement age.

**FR6 — Higher education (Step 8, pp. 85–103).** Four mutually exclusive tracks —
undergraduate university, military academy, graduate school, medical school — each
with an entry throw, a failure branch, an events table and a life-events table.
Graduate and medical school have their own prerequisites. A character may return
here later from Step 18 (p. 126).

**FR7 — Career entry (Steps 9–11, pp. 104–112).** Thirty-four careers, a count
the book states itself (p. 107). Enlistment
is a 2d6 characteristic check against a per-career target, plus per-career
situational modifiers and prerequisites that gate the attempt entirely (pp.
110–111). On failure the character must try a different career and may not retry
the failed one for two full terms. **Three consecutive failed enlistments, across
any careers, force the Vagabond career** (p. 110).

Two careers have no enlistment roll (p. 111): Prisoner, which cannot be chosen and
is entered only when an event or mishap sends the character there, and Vagabond,
which may be entered voluntarily at any time and is where the rules dump a
character with no other path.

Assignment is chosen on entry (p. 112). A character must hold an assignment for at
least two terms, after which they may switch within the career without a roll and
without losing rank (pp. 112, 126).

**FR8 — The term loop (Steps 12–16, pp. 112–121).** Per term: survival throw
against the assignment's target; on failure, a 2d6 mishap (p. 113) that ejects the
character from the career unless the mishap says otherwise; advancement throw,
with commission where the career has one, rank and rank benefits (pp. 114–116);
one skill roll on a player-chosen table (p. 117); and an event roll (p. 118).

Skill tables per career are Personal Development, Service Skills, Advanced
Education (gated on EDU, usually 8+), and one table per assignment; military
careers add Officer Skills. **On the first term in a career the character gains
every skill on the Service Skills table at level 0** that they do not already hold
(p. 117). Thereafter it is one table, one d6, one result per term.

Events are d66 — thirty-six entries, of which a block (31–36 in Colonist) defers to
the shared Life Events table (p. 120), and military careers route to the shared
Military Events table (p. 121). Events reach further into the character than
skills do: they grant and remove benefit rolls, attach +1 modifiers to future
benefit rolls, add allies, contacts, rivals and enemies, force injury throws
(p. 119), change assignment, and send characters to other careers.

**FR9 — Aging (Step 17, pp. 121–125).** Add four years per term; a character
ejected by a mishap adds 1d3 years instead (p. 121).

Aging throws are indexed by **term number**, and which table applies is a function
of species and the homeworld's tech level (pp. 122–123). Humans from a TL 9-or-less
world begin aging throws at term 6; from TL 10, at term 18; from TL 11, at term 33;
from TL 12–13, at term 44. Each failed check costs one point of the checked
characteristic.

A characteristic reduced to 0 by aging is an **Aging Crisis** (p. 123): the
character dies unless 1d6 × 1000 credits buys emergency treatment, which restores
it to 1. A survivor automatically fails all future enlistment checks and may only
continue in their current career, stop, or become a Vagabond. Two physical
characteristics at 0 ends generation with the character playable but immobile;
three is death. One mental characteristic at 0 bars further enlistment; two is
death and a restart (pp. 123–124).

**Apparent age is derived, never rolled** (pp. 124–125): a lookup of actual age
against homeworld tech level, where TL 9 and below means apparent age equals
actual age. It is the setting's signature number and belongs on the sheet.

**FR10 — Next term (Step 18, pp. 125–126).** Continue, change assignment, change
career, return to higher education, or exit. Two rules constrain the choice: a
character at their homeworld's maximum terms **must** stop, and a character who
rolled a natural 12 on survival **must** continue in the career (p. 125). A
character who arrived here from a mishap may not continue in that career (p. 126).
Changing careers routes through muster out first (p. 126).

**FR11 — Muster out (Step 19, pp. 126–129).** Two benefit rolls per term served in
the career being left, adjusted by any event or mishap that granted or removed
rolls. Per roll the player picks the Cash table or the Other Benefits table, then
rolls 1d6; **the Cash table may be used at most three times per career** (p. 127).

Career benefit tables print rows 1 through 7 (e.g. Colonist p. 173, National Navy
p. 233) although the roll is 1d6. Row 7 is reachable only through the +1
benefit-roll modifiers that events award, and those come in two scopes: some apply
to every benefit roll made in that career, others only to the specific rolls the
event granted ("gain three Benefit rolls at +1"). The engine carries benefit rolls
as a queue of scoped batches, not as a count plus a career modifier.

Benefits are a closed vocabulary of named kinds (pp. 127–129): characteristic
increases, Ally, Contact, Rival, Armor, Weapon, Land Plot, Medical Kit, Gambler
Kit, Holy Book, Rare Item, Pieces of Art, Stash, Captain's Guild membership,
Company Share, Prize Share, Pension, Church Pension, Producer. Several carry rolled
credit values; several (Stash, Holy Book, Armor) are explicitly referee-valued or
player-chosen. **The record names those and stores no invented value.**

**FR12 — Finishing touches (Step 20, pp. 129–130).** Name, gender, appearance,
personal goals. All player input or policy default; none is rolled. The engine
records them and leaves them empty rather than inventing them in auto mode, except
for name, which auto mode may take from the data file if it supplies name lists.

**FR13 — Dice engine.** 2d6, 1d6, 3d6-drop-lowest, 1d3, d66, and target-number
throws as a distinct, unit-tested package. Every roll drawn from a seeded stream
and logged. d66 is its own throw kind, not two d6 rolls, so the log reads as the
book does.

**FR14 — Data files.** Two schemas, both JSON, both versioned:

- **Setting data** — the external, user-supplied file, with `worlds` and `species`
  sections. The printed origin charts (pp. 43–56) give the world columns directly:
  D100 range, planet of origin, background skills, primary language, maximum
  age/terms, altrant/uplift status, tech level. Two need more structure than the
  page shows. Background skills are an expression, not a list — `and`/`or`
  composition with specialty choices (p. 40). Altrant/uplift status is not a
  boolean: entries read like "Achilles banned. Others are allowed and free", so it
  is a per-species permission plus a free/enslaved status. Settlement age is
  additionally required by FR5's teenage-path gate and is not a printed column, so
  the schema carries it as its own field. The repo ships `data/setting.sample.json`
  with invented worlds and `data/setting.schema.json`.
- **Careers** — in the repo, as open content. Per career: enlistment
  characteristic and target, modifiers, prerequisites, assignments, per-assignment
  survival and advancement targets, the four-or-more skill tables, rank tables per
  assignment, the 2d6 mishap table, the d66 event table, and the muster-out cash
  and benefit tables.

Species definitions live in the **same** external file, in a `species` section
beside `worlds`, for the PI reason in FR2 — one file, one content hash, one
validator, one flag. The `origin` and `species` packages read different sections
of it.

**FR15 — Character record and generation log.** The JSON record is the source of
truth; the Markdown sheet renders it. The record carries the full career history
term by term, the homeworld history, skills with specialties, relationships
(allies, contacts, rivals, enemies) with their origin, benefits, credits, actual
and apparent age, and the ordered event log.

The log records each step entered, each throw (kind, dice, target, modifiers,
result, page cite), each choice (decider, what was chosen, what the alternatives
were), and each consequence, with consequences referencing the sequence number of
whatever caused them. It serves audit (walk it against the book), replay
(verification data), and narrative (render it as a lifepath transcript).

## Replay and provenance contract

Every record carries `schema_version`, `ruleset` (pinned: Core Character Creation
Book, © 2026), `engine_version`, `policy_version`, `rng` (algorithm and seed), an
`inputs` block, and a `data_sources` block naming the setting data file by identifier and content
hash.

`data_sources` is the departure from `t5chargen`, and it is forced by the PI
decision: the engine's tables are no longer entirely inside the binary, so a seed
plus a policy version no longer determines a character. A record made against one
setting file and replayed against another would diverge at Step 3 and never
recover. Replay therefore verifies the hash and refuses on mismatch, with
`--ignore-provenance` to force it and report where the generation disagrees.

RNG: Go `math/rand/v2` PCG, named in the record. Replay re-runs the engine from the
recorded seed, inputs and choice events, recomputing every throw; the stored log is
verification data, not input. It exits non-zero at the first divergence, naming the
sequence number.

## Recorded readings

The house convention is that every place the text is silent or ambiguous gets a
recorded reading with its page cite. Four are already known:

1. **Benefit tables print seven rows for a 1d6 roll** (p. 173 and each career).
   Read as: row 7 is reachable only via event-granted +1 modifiers (p. 127 and the
   event tables), not via an unstated rank bonus.
2. **Mishap 9 in Colonist cites "the injury table (p.136)"** where every other
   reference in the career cites p. 119 (p. 174). Read as a typo for p. 119; the
   book has no injury table at p. 136.
3. **The aging tables for TL 11 and TL 12–13 have no open-ended final row** (p. 122)
   — TL 11 stops at terms 49–58, TL 12–13 at 44–58 — where the TL 9 and TL 10 tables
   end in `12+` and `49+`. This is not an omission. The highest maximum-terms value
   printed anywhere in the origin charts is 58, on Earth (250 years / 58 terms,
   TL 13, p. 43); every other world caps at 38 or below. The closed bands are
   bounded by the highest term count the rules can produce. No reading is needed,
   and the engine should treat a term past the last band as an error rather than
   silently skipping the check — reaching one means the data file is wrong.
4. **Step 18 says a natural 12 on survival forces continuation** (p. 125), and Step
   12's mishap branch says a mishap ends the career. Read as: the natural 12 is a
   success, so the two never collide.

These go in `ERRATA.md` on first implementation, not here.

## CLI sketch

```
cschargen new [--seed N] [--auto] [--data file] [--career colonist]
              [--species human] [--homeworld <id>] [--name X] [-o file]
cschargen batch --count 20 --auto --data file [-o dir|file.jsonl]
cschargen render [--history] character.json
cschargen replay [--ignore-provenance] character.json
cschargen data validate <file>
cschargen version
```

`new` writes JSON to stdout unless `-o`. `batch` emits JSONL, requires `--auto`,
and derives each member's seed from the base seed plus index, recorded per record.
Without `--data`, the engine uses the bundled sample and stamps the record
accordingly, so a sample-data character is never mistaken for a real one.
`data validate` exists because the file is user-authored and a bad one should fail
with a line number, not a panic mid-lifepath.

Interactive mode walks the twenty steps; auto mode applies a fixed policy recorded
in `POLICY.md` and identified by `policy_version`.

**The auto policy needs a term-count stop rule, and it is a policy decision rather
than a rule.** A character from a long-settled world may serve thirty-five terms
(p. 42, Toku), and the rules permit running to that cap. Thirty-five terms of
Colonist is not a usable NPC. Proposal: auto mode stops at four terms by default,
`--terms N` overrides, and `--terms max` runs to the homeworld cap. The rules
impose no such limit and the PRD should not pretend otherwise.

## Architecture notes

- Repo: new `philoserf/cschargen`, public. Code does not live in the Traveller
  collection, per the `t5chargen` and `ctworldgen` precedent.
- Packages: `dice`, `chargen` (engine; a `Decider` interface owns every choice
  point), `career` (data-driven definitions, in-repo), `setting` (the external file:
  loading, validation, hashing), `origin` and `species` (readers over its two
  sections), `render`, `cmd/cschargen`.
- Data/logic boundary: tables, thresholds and labels are data; orchestration and
  the career-specific exceptions are typed Go code. No rules language.
- Zero external dependencies, as in `t5chargen`. (`ctchargen` takes a JSON Schema
  library; this repo follows `t5chargen`'s hand-written checker instead, per the
  validation note above.)
- Testing: golden fixtures per career, replay round-trips, schema validation,
  property tests on the dice engine, and a table-completeness test asserting each
  career's tables match the shape its data declares. The shape genuinely varies and
  the test must not assume otherwise: National Navy (p. 233) has five named skill
  tables plus four assignment tables and a Commission target, while Vagabond
  (p. 296) has no Advanced Education table at all and no enlistment throw. What is
  invariant is narrower — a d66 event table has 36 entries, a 2d6 mishap table has
  11, any skill table that exists has 6 rows, a benefit table has 7.
- Coverage ratchet, `THEORY.md`, `WALKTHROUGH.md`, and a local gate CI runs
  verbatim — the `dusk`/`ctchargen` house pattern.

## Milestones

1. **Walking skeleton.** Dice engine, characteristics (FR1), the term loop for one
   career (Colonist, pp. 173–176), muster out, character record and render, with
   the generation log wired in from the start. Human only; homeworld supplied by
   flag as a literal tech level and term cap; Steps 1, 3–8 stubbed. This is the
   whole engine shape in one career.
2. **Origins.** The setting data format, schema, validator, sample file, and Steps
   3–4 (FR3) including background skills, primary language, and the term and age
   caps feeding FR9 and FR10.
3. **All thirty-four careers** (FR7, FR8), with the career data format proven
   against the awkward ones first, since they are what sets the format: Vagabond
   (no enlistment, no Advanced Education table), Prisoner (unenterable by choice),
   National Navy (commission, officer table, military events), Altrant/Uplift Slave.
4. **Aging** (FR9) — the aging tables, the crisis rules, apparent age.
5. **Pre-career.** Family (FR4), youth and teenage events (FR5), higher education
   (FR6). Roughly ninety pages of tables; the largest milestone by volume and the
   one that can be deferred longest without the tool being useless.
6. **Species** (FR2) — altrants and uplifts, their characteristic methods, caps,
   aging tables, and the homeworld permission and slavery interactions.
7. **Finishing** (FR12), interactive mode, batch, beta.

Milestones 1–4 produce a tool that generates a human adult character with a real
career history. That is the useful midpoint, and it is reachable without touching
the pre-career ninety pages.
