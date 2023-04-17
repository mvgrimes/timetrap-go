package models

import (
	"database/sql"
)

type Entry struct {
	ID    int          `json:"id"`
	Sheet string       `json:"sheet"`
	Start sql.NullTime `json:"start"`
	End   sql.NullTime `json:"end"`
	Note  string       `json:"note"`
}
