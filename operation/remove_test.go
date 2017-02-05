package operation

import (
	"testing"
	"github.com/kristenfelch/pkgindexer/data"
	"github.com/kristenfelch/pkgindexer/err"
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
	store := data.NewTestStore(true, nil, true, nil, false, err.New("Error looking up library"), true, nil)
	remover := &SimpleRemover{store}

	removed, err := remover.Remove("lib")
	if (err.Error() != "Error looking up library" || removed) {
		t.Error("When we have an error looking up the library, it is propagated")
	}
}

// Tests case where determining if library has parents throws an error, which is propagated.
func TestRemoveErrorHasParents(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, true, nil, true, err.New("Error looking up parents"))
	remover := &SimpleRemover{store}

	removed, err := remover.Remove("lib")
	if (err.Error() != "Error looking up parents" || removed) {
		t.Error("When we have an error looking up the library's parents, it is propagated")
	}
}

// Tests case where removal of library causes an error.
func TestRemoveError(t *testing.T) {
	store := data.NewTestStore(true, nil, false, err.New("Error removing"), true, nil, false, nil)
	remover := &SimpleRemover{store}

	removed, err := remover.Remove("lib")
	if (err.Error() != "Error removing" || removed) {
		t.Error("Error removing should be thrown")
	}
}
