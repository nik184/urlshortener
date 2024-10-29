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

func (s *MapStorage) GetByShort(short string) (string, error) {
	url, exists := s.urlStorage[short]

	if !exists {
		return "", fmt.Errorf("cannot find url by id")
	}

	return url, nil
}

func (s *MapStorage) GetByURL(url string) (string, error) {
	for k, v := range s.urlStorage {
		if v == url {
			return k, nil
		}
	}
	return "", fmt.Errorf("cannot find url by id")
}

func (s *MapStorage) SetBatch(banch []URLWithShort) error {
	return baseSaveBanch(banch)
}
