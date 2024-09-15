package storage

import (
	"math/rand"
)

type UrlsStorage map[string]string

var storage UrlsStorage = UrlsStorage{}

func Set(url string) (encode string) {
	encode = randStringBytes(12)
	storage[encode] = url

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
	url, exists := storage[id]
	return url, exists
}
