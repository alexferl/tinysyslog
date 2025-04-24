package config

import (
	"fmt"

	libConfig "github.com/alexferl/golib/config"
	libLog "github.com/alexferl/golib/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/alexferl/tinysyslog/constants"
	"github.com/alexferl/tinysyslog/filters"
	"github.com/alexferl/tinysyslog/mutators"
	"github.com/alexferl/tinysyslog/sinks"
)

// Config holds all configuration for our program
type Config struct {
	Config         *libConfig.Config
	Logging        *libLog.Config
	BindAddr       string
	ConsoleSink    ConsoleSink
	FilesystemSink FilesystemSink
	FilterType     string
	LogFile        string
	LogFormat      string
	LogLevel       string
	MutatorType    string
	RegexFilter    RegexFilter
	SinkTypes      []string
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

// New creates a Config instance
func New() *Config {
	c := libConfig.New("TINYSYSLOG")
	c.AppName = "tinysyslog"
	c.EnvName = "PROD"
	return &Config{
		Config:   c,
		Logging:  libLog.DefaultConfig,
		BindAddr: "127.0.0.1:5140",
		ConsoleSink: ConsoleSink{
			Output: constants.ConsoleStdOut,
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
		MutatorType: mutators.TextKind.String(),
		RegexFilter: RegexFilter{
			Regex: "",
		},
		SinkTypes:  []string{sinks.ConsoleKind.String()},
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
		fmt.Sprintf("Filter to filter logs with. Valid filters: %s", filters.Kinds),
	)
	fs.StringVar(&c.RegexFilter.Regex, FilterRegex, c.RegexFilter.Regex, "Regex to filter with.")
	fs.StringVar(&c.MutatorType, Mutator, c.MutatorType,
		fmt.Sprintf("Mutator type to use. Valid mutators: %s", mutators.Kinds),
	)
	fs.StringSliceVar(&c.SinkTypes, Sinks, c.SinkTypes,
		fmt.Sprintf("Sinks to save syslogs to. Valid sinks: %s", sinks.Kinds),
	)
	fs.StringVar(&c.ConsoleSink.Output, SinkConsoleOutput, c.ConsoleSink.Output,
		fmt.Sprintf("Console to output to. Valid outputs: %s", constants.ConsoleOutputs))
	fs.StringVar(&c.FilesystemSink.Filename, SinkFilesystemFilename, c.FilesystemSink.Filename,
		"File path to write incoming logs to.")
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
