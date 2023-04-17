package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mvgrimes/timetrap-go/internal/format"
	"github.com/mvgrimes/timetrap-go/internal/models"
	"github.com/mvgrimes/timetrap-go/internal/parse"
	"github.com/mvgrimes/timetrap-go/internal/tt"
)

var editCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"e", "ed"},
	Short:   "Alter an entry's note, start, or end time. Defaults to the active entry. Defaults to the last entry to be checked out of if no entry is active.",
	Long:    "Alter an entry's note, start, or end time. Defaults to the active entry. Defaults to the last entry to be checked out of if no entry is active.",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetInt("id")
		start, _ := cmd.Flags().GetString("start")
		end, _ := cmd.Flags().GetString("end")
		appendToNote, _ := cmd.Flags().GetBool("append")
		move, _ := cmd.Flags().GetString("move")
		runEdit(id, start, end, appendToNote, move, args)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.PersistentFlags().IntP("id", "i", 0, "Alter entry with id <id> instead of the running entry")
	editCmd.PersistentFlags().StringP("start", "s", "", "Change the start time to <time>")
	editCmd.PersistentFlags().StringP("end", "e", "", "Change the end time to <time>")
	editCmd.PersistentFlags().BoolP("append", "z", false, "Append to the current note instead of replacing it the delimiter between appends notes is configurable (see config)")
	editCmd.PersistentFlags().StringP("move", "m", "", "Move to another sheet")
}

func runEdit(id int, start string, end string, appendToNote bool, move string, args []string) {
	t := tt.New(viper.GetString("database_file"))

	var entry models.Entry
	if id > 0 {
		entry = t.DB.GetEntry(id)
	} else {
		entry = t.DB.GetCurrentEntry()
		id = entry.ID
	}

	if entry.ID == 0 {
		fmt.Println("Can't find entry")
		os.Exit(1)
	}

	sheet := entry.Sheet
	if move != "" {
		sheet = move
	}

	var err error

	startTime := entry.Start.Time
	if start != "" {
		startTime, err = parse.Time(start)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	endTime := entry.End.Time
	if end != "" {
		endTime, err = parse.Time(end)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(endTime)
	}

	note := entry.Note
	if len(args) > 0 {
		// TODO: respect append_notes_delimiter
		notes := []string{}
		if appendToNote {
			notes = []string{note}
		}
		notes = append(notes, args...)
		note = strings.Join(notes, " ")
	}
	err = t.DB.UpdateEntry(id, sheet, startTime, endTime, note)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Print the entry via display
	entry = t.DB.GetEntry(entry.ID)
	entries := []models.Entry{entry}
	format.DisplayEntries(entries, entry.Sheet, false)
}
