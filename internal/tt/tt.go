package tt

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/mvgrimes/timetrap-go/internal/db"
	"github.com/mvgrimes/timetrap-go/internal/models"
)

// TODO: Keep this only in the db/models file
const dateFmt = "2006-01-02 15:04:05.999999"

type TimeTrap struct {
	DB *db.Database
}

func New(filename string) *TimeTrap {
	db := db.New(filename)

	return &TimeTrap{
		DB: db,
	}
}

func (t *TimeTrap) ClockIn(sheet string, startTime time.Time, note string) (models.Entry, error) {
	entry := t.DB.GetCurrentEntry() // respect sheet?

	if entry.Start.Valid && !entry.End.Valid {
		return entry, errors.New("Timetrap is already running")
	}

	startTimeStr := startTime.Format(dateFmt)

	result, err := t.DB.Conn.Exec(
		`INSERT INTO entries
				(start, sheet, note)
				VALUES
				(?, ?, ?);`,
		startTimeStr, sheet, note)
	if err != nil {
		panic(err.Error())
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	log.Debugf("result id = %d\n", id)

	return t.DB.GetCurrentEntry(), nil
}

func (t *TimeTrap) ClockOut(stopTime time.Time) (models.Entry, error) {
	entry := t.DB.GetCurrentEntry()

	if entry.End.Valid {
		return entry, errors.New(fmt.Sprintf(`No running entry on sheet "%s".`, entry.Sheet))
	}

	stopTimeStr := stopTime.Format(dateFmt)

	res, err := t.DB.Conn.Exec(
		`UPDATE entries SET end = ?
		   WHERE id = ?;`,
		stopTimeStr, entry.ID)
	if err != nil {
		panic(err.Error())
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	if rowCnt != 1 {
		panic(fmt.Sprintf("wrong number of rows updated: %d\n", rowCnt))
	}

	return t.DB.GetCurrentEntry(), nil
}
