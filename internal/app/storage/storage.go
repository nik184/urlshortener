package storage

import (
	"math/rand"
)

type stor interface {
	Set(url string) (string, error)
	Get(id string) (string, error)
}

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func getHash() string {
	b := make([]byte, 12)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
