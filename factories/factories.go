package factories

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"tinysyslog/config"
	"tinysyslog/constants"
	"tinysyslog/filters"
	"tinysyslog/mutators"
	"tinysyslog/sinks"
)

// Mutator creates a new object with mutators.Mutator interface
func Mutator() mutators.Mutator {
	mutator := viper.GetString(config.Mutator)

	if mutator == mutators.TextKind.String() {
		log.Debug().Msgf("using mutator '%s'", mutator)
		return mutators.NewText()
	}

	if mutator == mutators.JSONKind.String() {
		log.Debug().Msgf("using mutator '%s'", mutator)
		return mutators.NewJSON()
	}

	log.Warn().Msgf("unknown mutator '%s'. Falling back to '%s'", mutator, mutators.TextKind)
	return mutators.NewText()
}

// Filter creates a new object with filters.Filter interface
func Filter() filters.Filter {
	filter := viper.GetString(config.Filter)

	if filter == "" {
		log.Debug().Msgf("using no filtering")
		return filters.NewNoOp()
	}

	if filter == filters.RegexKind.String() {
		regex := viper.GetString(config.FilterRegex)
		log.Debug().Msgf("using regex '%s' with regular expression '%s'", filter, regex)
		return filters.NewRegex(regex)
	}

	log.Warn().Msgf("unknown filter '%s', falling back to no filtering", filter)
	return filters.NewNoOp()
}

// Sinks creates a new slice of objects with sinks.Sink interface
func Sinks() []sinks.Sink {
	sinksSlice := viper.GetStringSlice(config.Sinks)
	mutatorType := viper.GetString(config.Mutator)

	var sinksList []sinks.Sink

	for _, s := range sinksSlice {
		switch s {
		case sinks.ConsoleKind.String():
			cOutput := viper.GetString(config.SinkConsoleOutput)
			var stdOutput *os.File

			if cOutput == constants.ConsoleStdOut {
				stdOutput = os.Stdout
			} else if cOutput == constants.ConsoleStdErr {
				stdOutput = os.Stderr
			} else {
				log.Warn().Msgf("unknown console output '%s', falling back to '%s'", cOutput, constants.ConsoleStdOut)
			}
			log.Debug().Msgf("adding sink '%s'", s)
			c := sinks.NewConsole(stdOutput)
			sinksList = append(sinksList, c)
		case sinks.ElasticsearchKind.String():
			if mutatorType != mutators.JSONKind.String() {
				log.Panic().Msgf("mutator must be '%s' when using '%s' sink", mutators.JSONKind, sinks.ElasticsearchKind)
			}

			cfg := sinks.ElasticsearchConfig{
				IndexName:    viper.GetString(config.SinkElasticsearchIndexName),
				Timeout:      time.Second * 10,
				Addresses:    viper.GetStringSlice(config.SinkElasticsearchAddresses),
				Username:     viper.GetString(config.SinkElasticsearchUsername),
				Password:     viper.GetString(config.SinkElasticsearchPassword),
				CloudID:      viper.GetString(config.SinkElasticsearchCloudID),
				APIKey:       viper.GetString(config.SinkElasticsearchAPIKey),
				ServiceToken: viper.GetString(config.SinkElasticsearchServiceToken),
			}

			log.Debug().Msgf("adding sink '%s'", s)
			es := sinks.NewElasticsearch(cfg)
			sinksList = append(sinksList, es)
		case sinks.FilesystemKind.String():
			fsFilename := viper.GetString(config.SinkFilesystemFilename)
			fsMaxAge := viper.GetInt(config.SinkFilesystemMaxAge)
			fsMaxBackups := viper.GetInt(config.SinkFilesystemMaxBackups)
			fsMaxSize := viper.GetInt(config.SinkFilesystemMaxSize)

			log.Debug().Msgf("adding sink '%s'", s)
			fs := sinks.NewFilesystem(fsFilename, fsMaxAge, fsMaxBackups, fsMaxSize)
			sinksList = append(sinksList, fs)
		default:
			log.Warn().Msgf("unknown sink '%s'.", s)
		}
	}
	return sinksList
}
