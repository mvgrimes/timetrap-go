package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mvgrimes/timetrap-go/internal/tt"
)

var killCmd = &cobra.Command{
	Use:     "kill",
	Aliases: []string{"k"},
	Short:   "Delete a timesheet or an entry.",
	Long:    "Delete a timesheet or an entry.",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetInt("id")
		runKill(id, args)
	},
}

func init() {
	rootCmd.AddCommand(killCmd)

	killCmd.PersistentFlags().IntP("id", "", 0, "Delete entry with id <id> instead of timesheet")
}

func runKill(id int, args []string) {
	if id == 0 {
		if len(args) != 1 {
			fmt.Printf("usage: t kill [SHEET]")
			os.Exit(1)
		}
		killSheet(args[0])
	} else {
		killEntry(id)
	}
}

func killEntry(id int) {
	t := tt.TimeTrap{}
	t.Connect(viper.GetString("database_file"))

	err := t.DeleteEntry(id)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
}

func killSheet(sheet string) {
	t := tt.TimeTrap{}
	t.Connect(viper.GetString("database_file"))

	err := t.DeleteSheet(sheet)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
}
