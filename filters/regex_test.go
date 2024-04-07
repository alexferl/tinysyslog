package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegex(t *testing.T) {
	msg := `<165>1 2016-01-01T12:01:21Z hostname appname 1234 ID47 [exampleSDID@32473 iut="9" eventSource="test" eventID="123"] message"`

	testCases := []struct {
		name string
		re   string
		res  string
		err  bool
	}{
		{"match", "appname", "", false},
		{"no match", "xyz", msg, false},
		{"no regex", "", msg, false},
		{"invalid regex", "(?=\"", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := NewRegex(tc.re)
			s, err := f.Filter(msg)
			assert.Equal(t, RegexKind, f.GetKind())
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.res, s)
			}
		})
	}
}
