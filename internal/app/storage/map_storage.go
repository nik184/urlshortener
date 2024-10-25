package storage

import "fmt"

type MapStorage struct {
	urlStorage map[string]string
}

func NewMapStorage() *MapStorage {
	s := MapStorage{
		urlStorage: make(map[string]string),
	}

	return &s
}

func (s *MapStorage) Set(url, short string) (err error) {
	s.urlStorage[short] = url

	return nil
}

func (s *MapStorage) Get(encode string) (string, error) {
	url, exists := s.urlStorage[encode]

	if !exists {
		return "", fmt.Errorf("cannot find url by id")
	}

	return url, nil
}

func (s *MapStorage) SetBatch(banch []URLWithShort) error {
	return baseSaveBanch(banch)
}
