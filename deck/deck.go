package deck

import (
	"github.com/andscoop/ancli/card"
	"github.com/andscoop/ancli/config"
)

// Deck is a deck of cards to be quizzed
type Deck struct {
	// Easy-fetching of any card by fp key
	Cards map[string]card.Card
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
	cs, err := config.GetSavedCards()
	if err != nil {
		panic(err)
	}

	d := Deck{State: Idle, Index: 0, Cards: cs}

	d.UpdateKeys()

	return &d
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

	fp := d.Keys[d.Index]
	c, _ := d.Cards[fp]

	c.UpdateQuizElems()

	return c
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
