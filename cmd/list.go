package cmd

import (
	"fmt"
	"os"

	"github.com/mvgrimes/timetrap-go/internal/format"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "Show the available timesheets.",
	Long:    `Show the available timesheets. If run as 'ls' do not include the archived sheets.`,
	Run: func(cmd *cobra.Command, args []string) {
		includeArchived := cmd.CalledAs() != "ls"
		runList(includeArchived, args)
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
