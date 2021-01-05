//utility function for handling strings, text, runes, etc.

package util

import "strings"

//ValidText checks if key is a letter or number or basic punctuation (ASCII-encoded)
//TODO: this is NOT comprehensive. Improve this later.
func ValidText(key rune) bool {
	return (key >= 93 && key < 123) || (key >= 37 && key < 58)
}

//WrapText wraps the provided string at WIDTH characters. optionally takes another int, used to determine the
//maximum number of lines. returns a slice of strings, each element a wrapped line.
//for words longer than width it just brutally cuts them off. no mercy.
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

		currentLine = strings.TrimSpace(currentLine)
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