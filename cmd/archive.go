package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

var archiveCmd = &cobra.Command{
	Use:     "archive",
	Aliases: []string{"a", "arch"},
	Short:   "Move entries to a hidden sheet (by default named '_[SHEET]') so they're out of the way.",
	Long:    "Move entries to a hidden sheet (by default named '_[SHEET]') so they're out of the way.",
	Run: func(cmd *cobra.Command, args []string) {
		runArchive(args)
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)

	archiveCmd.PersistentFlags().StringP("start", "s", "", "Include entries that start on this date or later")
	archiveCmd.PersistentFlags().StringP("end", "e", "", "Include entries that start on this date or earlier")
	archiveCmd.PersistentFlags().StringP("grep", "g", "", "Include entries where the note matches this regexp.")
}

func runArchive(args []string) {
	if len(args) != 0 {
		fmt.Printf("usage: t archive [SHEET]")
		os.Exit(1)
	}

	fmt.Println("archive command is not yet implemented")
	os.Exit(1)
}
