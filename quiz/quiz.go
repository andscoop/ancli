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

// Parse will take a string of bytes, looking for quiz card delimeters
// and return a quiz ready to be used
func Parse(b []byte) Quiz {
	q := Quiz{
		Type:     NoQuiz,
		question: b,
	}

	if len(b) == 0 {
		q.Type = NullQuiz
	}

	// check for ---
	if i := bytes.Index(b, []byte{'-', '-', '-'}); i >= 0 {
		if i+3 < len(b) {
			q.question = b[0:i]
			q.answer = b[i+3 : len(b)]
			q.Type = Card
		}

		return q
	}

	// check for inline ~~
	if s := bytes.Index(b, []byte{'~', '~'}); s >= 0 {
		// found opening ~~, try to close it out
		if e := bytes.Index(b[s+2:len(b)], []byte{'~', '~'}); e >= 0 {
			q.answer = b
			q.Type = Inline

			for s < e+2 {
				b[s] = '_'
				s++
			}

			q.question = b

		}
	}

	return q
}
