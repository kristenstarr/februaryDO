package integration

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

//ResponseCode is the code returned by the server as a response to our requests
type ResponseCode string

const (
	OK = "OK"
	FAIL = "FAIL"
	ERROR = "ERROR"
	UNKNOWN = "UNKNOWN"
)

//PackageIndexerClient sends messages to a running server.
//NOTE:  This interface has been taken from the load test test harness, to be used in integration tests.
type PackageIndexerClient interface {
	Close() error
	Send(msg string) (ResponseCode, error)
}

// TCPPackageIndexerClient connects to the running server via TCP
type TCPPackageIndexerClient struct {
	conn net.Conn
}

//Close closes the connection to the server.
func (client *TCPPackageIndexerClient) Close() error {
	return client.conn.Close()
}

//Send sends a message to the server using its line-oriented protocol
func (client *TCPPackageIndexerClient) Send(msg string) (ResponseCode, error) {
	_, err := fmt.Fprintln(client.conn, msg)

	if err != nil {
		return UNKNOWN, fmt.Errorf("Error sending message to server: %v", err)
	}

	responseMsg, err := bufio.NewReader(client.conn).ReadString('\n')
	if err != nil {
		return UNKNOWN, fmt.Errorf("Error reading response code from server: %v", err)
	}

	returnedString := strings.TrimRight(responseMsg, "\n")

	if returnedString == OK {
		return OK, nil
	}

	if returnedString == FAIL {
		return FAIL, nil
	}

	if returnedString == ERROR {
		return ERROR, nil
	}

	return UNKNOWN, fmt.Errorf("Error parsing message from server [%s]: %v", responseMsg, err)
}

// MakeTCPPackageIndexClient returns a new instance of the client
func MakeTCPPackageIndexClient(port int) (PackageIndexerClient, error) {
	host := fmt.Sprintf("localhost:%d", port)
	conn, err := net.Dial("tcp", host)

	if err != nil {
		return nil, fmt.Errorf("Failed to open connection to [%s]: %#v", host, err)
	}

	return &TCPPackageIndexerClient{
		conn: conn,
	}, nil
}
