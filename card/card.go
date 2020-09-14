package card

import (
	"bufio"
	"os"
	"strings"

	tm "github.com/buger/goterm"
)

const (
	passOutput = "===========\nPASS\n"
	failOutput = "===========\nFAIL\n"
)

// Card contains the necessary pieces of an Anki Card
type Card struct {
	Fp          string
	LastIndexed string
	LastQuizzed string
	Quiz        Quiz
}

// Quiz holds question/answer elems of a card
type Quiz struct {
	Question string
	Answer   string
	HasBlank bool
}

// PrintQ handles the printing of a Question
func (c *Card) PrintQ() {
	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Println(c.Quiz.Question)
	tm.Flush()
}

// PrintA handles the printing of an Answer
func (c *Card) PrintA() {
	tm.MoveCursorUp(1)
	tm.Println(c.Quiz.Answer)
	tm.Flush()
}

// QuizResult prints the output of a single card quiz
func (c *Card) QuizResult(pass bool, overwrite bool) {
	lineCount := 1
	output := failOutput
	if overwrite {
		lineCount = 3
	}

	if pass {
		output = passOutput
	}

	tm.MoveCursorUp(lineCount)
	tm.Println(output)
	tm.Flush()
}

// UpdateQuizElems updates the quiz elems for a card
func (c *Card) UpdateQuizElems() {
	q := extractQuizElem(c.Fp)
	c.Quiz = q
}

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

		// determine if card contains strike through
		b, e := strings.Index(t, "~"), strings.LastIndex(t, "~")

		if (b == -1) || (e == -1) { // no strike through
			q.HasBlank = false
		} else {
			q.HasBlank = true
			// todo add support for multi line answers
			q.Answer = t[b+1 : e]
			// replace strikethrough text with underscores
			q.Question = strings.Replace(t, t[b:e+1], strings.Repeat("_", e-b), 1)
			break
		}
		if strings.Trim(t, " ") == "---" {
			q.Question = strings.Join(scannedLines, "\n")

			// we know where the question and answer are
			// fast parse rest of card
			for scanner.Scan() {
				remainingLines = append(remainingLines, scanner.Text())
			}

			q.Answer = strings.Join(remainingLines, "\n")
			break
		}

		scannedLines = append(scannedLines, t)
	}

	// scrub q/a
	q.Answer = strings.Trim(q.Question, " \n")
	q.Question = strings.Trim(q.Question, " \n")

	return q
}
