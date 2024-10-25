package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/nik184/urlshortener/internal/app/handlers"
	"github.com/nik184/urlshortener/internal/app/middlewares"
)

var r = chi.NewRouter()

func Start() {
	r.Use(middlewares.Logger)

	r.Post("/", handlers.ShortURL)
	r.Get("/{id}", handlers.RedirectByURLID)
	r.Get("/ping", handlers.Ping)

	r.Post("/api/shorten", handlers.APIShortURL)
	r.Post("/api/shorten/batch", handlers.APIProcessBatchOfURLs)

	log.Fatal(http.ListenAndServe(config.ServerAddr, r))
}
