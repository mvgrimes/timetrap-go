package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "t",
	Short: "Timetrap - Simple Time Tracking (go)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root called")
		// TODO: print the debug and examples
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolP("round", "r", false, "Round output to 15 minute start and end times.")
	rootCmd.PersistentFlags().BoolP("yes", "y", false, "Noninteractive, assume yes as answer to all prompts.")
	rootCmd.PersistentFlags().Bool("debug", false, "Display stack traces for errors.")
	viper.BindPFlag("round", rootCmd.PersistentFlags().Lookup("round"))
	viper.BindPFlag("yes", rootCmd.PersistentFlags().Lookup("yes"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

func initConfig() {
	viper.SetConfigType("yaml")
	configFile := os.ExpandEnv("$TIMETRAP_CONFIG_FILE")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath("$HOME/.timetrap.yml")
	}

	viper.SetDefault("database_file", os.ExpandEnv("$HOME/.timetrap.db"))
	// Unit of time for rounding (-r) in seconds
	viper.SetDefault("round_in_seconds", 900)
	// delimiter used when appending notes with `t edit --append`
	viper.SetDefault("append_notes_delimiter", " ")
	// an array of directories to search for user defined fomatter classes
	// viper.SetDefault("formatter_search_paths" , [
	//   "$HOME/.timetrap/formatters"
	// ])
	// formatter to use when display is invoked without a --format option
	viper.SetDefault("default_formatter", "text")
	// the auto_sheet to use
	viper.SetDefault("auto_sheet", "dotfiles")
	// an array of directories to search for user defined auto_sheet classes
	// viper.SetDefault("auto_sheet_search_paths" , [
	//   "$HOME/.timetrap/auto_sheets"
	// ])
	// the default command to when you run `t`.  default to printing usage.
	viper.SetDefault("default_command", nil)
	// only allow one running entry at a time.
	// automatically check out of any running tasks when checking in.
	viper.SetDefault("auto_checkout", false)
	// interactively prompt for a note if one isn't passed when checking in.
	viper.SetDefault("require_note", false)
	// command to launch external editor (false if no external editor used)
	viper.SetDefault("note_editor", false)
	// set day of the week when determining start of the week.
	viper.SetDefault("week_start", "Monday")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// no config file, just use defaults
			log.Println("unable to load config file using defaults")
		} else {
			panic(fmt.Errorf("Fatal error reading config file: %w\n", err))
		}
	}
}
