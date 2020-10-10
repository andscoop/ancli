package cmd

import (
	"fmt"

	"github.com/andscoop/ancli/config"
	"github.com/andscoop/ancli/deck"
	tm "github.com/buger/goterm"
	"github.com/spf13/cobra"
)

var includeHiddenFlag bool
var deckAlgoFlag string

func init() {
	rootCmd.AddCommand(decksCmd)
	decksCmd.AddCommand(decksCreateCmd)

	quizAlgo := config.GetString("defaultAlgo")

	decksCreateCmd.Flags().StringVarP(&deckAlgoFlag, "algo", "a", quizAlgo, "reptition algo to use for the deck")
}

var decksCmd = &cobra.Command{
	Use:     "decks",
	Aliases: []string{"deck"},
	Short:   "List all decks and their configurations",
	Long:    `List all decks and their configurations`,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.GetConfig()

		decks := make(deck.Decks)

		err := c.UnmarshalKey("decks", &decks)
		if err != nil {
			panic(err)
		}

		deckStats := tm.NewTable(0, 10, 5, ' ', 0)
		fmt.Fprintf(deckStats, "Name\tLastUsed\tRootDir\n")
		fmt.Fprintf(deckStats, "----\t--------\t-------\n")

		for _, d := range decks {
			fmt.Fprintf(deckStats, "%s\t%s\t%s\n", d.Name, d.LastIndexed, d.RootDir)
		}

		tm.Println(deckStats)
		tm.Flush()
	},
}

var decksCreateCmd = &cobra.Command{
	Use:   "create [name] [regex] [filepath]",
	Short: "Create a new deck",
	Long:  `Create a new deck`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		d := deck.Deck{
			Name:      args[0],
			DeckRegex: args[1],
			RootDir:   args[2],
			QuizAlgo:  deckAlgoFlag,
		}

		err := d.IndexAndSave(false)
		if err != nil {
			panic(err)
		}
	},
}
