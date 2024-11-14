package util

import (
	"os"
)

func FileExists(path string) (bool, error) {
	_, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return false, err
	}

	if _, err := os.Stat(path); err != nil {
		return false, err
	}

	return true, nil
}
