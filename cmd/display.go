package cmd

import (
	"fmt"
	"os"

	"github.com/mvgrimes/timetrap-go/internal/format"
	"github.com/mvgrimes/timetrap-go/internal/tt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var displayCmd = &cobra.Command{
	Use:     "display",
	Aliases: []string{"d"},
	Short:   "Display the current timesheet or a specific.",
	Long: `Display the current timesheet or a specific. Pass 'all' as SHEET
      to display all unarchived sheets or 'full' to display archived and
      unarchived sheets.`,
	Run: func(cmd *cobra.Command, args []string) {
		includeIds, _ := cmd.Flags().GetBool("ids")
		runDisplay(includeIds, args)
	},
}

func init() {
	rootCmd.AddCommand(displayCmd)

	displayCmd.PersistentFlags().BoolP("ids", "v", false, "Print database ids (for use with edit)")
	displayCmd.PersistentFlags().StringP("start", "s", "", "Include entries that start on this date or later")
	displayCmd.PersistentFlags().StringP("end", "e", "", "Include entries that start on this date or earlier")
	displayCmd.PersistentFlags().StringP("format", "f", "", `The output format.
Valid built-in formats are ical, csv, json, ids, factor, and text (default).
Documentation on defining custom formats can be found in the README included
in this`)
	displayCmd.PersistentFlags().StringP("grep", "g", "", "Include entries where the note matches this regexp.")
}

// TODO: add start/end filtering
// TODO: add output formatting
// TODO: add grep filtering
func runDisplay(includeIds bool, args []string) {
	if len(args) > 0 {
		fmt.Println("usage: t display")
		os.Exit(1)
	}

	t := tt.TimeTrap{}
	t.Connect(viper.GetString("database_file"))

	meta := t.GetMeta()
	entries := t.GetEntries(meta.CurrentSheet)

	format.DisplayEntries(entries, meta.CurrentSheet, includeIds)
}
