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

func (s *MapStorage) Set(url string) (string, error) {
	encode := getHash()
	s.urlStorage[encode] = url

	return encode, nil
}

func (s *MapStorage) Get(encode string) (string, error) {
	url, exists := s.urlStorage[encode]

	if !exists {
		return "", fmt.Errorf("cannot find url by id")
	}

	return url, nil
}
