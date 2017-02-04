package input

import (
	"github.com/kristenfelch/pkgindexer/err"
	"strings"
	"regexp"
	"fmt"
)

// Validator is responsible for validating input format of received messages.
type Validator interface {
	ValidateInput(input string) (validMessage *InputMessage, err error)
}

type SimpleValidator struct{}

type InputMessage struct {
	Verb         string
	Library      string
	Dependencies string
}

func (s *SimpleValidator) ValidateInput(input string) (validMessage *InputMessage, error error) {
	pieces := strings.Split(input, "|")

	// First ensure that we have the 3 required parts to our input.
	if len(pieces) != 3 {
		return nil, err.New(fmt.Sprintf("Input does not have 3 arguments : %s", input))
	}

	//Ensure that our request type is REMOVE/INDEX/QUERY
	method := pieces[0]
	if method != "REMOVE" && method != "INDEX" && method != "QUERY" {
		return nil, err.New(fmt.Sprintf("Input method should be REMOVE/INDEX/QUERY, not : %s", method))
	}

	//Make sure our lib name is >1 alphanumeric character
	lib := pieces[1]
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_\-\+]+$`, lib)
	if (!match) {
		return nil, err.New(fmt.Sprintf("Library name missing or incorrect : %s", lib))
	}

	//Make sure that our dependencies list is a comma delimited list of alphanumeric words.
	dependencies := pieces[2][:len(pieces[2])-1]
	match, _ = regexp.MatchString(`^[a-zA-Z0-9_,\-\+]*$`, dependencies)
	if !match {
		return nil, err.New(fmt.Sprintf("Dependencies are incorrectly formatted : %s", dependencies))
	}

	return &InputMessage{
		method,
		lib,
		dependencies,
	}, nil
}

func NewValidator() Validator {
	return &SimpleValidator{}
}
