package handlers

import (
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok && pgErr.Code == pgerrcode.UniqueViolation {
			short, err = storage.Stor().GetByURL(url)

			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			} else {
				status = http.StatusConflict
			}
		} else {
			http.Error(rw, "failed to save url!", http.StatusInternalServerError)
			return
		}
	}

	result := concatPathToAddr(short)
	rw.WriteHeader(status)
	rw.Write([]byte(result))
}
