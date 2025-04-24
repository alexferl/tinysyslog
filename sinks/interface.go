package sinks

type Kind int8

const (
	ConsoleKind Kind = iota + 1
	FilesystemKind
)

func (k Kind) String() string {
	return [...]string{"console", "filesystem"}[k-1]
}

var Kinds = []string{ConsoleKind.String(), FilesystemKind.String()}

// Sink a common interface for all sinks
type Sink interface {
	Write([]byte) error
	GetKind() Kind
}
