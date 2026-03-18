package server

import (
	"bytes"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/leodido/go-syslog/v4"
	"github.com/leodido/go-syslog/v4/rfc5424"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/alexferl/tinysyslog/config"
	"github.com/alexferl/tinysyslog/factories"
	"github.com/alexferl/tinysyslog/mutators"
	"github.com/alexferl/tinysyslog/sinks"
)

// Server holds the config
type Server struct {
	listener net.Listener
	udpConn  net.PacketConn
}

// New creates a Server instance
func New() (*Server, error) {
	address := viper.GetString(config.BindAddr)
	socketType := strings.ToLower(viper.GetString(config.SocketType))

	server := &Server{}

	switch socketType {
	case "tcp":
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return nil, err
		}
		server.listener = listener
		go server.handleTCP(listener)
	case "udp":
		conn, err := net.ListenPacket("udp", address)
		if err != nil {
			return nil, err
		}
		server.udpConn = conn
		go server.handleUDP(conn)
	default:
		// Both TCP and UDP
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return nil, err
		}
		server.listener = listener
		go server.handleTCP(listener)

		conn, err := net.ListenPacket("udp", address)
		if err != nil {
			return nil, err
		}
		server.udpConn = conn
		go server.handleUDP(conn)
	}

	log.Info().Msg("tinysyslog starting")
	log.Info().Msgf("tinysyslog listening on %s", address)

	return server, nil
}

func (s *Server) handleTCP(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error().Err(err).Msg("failed to accept TCP connection")
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close connection")
		}
	}()
	client := conn.RemoteAddr().String()

	parser := rfc5424.NewParser(rfc5424.WithBestEffort())

	buf := make([]byte, 65536)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		// Handle multiple messages in one read (newline separated for TCP)
		messages := bytes.Split(buf[:n], []byte("\n"))
		for _, msg := range messages {
			if len(msg) == 0 {
				continue
			}
			s.processMessage(msg, client, "", parser)
		}
	}
}

func (s *Server) handleUDP(conn net.PacketConn) {
	parser := rfc5424.NewParser(rfc5424.WithBestEffort())

	buf := make([]byte, 65536)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Error().Err(err).Msg("failed to read UDP packet")
			continue
		}

		client := addr.String()
		message := make([]byte, n)
		copy(message, buf[:n])

		// For UDP, parse the message directly
		s.processMessage(message, client, "", parser)
	}
}

func (s *Server) processMessage(data []byte, client, tlsPeer string, parser syslog.Machine) {
	parsed, err := parser.Parse(data)
	if err != nil {
		log.Debug().Err(err).Str("msg", string(data)).Msg("failed to parse message")
		return
	}

	s.processParsed(parsed, client, tlsPeer)
}

func (s *Server) processParsed(parsed syslog.Message, client, tlsPeer string) {
	// Cast to RFC5424 message
	rfcMsg, ok := parsed.(*rfc5424.SyslogMessage)
	if !ok {
		log.Debug().Msg("message is not RFC5424")
		return
	}

	// Convert to logParts map with defaults
	logParts := make(map[string]interface{})

	// Set defaults first
	logParts["priority"] = 0
	logParts["facility"] = 0
	logParts["severity"] = 0
	logParts["version"] = 1
	logParts["timestamp"] = time.Now()
	logParts["hostname"] = "-"
	logParts["app_name"] = "-"
	logParts["proc_id"] = "-"
	logParts["msg_id"] = "-"
	logParts["structured_data"] = "-"
	logParts["message"] = ""
	logParts["client"] = client
	logParts["tls_peer"] = tlsPeer

	// Override with actual values if present
	if rfcMsg.Priority != nil {
		logParts["priority"] = int(*rfcMsg.Priority)
		logParts["facility"] = int(*rfcMsg.Priority) / 8
		logParts["severity"] = int(*rfcMsg.Priority) % 8
	}
	logParts["version"] = int(rfcMsg.Version)
	if rfcMsg.Timestamp != nil {
		logParts["timestamp"] = *rfcMsg.Timestamp
	}
	if rfcMsg.Hostname != nil {
		logParts["hostname"] = *rfcMsg.Hostname
	}
	if rfcMsg.Appname != nil {
		logParts["app_name"] = *rfcMsg.Appname
	}
	if rfcMsg.ProcID != nil {
		logParts["proc_id"] = *rfcMsg.ProcID
	}
	if rfcMsg.MsgID != nil {
		logParts["msg_id"] = *rfcMsg.MsgID
	}
	if rfcMsg.StructuredData != nil && len(*rfcMsg.StructuredData) > 0 {
		logParts["structured_data"] = formatStructuredData(*rfcMsg.StructuredData)
	}
	if rfcMsg.Message != nil {
		logParts["message"] = *rfcMsg.Message
	}

	// Process the log
	mutator := factories.Mutator()
	filter := factories.Filter()
	sinksSlice := factories.Sinks()

	newLog := mutators.NewLog(logParts)

	mutated, err := mutator.Mutate(newLog)
	if err != nil {
		log.Err(err).Msg("failed mutating log")
		return
	}

	filtered := mutated
	if viper.GetString(config.Filter) != "" {
		filtered, err = filter.Filter(mutated)
		if err != nil {
			log.Err(err).Msg("failed filtering log")
			return
		}
	}

	if len(filtered) > 0 {
		for _, sink := range sinksSlice {
			go write(sink, filtered)
		}
	}
}

func formatStructuredData(sd map[string]map[string]string) string {
	// Format as [sdID param1="value1" param2="value2"][sdID2 ...]
	var result strings.Builder

	// Sort SD IDs for consistent ordering
	sdIDs := make([]string, 0, len(sd))
	for sdID := range sd {
		sdIDs = append(sdIDs, sdID)
	}
	sort.Strings(sdIDs)

	for _, sdID := range sdIDs {
		params := sd[sdID]
		result.WriteString("[")
		result.WriteString(sdID)

		// Sort parameter keys for consistent ordering
		keys := make([]string, 0, len(params))
		for k := range params {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			result.WriteString(" ")
			result.WriteString(key)
			result.WriteString(`="`)
			result.WriteString(params[key])
			result.WriteString(`"`)
		}
		result.WriteString("]")
	}
	return result.String()
}

// Run runs the server
func (s *Server) Run() {
	// Block forever
	select {}
}

func write(sink sinks.Sink, msg string) {
	sinkName := sink.GetKind().String()
	if err := sink.Write([]byte(msg + "\n")); err != nil {
		log.Err(err).Str("sink", sinkName).Msgf("failed writing log to sink: %s", sinkName)
	} else {
		log.Debug().Msgf("wrote log to %s sink", sinkName)
	}
}
