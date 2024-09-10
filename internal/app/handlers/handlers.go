package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/nik184/urlshortener/internal/app/storage"
)

var host string
var port string

func GetMainHadler(h string, p string) func(http.ResponseWriter, *http.Request) {
	host = h
	port = p

	return mainHandler
}

func mainHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		generateURL(rw, r)
	case http.MethodGet:
		redirectByURLID(rw, r)
	default:
		http.Error(rw, "bad request!", http.StatusBadRequest)
	}
}

func generateURL(rw http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "cannot read payload!", http.StatusBadRequest)
		return
	}

	url := string(body)
	if !isURLValid(url) {
		http.Error(rw, "incorrect url was received!", http.StatusBadRequest)
		return
	}

	result := host + port + "/" + storage.Set(string(url))
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(result))
}

func isURLValid(u string) bool {
	parsedURL, err := url.ParseRequestURI(u)

	return err == nil && parsedURL.Host != "" && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

func redirectByURLID(rw http.ResponseWriter, r *http.Request) {
	id := strings.TrimLeft(r.URL.Path, "/")
	url, exists := storage.Get(id)

	if !exists {
		http.Error(rw, "wrong id was received!", http.StatusNotFound)
		return
	}

	rw.Header().Add("Location", url)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
