package filters

// NoOp represents a no operation filter
type NoOp struct {
	kind Kind
}

// NewNoOp creates a NoOp instance
func NewNoOp() Filter {
	return Filter(&NoOp{kind: NoOpKind})
}

// Filter filters a log entry
func (n *NoOp) Filter(data string) (string, error) {
	return data, nil
}

func (n *NoOp) GetKind() Kind {
	return n.kind
}
