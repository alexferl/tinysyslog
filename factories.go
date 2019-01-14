package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/admiralobvious/tinysyslog/filters"
	"github.com/admiralobvious/tinysyslog/mutators"
	"github.com/admiralobvious/tinysyslog/sinks"
)

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
				logrus.Warningf("Unknown console output type '%s'. Falling back to 'stdout'", cOutput)
			}
			logrus.Debugf("Adding sink type '%s'", sink)
			cs := sinks.NewConsoleSink(stdOutput)
			sinksList = append(sinksList, cs)
		case "elasticsearch":
			if mutatorType != "json" {
				m := fmt.Sprint("Mutator must be 'json' when using 'elasticsearch' sink")
				logrus.Panic(m)
				panic(m)
			}

			esAddress := viper.GetString("sink-elasticsearch-address")
			esIndexName := viper.GetString("sink-elasticsearch-index-name")

			logrus.Debugf("Adding sink type '%s'", sink)
			es := sinks.NewElasticsearchSink(esAddress, esIndexName)
			sinksList = append(sinksList, es)
		case "filesystem":
			fsFilename := viper.GetString("sink-filesystem-filename")
			fsMaxAge := viper.GetInt("sink-filesystem-max-age")
			fsMaxBackups := viper.GetInt("sink-filesystem-max-backups")
			fsMaxSize := viper.GetInt("sink-filesystem-max-size")

			logrus.Debugf("Adding sink type '%s'", sink)
			fs := sinks.NewFilesystemSink(fsFilename, fsMaxAge, fsMaxBackups, fsMaxSize)
			sinksList = append(sinksList, fs)
		default:
			logrus.Warningf("Unknown sink type '%s'.", sink)
		}
	}
	return sinksList
}
