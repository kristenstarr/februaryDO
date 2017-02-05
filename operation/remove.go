package operation

import (
	"github.com/kristenfelch/pkgindexer/data"
	"github.com/kristenfelch/pkgindexer/logging"
)

// Remover is responsible for removing Libraries from the index.
// It encapsulates the business logic behind removal and conditions required for removal,
// leaving removal itself to IndexStore.
type Remover interface {
	//removed indicates if element was removed, err if we tried and failed.
	Remove(name string) (removed bool, err error)
}

type SimpleRemover struct {
	store data.IndexStore
	logger logging.Logger
}

func (s *SimpleRemover) Remove(name string) (removed bool, err error) {
	lib, libError := s.store.HasPackage(name)
	if libError != nil {
		s.logger.Error(libError.Error())
		return false, libError
	}
	if lib {
		hasParents, hasParentsError := s.store.HasParents(name)
		if hasParentsError != nil {
			s.logger.Error(hasParentsError.Error())
			return false, hasParentsError
		}
		if hasParents {
			return false, nil
		} else {
			removed, removedErr := s.store.RemovePackage(name)
			if removedErr != nil {
				s.logger.Error(removedErr.Error())
				return false, removedErr
			}
			return removed, nil
		}
	} else {
		return true, nil
	}
}

// NewRemover creates a new Remover referencing our Index data store and logger.
func NewRemover(store data.IndexStore, logger logging.Logger) Remover {
	return &SimpleRemover{
		store,
		logger,
	}
}
