package pkg

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"tinysyslog/internal/pkg/filters"
	"tinysyslog/internal/pkg/mutators"
	"tinysyslog/internal/pkg/sinks"
)

// MutatorFactory creates a new object with mutators.Mutator interface
func MutatorFactory() mutators.Mutator {
	mutatorType := viper.GetString("mutator")

	if mutatorType == "text" {
		log.Debug().Msgf("Using mutator type '%s'", mutatorType)
		return mutators.NewTextMutator()
	}

	if mutatorType == "json" {
		log.Debug().Msgf("Using mutator type '%s'", mutatorType)
		return mutators.NewJSONMutator()
	}

	log.Warn().Msgf("Unknown mutator type '%s'. Falling back to 'text'", mutatorType)
	return mutators.NewTextMutator()
}

// FilterFactory creates a new object with filters.Filter interface
func FilterFactory() filters.Filter {
	filterType := viper.GetString("filter")

	if filterType == "" || filterType == "null" {
		log.Debug().Msgf("Using filter type '%s'", filterType)
		return filters.NewNullFilter()
	}

	if filterType == "regex" {
		filter := viper.GetString("filter-regex-filter")
		log.Debug().Msgf("Using filter type '%s' with filter '%s'", filterType, filter)
		return filters.NewRegexFilter(filter)
	}

	if filterType == "grok" {
		pattern := viper.GetString("filter-grok-pattern")
		fields := viper.GetStringSlice("filter-grok-fields")
		log.Debug().Msgf("Using filter type '%s' with pattern '%s'", filterType, pattern)
		return filters.NewGrokFilter(pattern, fields)
	}

	log.Warn().Msgf("Unknown filter type '%s'. Falling back to 'null'", filterType)
	return filters.NewNullFilter()
}

// SinksFactory creates a new slice of objects with sinks.Sink interface
func SinksFactory() []sinks.Sink {
	sinkTypes := viper.GetStringSlice("sinks")
	mutatorType := viper.GetString("mutator")

	var sinksList []sinks.Sink

	for _, sink := range sinkTypes {
		switch sink {
		case "console":
			cOutput := viper.GetString("sink-console-output")
			var stdOutput *os.File

			if cOutput == "stdout" {
				stdOutput = os.Stdout
			} else if cOutput == "stderr" {
				stdOutput = os.Stderr
			} else {
				log.Warn().Msgf("Unknown console output type '%s'. Falling back to 'stdout'", cOutput)
			}
			log.Debug().Msgf("Adding sink type '%s'", sink)
			cs := sinks.NewConsoleSink(stdOutput)
			sinksList = append(sinksList, cs)
		case "elasticsearch":
			if mutatorType != "json" {
				m := fmt.Sprint("Mutator must be 'json' when using 'elasticsearch' sink")
				log.Panic().Msg(m)
			}

			esAddress := viper.GetString("sink-elasticsearch-address")
			esIndexName := viper.GetString("sink-elasticsearch-index-name")
			esUsername := viper.GetString("sink-elasticsearch-username")
			esPassword := viper.GetString("sink-elasticsearch-password")
			esInsecure := viper.GetBool("sink-elasticsearch-insecure-skip-verify")
			esSniff := viper.GetBool("sink-elasticsearch-disable-sniffing")

			log.Debug().Msgf("Adding sink type '%s'", sink)
			es := sinks.NewElasticsearchSink(esAddress, esIndexName, esUsername, esPassword, esInsecure, esSniff)
			sinksList = append(sinksList, es)
		case "filesystem":
			fsFilename := viper.GetString("sink-filesystem-filename")
			fsMaxAge := viper.GetInt("sink-filesystem-max-age")
			fsMaxBackups := viper.GetInt("sink-filesystem-max-backups")
			fsMaxSize := viper.GetInt("sink-filesystem-max-size")

			log.Debug().Msgf("Adding sink type '%s'", sink)
			fs := sinks.NewFilesystemSink(fsFilename, fsMaxAge, fsMaxBackups, fsMaxSize)
			sinksList = append(sinksList, fs)
		default:
			log.Warn().Msgf("Unknown sink type '%s'.", sink)
		}
	}
	return sinksList
}
