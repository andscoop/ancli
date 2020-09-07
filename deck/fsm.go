package deck

import (
	"fmt"
	"strings"
)

// State type for readability of state machine
type State uint8

const (
	// Idle is init state for FSM
	Idle State = iota
	// DisplayQuestion until user commands otherwise
	DisplayQuestion
	// DisplayAnswer until user commands otherwise
	DisplayAnswer
)

const (
	// CmdNext execs a "next" transition
	CmdNext = "next"
	// CmdBack execs a "prev" transition
	CmdBack = "back"
	// // CmdCorrect execs a "correct" transition"
	// CmdCorrect = "correct"
	// //CmdIncorrect execs a "incorrect" transition
	// CmdIncorrect = "incorrect"
)

// CmdStateTupple tupple for state-command combination
type CmdStateTupple struct {
	Cmd   string
	State State
}

// TransitionFunc transition function
type TransitionFunc func(deck *Deck)

// Exec will attempt to transition the state machine
func (d *Deck) Exec(cmd string) {
	// get function from transition table
	tupple := CmdStateTupple{strings.TrimSpace(cmd), d.State}
	if f := StateTransitionTable[tupple]; f == nil {
		fmt.Println("unknown command, try again please")
	} else {
		f(d)
	}
}

// StateTransitionTable transition table
var StateTransitionTable = map[CmdStateTupple]TransitionFunc{
	{CmdNext, Idle}: func(d *Deck) {
		c := d.getCard()
		c.PrintQ()
		d.State = DisplayQuestion
	},
	{CmdNext, DisplayQuestion}: func(d *Deck) {
		c := d.getCard()
		c.PrintA()
		d.State = DisplayAnswer
	},
	{CmdNext, DisplayAnswer}: func(d *Deck) {
		c := d.getNextCard()
		c.PrintQ()
		d.State = DisplayQuestion
	},
	{CmdBack, DisplayQuestion}: func(d *Deck) {
		c := d.getPrevCard()
		c.PrintQ()
		d.State = DisplayQuestion
	},
	{CmdBack, DisplayAnswer}: func(d *Deck) {
		c := d.getCard()
		c.PrintQ()
		d.State = DisplayQuestion
	},
}
