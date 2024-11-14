package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/nik184/urlshortener/internal/app/logger"
	"github.com/nik184/urlshortener/internal/app/util"

	uuid "github.com/satori/go.uuid"
)

type FileStorage struct {
}

func NewFileStorage() (*FileStorage, error) {
	if config.FileStoragePath == "" {
		return nil, errors.New("file storage credantials was not given")
	}

	if err := createStorageIfNotExisits(); err != nil {
		return nil, err
	}

	if exists, err := util.FileExists(config.FileStoragePath); err != nil || !exists {
		return nil, err
	}

	s := FileStorage{}

	return &s, nil
}

func (s *FileStorage) Set(url, short string) (err error) {
	err = saveToStorage(url, short)
	if err != nil {
		logger.Zl.Errorln("save to storage | ",
			"url:", url,
			"short:", short,
			"file:", config.FileStoragePath,
			"error:", err.Error(),
		)
	}

	return
}

func (s *FileStorage) GetByShort(short string) (*ShortenURLRow, error) {
	file, err := os.OpenFile(config.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	defer file.Close()

	for {
		data, err := reader.ReadBytes('\n')

		if err == io.EOF {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		shortenURLRow := ShortenURLRow{}
		if err = json.Unmarshal(data, &shortenURLRow); err != nil {
			return nil, err
		}

		if shortenURLRow.Short == short {
			return &shortenURLRow, nil
		}
	}
}

func (s *FileStorage) GetByURL(url string) (*ShortenURLRow, error) {
	file, err := os.OpenFile(config.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	defer file.Close()

	for {
		data, err := reader.ReadBytes('\n')

		if err == io.EOF {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		event := ShortenURLRow{}
		if err = json.Unmarshal(data, &event); err != nil {
			return nil, err
		}

		if event.URL == url {
			return &event, nil
		}
	}
}

func (s *FileStorage) SetBatch(banch []ShortenURLRow) error {
	return baseSaveBanch(banch)
}

func saveToStorage(url string, enc string) error {
	if err := createStorageIfNotExisits(); err != nil {
		return err
	}

	file, fopErr := os.OpenFile(config.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if fopErr != nil {
		return fopErr
	}
	writer := bufio.NewWriter(file)

	newUUID := uuid.NewV4().String()

	row := ShortenURLRow{
		UUID:  newUUID,
		Short: enc,
		URL:   url,
	}

	jsonRaw, encErr := json.Marshal(row)

	if encErr != nil {
		return encErr
	}

	if _, writeErr := writer.Write(jsonRaw); writeErr != nil {
		return writeErr
	}

	if err := writer.WriteByte('\n'); err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

func createStorageIfNotExisits() error {
	if mkdirErr := os.MkdirAll(filepath.Dir(config.FileStoragePath), 0666); mkdirErr != nil {
		return mkdirErr
	}

	return nil
}

// func checkFile() error {
// 	_, err := os.OpenFile(config.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

// 	if err != nil {
// 		return err
// 	}

// 	if _, err := os.Stat(config.FileStoragePath); err != nil {
// 		return err
// 	}

// 	return nil
// }
