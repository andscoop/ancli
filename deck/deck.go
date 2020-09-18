package deck

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/andscoop/ancli/config"
)

// Deck is a deck of cards to be quizzed
type Deck struct {
	// Easy-fetching of any card by fp key
	Cards map[string]Card
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
	var cs = make(map[string]Card)
	c := config.GetConfig()

	err := c.UnmarshalKey("decks", &cs)
	if err != nil {
		panic(err)
	}

	d := Deck{State: Idle, Index: 0, Cards: cs}

	d.UpdateKeys()

	return &d
}

// SubmitAnswer updates the card to have an answer
func (d *Deck) SubmitAnswer(score int) {
	c := d.PullCard()
	n := time.Now().String()
	c.LastQuizzed = n

	useEF := config.GetBool("useSM2")
	if useEF {
		// todo pass for now
	}

	if score > 0 {
		c.LastPassed = n
	}

}

// PullCard pulls the current card of the deck
func (d *Deck) PullCard() *Card {
	fmt.Println(d.Index)
	fmt.Println(len(d.Cards))
	fp := d.Keys[d.Index]
	c, _ := d.Cards[fp]

	c.UpdateQuizElems()

	return &c
}

// NextCard shifts deck index up for later pulling
func (d *Deck) NextCard() {
	d.Index = d.Index + 1

	// take first card
	// if it doesn't meet criteria
	// sort cards by cards ready to be quizzed

	// end of deck, go to beginning
	if d.Index >= len(d.Cards) {
		d.Index = 0
	}
}

// LastCard shifts deck index down for later pulling
func (d *Deck) LastCard() {
	d.Index = d.Index - 1

	// no more cards on top of deck, go to end
	if d.Index < 0 {
		d.Index = len(d.Cards) - 1
	}
}

// RandCard will return a random card from the deck. TODO
func (d *Deck) RandCard() {
	v := rand.Intn(len(d.Cards) - 1)

	d.Index = v
}
