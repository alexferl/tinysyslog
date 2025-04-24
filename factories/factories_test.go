package factories

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexferl/tinysyslog/filters"
	"github.com/alexferl/tinysyslog/mutators"
	"github.com/alexferl/tinysyslog/sinks"
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
