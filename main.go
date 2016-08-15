package main

import (
	"strings"

	"github.com/admiralobvious/tinysyslog/config"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	server.Run()
}
