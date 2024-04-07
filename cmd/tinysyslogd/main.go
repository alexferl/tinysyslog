package main

import (
	"github.com/rs/zerolog/log"

	"tinysyslog/config"
	"tinysyslog/server"
)

func main() {
	c := config.New()
	c.BindFlags()

	s, err := server.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed staring server")
	}

	s.Run()
}
