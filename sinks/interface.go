package sinks

import (
	"strings"

	"tinysyslog/util"
)

// Sink a common interface for all sinks
type Sink interface {
	Write([]byte) error
}

func GetSinkName(sink Sink) string {
	return strings.ToLower(util.GetType(sink))
}
