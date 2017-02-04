package input

import (
	"testing"
	"time"
)

// Tests our default Throttle behavior when no rate limiting should be applied.
func TestNoThrottling(t *testing.T) {
	rate := 0
	throttler := NewThrottler(&rate)

	startTime := time.Now().UnixNano()
	throttler.Next()
	endTime := time.Now().UnixNano()

	next := throttler.Next()
	if (next) {
		t.Error("Next should return false when rate limiting is not in place")
	}
	noThrottleTime := endTime - startTime
	if (noThrottleTime > 10000) {
		t.Error("With no rate limit, processing time should be very small")
	}
}

// Tests one hundred per second throttling
func TestTenPerSecond(t *testing.T) {
	rate := 100
	throttler := NewThrottler(&rate)

	startTime := time.Now().UnixNano()
	throttler.Next()
	endTime := time.Now().UnixNano()

	next := throttler.Next()
	if (!next) {
		t.Error("Next should be called successfully when throttler is active")
	}
	onePerSecondThrottleTime := endTime - startTime
	if (onePerSecondThrottleTime < 10000000 || onePerSecondThrottleTime > 20000000) {
		t.Error("We should see a .01 second delay with one per second throttling")
	}
}

// Tests that the processing time for no rate limiting is less than that of 100 per second limiting.
func TestThrottlingComparison(t *testing.T) {
	rate := 0
	throttler := NewThrottler(&rate)

	startTime := time.Now().UnixNano()
	throttler.Next()
	endTime := time.Now().UnixNano()

	noThrottleTime := endTime - startTime
	rate = 100
	throttler = NewThrottler(&rate)

	startTime = time.Now().UnixNano()
	throttler.Next()
	endTime = time.Now().UnixNano()

	onePerSecondThrottleTime := endTime - startTime

	if(noThrottleTime > onePerSecondThrottleTime) {
		t.Error("Adding rate limiting at 100/second should increase processing time")
	}

}

// Tests that once throttle is stopped, Next() calls return false.
func TestThrottleStop(t *testing.T) {
	rate := 100
	throttler := NewThrottler(&rate)
	throttler.Stop()
	next := throttler.Next()
	if (next) {
		t.Error("Once throttler is stopped, we should not be able to call Next")
	}
}

