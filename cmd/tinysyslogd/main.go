package main

import (
	"github.com/rs/zerolog/log"

	"tinysyslog"
	"tinysyslog/config"
)

func main() {
	c := config.NewConfig()
	c.BindFlags()

	server := tinysyslog.NewServer()
	err := server.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("error staring server")
	}
}
