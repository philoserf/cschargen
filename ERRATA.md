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

**The rating is now the mechanic p. 320 describes**, not a field nothing sets.
Milestone 5 made it one, because the youth and teenage tables move ratings on
almost every row and could not be transcribed on top of a number that never
changed. The four Life Events above carry it out rather than recording it. See
[E-22](#e-22-reading--a-tie-is-born-with-the-kind-its-result-names) and
[E-23](#e-23-reading--where-a-tie-with-no-printed-rating-starts).

---

## E-6 (reading) — "choose another career" is the player's choice, not a transfer

**Where:** Colonist mishaps 3, 4, 5, 6 and 8 (p. 174).

Five Colonist mishaps end "choose another career and a new homeworld". They do
not name a career, and they are not the same as the mishaps that do — Colonist 10
and 12 send the character to Vagabond by name.

**Reading:** an unnamed career change is a choice point, distinct from a named
transfer, and the two are separate effects in the career data for that reason.

**Applied since all thirty-four careers are transcribed.** The effect
(`EffectChooseCareer`) ends the career without naming the next one; Step 18
(p. 125) then offers the career list and the character enlists on the ordinary
terms, enlistment throw and all. That is the difference from a transfer, which
places the character in a named career without one.

Until milestone 3 the engine recorded these as unimplemented and left the
character where they were, because three careers to choose between was not a
choice. Records generated before that stamp no E-6; records generated after it
stamp E-6 wherever one of these results is reached.

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

---

## E-15 (reading) — two aging tables stop, and nothing is printed past them

**Where:** the aging tables for tech level 11 and tech levels 12-13 (p. 122).

Four human aging tables are printed. Two of them end on an open band — TL 9's
last row is "12+", TL 10's is "49+" — and two do not. TL 11's last row is 49-58
and TL 12-13's is 44-58, and the page prints nothing for a character past term
58 on either.

**Engine:** no check is made there. A TL 11 character who serves a 59th term
makes no aging throw at all, which is what the page says rather than what it
probably means.

**If this is wrong:** the intended reading is almost certainly that the last
band stays open, the way the other two tables' do, and a very long-lived
character is being let off. That direction is the forgiving one — a character
generated under this reading is never harmed by it — and it is reachable only
past term 58, which no world in the sample setting data allows.

---

## E-16 (reading) — nothing is printed above tech level 13

**Where:** the aging tables (pp. 122-123), against `setting`'s world schema.

The tables are headed "Tech Level 9 or less", "Tech Level 10", "Tech Level 11"
and "Tech Level 12-13". The setting data format accepts a world of any tech
level up to 20, and the book's own setting reaches TL 13 but the schema does not
stop there.

**Engine:** a homeworld above TL 13 reads the TL 12-13 table.

**Reasoning:** each band in the sequence delays the onset of aging and none of
them abolishes it — TL 9 begins at term 6, TL 10 at 18, TL 11 at 33, TL 12-13 at 44. Reading the highest printed band forward continues that trend; reading
"nothing printed" as "no aging" invents a rule the trend does not support, and
would make a TL 14 homeworld strictly better than any the book describes.

**If this is wrong:** a TL 14+ character ages slightly sooner than intended. No
world in the sample setting data is above TL 12, so this is reachable only in a
user-supplied file.

---

## E-17 (reading) — which characteristics are physical and which are mental

**Where:** the Aging Crisis rules (pp. 123-124), read against p. 14.

Four of the terminal states turn on the words "physical" and "mental":

> If two or more physical characteristics reach 0 ... If all three physical
> characteristics are reduced to 0 ... If any one mental characteristic is
> reduced to 0 ... If two or more mental characteristics are reduced to 0

That section enumerates neither group. **p. 14 enumerates one of them**, in the
rule this one restates: "In the highly unlikely event that all three physical
characteristics (Strength, Dexterity, and Endurance) are reduced to 0 during
character creation, the character is considered to have died before play
begins."

**Engine:** STR, DEX and END are physical; INT, EDU and CHA are mental. There are
six characteristics and the book uses only these two words for them, so naming
the physical three names the other three by subtraction.

This is closer to a cross-reference than a reading, and it is here because the
aging section can be read on its own and does not say it. The aging tables agree
independently: every band begins with STR, DEX and END and adds INT, EDU and CHA
as the character gets older.

---

## E-18 (reading) — two crises in one term is two payments

**Where:** the Aging Crisis (p. 123).

A term can reduce more than one characteristic to 0: the last band of the TL 9
and TL 10 tables makes six checks. The crisis rule is written for one — "the
affected characteristic is restored to 1" — and says nothing about two.

**Engine:** each characteristic that reaches 0 is its own crisis, with its own
1d6 × 1000 credits and its own choice.

**Reasoning:** the price is for treating a characteristic, not for surviving a
term. One payment covering both would make the second free, which no reading of
"emergency treatment" supports.

---

## E-19 (reading) — a character who cannot afford the treatment

**Where:** the Aging Crisis (p. 123).

> The player may pay 1D6 × 1000 credits for emergency treatment. ... If payment
> is made, the affected characteristic is restored to 1.

The book sets a price and treats payment as a choice. It does not say what
happens to a character who would pay and cannot: cash is the one benefit a
character can be entirely without, and 1d6 × 1000 is a real sum early in a
lifepath.

**Engine:** treatment that is not paid for is treatment that is not given, so
this resolves the same way as declining it — the character dies.

**If this is wrong:** the intended reading is probably that the price is a
formality and the treatment is always available, in which case a character the
engine killed should have lived. That is the direction to check first if a
generated record looks unfair: the log names the price and the character's
balance at the moment they died.

---

## E-20 (reading) — the apparent age chart's two ends

**Where:** the Apparent Age chart (p. 125).

The chart's first row is 30-40 and its last is 281-290. A character below 30 is
off the top of it and one above 290 is off the bottom, and the page says nothing
about either.

**Engine:** below age 30, apparent age is the character's age; above 290, the
last printed row holds.

**Reasoning:** p. 124 already gives the rule for the other case where the chart
does not apply — "If, by chance, your character happens to be from a world of
Tech Level 9 and lower, your apparent age and your real age are the same" — and
below 30 the same thing is true for a reason the chart itself shows: every
column's first row is 20-25 against an actual age of 30-40, so the divergence has
barely begun. A character of 18 who "appears to be" 20-25 would be made older by
the chart, which is the opposite of what it does.

Above 290 the reading is simply that the last row holds, because the alternative
is a character with no apparent age at all.

---

## E-21 (reading) — "over 40" against a chart that gives bands

**Where:** twelve careers' enlistment throws, e.g. Exotic (p. 190), Organized
Crime (p. 246), Investigator (p. 216).

> If you have an apparent age of over 40, take a -2 modifier to this roll.

The Apparent Age chart does not give a number. It gives a five-year band — 35-40,
40-45 — so "over 40" has to be read against a range that can straddle it.

**Engine:** the modifier applies where the band's **lower** bound is 40 or more.
A character whose apparent age is 35-40 is not over 40; one at 40-45 is.

**Reasoning:** the alternative — any band that reaches 40 — would catch 35-40,
and a character at the bottom of that band appears to be 35. Reading the lower
bound makes the modifier apply to characters who are over 40 on any reading of
their band, which is the conservative direction for a penalty.

Sports (p. 273) prints "If the character's apparent age is 40+" where the other
eleven print "over 40". The engine treats them as one modifier at one bar: 40+
and "over 40" differ only for a band whose lower bound is exactly 40, and that
band is 40-45, which both readings catch.

---

## E-22 (reading) — a tie is born with the kind its result names

**Where:** the Relationship Rating bands (p. 320), against the results that grant
ties (pp. 57-61, 68-84, 120, and every career's tables).

p. 320 divides the scale into four bands and names each: Contacts are 1-100,
Allies 101-200, Rivals -1 to -100, Enemies -101 and below. Read as a definition,
the band decides what a tie is.

The results that grant ties do not read that way. Step 5 grants a grandparent
"as an Ally with a Relationship Rating of 100" (p. 61) — and 100 is the top of the
Contact band. Youth Path 1 result 18 grants "a Contact with a Relationship Rating
of 90", which agrees; result 12 grants "an Ally with a Relationship Rating of 125",
which also agrees. Only the family numbers disagree, and they disagree twice over:
p. 320 also says "Family members will always begin as Allies with a Relationship
Rating of 150", where Step 5 says 125 for parents and siblings and 100 for
everyone else.

**Engine:** the kind is stored, and the result that grants a tie decides it. The
bands govern movement.

**Reasoning:** every band sentence on p. 320 is about motion — "Contacts which
fall to", "Allies whose Relationship Rating drops to", "Enemies whose rating
rises" — and none of them is about creation. A specific instruction beats a
general rule, and a named kind beats an inferred one. Where the two family
numbers disagree, the step that grants the tie wins over the summary on p. 320.

**If this is wrong:** a character's grandparents are Contacts rather than Allies,
one point of Relationship Rating below the line.

---

## E-23 (reading) — where a tie with no printed rating starts

**Where:** most results in the book that grant a relationship, e.g. "Gain an
Ally" (Adventurer event 13, p. 148).

The book gives a rating for perhaps a dozen ties and not for the rest. A rating
is now load-bearing — it is what decides whether a Contact survives a bad term —
so a tie with no printed one has to start somewhere.

**Engine:** the middle of its band. An Ally starts at 150, a Contact at 50, a
Rival at -50, an Enemy at -150.

**Reasoning:** it is the only choice that is not an argument about which end of
the band a nameless relationship belongs at, and it happens to agree with p. 320's
"Family members will always begin as Allies with a Relationship Rating of 150".
The alternative — starting every tie at a band edge — would make one direction of
movement free and the other instantly fatal.

**If this is wrong:** ties are a little more robust or a little more fragile than
intended, uniformly. The ratings the book does print are used as printed, so this
never overrides a number on the page.

---

## E-24 (reading) — a relationship ends rather than inverting

**Where:** the Relationship Rating bands (p. 320).

> Contacts which fall to a Relationship Rating of 0 have decided to no longer be
> associated with the character and are lost. ... Allies which drop immediately to
> 0 or less are no longer associated with the character and are lost. ... Enemies
> whose rating rises immediately to 0 or higher have decided they no longer care
> about the character and are lost.

Read together, those three sentences say a tie never crosses zero: a positive
relationship that goes negative is lost rather than becoming a Rival, and a
negative one that goes positive is lost rather than becoming a Contact. The book
states it for Contacts, Allies and Enemies and not for Rivals, whose sentence
instead reads "Rivals which rise to 0 or higher have decided that their rivalry is
no longer valid, and they are lost to the character" — which is the same rule.

**Engine:** a rating change that would take a tie to zero or across it loses the
tie, and the change is classified once on the final rating rather than band by
band. An Ally at 110 who takes -120 is gone; they are not a Rival at -10.

This is closer to a reading of three sentences together than an interpretation of
any one of them, and it is here because the obvious implementation — move the
number, look up the band — gets it wrong in exactly the cases the page bothered to
call out.

---

## E-25 (limit) — a sibling's age is not checked against the parents'

**Where:** the Sibling Age table (p. 60).

> Always compare the result to the age of the parents at that point. If one of
> the parents will have been less than 16 years of age when the child will have
> been born, then reduce the age difference to account for this.

An older sibling pushes a parent's age at that birth back by the age difference,
and a 3d10 result can push it back by thirty years. The rule is a floor on how
young a parent may have been.

**Engine:** the age difference is generated and not adjusted.

**Why it is a limit rather than a reading:** the check needs a parent's age at
the sibling's birth, and Step 5 records their age at the **character's** birth
for each parent separately. With more than two parents — a communal household can
have a dozen — "one of the parents" does not identify which, and the rule gives no
way to choose. A household where every parent was over 46 satisfies it
automatically; one where a parent was 17 cannot have a sibling more than a year
older, whichever parent that was.

The consequence is cosmetic: a sibling may be recorded as older than a parent
could plausibly have made them. Nothing downstream reads a sibling's age.

---

## E-26 (limit) — "the characteristic that you used to get into this path"

**Where:** Youth Path 2 result 6 (p. 70).

> Gain +1 in the characteristic that you used to get into this path but lose 1
> from another physical characteristic. For instance, if you gained access to this
> path by your STR being 8+, lose one from DEX or END.

Path 2 is open on "STR 8+ or END 8+", and a character may satisfy both. The
result assumes exactly one characteristic let them in.

**Engine:** both halves are a choice. The character picks which of STR and END
opened the path, and then which of the other two physical characteristics to
lose.

**Why it is a limit rather than a reading:** nothing in the record says which
requirement was checked, because the path was open and the engine did not need to
know which clause opened it — and for a character with both at 8 or higher there
is no fact of the matter to record.

The choice is the honest form of the rule. Under the auto policy the first option
printed wins, which is STR, the first clause of the requirement.

---

## E-27 (reading) — the campaign's present year is given twice

**Where:** p. 43 and p. 124.

Two of Step 7's four paths are gated on whether the homeworld "has been
colonized or established for 100+ standard years" (p. 76), which needs a year to
measure a world's settlement date against. The book gives two:

> the character could not have immigrated to Clement Sector after 2331 when the
> Conduit collapsed and they must have present in Clement Sector for the 19 years
> between the 2331 collapse and the present day of 2350 (p. 43)

> To the people living in 2345, it is perfectly normal for an average 75-year-old
> to appear to be and have the general health of an early 21st century person in
> their 30s (p. 124)

**Engine:** 2350, which is p. 43's. That page is about dates and their
consequences for character generation, and it derives the year arithmetically
from 2331 plus nineteen; p. 124's is an aside in a passage about how people
perceive age.

A setting data file may state its own `presentYear`, because a campaign is
entitled to be set somewhere else in the timeline, and the settlement dates in
that file are the ones it will be measured against.

**If this is wrong:** a world settled in exactly 2245 through 2250 falls on the
other side of the hundred-year line, and its characters take Teenage Path 2
rather than Path 1. Nothing else in the engine reads the year.

---

## E-28 (typo) — medical school's success throw is printed twice, differently

**Where:** the Medical School box (p. 100) and the prose beneath it (p. 101).

The box reads:

> Admission INT 8+ / Success EDU 9+ / Honors INT 10+

and the paragraph on the next page reads:

> The player must roll 8 or higher on 2d6 for the character to succeed in medical
> school. Add any EDU characteristic modifier that the character may have to this
> roll.

**Engine:** 8+, which the prose gives.

**Reasoning:** every other institution's prose and box agree, and where they do it
is the prose that spells out which characteristic modifies the throw and what the
number means. The three other tracks all print 8+ for success. A 9+ here would
make medical school harder to finish than the military academy, which is not
something anything else on the page suggests.

**If this is wrong:** medical school is one point easier than intended, and a
character who scraped through on an 8 should have washed out. The log names the
throw and its target, so a record can be checked against either reading.

---

## E-29 (reading) — a character does not take two bachelor's degrees

**Where:** Step 8 as a whole (pp. 85-104).

The step is written as one decision made once. It never contemplates a
character attempting more than one institution, and its only explicit sequel is
the graduate tracks, which require "Success in an Undergraduate University or
Military Academy" (pp. 97, 101) and so cannot be the first thing attempted.

The engine runs the step as a loop, because otherwise Graduate School and
Medical School are unreachable at Step 8 — the degree they need is one the
character earns in the same step. That loop then raises a question the book does
not: may a graduate of Undergraduate College go straight on to the Military
Academy for a second bachelor's?

**Engine:** no. An institution whose degree is a bachelor's is closed to a
character who holds one.

**Reasoning:** the alternative is a character who collects degrees rather than
progressing, and the tracks that follow a bachelor's are the ones the book gives
as what comes next. It also keeps the Military Academy's obligation coherent — a
graduate "must enter a military career in Step 9" (p. 94), which a character who
then spent four years at Undergraduate College plainly has not.

**If this is wrong:** a character who wanted both is offered only one, and can
have the other by attending at Step 18 instead, where "return to higher
education" is one of the choices (p. 125).

---

## E-30 (reading) — one aging table has a term in no band

**Where:** the aging table shared by two uplift species (p. 122), which this
repository calls the `sturdy` profile.

Its three rows are 6-8, 9-11 and 13+. **Term 12 is in none of them.** Every other
table on pp. 122-123 is contiguous: 6-8, 9-10, 11, 12+; 18-28, 29-38, 39-48, 49+;
4-5, 6-7, 8+.

**Engine:** the middle band holds through the gap, so terms 9 to 12 make the same
five checks and term 13 moves to six.

**Reasoning:** the alternative reading is that term 12 makes no checks at all, on
a table whose whole shape is that ageing accelerates. A character would age at 11,
not at 12, and then harder at 13 — which is not something any other table on the
page does.

**If this is wrong:** the intended row is "12+" rather than "13+", and a character
in their twelfth term faces six checks rather than five. That is the only other
plausible transcription, and it is the harsher one — so this reading never treats
a character worse than the book intended.
