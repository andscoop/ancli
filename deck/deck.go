package deck

import (
	"github.com/andscoop/ancli/card"
)

// Deck is a deck of cards to be quizzed
type Deck struct {
	// Easy-fetching of any card by fp key
	Cards map[string]*card.Card
	// Track order of deck and allow for easy sorting
	Keys []string
	// Track place in deck
	Index     int
	UserInput rune
	State     State
}

// UpdateKeys because they can fall out of sync
func (d *Deck) UpdateKeys() error {
	for k := range d.Cards {
		d.Keys = append(d.Keys, k)
	}

	return nil
}

// NewDeck creates a new Deck and loads cards
func NewDeck() *Deck {
	// index, err := config.GetIndex()
	// if err != nil {
	// 	panic(err)
	// }

	// todo Card and Index should eventually be merged
	// that will avoid needing to loop the index when loading a deck
	// https://github.com/spf13/viper#unmarshaling unmardhalling directly to struct should help
	// but needs tested
	// var cards map[string]*card.Card
	// for fp := range index {
	// 	c, err := card.ParseCard(fp)
	// 	if err != nil {
	// 		fmt.Println("Error Parsing ", fp)
	// 		fmt.Println("Check if file exists")
	// 	}

	// 	cards[fp] = c
	// }
	return &Deck{State: Idle, Index: 0, Cards: make(map[string]*card.Card)}
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

	return *d.Cards[d.Keys[d.Index]]
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
