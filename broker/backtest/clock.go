package backtest

import (
	"time"
)

// Clock is the default Clocker implementation for Simulation.
// When Now is called an incrementally later time is returned.
// Each increment is a 'tock' which defaults to 1 millisecond.
// Tock term is used to avoid confusion with 'tick' which has a defined meaning in trading.
// Clock helps ensure orders are processed in the sequence they are submitted.
type Clock struct {
	now      time.Time
	interval time.Duration
	elapsed  time.Duration
}

// NewClock sets the start to the zero time and tock interval to 1 millisecond.
func NewClock() *Clock {
	return &Clock{
		interval: 1 * time.Millisecond,
	}
}

// Start initializes the clock and resets all state.
func (c *Clock) Start(start time.Time, tock time.Duration) {
	c.now = start
	c.interval = tock
	c.elapsed = 0
}

// Advance advances to the next epoch at the given time.
// When Now is next called it will be epoch + 1 tock interval.
// Undefined behaviour if the given epoch is earlier than the current.
func (c *Clock) Advance(epoch time.Time) {
	c.elapsed += epoch.Sub(c.now)
	c.now = epoch
}

// Now returns the time incremented by a tock,
// which by default is  1 * time.millisecond later than the last call.
func (c *Clock) Now() time.Time {
	c.now = c.now.Add(c.interval)
	return c.now
}

// Peek returns the time without mutation.
func (c *Clock) Peek() time.Time {
	return c.now
}

// Elapsed returns the total elapsed duration since the start.
// Elapsed time is calculated on each call to Advance.
// Primarily used for calculating funding charges.
func (c *Clock) Elapsed() time.Duration {
	return c.elapsed
}
