package formatter

import (
	"fmt"
	"strings"

	"github.com/mvgrimes/timetrap-go/internal/models"
)

func FormatAsJson(entries []models.Entry, sheet string, includeIds bool) {
	// Ruby version defaults to including id

	fmtEntries := []string{}
	dateFmt := "2006-01-02 15:04:05 -0500"
	for _, entry := range entries {
		// Ruby version does not include any currently running entries
		if !entry.End.Valid {
			continue
		}

		// TODO: escape quotes
		fmtEntry := fmt.Sprintf(`{"id":%d,"note":"%s","start":"%s","end":"%s","sheet":"%s"}`,
			entry.ID,
			entry.Note,
			entry.Start.Time.Format(dateFmt),
			entry.End.Time.Format(dateFmt),
			entry.Sheet,
		)

		fmtEntries = append(fmtEntries, fmtEntry)
	}

	fmt.Printf("[%s]\n", strings.Join(fmtEntries, ","))
}
