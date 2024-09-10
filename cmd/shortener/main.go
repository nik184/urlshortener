package main

import (
	"log"
	"net/http"

	"github.com/nik184/urlshortener/internal/app/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.MainHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
