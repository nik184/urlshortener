package handlers

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/nik184/urlshortener/internal/app/database"
	"github.com/nik184/urlshortener/internal/app/storage"
)

type Req struct {
	URL string `json:"url"`
}

type Resp struct {
	Result string `json:"result"`
}

func APIGenerateURL(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	body, err := readBody(rw, r)
	if err != nil {
		http.Error(rw, "cannot read payload!", http.StatusBadRequest)
		return
	}

	req := Req{}
	decodeErr := json.Unmarshal(body, &req)
	if decodeErr != nil {
		http.Error(rw, "cannot decode json!", http.StatusBadRequest)
	}

	if !isURLValid(req.URL) {
		http.Error(rw, "incorrect url was received!", http.StatusBadRequest)
		return
	}

	hash, err := storage.Stor().Set(string(req.URL))
	if err != nil {
		http.Error(rw, "incorrect url was received!", http.StatusBadRequest)
		return
	}

	result := config.BaseURL + "/" + hash
	if !strings.Contains(config.BaseURL, "") {
		result = "http://" + result
	}

	resp := Resp{Result: result}
	encodedResp, encodeErr := json.Marshal(resp)
	if encodeErr != nil {
		http.Error(rw, "cannot encode response!", http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(encodedResp)
}

func GenerateURL(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "method not allowed!", http.StatusMethodNotAllowed)
		return
	}

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

	hash, err := storage.Stor().Set(string(url))
	if err != nil {
		http.Error(rw, "incorrect url was received!", http.StatusInternalServerError)
		return
	}

	result := config.BaseURL + "/" + hash
	if !strings.Contains(config.BaseURL, "") {
		result = "http://" + result
	}
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(result))
}

func readBody(rw http.ResponseWriter, r *http.Request) ([]byte, error) {

	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}

		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	return io.ReadAll(reader)
}

func isURLValid(u string) bool {
	parsedURL, err := url.ParseRequestURI(u)

	return err == nil && parsedURL.Host != "" && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

func RedirectByURLID(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(rw, "method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimLeft(r.URL.Path, "/")
	url, err := storage.Stor().Get(id)
	if err != nil {
		http.Error(rw, "cannot find url by id", http.StatusNotFound)
		return
	}

	rw.Header().Add("Location", url)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

func Ping(rw http.ResponseWriter, r *http.Request) {
	database.ConnectIfNeeded()

	if !database.IsConnected() {
		http.Error(rw, "cannot connect to database", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
