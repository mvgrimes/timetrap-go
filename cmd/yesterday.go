package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var yesterdayCmd = &cobra.Command{
	Use:     "yesterday",
	Aliases: []string{""},
	Short:   "",
	Long:    "",
	Run: func(cmd *cobra.Command, args []string) {
		runYesterday(args)
	},
}

func init() {
	rootCmd.AddCommand(yesterdayCmd)

	yesterdayCmd.PersistentFlags().BoolP("ids", "v", false, "Print database ids (for use with edit)")
	yesterdayCmd.PersistentFlags().StringP("format", "f", "", `The output format.
Valid built-in formats are ical, csv, json, ids, factor, and text (default).
Documentation on defining custom formats can be found in the README included
in this`)

	viper.BindPFlag("ids", yesterdayCmd.PersistentFlags().Lookup("ids"))
	viper.BindPFlag("format", yesterdayCmd.PersistentFlags().Lookup("format"))
}

func runYesterday(args []string) {
	if len(args) != 0 {
		fmt.Printf("usage: t yesterday [SHEET | all]")
		os.Exit(1)
	}

	fmt.Println("yesterday command is not yet implemented")
	os.Exit(1)
}
