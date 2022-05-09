package ta

// Peak returns true if the latest values in the series have formed a peak.
// Series must be in chronological order, with the earliest value at slice index 0.
// Series must have at least 3 values.
func Peak(series []float64) bool {
	if len(series) < 3 {
		return false
	}
	t2, t1, t0 := Lookback(series, 2), Lookback(series, 1), Lookback(series, 0)
	prev := Slope(t2, t1)
	curr := Slope(t1, t0)
	return prev == 1 && curr == -1
}

// Valley returns true if the latest values in the series have formed a valley.
// Series must be in chronological order, with the earliest value at slice index 0.
// Series must have at least 3 values.
func Valley(series []float64) bool {
	if len(series) < 3 {
		return false
	}
	t2, t1, t0 := Lookback(series, 2), Lookback(series, 1), Lookback(series, 0)
	prev := Slope(t2, t1)
	curr := Slope(t1, t0)
	return prev == -1 && curr == 1
}
