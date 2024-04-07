package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	_ "tinysyslog/testing"
)

func TestNew(t *testing.T) {
	_, err := New()
	assert.NoError(t, err)
}
