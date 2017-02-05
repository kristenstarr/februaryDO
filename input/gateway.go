package input

import (
	"bufio"
	"net"
	"github.com/kristenfelch/pkgindexer/logging"
)

// MessageGateway is responsible for communication with the client through sockets.
// It is passed a channel when Opened, and is responsible for listening on a port and accepting
// TCP socket connections.  It handles these connections each in a separate GoRoutine,
// reading messages received through the socket and passing them back through the channel for processing.
// Gateway also adds a small (1 element limit) channel to each message to receive its return value,
// to be passed back to the client.
type MessageGateway interface {
	Open(chan<- *ValidatedMessage) (opened bool, err error)

	Close() (closed bool, err error)
}

// SimpleMessageGateway is a MessageGateway that has an optional rate limit.
type SimpleMessageGateway struct {
	validator Validator
	rate *int
	logger logging.Logger
}

// ValidatedMessage contains an input message as well as a channel created to receive the
// result of processing this message.
type ValidatedMessage struct {
	*InputMessage
	ResponseChannel chan<- string
}

// Open starts listening on a Port and accepting connections.
func (s *SimpleMessageGateway) Open(c chan<- *ValidatedMessage) (opened bool, err error) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		s.logger.Error("Error starting on 8080")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			s.logger.Error(err.Error())
		} else {
			go s.handleConnection(conn, c)
		}
	}
}

// handleConnection reads messages through the TCP connection, rate limiting if desired.
func (s *SimpleMessageGateway) handleConnection(conn net.Conn, c chan<- *ValidatedMessage) {
	throttler := NewThrottler(s.rate)
	for {
		throttler.Next()
		message, msgError := bufio.NewReader(conn).ReadString('\n')
		if (msgError != nil) {
			s.logger.Debug("No more messages available from connection")
			conn.Close()
			throttler.Stop()
			break;
		}
		s.handleMessage(conn, message, c)
	}
}

// handleMessage validates our input message.  If it is valid, it is returned to the ValidatedMessage
// channel with it's own length-1 channel to contain the final result of processing the message.
func (s *SimpleMessageGateway) handleMessage(conn net.Conn, message string, c chan<- *ValidatedMessage) {
	validated, validatedError := s.validator.ValidateInput(message)
	if validatedError != nil {
		s.logger.Debug(validatedError.Error())
		conn.Write(s.formatResponse("error"))
	} else {
		ch := make(chan string, 1)
		validMessage := &ValidatedMessage{
			validated,
			ch,
		}
		c <- validMessage
		returned := <-ch
		close(ch)
		conn.Write(s.formatResponse(returned))
	}
}

// formatResponse formats our generic 'ok', 'fail', and 'error' into format that clients receive.
func (s *SimpleMessageGateway) formatResponse(str string) (resp []byte) {
	switch str {
	case "ok":
		return []byte("OK\n")
	case "fail":
		return []byte("FAIL\n")
	case "error":
		return []byte("ERROR\n")
	}
	return []byte("ERROR\n")
}

func (s *SimpleMessageGateway) Close() (closed bool, err error) {
	return true, nil
}

// NewMessageGateway create an instance of MessageGateway including validator and throttler.
func NewMessageGateway(throttle *int, logger logging.Logger) MessageGateway {
	return &SimpleMessageGateway{
		NewValidator(),
		throttle,
		logger,
	}
}
