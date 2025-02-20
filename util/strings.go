// utility functions for handling strings, text, runes, etc.
package util

import (
	"strings"
)

// WrapText wraps the provided string at WIDTH characters. optionally takes another int, used to determine the maximum
// number of lines. returns a slice of strings, each element a wrapped line. for words longer than width it just brutally
// cuts them off. no mercy.
func WrapText(str string, width int, maxlines ...int) (lines []string) {
	capped := false
	if len(maxlines) == 1 {
		lines = make([]string, 0, maxlines[0])
		capped = true
	} else {
		lines = make([]string, 0)
	}

	currentLine := ""

	for _, broken := range strings.Split(str, "/n") {
		for _, s := range strings.Split(broken, " ") {
			s = strings.TrimSpace(s) //get rid of nasty tabs and other weird whitespace.

			//super long word make-it-not-break hack.
			if len(s) > width {
				s = s[:width]
			}

			//add a line if current word won't fit
			if len(currentLine)+len(s) > width {
				currentLine = strings.TrimSpace(currentLine)
				lines = append(lines, currentLine)
				currentLine = ""

				//break if number of lines == height
				if capped && len(lines) == cap(lines) {
					break
				}
			}

			currentLine += s
			if len(currentLine) != width {
				currentLine += " "
			}
		}

		currentLine = strings.TrimSuffix(currentLine, " ")
		lines = append(lines, currentLine)
		currentLine = ""

		if capped && len(lines) == cap(lines) {
			break
		}
	}
	//append last line if needed after we're done looping through text
	if currentLine != "" {
		currentLine = strings.TrimSpace(currentLine)
		lines = append(lines, currentLine)
	}

	return
}

// Returns a string of lorem ipsum test text with the requested number of words.
func LoremIpsum(words int) string {
	l := strings.Split("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", " ")
	if words <= len(l) {
		return strings.Join(l[0:words], " ")
	}

	str := make([]string, 0, words)
	for len(str) <= words {
		str = append(str, l...)
	}

	return strings.Join(str[0:words], " ")
}
