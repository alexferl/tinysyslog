package main

import (
	"os"

	"github.com/admiralobvious/tinysyslog/filters"
	"github.com/admiralobvious/tinysyslog/mutators"
	"github.com/admiralobvious/tinysyslog/sinks"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// SinkFactory creates a new object with sinks.Sink interface
func SinkFactory() sinks.Sink {
	sinkType := viper.GetString("sink-type")
	filename := viper.GetString("filesystem-filename")
	maxAge := viper.GetInt("filesystem-max-age")
	maxBackups := viper.GetInt("filesystem-max-backups")
	maxSize := viper.GetInt("filesystem-max-size")

	if sinkType == "filesystem" {
		return sinks.NewFilesystemSink(filename, maxAge, maxBackups, maxSize)
	}

	output := viper.GetString("console-output")
	var stdOutput *os.File

	if sinkType == "console" {
		if output == "stdout" {
			stdOutput = os.Stdout
		} else if output == "stderr" {
			stdOutput = os.Stderr
		} else {
			log.Warningf("Unknown console output type '%s'. Falling back to 'stdout'", output)
		}
		return sinks.NewConsoleSink(stdOutput)
	}

	log.Warningf("Unknown sink type '%s'. Falling back to 'filesystem'", sinkType)
	return sinks.NewFilesystemSink(filename, maxAge, maxBackups, maxSize)
}

// MutatorFactory creates a new object with mutators.Mutator interface
func MutatorFactory() mutators.Mutator {
	mutatorType := viper.GetString("mutator-type")

	if mutatorType == "text" {
		return mutators.NewTextMutator()
	}

	if mutatorType == "json" {
		return mutators.NewJSONMutator()
	}

	log.Warningf("Unknown mutator type '%s'. Falling back to 'text'", mutatorType)
	return mutators.NewTextMutator()
}

// FilterFactory creates a new object with filters.Filter interface
func FilterFactory() filters.Filter {
	filterType := viper.GetString("filter-type")

	if filterType == "regex" {
		filter := viper.GetString("regex-filter")
		return filters.NewRegexFilter(filter)
	}

	log.Warningf("Unknown filter type '%s'. Falling back to 'regex'", filterType)
	return filters.NewRegexFilter("")
}
