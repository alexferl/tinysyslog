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

func NewLog(logParts map[string]interface{}) Log {
	return Log{
		Timestamp: logParts["timestamp"].(time.Time),
		Hostname:  logParts["hostname"].(string),
		AppName:   logParts["app_name"].(string),
		ProcId:    logParts["proc_id"].(string),
		Severity:  logParts["severity"].(int),
		Message:   logParts["message"].(string),
	}
}
