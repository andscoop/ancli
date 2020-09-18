package deck

// Card contains the necessary pieces of an Anki Card
type Card struct {
	Fp          string
	LastIndexed string
	LastQuizzed string
	LastPassed  string
	EasyFactor  float32
	Quiz        Quiz
}

// Quiz holds question/answer elems of a card
type Quiz struct {
	Question string
	Answer   string
	HasBlank bool
}

// UpdateQuizElems updates the quiz elems for a card
func (c *Card) UpdateQuizElems() {
	q := extractQuizElem(c.Fp)
	c.Quiz = q
}
