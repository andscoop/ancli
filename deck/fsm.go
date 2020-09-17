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
	// PassAnswer lets the user know they passed the quiz
	PassAnswer
	// FailAnswer lets the user know they failed the quiz
	FailAnswer
)

const (
	// CmdNext execs a "next" transition
	CmdNext = "next"
	// CmdBack execs a "prev" transition
	CmdBack = "back"
	// CmdPass marks an card quiz answer as correct
	CmdPass = "pass"
	// CmdFail marks an card quiz answer as incorrect
	CmdFail = "fail"
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
		d.ToScreen()
	}
}

// StateTransitionTable transition table
var StateTransitionTable = map[CmdStateTupple]TransitionFunc{
	// Idle state transitions
	{CmdNext, Idle}: func(d *Deck) {
		d.State = DisplayQuestion
	},
	// Question state transitions
	{CmdNext, DisplayQuestion}: func(d *Deck) {
		d.State = DisplayAnswer
	},
	{CmdBack, DisplayQuestion}: func(d *Deck) {
		d.LastCard()
		d.State = DisplayQuestion
	},
	// Answer state transitions
	{CmdNext, DisplayAnswer}: func(d *Deck) {
		d.NextCard()
		d.State = DisplayQuestion
	},
	{CmdBack, DisplayAnswer}: func(d *Deck) {
		d.State = DisplayQuestion
	},
	{CmdPass, DisplayAnswer}: func(d *Deck) {
		d.State = PassAnswer
	},
	{CmdFail, DisplayAnswer}: func(d *Deck) {
		d.State = FailAnswer
	},
	// Pass/Fail state transitions
	{CmdNext, PassAnswer}: func(d *Deck) {
		d.LastCard()
		d.State = DisplayQuestion
	},
	{CmdNext, FailAnswer}: func(d *Deck) {
		d.NextCard()
		d.State = DisplayQuestion
	},
	{CmdFail, PassAnswer}: func(d *Deck) {
		d.State = FailAnswer
	},
	{CmdPass, FailAnswer}: func(d *Deck) {
		d.State = PassAnswer
	},
	{CmdBack, PassAnswer}: func(d *Deck) {
		d.LastCard()
		d.State = DisplayQuestion
	},
	{CmdBack, FailAnswer}: func(d *Deck) {
		d.LastCard()
		d.State = DisplayQuestion
	},
}
