package main

import (
	"github.com/rs/zerolog/log"

	"tinysyslog"
	"tinysyslog/config"
)

func main() {
	cnf := config.NewConfig()
	cnf.BindFlags()

	server := tinysyslog.NewServer()
	err := server.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("error staring server")
	}
}
