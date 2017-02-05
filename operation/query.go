package operation

import (
	"github.com/kristenfelch/pkgindexer/data"
)

// Querier is responsible for Querying if a Package is indexed.
type Querier interface {
	// indicates if element is currently indexed.
	Query(name string) (indexed bool, err error)
}

type SimpleQuerier struct {
	store data.IndexStore
}

func (s *SimpleQuerier) Query(name string) (indexed bool, err error) {
	return s.store.HasPackage(name)
}

// NewQuerier creates a new Querier referencing our Index data store.
func NewQuerier(store data.IndexStore) Querier {
	return &SimpleQuerier{store}
}
