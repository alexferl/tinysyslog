package sinks

// Sink a common interface for all sinks
type Sink interface {
	Write([]byte) error
}
