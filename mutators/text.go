package mutators

import (
	"fmt"
	"time"
)

// TextMutator represents a text mutator
type TextMutator struct{}

// NewTextMutator creates TextMutator instance
func NewTextMutator() Mutator {
	return Mutator(&TextMutator{})
}

// Mutate mutates a slice of bytes
func (tm *TextMutator) Mutate(logParts map[string]interface{}) string {
	t := logParts["timestamp"].(time.Time)
	// will eventually need to support user-defined format
	formatted := fmt.Sprintf("%s %s %s[%s]: %s",
		t.Format("Jan _2 15:04:05"),
		logParts["hostname"],
		logParts["app_name"],
		logParts["proc_id"],
		logParts["message"])
	return formatted
}
