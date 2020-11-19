package config

import (
	"fmt"
	"strings"

	xlog "github.com/alexferl/x/log"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all configuration for our program
type Config struct {
	BindAddress       string
	ConsoleSink       ConsoleSink
	ElasticSearchSink ElasticSearchSink
	FilesystemSink    FilesystemSink
	FilterType        string
	GrokFilter        GrokFilter
	Logging           xlog.Config
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
	Address            string
	IndexName          string
	Username           string
	Password           string
	InsecureSkipVerify bool
	DisableSniffing    bool
}

// FilesystemSink holds all configuration for the FilesystemSink sink
type FilesystemSink struct {
	Filename     string
	MaxAge       int
	MaxBackups   int
	MaxSize      int
	OutputFormat string
}

// GrokFilter holds grok configuration
type GrokFilter struct {
	Fields  []string
	Pattern string
}

// RegexFilter holds regex configuration
type RegexFilter struct {
	Regex string
}

// NewConfig creates a Config instance
func NewConfig() *Config {
	cnf := Config{
		BindAddress: "127.0.0.1:5140",
		ConsoleSink: ConsoleSink{
			Output: "stdout",
		},
		ElasticSearchSink: ElasticSearchSink{
			Address:            "http://127.0.0.1:9200",
			IndexName:          "tinysyslog",
			InsecureSkipVerify: false,
			DisableSniffing:    false,
		},
		FilesystemSink: FilesystemSink{
			Filename:   "syslog.log",
			MaxAge:     30,
			MaxBackups: 10,
			MaxSize:    100,
		},
		FilterType: "null",
		GrokFilter: GrokFilter{
			Fields:  []string{},
			Pattern: "",
		},
		Logging:     xlog.NewConfig(),
		MutatorType: "text",
		RegexFilter: RegexFilter{
			Regex: "",
		},
		SinkTypes:  []string{"console"},
		SocketType: "",
	}
	return &cnf
}

// addFlags adds all the flags from the command line and the config file
func (cnf *Config) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cnf.BindAddress, "bind-address", cnf.BindAddress, "IP and port to listen on.")
	fs.StringVar(&cnf.FilterType, "filter", cnf.FilterType,
		"Filter to filter logs with. Valid filters are: 'null', 'regex'. Null doesn't do any filtering.")
	fs.StringSliceVar(&cnf.GrokFilter.Fields, "filter-grok-fields", cnf.GrokFilter.Fields,
		"Grok fields to keep.")
	fs.StringVar(&cnf.GrokFilter.Pattern, "filter-grok-pattern", cnf.GrokFilter.Pattern,
		"Grok pattern to filter with.")
	fs.StringVar(&cnf.RegexFilter.Regex, "filter-regex", cnf.RegexFilter.Regex, "Regex to filter with.")
	fs.StringVar(&cnf.MutatorType, "mutator", cnf.MutatorType,
		"Mutator type to use. Valid mutators are: 'text', 'json'.")
	fs.StringSliceVar(&cnf.SinkTypes, "sinks", cnf.SinkTypes,
		"Sinks to save syslogs to. Valid sinks are: 'console', 'elasticsearch', 'filesystem'.")
	fs.StringVar(&cnf.ConsoleSink.Output, "sink-console-output", cnf.ConsoleSink.Output,
		"Console to output to. Valid outputs are: 'stdout', 'stderr'.")
	fs.StringVar(&cnf.ElasticSearchSink.Address, "sink-elasticsearch-address", cnf.ElasticSearchSink.Address,
		"Elasticsearch server address.")
	fs.StringVar(&cnf.ElasticSearchSink.IndexName, "sink-elasticsearch-index-name",
		cnf.ElasticSearchSink.IndexName, "Elasticsearch index name.")
	fs.StringVar(&cnf.ElasticSearchSink.Username, "sink-elasticsearch-username", cnf.ElasticSearchSink.Username,
		"Elasticsearch username.")
	fs.StringVar(&cnf.ElasticSearchSink.Password, "sink-elasticsearch-password", cnf.ElasticSearchSink.Password,
		"Elasticsearch password.")
	fs.BoolVar(&cnf.ElasticSearchSink.InsecureSkipVerify, "sink-elasticsearch-insecure-skip-verify",
		cnf.ElasticSearchSink.InsecureSkipVerify, "Elasticsearch skip verifying TLS certificates.")
	fs.BoolVar(&cnf.ElasticSearchSink.DisableSniffing, "sink-elasticsearch-disable-sniffing",
		cnf.ElasticSearchSink.DisableSniffing, "Elasticsearch disable sniffing process.")
	fs.StringVar(&cnf.FilesystemSink.Filename, "sink-filesystem-filename", cnf.FilesystemSink.Filename,
		"File to write incoming logs to.")
	fs.IntVar(&cnf.FilesystemSink.MaxAge, "sink-filesystem-max-age", cnf.FilesystemSink.MaxAge,
		"Maximum age (in days) before a log is deleted.")
	fs.IntVar(&cnf.FilesystemSink.MaxBackups, "sink-filesystem-max-backups", cnf.FilesystemSink.MaxBackups,
		"Maximum backups to keep.")
	fs.IntVar(&cnf.FilesystemSink.MaxSize, "sink-filesystem-max-size", cnf.FilesystemSink.MaxSize,
		"Maximum log size (in megabytes) before it's rotated.")
	fs.StringVar(&cnf.SocketType, "socket-type", cnf.SocketType,
		"Type of socket to use, TCP or UDP. If no type is specified, both are used.")
}

// wordSepNormalizeFunc changes all flags that contain "_" separators
func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}

// BindFlags normalizes and parses the command line flags
func (cnf *Config) BindFlags() {
	cnf.Logging.AddFlags(pflag.CommandLine)

	cnf.addFlags(pflag.CommandLine)
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		m := fmt.Sprintf("Error binding flags: '%v'", err)
		log.Panic().Msg(m)
	}

	pflag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
	pflag.Parse()

	viper.SetEnvPrefix("tinysyslog")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
}
