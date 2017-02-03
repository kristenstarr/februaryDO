package err

// ErrorString is a custom Error type used throughout application.
type ErrorString struct {
	s string
}

func (e *ErrorString) Error() string {
	return e.s
}

func New(text string) error {
	return &ErrorString{text}
}
