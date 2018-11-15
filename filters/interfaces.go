package filters

// Filter is a common interface for all filters
type Filter interface {
	Filter(string) (string, error)
}
