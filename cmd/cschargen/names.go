package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
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
// generated with.
//
// By the seed rather than by position in the batch, so that the PRD's
// promise holds: batch member i and `new --seed base+i` are the same
// character, name included. Two characters in a batch may draw the same
// name, which is what a list shorter than the batch has to mean and is true
// of people anyway.
func nameFor(names []string, seed uint64) string {
	if len(names) == 0 {
		return ""
	}

	return names[rand.New(rand.NewPCG(seed, nameDrawStream)).IntN(len(names))]
}
