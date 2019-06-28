package main

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/mcuadros/go-syslog.v2"

	"tinysyslog/mutators"
	"tinysyslog/sinks"
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

	logrus.Infof("tinysyslog starting")

	err := server.Boot()
	if err != nil {
		return err
	}

	mutator := MutatorFactory()
	filter := FilterFactory()
	sinksf := SinksFactory()

	logrus.Infof("tinysyslog listening on %s", address)

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			logrus.Debugf("Received log: %v", logParts)
			log := mutators.NewLog(logParts)

			mutated, err := mutator.Mutate(log)
			logrus.Debugf("Mutated log: %v", mutated)
			if err != nil {
				logrus.Errorf("Error mutating log: %v", err)
			}

			filtered, err := filter.Filter(mutated)
			logrus.Debugf("Filtered log: %v", filtered)
			if err != nil {
				logrus.Errorf("Error filtering log: %v", err)
			}

			if len(filtered) > 0 {
				for _, sink := range sinksf {
					go write(sink, filtered) // should probably be a worker pool
				}
			}
		}
	}(channel)

	server.Wait()
	return nil
}

func write(sink sinks.Sink, msg string) {
	sinkName := sinks.GetSinkName(sink)
	if err := sink.Write([]byte(msg + "\n")); err != nil {
		logrus.Errorf("Error writing log to %s sink: %v", sinkName, err)
	} else {
		logrus.Debugf("Wrote log to %s sink", sinkName)
	}
}
