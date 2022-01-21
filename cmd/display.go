package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/mvgrimes/timetrap-go/internal/formatter"
	"github.com/mvgrimes/timetrap-go/internal/parse"
	"github.com/mvgrimes/timetrap-go/internal/tt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var displayCmd = &cobra.Command{
	Use:     "display",
	Aliases: []string{"d"},
	Short:   "Display the current timesheet or a specific.",
	Long: `Display the current timesheet or a specific. Pass 'all' as SHEET
      to display all unarchived sheets or 'full' to display archived and
      unarchived sheets.`,
	Run: func(cmd *cobra.Command, args []string) {
		includeIds, _ := cmd.Flags().GetBool("ids")
		start, _ := cmd.Flags().GetString("start")
		end, _ := cmd.Flags().GetString("end")
		format, _ := cmd.Flags().GetString("format")
		grep, _ := cmd.Flags().GetString("grep")

		if format == "" {
			format = viper.GetString("default_formatter")
		}

		runDisplay(includeIds, start, end, grep, format, args)
	},
}

func init() {
	rootCmd.AddCommand(displayCmd)

	displayCmd.PersistentFlags().BoolP("ids", "v", false, "Print database ids (for use with edit)")
	displayCmd.PersistentFlags().StringP("start", "s", "", "Include entries that start on this date or later")
	displayCmd.PersistentFlags().StringP("end", "e", "", "Include entries that start on this date or earlier")
	displayCmd.PersistentFlags().StringP("format", "f", "", `The output format.
Valid built-in formats are ical, csv, json, ids, factor, and text (default).
Documentation on defining custom formats can be found in the README included
in this`)
	displayCmd.PersistentFlags().StringP("grep", "g", "", "Include entries where the note matches this regexp.")
}

// TODO: add output formatting
func runDisplay(includeIds bool, startStr string, endStr string, grep string, format string, args []string) {
	if len(args) > 0 {
		fmt.Println("usage: t display")
		os.Exit(1)
	}

	t := tt.TimeTrap{}
	t.Connect(viper.GetString("database_file"))

	var err error
	var start sql.NullTime
	if startStr != "" {
		start.Time, err = parse.Time(startStr)
		if err != nil {
			panic(err.Error())
		}
		start.Valid = true
	}

	var end sql.NullTime
	if endStr != "" {
		end.Time, err = parse.Time(endStr)
		if err != nil {
			panic(err.Error())
		}
		end.Valid = true
	}

	meta := t.GetMeta()
	entries := t.GetFilteredEntries(meta.CurrentSheet, start, end, grep)

	switch format {
	case "text":
		formatter.FormatAsText(entries, meta.CurrentSheet, includeIds)
	case "csv":
		formatter.FormatAsCsv(entries, meta.CurrentSheet, includeIds)
	case "json":
		formatter.FormatAsJson(entries, meta.CurrentSheet, includeIds)
	case "ical":
		formatter.FormatAsIcal(entries, meta.CurrentSheet, includeIds)
	case "ids":
		formatter.FormatAsIds(entries, meta.CurrentSheet, includeIds)
	default:
		fmt.Printf("Don't recognize that formatter: %s\n", format)
		os.Exit(1)
	}
}
