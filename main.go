package main

import (
	"github.com/admiralobvious/tinysyslog/config"

	"github.com/spf13/pflag"
)

func main() {
	cnf := config.NewConfig()
	cnf.AddFlags(pflag.CommandLine)
	cnf.InitFlags()
	InitLogging(cnf)

	server := NewServer(cnf)
	server.Run(pflag.CommandLine.Args())
}
