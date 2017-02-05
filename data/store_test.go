package data

import (
	"testing"
)

// Tests simple addition of library to index.
func TestAddLibrary(t *testing.T) {
	store := NewIndexStore()
	added, err := store.AddLibrary("library", []string{"dep1", "dep2"})
	if (err != nil || !added) {
		t.Errorf("Error encountered adding library : %s", err.Error())
	}
}

// Tests adding of library and that new library is added to parents' lists of its dependencies.
func TestAddLibraryToDependenciesParents(t *testing.T) {
	store := NewIndexStore()
	store.AddLibrary("dep1", nil)
	store.AddLibrary("dep2", nil)

	hasParents, err := store.HasParents("dep1")
	if (err != nil || hasParents) {
		t.Errorf("Library should not have any parents when initially indexed")
	}
	store.AddLibrary("library", []string{"dep1", "dep2"})
	hasParents, err = store.HasParents("dep1")
	if (err != nil || !hasParents) {
		t.Error("Library should have parents once another library depends on it")
	}
	hasParents, err = store.HasParents("dep2")
	if (err != nil || !hasParents) {
		t.Error("Library should have parents once another library depends on it")
	}
}

// Tests simple removal of a library that is already not indexed.
func TestRemoveNonIndexedLibrary(t *testing.T) {
	store := NewIndexStore()
	removed, err := store.RemoveLibrary("dep1")
	if (err != nil || !removed) {
		t.Error("Remove should return with no error if library is not present")
	}
}

// Tests that when library is removed, it is also removed from parents list of its dependencies.
func TestRemoveIndexedWithDeps(t *testing.T) {
	store := NewIndexStore()
	store.AddLibrary("dep1", nil)
	store.AddLibrary("library", []string{"dep1"})
	//dependency should have library as a parent
	hasParents, err := store.HasParents("dep1")
	if (err != nil || !hasParents) {
		t.Error("Library should have parents once another library depends on it")
	}
	removed, err := store.RemoveLibrary("library")
	if (err != nil || !removed) {
		t.Error("Remove should return with no error")
	}
	//once parent is removed, dependency should have it removed from its parent list.
	hasParents, err = store.HasParents("dep1")
	if (err != nil || hasParents) {
		t.Error("Library should have had parent removed")
	}
}

// Tests simple querying for library once it has been added.
func TestHasLibrary(t *testing.T) {
	store := NewIndexStore()
	exists, err := store.HasLibrary("library")
	if (err != nil || exists) {
		t.Error("Has Library should return false before library has been indexed")
	}
	store.AddLibrary("library", nil)
	exists, err = store.HasLibrary("library")
	if (err != nil || !exists) {
		t.Error("Has Library should return true once library has been indexed")
	}

}
