package handlers

import (
	"net/http"

	"github.com/nik184/urlshortener/internal/app/storage"
	"github.com/nik184/urlshortener/internal/app/urlservice"
)

func ShortURL(rw http.ResponseWriter, r *http.Request) {
	body, err := readBody(rw, r)
	if err != nil {
		http.Error(rw, "cannot read payload!", http.StatusBadRequest)
		return
	}

	url := string(body)
	if !isURLValid(url) {
		http.Error(rw, "incorrect url was received!", http.StatusBadRequest)
		return
	}

	hash := urlservice.GenShort()
	if err = storage.Stor().Set(string(url), hash); err != nil {
		http.Error(rw, "failed to save url!", http.StatusInternalServerError)
		return
	}

	result := concatPathToAddr(hash)
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(result))
}
