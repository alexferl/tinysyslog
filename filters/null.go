package filters

// NullFilter represents a null filter
type NullFilter struct{}

// NewNullFilter creates a NullFilter instance
func NewNullFilter() Filter {
	return Filter(&NullFilter{})
}

// Filter filters a log entry
func (nf *NullFilter) Filter(data string) (string, error) {
	return data, nil
}
