package filters

import (
	"regexp"
)

// RegexFilter represents a regex filter
type RegexFilter struct {
	regex string
}

// NewRegexFilter creates a RegexFilter instance
func NewRegexFilter(s string) Filter {
	return Filter(&RegexFilter{regex: s})
}

// Filter filters a log entry
func (rf *RegexFilter) Filter(data string) (string, error) {
	if len(rf.regex) > 0 {
		m, err := regexp.MatchString(rf.regex, data)
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
