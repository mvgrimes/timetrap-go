package db

import (
	"github.com/mvgrimes/timetrap-go/internal/models"

	log "github.com/sirupsen/logrus"
)

func (db *Database) GetMeta() models.Meta {
	meta := models.Meta{}

	err := db.Conn.QueryRow("select value from meta where id = 1;").Scan(&meta.CurrentSheet)
	if err != nil {
		panic(err.Error())
	}
	err = db.Conn.QueryRow("select value from meta where id = 2;").Scan(&meta.LastSheet)
	if err != nil {
		panic(err.Error())
	}
	err = db.Conn.QueryRow("select value from meta where id = 3;").Scan(&meta.LastCheckout)
	if err != nil {
		panic(err.Error())
	}

	log.Debugf("meta: %v", meta)

	return meta
}
