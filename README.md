# cschargen

A Go CLI that generates rules-accurate Clement Sector characters.

Ruleset baseline: **_Clement Sector Core Character Creation Book_**, © 2026
Independence Games, author John Watts. Every implemented rule carries its printed
page cite, and every place the text is silent or ambiguous has a recorded reading.

Clement Sector is Cepheus Engine-compatible, not Traveller: the lineage runs
Mongoose Traveller SRD (2008) → Cepheus Engine SRD (Samardan Press, 2016) →
_Clement Sector: The Rules_ (2016) → this book. This repository is not affiliated
with Independence Games, Samardan Press, Mongoose Publishing or Far Future
Enterprises.

## Status

Design only. `docs/PRD.md` is the v1 contract; no engine yet. Milestone 1 is a
walking skeleton — dice engine, characteristics, the term loop for one career.

## What generation involves

The book calls its own character creation a minigame that "often requires a game
session of its own" (p. 16), and it is considerably deeper than Classic Traveller
or baseline Cepheus:

- Twenty steps, eight of them before the character enters a first career —
  species, characteristics, origin, family, youth events, teenage events, and
  four tracks of higher education.
- Thirty-four careers, each with assignment tracks, a 2d6 mishap table, a d66
  event table, rank tables, and four or more skill tables.
- Aging indexed by **term number** and gated by the homeworld's tech level, so a
  high-tech character makes no aging check until term 44 and may serve 58 terms.
- Apparent age tracked separately from actual age — the setting's signature
  number, derived rather than rolled.

## Setting data is not in this repository

The book declares its game mechanics Open Game Content and declares as Product
Identity "all subsector names, world maps, world names, system names, system
maps, vehicle names, starship names, starship classes, artistic depictions of
ships and vehicles, and organizations", adding that the term "altrant" is not
open content (OGL §16, p. 335).

Those names are load-bearing here: a homeworld determines background skills,
primary language, maximum age, maximum terms, and whether engineered characters
may be born there at all.

**So this repository ships no Product Identity.** World, subsector and species
names live in an external data file you supply from your own copy of the book:

```sh
cschargen new --data ~/path/to/setting.json --seed 7
```

A small invented sample (`data/setting.sample.json`) ships with the repo so the
tests and the demo run without it. Characters generated against the sample are
valid and replayable, and are stamped as such in their record so they are never
mistaken for characters set on real worlds.

The career tables *are* here, because they are mechanics rather than names.

## License

MIT for the code — see `LICENSE`. The rules it implements are the publisher's;
this tool reproduces mechanics, not text.
