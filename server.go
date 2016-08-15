package main

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/mcuadros/go-syslog.v2"
)

// Server holds the config
type Server struct {
}

// NewServer creates a Server instance
func NewServer() *Server {
	return &Server{}
}

// Run runs the server
func (s *Server) Run() error {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC5424)
	server.SetHandler(handler)

	address := viper.GetString("address")

	switch strings.ToLower(viper.GetString("socket-type")) {
	case "tcp":
		if err := server.ListenTCP(address); err != nil {
			log.Fatalln(err)
		}
	case "udp":
		if err := server.ListenUDP(address); err != nil {
			log.Fatalln(err)
		}
	default:
		if err := server.ListenTCP(address); err != nil {
			log.Fatalln(err)
		}
		if err := server.ListenUDP(address); err != nil {
			log.Fatalln(err)
		}
	}

	server.Boot()
	log.Infof("tinysyslog listening on %s", address)

	filter := FilterFactory()
	sink := SinkFactory()
	mutator := MutatorFactory()

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			formatted := mutator.Mutate(logParts)
			filtered := formatted
			if viper.GetString("filter-type") == "regex" {
				filtered = filter.Filter(formatted)
			}
			if len(filtered) > 0 {
				if err := sink.Write([]byte(filtered + "\n")); err != nil {
					log.Errorln(err)
				}
			}
		}
	}(channel)

	server.Wait()
	return nil
}
