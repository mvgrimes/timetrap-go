package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mvgrimes/timetrap-go/internal/tt"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var inCmd = &cobra.Command{
	Use:   "in",
	Short: "Start the timer for the current timesheet.",
	Long:  `Start the timer for the current timesheet.`,
	Run: func(cmd *cobra.Command, args []string) {
		runIn(args)
	},
}

func init() {
	rootCmd.AddCommand(inCmd)

	inCmd.PersistentFlags().StringP("at", "a", "", "Use this time instead of now <time:qs>")
	viper.BindPFlag("at", inCmd.PersistentFlags().Lookup("at"))
}

func runIn(args []string) {
	note := ""

	if len(args) > 0 {
		note = strings.Join(args, " ")
	}

	atTime := time.Now()
	atTimeStr := viper.GetString("at")
	if atTimeStr != "" {
		w := when.New(nil)
		w.Add(en.All...)
		w.Add(common.All...)

		r, err := w.Parse(atTimeStr, time.Now())
		if err != nil || r == nil {
			fmt.Printf("Unable to parse time: %s\n", atTimeStr)
			os.Exit(1)
		}

		atTime = r.Time
	}

	meta := tt.GetMeta(viper.GetString("database_file"))
	log.Printf("checking into sheet: %s\n", meta.CurrentSheet)
	log.Printf("in at: %v\n", atTime)
	log.Printf("note is: %s\n", note)

	entry := tt.GetCurrentEntry(viper.GetString("database_file"))
	log.Printf("start: %v\n", entry.Start)
	log.Printf("start: %t\n", entry.Start)
	log.Printf("end: %v\n", entry.End)
	log.Printf("end: %t\n", entry.End)
	if !entry.End.Valid {
		fmt.Println("Timetrap is already running")
		os.Exit(1)
	}

	tt.Start(viper.GetString("database_file"), meta.CurrentSheet, atTime, note)
	fmt.Printf(`Checked into sheet "%s"`, meta.CurrentSheet)
}
