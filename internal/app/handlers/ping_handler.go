package handlers

import (
	"net/http"

	"github.com/nik184/urlshortener/internal/app/database"
)

func Ping(rw http.ResponseWriter, r *http.Request) {
	database.ConnectIfNeeded()

	if !database.IsConnected() {
		http.Error(rw, "cannot connect to database", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
