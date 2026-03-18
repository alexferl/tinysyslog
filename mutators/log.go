package mutators

import (
	"strings"
	"time"
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
	sdString := "-"
	if sd, ok := logParts["structured_data"]; ok && sd != nil {
		sdString = sd.(string)
	}

	sd := parseStructuredData(sdString)
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

	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")

	// Find first space to separate SD-ID from params
	// SD-ID can contain @ for enterprise ID (e.g., "req@32473")
	firstSpace := strings.Index(s, " ")
	if firstSpace < 0 {
		// No params, just SD-ID
		return m
	}

	sdID := s[:firstSpace] // "req@32473"
	paramsStr := s[firstSpace+1:]

	m["sd_id"] = sdID

	items := strings.Split(paramsStr, " ")
	for _, i := range items {
		equal := strings.Index(i, "=")
		if equal >= 0 {
			key := i[:equal]
			value := strings.Trim(i[equal+1:], `"`)
			m[key] = value
		}
	}
	return m
}
