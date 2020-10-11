package deck

import (
	"fmt"
	"os"
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
	// RequestRestart indicates a user has reached of deck,
	// notifies them, and checks if they want to restart
	RequestRestart
	// ExitProgram is a quitter
	ExitProgram
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
	// CmdUnknown is a catchup and should throw an error
	CmdUnknown
	// CmdYes signals an affirmative
	CmdYes
	//CmdNo signals a negative ghostrider
	CmdNo
)

// CmdStateTupple tupple for state-command combination
type CmdStateTupple struct {
	Cmd   Cmd
	State State
}

// TransitionFunc transition function
type TransitionFunc func(deck *Deck)

// StateTransitionTable transition table
var StateTransitionTable = map[CmdStateTupple]TransitionFunc{
	{CmdNext, Idle}:               cmdNextFromIdle,            // Transitions from Idle
	{CmdNext, DisplayQuestion}:    cmdNextFromDisplayQuestion, // Transitions from DisplayQuestion
	{CmdBack, DisplayQuestion}:    cmdBackFromDisplayQuestion,
	{CmdArchive, DisplayQuestion}: archiveTranstitionFunc,
	{CmdNext, DisplayAnswer}:      cmdNextFromDisplayAnswer, // Transitions from DisplayAnswer
	{CmdBack, DisplayAnswer}:      cmdBackFromDisplayQuestion,
	{CmdScore, DisplayAnswer}:     cmdScoreFromDisplayAnswer,
	{CmdArchive, DisplayAnswer}:   archiveTranstitionFunc,
	{CmdNext, ScoreAnswer}:        cmdNextFromScoreAnswer, // Transitions from ScoreAnswer
	{CmdScore, ScoreAnswer}:       cmdScoreFromScoreAnswer,
	{CmdBack, ScoreAnswer}:        cmdBackFromScoreAnswer,
	{CmdArchive, ScoreAnswer}:     archiveTranstitionFunc, // Transitions from RequestRestart
	{CmdYes, RequestRestart}:      cmdYesFromRequestRestart,
	{CmdNo, RequestRestart}:       cmdNoFromRequestRestart,
}

// Exec will attempt to transition the state machine
func (d *Deck) Exec(cmd Cmd) {
	// get function from transition table
	tupple := CmdStateTupple{cmd, d.state}
	if f := StateTransitionTable[tupple]; f == nil {
		fmt.Println("unknown command, try again please")
	} else {
		f(d)

		if d.state == ExitProgram {
			os.Exit(0)
		}

		d.toScreen()
	}
}

// ArchiveTranstitionFunc is a commonly repeated archive command
func archiveTranstitionFunc(d *Deck) {
	d.ArchiveCard()
	deckEmpty := d.NextCard(1)
	if deckEmpty {
		d.state = RequestRestart
	} else {
		d.state = DisplayQuestion
	}
}

func cmdNextFromIdle(d *Deck) {
	d.state = DisplayQuestion
}

func cmdNextFromDisplayQuestion(d *Deck) {
	d.state = DisplayAnswer
}

func cmdNextFromDisplayAnswer(d *Deck) {
	deckEmpty := d.NextCard(1)
	if deckEmpty {
		d.state = RequestRestart
	} else {
		d.state = DisplayQuestion
	}
}

func cmdNextFromScoreAnswer(d *Deck) {
	deckEmpty := d.NextCard(0)
	if deckEmpty {
		d.state = RequestRestart
	} else {
		d.state = DisplayQuestion
	}
}

func cmdBackFromDisplayQuestion(d *Deck) {
	deckEmpty := d.NextCard(-1)
	if deckEmpty {
		d.state = RequestRestart
	} else {
		d.state = DisplayQuestion
	}
}

func cmdBackFromDisplayAnswer(d *Deck) {
	d.state = DisplayQuestion
}

func cmdBackFromScoreAnswer(d *Deck) {
	d.state = DisplayAnswer

}

func cmdScoreFromDisplayAnswer(d *Deck) {
	d.SubmitCardAnswer()
	d.state = ScoreAnswer
}

func cmdScoreFromScoreAnswer(d *Deck) {
	d.SubmitCardAnswer()
	d.state = ScoreAnswer
}

func cmdYesFromRequestRestart(d *Deck) {
	d.resetQuizHistory()
	d.NextCard(1)
	d.state = DisplayQuestion
}

func cmdNoFromRequestRestart(d *Deck) {
	d.state = ExitProgram
}
