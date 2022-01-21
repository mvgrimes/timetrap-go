package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var backendCmd = &cobra.Command{
	Use:     "backend",
	Aliases: []string{"b", "back"},
	Short:   "Open an sqlite shell to the database.",
	Long:    "Open an sqlite shell to the database.",
	Run: func(cmd *cobra.Command, args []string) {
		runBackend(args)
	},
}

func init() {
	rootCmd.AddCommand(backendCmd)
}

func runBackend(args []string) {
	if len(args) != 0 {
		fmt.Printf("usage: t backend [SHEET]")
		os.Exit(1)
	}

	// TODO: open the sqlite file
	fmt.Println("backend command is not yet implemented")
	os.Exit(1)
}
