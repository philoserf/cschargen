#!/bin/sh
# Issue references cited in comments and documents, checked two ways.
#
# It FAILS on a reference that does not resolve -- a typo'd or deleted number
# is a pointer to nothing and there is no legitimate reason to have one.
#
# It REPORTS, without failing, every reference to a closed issue. Most are
# correct: "#63 gave the flag a note", "the finding that opened issue #42"
# are history, and history is the good use of a number. The bad use looks
# almost the same -- career/build.go's skillAt said "the engine's reading of
# it is issue #20" and went on saying it for four releases after #20 was
# fixed, telling every reader the code beneath it was known-wrong when it was
# not. No pattern separates the two, so this prints the list and a person
# decides. Five lines once per release is cheaper than a check nobody trusts.
#
# Not part of `task`: it needs the network and a GitHub token, and the gate
# runs offline and identically in CI. It belongs to the release pass.
set -eu

refs=$(git grep -ho '#[0-9][0-9]*' -- '*.go' '*.md' ':!CHANGELOG.md' |
	tr -d '#' | sort -un)

[ -n "$refs" ] || { echo "no issue references"; exit 0; }

closed="" missing=""
for n in $refs; do
	state=$(gh issue view "$n" --json state --jq .state 2>/dev/null || echo MISSING)
	case "$state" in
	OPEN) ;;
	CLOSED) closed="$closed $n" ;;
	*) missing="$missing $n" ;;
	esac
done

if [ -n "$closed" ]; then
	echo "closed issues cited -- check each still reads as history, not as open work:"
	for n in $closed; do
		git grep -n "#$n" -- '*.go' '*.md' ':!CHANGELOG.md' | sed 's/^/  /'
	done
	echo ""
fi

if [ -n "$missing" ]; then
	echo "references that do not resolve:$missing" >&2
	for n in $missing; do
		git grep -n "#$n" -- '*.go' '*.md' ':!CHANGELOG.md' | sed 's/^/  /' >&2
	done
	exit 1
fi

echo "$(echo "$refs" | wc -l | tr -d ' ') references, all resolve"
