package common_test

import (
	"testing"
	"time"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules"
	// "github.com/olebedev/when/rules/common"
	"github.com/mvgrimes/timetrap-go/internal/when/rules/common"
)

// Reference date is 7/15/2016

func TestSlashMDY(t *testing.T) {
	fixt := []Fixture{
		// {"The Deadline is 10/10/2016", 16, "10/10/2016", (284 - OFFSET) * 24 * time.Hour},
		// {"The Deadline is 2/1/2016", 16, "2/1/2016", (32 - OFFSET) * 24 * time.Hour},
		// {"The Deadline is 2/29/2016", 16, "2/29/2016", (60 - OFFSET) * 24 * time.Hour},
		{"The Deadline is 10/10/2016", 16, "10/10/2016", time.Date(2016, 10, 10, 0, 0, 0, 0, time.UTC)},
		{"The Deadline is 2/1/2016", 16, "2/1/2016", time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC)},
		{"The Deadline is 2/29/2016", 16, "2/29/2016", time.Date(2016, 2, 29, 0, 0, 0, 0, time.UTC)},

		// next year
		// {"The Deadline is 2/28", 16, "2/28", (59 + 366 - OFFSET) * 24 * time.Hour},
		// {"The Deadline is 02/28/2017", 16, "02/28/2017", (59 + 366 - OFFSET) * 24 * time.Hour},
		// {"The Deadline is 2/28", 16, "2/28", time.Date(2017, 2, 28, 0, 0, 0, 0, time.UTC)},
		{"The Deadline is 02/28/2017", 16, "02/28/2017", time.Date(2017, 2, 28, 0, 0, 0, 0, time.UTC)},

		// right after w/o a year
		// {"The Deadline is 07/28", 16, "07/28", (210 - OFFSET) * 24 * time.Hour},
		{"The Deadline is 07/28", 16, "07/28", time.Date(2015, 7, 28, 0, 0, 0, 0, time.UTC)},

		// before w/o a year
		// {"The Deadline is 06/30", 16, "06/30", (181 + 366 - OFFSET) * 24 * time.Hour},
		// {"The Deadline is 06/30", 16, "06/30", time.Date(2017, 6, 30, 0, 0, 0, 0, time.UTC)},
		{"The Deadline is 06/30", 16, "06/30", time.Date(2016, 6, 30, 0, 0, 0, 0, time.UTC)},

		// prev day will be added to the future
		// {"The Deadline is 07/14", 16, "07/14", (195 + 366 - OFFSET) * 24 * time.Hour},
		// {"The Deadline is 07/14", 16, "07/14", time.Date(2017, time.July, 14, 0, 0, 0, 0, time.UTC)},
		{"The Deadline is 07/14", 16, "07/14", time.Date(2016, time.July, 14, 0, 0, 0, 0, time.UTC)},
	}

	w := when.New(nil)
	w.Add(common.SlashMDY(rules.Skip))

	ApplyFixtures(t, "common.SlashMDY", w, fixt)

}
