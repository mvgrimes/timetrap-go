package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

var weekCmd = &cobra.Command{
	Use:     "week",
	Aliases: []string{""},
	Short:   "",
	Long:    "",
	Run: func(cmd *cobra.Command, args []string) {
		runWeek(args)
	},
}

func init() {
	rootCmd.AddCommand(weekCmd)

	weekCmd.PersistentFlags().BoolP("ids", "v", false, "Print database ids (for use with edit)")
	weekCmd.PersistentFlags().StringP("end", "e", "", "Include entries that start on this date or earlier")
	weekCmd.PersistentFlags().StringP("format", "f", "", `The output format.
Valid built-in formats are ical, csv, json, ids, factor, and text (default).
Documentation on defining custom formats can be found in the README included
in this`)
}

func runWeek(args []string) {
	if len(args) != 0 {
		fmt.Printf("usage: t week [SHEET | all]")
		os.Exit(1)
	}

	fmt.Println("week command is not yet implemented")
	os.Exit(1)
}
