package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Filename string
	Conn     *sql.DB
}

func New(filename string) *Database {
	conn, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic(err.Error())
	}

	db := Database{}
	db.Filename = filename
	db.Conn = conn

	return &db
}

func (db *Database) Close() {
	db.Conn.Close()
}

func (db *Database) Create() {
	query := "CREATE TABLE `entries` (`id` integer PRIMARY KEY AUTOINCREMENT, `note` string, `start` timestamp, `end` timestamp, `sheet` string);"
	_, err := db.Conn.Exec(query)
	if err != nil {
		panic(err.Error())
	}

	query = "CREATE TABLE `meta` (`id` integer PRIMARY KEY AUTOINCREMENT, `key` string, `value` string);"
	_, err = db.Conn.Exec(query)
	if err != nil {
		panic(err.Error())
	}
}
