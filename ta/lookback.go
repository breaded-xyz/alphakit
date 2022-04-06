package ta

func Lookback(series []float64, n int) float64 {
	i := (len(series) - n) - 1
	if i < 0 {
		return 0
	}
	return series[i]
}

func Window(series []float64, n int) []float64 {
	window := make([]float64, len(series))
	copy(window, series)

	ln := len(series)
	if ln <= n {
		return window
	}

	i := (ln - n) - 1
	return window[i:]
}

func AppendWindow(series []float64, n int, v float64) []float64 {
	return Window(append(series, v), n-1)
}
