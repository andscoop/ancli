package cmd

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/andscoop/ancli/config"
	"github.com/andscoop/ancli/deck"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(quizCmd)
}

var quizCmd = &cobra.Command{
	Use:   "quiz",
	Short: "Start an ancli quiz session",
	Long:  `Starts a new quiz session where you can job your memory`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cNext := config.GetString("cmdShortcuts.next")
		cBack := config.GetString("cmdShortcuts.back")
		cPass := config.GetString("cmdShortcuts.pass")
		cFail := config.GetString("cmdShortcuts.fail")
		cArchive := config.GetString("cmdShortcuts.archive")

		d := deck.LoadDeck(args[0])

		reader := bufio.NewReader(os.Stdin)

		// kick it out of idle state
		d.Exec(deck.CmdNext)

		for {
			var fsmCmd deck.Cmd

			// read command from stdin
			rawInput, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalln(err)
			}

			scrubbedInput := strings.Trim(rawInput, " \n")

			switch scrubbedInput {
			case cNext:
				fsmCmd = deck.CmdNext
			case cBack:
				fsmCmd = deck.CmdBack
			case cPass:
				fsmCmd = deck.CmdScore
				d.LastScoreSubmitted = 1
			case cFail:
				fsmCmd = deck.CmdScore
				d.LastScoreSubmitted = 0
			case cArchive:
				fsmCmd = deck.CmdArchive
			default:
				// attempt to convert to int
				value, err := strconv.ParseInt(scrubbedInput, 0, 64)
				if err == nil {
					fsmCmd = deck.CmdScore
					d.LastScoreSubmitted = value
				}
			}

			d.Exec(fsmCmd)
		}
	},
}
