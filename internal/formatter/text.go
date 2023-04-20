package formatter

import (
	"fmt"
	"time"

	"github.com/mvgrimes/timetrap-go/internal/format"
	"github.com/mvgrimes/timetrap-go/internal/models"
)

func FormatAsText(entries []models.Entry, sheet string, includeIds bool) {
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

		// Set the dayHeader to the current day if it is a new day
		dayHeader := ""
		if lastDay != entryDay {
			dayHeader = entryDay
		}

		id := ""
		if includeIds {
			id = fmt.Sprintf("%d", entry.ID)
		}

		endTimeStr := ""
		if entry.End.Valid {
			endTimeStr = entry.End.Time.Format("15:04:05")
			// endTimeStr = entry.End.Time.Format("15:04:05 -0700")
		}

		// Get the current time but treat UTC as localtime
		endTime := time.Now()
		_, offset := endTime.Zone()
		endTime = endTime.Add(time.Second * time.Duration(offset))

		if entry.End.Valid {
			endTime = entry.End.Time
		}
		duration := endTime.Sub(entry.Start.Time)

		fmt.Printf(
			"%-4s %-18s %-8s - %-8s   %8s   %s\n",
			id,
			dayHeader,
			entry.Start.Time.Format("15:04:05"),
			endTimeStr,
			format.Duration(duration),
			entry.Note,
		)

		// fmt.Printf(
		// 	"%-4s %-18s %-14s - %-14s   %8s   %s\n",
		// 	id,
		// 	dayHeader,
		// 	entry.Start.Time.Format("15:04:05 -0700"),
		// 	endTimeStr,
		// 	format.Duration(duration),
		// 	entry.Note,
		// )

		// if endTimeStr == "" {
		// 	fmt.Printf("%s\n", endTime.Format("15:04:05 -0700"))
		// }

		// Save the dayHeader if it wasn't ""
		if dayHeader != "" {
			lastDay = dayHeader
		}

		// Add to the total duration
		total += duration
		// TODO: add daily total
	}

	fmt.Printf("    --------------------------------------------------\n")
	fmt.Printf("    Total %43s\n", format.Duration(total))

}
