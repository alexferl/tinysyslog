package sinks

import (
	"bufio"
	"os"
)

// Console represents a filesystem sink
type Console struct {
	output *os.File
	kind   Kind
}

// NewConsole creates a Console instance
func NewConsole(output *os.File) Sink {
	return Sink(&Console{
		output: output,
		kind:   ConsoleKind,
	})
}

// Write writes to the specified output
func (c *Console) Write(stdOutput []byte) error {
	w := bufio.NewWriter(c.output)
	defer w.Flush()
	_, err := w.Write(stdOutput)
	if err != nil {
		return err
	}
	return nil
}

func (c *Console) GetKind() Kind {
	return c.kind
}
