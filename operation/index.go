package operation

import (
	"github.com/kristenfelch/pkgindexer/data"
	"github.com/kristenfelch/pkgindexer/logging"
)

// Indexer is responsible for adding a Package to our Index.
// It encapsulates business logic for knowing when a Package can be added,
// leaving actual storage of index to IndexStore.
type Indexer interface {
	// indicates if element was Indexed, err if we tried and failed.
	Index(name string, dependencies []string) (Indexed bool, err error)
}

type SimpleIndexer struct {
	store data.IndexStore
	logger logging.Logger
}

func (s *SimpleIndexer) Index(name string, dependencies []string) (Indexed bool, err error) {
	//// Check to see if all dependencies are present
	for _, dep := range dependencies {
		lib, libError := s.store.HasPackage(dep)
		if libError != nil {
			//error determining if dependency is there, for indexing
			s.logger.Error(libError.Error())
			return false, libError
		}
		if !lib {
			//dependency is missing, not indexed
			return false, nil
		}
	}

	exists, existsErr := s.store.HasPackage(name)
	if existsErr != nil {
		// error looking up existing indexed package
		s.logger.Error(existsErr.Error())
		return false, existsErr
	}
	if exists {
		//remove package with old dependencies
		s.store.RemovePackage(name)

	}

	return s.store.AddPackage(name, dependencies)

}

// NewIndexer creates a new Indexer referencing our Index data store and a logger.
func NewIndexer(store data.IndexStore, logger logging.Logger) Indexer {
	return &SimpleIndexer{
		store,
		logger,
	}
}
