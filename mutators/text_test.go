package mutators

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	m := NewText()

	p, err := time.Parse(time.RFC3339, "2016-01-01T12:01:21Z")
	assert.NoError(t, err)

	l := Log{
		AppName:        "appname",
		Client:         "127.0.0.1:64844",
		Facility:       20,
		Hostname:       "hostname",
		Message:        "message",
		MsgID:          "ID47",
		Priority:       165,
		ProcId:         "1234",
		Severity:       5,
		StructuredData: map[string]string{"eventID": "123", "eventSource": "test", "exampleSDID": "32473", "iut": "9"},
		Timestamp:      p,
		TLSPeer:        "",
		Version:        1,
	}

	res, err := m.Mutate(l)
	assert.NoError(t, err)
	assert.Equal(t, "Jan  1 12:01:21 hostname appname[1234]: message", res)
	assert.Equal(t, TextKind, m.GetKind())
}
