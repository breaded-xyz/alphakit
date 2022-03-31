package backtest

import (
	"time"
)

// Compiler hint that Clock must implement Clocker.
var _ Clocker = (*Clock)(nil)

// Clocker defines a clock to be used by Simulation to timestamp events.
type Clocker interface {

	// Resets state with the first epoch at the given time.
	Start(time.Time)

	// NextEpoch advances the simulation clock to a new time.
	NextEpoch(time.Time)

	// Now returns an incrementally later each time it is called.
	// The returned time should always be >= the epoch start time.
	Now() time.Time

	// Epoch returns the start time of the current epoch.
	Epoch() time.Time

	// Elapsed returns the total duration since the first epoch.
	Elapsed() time.Duration
}

// Clock is the default Clocker implementation for Simulation.
// When Now is called an incrementally later time is returned.
// Each increment is a 'tock' and equals 1 * time.Millisecond.
// Tock term is used to avoid confusion with 'tick' which has a defined meaning in trading.
// Clock helps ensure orders are processed in the sequence they are submitted.
type Clock struct {
	tock    int64
	epoch   time.Time
	elapsed time.Duration
}

// NewClock sets the starting epoch to the zero time.
func NewClock() Clocker {
	return &Clock{}
}

// Start initializes the clock and resets all state.
func (c *Clock) Start(epoch time.Time) {
	c.epoch = epoch
	c.tock = 0
	c.elapsed = 0
}

// NextEpoch advances to the next epoch with the start at the given time.
// Tock counter is reset, when Now is called it will be epoch + 1 tock.
// Undefined behaviour if the given epoch is earlier than the current.
func (c *Clock) NextEpoch(epoch time.Time) {
	c.elapsed += epoch.Sub(c.epoch)
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

// Elapsed returns the total elapsed duration since the first epoch.
// This is primarily used for calculating funding charges.
func (c *Clock) Elapsed() time.Duration {
	return c.elapsed
}
