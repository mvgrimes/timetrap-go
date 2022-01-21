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

func runList(args []string) {
	if len(args) > 0 {
		fmt.Println("usage: t list")
		os.Exit(1)
	}

	t := tt.TimeTrap{}
	t.Connect(viper.GetString("database_file"))

	fmt.Printf(" %-15s  %-10s  %-10s  %s\n",
		"Timesheet", "Running", "Today", "Total Time")
	summaries := t.List()
	for _, summary := range summaries {
		active := " "
		if summary.LastActive {
			active = "-"
		}
		if summary.Active {
			active = "*"
		}
		fmt.Printf(
			"%s%-15s% 10s  % 10s%12s\n",
			active,
			summary.Sheet,
			tt.FmtDuration(summary.Running),
			tt.FmtDuration(summary.Today),
			tt.FmtDuration(summary.Total),
		)
	}
}
