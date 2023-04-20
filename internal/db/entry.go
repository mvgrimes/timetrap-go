package db

import (
	"github.com/mvgrimes/timetrap-go/internal/models"

	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"

	"database/sql"
)

const dateFmt = "2006-01-02 15:04:05.999999"

func (db *Database) GetEntry(id int) models.Entry {
	entry := models.Entry{}

	err := db.Conn.QueryRow(`SELECT id, sheet, start, end, note
						FROM entries
						WHERE id = ?
						ORDER BY id DESC
						LIMIT 1;`,
		id,
	).Scan(&entry.ID, &entry.Sheet, &entry.Start, &entry.End, &entry.Note)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Entry{}
		}
		panic(err.Error())
	}

	log.Debugf("entry: %v", entry)
	return entry
}

func (db *Database) GetCurrentEntry() models.Entry {
	meta := db.GetMeta()
	entry := models.Entry{}

	err := db.Conn.QueryRow(`SELECT id, sheet, start, end, note
						FROM entries
						WHERE sheet = ?
						ORDER BY id DESC
						LIMIT 1;`,
		meta.CurrentSheet,
	).Scan(&entry.ID, &entry.Sheet, &entry.Start, &entry.End, &entry.Note)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Entry{}
		}
		panic(err.Error())
	}

	log.Debugf("entry: %v", entry)
	return entry
}

func (db *Database) GetEntries(sheet string) []models.Entry {
	results, err := db.Conn.Query(
		`SELECT id, sheet, start, end, note
				FROM entries
				WHERE sheet = ?;`,
		sheet)
	if err != nil {
		panic(err.Error())
	}

	entries := []models.Entry{}
	for results.Next() {
		var entry models.Entry
		err = results.Scan(&entry.ID, &entry.Sheet, &entry.Start, &entry.End, &entry.Note)
		if err != nil {
			panic(err.Error())
		}
		entries = append(entries, entry)
	}

	return entries
}

func (db *Database) GetFilteredEntries(sheet string, start sql.NullTime, end sql.NullTime, grep string) []models.Entry {
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

	results, err := db.Conn.Query(sql, params...)
	if err != nil {
		panic(err.Error())
	}

	entries := []models.Entry{}
	for results.Next() {
		var entry models.Entry
		err = results.Scan(&entry.ID, &entry.Sheet, &entry.Start, &entry.End, &entry.Note)
		if err != nil {
			panic(err.Error())
		}
		entries = append(entries, entry)
	}

	return entries
}

func (db *Database) DeleteEntry(id int) error {
	var count int
	err := db.Conn.QueryRow("SELECT COUNT(*) FROM entries WHERE id = ?;", id).Scan(&count)
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

	return db.deleteEntry(id)
}

func (db *Database) deleteEntry(id int) error {
	_, err := db.Conn.Exec(`DELETE FROM entries WHERE id = ?;`, id)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("it's dead")

	return nil
}

func (db *Database) UpdateEntry(id int, sheet string, startTime time.Time, endTime time.Time, note string) error {
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

	res, err := db.Conn.Exec(
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
