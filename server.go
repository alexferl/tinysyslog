package main

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/mcuadros/go-syslog.v2"

	"github.com/admiralobvious/tinysyslog/mutators"
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
	sinks := SinksFactory()
	mutator := MutatorFactory()

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			log := makeLog(logParts)
			formatted, err := mutator.Mutate(log)
			if err != nil {
				logrus.Errorf("Error mutating log: %v", err)
			}
			filtered, err := filter.Filter(formatted)
			if err != nil {
				logrus.Errorf("Error filtering log: %v", err)
			}
			if len(filtered) > 0 {
				for _, sink := range sinks {
					if err := sink.Write([]byte(filtered + "\n")); err != nil {
						logrus.Errorf("Error writing log: %v", err)
					}
				}
			}
		}
	}(channel)

	server.Wait()
	return nil
}

func makeLog(logParts map[string]interface{}) mutators.Log {
	return mutators.Log{
		Timestamp: logParts["timestamp"].(time.Time),
		Hostname:  logParts["hostname"].(string),
		AppName:   logParts["app_name"].(string),
		ProcId:    logParts["proc_id"].(string),
		Severity:  logParts["severity"].(int),
		Message:   logParts["message"].(string),
	}
}
