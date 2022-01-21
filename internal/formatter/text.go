package formatter

import (
	"fmt"
	"time"

	"github.com/mvgrimes/timetrap-go/internal/format"
	"github.com/mvgrimes/timetrap-go/internal/tt"
)

func FormatAsText(entries []tt.Entry, sheet string, includeIds bool) {
	idHeader := ""
	if includeIds {
		idHeader = "Id"
	}
	fmt.Printf("Timesheet: %s\n", sheet)
	fmt.Printf("%-4s %-18s %-8s   %-8s   %8s   %s\n", idHeader, "Day", "Start", "End", "Duration", "Notes")

	lastDay := ""
	var total time.Duration
	for _, entry := range entries {
		entryDay := entry.Start.Time.Format("Mon Jan 02, 2006")

		day := ""
		if lastDay != entryDay {
			day = entryDay
		}

		id := ""
		if includeIds {
			id = fmt.Sprintf("%d", entry.ID)
		}

		endTimeStr := ""
		if entry.End.Valid {
			endTimeStr = entry.End.Time.Format("15:04:05")
		}

		endTime := time.Now()
		if entry.End.Valid {
			endTime = entry.End.Time
		}
		duration := endTime.Sub(entry.Start.Time)

		fmt.Printf(
			"%-4s %-18s %-8s - %-8s   %8s   %s\n",
			id,
			day,
			entry.Start.Time.Format("15:04:05"),
			endTimeStr,
			format.Duration(duration),
			entry.Note,
		)

		if day != "" {
			lastDay = day
		}
		total += duration
		// TODO: add daily total
	}

	fmt.Printf("    --------------------------------------------------\n")
	fmt.Printf("    Total %43s\n", format.Duration(total))
}
