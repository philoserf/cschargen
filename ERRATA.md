# Errata and recorded readings

Every place the _Clement Sector Core Character Creation Book_ (© 2026
Independence Games) is silent, ambiguous, or wrong, and what this engine does
about it. Page cites are to that book's printed page numbers.

Two kinds of entry, and the difference matters:

- **Typo** — the book plainly means something other than what it says, and the
  engine follows the meaning. No judgment is being exercised.
- **Reading** — the book does not say, and the engine had to choose. The entry
  carries the argument, so that a reader who disagrees can see what changes.

Each entry has an identifier. Records stamp the identifiers of every deviation
that applied to them, so a character generated under one reading stays auditable
after the reading changes.

---

## E-1 (typo) — the Injury table is on p. 119, not p. 136

**Where:** Colonist mishap 9 (p. 174); Prisoner events 26 and 44 (pp. 264, 264).

Those three say "roll on the Injury table (p.136)". Every other reference in the
book — Colonist mishaps 2 and 7, Vagabond mishaps 2, 3, 7, 10, the Military
Events table, the Life Events table — cites p. 119, and p. 136 carries no table.

**Engine:** all Injury-table references resolve to p. 119.

---

## E-2 (typo) — the Vagabond career is on p. 296, not p. 290

**Where:** the Injury table's own closing paragraph (p. 119).

It offers a character reduced to 1 in any characteristic the option to "enter the
Vagabond career (p.290)". The Vagabond career begins on p. 296, and p. 290 is
inside the Thief career.

**Engine:** cosmetic; recorded because the cite appears in a rule the engine
implements and a reader checking the page would find the wrong career.

---

## E-3 (reading) — a benefit table's seventh row is reached only by an event's +1

**Where:** every career's Mustering Out Benefits table (Colonist p. 173, National
Navy p. 233, Prisoner p. 261, Vagabond p. 296); the roll is specified on p. 127.

The tables print rows 1 through 7. The roll is "roll 1d6, consult the selected
table" (p. 127), which cannot produce a 7. No rank benefit anywhere in the book
adds to a benefit roll.

What does add is events: several award benefit rolls carrying a +1, and they do
it in two different scopes — "a +1 modifier to all Benefit rolls made in this
career" against "gain three Benefit rolls at +1", which modifies only those
three.

**Reading:** row 7 exists for exactly those modifiers. The engine therefore holds
benefit rolls as a queue of scoped batches rather than as a count plus a career
modifier, and a modified roll above 7 is clamped to 7 — the table has no eighth
row and the book offers no other reading.

**If this is wrong:** the alternative would be an unstated rank bonus, which
would make row 7 routine rather than exceptional and would change the expected
payout of every career.

---

## E-4 (reading) — an injury becomes permanent at the moment it happens

**Where:** the Injury table's recovery clause (p. 119); mustering out (p. 126).

A character who loses a characteristic to injury "may attempt to recover their
score through medical care. The charge for this recovery will be 1d6x10,000
HFCredits. If the character can afford this cost, then the character recovers the
lost amount in full. If the character cannot afford this care at the time of the
injury, the injury becomes permanent."

But credits come from mustering out (p. 126), which happens when the character
leaves the career — after the term the injury occurred in. A character injured in
their first term has no money by construction, so the clause reads as though it
can never fire.

**Reading:** the engine charges against credits held at the moment of the injury,
which is what "at the time of the injury" says. An early injury is therefore
permanent, and the injury record carries a `permanent` flag set then rather than
recomputed later.

**If this is wrong:** the alternative is to defer the decision to mustering out,
which would let a benefit roll retroactively heal an injury the mishap table has
already acted on — and the mishap may have ejected the character from the career
on the strength of it.

---

## E-5 (reading) — Relationship Ratings are part of the term loop

**Where:** Life Events 4, 5, 6 and 7 (p. 120); the NPC rules (pp. 315-323).

Relationship Ratings are defined in the NPC chapter, which this engine does not
implement. But four Life Events write to them: Life Event 5 removes 1d6 × 20
points from an Ally or Contact, or grants an Enemy at -110; Life Event 6 moves
any NPC by 1d6 × 10 in either direction; Life Event 7 grants a Contact at
exactly 40; Life Event 4 orders the four relationship types for removal.

**Reading:** an Ally, Contact, Rival or Enemy carries a numeric rating from the
first term, because the term loop cannot resolve those four results without one.
Generating an NPC around a rating stays out of scope.

---

## E-6 (reading) — "choose another career" is the player's choice, not a transfer

**Where:** Colonist mishaps 3, 4, 5, 6 and 8 (p. 174).

Five Colonist mishaps end "choose another career and a new homeworld". They do
not name a career, and they are not the same as the mishaps that do — Colonist 10
and 12 send the character to Vagabond by name.

**Reading:** an unnamed career change is a choice point, distinct from a named
transfer, and the two are separate effects in the career data for that reason.

**Not yet applied.** With three careers implemented, offering the choice would
offer almost nothing, so the engine records these as unimplemented and leaves the
character where they are. The entry is here because the distinction is already in
the data and will be acted on in milestone 3, when there are careers to choose
between. No record stamps E-6, and none should until then.

---

## E-7 (reading) — an event's skill check is 2d6 plus the skill's level

**Where:** every career's event and mishap tables, e.g. "Roll Gambler 8+"
(Colonist event 24, p. 175), "Roll Melee (Any) 8+" (Prisoner mishap 3, p. 262).

Dozens of table results turn on a skill check, and this book never says how to
roll one. It states the _characteristic_ check plainly — "the player rolls 2d6 and
adds the Characteristic Modifier associated with the relevant characteristic"
(p. 110) — and for anything more points at another book: "See p.20-27 of the
Clement Sector Core Rulebook for more detail on the task system" (p. 13).

**Reading:** the engine rolls 2d6 and adds the character's level in the named
skill, taking the best level across specialties where the book writes "(Any)". A
skill the character does not hold contributes nothing.

**If this is wrong:** the Core Rulebook's task system almost certainly adds a
characteristic modifier as well, which would raise the success rate of every
event check by roughly the average characteristic modifier. It may also apply an
unskilled penalty, which would lower it for a skill the character lacks. Both
would change outcomes without changing the shape of the engine — `rollCheck` is
one function.

---

## E-8 (reading) — the Rank 0 benefit applies on entering a career

**Where:** rank benefits (p. 116); starting rank (p. 115).

Every career's rank table prints a Rank 0 row with a benefit in it. p. 116 says
benefits "are applied immediately upon achieving that Rank", and lists four ways
of achieving one — an advancement roll, an event, a mishap, a career transfer
that retains rank. Entering a career is not among them. But p. 115 says "All
characters begin a career at Rank 0 unless otherwise stated", and the printed
Rank 0 row would otherwise be unreachable.

**Reading:** beginning a career at Rank 0 is achieving Rank 0, so the row applies
on entry.

**If this is wrong:** every character would be one rank benefit poorer per career
entered — for a Colonist Settler, Animals (Farming) 1.

---

## E-9 (reading) — a homeworld reassigned in adulthood changes the tech level and nothing else

**Where:** Colonist mishaps 3, 4, 5, 6, 8, 9 and 12 (p. 174); background skills
(p. 40); maximum age and terms (p. 42); the aging tables (pp. 122-123).

Six of Colonist's eleven mishaps end "choose another career and a new homeworld",
and several other results reassign one. A homeworld carries four things — the
background skills, the primary language, the maximum age and terms, and the tech
level the aging throws are gated by. The book never says which of them a
reassignment changes.

**Reading:** the tech level, and nothing else.

Background skills are "skills which your character is assumed to have learned
through living their normal life on this world" (p. 40), which is a childhood; an
adult who is deported does not acquire a new one. The maximum terms "represents
the longest possible career history that character could have accumulated before
play begins" and is drawn from the settlement history of where they were **born**
(p. 42). But the aging tables are indexed by "the tech level of the character's
homeworld" (p. 121), present tense, and a character living on a lower-tech world
ages by its medicine rather than by the one they left.

So the record holds a homeworld history: which world, in which subsector, at
which tech level, from which term, and why they moved.

**If this is wrong:** granting background skills on every reassignment would make
a much-deported Colonist the most broadly skilled character in the game, from
mishaps. Carrying the birth world's tech level instead would make a character's
aging independent of where they actually live, which is the opposite of what
p. 121 says.

---

## E-10 (reading) — System Defense Forces (Navy) prints no Mishap table

**Where:** pp. 278-281, against Troopers (p. 284) and Wet Navy (p. 289).

Every career in the book prints a 2d6 Mishap table. This one does not: p. 279 ends
with its rank tables and a note about rank titles, and p. 280 begins its d66
events. Both of its sibling System Defense Forces careers print one.

A failed survival throw sends a character to the Mishap table (p. 113), so
without one the career has no failure path at all.

**Reading:** the career uses National Navy's mishap table (p. 235).

The two careers are otherwise the same tables cell for cell — the same four
assignments with the same skill rows, the same Service, Advanced Education and
Officer tables, and the same two rank tables (compare pp. 233-234 with pp.
278-279). What differs is the enlistment throw, the commission target and cash
figures at roughly half. A career that shares every other table with National
Navy almost certainly shares the missing one.

**If this is wrong:** the alternative is that the omission is deliberate and the
career cannot suffer a mishap, which would make it the only career in the book
where failing survival costs nothing — a strictly better career than its two
siblings, for no stated reason.

---

## E-11 (limit) — one career's prerequisite cannot be checked

**Where:** System Defense Forces (Wet Navy), p. 287.

"Your homeworld must have a Hydrographics score of 4+ to join this career."

Hydrographics is a digit of a world's Universal World Profile, which the origin
charts of pp. 43-56 do not print — they carry tech level and nothing else of the
UWP. The setting data therefore has no hydrographics figure to check against.

**The engine records the prerequisite on the career and does not enforce it.**
Refusing a career on a rule this engine cannot evaluate would be worse than
admitting it cannot. This is a limit rather than a reading: nothing is being
interpreted, and it closes when the setting data carries a UWP.

---

## E-12 (limit) — a result that names a difficulty cannot be resolved

**Where:** Arts event 63 (p. 158), and every later result that does the same.

Most checks in a career's tables name a target number: "Roll Gambler 8+". A few
name a **difficulty** instead — "Make a Diplomat check at Very Difficult",
"Routine", "Easy" — which are the Core Rulebook's task system (p. 13 sends the
reader there), and that book is outside this ruleset.

**The engine records these results rather than resolving them.** A difficulty
cannot be turned into a target number without the table that maps them, and
guessing one would be inventing a rule.

Like E-11 this is a limit rather than a reading: nothing is being interpreted,
and it closes if the difficulty ladder is ever brought into scope.

---

## E-13 (typo) — one d66 result is printed blank

**Where:** Scientist event 26 (p. 271).

The table prints "26" with nothing beside it. Result 25 runs long — it carries a
1d6 sub-table of its own — and 26 sits on its own line between the end of that
and the 31-36 Life Event span.

**Engine:** the result is recorded as having no printed effect and nothing
happens. It is a transcription note rather than a reading: there is no rule to
interpret, and inventing one for a blank row would be worse than a term in which
nothing did.

**If this is wrong:** the likeliest explanation is that result 25's sub-table
overran its row and displaced 26's text, in which case the missing entry is lost
rather than absent. Either way the engine cannot supply it.

---

## E-14 (typo) — a result that prints a target number with nothing to throw

**Where:** Scavenger event 41 (p. 267).

The row reads "You must stay aware of all the local laws and regulations
concerning your business. Gain a level in Advocate (Legal) 8+". There is no
throw in it: nothing succeeds or fails, and no consequence hangs on either.
Every other result in the book that prints "8+" prints "Roll <skill> 8+" and
two branches beneath it.

**Engine:** the result grants a level in Advocate (Legal) and the "8+" is
dropped.

**If this is wrong:** the row lost a check in layout, and the level is the
reward for passing rather than the result itself. That would make the result
strictly worse than transcribed, never better, so a character generated under
this reading is never owed a skill they did not earn.

A second, smaller slip runs across several tables: Scavenger event 26 (p. 267)
and Exotic event 63 (p. 193) both say "Engineering (Any)" where the skill list
(p. 304) has Engineer. The engine uses Engineer, because no skill named
Engineering exists to grant.
