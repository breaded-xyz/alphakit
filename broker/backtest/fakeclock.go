package backtest

import "time"

var _ Clocker = (*FakeClock)(nil)

type FakeClock struct {
	Fixed time.Time
}

func (c *FakeClock) NextEpoch(epoch time.Time) {
	// Fake does nothing
}

func (c *FakeClock) Now() time.Time {
	return c.Fixed
}

func (c *FakeClock) Epoch() time.Time {
	return c.Fixed
}
