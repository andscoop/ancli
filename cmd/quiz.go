package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andscoop/ancli/config"
	"github.com/andscoop/ancli/deck"
	"github.com/spf13/cobra"
)

var randOrderFlag bool

func init() {
	rootCmd.AddCommand(quizCmd)

	quizCmd.Flags().BoolVarP(&randOrderFlag, "random", "r", false, "shuffle card order in quiz")

}

var quizCmd = &cobra.Command{
	Use:   "quiz",
	Short: "Start an ancli quiz session",
	Long:  `Starts a new quiz session where you can job your memory`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmds := map[string]deck.CmdScoreTupple{
			config.GetString("cmdShortcuts.next"):    {Cmd: deck.CmdNext, Score: 0},
			config.GetString("cmdShortcuts.back"):    {Cmd: deck.CmdBack, Score: 0},
			config.GetString("cmdShortcuts.archive"): {Cmd: deck.CmdArchive, Score: 0},
			"y":                                      {Cmd: deck.CmdYes, Score: 0},
			"n":                                      {Cmd: deck.CmdNo, Score: 0},
			config.GetString("cmdShortcuts.pass"):    {Cmd: deck.CmdScore, Score: 1},
			config.GetString("cmdShortcuts.fail"):    {Cmd: deck.CmdScore, Score: 0},
			"0":                                      {Cmd: deck.CmdScore, Score: 0},
			"1":                                      {Cmd: deck.CmdScore, Score: 1},
			"2":                                      {Cmd: deck.CmdScore, Score: 2},
			"3":                                      {Cmd: deck.CmdScore, Score: 3},
			"4":                                      {Cmd: deck.CmdScore, Score: 4},
			"5":                                      {Cmd: deck.CmdScore, Score: 5},
		}

		d, err := deck.LoadDeck(args[0], randOrderFlag)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Hint: use `ancli decks` to see existing decks")
			os.Exit(1)
		}

		reader := bufio.NewReader(os.Stdin)

		// kick it out of idle state
		d.Exec(deck.CmdNext)

		for {
			rawInput, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalln(err)
			}
			scrubbedInput := strings.Trim(rawInput, " \n")

			v, ok := cmds[scrubbedInput]
			if !ok {
				d.Exec(deck.CmdUnknown)
			}

			d.LastScoreSubmitted = v.Score
			d.Exec(v.Cmd)
		}
	},
}
