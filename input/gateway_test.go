package input

import (
	"testing"
	"bytes"
	"github.com/kristenfelch/pkgindexer/logging"
)

// TestingGateway mocks the portion of our Gateway that deals with the network.
// Since this gateway is only responsible for accepting messages and passing them back
// through channels, we test the send of messages and that they are validated, as well
// as the responses sent back.
type TestingGateway struct {
	SimpleMessageGateway
}

func (s *TestingGateway) Open(c chan<- *ValidatedMessage) (opened bool, err error) {
	conn := NewTestConnection()
	s.handleMessage(conn, "QUERY|lib|dep\n", c)
	return true, nil
}

func NewTestingGateway() MessageGateway {
	throttle := 0
	logLevel := "FATAL"
	return &TestingGateway{
		SimpleMessageGateway{
			NewValidator(),
			&throttle,
			logging.NewIndexLogger(&logLevel),
		},
	}
}

// Tests that when a single message is sent, it is sent back in validated message format.
func TestGatewayReceiveMessage(t *testing.T) {
	gateway := NewTestingGateway()
	msgChannel := make(chan *ValidatedMessage, 1)
	go gateway.Open(msgChannel)
	val := <- msgChannel
	val.ResponseChannel <- "OK"
	close(msgChannel)
	if (val.Dependencies != "dep" || val.Package != "lib" || val.Verb != "QUERY") {
		t.Errorf("Invalid message returned through channel")
	}
}

// Tests basic formatting of responses to client.
func TestGatewayFormatResponse(t *testing.T) {
	throttle := 0
	logLevel := "FATAL"
	gateway := &SimpleMessageGateway{
		NewValidator(),
		&throttle,
		logging.NewIndexLogger(&logLevel),
	}
	formatted := gateway.formatResponse("ok")
	if (!bytes.Equal(formatted, []byte("OK\n"))) {
		t.Error("Incorrect OK response formatting")
	}
	formatted = gateway.formatResponse("error")
	if (!bytes.Equal(formatted, []byte("ERROR\n"))) {
		t.Error("Incorrect ERROR response formatting")
	}
	formatted = gateway.formatResponse("fail")
	if (!bytes.Equal(formatted, []byte("FAIL\n"))) {
		t.Error("Incorrect FAIL response formatting")
	}
	formatted = gateway.formatResponse("garbage")
	if (!bytes.Equal(formatted, []byte("ERROR\n"))) {
		t.Error("Incorrect ERROR response formatting")
	}
}
