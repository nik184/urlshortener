package storage

import "github.com/nik184/urlshortener/internal/app/config"

type storDriver string

const (
	MapStor  storDriver = `map`
	FileStor storDriver = `file`
)

var s stor

func Stor() stor {
	if s == nil {
		switch storDriver(config.StorageDriver) {
		case MapStor:
			st := NewMapStorage()
			s = &st
		case FileStor:
			st, _ := NewFileStorage()
			s = &st
		}
	}

	return s
}
