package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nik184/urlshortener/internal/app/handlers"
)

const Host = "http://localhost"
const Port = ":8080"

func Start() {
	r := chi.NewRouter()

	r.Post("/", handlers.GetMainHadler(Host, Port))
	r.Get("/{id}", handlers.GetMainHadler(Host, Port))

	log.Fatal(http.ListenAndServe(Port, r))
}
