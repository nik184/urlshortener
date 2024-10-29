package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nik184/urlshortener/internal/app/storage"
	"github.com/nik184/urlshortener/internal/app/urlservice"
)

type URLReq struct {
	URL string `json:"url"`
}

type URLResp struct {
	Result string `json:"result"`
}

func APIShortURL(rw http.ResponseWriter, r *http.Request) {
	status := http.StatusCreated

	body, err := readBody(rw, r)
	if err != nil {
		return
	}

	req := URLReq{}
	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if !isURLValid(req.URL) {
		http.Error(rw, "incorrect url was received", http.StatusBadRequest)
		return
	}

	short := urlservice.GenShort()
	err = storage.Stor().Set(string(req.URL), short)

	var pgErr *pgconn.PgError
	if ok := errors.As(err, &pgErr); ok && pgErr.Code == pgerrcode.UniqueViolation {
		short, err = storage.Stor().GetByURL(req.URL)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		} else {
			status = http.StatusConflict
		}
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	result := concatPathToAddr(short)

	resp := URLResp{Result: result}
	encodedResp, encodeErr := json.Marshal(resp)
	if encodeErr != nil {
		http.Error(rw, "cannot encode response!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(encodedResp)
}
