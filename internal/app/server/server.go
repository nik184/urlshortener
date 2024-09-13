package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/nik184/urlshortener/internal/app/handlers"
)

var r = chi.NewRouter()

func Start() {
	r.Post("/", handlers.GenerateURL)
	r.Get("/{id}", handlers.RedirectByURLID)

	log.Fatal(http.ListenAndServe(config.ServerAddr, r))
}
