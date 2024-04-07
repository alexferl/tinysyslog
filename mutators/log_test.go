package mutators

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLog(t *testing.T) {
	p, err := time.Parse(time.RFC3339, "2016-01-01T12:01:21Z")
	assert.NoError(t, err)

	logParts := map[string]interface{}{
		"app_name":        "appname",
		"client":          "127.0.0.1:64844",
		"facility":        20,
		"hostname":        "hostname",
		"message":         "message",
		"msg_id":          "ID47",
		"priority":        165,
		"proc_id":         "1234",
		"severity":        5,
		"structured_data": `{"eventID":"123","eventSource":"test","exampleSDID":"32473","iut":"9"}`,
		"timestamp":       p,
		"tls_peer":        "",
		"version":         1,
	}

	l := NewLog(logParts)
	assert.Equal(t, logParts["app_name"], l.AppName)
	assert.Equal(t, logParts["client"], l.Client)
	assert.Equal(t, logParts["facility"], l.Facility)
	assert.Equal(t, logParts["hostname"], l.Hostname)
	assert.Equal(t, logParts["message"], l.Message)
	assert.Equal(t, logParts["msg_id"], l.MsgID)
	assert.Equal(t, logParts["priority"], l.Priority)
	assert.Equal(t, logParts["proc_id"], l.ProcId)
	assert.Equal(t, logParts["severity"], l.Severity)
	// assert.Equal(t, logParts["structured_data"], l.StructuredData)
	assert.Equal(t, logParts["timestamp"], l.Timestamp)
	assert.Equal(t, logParts["tls_peer"], l.TLSPeer)
	assert.Equal(t, logParts["version"], l.Version)
}
