package data

import (
	"github.com/kristenfelch/pkgindexer/err"
	"github.com/kristenfelch/pkgindexer/logging"
	"fmt"
)

// IndexStore is responsible for storing the current state of our index.
// It encapsulates the lower level operations of adding, removing, and querying for
// packages, and abstracts the data storage choice for our index.
type IndexStore interface {

	// Adds a Package to our Index.
	AddPackage(name string, deps []string) (added bool, error error)

	// Removes a Package from our Index.
	RemovePackage(name string) (removed bool, error error)

	// Determines if Package has already been indexed.
	HasPackage(name string) (exists bool, error error)

	// Determines if a Package has Parents - other packages that depend on it.
	HasParents(name string) (hasParents bool, error error)
}

type MapsIndexStore struct {
	store map[string]*Package
	logger logging.Logger
}

// Package is a type of struct used to store our packages that have been indexed.
// We store parents (packages that depend on package) so that we can efficiently determine if
// a package can be removed, without iteration to determine if any packages depend on the one in question.
// This provides the performance optimization of O(1) rather than O(n) for Remove operation.
// We store dependencies (packages that this package depends on) also for efficient removal,
// so that when a package is removed we can remove it also from the Parents list of its dependencies.
// We are using map[string]bool for Dependencies and Parents for faster lookup time than a []string would provide.
type Package struct {
	Dependencies map[string]bool
	Parents      map[string]bool
}

func (l *Package) HasParents() bool {
	return len(l.Parents) > 0
}

func (m *MapsIndexStore) AddPackage(name string, deps []string) (added bool, error error) {
	dependencies := make(map[string]bool, len(deps))
	for v := range deps {
		dependencies[deps[v]] = true
		if depPackage, _ := m.getPackage(deps[v]); depPackage != nil {
			// add this package to each dependency's parents, so that we know we
			// cannot remove the dependency.
			m.logger.Trace(fmt.Sprintf("Package %s added to dependencies of %s", name, deps[v]))
			depPackage.Parents[name] = true
		}
	}
	m.store[name] = &Package{
		dependencies,
		// No packages can depend on this one until after this one has been created
		// thus initialize with an empty list.
		make(map[string]bool),
	}
	m.logger.Trace(fmt.Sprintf("Package %s added to Index", name))
	return true, nil
}

func (m *MapsIndexStore) RemovePackage(name string) (removed bool, error error) {
	if lib, _ := m.getPackage(name); lib != nil {
		delete(m.store, name)
		for key := range lib.Dependencies {
			if dependentPackage, _ := m.getPackage(key); dependentPackage != nil {
				// remove this package from each dependency's parents, so that we know
				// we can remove the dependency if no others depend on it.
				m.logger.Trace(fmt.Sprintf("Package %s removed as parent of %s", name, key))
				delete(dependentPackage.Parents, name)
			}
		}
	}
	m.logger.Trace(fmt.Sprintf("Package %s removed from Index", name))
	return true, nil
}

func (m *MapsIndexStore) getPackage(name string) (lib *Package, error error) {
	if lib, ok := m.store[name]; ok {
		return lib, nil
	} else {
		return nil, nil
	}
}

func (m *MapsIndexStore) HasPackage(name string) (exists bool, error error) {
	if _, ok := m.store[name]; ok {
		return true, nil
	} else {
		return false, nil
	}
}

func (m *MapsIndexStore) HasParents(name string) (hasParents bool, error error) {
	if lib, _ := m.getPackage(name); lib != nil {
		return lib.HasParents(), nil
	} else {
		return false, err.NewIndexError("Unable to determined if Unindexed package has parents")
	}
}

func NewIndexStore(logger logging.Logger) IndexStore {
	return &MapsIndexStore{
		make(map[string]*Package),
		logger,
	}
}
