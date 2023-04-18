package parse

import (
	"errors"
	"fmt"
	"time"

	"github.com/mvgrimes/when/rules"
	"github.com/mvgrimes/when/rules/common"
	"github.com/mvgrimes/when/rules/en"
	"github.com/mvgrimes/when"
)

func Time(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Now(), nil
	}

	w := when.New(&rules.Options{
		Distance:     5,
		MatchByOrder: true,
		WantPast:     true})
	w.Add(en.All...)
	w.Add(common.SlashMDY(rules.Override))

	r, err := w.Parse(timeStr, time.Now())
	if err != nil || r == nil {
		return time.Now(), errors.New(fmt.Sprintf("Unable to parse time: %s\n", timeStr))
	}

	return r.Time, nil
}
