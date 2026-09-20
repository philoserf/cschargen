package chargen

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Player is the interactive decider: it prints a choice and reads an
// answer.
//
// It writes to a separate stream from the record, because `cschargen new`
// writes the record to stdout so that it can be piped. So the prompts go to
// stderr and the answers come from stdin, and a record generated
// interactively is the same file an auto one is.
type Player struct {
	in  *bufio.Reader
	out io.Writer

	// failed remembers a write that did not land. A prompt that cannot be
	// printed is a question that was never asked, and answering it would
	// be answering something the player never saw.
	failed error
}

// NewPlayer builds an interactive decider over a reader and a writer.
func NewPlayer(in io.Reader, out io.Writer) *Player {
	return &Player{in: bufio.NewReader(in), out: out}
}

// Kind implements [Decider].
func (*Player) Kind() DeciderKind { return DeciderPlayer }

// Choose implements [Decider]. It prints the prompt, the options and the
// page they come from, and reads a number.
//
// An answer that is not a number or not in range is re-asked rather than
// refused: a typo is not an abandoned session. End of input is, and comes
// back as an error, which is what ends generation.
func (p *Player) Choose(ask Choice) (int, error) {
	if len(ask.Options) == 0 {
		return 0, ErrNoOptions
	}

	// A choice with one option is not a choice. The engine mostly declines
	// to ask those, but a table that narrowed to one still arrives here.
	if len(ask.Options) == 1 {
		p.printf("\n%s\n  %s (the only option)\n", p.heading(ask), ask.Options[0])

		return 0, p.failed
	}

	p.printf("\n%s\n", p.heading(ask))

	for i, option := range ask.Options {
		p.printf("  %d. %s\n", i+1, option)
	}

	return p.read(len(ask.Options))
}

// Ask implements [Asker]: an open question with no list to answer from,
// which Step 20 is four of.
//
// An empty answer is an answer. "The engine records them and leaves them
// empty rather than inventing them" applies to the player too: a player who
// does not want to name their character yet presses return.
func (p *Player) Ask(question Question) (string, error) {
	heading := question.Prompt

	if question.Of > 1 {
		heading += fmt.Sprintf(" (%d of %d)", question.Nth, question.Of)
	}

	if question.Cite != "" {
		heading += "  [" + question.Cite + "]"
	}

	p.printf("\n%s\n> ", heading)

	if p.failed != nil {
		return "", p.failed
	}

	line, err := p.in.ReadString('\n')
	if err != nil && line == "" {
		// The end of input is not a refusal here. Step 20 is the last step
		// and its fields are allowed to be empty, so a session that ends
		// on them ends with a finished character.
		return "", nil
	}

	return strings.TrimSpace(line), nil
}

// printf writes to the prompt stream and remembers the first failure.
func (p *Player) printf(format string, args ...any) {
	_, err := fmt.Fprintf(p.out, format, args...)
	if err != nil && p.failed == nil {
		p.failed = err
	}
}

// heading is the prompt, the page it comes from, and where in a run of
// identical questions this one falls.
func (p *Player) heading(ask Choice) string {
	heading := ask.Prompt

	if ask.Of > 1 {
		heading += fmt.Sprintf(" (%d of %d)", ask.Nth, ask.Of)
	}

	if ask.Cite != "" {
		heading += "  [" + ask.Cite + "]"
	}

	return heading
}

// read asks until it gets a number in range, or until the input ends.
func (p *Player) read(options int) (int, error) {
	for {
		p.printf("> ")

		if p.failed != nil {
			return 0, p.failed
		}

		line, err := p.in.ReadString('\n')
		if err != nil && line == "" {
			return 0, ErrPlayerGone
		}

		chosen, ok := numberIn(strings.TrimSpace(line), options)
		if ok {
			return chosen, nil
		}

		p.printf("  a number from 1 to %d, please\n", options)

		if err != nil {
			return 0, ErrPlayerGone
		}
	}
}

// numberIn parses an answer and holds it against the list it answers.
func numberIn(answer string, options int) (int, bool) {
	chosen, err := strconv.Atoi(answer)
	if err != nil || chosen < 1 || chosen > options {
		return 0, false
	}

	return chosen - 1, true
}
