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

	// if (config.MainAddr.Port != config.RedirAddr.Port) || (config.MainAddr.Host != config.RedirAddr.Host) {
	// 	startMultiserver()
	// } else {
	// 	log.Fatal(http.ListenAndServe(config.MainAddr.AddrWithOnlyPort(), r))
	// }

	log.Fatal(http.ListenAndServe(config.MainAddr.AddrWithOnlyPort(), r))

}

func startMultiserver() {
	finish := make(chan bool)

	go func() { log.Fatal(http.ListenAndServe(config.MainAddr.AddrWithOnlyPort(), r)) }()
	// go func() { log.Fatal(http.ListenAndServe(config.RedirAddr.AddrWithOnlyPort(), r)) }()

	<-finish
}
