package format

import (
	"fmt"

	"github.com/mvgrimes/timetrap-go/internal/tt"
)

func DisplayList(summaries []tt.SheetSummary, includeArchived bool) {
	fmt.Printf(" %-15s  %-10s  %-10s  %s\n",
		"Timesheet", "Running", "Today", "Total Time")
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
