package backtest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClockStart(t *testing.T) {
	giveStart := time.Now()
	giveInterval := time.Second * 10

	want := &Clock{
		now:      giveStart,
		interval: giveInterval,
		elapsed:  0,
	}

	act := NewClock()
	act.Start(giveStart, giveInterval)
	assert.Equal(t, want, act)
}

func TestClockAdvance(t *testing.T) {
	start := time.Now()
	giveEpoch := start.Add(time.Hour * 2)

	want := Clock{
		now:     giveEpoch,
		elapsed: time.Hour * 2,
	}

	act := NewClock()
	act.Start(start, 0)
	act.Advance(giveEpoch)
	assert.Equal(t, want, *act)
}

func TestClockNow(t *testing.T) {
	start := time.Now()

	want := start.Add(1 * time.Hour)

	clock := NewClock()
	clock.Start(start, time.Hour)
	act := clock.Now()
	assert.Equal(t, want, act)
}
