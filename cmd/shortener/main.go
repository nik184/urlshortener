package main

import (
	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/nik184/urlshortener/internal/app/server"
)

func main() {
	config.Configure()

	server.Start()
}
