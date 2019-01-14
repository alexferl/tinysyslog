package mutators

import (
	"fmt"
)

// TextMutator represents a text mutator
type TextMutator struct{}

// NewTextMutator creates TextMutator instance
func NewTextMutator() Mutator {
	return Mutator(&TextMutator{})
}

// Mutate mutates a Log
func (tm *TextMutator) Mutate(log Log) (string, error) {
	formatted := fmt.Sprintf("%s %s %s[%s]: %s",
		log.Timestamp.Format("Jan _2 15:04:05"),
		log.Hostname,
		log.AppName,
		log.ProcId,
		log.Message)
	return formatted, nil
}
