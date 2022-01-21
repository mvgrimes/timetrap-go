package tt

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
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

const dateFmt = "2006-01-02 15:04:05.999999"

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

	log.Debugf("meta: %v", meta)

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
		if err == sql.ErrNoRows {
			return Entry{}
		}
		panic(err.Error())
	}

	log.Debugf("entry: %v", entry)
	return entry
}

func (t *TimeTrap) GetEntry(id int) Entry {
	entry := Entry{}

	err := t.db.QueryRow(`SELECT id, sheet, start, end, note
						FROM entries
						WHERE id = ?
						ORDER BY id DESC
						LIMIT 1;`,
		id,
	).Scan(&entry.ID, &entry.Sheet, &entry.Start, &entry.End, &entry.Note)
	if err != nil {
		if err == sql.ErrNoRows {
			return Entry{}
		}
		panic(err.Error())
	}

	log.Debugf("entry: %v", entry)
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

func (t *TimeTrap) GetFilteredEntries(sheet string, start sql.NullTime, end sql.NullTime, grep string) []Entry {
	var params []interface{}
	params = append(params, sheet)

	sql := "SELECT id, sheet, start, end, note FROM entries WHERE sheet = ? "
	if start.Valid {
		sql = sql + "AND start >= ? "
		params = append(params, start.Time.Format("2006-01-02"))
	}
	if end.Valid {
		sql = sql + "AND start <= ? "
		params = append(params, end.Time.Add(time.Hour*24).Format("2006-01-02"))
	}
	if grep != "" {
		sql = sql + "AND note LIKE ? "
		params = append(params, grep)
	}
	sql = sql + ";"

	results, err := t.db.Query(sql, params...)
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
	meta := t.GetMeta()
	entry := t.GetCurrentEntry()

	if entry.Start.Valid && !entry.End.Valid {
		return entry, errors.New("Timetrap is already running")
	}

	startTimeStr := startTime.Format("2006-01-02 15:04:05.999999")

	result, err := t.db.Exec(
		`INSERT INTO entries
				(start, sheet, note)
				VALUES
				(?, ?, ?);`,
		startTimeStr, meta.CurrentSheet, note)
	if err != nil {
		panic(err.Error())
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	log.Debugf("result id = %d\n", id)

	return t.GetCurrentEntry(), nil
}

func (t *TimeTrap) Stop(stopTime time.Time) (Entry, error) {
	entry := t.GetCurrentEntry()

	if entry.End.Valid {
		return entry, errors.New(fmt.Sprintf(`No running entry on sheet "%s".`, entry.Sheet))
	}

	stopTimeStr := stopTime.Format(dateFmt)

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

	return t.GetCurrentEntry(), nil
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
		FROM entries
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

func (t *TimeTrap) DeleteSheet(sheet string) error {
	var count int
	err := t.db.QueryRow("SELECT COUNT(*) FROM entries WHERE sheet = ?;", sheet).Scan(&count)
	if err != nil {
		panic(err.Error())
	}

	if count == 0 {
		fmt.Printf("can't find \"%s\" to kill\n", sheet)
		os.Exit(1)
	}

	fmt.Printf("are you sure you want to delete %d entries on sheet \"%s\"? ", count, sheet)
	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		panic(err.Error())
	}

	if !(strings.EqualFold(input, "y") || strings.EqualFold(input, "yes")) {
		fmt.Println("will not kill")
		return nil
	}

	return t.deleteSheet(sheet)
}

func (t *TimeTrap) deleteSheet(sheet string) error {
	res, err := t.db.Exec(`DELETE FROM entries WHERE sheet = ?;`, sheet)
	if err != nil {
		panic(err.Error())
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("killed %d entries\n", rowCnt)

	return nil
}

func (t *TimeTrap) DeleteEntry(id int) error {
	var count int
	err := t.db.QueryRow("SELECT COUNT(*) FROM entries WHERE id = ?;", id).Scan(&count)
	if err != nil {
		panic(err.Error())
	}

	if count != 1 {
		fmt.Printf("couldn't find entry with id: %d\n", id)
		os.Exit(1)
	}

	fmt.Printf("are you sure you want to delete entry %d? ", id)
	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		panic(err.Error())
	}

	if !(strings.EqualFold(input, "y") || strings.EqualFold(input, "yes")) {
		fmt.Println("will not kill")
		return nil
	}

	return t.deleteEntry(id)
}

func (t *TimeTrap) deleteEntry(id int) error {
	_, err := t.db.Exec(`DELETE FROM entries WHERE id = ?;`, id)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("it's dead")

	return nil
}

func (t *TimeTrap) UpdateEntry(id int, sheet string, startTime time.Time, endTime time.Time, note string) error {
	var endTimeOrNull sql.NullString
	if !endTime.Equal(time.Time{}) {
		endTimeOrNull.String = endTime.Format(dateFmt)
		endTimeOrNull.Valid = true
	}

	var startTimeOrNull sql.NullString
	if !startTime.Equal(time.Time{}) {
		startTimeOrNull.String = startTime.Format(dateFmt)
		startTimeOrNull.Valid = true
	}

	res, err := t.db.Exec(
		"UPDATE entries SET sheet = ?, start = ?, end = ?, note = ? WHERE id = ?",
		sheet, startTimeOrNull, endTimeOrNull, note, id,
	)
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

	return nil
}
