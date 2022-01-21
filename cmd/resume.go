package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:     "resume",
	Aliases: []string{""},
	Short:   "",
	Long:    "",
	Run: func(cmd *cobra.Command, args []string) {
		runResume(args)
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)

	resumeCmd.PersistentFlags().StringP("start", "s", "", "Include entries that start on this date or later")
	resumeCmd.PersistentFlags().StringP("end", "e", "", "Include entries that start on this date or earlier")
	resumeCmd.PersistentFlags().StringP("grep", "g", "", "Include entries where the note matches this regexp.")
}

// TODO: implement resume command
func runResume(args []string) {
	if len(args) != 0 {
		fmt.Printf("usage: t resume [SHEET]")
		os.Exit(1)
	}

	fmt.Println("resume command is not yet implemented")
	os.Exit(1)
}
