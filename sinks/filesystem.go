package sinks

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

// FilesystemSink represents a filesystem sink
type FilesystemSink struct {
	logger *lumberjack.Logger
}

// NewFilesystemSink creates a FilesystemSink instance
func NewFilesystemSink(filename string, maxAge, maxBackups, maxSize int) Sink {
	logger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
	}

	return Sink(&FilesystemSink{
		logger: logger,
	})
}

// Write writes to a file
func (fs *FilesystemSink) Write(output []byte) error {
	_, err := fs.logger.Write(output)
	if err != nil {
		return err
	}

	return nil
}
