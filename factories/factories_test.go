package factories

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/alexferl/tinysyslog/config"
	"github.com/alexferl/tinysyslog/constants"
	"github.com/alexferl/tinysyslog/filters"
	"github.com/alexferl/tinysyslog/mutators"
	"github.com/alexferl/tinysyslog/sinks"
)

func TestMutator(t *testing.T) {
	tests := []struct {
		name         string
		mutatorType  string
		expectedKind mutators.Kind
	}{
		{
			name:         "default text mutator",
			mutatorType:  "",
			expectedKind: mutators.TextKind,
		},
		{
			name:         "explicit text mutator",
			mutatorType:  "text",
			expectedKind: mutators.TextKind,
		},
		{
			name:         "json mutator",
			mutatorType:  "json",
			expectedKind: mutators.JSONKind,
		},
		{
			name:         "unknown mutator falls back to text",
			mutatorType:  "unknown",
			expectedKind: mutators.TextKind,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			viper.Set(config.Mutator, tt.mutatorType)

			m := Mutator()
			assert.Equal(t, tt.expectedKind, m.GetKind())
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name         string
		filterType   string
		regex        string
		expectedKind filters.Kind
	}{
		{
			name:         "no filter",
			filterType:   "",
			regex:        "",
			expectedKind: filters.NoOpKind,
		},
		{
			name:         "regex filter",
			filterType:   "regex",
			regex:        "^test.*",
			expectedKind: filters.RegexKind,
		},
		{
			name:         "unknown filter falls back to noop",
			filterType:   "unknown",
			regex:        "",
			expectedKind: filters.NoOpKind,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			viper.Set(config.Filter, tt.filterType)
			viper.Set(config.FilterRegex, tt.regex)

			f := Filter()
			assert.Equal(t, tt.expectedKind, f.GetKind())
		})
	}
}

func TestSinks(t *testing.T) {
	tests := []struct {
		name          string
		sinksConfig   []string
		consoleOut    string
		expectedLen   int
		expectedKinds []sinks.Kind
	}{
		{
			name:          "no sinks configured",
			sinksConfig:   []string{},
			consoleOut:    "stdout",
			expectedLen:   0,
			expectedKinds: nil,
		},
		{
			name:          "console sink stdout",
			sinksConfig:   []string{"console"},
			consoleOut:    "stdout",
			expectedLen:   1,
			expectedKinds: []sinks.Kind{sinks.ConsoleKind},
		},
		{
			name:          "console sink stderr",
			sinksConfig:   []string{"console"},
			consoleOut:    "stderr",
			expectedLen:   1,
			expectedKinds: []sinks.Kind{sinks.ConsoleKind},
		},
		{
			name:          "console sink unknown output",
			sinksConfig:   []string{"console"},
			consoleOut:    "unknown",
			expectedLen:   1,
			expectedKinds: []sinks.Kind{sinks.ConsoleKind},
		},
		{
			name:          "filesystem sink",
			sinksConfig:   []string{"filesystem"},
			consoleOut:    "stdout",
			expectedLen:   1,
			expectedKinds: []sinks.Kind{sinks.FilesystemKind},
		},
		{
			name:          "multiple sinks",
			sinksConfig:   []string{"console", "filesystem"},
			consoleOut:    "stdout",
			expectedLen:   2,
			expectedKinds: []sinks.Kind{sinks.ConsoleKind, sinks.FilesystemKind},
		},
		{
			name:          "unknown sink type",
			sinksConfig:   []string{"unknown"},
			consoleOut:    "stdout",
			expectedLen:   0,
			expectedKinds: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			viper.Set(config.Sinks, tt.sinksConfig)
			viper.Set(config.SinkConsoleOutput, tt.consoleOut)
			viper.Set(config.SinkFilesystemFilename, "test.log")
			viper.Set(config.SinkFilesystemMaxAge, 30)
			viper.Set(config.SinkFilesystemMaxBackups, 10)
			viper.Set(config.SinkFilesystemMaxSize, 100)

			s := Sinks()
			assert.Len(t, s, tt.expectedLen)

			for i, kind := range tt.expectedKinds {
				assert.Equal(t, kind, s[i].GetKind())
			}
		})
	}
}

func TestSinks_ConsoleOutputValues(t *testing.T) {
	tests := []struct {
		name       string
		output     string
		shouldWork bool
	}{
		{
			name:       "stdout",
			output:     constants.ConsoleStdOut,
			shouldWork: true,
		},
		{
			name:       "stderr",
			output:     constants.ConsoleStdErr,
			shouldWork: true,
		},
		{
			name:       "unknown falls back",
			output:     "unknown",
			shouldWork: true, // Creates console but with nil output
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			viper.Set(config.Sinks, []string{"console"})
			viper.Set(config.SinkConsoleOutput, tt.output)

			s := Sinks()
			assert.Len(t, s, 1)
			assert.Equal(t, sinks.ConsoleKind, s[0].GetKind())
		})
	}
}
