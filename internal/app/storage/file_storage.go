package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/nik184/urlshortener/internal/app/logger"

	uuid "github.com/satori/go.uuid"
)

type FileStorage struct {
}

type ShortenURLRow struct {
	UUID        string `json:"uuid"`
	ShortenURL  string `json:"shorten_url"`
	OriginalURL string `json:"original_url"`
}

func NewFileStorage() (FileStorage, error) {
	s := FileStorage{}

	return s, nil
}

func (s *FileStorage) Set(url string) (encode string, err error) {
	encode = getHash()

	err = SaveToStorage(url, encode)
	if err != nil {
		logger.Zl.Errorln("save to storage |",
			"url:", url,
			"encode:", encode,
			"file:", config.FileStoragePath,
			"error:", err.Error(),
		)
	}

	return
}

func (s *FileStorage) Get(id string) (string, error) {
	file, err := os.OpenFile(config.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(file)

	defer file.Close()

	for {
		data, err := reader.ReadBytes('\n')

		if err == io.EOF {
			return "", err
		}

		if err != nil {
			return "", err
		}

		event := ShortenURLRow{}
		if err = json.Unmarshal(data, &event); err != nil {
			return "", err
		}

		if event.ShortenURL == id {
			return event.OriginalURL, nil
		}
	}
}

func SaveToStorage(url string, enc string) error {
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
		UUID:        newUUID,
		ShortenURL:  enc,
		OriginalURL: url,
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
