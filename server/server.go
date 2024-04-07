package server

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/mcuadros/go-syslog.v2"

	"tinysyslog/config"
	"tinysyslog/factories"
	"tinysyslog/mutators"
	"tinysyslog/sinks"
)

// Server holds the config
type Server struct {
	server *syslog.Server
}

// New creates a Server instance
func New() (*Server, error) {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC5424)
	server.SetHandler(handler)

	address := viper.GetString(config.BindAddr)
	switch strings.ToLower(viper.GetString(config.SocketType)) {
	case "tcp":
		if err := server.ListenTCP(address); err != nil {
			return nil, err
		}
	case "udp":
		if err := server.ListenUDP(address); err != nil {
			return nil, err
		}
	default:
		if err := server.ListenTCP(address); err != nil {
			return nil, err
		}
		if err := server.ListenUDP(address); err != nil {
			return nil, err
		}
	}

	log.Info().Msg("tinysyslog starting")

	err := server.Boot()
	if err != nil {
		return nil, err
	}

	mutator := factories.Mutator()
	filter := factories.Filter()
	sinksSlice := factories.Sinks()

	log.Info().Msgf("tinysyslog listening on %s", address)

	go func(ch syslog.LogPartsChannel) {
		for logParts := range ch {
			log.Debug().Msgf("received log: %v", logParts)
			newLog := mutators.NewLog(logParts)

			mutated, err := mutator.Mutate(newLog)
			if err != nil {
				log.Err(err).Msg("failed mutating log")
			} else {
				log.Debug().Msgf("mutated log: %v", mutated)
			}

			filtered := mutated
			if viper.GetString(config.Filter) != "" {
				filtered, err = filter.Filter(mutated)
				if err != nil {
					log.Err(err).Msg("failed filtering log")
				} else {
					log.Debug().Msgf("filtered log: %v", filtered)
				}
			}

			if len(filtered) > 0 {
				for _, sink := range sinksSlice {
					go write(sink, filtered)
				}
			}
		}
	}(channel)

	return &Server{server: server}, nil
}

// Run runs the server
func (s *Server) Run() {
	s.server.Wait()
}

func write(sink sinks.Sink, msg string) {
	sinkName := sink.GetKind().String()
	if err := sink.Write([]byte(msg + "\n")); err != nil {
		log.Err(err).Str("sink", sinkName).Msgf("failed writing log to sink: %s", sinkName)
	} else {
		log.Debug().Msgf("wrote log to %s sink", sinkName)
	}
}
