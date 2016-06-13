package config

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all configuration for our program
type Config struct {
	Address     string
	Console     Console
	Filesystem  Filesystem
	LogFile     string
	LogFormat   string
	LogLevel    string
	MutatorType string
	SinkType    string
	SocketType  string
}

// Filesystem holds all configuration for the filesystem sink
type Filesystem struct {
	Filename     string
	MaxAge       int
	MaxBackups   int
	MaxSize      int
	OutputFormat string
}

// Console holds all configuration for the console sink
type Console struct {
	Output string
}

// NewConfig creates a Config instance
func NewConfig() *Config {
	cnf := Config{
		Address: "127.0.0.1:5140",
		Console: Console{
			Output: "stdout",
		},
		Filesystem: Filesystem{
			Filename:   "syslog.log",
			MaxAge:     30,
			MaxBackups: 10,
			MaxSize:    100,
		},
		LogFile:     "tinysyslog.log",
		LogFormat:   "text",
		LogLevel:    "info",
		MutatorType: "text",
		SinkType:    "filesystem",
		SocketType:  "",
	}
	return &cnf
}

// AddFlags adds all the flags from the command line and the config file
func (cnf *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cnf.Address, "address", cnf.Address, "IP and port to listen on.")
	fs.StringVar(&cnf.Console.Output, "console-output", cnf.Console.Output, "Console to output too. "+
		"Valid outputs are: stdout, stderr.")
	fs.StringVar(&cnf.Filesystem.Filename, "filesystem-filename", cnf.Filesystem.Filename, "File to write incoming logs to.")
	fs.IntVar(&cnf.Filesystem.MaxAge, "filesystem-max-age", cnf.Filesystem.MaxAge,
		"Maximum age (in days) before a log is deleted. Set to '0' to disable.")
	fs.IntVar(&cnf.Filesystem.MaxBackups, "filesystem-max-backups", cnf.Filesystem.MaxBackups,
		"Maximum backups to keep. Set to '0' to disable.")
	fs.IntVar(&cnf.Filesystem.MaxSize, "filesystem-max-size", cnf.Filesystem.MaxSize,
		"Maximum log size (in megabytes) before it's rotated.")
	fs.StringVar(&cnf.LogFile, "log-file", cnf.LogFile, "The log file to write to. "+
		"'stdout' means log to stdout and 'stderr' means log to stderr.")
	fs.StringVar(&cnf.LogFormat, "log-format", cnf.LogFormat,
		"The log format. Valid format values are: text, json.")
	fs.StringVar(&cnf.LogLevel, "log-level", cnf.LogLevel, "The granularity of log outputs. "+
		"Valid level names are: debug, info, warning, error and critical.")
	fs.StringVar(&cnf.MutatorType, "mutator-type", cnf.MutatorType, "Mutator to transform logs as.")
	fs.StringVar(&cnf.SinkType, "sink-type", cnf.SinkType, "Sink to save logs to.")
	fs.StringVar(&cnf.SocketType, "socket-type", cnf.SocketType, "Type of socket to use, TCP or UDP."+
		" If no type is specified, both are used.")
}

// wordSepNormalizeFunc changes all flags that contain "_" separators
func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}

// InitFlags normalizes and parses the command line flags
func (cnf *Config) InitFlags() {
	viper.BindPFlags(pflag.CommandLine)
	pflag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
	pflag.Parse()
}
