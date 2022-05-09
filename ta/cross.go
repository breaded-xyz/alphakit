package ta

// CrossUp returns true if the latest series values cross above the given value.
// The series must be in chronological order, with the earliest value at index 0.
// The series must have at least 2 values
func CrossUp(series []float64, x float64) bool {
	curr := Lookback(series, 0)
	prev := Lookback(series, 1)
	if prev <= x && curr > x {
		return true
	}
	return false
}

// CrossDown returns true if the latest series values cross below the given value.
// The series must be in chronological order, with the earliest value at index 0.
// The series must have at least 2 values
func CrossDown(series []float64, x float64) bool {
	curr := Lookback(series, 0)
	prev := Lookback(series, 1)
	if prev >= x && curr < x {
		return true
	}
	return false
}
