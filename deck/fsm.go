package deck

import (
	"fmt"
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
)

// Cmd type for readability of state machine
type Cmd uint8

const (
	// CmdNext execs a "next" transition
	CmdNext Cmd = iota
	// CmdBack execs a "prev" transition
	CmdBack
	// CmdScore marks an card quiz answer as correct
	CmdScore
	// CmdArchive will take a card out of quiz mode
	CmdArchive
)

// CmdStateTupple tupple for state-command combination
type CmdStateTupple struct {
	Cmd   Cmd
	State State
}

// TransitionFunc transition function
type TransitionFunc func(deck *Deck)

// Exec will attempt to transition the state machine
func (d *Deck) Exec(cmd Cmd) {
	// get function from transition table

	tupple := CmdStateTupple{cmd, d.state}
	if f := StateTransitionTable[tupple]; f == nil {
		fmt.Println("unknown command, try again please")
	} else {
		f(d)
		d.ToScreen()
	}
}

// ArchiveTranstitionFunc is a commonly repeated archive command
func ArchiveTranstitionFunc(d *Deck) {
	d.ArchiveCard()
	d.NextCard()
	d.state = DisplayQuestion
}

// StateTransitionTable transition table
var StateTransitionTable = map[CmdStateTupple]TransitionFunc{
	// Transitions from Idle
	{CmdNext, Idle}: func(d *Deck) {
		if d.ShouldRandom() {
			d.RandCard()
		}

		d.NextCard()
		d.state = DisplayQuestion
	},
	// Transitions from DisplayQuestion
	{CmdNext, DisplayQuestion}: func(d *Deck) {
		if d.ShouldRandom() {
			d.RandCard()
		}

		d.NextCard()
		d.state = DisplayAnswer
	},
	{CmdBack, DisplayQuestion}: func(d *Deck) {
		d.LastCard()
		d.state = DisplayQuestion
	},
	{CmdArchive, DisplayQuestion}: ArchiveTranstitionFunc,
	// Transitions from DisplayAnswer
	{CmdNext, DisplayAnswer}: func(d *Deck) {
		if d.ShouldRandom() {
			d.RandCard()
		}

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
	{CmdArchive, DisplayAnswer}: ArchiveTranstitionFunc,
	// Transitions from ScoreAnswer
	{CmdNext, ScoreAnswer}: func(d *Deck) {
		if d.ShouldRandom() {
			d.RandCard()
		}

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
	{CmdArchive, ScoreAnswer}: ArchiveTranstitionFunc,
}
