package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

var todayCmd = &cobra.Command{
	Use:     "today",
	Aliases: []string{""},
	Short:   "",
	Long:    "",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("%#v\n", cmd)
		// x, _ := cmd.Flags().GetBool("ids")
		// fmt.Printf("%#v\n", x)
		runToday(args)
	},
}

func init() {
	rootCmd.AddCommand(todayCmd)

	todayCmd.PersistentFlags().BoolP("ids", "v", false, "Print database ids (for use with edit)")
	todayCmd.PersistentFlags().StringP("format", "f", "", `The output format.
Valid built-in formats are ical, csv, json, ids, factor, and text (default).
Documentation on defining custom formats can be found in the README included
in this`)
}

// TODO: implement the today command
func runToday(args []string) {
	if len(args) != 0 {
		fmt.Printf("usage: t today [SHEET | all]")
		os.Exit(1)
	}

	fmt.Println("today command is not yet implemented")
	os.Exit(1)
}
