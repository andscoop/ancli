package card

import (
	"bufio"
	"os"
	"strings"

	tm "github.com/buger/goterm"
)

// Card contains the necessary pieces of an Anki Card
type Card struct {
	Fp, Question, Answer string
	HasBlank             bool
}

// PrintQ handles the printing of a Question
func (c *Card) PrintQ() {
	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Println(c.Question)
	tm.Flush()
}

// PrintA handles the printing of an Answer
func (c *Card) PrintA() {
	tm.MoveCursorUp(1)
	tm.Println(c.Answer)
	tm.Flush()
}

func scrub(a string) string {
	return strings.Trim(a, " \n")
}

// ParseCard will break out an Anki card into its necessary parts
func ParseCard(fp string) (*Card, error) {
	scannedLines := make([]string, 1)
	remainingLines := make([]string, 1)
	card := Card{Fp: fp}

	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// todo clean this up. FSM?
	for scanner.Scan() {
		t := scanner.Text()

		b, e := indexStrikethrough(t)
		if (b == -1) || (e == -1) {
			card.HasBlank = false
		} else {
			card.HasBlank = true
			card.Answer = scrub(t[b+1 : e])
			// replace strikethrough text with underscores
			card.Question = scrub(strings.Replace(t, t[b:e+1], strings.Repeat("_", e-b), 1))
			break
		}
		if strings.Trim(t, " ") == "---" {
			card.Question = scrub(strings.Join(scannedLines, "\n"))

			// we know where the question and answer are
			// fast parse rest of card
			for scanner.Scan() {
				remainingLines = append(remainingLines, scanner.Text())
			}

			card.Answer = scrub(strings.Join(remainingLines, "\n"))
			break
		}

		scannedLines = append(scannedLines, t)
	}

	return &card, nil
}

func indexStrikethrough(s string) (int, int) {
	b := strings.Index(s, "~")
	e := strings.LastIndex(s, "~")
	return b, e
}
