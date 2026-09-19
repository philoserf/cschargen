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

**Reading:** an unnamed career change is a choice point offered to the Decider,
resolved against the careers the engine implements. While only three are
implemented it will usually have nothing to offer, and the engine records a
`career_transfer_unimplemented` consequence rather than silently picking one.
This distinction is why the two shapes are separate effects in the career data.
