package main

import (
	"strings"

	"github.com/sirupsen/logrus"
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
			return err
		}
	case "udp":
		if err := server.ListenUDP(address); err != nil {
			return err
		}
	default:
		if err := server.ListenTCP(address); err != nil {
			return err
		}
		if err := server.ListenUDP(address); err != nil {
			return err
		}
	}

	err := server.Boot()
	if err != nil {
		return err
	}
	logrus.Infof("tinysyslog listening on %s", address)

	filter := FilterFactory()
	sink := SinkFactory()
	mutator := MutatorFactory()

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			formatted, err := mutator.Mutate(logParts)
			if err != nil {
				logrus.Errorf("Error mutating log: %v", err)
			}
			filtered, err := filter.Filter(formatted)
			if err != nil {
				logrus.Errorf("Error filtering log: %v", err)
			}
			if len(filtered) > 0 {
				if err := sink.Write([]byte(filtered + "\n")); err != nil {
					logrus.Errorf("Error writing log: %v", err)
				}
			}
		}
	}(channel)

	server.Wait()
	return nil
}
