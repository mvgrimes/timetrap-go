package cmd

import (
	"fmt"
	"os"

	"github.com/mvgrimes/timetrap-go/internal/tt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Start the timer for the current timesheet.",
	Long:  `Start the timer for the current timesheet.`,
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

	fmt.Printf(" %-18s%-12s%-12s%s\n", "Timesheet", "Running", "Today", "Total Time")
	summaries := t.List()
	for _, summary := range summaries {
		active := " "
		if summary.LastActive {
			active = "-"
		}
		if summary.Active {
			active = "*"
		}
		fmt.Printf("%s%-18s%-12s%-12s%s\n", active, summary.Sheet, summary.Running, summary.Today, summary.Total)
	}
}
