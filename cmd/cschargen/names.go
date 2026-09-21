package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
)

// nameDrawStream keeps the name draw away from the dice, the way the batch
// policy's career draw is kept away: a name is chosen before generation
// begins and must not shift a single throw.
const nameDrawStream = 0x6e_61_6d_65_73_00_00_01

// readNames loads a list of names, one per line. Blank lines are skipped
// and a line beginning with # is a comment, so a referee can keep their
// sector's names in one file with notes in it.
//
// FR12 leaves Step 20's four fields empty in auto mode rather than
// inventing them, and that is right for a player character: the name is
// theirs to choose. For a cast of NPCs it is the one thing a referee cannot
// script, so the engine will use names it is given -- never names it made
// up.
func readNames(path string) ([]string, error) {
	if path == "" {
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	// Closing a file that was only read cannot fail in a way that changes
	// the names already scanned out of it.
	defer func() { _ = file.Close() }()

	var names []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		names = append(names, line)
	}

	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	if len(names) == 0 {
		return nil, usagef("%s holds no names", path)
	}

	return names, nil
}

// nameFor draws one name from the list, by the seed the character is
// generated with. It is `new`'s draw: one character, one name, and the seed
// reproduces it.
//
// A batch does not use it. It used to, on the reasoning that batch member i
// and `new --seed base+i` had to be the same character, name included --
// but v0.1.0-alpha.3 gave that up on purpose ("a batch member is no longer
// `new --seed base+i`"), and drawing each name independently of the others
// was the only thing still holding the duplicates in place.
func nameFor(names []string, seed uint64) string {
	if len(names) == 0 {
		return ""
	}

	return names[rand.New(rand.NewPCG(seed, nameDrawStream)).IntN(len(names))]
}

// dealNames puts the list in an order and deals from it, so that a cast of
// twelve drawn from twelve names is twelve people rather than nine.
//
// A list shorter than the batch is dealt again from the top, reshuffled, so
// a hundred NPCs from twenty names still reads as a crowd rather than as
// five copies of the same twenty in the same order.
//
// The shuffle runs on nameDrawStream, which is the point of that stream: a
// name must not move a single throw.
func dealNames(names []string, count int, base uint64) []string {
	if len(names) == 0 {
		return make([]string, count)
	}

	stream := rand.New(rand.NewPCG(base, nameDrawStream))
	dealt := make([]string, 0, count)

	for len(dealt) < count {
		hand := slices.Clone(names)
		stream.Shuffle(len(hand), func(i, j int) { hand[i], hand[j] = hand[j], hand[i] })

		dealt = append(dealt, hand...)
	}

	return dealt[:count]
}
