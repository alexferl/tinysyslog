package mutators

type Kind int8

const (
	TextKind Kind = iota + 1
	JSONKind
)

func (k Kind) String() string {
	return [...]string{"text", "json"}[k-1]
}

var Kinds = []string{TextKind.String(), JSONKind.String()}

// Mutator is a common interface for all mutators
type Mutator interface {
	Mutate(Log) (string, error)
	GetKind() Kind
}
