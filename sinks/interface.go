package sinks

type Kind int8

const (
	ConsoleKind Kind = iota + 1
	ElasticsearchKind
	FilesystemKind
)

func (k Kind) String() string {
	return [...]string{"console", "elasticsearch", "filesystem"}[k-1]
}

var Kinds = []string{ConsoleKind.String(), ElasticsearchKind.String(), FilesystemKind.String()}

// Sink a common interface for all sinks
type Sink interface {
	Write([]byte) error
	GetKind() Kind
}
