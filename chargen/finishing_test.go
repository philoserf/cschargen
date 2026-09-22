package chargen_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/philoserf/cschargen/chargen"
)

// TestAutoModeInventsNoFinishingTouches. Step 20's four fields are player
// input or policy default and none is rolled, so an auto run records them
// empty rather than inventing them.
func TestAutoModeInventsNoFinishingTouches(t *testing.T) {
	t.Parallel()

	opts := options(t, 5)

	opts.Inputs.TermLimit = 1

	character := generate(t, opts)
	if !character.State.Finishing.Empty() {
		t.Errorf("auto mode invented %+v", character.State.Finishing)
	}

	// And the record says so, rather than saying nothing.
	said := false

	for _, event := range character.Events {
		if event.Consequence == nil {
			continue
		}

		if strings.Contains(event.Consequence.Detail, "the engine invents none") {
			said = true
		}
	}

	if !said {
		t.Error("the record does not say the finishing touches were left empty")
	}
}

// TestTheFinishingTouchesComeFromTheInputs. Four flags, four fields, and
// the name is the one that was already an input.
func TestTheFinishingTouchesComeFromTheInputs(t *testing.T) {
	t.Parallel()

	opts := options(t, 5)

	opts.Inputs.TermLimit = 1
	opts.Inputs.Name = "Vela Ashgrove"
	opts.Inputs.Gender = "she/her"
	opts.Inputs.Appearance = "tall, with a spacer's stoop"
	opts.Inputs.Goals = "to find out what happened to the Kestrel"

	got := generate(t, opts).State.Finishing

	want := chargen.Finishing{
		Name: "Vela Ashgrove", Gender: "she/her",
		Appearance: "tall, with a spacer's stoop",
		Goals:      "to find out what happened to the Kestrel",
	}

	if got != want {
		t.Errorf("= %+v, want %+v", got, want)
	}
}

// scripted answers every choice with the first option and every open
// question from a list, which is what lets a test drive Step 20 without
// counting how many choices the lifepath before it made.
type scripted struct {
	answers []string
	asked   []string
}

func (s *scripted) Choose(chargen.Choice) (int, error) { return 0, nil }

func (*scripted) Kind() chargen.DeciderKind { return chargen.DeciderPlayer }

func (s *scripted) Ask(q chargen.Question) (string, error) {
	s.asked = append(s.asked, q.Prompt)

	if len(s.answers) == 0 {
		return "", nil
	}

	answer := s.answers[0]

	s.answers = s.answers[1:]

	return answer, nil
}

// TestThePlayerIsAskedForTheFinishingTouches. Step 20 is four open
// questions, and only a decider that can be asked answers them.
func TestThePlayerIsAskedForTheFinishingTouches(t *testing.T) {
	t.Parallel()

	player := &scripted{answers: []string{
		"Vela Ashgrove", "she/her", "tall", "to find the Kestrel",
	}}

	opts := options(t, 5)

	opts.Decider = player
	opts.Inputs.TermLimit = 1

	got := generate(t, opts).State.Finishing
	if got.Name != "Vela Ashgrove" || got.Goals != "to find the Kestrel" {
		t.Errorf("= %+v", got)
	}

	if len(player.asked) != 4 {
		t.Errorf("the player was asked %d questions, want 4", len(player.asked))
	}
}

// TestThePlayerSeesTheQuestion. The prompt carries the page and where in
// the four it falls.
func TestThePlayerSeesTheQuestion(t *testing.T) {
	t.Parallel()

	var out strings.Builder

	player := chargen.NewPlayer(strings.NewReader("Vela Ashgrove\n"), &out)

	answer, err := player.Ask(chargen.Question{
		Prompt: "What is the character's name?", Cite: "pp. 129-130", Nth: 1, Of: 4,
	})
	if err != nil {
		t.Fatalf("Ask: %v", err)
	}

	if answer != "Vela Ashgrove" {
		t.Errorf("answer = %q", answer)
	}

	for _, want := range []string{
		"What is the character's name?", "(1 of 4)", "[pp. 129-130]",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("the prompt does not contain %q:\n%s", want, out.String())
		}
	}
}

// TestAQuestionNobodyAnswers. An empty answer is an answer, and so is the
// end of input: Step 20 is the last step and its fields may be empty.
func TestAQuestionNobodyAnswers(t *testing.T) {
	t.Parallel()

	player := chargen.NewPlayer(strings.NewReader(""), &strings.Builder{})

	answer, err := player.Ask(chargen.Question{Prompt: "What do they want?"})
	if err != nil || answer != "" {
		t.Errorf("Ask = %q, %v; want an empty answer and no error", answer, err)
	}
}

// TestAQuestionThatCannotBePrinted ends generation, for the reason a
// choice that cannot be printed does.
func TestAQuestionThatCannotBePrinted(t *testing.T) {
	t.Parallel()

	player := chargen.NewPlayer(strings.NewReader("something\n"), brokenWriter{})

	_, err := player.Ask(chargen.Question{Prompt: "What do they want?"})
	if err == nil {
		t.Error("a question the player never saw was answered anyway")
	}
}

// TestAPlayerNeedNotAnswerStep20. "The engine records them and leaves them
// empty rather than inventing them" is about the player too: an empty
// answer is an answer, and a session that ends on Step 20 ends with a
// finished character.
func TestAPlayerNeedNotAnswerStep20(t *testing.T) {
	t.Parallel()

	opts := options(t, 5)

	opts.Decider = &scripted{}
	opts.Inputs.TermLimit = 1

	character := generate(t, opts)
	if !character.State.Finishing.Empty() {
		t.Errorf("a player who answered nothing got %+v", character.State.Finishing)
	}
}

// TestAFlagAnsweredIsNotAsked. A field the command line supplied is not a
// question, so a player who named their character on the command line is
// not asked for a name.
func TestAFlagAnsweredIsNotAsked(t *testing.T) {
	t.Parallel()

	player := &scripted{answers: []string{"she/her", "tall", "to find the Kestrel"}}

	opts := options(t, 5)

	opts.Decider = player
	opts.Inputs.TermLimit = 1
	opts.Inputs.Name = "Vela Ashgrove"

	got := generate(t, opts).State.Finishing
	if got.Name != "Vela Ashgrove" || got.Gender != "she/her" {
		t.Errorf("= %+v", got)
	}

	if len(player.asked) != 3 {
		t.Errorf("the player was asked %d questions; the name was already given",
			len(player.asked))
	}
}

// refusingAsker declines every open question, which is what an interactive
// session whose input stream broke looks like at Step 20.
type refusingAsker struct{ scripted }

func (refusingAsker) Ask(chargen.Question) (string, error) {
	return "", errNobodyThere
}

var errNobodyThere = errors.New("nobody to ask")

// TestARefusedFinishingQuestionEndsGeneration. Step 20's questions go
// through the Decider like every choice does, and a refusal is a refusal.
func TestARefusedFinishingQuestionEndsGeneration(t *testing.T) {
	t.Parallel()

	opts := options(t, 5)

	opts.Decider = &refusingAsker{}
	opts.Inputs.TermLimit = 1

	_, err := chargen.New(opts).Run()
	if !errors.Is(err, errNobodyThere) {
		t.Errorf("err = %v, want the refusal", err)
	}

	// And the error says which of the four it was, because four identical
	// failures are not four identical questions.
	if err != nil && !strings.Contains(err.Error(), "asking for the name") {
		t.Errorf("err = %v, which does not say which question", err)
	}
}

// TestAnInteractiveCharacterSurvivesReplay is #109. Step 20's four answers
// arrive through Ask, and Replay implements Choose and not Ask -- so a
// replay cannot re-ask them and has to read them out of Inputs. Only the
// name was written back, and the other three came back empty.
//
// It failed two ways, and both are held here. With a name given, the log
// agreed and `replay` reported "identical" while the character had lost
// three fields. With the name left empty, Finishing.Empty() flipped between
// the run and the replay and the log itself diverged.
func TestAnInteractiveCharacterSurvivesReplay(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		answers []string
	}{
		{
			"all four answered",
			[]string{"Vela Ashgrove", "she/her", "tall, grey-eyed", "to find her brother"},
		},
		{
			"no name, the other three answered",
			[]string{"", "she/her", "tall, grey-eyed", "to find her brother"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := options(t, 7)

			opts.Decider = &scripted{answers: tc.answers}
			opts.Inputs.Interactive = true

			original := generate(t, opts)

			// Replay is handed the record's own inputs, exactly as
			// cmd/cschargen does, and a decider that cannot be asked.
			again := options(t, original.Provenance.RNG.Seed)

			again.Decider = chargen.NewReplay(original.Events)
			again.Inputs = original.Provenance.Inputs

			replayed := generate(t, again)

			if replayed.State.Finishing != original.State.Finishing {
				t.Errorf("finishing touches did not survive replay:\n  was  %+v\n  now  %+v",
					original.State.Finishing, replayed.State.Finishing)
			}

			if len(replayed.Events) != len(original.Events) {
				t.Errorf("replay ran to %d events, the record holds %d",
					len(replayed.Events), len(original.Events))
			}
		})
	}
}
