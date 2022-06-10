package ta

import (
	"github.com/mpiannucci/peakdetect"
)

// Peak returns true if the latest values in the series have formed a peak.
// Series must be in chronological order, with the earliest value at slice index 0.
// Series must have at least 3 values.
// Arg delta is the threshold change in the series values required to detect a peak.
func Peak(series []float64, delta float64) bool {
	if len(series) < 3 {
		return false
	}
	_, _, maxIdx, _ := peakdetect.PeakDetect(Window(series, 2), delta)

	if len(maxIdx) == 0 {
		return false
	}

	return maxIdx[0] == 1
}

// Valley returns true if the latest values in the series have formed a valley.
// Series must be in chronological order, with the earliest value at slice index 0.
// Series must have at least 3 values.
// Arg delta is the threshold change in the series values required to detect a valley.
func Valley(series []float64, delta float64) bool {
	if len(series) < 3 {
		return false
	}
	minIdx, _, _, _ := peakdetect.PeakDetect(Window(series, 2), delta)

	if len(minIdx) == 0 {
		return false
	}

	return minIdx[0] == 1
}
