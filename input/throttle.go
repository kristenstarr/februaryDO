package input

import "time"

// Throttler is responsible for rate limiting our traffic.
// It has 2 methods - Next() will return when it is a permissible time for the next message to
// be processed, Stop will stop the throttler from generating any more available opportunities
// for messages to be processed.
type Throttler interface {

	// Next() will complete as soon as it is an appropriate time for the next message to be handled.
	// Will return true if rate limiting is in place and Next completes successfully.
	// Will return false if we have no rate limiting or our throttler is stopped
	Next() bool

	// Stop() is primarily for clean up purposes and stops our throttler from generating opportunities.
	Stop()
}

// SimpleThrottler is a Throttler with a Ticker and a boolean to indicate state.
type SimpleThrottler struct{
	ticker *time.Ticker
	stopped bool
}

func (s *SimpleThrottler) Next() bool {
	if (s.ticker != nil && !s.stopped) {
		<-s.ticker.C
		return true
	}
	return false
}

func (s *SimpleThrottler) Stop() {
	if (s.ticker != nil) {
		s.stopped = true
		s.ticker.Stop()
	}
}

// NewThrottler creates a new instance of a Throttler for our traffic, starting a time Ticker.
func NewThrottler(rate *int) Throttler {
	// We are implementing a throttle using a 'ticker' that adds a message to a channel at a set rate.
	var ticker *time.Ticker
	if (*rate > 0) {
		ticker = time.NewTicker(time.Second / time.Duration(*rate))
	}
	return &SimpleThrottler{
		ticker,
		false,
	}
}
