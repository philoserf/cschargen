package chargen

import "fmt"

// Step 20: Finishing Touches (pp. 129-130).
//
// "Once you have completed the career process, your character is
// mechanically ready to enter play ... However, numbers alone do not create
// a compelling character."
//
// None of the four is rolled and none is mechanical. The engine records
// them and leaves them empty rather than inventing them, which is why this
// step reads nothing from the dice.

// finishingCite is where Step 20 begins.
const finishingCite = "pp. 129-130"

// Finishing is what Step 20 established.
type Finishing struct {
	// Name: "Every character requires a name ... Consider the character's
	// homeworld and background when selecting a name" (p. 129).
	Name string `json:"name,omitempty"`

	// Gender: "choose the character's gender identity. There is no
	// requirement that it match your own" (p. 130).
	Gender string `json:"gender,omitempty"`

	// Appearance: "height, build, posture, clothing style, hair, scars, and
	// other distinctive features ... Even a few details can make the
	// character memorable" (p. 130).
	Appearance string `json:"appearance,omitempty"`

	// Goals: "at least one personal goal ... A goal should extend beyond a
	// single adventure and provide direction for future decisions" (p. 130).
	Goals string `json:"goals,omitempty"`
}

// Empty reports whether nobody supplied anything, which is what an auto run
// with no flags produces.
func (f Finishing) Empty() bool {
	return f.Name == "" && f.Gender == "" && f.Appearance == "" && f.Goals == ""
}

// finishingTouches is Step 20. It takes what the inputs carry and, in an
// interactive run, asks for the rest.
func (g *Generator) finishingTouches() error {
	step := g.log.Step("Step 20: Finishing Touches", finishingCite)

	// The cite is the step's own: whatever career the character left last
	// set g.cite, and Step 20 is not in it.
	g.cite = finishingCite

	inputs := g.char.Provenance.Inputs

	finishing := Finishing{
		Name:       inputs.Name,
		Gender:     inputs.Gender,
		Appearance: inputs.Appearance,
		Goals:      inputs.Goals,
	}

	err := g.askFinishing(&finishing, step)
	if err != nil {
		return err
	}

	g.char.State.Finishing = finishing

	// The name is the one field the record already carried, because it was
	// an input from the start. Keeping the two in step is what stops a
	// sheet headed by one name and a record stamped with another.
	g.char.Provenance.Inputs.Name = finishing.Name

	if finishing.Empty() {
		g.consequence(ConsequenceFinishing, step,
			"no finishing touches were supplied; the engine invents none (p. 129)", "")

		return nil
	}

	g.consequence(ConsequenceFinishing, step, "finishing touches recorded", "")

	return nil
}

// finishingFields is the four, in the order p. 129 prints them, with the
// question each is asked by.
var finishingFields = [...]struct {
	label  string
	prompt string
	get    func(*Finishing) *string
}{
	{"name", "What is the character's name?", func(f *Finishing) *string { return &f.Name }},
	{"gender", "What is their gender identity?", func(f *Finishing) *string { return &f.Gender }},
	{"appearance", "What do they look like?", func(f *Finishing) *string { return &f.Appearance }},
	{"goals", "What do they want?", func(f *Finishing) *string { return &f.Goals }},
}

// askFinishing asks for whatever the inputs did not supply, where there is
// somebody to ask. An auto run leaves them empty: the engine records these
// four and never invents them.
func (g *Generator) askFinishing(finishing *Finishing, step int) error {
	asker, ok := g.decider.(Asker)
	if !ok {
		return nil
	}

	for nth, field := range finishingFields {
		held := field.get(finishing)
		if *held != "" {
			continue
		}

		answer, err := asker.Ask(Question{
			Point:  "finishing",
			Prompt: field.prompt,
			Cite:   finishingCite,
			Nth:    nth + 1,
			Of:     len(finishingFields),
		})
		if err != nil {
			return fmt.Errorf("asking for the %s: %w", field.label, err)
		}

		*held = answer
	}

	_ = step

	return nil
}
