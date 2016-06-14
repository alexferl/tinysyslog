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

// Mutate mutates a slice of bytes
func (jm *JSONMutator) Mutate(logParts map[string]interface{}) string {
	t := logParts["timestamp"].(time.Time)
	// will eventually need to support user-defined format
	m := map[string]interface{}{
		"timestamp": t.Format(time.RFC3339Nano),
		"hostname":  logParts["hostname"].(string),
		"app_name":  logParts["app_name"].(string),
		"proc_id":   logParts["proc_id"].(string),
		"severity":  util.SeverityNumToString(logParts["severity"].(int)),
		"message":   logParts["message"].(string),
	}
	formatted, _ := json.Marshal(m)
	return string(formatted)
}
