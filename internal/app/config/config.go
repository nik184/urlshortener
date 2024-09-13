package config

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

type Adds struct {
	Host string
	Port int
}

var (
	MainAddr  = Adds{Host: "localhost", Port: 8080}
	RedirAddr = "http://localhost:8080"
)

func (a Adds) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a Adds) AddrWithOnlyPort() string {
	return ":" + strconv.Itoa(a.Port)
}

func (a *Adds) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("тебуемый формат - host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = port
	return nil
}

func init() {
	flag.Var(&MainAddr, "a", "основной адрес сервера")
	flag.StringVar(&RedirAddr, "b", "http://localhost:8080", "адрес результирующего сокращенного url")
}

func ParceFlags() {
	flag.Parse()
}
