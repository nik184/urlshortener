package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

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

	if err = storage.Stor().Set(string(req.URL), short); err != nil {
		var notUniqErr *storage.NotUniqErr
		if ok := errors.As(err, &notUniqErr); ok {
			row, err := storage.Stor().GetByURL(req.URL)
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
