package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/alexferl/tinysyslog/constants"
	"github.com/alexferl/tinysyslog/filters"
	"github.com/alexferl/tinysyslog/mutators"
	"github.com/alexferl/tinysyslog/sinks"
)

const (
	// LogLevel is the flag name for log level
	LogLevel = "log-level"
	// LogOutput is the flag name for log output
	LogOutput = "log-output"
	// LogWriter is the flag name for log writer
	LogWriter = "log-writer"
)

// Config holds all configuration for our program
type Config struct {
	BindAddr       string
	ConsoleSink    ConsoleSink
	FilesystemSink FilesystemSink
	FilterType     string
	LogFile        string
	LogFormat      string
	LogLevel       string
	LogOutput      string
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
	return &Config{
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
		LogFile:     "",
		LogFormat:   "text",
		LogLevel:    "info",
		LogOutput:   "stdout",
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

// addFlags adds all the flags from the command line and the config file
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

	fs.StringVar(&c.LogLevel, LogLevel, c.LogLevel, "Log level (debug, info, warn, error, fatal, panic).")
	fs.StringVar(&c.LogOutput, LogOutput, c.LogOutput, "Log output (stdout, stderr, file).")
	fs.StringVar(&c.LogFormat, LogWriter, c.LogFormat, "Log writer format (text, json).")
	fs.StringVar(&c.LogFile, "log-file", c.LogFile, "Log file path (when log-output is file).")
}

// BindFlags normalizes and parses the command line flags
func (c *Config) BindFlags() {
	if pflag.Parsed() {
		return
	}

	viper.SetEnvPrefix("TINYSYSLOG")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	c.addFlags(pflag.CommandLine)
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(fmt.Errorf("failed binding flags: %v", err))
	}

	c.setupLogger()
}

// setupLogger configures zerolog based on config
func (c *Config) setupLogger() {
	level, err := zerolog.ParseLevel(viper.GetString(LogLevel))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	logOutput := viper.GetString(LogOutput)
	logFile := viper.GetString("log-file")

	var output io.Writer
	switch logOutput {
	case "stderr":
		output = os.Stderr
	case "file":
		if logFile != "" && logFile != "stdout" && logFile != "stderr" {
			file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				output = os.Stdout
			} else {
				output = file
			}
		} else {
			output = os.Stdout
		}
	default:
		output = os.Stdout
	}

	logFormat := viper.GetString(LogWriter)
	if logFormat == "text" || logFormat == "" {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: "2006-01-02T15:04:05Z07:00",
		}
	}

	logger := zerolog.New(output).With().Timestamp().Logger()
	log.Logger = logger
}
