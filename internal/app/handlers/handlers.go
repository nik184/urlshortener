package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/nik184/urlshortener/internal/app/storage"
)

func GenerateURL(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "bad request!", http.StatusBadRequest)
		return
	}

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

	result := config.BaseURL + "/" + storage.Set(string(url))
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(result))
}

func isURLValid(u string) bool {
	parsedURL, err := url.ParseRequestURI(u)

	return err == nil && parsedURL.Host != "" && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

func RedirectByURLID(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(rw, "bad request!", http.StatusBadRequest)
		return
	}

	id := strings.TrimLeft(r.URL.Path, "/")
	url, exists := storage.Get(id)

	if !exists {
		http.Error(rw, "wrong id was received!", http.StatusNotFound)
		return
	}

	rw.Header().Add("Location", url)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
