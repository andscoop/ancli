package cmd

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/andscoop/ancli/config"
	"github.com/andscoop/ancli/deck"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "quiz",
	Short: "Start an ancli quiz session",
	Long:  `Starts a new quiz session where you can job your memory`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {

		cNext := config.GetString("cmdShortcuts.next")
		cBack := config.GetString("cmdShortcuts.back")
		cPass := config.GetString("cmdShortcuts.pass")
		cFail := config.GetString("cmdShortcuts.fail")

		cmds := map[string]string{cNext: deck.CmdNext, cBack: deck.CmdBack, cPass: deck.CmdPass, cFail: deck.CmdFail}

		d := deck.NewDeck()
		reader := bufio.NewReader(os.Stdin)

		d.Exec(deck.CmdNext)

		for {
			// read command from stdin
			c, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalln(err)
			}

			cmd, ok := cmds[strings.Trim(c, " \n")]
			if !ok {
				cmd = "unknown"
			}
			d.Exec(cmd)
		}
	},
}
