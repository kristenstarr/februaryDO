package input

import (
	"github.com/kristenfelch/pkgindexer/err"
	"strings"
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
	if len(pieces) != 3 {
		return nil, err.New("Input does not have 3 arguments")
	}
	method := pieces[0]
	if method != "REMOVE" && method != "INDEX" && method != "QUERY" {
		return nil, err.New("Input method should be REMOVE/INDEX/QUERY")
	}
	lib := pieces[1]
	if len(lib) < 1 {
		return nil, err.New("Library name missing from input")
	}
	return &InputMessage{
		method,
		lib,
		pieces[2][:len(pieces[2])-1],
	}, nil
}

func NewValidator() Validator {
	return &SimpleValidator{}
}
