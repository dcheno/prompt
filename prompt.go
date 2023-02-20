package prompt

import (
	"bufio"
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

func Prompt(in io.Reader, out io.Writer, text string, options []Answer) (Answer, error) {
	if len(options) == 0 {
		panic("Cannot create prompt without any options.")
	}

	defaultOption := options[0]

	var optionDisplays []string
	for _, option := range options {
		optionDisplays = append(optionDisplays, option.display())
	}

	_, err := fmt.Fprintf(out, "%s (%s) [%s]\n", text, strings.Join(optionDisplays, ", "), defaultOption.Name)

	if err != nil {
		return Answer{}, err
	}

	for {
		scanner := bufio.NewScanner(in)
		var reply string
		if scanner.Scan() {
			reply = strings.ToLower(scanner.Text())
		} else {
			return Answer{}, scanner.Err()
		}

		if reply == "" {
			return defaultOption, nil
		}

		for _, option := range options {
			if option.isMatch(reply) {
				return option, nil
			}
		}

		fmt.Fprintln(out, "Sorry, that didn't match any of the prompt options.")
	}
}

func (a Answer) display() string {
	if strings.ContainsRune(a.Name, a.Key) {
		return strings.Replace(a.Name, string(a.Key), emphasize(a.Key), 1)
	} else {
		return fmt.Sprintf("%s %s", emphasize(a.Key), a.Name)
	}
}

func (a Answer) isMatch(s string) bool {
	return strings.EqualFold(s, string(a.Key)) || strings.EqualFold(s, a.Name)
}

func emphasize(key rune) string {
	return "\033[1m" + string(key) + "\033[22m"
}
