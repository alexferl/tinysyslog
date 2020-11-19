package sinks

import (
	"strings"

	"tinysyslog/internal/pkg/util"
)

// Sink a common interface for all sinks
type Sink interface {
	Write([]byte) error
}

func GetSinkName(sink Sink) string {
	name := util.GetType(sink)
	s := strings.Split(name, "Sink")
	return strings.ToLower(s[0])
}
