package mutators

// Mutator is a common interface for all mutators
type Mutator interface {
	Mutate(map[string]interface{}) string
}
