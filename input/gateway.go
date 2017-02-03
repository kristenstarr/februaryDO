package input

import (
	"fmt"
	"net"
	//"time"
	"bufio"
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
	rate *int
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
			conn.Close()
			fmt.Println(err)
		} else {
			go s.handleConnection(conn, c)
		}
	}
}

func (s *SimpleMessageGateway) handleConnection(conn net.Conn, c chan<- *ValidatedMessage) {

	//var ticker *time.Ticker
	//if (*s.rate > 0) {
	//	ticker = time.NewTicker(time.Second / time.Duration(*s.rate))
	//}

	for {
		//if ticker != nil {
		//	<-ticker.C
		//}
		_, messageError := bufio.NewReader(conn).ReadString('\n')

		// If we have an error reading from socket, close and break loop/goroutine
		if (messageError != nil) {
			conn.Write([]byte("ERROR\n"))
			conn.Close()
			//if (ticker != nil) {
			//	ticker.Stop()
			//}
			break
		} else {

			//validated, validatedError := s.validator.ValidateInput(message)
			//if validatedError != nil {
				conn.Write(s.formatResponse("ok"))
			//} else {
			//	ch := make(chan string, 1)
			//	validMessage := &ValidatedMessage{
			//		validated,
			//		ch,
			//	}
			//	c <- validMessage
			//	returned := <-ch
			//	close(ch)
			//	conn.Write(s.formatResponse(returned))
			//}
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

func NewMessageGateway(throttle *int) MessageGateway {
	return &SimpleMessageGateway{
		NewValidator(),
		throttle,
	}
}
