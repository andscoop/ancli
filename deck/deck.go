package deck

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/andscoop/ancli/config"
	"github.com/andscoop/ancli/quiz"
	tm "github.com/buger/goterm"
)

const (
	passOutput = "===========\nPASS\n"
	failOutput = "===========\nFAIL\n"
)

// Deck holds one to many Cards and provides an interface for
// interacting with those Cards.
//
// Private fields are ignored by mapstructure lib when saving
// decks to config files.
//
// In order to simplify iterating through a deck and sorting
// Deck uses both a map (Cards) and a slice (Keys).
type Deck struct {
	Cards              map[string]*Card // fetching of any card by fp key
	state              State
	keys               []string
	index              int
	shouldQuiz         shouldQuizFunc
	quizzedKeys        map[int]bool
	LastScoreSubmitted int64
	QuizAlgo           string
	DeckRegex          string
	LastIndexed        string
	Name               string
	RootDir            string
}

// Decks on decks
type Decks map[string]*Deck

// Card holds all data relating to the actual quizzing.
//
// Interfacing with a card should happen through the Deck since the
// Deck tracks state.
// Private fields are ignored by mapstructure lib when saving
// decks to config files.
type Card struct {
	Fp           string
	LastIndexed  string
	LastAnswered string
	Reptitions   int
	EasyFactor   float64 // sm2 only
	FibIteration int     // fib only
	quiz         quiz.Quiz
	IsArchived   bool
}

// LoadDeck will load a deck from a saved config
func LoadDeck(name string, shouldShuffle bool) (*Deck, error) {
	quizAlgo := config.GetString("defaultAlgo")

	if !DeckIsSet(name) {
		return nil, fmt.Errorf("config: deck not found in saved config, %s", name)
	}

	var d = Deck{
		Name:        name,
		QuizAlgo:    quizAlgo,
		index:       0,
		quizzedKeys: make(map[int]bool),
	}

	c := config.GetConfig()

	err := c.UnmarshalKey("decks."+name, &d)
	if err != nil {
		panic(err)
	}

	d.syncQuizzableCards()
	if shouldShuffle {
		rand.Shuffle(len(d.keys), func(i, j int) {
			d.keys[i], d.keys[j] = d.keys[j], d.keys[i]
		})
	}

	return &d, nil
}

// DeckIsSet checks if deck exist in config
func DeckIsSet(name string) bool {
	c := config.GetConfig()
	return c.IsSet("decks." + name)
}

// Save saves current deck back to index
func (d *Deck) Save() {
	config.SetAndSave("decks."+d.Name, d)
}

// shouldQuizFuncs are responsible for determining if a card is due
// to be quizzed in accordance with the quiz algorithm in the config or
// passed through via flags
//
// a full list of supported algos can be found in shouldQuizFuncs map
type shouldQuizFunc func(c *Card) bool

// shouldQuizFuncs simplifies fetching of the shouldQuizFunc for a deck
var shouldQuizFuncs = map[string]shouldQuizFunc{
	"simple": shouldQuizSimple,
	"sm2":    shouldQuizSM2,
	"fib":    shouldQuizFib,
	"never":  shouldQuizNever,
}

// shouldQuizSimple always returns true
func shouldQuizSimple(c *Card) bool {
	return true
}

// shouldQuizSM2 returns true according to the SM2 algo
// https://www.supermemo.com/en/archives1990-2015/english/ol/sm2
func shouldQuizSM2(c *Card) bool {
	var since time.Duration
	reps := c.Reptitions
	ef := c.getEF()

	lq, err := time.Parse(time.RFC3339, c.LastAnswered)
	if err != nil {
		// todo I eat errors
		// currently assumes bad or null timestamp string and resets
		//
	}

	since = time.Since(lq)

	switch reps {
	case 0:
		return true
	case 1:
		if (since / 24) >= time.Duration(24)*time.Hour {
			return true
		}
		return false
	case 2:
		if (since / 24) >= time.Duration(6*42)*time.Hour {
			return true
		}
		return false
	default:
		calc := float64(reps-1) * ef
		expectedIntervalHours := math.Ceil(calc)
		expectedIntervalDays := (time.Duration(expectedIntervalHours) * time.Hour) / 24

		if (since / 24 * time.Hour) >= expectedIntervalDays {
			return true
		}
		return false
	}
}

// shouldQuizFib is based on the Fibonacci Sequence
// When a quiz is answered successfully, it increases
// the time between quiz repititions in accordance with the sequence.
// When a quiz is answered incorrectly, it decrements the sequence by one.
// todo
func shouldQuizFib(c *Card) bool {
	return true
}

// shouldQuizNever is for procrastinators
func shouldQuizNever(c *Card) bool {
	return false
}

// getEF fetches current EF of card, falling back to config default
func (c *Card) getEF() float64 {
	ef := c.EasyFactor
	if ef == 0 {
		ef = config.GetFloat("defaultEasyFactor")
	}

	return ef
}

func (c *Card) open() ([]byte, error) {
	file, err := os.Open(c.Fp)
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (c *Card) parseQuiz() error {
	b, err := c.open()
	if err != nil {
		return err
	}

	q := quiz.Parse(b)

	c.quiz = q

	return nil
}

// getQuizAlgo fetches set quizAlgo, falls back to global default
func (d *Deck) getQuizAlgo() string {
	dqa := d.QuizAlgo

	if len(dqa) == 0 {
		dqa = config.GetString("defaultAlgo")
	}

	return dqa
}

// syncQuizzableCards adds quizzable cards to Deck.Keys
// a card is quizzable if it ShouldQuiz()
func (d *Deck) syncQuizzableCards() error {
	qa := d.getQuizAlgo()

	shouldQuiz := shouldQuizFuncs[qa]

	for k, c := range d.Cards {
		if shouldQuiz(c) && !c.IsArchived {
			d.keys = append(d.keys, k)
		}
	}
	return nil
}

// SubmitCardAnswer will update any algo specific components for a card
func (d *Deck) SubmitCardAnswer() error {
	score := d.LastScoreSubmitted
	c, err := d.PullCard()
	if err != nil {
		return err
	}

	// common attrs that need updated regardless of algo
	c.LastAnswered = time.Now().Format(time.RFC3339)
	c.Reptitions++

	// handle algo specific fields
	switch d.QuizAlgo {
	case "sm2":
		minEF := config.GetFloat("minEasyFactor")
		ef := c.EasyFactor
		ef = ef + (.1 - (float64(5)-float64(score))*(.08+(float64(5)-float64(score))*.02))
		if ef < minEF {
			ef = minEF
		}

		c.EasyFactor = ef
	case "fib":
		if score == 0 {
			c.FibIteration--
		} else { // assume any score above 0 is a simple pass
			c.FibIteration++
		}
	}

	// track cards that are already quizzed for a session
	d.quizzedKeys[d.index] = true

	d.Save()

	return nil
}

// PullCard pulls the current card of the deck
func (d *Deck) PullCard() (*Card, error) {
	k := d.keys[d.index]
	c, _ := d.Cards[k]

	err := c.parseQuiz()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (d *Deck) quizzableCardsRemaining() int {
	c := 0
	for i := range d.keys {
		// found quizzable card
		quizzed, ok := d.quizzedKeys[i]
		if !ok || !quizzed {
			c++
		}
	}
	return c
}

// NextCard shifts the deck index forward or backwards
// returns a bool if there are no more quizzable cards in deck
func (d *Deck) NextCard(seek int) bool {
	var offsetI int
	directionShift := 1
	if seek < 0 {
		directionShift = -1
	}

	offset := d.index + seek

	// loop at most all keys, starting from root,
	// checking each key if it has been quizzed
	for i := range d.keys {
		offsetI = offset + (i * directionShift)
		// reached end of deck
		if offsetI > len(d.keys)-1 {
			// ensures next loop that offset = 0
			offset = -1 - i
			continue
		}

		// reached begginning of deck
		if offsetI < 0 {
			offset = len(d.keys) + i
			continue
		}

		// found quizzable card
		quizzed, ok := d.quizzedKeys[offsetI]
		if !ok || !quizzed {
			d.index = offsetI
			return false
		}
	}

	return true
}

// ArchiveCard will mark a card as archived so it won't be quizzed
func (d *Deck) ArchiveCard() {
	c, _ := d.PullCard()
	c.IsArchived = true

	d.Save()
}

func hashFp(fp string) string {
	data := []byte(fp)
	return fmt.Sprintf("%x", md5.Sum(data))
}

// ToScreen handles printing of deck quiz given current state
func (d *Deck) toScreen() {
	if d.state == RequestRestart {
		d.printRestartRequestScreen()
		return
	}

	d.printCardScreen()
}

func (d *Deck) resetQuizHistory() {
	d.quizzedKeys = make(map[int]bool)
}

func (d *Deck) printRestartRequestScreen() {
	screen := "End of deck reached, would you like to restart?\n\nyes (y) no (n)\n>"

	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Println(screen)
	tm.Flush()
}

func (d *Deck) printCardScreen() {
	c, _ := d.PullCard()

	screen := c.quiz.Question()

	cNext := config.GetString("cmdShortcuts.next")
	cBack := config.GetString("cmdShortcuts.back")
	cPass := config.GetString("cmdShortcuts.pass")
	cFail := config.GetString("cmdShortcuts.fail")
	cArchive := config.GetString("cmdShortcuts.archive")

	answer := c.quiz.Answer()

	switch d.state {
	case DisplayAnswer:
		if c.quiz.Type == quiz.Inline {
			screen = answer
		} else {
			screen = screen + "\n\n" + answer
		}
	case ScoreAnswer:
		if c.quiz.Type == quiz.Inline {
			screen = answer
		} else {
			screen = screen + "\n\n" + answer
		}

		if d.LastScoreSubmitted == 0 {
			screen = screen + "\n\n" + failOutput
		} else {
			screen = screen + "\n\n" + passOutput
		}
	}

	tm.Clear()
	tm.MoveCursor(1, 1)

	screen = screen + "\n" + string(d.LastScoreSubmitted)
	tm.Print(screen)
	tm.MoveCursor(1, 30)
	screen = "======================="
	screen = screen + fmt.Sprintf("\nnext (%s)  back (%s)  pass (%s)  fail (%s)  archive (%s)\n",
		cNext, cBack, cPass, cFail, cArchive)
	screen = screen + "======================="
	screen = screen + fmt.Sprintf("\nPath: %s\n", c.Fp)
	screen = screen + fmt.Sprintf("\nCards Remaining: %d", d.quizzableCardsRemaining())

	tm.Print(screen)
	tm.Print("\n> ")
	tm.Flush()
}
