// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package backtest

import "time"

// Compiler hint that Clock must implement Clocker.
var _ Clocker = (*Clock)(nil)

// Clocker defines a clock to be used by Simulation to timestamp events.
type Clocker interface {

	// Resets the clock with a new start time and tock interval.
	Start(time.Time, time.Duration)

	// Advance advances the simulation clock to a future time.
	Advance(time.Time)

	// Now returns an incrementally later time each call.
	// Increments are defined by the tock interval given to Start.
	// Returned time should always be >= the start time or latest advance time.
	Now() time.Time

	// Peek returns the current time (last value returned by Now())
	// but does not advance the time by the tock interval.
	Peek() time.Time

	// Elapsed returns the total duration since the start time.
	Elapsed() time.Duration
}
