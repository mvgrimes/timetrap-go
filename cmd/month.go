package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var monthCmd = &cobra.Command{
	Use:     "month",
	Aliases: []string{""},
	Short:   "",
	Long:    "",
	Run: func(cmd *cobra.Command, args []string) {
		runMonth(args)
	},
}

func init() {
	rootCmd.AddCommand(monthCmd)

	monthCmd.PersistentFlags().BoolP("ids", "v", false, "Print database ids (for use with edit)")
	monthCmd.PersistentFlags().StringP("start", "s", "", "Include entries that start on this date or later")
	monthCmd.PersistentFlags().StringP("format", "f", "", `The output format.
Valid built-in formats are ical, csv, json, ids, factor, and text (default).
Documentation on defining custom formats can be found in the README included
in this`)
}

// TODO: implement month command
func runMonth(args []string) {
	if len(args) != 0 {
		fmt.Printf("usage: t month [SHEET | all]")
		os.Exit(1)
	}

	fmt.Println("month command is not yet implemented")
	os.Exit(1)
}
