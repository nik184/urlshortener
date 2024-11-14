package storage

import (
	"os"
	"testing"

	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/nik184/urlshortener/internal/app/urlservice"
	"github.com/stretchr/testify/assert"
)

func TestStorages(t *testing.T) {
	mapStor := NewMapStorage()
	fileStor, _ := NewFileStorage()

	t.Run("map storage test", func(t *testing.T) {
		testSetAndGet(t, mapStor)
		testGetNonexists(t, mapStor)
	})

	t.Run("file storage test", func(t *testing.T) {
		testSetAndGet(t, fileStor)
		testGetNonexists(t, fileStor)
		testFewFileStorages(t, fileStor)
	})
}

func testSetAndGet(t *testing.T, stor stor) {
	tests := []string{
		"one", "two", "three",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			hash := urlservice.GenShort()
			setErr := stor.Set(tt, hash)
			row, getErr := stor.GetByShort(hash)

			assert.Nil(t, setErr)
			assert.Nil(t, getErr)
			assert.Equal(t, row.URL, tt)
		})
	}
}

func testGetNonexists(t *testing.T, stor stor) {
	url, err := stor.GetByShort("wrong_hash")
	assert.NotNil(t, err)
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

func testFewFileStorages(t *testing.T, stor stor) {

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

			hash := urlservice.GenShort()
			stor.Set(tt.url, hash)

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

				url, err := stor.GetByShort(pt.hash)

				assert.Nil(t, err)
				assert.Equal(t, pt.url, url.URL)
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
