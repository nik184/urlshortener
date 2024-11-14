package handlers

import (
	"errors"
	"net/http"

	"github.com/nik184/urlshortener/internal/app/storage"
	"github.com/nik184/urlshortener/internal/app/urlservice"
)

func ShortURL(rw http.ResponseWriter, r *http.Request) {
	status := http.StatusCreated

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

	short := urlservice.GenShort()
	if err = storage.Stor().Set(string(url), short); err != nil {
		var notUniqErr *storage.NotUniqErr
		if ok := errors.As(err, &notUniqErr); ok {
			row, err := storage.Stor().GetByURL(url)
			short = row.Short

			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			} else {
				status = http.StatusConflict
			}
		} else {
			http.Error(rw, "failed to save url: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	result := concatPathToAddr(short)
	rw.WriteHeader(status)
	rw.Write([]byte(result))
}
