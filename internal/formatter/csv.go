package formatter

import (
	"fmt"
	"strings"

	"github.com/mvgrimes/timetrap-go/internal/tt"
)

func FormatAsCsv(entries []tt.Entry, sheet string, includeIds bool) {
	header := []string{"start", "end", "note", "sheet"}
	if includeIds {
		header = append([]string{"id"}, header...)
	}
	fmt.Println(strings.Join(header, ","))

	dateFmt := "2006-01-02 15:04:05"
	for _, entry := range entries {
		// Ruby version does not include any currently running entries
		if !entry.End.Valid {
			continue
		}

		fields := []string{
			entry.Start.Time.Format(dateFmt),
			entry.End.Time.Format(dateFmt),
			entry.Sheet,
			entry.Note,
		}
		if includeIds {
			fields = append([]string{fmt.Sprintf("%d", entry.ID)}) // Is this the most efficient way to conver to a string?
		}

		for i, v := range fields {
			fields[i] = fmt.Sprintf(`"%s"`, v)
		}

		fmt.Println(strings.Join(fields, ","))
	}
}
