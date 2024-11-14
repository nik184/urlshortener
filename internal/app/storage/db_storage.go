package storage

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

func (s *DBStorage) Set(url, short string) (err error) {
	return insertStorten(url, short, database.DB)
}

func (s *DBStorage) GetByShort(short string) (*ShortenURLRow, error) {
	row := database.DB.QueryRow(`SELECT id, url, short FROM url WHERE short = $1;`, short)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var r ShortenURLRow
	if err := row.Scan(&r.UUID, &r.URL, &r.Short); err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *DBStorage) GetByURL(url string) (*ShortenURLRow, error) {
	row := database.DB.QueryRow(`SELECT id, url, short FROM url WHERE url = $1;`, url)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var r ShortenURLRow
	if err := row.Scan(&r.UUID, &r.URL, &r.Short); err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *DBStorage) SetBatch(banch []ShortenURLRow) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	for _, pair := range banch {
		err := insertStorten(pair.URL, pair.Short, tx)

		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return err
		}
	}

	return tx.Commit()
}

func insertStorten(url, short string, db database.QueryAble) error {
	newUUID := uuid.NewV4().String()

	res, err := db.Exec(`INSERT INTO url (id, url, short) VALUES ($1, $2, $3);`, newUUID, url, short)

	var pgErr *pgconn.PgError

	if ok := errors.As(err, &pgErr); ok && pgErr.Code == pgerrcode.UniqueViolation {
		shortenURLRow, getErr := s.GetByURL(url)

		if getErr != nil {
			logger.Zl.Infoln("db get | ", getErr.Error())
			return getErr
		}

		return NewNotUniqErr(pgErr, shortenURLRow)
	}

	if err != nil {
		logger.Zl.Infoln("db insert | ", err.Error())
		return err
	}

	id, err := res.RowsAffected()

	if err != nil {
		logger.Zl.Infoln("db insert | ", err.Error())
		return err
	}

	logger.Zl.Infoln(
		"db insert | ",
		"rows created: ", id,
	)

	return nil
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
			url		VARCHAR(255) NOT NULL UNIQUE,
			short	VARCHAR(255) NOT NULL
		);
	`)

	return err
}
