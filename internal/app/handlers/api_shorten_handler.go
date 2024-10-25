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
	hash, err := parceURLReq(rw, r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	result := concatPathToAddr(hash)

	resp := URLResp{Result: result}
	encodedResp, encodeErr := json.Marshal(resp)
	if encodeErr != nil {
		http.Error(rw, "cannot encode response!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(encodedResp)
}

func parceURLReq(rw http.ResponseWriter, r *http.Request) (hash string, err error) {
	body, err := readBody(rw, r)
	if err != nil {
		return
	}

	req := URLReq{}
	if err = json.Unmarshal(body, &req); err != nil {
		return
	}

	if !isURLValid(req.URL) {
		err = errors.New("incorrect url was received")
		return
	}

	hash = urlservice.GenShort()
	err = storage.Stor().Set(string(req.URL), hash)

	return
}
