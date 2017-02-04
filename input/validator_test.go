package input

import "testing"

// Tests correct Query message.
func TestCorrectQuery(t *testing.T) {
	validator := NewValidator()
	validQuery := "QUERY|lib|deps\n"
	result, err := validator.ValidateInput(validQuery)
	validateMessage(t, result, err, "QUERY", "lib", "deps")

}

// Tests correct Index message.
func TestCorrectIndex(t *testing.T) {
	validator := NewValidator()
	validQuery := "INDEX|lib|\n"
	result, err := validator.ValidateInput(validQuery)
	validateMessage(t, result, err, "INDEX", "lib", "")
}

// Tests correct Remove message.
func TestCorrectRemove(t *testing.T) {
	validator := NewValidator()
	validQuery := "REMOVE|lib|dep1,dep2\n"
	result, err := validator.ValidateInput(validQuery)
	validateMessage(t, result, err, "REMOVE", "lib", "dep1,dep2")
}

// Tests incorrect piping in input.
func TestBadFormat(t *testing.T) {
	validator := NewValidator()
	badQuery := "NO_PIPES"
	_, err := validator.ValidateInput(badQuery)
	if (err.Error() != "Input does not have 3 arguments : NO_PIPES") {
		t.Errorf("Incorrect error message : %s", err.Error())
	}
}

// Tests Bad Verb.
func TestBadVerb(t *testing.T) {
	validator := NewValidator()
	badQuery := "FAKE|lib|\n"
	_, err := validator.ValidateInput(badQuery)
	if (err.Error() != "Input method should be REMOVE/INDEX/QUERY, not : FAKE") {
		t.Errorf("Incorrect error message : %s", err.Error())
	}
}

// Tests Bad Library.
func TestBadLibrary(t *testing.T) {
	validator := NewValidator()
	badQuery := "QUERY|l*ib|\n"
	_, err := validator.ValidateInput(badQuery)
	if (err.Error() != "Library name missing or incorrect : l*ib") {
		t.Errorf("Incorrect error message : %s", err.Error())
	}
}

// Tests Bad Dependencies.
func TestBadDependencies(t *testing.T) {
	validator := NewValidator()
	badQuery := "QUERY|lib|c&\n"
	_, err := validator.ValidateInput(badQuery)
	if (err.Error() != "Dependencies are incorrectly formatted : c&") {
		t.Errorf("Incorrect error message : %s", err.Error())
	}
}

func validateMessage(t *testing.T, result *InputMessage, err error, expectedVerb, expectedLib, expectedDeps string) {
	if (err != nil) {
		t.Error(err)
	}
	if (result.Verb != expectedVerb) {
		t.Errorf("Incorrect message verb parsed : %s", result.Verb)
	}
	if (result.Library != expectedLib) {
		t.Errorf("Incorrect message library parsed : %s", result.Library)
	}
	if (result.Dependencies != expectedDeps) {
		t.Errorf("Incorrect message dependencies parsed : %s", result.Dependencies)
	}
}
