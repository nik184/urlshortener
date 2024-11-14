package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nik184/urlshortener/internal/app/logger"
	"github.com/nik184/urlshortener/internal/app/storage"
	"github.com/nik184/urlshortener/internal/app/urlservice"
)

type BanchReqElem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BanchReq []BanchReqElem

type BanchRespElem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BanchResp []BanchRespElem

func APIProcessBatchOfURLs(rw http.ResponseWriter, r *http.Request) {
	resp, err := parceBanchReq(rw, r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	encodedResp, encodeErr := json.Marshal(resp)
	if encodeErr != nil {
		http.Error(rw, "cannot encode response!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(encodedResp)
}

func parceBanchReq(rw http.ResponseWriter, r *http.Request) (*BanchResp, error) {
	body, err := readBody(rw, r)
	if err != nil {
		return nil, err
	}

	reqs := BanchReq{}
	if err = json.Unmarshal(body, &reqs); err != nil {
		return nil, err
	}

	resp := BanchResp{}
	banch := []storage.ShortenURLRow{}
	for _, req := range reqs {
		logger.Zl.Info(
			"parce batch | ",
			"correlationID: ", req.CorrelationID,
			"originalURL: ", req.OriginalURL,
		)

		if !isURLValid(req.OriginalURL) {
			err = errors.New("incorrect url was received")
			return nil, err
		}

		hash := urlservice.GenShort()
		banch = append(banch, storage.ShortenURLRow{
			URL:   req.OriginalURL,
			Short: hash,
		})

		shortURL := concatPathToAddr(hash)

		resp = append(resp, BanchRespElem{
			CorrelationID: req.CorrelationID,
			ShortURL:      shortURL,
		})
	}

	if len(banch) >= 1 {
		err := storage.Stor().SetBatch(banch)

		if err != nil {
			logger.Zl.Infoln(
				"save batch | ",
				"error: ", err.Error(),
			)
		}
	}

	return &resp, nil
}
