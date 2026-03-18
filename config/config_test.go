package config

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/alexferl/tinysyslog/sinks"
)

func TestNew(t *testing.T) {
	c := New()

	assert.NotNil(t, c)
	assert.Equal(t, "127.0.0.1:5140", c.BindAddr)
	assert.Equal(t, "stdout", c.ConsoleSink.Output)
	assert.Equal(t, "syslog.log", c.FilesystemSink.Filename)
	assert.Equal(t, 30, c.FilesystemSink.MaxAge)
	assert.Equal(t, 10, c.FilesystemSink.MaxBackups)
	assert.Equal(t, 100, c.FilesystemSink.MaxSize)
	assert.Empty(t, c.FilterType)
	assert.Empty(t, c.LogFile)
	assert.Equal(t, "text", c.LogFormat)
	assert.Equal(t, "info", c.LogLevel)
	assert.Equal(t, "stdout", c.LogOutput)
	assert.Equal(t, "text", c.MutatorType)
	assert.Empty(t, c.RegexFilter.Regex)
	assert.Equal(t, []string{sinks.ConsoleKind.String()}, c.SinkTypes)
	assert.Empty(t, c.SocketType)
}

func TestConstants(t *testing.T) {
	assert.Equal(t, "bind-addr", BindAddr)
	assert.Equal(t, "filter", Filter)
	assert.Equal(t, "filter-regex", FilterRegex)
	assert.Equal(t, "mutator", Mutator)
	assert.Equal(t, "sinks", Sinks)
	assert.Equal(t, "sink-console-output", SinkConsoleOutput)
	assert.Equal(t, "sink-filesystem-filename", SinkFilesystemFilename)
	assert.Equal(t, "sink-filesystem-max-age", SinkFilesystemMaxAge)
	assert.Equal(t, "sink-filesystem-max-backups", SinkFilesystemMaxBackups)
	assert.Equal(t, "sink-filesystem-max-size", SinkFilesystemMaxSize)
	assert.Equal(t, "socket-type", SocketType)
	assert.Equal(t, "log-level", LogLevel)
	assert.Equal(t, "log-output", LogOutput)
	assert.Equal(t, "log-writer", LogWriter)
}

func TestAddFlags(t *testing.T) {
	c := New()
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	c.addFlags(fs)

	// Check all expected flags are registered
	flags := []string{
		BindAddr,
		Filter,
		FilterRegex,
		Mutator,
		Sinks,
		SinkConsoleOutput,
		SinkFilesystemFilename,
		SinkFilesystemMaxAge,
		SinkFilesystemMaxBackups,
		SinkFilesystemMaxSize,
		SocketType,
		LogLevel,
		LogOutput,
		LogWriter,
	}

	for _, flag := range flags {
		f := fs.Lookup(flag)
		assert.NotNil(t, f, "flag %s should be registered", flag)
	}
}

func TestAddFlags_BindValues(t *testing.T) {
	c := New()
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	c.addFlags(fs)

	// Test that default values are set
	bindAddr, _ := fs.GetString(BindAddr)
	assert.Equal(t, "127.0.0.1:5140", bindAddr)

	logLevel, _ := fs.GetString(LogLevel)
	assert.Equal(t, "info", logLevel)

	logOutput, _ := fs.GetString(LogOutput)
	assert.Equal(t, "stdout", logOutput)
}

func TestSetupLogger(t *testing.T) {
	// Store original values to restore later
	origLevel := zerolog.GlobalLevel()
	origLogger := zerolog.Logger{}
	defer func() {
		zerolog.SetGlobalLevel(origLevel)
		zerolog.DefaultContextLogger = &origLogger
	}()

	tests := []struct {
		name      string
		logLevel  string
		logOutput string
		logFile   string
		wantLevel zerolog.Level
	}{
		{
			name:      "default info level stdout",
			logLevel:  "info",
			logOutput: "stdout",
			logFile:   "",
			wantLevel: zerolog.InfoLevel,
		},
		{
			name:      "debug level",
			logLevel:  "debug",
			logOutput: "stdout",
			logFile:   "",
			wantLevel: zerolog.DebugLevel,
		},
		{
			name:      "error level stderr",
			logLevel:  "error",
			logOutput: "stderr",
			logFile:   "",
			wantLevel: zerolog.ErrorLevel,
		},
		{
			name:      "warn level file",
			logLevel:  "warn",
			logOutput: "file",
			logFile:   os.TempDir() + "/test.log",
			wantLevel: zerolog.WarnLevel,
		},
		{
			name:      "invalid level falls back to info",
			logLevel:  "invalid",
			logOutput: "stdout",
			logFile:   "",
			wantLevel: zerolog.InfoLevel,
		},
		{
			name:      "fatal level",
			logLevel:  "fatal",
			logOutput: "stdout",
			logFile:   "",
			wantLevel: zerolog.FatalLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper for each test
			viper.Reset()

			c := New()
			c.LogLevel = tt.logLevel
			c.LogOutput = tt.logOutput
			c.LogFile = tt.logFile

			// Set viper values
			viper.Set(LogLevel, tt.logLevel)
			viper.Set(LogOutput, tt.logOutput)
			viper.Set(LogWriter, tt.logFile)

			c.setupLogger()

			assert.Equal(t, tt.wantLevel, zerolog.GlobalLevel())

			// Cleanup temp file if created
			if tt.logFile != "" {
				_ = os.Remove(tt.logFile)
			}
		})
	}
}

func TestSetupLogger_FileOutputFallback(t *testing.T) {
	viper.Reset()

	c := New()
	viper.Set(LogLevel, "info")
	viper.Set(LogOutput, "file")
	viper.Set(LogWriter, "/nonexistent/path/test.log")

	// Should fallback to stdout without error
	c.setupLogger()

	assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())
}

func TestConfigStructs(t *testing.T) {
	t.Run("ConsoleSink", func(t *testing.T) {
		cs := ConsoleSink{Output: "stderr"}
		assert.Equal(t, "stderr", cs.Output)
	})

	t.Run("FilesystemSink", func(t *testing.T) {
		fs := FilesystemSink{
			Filename:   "test.log",
			MaxAge:     7,
			MaxBackups: 5,
			MaxSize:    50,
		}
		assert.Equal(t, "test.log", fs.Filename)
		assert.Equal(t, 7, fs.MaxAge)
		assert.Equal(t, 5, fs.MaxBackups)
		assert.Equal(t, 50, fs.MaxSize)
	})

	t.Run("RegexFilter", func(t *testing.T) {
		rf := RegexFilter{Regex: "^test.*"}
		assert.Equal(t, "^test.*", rf.Regex)
	})
}
