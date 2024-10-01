package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/nik184/urlshortener/internal/app/logger"
)

type URLStorage map[string]string

var storage URLStorage
var storageLen int = 0

func InitStorage() {
	storage = URLStorage{}

	if err := ReadEvents(); err != nil {
		panic("Storage reading: error " + err.Error())
	}
}

func Set(url string) (encode string, err error) {
	if storage == nil {
		InitStorage()
	}

	encode = randStringBytes(12)
	storage[encode] = url

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

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Get(id string) (string, bool) {
	if storage == nil {
		InitStorage()
	}

	url, exists := storage[id]
	return url, exists
}

type ShortenURLRow struct {
	UUID        int    `json:"uuid"`
	ShortenURL  string `json:"shorten_url"`
	OriginalURL string `json:"original_url"`
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

	storageLen++
	row := ShortenURLRow{
		UUID:        storageLen,
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

func ReadEvents() error {
	if err := createStorageIfNotExisits(); err != nil {
		return err
	}

	file, err := os.OpenFile(config.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil
	}

	reader := bufio.NewReader(file)

	defer file.Close()

	for {
		data, err := reader.ReadBytes('\n')
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		event := ShortenURLRow{}
		if err = json.Unmarshal(data, &event); err != nil {
			return err
		}

		storage[event.ShortenURL] = event.OriginalURL
		storageLen++
	}
}

func createStorageIfNotExisits() error {
	if mkdirErr := os.MkdirAll(filepath.Dir(config.FileStoragePath), 0666); mkdirErr != nil {
		return mkdirErr
	}

	return nil
}
