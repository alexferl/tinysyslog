package main

import (
	"os"

	"github.com/admiralobvious/tinysyslog/config"
	"github.com/admiralobvious/tinysyslog/mutators"
	"github.com/admiralobvious/tinysyslog/sinks"

	log "github.com/Sirupsen/logrus"
)

// SinkFactory creates a new object with sinks.Sink interface
func SinkFactory(cnf *config.Config) sinks.Sink {
	sinkType := cnf.SinkType
	filename := cnf.Filesystem.Filename
	maxAge := cnf.Filesystem.MaxAge
	maxBackups := cnf.Filesystem.MaxBackups
	maxSize := cnf.Filesystem.MaxSize

	if sinkType == "filesystem" {
		return sinks.NewFilesystemSink(filename, maxAge, maxBackups, maxSize)
	}

	output := cnf.Console.Output
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
func MutatorFactory(cnf *config.Config) mutators.Mutator {
	mutatorType := cnf.MutatorType

	if mutatorType == "text" {
		return mutators.NewTextMutator()
	}

	if mutatorType == "json" {
		return mutators.NewJSONMutator()
	}

	log.Warningf("Unknown mutator type '%s'. Falling back to 'text'", mutatorType)
	return mutators.NewTextMutator()
}
