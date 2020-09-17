package deck

import (
	"github.com/andscoop/ancli/config"
	tm "github.com/buger/goterm"
)

const (
	passOutput = "===========\nPASS\n"
	failOutput = "===========\nFAIL\n"
)

// ToScreen handles printing of deck quiz given current state
func (d *Deck) ToScreen() {
	var screen string
	state := d.State
	c := d.PullCard()

	cNext := config.GetString("cmdShortcuts.next")
	cBack := config.GetString("cmdShortcuts.back")
	cPass := config.GetString("cmdShortcuts.pass")
	cFail := config.GetString("cmdShortcuts.fail")

	switch state {
	case Idle:
		return
	case DisplayQuestion:
		screen = screen + c.Quiz.Question
	case DisplayAnswer:
		screen = screen + c.Quiz.Question
		screen = screen + "\n" + c.Quiz.Answer
	case PassAnswer:
		screen = screen + c.Quiz.Question
		screen = screen + "\n" + c.Quiz.Answer
		screen = screen + "\n" + passOutput
	case FailAnswer:
		screen = screen + c.Quiz.Question
		screen = screen + "\n" + c.Quiz.Answer
		screen = screen + "\n" + failOutput
	}

	tm.Clear()

	tm.MoveCursor(1, 1)

	tm.Println(screen)

	tm.Printf(
		"\n\nnext (%s)  back (%s)  pass(%s)  fail(%s)\n",
		cNext, cBack, cPass, cFail,
	)

	tm.Flush()

}
