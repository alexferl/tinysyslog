package tinysyslog

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

// MutatorFactory creates a new object with mutators.Mutator interface
func MutatorFactory() mutators.Mutator {
	mutatorType := viper.GetString(config.Mutator)

	if mutatorType == constants.MutatorText {
		log.Debug().Msgf("using mutator '%s'", mutatorType)
		return mutators.NewText()
	}

	if mutatorType == constants.MutatorJSON {
		log.Debug().Msgf("using mutator '%s'", mutatorType)
		return mutators.NewJSON()
	}

	log.Warn().Msgf("unknown mutator '%s'. Falling back to '%s'", mutatorType, constants.MutatorText)
	return mutators.NewText()
}

// FilterFactory creates a new object with filters.Filter interface
func FilterFactory() filters.Filter {
	filterType := viper.GetString(config.Filter)

	if filterType == "" {
		log.Debug().Msgf("using no filtering")
		return filters.NewNoOp()
	}

	if filterType == constants.FilterRegex {
		filter := viper.GetString(config.FilterRegex)
		log.Debug().Msgf("using filter '%s' with regular expression '%s'", filterType, filter)
		return filters.NewRegex(filter)
	}

	log.Warn().Msgf("unknown filter '%s', falling back to no filtering", filterType)
	return filters.NewNoOp()
}

// SinksFactory creates a new slice of objects with sinks.Sink interface
func SinksFactory() []sinks.Sink {
	sinkTypes := viper.GetStringSlice(config.Sinks)
	mutatorType := viper.GetString(config.Mutator)

	var sinksList []sinks.Sink

	for _, sink := range sinkTypes {
		switch sink {
		case constants.SinkConsole:
			cOutput := viper.GetString(config.SinkConsoleOutput)
			var stdOutput *os.File

			if cOutput == constants.ConsoleStdOut {
				stdOutput = os.Stdout
			} else if cOutput == constants.ConsoleStdErr {
				stdOutput = os.Stderr
			} else {
				log.Warn().Msgf("unknown console output '%s', falling back to '%s'", cOutput, constants.ConsoleStdOut)
			}
			log.Debug().Msgf("adding sink '%s'", sink)
			c := sinks.NewConsole(stdOutput)
			sinksList = append(sinksList, c)
		case constants.SinkElasticsearch:
			if mutatorType != constants.MutatorJSON {
				log.Panic().Msgf("mutator must be '%s' when using '%s' sink", constants.MutatorJSON, constants.SinkElasticsearch)
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

			log.Debug().Msgf("adding sink '%s'", sink)
			es := sinks.NewElasticsearch(cfg)
			sinksList = append(sinksList, es)
		case constants.SinkFilesystem:
			fsFilename := viper.GetString(config.SinkFilesystemFilename)
			fsMaxAge := viper.GetInt(config.SinkFilesystemMaxAge)
			fsMaxBackups := viper.GetInt(config.SinkFilesystemMaxBackups)
			fsMaxSize := viper.GetInt(config.SinkFilesystemMaxSize)

			log.Debug().Msgf("adding sink '%s'", sink)
			fs := sinks.NewFilesystem(fsFilename, fsMaxAge, fsMaxBackups, fsMaxSize)
			sinksList = append(sinksList, fs)
		default:
			log.Warn().Msgf("unknown sink '%s'.", sink)
		}
	}
	return sinksList
}
