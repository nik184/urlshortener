package storage

import (
	"errors"

	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/nik184/urlshortener/internal/app/database"
	"github.com/nik184/urlshortener/internal/app/logger"
	uuid "github.com/satori/go.uuid"
)

type DBStorage struct {
}

func NewDBStorage() (*DBStorage, error) {
	if config.DatabaseDSN == "" {
		return nil, errors.New("db credantials was not given")
	}

	if err := database.ConnectIfNeeded(); err != nil {
		return nil, err
	}

	if !database.IsConnected() {
		return nil, errors.New("connection unsuccessful")
	}

	if err := prepareDB(); err != nil {
		return nil, err
	}

	s := DBStorage{}

	return &s, nil
}

func (s *DBStorage) Set(url string) (string, error) {
	encode := getHash()
	newUUID := uuid.NewV4().String()

	res, err := database.DB.Exec(`INSERT INTO url (id, url, encode) VALUES ($1, $2, $3);`, newUUID, url, encode)

	if err != nil {
		logger.Zl.Infoln("db insert | ", err.Error())
		return "", err
	}

	id, err := res.RowsAffected()

	if err != nil {
		logger.Zl.Infoln("db insert | ", err.Error())
		return "", err
	}

	logger.Zl.Infoln(
		"db insert | ",
		"rows created: ", id,
	)

	return encode, nil
}

func (s *DBStorage) Get(encode string) (string, error) {
	row := database.DB.QueryRow(`SELECT url FROM url WHERE encode = $1;`, encode)

	if row.Err() != nil {
		return "", row.Err()
	}

	var url string
	if err := row.Scan(&url); err != nil {
		return "", err
	}

	return url, nil
}

func prepareDB() error {
	if isURLTablesExisis() {
		return nil
	}

	return createTable()
}

func isURLTablesExisis() bool {
	row := database.DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'urls') AS exists;`)

	var exists bool
	row.Scan(&exists)

	return exists
}

func createTable() error {
	_, err := database.DB.Exec(`
		CREATE TABLE IF NOT EXISTS url (
			id		VARCHAR(255) NOT NULL PRIMARY KEY,
			url		VARCHAR(255) NOT NULL,
			encode	VARCHAR(255) NOT NULL
		);
	`)

	return err
}
