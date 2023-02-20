package prompt_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dcheno/prompt"
	"github.com/dcheno/scripter"
)

func TestPromptReturnsValidAnswer_ShortCode(t *testing.T) {
	options := []prompt.Answer{
		{
			"good",
			'g',
		},
		{
			"fine",
			'f',
		},
	}

	script := scripter.NewScript(
		t,
		scripter.Expect("how are you? (\033[1mg\033[22mood, \033[1mf\033[22mine)\n"),
		scripter.Reply("f\n"),
	)

	prompter := prompt.Prompter{script.In(), script.Out()}
	answer, _ := prompter.Prompt("how are you?", options)

	expectedAnswer := prompt.Answer{"fine", 'f'}
	if answer != expectedAnswer {
		t.Errorf("%v != %v", answer, expectedAnswer)
	}
}

func TestPromptReturnsValidAnswer_FullName(t *testing.T) {
	options := []prompt.Answer{
		{
			"good",
			'g',
		},
		{
			"fine",
			'f',
		},
	}

	script := scripter.NewScript(
		t,
		scripter.Expect("how are you? (\033[1mg\033[22mood, \033[1mf\033[22mine)\n"),
		scripter.Reply("fine\n"),
	)

	prompter := prompt.Prompter{script.In(), script.Out()}
	answer, _ := prompter.Prompt("how are you?", options)

	expectedAnswer := prompt.Answer{"fine", 'f'}
	if answer != expectedAnswer {
		t.Errorf("%v != %v", answer, expectedAnswer)
	}
}

func TestPromptEmphasizesFirstMatchingCharacter(t *testing.T) {
	options := []prompt.Answer{
		{
			"I don't know",
			'w',
		},
	}

	script := scripter.NewScript(
		t,
		scripter.Expect("how are you? (I don't kno\033[1mw\033[22m)\n"),
		scripter.Reply("w\n"),
	)

	prompter := prompt.Prompter{script.In(), script.Out()}
	prompter.Prompt("how are you?", options)

	script.AssertFinished()
}

func TestPromptAddsLeadingCharacterIfNoMatching(t *testing.T) {
	options := []prompt.Answer{
		{
			"alright",
			'K',
		},
	}

	script := scripter.NewScript(
		t,
		scripter.Expect("how are you? (\033[1mK\033[22m alright)\n"),
		scripter.Reply("k\n"),
	)

	prompter := prompt.Prompter{script.In(), script.Out()}
	prompter.Prompt("how are you?", options)

	script.AssertFinished()
}

func TestPromoptWritesAllPromptOptions(t *testing.T) {
	options := []prompt.Answer{
		{
			"good",
			'g',
		},
		{
			"fine",
			'f',
		},
		{
			"otherwise",
			'o',
		},
		{
			"alright",
			'K',
		},
	}

	script := scripter.NewScript(
		t,
		scripter.Expect("how are you? (\033[1mg\033[22mood, \033[1mf\033[22mine, \033[1mo\033[22mtherwise, \033[1mK\033[22m alright)\n"),
		scripter.Reply("fine\n"),
	)

	prompter := prompt.Prompter{script.In(), script.Out()}
	prompter.Prompt("how are you?", options)

	script.AssertFinished()
}

func TestPromptRetriesOnBadAnswer(t *testing.T) {
	options := []prompt.Answer{
		{
			"good",
			'g',
		},
		{
			"fine",
			'f',
		},
	}

	script := scripter.NewScript(
		t,
		scripter.Expect("how are you? (\033[1mg\033[22mood, \033[1mf\033[22mine)\n"),
		scripter.Reply("not great\n"),
		scripter.Expect("Sorry, that didn't match any of the prompt options.\n"),
		scripter.Reply("f\n"),
	)

	prompter := prompt.Prompter{script.In(), script.Out()}
	prompter.Prompt("how are you?", options)

	script.AssertFinished()
}

type errorWriter struct{}

func (ew errorWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("An error!")
}

func TestPromptPropagatesWriteError(t *testing.T) {
	r := strings.NewReader("some value")
	w := errorWriter{}

	options := []prompt.Answer{
		{
			"hopefully it errors",
			'h',
		},
	}

	prompter := prompt.Prompter{r, w}
	answer, err := prompter.Prompt("What happens if the write fails?", options)

	blankAnswer := prompt.Answer{}
	if answer != blankAnswer {
		t.Error("Should return blank answer.")
	}

	if err.Error() != "An error!" {
		t.Error("Should have forwarded error from writer!")
	}
}

type errorReader struct{}

func (er errorReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("An error!")
}

type fakeWriter struct{}

func (fw fakeWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func TestPromptPropagatesReadError(t *testing.T) {
	r := errorReader{}
	w := fakeWriter{}

	options := []prompt.Answer{
		{
			"hopefully it errors",
			'h',
		},
	}

	prompter := prompt.Prompter{r, w}
	answer, err := prompter.Prompt("What happens if the read fails?", options)

	blankAnswer := prompt.Answer{}
	if answer != blankAnswer {
		t.Error("Should return blank answer.")
	}

	if err.Error() != "An error!" {
		t.Error("Should have forwarded error from writer!")
	}
}

func TestPromptAcceptsCaseInsensitiveShortCode(t *testing.T) {
	expectedAnswer := prompt.Answer{"yes", 'y'}
	options := []prompt.Answer{
		expectedAnswer,
		{
			"no",
			'n',
		},
	}

	script := scripter.NewScript(
		t,
		scripter.Expect("Are you enjoying 'prompt'? (\033[1my\033[22mes, \033[1mn\033[22mo)\n"),
		scripter.Reply("Y\n"),
	)

	prompter := prompt.Prompter{script.In(), script.Out()}
	answer, err := prompter.Prompt("Are you enjoying 'prompt'?", options)

	if err != nil {
		t.Error("Unexpected error!")
	}

	if answer != expectedAnswer {
		t.Error("Did not return expected answer.")
	}

	script.AssertFinished()
}

func TestPromptAcceptsCaseInsensitiveLongAnswer(t *testing.T) {

	expectedAnswer := prompt.Answer{"yes", 'y'}
	options := []prompt.Answer{
		expectedAnswer,
		{
			"no",
			'n',
		},
	}

	script := scripter.NewScript(
		t,
		scripter.Expect("Are you enjoying 'prompt'? (\033[1my\033[22mes, \033[1mn\033[22mo)\n"),
		scripter.Reply("YeS\n"),
	)

	prompter := prompt.Prompter{script.In(), script.Out()}
	answer, err := prompter.Prompt("Are you enjoying 'prompt'?", options)

	if err != nil {
		t.Error("Unexpected error!")
	}

	if answer != expectedAnswer {
		t.Error("Did not return expected answer.")
	}

	script.AssertFinished()
}

func TestPromptRetriesOnEmptyAnswer(t *testing.T) {
	options := []prompt.Answer{
		{
			"good",
			'g',
		},
		{
			"fine",
			'f',
		},
	}

	script := scripter.NewScript(
		t,
		scripter.Expect("how are you? (\033[1mg\033[22mood, \033[1mf\033[22mine)\n"),
		scripter.Reply("\n"),
		scripter.Expect("Sorry, that didn't match any of the prompt options.\n"),
		scripter.Reply("f\n"),
	)

	prompter := prompt.Prompter{script.In(), script.Out()}
	prompter.Prompt("how are you?", options)

	script.AssertFinished()
}

func TestPromptWithDefaultUsesDefaultWhenBlank(t *testing.T) {
	expectedAnswer := prompt.Answer{"fine", 'f'}
	options := []prompt.Answer{
		{
			"good",
			'g',
		},
		expectedAnswer,
	}

	script := scripter.NewScript(
		t,
		scripter.Expect("how are you? (\033[1mg\033[22mood, \033[1mf\033[22mine) [fine]\n"),
		scripter.Reply("\n"),
	)

	prompter := prompt.Prompter{script.In(), script.Out()}
	answer, err := prompter.PromptWithDefault("how are you?", options, &expectedAnswer)

	if err != nil {
		t.Error("Unexpected error!")
	}

	if answer != expectedAnswer {
		t.Error("Got the wrong answer", answer)
	}

	script.AssertFinished()
}

func TestPromptWithDefaultStillTakesOtherAnswers(t *testing.T) {
	expectedAnswer := prompt.Answer{"fine", 'f'}
	defaultAnswer := prompt.Answer{"good", 'g'}
	options := []prompt.Answer{
		defaultAnswer,
		expectedAnswer,
	}

	script := scripter.NewScript(
		t,
		scripter.Expect("how are you? (\033[1mg\033[22mood, \033[1mf\033[22mine) [good]\n"),
		scripter.Reply("f\n"),
	)

	prompter := prompt.Prompter{script.In(), script.Out()}
	answer, err := prompter.PromptWithDefault("how are you?", options, &defaultAnswer)

	if err != nil {
		t.Error("Unexpected error!")
	}

	if answer != expectedAnswer {
		t.Error("Got the wrong answer", answer)
	}

	script.AssertFinished()
}

func TestPromptWithDefaultStillAcceptsExplicitDefaultAnswer(t *testing.T) {
	expectedAnswer := prompt.Answer{"fine", 'f'}
	otherAnswer := prompt.Answer{"good", 'g'}
	options := []prompt.Answer{
		otherAnswer,
		expectedAnswer,
	}

	script := scripter.NewScript(
		t,
		scripter.Expect("how are you? (\033[1mg\033[22mood, \033[1mf\033[22mine) [fine]\n"),
		scripter.Reply("f\n"),
	)

	prompter := prompt.Prompter{script.In(), script.Out()}
	answer, err := prompter.PromptWithDefault("how are you?", options, &expectedAnswer)

	if err != nil {
		t.Error("Unexpected error!")
	}

	if answer != expectedAnswer {
		t.Error("Got the wrong answer", answer)
	}

	script.AssertFinished()
}
