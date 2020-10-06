package quiz

import (
	"fmt"
	"testing"
)

func testEq(a, b []byte) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		fmt.Printf("a: %v", (a == nil))
		fmt.Printf("b: %v", (b == nil))
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestParseQuizTypes(t *testing.T) {
	var tests = []struct {
		b    []byte
		name string
		want QuizType
	}{
		// Basic Type Checks
		{[]byte{}, "TypeCheck basic NullQuiz", NullQuiz},
		{[]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g'}, "TypeCheck basic NoQuiz", NoQuiz},
		{[]byte{'a', 'b', 'c', '-', '-', '-', 'd', 'e', 'f', 'g'}, "TypeCheck basic Card", Card},
		{[]byte{'a', 'b', '~', '~', 'c', 'd', '~', '~', 'e', 'f', 'g'}, "TypeCheck basic Inline", Inline},
		{[]byte{'a', 'b', '~', '~', 'c', 'd'}, "TypeCheck NoQuiz with opening tildes only", NoQuiz},
		{[]byte{'a', 'b', 'c', '-', '-', '-'}, "TypeCheck NoQuiz with no bytes following dashes", NoQuiz},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := Parse(tt.b)
			if out.Type != tt.want {
				t.Errorf("got %d, want %d", out.Type, tt.want)
			}
		})
	}
}

func TestParseQuestionsAndAnswers(t *testing.T) {
	var tests = []struct {
		b            []byte
		name         string
		wantQuestion []byte
		wantAnswer   []byte
	}{
		// Q/A Checks
		{[]byte{}, "Q/A NullQuiz", nil, nil},
		{[]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g'}, "Q/A NoQuiz", []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g'}, nil},
		{[]byte{'a', 'b', 'c', '-', '-', '-', 'd', 'e', 'f', 'g'}, "Q/A Card", []byte{'a', 'b', 'c'}, []byte{'d', 'e', 'f', 'g'}},
		{[]byte{'a', 'b', '~', '~', 'c', 'd', '~', '~', 'e', 'f', 'g'}, "Q/A Inline", []byte{'a', 'b', '_', '_', 'e', 'f', 'g'}, []byte{'a', 'b', '~', '~', 'c', 'd', '~', '~', 'e', 'f', 'g'}},
		{[]byte{'a', 'b', '~', '~', 'c', 'd'}, "Q/A Missing closing tildes", []byte{'a', 'b', '~', '~', 'c', 'd'}, nil},
		{[]byte{'a', 'b', 'c', '-', '-', '-'}, "Q/A Card Missing answer", []byte{'a', 'b', 'c', '-', '-', '-'}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			have := Parse(tt.b)

			if !testEq(have.question, tt.wantQuestion) {
				t.Errorf("got question %v, want %v", have.question, tt.wantQuestion)
			}

			if !testEq(have.answer, tt.wantAnswer) {
				t.Errorf("got answer %v, want %v", have.answer, tt.wantAnswer)
			}
		})
	}
}
