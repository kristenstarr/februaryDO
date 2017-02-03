package input

import (
	"bufio"
	"fmt"
	"net"
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

type SimpleMessageGateway struct {
	validator Validator
}

type ValidatedMessage struct {
	*InputMessage
	ResponseChannel chan<- string
}

func (s *SimpleMessageGateway) Open(c chan<- *ValidatedMessage) (opened bool, err error) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting on 8080")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			go s.handleConnection(conn, c)
		}
	}
}

func (s *SimpleMessageGateway) handleConnection(conn net.Conn, c chan<- *ValidatedMessage) {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		validated, validatedError := s.validator.ValidateInput(message)
		if validatedError != nil {
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
}

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

func NewMessageGateway() MessageGateway {
	return &SimpleMessageGateway{
		NewValidator(),
	}
}
