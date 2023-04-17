package models

import (
	"time"
)

type Sheet struct {
	Sheet      string        `json:"sheet"`
	Running    time.Duration `json:"running"`
	Today      time.Duration `json:"today"`
	Total      time.Duration `json:"total"`
	Active     bool
	LastActive bool
	// puts " %-#{width}s%-12s%-12s%s" % ["Timesheet", "Running", "Today", "Total Time"]
}
