package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	serverArrd string `env:"SERVER_ADDRESS"`
	baseURL    string `env:"BASE_URL"`
}

var (
	ServerAddr = "localhost:8080"
	BaseURL    = "http://localhost:8080"
)

func Configure() {
	parceConf()
	parceFlag()
}

func parceConf() {
	var conf Config

	err := env.Parse(&conf)

	if err != nil {
		log.Fatal(err)
	}

	if conf.serverArrd != "" {
		ServerAddr = conf.serverArrd
	}

	if conf.baseURL != "" {
		BaseURL = conf.baseURL
	}
}

func parceFlag() {
	a := flag.String("a", "", "основной адрес сервера")
	b := flag.String("b", "", "адрес результирующего сокращенного url")

	flag.Parse()

	if *a != "" {
		ServerAddr = *a
	}

	if *b != "" {
		BaseURL = *b
	}
}
