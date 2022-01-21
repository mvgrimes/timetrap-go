package format

import (
	"fmt"

	"github.com/mvgrimes/timetrap-go/internal/tt"
	"github.com/spf13/viper"
)

func DisplayList(includeArchived bool) {
	t := tt.TimeTrap{}
	t.Connect(viper.GetString("database_file"))

	fmt.Printf(" %-15s  %-10s  %-10s  %s\n",
		"Timesheet", "Running", "Today", "Total Time")
	summaries := t.List()
	for _, summary := range summaries {
		if (!includeArchived) && len(summary.Sheet) > 0 && summary.Sheet[0:1] == "_" {
			continue
		}

		active := " "
		if summary.LastActive {
			active = "-"
		}
		if summary.Active {
			active = "*"
		}

		fmt.Printf(
			"%s%-15s% 10s  % 10s%12s\n",
			active,
			summary.Sheet,
			Duration(summary.Running),
			Duration(summary.Today),
			Duration(summary.Total),
		)
	}
}
