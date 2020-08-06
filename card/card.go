package card

import (
	"fmt"
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

func ParseCard(fp string, s []byte) *Card {
	prevLines := make([]string, 1)

	raw := string(s)
	card := Card{Fp: fp}
	lines := strings.Split(raw, "\n")

	for i, line := range lines {
		b, e := indexStrikethrough(raw)
		if (b == -1) || (e == -1) {
			card.HasBlank = false
		} else {
			card.HasBlank = true
			card.Answer = line[b+1 : e+1]
			card.Question = strings.Replace(raw, raw[b:e], strings.Repeat("_", e-b), 1)
			break
		}

		if strings.Trim(line, " ") == "---" {
			// todo check to make sure '---' is not first line
			fmt.Println("found line break")
			card.Question = strings.Join(prevLines, "\n")
			card.Answer = strings.Join(lines[i+1:], "\n")
			break
		}

		prevLines = append(prevLines, line)
	}

	return &card
}
