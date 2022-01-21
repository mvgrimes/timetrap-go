package format

import (
	"fmt"
	"time"
)

func Duration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("% 2d:%02d:%02d", h, m, s)
}
