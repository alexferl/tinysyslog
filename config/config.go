package config

import (
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all configuration for our program
type Config struct {
	Address           string
	ConsoleSink       ConsoleSink
	ElasticSearchSink ElasticSearchSink
	FilesystemSink    FilesystemSink
	FilterType        string
	LogFile           string
	LogFormat         string
	LogLevel          string
	MutatorType       string
	RegexFilter       RegexFilter
	SinkTypes         []string
	SocketType        string
}

// ConsoleSink holds all configuration for the ConsoleSink sink
type ConsoleSink struct {
	Output string
}

type ElasticSearchSink struct {
	Address   string
	IndexName string
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
		ElasticSearchSink: ElasticSearchSink{
			Address:   "http://127.0.0.1:9200",
			IndexName: "tinysyslog",
		},
		FilesystemSink: FilesystemSink{
			Filename:   "syslog.log",
			MaxAge:     30,
			MaxBackups: 10,
			MaxSize:    100,
		},
		FilterType:  "null",
		LogFile:     "stdout",
		LogFormat:   "text",
		LogLevel:    "info",
		MutatorType: "text",
		RegexFilter: RegexFilter{
			Regex: "",
		},
		SinkTypes:  []string{"console"},
		SocketType: "",
	}
	return &cnf
}

// AddFlags adds all the flags from the command line and the config file
func (cnf *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cnf.Address, "address", cnf.Address, "IP and port to listen on.")
	fs.StringVar(&cnf.FilterType, "filter", cnf.FilterType, "Filter to filter logs with. Valid filters are: null and regex. "+
		"Null doesn't do any filtering.")
	fs.StringVar(&cnf.RegexFilter.Regex, "filter-regex", cnf.RegexFilter.Regex, "Regex to filter with.")
	fs.StringVar(&cnf.LogFile, "log-file", cnf.LogFile, "The log file to write to. "+
		"'stdout' means log to stdout and 'stderr' means log to stderr.")
	fs.StringVar(&cnf.LogFormat, "log-format", cnf.LogFormat, "The log format. Valid format values are: text, json.")
	fs.StringVar(&cnf.LogLevel, "log-level", cnf.LogLevel, "The granularity of log outputs. "+
		"Valid level names are: debug, info, warning, error and critical.")
	fs.StringVar(&cnf.MutatorType, "mutator", cnf.MutatorType, "Mutator type to use. Valid mutators are: text, json.")
	fs.StringSliceVar(&cnf.SinkTypes, "sinks", cnf.SinkTypes, "Sinks to save syslogs to. Valid sinks are: console, elasticsearch and filesystem.")
	fs.StringVar(&cnf.ConsoleSink.Output, "sink-console-output", cnf.ConsoleSink.Output, "Console to output too. "+
		"Valid outputs are: stdout, stderr.")
	fs.StringVar(&cnf.ElasticSearchSink.Address, "sink-elasticsearch-address", cnf.ElasticSearchSink.Address, "Elasticsearch server address.")
	fs.StringVar(&cnf.ElasticSearchSink.IndexName, "sink-elasticsearch-index-name", cnf.ElasticSearchSink.IndexName, "Elasticsearch index name.")
	fs.StringVar(&cnf.FilesystemSink.Filename, "sink-filesystem-filename", cnf.FilesystemSink.Filename, "File to write incoming logs to.")
	fs.IntVar(&cnf.FilesystemSink.MaxAge, "sink-filesystem-max-age", cnf.FilesystemSink.MaxAge,
		"Maximum age (in days) before a log is deleted.")
	fs.IntVar(&cnf.FilesystemSink.MaxBackups, "sink-filesystem-max-backups", cnf.FilesystemSink.MaxBackups, "Maximum backups to keep.")
	fs.IntVar(&cnf.FilesystemSink.MaxSize, "sink-filesystem-max-size", cnf.FilesystemSink.MaxSize,
		"Maximum log size (in megabytes) before it's rotated.")
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
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		logrus.Fatalf("Error binding flags: %v", err)
		panic(err)
	}

	pflag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
	pflag.Parse()
}
