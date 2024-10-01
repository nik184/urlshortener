package storage

import (
	"os"
	"testing"

	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/stretchr/testify/assert"
)

func TestSetAndGet(t *testing.T) {
	tests := []string{
		"one", "two", "three",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			hash, _ := Set(tt)
			url, success := Get(hash)

			assert.True(t, success)
			assert.Equal(t, url, tt)
		})
	}
}

func TestGetNonexists(t *testing.T) {
	url, success := Get("wrong_hash")
	assert.False(t, success)
	assert.Empty(t, url)
}

type tc struct {
	file string
	url  string
}

type ptc struct {
	file string
	url  string
	hash string
}

func TestFewStorages(t *testing.T) {

	tests := []tc{
		{
			file: "strg.file",
			url:  "abc",
		},
		{
			file: "/tmp/mtp",
			url:  "111111111",
		},
		{
			file: "/tmp/some_file.file",
			url:  "def",
		},
		{
			file: "strg.file",
			url:  "",
		},
		{
			file: "strg.file",
			url:  "123456789",
		},
		{
			file: "qwerty",
			url:  "ghi",
		},
		{
			file: "strg.file",
			url:  "jkl",
		},
		{
			file: "/tmp/some_file.file",
			url:  "xyz",
		},
		{
			file: "/tmp/mtp",
			url:  "qazwsxedc",
		},
		{
			file: "/tmp/some_file.file",
			url:  "def",
		},
		{
			file: "qwerty",
			url:  "o_o",
		},
	}

	passedTests := []ptc{}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			config.FileStoragePath = tt.file

			hash, _ := Set(tt.url)

			if _, err := os.Stat(tt.file); err != nil {
				panic("file was not created")
			}

			newPt := ptc{
				file: tt.file,
				url:  tt.url,
				hash: hash,
			}

			passedTests = append(passedTests, newPt)

			for _, pt := range passedTests {
				config.FileStoragePath = pt.file

				url, success := Get(pt.hash)

				assert.True(t, success)
				assert.Equal(t, url, pt.url)
			}
		})
	}

	clearFiles(passedTests)
}

func clearFiles(passedTests []ptc) {
	for _, pt := range passedTests {
		os.Remove(pt.file)
	}
}
