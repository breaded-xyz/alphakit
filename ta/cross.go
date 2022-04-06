package ta

func CrossUp(series []float64, x float64) bool {
	curr := Lookback(series, 0)
	prev := Lookback(series, 1)
	if prev <= x && curr > x {
		return true
	}
	return false
}

func CrossDown(series []float64, x float64) bool {
	curr := Lookback(series, 0)
	prev := Lookback(series, 1)
	if prev >= x && curr < x {
		return true
	}
	return false
}
