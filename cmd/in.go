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
		runIn(atTime, args)
	},
}

func init() {
	rootCmd.AddCommand(inCmd)

	inCmd.PersistentFlags().StringP("at", "a", "", "Use this time instead of now <time:qs>")
	viper.BindPFlag("at", inCmd.PersistentFlags().Lookup("at"))
}

func runIn(atTimeStr string, args []string) {
	note := ""
	if len(args) > 0 {
		note = strings.Join(args, " ")
	}

	atTime, err := parse.Time(atTimeStr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	t := tt.TimeTrap{}
	t.Connect(viper.GetString("database_file"))

	entry, err := t.Start(atTime, note)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Checked into sheet \"%s\"\n", entry.Sheet)
}
