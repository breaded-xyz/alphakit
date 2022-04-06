package ta

func CrossUp(vals []float64, x float64) bool {
	curr := Lookback(vals, 0)
	prev := Lookback(vals, 1)
	if prev <= x && curr > x {
		return true
	}
	return false
}

func CrossDown(vals []float64, x float64) bool {
	curr := Lookback(vals, 0)
	prev := Lookback(vals, 1)
	if prev >= x && curr < x {
		return true
	}
	return false
}
