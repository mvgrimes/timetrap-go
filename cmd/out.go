package cmd

import (
	"fmt"
	"os"

	"github.com/mvgrimes/timetrap-go/internal/parse"
	"github.com/mvgrimes/timetrap-go/internal/tt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var outCmd = &cobra.Command{
	Use:     "out",
	Aliases: []string{"o"},
	Short:   "Stop the timer for a timesheet.",
	Long:    `Stop the timer for a timesheet.`,
	Run: func(cmd *cobra.Command, args []string) {
		atTimeStr, _ := cmd.Flags().GetString("at")
		runOut(atTimeStr, args)
	},
}

func init() {
	rootCmd.AddCommand(outCmd)

	outCmd.PersistentFlags().StringP("at", "a", "", "Use this time instead of now <time:qs>")
}

func runOut(atTimeStr string, args []string) {
	// TODO: clock out of a particular sheet
	// if len(args) == 1 {
	// sheet = args[0]
	// }

	atTime, err := parse.Time(atTimeStr)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	t := tt.TimeTrap{}
	t.Connect(viper.GetString("database_file"))

	entry, err := t.Stop(atTime)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Checked out of sheet \"%s\"\n", entry.Sheet)
}
