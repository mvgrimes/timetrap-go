package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/mvgrimes/timetrap-go/internal/format"
	"github.com/mvgrimes/timetrap-go/internal/tt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var nowCmd = &cobra.Command{
	Use:     "now",
	Aliases: []string{"n"},
	Short:   "Show all running entries.",
	Long:    `Show all running entries.`,
	Run: func(cmd *cobra.Command, args []string) {
		runNow(args)
	},
}

func init() {
	rootCmd.AddCommand(nowCmd)
}

func runNow(args []string) {
	if len(args) != 0 {
		fmt.Printf("usage: now\n")
		os.Exit(1)
	}

	t := tt.New(viper.GetString("database_file"))
	meta := t.DB.GetMeta()

	entry := t.DB.GetCurrentEntry()
	state := "not running"
	if entry.Start.Valid && !entry.End.Valid {
		// Get the current time but treat UTC as localtime
		endTime := time.Now()
		_, offset := endTime.Zone()
		endTime = endTime.Add(time.Second * time.Duration(offset))

		state = format.Duration(endTime.Sub(entry.Start.Time))
	}

	fmt.Printf("*%s: %s\n", meta.CurrentSheet, state)

	// TODO: support multiple running sheets
}
