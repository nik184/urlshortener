package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/nik184/urlshortener/internal/app/logger"
)

var db *sql.DB

func Q(q string) (*sql.Rows, error) {
	if err := ConnectIfNeeded(); err != nil {
		return nil, err
	}

	rows, err := db.Query(q)

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func ConnectIfNeeded() error {
	if IsConnected() {
		return nil
	}

	return connect()
}

func connect() error {
	newConnect, err := sql.Open("pgx", config.DatabaseDSN)

	if err != nil {
		logger.Zl.Errorln("sql connection error |", err.Error())
		return err
	}

	db = newConnect

	return nil
}

func IsConnected() bool {
	if db == nil {
		return false
	}

	if err := db.Ping(); err != nil {
		logger.Zl.Errorln("sql connection error |", err.Error())
		return false
	}

	return true
}

func CloseIfConnected() {
	if IsConnected() {
		if err := db.Close(); err != nil {
			logger.Zl.Errorln("sql close conn err | ", err.Error())
		}
	}
}
