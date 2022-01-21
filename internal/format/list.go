package format

import (
	"fmt"
	"time"

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

func DisplayEntries(entries []tt.SheetDetails, sheet string, includeIds bool) {
	idHeader := ""
	if includeIds {
		idHeader = "Id"
	}
	fmt.Printf("Timesheet: %s\n", sheet)
	fmt.Printf("%-4s %-18s %-8s   %-8s   %8s   %s\n", idHeader, "Day", "Start", "End", "Duration", "Notes")

	lastDay := ""
	var total time.Duration
	for _, entry := range entries {
		day := ""
		if lastDay != entry.Day {
			day = entry.Day
		}
		id := ""
		if includeIds {
			id = fmt.Sprintf("%d", entry.ID)
		}
		fmt.Printf(
			"%-4s %-18s %-8s - %-8s   %8s   %s\n",
			id,
			day,
			entry.Start.Time.Format("15:04:05"),
			entry.End.Time.Format("15:04:05"),
			Duration(entry.Duration),
			entry.Note,
		)
		// TODO: display "" for empty end
		lastDay = day
		total += entry.Duration
		// TODO: add daily total
	}
	fmt.Printf("    -------------------------------------------------\n")
	fmt.Printf("    Total %43s\n", Duration(total))
}
