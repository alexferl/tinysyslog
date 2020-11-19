package sinks

import (
	"bufio"
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
func (cs *ConsoleSink) Write(stdOutput []byte) error {
	w := bufio.NewWriter(cs.output)
	defer w.Flush()
	_, err := w.Write(stdOutput)
	if err != nil {
		return err
	}
	return nil
}
