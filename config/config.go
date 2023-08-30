package config

import (
	"fmt"

	libConfig "github.com/alexferl/golib/config"
	libLog "github.com/alexferl/golib/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all configuration for our program
type Config struct {
	Config            *libConfig.Config
	Logging           *libLog.Config
	BindAddr          string
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
	Addresses    []string
	IndexName    string
	Username     string
	Password     string
	CloudID      string
	APIKey       string
	ServiceToken string
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
	c := libConfig.New("TINYSYSLOG")
	c.AppName = "tinysyslog"
	c.EnvName = "PROD"
	return &Config{
		Config:   c,
		Logging:  libLog.DefaultConfig,
		BindAddr: "127.0.0.1:5140",
		ConsoleSink: ConsoleSink{
			Output: "stdout",
		},
		ElasticSearchSink: ElasticSearchSink{
			Addresses: []string{"http://127.0.0.1:9200"},
			IndexName: "tinysyslog",
		},
		FilesystemSink: FilesystemSink{
			Filename:   "syslog.log",
			MaxAge:     30,
			MaxBackups: 10,
			MaxSize:    100,
		},
		FilterType:  "",
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
}

const (
	BindAddr = "bind-addr"

	Filter      = "filter"
	FilterRegex = "filter-regex"

	Mutator = "mutator"

	Sinks = "sinks"

	SinkConsoleOutput = "sink-console-output"

	SinkElasticsearchAddresses    = "sink-elasticsearch-addresses"
	SinkElasticsearchIndexName    = "sink-elasticsearch-index-name"
	SinkElasticsearchUsername     = "sink-elasticsearch-username"
	SinkElasticsearchPassword     = "sink-elasticsearch-password"
	SinkElasticsearchCloudID      = "sink-elasticsearch-cloud-id"
	SinkElasticsearchAPIKey       = "sink-elasticsearch-api-key"
	SinkElasticsearchServiceToken = "sink-elasticsearch-service-token"

	SinkFilesystemFilename   = "sink-filesystem-filename"
	SinkFilesystemMaxAge     = "sink-filesystem-max-age"
	SinkFilesystemMaxBackups = "sink-filesystem-max-backups"
	SinkFilesystemMaxSize    = "sink-filesystem-max-size"

	SocketType = "socket-type"
)

// AddFlags adds all the flags from the command line and the config file
func (c *Config) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.BindAddr, BindAddr, c.BindAddr, "IP and port to listen on.")
	fs.StringVar(&c.FilterType, Filter, c.FilterType,
		"Filter to filter logs with. Valid filters are: regex.")
	fs.StringVar(&c.RegexFilter.Regex, FilterRegex, c.RegexFilter.Regex, "Regex to filter with.")
	fs.StringVar(&c.MutatorType, Mutator, c.MutatorType, "Mutator type to use. Valid mutators are: text, json.")
	fs.StringSliceVar(&c.SinkTypes, Sinks, c.SinkTypes, "Sinks to save syslogs to. "+
		"Valid sinks are: console, elasticsearch and filesystem.")
	fs.StringVar(&c.ConsoleSink.Output, SinkConsoleOutput, c.ConsoleSink.Output, "Console to output too. "+
		"Valid outputs are: stdout, stderr.")
	fs.StringSliceVar(&c.ElasticSearchSink.Addresses, SinkElasticsearchAddresses, c.ElasticSearchSink.Addresses,
		"Elasticsearch server address.")
	fs.StringVar(&c.ElasticSearchSink.IndexName, SinkElasticsearchIndexName, c.ElasticSearchSink.IndexName,
		"Elasticsearch index name.")
	fs.StringVar(&c.ElasticSearchSink.Username, SinkElasticsearchUsername, c.ElasticSearchSink.Username,
		"Elasticsearch username.")
	fs.StringVar(&c.ElasticSearchSink.Password, SinkElasticsearchPassword, c.ElasticSearchSink.Password,
		"Elasticsearch password.")
	fs.StringVar(&c.ElasticSearchSink.CloudID, SinkElasticsearchCloudID, c.ElasticSearchSink.CloudID,
		"Elasticsearch cloud id.")
	fs.StringVar(&c.ElasticSearchSink.APIKey, SinkElasticsearchAPIKey, c.ElasticSearchSink.APIKey,
		"Elasticsearch api key.")
	fs.StringVar(&c.ElasticSearchSink.ServiceToken, SinkElasticsearchServiceToken, c.ElasticSearchSink.ServiceToken,
		"Elasticsearch service token.")
	fs.StringVar(&c.FilesystemSink.Filename, SinkFilesystemFilename, c.FilesystemSink.Filename,
		"File to write incoming logs to.")
	fs.IntVar(&c.FilesystemSink.MaxAge, SinkFilesystemMaxAge, c.FilesystemSink.MaxAge,
		"Maximum age (in days) before a log is deleted.")
	fs.IntVar(&c.FilesystemSink.MaxBackups, SinkFilesystemMaxBackups, c.FilesystemSink.MaxBackups,
		"Maximum backups to keep.")
	fs.IntVar(&c.FilesystemSink.MaxSize, SinkFilesystemMaxSize, c.FilesystemSink.MaxSize,
		"Maximum log size (in megabytes) before it's rotated.")
	fs.StringVar(&c.SocketType, SocketType, c.SocketType, "Type of socket to use, TCP or UDP."+
		" If no type is specified, both are used.")
}

// BindFlags normalizes and parses the command line flags
func (c *Config) BindFlags() {
	if pflag.Parsed() {
		return
	}

	c.addFlags(pflag.CommandLine)
	c.Logging.BindFlags(pflag.CommandLine)

	err := c.Config.BindFlags()
	if err != nil {
		panic(fmt.Errorf("failed binding flags: %v", err))
	}

	err = libLog.New(&libLog.Config{
		LogLevel:  viper.GetString(libLog.LogLevel),
		LogOutput: viper.GetString(libLog.LogOutput),
		LogWriter: viper.GetString(libLog.LogWriter),
	})
	if err != nil {
		panic(fmt.Errorf("failed creating logger: %v", err))
	}
}
