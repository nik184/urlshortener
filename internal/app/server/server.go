package server

import (
	"log"
	"net/http"

	"github.com/nik184/urlshortener/internal/app/handlers"
)

const Host = "http://localhost"
const Port = ":8080"

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.GetMainHadler(Host, Port))

	log.Fatal(http.ListenAndServe(Port, mux))
}
