package cmd

import (
	"fmt"
	"os"
	"time"

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
		runDisplay(args)
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

	viper.BindPFlag("ids", displayCmd.PersistentFlags().Lookup("ids"))
	viper.BindPFlag("start", displayCmd.PersistentFlags().Lookup("start"))
	viper.BindPFlag("end", displayCmd.PersistentFlags().Lookup("end"))
	viper.BindPFlag("format", displayCmd.PersistentFlags().Lookup("format"))
	viper.BindPFlag("grep", displayCmd.PersistentFlags().Lookup("grep"))
}

func runDisplay(args []string) {
	if len(args) > 0 {
		fmt.Println("usage: t display")
		os.Exit(1)
	}

	t := tt.TimeTrap{}
	t.Connect(viper.GetString("database_file"))

	meta := t.GetMeta()
	fmt.Printf("Timesheet: %s\n", meta.CurrentSheet)
	fmt.Printf("    %-18s %-8s   %-8s   %8s   %s\n", "Day", "Start", "End", "Duration", "Notes")
	entries := t.Display()
	lastDay := ""
	var total time.Duration
	for _, entry := range entries {
		day := ""
		if lastDay != entry.Day {
			day = entry.Day
		}
		fmt.Printf(
			"    %-18s %-8s - %-8s   %8s   %s\n",
			day,
			entry.Start.Format("15:04:05"),
			entry.End.Format("15:04:05"),
			entry.Duration,
			entry.Note,
		)
		lastDay = day
		total += entry.Duration
		// TODO: add daily total
	}
	fmt.Printf("    -------------------------------------------------\n")
	fmt.Printf("    Total %43s\n", total)
}
