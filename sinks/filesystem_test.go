package sinks

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesystem(t *testing.T) {
	msg := "test"
	path := "/tmp/syslog.log"
	_ = os.Remove(path)

	fs := NewFilesystem(path, 30, 10, 100)

	err := fs.Write([]byte(msg))
	assert.NoError(t, err)

	f, err := os.ReadFile(path)
	assert.NoError(t, err)

	assert.Equal(t, FilesystemKind, fs.GetKind())
	assert.Equal(t, msg, string(f))
}
