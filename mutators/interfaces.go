package mutators

import "time"

// Mutator is a common interface for all mutators
type Mutator interface {
	Mutate(Log) (string, error)
}

type Log struct {
	Timestamp time.Time `json:"timestamp"`
	Hostname  string    `json:"hostname"`
	AppName   string    `json:"app_name"`
	ProcId    string    `json:"proc_id"`
	Severity  int       `json:"severity"`
	Message   string    `json:"message"`
}
