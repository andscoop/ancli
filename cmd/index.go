package cmd

import (
	"github.com/spf13/cobra"

	"github.com/andscoop/ancli/indexer"
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
		indexer.Walk(args[0], false)
	},
}
