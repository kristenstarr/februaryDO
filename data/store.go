package data

import (
	"errors"
)

// IndexStore is responsible for storing the current state of our index.
// It encapsulates the lower level operations of adding, removing, and querying for
// libraries, and abstracts the data storage choice for our index.
type IndexStore interface {

	// Adds a Library to our Index.
	AddLibrary(name string, deps []string) (added bool, err error)

	// Removes a Library from our Index.
	RemoveLibrary(name string) (removed bool, err error)

	// Determines if Library has already been indexed.
	HasLibrary(name string) (exists bool, err error)

	// Determines if a Library has Parents - other libraries that depend on it.
	HasParents(name string) (hasParents bool, err error)
}

type MapsIndexStore struct {
	store map[string]*Library
}

// Library is a type of struct used to store our libraries that have been indexed.
// We store parents (libraries that depend on library) so that we can efficiently determine if
// a library can be removed, without iteration to determine if any libraries depend on the one in question.
// This provides the performance optimization of O(1) rather than O(n) for Remove operation.
// We store dependencies (libraries that this library depends on) also for efficient removal,
// so that when a library is removed we can remove it also from the Parents list of its dependencies.
// We are using map[string]bool for Dependencies and Parents for faster lookup time than a []string would provide.
type Library struct {
	Dependencies map[string]bool
	Parents      map[string]bool
}

func (l *Library) HasParents() bool {
	return len(l.Parents) > 0
}

func (m *MapsIndexStore) AddLibrary(name string, deps []string) (added bool, err error) {
	dependencies := make(map[string]bool, len(deps))
	for v := range deps {
		dependencies[deps[v]] = true
		if depLibrary, _ := m.getLibrary(deps[v]); depLibrary != nil {
			// add this library to each dependency's parents, so that we know we
			// cannot remove the dependency.
			depLibrary.Parents[name] = true
		}
	}
	m.store[name] = &Library{
		dependencies,
		// No libraries can depend on this one until after this one has been created
		// thus initialize with an empty list.
		make(map[string]bool),
	}
	return true, nil
}

func (m *MapsIndexStore) RemoveLibrary(name string) (removed bool, err error) {
	if lib, _ := m.getLibrary(name); lib != nil {
		delete(m.store, name)
		for key := range lib.Dependencies {
			if dependentLibrary, _ := m.getLibrary(key); dependentLibrary != nil {
				// remove this library from each dependency's parents, so that we know
				// we can remove the dependency if no others depend on it.
				delete(dependentLibrary.Parents, name)
			}
		}
	}
	return true, nil
}

func (m *MapsIndexStore) getLibrary(name string) (lib *Library, err error) {
	if lib, ok := m.store[name]; ok {
		return lib, nil
	} else {
		return nil, nil
	}
}

func (m *MapsIndexStore) HasLibrary(name string) (exists bool, err error) {
	if _, ok := m.store[name]; ok {
		return true, nil
	} else {
		return false, nil
	}
}

func (m *MapsIndexStore) HasParents(name string) (hasParents bool, err error) {
	if lib, _ := m.getLibrary(name); lib != nil {
		return lib.HasParents(), nil
	} else {
		return false, errors.New("Unable to determined if Unindexed library has parents")
	}
}

func NewIndexStore() IndexStore {
	return &MapsIndexStore{
		make(map[string]*Library),
	}
}
