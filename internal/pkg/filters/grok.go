package filters

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/vjeantet/grok"
)

// GrokFilter represents a grok filter
type GrokFilter struct {
	g       *grok.Grok
	fields  []string
	pattern string
}

// NewGrokFilter creates a GrokFilter instance
func NewGrokFilter(pattern string, fields []string) Filter {
	g, err := grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	if err != nil {
		log.Panic().Msgf("Failed to initialize grok: '%v'", err)
	}

	return Filter(&GrokFilter{
		g:       g,
		fields:  fields,
		pattern: pattern,
	})
}

// Filter filters a log entry
func (gf *GrokFilter) Filter(data string) (string, error) {
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(data), &logEntry)
	if err != nil {
		return "", err
	}

	values, err := gf.g.Parse(gf.pattern, logEntry["message"].(string))
	if err != nil {
		return "", err
	}

	if len(gf.fields) <= 0 {
		for k, v := range values {
			if k != "timestamp" {
				logEntry[k] = v
			}
		}
	} else if len(gf.fields) >= 1 {
		for k, v := range values {
			if k != "timestamp" && stringInSlice(k, gf.fields) {
				logEntry[k] = v
			}
		}
	}

	out, err := json.Marshal(logEntry)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
