package prompt

import (
	"fmt"
	"io"
	"strings"
)

// change to interface (key(), display()), have Short and Long Answers?
// add default answer option
type Answer struct {
	Name string
	Key  rune
}

func Prompt(rw io.ReadWriter, text string, options []Answer) (Answer, error) {
	var optionDisplays []string
	for _, option := range options {
		optionDisplays = append(optionDisplays, option.display())
	}

	display := fmt.Sprintf("%s (%s)", text, strings.Join(optionDisplays, ", "))

	rw.Write([]byte(display))

	return Answer{}, nil
}

func (a Answer) display() string {
	if strings.ContainsRune(a.Name, a.Key) {
		return strings.Replace(a.Name, string(a.Key), emphasize(a.Key), 1)
	} else {
		return fmt.Sprintf("%s %s", emphasize(a.Key), a.Name)
	}
}

func emphasize(key rune) string {
	return "\033[1m" + string(key) + "\033[22m"
}
