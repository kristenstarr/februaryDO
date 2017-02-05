package err

import (
	"time"
	"fmt"
)

// IndexError is a custom Error type used throughout application,
// which includes a timestamp on all error messages.
type IndexError struct {
	msg string
	time time.Time
}

func (e *IndexError) Error() string {
	return fmt.Sprintf("%s : %s", e.msg, e.time.Format(time.UnixDate))
}

// Creates a new IndexError including the current timestamp
func NewIndexError(text string) error {
	return &IndexError{
		text,
		time.Now(),
	}
}
