package filters

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
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
func (rf *RegexFilter) Filter(data string) string {
	if len(rf.regex) > 0 {
		m, err := regexp.MatchString(rf.regex, data)
		if err != nil {
			log.Error(err)
		}
		if m != true {
			return data
		}
		return ""
	}
	return data
}
