package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"a", "arch"},
	Short:   "Alter an entry's note, start, or end time. Defaults to the active entry. Defaults to the last entry to be checked out of if no entry is active.",
	Long:    "Alter an entry's note, start, or end time. Defaults to the active entry. Defaults to the last entry to be checked out of if no entry is active.",
	Run: func(cmd *cobra.Command, args []string) {
		runEdit(args)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.PersistentFlags().Int32P("id", "i", 0, "lter entry with id <id> instead of the running entry")
	editCmd.PersistentFlags().StringP("start", "s", "", "Change the start time to <time>")
	editCmd.PersistentFlags().StringP("end", "e", "", "Change the end time to <time>")
	editCmd.PersistentFlags().BoolP("append", "z", false, "Append to the current note instead of replacing it the delimiter between appends notes is configurable (see config)")
	editCmd.PersistentFlags().StringP("move", "m", "", "Move to another sheet")
}

// TODO: implment edit command
func runEdit(args []string) {
	if len(args) != 0 {
		fmt.Printf("usage: t edit ...")
		os.Exit(1)
	}

	fmt.Println("edit command is not yet implemented")
	os.Exit(1)
}
