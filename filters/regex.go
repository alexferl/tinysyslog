package filters

import (
	"regexp"
)

// Regex represents a regex filter
type Regex struct {
	regex string
	kind  Kind
}

// NewRegex creates a Regex instance
func NewRegex(s string) Filter {
	return Filter(&Regex{regex: s, kind: RegexKind})
}

// Filter filters a log entry
func (r *Regex) Filter(data string) (string, error) {
	if len(r.regex) > 0 {
		m, err := regexp.MatchString(r.regex, data)
		if err != nil {
			return "", err
		}
		if m != true {
			return data, nil
		}
		return "", nil
	}
	return data, nil
}

func (r *Regex) GetKind() Kind {
	return r.kind
}
