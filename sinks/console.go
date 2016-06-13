package sinks

import (
	"os"
)

// ConsoleSink represents a filesystem sink
type ConsoleSink struct {
	output *os.File
}

// NewConsoleSink creates a ConsoleSink instance
func NewConsoleSink(output *os.File) Sink {
	return Sink(&ConsoleSink{
		output: output,
	})
}

// Write writes to the specified output
func (fs *ConsoleSink) Write(output []byte) error {
	_, err := fs.output.Write(output)
	if err != nil {
		return err
	}

	return nil
}
