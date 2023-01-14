package prompt_test

import (
	"io"
	"testing"

	"github.com/dcheno/prompt"
)

type readWriter struct {
	readPos      int
	bytesToRead  []byte
	bytesWritten []byte
}

func (rw *readWriter) Read(p []byte) (int, error) {
	bytesRead := 0

	if rw.readPos == len(rw.bytesToRead) {
		return 0, io.EOF
	}

	for bytesRead < len(p) && rw.readPos < len(rw.bytesToRead) {
		p[bytesRead] = rw.bytesToRead[rw.readPos]
		bytesRead++
		rw.readPos++
	}
	return bytesRead, nil
}

func (rw *readWriter) Write(p []byte) (int, error) {
	rw.bytesWritten = append(rw.bytesWritten, p...)
	return len(p), nil
}

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

	rw := readWriter{}
	rw.bytesToRead = []byte("f\n")

	answer, _ := prompt.Prompt(&rw, "how are you?", options)

	expectedAnswer := prompt.Answer{"fine", 'f'}
	if answer != expectedAnswer {
		// t.Errorf("%v != %v", answer, expectedAnswer)
	}
}

func TestPromptReturnsValidAnswer_FullName(t *testing.T) {

}

func TestPromptEmphasizesFirstMatchingCharacter(t *testing.T) {
	options := []prompt.Answer{
		{
			"I don't know",
			'w',
		},
	}

	rw := readWriter{}
	rw.bytesToRead = []byte("w\n")

	prompt.Prompt(&rw, "how are you?", options)

	expectedWritten := []byte("how are you? (I don't kno\033[1mw\033[22m)")
	if string(rw.bytesWritten) != string(expectedWritten) {
		t.Errorf("%v != %v", string(rw.bytesWritten), string(expectedWritten))
	}

}

func TestPromptAddsLeadingCharacterIfNoMatching(t *testing.T) {
	options := []prompt.Answer{
		{
			"alright",
			'K',
		},
	}

	rw := readWriter{}
	rw.bytesToRead = []byte("K\n")

	prompt.Prompt(&rw, "how are you?", options)

	expectedWritten := []byte("how are you? (\033[1mK\033[22m alright)")
	if string(rw.bytesWritten) != string(expectedWritten) {
		t.Errorf("%v != %v", string(rw.bytesWritten), string(expectedWritten))
	}
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

	rw := readWriter{}
	rw.bytesToRead = []byte("g\n")

	prompt.Prompt(&rw, "how are you?", options)

	expectedWritten := []byte("how are you? (\033[1mg\033[22mood, \033[1mf\033[22mine, \033[1mo\033[22mtherwise, \033[1mK\033[22m alright)")
	if string(rw.bytesWritten) != string(expectedWritten) {
		t.Errorf("%v != %v", string(rw.bytesWritten), string(expectedWritten))
	}
}

func TestPromptRetriesOnBadAnswer(t *testing.T) {

}

func TestPromptPropagatesWriteError(t *testing.T) {

}

func TestPromptPropagatesReadError(t *testing.T) {

}
