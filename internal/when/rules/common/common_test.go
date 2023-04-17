package common_test

import (
	// "log"
	"testing"
	"time"

	"github.com/olebedev/when"
	// "github.com/olebedev/when/rules/common"
	"github.com/mvgrimes/timetrap-go/internal/when/rules/common"
	"github.com/stretchr/testify/require"
)

var null = time.Date(2016, time.July, 15, 0, 0, 0, 0, time.UTC)

// July 15 days offset from the begining of the year
const OFFSET = 197

type Fixture struct {
	Text   string
	Index  int
	Phrase string
	// Diff   time.Duration
	Want time.Time
}

func ApplyFixtures(t *testing.T, name string, w *when.Parser, fixt []Fixture) {
	for i, f := range fixt {
		res, err := w.Parse(f.Text, null)
		require.Nil(t, err, "[%s] err #%d", name, i)
		require.NotNil(t, res, "[%s] res #%d", name, i)
		require.Equal(t, f.Index, res.Index, "[%s] index #%d", name, i)
		require.Equal(t, f.Phrase, res.Text, "[%s] text #%d", name, i)
		// log.Printf("%v", res.Time)
		// log.Printf("%v", res.Time.Sub(null))
		// require.Equal(t, f.Diff, res.Time.Sub(null), "[%s] diff #%d %s", name, i, f.Phrase)
		require.Equal(t, f.Want, res.Time, "[%s] diff #%d %s", name, i, f.Phrase)
	}
}

func ApplyFixturesNil(t *testing.T, name string, w *when.Parser, fixt []Fixture) {
	for i, f := range fixt {
		res, err := w.Parse(f.Text, null)
		require.Nil(t, err, "[%s] err #%d", name, i)
		require.Nil(t, res, "[%s] res #%d", name, i)
	}
}

func ApplyFixturesErr(t *testing.T, name string, w *when.Parser, fixt []Fixture) {
	for i, f := range fixt {
		_, err := w.Parse(f.Text, null)
		require.NotNil(t, err, "[%s] err #%d", name, i)
		require.Equal(t, f.Phrase, err.Error(), "[%s] err text #%d", name, i)
	}
}

func TestAll(t *testing.T) {
	w := when.New(nil)
	w.Add(common.All...)

	// res, _ := w.Parse("7/14", null)
	// log.Printf("7/14 = %v", res)
	// res, _ = w.Parse("7/15", null)
	// log.Printf("7/15 = %v", res)
	// res, _ = w.Parse("7/16", null)
	// log.Printf("7/16 = %v", res)

	// complex cases
	fixt := []Fixture{}
	ApplyFixtures(t, "common.All...", w, fixt)
}
