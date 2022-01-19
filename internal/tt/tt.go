package tt

import (
	"errors"
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

func getDatabaseConn(TimetrapDB string) *sql.DB {
	db, err := sql.Open("sqlite3", TimetrapDB)
	if err != nil {
		panic(err.Error())
	}

	return db
}

func GetMeta(databaseFile string) Meta {
	db := getDatabaseConn(databaseFile)
	defer db.Close()

	meta := Meta{}

	err := db.QueryRow("select value from meta where id = 1;").Scan(&meta.CurrentSheet)
	if err != nil {
		panic(err.Error())
	}
	err = db.QueryRow("select value from meta where id = 2;").Scan(&meta.LastSheet)
	if err != nil {
		panic(err.Error())
	}
	err = db.QueryRow("select value from meta where id = 3;").Scan(&meta.LastCheckout)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("meta: %v", meta)

	return meta
}

func GetCurrentEntry(databaseFile string) Entry {
	db := getDatabaseConn(databaseFile)
	defer db.Close()

	meta := GetMeta(databaseFile)
	entry := Entry{}

	err := db.QueryRow(`SELECT id, sheet, start, end, note
						FROM entries
						WHERE sheet = ?
						ORDER BY id DESC
						LIMIT 1;`,
		meta.CurrentSheet,
	).Scan(&entry.ID, &entry.Sheet, &entry.Start, &entry.End, &entry.Note)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("entries: %v", entry)

	return entry
}

func GetEntries(databaseFile string, sheet string) []Entry {
	db := getDatabaseConn(databaseFile)
	defer db.Close()

	results, err := db.Query(
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

func Start(databaseFile string, sheet string, startTime time.Time, note string) {
	db := getDatabaseConn(databaseFile)
	defer db.Close()

	startTimeStr := startTime.Format("2006-01-02 15:04:05.999999")

	result, err := db.Exec(
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
	log.Printf("result id = %d\n", id)
}

func Stop(databaseFile string, sheet string, stopTime time.Time) error {
	db := getDatabaseConn(databaseFile)
	defer db.Close()

	entry := GetCurrentEntry(databaseFile)

	if entry.End.Valid {
		return errors.New("No running...")
	}

	stopTimeStr := stopTime.Format("2006-01-02 15:04:05.999999")

	update, err := db.Query(
		`UPDATE entries SET end = ?
		   WHERE id = ?;`,
		stopTimeStr, entry.ID)
	if err != nil {
		panic(err.Error())
	}
	defer update.Close()

	return nil
}
