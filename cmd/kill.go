package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var killCmd = &cobra.Command{
	Use:     "kill",
	Aliases: []string{""},
	Short:   "",
	Long:    "",
	Run: func(cmd *cobra.Command, args []string) {
		runKill(args)
	},
}

func init() {
	rootCmd.AddCommand(killCmd)

	killCmd.PersistentFlags().StringP("start", "s", "", "Include entries that start on this date or later")
	killCmd.PersistentFlags().StringP("end", "e", "", "Include entries that start on this date or earlier")
	killCmd.PersistentFlags().StringP("grep", "g", "", "Include entries where the note matches this regexp.")

	viper.BindPFlag("start", killCmd.PersistentFlags().Lookup("start"))
	viper.BindPFlag("end", killCmd.PersistentFlags().Lookup("end"))
	viper.BindPFlag("grep", killCmd.PersistentFlags().Lookup("grep"))
}

func runKill(args []string) {
	if len(args) != 0 {
		fmt.Printf("usage: t kill [SHEET]")
		os.Exit(1)
	}

	fmt.Println("kill command is not yet implemented")
	os.Exit(1)
}
