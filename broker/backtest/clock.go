package backtest

import "time"

type Clock struct {
	tock  int64
	epoch time.Time
}

func NewClock() Clock {
	return Clock{
		epoch: time.Now(),
	}
}

func (c *Clock) NextEpoch(epoch time.Time) {
	c.epoch = epoch
	c.tock = 0
}

func (c *Clock) Now() time.Time {
	c.tock++
	return c.epoch.Add(time.Duration(c.tock) * time.Millisecond)
}

func (c *Clock) Epoch() time.Time {
	return c.epoch
}
