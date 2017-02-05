package operation

import (
	"testing"
	"github.com/kristenfelch/pkgindexer/data"
	"github.com/kristenfelch/pkgindexer/err"
	"strings"
	"github.com/kristenfelch/pkgindexer/logging"
)

// Tests case where remove is successful because no packages depend on this one.
func TestRemoveSuccessfulNoParents(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, true, nil, false, nil)
	logLevel := "FATAL"
	remover := &SimpleRemover{store, logging.NewIndexLogger(&logLevel)}

	removed, err := remover.Remove("lib")
	if (err != nil || !removed) {
		t.Error("When package is present with no others that depend on it, it can be removed")
	}
}

// Tests case where remove request is successful because package was not present to start with.
func TestRemoveSuccessfulAlreadyRemoved(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, false, nil, true, nil)
	logLevel := "FATAL"
	remover := &SimpleRemover{store, logging.NewIndexLogger(&logLevel)}

	removed, err := remover.Remove("lib")
	if (err != nil || !removed) {
		t.Error("When package has already been removed, remove should return true with no error")
	}
}

// Tests case where remove request fails because other packages depend on the one being removed.
func TestRemoveFailDependenciesPresent(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, true, nil, true, nil)
	logLevel := "FATAL"
	remover := &SimpleRemover{store, logging.NewIndexLogger(&logLevel)}

	removed, err := remover.Remove("lib")
	if (err != nil || removed) {
		t.Error("When package has others that depend on it, it cannot be removed")
	}
}

// Tests case where determining if package is present throws an error, which is propagated.
func TestRemoveErrorHasPackage(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, false, err.NewIndexError("Error looking up package"), true, nil)
	logLevel := "FATAL"
	remover := &SimpleRemover{store, logging.NewIndexLogger(&logLevel)}

	removed, err := remover.Remove("lib")
	if (strings.Index(err.Error(), "Error looking up package") == -1 || removed) {
		t.Error("When we have an error looking up the package, it is propagated")
	}
}

// Tests case where determining if package has parents throws an error, which is propagated.
func TestRemoveErrorHasParents(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, true, nil, true, err.NewIndexError("Error looking up parents"))
	logLevel := "FATAL"
	remover := &SimpleRemover{store, logging.NewIndexLogger(&logLevel)}

	removed, err := remover.Remove("lib")
	if (strings.Index(err.Error(), "Error looking up parents") == -1 || removed) {
		t.Error("When we have an error looking up the package's parents, it is propagated")
	}
}

// Tests case where removal of package causes an error.
func TestRemoveError(t *testing.T) {
	store := data.NewTestStore(true, nil, false, err.NewIndexError("Error removing"), true, nil, false, nil)
	logLevel := "FATAL"
	remover := &SimpleRemover{store, logging.NewIndexLogger(&logLevel)}

	removed, err := remover.Remove("lib")
	if (strings.Index(err.Error(), "Error removing") == -1 || removed) {
		t.Error("Error removing should be thrown")
	}
}
