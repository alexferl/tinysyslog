package mutators

import (
	"fmt"
)

// Text represents a text mutator
type Text struct {
	kind Kind
}

// NewText creates Text instance
func NewText() Mutator {
	return Mutator(&Text{kind: TextKind})
}

// Mutate mutates a Log
func (t *Text) Mutate(log Log) (string, error) {
	formatted := fmt.Sprintf("%s %s %s[%s]: %s",
		log.Timestamp.Format("Jan _2 15:04:05"),
		log.Hostname,
		log.AppName,
		log.ProcId,
		log.Message)
	return formatted, nil
}

func (t *Text) GetKind() Kind {
	return t.kind
}
