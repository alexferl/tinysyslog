package sinks

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

// Filesystem represents a filesystem sink
type Filesystem struct {
	logger *lumberjack.Logger
}

// NewFilesystem creates a Filesystem instance
func NewFilesystem(filename string, maxAge, maxBackups, maxSize int) Sink {
	logger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
	}

	return Sink(&Filesystem{
		logger: logger,
	})
}

// Write writes to a file
func (fs *Filesystem) Write(output []byte) error {
	_, err := fs.logger.Write(output)
	if err != nil {
		return err
	}

	return nil
}
