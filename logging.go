package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// InitLogging initializes the logger based on the config
func InitLogging() {
	logFile := viper.GetString("log-file")
	logFormat := viper.GetString("log-format")
	logLevel := viper.GetString("log-level")

	switch logFile {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		file, err := os.Create(logFile)
		if err != nil {
			log.Warnf("Couldn't open log-file '%s', falling back to stdout: %s", logFile, err)
			log.SetOutput(os.Stdout)
		} else {
			log.SetOutput(file)
		}

	}

	switch logFormat {
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.Warnf("Unknown log-format '%s', falling back to 'text' format.", logFormat)
		log.SetFormatter(&log.TextFormatter{})
	}

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.Warnf("Unknown log-level '%s', falling back to 'warning' level.", logLevel)
		log.SetLevel(log.WarnLevel)
	}
}
