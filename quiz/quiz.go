package quiz

import (
	"bytes"
)

// Quiz holds elements for quiz cards
type Quiz struct {
	question []byte
	answer   []byte
	Type     QuizType
}

// Question will fetch ready to ask Question string
func (q *Quiz) Question() string {
	return string(q.question)
}

// Answer to the questions you seek
func (q *Quiz) Answer() string {
	return string(q.answer)
}

// QuizType for readability
type QuizType uint8

const (
	// Unknown quiz types should not happen. Handle accordingly
	Unknown QuizType = iota
	// Inline quizes use ~~strikethrough~~ syntax
	Inline
	// Card quizzes use --- line break style sytnax to separate cards and answers
	Card
	// NoQuiz will return the whole note and not have an answer
	NoQuiz
	// NullQuiz has no text to share
	NullQuiz
)

func getCopy(s []byte) []byte {
	c := make([]byte, len(s))
	copy(c, s)

	return c
}

// Parse will take a string of bytes, looking for quiz card delimeters
// and return a quiz ready to be used
func Parse(b []byte) Quiz {
	l := len(b)
	qbs := getCopy(b)
	abs := getCopy(b)

	q := Quiz{
		Type:     NoQuiz,
		question: qbs,
	}

	if l == 0 {
		q.Type = NullQuiz
		q.question = nil
		return q
	}

	// check for ---
	if i := bytes.Index(b, []byte{'-', '-', '-'}); i >= 0 {
		if i+3 < l {
			q.question = qbs[0:i]
			q.answer = abs[i+3 : l]
			q.Type = Card
		}

		return q
	}

	// check for inline ~~
	if s := bytes.Index(b, []byte{'~', '~'}); s >= 0 {
		// found opening ~~, try to close it out
		if e := bytes.Index(b[s+2:l], []byte{'~', '~'}); e >= 0 {
			q.answer = abs
			q.Type = Inline

			blanks := bytes.Repeat([]byte{'_'}, (s+2)-e)

			// put it back together
			qbs2 := getCopy(qbs[0:s])
			qbs2 = append(qbs2, blanks...)
			// +4 accounts for removing opening and closing tildes
			qbs2 = append(qbs2, qbs[s+e+4:l]...)

			q.question = qbs2

			return q
		}
	}

	return q
}
