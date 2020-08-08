package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"

	"github.com/andscoop/ancli/card"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "test",
	Short: "Start an ancli test session",
	Long:  `Starts a new test session where you can job your memory`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Files: ", args)
		for _, file := range args {
			info, err := os.Stat(file)
			if err != nil {
				panic(err)
			}

			fmt.Println(info.ModTime())

			dat, err := ioutil.ReadFile(file)
			if err != nil {
				panic(err)
			}

			c := card.ParseCard(file, dat)

			fmt.Println(c.Question)
			fmt.Println(c.Answer)
		}
	},
}
