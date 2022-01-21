package cmd

import (
	"fmt"
	"os"

	"github.com/mvgrimes/timetrap-go/internal/tt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "Show the available timesheets.",
	Long:    `Show the available timesheets.`,
	Run: func(cmd *cobra.Command, args []string) {
		runList(args)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(includeArchived bool, args []string) {
	if len(args) > 0 {
		fmt.Println("usage: t list")
		os.Exit(1)
	}

	format.DisplayList(includeArchived)
}
