// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

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

// Start not implemented.
func (c *StubClock) Start(start time.Time, tock time.Duration) {
	// Intentionally empty
}

// Advance not implemented.
func (c *StubClock) Advance(epoch time.Time) {
	// Intentionally empty
}

// Now returns Fixed.
func (c *StubClock) Now() time.Time {
	return c.Fixed
}

// Peek returns Fixed.
func (c *StubClock) Peek() time.Time {
	return c.Fixed
}

// Elapsed returns 0.
func (c *StubClock) Elapsed() time.Duration {
	return 0
}
