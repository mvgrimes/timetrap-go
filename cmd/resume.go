package cmd

import (
	"fmt"
	"os"

	"github.com/mvgrimes/timetrap-go/internal/parse"
	"github.com/mvgrimes/timetrap-go/internal/tt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var resumeCmd = &cobra.Command{
	Use:     "resume",
	Aliases: []string{"r", "res"},
	Short:   "Start the timer for the current time sheet for an entry. Defaults to the active entry.",
	Long:    "Start the timer for the current time sheet for an entry. Defaults to the active entry.",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetInt("id")
		atTimeStr, _ := cmd.Flags().GetString("at")
		runResume(id, atTimeStr, args)
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)

	resumeCmd.PersistentFlags().IntP("id", "i", 0, "Resume entry with id <id> instead of the last entry")
	resumeCmd.PersistentFlags().StringP("at", "a", "", "Use this time instead of now")
}

func runResume(id int, atTimeStr string, args []string) {
	if len(args) != 0 {
		fmt.Printf("usage: t resume")
		os.Exit(1)
	}

	atTime, err := parse.Time(atTimeStr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	t := tt.TimeTrap{}
	t.Connect(viper.GetString("database_file"))

	meta := t.GetMeta()
	entry := t.GetCurrentEntry()
	if entry.ID == 0 {
		fmt.Printf("No entry yet on this sheet yet. Started a new entry.")
		// this is basically the Ruby implementation, but it could be improved
	}

	entry = t.GetEntry(id)
	if entry.ID == 0 {
		fmt.Printf("Can't find entry")
		os.Exit(1)
	}

	fmt.Printf("Resuming \"%s\" from entry #%d\n", meta.CurrentSheet, entry.ID)

	// Already running
	if !entry.End.Valid {
		fmt.Println("Timetrap is already running")
		os.Exit(1)
	}

	entry, err = t.ClockIn(meta.CurrentSheet, atTime, entry.Note)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Checked into sheet \"%s\"\n", entry.Sheet)
}
