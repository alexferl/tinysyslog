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
		"timestamp": log.Timestamp.Format(time.RFC3339Nano),
		"hostname":  log.Hostname,
		"app_name":  log.AppName,
		"proc_id":   log.ProcId,
		"severity":  util.SeverityNumToString(log.Severity),
		"message":   log.Message,
	}
	formatted, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}
