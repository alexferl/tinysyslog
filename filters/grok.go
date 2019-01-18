package filters

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
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
		logrus.Panicf("Failed to initialize grok: %v", err)
		panic(err)
	}

	return Filter(&GrokFilter{
		g:       g,
		fields:  fields,
		pattern: pattern,
	})
}

// Filter filters a log entry
func (gf *GrokFilter) Filter(data string) (string, error) {
	var log map[string]interface{}
	err := json.Unmarshal([]byte(data), &log)
	if err != nil {
		return "", err
	}

	values, err := gf.g.Parse(gf.pattern, log["message"].(string))
	if err != nil {
		return "", err
	}

	if len(gf.fields) <= 0 {
		for k, v := range values {
			if k != "timestamp" {
				log[k] = v
			}
		}
	} else if len(gf.fields) >= 1 {
		for k, v := range values {
			if k != "timestamp" && stringInSlice(k, gf.fields) {
				log[k] = v
			}
		}
	}

	out, err := json.Marshal(log)
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
