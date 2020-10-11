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
		{"TestNthRep", &Card{Reptitions: 2000, LastAnswered: "1900-01-01T00:00:00.00+00:00"}, true},
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

func TestNextCard(t *testing.T) {
	cardKeys := []string{"a", "b", "c", "d", "e"}

	var tests = []struct {
		name          string
		deck          *Deck
		forward       bool
		wantIndex     int
		wantEndOfDeck bool
	}{
		{"TestSimpleNextOne", &Deck{index: 0, keys: cardKeys}, true, 1, false},
		{"TestSimpleLastOne", &Deck{index: 1, keys: cardKeys}, false, 0, false},
		{"TestNextCardAnswered", &Deck{index: 1, keys: cardKeys, quizzedKeys: map[int]bool{1: true}}, true, 2, false},
		{"TestLastCardAnswered", &Deck{index: 1, keys: cardKeys, quizzedKeys: map[int]bool{1: true}}, false, 0, false},
		{"TestNextTwoAnswered", &Deck{index: 1, keys: cardKeys, quizzedKeys: map[int]bool{1: true, 2: true}}, true, 3, false},
		{"TestLastTwoAnswered", &Deck{index: 2, keys: cardKeys, quizzedKeys: map[int]bool{1: true, 2: true}}, false, 0, false},
		{"TestSimpleEndOfDeckForward", &Deck{index: 2, keys: cardKeys, quizzedKeys: map[int]bool{0: true, 1: true, 2: true, 3: true, 4: true}}, true, 2, true},
		{"TestSimpleEndOfDeckBackward", &Deck{index: 1, keys: cardKeys, quizzedKeys: map[int]bool{0: true, 1: true, 2: true, 3: true, 4: true}}, false, 1, true},
		{"TestLoopToTopOfDeck", &Deck{index: 4, keys: cardKeys, quizzedKeys: map[int]bool{4: true}}, true, 0, false},
		{"TestLoopToBottomOfDeck", &Deck{index: 0, keys: cardKeys, quizzedKeys: map[int]bool{0: true}}, false, 4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haveEndOfDeck := tt.deck.NextCard(tt.forward)
			haveIndex := tt.deck.index

			if haveIndex != tt.wantIndex {
				t.Errorf("expected deck.index=%d but got deck.index=%d", tt.wantIndex, haveIndex)
			}

			if haveEndOfDeck != tt.wantEndOfDeck {
				t.Errorf("expected endOfDeck=%v but got endOfDeck=%v", tt.wantEndOfDeck, haveEndOfDeck)
			}
		})
	}

}
