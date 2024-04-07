package mutators

import (
	"encoding/json"
	"time"

	"tinysyslog/util"
)

// JSON represents a JSON mutator
type JSON struct {
	kind Kind
}

// NewJSON creates a JSON instance
func NewJSON() Mutator {
	return Mutator(&JSON{kind: JSONKind})
}

// Mutate mutates a Log
func (j *JSON) Mutate(log Log) (string, error) {
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

func (j *JSON) GetKind() Kind {
	return j.kind
}
