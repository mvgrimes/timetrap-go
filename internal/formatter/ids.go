package formatter

import (
	"fmt"
	"strings"

	"github.com/mvgrimes/timetrap-go/internal/models"
)

func FormatAsIds(entries []models.Entry, sheet string, includeIds bool) {
	ids := []string{}

	for _, entry := range entries {
		ids = append(ids, fmt.Sprintf("%d", entry.ID))
	}

	fmt.Println(strings.Join(ids, " "))
}
