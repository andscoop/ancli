package deck

import (
	"testing"
)

func TestSM2Algo(t *testing.T) {
	var tests = []struct {
		name           string
		card           *Card
		wantShouldQuiz bool
	}{
		{"AlwaysTestFirstReptition", &Card{Reptitions: 0, LastAnswered: "1900-01-01T00:00:00.00+00:00"}, true},
		{"Test2ndRep", &Card{Reptitions: 1, LastAnswered: "1900-01-01T00:00:00.00+00:00"}, true},
		{"Test3rdRep", &Card{Reptitions: 2, LastAnswered: "1900-01-01T00:00:00.00+00:00"}, true},
		{"TestNthRep", &Card{Reptitions: 1000, LastAnswered: "1900-01-01T00:00:00.00+00:00"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haveShouldQuiz := shouldQuizSM2(tt.card)

			if haveShouldQuiz != tt.wantShouldQuiz {
				t.Errorf("expected shouldQuiz=%v but got shouldQuiz=%v", tt.wantShouldQuiz, haveShouldQuiz)
			}
		})
	}
}
