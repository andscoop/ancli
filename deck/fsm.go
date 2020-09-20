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
	// ScoreAnswer updates score factors and tells user quizoutcome
	ScoreAnswer
	// // PassAnswer lets the user know they passed the quiz
	// PassAnswer
	// // FailAnswer lets the user know they failed the quiz
	// FailAnswer
)

const (
	// CmdNext execs a "next" transition
	CmdNext = "next"
	// CmdBack execs a "prev" transition
	CmdBack = "back"
	// CmdScore marks an card quiz answer as correct
	CmdScore = "score"
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

	tupple := CmdStateTupple{strings.TrimSpace(cmd), d.state}
	if f := StateTransitionTable[tupple]; f == nil {
		fmt.Println("unknown command, try again please")
	} else {
		f(d)
		d.ToScreen()
	}
}

// StateTransitionTable transition table
var StateTransitionTable = map[CmdStateTupple]TransitionFunc{
	// Transitions from Idle
	{CmdNext, Idle}: func(d *Deck) {
		d.state = DisplayQuestion
	},
	// Transitions from DisplayQuestion
	{CmdNext, DisplayQuestion}: func(d *Deck) {
		d.state = DisplayAnswer
	},
	{CmdBack, DisplayQuestion}: func(d *Deck) {
		d.LastCard()
		d.state = DisplayQuestion
	},
	// Transitions from DisplayAnswer
	{CmdNext, DisplayAnswer}: func(d *Deck) {
		d.NextCard()
		d.state = DisplayQuestion
	},
	{CmdBack, DisplayAnswer}: func(d *Deck) {
		d.state = DisplayQuestion
	},
	{CmdScore, DisplayAnswer}: func(d *Deck) {
		d.SubmitCardAnswer()
		d.state = ScoreAnswer
	},
	// Transitions from ScoreAnswer
	{CmdNext, ScoreAnswer}: func(d *Deck) {
		d.NextCard()
		d.state = DisplayQuestion
	},
	{CmdScore, ScoreAnswer}: func(d *Deck) {
		d.SubmitCardAnswer()
		d.state = ScoreAnswer
	},
	{CmdBack, ScoreAnswer}: func(d *Deck) {
		d.LastCard()
		d.state = DisplayQuestion
	},
}
