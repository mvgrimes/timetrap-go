package cmd

import (
	"fmt"
	"os"

	"github.com/mvgrimes/timetrap-go/internal/format"
	"github.com/mvgrimes/timetrap-go/internal/tt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var sheetCmd = &cobra.Command{
	Use:     "sheet",
	Aliases: []string{"s", "sh"},
	Short:   "Stop the timer for a timesheet.",
	Long: `Switch to a timesheet creating it if necessary. When no sheet is
			specified list all sheets. The special sheetname '-' will switch to the
			last active sheet`,
	Run: func(cmd *cobra.Command, args []string) {
		runSheet(args)
	},
}

func init() {
	rootCmd.AddCommand(sheetCmd)
}

func runSheet(args []string) {
	if len(args) == 0 {
		t := tt.New(viper.GetString("database_file"))
		summaries := t.DB.GetSheets()
		format.DisplayList(summaries, true)
		return
	} else if len(args) > 1 {
		fmt.Println("usage: t sheet [TIMESHEET]")
		os.Exit(1)
	}

	sheet := args[0]

	t := tt.New(viper.GetString("database_file"))

	meta, err := t.DB.SwitchSheet(sheet)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Switching to sheet \"%s\"\n", meta.CurrentSheet)
}
