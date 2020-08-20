package card

import (
	"bufio"
	"os"
	"strings"
)

type Card struct {
	Fp, Question, Answer string
	HasBlank             bool
}

func indexStrikethrough(s string) (int, int) {
	b := strings.Index(s, "~")
	e := strings.LastIndex(s, "~")
	return b, e
}

func ParseCard(fp string) (*Card, error) {
	prevLines := make([]string, 1)
	remainingLines := make([]string, 1)
	card := Card{Fp: fp}

	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()

		b, e := indexStrikethrough(t)
		if (b == -1) || (e == -1) {
			card.HasBlank = false
		} else {
			card.HasBlank = true
			card.Answer = t[b+1 : e]
			// replace strikethrough text with underscores
			card.Question = strings.Replace(t, t[b:e+1], strings.Repeat("_", e-b), 1)
			break
		}
		if strings.Trim(t, " ") == "---" {
			card.Question = strings.Join(prevLines, "\n")

			// we know where the question and answer are
			// fast parse rest of card
			for scanner.Scan() {
				remainingLines = append(remainingLines, scanner.Text())
			}

			card.Answer = strings.Join(remainingLines, "\n")
			break
		}

		prevLines = append(prevLines, t)
	}

	return &card, nil
}
