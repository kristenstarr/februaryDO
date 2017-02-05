package operation

import (
	"testing"
	"github.com/kristenfelch/pkgindexer/data"
	"github.com/kristenfelch/pkgindexer/err"
	"strings"
)

// Tests case where remove is successful because no libraries depend on this one.
func TestRemoveSuccessfulNoParents(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, true, nil, false, nil)
	remover := &SimpleRemover{store}

	removed, err := remover.Remove("lib")
	if (err != nil || !removed) {
		t.Error("When library is present with no others that depend on it, it can be removed")
	}
}

// Tests case where remove request is successful because library was not present to start with.
func TestRemoveSuccessfulAlreadyRemoved(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, false, nil, true, nil)
	remover := &SimpleRemover{store}

	removed, err := remover.Remove("lib")
	if (err != nil || !removed) {
		t.Error("When library has already been removed, remove should return true with no error")
	}
}

// Tests case where remove request fails because other libraries depend on the one being removed.
func TestRemoveFailDependenciesPresent(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, true, nil, true, nil)
	remover := &SimpleRemover{store}

	removed, err := remover.Remove("lib")
	if (err != nil || removed) {
		t.Error("When library has others that depend on it, it cannot be removed")
	}
}

// Tests case where determining if library is present throws an error, which is propagated.
func TestRemoveErrorHasLibrary(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, false, err.NewIndexError("Error looking up library"), true, nil)
	remover := &SimpleRemover{store}

	removed, err := remover.Remove("lib")
	if (strings.Index(err.Error(), "Error looking up library") != 0 || removed) {
		t.Error("When we have an error looking up the library, it is propagated")
	}
}

// Tests case where determining if library has parents throws an error, which is propagated.
func TestRemoveErrorHasParents(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, true, nil, true, err.NewIndexError("Error looking up parents"))
	remover := &SimpleRemover{store}

	removed, err := remover.Remove("lib")
	if (strings.Index(err.Error(), "Error looking up parents") != 0 || removed) {
		t.Error("When we have an error looking up the library's parents, it is propagated")
	}
}

// Tests case where removal of library causes an error.
func TestRemoveError(t *testing.T) {
	store := data.NewTestStore(true, nil, false, err.NewIndexError("Error removing"), true, nil, false, nil)
	remover := &SimpleRemover{store}

	removed, err := remover.Remove("lib")
	if (strings.Index(err.Error(), "Error removing") != 0 || removed) {
		t.Error("Error removing should be thrown")
	}
}
