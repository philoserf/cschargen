package chargen

import "testing"

// TestAsPromptOnlyLiftsTheFirstLetter. Four choice points take their prompt
// from the result that raised them, and a result is written as the page
// writes it. Beside the prompts the engine composes itself, which are
// sentences, a lowercase one reads as a fragment of the book rather than as
// a question being asked.
func TestAsPromptOnlyLiftsTheFirstLetter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		detail string
		want   string
	}{
		{"a result the page writes in lower case", "what surviving taught you", "What surviving taught you"},
		{"and one naming skills", "choose Broker, Carouse or Streetwise", "Choose Broker, Carouse or Streetwise"},
		{"a sentence the engine wrote is left alone", "Choose a specialty", "Choose a specialty"},
		{"and so is a name", "STR or END", "STR or END"},
		{"nothing to lift", "", ""},
		// The rest of the string is untouched, so a skill or a proper noun
		// inside it keeps its own spelling.
		{"only the first letter moves", "gain a level in Gun Combat", "Gain a level in Gun Combat"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := asPrompt(test.detail); got != test.want {
				t.Errorf("asPrompt(%q) = %q, want %q", test.detail, got, test.want)
			}
		})
	}
}
