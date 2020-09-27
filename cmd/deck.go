package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/andscoop/ancli/config"
	"github.com/andscoop/ancli/deck"
	tm "github.com/buger/goterm"
	"github.com/spf13/cobra"
)

var rootPathFlag string
var includeHiddenFlag bool
var deckAlgoFlag string

func init() {
	rootCmd.AddCommand(decksCmd)
	decksCmd.AddCommand(decksCreateCmd)

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	quizAlgo := config.GetString("defaultAlgo")

	decksCreateCmd.Flags().StringVarP(&rootPathFlag, "filepath", "f", pwd, "filepath to index for deck")
	// decksCreateCmd.Flags().BoolVarP(&includeHiddenFlag, "include-hidden", "h", false, "maybe include hidden files")
	decksCreateCmd.Flags().StringVarP(&deckAlgoFlag, "algo", "a", quizAlgo, "reptition algo to use for the deck")
}

var decksCmd = &cobra.Command{
	Use: "decks",
	// todo alias to deck
	Short: "List all decks and their configurations",
	Long:  `List all decks and their configurations`,
	// Args:  cobra.ExactArgs(1),
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
	Use:   "create [name] [prefix] [filepath]",
	Short: "Create a new deck",
	Long:  `Create a new deck`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		d := deck.Deck{
			Name:       args[0],
			DeckPrefix: args[1],
			RootDir:    rootPathFlag,
			QuizAlgo:   deckAlgoFlag,
		}

		err := d.IndexAndSave(false)
		if err != nil {
			panic(err)
		}
	},
}
