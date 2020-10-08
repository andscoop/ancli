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

	q := Quiz{
		Type:     NoQuiz,
		question: b,
	}

	if l == 0 {
		q.Type = NullQuiz
		q.question = nil
		return q
	}

	// check for ---
	parts := bytes.SplitN(b, []byte{'-', '-', '-'}, 3)
	if len(parts) >= 2 && len(parts[1]) != 0 {
		q.question = parts[0]
		q.answer = parts[1]
		q.Type = Card
		return q
	}

	// check for inline ~~
	parts = bytes.SplitN(b, []byte{'~', '~'}, 4)

	if len(parts) >= 3 {
		blanks := getCopy(parts[1])
		for i := range blanks {
			blanks[i] = '_'
		}
		parts[1] = blanks

		q.answer = getCopy(b)
		q.Type = Inline
		q.question = bytes.Join(parts, nil)
	}

	return q
}
