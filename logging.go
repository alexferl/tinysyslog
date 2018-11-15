package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// InitLogging initializes the logger based on the config
func InitLogging() {
	logFile := viper.GetString("log-file")
	logFormat := viper.GetString("log-format")
	logLevel := viper.GetString("log-level")

	switch logFile {
	case "stdout":
		logrus.SetOutput(os.Stdout)
	case "stderr":
		logrus.SetOutput(os.Stderr)
	default:
		file, err := os.Create(logFile)
		if err != nil {
			logrus.Warnf("Couldn't open log-file '%s', falling back to stdout: %s", logFile, err)
			logrus.SetOutput(os.Stdout)
		} else {
			logrus.SetOutput(file)
		}

	}

	switch logFormat {
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.Warnf("Unknown log-format '%s', falling back to 'text' format.", logFormat)
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	switch logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.Warnf("Unknown log-level '%s', falling back to 'warning' level.", logLevel)
		logrus.SetLevel(logrus.WarnLevel)
	}
}
