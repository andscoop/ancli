package cmd

import (
	"fmt"
	"os"

	"github.com/andscoop/ancli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ancli",
		Short: "anCLI is a command line interface for repetition based learning",
		Long:  "A CLI for repeition based learning",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("welcome to ancli")
		},
	}
)

// Execute the rootcmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.Init)

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))

}
