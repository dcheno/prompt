package prompt_test

import (
	"github.com/dcheno/prompt"
	"github.com/dcheno/scripter"
	"testing"
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

	answer, _ := prompt.Prompt(script.In(), script.Out(), "how are you?", options)

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

	answer, _ := prompt.Prompt(script.In(), script.Out(), "how are you?", options)

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

	prompt.Prompt(script.In(), script.Out(), "how are you?", options)

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

	prompt.Prompt(script.In(), script.Out(), "how are you?", options)

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

	prompt.Prompt(script.In(), script.Out(), "how are you?", options)

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

	prompt.Prompt(script.In(), script.Out(), "how are you?", options)

	script.AssertFinished()
}

func TestPromptPropagatesWriteError(t *testing.T) {

}

func TestPromptPropagatesReadError(t *testing.T) {

}

func TestPromptAcceptsCaseInsensitiveShortCode(t *testing.T) {

}

func TestPromptAcceptsCaseInsensitiveLongAnswer(t *testing.T) {

}
