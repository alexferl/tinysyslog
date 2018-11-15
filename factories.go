package main

import (
	"os"

	"github.com/admiralobvious/tinysyslog/filters"
	"github.com/admiralobvious/tinysyslog/mutators"
	"github.com/admiralobvious/tinysyslog/sinks"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// SinkFactory creates a new object with sinks.Sink interface
func SinkFactory() sinks.Sink {
	sinkType := viper.GetString("sink")
	filename := viper.GetString("sink-filesystem-filename")
	maxAge := viper.GetInt("sink-filesystem-max-age")
	maxBackups := viper.GetInt("sink-filesystem-max-backups")
	maxSize := viper.GetInt("sink-filesystem-max-size")

	if sinkType == "filesystem" {
		logrus.Debugf("Using sink type '%s'", sinkType)
		return sinks.NewFilesystemSink(filename, maxAge, maxBackups, maxSize)
	}

	output := viper.GetString("sink-console-output")
	var stdOutput *os.File

	if sinkType == "console" {
		if output == "stdout" {
			stdOutput = os.Stdout
		} else if output == "stderr" {
			stdOutput = os.Stderr
		} else {
			logrus.Warningf("Unknown console output type '%s'. Falling back to 'stdout'", output)
		}
		logrus.Debugf("Using sink type '%s'", sinkType)
		return sinks.NewConsoleSink(stdOutput)
	}

	logrus.Warningf("Unknown sink type '%s'. Falling back to 'filesystem'", sinkType)
	return sinks.NewFilesystemSink(filename, maxAge, maxBackups, maxSize)
}

// MutatorFactory creates a new object with mutators.Mutator interface
func MutatorFactory() mutators.Mutator {
	mutatorType := viper.GetString("mutator")

	if mutatorType == "text" {
		logrus.Debugf("Using mutator type '%s'", mutatorType)
		return mutators.NewTextMutator()
	}

	if mutatorType == "json" {
		logrus.Debugf("Using mutator type '%s'", mutatorType)
		return mutators.NewJSONMutator()
	}

	logrus.Warningf("Unknown mutator type '%s'. Falling back to 'text'", mutatorType)
	return mutators.NewTextMutator()
}

// FilterFactory creates a new object with filters.Filter interface
func FilterFactory() filters.Filter {
	filterType := viper.GetString("filter")

	if filterType == "" || filterType == "null" {
		logrus.Debugf("Using filter type '%s'", filterType)
		return filters.NewNullFilter()
	}

	if filterType == "regex" {
		filter := viper.GetString("filter-regex-filter")
		logrus.Debugf("Using filter type '%s' with filter '%s'", filterType, filter)
		return filters.NewRegexFilter(filter)
	}

	logrus.Warningf("Unknown filter type '%s'. Falling back to 'null'", filterType)
	return filters.NewNullFilter()
}
