package tinysyslogd

import (
	"strings"

	xlog "github.com/alexferl/x/log"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/mcuadros/go-syslog.v2"

	"tinysyslog/internal/pkg"
	"tinysyslog/internal/pkg/config"
	"tinysyslog/internal/pkg/mutators"
	"tinysyslog/internal/pkg/sinks"
)

func init() {
	c := config.NewConfig()
	c.BindFlags()
	lc := xlog.Config{
		LogLevel:  viper.GetString("log-level"),
		LogOutput: viper.GetString("log-output"),
		LogWriter: viper.GetString("log-writer"),
	}
	err := xlog.Init(lc)
	if err != nil {
		log.Panic().Msgf("Panic initializing logger: '%v'", err)
	}
}

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

	address := viper.GetString("bind-address")

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

	log.Info().Msgf("tinysyslog starting")

	err := server.Boot()
	if err != nil {
		return err
	}

	mutator := pkg.MutatorFactory()
	filter := pkg.FilterFactory()
	sinksf := pkg.SinksFactory()

	log.Info().Msgf("tinysyslog listening on '%s'", address)

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			log.Debug().Msgf("Received log: '%v'", logParts)
			l := mutators.NewLog(logParts)

			mutated, err := mutator.Mutate(l)
			log.Debug().Msgf("Mutated log: '%v'", mutated)
			if err != nil {
				log.Error().Msgf("Error mutating log: '%v'", err)
			}

			filtered, err := filter.Filter(mutated)
			log.Debug().Msgf("Filtered log: '%v'", filtered)
			if err != nil {
				log.Error().Msgf("Error filtering log: '%v'", err)
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
		log.Error().Msgf("Error writing log to '%s' sink: '%v'", sinkName, err)
	} else {
		log.Debug().Msgf("Wrote log to '%s' sink", sinkName)
	}
}
