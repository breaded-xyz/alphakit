package backtest

import (
	"time"
)

// Compiler hint that Clock must implement Clocker.
var _ Clocker = (*Clock)(nil)

// Clocker defines a clock to be used by Simulation to timestamp events.
type Clocker interface {

	// NextEpoch advances the simulation clock to a new time.
	NextEpoch(time.Time)

	// Now returns a time that is the same as the epoch start time
	// or incrementally later each time it is called.
	Now() time.Time

	// Epoch returns the start time of the current epoch.
	Epoch() time.Time
}

// Clock is the default Clocker implementation for Simulation.
// When Now is called an incrementally later time is returned.
// Each increment is a 'tock' and equals 1 * time.Millisecond.
// Tock term is used to avoid confusion with 'tick' which has a defined meaning in trading.
// Clock helps ensure orders are processed in the sequence they are submitted.
type Clock struct {
	tock  int64
	epoch time.Time
}

// NewClock sets the epoch to the current time.
func NewClock() Clocker {
	return &Clock{
		epoch: time.Now(),
	}
}

// NextEpoch advances to the next epoch with the given time.
// Tock counter is reset. Next time Now is called it will be epoch + 1 tock.
func (c *Clock) NextEpoch(epoch time.Time) {
	c.epoch = epoch
	c.tock = 0
}

// Now returns the next tock, which is 1 * time.millisecond later than the last call.
// To reset tock call NextEpoch.
func (c *Clock) Now() time.Time {
	c.tock++
	return c.epoch.Add(time.Duration(c.tock) * time.Millisecond)
}

// Epoch returns the current epoch start time.
func (c *Clock) Epoch() time.Time {
	return c.epoch
}
