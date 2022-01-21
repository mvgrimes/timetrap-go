package parse

import (
	"errors"
	"fmt"
	"time"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

func Time(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Now(), nil
	}

	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	r, err := w.Parse(timeStr, time.Now())
	if err != nil || r == nil {
		return time.Now(), errors.New(fmt.Sprintf("Unable to parse time: %s\n", timeStr))
	}

	return r.Time, nil
}
