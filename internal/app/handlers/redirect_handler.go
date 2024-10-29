package handlers

import (
	"net/http"
	"strings"

	"github.com/nik184/urlshortener/internal/app/storage"
)

func RedirectByURLID(rw http.ResponseWriter, r *http.Request) {
	hash := strings.TrimLeft(r.URL.Path, "/")
	url, err := storage.Stor().GetByShort(hash)
	if err != nil {
		http.Error(rw, "cannot find url by id", http.StatusNotFound)
		return
	}

	rw.Header().Add("Location", url)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
