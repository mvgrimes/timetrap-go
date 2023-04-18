package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mvgrimes/timetrap-go/internal/parse"
	"github.com/mvgrimes/timetrap-go/internal/tt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var inCmd = &cobra.Command{
	Use:     "in",
	Aliases: []string{"i"},
	Short:   "Start the timer for the current timesheet.",
	Long:    `Start the timer for the current timesheet.`,
	Run: func(cmd *cobra.Command, args []string) {
		atTime, _ := cmd.Flags().GetString("at")
		sheet, _ := cmd.Flags().GetString("sheet")
		runIn(atTime, sheet, args)
	},
}

func init() {
	rootCmd.AddCommand(inCmd)

	inCmd.PersistentFlags().StringP("at", "a", "", "Use this time instead of now <time:qs>")
	inCmd.PersistentFlags().StringP("sheet", "s", "", "Swith to this sheet frist")

	// Not sure why we had this?
	// viper.BindPFlag("at", inCmd.PersistentFlags().Lookup("at"))
}

func runIn(atTimeStr string, sheet string, args []string) {
	note := ""
	if len(args) > 0 {
		note = strings.Join(args, " ")
	}

	atTime, err := parse.Time(atTimeStr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	t := tt.New(viper.GetString("database_file"))
	meta := t.DB.GetMeta()

	if len(sheet) > 0 && sheet != meta.CurrentSheet {
		// TODO: check that sheet already exists

		meta, err = t.DB.SwitchSheet(sheet)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("Switched to sheet \"%s\"\n", meta.CurrentSheet)
	} else {
		sheet = meta.CurrentSheet
	}

	entry, err := t.ClockIn(sheet, atTime, note)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Checked into sheet \"%s\"\n", entry.Sheet)
}
