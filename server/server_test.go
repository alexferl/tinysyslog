package server

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"

	"github.com/leodido/go-syslog/v4/rfc5424"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/alexferl/tinysyslog/config"
	"github.com/alexferl/tinysyslog/sinks"
)

func init() {
	// Suppress logs during tests - expected "closed connection" errors are noisy
	zerolog.SetGlobalLevel(zerolog.FatalLevel)
	log.Logger = zerolog.New(io.Discard)
}

func getFreePort() string {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "127.0.0.1:0"
	}
	defer func() { _ = l.Close() }()
	return l.Addr().String()
}

func TestNew_TCPOnly(t *testing.T) {
	viper.Reset()
	viper.Set(config.BindAddr, getFreePort())
	viper.Set(config.SocketType, "tcp")
	viper.Set(config.Sinks, []string{})

	s, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotNil(t, s.listener)
	assert.Nil(t, s.udpConn)

	if s.listener != nil {
		_ = s.listener.Close()
	}
}

func TestNew_UDPOnly(t *testing.T) {
	viper.Reset()
	viper.Set(config.BindAddr, getFreePort())
	viper.Set(config.SocketType, "udp")
	viper.Set(config.Sinks, []string{})

	s, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Nil(t, s.listener)
	assert.NotNil(t, s.udpConn)

	if s.udpConn != nil {
		_ = s.udpConn.Close()
	}
}

func TestNew_BothTCPAndUDP(t *testing.T) {
	viper.Reset()
	viper.Set(config.BindAddr, getFreePort())
	viper.Set(config.SocketType, "")
	viper.Set(config.Sinks, []string{})

	s, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotNil(t, s.listener)
	assert.NotNil(t, s.udpConn)

	if s.listener != nil {
		_ = s.listener.Close()
	}
	if s.udpConn != nil {
		_ = s.udpConn.Close()
	}
}

func TestNew_InvalidAddress(t *testing.T) {
	viper.Reset()
	viper.Set(config.BindAddr, "invalid:address:here")
	viper.Set(config.SocketType, "tcp")

	s, err := New()
	assert.Error(t, err)
	assert.Nil(t, s)
}

func TestProcessMessage_ValidMessage(t *testing.T) {
	viper.Reset()
	viper.Set(config.Sinks, []string{})

	s := &Server{}
	parser := rfc5424.NewParser(rfc5424.WithBestEffort())

	// Valid RFC5424 message
	msg := []byte("<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut=\"3\" eventSource=\"Application\" eventID=\"1011\"] BOMAn application event log entry\n")

	// Should not panic
	s.processMessage(msg, "127.0.0.1:12345", "", parser)
}

func TestProcessMessage_InvalidMessage(t *testing.T) {
	viper.Reset()
	viper.Set(config.Sinks, []string{})

	s := &Server{}
	parser := rfc5424.NewParser(rfc5424.WithBestEffort())

	// Invalid message - should not panic, just log and return
	msg := []byte("not a valid syslog message")
	s.processMessage(msg, "127.0.0.1:12345", "", parser)
}

func TestProcessMessage_EmptyMessage(t *testing.T) {
	viper.Reset()
	viper.Set(config.Sinks, []string{})

	s := &Server{}
	parser := rfc5424.NewParser(rfc5424.WithBestEffort())

	// Empty message
	msg := []byte{}
	s.processMessage(msg, "127.0.0.1:12345", "", parser)
}

func TestProcessParsed_MinimalMessage(t *testing.T) {
	viper.Reset()
	viper.Set(config.Sinks, []string{})

	s := &Server{}
	parser := rfc5424.NewParser(rfc5424.WithBestEffort())

	// Use parser to create a valid message with minimal fields
	msg := []byte("<165>1 2003-10-11T22:14:15.003Z - - - - - -")
	parsed, err := parser.Parse(msg)
	assert.NoError(t, err)
	assert.NotNil(t, parsed)

	s.processParsed(parsed, "127.0.0.1:12345", "")
}

func TestProcessParsed_FullMessage(t *testing.T) {
	viper.Reset()
	viper.Set(config.Sinks, []string{})

	s := &Server{}
	parser := rfc5424.NewParser(rfc5424.WithBestEffort())

	// Full RFC5424 message with all fields
	msg := []byte("<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog 1234 ID47 [exampleSDID@32473 iut=\"3\" eventSource=\"Application\" eventID=\"1011\"] An application event log entry")
	parsed, err := parser.Parse(msg)
	assert.NoError(t, err)
	assert.NotNil(t, parsed)

	s.processParsed(parsed, "127.0.0.1:12345", "")
}

func TestProcessParsed_WithFilter(t *testing.T) {
	viper.Reset()
	viper.Set(config.Sinks, []string{})
	viper.Set(config.Filter, "")

	s := &Server{}
	parser := rfc5424.NewParser(rfc5424.WithBestEffort())

	// Create a valid message using parser
	msg := []byte("<165>1 2003-10-11T22:14:15.003Z test test - - - test message")
	parsed, err := parser.Parse(msg)
	assert.NoError(t, err)
	assert.NotNil(t, parsed)

	s.processParsed(parsed, "127.0.0.1:12345", "")
}

func TestProcessParsed_NonRFC5424Message(t *testing.T) {
	viper.Reset()
	viper.Set(config.Sinks, []string{})

	s := &Server{}

	// Pass a nil message - should just return
	s.processParsed(nil, "127.0.0.1:12345", "")
}

func TestFormatStructuredData(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]map[string]string
		expected string
	}{
		{
			name:     "nil map",
			input:    nil,
			expected: "",
		},
		{
			name:     "empty map",
			input:    map[string]map[string]string{},
			expected: "",
		},
		{
			name: "single SD element",
			input: map[string]map[string]string{
				"exampleSDID@32473": {
					"iut":         "3",
					"eventSource": "Application",
					"eventID":     "1011",
				},
			},
			expected: `[exampleSDID@32473 eventID="1011" eventSource="Application" iut="3"]`,
		},
		{
			name: "multiple SD elements",
			input: map[string]map[string]string{
				"sd1": {
					"param1": "value1",
				},
				"sd2": {
					"param2": "value2",
				},
			},
			expected: `[sd1 param1="value1"][sd2 param2="value2"]`,
		},
		{
			name: "SD without params",
			input: map[string]map[string]string{
				"emptySD": {},
			},
			expected: `[emptySD]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatStructuredData(tt.input)
			// Note: map iteration order is non-deterministic, so for multi-element tests
			// we just check that the result contains expected parts
			if tt.name == "multiple SD elements" {
				assert.Contains(t, result, `[sd1 param1="value1"]`)
				assert.Contains(t, result, `[sd2 param2="value2"]`)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestRun(t *testing.T) {
	s := &Server{}

	// Run blocks forever, so we test it in a goroutine and timeout
	done := make(chan struct{})
	go func() {
		s.Run()
		close(done)
	}()

	select {
	case <-done:
		t.Error("Run should block forever")
	case <-time.After(10 * time.Millisecond):
		// Expected - Run blocks
	}
}

func TestWrite(t *testing.T) {
	viper.Reset()
	viper.Set(config.Sinks, []string{"console"})
	viper.Set(config.SinkConsoleOutput, "stdout")

	// Create a console sink
	consoleSink := sinks.NewConsole(nil)

	// Should not panic
	write(consoleSink, "test message")
}

func TestHandleConnection(t *testing.T) {
	viper.Reset()
	viper.Set(config.Sinks, []string{})

	s := &Server{}

	// Create a pipe to simulate a connection
	client, server := net.Pipe()
	defer func() { _ = client.Close() }()

	// Write a valid RFC5424 message to the client side
	go func() {
		msg := "<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 - BOMAn application event log entry\n"
		_, _ = client.Write([]byte(msg))
		_ = client.Close()
	}()

	// Handle the connection (server side)
	done := make(chan struct{})
	go func() {
		s.handleConnection(server)
		close(done)
	}()

	// Wait for handling to complete or timeout
	select {
	case <-done:
		// Expected
	case <-time.After(100 * time.Millisecond):
		// Also acceptable - connection may still be open
	}
}

func TestHandleConnection_MultipleMessages(t *testing.T) {
	viper.Reset()
	viper.Set(config.Sinks, []string{})

	s := &Server{}

	// Create a pipe
	client, server := net.Pipe()
	defer func() { _ = client.Close() }()

	// Write multiple messages
	go func() {
		msg1 := "<165>1 2003-10-11T22:14:15.003Z host1 app1 - ID1 - Message1\n"
		msg2 := "<166>1 2003-10-11T22:14:15.003Z host2 app2 - ID2 - Message2\n"
		_, _ = client.Write([]byte(msg1 + msg2))
		_ = client.Close()
	}()

	// Handle the connection
	done := make(chan struct{})
	go func() {
		s.handleConnection(server)
		close(done)
	}()

	select {
	case <-done:
		// Expected
	case <-time.After(100 * time.Millisecond):
	}
}

func TestHandleConnection_EmptyMessages(t *testing.T) {
	viper.Reset()
	viper.Set(config.Sinks, []string{})

	s := &Server{}

	client, server := net.Pipe()
	defer func() { _ = client.Close() }()

	// Write empty lines and newlines only
	go func() {
		_, _ = client.Write([]byte("\n\n\n"))
		_ = client.Close()
	}()

	done := make(chan struct{})
	go func() {
		s.handleConnection(server)
		close(done)
	}()

	select {
	case <-done:
		// Expected
	case <-time.After(100 * time.Millisecond):
	}
}

func TestBytesSplit(t *testing.T) {
	// Test the bytes.Split behavior used in handleConnection
	input := []byte("msg1\nmsg2\nmsg3\n")
	parts := bytes.Split(input, []byte("\n"))

	assert.Equal(t, 4, len(parts))
	assert.Equal(t, []byte("msg1"), parts[0])
	assert.Equal(t, []byte("msg2"), parts[1])
	assert.Equal(t, []byte("msg3"), parts[2])
	assert.Equal(t, []byte{}, parts[3])
}
