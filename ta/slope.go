// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package ta

// Slope indicates the direction between two points.
// t1 is the first point, t2 is the second point.
// Returns 1 for up, -1 for down, 0 for flat.
func Slope(t1, t2 float64) int {
	s := t2 - t1
	switch {
	case s < 0:
		return -1
	case s > 0:
		return 1
	}
	return 0
}
