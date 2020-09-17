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

// PullCard pulls the current card of the deck
func (d *Deck) PullCard() card.Card {
	fp := d.Keys[d.Index]
	c, _ := d.Cards[fp]

	c.UpdateQuizElems()

	return c
}

// NextCard shifts deck index up for later pulling
func (d *Deck) NextCard() {
	// end of deck, go to beginning
	if d.Index >= len(d.Cards) {
		d.Index = 0
	}

	d.Index = d.Index + 1
}

// LastCard shifts deck index down for later pulling
func (d *Deck) LastCard() {
	// no more cards on top of deck, go to end
	if d.Index < 0 {
		d.Index = len(d.Cards) - 1
	}

	d.Index = d.Index - 1
}

// todo implement
// func (d *Deck) RandCard {}
