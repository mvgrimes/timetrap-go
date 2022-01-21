package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configureCmd = &cobra.Command{
	Use:     "configure",
	Aliases: []string{"c", "config", "conf"},
	Short:   "Write out a YAML config file. Print path to config file.  The file may contain ERB.",
	Long:    "Write out a YAML config file. Print path to config file.  The file may contain ERB.",
	Run: func(cmd *cobra.Command, args []string) {
		runConfigure(args)
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

func runConfigure(args []string) {
	if len(args) != 0 {
		fmt.Printf("usage: t configure")
		os.Exit(1)
	}

	fmt.Printf("Config file is at \"%s\"\n", viper.GetViper().ConfigFileUsed())
	// TODO: does configuration command do anything else?

	os.Exit(0)
}
