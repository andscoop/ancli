package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/andscoop/ancli/card"
	"github.com/andscoop/ancli/config"
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
		// todo make decision about whether or not to re-index
		// todo reindex
		index, err := config.GetIndex()
		if err != nil {
			panic(err)
		}

		// todo load cards from index
		for fp, _ := range index {
			c, err := card.ParseCard(fp)
			if err != nil {
				fmt.Println("Error Parsing ", fp)
				fmt.Println("Check if file exists")
				continue
			}
			fmt.Println(c.Question)
			fmt.Println(c.Answer)
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter to continue!")
			_, _ = reader.ReadString('\n')
		}
	},
}
