package main

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"tinysyslog/config"
)

func main() {
	cnf := config.NewConfig()
	cnf.AddFlags(pflag.CommandLine)
	cnf.InitFlags()

	viper.SetEnvPrefix("tinysyslog")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	InitLogging()

	server := NewServer()
	err := server.Run()
	if err != nil {
		logrus.Fatalf("Error starting server: %v", err)
	}
}
