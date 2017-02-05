package operation

import (
	"testing"
	"github.com/kristenfelch/pkgindexer/data"
	"github.com/kristenfelch/pkgindexer/err"
	"strings"
	"github.com/kristenfelch/pkgindexer/logging"
)

// Tests case where index is successful because all dependencies are present.
func TestIndexSuccessful(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, true, nil, true, nil)
	logLevel := "FATAL"
	indexer := &SimpleIndexer{store, logging.NewIndexLogger(&logLevel)}

	indexed, err := indexer.Index("lib", []string{"dep1", "dep2"})
	if (err != nil || !indexed) {
		t.Error("When all dependencies are present, package should be successfully indexed")
	}
}

// Tests case where indexing does not occur because dependencies are missing.
func TestDependencyMissing(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, false, nil, true, nil)
	logLevel := "FATAL"
	indexer := &SimpleIndexer{store, logging.NewIndexLogger(&logLevel)}

	indexed, err := indexer.Index("lib", []string{"dep1", "dep2"})
	if (err != nil || indexed) {
		t.Error("When dependency is missing, should not be able to index")
	}
}

// Tests case where there is an error checking dependencies, and therefore an error indexing is propagated.
func TestErrorCheckingDependencies(t *testing.T) {
	store := data.NewTestStore(false, nil, true, nil, false, err.NewIndexError("Error checking for dependencies"), true, nil)
	logLevel := "FATAL"
	indexer := &SimpleIndexer{store, logging.NewIndexLogger(&logLevel)}

	indexed, err := indexer.Index("lib", []string{"dep1", "dep2"})
	if (strings.Index(err.Error(), "Error checking for dependencies") == -1 || indexed) {
		t.Error("Error checking for dependencies should be propagated, and indexing should not take place")
	}
}
