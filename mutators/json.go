package mutators

import (
	"encoding/json"
	"time"

	"github.com/admiralobvious/tinysyslog/util"
)

// JSONMutator represents a JSON mutator
type JSONMutator struct{}

// NewJSONMutator creates a JSONMutator instance
func NewJSONMutator() Mutator {
	return Mutator(&JSONMutator{})
}

// Mutate mutates a Log
func (jm *JSONMutator) Mutate(log Log) (string, error) {
	m := map[string]interface{}{
		"app_name":        log.AppName,
		"client":          log.Client,
		"facility":        log.Facility,
		"hostname":        log.Hostname,
		"message":         log.Message,
		"msg_id":          log.MsgID,
		"priority":        log.Priority,
		"proc_id":         log.ProcId,
		"severity":        util.SeverityNumToString(log.Severity),
		"structured_data": log.StructuredData,
		"timestamp":       log.Timestamp.Format(time.RFC3339Nano),
		"tls_peer":        log.TLSPeer,
		"version":         log.Version,
	}
	formatted, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}
