package filters

// NoOp represents a no operation filter
type NoOp struct{}

// NewNoOp creates a NoOp instance
func NewNoOp() Filter {
	return Filter(&NoOp{})
}

// Filter filters a log entry
func (n *NoOp) Filter(data string) (string, error) {
	return data, nil
}
