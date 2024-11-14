package storage

import (
	"github.com/nik184/urlshortener/internal/app/logger"
)

var s stor

func Stor() stor {
	if s != nil {
		return s
	}

	if storage := createDBStor(); storage != nil {
		s = storage
	} else if storage := createFileStor(); storage != nil {
		s = storage
	} else if storage := createMapStor(); storage != nil {
		s = storage
	} else {
		panic("cannot open storage")
	}

	return s
}

func createDBStor() *DBStorage {
	storage, err := NewDBStorage()
	if err != nil {
		logger.Zl.Infoln("storage factory | ", err.Error())
		return nil
	}

	return storage
}

func createFileStor() *FileStorage {
	storage, err := NewFileStorage()

	if err != nil {
		logger.Zl.Infoln("storage factory | ", err.Error())
	}

	return storage
}

func createMapStor() *MapStorage {
	return NewMapStorage()
}

func baseSaveBanch(banch []ShortenURLRow) error {
	for _, pair := range banch {
		err := s.Set(pair.URL, pair.Short)
		if err != nil {
			return err
		}
	}

	return nil
}
