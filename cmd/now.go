package cmd

import (
	"fmt"
	"os"

	"github.com/mvgrimes/timetrap-go/internal/tt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var nowCmd = &cobra.Command{
	Use:     "now",
	Aliases: []string{"n"},
	Short:   "Show all running entries.",
	Long:    `Show all running entries.`,
	Run: func(cmd *cobra.Command, args []string) {
		runNow(args)
	},
}

func init() {
	rootCmd.AddCommand(nowCmd)
}

func runNow(args []string) {
	if len(args) != 0 {
		fmt.Printf("usage: now\n")
		os.Exit(1)
	}

	t := tt.TimeTrap{}
	t.Connect(viper.GetString("database_file"))
	meta := t.GetMeta()

	fmt.Printf("%s\n", meta.CurrentSheet)
}
