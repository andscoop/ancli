package card

import (
	"bufio"
	"os"
	"strings"
)

func extractQuizElem(fp string) Quiz {
	scannedLines := make([]string, 1)
	remainingLines := make([]string, 1)
	q := Quiz{}

	file, err := os.Open(fp)
	if err != nil {
		return q
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// todo clean this up. FSM?
	for scanner.Scan() {
		t := scanner.Text()

		b, e := indexStrikethrough(t)
		if (b == -1) || (e == -1) {
			q.HasBlank = false
		} else {
			q.HasBlank = true
			q.Answer = scrub(t[b+1 : e])
			// replace strikethrough text with underscores
			q.Question = scrub(strings.Replace(t, t[b:e+1], strings.Repeat("_", e-b), 1))
			break
		}
		if strings.Trim(t, " ") == "---" {
			q.Question = scrub(strings.Join(scannedLines, "\n"))

			// we know where the question and answer are
			// fast parse rest of card
			for scanner.Scan() {
				remainingLines = append(remainingLines, scanner.Text())
			}

			q.Answer = scrub(strings.Join(remainingLines, "\n"))
			break
		}

		scannedLines = append(scannedLines, t)
	}

	return q
}

func indexStrikethrough(s string) (int, int) {
	b := strings.Index(s, "~")
	e := strings.LastIndex(s, "~")
	return b, e
}

func scrub(a string) string {
	return strings.Trim(a, " \n")
}
