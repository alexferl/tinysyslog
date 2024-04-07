package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeverityNumToString(t *testing.T) {
	assert.Equal(t, severities[0], SeverityNumToString(0))
	assert.Equal(t, severities[1], SeverityNumToString(1))
	assert.Equal(t, severities[2], SeverityNumToString(2))
	assert.Equal(t, severities[3], SeverityNumToString(3))
	assert.Equal(t, severities[4], SeverityNumToString(4))
	assert.Equal(t, severities[5], SeverityNumToString(5))
	assert.Equal(t, severities[6], SeverityNumToString(6))
	assert.Equal(t, severities[7], SeverityNumToString(7))
	assert.Equal(t, "DEFAULT", SeverityNumToString(8))
}
