package err

import (
	"testing"
	"strings"
)

// Tests basic error creation.
func TestErrCreation(t *testing.T) {
	err := NewIndexError("error test")
	if (strings.Index(err.Error(), "error test") == -1) {
		t.Error("Index Error should include error text")
	}
}
