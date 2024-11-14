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

func (s *MapStorage) GetByShort(short string) (*ShortenURLRow, error) {
	url, exists := s.urlStorage[short]

	if !exists {
		return nil, fmt.Errorf("cannot find url by id")
	}

	return &ShortenURLRow{
		URL:   url,
		Short: short,
	}, nil
}

func (s *MapStorage) GetByURL(url string) (*ShortenURLRow, error) {
	for _, v := range s.urlStorage {
		if v == url {
			return &ShortenURLRow{
				URL:   url,
				Short: v,
			}, nil
		}
	}

	return nil, fmt.Errorf("cannot find url by id")
}

func (s *MapStorage) SetBatch(banch []ShortenURLRow) error {
	return baseSaveBanch(banch)
}
