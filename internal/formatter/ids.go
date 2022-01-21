package formatter

import (
	"fmt"
	"strings"

	"github.com/mvgrimes/timetrap-go/internal/tt"
)

func FormatAsIds(entries []tt.Entry, sheet string, includeIds bool) {
	ids := []string{}

	for _, entry := range entries {
		ids = append(ids, fmt.Sprintf("%d", entry.ID))
	}

	fmt.Println(strings.Join(ids, " "))
}
