package mutators

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	m := NewJSON()

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
	assert.Equal(t, `{"app_name":"appname","client":"127.0.0.1:64844","facility":20,"hostname":"hostname","message":"message","msg_id":"ID47","priority":165,"proc_id":"1234","severity":"NOTICE","structured_data":{"eventID":"123","eventSource":"test","exampleSDID":"32473","iut":"9"},"timestamp":"2016-01-01T12:01:21Z","tls_peer":"","version":1}`, res)
	assert.Equal(t, JSONKind, m.GetKind())
}
