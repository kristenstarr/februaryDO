package data

import (
	"testing"
	"github.com/kristenfelch/pkgindexer/logging"
)

// Tests simple addition of package to index.
func TestAddPackage(t *testing.T) {
	logLevel := "FATAL"
	store := NewIndexStore(logging.NewIndexLogger(&logLevel))
	added, err := store.AddPackage("package", []string{"dep1", "dep2"})
	if (err != nil || !added) {
		t.Errorf("Error encountered adding package : %s", err.Error())
	}
}

// Tests adding of package and that new package is added to parents' lists of its dependencies.
func TestAddPackageToDependenciesParents(t *testing.T) {
	logLevel := "FATAL"
	store := NewIndexStore(logging.NewIndexLogger(&logLevel))
	store.AddPackage("dep1", nil)
	store.AddPackage("dep2", nil)

	hasParents, err := store.HasParents("dep1")
	if (err != nil || hasParents) {
		t.Errorf("Package should not have any parents when initially indexed")
	}
	store.AddPackage("package", []string{"dep1", "dep2"})
	hasParents, err = store.HasParents("dep1")
	if (err != nil || !hasParents) {
		t.Error("Package should have parents once another package depends on it")
	}
	hasParents, err = store.HasParents("dep2")
	if (err != nil || !hasParents) {
		t.Error("Package should have parents once another package depends on it")
	}
}

// Tests simple removal of a package that is already not indexed.
func TestRemoveNonIndexedPackage(t *testing.T) {
	logLevel := "FATAL"
	store := NewIndexStore(logging.NewIndexLogger(&logLevel))
	removed, err := store.RemovePackage("dep1")
	if (err != nil || !removed) {
		t.Error("Remove should return with no error if package is not present")
	}
}

// Tests that when package is removed, it is also removed from parents list of its dependencies.
func TestRemoveIndexedWithDeps(t *testing.T) {
	logLevel := "FATAL"
	store := NewIndexStore(logging.NewIndexLogger(&logLevel))
	store.AddPackage("dep1", nil)
	store.AddPackage("package", []string{"dep1"})
	//dependency should have package as a parent
	hasParents, err := store.HasParents("dep1")
	if (err != nil || !hasParents) {
		t.Error("Package should have parents once another package depends on it")
	}
	removed, err := store.RemovePackage("package")
	if (err != nil || !removed) {
		t.Error("Remove should return with no error")
	}
	//once parent is removed, dependency should have it removed from its parent list.
	hasParents, err = store.HasParents("dep1")
	if (err != nil || hasParents) {
		t.Error("Package should have had parent removed")
	}
}

// Tests simple querying for package once it has been added.
func TestHasPackage(t *testing.T) {
	logLevel := "FATAL"
	store := NewIndexStore(logging.NewIndexLogger(&logLevel))
	exists, err := store.HasPackage("package")
	if (err != nil || exists) {
		t.Error("Has Package should return false before package has been indexed")
	}
	store.AddPackage("package", nil)
	exists, err = store.HasPackage("package")
	if (err != nil || !exists) {
		t.Error("Has Package should return true once package has been indexed")
	}

}
