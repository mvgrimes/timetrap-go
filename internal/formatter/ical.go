package formatter

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mvgrimes/timetrap-go/internal/models"
)

func FormatAsIcal(entries []models.Entry, sheet string, includeIds bool) {
	fmt.Println("BEGIN:VCALENDAR")
	fmt.Println("CALSCALE:GREGORIAN")
	fmt.Println("METHOD:PUBLISH")
	fmt.Println("PRODID:iCalendar-Ruby")
	fmt.Println("VERSION:2.0")

	dateFmt := "20060102T150405"
	for _, entry := range entries {
		// Ruby version does not include any currently running entries
		if !entry.End.Valid {
			continue
		}

		fmt.Println("BEGIN:VEVENT")
		fmt.Printf("DESCRIPTION:%s\n", entry.Note)
		fmt.Printf("DTEND:%s\n", entry.End.Time.Format(dateFmt))
		fmt.Printf("DTSTAMP:%s\n", time.Now().Format(dateFmt))
		fmt.Printf("DTSTART:%s\n", entry.Start.Time.Format(dateFmt))
		fmt.Printf("SEQUENCE:%d\n", 0)
		fmt.Printf("SUMMARY:%s\n", entry.Note)
		fmt.Printf("UID:%s\n", Uid())
		fmt.Println("END:VEVENT")

	}
	fmt.Println("END:VCALENDAR")
}

func Uid() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s_%d@%s",
		time.Now().Format("2006-01-02T15:04:05-0700"), // -05:00
		rand.Int31(),
		hostname,
	)
}
