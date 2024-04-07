package sinks

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsole(t *testing.T) {
	msg := "test"
	f, err := os.CreateTemp("/tmp", "")
	assert.NoError(t, err)

	c := NewConsole(f)
	err = c.Write([]byte(msg))
	assert.NoError(t, err)

	_, err = f.Seek(0, 0)
	assert.NoError(t, err)

	all, err := io.ReadAll(f)
	assert.NoError(t, err)

	assert.Equal(t, ConsoleKind, c.GetKind())
	assert.Equal(t, msg, string(all))
}
