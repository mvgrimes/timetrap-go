package tt

import (
	"errors"
	"fmt"
	"log"
	"time"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Meta struct {
	CurrentSheet string `json:"current_sheet"`
	LastSheet    string `json:"last_sheet"`
	LastCheckout int    `json:"last_checkout"`
}

type Entry struct {
	ID    int          `json:"id"`
	Sheet string       `json:"sheet"`
	Start sql.NullTime `json:"start"`
	End   sql.NullTime `json:"end"`
	Note  string       `json:"note"`
}

type TimeTrap struct {
	Filename string
	db       *sql.DB
}

func (t *TimeTrap) Connect(filename string) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic(err.Error())
	}

	t.Filename = filename
	t.db = db
}

func (t *TimeTrap) Close() {
	t.db.Close()
}

func (t *TimeTrap) GetMeta() Meta {
	meta := Meta{}

	err := t.db.QueryRow("select value from meta where id = 1;").Scan(&meta.CurrentSheet)
	if err != nil {
		panic(err.Error())
	}
	err = t.db.QueryRow("select value from meta where id = 2;").Scan(&meta.LastSheet)
	if err != nil {
		panic(err.Error())
	}
	err = t.db.QueryRow("select value from meta where id = 3;").Scan(&meta.LastCheckout)
	if err != nil {
		panic(err.Error())
	}

	// log.Printf("meta: %v", meta)

	return meta
}

func (t *TimeTrap) GetCurrentEntry() Entry {
	meta := t.GetMeta()
	entry := Entry{}

	err := t.db.QueryRow(`SELECT id, sheet, start, end, note
						FROM entries
						WHERE sheet = ?
						ORDER BY id DESC
						LIMIT 1;`,
		meta.CurrentSheet,
	).Scan(&entry.ID, &entry.Sheet, &entry.Start, &entry.End, &entry.Note)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("entry: %v", entry)

	return entry
}

func (t *TimeTrap) GetEntries(sheet string) []Entry {
	results, err := t.db.Query(
		`SELECT id, sheet, start, end, note
				FROM entries
				WHERE sheet = ?;`,
		sheet)
	if err != nil {
		panic(err.Error())
	}

	entries := []Entry{}
	for results.Next() {
		var entry Entry
		err = results.Scan(&entry.ID, &entry.Sheet, &entry.Start, &entry.End, &entry.Note)
		if err != nil {
			panic(err.Error())
		}
		entries = append(entries, entry)
	}

	return entries
}

func (t *TimeTrap) Start(startTime time.Time, note string) (Entry, error) {
	entry := t.GetCurrentEntry()

	if !entry.End.Valid {
		return entry, errors.New("Timetrap is already running")
	}

	startTimeStr := startTime.Format("2006-01-02 15:04:05.999999")

	result, err := t.db.Exec(
		`INSERT INTO entries
				(start, sheet, note)
				VALUES
				(?, ?, ?);`,
		startTimeStr, entry.Sheet, note)
	if err != nil {
		panic(err.Error())
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	log.Printf("result id = %d\n", id)

	return entry, nil
}

func (t *TimeTrap) Stop(stopTime time.Time) (Entry, error) {
	entry := t.GetCurrentEntry()

	if entry.End.Valid {
		return entry, errors.New(fmt.Sprintf(`No running entry on sheet "%s".`, entry.Sheet))
	}

	stopTimeStr := stopTime.Format("2006-01-02 15:04:05.999999")

	res, err := t.db.Exec(
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

	return entry, nil
}

func (t *TimeTrap) SwitchSheet(sheet string) (Meta, error) {
	meta := t.GetMeta()

	if sheet == "-" {
		sheet = meta.LastSheet
	}

	if sheet == meta.CurrentSheet {
		return Meta{}, errors.New(fmt.Sprintf(`Already on sheet "%s"`, sheet))
	}

	res, err := t.db.Exec(
		`UPDATE meta SET value = ?
		   WHERE id = ?;`,
		meta.CurrentSheet, 2)
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

	res, err = t.db.Exec(
		`UPDATE meta SET value = ?
		   WHERE id = ?;`,
		sheet, 1)
	if err != nil {
		panic(err.Error())
	}

	rowCnt, err = res.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	if rowCnt != 1 {
		panic(fmt.Sprintf("wrong number of rows updated: %d\n", rowCnt))
	}

	return t.GetMeta(), nil
}

type SheetSummary struct {
	Sheet      string        `json:"sheet"`
	Running    time.Duration `json:"running"`
	Today      time.Duration `json:"today"`
	Total      time.Duration `json:"total"`
	Active     bool
	LastActive bool
	// puts " %-#{width}s%-12s%-12s%s" % ["Timesheet", "Running", "Today", "Total Time"]
}

func (t *TimeTrap) List() []SheetSummary {
	meta := t.GetMeta()

	results, err := t.db.Query(
		`SELECT
				sheet,
				sum(case when end is null then strftime("%s",'now', 'localtime')-strftime("%s",start) else 0 end)*1000000000 as running,
				sum(case when strftime("%Y%j",start)=strftime("%Y%j",datetime('now','localtime')) then
						strftime("%s",ifnull(end,datetime('now','localtime')))-strftime("%s",start)
					else 0 end)*1000000000 as today,
				sum(strftime("%s",ifnull(end,datetime('now','localtime')))-strftime("%s",start))*1000000000 as total
		FROM ENTRIES
		GROUP BY sheet
		ORDER BY sheet;`,
	)
	if err != nil {
		panic(err.Error())
	}

	summaries := []SheetSummary{}
	for results.Next() {
		var summary SheetSummary
		err = results.Scan(&summary.Sheet, &summary.Running, &summary.Today, &summary.Total)
		if err != nil {
			panic(err.Error())
		}
		summary.Active = summary.Sheet == meta.CurrentSheet
		summary.LastActive = summary.Sheet == meta.LastSheet
		// fmt.Printf("s: %v\n", summary)
		summaries = append(summaries, summary)
	}

	return summaries
}

type SheetDetails struct {
	ID       string `json:"id"`
	Day      string
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
	Note     string        `json:"note"`
}

func (t *TimeTrap) Display() []SheetDetails {
	meta := t.GetMeta()

	results, err := t.db.Query(`
		SELECT
				id,
				start,
				end,
				(strftime("%s",ifnull(end,datetime('now','localtime')))-strftime("%s",start))*1000000000 as duration,
				note
		FROM ENTRIES
		WHERE sheet = ?
		ORDER BY start;
		`, meta.CurrentSheet)
	if err != nil {
		panic(err.Error())
	}

	summaries := []SheetDetails{}
	for results.Next() {
		var summary SheetDetails
		err = results.Scan(&summary.ID, &summary.Start, &summary.End, &summary.Duration, &summary.Note)
		summary.Day = summary.Start.Format("Mon Jan 02, 2006")
		if err != nil {
			panic(err.Error())
		}
		summaries = append(summaries, summary)
	}

	return summaries
}
