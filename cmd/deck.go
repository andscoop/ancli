package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/andscoop/ancli/config"
	"github.com/andscoop/ancli/deck"
	"github.com/spf13/cobra"
)

var rootPathFlag string
var includeHiddenFlag bool

func init() {
	rootCmd.AddCommand(decksCmd)
	decksCmd.AddCommand(decksCreateCmd)

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	decksCreateCmd.Flags().StringVarP(&rootPathFlag, "filepath", "f", pwd, "filepath to index for deck")
	decksCreateCmd.Flags().BoolVarP(&includeHiddenFlag, "include-hidden", "h", false, "maybe include hidden files")
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

		for _, d := range decks {
			// todo print deck stats
			fmt.Println(d.Name)
		}
	},
}

var decksCreateCmd = &cobra.Command{
	Use:   "create [name] [prefix] [filepath]",
	Short: "Create a new deck",
	Long:  `Create a new deck`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		quizAlgo := config.GetString("defaultAlgo")
		d := deck.Deck{
			Name:       args[0],
			DeckPrefix: args[1],
			RootDir:    rootPathFlag,
			QuizAlgo:   quizAlgo,
		}

		err := d.IndexAndSave(includeHiddenFlag)
		if err != nil {
			panic(err)
		}
	},
}

// var deckDeleteCmd = &cobra.Command{
// 	Use:   "deck delete",
// 	Short: "Refresh the index of your anki cards",
// 	Long:  `Refresh the index of your anki cards`,
// 	Args:  cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		err := deck.Walk(args[0], false)
// 		if err != nil {
// 			panic(err)
// 		}
// 	},
// }
