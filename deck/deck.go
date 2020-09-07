package deck

import (
	"fmt"

	"github.com/andscoop/ancli/card"
	"github.com/andscoop/ancli/config"
)

// Deck is a deck of cards to be quizzed
type Deck struct {
	Cards     []*card.Card
	Index     int
	UserInput rune
	State     State
}

// NewDeck creates a new Deck and loads cards
func NewDeck() *Deck {
	index, err := config.GetIndex()
	if err != nil {
		panic(err)
	}

	// todo Card and Index should eventually be merged
	// that will avoid needing to loop the index when loading a deck
	// https://github.com/spf13/viper#unmarshaling unmardhalling directly to struct should help
	// but needs tested
	var cards []*card.Card
	for fp := range index {
		c, err := card.ParseCard(fp)
		if err != nil {
			fmt.Println("Error Parsing ", fp)
			fmt.Println("Check if file exists")
		}

		cards = append(cards, c)
	}

	return &Deck{State: Idle, Index: 0, Cards: cards}
}

func (d *Deck) safeGet(i int) card.Card {
	d.Index = d.Index + i

	// no more cards on top of deck, go to end
	if d.Index < 0 {
		d.Index = len(d.Cards) - 1
	}

	// end of deck, go to beginning
	if d.Index >= len(d.Cards) {
		d.Index = 0
	}

	return *d.Cards[d.Index]
}

// GetCard is a func
func (d *Deck) getCard() card.Card {
	return d.safeGet(0)
}

// GetNextCard is a func
func (d *Deck) getNextCard() card.Card {
	return d.safeGet(1)
}

func (d *Deck) getPrevCard() card.Card {
	return d.safeGet(-1)
}
