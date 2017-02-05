package input

import (
	"time"
	"net"
)

// TestConnection implements net.Conn, to be used for testing purposes.
type TestConnection struct {}

func (t *TestConnection) Read(b []byte) (n int, err error) {
	return 100, nil
}

func (t *TestConnection) Write(b []byte) (n int, err error) {
	return 100, nil
}

func (t *TestConnection) Close() error {
	return nil
}

func (t *TestConnection) LocalAddr() net.Addr {
	return nil
}

func (t *TestConnection) RemoteAddr() net.Addr {
	return nil
}

func (t *TestConnection) SetDeadline(ti time.Time) error {
	return nil
}

func (t *TestConnection) SetReadDeadline(ti time.Time) error {
	return nil
}

func (t *TestConnection) SetWriteDeadline(ti time.Time) error {
	return nil
}

func NewTestConnection() net.Conn {
	return &TestConnection{}
}
