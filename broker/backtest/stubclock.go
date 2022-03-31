package backtest

import "time"

var _ Clocker = (*StubClock)(nil)

type StubClock struct {
	Fixed time.Time
}

func (c *StubClock) NextEpoch(epoch time.Time) {
	// Fake does nothing
}

func (c *StubClock) Now() time.Time {
	return c.Fixed
}

func (c *StubClock) Epoch() time.Time {
	return c.Fixed
}
