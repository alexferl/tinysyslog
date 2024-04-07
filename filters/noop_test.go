package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoOp(t *testing.T) {
	msg := `<165>1 2016-01-01T12:01:21Z hostname appname 1234 ID47 [exampleSDID@32473 iut="9" eventSource="test" eventID="123"] message"`

	f := NewNoOp()
	s, err := f.Filter(msg)
	assert.NoError(t, err)
	assert.Equal(t, msg, s)
}
