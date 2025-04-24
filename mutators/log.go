package mutators

import (
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type Log struct {
	AppName        string            `json:"app_name"`
	Client         string            `json:"client"`
	Facility       int               `json:"facility"`
	Hostname       string            `json:"hostname"`
	Message        string            `json:"message"`
	MsgID          string            `json:"msg_id"`
	Priority       int               `json:"priority"`
	ProcId         string            `json:"proc_id"`
	Severity       int               `json:"severity"`
	StructuredData map[string]string `json:"structured_data"`
	Timestamp      time.Time         `json:"timestamp"`
	TLSPeer        string            `json:"tls_peer"`
	Version        int               `json:"version"`
}

// NewLog creates a Log instance
func NewLog(logParts map[string]interface{}) Log {
	sd := parseStructuredData(logParts["structured_data"].(string))
	return Log{
		AppName:        logParts["app_name"].(string),
		Client:         logParts["client"].(string),
		Facility:       logParts["facility"].(int),
		Hostname:       logParts["hostname"].(string),
		Message:        logParts["message"].(string),
		MsgID:          logParts["msg_id"].(string),
		Priority:       logParts["priority"].(int),
		ProcId:         logParts["proc_id"].(string),
		Severity:       logParts["severity"].(int),
		StructuredData: sd,
		Timestamp:      logParts["timestamp"].(time.Time),
		TLSPeer:        logParts["tls_peer"].(string),
		Version:        logParts["version"].(int),
	}
}

func parseStructuredData(s string) map[string]string {
	m := make(map[string]string)

	replacer := strings.NewReplacer("[", "", "]", "")
	s = replacer.Replace(s)
	items := strings.Split(s, " ")

	for _, i := range items {
		at := strings.Index(i, "@")
		if at >= 0 {
			kv := strings.Split(i, "@")
			if len(kv) < 2 {
				log.Error().Msgf("failed parsing structured data item: %v", i)
			} else {
				m[kv[0]] = kv[1]
			}
		}

		equal := strings.Index(i, "=")
		if equal >= 0 {
			kv := strings.Split(i, "=")
			if len(kv) < 2 {
				log.Error().Msgf("failed parsing structured data item: %v", i)
			} else {
				m[kv[0]] = strings.ReplaceAll(kv[1], "\"", "")
			}
		}
	}
	return m
}
