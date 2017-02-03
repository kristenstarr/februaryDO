package input

import (
	//"time"
	"bufio"
	"net"
)

type Throttler interface {

	Throttle(conn net.Conn, output chan string)
}

type SimpleThrottler struct{
	rate *int
}

func (s *SimpleThrottler) Throttle(conn net.Conn, output chan string) {

	//A rate less than zero is interpreted as an error and defaults to zero.
	//if (*s.rate <= 0) {
		for {
			message, messageError := bufio.NewReader(conn).ReadString('\n')

			// If we have an error reading from socket, close and break loop/goroutine
			if (messageError != nil) {
				conn.Write([]byte("ERROR\n"))
				conn.Close()
				close(output)
				break
			} else {
				output <- message
			}

		}
	//} else {
	//	throttle := time.NewTicker(time.Duration(*s.rate))
	//	for _ = range throttle.C {
	//		message, messageError := bufio.NewReader(conn).ReadString('\n')
	//		// If we have an error reading from socket, close and break loop/goroutine
	//		if (messageError != nil) {
	//			conn.Write([]byte("ERROR\n"))
	//			conn.Close()
	//			throttle.Stop()
	//			break
	//		} else {
	//			output <- message
	//		}
	//	}
	//}
}

func NewThrottler(rate *int) Throttler {
	return &SimpleThrottler{
		rate,
	}
}
