package ta

func Peak(series []float64) bool {
	t2, t1, t0 := Lookback(series, 2), Lookback(series, 1), Lookback(series, 0)
	slo1 := Slope(t2, t1)
	slo2 := Slope(t1, t0)
	return (slo1 == 1 || slo1 == 0) && slo2 == -1
}

func Valley(series []float64) bool {
	t2, t1, t0 := Lookback(series, 2), Lookback(series, 1), Lookback(series, 0)
	slo1 := Slope(t2, t1)
	slo2 := Slope(t1, t0)
	return (slo1 == -1 || slo1 == 0) && slo2 == 1
}
