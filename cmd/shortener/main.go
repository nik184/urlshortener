package main

import (
	"github.com/nik184/urlshortener/internal/app/config"
	"github.com/nik184/urlshortener/internal/app/database"
	"github.com/nik184/urlshortener/internal/app/server"
)

func main() {
	config.Configure()

	database.ConnectIfNeeded()
	defer database.CloseIfConnected()

	server.Start()
}
