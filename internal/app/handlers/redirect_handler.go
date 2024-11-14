package handlers

import (
	"net/http"
	"strings"

	"github.com/nik184/urlshortener/internal/app/logger"
	"github.com/nik184/urlshortener/internal/app/storage"
)

func RedirectByURLID(rw http.ResponseWriter, r *http.Request) {
	hash := strings.TrimLeft(r.URL.Path, "/")
	row, err := storage.Stor().GetByShort(hash)
	if err != nil {
		logger.Zl.Error("redirect handler 404 | ", err.Error())
		http.Error(rw, "cannot find url by id", http.StatusNotFound)
		return
	}

	rw.Header().Add("Location", row.URL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
