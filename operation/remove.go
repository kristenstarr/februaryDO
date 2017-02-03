package operation

import "github.com/kristenfelch/pkgindexer/data"

// Remover is responsible for removing Libraries from the index.
// It encapsulates the business logic behind removal and conditions required for removal.
type Remover interface {
	//removed indicates if element was removed, err if we tried and failed.
	Remove(name string) (removed bool, err error)
}

type SimpleRemover struct {
	store data.IndexStore
}

func (s *SimpleRemover) Remove(name string) (removed bool, err error) {
	lib, libError := s.store.HasLibrary(name)
	if libError != nil {
		return false, libError
	}
	if lib {
		hasParents, hasParentsError := s.store.HasParents(name)
		if hasParentsError != nil {
			return false, hasParentsError
		}
		if hasParents {
			return false, nil
		} else {
			removed, removedErr := s.store.RemoveLibrary(name)
			if removedErr != nil {
				return false, removedErr
			}
			return removed, nil
		}
	} else {
		return true, nil
	}
}

func NewRemover(store data.IndexStore) Remover {
	return &SimpleRemover{store}
}
