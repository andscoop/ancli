package cmd

import (
	"github.com/andscoop/ancli/deck"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(indexCmd)
}

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Refresh the index of your anki cards",
	Long:  `Refresh the index of your anki cards`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := deck.Walk(args[0], false)
		if err != nil {
			panic(err)
		}
	},
}
