package operation

import (
	"testing"
	"github.com/kristenfelch/pkgindexer/data"
)

// Tests case where query returns true
func TestQueryTrue(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, true, nil, true, nil)
	querier := &SimpleQuerier{store}

	hasIndex, err := querier.Query("lib")
	if (err != nil || !hasIndex) {
		t.Error("When library is present in index, true should be returned with no error")
	}
}

// Tests case where query returns false
func TestQueryFalse(t *testing.T) {
	store := data.NewTestStore(true, nil, true, nil, false, nil, true, nil)
	querier := &SimpleQuerier{store}

	hasIndex, err := querier.Query("lib")
	if (err != nil || hasIndex) {
		t.Error("When library is not present in index, false should be returned with no error")
	}
}
