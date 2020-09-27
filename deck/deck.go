package deck

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/andscoop/ancli/config"
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
	LastScoreSubmitted int64
	QuizAlgo           string
	DeckPrefix         string
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
	quiz         Quiz
}

// Quiz holds question/answer elems of a card
// It is loaded JIT by parser
type Quiz struct {
	Question string
	Answer   string
	HasBlank bool
}

// LoadDeck TODO
func LoadDeck(name string) *Deck {
	quizAlgo := config.GetString("defaultAlgo")

	var d = Deck{
		Name:     name,
		QuizAlgo: quizAlgo,
	}

	c := config.GetConfig()

	err := c.UnmarshalKey("decks."+name, &d)
	if err != nil {
		panic(err)
	}

	d.syncQuizzableCards()

	return &d
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
		// todo double check conversion and cleanup
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

// syncQuizzableCards adds quizzable cards to Deck.Keys
// a card is quizzable if it ShouldQuiz()
func (d *Deck) syncQuizzableCards() error {
	shouldQuiz := shouldQuizFuncs[d.QuizAlgo]

	for k, c := range d.Cards {
		if shouldQuiz(c) {
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

	config.SetAndSave("decks", d.Cards)

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

// NextCard shifts deck index up for later pulling
func (d *Deck) NextCard() {
	d.index = d.index + 1

	// end of deck, go to beginning
	if d.index >= len(d.Cards) {
		d.index = 0
	}
}

// LastCard shifts deck index down for later pulling
func (d *Deck) LastCard() {
	d.index = d.index - 1

	// no more cards on top of deck, go to end
	if d.index < 0 {
		d.index = len(d.Cards) - 1
	}
}

// RandCard will return a random card from the deck. TODO
func (d *Deck) RandCard() {
	v := rand.Intn(len(d.Cards) - 1)

	d.index = v
}

func (c *Card) parseQuiz() error {
	scannedLines := make([]string, 1)
	remainingLines := make([]string, 1)
	q := Quiz{}

	file, err := os.Open(c.Fp)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// todo clean this up. FSM?
	for scanner.Scan() {
		t := scanner.Text()

		b, e := indexStrikethrough(t)
		if (b == -1) || (e == -1) {
			q.HasBlank = false
		} else {
			q.HasBlank = true
			q.Answer = scrub(t[b+1 : e])
			// replace strikethrough text with underscores
			q.Question = scrub(strings.Replace(t, t[b:e+1], strings.Repeat("_", e-b), 1))
			break
		}
		if strings.Trim(t, " ") == "---" {
			q.Question = scrub(strings.Join(scannedLines, "\n"))

			// we know where the question and answer are
			// fast parse rest of card
			for scanner.Scan() {
				remainingLines = append(remainingLines, scanner.Text())
			}

			q.Answer = scrub(strings.Join(remainingLines, "\n"))
			break
		}

		scannedLines = append(scannedLines, t)
	}

	c.quiz = q

	return nil
}

func indexStrikethrough(s string) (int, int) {
	b := strings.Index(s, "~")
	e := strings.LastIndex(s, "~")
	return b, e
}

func scrub(a string) string {
	return strings.Trim(a, " \n")
}

func hashFp(fp string) string {
	data := []byte(fp)
	return fmt.Sprintf("%x", md5.Sum(data))
}

// ToScreen handles printing of deck quiz given current state
func (d *Deck) ToScreen() error {
	c, err := d.PullCard()
	if err != nil {
		return err
	}
	screen := c.quiz.Question

	cNext := config.GetString("cmdShortcuts.next")
	cBack := config.GetString("cmdShortcuts.back")
	cPass := config.GetString("cmdShortcuts.pass")
	cFail := config.GetString("cmdShortcuts.fail")

	switch d.state {
	case DisplayAnswer:
		screen = screen + "\n" + c.quiz.Answer
	case ScoreAnswer:
		screen = screen + "\n" + c.quiz.Answer
		if d.LastScoreSubmitted == 0 {
			screen = screen + "\n" + failOutput
		} else {
			screen = screen + "\n" + passOutput
		}
	}

	tm.Clear()
	tm.MoveCursor(1, 1)

	screen = screen + "\n" + string(d.LastScoreSubmitted)
	tm.Println(screen)

	tm.Printf(
		"\n\nnext (%s)  back (%s)  pass(%s)  fail(%s)\n",
		cNext, cBack, cPass, cFail,
	)

	tm.Flush()

	return nil
}
