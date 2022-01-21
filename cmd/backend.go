package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	cmd := exec.Command("sqlite3", viper.GetString("database_file"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
