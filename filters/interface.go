package filters

type Kind int8

const (
	NoOpKind Kind = iota + 1
	RegexKind
)

func (k Kind) String() string {
	return [...]string{"noop", "regex"}[k-1]
}

var Kinds = []string{NoOpKind.String(), RegexKind.String()}

// Filter is a common interface for all filters
type Filter interface {
	Filter(string) (string, error)
	GetKind() Kind
}
