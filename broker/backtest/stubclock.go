package backtest

import "time"

// Compiler hint that StubClock must implement Clocker.
var _ Clocker = (*StubClock)(nil)

// StubClock provides a canned time to Simulator for testing.
// For a discussion of various types of 'test double' see:
// https://www.martinfowler.com/articles/mocksArentStubs.html
type StubClock struct {

	// Fixed is the time to return
	Fixed time.Time
}

// NextEpoch not implemented.
func (c *StubClock) NextEpoch(epoch time.Time) {
	// Intentionally empty
}

// Now returns Fixed.
func (c *StubClock) Now() time.Time {
	return c.Fixed
}

// Epoch returns Fixed.
func (c *StubClock) Epoch() time.Time {
	return c.Fixed
}
