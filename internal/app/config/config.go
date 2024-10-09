package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/nik184/urlshortener/internal/app/logger"
)

type Config struct {
	ServerArrd      string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

var (
	ServerAddr      = "localhost:8080"
	BaseURL         = "http://localhost:8080"
	FileStoragePath = "storage.tmpstorage"
	DatabaseDSN     = "postgres://urlshortener:urlshortener@localhost:5433/urlshortener"
	StorageDriver   = "file"
)

func Configure() {
	parceConf()
	parceFlag()
}

func parceConf() {
	var conf Config

	err := env.Parse(&conf)

	logger.Zl.Infoln("env |", conf)

	if err != nil {
		log.Fatal(err)
	}

	if conf.ServerArrd != "" {
		ServerAddr = conf.ServerArrd
	}

	if conf.BaseURL != "" {
		BaseURL = conf.BaseURL
	}

	if conf.FileStoragePath != "" {
		FileStoragePath = conf.FileStoragePath
	}

	if conf.DatabaseDSN != "" {
		DatabaseDSN = conf.DatabaseDSN
	}
}

func parceFlag() {
	a := flag.String("a", "", "основной адрес сервера")
	b := flag.String("b", "", "адрес результирующего сокращенного url")
	f := flag.String("f", "", "адрес файла - хранилища сокращенных url адресов")
	d := flag.String("d", "", "адрес для подключения к базе данных")

	flag.Parse()

	logger.Zl.Infoln(
		"flags |",
		"a: ", *a,
		"b: ", *b,
		"f: ", *f,
		"d: ", *d,
	)

	if *a != "" {
		ServerAddr = *a
	}

	if *b != "" {
		BaseURL = *b
	}

	if *f != "" {
		FileStoragePath = *f
	}

	if *d != "" {
		DatabaseDSN = *d
	}
}
