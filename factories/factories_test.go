package factories

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"tinysyslog/filters"
	"tinysyslog/mutators"
	"tinysyslog/sinks"
)

func TestMutator(t *testing.T) {
	m := Mutator()
	assert.Equal(t, mutators.TextKind, m.GetKind())
}

func TestFilter(t *testing.T) {
	f := Filter()
	assert.Equal(t, filters.NoOpKind, f.GetKind())
}

func TestSinks(t *testing.T) {
	s := Sinks()
	assert.Equal(t, []sinks.Sink(nil), s)
}
