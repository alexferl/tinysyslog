package config

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all configuration for our program
type Config struct {
	Address        string
	ConsoleSink    ConsoleSink
	FilesystemSink FilesystemSink
	FilterType     string
	LogFile        string
	LogFormat      string
	LogLevel       string
	MutatorType    string
	RegexFilter    RegexFilter
	SinkType       string
	SocketType     string
}

// ConsoleSink holds all configuration for the ConsoleSink sink
type ConsoleSink struct {
	Output string
}

// FilesystemSink holds all configuration for the FilesystemSink sink
type FilesystemSink struct {
	Filename     string
	MaxAge       int
	MaxBackups   int
	MaxSize      int
	OutputFormat string
}

// RegexFilter holds regex configuration
type RegexFilter struct {
	Regex string
}

// NewConfig creates a Config instance
func NewConfig() *Config {
	cnf := Config{
		Address: "127.0.0.1:5140",
		ConsoleSink: ConsoleSink{
			Output: "stdout",
		},
		FilesystemSink: FilesystemSink{
			Filename:   "syslog.log",
			MaxAge:     30,
			MaxBackups: 10,
			MaxSize:    100,
		},
		FilterType:  "regex",
		LogFile:     "tinysyslog.log",
		LogFormat:   "text",
		LogLevel:    "info",
		MutatorType: "text",
		RegexFilter: RegexFilter{
			Regex: "",
		},
		SinkType:   "filesystem",
		SocketType: "",
	}
	return &cnf
}

// AddFlags adds all the flags from the command line and the config file
func (cnf *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cnf.Address, "address", cnf.Address, "IP and port to listen on.")
	fs.StringVar(&cnf.ConsoleSink.Output, "console-output", cnf.ConsoleSink.Output, "Console to output too. "+
		"Valid outputs are: stdout, stderr.")
	fs.StringVar(&cnf.FilesystemSink.Filename, "filesystem-filename", cnf.FilesystemSink.Filename, "File to write incoming logs to.")
	fs.IntVar(&cnf.FilesystemSink.MaxAge, "filesystem-max-age", cnf.FilesystemSink.MaxAge,
		"Maximum age (in days) before a log is deleted.")
	fs.IntVar(&cnf.FilesystemSink.MaxBackups, "filesystem-max-backups", cnf.FilesystemSink.MaxBackups, "Maximum backups to keep.")
	fs.IntVar(&cnf.FilesystemSink.MaxSize, "filesystem-max-size", cnf.FilesystemSink.MaxSize,
		"Maximum log size (in megabytes) before it's rotated.")
	fs.StringVar(&cnf.FilterType, "filter-type", cnf.FilterType, "Filter to filter logs with. Valid filters are: regex.")
	fs.StringVar(&cnf.LogFile, "log-file", cnf.LogFile, "The log file to write to. "+
		"'stdout' means log to stdout and 'stderr' means log to stderr.")
	fs.StringVar(&cnf.LogFormat, "log-format", cnf.LogFormat, "The log format. Valid format values are: text, json.")
	fs.StringVar(&cnf.LogLevel, "log-level", cnf.LogLevel, "The granularity of log outputs. "+
		"Valid level names are: debug, info, warning, error and critical.")
	fs.StringVar(&cnf.MutatorType, "mutator-type", cnf.MutatorType, "Mutator type to use. Valid mutators are: text, json.")
	fs.StringVar(&cnf.RegexFilter.Regex, "regex-filter", cnf.RegexFilter.Regex, "Regex to filter with. No filtering by default.")
	fs.StringVar(&cnf.SinkType, "sink-type", cnf.SinkType, "Sink to save logs to. Valid sinks are: console, filesystem.")
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
