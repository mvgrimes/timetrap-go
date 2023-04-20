package db

import (
	"github.com/mvgrimes/timetrap-go/internal/models"

	"errors"
	"fmt"
	"os"
	"strings"
)

func (db *Database) GetSheets() []models.Sheet {
	meta := db.GetMeta()

	results, err := db.Conn.Query(
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

	summaries := []models.Sheet{}
	for results.Next() {
		var summary models.Sheet
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

func (db *Database) DeleteSheet(sheet string) error {
	var count int
	err := db.Conn.QueryRow("SELECT COUNT(*) FROM entries WHERE sheet = ?;", sheet).Scan(&count)
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

	return db.deleteSheet(sheet)
}

func (db *Database) deleteSheet(sheet string) error {
	res, err := db.Conn.Exec(`DELETE FROM entries WHERE sheet = ?;`, sheet)
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

func (db *Database) SwitchSheet(sheet string) (models.Meta, error) {
	meta := db.GetMeta()

	if sheet == "-" {
		sheet = meta.LastSheet
	}

	if sheet == meta.CurrentSheet {
		return models.Meta{}, errors.New(fmt.Sprintf(`Already on sheet "%s"`, sheet))
	}

	res, err := db.Conn.Exec(
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

	res, err = db.Conn.Exec(
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

	return db.GetMeta(), nil
}
