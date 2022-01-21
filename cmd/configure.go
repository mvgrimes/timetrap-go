package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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

	// TODO: print the location of the config file
	// TODO: does this do anything else?
	fmt.Println("configure command is not yet implemented")
	os.Exit(1)
}
